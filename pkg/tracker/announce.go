package tracker

import "net"

type Event uint32

const (
	Empty     Event = 0
	Completed       = 1
	Started         = 2
	Stopped         = 3
)

type AnnounceRequest struct {
	InfoHash   [20]byte
	PeerId     [20]byte
	Ip         *string
	Port       uint16
	Uploaded   uint64
	Downloaded uint64
	Left       uint64
	Event      Event
	NumWant    int32
}

type AnnounceResponse struct {
	FailureReason *string // if present, no other keys will be shown

	Interval    uint32
	MinInterval *uint32
	TrackerId   *string
	Complete    uint32
	Incomplete  uint32
	Peers       []Peer

	WarningMessage *string
}

type Peer struct {
	PeerId string
	Ip     string
	Port   uint16
}

func (p Peer) GetIPs() (i4 []net.IP, i6 []net.IP, err error) {
	var ips []net.IP
	ip := net.ParseIP(p.Ip)
	if ip == nil {
		ips, err = net.LookupIP(p.Ip)
		if err != nil {
			return nil, nil, err
		}
	} else {
		ips = append(ips, ip)
	}
	for _, ip := range ips {
		if ip.To4() != nil {
			i4 = append(i4, ip.To4())
		} else {
			i6 = append(i6, ip.To16())
		}
	}
	return i4, i6, nil
}
