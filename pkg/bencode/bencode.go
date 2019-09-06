package bencode

import (
	"github.com/MoF-Dev/go-bt-tracker/pkg/bencode/types"
	"math/big"
)

type BValue interface {
	Encode() string
	AsString() (string, error)
	AsInteger() (*big.Int, error)
	AsList() ([]BValue, error)
	AsDictionary() (map[string]BValue, error)
	Type() types.Type
	Equals(another BValue) bool
}
