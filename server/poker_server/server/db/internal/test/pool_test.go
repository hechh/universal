package test

import (
	"fmt"
	"os"
	"poker_server/common/config"
	"poker_server/common/dao"
	"poker_server/common/dao/domain"
	"poker_server/common/pb"
	"poker_server/common/yaml"
	"poker_server/framework"
	"poker_server/framework/token"
	"poker_server/library/async"
	"poker_server/library/mlog"
	"strings"
	"testing"
	"time"
)

// 实例化一个可以访问的actor的client
func init() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("获取工作目录失败:", err)
		return
	}

	parts := strings.SplitAfter(dir, "poker_server")

	fmt.Println("当前工作目录:", dir)
	// 加载游戏配置
	yamlcfg, node, err := yaml.LoadConfig(parts[0]+"/env/mason/local.yaml", pb.NodeType_NodeTypeClient, int32(1))
	if err != nil {
		panic(fmt.Sprintf("游戏配置加载失败: %v", err))
	}
	nodeCfg := yamlcfg.Client[node.Id]

	// 初始化日志库
	if err := mlog.Init(yamlcfg.Common.Env, nodeCfg.LogLevel, nodeCfg.LogFile); err != nil {
		panic(fmt.Sprintf("日志库初始化失败: %v", err))
	}

	// 初始化mysql
	mlog.Infof("初始化mysql配置")
	if err := dao.InitMysql(yamlcfg.Mysql); err != nil {
		panic(fmt.Sprintf("mysql初始化失败: %v", err))
	}
	async.Init(mlog.Errorf)
	token.Init(yamlcfg.Common.SecretKey)

	// 配置初始化
	if err := config.Init(yamlcfg.Etcd, yamlcfg.Common); err != nil {
		panic(fmt.Sprintf("配置初始化失败: %v", err))
	}
	if err != nil {
		return
	}
	// 初始化框架
	mlog.Infof("启动框架服务: %v", node)
	if err := framework.InitDefault(node, nodeCfg, yamlcfg); err != nil {
		panic(fmt.Sprintf("框架初始化失败: %v", err))
	}
	//select {}
}

func TestRummySettleMatch(t *testing.T) {
	dst := framework.NewDbRouter(uint64(pb.DataType_DataTypeRummySettle), "RummySettleMatchPool", "Insert")
	head := framework.NewHead(dst, pb.RouterType_RouterTypeDataType, 1)
	req := &pb.RummySettleMatchInsertReq{
		Data: []*pb.RummySettleMatch{
			{
				RoomId:     12345,
				PlayerId:   []uint64{23456, 23457},
				SettleCoin: 500,
				CreatedAt:  time.Now().UnixMilli(),
			},
			{
				RoomId:     12345,
				PlayerId:   []uint64{23456, 23457},
				SettleCoin: 500,
				CreatedAt:  time.Now().UnixMilli(),
			},
		},
	}

	err := framework.Send(head, req)
	mlog.Infof("<RummySettleMatchPool Insert recv>: %v <head.src>: %v", err, *head.Src)
}

func TestRummySettle(t *testing.T) {
	// todo send msg to dbsvr
	dst := framework.NewDbRouter(uint64(pb.DataType_DataTypeRummySettle), "RummySettlePool", "Insert")
	head := framework.NewHead(dst, pb.RouterType_RouterTypeDataType, 1)
	req := &pb.RummySettleInsertReq{
		Data: []*pb.RummySettleData{
			{
				RoomId:    12345,
				PlayerId:  23456,
				Coin:      100,
				HandScore: 50,
				CreatedAt: time.Now().UnixMilli(),
			},
			{
				RoomId:    12345,
				PlayerId:  23456,
				Coin:      100,
				HandScore: 50,
				CreatedAt: time.Now().UnixMilli(),
			},
		},
	}
	//rsp := &pb.RummySettleInsertRsp{}
	err := framework.Send(head, req)
	mlog.Infof("<RummySettlePool Insert recv>: %v <head.src>: %v", err, *head.Src)
}

func TestRummySettleSelect(t *testing.T) {
	dst := framework.NewDbRouter(uint64(pb.DataType_DataTypeRummySettle), "RummySettlePool", "Select")
	head := framework.NewHead(dst, pb.RouterType_RouterTypeDataType, 1)
	head.Uid = 23456
	req := &pb.RummySettleSelectReq{
		Page:     1,
		PageSize: 10,
	}
	rsp := &pb.RummySettleSelectRsp{}

	err := framework.Request(head, req, rsp)
	mlog.Infof("rsp : %v err:%v", rsp, err)
}

func TestRummySettleOrm(t *testing.T) {
	session := dao.GetMysql(domain.MYSQL_DB_PLAYER_DATA).NewSession()
	defer session.Close()

	data := []*pb.RummySettleData{
		{
			RoomId:    12345,
			PlayerId:  23456,
			Coin:      100,
			HandScore: 50,
			Groups: []*pb.RummyCardGroup{
				{
					Cards: []uint32{1, 2, 3},
				},
				{
					Cards: []uint32{3, 2, 3},
				},
			},
			CreatedAt: time.Now().UnixMilli(),
		},
		{
			RoomId:    12345,
			PlayerId:  23456,
			Coin:      100,
			HandScore: 50,
			CreatedAt: time.Now().UnixMilli(),
		},
	}
	affectNum, err := session.Insert(data)
	mlog.Infof("affect count: %v, err: %v", affectNum, err)

	var result []*pb.RummySettleData
	count, err := session.Where("player_id = ?", 23456).Desc("created_at").Limit(10, 0).FindAndCount(&result)
	mlog.Infof("err: %v result:%v count:%v", err, result, count)
}
