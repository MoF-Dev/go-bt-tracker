package bencode

import (
	"sort"
	"strings"
)

type Dictionary map[string]BValue

func (dict Dictionary) Encode() string {
	var out strings.Builder
	keys := make([]string, 0)
	for k := range dict {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out.WriteString("d")
	for _, k := range keys {
		out.WriteString(ByteString(k).Encode())
		out.WriteString(dict[k].Encode())
	}
	out.WriteString("e")
	return out.String()
}
