package base

import (
	"fmt"
	"strings"
)

type Param struct {
	Name string
	Type string
}

type Params []*Param

func (d *Param) UniqueID() string {
	return d.Name + "@" + d.Type
}

func (d *Param) Format() string {
	if d.Type == "string" {
		return "%s"
	}
	return "%d"
}

func (d *Param) Arg() string {
	return d.Name + " " + d.Type
}

func (d *Param) Val(str string) string {
	if len(str) > 0 {
		return str + "." + d.Name
	}
	return d.Name
}

func (d *Param) Init(str string) string {
	if len(str) > 0 {
		return d.Name + ": " + str + "." + d.Name
	}
	return d.Name + ": " + d.Name
}

// 类型转换，解决interface{}类型转换地问题
func (d *Param) Cast(str string) string {
	if strings.HasPrefix(d.Type, "pb.") {
		return fmt.Sprintf("%s(cast.ToInt32(%s))", d.Type, str)
	}
	switch d.Type {
	case "int":
		return fmt.Sprintf("cast.ToInt(%s)", str)
	case "int32":
		return fmt.Sprintf("cast.ToInt32(%s)", str)
	case "int64":
		return fmt.Sprintf("cast.ToInt64(%s)", str)
	case "uint":
		return fmt.Sprintf("cast.ToUint(%s)", str)
	case "uint32":
		return fmt.Sprintf("cast.ToUint32(%s)", str)
	case "uint64":
		return fmt.Sprintf("cast.ToUint64(%s)", str)
	case "string":
		return fmt.Sprintf("cast.ToString(%s)", str)
	}
	return fmt.Sprintf("%s.(%s)", str, d.Type)
}

func (d Params) Count() int {
	return len(d)
}

func (d Params) HasUID() bool {
	for _, item := range d {
		if strings.ToLower(item.UniqueID()) == "uid@string" {
			return true
		}
	}
	return false
}

func (d Params) UniqueID() string {
	return strings.Join(d.Vals(""), "")
}

func (d Params) Index() string {
	return "Index" + d.UniqueID()
}

// 函数、结构体参数格式；func(a int, b int)
func (d Params) Args() (result []string) {
	for _, item := range d {
		result = append(result, item.Arg())
	}
	return
}

func (d Params) Arg() string {
	return strings.Join(d.Args(), ",")
}

func (d Params) Struct() string {
	return strings.Join(d.Args(), ";")
}

// 函数调用、值引用参数格式：swap(a, b)
func (d Params) Vals(str string) (result []string) {
	for _, item := range d {
		result = append(result, item.Val(str))
	}
	return
}

func (d Params) Val(str string) string {
	return strings.Join(d.Vals(str), ",")
}

// 结构体初始化参数格式
func (d Params) Inits(str string) (result []string) {
	for _, item := range d {
		result = append(result, item.Init(str))
	}
	return
}

func (d Params) Init(str string) string {
	return strings.Join(d.Inits(str), ",")
}

// func(args ...interface{}) 类型参数使用
func (d Params) Casts(str string) (result []string) {
	for index, elem := range d {
		result = append(result, elem.Cast(fmt.Sprintf("%s[%d]", str, index)))
	}
	return
}

func (d Params) Cast(str string) string {
	return strings.Join(d.Casts(str), ",")
}

func (d Params) Formats() (rets []string) {
	for _, item := range d {
		rets = append(rets, item.Format())
	}
	return
}

func (d Params) Format() string {
	return strings.Join(d.Formats(), ":")
}

func (d Params) Join(as ...Params) (rets Params) {
	rets = append(rets, d...)
	for _, a := range as {
		rets = append(rets, a...)
	}
	return rets
}
