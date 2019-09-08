package request

import (
	"encoding/binary"
	"errors"
	"github.com/MoF-Dev/go-bt-tracker/pkg/bencode"
	"net"
	"strings"
)

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
			peers[i] = peer.Encode()
		}
		dict["peers"] = peers
	}

	if r.WarningMessage != nil {
		dict["warning message"] = bencode.String(*r.WarningMessage)
	}
	return dict
}

type Peer struct {
	PeerId string
	Ip     string
	Port   uint16
}

func (p Peer) Encode() bencode.Dictionary {
	dict := make(bencode.Dictionary)
	dict["peer id"] = bencode.String(p.PeerId)
	dict["ip"] = bencode.String(p.Ip)
	dict["port"] = bencode.NewUInteger(uint64(p.Port))
	return dict
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

func GetAnnounce(r *AnnounceRequest) (*AnnounceResponse, error) {
	// TODO
	return nil, errors.New("not yet implemented")
}
