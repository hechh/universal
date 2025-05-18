{{$pbname := .Type.Name}}
type {{$pbname}}S struct {
{{range $item := .Idents}} {{$item.GetNameS}} {{$item.Type.GetTypeS}}
{{end}} {{range $item := .Arrays}} {{$item.GetNameS}} {{$item.Type.GetTypeS}}
{{end}} {{range $item := .Maps}} {{$item.GetNameS}} map[{{$item.KType.GetTypeS}}]{{$item.VType.GetTypeS}}
{{end}} }

{{range $item := .Idents}} func (d *{{$pbname}}S) Get{{$item.Name}}() {{$item.Type.GetTypeS}} {
    return d.{{$item.GetNameS}}
} 
{{end}} 
{{range $item := .Arrays}} func (d *{{$pbname}}S) Get{{$item.Name}}() {{$item.Type.GetTypeS}} {
    return d.{{$item.GetNameS}}
} 
{{end}}
{{range $item := .Maps}} func (d *{{$pbname}}S) Get{{$item.Name}}() map[{{$item.KType.GetTypeS}}]{{$item.VType.GetTypeS}} {
    return d.{{$item.GetNameS}}
} 
{{end}}

func (val *{{$pbname}}S) ToProto() (ret *{{$pbname}}) {
    ret = &{{$pbname}}{
        {{range $item := .Idents}} {{$item.Name}}: {{if $item.Type.IsStruct}} val.{{$item.GetNameS}}.ToProto() {{else}} val.{{$item.GetNameS}} {{end}},
        {{end}} {{range $item := .Arrays}} {{$item.Name}}: {{if $item.Type.IsStruct}} To{{$item.Type.Name}}(val.{{$item.GetNameS}}) {{else}} val.{{$item.GetNameS}} {{end}},
        {{end}} {{range $item := .Maps}} {{$item.Name}}: make(map[{{$item.KType.GetType}}]{{$item.VType.GetType}}), 
        {{end}} }
    {{range $item := .Maps}} for key, val := range val.{{$item.GetNameS}} {
        ret.{{$item.Name}}[key] = {{if $item.Type.IsStruct}} val.{{$item.GetNameS}}.ToProto() {{else}} val.{{$item.GetNameS}} {{end}}
    } {{end}}
    return
}

func (val *{{$pbname}}) ToSecure() (ret *{{$pbname}}S) {
    ret = &{{$pbname}}S{
        {{range $item := .Idents}}  {{$item.GetNameS}}: {{if $item.Type.IsStruct}} val.{{$item.Name}}.ToSecure() {{else}} val.{{$item.Name}} {{end}},
        {{end}} {{range $item := .Arrays}} {{$item.GetNameS}}: {{if $item.Type.IsStruct}} To{{$item.Type.Name}}S(val.{{$item.Name}}) {{else}} val.{{$item.Name}} {{end}},
        {{end}} {{range $item := .Maps}} {{$item.GetNameS}}: make(map[{{$item.KType.GetTypeS}}]{{$item.VType.GetTypeS}}), 
        {{end}} }
    {{range $item := .Maps}} for key, val := range val.{{$item.Name}} {
        ret.{{$item.GetNameS}}[key] = {{if $item.Type.IsStruct}} val.{{$item.Name}}.ToSecure() {{else}} val.{{$item.Name}} {{end}}
    } {{end}}
    return
}

func To{{$pbname}}(vals []*{{$pbname}}S) (rets []*{{$pbname}}) {
    rets = make([]*{{$pbname}}, len(vals))
    for i, val := range vals {
        rets[i] = val.ToProto()
    }
    return
}

func To{{$pbname}}S(vals []*{{$pbname}}) (rets []*{{$pbname}}S) {
    rets = make([]*{{$pbname}}S, len(vals))
    for i, val := range vals {
        rets[i] = val.ToSecure()
    }
    return
}
