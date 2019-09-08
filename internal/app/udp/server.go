package udp

import (
	"encoding/binary"
	"github.com/MoF-Dev/go-bt-tracker/pkg/tracker"
	"net"
)

func Handler(server tracker.UdpServer, addr net.Addr, buf []byte) {
	// all packets have minimum len(connId)+len(action)+len(transId)=16
	if len(buf) < 16 {
		// TODO log warning but do not reply
		return
	}

	// we prolly do not want to reply to a bad request because we don't want to amplify attacks
	var BE = binary.BigEndian
	var request basePacket
	request.ConnectionId = BE.Uint64(buf[0:])
	request.Action = Action(BE.Uint32(buf[8:]))
	request.TransactionId = BE.Uint32(buf[12:])
	switch request.Action {
	case Connect:
		if request.ConnectionId != ProtocolMagic {
			// TODO lwbdnr
			return
		}
		connectHandler(server, addr, &request, buf)
	case Announce:
		valid, err := server.CheckSession(request.ConnectionId)
		if err != nil {
			// TODO lwbdnr
			return
		}
		if !valid {
			// TODO lwbdnr
			return
		}
		announceHandler(server, addr, &request, buf)
	case Scrape:
		valid, err := server.CheckSession(request.ConnectionId)
		if err != nil {
			// TODO lwbdnr
			return
		}
		if !valid {
			// TODO lwbdnr
			return
		}
		scrapeHandler(server, addr, &request, buf)
	default:
		// TODO log warning but do not reply
	}
}

func writeError(s tracker.UdpServer, addr net.Addr, tid uint32, msg string, cause error) (err error) {
	if cause != nil {
		// TODO
	}
	var base [8]byte
	binary.BigEndian.PutUint32(base[0:], Error)
	binary.BigEndian.PutUint32(base[4:], tid)
	_, err = s.WriteTo(append(base[:], []byte(msg)...), addr)
	return err
}
