package manager

import (
	"flag"
	"fmt"
	"universal/framework/fbasic"
)

type GenFunc func(action string, src string, params string) error

type GenInfo struct {
	action string
	help   string
	gen    GenFunc
}

var (
	genMgr = make(map[string]*GenInfo)
)

func Register(action string, f GenFunc, helps ...string) {
	if _, ok := genMgr[action]; ok {
		panic(fmt.Sprintf("%s has already registered", action))
	}
	genMgr[action] = &GenInfo{
		action: action,
		gen:    f,
		help: func() string {
			if len(helps) <= 0 {
				return ""
			}
			return helps[0]
		}(),
	}
}

func Gen(action string, src string, params string) error {
	val, ok := genMgr[action]
	if !ok {
		return fbasic.NewUError(2, -1, fmt.Sprintf("%s is not suppported", action))
	}
	return val.gen(action, src, params)
}

func HelpAction() {
	fmt.Fprintf(flag.CommandLine.Output(), "-a 参数说明：\n")
	for _, item := range genMgr {
		fmt.Fprintf(flag.CommandLine.Output(), "\t %s: %s \n", item.action, item.help)
	}
}
