package test

import (
	"testing"
	"time"
	"universal/common/dao/internal/manager"
	"universal/common/dao/internal/orm"
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
		return
	}

	// 初始化mysql
	if err := orm.RegisterDriver("mysql", orm.DRMySQL); err != nil {
		panic(err)
		return
	}

	orm.RegisterDataBase("corps_game_1", "mysql", "root:link123!@tcp(172.16.126.208:3306)/corps_game_1?charset=utf8")
	if err := orm.RegisterDataBase("default", "mysql", "root:link123!@tcp(172.16.126.208:3306)/corps_common?charset=utf8"); err != nil {
		panic(err)
		return
	}

	// 注册table
	orm.RegisterModel(new(TPlayerName), new(UserData))
	//orm.RegisterModel(new(TPlayerName))

	m.Run()
}

type TPlayerName struct {
	AccountID  uint64    `orm:"column(AccountID);pk"` // 使用空格分隔标签
	Name       string    `orm:"column(PlayerName)"`   // 可以指定大小
	UpdateTime time.Time `orm:"column(UpdateTime)"`   // 指定类型
}

func (d *TPlayerName) TableName() string {
	return "t_player_name"
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

func TestMysql(t *testing.T) {
	t.Run("t_player_name", func(t *testing.T) {
		player := &TPlayerName{AccountID: 100100596}
		o := orm.NewOrm()
		if err := o.Read(player); err != nil {
			t.Log(err)
		} else {
			t.Log("------->", player)
		}
	})

	t.Run("t_player", func(t *testing.T) {
		data := &UserData{AccountID: 100000002}
		o := orm.NewOrm()
		o.Using("corps_game_1")
		if err := o.Read(data); err != nil {
			t.Log(err)
		} else {
			t.Log("------->", data)
		}
	})
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
