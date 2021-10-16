package action

import gcrc "github.com/shipengqi/keel-pkg/app/synchronizer/pkg/registry/gcr/client"

type checka struct {
	*action

	gcr  *gcrc.Client
}

const (
	NameCheck = "check"
)

func NewCheckAction(opts *gcrc.Options) Interface {
	a := &checka{
		action: &action{
			name: NameCheck,
		},
		gcr:  gcrc.New(opts),
	}
	return a
}
