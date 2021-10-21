package cmd

import (
	"fmt"

	"github.com/shipengqi/keel-pkg/lib/deps"
)

const (
	CorednsImageName = "coredns"
	ContainerdName   = "containerd"
	RuncName         = "runc"
	CrictlName       = "crictl"
	KubectlName      = "kubectl"
	KubeletName      = "kubelet"
)

var (
	containerdUriTmpl = "https://github.com/containerd/containerd/releases/download/v%s/containerd-%s-linux-%s.tar.gz"
	runcUriTmpl       = "https://github.com/opencontainers/runc/releases/download/v%s/runc.%s"
	crictlUriTmpl     = "https://github.com/kubernetes-sigs/cri-tools/releases/download/v%s/crictl-v%s-linux-%s.tar.gz"
	kubectlUriTmpl    = "https://dl.k8s.io/v%s/bin/linux/%s/kubectl"
	kubeletUriTmpl    = "https://dl.k8s.io/v%s/bin/linux/%s/kubelet"
)

func normalizeImgName(name, tag, arch string) string {
	// image name needs to add arch string, except for the coredns, metrics image
	if name == CorednsImageName {
		return fmt.Sprintf("%s:%s", name, tag)
	}
	return fmt.Sprintf("%s-%s:%s", name, arch, tag)
}

func uriTmplList(vs *deps.Versions) []string {
	var list []string
	for i := range vs.Components {
		c := vs.Components[i]
		switch c.Name {
		case ContainerdName:
			list = append(list, fmt.Sprintf(containerdUriTmpl, c.Tag, c.Tag, vs.Arch))
		case RuncName:
			list = append(list, fmt.Sprintf(runcUriTmpl, c.Tag, vs.Arch))
		case CrictlName:
			list = append(list, fmt.Sprintf(crictlUriTmpl, c.Tag, c.Tag, vs.Arch))
		case KubectlName:
			list = append(list, fmt.Sprintf(kubectlUriTmpl, c.Tag, vs.Arch))
		case KubeletName:
			list = append(list, fmt.Sprintf(kubeletUriTmpl, c.Tag, vs.Arch))
		}
	}
	return list
}
