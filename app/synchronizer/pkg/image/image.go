package image

import (
	"fmt"

	"github.com/containers/image/v5/manifest"
)

const (
	K8SRegistryUri = "k8s.gcr.io"
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

func (i *Image) Manifest() string {
	return fmt.Sprintf("%s:%s", i.Name, i.Tag)
}
