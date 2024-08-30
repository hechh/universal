package player_name

import (
	"time"
	"universal/common/dao/internal/manager"

	"github.com/astaxie/beego/orm"
)

type PlayerName struct {
	Name       string    `orm:"column(PlayerName);pk"` // 可以指定大小
	AccountID  uint64    `orm:"column(AccountID)"`     // 使用空格分隔标签
	UpdateTime time.Time `orm:"column(UpdateTime)"`    // 指定类型
}

func (d *PlayerName) TableName() string {
	return "t_player_name"
}

func init() {
	orm.RegisterModel(new(PlayerName))
}

func Get(dbname string, name string) (*PlayerName, error) {
	cli, err := manager.NewOrmer(dbname)
	if err != nil {
		return nil, err
	}
	result := &PlayerName{Name: name}
	if err := cli.Read(result); err != nil {
		return nil, err
	}
	return result, nil
}

func Query(dbname string, uid uint64) (*PlayerName, error) {
	cli, err := manager.NewOrmer(dbname)
	if err != nil {
		return nil, err
	}
	result := &PlayerName{}
	if err := cli.QueryTable(result.TableName()).Filter("AccountID", uid).One(result); err != nil {
		return nil, err
	}
	return result, nil
}

func Update(dbname string, data *PlayerName) error {
	cli, err := manager.NewOrmer(dbname)
	if err != nil {
		return err
	}
	_, err = cli.InsertOrUpdate(data)
	return err
}
