package bencode

import "strings"

type List []BValue

func (list List) Encode() string {
	var out strings.Builder
	out.WriteString("l")
	for _, e := range list {
		out.WriteString(e.Encode())
	}
	out.WriteString("e")
	return out.String()
}
