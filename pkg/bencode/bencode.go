package bencode

type BValue interface {
	Encode() string
	AsString() (string, error)
	AsInteger() (int, error)
	AsList() ([]BValue, error)
	AsDictionary() (map[string]BValue, error)
}

