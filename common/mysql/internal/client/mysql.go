package client

import (
	"fmt"
	"sync/atomic"
	"universal/common/yaml"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var (
	url = "%s:%s@tcp(%s)/%s?timeout=3s&parseTime=true&loc=Local&charset=utf8"
)

type MysqlClient struct {
	engine  *xorm.EngineGroup // 数据库引擎组
	dsn     []string          // 数据库连接字符串
	dbname  string            // 数据库名称
	isAlive int32             // 连接是否正常
}

func NewMysqlClient(cfg *yaml.DbConfig, tables ...interface{}) *MysqlClient {
	dsn := []string{}
	// 主节点
	dsn = append(dsn, fmt.Sprintf(url, cfg.User, cfg.Password, cfg.Host, cfg.DbName))
	// 从节点配置
	for _, scfg := range cfg.Slave {
		dsn = append(dsn, fmt.Sprintf(url, scfg.User, scfg.Password, scfg.Host, scfg.DbName))
	}
	return &MysqlClient{dsn: dsn, dbname: cfg.DbName}
}

func (o *MysqlClient) Connect(tables ...interface{}) error {
	eng, err := xorm.NewEngineGroup("mysql", o.dsn)
	if err != nil {
		return err
	}
	eng.SetMaxIdleConns(10)
	eng.SetMaxOpenConns(200)
	if len(tables) > 0 {
		eng.Sync2(tables...)
	}

	// 查看连接是否联通
	if err := eng.Ping(); err != nil {
		eng.Close()
		return err
	}
	if o.engine != nil {
		o.engine.Close()
	}
	o.engine = eng
	atomic.StoreInt32(&o.isAlive, 1)
	return nil
}

// 检测连接是否联通
func (o *MysqlClient) Ping() error {
	return o.engine.Ping()
}

func (o *MysqlClient) IsAlive() bool {
	return atomic.LoadInt32(&o.isAlive) > 0
}

// 创建session
func (o *MysqlClient) NewSession() *xorm.Session {
	return o.engine.NewSession()
}

func (o *MysqlClient) GetEngine() *xorm.Engine {
	return o.engine.Engine
}
