package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common/orm/redis"
	"corps/framework/cluster"
	"corps/pb"
	"github.com/golang/protobuf/proto"
)

type (
	PlayerClientFun struct {
		PlayerFun
		mapData           map[string]string
		ClientUpdateOrder uint32
	}
)

func (this *PlayerClientFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
	this.mapData = make(map[string]string)
	this.ClientUpdateOrder = 0
}

func (this *PlayerClientFun) Load(pData []byte) {
	this.BSave = false
	pbData := &pb.PBPlayerClientData{}
	proto.Unmarshal(pData, pbData)
	this.loadData(pbData)
	this.UpdateSave(false)
}
func (this *PlayerClientFun) loadData(pbData *pb.PBPlayerClientData) {
	if pbData == nil {
		pbData = &pb.PBPlayerClientData{}
	}

	this.mapData = make(map[string]string)
	for i := 0; i < len(pbData.ClientDataList); i++ {
		this.mapData[pbData.ClientDataList[i].Type] = pbData.ClientDataList[i].Data
	}
	this.UpdateSave(true)
}
func (this *PlayerClientFun) LoadComplete() {
	this.ClientUpdateOrder = 0
}

func (this *PlayerClientFun) Save(bNow bool) {
	if !this.BSave {
		return
	}

	this.BSave = false

	pbData := &pb.PBPlayerClientData{}
	this.SavePb(pbData)

	//通知db保存玩家数据
	buff, _ := proto.Marshal(pbData)
	cluster.SendToDb(&pb.RpcHead{Id: this.AccountId}, "DbPlayerMgr", "SavePlayerDB", this.PbType, buff, bNow)
}

// 保存
func (this *PlayerClientFun) SavePb(pbData *pb.PBPlayerClientData) {
	if pbData == nil {
		pbData = &pb.PBPlayerClientData{}
	}
	for key, value := range this.mapData {
		pbData.ClientDataList = append(pbData.ClientDataList, &pb.PBClientData{
			Type: key,
			Data: value,
		})
	}
}

func (this *PlayerClientFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	if pbData.Client == nil {
		pbData.Client = &pb.PBPlayerClientData{}
	}
	this.SavePb(pbData.Client)
}

// 设置数据
func (this *PlayerClientFun) SetClientRequest(head *pb.RpcHead, pbRequest *pb.SetClientRequest) {
	this.mapData[pbRequest.ClientData.Type] = pbRequest.ClientData.Data

	this.UpdateSave(true)
	cluster.SendToClient(head, &pb.SetClientResponse{
		PacketHead: &pb.IPacket{},
	}, cfgEnum.ErrorCode_Success)
}
func (this *PlayerClientFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerClientData{}
}

// 设置玩家数据
func (this *PlayerClientFun) SetUserTypeInfo(pbData proto.Message) bool {
	if pbData == nil {
		return false
	}

	pbSystem := pbData.(*pb.PBPlayerClientData)
	if pbData == nil {
		return false
	}

	this.loadData(pbSystem)
	return true
}

// 定时读取客户端数据
func (this *PlayerClientFun) TimeReadClientData() {
	//先读取个人邮件
	redisCommon := redis.GetCommonRedis()
	if redisCommon == nil {
		return
	}

	//遍历邮件
	uOrder := redisCommon.GetUint32(base.ERK_ClientUpdateOrder)
	if this.ClientUpdateOrder == uOrder {
		return
	}

	this.ClientUpdateOrder = uOrder
	pbNotify := &pb.ClientJsonNotify{PacketHead: &pb.IPacket{Id: this.AccountId}}
	mapData := redisCommon.HGetAll(base.ERK_ClientUpdate)
	for name, data := range mapData {
		pbNotify.JsonList = append(pbNotify.JsonList, &pb.PBJsonInfo{JsonName: name, JsonData: data})
	}

	cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, pbNotify, cfgEnum.ErrorCode_Success)
}
