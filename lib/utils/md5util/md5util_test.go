package md5util

import "testing"

func TestEncode(t *testing.T) {
	text := "hello"
	encoded1 := Encode(text)
	encoded2 := Encode(text)
	if encoded1 != encoded2 {
		t.Fatal(encoded1, encoded2)
	}
	t.Log(encoded1)
	t.Log(encoded2)
}
