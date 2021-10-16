package action

import (
	"context"
	"strings"

	"github.com/shipengqi/keel-pkg/lib/log"
)

const (
	K8SRegistryUri = "k8s.gcr.io"
)

type Interface interface {
	Name() string
	PreRun() error
	Run() error
	PostRun() error
	Ctx() context.Context
	Cancel() context.CancelFunc
}

type action struct {
	name   string
	ctx    context.Context
	cancel context.CancelFunc
}

func (a *action) Name() string {
	return a.name
}

func (a *action) PreRun() error {
	log.Debugf("***** [%s] PreRun *****", strings.ToUpper(a.name))
	return nil
}

func (a *action) Run() error {
	log.Debugf("***** [%s] Run *****", strings.ToUpper(a.name))
	return nil
}

func (a *action) PostRun() error {
	log.Debugf("***** [%s] PostRun *****", strings.ToUpper(a.name))
	return nil
}

func (a *action) Ctx() context.Context {
	return a.ctx
}

func (a *action) Cancel() context.CancelFunc {
	return a.cancel
}

func Execute(a Interface) error {
	defer func() {
		if a.Cancel() != nil {
			a.Cancel()
		}
	}()
	err := a.PreRun()
	if err != nil {
		return err
	}
	err = a.Run()
	if err != nil {
		return err
	}
	return a.PostRun()
}