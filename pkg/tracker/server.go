package tracker

type Server interface {
	HandleAnnounce(request *AnnounceRequest) (*AnnounceResponse, error)
	HandleScrape(request *ScrapeRequest) (*ScrapeResponse, error)
	ChooseLimitedPeers(peers [][]byte, n int32) (limit int, randomPeers [][]byte)
}

type HttpUdpServer interface {
	Server
	udpServerExtras
	httpServerExtras
}
