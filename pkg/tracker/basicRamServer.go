package tracker

import (
	srand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"github.com/patrickmn/go-cache"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type brs struct {
	conn     *net.UDPConn
	sessions *cache.Cache
	torrents *cache.Cache
}

type torrent struct {
	peers      *cache.Cache
	downloaded uint32
}

func (s *brs) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

type peer struct {
	port       uint16
	uploaded   uint64
	downloaded uint64
	left       uint64
	ip         string
}

func (s *brs) HandleAnnounce(request *AnnounceRequest) (*AnnounceResponse, error) {
	// accepts all torrents
	key := hex.EncodeToString(request.InfoHash[:])
	var tor *torrent
	cachedTor, has := s.torrents.Get(key)
	if has {
		tor = cachedTor.(*torrent)
	} else {
		tor = &torrent{}
		tor.peers = cache.New(1*time.Minute, 30*time.Second)
		s.torrents.Set(key, tor, cache.DefaultExpiration)
	}

	peerId := hex.EncodeToString(request.PeerId[:])
	if request.Event == Stopped {
		tor.peers.Delete(peerId)
	} else {
		cachedPeer, has := tor.peers.Get(peerId)
		var p *peer
		if has {
			p = cachedPeer.(*peer)
		} else {
			p = &peer{}
		}
		p.port = request.Port
		p.downloaded = request.Downloaded
		p.uploaded = request.Uploaded
		p.left = request.Left
		p.ip = *request.Ip

		tor.peers.Set(peerId, p, cache.DefaultExpiration)
	}
	if request.Event == Completed {
		tor.downloaded++
	}

	var res AnnounceResponse
	res.Interval = 30
	var minInterval uint32 = 10
	res.MinInterval = &minInterval
	for peerId, rp := range tor.peers.Items() {
		p := rp.Object.(*peer)
		if p.left > 0 {
			res.Incomplete += 1
		} else {
			res.Complete += 1
		}
		peerIdBytes, err := hex.DecodeString(peerId)
		if err != nil {
			// TODO maybe soft error?
			continue
		}
		res.Peers = append(res.Peers, Peer{
			PeerId: string(peerIdBytes),
			Ip:     p.ip,
			Port:   p.port,
		})
	}

	return &res, nil
}

func (s *brs) HandleScrape(request *ScrapeRequest) (*ScrapeResponse, error) {
	var res ScrapeResponse
	// TODO should we allow full scraping?
	filter := request.InfoHashes
	if len(filter) == 0 {
		for hashS, _ := range s.torrents.Items() {
			var hashB [20]byte
			copy(hashB[:], hashS)
			filter = append(filter, hashB)
		}
	}
	for _, torIdB := range filter {
		torIdHex := hex.EncodeToString(torIdB[:])
		cTor, has := s.torrents.Get(torIdHex)
		if !has {
			continue
		}
		tor := cTor.(*torrent)
		f := File{
			Completed:  0, // seeders
			Downloaded: tor.downloaded,
			Incomplete: 0,
			Name:       nil,
		}
		for _, rp := range tor.peers.Items() {
			p := rp.Object.(*peer)
			if p.left > 0 {
				f.Incomplete++
			} else {
				f.Completed++
			}
		}
		res.Files = append(res.Files, f)
	}
	return &res, nil
}

func (s *brs) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	return s.conn.ReadFrom(p)
}

func (s *brs) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	return s.conn.WriteTo(p, addr)
}

func (s *brs) NewSession() (sid uint64, _ error) {
	var has = true
	var sidText string
	for has {
		var b [4]byte
		if _, err := srand.Read(b[:]); err != nil {
			return 0, err
		}
		sid = binary.BigEndian.Uint64(b[:])
		sidText = strconv.FormatUint(sid, 16)
		_, has = s.sessions.Get(sidText)
	}
	s.sessions.Set(sidText, '0', cache.DefaultExpiration)
	return sid, nil
}

func (s *brs) CheckSession(connId uint64) (validSession bool, err error) {
	sidText := strconv.FormatUint(connId, 16)
	_, has := s.sessions.Get(sidText)
	return has, nil
}

func NewBasicRamServer(listen string) HttpUdpServer {
	var server brs
	pc, err := net.ListenPacket("udp", listen)
	if err != nil {
		log.Fatal(err)
	}
	server.conn = pc.(*net.UDPConn)
	// TODO make sure default times are right + configurable
	server.sessions = cache.New(1*time.Minute, 1*time.Minute)
	server.torrents = cache.New(7*24*time.Hour, 1*time.Hour)
	return &server
}

func (*brs) ChooseLimitedPeers(peers [][]byte, n int32) (limit int, randomPeers [][]byte) {
	limit = 50 // TODO default limit in config
	if n > 0 {
		limit = int(n)
	}
	const UpperLimit = 100 // TODO should be set via config
	if limit > UpperLimit {
		limit = UpperLimit
	}
	if limit >= len(peers) {
		return len(peers), peers
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	randomPeers = make([][]byte, limit)
	for random, original := range r.Perm(len(peers)) {
		if random >= limit {
			break
		}
		randomPeers[random] = peers[original]
	}
	return limit, randomPeers
}
