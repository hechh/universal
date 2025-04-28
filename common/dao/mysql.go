package dao

import (
	"universal/common/dao/internal/manager"
	"universal/common/global"

	"github.com/astaxie/beego/orm"
)

func InitMysql(cfgs map[uint32]*global.DbConfig) error {
	return manager.InitMysql(cfgs)
}

func NewOrmer(dbname string) (orm.Ormer, error) {
	return manager.NewOrmer(dbname)
}
