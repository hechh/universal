package manager

import (
	"database/sql"
	"fmt"
	"universal/common/global"
	"universal/framework/basic/util"

	"github.com/astaxie/beego/orm"
)

const (
	DriverName = "mysql"
)

var (
	mysqlPool = make(map[string]*MysqlConn)
)

type MysqlConn struct {
	conn       *sql.DB
	driverName string
	dbName     string
	alias      string
}

func InitMysql(cfgs map[uint32]*global.DbConfig) error {
	if err := orm.RegisterDriver(DriverName, orm.DRMySQL); err != nil {
		return err
	}
	// 建立连接
	for _, cfg := range cfgs {
		// 建立连接
		str := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", cfg.User, cfg.Password, cfg.Host, cfg.DbName)
		alias := util.GetString(cfg.Alias, cfg.DbName)
		if err := orm.RegisterDataBase(alias, DriverName, str); err != nil {
			return err
		}

		// 获取*sql.DB
		db, err := orm.GetDB(alias)
		if err != nil {
			return err
		}
		// 测试联通
		if err := db.Ping(); err != nil {
			return err
		}
		mysqlPool[cfg.DbName] = &MysqlConn{
			conn:       db,
			driverName: DriverName,
			dbName:     cfg.DbName,
			alias:      alias,
		}
	}
	return nil
}

func GetMysql(dbname string) (orm.Ormer, error) {
	if val, ok := mysqlPool[dbname]; ok {
		o, err := orm.NewOrmWithDB(DriverName, val.alias, val.conn)
		if err != nil {
			return nil, err
		}
		return o, nil
	}
	return nil, fmt.Errorf("数据库(%s)不存在", dbname)
}

func CreateTable(dbname string, model interface{}) error {
	if _, ok := mysqlPool[dbname]; !ok {
		return fmt.Errorf("数据库(%s)不存在", dbname)
	}
	orm.RegisterModel(model)
	// 如果表不存在，则创建表
	//orm.RunSyncdb(dbname, false, true)
	return nil
}

func CreateTableWithPrefix(dbname string, model interface{}, prefix string) error {
	if _, ok := mysqlPool[dbname]; !ok {
		return fmt.Errorf("数据库(%s)不存在", dbname)
	}
	orm.RegisterModelWithPrefix(prefix, model)
	// 如果表不存在，则创建表
	return orm.RunSyncdb(dbname, false, true)
}

func CreateTableWithSuffix(dbname string, model interface{}, suffix string) error {
	if _, ok := mysqlPool[dbname]; !ok {
		return fmt.Errorf("数据库(%s)不存在", dbname)
	}
	orm.RegisterModelWithSuffix(suffix, model)
	// 如果表不存在，则创建表
	return orm.RunSyncdb(dbname, false, true)
}
