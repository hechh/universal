package tpl

import "text/template"

const packageTpl = `
/*
* 本代码由pbtool工具生成，请勿手动修改
*/

package {{.}}

`

const memberTpl = `
{{$stname := .GetName}}
{{range $field := .GetAll -}}
func (d *{{$stname}}) Set{{$field.GetName}}(v {{$field.FullName "pb"}}) {
	d.{{$field.GetName}} = v
}
{{end}}

`

var (
	PackageTpl *template.Template = template.Must(template.New("pb").Parse(packageTpl))
	MemberTpl  *template.Template = template.Must(template.New("pb").Parse(memberTpl))
)
