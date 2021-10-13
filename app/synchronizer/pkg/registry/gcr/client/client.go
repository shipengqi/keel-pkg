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
	repo = "k8s.gcr.io"
)

type Options struct {
	Username      string
	Password      string
	Token         string
	Repo          string
	Retry         int
	RetryInterval time.Duration
	Timeout       time.Duration
	PushTimeout   time.Duration
	Ctx           context.Context
}

type Client struct {
	*resty.Client

	opts *Options
}

func New(options *Options) *Client {
	r := resty.New()
	r.SetTimeout(options.Timeout)
	r.SetRetryCount(options.Retry)
	r.SetRetryWaitTime(options.RetryInterval)
	if options.Ctx == nil {
		options.Ctx = context.Background()
	}
	if options.Repo == "" {
		options.Repo = repo
	}
	return &Client{
		Client: r,
		opts:   options,
	}
}

func (c *Client) AllImages() ([]string, error) {
	var baseNames []string
	res, err := c.R().Get("https://k8s.gcr.io/v2/tags/list")
	if err != nil {
		return nil, err
	}
	switch res.StatusCode() {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		// do nothing
	default:
		return nil, errors.Errorf("/tags/list status: %d", res.StatusCode())
	}
	err = jsoniter.UnmarshalFromString(jsoniter.Get(res.Body(), "child").ToString(), &baseNames)
	if err != nil {
		return nil, err
	}
	return baseNames, nil
}

func (c *Client) AllTags(baseName string) ([]string, error) {
	imageName := fmt.Sprintf("%s/%s", c.opts.Repo, baseName)
	ref, err := docker.ParseReference("//" + imageName)
	if err != nil {
		return nil, err
	}
	authCtx := &types.SystemContext{DockerAuthConfig: &types.DockerAuthConfig{}}
	ctx, cancel := context.WithTimeout(c.opts.Ctx, c.opts.Timeout)
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
