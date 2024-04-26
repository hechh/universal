
{{$pkgName := .PkgName}}
{{$name := .Name}}
{{$field := .Field}}
{{$st := .Struct}}
{{$pbname := .Struct.GetTypeString $pkgName}}
{{$ftype := .Field.GetTypeString $pkgName}}

type {{$name}}Entity struct {
/*
{{range $f := $st.List}} {{if ne $f.Name $field.Name}} {{$.FirstCharToLower $f.Name}} {{$f.GetTypeString $pkgName}}
{{end}} {{end}} 
*/
    datas map[{{$ftype}}]*{{$pbname}}
}

func New{{$name}}Entity(list ...*{{$pbname}}) *{{$name}}Entity {
    ret := &{{$name}}Entity{
        datas: make(map[{{$ftype}}]*{{$pbname}}), 
    }
    for _, item := range list {
        ret.datas[item.Get{{$field.Name}}()] = item
    }
    return ret
}

func (e *{{$name}}Entity) ToProto() (rets []*{{$pbname}}) {
    for _, item := range e.datas {
        rets = append(rets, item)
    }
    return
}

