package player

import (
	"context"
	basic "corps/framework/base"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"
)

// 设置玩家数据
func (this *Player) DipSetUserTypeInfo(ctx context.Context, emType pb.PlayerDataType, strData string) {
	head := this.GetRpcHead(ctx)
	if err := this.GetIPlayerFun(emType).SetUserTypeInfo(basic.StringToBytes(strData)); err != nil {
		plog.Error("DipSetUserTypeInfo head: %v, strData: %s", head, strData)
		return
	}
	//返回结果
	head.DestServerType = pb.SERVICE_Dip
	cluster.ReplyMsgTo(head, emType)
}

// 查询玩家缓存数据
func (this *Player) DipGetUserInfo(ctx context.Context, dataType pb.PlayerDataType) {
	pbPlayerData := &pb.PBPlayerData{}
	// 拷贝数据
	this.GetIPlayerFun(dataType).CopyTo(pbPlayerData)
	if dataType == pb.PlayerDataType_Crystal {
		pbPlayerData.Crystal.Effects = nil
	}
	// 应答
	head := this.GetRpcHead(ctx)
	head.DestServerType = pb.SERVICE_Dip
	cluster.ReplyMsgTo(head, pbPlayerData)
}
