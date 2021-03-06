package client

import (
	"context"
	"fmt"
	"hash/crc32"
	"net/http"
	"time"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/manifest"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/types"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	specv1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"

	"github.com/shipengqi/keel-pkg/lib/deps"
	"github.com/shipengqi/keel-pkg/lib/log"
)

const (
	DefaultGcrRepo       = "k8s.gcr.io"
	DefaultRetryTimes    = 5
	DefaultRetryInterval = time.Second * 5
	DefaultTimeout       = time.Second * 20
	DefaultPushTimeout   = time.Minute * 15
	// LegacyConfigMediaType should be replaced by OCI image spec.
	// More detail: docker/distribution#1622
	LegacyConfigMediaType = "application/octet-stream"
)

const (
	_defaultGcrImagesListAPI = "https://k8s.gcr.io/v2/tags/list"
)

var (
	ErrUnsupportedMediaType = errors.New("Unsupported MediaType")
)

type ManifestConfig struct {
	MediaType string `json:"mediaType"`
}

type Manifest struct {
	MediaType     string          `json:"mediaType"`
	SchemaVersion int             `json:"schemaVersion"`
	Signatures    interface{}     `json:"signatures"`
	Config        *ManifestConfig `json:"config""`
}

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
	ImageSet      *deps.ImageSet
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
		for k := range baseNames {
			baseNames[k] = fmt.Sprintf("%s/%s", c.opts.AdditionalNS[i], baseNames[k])
		}
		allBaseNames = append(allBaseNames, baseNames...)
	}

	return allBaseNames, nil
}

func (c *Client) AllTagsWithRepo(repo, baseName string) ([]string, error) {
	return c.allTagsWithRepo(repo, baseName)
}

func (c *Client) AllTags(baseName string) ([]string, error) {
	return c.allTagsWithRepo(c.opts.Repo, baseName)
}

func (c *Client) Sync(src, dst string) error {
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
	return nil
}

func (c *Client) ManifestCheckSum(imageName string) (uint32, error) {
	ref, err := docker.ParseReference("//" + imageName)
	if err != nil {
		return 0, err
	}
	authCtx := &types.SystemContext{DockerAuthConfig: &types.DockerAuthConfig{}}
	ctx, cancel := context.WithTimeout(c.opts.Ctx, c.opts.ReqTimeout)
	defer cancel()
	src, err := ref.NewImageSource(ctx, authCtx)
	if err != nil {
		return 0, err
	}
	reqCtx, reqCancel := context.WithTimeout(context.Background(), c.opts.ReqTimeout)
	defer reqCancel()
	mbs, _, err := src.GetManifest(reqCtx, nil)
	if err != nil {
		return 0, err
	}
	if !mediaTypeSupported(mbs) {
		return 0, ErrUnsupportedMediaType
	}
	mType := manifest.GuessMIMEType(mbs)
	if mType == "" {
		return 0, errors.Errorf("parse image [%s] manifest type", imageName)
	}

	if mType != manifest.DockerV2ListMediaType && mType != specv1.MediaTypeImageIndex {
		_, err = manifest.FromBlob(mbs, mType)
		if err != nil {
			return 0, err
		}
	}
	return crc32.ChecksumIEEE(mbs), nil
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

func (c *Client) allTagsWithRepo(repo, baseName string) ([]string, error) {
	imageName := fmt.Sprintf("%s/%s", repo, baseName)
	ref, err := docker.ParseReference("//" + imageName)
	if err != nil {
		return nil, err
	}
	authCtx := &types.SystemContext{DockerAuthConfig: &types.DockerAuthConfig{}}
	ctx, cancel := context.WithTimeout(c.opts.Ctx, c.opts.ReqTimeout)
	defer cancel()
	return docker.GetRepositoryTags(ctx, authCtx, ref)
}

func mediaTypeSupported(mbs []byte) bool {
	meta := &Manifest{}
	if err := jsoniter.Unmarshal(mbs, &meta); err != nil {
		return false
	}
	if err := manifest.SupportedSchema2MediaType(meta.MediaType); err != nil {
		return false
	}
	if meta.Config != nil && meta.Config.MediaType != "" {
		if err := manifest.SupportedSchema2MediaType(meta.Config.MediaType); err != nil {
			return false
		}
	}
	return true
}
