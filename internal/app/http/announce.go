package http

import (
	"encoding/binary"
	"github.com/MoF-Dev/go-bt-tracker/pkg/bencode"
	"github.com/MoF-Dev/go-bt-tracker/pkg/tracker"
	"net/http"
	"strconv"
	"strings"
)

type AnnounceResponse tracker.AnnounceResponse

func (r AnnounceResponse) Encode(compact bool) bencode.Dictionary {
	dict := make(bencode.Dictionary)
	if r.FailureReason != nil {
		dict["failure reason"] = bencode.String(*r.FailureReason)
		return dict
	}
	dict["interval"] = bencode.NewUInteger(uint64(r.Interval))
	if r.MinInterval != nil {
		dict["min interval"] = bencode.NewUInteger(uint64(*r.MinInterval))
	}
	if r.TrackerId != nil {
		dict["tracker id"] = bencode.String(*r.TrackerId)
	}
	dict["complete"] = bencode.NewUInteger(uint64(r.Complete))
	dict["incomplete"] = bencode.NewUInteger(uint64(r.Incomplete))

	// TODO peer limit (numwant) is not yet implemented for HTTP, only for UDP
	// may require big changes
	if compact {
		peers := strings.Builder{}
		peers6 := strings.Builder{}
		for _, peer := range r.Peers {
			i4, i6, err := peer.GetIPs()
			if err != nil {
				if r.WarningMessage == nil {
					var str = err.Error()
					r.WarningMessage = &str
				} else {
					var str = *r.WarningMessage + " " + err.Error()
					r.WarningMessage = &str
				}
				continue
			}
			port := make([]byte, 2)
			binary.BigEndian.PutUint16(port, peer.Port)
			for _, ip := range i4 {
				peers.Write(ip[0:4])
				peers.Write(port)
			}
			for _, ip := range i6 {
				peers6.Write(ip[0:16])
				peers6.Write(port)
			}
		}
		dict["peers"] = bencode.String(peers.String())
		if peers6.Len() > 0 {
			dict["peers6"] = bencode.String(peers6.String())
		}
	} else {
		peers := make(bencode.List, len(r.Peers))
		for i, peer := range r.Peers {
			peers[i] = Peer(peer).Encode()
		}
		dict["peers"] = peers
	}

	if r.WarningMessage != nil {
		dict["warning message"] = bencode.String(*r.WarningMessage)
	}
	return dict
}

type Peer tracker.Peer

func (p Peer) Encode() bencode.Dictionary {
	dict := make(bencode.Dictionary)
	dict["peer id"] = bencode.String(p.PeerId)
	dict["ip"] = bencode.String(p.Ip)
	dict["port"] = bencode.NewUInteger(uint64(p.Port))
	return dict
}

func AnnounceHandler(server tracker.HttpServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var formalRequest tracker.AnnounceRequest

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

		response, err := server.HandleAnnounce(&formalRequest)
		if err != nil {
			tmpe(w, http.StatusInternalServerError, err.Error())
			return
		}
		_, _ = w.Write([]byte((*AnnounceResponse)(response).Encode(compact).Encode()))
	}
}
