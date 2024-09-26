package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common"
	"corps/common/cfgData"
	report2 "corps/common/report"
	"corps/framework/cluster"
	"corps/framework/common/uerror"
	"corps/framework/plog"
	"corps/pb"
	"corps/server/game/module/entry"
	"math/rand"

	"github.com/golang/protobuf/proto"
)

type (
	PlayerHeroFun struct {
		PlayerFun
		orderId               uint32                                       //道具索引
		uFightPower           uint32                                       //战斗力
		uMaxHistoryFightPower uint32                                       //历史最大战斗力
		mapUpStarCount        map[uint32]uint32                            //升星次数
		mapHero               map[uint32]*PlayerHero                       //英雄数据 key：sn
		mapStarHero           map[uint32]map[uint32]map[uint32]*PlayerHero //英雄数据 key：星级 key2:id key3:sn
		bCalcFightpower       bool                                         //是否计算战斗力
		mapHeroBook           map[uint32]*pb.PBHeroBook                    //英雄图鉴数据
		mapHeroBookProp       map[uint32]uint32                            //英雄图鉴属性
		mapHeroTeam           map[uint32]*pb.PBHeroTeamList                //编队列表
		GlobalRandHeroProf    []uint32                                     //全局前10个蓝色和全局前10个紫色 5个职业乱序排序为一轮，保护两系统前两轮（第二轮重新乱序）
	}

	PlayerHero struct {
		*pb.PBHero
		mapProp *common.Property
		uScore  uint32
		bCalc   bool
	}
)

func (this *PlayerHeroFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
	this.mapHero = make(map[uint32]*PlayerHero)
	this.mapStarHero = make(map[uint32]map[uint32]map[uint32]*PlayerHero)
	this.bCalcFightpower = false
	this.mapUpStarCount = make(map[uint32]uint32)
	this.mapHeroBook = make(map[uint32]*pb.PBHeroBook)
	this.mapHeroBookProp = make(map[uint32]uint32)
	this.mapHeroTeam = make(map[uint32]*pb.PBHeroTeamList)
	this.uMaxHistoryFightPower = 0
	this.GlobalRandHeroProf = make([]uint32, 0)

}

func (this *PlayerHeroFun) updateCalcFightpower(bAll bool) {
	if bAll {
		for _, sn := range this.GetTeamList(uint32(cfgEnum.EHeroTeam_BattleNormal)) {
			playerHero := this.getHero(sn)
			if playerHero == nil {
				continue
			}

			playerHero.bCalc = true
		}
	}
	this.bCalcFightpower = true
}

// 新玩家 需要初始化数据
func (this *PlayerHeroFun) NewPlayer() {
	//送初始英雄
	cfgGlobal := cfgData.GetCfgGlobalArrConfig(cfgEnum.GlobalArrConfig_GLOBAL_CFG_NewAddHero)
	uHeroLen := len(cfgData.GetCfgGlobalArrConfig(cfgEnum.GlobalArrConfig_GLOBAL_TEAMPOS_UNLOCKMAP).Value)
	if cfgGlobal != nil {
		heroList := make([]uint32, 0)
		for i := 0; i < len(cfgGlobal.Value)/2; i++ {
			_, pbHero := this.addHero(&pb.RpcHead{Id: this.AccountId}, cfgGlobal.Value[i*2], cfgGlobal.Value[i*2+1], pb.EmDoingType_EDT_Gm)
			if pbHero != nil && len(heroList) < uHeroLen {
				heroList = append(heroList, pbHero.Sn)
			}
		}

		this.HeroGameHeroList(&pb.RpcHead{Id: this.AccountId}, &pb.PBHeroTeamList{HeroSn: heroList, TeamType: uint32(cfgEnum.EHeroTeam_BattleNormal)})
	}

	this.UpdateSave(true)
	this.Save(true)
}

// 加载数据
func (this *PlayerHeroFun) Load(pData []byte) {
	pbData := &pb.PBPlayerHero{}
	proto.Unmarshal(pData, pbData)

	this.loadData(pbData)

	this.UpdateSave(false)
}
func (this *PlayerHeroFun) loadData(pbData *pb.PBPlayerHero) {
	if pbData == nil {
		pbData = &pb.PBPlayerHero{}
	}

	this.mapHero = make(map[uint32]*PlayerHero)
	this.mapStarHero = make(map[uint32]map[uint32]map[uint32]*PlayerHero)
	this.mapUpStarCount = make(map[uint32]uint32)

	this.orderId = pbData.OrderId
	this.uFightPower = pbData.FightPower
	this.uMaxHistoryFightPower = pbData.MaxHistoryFightPower
	if this.uFightPower > pbData.MaxHistoryFightPower {
		this.uMaxHistoryFightPower = pbData.MaxHistoryFightPower
	}
	plog.Trace("MaxHistoryFightPower loadData id:%d max:%d cur:%d", this.AccountId, pbData.MaxHistoryFightPower, pbData.FightPower)
	for i := 0; i < len(pbData.UpStarCount); i++ {
		this.mapUpStarCount[pbData.UpStarCount[i].Key] = pbData.UpStarCount[i].Value
	}

	for i := 0; i < len(pbData.HeroList); i++ {
		pbPlayerHero := &PlayerHero{
			PBHero:  pbData.HeroList[i],
			mapProp: common.NewProperty(),
			uScore:  0,
			bCalc:   true,
		}
		this.updateHero(pbPlayerHero, 0)
	}
	for i := 0; i < len(pbData.HeroBookList); i++ {
		this.mapHeroBook[pbData.HeroBookList[i].Id] = pbData.HeroBookList[i]
	}

	//编队
	for i := 0; i < len(pbData.TeamList); i++ {
		this.mapHeroTeam[pbData.TeamList[i].TeamType] = pbData.TeamList[i]
	}

	this.GlobalRandHeroProf = pbData.GlobalRandHeroProf
	if len(this.GlobalRandHeroProf) <= 0 {
		//初始化  5蓝职业+5蓝职业 5紫色职业 5紫色职业 当前蓝色索引  当前紫色索引 全局前10个蓝色和全局前10个紫色 5个职业乱序排序为一轮，保护两系统前两轮（第二轮重新乱序）
		arrProf := make([]uint32, 0)
		for i := cfgEnum.EHydraProf_Tank; i <= cfgEnum.EHydraProf_Machinist; i++ {
			arrProf = append(arrProf, uint32(i))
		}
		for it := 0; it < 2; it++ {
			uTmpIndex := uint32(len(this.GlobalRandHeroProf))
			for jt := 0; jt < 2; jt++ {
				rand.Shuffle(len(arrProf), func(i, j int) {
					arrProf[i], arrProf[j] = arrProf[j], arrProf[i]
				})
				this.GlobalRandHeroProf = append(this.GlobalRandHeroProf, arrProf...)
			}
			this.GlobalRandHeroProf = append(this.GlobalRandHeroProf, uTmpIndex)
		}
	}

	this.UpdateSave(true)
}

// 获取下一个职业英雄 10个蓝色 1蓝色职业索引 10个紫色 1蓝色职业索引紫色职业索引
func (this *PlayerHeroFun) GetNextGlobalRandHeroProf(uQuality uint32, mapProfCount map[uint32]map[uint32]uint32) int32 {
	uProfIndex := uint32(10)
	if uQuality == uint32(cfgEnum.EQuality_Purple) {
		uProfIndex = uint32(21)
	}

	for {
		uCurIndex := this.GlobalRandHeroProf[uProfIndex]
		if uCurIndex >= uProfIndex {
			return -1
		}

		uTmpCount := uint32(0)
		if tmpInfo, ok := mapProfCount[uQuality]; ok {
			if _, tok := tmpInfo[this.GlobalRandHeroProf[uCurIndex]]; tok {
				uTmpCount = mapProfCount[uQuality][this.GlobalRandHeroProf[uCurIndex]]
			}
		}

		// 0-4 1个 5-9 2个
		if this.GetQualityProfHeroCount(uQuality, this.GlobalRandHeroProf[uCurIndex])+uTmpCount >= base.CeilU32(uCurIndex+1, 5) {
			this.GlobalRandHeroProf[uProfIndex]++
			continue
		}
		this.GlobalRandHeroProf[uProfIndex]++
		if _, ok := mapProfCount[uQuality]; !ok {
			mapProfCount[uQuality] = make(map[uint32]uint32)
		}
		mapProfCount[uQuality][this.GlobalRandHeroProf[uCurIndex]]++
		this.UpdateSave(true)
		return int32(this.GlobalRandHeroProf[uCurIndex])
	}

	return -1
}

