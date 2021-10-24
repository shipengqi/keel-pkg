package cmd

import (
	"fmt"

	"github.com/shipengqi/keel-pkg/lib/deps"
	"github.com/shipengqi/keel-pkg/lib/utils/tmplutil"
)

const (
	CorednsImageName = "coredns"
)

func normalizeImgName(name, tag, arch string) string {
	// image name needs to add arch string, except for the coredns, metrics image
	if name == CorednsImageName {
		return fmt.Sprintf("%s:%s", name, tag)
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
