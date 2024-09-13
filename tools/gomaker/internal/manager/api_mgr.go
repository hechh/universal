package manager

import (
	"flag"
	"fmt"
	"text/template"
	"universal/tools/gomaker/domain"
)

var (
	apis = make(map[string]*ApiInfo)
)

type ApiInfo struct {
	fs   []domain.GenFunc
	help string
}

func Help() {
	fmt.Fprintf(flag.CommandLine.Output(), "action使用说明: \n")
	for action, info := range apis {
		fmt.Fprint(flag.CommandLine.Output(), fmt.Sprintf("\t -action=%s\t#%s\n", action, info.help))
	}
}

func Register(action, help string, fs ...domain.GenFunc) {
	apis[action] = &ApiInfo{help: help, fs: fs}
}

func IsAction(action string) bool {
	_, ok := apis[action]
	return ok
}

func Generator(action, dst, param string, tpls *template.Template) error {
	for _, f := range apis[action].fs {
		if err := f(dst, param, tpls); err != nil {
			return err
		}
	}
	return nil
}