// 查询职业品质英雄个数
func (this *PlayerHeroFun) GetQualityProfHeroCount(uQuality uint32, uProf uint32) uint32 {
	uCount := uint32(0)
	for _, pbHero := range this.mapHero {
		cfgHero := cfgData.GetCfgHero(pbHero.Id)
		if cfgHero == nil || cfgHero.Prof != uProf {
			continue
		}
		cfgStar := cfgData.GetCfgHeroStar(pbHero.Star)
		if cfgStar == nil || cfgStar.Quality != uQuality {
			continue
		}
		uCount++
	}
	return uCount

}

// 保存
func (this *PlayerHeroFun) Save(bNow bool) {
	if !this.BSave {
		return
	}

	this.BSave = false

	pbData := &pb.PBPlayerHero{}
	this.SavePb(pbData)

	//通知db保存玩家数据
	buff, _ := proto.Marshal(pbData)
	cluster.SendToDb(&pb.RpcHead{Id: this.AccountId}, "DbPlayerMgr", "SavePlayerDB", this.PbType, buff, bNow)
}

// 保存
func (this *PlayerHeroFun) SavePb(pbData *pb.PBPlayerHero) {
	for _, v := range this.mapHero {
		pbData.HeroList = append(pbData.HeroList, v.PBHero)
	}

	pbData.GlobalRandHeroProf = this.GlobalRandHeroProf
	pbData.OrderId = this.orderId
	pbData.FightPower = this.uFightPower
	pbData.MaxHistoryFightPower = this.uMaxHistoryFightPower
	plog.Trace("MaxHistoryFightPower  SavePb id:%d max:%d cur:%d", this.AccountId, this.uMaxHistoryFightPower, this.uFightPower)
	for k, v := range this.mapUpStarCount {
		pbData.UpStarCount = append(pbData.UpStarCount, &pb.PBU32U32{Key: k, Value: v})
	}

	for _, v := range this.mapHeroBook {
		pbData.HeroBookList = append(pbData.HeroBookList, v)
	}

	//编队
	for _, v := range this.mapHeroTeam {
		pbData.TeamList = append(pbData.TeamList, v)
	}
}

func (this *PlayerHeroFun) LoadPlayerDBFinish() {
	if _, ok := this.mapHeroTeam[uint32(cfgEnum.EHeroTeam_BattleNormal)]; !ok {
		this.mapHeroTeam[uint32(cfgEnum.EHeroTeam_BattleNormal)] = &pb.PBHeroTeamList{
			TeamType: uint32(cfgEnum.EHeroTeam_BattleNormal),
		}
	}
}

// 加载完成需要计算战斗力
func (this *PlayerHeroFun) LoadComplete() {
	if len(this.GlobalRandHeroProf) <= 0 {
		//初始化  5蓝职业+5蓝职业 5紫色职业 5紫色职业 当前蓝色索引  当前紫色索引 全局前10个蓝色和全局前10个紫色 5个职业乱序排序为一轮，保护两系统前两轮（第二轮重新乱序）
		arrProf := make([]uint32, 0)
		for i := cfgEnum.EHydraProf_Tank; i <= cfgEnum.EHydraProf_Machinist; i++ {
			arrProf = append(arrProf, uint32(i))
		}
		for it := 0; it < 2; it++ {
			uTmpIndex := uint32(len(this.GlobalRandHeroProf))
			for jt := 0; jt < 2; jt++ {
				rand.Shuffle(len(arrProf), func(i, j int) {
					arrProf[i], arrProf[j] = arrProf[j], arrProf[i]
				})
				this.GlobalRandHeroProf = append(this.GlobalRandHeroProf, arrProf...)
			}
			this.GlobalRandHeroProf = append(this.GlobalRandHeroProf, uTmpIndex)
		}
	}

	if _, ok := this.mapHeroTeam[uint32(cfgEnum.EHeroTeam_BattleNormal)]; !ok {
		this.mapHeroTeam[uint32(cfgEnum.EHeroTeam_BattleNormal)] = &pb.PBHeroTeamList{
			TeamType: uint32(cfgEnum.EHeroTeam_BattleNormal),
		}
	}
	this.CalcHeroBookProp(&pb.RpcHead{Id: this.AccountId})
}

func (this *PlayerHeroFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	pbData.Hero = new(pb.PBPlayerHero)
	this.SavePb(pbData.Hero)
}
func (this *PlayerHeroFun) getOrderId() uint32 {
	this.orderId++
	return this.orderId
}
func (this *PlayerHeroFun) GetFightPower() uint32 {
	return this.uFightPower
}
func (this *PlayerHeroFun) GetUpStarCount(uHeroId uint32) uint32 {
	uCount := uint32(0)
	if uHeroId == 0 {
		for _, v := range this.mapUpStarCount {
			uCount += v
		}
	} else {
		if _, ok := this.mapUpStarCount[uHeroId]; ok {
			uCount = this.mapUpStarCount[uHeroId]
		}
	}

	return uCount
}

// 增加英雄(或者转换成英雄碎片)
func (this *PlayerHeroFun) AddHeros(head *pb.RpcHead, uHeroId uint32, uStar uint32, uCount uint32, doingType pb.EmDoingType, bSend bool) cfgEnum.ErrorCode {
	if uCount <= 0 {
		return cfgEnum.ErrorCode_Success
	}
	// 加载配置
	cfgHero := cfgData.GetCfgHero(uHeroId)
	if cfgHero == nil {
		return cfgData.GetHeroErrorCode(uHeroId)
	}
	cfgStar := cfgData.GetCfgHeroStar(uStar)
	if cfgStar == nil {
		return cfgData.GetHeroStarErrorCode(uHeroId)
	}
	//成就类型（狗粮类型不触发成就）
	if cfgHero.IsMaterial <= 0 {
		this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_AddQualityHero, uCount, cfgStar.Quality)
	}
	// 判断图鉴
	bTips := false
	if _, ok := this.mapHeroBook[uHeroId]; !ok {
		bTips = true
	} else if uStar > this.mapHeroBook[uHeroId].MaxStar {
		bTips = true
	}
	// 判断是否已经解锁
	flag := false
	for _, hero := range this.mapHero {
		if hero.Id == uHeroId {
			flag = true
			break
		}
	}
	// 英雄转成碎片发放
	if !flag {
		uCount--
	}
	if uCount > 0 {
		errCode := this.getPlayerBagFun().AddItem(head, uint32(cfgEnum.ESystemType_Item), cfgHero.AwakenItemId, int64(cfgStar.PieceCount*uCount), doingType, true)
		if errCode != cfgEnum.ErrorCode_Success {
			return errCode
		}
	}
	// 英雄解锁
	var pbHero *pb.PBHero
	var errCode cfgEnum.ErrorCode
	if !flag {
		if errCode, pbHero = this.addHero(head, uHeroId, uStar, doingType); errCode != cfgEnum.ErrorCode_Success || pbHero == nil {
			return errCode
		}
		// 通知客户端
		this.HeroNotify(head, []*pb.PBHero{pbHero})
	}
	if bSend {
		//星级更新 图鉴更新
		if bTips && cfgHero.IsMaterial <= 0 {
			cluster.SendToClient(head, &pb.HeroBookNotify{
				PacketHead: &pb.IPacket{},
				HeroBook:   this.mapHeroBook[uHeroId],
			}, cfgEnum.ErrorCode_Success)

			if pbHero != nil {
				this.HeroNewStarNotify(head, []*pb.PBHero{pbHero})
			}
		}
		//自动上阵
		if cfgHero.IsMaterial <= 0 && pbHero != nil {
			this.AutoHeroGameHeroList(head, pbHero)
		}
	}
	return cfgEnum.ErrorCode_Success
}

func (this *PlayerHeroFun) GetSpareBag() uint32 {
	uMax := cfgData.GetCfgGlobalConfig(cfgEnum.GlobalConfig_GLOBAL_HERO_MAX_COUNT)
	if uint32(len(this.mapHero)) >= uMax {
		return 0
	}
	return uMax - uint32(len(this.mapHero))
}

