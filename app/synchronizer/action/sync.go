package action

import (
	"context"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/types"
	"github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"

	"github.com/shipengqi/keel-pkg/lib/log"
)

type SyncOptions struct {
	Username    string
	Password    string
	Registry    string
	Namespace   string
	DbFile      string
	Concurrency int
	LoginRetry  int
	SyncRetry   int
	Timeout     time.Duration
}

type synca struct {
	ctx     context.Context
	imgs    []string
	options *SyncOptions
}

func (s *synca) PreRun() (err error) {
	var cancel context.CancelFunc
	s.ctx, cancel = context.WithCancel(context.Background())
	if s.options.Timeout > 0 {
		s.ctx, cancel = context.WithTimeout(s.ctx, s.options.Timeout)
	}
	var cancelOnce sync.Once
	defer cancel()
	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		for {
			select {
			case sig := <-quit:
				log.Debugf("get a signal %s", sig.String())
				cancelOnce.Do(func() {
					log.Info("Shutdown!")
					cancel()
				})
			}
		}
	}()

	return
}

func (s *synca) Run() (err error) {
	processWg := new(sync.WaitGroup)
	processWg.Add(len(s.imgs))

	pool, err := ants.NewPool(s.options.Concurrency, ants.WithPreAlloc(true), ants.WithPanicHandler(func(i interface{}) {
		log.Error("WithPanicHandler")
	}))
	if err != nil {
		return errors.Wrap(err, "create goroutines pool")
	}
	sort.Sort(s.imgs)

	for i := 0; i < len(s.imgs); i++ {

	}
	return
}
