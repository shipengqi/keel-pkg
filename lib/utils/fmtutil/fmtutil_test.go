package fmtutil

import "testing"

func TestPretty(t *testing.T) {
	s := Pretty("fmtutil testing", "ok")
	t.Log(s)
}

func TestPrettyWithColor(t *testing.T) {
	s := PrettyWithColor("fmtutil testing", "ok", FgGreen)
	t.Log(s)
}
