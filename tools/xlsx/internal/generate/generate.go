package generate

import (
	"bytes"
	"fmt"
	"hego/Library/file"
	"hego/tools/xlsx/domain"
	"hego/tools/xlsx/internal/base"
	"strings"
)

func Generate(codePath string, cfg *base.Config, buf *bytes.Buffer) error {
	dataName := fmt.Sprintf("%sData", cfg.Name)
	cfgName := fmt.Sprintf("%s.%sConfig", domain.PkgName, cfg.Name)

	// 包
	buf.WriteString(fmt.Sprintf(pack, dataName))

	// 结构
	var members, makes, inits, sets []string
	for _, item := range cfg.IndexList {
		members = append(members, fmt.Sprintf("_%s %s", item.Name, item.GetType(cfgName)))
		if len(item.List) > 1 {
			buf.WriteString(fmt.Sprintf("\ntype %s struct {%s}\n", item.Name, item.Arg(domain.PkgName, ";")))
		}
	}
	buf.WriteString(fmt.Sprintf("\ntype %s struct {%s}\n", dataName, strings.Join(members, ";")))

	// 函数
	for _, item := range cfg.IndexList {
		arg := item.Arg(domain.PkgName, ",")
		val := item.IndexValue(item.Value("", ","))
		ref := item.IndexValue(item.Value("item", ","))
		switch item.Type.ValueOf {
		case domain.ValueOfList:
			inits = append(inits, fmt.Sprintf("_%s: ary,", item.Name))
			buf.WriteString(fmt.Sprintf(sget, cfgName, dataName, item.Name))
			buf.WriteString(fmt.Sprintf(lget, cfgName, dataName, cfgName, item.Name, item.Name))
		case domain.ValueOfMap:
			sets = append(sets, fmt.Sprintf("_%s[%s] = item", item.Name, ref))
			makes = append(makes, fmt.Sprintf("_%s:= make(%s)", item.Name, item.GetType(cfgName)))
			inits = append(inits, fmt.Sprintf("_%s: _%s,", item.Name, item.Name))
			buf.WriteString(fmt.Sprintf(mget, item.Name, arg, cfgName, dataName, item.Name, val))
		case domain.ValueOfGroup:
			sets = append(sets, fmt.Sprintf("_%s[%s] = append(_%s[%s], item)", item.Name, ref, item.Name, ref))
			makes = append(makes, fmt.Sprintf("_%s:= make(%s)", item.Name, item.GetType(cfgName)))
			inits = append(inits, fmt.Sprintf("_%s: _%s,", item.Name, item.Name))
			buf.WriteString(fmt.Sprintf(gget, item.Name, arg, cfgName, dataName, item.Name, val, cfgName))
		}
	}
	buf.WriteString(fmt.Sprintf(parse, cfgName, strings.Join(makes, "\n"), strings.Join(sets, "\n"), dataName, strings.Join(inits, "\n")))

	return file.SaveGo(codePath, fmt.Sprintf("%s.gen.go", cfg.FileName), buf.Bytes())
}

const (
	parse = `
func Parse(buf []byte) {
	ary := []*%s{}
	if err := json.Unmarshal(buf, &ary); err != nil {
		panic(err)
	}
	%s
	for _, item := range ary {
		%s
	}
	obj.Store(&%s{
		%s
	})
}
	`
	pack = `
package %s

import (
	"encoding/json"
	"hego/common/cfg"
	"sync/atomic"
)

var obj = atomic.Value{}
	`
	sget = `
func SGet() *%s {
	if d, ok := obj.Load().(*%s); ok {
		return d._%s[0]
	}
	return nil
}	
	`
	lget = `
func LGet() (rets []*%s) {
	if d, ok := obj.Load().(*%s); ok {
		rets = make([]*%s, len(d._%s))	
		copy(rets, d._%s)
	}
	return
}	
	`
	mget = `
func MGet%s(%s) *%s {
	if d, ok := obj.Load().(*%s); ok {
		val, ok := d._%s[%s]
		if ok {
			return val	
		}
	}
	return nil
}	
	`
	gget = `
func GGet%s(%s) (rets []*%s) {
	if d, ok := obj.Load().(*%s); ok {
		vals, ok := d._%s[%s]
		if ok {
			rets = make([]*%s, len(vals))	
			copy(rets, vals)
		}	
	}
	return
}
	`
)
