package player

import (
	"context"
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common"
	"corps/framework/cluster"
	"corps/framework/common/uerror"
	"corps/framework/message"
	"corps/framework/plog"
	"corps/pb"
	"corps/server/game/playerMgr/playerFun"
	"log"
	"reflect"

	"github.com/golang/protobuf/proto"
)

// 玩家登录
func (p *Player) Login(ctx context.Context) {
	head := p.GetRpcHead(ctx)
	p.offlineTime = 0

	//通知db加载玩家数据 失败
	if !cluster.SendToDb(head, "DbPlayerMgr", "LoadPlayerDB") {
		plog.Info("玩家登录失败 db找不到 Login id:%d", p.AccountId)
		p.SendToClient(head, &pb.LoginResponse{PacketHead: &pb.IPacket{Id: p.AccountId, Code: uint32(cfgEnum.ErrorCode_ServerBusy)}})
		return
	}

	plog.Info("玩家登录成功 Login id:%d", p.AccountId)
}

// 断线重连
func (p *Player) ReLogin(ctx context.Context) {
	head := p.GetRpcHead(ctx)
	if !p.isInGame {
		p.Login(ctx)
		return
	}

	//发送给网关
	p.loginSuccess(head)

	p.isNeedSendClient = false
	p.updateOffline()
	plog.Info("玩家重连登录成功 (p *Player) ReLogin id: %d", head.Id)
}

// 发送聊天消息请求
func (this *Player) PlayerFun(emType pb.PlayerDataType, pbRequest proto.Message) {
	pFun := this.getPlayerFun(emType)
	rType := reflect.TypeOf(pFun)
	funcName := message.GetMessageName(pbRequest)
	m, bEx := rType.MethodByName(funcName)
	if !bEx {
		log.Printf("PlayerSystemFun [%s] has no method", funcName)
		cluster.SendCommonToClient(&pb.RpcHead{Id: this.AccountId}, cfgEnum.ErrorCode_PARAM)
		return
	}
	in := make([]reflect.Value, 2)
	in[0] = reflect.ValueOf(pFun)
	in[1] = reflect.ValueOf(pbRequest)
	m.Func.Call(in)
}

// 心跳请求
func (this *Player) HeartbeatRequest(ctx context.Context, pbData *pb.HeartbeatRequest) {
	if this.ELoginState <= ELS_LoadComplete {
		plog.Info("心跳超时 HeartbeatRequest client no client id:%d is null", this.GetId())
		return
	}

	if pbData == nil {
		plog.Info("心跳超时 HeartbeatRequest id:%d is null", this.GetId())
		return
	}

	head := this.GetRpcHead(ctx)
	if base.GetNow()-pbData.Time > 1 {
		//plog.Info("心跳超时 HeartbeatRequest id:%d step:%d, sendtime: %d", this.GetId(), base.GetNow()-pbData.Time, pbData.Time)
	}
	if this.ELoginState != ELS_Init {
		this.updateOffline()

		//更新登出时间
		this.getPlayerSystemOfflineFun().UpdateLogoutTime()

		this.getPlayerSystemChargeFun().HeartbeatRequest(head)
	} else {
		plog.Error("HeartbeatRequest server no complete id:%d step:%d, sendtime: %d,seqid:%d", this.GetId(), base.GetNow()-pbData.Time, pbData.Time, head.SeqId)
	}

	//同步心跳包
	this.SendToClient(head, &pb.HeartbeatResponse{SendTime: pbData.Time, RecvTime: base.GetNow(), CurTime: base.GetNow(), PacketHead: &pb.IPacket{Id: this.AccountId}})
}

// gm命令
func (this *Player) GmFuncRequest(ctx context.Context, pbRequest *pb.GmFuncRequest) {
	this.getPlayerBaseFun().GmFuncRequest(this.GetRpcHead(ctx), pbRequest)
}

// 兑换码
func (this *Player) GiftCodeRequest(ctx context.Context, pbRequest *pb.GiftCodeRequest) {
	this.getPlayerSystemCommonFun().GiftCodeRequest(this.GetRpcHead(ctx), pbRequest)
}
func (this *Player) NoticeRequest(ctx context.Context, pbRequest *pb.NoticeRequest) {
	this.getPlayerSystemCommonFun().NoticeRequest(this.GetRpcHead(ctx))
}
func (this *Player) AdvertRequest(ctx context.Context, req *pb.AdvertRequest) {
	this.getPlayerSystemCommonFun().AdvertRequest(this.GetRpcHead(ctx), req)
}
func (this *Player) PageOpenRequest(ctx context.Context, req *pb.PageOpenRequest) {
	this.getPlayerSystemCommonFun().PageOpenRequest(this.GetRpcHead(ctx), req)
}
func (this *Player) SystemOpenPrizeRequest(ctx context.Context, req *pb.SystemOpenPrizeRequest) {
	this.getPlayerSystemCommonFun().SystemOpenPrizeRequest(this.GetRpcHead(ctx), req)
}
func (this *Player) EjectAdvertRequest(ctx context.Context, req *pb.EjectAdvertRequest) {
	this.getPlayerSystemCommonFun().EjectAdvertRequest(this.GetRpcHead(ctx), req)
}

