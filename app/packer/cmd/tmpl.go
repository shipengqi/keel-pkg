package cmd

import (
	"fmt"

	"github.com/shipengqi/keel-pkg/lib/deps"
	"github.com/shipengqi/keel-pkg/lib/utils/tmplutil"
)

const (
	CorednsImageName = "coredns"
	FlannelImageName = "flannel"
)

func normalizeImgName(name, tag, arch string) string {
	if name == FlannelImageName {
		return fmt.Sprintf("%s:%s-%s", name, tag, arch)
	}
	return fmt.Sprintf("%s-%s:%s", name, arch, tag)
}

func uriTmplList(vs *deps.Versions) ([]string, error) {
	var list []string
	for i := range vs.Components {
		c := vs.Components[i]
		uri, err := tmplutil.ReplaceString(
			c.Uri, c.Name, &deps.TmplData{Tag: c.Tag, Arch: vs.Arch})
		if err != nil {
			return nil, err
		}
		list = append(list, uri)
	}
	return list, nil
}
