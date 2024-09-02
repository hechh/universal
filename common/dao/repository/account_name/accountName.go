package account_name

import (
	"time"
	"universal/common/dao/internal/manager"

	"github.com/astaxie/beego/orm"
)

type AccountName struct {
	Name       string    `orm:"column(AccountName)"`  // 可以指定大小
	AccountID  uint64    `orm:"column(AccountID);pk"` // 使用空格分隔标签
	UpdateTime time.Time `orm:"column(CreateTime)"`   // 指定类型
}

func (d *AccountName) TableName() string {
	return "t_account_name"
}

func init() {
	orm.RegisterModel(new(AccountName))
}

func Get(dbname string, uid uint64) (*AccountName, error) {
	cli, err := manager.NewOrmer(dbname)
	if err != nil {
		return nil, err
	}
	result := &AccountName{AccountID: uid}
	if err := cli.Read(result); err != nil {
		return nil, err
	}
	return result, nil
}

func Update(dbname string, data *AccountName) error {
	cli, err := manager.NewOrmer(dbname)
	if err != nil {
		return err
	}
	_, err = cli.InsertOrUpdate(data)
	return err
}
