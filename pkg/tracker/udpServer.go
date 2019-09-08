package tracker

import (
	"github.com/MoF-Dev/go-bt-tracker/internal/app/udp"
	"net"
)

type udpServerExtras interface {
	ReadFrom(p []byte) (n int, addr net.Addr, err error)
	WriteTo(p []byte, addr net.Addr) (n int, err error)
	Close() error
	NewSession() (uint64, error)
	CheckSession(connId uint64) (validSession bool, err error)
}

type UdpServer interface {
	Server
	udpServerExtras
}

func ListenUdp(server UdpServer) {
	for {
		// TODO what size does the buffer have to be? can it be less?
		buf := make([]byte, 1024)
		n, addr, err := server.ReadFrom(buf)
		if err != nil {
			// TODO handle the err
			continue
		}
		go udp.Handler(server, addr, buf[:n])
	}
}
