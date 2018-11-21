package bencode

import "testing"

func TestDictionary_Encode_Wiki(t *testing.T) {
	expected := "d3:bar4:spam3:fooi42ee"
	dict := Dictionary(make(map[string]BValue))
	dict["foo"] = Integer(42)
	dict["bar"] = String("spam")
	actual := dict.Encode()
	if actual != expected {
		t.Fatalf("Expected %s but got %s", expected, actual)
	}
}
