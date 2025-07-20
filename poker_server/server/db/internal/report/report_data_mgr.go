package report

import (
	"poker_server/common/mysql"
	"poker_server/common/pb"
	"poker_server/framework/actor"
	"poker_server/library/mlog"
	"poker_server/library/uerror"
	"reflect"

	"github.com/golang/protobuf/proto"
)

func init() {
	mysql.Register(mysql.MYSQL_DB_PLAYER_DATA,
		&pb.TexasPlayerFlowReport{},
		&pb.TexasRoomReport{},
		&pb.TexasGameReport{})
}

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
	mlog.Infof("ReportDataMgr关闭成功")
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
	session := mysql.GetClient(mysql.MYSQL_DB_PLAYER_DATA).NewSession()
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
	session := mysql.GetClient(mysql.MYSQL_DB_PLAYER_DATA).NewSession()
	defer session.Close()

	data := &pb.TexasGameReport{}
	has, err := session.Where("room_id = ? AND round = ?", req.RoomId, req.Round).Get(data)
	if err != nil {
		return uerror.New(1, pb.ErrorCode_READ_FAIELD, "数据库查询失败：%v", err)
	}
	if !has {
		return uerror.New(1, pb.ErrorCode_DB_TEXAS_GAME_DATA_NOT_FOUND, "数据不存在：%v", req)
	}
	rsp.Data = data
	return nil
}

func (d *ReportDataMgr) GetTexasGameReport(head *pb.Head, req *pb.GetTexasGameReportReq, rsp *pb.GetTexasGameReportRsp) error {
	head.FuncName = "GetTexasGameReportPool"
	return d.pool.SendMsg(head, req, rsp)
}

func (d *ReportDataMgr) GetTexasGameReportPool(head *pb.Head, req *pb.GetTexasGameReportReq, rsp *pb.GetTexasGameReportRsp) error {
	session := mysql.GetClient(mysql.MYSQL_DB_PLAYER_DATA).NewSession()
	defer session.Close()

	totalCount, err := session.Where("room_id = ?", req.RoomId).Count(&pb.TexasGameReport{})
	if err != nil {
		return uerror.New(1, pb.ErrorCode_MYSQL_READ_FAILED, "查询数据失败: %v", err)
	}
	rsp.TotalCount = totalCount

	var list []*pb.TexasGameReport
	query := session.Where("room_id = ?", req.RoomId)
	if req.Round > 0 {
		query = query.And("round = ?", req.Round)
	}
	// 设置分页
	if req.PageSize > 0 {
		query = query.Limit(int(req.PageSize), int((req.Page-1)*req.PageSize))
	}
	// 按时间降序排列
	query = query.Desc("begin_time")
	// 执行查询
	if err := query.Find(&list); err != nil {
		return uerror.New(1, pb.ErrorCode_MYSQL_READ_FAILED, "查询数据失败: %v", err)
	}
	// 填充响应
	rsp.List = list
	return nil
}
