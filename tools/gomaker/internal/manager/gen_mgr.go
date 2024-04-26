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
		var str string
		if len(item.Param) > 0 {
			str = fmt.Sprintf("    -action=%s -param=%s", item.Name, item.Param)
		} else {
			str = fmt.Sprintf("    -action=%s", item.Name)
		}
		fmt.Fprintf(flag.CommandLine.Output(), "%-80s #%s\n", str, item.Help)
	}
}
