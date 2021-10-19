package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/shipengqi/keel-pkg/app/packer/cmd"
	"github.com/shipengqi/keel-pkg/lib/log"
	"github.com/shipengqi/keel-pkg/lib/recovery"
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
	done := make(chan error, 1)
	go execute(done)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		select {
		case sig := <-quit:
			log.Debugf("get a signal %s", sig.String())
			code = ExitCodeSignal
			return
		case err := <-done:
			code = handleExecuteError(err)
			return
		}
	}
}

func execute(done chan error) {
	var err error
	defer func() {
		done <- err
	}()
	defer recovery.Recovery(mainRecover)()

	c := cmd.New()
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
	exitCode = ExitCodeException
	log.Errorf("Execute(): %+v", err)
	return
}
