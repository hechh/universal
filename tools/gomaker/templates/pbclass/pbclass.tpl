

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

func registerType(val interface{}) {
	ttt := reflect.TypeOf(val).Elem()
	types[ttt.Name()] = ttt
}

func init() {
{{range $v := .}} registerType((*{{$v.Type.GetType ""}})(nil))
{{end}}
}
