package bencode

import "strconv"

type Integer int

func (i Integer) Encode() string {
	return "i" + strconv.Itoa(int(i)) + "e"
}
