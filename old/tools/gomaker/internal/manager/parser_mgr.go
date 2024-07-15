package manager

import (
	"flag"
	"fmt"
	"universal/tools/gomaker/domain"
)

var (
	genMgr = make(map[string]domain.IMaker)
)

func Register(name string, act domain.IMaker) {
	if _, ok := genMgr[name]; ok {
		panic(fmt.Sprintf("%s has already registered", name))
	}
	genMgr[name] = act
}

func GetMaker(name string) domain.IMaker {
	return genMgr[name]
}

func Help() {
	fmt.Fprintf(flag.CommandLine.Output(), "action使用说明: \n")
	for name, item := range genMgr {
		fmt.Fprint(flag.CommandLine.Output(), item.GetHelp(name))
	}
}
