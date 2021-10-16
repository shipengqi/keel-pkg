package action

import (
	"context"
	"fmt"
	"sort"
	"sync"
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
}

const (
	NameSync = "sync"
)

func NewSyncAction(opts *SyncOptions) Interface {
	a := &synca{
		action: &action{
			name: NameSync,
		},
		opts: opts,
		gcr:  gcrc.New(opts.Options),
	}
	a.ctx, a.cancel = context.WithCancel(context.Background())
	if opts.CmdTimeout > 0 {
		a.ctx, a.cancel = context.WithTimeout(a.ctx, opts.CmdTimeout)
	}
	return a
}

func (s *synca) Run() (err error) {
	log.Infof("fetch all public images from %s", s.opts.Repo)
	pubs, err := s.gcr.AllImages()
	if err != nil {
		return err
	}
	log.Infof("found images count: %d in k8s.gcr.io", len(pubs))
	images, err := s.fetchImageTagList(pubs)
	if err != nil {
		return err
	}
	sort.Sort(images)
	log.Infof("sync images count: %d", len(images))
	err = s.syncImages(images)
	if err != nil {
		return err
	}
	log.Info("sync images done")
	return
}

func (s *synca) fetchImageTagList(pubs []string) (Images, error) {
	log.Infof("fetch all public images from %s", s.opts.Repo)

	pool, err := ants.NewPool(s.opts.QueryLimit,
		ants.WithPreAlloc(true),
		ants.WithPanicHandler(func(i interface{}) {
			log.Errors(i)
		}))
	if err != nil {
		return nil, errors.Wrap(err, "create pool")
	}
	defer pool.Release()

	var images Images
	imgC := make(chan Image, s.opts.QueryLimit)
	defer close(imgC)

	err = pool.Submit(func() {
		for image := range imgC {
			img := image
			images = append(images, &img)
		}
	})
	if err != nil {
		log.Error("submit image summary task failed!")
		return nil, err
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(pubs))
	for i := 0; i < len(pubs); i++ {
		pub := pubs[i]
		// fullName := fmt.Sprintf("%s/%s", s.opts.Repo, pub)
		err = pool.Submit(func() {
			defer wg.Done()
			select {
			case <-s.ctx.Done():
				log.Warnf("context done, fetch tags of %s", pub)
			default:
				log.Debugf("fetch tags of %s ...", pub)
				tags, err := s.gcr.AllTags(pub)
				if err != nil {
					log.Warnf("fetch tags of %s failed!", pub)
					return
				}
				log.Debugf("fetch tags count: %d, %s ...", len(tags), pub)
				for _, tag := range tags {
					imgC <- Image{
						Name: pub,
						Tag:  tag,
					}
				}
			}
		})
		if err != nil {
			log.Error("submit fetch image tags task failed!")
			return nil, err
		}
	}
	wg.Wait()
	log.Infof("fetched all images, total: %d", len(images))
	return images, nil
}

func (s *synca) syncImages(images Images) error {
	wg := new(sync.WaitGroup)
	wg.Add(len(images))

	pool, err := ants.NewPool(s.opts.Limit,
		ants.WithPreAlloc(true),
		ants.WithPanicHandler(func(i interface{}) {
			log.Errors(i)
		}))
	if err != nil {
		return errors.Wrap(err, "create pool")
	}
	defer pool.Release()

	for i := 0; i < len(images); i++ {
		k := i
		err = pool.Submit(func() {
			defer wg.Done()
			select {
			case <-s.ctx.Done():
				log.Warnf("context done, sync image: %s", images[k].String())
			default:
				log.Debugf("syncing image: %s ...", images[k].String())
				err := retry(s.opts.Retry, s.opts.RetryInterval, func() error {
					return s.pushOne(images[k])
				})
				if err != nil {
					log.Warnf("sync image %s: %s", images[k].String(), err)
					return
				}
				images[k].Success = true
				log.Debugf("sync image: %s done", images[k].String())
			}
		})
		if err != nil {
			log.Error("submit sync image task failed!")
			return err
		}
	}
	wg.Wait()
	return nil
}

func (s *synca) pushOne(image *Image) error {
	dst := fmt.Sprintf("%s/%s/%s:%s", s.opts.PushToRepo, s.opts.PushToNS, image.Name, image.Tag)
	log.Debugf("syncing %s to %s ...", image.String(), dst)
	return nil
	// return s.gcr.Sync(image.String(), dst)
}

func retry(count int, interval time.Duration, f func() error) error {
	var err error
	for ; count > 0; count-- {
		if err = f(); err != nil {
			if interval > 0 {
				<-time.After(interval)
			}
		} else {
			break
		}
	}
	return err
}
