package udp

import (
	"encoding/binary"
	"github.com/MoF-Dev/go-bt-tracker/pkg/tracker"
	"net"
)

func announceHandler(server tracker.UdpServer, addr net.Addr, base *basePacket, buf []byte) {
	if len(buf) != 98 {
		// TODO lw
		return
	}

	var req tracker.AnnounceRequest
	var BE = binary.BigEndian
	/*	Offset  Size    Name    Value
		0       64-bit integer  connection_id
		8       32-bit integer  action          1 // announce
		12      32-bit integer  transaction_id
		16      20-byte string  info_hash
		36      20-byte string  peer_id
		56      64-bit integer  downloaded
		64      64-bit integer  left
		72      64-bit integer  uploaded
		80      32-bit integer  event           0 // 0: none; 1: completed; 2: started; 3: stopped
		84      32-bit integer  IP address      0 // default
		88      32-bit integer  key
		92      32-bit integer  num_want        -1 // default
		96      16-bit integer  port
		98 */
	copy(req.InfoHash[:], buf[16:])
	copy(req.PeerId[:], buf[36:])
	req.Downloaded = BE.Uint64(buf[56:])
	req.Left = BE.Uint64(buf[64:])
	req.Uploaded = BE.Uint64(buf[72:])
	req.Event = tracker.Event(BE.Uint32(buf[80:]))
	if BE.Uint32(buf[84:]) != 0 {
		ip := net.IPv4(buf[84], buf[85], buf[86], buf[87])
		ipStr := ip.String()
		req.Ip = &ipStr
	}
	// req.Key???
	req.NumWant = int32(BE.Uint32(buf[92:]))
	req.Port = BE.Uint16(buf[96:])

	fr, err := server.HandleAnnounce(&req)
	if err != nil {
		// TODO
		return
	}

	res := make([]byte, 8)
	writeBaseReply(res, base)

	if fr.FailureReason != nil {
		err = writeError(server, addr, base.TransactionId, *fr.FailureReason, nil)
		// TODO deal with err
		return
	}

	var fixedBody [12]byte
	BE.PutUint32(fixedBody[0:], fr.Interval)
	BE.PutUint32(fixedBody[4:], fr.Incomplete)
	BE.PutUint32(fixedBody[8:], fr.Complete)

	var is4 bool
	var peerSize int
	if len(addr.(*net.UDPAddr).IP) == net.IPv4len {
		is4 = true
		peerSize = net.IPv4len + 2
	} else {
		is4 = false
		peerSize = net.IPv6len + 2
	}

	var peers [][]byte
	for _, peer := range fr.Peers {
		ips, i6, err := peer.GetIPs()
		if err != nil {
			// TODO log
			continue
		}
		var port [2]byte
		BE.PutUint16(port[:], peer.Port)
		if !is4 {
			ips = i6
		}
		for _, ip := range ips {
			peers = append(peers, append(ip[:peerSize-2], port[0], port[1]))
		}
	}

	numHave, peers := server.ChooseLimitedPeers(peers, req.NumWant)
	peersBuf := make([]byte, peerSize*numHave)
	for i, peer := range peers {
		start := i * peerSize
		copy(peersBuf[start:start+peerSize], peer)
	}

	packet := append(res, fixedBody[:]...)
	packet = append(packet, peersBuf...)
	// TODO is there a limit to UDP packet size?

	_, err = server.WriteTo(packet, addr)
	if err != nil {
		//TODO
	}
}
