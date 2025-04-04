package account

import (
	"hego/common/dao/internal/manager"
	"time"

	"github.com/astaxie/beego/orm"
)

type Account struct {
	AccountID      uint64    `orm:"column(AccountID);pk"` // 使用空格分隔标签
	Name           string    `orm:"column(AccountName)"`  // 可以指定大小
	Password       string    `orm:"column(Password)"`     // 可以指定大小
	PlatType       int32     `orm:"column(PlatType)"`
	PlatSystemType int32     `orm:"column(PlatSystemType)"`
	LoginIP        string    `orm:"column(LoginIP)"`
	LoginTime      time.Time `orm:"column(LoginTime)"`
	FirstLoginIP   string    `orm:"column(FirstLoginIP)"`
	FirstLoginTime time.Time `orm:"column(FirstLoginTime)"`
	LoginNumber    int32     `orm:"column(LoginNumber)"`
}

func (d *Account) TableName() string {
	return "t_account"
}

func init() {
	orm.RegisterModel(new(Account))
}

func Get(dbname string, uid uint64) (*Account, error) {
	cli, err := manager.NewOrmer(dbname)
	if err != nil {
		return nil, err
	}
	result := &Account{AccountID: uid}
	if err := cli.Read(result); err != nil {
		return nil, err
	}
	return result, nil
}

func Update(dbname string, data *Account) error {
	cli, err := manager.NewOrmer(dbname)
	if err != nil {
		return err
	}
	_, err = cli.InsertOrUpdate(data)
	return err
}
