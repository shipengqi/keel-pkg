package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTar(t *testing.T) {
	p, _ := os.Getwd()
	datap := filepath.Join(p, "../../../", "test/testdata/utils")
	t.Log(datap)
	dstp := filepath.Join(datap, "test.tar.gz")
	t.Log(dstp)
	err := Tar(dstp, datap)
	if err != nil {
		t.Fatal(err)
	}
	_ = os.Remove(dstp)
}

func TestUnTar(t *testing.T) {
	p, _ := os.Getwd()
	datap := filepath.Join(p, "../../../", "test/testdata/utils")
	t.Log(datap)
	tdstp := filepath.Join(datap, "test.tar.gz")
	t.Log(tdstp)
	err := Tar(tdstp, datap)
	if err != nil {
		t.Fatal(err)
	}
	utdstp := filepath.Join(datap, "testuntar")
	t.Log(utdstp)
	err = UnTar(utdstp, tdstp)
	if err != nil {
		t.Fatal(err)
	}
	_ = os.Remove(tdstp)
	_ = os.Remove(utdstp)
}

func TestMustCopyDir(t *testing.T) {
	p, _ := os.Getwd()
	datap := filepath.Join(p, "../../../", "test/testdata/utils/tardata")
	utdstp := filepath.Join(p, "../../../", "test/testdata/utils/copieddata")
	t.Log(datap)
	t.Log(utdstp)
	err := MustCopyDir(utdstp, datap)
	if err != nil {
		t.Fatal(err)
	}

	_ = os.Remove(utdstp)
}

func TestCleanDir(t *testing.T) {
	p, _ := os.Getwd()
	datap := filepath.Join(p, "../../../", "test/testdata/utils/cleandir")
	// _ = MustMkDir(filepath.Join(datap, "testdir1"))
	// _ = MustMkDir(filepath.Join(datap, "testdir2"))
	err := CleanDir(datap)
	if err != nil {
		t.Fatal(err)
	}
}
