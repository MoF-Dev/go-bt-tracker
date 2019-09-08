package http

import (
	"github.com/MoF-Dev/go-bt-tracker/pkg/bencode"
	"github.com/MoF-Dev/go-bt-tracker/pkg/tracker"
	"net/http"
)

type ScrapeResponse tracker.ScrapeResponse

func (r ScrapeResponse) Encode() bencode.Dictionary {
	dict := make(bencode.Dictionary)
	files := make(bencode.List, len(r.Files))
	for i, file := range r.Files {
		files[i] = File(file).Encode()
	}
	dict["files"] = files
	if r.FailureReason != nil {
		dict["failure reason"] = bencode.String(*r.FailureReason)
	}
	if r.Flags != nil {
		dict["flags"] = *r.Flags
	}
	return dict
}

type File tracker.File

func (f File) Encode() bencode.Dictionary {
	dict := make(bencode.Dictionary)
	dict["completed"] = bencode.NewUInteger(uint64(f.Completed))
	dict["downloaded"] = bencode.NewUInteger(uint64(f.Downloaded))
	dict["incomplete"] = bencode.NewUInteger(uint64(f.Incomplete))
	if f.Name != nil {
		dict["name"] = bencode.String(*f.Name)
	}
	return dict
}

func ScrapeHandler(server tracker.HttpServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var formalRequest tracker.ScrapeRequest
		formalRequest.InfoHashes = make([][20]byte, 0)

		q := r.URL.Query()
		for _, infoHashR := range q["info_hash"] {
			if len(infoHashR) != 20 {
				tmpe(w, http.StatusBadRequest, "info_hash must be exactly 20 bytes")
				return
			}
			var infoHash [20]byte
			copy(infoHash[:], infoHashR)
			formalRequest.InfoHashes = append(formalRequest.InfoHashes, infoHash)
		}

		response, err := server.HandleScrape(&formalRequest)
		if err != nil {
			tmpe(w, http.StatusInternalServerError, err.Error())
			return
		}
		_, _ = w.Write([]byte((*ScrapeResponse)(response).Encode().Encode()))
	}
}
