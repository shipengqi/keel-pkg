package fsutil

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Name return file/dir name
func Name(fpath string) string {
	return filepath.Base(fpath)
}

// Dir return the directory path without base name.
func Dir(fpath string) string {
	return filepath.Dir(fpath)
}

func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func MustMkDir(fpath string) error {
	if exist := PathExists(fpath); !exist {
		if err := MkDirAll(fpath); err != nil {
			return err
		}
	}

	return nil
}

func MkDirAll(fpath string) error {
	return os.MkdirAll(fpath, os.ModePerm)
}

func UnTar(dst, src string) (err error) {
	fr, err := os.Open(src)
	if err != nil {
		return
	}
	defer fr.Close()

	// uncompress
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, header.Name)
		switch header.Typeflag {
		case tar.TypeDir: // directory
			if err = MustMkDir(dstPath); err != nil {
				return err
			}
		case tar.TypeReg: // file
			file, err := os.OpenFile(dstPath, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			_, err = io.Copy(file, tr)
			if err != nil {
				return err
			}
			file.Close()
		}
	}

	return
}

func Tar(dst, src string) error {
	fw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fw.Close()

	gw := gzip.NewWriter(fw)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()
	return filepath.Walk(src, func(filename string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(filename, string(filepath.Separator))
		// write file info
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// whether info describes a regular file.
		if !info.Mode().IsRegular() {
			return nil
		}
		fr, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer fr.Close()
		_, err = io.Copy(tw, fr)
		if err != nil {
			return err
		}
		return nil
	})
}

func CopyFile(dst, src string) (err error) {
	srcFd, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFd.Close()

	dstFd, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer dstFd.Close()

	if _, err = io.Copy(dstFd, srcFd); err != nil {
		return err
	}
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, info.Mode())
}

func MustCopyDir(dst, src string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	if !srcinfo.IsDir() {
		return CopyFile(dst, src)
	}
	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}
	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = MustCopyDir(dstfp, srcfp); err != nil {
				return err
			}
		} else {
			if err = CopyFile(dstfp, srcfp); err != nil {
				return err
			}
		}
	}
	return nil
}

func Scanner(file string) (*bufio.Scanner, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	s := bufio.NewScanner(fd)
	return s, err
}


func CleanDir(fpath string) error {
	entries, err := os.ReadDir(fpath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		err = os.RemoveAll(filepath.Join(fpath, entry.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}
