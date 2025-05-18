{{/* 定义变量 */}}
{{$package := .Package}}
{{$pdatahis := .DaoPackage}}
{{$datahis := .DataHis}}
{{$data := .Data}}
{{$act := .ActivityType}}
{{$property := .PropertyType}}

package {{$package}}

import (
	"forevernine.com/base/srvcore/framework"
	pb "forevernine.com/planet/server/common/pbclass"
	"forevernine.com/planet/server/common/xerrors"

	"forevernine.com/planet/server/srv/gamesrv/internal/report"
	"forevernine.com/planet/server/srv/gamesrv/internal/reward"
)

// 道具转换接口
func {{$property}}Convert(ctx *framework.Context, item *pb.Reward) (result *pb.Reward, err error) {
    // todo: 请自行实现接口逻辑
    result = item
	return
}

// new接口
func New{{$property}}(data interface{}) interface{} {
	if v, ok := data.(*Entity); ok {
		return &{{$property}}{Entity: v}
	}
	return reward.NewInvalidProperty(xerrors.ErrPublicParameter().Format("*{{$property}}.Entity is exptected"))
}

type {{$property}} struct {
	*Entity
}

func (d *{{$property}}) GetStock(ctx *framework.Context, p pb.PropertyType, result map[pb.PropertyType]*pb.StockInfo, extra interface{}) {
	// todo：请实现业务逻辑
	return
}

func (d *{{$property}}) Reward(ctx *framework.Context, item *pb.Reward, scene pb.RewardScene, r *pb.BaseReport, extra interface{}) (result *pb.RewardResult, re report.IReport, err error) {
	// todo：请实现业务逻辑
	return
}

func (d *{{$property}}) Consume(ctx *framework.Context, item *pb.Reward, scene pb.RewardScene, r *pb.BaseReport, extra interface{}) (result *pb.RewardResult, re report.IReport, err error) {
	// todo：请实现业务逻辑
	return
}
