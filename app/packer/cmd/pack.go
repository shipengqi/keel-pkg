package cmd

import (
	"fmt"
	"github.com/shipengqi/keel-pkg/lib/deps"
	"github.com/shipengqi/keel-pkg/lib/log"
	"github.com/shipengqi/keel-pkg/lib/utils/cliutil"
	"path/filepath"
)

const (
	DefaultRegistry     = "registry.cn-hangzhou.aliyuncs.com"
	DefaultRegistryNs   = "keel"
	DefaultImagesOutput = "/var/run/keel/pack/images"
)

func pack(o *packOptions, set *deps.Versions) error {
	list := uriTmplList(set)
	for i := range list {
		err := download(list[i])
		if err != nil {
			return err
		}
	}
	err := login(o.RegistryUser, o.RegistryPass)
	if err != nil {
		return err
	}
	requiredImages := set.Images
	for i := range requiredImages {
		normalized := normalizeImgName(requiredImages[i].Name, requiredImages[i].Tag, set.Arch)
		err = pull(o.ImagesOutput, normalized)
		if err != nil {
			log.Warnf(err.Error())
			continue
		}
	}
	return nil
}

func login(user, pass string) error {
	log.Debugf("login to [%s] ...", DefaultRegistry)
	_, stderr, _, err := cliutil.Exec("docker",
		[]string{"login", "--username", user, "--password", pass, DefaultRegistry})
	if err != nil {
		log.Debugf("login to [%s]: %v", DefaultRegistry, stderr)
		return err
	}
	log.Debugf("login [%s] ok!", DefaultRegistry)
	return nil
}

func pull(output, imgName string) error {
	imgFullName := fmt.Sprintf("%s/%s/%s", DefaultRegistry, DefaultRegistryNs, imgName)
	log.Debugf("pulling [%s] ...", imgFullName)
	_, stderr, _, err := cliutil.Exec("docker", []string{"pull", imgFullName})
	if err != nil {
		log.Debugf("pull [%s]: %v", imgFullName, stderr)
		return err
	}
	log.Debugf("pull [%s] done!", imgFullName)
	tarName := filepath.Join(output, fmt.Sprintf("%s.tar", imgName))
	log.Debugf("saving [%s] to [%s] ...", imgFullName, tarName)

	_, stderr, _, err = cliutil.Exec("docker", []string{"save", "-o", tarName, imgFullName})
	if err != nil {
		log.Debugf("save [%s]: %v", tarName, stderr)
		return err
	}
	log.Debugf("save [%s] done!", tarName)
	return nil
}

func download(uri string) error {
	log.Debugf("wget [%s] ...", uri)
	_, stderr, _, err := cliutil.Exec("wget", []string{uri})
	if err != nil {
		log.Debugf("wget [%s]: %v", uri, stderr)
		return err
	}
	return nil
}
