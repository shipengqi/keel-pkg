package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/types"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/shipengqi/keel-pkg/lib/log"
)

const (
	DefaultGcrRepo       = "k8s.gcr.io"
	DefaultRetryTimes    = 5
	DefaultRetryInterval = time.Second * 5
	DefaultTimeout       = time.Second * 20
	DefaultPushTimeout   = time.Minute * 15
)

const (
	_defaultGcrImagesListAPI = "https://k8s.gcr.io/v2/tags/list"
)

type Options struct {
	Username      string
	Password      string
	Token         string
	Repo          string
	Retry         int
	RetryInterval time.Duration
	ReqTimeout    time.Duration
	PushTimeout   time.Duration
	Ctx           context.Context
	AdditionalNS  []string
}

func NewDefaultOptions() *Options {
	return &Options{
		Repo:          DefaultGcrRepo,
		Retry:         DefaultRetryTimes,
		RetryInterval: DefaultRetryInterval,
		ReqTimeout:    DefaultTimeout,
		PushTimeout:   DefaultPushTimeout,
	}
}

type Client struct {
	*resty.Client

	opts *Options
}

func New(options *Options) *Client {
	r := resty.New()
	r.SetTimeout(options.ReqTimeout)
	r.SetRetryCount(options.Retry)
	r.SetRetryWaitTime(options.RetryInterval)
	if options.Ctx == nil {
		options.Ctx = context.Background()
	}
	return &Client{
		Client: r,
		opts:   options,
	}
}

func (c *Client) AllImages() ([]string, error) {
	var allBaseNames []string
	var err error

	allBaseNames, err = c.allImages(_defaultGcrImagesListAPI)
	if err != nil {
		return nil, err
	}
	for i := range c.opts.AdditionalNS {
		baseNames, err := c.allImages(fmt.Sprintf("https://k8s.gcr.io/v2/%s/tags/list", c.opts.AdditionalNS[i]))
		if err != nil {
			log.Warnf("Additional namespace error: %v", err)
			continue
		}
		allBaseNames = append(allBaseNames, baseNames...)
	}
	return allBaseNames, nil
}

func (c *Client) AllTags(baseName string) ([]string, error) {
	imageName := fmt.Sprintf("%s/%s", c.opts.Repo, baseName)
	ref, err := docker.ParseReference("//" + imageName)
	if err != nil {
		return nil, err
	}
	authCtx := &types.SystemContext{DockerAuthConfig: &types.DockerAuthConfig{}}
	ctx, cancel := context.WithTimeout(c.opts.Ctx, c.opts.ReqTimeout)
	defer cancel()
	return docker.GetRepositoryTags(ctx, authCtx, ref)
}

func (c *Client) Sync(src, dst string) error {
	log.Debugf("syncing %s to %s ...", src, dst)
	srcRef, err := docker.ParseReference("//" + src)
	if err != nil {
		return err
	}
	destRef, err := docker.ParseReference("//" + dst)
	if err != nil {
		return err
	}

	policyCtx, err := signature.NewPolicyContext(
		&signature.Policy{
			Default: []signature.PolicyRequirement{signature.NewPRInsecureAcceptAnything()},
		},
	)

	if err != nil {
		return err
	}
	defer func() { _ = policyCtx.Destroy() }()

	srcCtx := &types.SystemContext{
		DockerAuthConfig:         &types.DockerAuthConfig{},
		OCIInsecureSkipTLSVerify: true,
	}
	dstCtx := &types.SystemContext{
		DockerAuthConfig: &types.DockerAuthConfig{
			Username: c.opts.Username,
			Password: c.opts.Password,
		},
		OCIInsecureSkipTLSVerify: true,
	}
	ctx, cancel := context.WithTimeout(c.opts.Ctx, c.opts.PushTimeout)
	defer cancel()

	_, err = copy.Image(ctx, policyCtx, destRef, srcRef, &copy.Options{
		SourceCtx:          srcCtx,
		DestinationCtx:     dstCtx,
		ImageListSelection: copy.CopyAllImages,
	})
	if err != nil {
		log.Debugf("sync %s error: %s", err)
		return err
	}
	log.Debugf("sync %s done", src)
	return nil
}

func (c *Client) allImages(imagesUri string) ([]string, error) {
	var baseNames []string
	res, err := c.R().Get(imagesUri)
	if err != nil {
		return nil, err
	}
	switch res.StatusCode() {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		// do nothing
	default:
		return nil, errors.Errorf("%s status: %d", imagesUri, res.StatusCode())
	}
	err = jsoniter.UnmarshalFromString(jsoniter.Get(res.Body(), "child").ToString(), &baseNames)
	if err != nil {
		return nil, err
	}
	return baseNames, nil
}
