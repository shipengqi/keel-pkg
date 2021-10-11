package tmplutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReplaceString(t *testing.T) {
	var testTemplate = `
  ./{{.}} cert renew -V 365         Renew the certificates.
  ./{{.}} cert create -V 365        Create the certificates.
  ./{{.}} cert apply                Apply the certificates.`

	s, err := ReplaceString(testTemplate, "test1", "keel")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}

func TestReplaceFile(t *testing.T) {
	p, _ := os.Getwd()
	templp := filepath.Join(p, "../../../", "test/testdata/templs/containerd.service")
	t.Log(templp)
	data := struct {
		KeelRuntimeDataHome string
		KeelHome            string
	}{"/opt/keel/data", "/opt/keel"}
	s, err := ReplaceFile(templp, data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)

	data2 := struct {
		Host    string
		OrgName string
		ImgName string
	}{"localhost", "keel", "test:v1"}
	templp = filepath.Join(p, "../../../", "test/testdata/templs/kube.yaml")
	t.Log(templp)
	s, err = ReplaceFile(templp, data2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}

func TestReplaceAndWriteFile(t *testing.T) {
	p, _ := os.Getwd()
	templp := filepath.Join(p, "../../../", "test/testdata/templs/containerd.service")
	output := filepath.Join(p, "../../../", "test/testdata/templs/output/containerd.service")
	t.Log(templp)
	data := struct {
		KeelRuntimeDataHome string
		KeelHome            string
	}{"/opt/keel/data", "/opt/keel"}
	err := ReplaceAndWriteFile(templp, output, data)
	if err != nil {
		t.Fatal(err)
	}
	_ = os.Remove(output)

	data2 := struct {
		Host    string
		OrgName string
		ImgName string
	}{"localhost", "keel", "test:v1"}
	templp = filepath.Join(p, "../../../", "test/testdata/templs/kube.yaml")
	output = filepath.Join(p, "../../../", "test/testdata/templs/output/kube.yaml")
	t.Log(templp)
	err = ReplaceAndWriteFile(templp, output, data2)
	if err != nil {
		t.Fatal(err)
	}
	_ = os.Remove(output)
}
