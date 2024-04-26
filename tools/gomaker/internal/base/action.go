package base

type GenFunc func(action string, src string, params string) error

type Action struct {
	Name  string
	Help  string
	Param string
	Gen   GenFunc
}

func NewAction(f GenFunc, name, param, help string) *Action {
	return &Action{name, help, param, f}
}
