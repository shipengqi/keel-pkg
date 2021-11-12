package cmd

import (
	"fmt"

	"github.com/shipengqi/keel-pkg/lib/deps"
	"github.com/shipengqi/keel-pkg/lib/utils/tmplutil"
)

const (
	ImageNameCoredns       = "coredns"
	ImageNameFlannel       = "flannel"
	ImageNameBusybox       = "busybox"
	ImageNameMetricsServer = "metrics-server"
)

func normalizeImgName(name, tag, arch string) string {
	var normalized string
	switch name {
	case ImageNameFlannel:
		normalized = fmt.Sprintf("%s:%s-%s", name, tag, arch)
	case ImageNameMetricsServer, ImageNameBusybox:
		normalized = fmt.Sprintf("%s:%s", name, tag)
	default:
		normalized = fmt.Sprintf("%s-%s:%s", name, arch, tag)
	}

	return normalized
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
