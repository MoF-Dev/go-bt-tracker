package bencode

import "testing"

// Wikipedia examples

func TestInteger_Encode_Positive(t *testing.T) {
	expected := "i42e"
	actual := Integer(42).Encode()
	if actual != expected {
		t.Fatalf("Expected %s but got %s", expected, actual)
	}
}
func TestInteger_Encode_Zero(t *testing.T) {
	expected := "i0e"
	actual := Integer(0).Encode()
	if actual != expected {
		t.Fatalf("Expected %s but got %s", expected, actual)
	}
}
func TestInteger_Encode_Negative(t *testing.T) {
	expected := "i-42e"
	actual := Integer(-42).Encode()
	if actual != expected {
		t.Fatalf("Expected %s but got %s", expected, actual)
	}
}
