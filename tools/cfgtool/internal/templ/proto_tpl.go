package templ

import (
	"text/template"
	"universal/library/util"
)

const protoTpl = `
/*
* 本代码由cfgtool工具生成，请勿手动修改
*/

syntax = "proto3";

package universal;

option go_package = "universal/common/pb";

{{range $item := .ReferenceList -}}
import "{{$item}}.proto";
{{end}}

{{- range $item := .Enums}}
enum {{$item.Name}} {
	{{- range $field := $item.ValueList}}
	{{$field.Name}} = {{$field.Value}}; // {{$field.Desc}}
	{{- end}}
}
{{end}}

{{- range $item := .Structs}}
message {{$item.Name}} {
	{{- range $pos, $field := $item.FieldList}}
		{{- if eq $field.Type.ValueOf 1}}
	{{$field.Type.Name}} {{$field.Name}} = {{add $pos 1}}; // {{$field.Desc}}
		{{- else if eq $field.Type.ValueOf 2}} 
	repeated {{$field.Type.Name}} {{$field.Name}} = {{add $pos 1}}; // {{$field.Desc}}
		{{- end}} 
{{- end}}
}
{{end}}

{{- range $item := .Configs}}
message {{$item.Name}} {
	{{- range $pos, $field := $item.FieldList}}
		{{- if eq $field.Type.ValueOf 1}}
	{{$field.Type.Name}} {{$field.Name}} = {{add $pos 1}}; // {{$field.Desc}}
		{{- else if eq $field.Type.ValueOf 2}} 
	repeated {{$field.Type.Name}} {{$field.Name}} = {{add $pos 1}}; // {{$field.Desc}}
		{{- end}} 
{{- end}}
}

message {{$item.Name}}Ary { repeated {{$item.Name}} Ary = 1; }
{{end}}
`

var (
	ProtoTpl *template.Template
)

func init() {
	funcs := template.FuncMap{
		"sub": util.Sub[int],
		"add": util.Add[int],
	}
	ProtoTpl = template.Must(template.New("ProtoTpl").Funcs(funcs).Parse(protoTpl))
}
