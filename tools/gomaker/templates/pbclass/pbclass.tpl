

import reflect "reflect"

var (
    types = make(map[string]reflect.Type)
)

func init() {
{{range $v := .}} types["{{$v.Type.Name}}"] = reflect.TypeOf((*{{$v.Type.GetType ""}})(nil)).Elem()
{{end}}
}
