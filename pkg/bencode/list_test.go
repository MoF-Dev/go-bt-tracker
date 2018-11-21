package bencode

import "testing"

func TestList_Encode_42Spam(t *testing.T) {
	expected := "l4:spami42ee"
	spam := String("spam")
	forty2 := Integer(42)
	list := List(make([]BValue, 2))
	list[0] = spam
	list[1] = forty2
	actual := list.Encode()
	if actual != expected {
		t.Fatalf("Expected %s but got %s", expected, actual)
	}
}
