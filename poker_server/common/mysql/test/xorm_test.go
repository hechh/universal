package test

import (
	"poker_server/common/config"
	"poker_server/common/mysql"
	"poker_server/common/pb"
	"poker_server/common/yaml"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	mysql.Register(mysql.MYSQL_DB_PLAYER_DATA, &pb.TexasPlayerFlowReport{}, &pb.TexasRoomReport{}, &pb.TexasGameReport{})
	mysql.Register(mysql.MYSQL_DB_PLAYER_DATA, &pb.PlayerData{})
}

func TestMain(m *testing.M) {
	if err := config.InitConfig("../../../../poker_gameconf/data", nil); err != nil {
		panic(err)
	}
	cfg, err := yaml.NewConfig("../../../env/local/config.yaml")
	if err != nil {
		panic(err)
	}
	if err := mysql.Init(cfg.Mysql); err != nil {
		panic(err)
	}
	m.Run()
}

func TestPlayer(t *testing.T) {
	session := mysql.GetClient(mysql.MYSQL_DB_PLAYER_DATA).NewSession()
	defer session.Close()

	query := session.Where("room_id = ?", 4311744514)
	query = query.And("round = ?", 2)
	// 设置分页
	query = query.Limit(int(10), 0)

	// 按时间降序排列
	query = query.Desc("begin_time")
	// 执行查询
	var list []*pb.TexasGameReport
	err := query.Find(&list)
	t.Log(err, list)
	/*
		// 查询
		usr := &pb.TexasGameReport{RoomId: 4311744514}
		ok, err := session.Get(usr)
		t.Log("查询", ok, err, usr)
	*/
}
