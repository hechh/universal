package rummy_room

import (
	"poker_server/common/mysql"
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
	mysql.Register(mysql.MYSQL_DB_PLAYER_DATA, &pb.RummySettleMatch{})
}

// RummySettleMatch 比赛抽水数据
type RummySettleMatchPool struct {
	actor.Actor
	pool *actor.ActorPool
	data []*pb.RummySettleMatch
}

func NewRummySettleMatchPool() *RummySettleMatchPool {
	ret := &RummySettleMatchPool{
		pool: new(actor.ActorPool),
		data: make([]*pb.RummySettleMatch, 0, 128),
	}
	ret.pool.Register(ret, POOLS)
	ret.pool.ParseFunc(reflect.TypeOf(ret))
	ret.pool.Start()

	ret.Register(ret)
	ret.ParseFunc(reflect.TypeOf(ret))
	ret.SetId(uint64(pb.DataType_DataTypeRummyMatch))
	ret.Start()
	return ret
}

func (d *RummySettleMatchPool) Init() error {
	actor.Register(d)
	// 启动定时器
	head := &pb.Head{SendType: pb.SendType_POINT, ActorName: "RummySettleMatchPool", FuncName: "OnTick"}
	err := d.RegisterTimer(head, 5000*time.Millisecond, -1)
	if err != nil {
		mlog.Info(head, "register timer err: %v", err)
	}
	return nil
}

// OnTick 定时清空缓冲队列
func (d *RummySettleMatchPool) OnTick() {
	if len(d.data) > 0 {
		head := &pb.Head{FuncName: "InsertPool"}
		framework.StopAutoSendToClient(head) //禁止自动回包
		tmp := make([]*pb.RummySettleMatch, 0, len(d.data))
		for i := range d.data {
			tmp = append(tmp, proto.Clone(d.data[i]).(*pb.RummySettleMatch))
		}
		d.pool.SendMsg(head, tmp)
		d.data = d.data[:0]
	}
}

// Insert 异步存db
func (d *RummySettleMatchPool) Insert(head *pb.Head, req *pb.RummySettleMatchInsertReq, rsp *pb.RummySettleMatchInsertRsp) error {
	if len(d.data)+len(req.Data) >= 0 {
		head.FuncName = "InsertPool"
		framework.StopAutoSendToClient(head) //InsertPool 不满足接口范式 禁止自动回包

		tmp := make([]*pb.RummySettleMatch, 0, len(d.data)+len(req.Data))
		for i := range d.data {
			tmp = append(tmp, proto.Clone(d.data[i]).(*pb.RummySettleMatch))
		}
		for i := range req.Data {
			tmp = append(tmp, proto.Clone(req.Data[i]).(*pb.RummySettleMatch))
		}

		ret := d.pool.SendMsg(head, tmp)
		d.data = d.data[:0]
		return ret
	}
	d.data = append(d.data, req.Data...)
	return nil
}

// InsertPool 每个子线程执行 todo type error
func (d *RummySettleMatchPool) InsertPool(head *pb.Head, data []*pb.RummySettleMatch) error {
	mlog.Debug(head, "InsertPool data:%v", data)
	session := mysql.GetClient(mysql.MYSQL_DB_PLAYER_DATA).NewSession()
	defer session.Close()
	mysql.GetClient(mysql.MYSQL_DB_PLAYER_DATA).GetEngine().ShowSQL(true)
	affectNum, err := session.Insert(data)
	if err != nil {
		mlog.Info(head, "affect count: %v, err: %v", affectNum, err)
	}
	return nil
}

func (d *RummySettleMatchPool) SelectPool(head *pb.Head, req *pb.RummyMatchSelectReq, rsp *pb.RummyMatchSelectRsp) error {
	session := mysql.GetClient(mysql.MYSQL_DB_PLAYER_DATA).NewSession()
	defer session.Close()

	var data []*pb.RummySettleMatch
	total, err := session.Where("room_id = ?", req.RoomId).OrderBy("created_at desc").Limit(int(req.PageSize), int((req.Page-1)*req.PageSize)).FindAndCount(&data)

	if err != nil {
		mlog.Errorf("SelectPool Find err: %v", err)
		return err
	}

	rsp.Count = total
	rsp.Data = data
	return nil
}

// Select 未优化
func (d *RummySettleMatchPool) Select(head *pb.Head, req *pb.RummyMatchSelectReq, rsp *pb.RummyMatchSelectRsp) error {
	head.FuncName = "SelectPool"

	if req.Page <= 0 || req.PageSize <= 0 {
		return uerror.New(1, pb.ErrorCode_RUMMY_PLAYER_DATA_INVAILD, "Page or PageSize numbering starts from 1")
	}

	return d.pool.SendMsg(head, req, rsp)
}

func (d *RummySettleMatchPool) Save() error {
	return nil
}
func (d *RummySettleMatchPool) Load() error {
	return nil
}
func (d *RummySettleMatchPool) Delete() error {
	return nil
}

func (d *RummySettleMatchPool) Close() {
	d.Actor.Stop()
	d.pool.Stop()
}
