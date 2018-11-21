package bencode

import (
	"errors"
	"fmt"
)

type IncorrectTypeError struct {
	have string
	want string
}

func (e *IncorrectTypeError) Error() string {
	return fmt.Sprintf("Incorrect bencode type: requested %s, got %s", e.want, e.have)
}

func Decode(str string) (BValue, error) {
	// TODO
	return nil, errors.New("not yet implemented")
}
