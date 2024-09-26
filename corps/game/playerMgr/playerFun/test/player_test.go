package test

import (
	"corps/common"
	"corps/common/cfgData"
	"corps/common/orm/redis"
	"corps/common/report"
	"corps/framework/plog"
	"corps/pb"
	"corps/server/game/playerMgr/playerFun"
	"testing"
)

var (
	player = &TestPlayer{
		FunCommon:     playerFun.FunCommon{MapFun: nil},
		mapPlayerData: make(map[pb.PlayerDataType]playerFun.IPlayerFun),
	}
)

type TestPlayer struct {
	playerFun.FunCommon
	mapPlayerData map[pb.PlayerDataType]playerFun.IPlayerFun //玩家数据
	list          []playerFun.IPlayerFun
}

func (d *TestPlayer) Init() {
	d.FunCommon.MapFun = &d.mapPlayerData
}

// 注册
func (d *TestPlayer) Register(pbType pb.PlayerDataType, fun playerFun.IPlayerFun) {
	d.list = append(d.list, fun)
	fun.Init(pbType, &d.FunCommon)
	d.mapPlayerData[pbType] = fun
}

func (d *TestPlayer) GetIPlayerFun(typ pb.PlayerDataType) playerFun.IPlayerFun {
	return d.mapPlayerData[typ]
}

func (d *TestPlayer) NewPlayer() {
	for _, val := range d.list {
		val.LoadPlayerDBFinish()
		val.NewPlayer()
	}
}

func TestMain(m *testing.M) {
	// 初始化日志
	plog.Init("game", plog.WithServerId(1))
	plog.SetLevel(uint32(plog.LOG_WARN | plog.LOG_DEBUG | plog.LOG_DEFAULT))
	// 初始化数据上报
	report.InitMock(&TestReportClient{})
	// 初始化packet
	// 初始化redis
	redis.Init(common.RedisCommon{
		Port:     6379,
		PreKey:   "corps",
		Password: "link123!",
		Redis:    map[uint32]common.RedisInfo{1: {Ip: "172.16.126.208"}},
	})

	// 初始化配置文件
	cfgData.LOADDATAMGR.Init("../../../../../../../share/serverVersion/linux/corps/data/")
	if !cfgData.LOADDATAMGR.LoadFile() {
		panic("CfgData load is failed")
	}
	// 执行
	m.Run()
}
