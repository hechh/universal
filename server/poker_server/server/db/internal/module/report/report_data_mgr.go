package report

import (
	"fmt"
	"poker_server/common/dao"
	"poker_server/common/dao/domain"
	"poker_server/common/pb"
	"poker_server/framework/actor"
	"poker_server/library/mlog"
	"poker_server/library/uerror"
	"reflect"

	"github.com/golang/protobuf/proto"
)

type ReportDataMgr struct {
	actor.Actor
	pool *actor.ActorPool
}

func NewReportDataMgr() *ReportDataMgr {
	ret := &ReportDataMgr{pool: new(actor.ActorPool)}
	ret.pool.Register(ret, 50)
	ret.pool.ParseFunc(reflect.TypeOf(ret))
	ret.pool.SetId(uint64(pb.DataType_DataTypeReport))
	ret.pool.Start()

	ret.Actor.Register(ret)
	ret.Actor.ParseFunc(reflect.TypeOf(ret))
	ret.SetId(uint64(pb.DataType_DataTypeReport))
	ret.Start()
	actor.Register(ret)
	return ret
}

func (d *ReportDataMgr) Init() error {
	return nil
}

func (d *ReportDataMgr) Close() {
	d.Actor.Stop()
	d.pool.Stop()
}

func (d *ReportDataMgr) TexasGameReport(head *pb.Head, data *pb.TexasGameReport) error {
	head.FuncName = "Insert"
	return d.pool.SendMsg(head, data)
}

func (d *ReportDataMgr) TexasPlayerFlowReport(head *pb.Head, data *pb.TexasPlayerFlowReport) error {
	head.FuncName = "Insert"
	return d.pool.SendMsg(head, data)
}

func (d *ReportDataMgr) TexasRoomReport(head *pb.Head, data *pb.TexasRoomReport) error {
	head.FuncName = "Insert"
	return d.pool.SendMsg(head, data)
}

func (d *ReportDataMgr) Insert(head *pb.Head, data proto.Message) error {
	session := dao.GetMysql(domain.MYSQL_DB_PLAYER_DATA).NewSession()
	defer session.Close()

	_, err := session.Insert(data)
	if err != nil {
		mlog.Error(head, "数据落地失败：data:%v, error:%v", data, err)
	}
	return err
}

func (d *ReportDataMgr) GameReportReq(head *pb.Head, req *pb.TexasGameReportReq, rsp *pb.TexasGameReportRsp) error {
	head.FuncName = "GameReportPool"
	return d.pool.SendMsg(head, req, rsp)
}

func (d *ReportDataMgr) GameReportPool(head *pb.Head, req *pb.TexasGameReportReq, rsp *pb.TexasGameReportRsp) error {
	session := dao.GetMysql(domain.MYSQL_DB_PLAYER_DATA).NewSession()
	defer session.Close()

	data := &pb.TexasGameReport{}
	has, err := session.Where("room_id = ? AND round = ?", req.RoomId, req.Round).Get(data)
	if err != nil {
		return uerror.NEW(pb.ErrorCode_READ_FAIELD, head, "数据库查询失败：%v", err)
	}
	if !has {
		return uerror.NEW(pb.ErrorCode_DB_TEXAS_GAME_DATA_NOT_FOUND, head, "数据不存在：%v", req)
	}
	rsp.Data = data
	return nil
}

func (d *ReportDataMgr) GetTexasGameReport(head *pb.Head, req *pb.GetTexasGameReportReq, rsp *pb.GetTexasGameReportRsp) error {
	head.FuncName = "GetTexasGameReportPool"
	return d.pool.SendMsg(head, req, rsp)
}

func (d *ReportDataMgr) GetTexasGameReportPool(head *pb.Head, req *pb.GetTexasGameReportReq, rsp *pb.GetTexasGameReportRsp) error {
	session := dao.GetMysql(domain.MYSQL_DB_PLAYER_DATA).NewSession()
	defer session.Close()

	var list []*pb.TexasGameReport
	query := session.Where("room_id = ?", req.RoomId)
	if req.Round > 0 {
		query = query.Where("round = ?", req.Round)
	}
	// 设置分页
	if req.PageSize > 0 {
		query = query.Limit(int(req.PageSize), int((req.Page-1)*req.PageSize))
	}
	// 按时间降序排列
	query = query.Desc("begin_time")
	// 执行查询
	err := query.Find(&list)
	if err != nil {
		return fmt.Errorf("GetTexasGameReportPool query failed: %v", err)
	}
	// 填充响应
	rsp.List = list
	return nil
}
