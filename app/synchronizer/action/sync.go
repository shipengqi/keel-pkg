package action

import (
	"context"
	"fmt"
	"github.com/docker/docker/registry"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"

	"github.com/shipengqi/keel-pkg/app/synchronizer/pkg/boltdb"
	gcrc "github.com/shipengqi/keel-pkg/app/synchronizer/pkg/registry/gcr/client"
	"github.com/shipengqi/keel-pkg/lib/log"
)

const (
	NameSync = "sync"
)

type SyncOptions struct {
	*gcrc.Options
	Db           string
	PushToRepo   string
	PushToNS     string
	ImageSetFile string
	Limit        int
	QueryLimit   int
	CmdTimeout   time.Duration
}

type reports struct {
	total   int
	success int
	failed  int
	synced  int
}

type synca struct {
	*action

	r          *reports
	opts       *SyncOptions
	gcr        *gcrc.Client
	db         *boltdb.Boltdb
	cancelFunc context.CancelFunc
}

func NewSyncAction(opts *SyncOptions) Interface {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithCancel(context.Background())
	if opts.CmdTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, opts.CmdTimeout)
	}

	db, err := boltdb.New(opts.Db)
	if err != nil {
		panic(err)
	}

	a := &synca{
		action: &action{
			name: NameSync,
			ctx:  ctx,
		},
		r:          &reports{},
		opts:       opts,
		gcr:        gcrc.New(opts.Options),
		db:         db,
		cancelFunc: cancel,
	}

	return a
}

func (s *synca) Close() error {
	log.Debugf("action [%s] closing ...", s.name)
	s.cancelFunc()
	return s.db.Close()
}

func (s *synca) PreRun() error {
	if err := s.auth(); err != nil {
		return err
	}
	if err := s.db.CreatBucket(gcrc.DefaultGcrRepo); err != nil {
		return err
	}
	return nil
}

func (s *synca) Run() (err error) {
	log.Infof("fetch all public images from %s", s.opts.Repo)
	pubs, err := s.gcr.AllImages()
	if err != nil {
		return err
	}

	log.Infof("found images count: %d in %s", len(pubs), gcrc.DefaultGcrRepo)
	images, err := s.fetchImageTagList(pubs)
	if err != nil {
		return err
	}

	log.Debugf("add [%d] patch images", len(s.opts.ImageSet.Patches))
	patches, err := s.fetchPatchTagList()
	if err != nil {
		return err
	}
	if len(patches) > 0 {
		log.Debugf("found [%d] patch image tags", len(patches))
		images = append(images, patches...)
	}
	sort.Sort(images)

	log.Infof("sync images count: %d", len(images))
	err = s.syncImages(images)
	if err != nil {
		return err
	}
	log.Info("sync images done!!!")
	return
}

func (s *synca) fetchPatchTagList() (Images, error) {
	patches := s.opts.ImageSet.Patches
	var images Images
	for i := 0; i < len(patches); i++ {
		patch := patches[i]
		log.Debugf("fetch tags of [%s] ...", patch)
		words := strings.Split(patch, "/")
		if len(words) < 2 || len(words) > 3 {
			log.Warnf("invalid image format [%s]", patch)
			continue
		}
		baseName := strings.Join(words[1:], "/")
		log.Debugf("parse repo [%s], basename [%s]", words[0], baseName)
		tags, err := s.gcr.AllTagsWithRepo(words[0], baseName)
		if err != nil {
			log.Warnf("fetch tags of [%s] failed!", patch)
			continue
		}
		var ns string
		name := baseName
		repo := words[0]
		if len(words) == 3 {
			ns = words[1]
			name = words[2]
		}
		for _, tag := range tags {
			images = append(images, &Image{
				Repo: repo,
				Ns:   ns,
				Name: name,
				Tag:  tag,
			})
		}
	}
	return images, nil
}

