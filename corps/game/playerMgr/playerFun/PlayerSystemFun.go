package playerFun

import (
	"corps/base"
	"corps/pb"

	"github.com/golang/protobuf/proto"
)

type (
	PlayerSystemFun struct {
		PlayerFun
		pbData pb.PBPlayerSystem //系统
	}
)

func (this *PlayerSystemFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
}

func (this *PlayerSystemFun) Load(pData []byte) {
	this.BSave = false
	proto.Unmarshal(pData, &this.pbData)
}

// 保存
func (this *PlayerSystemFun) SavePb(pbData *pb.PBPlayerSystem) {
	base.DeepCopy(&this.pbData, pbData)
}

func (this *PlayerSystemFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	base.DeepCopy(&this.pbData, pbData.System)
}
