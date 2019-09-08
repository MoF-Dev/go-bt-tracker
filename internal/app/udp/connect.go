package udp

import (
	"encoding/binary"
	"github.com/MoF-Dev/go-bt-tracker/pkg/tracker"
	"net"
)

func connectHandler(server tracker.UdpServer, addr net.Addr, base *basePacket, buf []byte) {
	connId, err := server.NewSession()
	if err != nil {
		// TODO
		return
	}
	var res [16]byte
	writeBaseReply(res[:], base)
	binary.BigEndian.PutUint64(res[8:], connId)
	_, err = server.WriteTo(res[:], addr)
	if err != nil {
		//TODO
	}
}
