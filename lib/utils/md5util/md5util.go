package md5util

import (
	"crypto/md5"
	"encoding/hex"
)

func Encode(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}
