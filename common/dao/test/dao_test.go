package test

import (
	"sync"
	"testing"
	"time"
	"universal/common/dao/internal/manager"
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
	// 注册table
	if err := manager.CreateTable("corps_game_1", new(UserData)); err != nil {
		panic(err)
	}
	if err := manager.CreateTable("corps_common", new(TPlayerName)); err != nil {
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

type TPlayerName struct {
	PlayerName string    `orm:"column(PlayerName);pk"` // 可以指定大小
	AccountID  uint64    `orm:"column(AccountID)"`     // 使用空格分隔标签
	UpdateTime time.Time `orm:"column(UpdateTime)"`    // 指定类型
}

func (d *TPlayerName) TableName() string {
	return "t_player_name"
}

// go test -race -run=TestPlayerName ./* -count=1 -v
func TestPlayerName(t *testing.T) {
	cli, err := manager.GetMysql("corps_common")
	if err != nil {
		t.Log("----error----->", err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		player := &TPlayerName{PlayerName: "asdwewqr"}
		if err := cli.Read(player); err != nil {
			t.Log(err)
		} else {
			t.Log("------->", player)
		}
	}()
	go func() {
		defer wg.Done()
		newer := &TPlayerName{}
		cli.QueryTable(newer.TableName()).Filter("AccountID", 100100120).All(newer)
		t.Log("----update---->", newer)
	}()
	wg.Wait()
}

type UserData struct {
	AccountID     uint64 `orm:"column(AccountID);pk"`
	PlayerLevel   uint64 `orm:"column(PlayerLevel)"`
	PlayerName    string `orm:"column(PlayerName)"`
	CreateTime    uint64 `orm:"column(CreateTime)"`
	LastDailyTime uint64 `orm:"column(LastDailyTime)"`
	BaseData      string `orm:"column(BaseData);type(mediumblob)"`
	SystemData    string `orm:"column(SystemData);type(mediumblob)"`
	BagData       string `orm:"column(BagData);type(mediumblob)"`
	EquipmentData string `orm:"column(EquipmentData);type(mediumblob)"`
	ClientData    string `orm:"column(ClientData);type(mediumblob)"`
	MailData      string `orm:"column(MailData);type(mediumblob)"`
	CrystalData   string `orm:"column(CrystalData);type(mediumblob)"`
	HeroData      string `orm:"column(HeroData);type(mediumblob)"`
}

func (d *UserData) TableName() string {
	return "t_player"
}

func TestUserData(t *testing.T) {
	o, err := manager.GetMysql("corps_game_1")
	if err != nil {
		t.Log("----error----->", err)
		return
	}
	data := &UserData{AccountID: 100000023}
	if err := o.Read(data); err != nil {
		t.Log(err)
	}
	t.Log("------->", data.AccountID, data.PlayerLevel, data.PlayerName, data.CreateTime, data.LastDailyTime, []byte(data.BaseData))
}
