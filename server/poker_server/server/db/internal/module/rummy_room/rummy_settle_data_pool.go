package rummy_room

import (
	"poker_server/common/dao"
	"poker_server/common/dao/domain"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/actor"
	"poker_server/library/mlog"
	"poker_server/library/uerror"
	"reflect"
	"time"

	"github.com/golang/protobuf/proto"
)

func init() {
	// todo 联合索引初始化
	dao.RegisterMysqlTable(domain.MYSQL_DB_PLAYER_DATA, &pb.RummySettleData{})
}

const (
	POOLS = 8 //协程数
)

type RummySettlePool struct {
	actor.Actor
	pool *actor.ActorPool
	data []*pb.RummySettleData
}

func NewRummySettlePool() *RummySettlePool {
	ret := &RummySettlePool{
		pool: new(actor.ActorPool),
		data: make([]*pb.RummySettleData, 0, 128),
	}
	ret.pool.Register(ret, POOLS)
	ret.pool.ParseFunc(reflect.TypeOf(ret))
	ret.pool.Start()

	ret.Register(ret)
	ret.ParseFunc(reflect.TypeOf(ret))
	ret.SetId(uint64(pb.DataType_DataTypeRummySettle))
	ret.Start()

	return ret
}

func (d *RummySettlePool) Init() error {
	actor.Register(d)
	// 启动定时器
	head := &pb.Head{SendType: pb.SendType_POINT, ActorName: "RummySettlePool", FuncName: "OnTick"}
	err := d.RegisterTimer(head, 5000*time.Millisecond, -1)
	if err != nil {
		mlog.Infof("register timer err: %v", err)
	}
	return nil
}

// OnTick 定时清空缓冲队列
func (d *RummySettlePool) OnTick() {
	if len(d.data) > 0 {
		head := &pb.Head{FuncName: "InsertPool"}
		framework.StopAutoSendToClient(head) //禁止自动回包
		tmp := make([]*pb.RummySettleData, 0, len(d.data))
		for i := range d.data {
			tmp = append(tmp, proto.Clone(d.data[i]).(*pb.RummySettleData))
		}
		d.pool.SendMsg(head, tmp)
		d.data = d.data[:0]
	}
}

// Insert 异步存db
func (d *RummySettlePool) Insert(head *pb.Head, req *pb.RummySettleInsertReq, rsp *pb.RummySettleInsertRsp) error {
	mlog.Infof("<RummySettlePool Insert recv>: %v", head)
	framework.StopAutoSendToClient(head) //InsertPool 不满足接口范式 禁止自动回包
	if len(d.data)+len(req.Data) >= 128 {
		head.FuncName = "InsertPool"

		tmp := make([]*pb.RummySettleData, 0, len(d.data)+len(req.Data))
		for i := range d.data {
			tmp = append(tmp, proto.Clone(d.data[i]).(*pb.RummySettleData))
		}
		for i := range req.Data {
			tmp = append(tmp, proto.Clone(req.Data[i]).(*pb.RummySettleData))
		}

		ret := d.pool.SendMsg(head, tmp)
		d.data = d.data[:0]
		return ret
	}
	d.data = append(d.data, req.Data...)
	rsp.Head = &pb.RspHead{}
	return nil
}

// InsertPool 每个子线程执行
func (d *RummySettlePool) InsertPool(head *pb.Head, data []*pb.RummySettleData) error {
	session := dao.GetMysql(domain.MYSQL_DB_PLAYER_DATA).NewSession()
	defer session.Close()
	affectNum, err := session.Insert(data)
	if err != nil {
		mlog.Infof("affect count: %v, err: %v", affectNum, err)
	}
	return nil
}

func (d *RummySettlePool) SelectPool(head *pb.Head, req *pb.RummySettleSelectReq, rsp *pb.RummySettleSelectRsp) error {
	session := dao.GetMysql(domain.MYSQL_DB_PLAYER_DATA).NewSession()
	defer session.Close()

	var data []*pb.RummySettleData
	total, err := session.Where("player_id = ?", head.Uid).OrderBy("created_at desc").Limit(int(req.PageSize), int((req.Page-1)*req.PageSize)).FindAndCount(&data)

	if err != nil {
		mlog.Errorf("SelectPool Find err: %v", err)
		return err
	}

	rsp.Count = total
	rsp.Data = data
	return nil
}

// Select 未优化
func (d *RummySettlePool) Select(head *pb.Head, req *pb.RummySettleSelectReq, rsp *pb.RummySettleSelectRsp) error {
	head.FuncName = "SelectPool"

	if req.Page <= 0 || req.PageSize <= 0 {
		return uerror.NEW(pb.ErrorCode_RUMMY_PLAYER_DATA_INVAILD, head, "Page or PageSize numbering starts from 1")
	}

	return d.pool.SendMsg(head, req, rsp)
}

func (d *RummySettlePool) Close() {
	d.Actor.Stop()
	d.pool.Stop()
}
