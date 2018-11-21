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
		out.WriteString(String(k).Encode())
		out.WriteString(dict[k].Encode())
	}
	out.WriteString("e")
	return out.String()
}

func (dict Dictionary) AsString() (string, error) {
	return "", &IncorrectTypeError{"Dictionary", "String"}
}
func (dict Dictionary) AsInteger() (int, error) {
	return -1, &IncorrectTypeError{"Dictionary", "Integer"}
}
func (dict Dictionary) AsList() ([]BValue, error) {
	return nil, &IncorrectTypeError{"Dictionary", "List"}
}
func (dict Dictionary) AsDictionary() (map[string]BValue, error) {
	return map[string]BValue(dict), nil
}
