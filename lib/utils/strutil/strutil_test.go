package strutil

import "testing"

func TestReplaceSpace(t *testing.T) {
	text := "/dev/mapper/centos7-root /boot                   xfs     defaults        0 0"
	expected := "/dev/mapper/centos7-root /boot xfs defaults 0 0"
	got := ReplaceSpace(text,  " ")
	if got != expected {
		t.Fatalf("expected: %s, got: %s", expected, got)
	}
	t.Logf("expected: %s, got: %s", expected, got)
}
