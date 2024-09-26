package playerFun

import (
	"corps/base"
	"corps/framework/plog"
	"corps/pb"

	"github.com/golang/protobuf/proto"
)

type (
	PlayerSystemRepairFun struct {
		PlayerFun
		pbData        *pb.PBPlayerSystemRepairData
		listRepairFun []*RepairFun
		mapRepairFun  map[uint32]*RepairFun
	}

	RepairFun struct {
		versionID uint32
		fun       func()
	}
)

func (this *PlayerSystemRepairFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
	this.pbData = &pb.PBPlayerSystemRepairData{}
	this.mapRepairFun = make(map[uint32]*RepairFun)
	this.registerRepairFun(1, this.repairTest)
	this.registerRepairFun(2, this.repairShopRed)
	this.registerRepairFun(3, this.repairMainTask0828)
}

// 从数据库中加载
func (this *PlayerSystemRepairFun) registerRepairFun(id uint32, fun func()) {
	this.mapRepairFun[id] = &RepairFun{versionID: id, fun: fun}
	this.listRepairFun = append(this.listRepairFun, this.mapRepairFun[id])
}

// 从数据库中加载
func (this *PlayerSystemRepairFun) LoadSystem(pbSystem *pb.PBPlayerSystem) {
	this.loadData(pbSystem.RepairData)
	this.UpdateSave(false)
}

func (this *PlayerSystemRepairFun) loadData(pbData *pb.PBPlayerSystemRepairData) {
	if pbData == nil {
		pbData = &pb.PBPlayerSystemRepairData{}
	}

	this.pbData = pbData

	this.UpdateSave(true)
}

// 加载完成
func (this *PlayerSystemRepairFun) LoadComplete() {
	if len(this.listRepairFun) <= 0 {
		return
	}

	if this.pbData.Version >= this.listRepairFun[len(this.listRepairFun)-1].versionID {
		return
	}

	if this.pbData.Version == 0 {
		this.pbData.Version = this.listRepairFun[len(this.listRepairFun)-1].versionID
		this.pbData.VersionTime = base.GetNow()
		this.UpdateSave(true)
		return
	}

	for _, repair := range this.listRepairFun {
		if this.pbData.Version >= repair.versionID {
			continue
		}

		repair.fun()
		plog.Info("(this *PlayerSystemRepairFun) LoadComplete id:%d version:%d", this.AccountId, repair.versionID)
		this.pbData.Version = repair.versionID
	}

	this.pbData.VersionTime = base.GetNow()
	this.UpdateSave(true)
}

// 存储到数据库
func (this *PlayerSystemRepairFun) SaveSystem(pbSystem *pb.PBPlayerSystem) bool {
	if pbSystem.RepairData == nil {
		pbSystem.RepairData = new(pb.PBPlayerSystemRepairData)
	}

	pbSystem.RepairData = this.pbData

	return this.BSave
}
func (this *PlayerSystemRepairFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerSystemRepairData{}
}

// 设置玩家数据
func (this *PlayerSystemRepairFun) SetUserTypeInfo(pbData proto.Message) bool {
	if pbData == nil {
		return false
	}
	pbSystem := pbData.(*pb.PBPlayerSystemRepairData)
	if pbSystem == nil {
		return false
	}
	this.loadData(pbSystem)
	this.LoadComplete()
	return true
}

func (this *PlayerSystemRepairFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	if pbData.System == nil {
		pbData.System = new(pb.PBPlayerSystem)
	}

	this.SaveSystem(pbData.System)
}
func (this *PlayerSystemRepairFun) repairTest() {
	plog.Info("(this *PlayerSystemRepairFun) repairTest id:%d version:%d", this.AccountId, this.pbData.Version)

}

// 修复商店红点
func (this *PlayerSystemRepairFun) repairShopRed() {
	this.getPlayerSystemShopFun().repairShopRed()
	plog.Info("(this *PlayerSystemRepairFun) repairShopRed id:%d version:%d", this.AccountId, this.pbData.Version)

}

func (this *PlayerSystemRepairFun) repairMainTask0828() {
	this.getPlayerSystemTaskFun().repairMainTask0828()
	plog.Info("(this *PlayerSystemRepairFun) repairTest id:%d version:%d", this.AccountId, this.pbData.Version)

}
