package generate

import (
	"fmt"
	"hego/Library/util"
	"hego/tools/xlsx/domain"
	"hego/tools/xlsx/internal/base"
	"strings"
)

type FunInfo struct {
	Pos     int
	ValueOf uint32
	Name    string
	Index   string
	Args    []string
	Values  []string
	Refs    []string
}

type StInfo struct {
	Name    string
	Members []string
	Makes   []string
}

type StInfoMgr struct {
	PkgName  string
	DataName string
	CType    string
	Data     map[string]*StInfo
	List     []*StInfo
	Funs     []*FunInfo
}

func NewStInfoMgr(pkg string, cfg *base.Config) *StInfoMgr {
	return &StInfoMgr{
		PkgName:  pkg,
		DataName: fmt.Sprintf("%sData", cfg.Name),
		CType:    fmt.Sprintf("%s.%sConfig", pkg, cfg.Name),
		Data:     map[string]*StInfo{},
	}
}

func (d *StInfo) Push(key, val string) {
	d.Members = append(d.Members, fmt.Sprintf("%s %s", key, val))
}

func (d *StInfoMgr) Push(Data string, key, val string) {
	if _, ok := d.Data[Data]; !ok {
		d.Data[Data] = &StInfo{Name: Data}
		d.List = append(d.List, d.Data[Data])
	}
	d.Data[Data].Push(key, val)
}

func (d *StInfo) Add(key, val string) {
	d.Push(key, val)
	d.Makes = append(d.Makes, fmt.Sprintf("%s: make(%s),", key, val))
}

func (d *StInfoMgr) Add(data string, key, val string) {
	if _, ok := d.Data[data]; !ok {
		d.Data[data] = &StInfo{Name: data}
		d.List = append(d.List, d.Data[data])
	}
	d.Data[data].Add(key, val)
}

func (d *StInfoMgr) AddFun(f *FunInfo) {
	d.Funs = append(d.Funs, f)
}

func (d *StInfoMgr) Package() string {
	return fmt.Sprintf(pack, d.DataName)
}

func (d *StInfoMgr) Define() string {
	strs := []string{}
	for _, item := range d.List {
		strs = append(strs, fmt.Sprintf("\ntype %s struct {%s}\n", item.Name, strings.Join(item.Members, ";")))
	}
	return strings.Join(strs, "\n")
}

func (d *StInfoMgr) Func() string {
	strs := []string{}
	strs = append(strs, fmt.Sprintf(sget, d.CType, d.DataName))
	strs = append(strs, fmt.Sprintf(lget, d.CType, d.DataName, d.CType))
	for _, item := range d.Funs {
		arg := strings.Join(item.Args, ", ")
		key := util.Ifelse(len(item.Values) == 1, item.Values[0], fmt.Sprintf("%s{%s}", item.Index, strings.Join(item.Values, ", ")))
		switch item.ValueOf {
		case domain.VALUE_OF_MAP:
			strs = append(strs, fmt.Sprintf(mget, item.Pos, arg, d.CType, d.DataName, item.Pos, key))
		case domain.VALUE_OF_GROUP:
			strs = append(strs, fmt.Sprintf(gget, item.Pos, arg, d.CType, d.DataName, item.Pos, key, d.CType))
		}
	}
	return strings.Join(strs, "\n")
}

func (d *StInfoMgr) Parse() string {
	data := d.Data[d.DataName]
	strs := []string{}
	for _, item := range d.Funs {
		key := util.Ifelse(len(item.Refs) == 1, item.Refs[0], fmt.Sprintf("%s{%s}", item.Index, strings.Join(item.Refs, ", ")))
		switch item.ValueOf {
		case domain.VALUE_OF_MAP:
			strs = append(strs, fmt.Sprintf("data.%s[%s] = item", item.Name, key))
		case domain.VALUE_OF_GROUP:
			strs = append(strs, fmt.Sprintf("key := %s", key))
			strs = append(strs, fmt.Sprintf("data.%s[key] = append(data.%s[key], item)", item.Name, item.Name))
		}
	}
	return fmt.Sprintf(parse, d.CType, d.DataName, strings.Join(data.Makes, "\n"), strings.Join(strs, "\n"))
}

const (
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
		return d.listData[0]
	}
	return nil
}	
	`
	lget = `
func LGet() (rets []*%s) {
	if d, ok := obj.Load().(*%s); ok {
		rets = make([]*%s, len(d.listData))	
		copy(rets, d.listData)
	}
	return
}	
	`
	mget = `
func MGet%d(%s) *%s {
	if d, ok := obj.Load().(*%s); ok {
		val, ok := d.mapData%d[%s]
		if ok {
			return val	
		}
	}
	return nil
}	
	`
	gget = `
func GGet%d(%s) (rets []*%s) {
	if d, ok := obj.Load().(*%s); ok {
		vals, ok := d.groupData%d[%s]
		if ok {
			rets = make([]*%s, len(vals))	
			copy(rets, vals)
		}	
	}
	return
}
	`
	parse = `
func Parse(buf []byte) {
	ary := []*%s{}
	if err := json.Unmarshal(buf, &ary); err != nil {
		panic(err)
	}
	data := &%s{
		%s
	}
	for _, item := range ary {
		data.listData = append(data.listData, item)
		%s
	}
	obj.Store(data)
}
	`
)
