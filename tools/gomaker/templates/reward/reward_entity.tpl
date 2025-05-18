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
	"forevernine.com/planet/server/common/dao/repository/redis/{{$pdatahis}}"
	pb "forevernine.com/planet/server/common/pbclass"
	"forevernine.com/planet/server/common/xerrors"
	{{if .ActivityType}} "forevernine.com/planet/server/common/action" {{end}}
	"forevernine.com/planet/server/srv/gamesrv/internal/online"
)

// 公共数据模块
type Entity struct {
	session *online.UserSession
{{if .ActivityType}} actionCfg *pb.ActivityConfigS {{end}}
	dataHis *pb.{{$datahis}}
{{if $data}} data *pb.{{$data}} {{end}}
	isnew bool
}

// 发奖回调接口
func (d *Entity) AfterReward(ctx *framework.Context, item *pb.RewardResult, scene pb.RewardScene, r *pb.BaseReport, extra interface{}) (err error) {
	// todo：请实现业务逻辑
	return
}

// 消耗接口回调
func (d *Entity) AfterConsume(ctx *framework.Context, item *pb.RewardResult, scene pb.RewardScene, r *pb.BaseReport, extra interface{}) (err error) {
	return
}

{{if .ActivityType}}
// 取数据接口
func GetEntity(ctx *framework.Context, extra interface{}) (result interface{}) {
	sess, err := online.GetUserSessionByCtx(ctx)
	if err != nil {
		return xerrors.ErrPublicDaoGet().Format("Getting UserSession is failed, error: %v", err)
	}
	actionCfg := action.GetOnline(ctx, pb.{{$act}})
	if actionCfg == nil {
		return xerrors.ErrPublicActivityNotOnline().Format("{{$act}} is not online")
	}
	{{if .Data}} 
	datahis, data, isnew, err := {{$pdatahis}}.GetByUUID(ctx.UID, actionCfg.GetUniqueID())
	if err != nil {
		return xerrors.ErrPublicDaoGet().Format("Getting {{$datahis}} is failed, error: %v", err)
	}
	return &Entity{actionCfg: actionCfg, dataHis: datahis, data: data, isnew: isnew, session: sess}
	{{else}}
	datahis, isnew, err := {{$pdatahis}}.GetByUUID(ctx.UID, actionCfg.GetUniqueID())
	if err != nil {
		return xerrors.ErrPublicDaoGet().Format("Getting {{$datahis}} is failed, error: %v", err)
	}
	return &Entity{actionCfg: actionCfg, dataHis: datahis, isnew: isnew, session: sess}
	{{end}}
}
{{else}}
// 取数据接口
func GetEntity(ctx *framework.Context, extra interface{}) (result interface{}) {
	sess, err := online.GetUserSessionByCtx(ctx)
	if err != nil {
		return xerrors.ErrPublicDaoGet().Format("Getting UserSession is failed, error: %v", err)
	}
	{{if .IsHash}} 
	datahis, isexist, err := {{$pdatahis}}.HGet(ctx.UID)
	if err != nil {
		return xerrors.ErrPublicDaoGet().Format("Getting {{$datahis}} is failed, error: %v", err)
	}
	return &Entity{dataHis: datahis, isnew: !isexist, session: sess}
	{{else}}
	datahis, isexist, err := {{$pdatahis}}.Get(ctx.UID)
	if err != nil {
		return xerrors.ErrPublicDaoGet().Format("Getting {{$datahis}} is failed, error: %v", err)
	}
	return &Entity{dataHis: datahis, isnew: !isexist, session: sess}
	{{end}}
}
{{end}}

{{if .IsHash}} 
// 保存数据接口
func SetEntity(ctx *framework.Context, extra interface{}, data interface{}) (err error) {
	entity, ok := data.(*Entity)
	if !ok {
		return xerrors.ErrPublicParameter().Format("{{$package}}.Entity is expected")
	}
	if entity == nil || entity.dataHis == nil {
		return xerrors.ErrPublicParameter().Format("{{$package}}.Entity is nil")
	}
	if retErr := {{$pdatahis}}.HSet(ctx.UID, entity.dataHis); retErr != nil {
		return xerrors.ErrPublicDaoSet().Format("Setting {{$datahis}} is failed, error: %v", retErr)
	}
	return
}
{{else}}
// 保存数据接口
func SetEntity(ctx *framework.Context, extra interface{}, data interface{}) (err error) {
	entity, ok := data.(*Entity)
	if !ok {
		return xerrors.ErrPublicParameter().Format("{{$package}}.Entity is expected")
	}
	if entity == nil || entity.dataHis == nil {
		return xerrors.ErrPublicParameter().Format("{{$package}}.Entity is nil")
	}
	if retErr := {{$pdatahis}}.Set(ctx.UID, entity.dataHis); retErr != nil {
		return xerrors.ErrPublicDaoSet().Format("Setting {{$datahis}} is failed, error: %v", retErr)
	}
	return
}
{{end}}
