package bencode

import (
	"github.com/MoF-Dev/go-bt-tracker/pkg/bencode/types"
	"math/big"
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
		out.WriteString(String(k).Encode())
		out.WriteString(dict[k].Encode())
	}
	out.WriteString("e")
	return out.String()
}

func (dict Dictionary) AsString() (string, error) {
	return "", &IncorrectTypeError{"Dictionary", "String"}
}
func (dict Dictionary) AsInteger() (*big.Int, error) {
	return big.NewInt(0), &IncorrectTypeError{"Dictionary", "Integer"}
}
func (dict Dictionary) AsList() ([]BValue, error) {
	return nil, &IncorrectTypeError{"Dictionary", "List"}
}
func (dict Dictionary) AsDictionary() (map[string]BValue, error) {
	return map[string]BValue(dict), nil
}
func (Dictionary) Type() types.Type {
	return types.Dictionary
}
func (dict Dictionary) Equals(another BValue) bool {
	if another == nil || another.Type() != types.Dictionary {
		return false
	}
	a := another.(Dictionary)
	if len(dict) != len(a) {
		return false
	}
	for k, v := range dict {
		if !v.Equals(a[k]) {
			return false
		}
	}
	return true
}

func decodeDictionary(str string) (BValue, string, error) {
	var dict Dictionary
	dict = make(Dictionary)

	var key BValue
	var value BValue
	var err error
	for {
		if str[0] == 'e' {
			return dict, str[1:], nil
		}
		key, str, err = decodeString(str)
		if err != nil {
			return dict, str, err
		}
		value, str, err = decode(str)
		if err != nil {
			return dict, str, err
		}
		dictKey, _ := key.AsString()
		dict[dictKey] = value
	}
}
