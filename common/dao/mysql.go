package dao

import (
	"hego/common/dao/internal/manager"
	"hego/common/global"

	"github.com/astaxie/beego/orm"
)

func InitMysql(cfgs map[uint32]*global.DbConfig) error {
	return manager.InitMysql(cfgs)
}

func NewOrmer(dbname string) (orm.Ormer, error) {
	return manager.NewOrmer(dbname)
}
