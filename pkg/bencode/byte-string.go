package bencode

import (
	"strconv"
	"strings"
)

type ByteString string

func (bs ByteString) Encode() string {
	var out strings.Builder
	out.WriteString(strconv.Itoa(len(bs)))
	out.WriteByte(':')
	out.WriteString(string(bs))
	return out.String()
}