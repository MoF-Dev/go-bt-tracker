package bencode

type BValue interface {
	Encode() string
}