// 增加英雄
func (this *PlayerHeroFun) addHero(head *pb.RpcHead, uHeroId uint32, uStar uint32, doingType pb.EmDoingType) (cfgEnum.ErrorCode, *pb.PBHero) {
	if uint32(len(this.mapHero)) >= cfgData.GetCfgGlobalConfig(cfgEnum.GlobalConfig_GLOBAL_HERO_MAX_COUNT) {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_HeroBagFull, uHeroId, uStar, doingType), nil
	}

	cfgHeroStar := cfgData.GetCfgHeroStar(uStar)
	if cfgHeroStar == nil {
		return plog.Print(this.AccountId, cfgData.GetHeroStarErrorCode(uStar), uHeroId, uStar, doingType), nil
	}

	cfgHero := cfgData.GetCfgHero(uHeroId)
	if cfgHero == nil {
		return plog.Print(this.AccountId, cfgData.GetHeroErrorCode(uHeroId), uHeroId, uStar, doingType), nil
	}

	playerHero := &PlayerHero{
		PBHero: &pb.PBHero{
			Sn:          this.getOrderId(),
			Id:          uHeroId,
			Star:        uStar,
			AwakenLevel: 0,
		},
		mapProp: common.NewProperty(),
		bCalc:   true,
	}

	//判断是否激活图鉴
	if cfgHero.IsMaterial <= 0 {
		if _, ok := this.mapHeroBook[uHeroId]; !ok {
			this.mapHeroBook[uHeroId] = &pb.PBHeroBook{
				Id:      uHeroId,
				Star:    0,
				MaxStar: uStar,
			}
		} else {
			if uStar > this.mapHeroBook[uHeroId].MaxStar {
				this.mapHeroBook[uHeroId].MaxStar = uStar
			}
		}
	} else {
		playerHero.bCalc = false
	}

	//算属性
	playerHero = this.updateHero(playerHero, 0)

	// 掉落头像
	if cfgHero.IsMaterial <= 0 {
		if errCode := this.getPlayerSystemCommonFun().AddHead(head, cfgHero.HeadId, doingType); errCode != cfgEnum.ErrorCode_Success {
			plog.Error("head: %v, hero: %v, dointType: %d", head, cfgHero, doingType)
		}
	}

	// 数据上报
	report2.Send(head, &report2.ReportAddItem{
		Kind:   uint32(cfgEnum.ESystemType_Hero),
		Doing:  uint32(doingType),
		ItemID: playerHero.Id,
		Add:    1,
		Total:  1,
		Params: []uint32{playerHero.Star},
		Sn:     playerHero.Sn,
	})
	return cfgEnum.ErrorCode_Success, playerHero.PBHero
}

// 通知客户端
func (this *PlayerHeroFun) HeroNotify(head *pb.RpcHead, arrList []*pb.PBHero) {
	cluster.SendToClient(head, &pb.HeroNotify{
		PacketHead: &pb.IPacket{},
		Info:       arrList,
	}, cfgEnum.ErrorCode_Success)
}

// 通知客户端
func (this *PlayerHeroFun) HeroNewStarNotify(head *pb.RpcHead, arrList []*pb.PBHero) {
	cluster.SendToClient(head, &pb.HeroNewStarNotify{
		PacketHead: &pb.IPacket{},
		Info:       arrList,
	}, cfgEnum.ErrorCode_Success)
}

// 英雄觉醒请求
func (this *PlayerHeroFun) HeroAwakenLevelRequest(head *pb.RpcHead, pbRequest *pb.HeroAwakenLevelRequest) {
	uCode := this.HeroAwakenLevel(head, pbRequest.Sn, pbRequest.CurLevel)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.HeroAwakenLevelResponse{
			PacketHead: &pb.IPacket{},
			Sn:         pbRequest.Sn,
			CurLevel:   pbRequest.CurLevel,
		}, uCode)
	}
}

// 英雄觉醒请求
func (this *PlayerHeroFun) HeroAwakenLevel(head *pb.RpcHead, uSn uint32, uCurLevel uint32) cfgEnum.ErrorCode {
	playerHero := this.getHero(uSn)
	if playerHero == nil {
		return plog.Print(head.Id, cfgEnum.ErrorCode_HeroSnNotFound, uSn, uCurLevel)
	}

	cfgAwaken := cfgData.GetCfgHeroAwaken(playerHero.AwakenLevel)
	if cfgAwaken == nil {
		return plog.Print(head.Id, cfgData.GetHeroAwakenErrorCode(playerHero.AwakenLevel), uSn, uCurLevel, playerHero.AwakenLevel)
	}

	if cfgData.GetCfgHeroAwaken(playerHero.AwakenLevel+1) == nil {
		return plog.Print(head.Id, cfgEnum.ErrorCode_MaxLevel, uSn, uCurLevel, playerHero.AwakenLevel)
	}

	if playerHero.AwakenLevel != uCurLevel {
		return plog.Print(head.Id, cfgEnum.ErrorCode_HeroAwakenLevelParam, uSn, uCurLevel, playerHero.AwakenLevel)
	}

	if playerHero.Star < cfgAwaken.NeedStar {
		return plog.Print(head.Id, cfgEnum.ErrorCode_HeroAwakenStarNotEnough, uSn, uCurLevel, playerHero.Star, cfgAwaken.NeedStar)

	}

	cfgHero := cfgData.GetCfgHero(playerHero.Id)
	if cfgHero == nil {
		return plog.Print(head.Id, cfgData.GetHeroErrorCode(playerHero.Id), uSn, uCurLevel, playerHero.Id)
	}

	//扣道具
	arrItem := make([]*common.ItemInfo, 0)
	uCurItemCount := this.getPlayerBagFun().GetItemCount(uint32(cfgEnum.ESystemType_Item), cfgHero.AwakenItemId)
	if uCurItemCount < int64(cfgAwaken.NeedItemCount) {
		uReplaceItemCount := this.getPlayerBagFun().GetItemCount(uint32(cfgEnum.ESystemType_Item), cfgAwaken.ReplaceItemId)
		if uCurItemCount+uReplaceItemCount < int64(cfgAwaken.NeedItemCount) {
			return plog.Print(head.Id, cfgEnum.ErrorCode_ItemNotEnough, uSn, uCurLevel, cfgHero.AwakenItemId, uCurItemCount, uReplaceItemCount, cfgAwaken.NeedItemCount)
		}

		if uCurItemCount > 0 {
			arrItem = append(arrItem, &common.ItemInfo{Kind: uint32(cfgEnum.ESystemType_Item), Id: cfgHero.AwakenItemId, Count: uCurItemCount})
		}

		arrItem = append(arrItem, &common.ItemInfo{Kind: uint32(cfgEnum.ESystemType_Item), Id: cfgAwaken.ReplaceItemId, Count: int64(cfgAwaken.NeedItemCount) - uCurItemCount})
	} else {
		arrItem = append(arrItem, &common.ItemInfo{Kind: uint32(cfgEnum.ESystemType_Item), Id: cfgHero.AwakenItemId, Count: int64(cfgAwaken.NeedItemCount)})
	}

	uCode := this.getPlayerBagFun().DelArrItem(head, arrItem, pb.EmDoingType_EDT_HeroAwaken)
	if uCode != cfgEnum.ErrorCode_Success {
		return plog.Print(head.Id, uCode, uSn, uCurLevel, *arrItem[0])
	}

	playerHero.AwakenLevel++

	playerHero = this.updateHero(playerHero, playerHero.Star)

	cluster.SendToClient(head, &pb.HeroAwakenLevelResponse{
		PacketHead: &pb.IPacket{},
		Sn:         uSn,
		CurLevel:   playerHero.AwakenLevel,
	}, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}

// 英雄上阵请求
func (this *PlayerHeroFun) HeroGameHeroListRequest(head *pb.RpcHead, pbRequest *pb.HeroGameHeroListRequest) {
	uCode := this.HeroGameHeroList(head, pbRequest.Team)
	cluster.SendToClient(head, &pb.HeroGameHeroListResponse{
		PacketHead: &pb.IPacket{},
		Team:       this.mapHeroTeam[pbRequest.Team.TeamType],
	}, uCode)
}

// 是否获取过这个伙伴
func (this *PlayerHeroFun) IsGetHeroId(uId uint32) bool {
	_, ok := this.mapHeroBookProp[uId]
	return ok
}

func (this *PlayerHeroFun) getHero(uSn uint32) *PlayerHero {
	info, ok := this.mapHero[uSn]
	if !ok {
		return nil
	}

	return info
}
func (this *PlayerHeroFun) updateHero(playerHero *PlayerHero, uOldStar uint32) *PlayerHero {
	if playerHero == nil {
		return nil
	}

	this.mapHero[playerHero.Sn] = playerHero

	//更新数据
	if uOldStar != playerHero.Star {
		//删除老的
		if uOldStar > 0 {
			delete(this.mapStarHero[uOldStar][playerHero.Id], playerHero.Sn)
		}

		if _, ok := this.mapStarHero[playerHero.Star]; !ok {
			this.mapStarHero[playerHero.Star] = make(map[uint32]map[uint32]*PlayerHero)
		}

		if _, ok := this.mapStarHero[playerHero.Star][playerHero.Id]; !ok {
			this.mapStarHero[playerHero.Star][playerHero.Id] = make(map[uint32]*PlayerHero)
		}

		this.mapStarHero[playerHero.Star][playerHero.Id][playerHero.Sn] = playerHero
	}

	//如果是出站中的 需要更新数据
	if this.IsFight(playerHero.Sn) {
		this.updateCalcFightpower(false)
	}

	this.UpdateSave(true)
	return playerHero
}

// 删除英雄
func (this *PlayerHeroFun) DelHeroList(arrList []uint32) bool {
	for _, sn := range arrList {
		this.delHero(sn)
	}

	return true
}

// 删除英雄
func (this *PlayerHeroFun) delHero(sn uint32) bool {
	playerHero := this.getHero(sn)
	if playerHero == nil {
		return false
	}

	delete(this.mapStarHero[playerHero.Star][playerHero.Id], sn)
	delete(this.mapHero, sn)
	this.UpdateSave(true)

	return true
}

// 上阵星星变更接口
func (this *PlayerHeroFun) HeroBattleStarChangeRequest(head *pb.RpcHead, req, rsp proto.Message) cfgEnum.ErrorCode {
	request := req.(*pb.HeroBattleStarChangeRequest)
	if len(request.Heros) <= 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_HeroBattleStarChangeParam, *request)
	}
	plog.Info("head: %v, req: %v", head, request)
	// 上阵列表
	mapHeros := map[uint32]*PlayerHero{}
	arrHeroSn := this.GetTeamList(uint32(cfgEnum.EHeroTeam_BattleNormal))
	for _, sn := range arrHeroSn {
		if sn > 0 && sn != 0xffffffff {
			pHero := this.getHero(sn)
			if pHero == nil {
				return plog.Print(this.AccountId, cfgEnum.ErrorCode_HeroNotFound, sn)
			}
			mapHeros[sn] = pHero
		}
	}
	// 所有英雄一定存在
	for _, item := range request.Heros {
		if _, ok := mapHeros[item.HeroID]; !ok {
			return plog.Print(this.AccountId, cfgEnum.ErrorCode_HeroNotFound, item.HeroID)
		}
	}
	// 计算上阵星星变更
	uTotalStar := uint32(0)
	notifys := &pb.HeroBattleStarChangeNotify{PacketHead: &pb.IPacket{}}
	for _, item := range request.Heros {
		hero := mapHeros[item.HeroID]
		// 是否变更
		if item.Total != hero.BattleStar {
			notifys.Heros = append(notifys.Heros, &pb.HeroBattleStarInfo{HeroID: item.HeroID, Total: item.Total})
		}
		// 设置英雄上阵星星数量
		hero.BattleStar = item.Total
		uTotalStar += item.Total
	}

	//成就统计
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_BattleStarLevel, uTotalStar)

	// 通知客户端(上阵星星)
	cluster.SendToClient(head, notifys, cfgEnum.ErrorCode_Success)

	this.UpdateSave(true)
	return cfgEnum.ErrorCode_Success
}

