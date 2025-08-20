package tpl

import "text/template"

const packageTpl = `
/*
* 本代码由pbtool工具生成，请勿手动修改
*/

package {{.}}

import (
	"github.com/golang/protobuf/proto"
)

`

const memberTpl = `
{{$stname := .GetName}}
{{range $field := .GetAll -}}
func (d *{{$stname}}) Set{{$field.GetName}}(v {{$field.FullName "pb"}}) {
	d.{{$field.GetName}} = v
}
{{end}}

`

const factoryTpl = `
var (
	factorys = make(map[string]func() proto.Message)
)

func init() {
{{- range $cls := .}}
{{- if eq $cls.GetKind 3}}
	factorys["{{$cls.GetName}}"] = func() proto.Message { return &{{$cls.GetName}}{} }
{{- end}}
{{- end}}
}

`

var (
	PackageTpl *template.Template = template.Must(template.New("pb").Parse(packageTpl))
	MemberTpl  *template.Template = template.Must(template.New("pb").Parse(memberTpl))
	FactoryTpl *template.Template = template.Must(template.New("pb").Parse(factoryTpl))
)
