package bencode

import "testing"

// Wikipedia examples

func TestInteger_Encode_Positive(t *testing.T) {
	expected := "i42e"
	actual := NewInteger(42).Encode()
	if actual != expected {
		t.Fatalf("Expected %s but got %s", expected, actual)
	}
}
func TestInteger_Encode_Zero(t *testing.T) {
	expected := "i0e"
	actual := NewInteger(0).Encode()
	if actual != expected {
		t.Fatalf("Expected %s but got %s", expected, actual)
	}
}
func TestInteger_Encode_Negative(t *testing.T) {
	expected := "i-42e"
	actual := NewInteger(-42).Encode()
	if actual != expected {
		t.Fatalf("Expected %s but got %s", expected, actual)
	}
}

func TestInteger_decodeInteger_Positive(t *testing.T) {
	basicDecodeTest(t, decodeInteger, "42e", NewInteger(42))
}
func TestInteger_decodeInteger_Zero(t *testing.T) {
	basicDecodeTest(t, decodeInteger, "0e", NewInteger(0))
}
func TestInteger_decodeInteger_Negative(t *testing.T) {
	basicDecodeTest(t, decodeInteger, "-42e", NewInteger(-42))
}
