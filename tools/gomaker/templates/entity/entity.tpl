
{{$pkgName := .PkgName}}
{{$name := .Name}}
{{$field := .Field}}
{{$st := .Struct}}
{{$pbname := .Struct.GetType .PkgName}}
{{$ftype := .Field.GetType .PkgName}}

type {{$name}}Entity struct {
    datas map[{{$ftype}}]*{{$pbname}}
    changes []{{$ftype}}
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

func (e *{{$name}}Entity) IsChange() bool {
    return len(e.changes) > 0
}
