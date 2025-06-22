package rummy

import (
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/actor"
	"poker_server/library/timer"
	"poker_server/library/uerror"
	"reflect"
	"time"
)

type BuilderRummyGenerator struct {
	actor.Actor
	datas    map[pb.GeneratorType]uint32 // 房间ID生成器数据
	isChange bool                        // 是否有变更
	loadFlag uint64                      // 是否加载完成
}

func NewBuilderRummyGenerator() *BuilderRummyGenerator {
	ret := &BuilderRummyGenerator{}
	ret.Actor.Register(ret)
	ret.Actor.ParseFunc(reflect.TypeOf(ret))
	ret.Actor.SetId(uint64(pb.GeneratorType_GeneratorTypeRummy))
	ret.Actor.Start()
	actor.Register(ret)
	return ret
}

func (g *BuilderRummyGenerator) Load() error {
	g.loadFlag = 1
	return timer.Register(&g.loadFlag, func() { g.SendMsg(&pb.Head{FuncName: "LoadRequest"}) }, 5*time.Second, -1)
}

// 发送数据加载请求
func (g *BuilderRummyGenerator) LoadRequest() error {
	if g.loadFlag <= 0 {
		return nil
	}
	dst := framework.NewDbRouter(uint64(pb.DataType_DataTypeGenerator), "DbGeneratorMgr", "Query")
	head := framework.NewHead(dst, pb.RouterType_RouterTypeGeneratorType, uint64(pb.GeneratorType_GeneratorTypeRummy), "BuilderRummyGenerator", "LoadData")
	return framework.Send(head)
}

// 加载数据
func (g *BuilderRummyGenerator) LoadData(head *pb.Head, rsp *pb.GetGeneratorDataRsp) error {
	if rsp.Head != nil {
		return uerror.ToError(rsp.Head)
	}
	if g.loadFlag > 0 {
		// 加载数据
		g.loadFlag = 0
		g.datas = make(map[pb.GeneratorType]uint32)
		for _, item := range rsp.List {
			g.datas[item.Id] = item.Incr
		}
		return g.RegisterTimer(&pb.Head{FuncName: "OnTick"}, 5*time.Second, -1)
	}
	return nil
}

/*
// 加载数据
func (g *BuilderRummyGenerator) Load() error {
	// 请求数据
	head := framework.NewDbHead("DbGeneratorMgr", "Query", uint64(pb.DataType_DataTypeGenerator), g)
	req := &pb.GetGeneratorDataReq{DataType: pb.DataType_DataTypeGenerator}
	rsp := &pb.GetGeneratorDataRsp{}
	if err := framework.Request(head, req, rsp); err != nil {
		return err
	}
	if rsp.Head != nil {
		return uerror.ToError(rsp.Head)
	}

	// 加载数据
	g.datas = make(map[pb.GeneratorType]uint32)
	for _, item := range rsp.List {
		g.datas[item.Id] = item.Incr
	}

	return g.RegisterTimer(&pb.Head{
		SendType:  pb.SendType_POINT,
		ActorName: g.GetActorName(),
		FuncName:  "OnTick",
	}, 5*time.Second, -1)
}
*/

func (g *BuilderRummyGenerator) OnTick() {
	if !g.isChange {
		return
	}
	g.SendMsg(&pb.Head{FuncName: "Save"})
}

// 保存数据
func (g *BuilderRummyGenerator) Save() error {
	notify := &pb.UpdateGeneratorDataNotify{DataType: pb.DataType_DataTypeGenerator}
	for id, incr := range g.datas {
		notify.List = append(notify.List, &pb.GeneratorInfo{Id: id, Incr: incr})
	}

	dst := framework.NewDbRouter(uint64(pb.DataType_DataTypeGenerator), "DbGeneratorMgr", "Update")
	head := framework.NewHead(dst, pb.RouterType_RouterTypeGeneratorType, uint64(pb.GeneratorType_GeneratorTypeRummy))

	if err := framework.Send(head, notify); err != nil {
		return err
	}
	g.isChange = false
	return nil
}

// 生成房间ID请求（同步+异步）
func (g *BuilderRummyGenerator) GenRoomIdReq(head *pb.Head, req *pb.GenRoomIdReq, rsp *pb.GenRoomIdRsp) error {
	if req.Count <= 1 {
		req.Count = 1
	}

	// 随便使用简单方法生成唯一ID
	for i := int32(0); i < req.Count; i++ {
		switch req.GeneratorType {
		case pb.GeneratorType_GeneratorTypeTexas, pb.GeneratorType_GeneratorTypeRummy, pb.GeneratorType_GeneratorTypeSng:
			g.datas[req.GeneratorType]++
			incr := g.datas[req.GeneratorType]
			id := uint64(req.MatchType&0xFF)<<40 | uint64(req.GameType&0xFF)<<32 | uint64(req.CoinType&0xFF)<<24 | uint64(incr&0xFFFFFF)
			rsp.RoomIdList = append(rsp.RoomIdList, id)
			g.isChange = true
		default:
			return uerror.New(1, pb.ErrorCode_TYPE_NOT_SUPPORTED, "不支持的生成器类型: %s", req.GeneratorType.String())
		}
	}
	return nil
}
