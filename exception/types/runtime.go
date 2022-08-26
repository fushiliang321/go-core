package types

type Runtime struct {
	error
	Trace []string
	Msg   string
	Mark  string
}

func (e Runtime) Error() string {
	return e.Msg
}
