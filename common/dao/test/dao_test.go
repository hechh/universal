package test

import (
	"testing"
	"universal/common/dao/internal/manager"
	"universal/common/dao/repository/player_data"
	"universal/common/dao/repository/player_name"
	"universal/common/global"

	_ "github.com/go-sql-driver/mysql"
)

func TestMain(m *testing.M) {
	cfg := &global.Config{}
	global.LoadFile("../../../env/yaml/common.yaml", cfg)
	global.LoadFile("../../../env/yaml/game.yaml", cfg)
	// 初始化redis
	if err := manager.InitRedis(cfg.Redis); err != nil {
		panic(err)
	}
	// 初始化mysql
	if err := manager.InitMysql(cfg.Mysql); err != nil {
		panic(err)
	}
	m.Run()
}

func TestRedis(t *testing.T) {
	cli := manager.GetRedis(1)
	if cli == nil {
		return
	}

	result, err := cli.Get("test-hch")
	t.Log(result, err)

	if err = cli.Set("test-hch", 123); err != nil {
		t.Log(err)
		return
	}

	result, err = cli.Get("test-hch")
	t.Log(result, err)
}

func TestPlayer(t *testing.T) {
	info, err := player_name.Get("corps_common", "asdwewqr")
	t.Log(err, info)

	info, err = player_name.Query("corps_common", 100100120)
	t.Log(err, info)

	data, err := player_data.Get("corps_game_2", 100101202)
	t.Log(err, data)
}
