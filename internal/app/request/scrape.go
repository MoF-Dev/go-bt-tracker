package request

import (
	"errors"
	"github.com/MoF-Dev/go-bt-tracker/pkg/bencode"
)

type ScrapeResponse struct {
	Files         []File
	FailureReason *string
	Flags         *bencode.Dictionary
}

func (r ScrapeResponse) Encode() bencode.Dictionary {
	dict := make(bencode.Dictionary)
	files := make(bencode.List, len(r.Files))
	for i, file := range r.Files {
		files[i] = file.Encode()
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

type File struct {
	Completed  uint32
	Downloaded uint32
	Incomplete uint32
	Name       *string
}

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

func GetScrape(r *ScrapeRequest) (*ScrapeResponse, error) {
	// TODO
	return nil, errors.New("not yet implemented")
}
