package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common"
	"corps/common/cfgData"
	"corps/common/orm/redis"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/proto"
)

type (
	PlayerMailFun struct {
		PlayerFun
		orderId        uint32
		allMailOrderId uint32                //全服邮件索引
		mapMail        map[uint32]*pb.PBMail //邮件数据 key:邮件id
	}
)

func (this *PlayerMailFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
	this.mapMail = make(map[uint32]*pb.PBMail)
}

// 加载背包数据
func (this *PlayerMailFun) Load(pData []byte) {
	this.BSave = false
	pbData := &pb.PBPlayerMail{}
	proto.Unmarshal(pData, pbData)
	this.loadData(pbData)
	this.UpdateSave(false)
}
func (this *PlayerMailFun) loadData(pbData *pb.PBPlayerMail) {
	if pbData == nil {
		pbData = &pb.PBPlayerMail{}
	}

	this.orderId = pbData.OrderId
	this.allMailOrderId = pbData.AllOrderId
	this.mapMail = make(map[uint32]*pb.PBMail)
	//过滤过期的邮件 已经领取的邮件
	uNow := base.GetNow()
	for _, v := range pbData.MailList {
		if v.State == pb.EmMailState_ReadRecieve && v.ExpireTime > 0 && v.ExpireTime < uNow {
			continue
		}

		this.mapMail[v.Id] = v
	}
	this.UpdateSave(true)
}

// 保存
func (this *PlayerMailFun) Save(bNow bool) {
	if !this.BSave {
		return
	}

	this.BSave = false

	pbData := &pb.PBPlayerMail{}
	this.SavePb(pbData)

	//通知db保存玩家数据
	buff, _ := proto.Marshal(pbData)
	cluster.SendToDb(&pb.RpcHead{Id: this.AccountId}, "DbPlayerMgr", "SavePlayerDB", this.PbType, buff, bNow)
}

// 保存
func (this *PlayerMailFun) SavePb(pbData *pb.PBPlayerMail) {
	pbData.OrderId = this.orderId
	pbData.AllOrderId = this.allMailOrderId
	for _, v := range this.mapMail {
		pbData.MailList = append(pbData.MailList, v)
	}
}

func (this *PlayerMailFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	pbData.Mail = new(pb.PBPlayerMail)
	this.SavePb(pbData.Mail)
}

// 心跳包
func (this *PlayerMailFun) Heat() {
	this.timerReadMail()
}

// 定时器读取全服邮件
func (this *PlayerMailFun) timerReadMail() {
	//遍历全服邮件
	redisCommon := redis.GetCommonRedis()
	if redisCommon == nil {
		return
	}

	//判断是否全部领取
	uMaxOrder := redisCommon.GetUint32(base.ERK_MailOrder)
	if this.allMailOrderId >= uMaxOrder {
		return
	}

	//遍历所有的
	mapKey := redisCommon.HGetAll(base.ERK_CommonAllMail)
	if len(mapKey) <= 0 {
		return
	}

	uCurTime := base.GetNow()
	for k, v := range mapKey {
		uCurOrder := base.StringToUInt32(k)
		if uCurOrder <= this.allMailOrderId {
			continue
		}

		stMail := &common.MailAllInfo{}
		err := json.Unmarshal([]byte(v), stMail)
		if err != nil {
			continue
		}

		if stMail.StartExpireTime > 0 && stMail.EndExpireTime > 0 {
			if uCurTime < stMail.StartExpireTime || uCurTime > stMail.EndExpireTime {
				continue
			}
		} else {
			//新玩家收不到邮件
			if this.getPlayerBaseFun().GetRegTime() > stMail.Time {
				continue
			}
		}

		//vip等级
		if stMail.VipMinLevel > 0 {
			uVipLevel := this.getPlayerBaseFun().GetVipLevel()
			if uVipLevel < stMail.VipMinLevel || uVipLevel > stMail.VipMaxLevel {
				continue
			}
		}

		//注册天数
		if stMail.RegMinDay > 0 {
			uRegDay := this.getPlayerBaseFun().GetRegDays()
			if uRegDay < stMail.RegMinDay || uRegDay > stMail.RegMaxDay {
				continue
			}
		}

		//组装道具
		pbItem := []*pb.PBAddItemData{}
		if stMail.Items != nil && len(stMail.Items) > 0 {
			for _, v := range stMail.Items {
				pbItem = append(pbItem, &pb.PBAddItemData{
					Kind:   v.Kind,
					Id:     v.Id,
					Count:  v.Count,
					Params: v.Params,
				})
			}
		}

		//判断邮件
		this.AddMail(&pb.RpcHead{Id: this.AccountId}, &pb.PBMail{
			Title:    stMail.Title,
			Content:  stMail.Content,
			Sender:   stMail.Sender,
			SendTime: stMail.Time,
			Item:     pbItem,
		})

	}

	this.allMailOrderId = uMaxOrder
	this.UpdateSave(true)
}

