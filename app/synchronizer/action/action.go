package action

import (
	"context"
	"strings"

	"github.com/shipengqi/keel-pkg/lib/log"
)

type CloseFunc func() error

type Interface interface {
	Name() string
	PreRun() error
	Run() error
	PostRun() error
	Ctx() context.Context
	Close() CloseFunc
}

type action struct {
	name  string
	ctx   context.Context
	close CloseFunc
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

func (a *action) Close() CloseFunc {
	return nil
}

func Execute(a Interface) error {
	defer func() {
		if a.Close() != nil {
			a.Close()
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
