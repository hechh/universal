package rummy_room

import (
	"poker_server/common/mysql"
	"poker_server/common/pb"
	"poker_server/framework/actor"
	"poker_server/library/mlog"
	"poker_server/library/uerror"
	"reflect"
)

func init() {
	mysql.Register(mysql.MYSQL_DB_PLAYER_DATA, &pb.RummyExtSettleMatch{}, &pb.RummyExtSettleUser{})
}

func (d *RummyExtSettleMatchPool) Init() error {
	return nil
}

// RummyExtSettleMatchPool 比赛抽水数据
type RummyExtSettleMatchPool struct {
	actor.Actor
	pool *actor.ActorPool
}

func NewRummyExtSettleMatchPool() *RummyExtSettleMatchPool {
	ret := &RummyExtSettleMatchPool{
		pool: new(actor.ActorPool),
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

func (d *RummyExtSettleMatchPool) Close() {
	d.Actor.Stop()
	d.pool.Stop()
	mlog.Infof("RummyExtSettleMatchPool关闭成功")
}

// Insert 异步存db
func (d *RummyExtSettleMatchPool) Insert(head *pb.Head, req *pb.RummyExtSettleMatchInsertReq) error {
	head.FuncName = "InsertPool"
	return d.pool.SendMsg(head, req)
}

func (d *RummyExtSettleMatchPool) InsertPool(head *pb.Head, req *pb.RummyExtSettleMatchInsertReq) error {
	session := mysql.GetClient(mysql.MYSQL_DB_PLAYER_DATA).NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		return err
	}
	match_id, err := session.Insert(req.Data)
	if err != nil {
		err1 := session.Rollback()
		mlog.Error(head, "数据落地失败 RummyExtSettleMatch：data:%v, error:%v trans err: %v", req, err, err1)
		return err
	}

	for i := range req.Bills {
		req.Bills[i].MatchId = uint64(match_id)
	}
	_, err = session.Insert(req.Bills)
	if err != nil {
		err1 := session.Rollback()
		mlog.Error(head, "数据落地失败 RummyExtSettleUser：data:%v, error:%v trans err: %v", req, err, err1)
		return err
	}
	return session.Commit()
}

func (d *RummyExtSettleMatchPool) Select(head *pb.Head, req *pb.RummyExtSettleSelectReq, rsp *pb.RummyExtSettleSelectRsp) error {
	head.FuncName = "SelectPool"

	if req.Page <= 0 || req.PageSize <= 0 {
		return uerror.New(1, pb.ErrorCode_RUMMY_PLAYER_DATA_INVAILD, "Page or PageSize numbering starts from 1")
	}

	return d.pool.SendMsg(head, req, rsp)
}

func (d *RummyExtSettleMatchPool) SelectPool(head *pb.Head, req *pb.RummyExtSettleSelectReq, rsp *pb.RummyExtSettleSelectRsp) error {
	session := mysql.GetClient(mysql.MYSQL_DB_PLAYER_DATA).NewSession()
	defer session.Close()

	var data []*pb.RummyExtSettleMatch
	total, err := session.Table("rummy_ext_settle_match").Join("INNER", "rummy_ext_settle_user", "rummy_ext_settle_user.match_id = rummy_ext_settle_match.id").Where("rummy_ext_settle_user.player_id = ?", head.Uid).OrderBy("rummy_ext_settle_match.created_at desc").Limit(int(req.PageSize), int((req.Page-1)*req.PageSize)).FindAndCount(&data)

	if err != nil {
		mlog.Errorf("SelectPool Find err: %v", err)
		return err
	}

	rsp.Count = total
	rsp.Data = data
	return nil
}
