

import (
    "reflect"
	"google.golang.org/protobuf/proto"
)

var (
    types = make(map[string]reflect.Type)
)

func NewType(name string) proto.Message {
    if vv, ok := types[name]; ok {
        return reflect.New(vv).Interface().(proto.Message)
    }
    return nil
}

func init() {
{{range $v := .}} types["{{$v.Type.Name}}"] = reflect.TypeOf((*{{$v.Type.GetType ""}})(nil)).Elem()
{{end}}
}
