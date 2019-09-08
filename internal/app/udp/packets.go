package udp

import "encoding/binary"

const ProtocolMagic uint64 = 0x41727101980

type Action uint32

const (
	Connect  Action = 0
	Announce        = 1
	Scrape          = 2
	Error           = 3
)

type basePacket struct {
	ConnectionId  uint64
	Action        Action
	TransactionId uint32
}

func writeBaseReply(response []byte, base *basePacket) {
	var BE = binary.BigEndian
	BE.PutUint32(response[0:], uint32(base.Action))
	BE.PutUint32(response[4:], base.TransactionId)
}
