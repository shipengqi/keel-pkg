package action

type suma struct {
	*action
}

const (
	NameSum = "sum"
)

func NewSumAction() Interface {
	a := &suma{
		action: &action{
			name: NameSum,
		},
	}
	return a
}
