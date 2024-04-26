package manager

import (
	"flag"
	"fmt"
	"universal/framework/fbasic"
	"universal/tools/gomaker/internal/base"
)

var (
	genMgr = make(map[string]*base.Action)
)

func Register(act *base.Action) {
	if _, ok := genMgr[act.Name]; ok {
		panic(fmt.Sprintf("%s has already registered", act.Name))
	}
	genMgr[act.Name] = act
}

func Gen(action string, src string, params string) error {
	val, ok := genMgr[action]
	if !ok {
		return fbasic.NewUError(2, -1, fmt.Sprintf("%s is not suppported", action))
	}
	return val.Gen(action, src, params)
}

func Help() {
	fmt.Fprintf(flag.CommandLine.Output(), "action使用说明: \n")
	for _, item := range genMgr {
		if len(item.Param) > 0 {
			fmt.Fprintf(flag.CommandLine.Output(), "\t-action=%s -param=%s  //%s\n", item.Name, item.Param, item.Help)
		} else {
			fmt.Fprintf(flag.CommandLine.Output(), "\t-action=%s  //%s\n", item.Name, item.Help)
		}
	}
}
