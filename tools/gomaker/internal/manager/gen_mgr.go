package manager

import (
	"flag"
	"fmt"
	"hego/framework/uerror"
	"hego/tools/gomaker/domain"
)

var (
	gens = make(map[string]*GenInfo)
)

type GenInfo struct {
	action string
	help   string
	gens   []domain.GenFunc
}

func (d *GenInfo) Run(dst string, extra ...string) error {
	for _, f := range d.gens {
		if err := f(dst, extra...); err != nil {
			return err
		}
	}
	return nil
}

func Help() {
	fmt.Fprintf(flag.CommandLine.Output(), "action使用说明: \n")
	for action, info := range gens {
		fmt.Fprint(flag.CommandLine.Output(), fmt.Sprintf("\t -action=%s\t#%s\n", action, info.help))
	}
}

func Register(action, help string, fs ...domain.GenFunc) {
	gens[action] = &GenInfo{action: action, help: help, gens: fs}
}

func Generator(action string, dst string, extra ...string) error {
	if val, ok := gens[action]; ok {
		return val.Run(dst, extra...)
	}
	return uerror.NewUError(1, -1, "生成模式(%s)不支持", action)
}
