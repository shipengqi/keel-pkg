package image

import (
	"context"
	"fmt"
	"hash/crc32"
	"time"

	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/manifest"
	"github.com/containers/image/v5/types"
	specv1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

const (
	DefaultCtxTimeout = 5 * time.Minute
	DefaultLimit      = 20
	K8SRegistryUri    = "k8s.gcr.io"
)

type Image struct {
	Registry string
	Name     string
	Tag      string
	layers   []string
}

func (i *Image) String() string {
	return fmt.Sprintf("%s/%s:%s", i.Registry, i.Name, i.Tag)
}

func (i *Image) Key() string {
	return fmt.Sprintf("%s:%s", i.Name, i.Tag)
}

func (i *Image) Manifest() (sum uint32, err error) {
	srcRef, err := docker.ParseReference("//" + i.Key())
	if err != nil {
		return
	}
	authCtx := &types.SystemContext{DockerAuthConfig: &types.DockerAuthConfig{}}
	imgCtx, cancel := context.WithTimeout(context.Background(), DefaultCtxTimeout)
	defer cancel()
	src, err := srcRef.NewImageSource(imgCtx, authCtx)
	if err != nil {
		return
	}
	defer src.Close()
	manifestCtx, manifestCancel := context.WithTimeout(context.Background(), DefaultCtxTimeout)
	defer manifestCancel()
	mbs, _, err := src.GetManifest(manifestCtx, nil)
	if err != nil {
		return 0, err
	}
	mType := manifest.GuessMIMEType(mbs)
	if len(mType) == 0 {
		return 0, errors.Errorf("parse image [%s] manifest type", i.Key())
	}
	if mType != manifest.DockerV2ListMediaType && mType != specv1.MediaTypeImageIndex {
		_, err = manifest.FromBlob(mbs, mType)
		if err != nil {
			return 0, err
		}
	}
	return crc32.ChecksumIEEE(mbs), nil
}
