package main

import (
	"github.com/shipengqi/keel-pkg/app/synchronizer/cmd"
	"github.com/shipengqi/keel-pkg/lib/log"
	"github.com/shipengqi/keel-pkg/lib/recovery"
	"os"
	"strings"
)

const (
	ExitCodeOk        = 0
	ExitCodeException = 1
	ExitCodePanic     = 2
	ExitCodeSignal    = 3
	ExitCodeSkipLog   = 4
)

const (
	reset = "\033[0m"
)

var code int

func main() {
	defer func() {
		os.Exit(code)
	}()

	done := make(chan error, 0)
	go execute(done)

	select {
	case err := <-done:
		code = handleExecuteError(err)
		return
	}
}

func execute(done chan error) {
	var err error
	defer func() {
		done <- err
	}()
	defer recovery.Recovery(mainRecover)()

	c := cmd.New(done)
	err = c.Execute()
	return
}

func mainRecover(err error, errStack string) {
	log.Errorf("[Recovery] panic:\n%s\n%s%s",
		err, errStack, reset)
	code = ExitCodePanic
}

func handleExecuteError(err error) (exitCode int) {
	if err == nil {
		exitCode = ExitCodeOk
		return
	}
	if strings.HasPrefix(err.Error(), "unknown command") {
		log.Warn(err.Error())
		exitCode = ExitCodeSkipLog
		return
	}
	if strings.HasPrefix(err.Error(), "unknown flag") {
		log.Warn(err.Error())
		exitCode = ExitCodeSkipLog
		return
	}
	if strings.HasPrefix(err.Error(), "signal") {
		log.Warn(err.Error())
		exitCode = ExitCodeSignal
		return
	}
	exitCode = ExitCodeException
	log.Errorf("Execute(): %+v", err)
	return
}
