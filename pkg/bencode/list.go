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

func (list List) AsString() (string, error) {
	return "", &IncorrectTypeError{"List", "String"}
}
func (list List) AsInteger() (int, error) {
	return -1, &IncorrectTypeError{"List", "Integer"}
}
func (list List) AsList() ([]BValue, error) {
	return []BValue(list), nil
}
func (list List) AsDictionary() (map[string]BValue, error) {
	return nil, &IncorrectTypeError{"List", "Dictionary"}
}
