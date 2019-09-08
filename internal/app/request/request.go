package request

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

type ScrapeRequest struct {
	InfoHashes [][20]byte
}
