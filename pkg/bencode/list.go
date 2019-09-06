package bencode

import (
	"github.com/MoF-Dev/go-bt-tracker/pkg/bencode/types"
	"math/big"
	"strings"
)

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

func (list List) AsString() (string, error) {
	return "", &IncorrectTypeError{"List", "String"}
}
func (list List) AsInteger() (*big.Int, error) {
	return big.NewInt(0), &IncorrectTypeError{"List", "Integer"}
}
func (list List) AsList() ([]BValue, error) {
	return []BValue(list), nil
}
func (list List) AsDictionary() (map[string]BValue, error) {
	return nil, &IncorrectTypeError{"List", "Dictionary"}
}
func (List) Type() types.Type {
	return types.List
}
func (list List) Equals(another BValue) bool {
	if another == nil || another.Type() != types.List {
		return false
	}
	a, _ := another.(List)
	if len(list) != len(a) {
		return false
	}
	for i, v := range list {
		if !v.Equals(a[i]) {
			return false
		}
	}
	return true
}

func decodeList(str string) (BValue, string, error) {
	var list List

	// is this necessary?
	list = make(List, 0)

	var item BValue
	var err error
	for {
		if str[0] == 'e' {
			return list, str[1:], nil
		}
		item, str, err = decode(str)
		if err != nil {
			return list, str, nil
		}
		list = append(list, item)
	}
}
