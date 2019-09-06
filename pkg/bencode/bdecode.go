package bencode

import (
	"fmt"
)

type IncorrectTypeError struct {
	have string
	want string
}

func (e *IncorrectTypeError) Error() string {
	return fmt.Sprintf("Incorrect bencode type: requested %s, got %s", e.want, e.have)
}

func decode(str string) (BValue, string, error) {
	switch str[0] {
	case 'i':
		return decodeInteger(str[1:])
	case 'l':
		return decodeList(str[1:])
	case 'd':
		return decodeDictionary(str[1:])
	default:
		return decodeString(str)
	}
}
