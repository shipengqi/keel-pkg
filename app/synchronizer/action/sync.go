package action

import (
	"context"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"

	gcrc "github.com/shipengqi/keel-pkg/app/synchronizer/pkg/registry/gcr/client"
	"github.com/shipengqi/keel-pkg/lib/log"
)

type SyncOptions struct {
	*gcrc.Options
	Db         string
	PushToRepo string
	PushToNS   string
	Limit      int
	QueryLimit int
	CmdTimeout time.Duration
}

type synca struct {
	*action

	opts *SyncOptions
	gcr  *gcrc.Client
	ctx  context.Context
}

func NewSyncAction(opts *SyncOptions) *synca {
	return &synca{
		opts: opts,
		gcr:  gcrc.New(opts.Options),
	}
}

func (s *synca) PreRun() (err error) {
	var cancel context.CancelFunc
	s.ctx, cancel = context.WithCancel(context.Background())
	if s.opts.CmdTimeout > 0 {
		s.ctx, cancel = context.WithTimeout(s.ctx, s.opts.CmdTimeout)
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
	log.Infof("fetch all public images from %s", s.opts.Repo)
	pubs, err := s.gcr.AllImages()
	if err != nil {
		return err
	}
	log.Infof("found images count: %d in k8s.gcr.io", len(pubs))


	sort.Sort(s.imgs)

	for i := 0; i < len(s.imgs); i++ {

	}
	return
}

func (s *synca) fetchImageTagList(pubs []string) ([]string, error) {
	log.Infof("fetch all public images from %s", s.opts.Repo)
	processWg := new(sync.WaitGroup)
	processWg.Add(len(pubs))

	pool, err := ants.NewPool(s.opts.QueryLimit,
		ants.WithPreAlloc(true),
		ants.WithPanicHandler(func(i interface{}) {
			log.Errors(i)
		}))
	if err != nil {
		return nil, errors.Wrap(err, "create pool")
	}
	return pubs, nil
}
