package mysql

import (
	"poker_server/common/mysql/internal/client"
	"poker_server/common/mysql/internal/manager"
	"poker_server/common/yaml"
)

func Register(dbname string, tables ...interface{}) {
	manager.RegisterTable(dbname, tables...)
}

func Init(cfgs map[int32]*yaml.DbConfig) error {
	return manager.InitMysql(cfgs)
}

func GetClient(dbname string) *client.OrmSql {
	return manager.GetMysql(dbname)
}
