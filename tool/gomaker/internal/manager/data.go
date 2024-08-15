package manager

import (
	"flag"
	"fmt"
	"text/template"
	"universal/tool/gomaker/domain"
	"universal/tool/gomaker/internal/typespec"
)

var (
	types   = make(map[string]*typespec.Type)
	enums   = make(map[string]*typespec.Enum)
	structs = make(map[string]*typespec.Struct)
	alias   = make(map[string]*typespec.Alias)
	apis    = make(map[string]*ApiInfo)
)

type ApiInfo struct {
	help       string
	generators []domain.GenFunc
}

func Help() {
	fmt.Fprintf(flag.CommandLine.Output(), "action使用说明: \n")
	for action, info := range apis {
		fmt.Fprint(flag.CommandLine.Output(), fmt.Sprintf("\t%s #%s\n", action, info.help))
	}
}

func Register(action, help string, infos ...domain.GenFunc) {
	if _, ok := apis[action]; !ok {
		apis[action] = &ApiInfo{help: help}
	}
	apis[action].generators = append(apis[action].generators, infos...)
}

func IsAction(action string) bool {
	_, ok := apis[action]
	return ok
}

func Generator(action, dst, param string, tpls map[string]*template.Template) error {
	api := apis[action]
	for _, ff := range api.generators {
		if err := ff(dst, param, tpls); err != nil {
			return err
		}
	}
	return nil
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

func GetStruct() map[string]*typespec.Struct {
	return structs
}

func GetEnum() map[string]*typespec.Enum {
	return enums
}

func GetAlias() map[string]*typespec.Alias {
	return alias
}
