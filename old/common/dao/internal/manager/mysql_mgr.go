package manager

import (
	"database/sql"
	"fmt"
	"universal/common/global"
	"universal/framework/basic"

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

func InitMysql(cfgs map[uint32]*global.DbConfig) (err error) {
	if len(cfgs) <= 0 {
		return fmt.Errorf("mysql配置为空")
	}
	if err = orm.RegisterDriver(DriverName, orm.DRMySQL); err != nil {
		return
	}
	// 建立连接
	for _, cfg := range cfgs {
		// 建立连接
		alias := basic.GetString(cfg.Alias, cfg.DbName)
		str := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", cfg.User, cfg.Password, cfg.Host, cfg.DbName)
		if err = orm.RegisterDataBase(alias, DriverName, str); err != nil {
			return
		}
		// 获取*sql.DB
		var db *sql.DB
		if db, err = orm.GetDB(alias); err != nil {
			return
		}
		// 测试联通
		if err = db.Ping(); err != nil {
			return
		}
		mysqlPool[cfg.DbName] = &MysqlConn{conn: db, driverName: DriverName, dbName: cfg.DbName, alias: alias}
	}
	return
}

func NewOrmer(dbname string) (orm.Ormer, error) {
	if val, ok := mysqlPool[dbname]; ok {
		o, err := orm.NewOrmWithDB(DriverName, val.alias, val.conn)
		if err != nil {
			return nil, err
		}
		return o, nil
	}
	return nil, fmt.Errorf("数据库(%s)不存在", dbname)
}
