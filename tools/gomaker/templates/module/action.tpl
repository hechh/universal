package {{.PackageName}}

import (
	"forevernine.com/base/srvcore/framework"
	"forevernine.com/base/srvcore/libs/xtime"
	pb "forevernine.com/planet/server/common/pbclass"
)

const (
	ActivityName = ""
	RedisData = ""
)

type Action struct{}

func (d *Action) GetProto() string {
	return RedisData
}

func (d *Action) GetName() string {
	return ActivityName
}

func (d *Action) GetType() pb.ActivityType {
	{{if .ActivityType}} 
	return pb.{{.ActivityType}}
	{{else}}
	return pb.ActivityTypeInit
	{{end}}
}

func (d *Action) GetInfo(ctx *framework.Context, act *pb.ActivityConfigS, redisData interface{}) *pb.ModActActivityCenterDetail {
	if act == nil {
		return nil
	}

	// 返回弹窗数据
	return &pb.ModActActivityCenterDetail{
		Id:         d.GetName(),
		Location:   act.GetLocation(),
		SkinIcon:   act.GetSkinIcon(),
		SkinName:   act.GetSkinName(),
		SkinBanner: act.GetSkinBanner(),
		Countdown:  act.GetEndTime() - xtime.Unix(),
		IsFirstRp:  1,
	}
}
