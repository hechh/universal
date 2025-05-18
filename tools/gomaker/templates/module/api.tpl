package {{.PackageName}}

import (
	"forevernine.com/base/srvcore/framework"
	"forevernine.com/planet/server/common/xerrors"
	"forevernine.com/planet/server/common/plog"
	pb "forevernine.com/planet/server/common/pbclass"
)

func {{.FuncName}}(ireq, irsp framework.IProto, ctx *framework.Context) (err error) {
	plog.Trace(ctx, "req: %v", ireq)
	req := ireq.(*pb.{{.Req}})
	rsp := irsp.(*pb.{{.Rsp}})
	if req == nil || rsp == nil {
		return xerrors.ErrPublicParameter().Format("req: %v, rsp: %v", req, rsp)
	}
	rsp.Head = &pb.RspHead{}

	return
}

