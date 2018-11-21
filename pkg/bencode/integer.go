package bencode

import "strconv"

type Integer int

func (i Integer) Encode() string {
	return "i" + strconv.Itoa(int(i)) + "e"
}

func (i Integer) AsString() (string, error) {
	return strconv.Itoa(int(i)), &IncorrectTypeError{"Integer", "String"}
}
func (i Integer) AsInteger() (int, error) {
	return int(i), nil
}
func (i Integer) AsList() ([]BValue, error) {
	return nil, &IncorrectTypeError{"Integer", "List"}
}
func (i Integer) AsDictionary() (map[string]BValue, error) {
	return nil, &IncorrectTypeError{"Integer", "Dictionary"}
}
