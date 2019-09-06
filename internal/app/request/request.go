package request

import "math/big"

type Event int

const (
	Empty Event = iota
	Started
	Completed
	Stopped
)

type AnnounceRequest struct {
	InfoHash   [20]byte
	PeerId     [20]byte
	Ip         *string
	Port       uint16
	Uploaded   *big.Int
	Downloaded *big.Int
	Left       *big.Int
	Event      Event
}

type ScrapeRequest struct {
	InfoHashes [][20]byte
}
