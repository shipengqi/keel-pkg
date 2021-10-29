package action

import "fmt"

type Images []*Image

type Image struct {
	Repo string // for patch images
	Ns   string // for patch images
	Name string // basename
	Tag  string
}

func (i *Image) String() string {
	var uri string
	if len(i.Repo) == 0 {
		i.Repo = "k8s.gcr.io"
	}
	uri = i.Repo
	if len(i.Ns) > 0 {
		uri = uri + "/" + i.Ns
	}

	return fmt.Sprintf("%s/%s:%s", uri, i.Name, i.Tag)
}

func (i *Image) Key() string {
	return fmt.Sprintf("%s:%s", i.Name, i.Tag)
}

func (is Images) Len() int           { return len(is) }
func (is Images) Less(i, j int) bool { return is[i].String() < is[j].String() }
func (is Images) Swap(i, j int)      { is[i], is[j] = is[j], is[i] }
