package game

import (
	"corps/base"
	common "corps/common"
	"corps/common/orm/mysql"
	"corps/common/orm/redis"
	"corps/common/report"
	"corps/common/serverCommon"
	"corps/framework/actor"
	"corps/framework/cluster"
	"corps/framework/network"
	"corps/framework/plog"
	"corps/framework/profiler"
	"corps/pb"
	"corps/server/game/cmd"
	"corps/server/game/dipPacket"
	"corps/server/game/playerMgr"
	"fmt"
	"log"
	"os"
	"os/signal"
)

type (
	ServerMgr struct {
		service  *network.ServerSocket
		isInited bool
		markCfg  uint32 //配置表标记
	}

	IServerMgr interface {
		Init() bool
		InitDB() bool
		GetServer() *network.ServerSocket
	}

	Config struct {
		IsDedug            uint32                `yaml:"IsDedug"`
		ServerId           uint32                `yaml:"ServerId"`
		StartPort          int                   `yaml:"StartPort"`
		ServerInfo         map[int]common.Server `yaml:"game"`
		RecordInfo         common.Server         `yaml:"record"`
		PlatformConfig     common.PlatformConfig `yaml:"PlatformConfig"`
		common.Db          `yaml:"mysql"`
		common.Etcd        `yaml:"etcd"`
		common.Nats        `yaml:"nats"`
		common.Stub        `yaml:"stub"`
		common.RedisCommon `yaml:"redis"`
		common.Server
	}
)

var (
	CONF     Config
	SERVER   ServerMgr
	RdID     int
	chSignal chan os.Signal //接收的信号量
	status   base.STAT_TYPE //服务器状态
)

func (s *ServerMgr) Init(id int, ip string, c chan os.Signal) bool {
	status = base.STAT_TYPE_Ing
	chSignal = c
	if s.isInited {
		return true
	}

	//初始配置文件
	base.ReadConf("config.yaml", &CONF)
	var ok any
	CONF.Server, ok = CONF.ServerInfo[id]
	if ok != true {
		log.Fatalf("id( %d) error: %v", id, nil)
	}

	CONF.Server.Port += CONF.StartPort
	CONF.Server.Pprof += CONF.StartPort
	CONF.Server.Gops += CONF.StartPort

	CONF.RecordInfo.Port += CONF.StartPort

	if CONF.IsDedug > 0 {
		plog.SetLevel(uint32(plog.LOG_ALL))
	}

	//修正ip 传参数
	if ip != "" {
		CONF.Server.Ip = ip
	}

	ShowMessage := func() {
		plog.Info("**********************************************************")
		plog.Info("\tGAME ServerId:\t%d", CONF.ServerId)
		plog.Info("\tGAME IsDebug:\t%d", CONF.IsDedug)
		plog.Info("\tGAME ID(%d), IP(LAN):\t%s:%d pprof:%d", id, CONF.Server.Ip, CONF.Server.Port, CONF.Server.Pprof)
		plog.Info("**********************************************************")
	}
	ShowMessage()

	//读平台数据
	serverCommon.SetPlatformConfig(CONF.PlatformConfig)

	//加载配置文件
	if !serverCommon.LoadCfgData(CONF.ServerId) {
		log.Fatalf("id( %d) error: %v", id, nil)
		return false
	}

	//初始化DB
	mysql.DBMGR.Init(CONF.Db)
	redis.Init(CONF.RedisCommon)

	//设置偏移时间
	serverCommon.UpdateOffsetTime()
	//本身game集群管理
	cluster.Init(pb.SERVICE_GAME, CONF.Server.Ip, CONF.Server.Port, CONF.Etcd.Endpoints, CONF.Nats.Endpoints, CONF.Stub)
	cluster.InitPointHandle(actor.MGR.DefaultActorHandle)
	cluster.InitTopicHandle(actor.MGR.DefaultActorHandle)
	cluster.InitQueryHandle(actor.MGR.DefaultActorHandle)
	cluster.InitRpcHandle(actor.MGR.DefaultActorHandle)
	cluster.SetClearExpire(900)
	plog.Info("初始化cluster id:%d", cluster.GetSelfClusterID())

	//初始化report数据上报
	err := report.Init(fmt.Sprintf("%s:%d", CONF.RecordInfo.Ip, CONF.RecordInfo.Port))
	if err != nil {
		panic(err)
	}

	//初始化socket
	s.service = new(network.ServerSocket)
	s.service.Init(CONF.Server.Ip, CONF.Server.Port)
	s.service.Start()

	var packet dipPacket.DipPacket
	packet.Init()
	cmd.Init()
	playerMgr.MGR.Init()

	//开启性能分析工具
	profiler.Init(CONF.Server.Gops, CONF.Server.Pprof)

	plog.Info("game server start success id")
	return false
}

func (s *ServerMgr) GetServer() *network.ServerSocket {
	return s.service
}

func (s *ServerMgr) Stop() {
	if status == base.STAT_TYPE_End {
		return
	}
	status = base.STAT_TYPE_End
	signal.Stop(chSignal)
}

// 停服
func (s *ServerMgr) ShutDown() {
	playerMgr.MGR.ShutDown()
}