// 自动上阵
func (this *PlayerHeroFun) AutoHeroGameHeroList(head *pb.RpcHead, pHero *pb.PBHero) {
	if pHero == nil {
		return
	}
	arrHeroUnlock := cfgData.GetCfgGlobalArrConfig(cfgEnum.GlobalArrConfig_GLOBAL_TEAMPOS_UNLOCKMAP).Value
	//默认最多自动上阵到五个人
	uHeroLen := len(arrHeroUnlock) - 3

	mapid, stageid := this.getPlayerSystemBattleFun().GetFinishMapIdAndStageId(pb.EmBattleType_EBT_Hook)
	uBattleId := uint32(0)
	cfgBattle := cfgData.GetCfgBattleHookStage(mapid, stageid)
	if cfgBattle != nil {
		uBattleId = cfgBattle.Id
	}
	arrOldGameHeroList := this.GetTeamList(uint32(cfgEnum.EHeroTeam_BattleNormal))
	//判断此英雄ID是否相同
	for _, sn := range arrOldGameHeroList {
		if sn <= 0 || sn == uint32(0xffffffff) {
			continue
		}

		pTmpHero := this.getHero(sn)
		if pTmpHero == nil {
			continue
		}

		if sn == pHero.Sn || pTmpHero.Id == pHero.Id {
			return
		}
	}

	arrGameHeroList := make([]uint32, 0)
	arrGameHeroList = append(arrGameHeroList, arrOldGameHeroList...)
	bUpdate := false
	for i := 0; i < uHeroLen; i++ {
		if arrHeroUnlock[i] > uBattleId {
			break
		}

		if len(arrGameHeroList) < i+1 {
			arrGameHeroList = append(arrOldGameHeroList, 0)
		}

		if arrGameHeroList[i] <= 0 {
			arrGameHeroList[i] = pHero.Sn
			bUpdate = true
			break
		}
	}

	if bUpdate {
		this.HeroGameHeroList(head, &pb.PBHeroTeamList{HeroSn: arrGameHeroList, TeamType: uint32(cfgEnum.EHeroTeam_BattleNormal)})
		// 发送给客户端
		cluster.SendToClient(head, &pb.HeroGameHeroListNotify{
			PacketHead: &pb.IPacket{},
			Team:       this.mapHeroTeam[uint32(cfgEnum.EHeroTeam_BattleNormal)],
		}, cfgEnum.ErrorCode_Success)
	}

}

