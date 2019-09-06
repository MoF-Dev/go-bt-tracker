package bencode

import "testing"

var wiki0TestString = "d3:bar4:spam3:fooi42ee"
var wiki0TestDict Dictionary

func init() {
	wiki0TestDict = make(Dictionary)
	wiki0TestDict["foo"] = NewInteger(42)
	wiki0TestDict["bar"] = String("spam")
}

func TestDictionary_Encode_Wiki(t *testing.T) {
	actual := wiki0TestDict.Encode()
	if actual != wiki0TestString {
		t.Fatalf("Expected %s but got %s", wiki0TestString, actual)
	}
}

func TestDictionary_decodeDictionary_Wiki(t *testing.T) {
	basicDecodeTest(t, decodeDictionary, wiki0TestString[1:], wiki0TestDict)
}
