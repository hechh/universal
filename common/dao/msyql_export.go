package dao

import (
	"universal/common/dao/internal/manager"
	"universal/common/dao/internal/mysql"
	"universal/common/yaml"
)

func RegisterMysqlTable(dbname string, tables ...interface{}) {
	manager.RegisterTable(dbname, tables...)
}

func InitMysql(cfgs map[int32]*yaml.DbConfig) error {
	return manager.InitMysql(cfgs)
}

func GetMysql(dbname string) *mysql.OrmSql {
	return manager.GetMysql(dbname)
}