// 英雄上阵请求
func (this *PlayerHeroFun) HeroGameHeroList(head *pb.RpcHead, pbTeam *pb.PBHeroTeamList) cfgEnum.ErrorCode {
	// 上一次上阵信息
	mapOld := map[uint32]struct{}{}
	for _, sn := range this.GetTeamList(pbTeam.TeamType) {
		if sn > 0 {
			mapOld[sn] = struct{}{}
		}
	}

	notifys := &pb.HeroBattleStarChangeNotify{PacketHead: &pb.IPacket{}} // 上阵，下阵星星变更
	listChange := []uint32{}                                             // 下阵列表
	mapHeroType := map[uint32]struct{}{}                                 // 英雄类型重复检测
	mapHeroSn := map[uint32]struct{}{}                                   // 英雄重复上阵检测
	mapQualityCount := make(map[uint32]uint32)
	mapLevelCount := make(map[uint32]uint32)
	mapIdCount := make(map[uint32]uint32)
	for _, sn := range pbTeam.HeroSn {
		if sn <= 0 || sn == uint32(0xffffffff) {
			continue
		}
		pHero := this.getHero(sn)
		// 英雄不存在时，不能上阵
		if pHero == nil {
			return plog.Print(this.AccountId, cfgEnum.ErrorCode_HeroNotFound, sn)
		}
		// 同一种类英雄，不能重复上阵
		if _, ok := mapHeroType[pHero.Id]; ok {
			return plog.Print(this.AccountId, cfgEnum.ErrorCode_HeroIdRepeated, pHero.Id)
		} else {
			mapHeroType[pHero.Id] = struct{}{}
		}
		// 所有英雄，不能重复sn。每一个英雄的sn唯一
		if _, ok := mapHeroSn[pHero.Sn]; ok {
			return plog.Print(this.AccountId, cfgEnum.ErrorCode_HeroIdRepeated, pHero.Sn)
		} else {
			mapHeroSn[pHero.Sn] = struct{}{}
		}

		cfgHeroStar := cfgData.GetCfgHeroStar(pHero.Star)
		cfgHero := cfgData.GetCfgHero(pHero.Id)
		mapIdCount[pHero.Id] += 1
		mapQualityCount[cfgHeroStar.Quality] += 1
		mapLevelCount[this.getPlayerSystemProfessionFun().GetProfLevel(cfgHero.Prof, false)] += 1

		// 计算下阵、上阵列表
		if _, ok := mapOld[sn]; ok {
			delete(mapOld, sn)
		} else {
			if pbTeam.TeamType == uint32(cfgEnum.EHeroTeam_BattleNormal) {
				// 需要下阵，重置英雄的上阵积分
				pHero.BattleStar = 0
				notifys.Heros = append(notifys.Heros, &pb.HeroBattleStarInfo{HeroID: sn, Total: pHero.BattleStar})
				listChange = append(listChange, sn)
			}
		}
	}
	// 判单刷新上阵列表
	if len(mapHeroSn) <= 0 {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_PARAM, *pbTeam)
	}
	//计算下阵
	for k, _ := range mapOld {
		listChange = append(listChange, k)
		// 需要下阵，重置英雄的上阵积分
		if pHero := this.getHero(k); pHero != nil {
			pHero.BattleStar = 0
			notifys.Heros = append(notifys.Heros, &pb.HeroBattleStarInfo{HeroID: k, Total: pHero.BattleStar})
		}
	}
	// 下阵英雄
	plog.Info("下阵英雄sn列表: %v", listChange)

	// 更新上阵列表
	this.mapHeroTeam[pbTeam.TeamType] = pbTeam

	//更新上下阵的英雄属性
	if pbTeam.TeamType == uint32(cfgEnum.EHeroTeam_BattleNormal) {
		for _, sn := range listChange {
			playerHero := this.getHero(sn)
			playerHero = this.CalcHeroProp(playerHero)
			playerHero = this.updateHero(playerHero, playerHero.Star)
			this.updateCalcFightpower(false)
		}

		// 通知客户端(上阵星星)
		if len(notifys.Heros) > 0 {
			cluster.SendToClient(head, notifys, cfgEnum.ErrorCode_Success)
		}

		//更新成就
		this.getPlayerSystemTaskFun().TriggerAchieveGameFightList(cfgEnum.AchieveType_FightQualityHero, mapQualityCount)
		this.getPlayerSystemTaskFun().TriggerAchieveGameFightList(cfgEnum.AchieveType_FightHeroLevel, mapLevelCount)
		this.getPlayerSystemTaskFun().TriggerAchieveGameFightList(cfgEnum.AchieveType_FightHeroId, mapIdCount)
	}

	this.UpdateSave(true)
	return cfgEnum.ErrorCode_Success
}

// 更新上阵成就
func (this *PlayerHeroFun) OnHeroGameHeroList() {
	mapLevelCount := make(map[uint32]uint32)
	arrHeroSn := this.GetTeamList(uint32(cfgEnum.EHeroTeam_BattleNormal))
	for _, sn := range arrHeroSn {
		if sn <= 0 || sn == uint32(0xffffffff) {
			continue
		}

		pHero := this.getHero(sn)
		// 英雄不存在时，不能上阵
		if pHero == nil {
			continue
		}

		cfgHero := cfgData.GetCfgHero(pHero.Id)
		mapLevelCount[this.getPlayerSystemProfessionFun().GetProfLevel(cfgHero.Prof, false)] += 1
	}

	this.getPlayerSystemTaskFun().TriggerAchieveGameFightList(cfgEnum.AchieveType_FightHeroLevel, mapLevelCount)
}

func (this *PlayerHeroFun) GetTeamList(uTeamType uint32) []uint32 {
	arrReturn := make([]uint32, 0)
	if _, ok := this.mapHeroTeam[uTeamType]; !ok {
		return arrReturn
	}
	return this.mapHeroTeam[uTeamType].HeroSn
}

// 清除队伍
func (this *PlayerHeroFun) ClearTeamList(uTeamType uint32) {
	if _, ok := this.mapHeroTeam[uTeamType]; !ok {
		return
	}

	this.mapHeroTeam[uTeamType] = &pb.PBHeroTeamList{
		TeamType: uTeamType,
		HeroSn:   make([]uint32, 0),
	}

	cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, &pb.HeroGameHeroListResponse{
		PacketHead: &pb.IPacket{},
		Team:       this.mapHeroTeam[uTeamType],
	}, cfgEnum.ErrorCode_Success)

	this.UpdateSave(true)
}

func (this *PlayerHeroFun) IsFight(uSn uint32) bool {
	return base.ArrayContainsValue(this.GetTeamList(uint32(cfgEnum.EHeroTeam_BattleNormal)), uSn)
}

// 计算属性
func (this *PlayerHeroFun) CalcHeroProp(playerHero *PlayerHero) *PlayerHero {
	if playerHero == nil {
		return playerHero
	}
	//未出战斗 清空数据
	if !this.IsFight(playerHero.Sn) {
		playerHero.mapProp = common.NewProperty()
		return this.updateHero(playerHero, playerHero.Star)
	}
	cfgHero := cfgData.GetCfgHero(playerHero.Id)
	if cfgHero == nil {
		return playerHero
	}
	cfgHeroStar := cfgData.GetCfgHeroStar(playerHero.Star)
	if cfgHeroStar == nil {
		return playerHero
	}
	cfgHeroGrade := cfgData.GetCfgHeroGrade(cfgHero.Grade)
	if cfgHeroGrade == nil {
		return playerHero
	}
	//获取装备属性
	pPlayerProf := this.getPlayerSystemProfessionFun().GetProf(cfgHero.Prof)
	if pPlayerProf == nil {
		return playerHero
	}

	//gradeStar := float64(cfgHeroGrade.GrowRate*cfgHeroStar.GrowRate) / (base.MIL_PERCENT * base.MIL_PERCENT) // 稀有度系数*星级系数
	// 所有属性
	level := base.MinUint32(pPlayerProf.Level, cfgHeroStar.MaxLevel)                                       // 等级
	talentProp := cfgData.GetCfgTalentProperty(playerHero.Id, playerHero.Star)                             // 英雄天赋属性
	gradeStar := float64(cfgHeroGrade.GrowRate) / base.MIL_PERCENT                                         // 稀有度系数
	breakthroughGrow := cfgData.GetCfgGlobalConfig(cfgEnum.GlobalConfig_GLOBAL_CFG_BREAKTHROUGHLEVEL_GROW) // 突破成长比例
	breakRate := 1 + float64(pPlayerProf.Grade*breakthroughGrow)/base.MIL_PERCENT                          // 突破系数
	playerHero.mapProp = common.NewProperty()
	for k, v := range cfgHero.MapProf {
		uGrow, ok := cfgHero.MapProfGrow[k]
		if !ok {
			uGrow = 0
		}
		talent := talentProp.GetProperty(k)         // 天赋固定值
		talentRate := talentProp.GetPropertyRate(k) // 天赋百分比
		// 英雄属性 = (英雄基础属性 + (等级-1)*成长系数*(1+突破次数*突破成长比例) + 天赋固定值)*稀有度系数*星级系数*(1+天赋百分比)
		playerHero.mapProp.UpdateScorePropertyVal(k, (float64(v)+float64(level-1)*uGrow*breakRate+talent)*gradeStar*(cfgHeroStar.MapGrow[k]/base.MIL_PERCENT)*(1+talentRate/base.MIL_PERCENT))
	}
	//plog.Info("propertyType heroID: %d, sn: %d, grade: %d, value: %f, %f, %f, talent: %f,%f,%f, talentRate: %f,%f,%f", playerHero.Id, playerHero.Sn, pPlayerProf.Grade, playerHero.mapProp.GetProperty(6), playerHero.mapProp.GetProperty(7), playerHero.mapProp.GetProperty(8), talentProp.GetProperty(6), talentProp.GetProperty(7), talentProp.GetProperty(8), talentProp.GetPropertyRate(6), talentProp.GetPropertyRate(7), talentProp.GetPropertyRate(8))
	//plog.Info("propertyType heroID: %d, sn: %d, value: %f, %f, %f", playerHero.Id, playerHero.Sn, playerHero.mapProp.GetProperty(6), playerHero.mapProp.GetProperty(7), playerHero.mapProp.GetProperty(8))

	// 晶核天赋属性
	for k, val := range this.getPlayerCrystalFun().GetCrystalProp(cfgHero.Element) {
		plog.Trace("uid: %d, CrystalProperty prof: %d, element: %d, %d  %f", this.AccountId, cfgHero.Prof, cfgHero.Element, k, val)
		playerHero.mapProp.UpdateScorePropertyVal(k, val)
	}
	// 词条属性
	tmpEntry := entry.KeyValueToMap(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_HeroProp), cfgHero.Element)...)
	for k, val := range entry.KeyValueToMap(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_ProfessionProp), cfgHero.Prof)...) {
		tmpEntry[k] += val
	}
	for k, val := range tmpEntry {
		plog.Trace("uid: %d, EntryProperty prof: %d, element: %d, %d  %f", this.AccountId, cfgHero.Prof, cfgHero.Element, k, float64(val)/100)
		playerHero.mapProp.UpdateScorePropertyVal(k, float64(val)/100)
	}
	// 职业属性(值还是不对)
	for k, v := range pPlayerProf.mapProp {
		plog.Trace("uid: %d, ProfessionProperty prof: %d, element: %d, %d  %f", this.AccountId, cfgHero.Prof, cfgHero.Element, k, float64(v))
		playerHero.mapProp.UpdateScorePropertyVal(k, float64(v))
	}
	// 加图鉴固定值
	for k, v := range this.mapHeroBookProp {
		plog.Trace("uid: %d, BookProperty prof: %d, element: %d, %d  %f", this.AccountId, cfgHero.Prof, cfgHero.Element, k, float64(v))
		playerHero.mapProp.UpdateScorePropertyVal(k, float64(v))
	}
	//加星源属性
	for k, v := range this.getPlayerSystemHookTechFun().GetHookTechEffect(cfgEnum.TechEffectType_HeroProp) {
		plog.Trace("uid: %d, HookTechProperty prof: %d, element: %d, %d  %f", this.AccountId, cfgHero.Prof, cfgHero.Element, k, float64(v))
		playerHero.mapProp.UpdateScorePropertyVal(k, float64(v))
	}
	// 打印日志
	//playerHero.mapProp.Print(cfgHero.Prof, cfgHero.Element)
	//转战斗力
	playerHero.uScore = uint32(cfgData.GetRankScoreCfgProfScore(playerHero.mapProp, cfgEnum.EHydraProf(cfgHero.Prof), cfgEnum.EHydraElementType(cfgHero.Element)))
	plog.Trace("uid: %d, propertyScore prof: %d, element: %d, score: %d", this.AccountId, cfgHero.Prof, cfgHero.Element, playerHero.uScore)
	playerHero.bCalc = false
	return playerHero
}

