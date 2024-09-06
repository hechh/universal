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
	f    domain.GenFunc
	help string
}

func Help() {
	fmt.Fprintf(flag.CommandLine.Output(), "action使用说明: \n")
	for action, info := range apis {
		fmt.Fprint(flag.CommandLine.Output(), fmt.Sprintf("\t -action=%s\t#%s\n", action, info.help))
	}
}

func Register(action string, info domain.GenFunc, help string) {
	apis[action] = &ApiInfo{help: help, f: info}
}

func IsAction(action string) bool {
	_, ok := apis[action]
	return ok
}

func Generator(action, dst, param string, tpls *template.Template) error {
	return apis[action].f(dst, param, tpls)
}