// 增加模版邮件
func (this *PlayerMailFun) AddTempMail(head *pb.RpcHead, mailId cfgEnum.EMailId, emDoing pb.EmDoingType, arrItem []*common.ItemInfo, arrContentParam ...interface{}) {
	cfgMail := cfgData.GetCfgMailConfig(uint32(mailId))
	if cfgMail == nil {
		plog.Error("(this *PlayerMailFun) AddTempMail %d", mailId)
		return
	}

	pbMail := &pb.PBMail{
		Title:    cfgMail.Title,
		Content:  cfgMail.Content,
		SendTime: base.GetNow(),
	}

	//拼接参数
	if len(arrContentParam) > 0 {
		pbMail.Content = fmt.Sprintf(pbMail.Content, arrContentParam...)
	}

	for _, v := range arrItem {
		pbMail.Item = append(pbMail.Item, &pb.PBAddItemData{
			Kind:      v.Kind,
			Id:        v.Id,
			Count:     v.Count,
			Params:    v.Params,
			DoingType: emDoing,
		})
	}
	this.AddMail(head, pbMail)
}

// 增加邮件
func (this *PlayerMailFun) AddMail(head *pb.RpcHead, pbMail *pb.PBMail) {

	this.orderId++
	pbMail.Id = this.orderId
	pbMail.ExpireTime = base.GetNow() + uint64(cfgData.GetCfgGlobalConfig(cfgEnum.GlobalConfig_GLOBAL_MAIL_EXPIRETIME))
	pbMail.State = pb.EmMailState_NoRead

	//转化装备
	//组装道具
	arrItem := make([]*common.ItemInfo, 0)
	emDoing := pb.EmDoingType_EDT_Mail
	if pbMail.Item != nil && len(pbMail.Item) > 0 {

		for _, v := range pbMail.Item {
			arrItem = append(arrItem, &common.ItemInfo{
				Kind:   v.Kind,
				Id:     v.Id,
				Count:  v.Count,
				Params: v.Params,
			})
			emDoing = v.DoingType
		}
	}

	pbMail.Item = this.getPlayerBagFun().GetPbItems(arrItem, emDoing)
	this.mapMail[pbMail.Id] = pbMail

	//通知客户端
	cluster.SendToClient(head, &pb.NewMailNotify{
		PacketHead: &pb.IPacket{},
		Mail:       pbMail,
	}, cfgEnum.ErrorCode_Success)

	this.UpdateSave(true)
}

// 一键删除邮件
func (this *PlayerMailFun) DipDelPlayerMail(head *pb.RpcHead, mailId uint32, bReward bool) bool {
	mail, ok := this.mapMail[mailId]
	if !ok {
		return false
	}

	//领取附件
	if bReward {
		if mail.State != pb.EmMailState_ReadRecieve {
			this.awardOneMail(head, mailId, false)
		}
	}

	if this.DeleteMail(head, mailId) == cfgEnum.ErrorCode_Success {
		return false
	}

	return true
}

