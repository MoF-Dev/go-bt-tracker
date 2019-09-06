package bencode

import "testing"

type decodeFunc func(string) (BValue, string, error)

func basicDecodeTest(t *testing.T, decoder decodeFunc, input string, expect BValue) {
	actual, str, err := decoder(input)
	if err != nil {
		t.Fatalf("error %s", err)
	}
	if str != "" {
		t.Fatalf("leftover unparsed string: %s", str)
	}
	if !expect.Equals(actual) {
		t.Fatalf("Expected %s but got %s", expect, actual)
	}
}
