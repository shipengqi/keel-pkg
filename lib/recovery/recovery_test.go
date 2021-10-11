package recovery

import (
	"fmt"
	"testing"
)

func TestRecovery(t *testing.T) {
	var testRecover = func(err error, errStack string) {
		fmt.Printf("[Recovery] panic:\n%s\n%s\n",
			err, errStack)
	}

	defer Recovery(testRecover)()
	t.Log("start panic")
	panic("test panic")
}
