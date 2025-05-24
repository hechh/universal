package templ

const codeTpl = `
{{/* 定义全局变量  */}}
{{$type := .Name}} 
{{$indexs := .IndexList}}
{{$indexMap := .Indexs}}
{{$pkg := .PbPkg}}

/*
* 本代码由cfgtool工具生成，请勿手动修改
*/

package {{.Pkg}}

import (
	"encoding/json"
	"sync/atomic"

	"poker_server/common/pb"
	"poker_server/common/config"
)

var obj = atomic.Value{}

type {{$type}}Data struct {
{{- range $index := $indexs -}}
    {{- if eq $index.Type.ValueOf 2 -}}         {{/*ValueOfList*/}}
        _{{$index.Name}} []*{{$pkg}}.{{$type}}
    {{- else if eq $index.Type.ValueOf 3 -}}    {{/*ValueOfMap*/}}
        _{{$index.Name}} map[{{$index.Type.Name}}]*{{$pkg}}.{{$type}}
    {{- else if eq $index.Type.ValueOf 4 -}}    {{/*ValueOfGroup*/}}
        _{{$index.Name}} map[{{$index.Type.Name}}][]*{{$pkg}}.{{$type}}
    {{- end -}}
{{- end -}}
}

// 注册函数
func init() {
    config.Register("{{$type}}", parse)
}

func parse(buf string) error {
    data := &{{$pkg}}.{{$type}}Ary{}
    if err := json.Unmarshal([]byte(buf), data); err != nil {
        return err
    }

{{if or (index $indexMap 3) (index $indexMap 4)}}
{{range $index := $indexs -}}
    {{- if eq $index.Type.ValueOf 3 -}}    {{/*ValueOfMap*/}}
        _{{$index.Name}} := make(map[{{$index.Type.Name}}]*{{$pkg}}.{{$type}})
    {{- else if eq $index.Type.ValueOf 4 -}}    {{/*ValueOfGroup*/}}
        _{{$index.Name}} := make(map[{{$index.Type.Name}}][]*{{$pkg}}.{{$type}})
    {{- end -}}  
{{- end}}
    for _, item :=range data.Ary {
{{- range $index := $indexs -}} 
    {{$key := $index.Value "item" ","}}
    {{- if eq $index.Type.ValueOf 3 -}}    {{/*ValueOfMap*/}}
        {{- if or (eq $index.Type.TypeOf 1) (eq $index.Type.TypeOf 2) -}} {{/*TypeOfBase*/}}
            _{{$index.Name}}[{{$key}}] = item
        {{- else if eq $index.Type.TypeOf 3 -}} {{/*TypeOfStruct*/}}
            _{{$index.Name}}[{{$index.Type.Name}}{ {{$key}} }] = item
        {{- end -}}
    {{- else if eq $index.Type.ValueOf 4 -}}    {{/*ValueOfGroup*/}}
        {{- if or (eq $index.Type.TypeOf 1) (eq $index.Type.TypeOf 2) -}} {{/*TypeOfBase*/}}
            _{{$index.Name}}[{{$key}}] = append(_{{$index.Name}}[{{$key}}], item)
        {{- else if eq $index.Type.TypeOf 3 -}} {{/*TypeOfStruct*/}}
            _{{$index.Name}}[{{$index.Type.Name}}{ {{$key}} }] = append(_{{$index.Name}}[{{$index.Type.Name}}{ {{$key}} }], item)
        {{- end -}}
    {{- end -}}  
{{- end -}}
    }
{{end}}
    obj.Store(&{{$type}}Data{
{{- range $index := $indexs}} 
    {{- if or (eq $index.Type.ValueOf 3) (eq $index.Type.ValueOf 4)}}
        _{{$index.Name}}: _{{$index.Name}},
    {{- else}}
        _{{$index.Name}}: data.Ary,
    {{- end -}}  
{{- end}}
    })
    return nil
}

{{$index := index (index $indexMap 2) 0}}
{{if $index -}}
func SGet() *{{$pkg}}.{{$type}} {
    obj, ok := obj.Load().(*{{$type}}Data)
    if !ok {
        return nil
    }
    return obj._{{$index.Name}}[0]
}

func LGet() (rets []*{{$pkg}}.{{$type}}) {
    obj, ok := obj.Load().(*{{$type}}Data)
    if !ok {
        return
    }
    rets = make([]*{{$pkg}}.{{$type}}, len(obj._{{$index.Name}}))
    copy(rets, obj._{{$index.Name}})
    return
}

func Walk(f func(*{{$pkg}}.{{$type}})bool) {
    obj, ok := obj.Load().(*{{$type}}Data)
    if !ok {
        return
    }
    for _, item := range obj._{{$index.Name}} {
        if !f(item) {
            return
        }
    }
}
{{- end}}

{{- range $index := $indexs -}} 
    {{$arg := $index.Arg ","}}
    {{$key := $index.Value "" ","}}
    {{- if eq $index.Type.ValueOf 3 -}}    {{/*ValueOfMap*/}}
func MGet{{$index.Name}}({{$arg}}) *{{$pkg}}.{{$type}} {
    obj, ok := obj.Load().(*{{$type}}Data)
    if !ok {
        return nil
    }
    {{if or (eq $index.Type.TypeOf 1) (eq $index.Type.TypeOf 2) -}}                       {{/*TypeOfBase*/}}
        if val, ok := obj._{{$index.Name}}[{{$key}}]; ok {
    {{- else if eq $index.Type.TypeOf 3 -}}                                                 {{/*TypeOfStruct*/}}
        if val, ok := obj._{{$index.Name}}[{{$index.Type.Name}}{ {{$key}} }]; ok {
    {{- end}}
        return val
    }
    return nil
}
    {{- else if eq $index.Type.ValueOf 4 -}}    {{/*ValueOfGroup*/}}
func GGet{{$index.Name}}({{$arg}}) (rets []*{{$pkg}}.{{$type}}) {
    obj, ok := obj.Load().(*{{$type}}Data)
    if !ok {
        return
    }
    {{- if or (eq $index.Type.TypeOf 1) (eq $index.Type.TypeOf 2) -}} {{/*TypeOfBase*/}}
        if vals, ok := obj._{{$index.Name}}[{{$key}}]; ok {
    {{- else if eq $index.Type.TypeOf 3 -}} {{/*TypeOfStruct*/}}
        if vals, ok := obj._{{$index.Name}}[{{$index.Type.Name}}{ {{$key}} }]; ok {
    {{- end -}}
        rets = make([]*{{$pkg}}.{{$type}}, len(vals))
        copy(rets, vals)
        return
    }
    return
}
    {{- end -}}  
{{- end -}}

`