// 计算战斗力
func (this *PlayerHeroFun) UpdateProfFightPower(uProfType uint32) {
	arrHeroSn := this.GetTeamList(uint32(cfgEnum.EHeroTeam_BattleNormal))
	for _, sn := range arrHeroSn {
		if sn <= 0 {
			continue
		}

		playerHero := this.getHero(sn)
		if playerHero == nil {
			continue
		}

		cfgHero := cfgData.GetCfgHero(playerHero.Id)
		if cfgHero.Prof != uProfType || playerHero.bCalc {
			continue
		}

		playerHero.bCalc = true

		this.updateHero(playerHero, playerHero.Star)
	}
}
func (this *PlayerHeroFun) UpdateElementFightPower(uElementType uint32) {
	arrHeroSn := this.GetTeamList(uint32(cfgEnum.EHeroTeam_BattleNormal))

	for _, sn := range arrHeroSn {
		if sn <= 0 {
			continue
		}

		playerHero := this.getHero(sn)
		if playerHero == nil {
			continue
		}

		cfgHero := cfgData.GetCfgHero(playerHero.Id)
		if cfgHero.Element != uElementType || playerHero.bCalc {
			continue
		}

		playerHero.bCalc = true

		this.updateHero(playerHero, playerHero.Star)
	}
}

// 计算战斗力
func (this *PlayerHeroFun) CalcFightPower() {
	uOldFightPower := this.uFightPower
	this.uFightPower = 0
	arrHeroSn := this.GetTeamList(uint32(cfgEnum.EHeroTeam_BattleNormal))
	for _, sn := range arrHeroSn {
		if sn <= 0 {
			continue
		}

		playerHero := this.getHero(sn)
		if playerHero == nil {
			continue
		}

		if !playerHero.bCalc {
			this.uFightPower += playerHero.uScore
			continue
		}

		playerHero = this.CalcHeroProp(playerHero)

		this.updateHero(playerHero, playerHero.Star)

		this.uFightPower += playerHero.uScore
	}

	this.bCalcFightpower = false

	//通知客户端
	if uOldFightPower != this.uFightPower {
		cluster.SendToClient(&pb.RpcHead{Id: this.AccountId}, &pb.HeroFightPowerNotify{
			PacketHead: &pb.IPacket{},
			FightPower: this.uFightPower,
		}, cfgEnum.ErrorCode_Success)
	}
	// 触发战斗力提升词条
	if this.uFightPower > this.uMaxHistoryFightPower {
		diff := this.uFightPower - this.uMaxHistoryFightPower
		this.uMaxHistoryFightPower = this.uFightPower

		if cfg := cfgData.GetCfgRankInfoConfig(uint32(cfgEnum.ERankType_ChampionshipPower)); cfg != nil {
			createTime := this.getPlayerBaseFun().GetServerStartTime()
			this.UpdateRank(
				cfgData.GetCfgRankActiveTime(cfg, createTime),
				uint32(cfgEnum.ERankType_ChampionshipPower),
				uint64(this.uMaxHistoryFightPower),
			)
		}
		//成就触发
		this.getPlayerSystemTaskFun().TriggerAchieve(&pb.RpcHead{}, cfgEnum.AchieveType_IncreasePower, diff)
		plog.Trace("MaxHistoryFightPower calc id:%d power:%d diff:%d total:%d", this.AccountId, this.uFightPower, diff,
			this.getPlayerSystemTaskFun().GetAchieveValue(uint32(cfgEnum.AchieveType_IncreasePower)))
		this.UpdateSave(true)
	}
}

func (this *PlayerHeroFun) GetMaxHistoryFightPower() uint64 {
	return uint64(this.uMaxHistoryFightPower)
}

// 更新英雄
func (this *PlayerHeroFun) DipUpdateHero(head *pb.RpcHead, pbRequest *pb.PBHero) cfgEnum.ErrorCode {
	if pbRequest == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_GmParam)
	}

	pHero := this.getHero(pbRequest.Sn)
	if pHero == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_HeroSnNotFound, pbRequest.Sn)
	}

	if pHero.Id != pbRequest.Id {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_HeroIdInconsistent, pHero.Id, pbRequest.Id)
	}

	uOldStar := pHero.Star
	*pHero.PBHero = *pbRequest

	this.updateHero(pHero, uOldStar)
	return cfgEnum.ErrorCode_Success
}

// 更新英雄
func (this *PlayerHeroFun) DipDelHero(head *pb.RpcHead, sn uint32) cfgEnum.ErrorCode {
	pHero := this.getHero(sn)
	if pHero == nil {
		return plog.Print(this.AccountId, cfgEnum.ErrorCode_HeroSnNotFound, sn)
	}

	//下阵
	arrHeroSn := this.GetTeamList(uint32(cfgEnum.EHeroTeam_BattleNormal))
	for index, v := range arrHeroSn {
		if v == sn {
			this.mapHeroTeam[uint32(cfgEnum.EHeroTeam_BattleNormal)].HeroSn[index] = 0
			this.CalcFightPower()
		}
	}

	this.delHero(sn)

	return cfgEnum.ErrorCode_Success
}

