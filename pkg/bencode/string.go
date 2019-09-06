package bencode

import (
	"errors"
	"github.com/MoF-Dev/go-bt-tracker/pkg/bencode/types"
	"math/big"
	"strconv"
	"strings"
)

type String string

func (s String) Encode() string {
	var out strings.Builder
	out.WriteString(strconv.Itoa(len(s)))
	out.WriteByte(':')
	out.WriteString(string(s))
	return out.String()
}

func (s String) AsString() (string, error) {
	return string(s), nil
}
func (s String) AsInteger() (*big.Int, error) {
	return big.NewInt(0), &IncorrectTypeError{"String", "Integer"}
}
func (s String) AsList() ([]BValue, error) {
	return nil, &IncorrectTypeError{"String", "List"}
}
func (s String) AsDictionary() (map[string]BValue, error) {
	return nil, &IncorrectTypeError{"String", "Dictionary"}
}
func (String) Type() types.Type {
	return types.String
}
func (s String) Equals(another BValue) bool {
	if another == nil || another.Type() != types.String {
		return false
	}
	return s == another
}

func decodeString(str string) (BValue, string, error) {
	e := int64(strings.IndexByte(str, ':'))
	if e < 0 {
		return nil, str, errors.New("string length not specified")
	}
	len, err := strconv.ParseInt(str[:e], 10, 64)
	if err != nil {
		return nil, str, err
	}
	return String(str[e+1 : e+len+1]), str[e+len+1:], nil
}
