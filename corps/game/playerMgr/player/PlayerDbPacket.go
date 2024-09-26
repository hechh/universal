package player

import (
	"context"
	"corps/base"
	"corps/framework/plog"
	"corps/pb"

	"github.com/golang/protobuf/proto"
)

// 加载玩家结束
func (this *Player) LoadPlayerDBFinish(ctx context.Context) {
	head := this.GetRpcHead(ctx)
	this.isInGame = true
	plog.Trace("玩家:%d 数据加载成功 LoadPlayerDBFinish", this.GetId())

	//需要通知各个系统加载数据库加载完成
	for typ, fun := range this.mapPlayerData {
		fun.LoadPlayerDBFinish()
		plog.Trace("uid: %d, type: %s LoadPlayerDBFinish", this.GetId(), typ.String())
	}
	//判断是否是新系统
	listNewPlayerType := this.getPlayerBaseFun().GetNewPlayerTypeList()
	for _, fun := range this.listPlayerData {
		if base.ArrayContainsValue(listNewPlayerType, uint32(fun.GetPbType())) {
			continue
		}
		fun.NewPlayer()
		plog.Trace("uid: %d, type: %s NewPlayer", this.GetId(), fun.GetPbType().String())
		this.getPlayerBaseFun().AddNewPlayerTypeList(uint32(fun.GetPbType()))
	}

	//发送给网关
	this.loginSuccess(head)

	this.RegisterTimers()
}

// 加载玩家结束
func (p *Player) LoadPlayerDBType(ctx context.Context, dataType pb.PlayerDataType, pData []byte) {
	head := p.GetRpcHead(ctx)
	plog.Trace("玩家数据加载成功 LoadPlayerDBType id:%d type:%d", head.Id, dataType)

	if dataType == pb.PlayerDataType_System {
		pbSystem := new(pb.PBPlayerSystem)
		proto.Unmarshal(pData, pbSystem)
		for i := pb.PlayerDataType_SystemCommon; i < pb.PlayerDataType_SystemMax; i++ {
			fun, ok := p.mapPlayerData[i]
			if ok {
				fun.LoadSystem(pbSystem)
			}
		}
		plog.Trace("LoadPlayerDBType uid: %d, PBPlayerSystem: %v", p.GetId(), pbSystem)
	} else {
		//如果是新玩家，需要初始化
		fun, ok := p.mapPlayerData[dataType]
		if !ok {
			plog.Info("玩家数据加载失败 LoadPlayerDBType id:%d type:%d", head.Id, dataType)
			return
		}
		fun.Load(pData)
	}
}
