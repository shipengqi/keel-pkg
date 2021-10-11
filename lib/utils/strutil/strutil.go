package strutil

import "regexp"

func ReplaceSpace(src, new string) string {
	return Replace(src, "\\s+", new)
}

func Replace(src, old, new string) string {
	if src == "" {
		return ""
	}
	reg := regexp.MustCompile(old)
	return reg.ReplaceAllString(src, new)
}
