package deps

import (
	jsoniter "github.com/json-iterator/go"
	"os"
	"path/filepath"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	p, _ := os.Getwd()
	datap := filepath.Join(p, "../../", "image_set.json")
	t.Log(datap)
	setBytes, err := os.ReadFile(datap)
	if err != nil {
		t.Fatal(err)
	}
	set := &ImageSet{}
	err = jsoniter.Unmarshal(setBytes, set)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", set)
}

func TestUnmarshal2(t *testing.T) {
	p, _ := os.Getwd()
	datap := filepath.Join(p, "../../", "versions.json")
	t.Log(datap)
	setBytes, err := os.ReadFile(datap)
	if err != nil {
		t.Fatal(err)
	}
	set := &Version{}
	err = jsoniter.Unmarshal(setBytes, set)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", set)
}
