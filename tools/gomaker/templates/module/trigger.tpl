package {{.PackageName}}

import (
	"forevernine.com/base/srvcore/framework"

	"forevernine.com/planet/server/common/plog"
	pb "forevernine.com/planet/server/common/pbclass"
	"forevernine.com/planet/server/srv/gamesrv/internal/module/shop/trigger"
)

func PayTrigger(ctx *framework.Context, typ pb.PayType, product *pb.PiggyProductConfig, localPrice int64, orderInfo trigger.OrderTrigger) error {
	plog.Trace(ctx, "PayTrigger PayType: %s, pid: %v", typ.String(), product)
	return nil
}

func ExchangeTrigger(ctx *framework.Context, pid string, num int64) error {
	plog.Trace(ctx, "ExchangeTrigger pid: %s", pid)
	return nil
}
