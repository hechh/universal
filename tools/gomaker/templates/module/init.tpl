package {{.PackageName}}

import (
	"forevernine.com/base/srvcore/framework"
	"forevernine.com/planet/server/common/action"
	pb "forevernine.com/planet/server/common/pbclass"
)

func Init() {
	action.RegisterActivity(&Action{})
{{range $val := .ApiList}} framework.RegisterCMD(pb.{{$val.CmdReq}}, &pb.{{$val.Req}}{}, {{$val.FuncName}}, &pb.{{$val.Rsp}}{})
{{end}} {{range $val := .EventList}} framework.RegisterCMD(pb.{{$val.CmdEvent}}, &pb.{{$val.Event}}{}, {{$val.FuncName}}, nil)
{{end}}
}
