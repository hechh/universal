package mysql

import (
	"universal/common/mysql/internal/client"
	"universal/common/mysql/internal/manager"
	"universal/common/yaml"
)

func Register(dbname string, tables ...interface{}) {
	manager.RegisterTable(dbname, tables...)
}

func Init(cfgs map[int32]*yaml.DbConfig) error {
	return manager.InitMysql(cfgs)
}

func GetClient(dbname string) *client.MysqlClient {
	return manager.GetMysql(dbname)
}
