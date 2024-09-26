package test

import (
	"corps/base/cfgEnum"
	"corps/common"
	"corps/common/cfgData"
	"corps/common/orm/redis"
	"corps/framework/plog"
	"corps/pb"
	"corps/server/game/module/entry"
	"encoding/json"
	"testing"
)

func TestMain(m *testing.M) {
	// 初始化日志
	plog.Init("game", plog.WithServerId(1))
	plog.SetLevel(uint32(plog.LOG_WARN | plog.LOG_DEBUG | plog.LOG_DEFAULT))
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
	entry.Init()
	// 执行
	m.Run()
}

func TestEntity(t *testing.T) {
	str := `{"Book":{"Stage":1},"Robots":[{"RobotID":1,"Stage":1,"RoleSkillID":1003,"Crystals":[101,102,103,104,105,106]},{"RobotID":2,"Stage":1,"RoleSkillID":1000,"Crystals":[201,202,203,204,205,206]}],"Crystals":[{"CrystalID":20402,"Quality":5,"Star":1,"RewardCoinTimes":1,"PassiveSkillIds":[204020],"Level":3},{"CrystalID":50302,"Element":4,"Quality":4,"Star":1,"RewardCoinTimes":1,"PassiveSkillIds":[503020],"Level":4},{"CrystalID":50402,"Element":4,"Quality":5,"Star":1,"RewardCoinTimes":1,"PassiveSkillIds":[504020],"Level":3}],"Conditions":[{"CfgID":246,"Times":1},{"CfgID":233,"Times":1},{"CfgID":90,"Times":1}]}`
	data := &pb.PBPlayerCrystal{}
	json.Unmarshal([]byte(str), data)

	tmps := make(map[uint32]*pb.PBCrystal)
	for _, item := range data.Crystals {
		tmps[item.CrystalID] = item
	}

	service := entry.NewEntryService(100101411, data, nil, nil)
	service.SetCrystal(tmps)

	//vals := entry.KeyValueToDMap(service.Get(uint32(cfgEnum.EntryEffectType_OpenBox), 1202)...)
	vals := entry.KeyValueToMap(service.Get(uint32(cfgEnum.EntryEffectType_LootGroupLootRate), 79)...)
	//vals := entry.KeyValueToMap(service.Get(uint32(cfgEnum.EntryEffectType_AdertiseReward), uint32(cfgEnum.EAdvertType_BoxOpen))...)
	//vals := entry.ToValue(service.Get(uint32(cfgEnum.EntryEffectType_AdertiseReward), uint32(cfgEnum.EAdvertType_BoxOpen))...)
	t.Log(vals)
}