/***************************道具开始****************************************************/
// 道具使用
func (this *Player) ItemUseRequest(ctx context.Context, pbRequest *pb.ItemUseRequest) {
	this.getPlayerBagFun().ItemUseRequest(this.GetRpcHead(ctx), pbRequest)
}
func (this *Player) ItemUseShowRequest(ctx context.Context, pbRequest *pb.ItemUseShowRequest) {
	this.getPlayerBagFun().ItemUseShowRequest(this.GetRpcHead(ctx), pbRequest)
}
func (this *Player) ItemSelectRequest(ctx context.Context, pbRequest *pb.ItemSelectRequest) {
	this.getPlayerBagFun().ItemSelectRequest(this.GetRpcHead(ctx), pbRequest)
}
func (this *Player) ItemBuyRequest(ctx context.Context, pbRequest *pb.ItemBuyRequest) {
	this.getPlayerBagFun().ItemBuyRequest(this.GetRpcHead(ctx), pbRequest)
}

/***************************道具结束****************************************************/

/***************************职业功能开始****************************************************/
// 更新职业数据
func (this *Player) DipUpdateProf(ctx context.Context, pbRequest *pb.PBPlayerSystemProfInfo) {

	this.getPlayerSystemProfessionFun().DipUpdateProf(this.GetRpcHead(ctx), pbRequest)
}

// 职业升级请求
func (this *Player) ProfessionLevelRequest(ctx context.Context, pbRequest *pb.ProfessionLevelRequest) {
	this.getPlayerSystemProfessionFun().ProfessionLevelRequest(this.GetRpcHead(ctx), pbRequest)
}

// 职业突破请求
func (this *Player) ProfessionGradeRequest(ctx context.Context, pbRequest *pb.ProfessionGradeRequest) {
	this.getPlayerSystemProfessionFun().ProfessionGradeRequest(this.GetRpcHead(ctx), pbRequest)
}

// 职业巅峰升级请求
func (this *Player) ProfessionPeakRequest(ctx context.Context, pbRequest *pb.ProfessionPeakRequest) {
	this.getPlayerSystemProfessionFun().ProfessionPeakRequest(this.GetRpcHead(ctx), pbRequest)
}

// 职业巅峰升级请求
func (this *Player) ProfessionPeakLevelRequest(ctx context.Context, pbRequest *pb.ProfessionPeakLevelRequest) {
	this.getPlayerSystemProfessionFun().ProfessionPeakLevelRequest(this.GetRpcHead(ctx), pbRequest)
}

// 职业部位升级请求
func (this *Player) ProfessionPartLevelRequest(ctx context.Context, pbRequest *pb.ProfessionPartLevelRequest) {
	this.getPlayerSystemProfessionFun().ProfessionPartLevelRequest(this.GetRpcHead(ctx), pbRequest)
}

// 职业部位一键升级请求
func (this *Player) ProfessionPartOnekeyLevelRequest(ctx context.Context, pbRequest *pb.ProfessionPartOnekeyLevelResponse) {
	this.getPlayerSystemProfessionFun().ProfessionPartOnekeyLevelRequest(this.GetRpcHead(ctx), pbRequest)
}

// 职业装备穿戴请求
func (this *Player) ProfessionEquipRequest(ctx context.Context, pbRequest *pb.ProfessionEquipRequest) {
	this.getPlayerSystemProfessionFun().ProfessionEquipRequest(this.GetRpcHead(ctx), pbRequest)
}

// 职业装备一键穿戴请求
func (this *Player) ProfessionOnekeyUnEquipRequest(ctx context.Context, pbRequest *pb.ProfessionOnekeyUnEquipRequest) {
	this.getPlayerSystemProfessionFun().ProfessionOnekeyUnEquipRequest(this.GetRpcHead(ctx), pbRequest)
}

// 职业部位精炼请求
func (this *Player) ProfessionPartRefineRequest(ctx context.Context, pbRequest *pb.ProfessionPartRefineRequest) {
	this.getPlayerSystemProfessionFun().ProfessionPartRefineRequest(this.GetRpcHead(ctx), pbRequest)
}

