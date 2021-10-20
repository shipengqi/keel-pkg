package cmd

import (
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/shipengqi/keel-pkg/app/synchronizer/action"
	"github.com/shipengqi/keel-pkg/lib/log"
	"github.com/shipengqi/keel-pkg/lib/utils/fmtutil"
)

const (
	DefaultPushToRegistry = "registry.cn-hangzhou.aliyuncs.com"
	DefaultPushToNs       = "keel"
	DefaultDBFile         = "sync.bolt.db"
	DefaultQueryLimit     = 10
	DefaultLimit          = 5
	DefaultRetryCount     = 5
)

var (
	Version   = "<unknown>"
	GitCommit = "<unknown>"
	BuildTime = "1970-01-01T00:00:00Z"
)

func New(done chan error) *cobra.Command {
	baseName := filepath.Base(os.Args[0])
	c := &cobra.Command{
		Use:     baseName + " [options]",
		Short:   "keel image synchronizer",
		Version: Version,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	c.AddCommand(
		NewSyncCommand(done),
		NewCheckCommand(done),
		NewSumCommand(done),
	)
	cobra.EnableCommandSorting = false
	c.CompletionOptions.DisableDefaultCmd = true
	c.DisableFlagsInUseLine = true
	c.SetVersionTemplate(fmtutil.VersionTmpl("Keel Synctl",
		Version, GitCommit, BuildTime))
	return c
}

func receiver(a action.Interface, done chan error) {
	var once sync.Once
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	select {
	case sig := <-quit:
		once.Do(func() {
			log.Warnf("action [%s], get a signal [%s], quit", a.Name(), sig.String())
			_ = a.Close()
		})
		log.Warnf("action [%s] stopped!", a.Name())
		done <- errors.Errorf("signal [%s]", sig.String())
		return
	}
}
