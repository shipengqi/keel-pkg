package action

import "fmt"

type Images []*Image

type Image struct {
	Name     string // basename
	Tag      string
	Success  bool
	CacheHit bool
	Err      error
}

func (i *Image) String() string {
	return fmt.Sprintf("k8s.gcr.io/%s:%s", i.Name, i.Tag)
}

func (i *Image) Key() string {
	return fmt.Sprintf("%s:%s", i.Name, i.Tag)
}

func (is Images) Len() int           { return len(is) }
func (is Images) Less(i, j int) bool { return is[i].String() < is[j].String() }
func (is Images) Swap(i, j int)      { is[i], is[j] = is[j], is[i] }