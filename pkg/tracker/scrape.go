package tracker

import "github.com/MoF-Dev/go-bt-tracker/pkg/bencode"

type ScrapeRequest struct {
	InfoHashes [][20]byte
}

type ScrapeResponse struct {
	Files         []File
	FailureReason *string
	Flags         *bencode.Dictionary
}

type File struct {
	Completed  uint32
	Downloaded uint32
	Incomplete uint32
	Name       *string
}
