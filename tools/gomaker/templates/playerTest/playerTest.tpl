{{$name := .Name}}

import (
	"corps/pb"
	"corps/server/game/playerMgr/playerFun"
	"testing"

	"github.com/stretchr/testify/assert"
)

{{range $req := .ReqList}} {{$rsp := $.Join ($.TrimSuffix $req "Request") "Response"}}
func Test_{{$req}}(t *testing.T) {
    obj := &playerFun.{{$name}}{}
    obj.NewPlayer()
	head := &pb.RpcHead{Id: 100000123}
    req := &pb.{{$req}}{}    
    rsp := &pb.{{$rsp}}{}
    err := obj.{{$req}}(head, req, rsp)
    assert.Equal(t, nil, err)
}
{{end}}