func (s *synca) fetchImageTagList(pubs []string) (Images, error) {
	log.Infof("fetch %d images tags from [%s]", len(pubs), s.opts.Repo)

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
				log.Warnf("context done, fetch tags of [%s]", pub)
			default:
				log.Debugf("fetch tags of [%s] ...", pub)
				tags, err := s.gcr.AllTags(pub)
				if err != nil {
					log.Warnf("fetch tags of [%s] failed!", pub)
					return
				}
				log.Debugf("fetch tags count: %d, [%s] ...", len(tags), pub)

				excludeStrs := s.opts.ImageSet.Exclude
				for _, tag := range tags {
					excluded := false
					for sk := range excludeStrs {
						if strings.Contains(tag, excludeStrs[sk]) {
							excluded = true
						}
					}
					if excluded {
						continue
					}
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
	return images, nil
}

func (s *synca) syncImages(images Images) error {
	var success, failed, synced int
	s.r.total = len(images)

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
				log.Warnf("context done, sync image: [%s]", images[k].String())
			default:
				log.Debugf("syncing image: [%s] ...", images[k].String())
				newSum, diff := s.check(images[k])
				if !diff {
					synced++
					return
				}
				err := retry(s.opts.Retry, s.opts.RetryInterval, func() error {
					return s.pushOne(images[k])
				})
				if err != nil {
					failed++
					log.Warnf("sync image [%s]: %s", images[k].String(), err)
					return
				}
				success++
				log.Debugf("sync image: [%s] done!", images[k].String())
				if err := s.db.SaveUint32(images[k].Key(), newSum); err != nil {
					log.Warnf("failed to save image [%s] checksum: %v", images[k].String(), err)
				}
				log.Debugf("save image [%s] checksum: %d", images[k].String(), newSum)
			}
		})
		if err != nil {
			log.Error("submit sync image task failed!")
			return err
		}
	}
	wg.Wait()
	s.r.failed = failed
	s.r.success = success
	s.r.synced = synced
	return nil
}

func (s *synca) PostRun() error {
	report := fmt.Sprintf(`========================================
>> Sync Repo: %s
>> Sync Total: %d
>> Sync Failed: %d
>> Sync Success: %d
>> Synced: %d`, gcrc.DefaultGcrRepo, s.r.total, s.r.failed, s.r.success, s.r.synced)
	fmt.Println(report)
	return nil
}

func (s *synca) check(image *Image) (uint32, bool) {
	var (
		bodySum uint32
		diff    bool
	)
	imgFullName := image.String()
	err := retry(s.opts.Retry, s.opts.RetryInterval, func() error {
		var mErr error
		bodySum, mErr = s.gcr.ManifestCheckSum(imgFullName)
		if mErr != nil {
			return mErr
		}
		if bodySum == 0 {
			return errors.New("checkSum is 0, maybe resp body is nil")
		}
		return nil
	})
	if err != nil {
		log.Errorf("failed to get image [%s] manifest, error: %s", imgFullName, err)
		return 0, false
	}
	diff, err = s.db.Diff(image.Key(), bodySum)
	if err != nil {
		log.Errorf("failed to get image [%s] checkSum, error: %s", imgFullName, err)
		return 0, false
	}
	log.Debugf("[%s] diff: %v", imgFullName, diff)
	if !diff {
		log.Debugf("image [%s] not changed, skip sync ...", imgFullName)
		return 0, false
	}
	return bodySum, true
}

func (s *synca) pushOne(image *Image) error {
	dst := fmt.Sprintf("%s/%s/%s:%s", s.opts.PushToRepo, s.opts.PushToNS, image.Name, image.Tag)
	log.Debugf("syncing [%s] to [%s] ...", image.String(), dst)
	return s.gcr.Sync(image.String(), dst)
}

func (s *synca) auth() error {
	authConf := &types.AuthConfig{
		Username: s.opts.Username,
		Password: s.opts.Password,
	}

	// https://github.com/moby/moby/blob/c3b3aedfa4ad51de0a2ebfd08a716f585390b512/daemon/daemon.go#L714
	// https://github.com/moby/moby/blob/master/daemon/auth.go

	if s.opts.PushToRepo == registry.IndexName {
		authConf.ServerAddress = registry.IndexServer
	} else {
		authConf.ServerAddress = s.opts.PushToRepo
	}
	if !strings.HasPrefix(authConf.ServerAddress, "https://") && !strings.HasPrefix(authConf.ServerAddress, "http://") {
		authConf.ServerAddress = "https://" + authConf.ServerAddress
	}
	service, err := registry.NewService(registry.ServiceOptions{})
	if err != nil {
		return err
	}
	var (
		status      string
		errContains = []string{"imeout", "dead"}
	)
	for count := s.opts.Retry; count > 0; count-- {
		status, _, err = service.Auth(s.ctx, authConf, "")
		if err != nil && contains(errContains, err.Error()) {
			<-time.After(time.Second * 1)
		} else {
			break
		}
	}
	if err != nil {
		return err
	}

	if !strings.Contains(status, "Succeeded") {
		return errors.Errorf("auth: %s", status)
	}
	return nil
}

func contains(s []string, searchterm string) bool {
	for _, v := range s {
		if strings.Contains(searchterm, v) {
			return true
		}
	}
	return false
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
