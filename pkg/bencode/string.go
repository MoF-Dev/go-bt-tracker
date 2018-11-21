package bencode

import (
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
func (s String) AsInteger() (int, error) {
	return -1, &IncorrectTypeError{"String", "Integer"}
}
func (s String) AsList() ([]BValue, error) {
	return nil, &IncorrectTypeError{"String", "List"}
}
func (s String) AsDictionary() (map[string]BValue, error) {
	return nil, &IncorrectTypeError{"String", "Dictionary"}
}
