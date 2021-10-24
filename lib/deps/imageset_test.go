package deps

import (
	"os"
	"path/filepath"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/shipengqi/keel-pkg/lib/utils/tmplutil"
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
	set := &Versions{}
	err = jsoniter.Unmarshal(setBytes, set)
	if err != nil {
		t.Fatal(err)
	}
	for i := range set.Components {
		uri, err := tmplutil.ReplaceString(
			set.Components[i].Uri, set.Components[i].Name, &TmplData{Tag: set.Components[i].Tag, Arch: set.Arch})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(uri)
	}
	t.Logf("%+v", set)
}