// 职业部位一键精炼请求
func (this *Player) ProfessionPartOnekeyRefineRequest(ctx context.Context, pbRequest *pb.ProfessionPartOnekeyRefineRequest) {
	this.getPlayerSystemProfessionFun().ProfessionPartOnekeyRefineRequest(this.GetRpcHead(ctx), pbRequest)
}

// 职业部位精炼突破请求
func (this *Player) ProfessionPartRefineTupoRequest(ctx context.Context, pbRequest *pb.ProfessionPartRefineTupoRequest) {
	this.getPlayerSystemProfessionFun().ProfessionPartRefineTupoRequest(this.GetRpcHead(ctx), pbRequest)
}

func (this *Player) ProfessionPartResetRequest(ctx context.Context, pbRequest *pb.ProfessionPartResetRequest) {
	this.getPlayerSystemProfessionFun().ProfessionPartResetRequest(this.GetRpcHead(ctx), pbRequest)
}

// 删除所有装备
func (this *Player) DipDelAllEquipment(ctx context.Context) {
	this.GetPlayerEquipmentFun().DipDelAllEquipment(this.GetRpcHead(ctx))
}

// 装备分解请求
func (this *Player) EquipmentSplitRequest(ctx context.Context, pbRequest *pb.EquipmentSplitRequest) {
	this.GetPlayerEquipmentFun().EquipmentSplitRequest(this.GetRpcHead(ctx), pbRequest)
}

// 装备格子购买请求
func (this *Player) EquipmentBuyPosRequest(ctx context.Context, pbRequest *pb.EquipmentBuyPosRequest) {
	this.GetPlayerEquipmentFun().EquipmentBuyPosRequest(this.GetRpcHead(ctx), pbRequest)
}

// 装备自动分解请求
func (this *Player) EquipmentAutoSplitRequest(ctx context.Context, pbRequest *pb.EquipmentAutoSplitRequest) {
	this.GetPlayerEquipmentFun().EquipmentAutoSplitRequest(this.GetRpcHead(ctx), pbRequest)
}

// 装备自动分解请求
func (this *Player) EquipmentLockRequest(ctx context.Context, pbRequest *pb.EquipmentLockRequest) {
	this.GetPlayerEquipmentFun().EquipmentLockRequest(this.GetRpcHead(ctx), pbRequest)
}

// 挂机装备领取请求
func (this *Player) HookEquipmentAwardRequest(ctx context.Context, pbRequest *pb.HookEquipmentAwardRequest) {
	this.GetPlayerEquipmentFun().HookEquipmentAwardRequest(this.GetRpcHead(ctx), pbRequest)
}

/***************************职业功能结束****************************************************/
/*************************** 邮件功能结束****************************************************/
// 一键领取邮件请求
func (this *Player) OnekeyAwardMailRequest(ctx context.Context, pbRequest *pb.OnekeyAwardMailRequest) {
	this.getPlayerMailFun().OnekeyAwardMailRequest(this.GetRpcHead(ctx), pbRequest)
}

// 领取邮件请求
func (this *Player) AwardMailRequest(ctx context.Context, pbRequest *pb.AwardMailRequest) {
	this.getPlayerMailFun().AwardMailRequest(this.GetRpcHead(ctx), pbRequest)
}

// 一键删除邮件请求
func (this *Player) OnekeyDeleteMailRequest(ctx context.Context, pbRequest *pb.OnekeyDeleteMailRequest) {
	this.getPlayerMailFun().OnekeyDeleteMailRequest(this.GetRpcHead(ctx), pbRequest)
}

// 删除邮件请求
func (this *Player) DeleteMailRequest(ctx context.Context, pbRequest *pb.DeleteMailRequest) {
	this.getPlayerMailFun().DeleteMailRequest(this.GetRpcHead(ctx), pbRequest)
}

// 删除邮件
func (this *Player) DipDelPlayerMail(ctx context.Context, mailId uint32, bReward bool) {
	this.getPlayerMailFun().DipDelPlayerMail(this.GetRpcHead(ctx), mailId, bReward)
}

/*************************** 邮件功能结束****************************************************/
/***************************英雄功能结束****************************************************/
// 更新英雄
func (this *Player) DipUpdateHero(ctx context.Context, pbRequest *pb.PBHero) {
	this.getPlayerHeroFun().DipUpdateHero(this.GetRpcHead(ctx), pbRequest)
}

// 删除英雄
func (this *Player) DipDelHero(ctx context.Context, sn uint32) {
	this.getPlayerHeroFun().DipDelHero(this.GetRpcHead(ctx), sn)
}

