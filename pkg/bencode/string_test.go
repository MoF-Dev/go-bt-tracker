package bencode

import "testing"

func TestString_Encode(t *testing.T) {
	expected := "4:spam"
	actual := String("spam").Encode()
	if actual != expected {
		t.Fatalf("Expected %s but got %s", expected, actual)
	}
}