// 英雄重生请求
func (this *PlayerHeroFun) HeroRebirthRequest(head *pb.RpcHead, pbRequest *pb.HeroRebirthRequest) {
	uCode := this.HeroRebirth(head, pbRequest.Sn)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.HeroRebirthResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 英雄重生请求
func (this *PlayerHeroFun) HeroRebirth(head *pb.RpcHead, uSn uint32) cfgEnum.ErrorCode {
	playerHero := this.getHero(uSn)
	if playerHero == nil {
		return plog.Print(head.Id, cfgEnum.ErrorCode_HeroSnNotFound, uSn)
	}

	cfgHero := cfgData.GetCfgHero(playerHero.Id)
	if cfgHero == nil {
		return plog.Print(head.Id, cfgData.GetHeroErrorCode(playerHero.Id), uSn)
	}

	cfgHeroStar := cfgData.GetCfgHeroStar(playerHero.Star)
	if cfgHeroStar == nil {
		return plog.Print(head.Id, cfgData.GetHeroStarErrorCode(playerHero.Star), uSn)
	}

	if len(cfgHeroStar.RebirthHero) <= 0 {
		return plog.Print(head.Id, cfgEnum.ErrorCode_RebirthHeroEmpty, uSn)
	}

	if cfgHeroStar.RebirthHeroCount <= 0 {
		return plog.Print(head.Id, cfgEnum.ErrorCode_PARAM, uSn)
	}

	//降低星级
	playerHero.Star = cfgHeroStar.RebirthStar
	this.updateHero(playerHero, cfgHeroStar.Id)

	//给狗粮
	this.AddHeros(head, cfgHeroStar.RebirthHero[cfgHero.Prof], cfgHeroStar.RebirthStar, cfgHeroStar.RebirthHeroCount, pb.EmDoingType_EDT_HeroRebirth, true)

	//弹窗恭喜获得 需要带上本体和重置的
	arrShowItems := make([]*common.ItemInfo, 0)
	arrShowItems = append(arrShowItems, &common.ItemInfo{
		Kind:   uint32(cfgEnum.ESystemType_Hero),
		Id:     playerHero.Id,
		Count:  1,
		Params: []uint32{cfgHeroStar.RebirthStar},
	})

	arrShowItems = append(arrShowItems, &common.ItemInfo{
		Kind:   uint32(cfgEnum.ESystemType_Hero),
		Id:     cfgHeroStar.RebirthHero[cfgHero.Prof],
		Count:  int64(cfgHeroStar.RebirthHeroCount),
		Params: []uint32{cfgHeroStar.RebirthStar},
	})
	arrPbItems := this.getPlayerBagFun().GetPbItems(arrShowItems, pb.EmDoingType_EDT_HeroRebirth)
	this.getPlayerBagFun().CommonPrizeNotify(head, arrPbItems, pb.EmDoingType_EDT_HeroRebirth)

	cluster.SendToClient(head, &pb.HeroRebirthResponse{
		PacketHead: &pb.IPacket{},
		Sn:         uSn,
		Star:       playerHero.Star,
	}, cfgEnum.ErrorCode_Success)

	return cfgEnum.ErrorCode_Success
}
func (this *PlayerHeroFun) Heat() {
	if this.bCalcFightpower {
		this.CalcFightPower()
	}
}
func (this *PlayerHeroFun) InGameListInfo(uHeroId uint32) bool {
	arrHeroSn := this.GetTeamList(uint32(cfgEnum.EHeroTeam_BattleNormal))
	for _, sn := range arrHeroSn {
		if sn <= 0 {
			continue
		}
		pPlayerHero := this.getHero(sn)
		if pPlayerHero == nil {
			continue
		}
		if pPlayerHero.Id == uHeroId {
			return true
		}
	}

	return false
}

func (this *PlayerHeroFun) GetBattleListInfo(heroIds []uint32) []*pb.PBHero {
	arrList := make([]*pb.PBHero, 0)
	if len(heroIds) <= 0 {
		return arrList
	}

	arrHeroSn := this.GetTeamList(uint32(cfgEnum.EHeroTeam_BattleNormal))
	for _, sn := range arrHeroSn {
		if sn <= 0 {
			continue
		}
		pPlayerHero := this.getHero(sn)
		if pPlayerHero == nil {
			continue
		}

		if !base.ArrayContainsValue(heroIds, pPlayerHero.Id) {
			continue
		}

		arrList = append(arrList, pPlayerHero.PBHero)
	}
	return arrList
}
func (this *PlayerHeroFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerHero{}
}

// 设置玩家数据
func (this *PlayerHeroFun) SetUserTypeInfo(pbData proto.Message) bool {
	if pbData == nil {
		return false
	}

	pbSystem := pbData.(*pb.PBPlayerHero)
	if pbSystem == nil {
		return false
	}

	this.loadData(pbSystem)

	this.CalcFightPower()
	return true
}

// 英雄图鉴激活请求
func (this *PlayerHeroFun) HeroBookActiveRequest(head *pb.RpcHead, pbRequest *pb.HeroBookActiveRequest) {
	uCode := this.HeroBookActive(head, pbRequest.HeroId)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.HeroBookActiveResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 英雄图鉴激活请求
func (this *PlayerHeroFun) HeroBookActive(head *pb.RpcHead, heroId uint32) cfgEnum.ErrorCode {
	bookInfo, ok := this.mapHeroBook[heroId]
	if !ok {
		return plog.Print(head.Id, cfgEnum.ErrorCode_HeroNoGet, heroId)
	}

	if bookInfo.Star > 0 {
		return plog.Print(head.Id, cfgEnum.ErrorCode_HeroBookHaveActive, heroId)
	}

	bookInfo.Star = 1
	this.CalcHeroBookProp(head)
	//给奖励
	this.getPlayerBagFun().AddItem(head, uint32(cfgEnum.ESystemType_Item), uint32(pb.EmItemExpendType_EIET_Cash), int64(cfgData.GetCfgGlobalConfig(cfgEnum.GlobalConfig_GLOBAL_HERO_BOOK_DIAMOND)), pb.EmDoingType_EDT_HeroBook, true)

	//成就统计
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_HeroBookAcivate, 1)

	this.UpdateSave(true)
	// 发送给客户端
	cluster.SendToClient(head, &pb.HeroBookActiveResponse{
		PacketHead: &pb.IPacket{},
		HeroBook:   bookInfo,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

// 英雄图鉴激活请求
func (this *PlayerHeroFun) HeroBookUpStarRequest(head *pb.RpcHead, pbRequest *pb.HeroBookUpStarRequest) {
	uCode := this.HeroBookUpStar(head, pbRequest.HeroId)
	if uCode != cfgEnum.ErrorCode_Success {
		cluster.SendToClient(head, &pb.HeroBookUpStarResponse{
			PacketHead: &pb.IPacket{},
		}, uCode)
	}
}

// 英雄图鉴激活请求
func (this *PlayerHeroFun) HeroBookUpStar(head *pb.RpcHead, heroId uint32) cfgEnum.ErrorCode {
	bookInfo, ok := this.mapHeroBook[heroId]
	if !ok {
		return plog.Print(head.Id, cfgEnum.ErrorCode_HeroNoGet, heroId)
	}

	if bookInfo.Star == 0 || bookInfo.Star >= bookInfo.MaxStar {
		return plog.Print(head.Id, cfgEnum.ErrorCode_HeroBookHaveActive, heroId)
	}

	bookInfo.Star++

	this.CalcHeroBookProp(head)

	this.UpdateSave(true)

	// 发送给客户端
	cluster.SendToClient(head, &pb.HeroBookUpStarResponse{
		PacketHead: &pb.IPacket{},
		HeroBook:   bookInfo,
	}, cfgEnum.ErrorCode_Success)
	return cfgEnum.ErrorCode_Success
}

func (this *PlayerHeroFun) CalcHeroBookProp(head *pb.RpcHead) {
	this.mapHeroBookProp = make(map[uint32]uint32)
	for _, info := range this.mapHeroBook {
		if info.Star <= 0 {
			continue
		}

		cfgStar := cfgData.GetCfgHeroStar(info.Star)
		this.mapHeroBookProp = base.MergeMapU32U32(this.mapHeroBookProp, cfgStar.MapBookProp)
	}

	//this.SystemPropNotify(head, pb.EmSyetemPropType_SPT_HeroBook, this.mapHeroBookProp)

	//更新战斗力
	this.updateCalcFightpower(true)
}

// 属性返回
func (this *PlayerHeroFun) SystemPropNotify(head *pb.RpcHead, propType pb.EmSyetemPropType, mapPropInfo map[uint32]uint32) {
	pbNotify := &pb.SystemPropNotify{
		PacketHead:     &pb.IPacket{},
		SyetemPropType: propType}

	for key, value := range mapPropInfo {
		pbNotify.PropInfoList = append(pbNotify.PropInfoList, &pb.PBPropInfo{PropId: key, Value: value})
	}

	// 发送给客户端
	cluster.SendToClient(head, pbNotify, cfgEnum.ErrorCode_Success)
}

// 万能碎片转换
func (this *PlayerHeroFun) HeroPieceRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.HeroPieceRequest)
	// 生成英雄碎片
	cfg := cfgData.GetCfgHero(req.HeroID)
	if cfg == nil {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_HeroNoGet, "head: %v, req: %v", head, req)
	}
	// 扣除万能碎片
	bag := this.getPlayerBagFun()
	errorCode := bag.DelItem(head, uint32(cfgEnum.ESystemType_Item), req.PieceID, int64(req.Count), pb.EmDoingType_EDT_HeroPiece)
	if errorCode != cfgEnum.ErrorCode_Success {
		return uerror.NewUErrorf(1, errorCode, "head: %v, req: %v", head, req)
	}
	// 生成英雄碎片
	errorCode = bag.AddItem(head, uint32(cfgEnum.ESystemType_Item), cfg.AwakenItemId, int64(req.Count), pb.EmDoingType_EDT_HeroPiece, true)
	if errorCode != cfgEnum.ErrorCode_Success {
		return uerror.NewUErrorf(1, errorCode, "head: %v, req: %v", head, req)
	}
	return nil
}

// 英雄升星请求(需求调整, 不消耗英雄本体了,改成消耗英雄碎片)
func (this *PlayerHeroFun) HeroUpStarRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.HeroUpStarRequest)
	rsp := response.(*pb.HeroUpStarResponse)
	// 判断英雄是否存在
	hero := this.getHero(req.Sn)
	if hero == nil {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_HeroSnNotFound, "head: %v, req: %v", head, req)
	}
	//判断是否最大星级
	if cfgData.GetCfgHeroStar(hero.Star+1) == nil {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_MaxLevel, "head: %v, req: %v, hero: %v", head, req, hero)
	}
	cfgStar := cfgData.GetCfgHeroStar(hero.Star)
	if cfgStar == nil {
		return uerror.NewUErrorf(1, cfgData.GetHeroStarErrorCode(hero.Star), "head: %v, req: %v, hero: %v", head, req, hero)
	}
	// 加载英雄配置
	heroCfg := cfgData.GetCfgHero(hero.Id)
	if heroCfg == nil {
		return uerror.NewUErrorf(1, cfgData.GetHeroErrorCode(hero.Id), "head: %v, req: %v, hero: %v", head, req, hero)
	}
	/*
		// 词条效果减免
		isReduce := false
		pbReduceItem := &pb.PBAddItemData{}
		if prob := entry.ToValue(this.getEntry().Get(uint32(cfgEnum.EntryEffectType_HeroUpgradeReduce), cfgStar.Quality)...); prob > 0 {
			isReduce = base.IsRadio(prob)
		}
	*/
	// 扣除英雄碎片
	errCode := this.getPlayerBagFun().DelItem(head, uint32(cfgEnum.ESystemType_Item), heroCfg.AwakenItemId, int64(cfgStar.NeedItemCnt), pb.EmDoingType_EDT_HeroUpStar)
	if errCode != cfgEnum.ErrorCode_Success {
		return uerror.NewUErrorf(1, errCode, "head: %v, req: %v, itemID: %d, itemCount: %d", head, req, heroCfg.AwakenItemId, cfgStar.NeedItemCnt)
	}
	// 提升英雄星级
	hero = this.innerUpdateHeroStar(head, hero)
	//成就触发
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_HeroStarUpgrade, 1, uint32(cfgEnum.EQuality_Any))
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_HeroStarUpgrade, 1, cfgStar.Quality)
	/*
		//恭喜获得
		if pbReduceItem.Id > 0 {
			this.getPlayerBagFun().CommonPrizeNotify(head, []*pb.PBAddItemData{pbReduceItem}, pb.EmDoingType_EDT_Entry)
		}
	*/
	rsp.Sn = hero.PBHero.Sn
	rsp.Star = hero.PBHero.Star
	return nil
}

