package udp

import (
	"encoding/binary"
	"github.com/MoF-Dev/go-bt-tracker/pkg/tracker"
	"net"
)

func scrapeHandler(server tracker.UdpServer, addr net.Addr, base *basePacket, buf []byte) {
	// packet size is 16+20*n with n>=1
	if len(buf) <= 36 {
		// TODO log
		return
	}
	numTorrent := len(buf) - 16
	if numTorrent%20 != 0 {
		// TODO log, must be 20n
		return
	}
	numTorrent /= 20

	var req tracker.ScrapeRequest
	req.InfoHashes = make([][20]byte, numTorrent)
	for i := 0; i < numTorrent; i++ {
		copy(req.InfoHashes[i][:], buf[16+20*i:])
	}

	fr, err := server.HandleScrape(&req)
	if err != nil {
		// TODO handle
		return
	}
	if fr.FailureReason != nil {
		err = writeError(server, addr, base.TransactionId, *fr.FailureReason, nil)
		if err != nil {
			// TODO
		}
		return
	}

	resLen := len(fr.Files)
	res := make([]byte, 8+12*resLen)
	writeBaseReply(res, base)

	var BE = binary.BigEndian
	for i, file := range fr.Files {
		start := 8 + 12*i
		BE.PutUint32(res[start:start+4], file.Downloaded)
		BE.PutUint32(res[start+4:start+8], file.Completed)
		BE.PutUint32(res[start+8:start+12], file.Incomplete)
	}

	_, err = server.WriteTo(res, addr)
	if err != nil {
		// TODO
	}
}
