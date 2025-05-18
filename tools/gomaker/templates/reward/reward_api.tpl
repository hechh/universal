{{/* 定义变量 */}}
{{$package := .Package}}
{{$pdatahis := .DaoPackage}}
{{$datahis := .DataHis}}
{{$data := .Data}}
{{$act := .ActivityType}}
{{$propertys := .PropertyTypes}}


package {{$package}}

import (
	pb "forevernine.com/planet/server/common/pbclass"
	"forevernine.com/planet/server/srv/gamesrv/internal/reward"
)

const ID = "{{$package}}"

// 业务发奖接口、道具转换接口 注册
func Init() {
	// 道具转换注册
	{{range $t := $propertys}} reward.RegisterConvert(pb.{{$t}}, {{$t}}Convert)
	{{end}}

	// 道具发奖注册
	propertyMgr := reward.NewBaseReward(ID, GetEntity, SetEntity)
	{{range $t := $propertys}} propertyMgr.Register(New{{$t}}, pb.{{$t}})
	{{end}}

	// 注册发奖接口
	reward.RegisterIReward(ID, propertyMgr)
}
