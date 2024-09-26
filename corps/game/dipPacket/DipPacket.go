package dipPacket

import (
	"context"
	"corps/base"
	"corps/common/serverCommon"
	"corps/framework/actor"
	"corps/framework/cluster"
	"corps/framework/plog"
	pb "corps/pb"
	"corps/server/game/playerMgr"
)

type (
	DipPacket struct {
		actor.Actor
	}

	IDipPacket interface {
		actor.IActor
	}
)

func (e *DipPacket) Init() {
	e.Actor.Init()
	actor.MGR.RegisterActor(e)
	e.Actor.Start()
}

// 热更配置文件
func (e *DipPacket) ReloadCfg(ctx context.Context) {
	serverCommon.ReLoadCfgData()
}

// 获取玩家数据
func (e *DipPacket) DipGetUserInfo(ctx context.Context, emType pb.PlayerDataType) {
	head := e.GetRpcHead(ctx)
	plog.Trace("%s.%s head: %v", head.ActorName, head.FuncName, head)

	if pPlayer := playerMgr.MGR.GetIActor(head.Id); pPlayer != nil {
		pPlayer.SendMsg(head, emType)
	} else {
		//发送到dbserver去取
		head.DestServerType = pb.SERVICE_DB
		head.ClusterId = 0
		head.ActorName = "DbPlayerMgr"
		head.FuncName = "DipGetUserInfo"
		cluster.RequestMsgTo(head, emType)
	}
}

// 获取玩家数据
func (e *DipPacket) DipSetUserTypeInfo(ctx context.Context, emType pb.PlayerDataType, strData string) {
	head := e.GetRpcHead(ctx)
	plog.Trace("%s.%s head: %v", head.ActorName, head.FuncName, head)

	//通知修改内存
	pPlayer := playerMgr.MGR.GetIActor(head.Id)
	if pPlayer != nil {
		head.FuncName = "DipSetUserTypeInfo"
		pPlayer.SendMsg(head, emType, strData)
	} else {
		//发送到dbserver去取
		head.DestServerType = pb.SERVICE_DB
		head.ActorName = "DbPlayerMgr"
		head.FuncName = "DipSetUserTypeInfo"
		head.ClusterId = 0
		cluster.RequestMsgTo(head, emType, strData)
	}
}

// 获取玩家邮件
func (e *DipPacket) DipSendPlayerMail(ctx context.Context, pbMail *pb.PBMail) {
	head := e.GetRpcHead(ctx)
	plog.Trace("%s.%s head: %v", head.ActorName, head.FuncName, head)

	serverCommon.SendPlayerMail(head.Id, pbMail)
}

// 删除玩家邮件
func (e *DipPacket) DipDelPlayerMail(ctx context.Context, mailId uint32, bReward bool) {
	head := e.GetRpcHead(ctx)
	head.FuncName = "DipDelPlayerMail"
	plog.Trace("%s.%s head: %v", head.ActorName, head.FuncName, head)

	//通知修改内存
	if pPlayer := playerMgr.MGR.GetIActor(head.Id); pPlayer != nil {
		pPlayer.SendMsg(head, mailId, bReward)
	} else {
		//去数据读取
		serverCommon.DipDelPlayerMail(head.Id, mailId, bReward)
	}
}

// 给玩家添加物品
func (e *DipPacket) DipAddPlayerItem(ctx context.Context, pbItem *pb.PBAddItemData) {
	head := e.GetRpcHead(ctx)
	plog.Trace("%s.%s head: %v", head.ActorName, head.FuncName, head)

	serverCommon.SendPlayerItem(head.Id, pbItem)
}

// 获取玩家数据
func (e *DipPacket) DipGetCfgDataInfo(ctx context.Context) {
	head := e.GetRpcHead(ctx)
	plog.Trace("%s.%s head: %v", head.ActorName, head.FuncName, head)

	mapInfo := serverCommon.GetCfgDataTime()
	head.DestServerType = pb.SERVICE_Dip
	cluster.ReplyMsgTo(head, mapInfo)
}

// 获取json数据
func (e *DipPacket) DipGetCfgDataJsonInfo(ctx context.Context, strJsonName string) {
	head := e.GetRpcHead(ctx)
	plog.Trace("%s.%s head: %v", head.ActorName, head.FuncName, head)

	head.DestServerType = pb.SERVICE_Dip
	cluster.ReplyMsgTo(head, serverCommon.GetFileContent(strJsonName))
}

// 获取json数据
func (e *DipPacket) DipSetCfgDataJsonInfo(ctx context.Context, strJsonName string, strJsonContent string) {
	head := e.GetRpcHead(ctx)
	head.DestServerType = pb.SERVICE_Dip
	cluster.ReplyMsgTo(head, serverCommon.SetFileContent(strJsonName, strJsonContent))
}

// 玩家拷贝
func (e *DipPacket) DipUpdatePlayerCopy(ctx context.Context, strPlayerName string) {
	head := e.GetRpcHead(ctx)
	plog.Trace("%s.%s head: %v", head.ActorName, head.FuncName, head)

	//去数据读取
	pbPlayer := new(pb.PBPlayerData)
	serverCommon.DipUpdatePlayerCopy(head.Id, pbPlayer, strPlayerName)

	head.DestServerType = pb.SERVICE_Dip
	cluster.ReplyMsgTo(head, pbPlayer)
}

// 更新充值数据
func (e *DipPacket) DipUpdatePlayerCharge(ctx context.Context, ProductId uint32) {
	head := e.GetRpcHead(ctx)
	head.FuncName = "DipUpdatePlayerCharge"
	plog.Trace("%s.%s head: %v", head.ActorName, head.FuncName, head)

	//通知修改内存
	pPlayer := playerMgr.MGR.GetIActor(head.Id)
	if pPlayer == nil {
		return
	}
	pPlayer.SendMsg(head, ProductId)
}

// 设置偏移时间
func (e *DipPacket) DipSetOffsetTime(ctx context.Context, uOffsetTime int32) {
	head := e.GetRpcHead(ctx)
	plog.Trace("%s.%s head: %v", head.ActorName, head.FuncName, head)
	base.SetOffsetTime(uOffsetTime)
}

// 调整日志等级
func (e *DipPacket) DipSetLogLevel(ctx context.Context, level uint32) {
	plog.SetLevel(level)
}