// 英雄升星请求
func (this *Player) HeroUpStarRequest(ctx context.Context, pbRequest *pb.HeroUpStarRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.HeroUpStarResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerHeroFun().HeroUpStarRequest(head, pbRequest, rsp)
	if err != nil {
		plog.Debug("code: %d, error: %v", uerror.GetCode(err), err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

// 英雄自动升星请求
func (this *Player) HeroAutoUpStarRequest(ctx context.Context, pbRequest *pb.HeroAutoUpStarRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.HeroAutoUpStarResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerHeroFun().HeroAutoUpStarRequest(head, pbRequest, rsp)
	if err != nil {
		plog.Debug("code: %d, error: %v", uerror.GetCode(err), err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

// 英雄觉醒请求
func (this *Player) HeroAwakenLevelRequest(ctx context.Context, pbRequest *pb.HeroAwakenLevelRequest) {
	this.getPlayerHeroFun().HeroAwakenLevelRequest(this.GetRpcHead(ctx), pbRequest)
}

// 英雄重生请求
func (this *Player) HeroRebirthRequest(ctx context.Context, pbRequest *pb.HeroRebirthRequest) {
	this.getPlayerHeroFun().HeroRebirthRequest(this.GetRpcHead(ctx), pbRequest)
}

// 英雄上阵请求
func (this *Player) HeroGameHeroListRequest(ctx context.Context, pbRequest *pb.HeroGameHeroListRequest) {
	this.getPlayerHeroFun().HeroGameHeroListRequest(this.GetRpcHead(ctx), pbRequest)
}

// 英雄上阵星星
func (this *Player) HeroBattleStarChangeRequest(ctx context.Context, req *pb.HeroBattleStarChangeRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.HeroBattleStarChangeResponse{PacketHead: &pb.IPacket{}}
	errCode := this.getPlayerHeroFun().HeroBattleStarChangeRequest(head, req, rsp)
	cluster.SendToClient(head, rsp, errCode)
}

// 英雄图鉴激活请求
func (this *Player) HeroBookActiveRequest(ctx context.Context, pbRequest *pb.HeroBookActiveRequest) {
	this.getPlayerHeroFun().HeroBookActiveRequest(this.GetRpcHead(ctx), pbRequest)
}

// 英雄图鉴升星请求
func (this *Player) HeroBookUpStarRequest(ctx context.Context, pbRequest *pb.HeroBookUpStarRequest) {
	this.getPlayerHeroFun().HeroBookUpStarRequest(this.GetRpcHead(ctx), pbRequest)
}

/***************************英雄功能结束****************************************************/
/***************************战斗功能开始****************************************************/
// 挑战开始请求
func (this *Player) BattleBeginRequest(ctx context.Context, pbRequest *pb.BattleBeginRequest) {
	this.getPlayerSystemBattleFun().BattleBeginRequest(this.GetRpcHead(ctx), pbRequest)
}

// 挑战结束请求
func (this *Player) BattleEndRequest(ctx context.Context, pbRequest *pb.BattleEndRequest) {
	this.getPlayerSystemBattleFun().BattleEndRequest(this.GetRpcHead(ctx), pbRequest)
}

// 挑战记录请求
func (this *Player) BattleRecordRequest(ctx context.Context, pbRequest *pb.BattleRecordRequest) {
	this.getPlayerSystemBattleFun().BattleRecordRequest(this.GetRpcHead(ctx), pbRequest)
}

func (this *Player) BattleFunBuyRequest(ctx context.Context, pbRequest *pb.BattleFunBuyRequest) {
	this.getPlayerSystemBattleFun().BattleFunBuyRequest(this.GetRpcHead(ctx), pbRequest)
}

// 领取精英关卡奖励请求
func (this *Player) NormalBattlePrizeRequest(ctx context.Context, pbRequest *pb.NormalBattlePrizeRequest) {
	this.getPlayerSystemBattleNormalFun().NormalBattlePrizeRequest(this.GetRpcHead(ctx), pbRequest)
}

// 挂机自动推关设置请求
func (this *Player) HookBattleAutoMapRequest(ctx context.Context, pbRequest *pb.HookBattleAutoMapRequest) {
	this.getPlayerSystemBattleHookFun().HookBattleAutoMapRequest(this.GetRpcHead(ctx))
}

func (this *Player) HookBattleLootRequest(ctx context.Context, pbRequest *pb.HookBattleLootRequest) {
	this.getPlayerSystemBattleHookFun().HookBattleLootRequest(this.GetRpcHead(ctx), pbRequest)
}
func (this *Player) BattleScheduleSaveRequest(ctx context.Context, pbRequest *pb.BattleScheduleSaveRequest) {
	this.getPlayerSystemBattleFun().BattleScheduleSaveRequest(this.GetRpcHead(ctx), pbRequest)
}

// 挑战开始请求
func (this *Player) BattleNormalCardRequest(ctx context.Context, pbRequest *pb.BattleNormalCardRequest) {
	this.getPlayerSystemBattleNormalFun().BattleNormalCardRequest(this.GetRpcHead(ctx), pbRequest)
}

/***************************战斗功能结束****************************************************/

/***************************宝箱系统功能****************************************************/
func (this *Player) BoxOpenRequest(ctx context.Context, req *pb.BoxOpenRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.BoxOpenResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerSystemBoxFun().BoxOpenRequest(head, req, rsp)
	if err != nil {
		plog.Error("error: %v", err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

func (this *Player) BoxProgressRewardRequest(ctx context.Context, req *pb.BoxProgressRewardRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.BoxProgressRewardResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerSystemBoxFun().BoxProgressRewardRequest(head, req, rsp)
	if err != nil {
		plog.Error("error: %v", err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

/***************************宝箱系统功能结束****************************************************/
// 设置主线任务数据
func (this *Player) MainTaskFinishRequest(ctx context.Context, pbRequest *pb.MainTaskFinishRequest) {
	this.doRequest("MainTaskFinishRequest", ctx, pbRequest)
}

func (this *Player) DailyTaskFinishRequest(ctx context.Context, pbRequest *pb.DailyTaskFinishRequest) {
	this.getPlayerSystemTaskFun().DailyTaskFinishRequest(this.GetRpcHead(ctx), pbRequest)
}

// 设置主线任务数据
func (this *Player) DailyTaskScorePrizeRequest(ctx context.Context, pbRequest *pb.DailyTaskScorePrizeRequest) {
	this.getPlayerSystemTaskFun().DailyTaskScorePrizeRequest(this.GetRpcHead(ctx))
}
func (this *Player) SetClientRequest(ctx context.Context, pbRequest *pb.SetClientRequest) {
	this.getPlayerSystemClientFun().SetClientRequest(this.GetRpcHead(ctx), pbRequest)
}

// *********************设置系统********************
// 修改玩家名字
func (this *Player) ChangePlayerNameRequest(ctx context.Context, pbData *pb.ChangePlayerNameRequest) {

	this.getPlayerBaseFun().ChangePlayerNameRequest(this.GetRpcHead(ctx), pbData.PlayerName)
}

// 修改头像
func (this *Player) ChangeAvatarRequest(ctx context.Context, req *pb.ChangeAvatarRequest) {
	this.getPlayerSystemCommonFun().ChangeAvatarRequest(this.GetRpcHead(ctx), req.AvatarID)
}

// 修改头像框
func (this *Player) ChangeAvatarFrameRequest(ctx context.Context, req *pb.ChangeAvatarFrameRequest) {
	this.getPlayerSystemCommonFun().ChangeAvatarFrameRequest(this.GetRpcHead(ctx), req.FrameID)
}

// *************************商店******************************
func (this *Player) ShopBuyRequest(ctx context.Context, req *pb.ShopBuyRequest) {
	this.getPlayerSystemShopFun().ShopBuyRequest(this.GetRpcHead(ctx), req)
}
func (this *Player) ShopRefreshRequest(ctx context.Context, req *pb.ShopRefreshRequest) {
	this.getPlayerSystemShopFun().ShopRefreshRequest(this.GetRpcHead(ctx), req)
}
func (this *Player) ShopOpenRequest(ctx context.Context, req *pb.ShopOpenRequest) {
	this.getPlayerSystemShopFun().ShopOpenRequest(this.GetRpcHead(ctx), req)
}
func (this *Player) ShopExchangeRequest(ctx context.Context, req *pb.ShopExchangeRequest) {
	this.getPlayerSystemShopFun().ShopExchangeRequest(this.GetRpcHead(ctx), req)
}

// *************************抽奖******************************
func (this *Player) DrawRequest(ctx context.Context, req *pb.DrawRequest) {
	this.getPlayerSystemDrawFun().DrawRequest(this.GetRpcHead(ctx), req)
}
func (this *Player) DrawPrizeInfoRequest(ctx context.Context, req *pb.DrawPrizeInfoRequest) {
	this.getPlayerSystemDrawFun().DrawPrizeInfoRequest(this.GetRpcHead(ctx), req)
}
func (this *Player) DrawScorePrizeRequest(ctx context.Context, req *pb.DrawScorePrizeRequest) {
	this.getPlayerSystemDrawFun().DrawScorePrizeRequest(this.GetRpcHead(ctx), req)
}

// *************************充值******************************
func (this *Player) ChargeOrderRequest(ctx context.Context, req *pb.ChargeOrderRequest) {
	this.getPlayerSystemChargeFun().ChargeOrderRequest(this.GetRpcHead(ctx), req)
}
func (this *Player) DipUpdatePlayerCharge(ctx context.Context, ProductId uint32) {
	this.getPlayerSystemChargeFun().DipUpdatePlayerCharge(this.GetRpcHead(ctx), ProductId)
}
func (this *Player) FirstChargePrizeRequest(ctx context.Context, req *pb.FirstChargePrizeRequest) {
	this.getPlayerSystemChargeFun().FirstChargePrizeRequest(this.GetRpcHead(ctx), req)
}
func (this *Player) BPPrizeRequest(ctx context.Context, req *pb.BPPrizeRequest) {
	this.GetChargeBPFun().BPPrizeRequest(this.GetRpcHead(ctx), req)
}
func (this *Player) ChargeCardPrizeRequest(ctx context.Context, req *pb.ChargeCardPrizeRequest) {
	this.GetChargeCardFun().ChargeCardPrizeRequest(this.GetRpcHead(ctx), req)
}
func (this *Player) ChargeQueryRequest(ctx context.Context, req *pb.ChargeQueryRequest) {
	this.getPlayerSystemChargeFun().ChargeQueryRequest(this.GetRpcHead(ctx))
}

// *************************星源******************************
func (this *Player) HookTechLevelRequest(ctx context.Context, req *pb.HookTechLevelRequest) {
	this.getPlayerSystemHookTechFun().HookTechLevelRequest(this.GetRpcHead(ctx), req)
}
func (this *Player) HookTechSpeedRequest(ctx context.Context, req *pb.HookTechSpeedRequest) {
	this.getPlayerSystemHookTechFun().HookTechSpeedRequest(this.GetRpcHead(ctx), req)
}

// *************************晶核系统******************************

/***************************七天活动开始****************************************************/
func (this *Player) SevenDayActivePrizeRequest(ctx context.Context, pbRequest *pb.SevenDayActivePrizeRequest) {
	this.getPlayerSystemSevenDayFun().SevenDayActivePrizeRequest(this.GetRpcHead(ctx), pbRequest)
}
func (this *Player) SevenDayTaskPrizeRequest(ctx context.Context, pbRequest *pb.SevenDayTaskPrizeRequest) {
	this.getPlayerSystemSevenDayFun().SevenDayTaskPrizeRequest(this.GetRpcHead(ctx), pbRequest)
}
func (this *Player) SevenDayBuyGiftRequest(ctx context.Context, pbRequest *pb.SevenDayBuyGiftRequest) {
	this.getPlayerSystemSevenDayFun().SevenDayBuyGiftRequest(this.GetRpcHead(ctx), pbRequest)
}

/***************************七天活动结束****************************************************/

func (this *Player) BookStageRewardRequest(ctx context.Context, req *pb.BookStageRewardRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.BookStageRewardResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerCrystalFun().BookStageRewardRequest(head, req, rsp)
	if err != nil {
		plog.Error("error: %v", err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

func (this *Player) BookCollectionCoinRequest(ctx context.Context, req *pb.BookCollectionCoinRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.BookCollectionCoinResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerCrystalFun().BookCollectionCoinRequest(head, req, rsp)
	if err != nil {
		plog.Error("error: %v", err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

func (this *Player) CrystalRedefineRequest(ctx context.Context, req *pb.CrystalRedefineRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.CrystalRedefineResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerCrystalFun().CrystalRedefineRequest(head, req, rsp)
	if err != nil {
		plog.Error("error: %v", err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

func (this *Player) CrystalRobotUpgradeRequest(ctx context.Context, req *pb.CrystalRobotUpgradeRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.CrystalRobotUpgradeResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerCrystalFun().CrystalRobotUpgradeRequest(head, req, rsp)
	if err != nil {
		plog.Error("error: %v", err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

func (this *Player) CrystalUpgradeRequest(ctx context.Context, req *pb.CrystalUpgradeRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.CrystalUpgradeResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerCrystalFun().CrystalUpgradeRequest(head, req, rsp)
	if err != nil {
		plog.Error("error: %v", err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

func (this *Player) CrystalGenerateRequest(ctx context.Context, req *pb.CrystalGenerateRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.CrystalGenerateResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerCrystalFun().CrystalGenerateRequest(head, req, rsp)
	if err != nil {
		plog.Debug("code: %d, error: %v", uerror.GetCode(err), err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

// 基因系统
func (this *Player) getPlayerSystemGeneFun() *playerFun.PlayerSystemGeneFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemGene).(*playerFun.PlayerSystemGeneFun)
}

func (this *Player) GeneSchemeChangeRequest(ctx context.Context, req *pb.GeneSchemeChangeRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.GeneSchemeChangeResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerSystemGeneFun().GeneSchemeChangeRequest(head, req, rsp)
	if err != nil {
		plog.Debug("code: %d, error: %v", uerror.GetCode(err), err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

func (this *Player) GeneSchemeResetRequest(ctx context.Context, req *pb.GeneSchemeResetRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.GeneSchemeResetResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerSystemGeneFun().GeneSchemeResetRequest(head, req, rsp)
	if err != nil {
		plog.Debug("code: %d, error: %v", uerror.GetCode(err), err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

func (this *Player) GeneChangeNameRequest(ctx context.Context, req *pb.GeneChangeNameRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.GeneChangeNameResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerSystemGeneFun().GeneChangeNameRequest(head, req, rsp)
	if err != nil {
		plog.Debug("code: %d, error: %v", uerror.GetCode(err), err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

func (this *Player) GeneCardActiveRequest(ctx context.Context, req *pb.GeneCardActiveRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.GeneCardActiveResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerSystemGeneFun().GeneCardActiveRequest(head, req, rsp)
	if err != nil {
		plog.Debug("code: %d, error: %v", uerror.GetCode(err), err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

func (this *Player) GeneRobotActiveRequest(ctx context.Context, req *pb.GeneRobotActiveRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.GeneRobotActiveResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerSystemGeneFun().GeneRobotActiveRequest(head, req, rsp)
	if err != nil {
		plog.Debug("code: %d, error: %v", uerror.GetCode(err), err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

func (this *Player) GeneRobotCardActiveRequest(ctx context.Context, req *pb.GeneRobotCardActiveRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.GeneRobotCardActiveResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerSystemGeneFun().GeneRobotCardActiveRequest(head, req, rsp)
	if err != nil {
		plog.Debug("code: %d, error: %v", uerror.GetCode(err), err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

// 离线收益系统
func (this *Player) getPlayerSystemOfflineFun() *playerFun.PlayerSystemOfflineFun {
	return this.getPlayerFun(pb.PlayerDataType_SystemOffline).(*playerFun.PlayerSystemOfflineFun)
}

func (this *Player) OfflineIncomeRewardRequest(ctx context.Context, req *pb.OfflineIncomeRewardRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.OfflineIncomeRewardResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerSystemOfflineFun().OfflineIncomeRewardRequest(head, req, rsp)
	if err != nil {
		plog.Debug("code: %d, error: %v", uerror.GetCode(err), err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

func (this *Player) HeroPieceRequest(ctx context.Context, req *pb.HeroPieceRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.HeroPieceResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerHeroFun().HeroPieceRequest(head, req, rsp)
	if err != nil {
		plog.Debug("code: %d, error: %v", uerror.GetCode(err), err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

// -----------------------------词条技能系统---------------------------------
func (this *Player) EntryUnlockRequest(ctx context.Context, req *pb.EntryUnlockRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.EntryUnlockResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerCrystalFun().EntryUnlockRequest(head, req, rsp)
	if err != nil {
		plog.Debug("code: %d, error: %v", uerror.GetCode(err), err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

func (this *Player) EntryTriggerRequest(ctx context.Context, req *pb.EntryTriggerRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.EntryTriggerResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerCrystalFun().EntryTriggerRequest(head, req, rsp)
	if err != nil {
		plog.Debug("code: %d, error: %v", uerror.GetCode(err), err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

func (this *Player) GetPlayerDataRequest(ctx context.Context, req *pb.GetPlayerDataRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.GetPlayerDataResponse{PacketHead: &pb.IPacket{}, PlayerData: &pb.PBPlayerData{}}
	if fun := this.getPlayerFun(pb.PlayerDataType(req.DataType)); fun != nil {
		fun.SaveDataToClient(rsp.PlayerData)
	}
	cluster.SendToClient(head, rsp, 0)
}

func (this *Player) SetPlayerDataRequest(ctx context.Context, req *pb.SetPlayerDataRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.SetPlayerDataResponse{PacketHead: &pb.IPacket{}}
	this.DipSetUserTypeInfo(ctx, pb.PlayerDataType(req.DataType), req.JsonData)
	cluster.SendToClient(head, rsp, 0)
}

func (this *Player) RewardRequest(ctx context.Context, req *pb.RankRewardRequest) {
	head := this.GetRpcHead(ctx)
	items := common.AddItemToItemInfo(req.Rewards...)
	code := this.getPlayerBagFun().AddArrItem(head, items, req.Doing, req.Notify)
	if code == cfgEnum.ErrorCode_Success {
		this.getPlayerSystemChampionshipFun().SetChampionshipFlag(req.RankType, 2)
	}
	rsp := &pb.RankRewardResponse{PacketHead: &pb.IPacket{}}
	cluster.SendToClient(head, rsp, code)
}

func (this *Player) ChampionshipTaskRewardRequest(ctx context.Context, req *pb.ChampionshipTaskRewardRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.ChampionshipTaskRewardResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerSystemChampionshipFun().ChampionshipTaskRewardRequest(head, req, rsp)
	if err != nil {
		plog.Error("error: %v", err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

func (this *Player) ChampionshipInfoRequest(ctx context.Context, req *pb.ChampionshipInfoRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.ChampionshipInfoResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerSystemChampionshipFun().ChampionshipInfoRequest(head, req, rsp)
	if err != nil {
		plog.Error("error: %v", err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

/***************************世界boss开始****************************************************/
func (this *Player) WorldBossStagePrizeRequest(ctx context.Context, pbRequest *pb.WorldBossStagePrizeRequest) {
	this.getPlayerSystemWorldBossFun().WorldBossStagePrizeRequest(this.GetRpcHead(ctx), pbRequest)
}
func (this *Player) WorldBossBuyCountRequest(ctx context.Context, pbRequest *pb.WorldBossBuyCountRequest) {
	this.getPlayerSystemWorldBossFun().WorldBossBuyCountRequest(this.GetRpcHead(ctx), pbRequest)
}
func (this *Player) WorldBossSweepRequest(ctx context.Context, pbRequest *pb.WorldBossSweepRequest) {
	this.getPlayerSystemWorldBossFun().WorldBossSweepRequest(this.GetRpcHead(ctx), pbRequest)
}
func (this *Player) WorldBossBattleBeginRequest(ctx context.Context, pbRequest *pb.WorldBossBattleBeginRequest) {
	this.getPlayerSystemWorldBossFun().WorldBossBattleBeginRequest(this.GetRpcHead(ctx), pbRequest)
}
func (this *Player) WorldBossBattleEndRequest(ctx context.Context, pbRequest *pb.WorldBossBattleEndRequest) {
	this.getPlayerSystemWorldBossFun().WorldBossBattleEndRequest(this.GetRpcHead(ctx), pbRequest)
}
func (this *Player) WorldBossRecordRequest(ctx context.Context, pbRequest *pb.WorldBossRecordRequest) {
	this.getPlayerSystemWorldBossFun().WorldBossRecordRequest(this.GetRpcHead(ctx), pbRequest)
}
func (this *Player) OpenBossRequest(ctx context.Context, pbRequest *pb.OpenBossRequest) {
	this.getPlayerSystemWorldBossFun().OpenBossRequest(this.GetRpcHead(ctx))
}

/***************************世界boss结束****************************************************/
/***************************活动开始****************************************************/
func (this *Player) GrowRoadTaskPrizeRequest(ctx context.Context, pbRequest *pb.GrowRoadTaskPrizeRequest) {
	this.getPlayerActivityGrowRoadFun().GrowRoadTaskPrizeRequest(this.GetRpcHead(ctx), pbRequest)
}
func (this *Player) ActivityOpenRequest(ctx context.Context, pbRequest *pb.ActivityOpenRequest) {
	this.getPlayerSystemActivityFun().ActivityOpenRequest(this.GetRpcHead(ctx), pbRequest)
}
func (this *Player) ActivityFreePrizeRequest(ctx context.Context, pbRequest *pb.ActivityFreePrizeRequest) {
	this.getPlayerSystemActivityFun().ActivityFreePrizeRequest(this.GetRpcHead(ctx), pbRequest)
}

/***************************活动结束****************************************************/
func (this *Player) AdventureRewardRequest(ctx context.Context, req *pb.AdventureRewardRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.AdventureRewardResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerSystemActivityFun().GetActivityAdventureFun().AdventureRewardRequest(head, req, rsp)
	if err != nil {
		plog.Error("error: %v", err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}

func (this *Player) BattleHookPassRewardRequest(ctx context.Context, req *pb.BattleHookPassRewardRequest) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.BattleHookPassRewardResponse{PacketHead: &pb.IPacket{}}
	err := this.getPlayerSystemBattleHookFun().BattleHookPassRewardRequest(head, req, rsp)
	if err != nil {
		plog.Debug("code: %d, error: %v", uerror.GetCode(err), err)
	}
	cluster.SendToClient(head, rsp, uerror.GetCode(err))
}
