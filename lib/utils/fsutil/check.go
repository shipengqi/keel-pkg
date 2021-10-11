package fsutil

import (
	"io/ioutil"
	"os"
	"syscall"
)

// PathExists whether the path exists.
func PathExists(fpath string) bool {
	if fpath == "" {
		return false
	}
	if _, err := os.Stat(fpath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// IsDir whether the path is a directory.
func IsDir(fpath string) bool {
	if fpath == "" {
		return false
	}
	if fi, err := os.Stat(fpath); err == nil {
		return fi.IsDir()
	}
	return false
}

// IsSymlink whether the path is a symbolic link.
func IsSymlink(fpath string) bool {
	if fpath == "" {
		return false
	}
	if fi, err := os.Stat(fpath); err == nil {
		return fi.Mode()&os.ModeSymlink != 0
	}
	return false
}

// IsEmpty whether the path is empty.
func IsEmpty(fpath string) bool {
	s, err := ioutil.ReadDir(fpath)
	if err != nil {
		return false
	}
	if len(s) == 0 {
		return true
	}
	return false
}


// FileExists whether the path exists.
func FileExists(fpath string) bool {
	return IsFile(fpath)
}

// IsFile whether the path exists.
func IsFile(fpath string) bool {
	if fpath == "" {
		return false
	}
	if fi, err := os.Stat(fpath); err == nil {
		return !fi.IsDir()
	}
	return false
}

func Owner(fpath string) (uid, gid int, err error) {
	uid = os.Getuid()
	gid = os.Getgid()
	if !FileExists(fpath) {
		return
	}
	info, err := os.Stat(fpath)
	if err != nil {
		return
	}
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		uid = int(stat.Uid)
		gid = int(stat.Gid)
	}
	return
}