type autoUpStar struct {
	hero    *PlayerHero
	cfgStar *cfgData.HeroStarCfg
	cfgHero *cfgData.HeroCfg
}

// 英雄自动升星请求
func (this *PlayerHeroFun) HeroAutoUpStarRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.HeroAutoUpStarRequest)
	rsp := response.(*pb.HeroAutoUpStarResponse)
	//白色=》绿色=》紫色
	list := []*autoUpStar{}
	items := []*common.ItemInfo{}
	for star := uint32(1); star <= 3; star++ {
		// 加载星级配置
		cfgStar := cfgData.GetCfgHeroStar(star)
		if cfgStar == nil {
			return uerror.NewUErrorf(1, cfgData.GetHeroStarErrorCode(star), "head: %v, req: %v, star: %d", head, req, star)
		}
		// key：星级 key2:id key3:sn
		for id, vals := range this.mapStarHero[star] {
			cfgHero := cfgData.GetCfgHero(id)
			if cfgHero == nil {
				return uerror.NewUErrorf(1, cfgData.GetHeroErrorCode(id), "head: %v, req: %v, id: %d", head, req, id)
			}
			for _, hero := range vals {
				// 判断是否有足够的英雄碎片数量
				count := this.getPlayerBagFun().GetItemCount(uint32(cfgEnum.ESystemType_Item), cfgHero.AwakenItemId)
				if count < int64(cfgStar.NeedItemCnt) {
					continue
				}
				list = append(list, &autoUpStar{hero: hero, cfgStar: cfgStar, cfgHero: cfgHero})
				items = append(items, &common.ItemInfo{
					Kind:  uint32(cfgEnum.ESystemType_Item),
					Id:    cfgHero.AwakenItemId,
					Count: int64(cfgStar.NeedItemCnt),
				})
			}
		}
	}
	if len(list) <= 0 {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_NoHeroAutoUpstar, "head: %v, req: %v", head, req)
	}
	// 扣除英雄碎片
	if errCode := this.getPlayerBagFun().DelArrItem(head, items, pb.EmDoingType_EDT_HeroAutoUpStar); errCode != cfgEnum.ErrorCode_Success {
		return uerror.NewUErrorf(1, errCode, "head: %v, req: %v, items: %v", head, req, items)
	}
	// 触发升级成就
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_HeroStarUpgrade, uint32(len(list)), uint32(cfgEnum.EQuality_Any))
	// 英雄升星
	for _, val := range list {
		val.hero = this.innerUpdateHeroStar(head, val.hero)
		rsp.HeroList = append(rsp.HeroList, &pb.PBU32U32{Key: val.hero.Sn, Value: val.hero.Star})
		this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_HeroStarUpgrade, 1, val.cfgStar.Quality)
	}
	return nil
}

// 英雄自动升星请求
func (this *PlayerHeroFun) innerUpdateHeroStar(head *pb.RpcHead, playerHero *PlayerHero) *PlayerHero {
	uOldStar := playerHero.Star
	playerHero.Star++
	this.mapUpStarCount[playerHero.Id]++
	//更新图鉴
	if _, ok := this.mapHeroBook[playerHero.Id]; ok {
		if playerHero.Star > this.mapHeroBook[playerHero.Id].MaxStar {
			this.mapHeroBook[playerHero.Id].MaxStar = playerHero.Star

			cluster.SendToClient(head, &pb.HeroBookNotify{
				PacketHead: &pb.IPacket{},
				HeroBook:   this.mapHeroBook[playerHero.Id],
			}, cfgEnum.ErrorCode_Success)
		}
	}
	//算属性
	playerHero = this.CalcHeroProp(playerHero)

	//升星获得英雄成就
	cfgStar := cfgData.GetCfgHeroStar(playerHero.Star)
	if cfgStar != nil && cfgStar.Quality > cfgData.GetCfgHeroStar(uOldStar).Quality {
		this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_AddQualityHero, 1, cfgStar.Quality)
	}
	return this.updateHero(playerHero, uOldStar)
}
