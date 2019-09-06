package bencode

import (
	"errors"
	"fmt"
	"github.com/MoF-Dev/go-bt-tracker/pkg/bencode/types"
	"math/big"
	"strings"
)

type Integer struct {
	*big.Int
}

func NewInteger(i int64) Integer {
	return Integer{big.NewInt(i)}
}
func NewUInteger(ui uint64) Integer {
	return Integer{big.NewInt(0).SetUint64(ui)}
}

func (i Integer) Encode() string {
	return "i" + i.Text(10) + "e"
}

func (i Integer) AsString() (string, error) {
	return i.Text(10), &IncorrectTypeError{"Integer", "String"}
}
func (i Integer) AsInteger() (*big.Int, error) {
	return i.Int, nil
}
func (i Integer) AsList() ([]BValue, error) {
	return nil, &IncorrectTypeError{"Integer", "List"}
}
func (i Integer) AsDictionary() (map[string]BValue, error) {
	return nil, &IncorrectTypeError{"Integer", "Dictionary"}
}
func (Integer) Type() types.Type {
	return types.Integer
}
func (i Integer) Equals(another BValue) bool {
	if another == nil || another.Type() != types.Integer {
		return false
	}
	a, _ := another.(Integer)
	return i.Int.Cmp(a.Int) == 0
}

func decodeInteger(str string) (BValue, string, error) {
	e := strings.IndexByte(str, 'e')
	if e == 0 {
		return nil, str, errors.New("parsing 0 length integer")
	}
	if e <= 0 {
		return nil, str, errors.New("parsing integer without ending e")
	}
	if i, success := big.NewInt(0).SetString(str[:e], 10); success {
		return Integer{i}, str[e+1:], nil
	} else {
		return nil, str[e:], fmt.Errorf("parsing integer error: %s", str[:e])
	}
}
