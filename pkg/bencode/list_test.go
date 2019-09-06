package bencode

import "testing"

var t0string = "l4:spami42ee"
var t0list List

func init() {
	spam := String("spam")
	forty2 := NewInteger(42)
	t0list = make(List, 2)
	t0list[0] = spam
	t0list[1] = forty2
}
func TestList_Encode_42Spam(t *testing.T) {
	expected := t0string
	actual := t0list.Encode()
	if actual != expected {
		t.Fatalf("Expected %s but got %s", expected, actual)
	}
}

func TestList_decodeList_42Spam(t *testing.T) {
	basicDecodeTest(t, decodeList, t0string[1:], t0list)
}
