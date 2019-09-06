package http

import (
	"fmt"
	"github.com/MoF-Dev/go-bt-tracker/internal/app/request"
	"github.com/MoF-Dev/go-bt-tracker/pkg/bencode"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
)

func tmpe(w http.ResponseWriter, statusCode int, failureReason string) {
	w.WriteHeader(statusCode)
	res := make(bencode.Dictionary)
	res["failure reason"] = bencode.String(failureReason)
	reply(w, res)
}

func reply(w http.ResponseWriter, val bencode.BValue) {
	_, err := fmt.Fprint(w, val.Encode())
	if err != nil {
		log.Println(err)
	}
}

func getQuery(w http.ResponseWriter, queries url.Values, key string, replyError bool) (string, bool) {
	raw := queries[key]
	if len(raw) != 1 {
		if replyError {
			tmpe(w, http.StatusBadRequest, "missing or too many parameter '"+key+"'")
		}
		return "", false
	}
	return raw[0], true
}

func getQueryBigInt(w http.ResponseWriter, queries url.Values, key string) (*big.Int, bool) {
	raw, success := getQuery(w, queries, key, true)
	if !success {
		return nil, false
	}
	return big.NewInt(0).SetString(raw, 10)
}

func AnnounceHandler(w http.ResponseWriter, r *http.Request) {
	var formalRequest request.AnnounceRequest

	q := r.URL.Query()

	// required params
	infoHashR, success := getQuery(w, q, "info_hash", true)
	if !success {
		return
	}
	infoHash := []byte(infoHashR)
	if len(infoHash) != 20 {
		tmpe(w, http.StatusBadRequest, "info_hash is not exactly 20 bytes")
	}
	copy(formalRequest.InfoHash[:], infoHash)

	peerIdR, success := getQuery(w, q, "peer_id", true)
	if !success {
		return
	}
	peerId := []byte(peerIdR)
	if len(peerId) != 20 {
		tmpe(w, http.StatusBadRequest, "peer_id is not exactly 20 bytes")
	}
	copy(formalRequest.PeerId[:], peerId)

	portR, success := getQuery(w, q, "port", true)
	if !success {
		return
	}
	portRI, err := strconv.ParseUint(portR, 10, 16)
	if err != nil {
		tmpe(w, http.StatusBadRequest, "port is not a valid number 1-65534")
		return
	}
	formalRequest.Port = uint16(portRI)

	formalRequest.Uploaded, success = getQueryBigInt(w, q, "uploaded")
	if !success {
		return
	}

	formalRequest.Downloaded, success = getQueryBigInt(w, q, "downloaded")
	if !success {
		return
	}

	formalRequest.Left, success = getQueryBigInt(w, q, "left")
	if !success {
		return
	}

	// optional params
	ip, success := getQuery(w, q, "ip", false)
	if success {
		// TODO sanity checks?
		formalRequest.Ip = &ip
	}

	compactR, success := getQuery(w, q, "compact", false)
	compact := true
	if success {
		if compactR == "0" {
			compact = false
		} else if compactR == "1" {
			compact = true
		} else {
			tmpe(w, http.StatusBadRequest, "compact must be 0 or 1 or omitted")
			return
		}
	}

	response, err := request.GetAnnounce(&formalRequest)
	if err != nil {
		tmpe(w, http.StatusInternalServerError, err.Error())
		return
	}
	_, _ = w.Write([]byte(response.Encode(compact).Encode()))
}

func ScrapeHandler(w http.ResponseWriter, r *http.Request) {
	var formalRequest request.ScrapeRequest
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

	response, err := request.GetScrape(&formalRequest)
	if err != nil {
		tmpe(w, http.StatusInternalServerError, err.Error())
		return
	}
	_, _ = w.Write([]byte(response.Encode().Encode()))
}
