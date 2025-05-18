package {{.PackageName}}

import (
	"forevernine.com/base/srvcore/framework"
	"forevernine.com/planet/server/common/plog"
	"forevernine.com/planet/server/common/xerrors"
	pb "forevernine.com/planet/server/common/pbclass"
)

func {{.FuncName}}(ireq, irsp framework.IProto, ctx *framework.Context) (err error) {
	plog.Trace(ctx, "event: %v", ireq)
	req := ireq.(*pb.{{.Event}})
	if req == nil {
		return xerrors.ErrPublicParameter().Format("req: %v", req)
	}

	return
}