// 一键领取邮件请求
func (this *PlayerMailFun) OnekeyAwardMailRequest(head *pb.RpcHead, pbRequest *pb.OnekeyAwardMailRequest) {
	uCode := this.OnekeyAwardMail(head)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.OnekeyAwardMailResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 一键领取邮件
func (this *PlayerMailFun) OnekeyAwardMail(head *pb.RpcHead) cfgEnum.ErrorCode {
	pbResponse := &pb.OnekeyAwardMailResponse{
		PacketHead: &pb.IPacket{},
	}

	bBagFull := false
	emBagCode := cfgEnum.ErrorCode_Success
	arrPbItems := make([]*pb.PBAddItemData, 0)
	for _, mail := range this.mapMail {
		if mail.State == pb.EmMailState_ReadRecieve {
			continue
		}

		emCode, pbMail := this.awardOneMail(head, mail.Id, true)
		if emCode != cfgEnum.ErrorCode_Success {
			if emCode == cfgEnum.ErrorCode_BagFull || emCode == cfgEnum.ErrorCode_HeroBagFull {
				bBagFull = true
			}
			emBagCode = emCode
			continue
		}

		if len(pbMail.Item) > 0 {
			arrPbItems = append(arrPbItems, pbMail.Item...)
		}

		pbResponse.Mails = append(pbResponse.Mails, pbMail)
	}
	if len(arrPbItems) <= 0 {
		if bBagFull {
			cluster.SendCommonToClient(head, emBagCode)
			return plog.Print(this.AccountId, emBagCode)
		}
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoMailAward)
	}

	//恭喜获得
	this.getPlayerBagFun().CommonPrizeNotify(head, arrPbItems, pb.EmDoingType_EDT_Mail)

	//通知客户端
	cluster.SendToClient(head, pbResponse, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}

// 领取邮件
func (this *PlayerMailFun) awardOneMail(head *pb.RpcHead, mailId uint32, isOnekey bool) (cfgEnum.ErrorCode, *pb.PBMail) {
	mail, ok := this.mapMail[mailId]
	if !ok {
		return plog.Print(head.Id, cfgEnum.ErrorCode_NoMail, mailId, isOnekey), nil
	}

	if mail.State == pb.EmMailState_ReadRecieve {
		return plog.Print(head.Id, cfgEnum.ErrorCode_HavePrize, mailId, mail.State), mail
	}

	//判断背包已满
	uNeedSpare := uint32(0)
	uNeedHeropSpare := uint32(0)
	for _, info := range mail.Item {
		if info.Equipment != nil && info.Equipment.Id > 0 {
			uNeedSpare++
		}

		if info.Kind == uint32(cfgEnum.ESystemType_Hero) {
			uNeedHeropSpare += uint32(info.Count)
		}
	}

	if uNeedSpare > 0 && this.GetPlayerEquipmentFun().GetSpareBag() < uNeedSpare {
		return plog.Print(head.Id, cfgEnum.ErrorCode_BagFull, mailId, uNeedSpare, this.GetPlayerEquipmentFun().GetSpareBag()), mail
	}

	if uNeedHeropSpare > 0 && this.getPlayerHeroFun().GetSpareBag() < uNeedHeropSpare {
		return plog.Print(head.Id, cfgEnum.ErrorCode_HeroBagFull, mailId, uNeedHeropSpare, this.getPlayerHeroFun().GetSpareBag()), mail
	}

	//领取附件
	if mail.Item != nil && len(mail.Item) > 0 {
		this.getPlayerBagFun().AddPbItems(head, mail.Item, pb.EmDoingType_EDT_Mail, false)
	} else if isOnekey {
		return cfgEnum.ErrorCode_NoMail, mail
	}

	mail.State = pb.EmMailState_ReadRecieve
	this.UpdateSave(true)

	return cfgEnum.ErrorCode_Success, mail
}

// 领取邮件请求
func (this *PlayerMailFun) AwardMailRequest(head *pb.RpcHead, pbRequest *pb.AwardMailRequest) {
	uCode := this.AwardMail(head, pbRequest.MailId)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.AwardMailResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 领取邮件请求
func (this *PlayerMailFun) AwardMail(head *pb.RpcHead, mailId uint32) cfgEnum.ErrorCode {
	emCode, pbMail := this.awardOneMail(head, mailId, false)
	if emCode != cfgEnum.ErrorCode_Success {
		if emCode == cfgEnum.ErrorCode_BagFull {
			cluster.SendCommonToClient(head, cfgEnum.ErrorCode_BagFull)
		}

		return plog.Print(this.AccountId, emCode, mailId)
	}

	if pbMail == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoMail, mailId)
	}

	if len(pbMail.Item) > 0 {
		this.getPlayerBagFun().CommonPrizeNotify(head, pbMail.Item, pb.EmDoingType_EDT_Mail)
	}

	//通知客户端
	cluster.SendToClient(head, &pb.AwardMailResponse{
		PacketHead: &pb.IPacket{},
		Mail:       pbMail,
	}, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}

// 一键删除邮件请求
func (this *PlayerMailFun) OnekeyDeleteMailRequest(head *pb.RpcHead, pbRequest *pb.OnekeyDeleteMailRequest) {
	uCode := this.OnekeyDeleteMail(head)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.OnekeyDeleteMailResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 一键删除邮件请求
func (this *PlayerMailFun) OnekeyDeleteMail(head *pb.RpcHead) cfgEnum.ErrorCode {
	pbResponse := &pb.OnekeyDeleteMailResponse{
		PacketHead: &pb.IPacket{},
	}

	for _, mail := range this.mapMail {
		if mail.State != pb.EmMailState_ReadRecieve {
			continue
		}

		pbResponse.MailIds = append(pbResponse.MailIds, mail.Id)
	}

	//真正删除
	for i := 0; i < len(pbResponse.MailIds); i++ {
		this.delOneMail(head, pbResponse.MailIds[i])
	}

	//通知客户端
	cluster.SendToClient(head, pbResponse, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}

// 内部删除邮件
func (this *PlayerMailFun) delOneMail(head *pb.RpcHead, mailId uint32) {
	delete(this.mapMail, mailId)
	this.UpdateSave(true)
}

// 删除邮件请求
func (this *PlayerMailFun) DeleteMailRequest(head *pb.RpcHead, pbRequest *pb.DeleteMailRequest) {
	uCode := this.DeleteMail(head, pbRequest.MailId)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.DeleteMailResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 删除邮件请求
func (this *PlayerMailFun) DeleteMail(head *pb.RpcHead, mailId uint32) cfgEnum.ErrorCode {
	mail, ok := this.mapMail[mailId]
	if !ok {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoMail, mailId)
	}

	if mail.State != pb.EmMailState_ReadRecieve {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_NoDelMail, mailId)
	}

	//真正删除
	this.delOneMail(head, mailId)

	//通知客户端
	cluster.SendToClient(head, &pb.DeleteMailResponse{
		PacketHead: &pb.IPacket{},
		MailId:     mailId,
	}, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}
func (this *PlayerMailFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerMail{}
}

// 设置玩家数据
func (this *PlayerMailFun) SetUserTypeInfo(pbData proto.Message) bool {
	if pbData == nil {
		return false
	}

	pbSystem := pbData.(*pb.PBPlayerMail)
	if pbSystem == nil {
		return false
	}

	this.loadData(pbSystem)
	return true
}
