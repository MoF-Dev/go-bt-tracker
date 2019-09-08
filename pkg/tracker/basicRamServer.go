package tracker

import (
	"math/rand"
	"time"
)

type hue struct {
}

func (s *hue) ChooseLimitedPeers(peers [][]byte, n int32) (limit int, randomPeers [][]byte) {
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
	randomPeers = make([]Peer, limit)
	for random, original := range r.Perm(len(peers)) {
		if random >= limit {
			break
		}
		randomPeers[random] = peers[original]
	}
	return limit, randomPeers
}
