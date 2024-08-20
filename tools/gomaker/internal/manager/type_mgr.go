package manager

import (
	"flag"
	"fmt"
	"sort"
	"strings"
	"text/template"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/typespec"
)

var (
	types   = make(map[string]*typespec.Type)
	enums   = make(map[string]*typespec.Enum)
	structs = make(map[string]*typespec.Struct)
	alias   = make(map[string]*typespec.Alias)
	apis    = make(map[string]*ApiInfo)
)

type ApiInfo struct {
	help      string
	generator domain.GenFunc
}

func Help() {
	fmt.Fprintf(flag.CommandLine.Output(), "action使用说明: \n")
	for action, info := range apis {
		fmt.Fprint(flag.CommandLine.Output(), fmt.Sprintf("\t -action=%s\t#%s\n", action, info.help))
	}
}

func Register(action, help string, info domain.GenFunc) {
	apis[action] = &ApiInfo{help: help, generator: info}
}

func IsAction(action string) bool {
	_, ok := apis[action]
	return ok
}

func Generator(action, dst, param string, tpls *template.Template) error {
	return apis[action].generator(dst, param, tpls)
}

func Print() {
	for key, val := range alias {
		fmt.Println(val.Format())
		if vv, ok := enums[key]; ok {
			fmt.Println(vv.Format())
		}
	}
	for _, val := range structs {
		fmt.Println(val.Format())
	}
}

func GetOrAddType(tt *typespec.Type) *typespec.Type {
	key := fmt.Sprintf("%s.%s", tt.Selector, tt.Name)
	if val, ok := types[key]; !ok {
		types[key] = tt
	} else {
		if len(tt.Doc) > 0 {
			val.Doc = tt.Doc
		}
		if tt.Kind > 0 {
			val.Kind = tt.Kind
		}
	}
	return types[key]
}

func AddEnum(vv *typespec.Enum) {
	if vv != nil {
		enums[fmt.Sprintf("%s.%s", vv.Type.Selector, vv.Type.Name)] = vv
	}
}

func AddStruct(vv *typespec.Struct) {
	if vv != nil {
		structs[fmt.Sprintf("%s.%s", vv.Type.Selector, vv.Type.Name)] = vv
	}
}

func AddAlias(vv *typespec.Alias) {
	if vv != nil {
		alias[fmt.Sprintf("%s.%s", vv.Type.Selector, vv.Type.Name)] = vv
	}
}

func GetStruct() (rets []*typespec.Struct) {
	for _, val := range structs {
		rets = append(rets, val)
	}
	sort.Slice(rets, func(i, j int) bool {
		return strings.Compare(rets[i].Type.Name, rets[j].Type.Name) < 0
	})
	return
}

func GetEnum() (rets []*typespec.Enum) {
	for _, val := range enums {
		rets = append(rets, val)
	}
	sort.Slice(rets, func(i, j int) bool {
		return strings.Compare(rets[i].Type.Name, rets[j].Type.Name) < 0
	})
	return
}

func GetAlias() (rets []*typespec.Alias) {
	for _, val := range alias {
		rets = append(rets, val)
	}
	sort.Slice(rets, func(i, j int) bool {
		return strings.Compare(rets[i].Type.Name, rets[j].Type.Name) < 0
	})
	return
}
