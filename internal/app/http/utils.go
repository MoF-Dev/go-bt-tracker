package http

import (
	"fmt"
	"github.com/MoF-Dev/go-bt-tracker/pkg/bencode"
	"log"
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

func getQueryBigInt(w http.ResponseWriter, queries url.Values, key string) (uint64, bool) {
	raw, success := getQuery(w, queries, key, true)
	if !success {
		return 0, false
	}
	x, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		// TODO log
		return 0, false
	}
	return x, true
}
