package httpkit

type AchieveTaskInfoNotify struct {
	PacketHead *IPacket           `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SystemType uint32             `protobuf:"varint,2,opt,name=SystemType,proto3" json:"SystemType"`
	TaskList   []*PBTaskStageInfo `protobuf:"bytes,3,rep,name=TaskList,proto3" json:"TaskList"`
}

type ActivityDataNewNotify struct {
	PacketHead   *IPacket                `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ActivityType uint32                  `protobuf:"varint,2,opt,name=ActivityType,proto3" json:"ActivityType"`
	Info         *PBPlayerSystemActivity `protobuf:"bytes,3,opt,name=Info,proto3" json:"Info"`
}

type ActivityFreePrizeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
}

type ActivityFreePrizeResponse struct {
	PacketHead         *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id                 uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	NextDailyPrizeTime uint64   `protobuf:"varint,4,opt,name=NextDailyPrizeTime,proto3" json:"NextDailyPrizeTime"`
}

type ActivityListNotify struct {
	PacketHead   *IPacket                `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ActivityList []*PBPlayerActivityInfo `protobuf:"bytes,2,rep,name=ActivityList,proto3" json:"ActivityList"`
	DelIdList    []uint32                `protobuf:"varint,3,rep,packed,name=DelIdList,proto3" json:"DelIdList"`
}

type ActivityOpenRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
}

type ActivityOpenResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	JsonData   []string `protobuf:"bytes,3,rep,name=JsonData,proto3" json:"JsonData"`
}

type ActivityRedNotify struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	IdList     []uint32 `protobuf:"varint,2,rep,packed,name=IdList,proto3" json:"IdList"`
}

type AdventureRewardRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ID         uint32   `protobuf:"varint,2,opt,name=ID,proto3" json:"ID"`
	CfgID      uint32   `protobuf:"varint,3,opt,name=CfgID,proto3" json:"CfgID"`
}

type AdventureRewardResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type AdvertNotify struct {
	PacketHead *IPacket        `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	AdvertList []*PBAdvertInfo `protobuf:"bytes,2,rep,name=AdvertList,proto3" json:"AdvertList"`
}

type AdvertRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	AdvestType uint32   `protobuf:"varint,2,opt,name=AdvestType,proto3" json:"AdvestType"`
}

type AdvertResponse struct {
	PacketHead *IPacket      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	AdvestInfo *PBAdvertInfo `protobuf:"bytes,2,opt,name=AdvestInfo,proto3" json:"AdvestInfo"`
}

type AllPlayerInfoNotify struct {
	PacketHead *IPacket      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Mark       uint32        `protobuf:"varint,2,opt,name=Mark,proto3" json:"Mark"`
	PlayerData *PBPlayerData `protobuf:"bytes,3,opt,name=PlayerData,proto3" json:"PlayerData"`
}

type AvatarFrameNotify struct {
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Frames     []*PBAvatarFrame `protobuf:"bytes,2,rep,name=Frames,proto3" json:"Frames"`
}

type AvatarNotify struct {
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Avatars    []*PBAvatar `protobuf:"bytes,2,rep,name=Avatars,proto3" json:"Avatars"`
}

type AwardMailRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	MailId     uint32   `protobuf:"varint,2,opt,name=MailId,proto3" json:"MailId"`
}

type AwardMailResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Mail       *PBMail  `protobuf:"bytes,2,opt,name=Mail,proto3" json:"Mail"`
}

type BPAcitiveNotify struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	BPType     uint32   `protobuf:"varint,2,opt,name=BPType,proto3" json:"BPType"`
	StageId    uint32   `protobuf:"varint,3,opt,name=StageId,proto3" json:"StageId"`
	ChargeTime uint64   `protobuf:"varint,4,opt,name=ChargeTime,proto3" json:"ChargeTime"`
}

type BPNewNotify struct {
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	BPInfo     *PBBPInfo `protobuf:"bytes,2,opt,name=BPInfo,proto3" json:"BPInfo"`
}

type BPNewStageNotify struct {
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	BPType     uint32           `protobuf:"varint,2,opt,name=BPType,proto3" json:"BPType"`
	StageList  []*PBBPStageInfo `protobuf:"bytes,3,rep,name=StageList,proto3" json:"StageList"`
	DelList    []uint32         `protobuf:"varint,4,rep,packed,name=DelList,proto3" json:"DelList"`
}

type BPPrizeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	BPType     uint32   `protobuf:"varint,2,opt,name=BPType,proto3" json:"BPType"`
	StageId    uint32   `protobuf:"varint,3,opt,name=StageId,proto3" json:"StageId"`
}

type BPPrizeResponse struct {
	PacketHead    *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	BPType        uint32   `protobuf:"varint,2,opt,name=BPType,proto3" json:"BPType"`
	StageId       uint32   `protobuf:"varint,3,opt,name=StageId,proto3" json:"StageId"`
	NoramlPrizeId uint32   `protobuf:"varint,4,opt,name=NoramlPrizeId,proto3" json:"NoramlPrizeId"`
	ExtralPrizeId uint32   `protobuf:"varint,5,opt,name=ExtralPrizeId,proto3" json:"ExtralPrizeId"`
}

type BPValueNotify struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	BPType     uint32   `protobuf:"varint,2,opt,name=BPType,proto3" json:"BPType"`
	Value      uint32   `protobuf:"varint,3,opt,name=Value,proto3" json:"Value"`
}

type BattleBeginRequest struct {
	PacketHead *IPacket     `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	BattleType EmBattleType `protobuf:"varint,2,opt,name=BattleType,proto3,enum=common.EmBattleType" json:"BattleType"`
	MapId      uint32       `protobuf:"varint,3,opt,name=MapId,proto3" json:"MapId"`
	StageId    uint32       `protobuf:"varint,4,opt,name=StageId,proto3" json:"StageId"`
	Params     []uint32     `protobuf:"varint,5,rep,packed,name=Params,proto3" json:"Params"`
}

type BattleBeginResponse struct {
	PacketHead *IPacket     `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	BattleType EmBattleType `protobuf:"varint,2,opt,name=BattleType,proto3,enum=common.EmBattleType" json:"BattleType"`
	MapId      uint32       `protobuf:"varint,3,opt,name=MapId,proto3" json:"MapId"`
	StageId    uint32       `protobuf:"varint,4,opt,name=StageId,proto3" json:"StageId"`
	FightCount uint32       `protobuf:"varint,5,opt,name=FightCount,proto3" json:"FightCount"`
	Params     []uint32     `protobuf:"varint,6,rep,packed,name=Params,proto3" json:"Params"`
}

type BattleEndRequest struct {
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Battle     *BattleInfo `protobuf:"bytes,2,opt,name=Battle,proto3" json:"Battle"`
}

type BattleEndResponse struct {
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	BattleType EmBattleType     `protobuf:"varint,2,opt,name=BattleType,proto3,enum=common.EmBattleType" json:"BattleType"`
	MapId      uint32           `protobuf:"varint,3,opt,name=MapId,proto3" json:"MapId"`
	StageId    uint32           `protobuf:"varint,4,opt,name=StageId,proto3" json:"StageId"`
	ItemInfo   []*PBAddItemData `protobuf:"bytes,5,rep,name=ItemInfo,proto3" json:"ItemInfo"`
}

type BattleFunBuyRequest struct {
	PacketHead    *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	BattleFunType uint32   `protobuf:"varint,2,opt,name=BattleFunType,proto3" json:"BattleFunType"`
}

type BattleFunBuyResponse struct {
	PacketHead    *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	BattleFunType uint32   `protobuf:"varint,2,opt,name=BattleFunType,proto3" json:"BattleFunType"`
}

type BattleHookPassRewardRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	MapID      uint32   `protobuf:"varint,2,opt,name=MapID,proto3" json:"MapID"`
	StageID    uint32   `protobuf:"varint,3,opt,name=StageID,proto3" json:"StageID"`
}

type BattleHookPassRewardResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type BattleInfo struct {
	BattleType  EmBattleType             `protobuf:"varint,1,opt,name=BattleType,proto3,enum=common.EmBattleType" json:"BattleType"`
	MapId       uint32                   `protobuf:"varint,2,opt,name=MapId,proto3" json:"MapId"`
	StageId     uint32                   `protobuf:"varint,3,opt,name=StageId,proto3" json:"StageId"`
	IsSucc      uint32                   `protobuf:"varint,4,opt,name=IsSucc,proto3" json:"IsSucc"`
	UseTime     uint32                   `protobuf:"varint,5,opt,name=UseTime,proto3" json:"UseTime"`
	StageRate   uint32                   `protobuf:"varint,6,opt,name=StageRate,proto3" json:"StageRate"`
	MonsterInfo []*BattleKillMonsterInfo `protobuf:"bytes,7,rep,name=MonsterInfo,proto3" json:"MonsterInfo"`
	ClientData  *PBBattleClientData      `protobuf:"bytes,8,opt,name=ClientData,proto3" json:"ClientData"`
	TotalDamage uint64                   `protobuf:"varint,9,opt,name=TotalDamage,proto3" json:"TotalDamage"`
}

type BattleKillMonsterInfo struct {
	MonsterType uint32 `protobuf:"varint,1,opt,name=MonsterType,proto3" json:"MonsterType"`
	KillCount   uint32 `protobuf:"varint,2,opt,name=KillCount,proto3" json:"KillCount"`
	MaxCount    uint32 `protobuf:"varint,3,opt,name=MaxCount,proto3" json:"MaxCount"`
}

type BattleMapNotify struct {
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	BattleType EmBattleType     `protobuf:"varint,2,opt,name=BattleType,proto3,enum=common.EmBattleType" json:"BattleType"`
	MapInfo    *PBBattleMapInfo `protobuf:"bytes,3,opt,name=MapInfo,proto3" json:"MapInfo"`
}

type BattleNormalCardRequest struct {
	PacketHead      *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Stage           uint32   `protobuf:"varint,2,opt,name=Stage,proto3" json:"Stage"`
	CardID          uint32   `protobuf:"varint,3,opt,name=CardID,proto3" json:"CardID"`
	BattleBeginTime uint64   `protobuf:"varint,4,opt,name=BattleBeginTime,proto3" json:"BattleBeginTime"`
}

type BattleNormalCardResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type BattleRecordRequest struct {
	PacketHead *IPacket     `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	BattleType EmBattleType `protobuf:"varint,2,opt,name=BattleType,proto3,enum=common.EmBattleType" json:"BattleType"`
	MapId      uint32       `protobuf:"varint,3,opt,name=MapId,proto3" json:"MapId"`
	StageId    uint32       `protobuf:"varint,4,opt,name=StageId,proto3" json:"StageId"`
}

type BattleRecordResponse struct {
	PacketHead *IPacket              `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	BattleType EmBattleType          `protobuf:"varint,2,opt,name=BattleType,proto3,enum=common.EmBattleType" json:"BattleType"`
	MapId      uint32                `protobuf:"varint,3,opt,name=MapId,proto3" json:"MapId"`
	StageId    uint32                `protobuf:"varint,4,opt,name=StageId,proto3" json:"StageId"`
	RecordList []*PBPlayerBattleData `protobuf:"bytes,5,rep,name=RecordList,proto3" json:"RecordList"`
}

type BattleReliveNotify struct {
	PacketHead *IPacket        `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Relive     *PBBattleRelive `protobuf:"bytes,2,opt,name=Relive,proto3" json:"Relive"`
}

type BattleReliveRequest struct {
	PacketHead *IPacket     `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	BattleType EmBattleType `protobuf:"varint,2,opt,name=BattleType,proto3,enum=common.EmBattleType" json:"BattleType"`
	HeroId     uint32       `protobuf:"varint,3,opt,name=HeroId,proto3" json:"HeroId"`
	AdvertType uint32       `protobuf:"varint,4,opt,name=AdvertType,proto3" json:"AdvertType"`
}

type BattleReliveResponse struct {
	PacketHead *IPacket        `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	BattleType EmBattleType    `protobuf:"varint,2,opt,name=BattleType,proto3,enum=common.EmBattleType" json:"BattleType"`
	HeroId     uint32          `protobuf:"varint,3,opt,name=HeroId,proto3" json:"HeroId"`
	Relive     *PBBattleRelive `protobuf:"bytes,4,opt,name=Relive,proto3" json:"Relive"`
}

type BattleScheduleSaveRequest struct {
	PacketHead     *IPacket          `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	BattleSchedule *PBBattleSchedule `protobuf:"bytes,2,opt,name=BattleSchedule,proto3" json:"BattleSchedule"`
}

type BattleScheduleSaveResponse struct {
	PacketHead     *IPacket          `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	BattleSchedule *PBBattleSchedule `protobuf:"bytes,2,opt,name=BattleSchedule,proto3" json:"BattleSchedule"`
}

type BookCollectionCoinRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	CrystalID  uint32   `protobuf:"varint,2,opt,name=CrystalID,proto3" json:"CrystalID"`
}

type BookCollectionCoinResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Coin       uint32   `protobuf:"varint,2,opt,name=Coin,proto3" json:"Coin"`
	Level      uint32   `protobuf:"varint,3,opt,name=Level,proto3" json:"Level"`
}

type BookStageRewardRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type BookStageRewardResponse struct {
	PacketHead *IPacket       `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	BookInfo   *PBCrystalBook `protobuf:"bytes,2,opt,name=BookInfo,proto3" json:"BookInfo"`
}

type BoxOpenRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ItemID     uint32   `protobuf:"varint,2,opt,name=ItemID,proto3" json:"ItemID"`
	ItemNum    uint32   `protobuf:"varint,3,opt,name=ItemNum,proto3" json:"ItemNum"`
	AdvestType uint32   `protobuf:"varint,4,opt,name=AdvestType,proto3" json:"AdvestType"`
}

type BoxOpenResponse struct {
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Score      uint32           `protobuf:"varint,2,opt,name=Score,proto3" json:"Score"`
	ItemInfo   []*PBAddItemData `protobuf:"bytes,3,rep,name=ItemInfo,proto3" json:"ItemInfo"`
}

type BoxProgressRewardRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type BoxProgressRewardResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	NeedScore  uint32   `protobuf:"varint,2,opt,name=NeedScore,proto3" json:"NeedScore"`
	Level      uint32   `protobuf:"varint,3,opt,name=Level,proto3" json:"Level"`
	Score      uint32   `protobuf:"varint,4,opt,name=Score,proto3" json:"Score"`
	Recycle    uint32   `protobuf:"varint,5,opt,name=Recycle,proto3" json:"Recycle"`
}

type BroadcastNotify struct {
	PacketHead    *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Channel       uint32   `protobuf:"varint,2,opt,name=Channel,proto3" json:"Channel"`
	BroadcastType uint32   `protobuf:"varint,3,opt,name=BroadcastType,proto3" json:"BroadcastType"`
	Content       string   `protobuf:"bytes,4,opt,name=Content,proto3" json:"Content"`
	Extends       []byte   `protobuf:"bytes,5,opt,name=Extends,proto3" json:"Extends"`
}

type C_Z_LoginCopyMap struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	DataId     int32    `protobuf:"varint,2,opt,name=DataId,proto3" json:"DataId"`
}

type C_Z_Move struct {
	PacketHead *IPacket       `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Move       *C_Z_Move_Move `protobuf:"bytes,2,opt,name=move,proto3" json:"move"`
}

type C_Z_Move_Move struct {
	Mode   int32                 `protobuf:"varint,1,opt,name=Mode,proto3" json:"Mode"`
	Normal *C_Z_Move_Move_Normal `protobuf:"bytes,2,opt,name=normal,proto3" json:"normal"`
}

type C_Z_Move_Move_Normal struct {
	Pos      *Point3F `protobuf:"bytes,1,opt,name=Pos,proto3" json:"Pos"`
	Yaw      float32  `protobuf:"fixed32,2,opt,name=Yaw,proto3" json:"Yaw"`
	Duration float32  `protobuf:"fixed32,3,opt,name=Duration,proto3" json:"Duration"`
}

type C_Z_Skill struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SkillId    int32    `protobuf:"varint,2,opt,name=SkillId,proto3" json:"SkillId"`
	TargetId   int64    `protobuf:"varint,3,opt,name=TargetId,proto3" json:"TargetId"`
}

type ChampionshipInfoRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type ChampionshipInfoResponse struct {
	PacketHead *IPacket                `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	List       []*ChampionshipRankInfo `protobuf:"bytes,2,rep,name=List,proto3" json:"List"`
}

type ChampionshipNotify struct {
	PacketHead *IPacket                `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	CreateTime uint64                  `protobuf:"varint,2,opt,name=CreateTime,proto3" json:"CreateTime"`
	Expire     uint64                  `protobuf:"varint,3,opt,name=Expire,proto3" json:"Expire"`
	List       []*ChampionshipTimeInfo `protobuf:"bytes,4,rep,name=List,proto3" json:"List"`
}

type ChampionshipRankInfo struct {
	RankType uint32    `protobuf:"varint,1,opt,name=RankType,proto3" json:"RankType"`
	First    *RankInfo `protobuf:"bytes,2,opt,name=First,proto3" json:"First"`
}

type ChampionshipTaskRewardRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RankType   uint32   `protobuf:"varint,2,opt,name=RankType,proto3" json:"RankType"`
	TaskID     uint32   `protobuf:"varint,3,opt,name=TaskID,proto3" json:"TaskID"`
}

type ChampionshipTaskRewardResponse struct {
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Task       *PBTaskStageInfo `protobuf:"bytes,2,opt,name=Task,proto3" json:"Task"`
}

type ChampionshipTimeInfo struct {
	RankType uint32 `protobuf:"varint,1,opt,name=RankType,proto3" json:"RankType"`
	Interval uint64 `protobuf:"varint,2,opt,name=Interval,proto3" json:"Interval"`
	Active   uint64 `protobuf:"varint,3,opt,name=Active,proto3" json:"Active"`
	Reward   uint64 `protobuf:"varint,4,opt,name=Reward,proto3" json:"Reward"`
	Show     uint64 `protobuf:"varint,5,opt,name=Show,proto3" json:"Show"`
}

type ChangeAvatarFrameRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	FrameID    uint32   `protobuf:"varint,2,opt,name=FrameID,proto3" json:"FrameID"`
}

type ChangeAvatarFrameResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	FrameID    uint32   `protobuf:"varint,2,opt,name=FrameID,proto3" json:"FrameID"`
}

type ChangeAvatarRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	AvatarID   uint32   `protobuf:"varint,2,opt,name=AvatarID,proto3" json:"AvatarID"`
}

type ChangeAvatarResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	AvatarID   uint32   `protobuf:"varint,2,opt,name=AvatarID,proto3" json:"AvatarID"`
}

type ChangePlayerNameRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PlayerName string   `protobuf:"bytes,2,opt,name=PlayerName,proto3" json:"PlayerName"`
}

type ChangePlayerNameResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PlayerName string   `protobuf:"bytes,2,opt,name=PlayerName,proto3" json:"PlayerName"`
}

type ChargeCardNewNotify struct {
	PacketHead *IPacket      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	CardInfo   *PBChargeCard `protobuf:"bytes,2,opt,name=CardInfo,proto3" json:"CardInfo"`
}

type ChargeCardPrizeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	CardType   uint32   `protobuf:"varint,2,opt,name=CardType,proto3" json:"CardType"`
}

type ChargeCardPrizeResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	CardType   uint32   `protobuf:"varint,2,opt,name=CardType,proto3" json:"CardType"`
	PrizeTime  uint64   `protobuf:"varint,3,opt,name=PrizeTime,proto3" json:"PrizeTime"`
}

type ChargeCardUpdateNotify struct {
	PacketHead *IPacket        `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	CardInfo   []*PBChargeCard `protobuf:"bytes,2,rep,name=CardInfo,proto3" json:"CardInfo"`
	DelList    []uint32        `protobuf:"varint,3,rep,packed,name=DelList,proto3" json:"DelList"`
}

type ChargeGiftBuyNotify struct {
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32    `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	BuyInfo    *PBU32U32 `protobuf:"bytes,3,opt,name=BuyInfo,proto3" json:"BuyInfo"`
}

type ChargeNotify struct {
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Charge     *PBCharge `protobuf:"bytes,2,opt,name=Charge,proto3" json:"Charge"`
}

type ChargeOrderRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProductId  uint32   `protobuf:"varint,2,opt,name=ProductId,proto3" json:"ProductId"`
	IsNeigou   bool     `protobuf:"varint,3,opt,name=IsNeigou,proto3" json:"IsNeigou"`
}

type ChargeOrderResponse struct {
	PacketHead     *IPacket                `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	BingchuanOrder *PBChargeBingchuanOrder `protobuf:"bytes,2,opt,name=BingchuanOrder,proto3" json:"BingchuanOrder"`
}

type ChargeQueryNotify struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProductId  uint32   `protobuf:"varint,2,opt,name=ProductId,proto3" json:"ProductId"`
}

type ChargeQueryRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type ChargeQueryResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProductIds []uint32 `protobuf:"varint,2,rep,packed,name=ProductIds,proto3" json:"ProductIds"`
}

type ClientJsonNotify struct {
	PacketHead *IPacket      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	JsonList   []*PBJsonInfo `protobuf:"bytes,2,rep,name=JsonList,proto3" json:"JsonList"`
}

type ClusterInfo struct {
	Type       SERVICE `protobuf:"varint,1,opt,name=Type,proto3,enum=common.SERVICE" json:"Type"`
	Ip         string  `protobuf:"bytes,2,opt,name=Ip,proto3" json:"Ip"`
	Port       int32   `protobuf:"varint,3,opt,name=Port,proto3" json:"Port"`
	Weight     int32   `protobuf:"varint,4,opt,name=Weight,proto3" json:"Weight"`
	SocketId   uint32  `protobuf:"varint,5,opt,name=SocketId,proto3" json:"SocketId"`
	Version    uint32  `protobuf:"varint,6,opt,name=Version,proto3" json:"Version"`
	ClusterID  uint32  `protobuf:"varint,7,opt,name=ClusterID,proto3" json:"ClusterID"`
	CreateTime uint64  `protobuf:"varint,8,opt,name=CreateTime,proto3" json:"CreateTime"`
}

type CommonNotify struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type CommonPrizeNotify struct {
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	DoingType  EmDoingType      `protobuf:"varint,2,opt,name=DoingType,proto3,enum=common.EmDoingType" json:"DoingType"`
	ItemInfo   []*PBAddItemData `protobuf:"bytes,3,rep,name=ItemInfo,proto3" json:"ItemInfo"`
}

type CrystalGenerateRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	CrystalID  uint32   `protobuf:"varint,2,opt,name=CrystalID,proto3" json:"CrystalID"`
}

type CrystalGenerateResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type CrystalNotify struct {
	PacketHead  *IPacket     `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	CrystalInfo []*PBCrystal `protobuf:"bytes,2,rep,name=CrystalInfo,proto3" json:"CrystalInfo"`
}

type CrystalRedefineRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	CrystalID  uint32   `protobuf:"varint,2,opt,name=CrystalID,proto3" json:"CrystalID"`
}

type CrystalRedefineResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	CrystalID  uint32   `protobuf:"varint,2,opt,name=CrystalID,proto3" json:"CrystalID"`
	CurStar    uint32   `protobuf:"varint,3,opt,name=CurStar,proto3" json:"CurStar"`
}

type CrystalRobotBattleRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RobotID    uint32   `protobuf:"varint,2,opt,name=RobotID,proto3" json:"RobotID"`
}

type CrystalRobotBattleResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RobotID    uint32   `protobuf:"varint,2,opt,name=RobotID,proto3" json:"RobotID"`
}

type CrystalRobotNotify struct {
	PacketHead *IPacket        `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RobotInfo  *PBCrystalRobot `protobuf:"bytes,2,opt,name=RobotInfo,proto3" json:"RobotInfo"`
}

type CrystalRobotUpgradeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RobotID    uint32   `protobuf:"varint,2,opt,name=RobotID,proto3" json:"RobotID"`
}

type CrystalRobotUpgradeResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RobotID    uint32   `protobuf:"varint,2,opt,name=RobotID,proto3" json:"RobotID"`
	CurLevel   uint32   `protobuf:"varint,3,opt,name=CurLevel,proto3" json:"CurLevel"`
}

type CrystalUpgradeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	CrystalID  uint32   `protobuf:"varint,2,opt,name=CrystalID,proto3" json:"CrystalID"`
}

type CrystalUpgradeResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	CrystalID  uint32   `protobuf:"varint,2,opt,name=CrystalID,proto3" json:"CrystalID"`
	CurLevel   uint32   `protobuf:"varint,3,opt,name=CurLevel,proto3" json:"CurLevel"`
}

type DailyTasFinishResponse struct {
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Task       *PBTaskStageInfo `protobuf:"bytes,2,opt,name=Task,proto3" json:"Task"`
	Score      uint32           `protobuf:"varint,3,opt,name=Score,proto3" json:"Score"`
}

type DailyTaskFinishRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	TaskId     uint32   `protobuf:"varint,2,opt,name=TaskId,proto3" json:"TaskId"`
}

type DailyTaskNotify struct {
	PacketHead *IPacket     `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	DailyTask  *PBDailyTask `protobuf:"bytes,2,opt,name=DailyTask,proto3" json:"DailyTask"`
}

type DailyTaskScorePrizeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type DailyTaskScorePrizeResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PrizeScore uint32   `protobuf:"varint,2,opt,name=PrizeScore,proto3" json:"PrizeScore"`
}

type DeleteMailRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	MailId     uint32   `protobuf:"varint,2,opt,name=MailId,proto3" json:"MailId"`
}

type DeleteMailResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	MailId     uint32   `protobuf:"varint,2,opt,name=MailId,proto3" json:"MailId"`
}

type DrawNotify struct {
	PacketHead  *IPacket      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	DrawList    []*PBDrawInfo `protobuf:"bytes,2,rep,name=DrawList,proto3" json:"DrawList"`
	DelDrawList []uint32      `protobuf:"varint,3,rep,packed,name=DelDrawList,proto3" json:"DelDrawList"`
}

type DrawPrizeInfoRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	DrawId     uint32   `protobuf:"varint,2,opt,name=DrawId,proto3" json:"DrawId"`
}

type DrawPrizeInfoResponse struct {
	PacketHead *IPacket           `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	DrawId     uint32             `protobuf:"varint,2,opt,name=DrawId,proto3" json:"DrawId"`
	PrizeList  []*PBDrawPrizeInfo `protobuf:"bytes,3,rep,name=PrizeList,proto3" json:"PrizeList"`
}

type DrawRequest struct {
	PacketHead   *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	DrawId       uint32   `protobuf:"varint,2,opt,name=DrawId,proto3" json:"DrawId"`
	DrawCount    uint32   `protobuf:"varint,3,opt,name=DrawCount,proto3" json:"DrawCount"`
	IsUseReplace bool     `protobuf:"varint,4,opt,name=IsUseReplace,proto3" json:"IsUseReplace"`
	AdvertType   uint32   `protobuf:"varint,5,opt,name=AdvertType,proto3" json:"AdvertType"`
}

type DrawResponse struct {
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	DrawInfo   *PBDrawInfo `protobuf:"bytes,2,opt,name=DrawInfo,proto3" json:"DrawInfo"`
}

type DrawScorePrizeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	DrawId     uint32   `protobuf:"varint,2,opt,name=DrawId,proto3" json:"DrawId"`
	Id         uint32   `protobuf:"varint,3,opt,name=Id,proto3" json:"Id"`
}

type DrawScorePrizeResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	DrawId     uint32   `protobuf:"varint,2,opt,name=DrawId,proto3" json:"DrawId"`
	Id         uint32   `protobuf:"varint,3,opt,name=Id,proto3" json:"Id"`
}

type EjectAdvertNotify struct {
	PacketHead      *IPacket           `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	EjectAdvertInfo *PBEjectAdvertInfo `protobuf:"bytes,2,opt,name=EjectAdvertInfo,proto3" json:"EjectAdvertInfo"`
}

type EjectAdvertRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
}

type EjectAdvertResponse struct {
	PacketHead      *IPacket           `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	EjectAdvertInfo *PBEjectAdvertInfo `protobuf:"bytes,2,opt,name=EjectAdvertInfo,proto3" json:"EjectAdvertInfo"`
}

type EntryCondition struct {
	CfgID      uint32 `protobuf:"varint,1,opt,name=CfgID,proto3" json:"CfgID"`
	Process    uint32 `protobuf:"varint,2,opt,name=Process,proto3" json:"Process"`
	Times      uint32 `protobuf:"varint,3,opt,name=Times,proto3" json:"Times"`
	UpdateTime uint64 `protobuf:"varint,4,opt,name=UpdateTime,proto3" json:"UpdateTime"`
}

type EntryConditionNotify struct {
	PacketHead *IPacket        `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Condition  *EntryCondition `protobuf:"bytes,2,opt,name=Condition,proto3" json:"Condition"`
}

type EntryEffect struct {
	ParamsType uint32             `protobuf:"varint,1,opt,name=ParamsType,proto3" json:"ParamsType"`
	Type       uint32             `protobuf:"varint,2,opt,name=Type,proto3" json:"Type"`
	List       []*EntryEffectData `protobuf:"bytes,3,rep,name=List,proto3" json:"List"`
}

type EntryEffectData struct {
	Object uint32              `protobuf:"varint,1,opt,name=Object,proto3" json:"Object"`
	Values []*EntryEffectValue `protobuf:"bytes,2,rep,name=Values,proto3" json:"Values"`
}

type EntryEffectNotify struct {
	PacketHead *IPacket     `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Effect     *EntryEffect `protobuf:"bytes,2,opt,name=Effect,proto3" json:"Effect"`
}

type EntryEffectValue struct {
	List []uint32 `protobuf:"varint,1,rep,packed,name=List,proto3" json:"List"`
}

type EntryTriggerRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	EntryType  uint32   `protobuf:"varint,2,opt,name=EntryType,proto3" json:"EntryType"`
	Times      uint32   `protobuf:"varint,3,opt,name=Times,proto3" json:"Times"`
	Params     []uint32 `protobuf:"varint,4,rep,packed,name=Params,proto3" json:"Params"`
}

type EntryTriggerResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type EntryUnlockRequest struct {
	PacketHead     *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PassiveSkillID uint32   `protobuf:"varint,2,opt,name=PassiveSkillID,proto3" json:"PassiveSkillID"`
}

type EntryUnlockResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type EquipmentAutoSplitRequest struct {
	PacketHead  *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	QualityList []uint32 `protobuf:"varint,2,rep,packed,name=QualityList,proto3" json:"QualityList"`
}

type EquipmentAutoSplitResponse struct {
	PacketHead  *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	QualityList []uint32 `protobuf:"varint,2,rep,packed,name=QualityList,proto3" json:"QualityList"`
}

type EquipmentBuyPosRequest struct {
	PacketHead     *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	CurPosBuyCount uint32   `protobuf:"varint,2,opt,name=CurPosBuyCount,proto3" json:"CurPosBuyCount"`
}

type EquipmentBuyPosResponse struct {
	PacketHead  *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PosBuyCount uint32   `protobuf:"varint,2,opt,name=PosBuyCount,proto3" json:"PosBuyCount"`
}

type EquipmentLockRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Sn         uint32   `protobuf:"varint,2,opt,name=Sn,proto3" json:"Sn"`
	IsLock     bool     `protobuf:"varint,3,opt,name=IsLock,proto3" json:"IsLock"`
}

type EquipmentLockResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Sn         uint32   `protobuf:"varint,2,opt,name=Sn,proto3" json:"Sn"`
	IsLock     bool     `protobuf:"varint,3,opt,name=IsLock,proto3" json:"IsLock"`
}

type EquipmentNotify struct {
	PacketHead *IPacket       `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Equipment  []*PBEquipment `protobuf:"bytes,2,rep,name=Equipment,proto3" json:"Equipment"`
	IsHook     bool           `protobuf:"varint,3,opt,name=IsHook,proto3" json:"IsHook"`
	DoingType  EmDoingType    `protobuf:"varint,4,opt,name=DoingType,proto3,enum=common.EmDoingType" json:"DoingType"`
}

type EquipmentSplitRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SnList     []uint32 `protobuf:"varint,2,rep,packed,name=SnList,proto3" json:"SnList"`
}

type EquipmentSplitResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SnList     []uint32 `protobuf:"varint,2,rep,packed,name=SnList,proto3" json:"SnList"`
}

type EquipmentSplitScoreNotify struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SplitScore uint32   `protobuf:"varint,2,opt,name=SplitScore,proto3" json:"SplitScore"`
}

type FirstChargeNotify struct {
	PacketHead      *IPacket       `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	FirstChargeInfo *PBFirstCharge `protobuf:"bytes,2,opt,name=FirstChargeInfo,proto3" json:"FirstChargeInfo"`
}

type FirstChargePrizeRequest struct {
	PacketHead    *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	FirstChargeId uint32   `protobuf:"varint,2,opt,name=FirstChargeId,proto3" json:"FirstChargeId"`
}

type FirstChargePrizeResponse struct {
	PacketHead    *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	FirstChargeId uint32   `protobuf:"varint,2,opt,name=FirstChargeId,proto3" json:"FirstChargeId"`
	PrizeDay      uint32   `protobuf:"varint,3,opt,name=PrizeDay,proto3" json:"PrizeDay"`
}

type GeneCardActiveInfo struct {
	IsActive bool   `protobuf:"varint,1,opt,name=IsActive,proto3" json:"IsActive"`
	CardID   uint32 `protobuf:"varint,2,opt,name=CardID,proto3" json:"CardID"`
}

type GeneCardActiveRequest struct {
	PacketHead     *IPacket              `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SchemeID       uint32                `protobuf:"varint,2,opt,name=SchemeID,proto3" json:"SchemeID"`
	Actives        []*GeneCardActiveInfo `protobuf:"bytes,3,rep,name=Actives,proto3" json:"Actives"`
	RobotPositions []uint32              `protobuf:"varint,4,rep,packed,name=RobotPositions,proto3" json:"RobotPositions"`
}

type GeneCardActiveResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type GeneChangeNameRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SchemeID   uint32   `protobuf:"varint,2,opt,name=SchemeID,proto3" json:"SchemeID"`
	Name       string   `protobuf:"bytes,3,opt,name=Name,proto3" json:"Name"`
}

type GeneChangeNameResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type GeneRobotActiveRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SchemeID   uint32   `protobuf:"varint,2,opt,name=SchemeID,proto3" json:"SchemeID"`
	Position   uint32   `protobuf:"varint,3,opt,name=Position,proto3" json:"Position"`
	RobotID    uint32   `protobuf:"varint,4,opt,name=RobotID,proto3" json:"RobotID"`
}

type GeneRobotActiveResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type GeneRobotCardActiveRequest struct {
	PacketHead *IPacket              `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SchemeID   uint32                `protobuf:"varint,2,opt,name=SchemeID,proto3" json:"SchemeID"`
	Position   uint32                `protobuf:"varint,3,opt,name=Position,proto3" json:"Position"`
	Actives    []*GeneCardActiveInfo `protobuf:"bytes,4,rep,name=Actives,proto3" json:"Actives"`
}

type GeneRobotCardActiveResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type GeneSchemeChangeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SchemeID   uint32   `protobuf:"varint,2,opt,name=SchemeID,proto3" json:"SchemeID"`
}

type GeneSchemeChangeResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type GeneSchemeResetRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SchemeID   uint32   `protobuf:"varint,2,opt,name=SchemeID,proto3" json:"SchemeID"`
}

type GeneSchemeResetResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type GetPlayerDataRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	DataType   int32    `protobuf:"varint,2,opt,name=DataType,proto3" json:"DataType"`
}

type GetPlayerDataResponse struct {
	PacketHead *IPacket      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PlayerData *PBPlayerData `protobuf:"bytes,2,opt,name=PlayerData,proto3" json:"PlayerData"`
}

type GiftCodeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Acode      string   `protobuf:"bytes,2,opt,name=Acode,proto3" json:"Acode"`
}

type GiftCodeResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Acode      string   `protobuf:"bytes,2,opt,name=Acode,proto3" json:"Acode"`
}

type GmFuncRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	GmType     uint32   `protobuf:"varint,2,opt,name=GmType,proto3" json:"GmType"`
	Param      []string `protobuf:"bytes,3,rep,name=Param,proto3" json:"Param"`
}

type GmFuncResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type GrowRoadTaskPrizeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	TaskId     uint32   `protobuf:"varint,3,opt,name=TaskId,proto3" json:"TaskId"`
}

type GrowRoadTaskPrizeResponse struct {
	PacketHead    *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id            uint32           `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	TaskStageInfo *PBTaskStageInfo `protobuf:"bytes,4,opt,name=TaskStageInfo,proto3" json:"TaskStageInfo"`
}

type HeardPacket struct {
}

type HeartbeatRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Time       uint64   `protobuf:"varint,2,opt,name=Time,proto3" json:"Time"`
}

type HeartbeatResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SendTime   uint64   `protobuf:"varint,2,opt,name=SendTime,proto3" json:"SendTime"`
	RecvTime   uint64   `protobuf:"varint,3,opt,name=RecvTime,proto3" json:"RecvTime"`
	CurTime    uint64   `protobuf:"varint,4,opt,name=CurTime,proto3" json:"CurTime"`
}

type HeroAutoUpStarRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type HeroAutoUpStarResponse struct {
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	HeroList   []*PBU32U32 `protobuf:"bytes,2,rep,name=HeroList,proto3" json:"HeroList"`
	DelSnList  []uint32    `protobuf:"varint,3,rep,packed,name=DelSnList,proto3" json:"DelSnList"`
}

type HeroAwakenLevelRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Sn         uint32   `protobuf:"varint,2,opt,name=Sn,proto3" json:"Sn"`
	CurLevel   uint32   `protobuf:"varint,3,opt,name=CurLevel,proto3" json:"CurLevel"`
}

type HeroAwakenLevelResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Sn         uint32   `protobuf:"varint,2,opt,name=Sn,proto3" json:"Sn"`
	CurLevel   uint32   `protobuf:"varint,3,opt,name=CurLevel,proto3" json:"CurLevel"`
}

type HeroBattleStarChangeNotify struct {
	PacketHead *IPacket              `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Heros      []*HeroBattleStarInfo `protobuf:"bytes,2,rep,name=Heros,proto3" json:"Heros"`
}

type HeroBattleStarChangeRequest struct {
	PacketHead *IPacket              `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Heros      []*HeroBattleStarInfo `protobuf:"bytes,2,rep,name=Heros,proto3" json:"Heros"`
}

type HeroBattleStarChangeResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type HeroBattleStarInfo struct {
	HeroID uint32 `protobuf:"varint,1,opt,name=HeroID,proto3" json:"HeroID"`
	Total  uint32 `protobuf:"varint,2,opt,name=Total,proto3" json:"Total"`
}

type HeroBookActiveRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	HeroId     uint32   `protobuf:"varint,2,opt,name=HeroId,proto3" json:"HeroId"`
}

type HeroBookActiveResponse struct {
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	HeroBook   *PBHeroBook `protobuf:"bytes,2,opt,name=HeroBook,proto3" json:"HeroBook"`
}

type HeroBookNotify struct {
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	HeroBook   *PBHeroBook `protobuf:"bytes,2,opt,name=HeroBook,proto3" json:"HeroBook"`
}

type HeroBookUpStarRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	HeroId     uint32   `protobuf:"varint,2,opt,name=HeroId,proto3" json:"HeroId"`
}

type HeroBookUpStarResponse struct {
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	HeroBook   *PBHeroBook `protobuf:"bytes,2,opt,name=HeroBook,proto3" json:"HeroBook"`
}

type HeroFightPowerNotify struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	FightPower uint32   `protobuf:"varint,2,opt,name=FightPower,proto3" json:"FightPower"`
}

type HeroGameHeroListNotify struct {
	PacketHead *IPacket        `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Team       *PBHeroTeamList `protobuf:"bytes,2,opt,name=Team,proto3" json:"Team"`
}

type HeroGameHeroListRequest struct {
	PacketHead *IPacket        `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Team       *PBHeroTeamList `protobuf:"bytes,2,opt,name=Team,proto3" json:"Team"`
}

type HeroGameHeroListResponse struct {
	PacketHead *IPacket        `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Team       *PBHeroTeamList `protobuf:"bytes,2,opt,name=Team,proto3" json:"Team"`
}

type HeroNewStarNotify struct {
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Info       []*PBHero `protobuf:"bytes,2,rep,name=Info,proto3" json:"Info"`
}

type HeroNotify struct {
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Info       []*PBHero `protobuf:"bytes,2,rep,name=Info,proto3" json:"Info"`
}

type HeroRebirthRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Sn         uint32   `protobuf:"varint,2,opt,name=Sn,proto3" json:"Sn"`
}

type HeroRebirthResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Sn         uint32   `protobuf:"varint,2,opt,name=Sn,proto3" json:"Sn"`
	Star       uint32   `protobuf:"varint,3,opt,name=Star,proto3" json:"Star"`
}

type HeroUpStarRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Sn         uint32   `protobuf:"varint,2,opt,name=Sn,proto3" json:"Sn"`
	UseSnList  []uint32 `protobuf:"varint,3,rep,packed,name=UseSnList,proto3" json:"UseSnList"`
}

type HeroUpStarResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Sn         uint32   `protobuf:"varint,2,opt,name=Sn,proto3" json:"Sn"`
	Star       uint32   `protobuf:"varint,3,opt,name=Star,proto3" json:"Star"`
	UseSnList  []uint32 `protobuf:"varint,4,rep,packed,name=UseSnList,proto3" json:"UseSnList"`
}

type HookBattleAutoMapRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type HookBattleAutoMapResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	AutoMap    bool     `protobuf:"varint,2,opt,name=AutoMap,proto3" json:"AutoMap"`
}

type HookBattleLootRequest struct {
	PacketHead  *IPacket                 `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	MonsterInfo []*BattleKillMonsterInfo `protobuf:"bytes,2,rep,name=MonsterInfo,proto3" json:"MonsterInfo"`
}

type HookBattleLootResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type HookEquipmentAwardRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Sn         uint32   `protobuf:"varint,2,opt,name=Sn,proto3" json:"Sn"`
}

type HookEquipmentAwardResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Sn         uint32   `protobuf:"varint,2,opt,name=Sn,proto3" json:"Sn"`
}

type HookTechLevelNotify struct {
	PacketHead   *IPacket      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	HookTechList []*PBHookTech `protobuf:"bytes,2,rep,name=HookTechList,proto3" json:"HookTechList"`
}

type HookTechLevelRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
}

type HookTechLevelResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	LevelTime  uint64   `protobuf:"varint,3,opt,name=LevelTime,proto3" json:"LevelTime"`
}

type HookTechSpeedRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	AdvertType uint32   `protobuf:"varint,3,opt,name=AdvertType,proto3" json:"AdvertType"`
}

type HookTechSpeedResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	Level      uint32   `protobuf:"varint,3,opt,name=Level,proto3" json:"Level"`
	LevelTime  uint64   `protobuf:"varint,4,opt,name=LevelTime,proto3" json:"LevelTime"`
}

type IPacket struct {
	Stx            uint32  `protobuf:"varint,1,opt,name=stx,proto3" json:"stx"`
	Ckx            uint32  `protobuf:"varint,2,opt,name=ckx,proto3" json:"ckx"`
	DestServerType SERVICE `protobuf:"varint,3,opt,name=destServerType,proto3,enum=common.SERVICE" json:"destServerType"`
	Id             uint64  `protobuf:"varint,4,opt,name=id,proto3" json:"id"`
	Seqid          uint32  `protobuf:"varint,5,opt,name=seqid,proto3" json:"seqid"`
	Code           uint32  `protobuf:"varint,6,opt,name=code,proto3" json:"code"`
}

type ItemBuyRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	Count      uint32   `protobuf:"varint,3,opt,name=Count,proto3" json:"Count"`
}

type ItemBuyResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	Count      uint32   `protobuf:"varint,3,opt,name=Count,proto3" json:"Count"`
}

type ItemSelectRequest struct {
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32      `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	SelectList []*PBU32U32 `protobuf:"bytes,3,rep,name=SelectList,proto3" json:"SelectList"`
}

type ItemSelectResponse struct {
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32      `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	SelectList []*PBU32U32 `protobuf:"bytes,3,rep,name=SelectList,proto3" json:"SelectList"`
}

type ItemUpdateNotify struct {
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ItemList   []*PBItem `protobuf:"bytes,2,rep,name=ItemList,proto3" json:"ItemList"`
}

type ItemUseRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	Count      uint32   `protobuf:"varint,3,opt,name=Count,proto3" json:"Count"`
}

type ItemUseResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	Count      uint32   `protobuf:"varint,3,opt,name=Count,proto3" json:"Count"`
}

type ItemUseShowInfo struct {
	Id   uint32     `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`
	Item *PBAddItem `protobuf:"bytes,2,opt,name=Item,proto3" json:"Item"`
}

type ItemUseShowRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
}

type ItemUseShowResponse struct {
	PacketHead *IPacket           `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32             `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	ItemList   []*ItemUseShowInfo `protobuf:"bytes,3,rep,name=ItemList,proto3" json:"ItemList"`
}

type LoginRequest struct {
	PacketHead  *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	AccountName string   `protobuf:"bytes,2,opt,name=AccountName,proto3" json:"AccountName"`
	TokenKey    string   `protobuf:"bytes,3,opt,name=TokenKey,proto3" json:"TokenKey"`
}

type LoginResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Time       uint64   `protobuf:"varint,2,opt,name=Time,proto3" json:"Time"`
}

type MailBox struct {
	Id        uint64 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`
	LeaseId   uint64 `protobuf:"varint,2,opt,name=LeaseId,proto3" json:"LeaseId"`
	MailType  MAIL   `protobuf:"varint,3,opt,name=MailType,proto3,enum=rpc3.MAIL" json:"MailType"`
	ClusterId uint32 `protobuf:"varint,5,opt,name=ClusterId,proto3" json:"ClusterId"`
}

type MainTaskFinishRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type MainTaskFinishResponse struct {
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	MainTask   *PBTaskStageInfo `protobuf:"bytes,2,opt,name=MainTask,proto3" json:"MainTask"`
}

type MainTaskNotify struct {
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	MainTask   *PBTaskStageInfo `protobuf:"bytes,2,opt,name=MainTask,proto3" json:"MainTask"`
}

type NewMailNotify struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Mail       *PBMail  `protobuf:"bytes,2,opt,name=Mail,proto3" json:"Mail"`
}

type NormalBattlePrizeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	MapId      uint32   `protobuf:"varint,2,opt,name=MapId,proto3" json:"MapId"`
	StageId    uint32   `protobuf:"varint,3,opt,name=StageId,proto3" json:"StageId"`
	PrizeStage uint32   `protobuf:"varint,4,opt,name=PrizeStage,proto3" json:"PrizeStage"`
}

type NormalBattlePrizeResponse struct {
	PacketHead   *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PrizeMapId   uint32   `protobuf:"varint,2,opt,name=PrizeMapId,proto3" json:"PrizeMapId"`
	PrizeStageId uint32   `protobuf:"varint,3,opt,name=PrizeStageId,proto3" json:"PrizeStageId"`
	PrizeStage   []uint32 `protobuf:"varint,4,rep,packed,name=PrizeStage,proto3" json:"PrizeStage"`
}

type NoticeNotify struct {
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	IsNew      bool      `protobuf:"varint,2,opt,name=IsNew,proto3" json:"IsNew"`
	Notice     *PBNotice `protobuf:"bytes,3,opt,name=Notice,proto3" json:"Notice"`
}

type NoticeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type NoticeResponse struct {
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	NoticeList []*PBNotice `protobuf:"bytes,2,rep,name=NoticeList,proto3" json:"NoticeList"`
}

type OfflineIncomeRewardRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	AdvertType uint32   `protobuf:"varint,2,opt,name=AdvertType,proto3" json:"AdvertType"`
}

type OfflineIncomeRewardResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type OnekeyAwardMailRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type OnekeyAwardMailResponse struct {
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Mails      []*PBMail `protobuf:"bytes,2,rep,name=Mails,proto3" json:"Mails"`
}

type OnekeyDeleteMailRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type OnekeyDeleteMailResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	MailIds    []uint32 `protobuf:"varint,2,rep,packed,name=MailIds,proto3" json:"MailIds"`
}

type OpenBossRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type OpenBossResponse struct {
	PacketHead    *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	WorldBossRank *PBU32U64 `protobuf:"bytes,2,opt,name=WorldBossRank,proto3" json:"WorldBossRank"`
}

type OpenServerGiftBuyNotify struct {
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32    `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	GiftId     uint32    `protobuf:"varint,3,opt,name=GiftId,proto3" json:"GiftId"`
	BuyInfo    *PBU32U32 `protobuf:"bytes,4,opt,name=BuyInfo,proto3" json:"BuyInfo"`
}

type OpenServerGiftNewNotify struct {
	PacketHead *IPacket              `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SId        uint32                `protobuf:"varint,2,opt,name=SId,proto3" json:"SId"`
	GiftInfo   *PBOpenServerGiftInfo `protobuf:"bytes,3,opt,name=GiftInfo,proto3" json:"GiftInfo"`
}

type PBAchieveInfo struct {
	AchieveType uint32   `protobuf:"varint,1,opt,name=AchieveType,proto3" json:"AchieveType"`
	Params      []uint32 `protobuf:"varint,2,rep,packed,name=Params,proto3" json:"Params"`
	Value       uint32   `protobuf:"varint,3,opt,name=Value,proto3" json:"Value"`
}

type PBActivityAdventure struct {
	Id           uint32   `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`
	BeginTime    uint64   `protobuf:"varint,2,opt,name=BeginTime,proto3" json:"BeginTime"`
	EndTime      uint64   `protobuf:"varint,3,opt,name=EndTime,proto3" json:"EndTime"`
	RegisterTime uint64   `protobuf:"varint,4,opt,name=RegisterTime,proto3" json:"RegisterTime"`
	PrizeIdList  []uint32 `protobuf:"varint,5,rep,packed,name=PrizeIdList,proto3" json:"PrizeIdList"`
}

type PBActivityChargeGift struct {
	Id        uint32      `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`
	BeginTime uint64      `protobuf:"varint,2,opt,name=BeginTime,proto3" json:"BeginTime"`
	EndTime   uint64      `protobuf:"varint,3,opt,name=EndTime,proto3" json:"EndTime"`
	BuyList   []*PBU32U32 `protobuf:"bytes,4,rep,name=BuyList,proto3" json:"BuyList"`
}

type PBActivityGrowRoadInfo struct {
	Id        uint32             `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`
	BeginTime uint64             `protobuf:"varint,2,opt,name=BeginTime,proto3" json:"BeginTime"`
	EndTime   uint64             `protobuf:"varint,3,opt,name=EndTime,proto3" json:"EndTime"`
	TaskList  []*PBTaskStageInfo `protobuf:"bytes,4,rep,name=TaskList,proto3" json:"TaskList"`
}

type PBActivityOpenServerGift struct {
	Id                 uint32                  `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`
	BeginTime          uint64                  `protobuf:"varint,2,opt,name=BeginTime,proto3" json:"BeginTime"`
	EndTime            uint64                  `protobuf:"varint,3,opt,name=EndTime,proto3" json:"EndTime"`
	GiftList           []*PBOpenServerGiftInfo `protobuf:"bytes,4,rep,name=GiftList,proto3" json:"GiftList"`
	NextDailyPrizeTime uint64                  `protobuf:"varint,5,opt,name=NextDailyPrizeTime,proto3" json:"NextDailyPrizeTime"`
}

type PBAddItem struct {
	Kind   uint32   `protobuf:"varint,1,opt,name=Kind,proto3" json:"Kind"`
	Id     uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	Count  int64    `protobuf:"varint,3,opt,name=Count,proto3" json:"Count"`
	Params []uint32 `protobuf:"varint,4,rep,packed,name=Params,proto3" json:"Params"`
}

type PBAddItemData struct {
	Kind      uint32       `protobuf:"varint,1,opt,name=Kind,proto3" json:"Kind"`
	Id        uint32       `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	Count     int64        `protobuf:"varint,3,opt,name=Count,proto3" json:"Count"`
	DoingType EmDoingType  `protobuf:"varint,4,opt,name=DoingType,proto3,enum=common.EmDoingType" json:"DoingType"`
	Params    []uint32     `protobuf:"varint,5,rep,packed,name=Params,proto3" json:"Params"`
	Equipment *PBEquipment `protobuf:"bytes,6,opt,name=Equipment,proto3" json:"Equipment"`
}

type PBAdvertInfo struct {
	Type       uint32 `protobuf:"varint,1,opt,name=Type,proto3" json:"Type"`
	DailyCount uint32 `protobuf:"varint,2,opt,name=DailyCount,proto3" json:"DailyCount"`
}

type PBAllChatMsgInfo struct {
	Msg []*PBChatMsgInfo `protobuf:"bytes,1,rep,name=Msg,proto3" json:"Msg"`
}

type PBAvatar struct {
	AvatarID uint32 `protobuf:"varint,1,opt,name=AvatarID,proto3" json:"AvatarID"`
	Type     uint32 `protobuf:"varint,2,opt,name=Type,proto3" json:"Type"`
}

type PBAvatarFrame struct {
	FrameID uint32 `protobuf:"varint,1,opt,name=FrameID,proto3" json:"FrameID"`
	Type    uint32 `protobuf:"varint,2,opt,name=Type,proto3" json:"Type"`
}

type PBBPInfo struct {
	BPType    uint32           `protobuf:"varint,1,opt,name=BPType,proto3" json:"BPType"`
	Value     uint32           `protobuf:"varint,2,opt,name=Value,proto3" json:"Value"`
	StageList []*PBBPStageInfo `protobuf:"bytes,3,rep,name=StageList,proto3" json:"StageList"`
	MaxStage  uint32           `protobuf:"varint,4,opt,name=MaxStage,proto3" json:"MaxStage"`
}

type PBBPStageInfo struct {
	StageId       uint32 `protobuf:"varint,1,opt,name=StageId,proto3" json:"StageId"`
	BeginTime     uint64 `protobuf:"varint,2,opt,name=BeginTime,proto3" json:"BeginTime"`
	EndTime       uint64 `protobuf:"varint,3,opt,name=EndTime,proto3" json:"EndTime"`
	NoramlPrizeId uint32 `protobuf:"varint,4,opt,name=NoramlPrizeId,proto3" json:"NoramlPrizeId"`
	ExtralPrizeId uint32 `protobuf:"varint,5,opt,name=ExtralPrizeId,proto3" json:"ExtralPrizeId"`
	ChargeTime    uint64 `protobuf:"varint,6,opt,name=ChargeTime,proto3" json:"ChargeTime"`
}

type PBBattleClientData struct {
	HeroBattleLevel []*PBU32U32 `protobuf:"bytes,1,rep,name=HeroBattleLevel,proto3" json:"HeroBattleLevel"`
	SelectCard      []uint32    `protobuf:"varint,2,rep,packed,name=SelectCard,proto3" json:"SelectCard"`
	LeaderId        uint32      `protobuf:"varint,3,opt,name=LeaderId,proto3" json:"LeaderId"`
	CryscalRobotId  []uint32    `protobuf:"varint,4,rep,packed,name=CryscalRobotId,proto3" json:"CryscalRobotId"`
	DropBoxCount    uint32      `protobuf:"varint,5,opt,name=DropBoxCount,proto3" json:"DropBoxCount"`
}

type PBBattleHookInfo struct {
	MapInfo        *PBBattleMapInfo `protobuf:"bytes,1,opt,name=MapInfo,proto3" json:"MapInfo"`
	CurMapId       uint32           `protobuf:"varint,2,opt,name=CurMapId,proto3" json:"CurMapId"`
	CurStageId     uint32           `protobuf:"varint,3,opt,name=CurStageId,proto3" json:"CurStageId"`
	AutoMap        bool             `protobuf:"varint,4,opt,name=AutoMap,proto3" json:"AutoMap"`
	BeginLootTime  uint64           `protobuf:"varint,5,opt,name=BeginLootTime,proto3" json:"BeginLootTime"`
	TotalLootCount uint32           `protobuf:"varint,6,opt,name=TotalLootCount,proto3" json:"TotalLootCount"`
	PassRewards    []uint64         `protobuf:"varint,7,rep,packed,name=PassRewards,proto3" json:"PassRewards"`
}

type PBBattleMapInfo struct {
	MapId        uint32 `protobuf:"varint,1,opt,name=MapId,proto3" json:"MapId"`
	StageId      uint32 `protobuf:"varint,2,opt,name=StageId,proto3" json:"StageId"`
	Time         uint64 `protobuf:"varint,3,opt,name=Time,proto3" json:"Time"`
	StageRate    uint32 `protobuf:"varint,4,opt,name=StageRate,proto3" json:"StageRate"`
	UseTime      uint32 `protobuf:"varint,5,opt,name=UseTime,proto3" json:"UseTime"`
	FightCount   uint32 `protobuf:"varint,6,opt,name=FightCount,proto3" json:"FightCount"`
	IsSuceess    uint32 `protobuf:"varint,7,opt,name=IsSuceess,proto3" json:"IsSuceess"`
	RebirthCount uint32 `protobuf:"varint,8,opt,name=RebirthCount,proto3" json:"RebirthCount"`
}

type PBBattleNormalInfo struct {
	MapInfo      *PBBattleMapInfo `protobuf:"bytes,1,opt,name=MapInfo,proto3" json:"MapInfo"`
	PrizeMapId   uint32           `protobuf:"varint,2,opt,name=PrizeMapId,proto3" json:"PrizeMapId"`
	PrizeStageId uint32           `protobuf:"varint,3,opt,name=PrizeStageId,proto3" json:"PrizeStageId"`
	PrizeStage   []uint32         `protobuf:"varint,4,rep,packed,name=PrizeStage,proto3" json:"PrizeStage"`
}

type PBBattleRelive struct {
	AdvestReliveCount uint32 `protobuf:"varint,1,opt,name=AdvestReliveCount,proto3" json:"AdvestReliveCount"`
	ShareReliveCount  uint32 `protobuf:"varint,2,opt,name=ShareReliveCount,proto3" json:"ShareReliveCount"`
}

type PBBattleSchedule struct {
	BattleType   uint32                   `protobuf:"varint,1,opt,name=BattleType,proto3" json:"BattleType"`
	StageRate    uint32                   `protobuf:"varint,2,opt,name=StageRate,proto3" json:"StageRate"`
	UseTime      uint32                   `protobuf:"varint,3,opt,name=UseTime,proto3" json:"UseTime"`
	RebirthCount uint32                   `protobuf:"varint,4,opt,name=RebirthCount,proto3" json:"RebirthCount"`
	ClientData   *PBBattleClientData      `protobuf:"bytes,5,opt,name=ClientData,proto3" json:"ClientData"`
	MonsterInfo  []*BattleKillMonsterInfo `protobuf:"bytes,6,rep,name=MonsterInfo,proto3" json:"MonsterInfo"`
}

type PBBlackShop struct {
	NextRefreshTime uint64             `protobuf:"varint,1,opt,name=NextRefreshTime,proto3" json:"NextRefreshTime"`
	Items           []*PBShopGoodInfo  `protobuf:"bytes,2,rep,name=Items,proto3" json:"Items"`
	RefreshInfo     *PBShopRefreshInfo `protobuf:"bytes,3,opt,name=RefreshInfo,proto3" json:"RefreshInfo"`
}

type PBBoxInfo struct {
	ItemID    uint32 `protobuf:"varint,1,opt,name=ItemID,proto3" json:"ItemID"`
	OpenTimes uint32 `protobuf:"varint,2,opt,name=OpenTimes,proto3" json:"OpenTimes"`
}

type PBCharge struct {
	OrderId     uint32 `protobuf:"varint,1,opt,name=OrderId,proto3" json:"OrderId"`
	TotalCharge uint32 `protobuf:"varint,2,opt,name=TotalCharge,proto3" json:"TotalCharge"`
	DailyCharge uint32 `protobuf:"varint,3,opt,name=DailyCharge,proto3" json:"DailyCharge"`
	WeekCharge  uint32 `protobuf:"varint,4,opt,name=WeekCharge,proto3" json:"WeekCharge"`
	MonthCharge uint32 `protobuf:"varint,5,opt,name=MonthCharge,proto3" json:"MonthCharge"`
}

type PBChargeBingchuanOrder struct {
	OrderItem        string `protobuf:"bytes,1,opt,name=OrderItem,proto3" json:"OrderItem"`
	OrderNo          string `protobuf:"bytes,2,opt,name=OrderNo,proto3" json:"OrderNo"`
	PayNum           string `protobuf:"bytes,3,opt,name=payNum,proto3" json:"payNum"`
	UserId           string `protobuf:"bytes,4,opt,name=userId,proto3" json:"userId"`
	OrderSign        string `protobuf:"bytes,5,opt,name=orderSign,proto3" json:"orderSign"`
	ActorId          string `protobuf:"bytes,6,opt,name=actorId,proto3" json:"actorId"`
	CurrencyType     string `protobuf:"bytes,7,opt,name=currencyType,proto3" json:"currencyType"`
	DeveloperPayload string `protobuf:"bytes,8,opt,name=developerPayload,proto3" json:"developerPayload"`
}

type PBChargeCard struct {
	CardType  uint32 `protobuf:"varint,1,opt,name=CardType,proto3" json:"CardType"`
	BeginTime uint64 `protobuf:"varint,2,opt,name=BeginTime,proto3" json:"BeginTime"`
	EndTime   uint64 `protobuf:"varint,3,opt,name=EndTime,proto3" json:"EndTime"`
	PrizeTime uint64 `protobuf:"varint,4,opt,name=PrizeTime,proto3" json:"PrizeTime"`
}

type PBChatMsgInfo struct {
	Index   uint64               `protobuf:"varint,1,opt,name=Index,proto3" json:"Index"`
	Display *PBPlayerBaseDisplay `protobuf:"bytes,2,opt,name=Display,proto3" json:"Display"`
	Msg     string               `protobuf:"bytes,3,opt,name=Msg,proto3" json:"Msg"`
	Time    uint64               `protobuf:"varint,4,opt,name=Time,proto3" json:"Time"`
}

type PBClientData struct {
	Type string `protobuf:"bytes,1,opt,name=Type,proto3" json:"Type"`
	Data string `protobuf:"bytes,2,opt,name=Data,proto3" json:"Data"`
}

type PBCrystal struct {
	CrystalID       uint32   `protobuf:"varint,1,opt,name=CrystalID,proto3" json:"CrystalID"`
	Element         uint32   `protobuf:"varint,2,opt,name=Element,proto3" json:"Element"`
	Quality         uint32   `protobuf:"varint,3,opt,name=Quality,proto3" json:"Quality"`
	Star            uint32   `protobuf:"varint,4,opt,name=Star,proto3" json:"Star"`
	RewardCoinTimes uint32   `protobuf:"varint,5,opt,name=RewardCoinTimes,proto3" json:"RewardCoinTimes"`
	PassiveSkillIds []uint32 `protobuf:"varint,6,rep,packed,name=PassiveSkillIds,proto3" json:"PassiveSkillIds"`
	Level           uint32   `protobuf:"varint,7,opt,name=Level,proto3" json:"Level"`
}

type PBCrystalBook struct {
	Coin          uint32 `protobuf:"varint,1,opt,name=Coin,proto3" json:"Coin"`
	Stage         uint32 `protobuf:"varint,2,opt,name=Stage,proto3" json:"Stage"`
	FinishedStage uint32 `protobuf:"varint,3,opt,name=FinishedStage,proto3" json:"FinishedStage"`
}

type PBCrystalProp struct {
	Key   uint32 `protobuf:"varint,1,opt,name=Key,proto3" json:"Key"`
	Value uint32 `protobuf:"varint,2,opt,name=Value,proto3" json:"Value"`
}

type PBCrystalRobot struct {
	RobotID          uint32   `protobuf:"varint,1,opt,name=RobotID,proto3" json:"RobotID"`
	Stage            uint32   `protobuf:"varint,2,opt,name=Stage,proto3" json:"Stage"`
	FinishStage      uint32   `protobuf:"varint,3,opt,name=FinishStage,proto3" json:"FinishStage"`
	RoleSkillID      uint32   `protobuf:"varint,4,opt,name=RoleSkillID,proto3" json:"RoleSkillID"`
	RoleSkillPercent uint32   `protobuf:"varint,5,opt,name=RoleSkillPercent,proto3" json:"RoleSkillPercent"`
	UnlockLinkages   []uint32 `protobuf:"varint,6,rep,packed,name=UnlockLinkages,proto3" json:"UnlockLinkages"`
	Crystals         []uint32 `protobuf:"varint,7,rep,packed,name=Crystals,proto3" json:"Crystals"`
}

type PBDailyTask struct {
	TaskList   []*PBTaskStageInfo `protobuf:"bytes,1,rep,name=TaskList,proto3" json:"TaskList"`
	Score      uint32             `protobuf:"varint,2,opt,name=Score,proto3" json:"Score"`
	PrizeScore uint32             `protobuf:"varint,3,opt,name=PrizeScore,proto3" json:"PrizeScore"`
}

type PBDrawInfo struct {
	DrawId         uint32   `protobuf:"varint,1,opt,name=DrawId,proto3" json:"DrawId"`
	DrawCount      uint32   `protobuf:"varint,2,opt,name=DrawCount,proto3" json:"DrawCount"`
	GuarCount      uint32   `protobuf:"varint,3,opt,name=GuarCount,proto3" json:"GuarCount"`
	Guar2Count     uint32   `protobuf:"varint,4,opt,name=Guar2Count,proto3" json:"Guar2Count"`
	AdvertNextTime uint64   `protobuf:"varint,5,opt,name=AdvertNextTime,proto3" json:"AdvertNextTime"`
	BeginTime      uint64   `protobuf:"varint,6,opt,name=BeginTime,proto3" json:"BeginTime"`
	EndTime        uint64   `protobuf:"varint,7,opt,name=EndTime,proto3" json:"EndTime"`
	Guar3Count     uint32   `protobuf:"varint,8,opt,name=Guar3Count,proto3" json:"Guar3Count"`
	ScorePrize     []uint32 `protobuf:"varint,9,rep,packed,name=ScorePrize,proto3" json:"ScorePrize"`
}

type PBDrawPrizeInfo struct {
	Name     string       `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name"`
	Rate     uint32       `protobuf:"varint,2,opt,name=Rate,proto3" json:"Rate"`
	ItemList []*PBAddItem `protobuf:"bytes,3,rep,name=ItemList,proto3" json:"ItemList"`
	Rate2    uint32       `protobuf:"varint,4,opt,name=Rate2,proto3" json:"Rate2"`
}

type PBEjectAdvertInfo struct {
	Id              uint32 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`
	NextRefreshTime uint64 `protobuf:"varint,2,opt,name=NextRefreshTime,proto3" json:"NextRefreshTime"`
}

type PBEquipment struct {
	Sn              uint32             `protobuf:"varint,1,opt,name=Sn,proto3" json:"Sn"`
	Id              uint32             `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	Quality         uint32             `protobuf:"varint,3,opt,name=Quality,proto3" json:"Quality"`
	Star            uint32             `protobuf:"varint,4,opt,name=Star,proto3" json:"Star"`
	MainProp        *PBEquipmentProp   `protobuf:"bytes,5,opt,name=MainProp,proto3" json:"MainProp"`
	MinorPropList   []*PBEquipmentProp `protobuf:"bytes,6,rep,name=MinorPropList,proto3" json:"MinorPropList"`
	VicePropList    []*PBEquipmentProp `protobuf:"bytes,7,rep,name=VicePropList,proto3" json:"VicePropList"`
	LinkPropList    []*PBEquipmentProp `protobuf:"bytes,8,rep,name=LinkPropList,proto3" json:"LinkPropList"`
	EquipProfession uint32             `protobuf:"varint,9,opt,name=EquipProfession,proto3" json:"EquipProfession"`
	IsLock          bool               `protobuf:"varint,10,opt,name=IsLock,proto3" json:"IsLock"`
}

type PBEquipmentProp struct {
	PropId uint32 `protobuf:"varint,1,opt,name=PropId,proto3" json:"PropId"`
	Value  uint32 `protobuf:"varint,2,opt,name=Value,proto3" json:"Value"`
	Score  uint32 `protobuf:"varint,3,opt,name=Score,proto3" json:"Score"`
}

type PBFirstCharge struct {
	FirstChargeId uint32 `protobuf:"varint,1,opt,name=FirstChargeId,proto3" json:"FirstChargeId"`
	ActiveTime    uint64 `protobuf:"varint,2,opt,name=ActiveTime,proto3" json:"ActiveTime"`
	PrizeDay      uint32 `protobuf:"varint,3,opt,name=PrizeDay,proto3" json:"PrizeDay"`
	OpenTime      uint64 `protobuf:"varint,4,opt,name=OpenTime,proto3" json:"OpenTime"`
}

type PBGeneRobot struct {
	RobotID  uint32       `protobuf:"varint,1,opt,name=RobotID,proto3" json:"RobotID"`
	Position uint32       `protobuf:"varint,2,opt,name=Position,proto3" json:"Position"`
	Tags     []*PBGeneTag `protobuf:"bytes,3,rep,name=Tags,proto3" json:"Tags"`
}

type PBGeneScheme struct {
	SchemeID uint32         `protobuf:"varint,1,opt,name=SchemeID,proto3" json:"SchemeID"`
	Name     string         `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name"`
	Tags     []*PBGeneTag   `protobuf:"bytes,3,rep,name=Tags,proto3" json:"Tags"`
	Robots   []*PBGeneRobot `protobuf:"bytes,4,rep,name=Robots,proto3" json:"Robots"`
}

type PBGeneTag struct {
	TagID uint32   `protobuf:"varint,1,opt,name=TagID,proto3" json:"TagID"`
	Cards []uint32 `protobuf:"varint,2,rep,packed,name=Cards,proto3" json:"Cards"`
}

type PBHero struct {
	Sn          uint32 `protobuf:"varint,1,opt,name=Sn,proto3" json:"Sn"`
	Id          uint32 `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	Star        uint32 `protobuf:"varint,3,opt,name=Star,proto3" json:"Star"`
	AwakenLevel uint32 `protobuf:"varint,4,opt,name=AwakenLevel,proto3" json:"AwakenLevel"`
	BattleStar  uint32 `protobuf:"varint,5,opt,name=BattleStar,proto3" json:"BattleStar"`
}

type PBHeroBook struct {
	Id      uint32 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`
	Star    uint32 `protobuf:"varint,2,opt,name=Star,proto3" json:"Star"`
	MaxStar uint32 `protobuf:"varint,3,opt,name=MaxStar,proto3" json:"MaxStar"`
}

type PBHeroTeamList struct {
	TeamType uint32   `protobuf:"varint,1,opt,name=TeamType,proto3" json:"TeamType"`
	HeroSn   []uint32 `protobuf:"varint,2,rep,packed,name=HeroSn,proto3" json:"HeroSn"`
}

type PBHookTech struct {
	Id        uint32 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`
	Level     uint32 `protobuf:"varint,2,opt,name=Level,proto3" json:"Level"`
	LevelTime uint64 `protobuf:"varint,3,opt,name=LevelTime,proto3" json:"LevelTime"`
}

type PBItem struct {
	Id    uint32 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`
	Count int64  `protobuf:"varint,2,opt,name=Count,proto3" json:"Count"`
}

type PBJsonInfo struct {
	JsonName string `protobuf:"bytes,1,opt,name=JsonName,proto3" json:"JsonName"`
	JsonData string `protobuf:"bytes,2,opt,name=JsonData,proto3" json:"JsonData"`
}

type PBMail struct {
	Sender     string           `protobuf:"bytes,1,opt,name=Sender,proto3" json:"Sender"`
	Receiver   uint64           `protobuf:"varint,2,opt,name=Receiver,proto3" json:"Receiver"`
	Id         uint32           `protobuf:"varint,3,opt,name=Id,proto3" json:"Id"`
	SendTime   uint64           `protobuf:"varint,4,opt,name=SendTime,proto3" json:"SendTime"`
	ExpireTime uint64           `protobuf:"varint,5,opt,name=ExpireTime,proto3" json:"ExpireTime"`
	AwardTime  uint64           `protobuf:"varint,6,opt,name=AwardTime,proto3" json:"AwardTime"`
	Title      string           `protobuf:"bytes,7,opt,name=Title,proto3" json:"Title"`
	Content    string           `protobuf:"bytes,8,opt,name=Content,proto3" json:"Content"`
	State      EmMailState      `protobuf:"varint,9,opt,name=State,proto3,enum=common.EmMailState" json:"State"`
	Item       []*PBAddItemData `protobuf:"bytes,10,rep,name=Item,proto3" json:"Item"`
}

type PBNotice struct {
	Id        uint32 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`
	Title     string `protobuf:"bytes,2,opt,name=Title,proto3" json:"Title"`
	Content   string `protobuf:"bytes,3,opt,name=Content,proto3" json:"Content"`
	BeginTime uint64 `protobuf:"varint,4,opt,name=BeginTime,proto3" json:"BeginTime"`
	EndTime   uint64 `protobuf:"varint,5,opt,name=EndTime,proto3" json:"EndTime"`
	Sender    string `protobuf:"bytes,6,opt,name=Sender,proto3" json:"Sender"`
}

type PBOfflineData struct {
	OfflineType EmPlayerOfflineType `protobuf:"varint,1,opt,name=OfflineType,proto3,enum=common.EmPlayerOfflineType" json:"OfflineType"`
	Mail        *PBMail             `protobuf:"bytes,2,opt,name=Mail,proto3" json:"Mail"`
	Item        []*PBAddItemData    `protobuf:"bytes,3,rep,name=Item,proto3" json:"Item"`
	DoingType   EmDoingType         `protobuf:"varint,4,opt,name=DoingType,proto3,enum=common.EmDoingType" json:"DoingType"`
	Notify      bool                `protobuf:"varint,5,opt,name=Notify,proto3" json:"Notify"`
}

type PBOpenServerGiftInfo struct {
	GiftId    uint32      `protobuf:"varint,1,opt,name=GiftId,proto3" json:"GiftId"`
	BeginTime uint64      `protobuf:"varint,2,opt,name=BeginTime,proto3" json:"BeginTime"`
	EndTime   uint64      `protobuf:"varint,3,opt,name=EndTime,proto3" json:"EndTime"`
	StageList []*PBU32U32 `protobuf:"bytes,4,rep,name=StageList,proto3" json:"StageList"`
}

type PBPlayerActivityInfo struct {
	ActivityId uint32 `protobuf:"varint,1,opt,name=ActivityId,proto3" json:"ActivityId"`
	BeginTime  uint64 `protobuf:"varint,2,opt,name=BeginTime,proto3" json:"BeginTime"`
	EndTime    uint64 `protobuf:"varint,3,opt,name=EndTime,proto3" json:"EndTime"`
}

type PBPlayerBag struct {
	ItemList     []*PBItem   `protobuf:"bytes,1,rep,name=ItemList,proto3" json:"ItemList"`
	DailyBuyItem []*PBU32U32 `protobuf:"bytes,2,rep,name=DailyBuyItem,proto3" json:"DailyBuyItem"`
}

type PBPlayerBase struct {
	Display           *PBPlayerBaseDisplay `protobuf:"bytes,1,opt,name=Display,proto3" json:"Display"`
	CreateTime        uint64               `protobuf:"varint,2,opt,name=CreateTime,proto3" json:"CreateTime"`
	LoginState        LoginState           `protobuf:"varint,3,opt,name=LoginState,proto3,enum=common.LoginState" json:"LoginState"`
	LastDailyTime     uint64               `protobuf:"varint,5,opt,name=LastDailyTime,proto3" json:"LastDailyTime"`
	LastModifyTime    uint64               `protobuf:"varint,6,opt,name=LastModifyTime,proto3" json:"LastModifyTime"`
	NewPlayerTypeList []uint32             `protobuf:"varint,7,rep,packed,name=NewPlayerTypeList,proto3" json:"NewPlayerTypeList"`
	PlatType          uint32               `protobuf:"varint,8,opt,name=PlatType,proto3" json:"PlatType"`
	PlatSystemType    uint32               `protobuf:"varint,9,opt,name=PlatSystemType,proto3" json:"PlatSystemType"`
	AccountName       string               `protobuf:"bytes,10,opt,name=AccountName,proto3" json:"AccountName"`
	SeverStartTime    uint64               `protobuf:"varint,11,opt,name=SeverStartTime,proto3" json:"SeverStartTime"`
}

type PBPlayerBaseDisplay struct {
	AccountId     uint64 `protobuf:"varint,1,opt,name=AccountId,proto3" json:"AccountId"`
	PlayerName    string `protobuf:"bytes,2,opt,name=PlayerName,proto3" json:"PlayerName"`
	PlayerLevel   uint32 `protobuf:"varint,3,opt,name=PlayerLevel,proto3" json:"PlayerLevel"`
	VipLevel      uint32 `protobuf:"varint,4,opt,name=vipLevel,proto3" json:"vipLevel"`
	AvatarID      uint32 `protobuf:"varint,5,opt,name=AvatarID,proto3" json:"AvatarID"`
	AvatarFrameID uint32 `protobuf:"varint,6,opt,name=AvatarFrameID,proto3" json:"AvatarFrameID"`
	SeverId       uint32 `protobuf:"varint,7,opt,name=SeverId,proto3" json:"SeverId"`
}

type PBPlayerBattleData struct {
	Display    *PBPlayerBaseDisplay `protobuf:"bytes,1,opt,name=Display,proto3" json:"Display"`
	HeroList   []*PBHero            `protobuf:"bytes,2,rep,name=HeroList,proto3" json:"HeroList"`
	Time       uint64               `protobuf:"varint,3,opt,name=Time,proto3" json:"Time"`
	UseTime    uint32               `protobuf:"varint,4,opt,name=UseTime,proto3" json:"UseTime"`
	FightPower uint32               `protobuf:"varint,5,opt,name=FightPower,proto3" json:"FightPower"`
	ClientData *PBBattleClientData  `protobuf:"bytes,6,opt,name=ClientData,proto3" json:"ClientData"`
}

type PBPlayerClientData struct {
	ClientDataList []*PBClientData `protobuf:"bytes,1,rep,name=ClientDataList,proto3" json:"ClientDataList"`
}

type PBPlayerCrystal struct {
	Book       *PBCrystalBook    `protobuf:"bytes,1,opt,name=Book,proto3" json:"Book"`
	Robots     []*PBCrystalRobot `protobuf:"bytes,2,rep,name=Robots,proto3" json:"Robots"`
	Crystals   []*PBCrystal      `protobuf:"bytes,3,rep,name=Crystals,proto3" json:"Crystals"`
	Conditions []*EntryCondition `protobuf:"bytes,4,rep,name=Conditions,proto3" json:"Conditions"`
	Effects    []*EntryEffect    `protobuf:"bytes,5,rep,name=Effects,proto3" json:"Effects"`
}

type PBPlayerData struct {
	Base      *PBPlayerBase       `protobuf:"bytes,1,opt,name=Base,proto3" json:"Base"`
	System    *PBPlayerSystem     `protobuf:"bytes,2,opt,name=System,proto3" json:"System"`
	Bag       *PBPlayerBag        `protobuf:"bytes,3,opt,name=Bag,proto3" json:"Bag"`
	Equipment *PBPlayerEquipment  `protobuf:"bytes,4,opt,name=Equipment,proto3" json:"Equipment"`
	Client    *PBPlayerClientData `protobuf:"bytes,5,opt,name=Client,proto3" json:"Client"`
	Hero      *PBPlayerHero       `protobuf:"bytes,6,opt,name=Hero,proto3" json:"Hero"`
	Mail      *PBPlayerMail       `protobuf:"bytes,7,opt,name=Mail,proto3" json:"Mail"`
	Crystal   *PBPlayerCrystal    `protobuf:"bytes,8,opt,name=Crystal,proto3" json:"Crystal"`
}

type PBPlayerEquipment struct {
	OrderId           uint32         `protobuf:"varint,1,opt,name=OrderId,proto3" json:"OrderId"`
	EquipmentList     []*PBEquipment `protobuf:"bytes,2,rep,name=equipmentList,proto3" json:"equipmentList"`
	AutoSplitQuality  []uint32       `protobuf:"varint,3,rep,packed,name=AutoSplitQuality,proto3" json:"AutoSplitQuality"`
	PosBuyCount       uint32         `protobuf:"varint,4,opt,name=PosBuyCount,proto3" json:"PosBuyCount"`
	SplitScore        uint32         `protobuf:"varint,5,opt,name=SplitScore,proto3" json:"SplitScore"`
	SplitAddBoxCount  uint32         `protobuf:"varint,6,opt,name=SplitAddBoxCount,proto3" json:"SplitAddBoxCount"`
	HookEquipmentList []*PBEquipment `protobuf:"bytes,7,rep,name=HookEquipmentList,proto3" json:"HookEquipmentList"`
	SplitEquipCount   uint32         `protobuf:"varint,8,opt,name=SplitEquipCount,proto3" json:"SplitEquipCount"`
	TotalSplitScore   uint64         `protobuf:"varint,9,opt,name=TotalSplitScore,proto3" json:"TotalSplitScore"`
}

type PBPlayerGiftCode struct {
	Acode string `protobuf:"bytes,1,opt,name=Acode,proto3" json:"Acode"`
	Count uint32 `protobuf:"varint,2,opt,name=Count,proto3" json:"Count"`
	Time  uint64 `protobuf:"varint,3,opt,name=Time,proto3" json:"Time"`
}

type PBPlayerHero struct {
	OrderId              uint32            `protobuf:"varint,1,opt,name=OrderId,proto3" json:"OrderId"`
	HeroList             []*PBHero         `protobuf:"bytes,2,rep,name=HeroList,proto3" json:"HeroList"`
	TeamList             []*PBHeroTeamList `protobuf:"bytes,3,rep,name=TeamList,proto3" json:"TeamList"`
	FightPower           uint32            `protobuf:"varint,4,opt,name=FightPower,proto3" json:"FightPower"`
	UpStarCount          []*PBU32U32       `protobuf:"bytes,5,rep,name=UpStarCount,proto3" json:"UpStarCount"`
	HeroBookList         []*PBHeroBook     `protobuf:"bytes,6,rep,name=HeroBookList,proto3" json:"HeroBookList"`
	MaxHistoryFightPower uint32            `protobuf:"varint,7,opt,name=MaxHistoryFightPower,proto3" json:"MaxHistoryFightPower"`
	GlobalRandHeroProf   []uint32          `protobuf:"varint,8,rep,packed,name=GlobalRandHeroProf,proto3" json:"GlobalRandHeroProf"`
}

type PBPlayerMail struct {
	MailList   []*PBMail `protobuf:"bytes,1,rep,name=MailList,proto3" json:"MailList"`
	OrderId    uint32    `protobuf:"varint,2,opt,name=OrderId,proto3" json:"OrderId"`
	AllOrderId uint32    `protobuf:"varint,3,opt,name=AllOrderId,proto3" json:"AllOrderId"`
}

type PBPlayerSystem struct {
	Common       *PBPlayerSystemCommon       `protobuf:"bytes,1,opt,name=Common,proto3" json:"Common"`
	Task         *PBPlayerSystemTask         `protobuf:"bytes,2,opt,name=Task,proto3" json:"Task"`
	Prof         *PBPlayerSystemProfession   `protobuf:"bytes,3,opt,name=Prof,proto3" json:"Prof"`
	Battle       *PBPlayerSystemBattle       `protobuf:"bytes,4,opt,name=Battle,proto3" json:"Battle"`
	Box          *PBPlayerSystemBox          `protobuf:"bytes,5,opt,name=Box,proto3" json:"Box"`
	Shop         *PBPlayerSystemShop         `protobuf:"bytes,6,opt,name=Shop,proto3" json:"Shop"`
	Draw         *PBPlayerSystemDraw         `protobuf:"bytes,7,opt,name=Draw,proto3" json:"Draw"`
	Charge       *PBPlayerSystemCharge       `protobuf:"bytes,8,opt,name=Charge,proto3" json:"Charge"`
	Gene         *PBPlayerSystemGene         `protobuf:"bytes,9,opt,name=Gene,proto3" json:"Gene"`
	Offline      *PBPlayerSystemOffline      `protobuf:"bytes,10,opt,name=Offline,proto3" json:"Offline"`
	HookTech     *PBPlayerSystemHookTech     `protobuf:"bytes,11,opt,name=HookTech,proto3" json:"HookTech"`
	SevenDay     *PBPlayerSystemSevenDay     `protobuf:"bytes,12,opt,name=SevenDay,proto3" json:"SevenDay"`
	WorldBoss    *PBPlayerSystemWorldBoss    `protobuf:"bytes,13,opt,name=WorldBoss,proto3" json:"WorldBoss"`
	Championship *PBPlayerSystemChampionship `protobuf:"bytes,14,opt,name=Championship,proto3" json:"Championship"`
	Activity     *PBPlayerSystemActivity     `protobuf:"bytes,15,opt,name=Activity,proto3" json:"Activity"`
	RepairData   *PBPlayerSystemRepairData   `protobuf:"bytes,16,opt,name=RepairData,proto3" json:"RepairData"`
}

type PBPlayerSystemActivity struct {
	ActivityList       []*PBPlayerActivityInfo     `protobuf:"bytes,1,rep,name=ActivityList,proto3" json:"ActivityList"`
	GrowRoadList       []*PBActivityGrowRoadInfo   `protobuf:"bytes,2,rep,name=GrowRoadList,proto3" json:"GrowRoadList"`
	GiftList           []*PBActivityChargeGift     `protobuf:"bytes,3,rep,name=GiftList,proto3" json:"GiftList"`
	AdventureList      []*PBActivityAdventure      `protobuf:"bytes,4,rep,name=AdventureList,proto3" json:"AdventureList"`
	OpenServerGiftList []*PBActivityOpenServerGift `protobuf:"bytes,5,rep,name=OpenServerGiftList,proto3" json:"OpenServerGiftList"`
}

type PBPlayerSystemBattle struct {
	BattleNormal  *PBBattleNormalInfo `protobuf:"bytes,1,opt,name=BattleNormal,proto3" json:"BattleNormal"`
	BattleHook    *PBBattleHookInfo   `protobuf:"bytes,2,opt,name=BattleHook,proto3" json:"BattleHook"`
	Battlechedule *PBBattleSchedule   `protobuf:"bytes,3,opt,name=Battlechedule,proto3" json:"Battlechedule"`
	Relive        *PBBattleRelive     `protobuf:"bytes,4,opt,name=Relive,proto3" json:"Relive"`
}

type PBPlayerSystemBox struct {
	BoxScore     uint32       `protobuf:"varint,1,opt,name=BoxScore,proto3" json:"BoxScore"`
	CurrentLevel uint32       `protobuf:"varint,2,opt,name=CurrentLevel,proto3" json:"CurrentLevel"`
	RecycleTimes uint32       `protobuf:"varint,3,opt,name=RecycleTimes,proto3" json:"RecycleTimes"`
	Boxs         []*PBBoxInfo `protobuf:"bytes,4,rep,name=Boxs,proto3" json:"Boxs"`
}

type PBPlayerSystemChampionship struct {
	HookHasReward   uint32           `protobuf:"varint,1,opt,name=HookHasReward,proto3" json:"HookHasReward"`
	BattleHasReward uint32           `protobuf:"varint,2,opt,name=BattleHasReward,proto3" json:"BattleHasReward"`
	DamageHasReward uint32           `protobuf:"varint,3,opt,name=DamageHasReward,proto3" json:"DamageHasReward"`
	PowerHasReward  uint32           `protobuf:"varint,4,opt,name=PowerHasReward,proto3" json:"PowerHasReward"`
	Hook            *PBTaskStageInfo `protobuf:"bytes,5,opt,name=Hook,proto3" json:"Hook"`
	Battle          *PBTaskStageInfo `protobuf:"bytes,6,opt,name=Battle,proto3" json:"Battle"`
	Damage          *PBTaskStageInfo `protobuf:"bytes,7,opt,name=Damage,proto3" json:"Damage"`
	Power           *PBTaskStageInfo `protobuf:"bytes,8,opt,name=Power,proto3" json:"Power"`
}

type PBPlayerSystemCharge struct {
	Charge          *PBCharge        `protobuf:"bytes,1,opt,name=Charge,proto3" json:"Charge"`
	FirstChargeList []*PBFirstCharge `protobuf:"bytes,2,rep,name=FirstChargeList,proto3" json:"FirstChargeList"`
	BPList          []*PBBPInfo      `protobuf:"bytes,3,rep,name=BPList,proto3" json:"BPList"`
	CardList        []*PBChargeCard  `protobuf:"bytes,4,rep,name=CardList,proto3" json:"CardList"`
}

type PBPlayerSystemCommon struct {
	LastChatTime        uint64              `protobuf:"varint,1,opt,name=LastChatTime,proto3" json:"LastChatTime"`
	GiftCode            []*PBPlayerGiftCode `protobuf:"bytes,2,rep,name=GiftCode,proto3" json:"GiftCode"`
	SystemOpenIds       []uint32            `protobuf:"varint,3,rep,packed,name=SystemOpenIds,proto3" json:"SystemOpenIds"`
	Avatars             []*PBAvatar         `protobuf:"bytes,4,rep,name=Avatars,proto3" json:"Avatars"`
	AvatarFrames        []*PBAvatarFrame    `protobuf:"bytes,5,rep,name=AvatarFrames,proto3" json:"AvatarFrames"`
	MaxNoticeId         uint32              `protobuf:"varint,6,opt,name=MaxNoticeId,proto3" json:"MaxNoticeId"`
	AdvertList          []*PBAdvertInfo     `protobuf:"bytes,7,rep,name=AdvertList,proto3" json:"AdvertList"`
	SystemOpenPrizeList []uint32            `protobuf:"varint,8,rep,packed,name=SystemOpenPrizeList,proto3" json:"SystemOpenPrizeList"`
	EjectAdvertInfo     *PBEjectAdvertInfo  `protobuf:"bytes,9,opt,name=EjectAdvertInfo,proto3" json:"EjectAdvertInfo"`
}

type PBPlayerSystemDraw struct {
	DrawList []*PBDrawInfo `protobuf:"bytes,1,rep,name=DrawList,proto3" json:"DrawList"`
}

type PBPlayerSystemGene struct {
	SchemeID uint32          `protobuf:"varint,1,opt,name=SchemeID,proto3" json:"SchemeID"`
	Schemes  []*PBGeneScheme `protobuf:"bytes,2,rep,name=Schemes,proto3" json:"Schemes"`
}

type PBPlayerSystemHookTech struct {
	HookTechList []*PBHookTech `protobuf:"bytes,1,rep,name=HookTechList,proto3" json:"HookTechList"`
}

type PBPlayerSystemOffline struct {
	LoginTime           uint64           `protobuf:"varint,1,opt,name=LoginTime,proto3" json:"LoginTime"`
	LogoutTime          uint64           `protobuf:"varint,2,opt,name=LogoutTime,proto3" json:"LogoutTime"`
	IncomTime           uint32           `protobuf:"varint,3,opt,name=IncomTime,proto3" json:"IncomTime"`
	Rewards             []*PBAddItemData `protobuf:"bytes,4,rep,name=Rewards,proto3" json:"Rewards"`
	TotalEquipment      uint32           `protobuf:"varint,5,opt,name=TotalEquipment,proto3" json:"TotalEquipment"`
	AddEquipmentBag     uint32           `protobuf:"varint,6,opt,name=AddEquipmentBag,proto3" json:"AddEquipmentBag"`
	SplitEquipmentScore uint64           `protobuf:"varint,7,opt,name=SplitEquipmentScore,proto3" json:"SplitEquipmentScore"`
	MaxIncomTime        uint32           `protobuf:"varint,8,opt,name=MaxIncomTime,proto3" json:"MaxIncomTime"`
}

type PBPlayerSystemProfInfo struct {
	ProfType     uint32                        `protobuf:"varint,1,opt,name=ProfType,proto3" json:"ProfType"`
	Level        uint32                        `protobuf:"varint,2,opt,name=Level,proto3" json:"Level"`
	Grade        uint32                        `protobuf:"varint,3,opt,name=Grade,proto3" json:"Grade"`
	PeakLevel    uint32                        `protobuf:"varint,4,opt,name=PeakLevel,proto3" json:"PeakLevel"`
	PartList     []*PBPlayerSystemProfPartInfo `protobuf:"bytes,5,rep,name=PartList,proto3" json:"PartList"`
	LastLinkStar uint32                        `protobuf:"varint,6,opt,name=LastLinkStar,proto3" json:"LastLinkStar"`
}

type PBPlayerSystemProfPartInfo struct {
	Part       uint32 `protobuf:"varint,1,opt,name=Part,proto3" json:"Part"`
	Level      uint32 `protobuf:"varint,2,opt,name=Level,proto3" json:"Level"`
	EquipSn    uint32 `protobuf:"varint,3,opt,name=EquipSn,proto3" json:"EquipSn"`
	Refine     uint32 `protobuf:"varint,4,opt,name=Refine,proto3" json:"Refine"`
	RefineTupo uint32 `protobuf:"varint,5,opt,name=RefineTupo,proto3" json:"RefineTupo"`
}

type PBPlayerSystemProfession struct {
	ProfList []*PBPlayerSystemProfInfo `protobuf:"bytes,1,rep,name=ProfList,proto3" json:"ProfList"`
}

type PBPlayerSystemRepairData struct {
	Version     uint32 `protobuf:"varint,1,opt,name=Version,proto3" json:"Version"`
	VersionTime uint64 `protobuf:"varint,2,opt,name=VersionTime,proto3" json:"VersionTime"`
}

type PBPlayerSystemSevenDay struct {
	SevenDayList []*PBSevenDayInfo `protobuf:"bytes,1,rep,name=SevenDayList,proto3" json:"SevenDayList"`
}

type PBPlayerSystemShop struct {
	BlackShop *PBBlackShop  `protobuf:"bytes,1,opt,name=BlackShop,proto3" json:"BlackShop"`
	ShopList  []*PBShopInfo `protobuf:"bytes,2,rep,name=ShopList,proto3" json:"ShopList"`
}

type PBPlayerSystemTask struct {
	MainTask    *PBTaskStageInfo `protobuf:"bytes,1,opt,name=MainTask,proto3" json:"MainTask"`
	DailyTask   *PBDailyTask     `protobuf:"bytes,2,opt,name=DailyTask,proto3" json:"DailyTask"`
	AchieveList []*PBAchieveInfo `protobuf:"bytes,3,rep,name=AchieveList,proto3" json:"AchieveList"`
}

type PBPlayerSystemWorldBoss struct {
	BossId            uint32 `protobuf:"varint,1,opt,name=BossId,proto3" json:"BossId"`
	DailyPrizeStageId uint32 `protobuf:"varint,2,opt,name=DailyPrizeStageId,proto3" json:"DailyPrizeStageId"`
	DailyMaxDamage    uint64 `protobuf:"varint,3,opt,name=DailyMaxDamage,proto3" json:"DailyMaxDamage"`
	DailyTotalDamage  uint64 `protobuf:"varint,4,opt,name=DailyTotalDamage,proto3" json:"DailyTotalDamage"`
	DailyBuyCount     uint32 `protobuf:"varint,5,opt,name=DailyBuyCount,proto3" json:"DailyBuyCount"`
	DailyEnterCount   uint32 `protobuf:"varint,6,opt,name=DailyEnterCount,proto3" json:"DailyEnterCount"`
	DailyPrizeCount   uint32 `protobuf:"varint,7,opt,name=DailyPrizeCount,proto3" json:"DailyPrizeCount"`
	MaxDamage         uint64 `protobuf:"varint,8,opt,name=MaxDamage,proto3" json:"MaxDamage"`
}

type PBPropInfo struct {
	PropId uint32 `protobuf:"varint,1,opt,name=PropId,proto3" json:"PropId"`
	Value  uint32 `protobuf:"varint,2,opt,name=Value,proto3" json:"Value"`
}

type PBSevenDayInfo struct {
	Id         uint32             `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`
	BeginTime  uint64             `protobuf:"varint,2,opt,name=BeginTime,proto3" json:"BeginTime"`
	EndTime    uint64             `protobuf:"varint,3,opt,name=EndTime,proto3" json:"EndTime"`
	Value      uint32             `protobuf:"varint,4,opt,name=Value,proto3" json:"Value"`
	PrizeValue uint32             `protobuf:"varint,5,opt,name=PrizeValue,proto3" json:"PrizeValue"`
	TaskList   []*PBTaskStageInfo `protobuf:"bytes,6,rep,name=TaskList,proto3" json:"TaskList"`
	GiftList   []*PBU32U32        `protobuf:"bytes,7,rep,name=GiftList,proto3" json:"GiftList"`
}

type PBShopGoodCfg struct {
	GoodsID     uint32           `protobuf:"varint,1,opt,name=GoodsID,proto3" json:"GoodsID"`
	BuyTimes    uint32           `protobuf:"varint,2,opt,name=BuyTimes,proto3" json:"BuyTimes"`
	MaxTimes    uint32           `protobuf:"varint,3,opt,name=MaxTimes,proto3" json:"MaxTimes"`
	Discount    uint32           `protobuf:"varint,4,opt,name=Discount,proto3" json:"Discount"`
	NeedItem    *PBAddItem       `protobuf:"bytes,5,opt,name=NeedItem,proto3" json:"NeedItem"`
	AddItem     []*PBAddItemData `protobuf:"bytes,6,rep,name=AddItem,proto3" json:"AddItem"`
	ProductId   uint32           `protobuf:"varint,7,opt,name=ProductId,proto3" json:"ProductId"`
	ProductName string           `protobuf:"bytes,8,opt,name=ProductName,proto3" json:"ProductName"`
	Price       uint32           `protobuf:"varint,9,opt,name=Price,proto3" json:"Price"`
	ValueTips   string           `protobuf:"bytes,10,opt,name=ValueTips,proto3" json:"ValueTips"`
	SortTag     uint32           `protobuf:"varint,11,opt,name=SortTag,proto3" json:"SortTag"`
}

type PBShopGoodInfo struct {
	GoodsID   uint32       `protobuf:"varint,1,opt,name=GoodsID,proto3" json:"GoodsID"`
	Discount  uint32       `protobuf:"varint,2,opt,name=Discount,proto3" json:"Discount"`
	BuyTimes  uint32       `protobuf:"varint,3,opt,name=BuyTimes,proto3" json:"BuyTimes"`
	FreeTimes uint32       `protobuf:"varint,4,opt,name=FreeTimes,proto3" json:"FreeTimes"`
	Equipment *PBEquipment `protobuf:"bytes,5,opt,name=Equipment,proto3" json:"Equipment"`
}

type PBShopInfo struct {
	ShopType        uint32      `protobuf:"varint,1,opt,name=ShopType,proto3" json:"ShopType"`
	NextRefreshTime uint64      `protobuf:"varint,2,opt,name=NextRefreshTime,proto3" json:"NextRefreshTime"`
	Items           []*PBU32U32 `protobuf:"bytes,3,rep,name=Items,proto3" json:"Items"`
	HaveRed         uint32      `protobuf:"varint,4,opt,name=HaveRed,proto3" json:"HaveRed"`
}

type PBShopRefreshInfo struct {
	DailyBuyCount       uint32 `protobuf:"varint,1,opt,name=DailyBuyCount,proto3" json:"DailyBuyCount"`
	DailyFreeUseCount   uint32 `protobuf:"varint,2,opt,name=DailyFreeUseCount,proto3" json:"DailyFreeUseCount"`
	NextFreeRefreshTime uint64 `protobuf:"varint,3,opt,name=NextFreeRefreshTime,proto3" json:"NextFreeRefreshTime"`
	DailyFreeMaxCount   uint32 `protobuf:"varint,4,opt,name=DailyFreeMaxCount,proto3" json:"DailyFreeMaxCount"`
}

type PBStringInt64 struct {
	Key   string `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key"`
	Value int64  `protobuf:"varint,2,opt,name=value,proto3" json:"value"`
}

type PBTaskStageInfo struct {
	Id       uint32      `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`
	Value    uint32      `protobuf:"varint,2,opt,name=Value,proto3" json:"Value"`
	MaxValue uint32      `protobuf:"varint,3,opt,name=MaxValue,proto3" json:"MaxValue"`
	State    EmTaskState `protobuf:"varint,4,opt,name=State,proto3,enum=common.EmTaskState" json:"State"`
}

type PBU32U32 struct {
	Key   uint32 `protobuf:"varint,1,opt,name=Key,proto3" json:"Key"`
	Value uint32 `protobuf:"varint,2,opt,name=Value,proto3" json:"Value"`
}

type PBU32U64 struct {
	Key   uint32 `protobuf:"varint,1,opt,name=Key,proto3" json:"Key"`
	Value uint64 `protobuf:"varint,2,opt,name=Value,proto3" json:"Value"`
}

type Packet struct {
	Id        uint32     `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`
	Reply     string     `protobuf:"bytes,2,opt,name=Reply,proto3" json:"Reply"`
	Buff      []byte     `protobuf:"bytes,3,opt,name=Buff,proto3" json:"Buff"`
	RpcPacket *RpcPacket `protobuf:"bytes,4,opt,name=RpcPacket,proto3" json:"RpcPacket"`
}

type PageOpenRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PageType   uint32   `protobuf:"varint,2,opt,name=PageType,proto3" json:"PageType"`
}

type PageOpenResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type PassTimeNotify struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	IsDay      bool     `protobuf:"varint,2,opt,name=IsDay,proto3" json:"IsDay"`
	IsWeek     bool     `protobuf:"varint,3,opt,name=IsWeek,proto3" json:"IsWeek"`
	IsMonth    bool     `protobuf:"varint,4,opt,name=IsMonth,proto3" json:"IsMonth"`
	CurTime    uint64   `protobuf:"varint,5,opt,name=CurTime,proto3" json:"CurTime"`
}

type PlayerData struct {
	PlayerID   uint64 `protobuf:"varint,1,opt,name=PlayerID,proto3" json:"PlayerID"`
	PlayerName string `protobuf:"bytes,2,opt,name=PlayerName,proto3" json:"PlayerName"`
	PlayerGold int32  `protobuf:"varint,3,opt,name=PlayerGold,proto3" json:"PlayerGold"`
	Ok         bool   `protobuf:"varint,4,opt,name=Ok,proto3" json:"Ok"`
}

type PlayerUpdateKvNotify struct {
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ListInfo   []*PBStringInt64 `protobuf:"bytes,2,rep,name=ListInfo,proto3" json:"ListInfo"`
}

type Point3F struct {
	X float32 `protobuf:"fixed32,1,opt,name=X,proto3" json:"X"`
	Y float32 `protobuf:"fixed32,2,opt,name=Y,proto3" json:"Y"`
	Z float32 `protobuf:"fixed32,3,opt,name=Z,proto3" json:"Z"`
}

type ProfessionEquipRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
	PartType   uint32   `protobuf:"varint,3,opt,name=PartType,proto3" json:"PartType"`
	Sn         uint32   `protobuf:"varint,4,opt,name=Sn,proto3" json:"Sn"`
}

type ProfessionEquipResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
	PartType   uint32   `protobuf:"varint,3,opt,name=PartType,proto3" json:"PartType"`
	Sn         uint32   `protobuf:"varint,4,opt,name=Sn,proto3" json:"Sn"`
}

type ProfessionGradeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
}

type ProfessionGradeResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
	Grade      uint32   `protobuf:"varint,3,opt,name=Grade,proto3" json:"Grade"`
}

type ProfessionLevelRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
	CurLevel   uint32   `protobuf:"varint,3,opt,name=CurLevel,proto3" json:"CurLevel"`
	AddLevel   uint32   `protobuf:"varint,4,opt,name=AddLevel,proto3" json:"AddLevel"`
}

type ProfessionLevelResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
	CurLevel   uint32   `protobuf:"varint,3,opt,name=CurLevel,proto3" json:"CurLevel"`
}

type ProfessionOnekeyUnEquipRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
}

type ProfessionOnekeyUnEquipResponse struct {
	PacketHead  *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType    uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
	EquipSnList []uint32 `protobuf:"varint,3,rep,packed,name=EquipSnList,proto3" json:"EquipSnList"`
}

type ProfessionPartLevelRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
	PartType   uint32   `protobuf:"varint,3,opt,name=PartType,proto3" json:"PartType"`
	CurLevel   uint32   `protobuf:"varint,4,opt,name=CurLevel,proto3" json:"CurLevel"`
}

type ProfessionPartLevelResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
	PartType   uint32   `protobuf:"varint,3,opt,name=PartType,proto3" json:"PartType"`
	CurLevel   uint32   `protobuf:"varint,4,opt,name=CurLevel,proto3" json:"CurLevel"`
}

type ProfessionPartOnekeyLevelRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
}

type ProfessionPartOnekeyLevelResponse struct {
	PacketHead *IPacket                      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32                        `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
	PartList   []*PBPlayerSystemProfPartInfo `protobuf:"bytes,3,rep,name=PartList,proto3" json:"PartList"`
}

type ProfessionPartOnekeyRefineRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
}

type ProfessionPartOnekeyRefineResponse struct {
	PacketHead *IPacket                      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32                        `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
	PartList   []*PBPlayerSystemProfPartInfo `protobuf:"bytes,3,rep,name=PartList,proto3" json:"PartList"`
}

type ProfessionPartRefineRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
	PartType   uint32   `protobuf:"varint,3,opt,name=PartType,proto3" json:"PartType"`
	CurLevel   uint32   `protobuf:"varint,4,opt,name=CurLevel,proto3" json:"CurLevel"`
}

type ProfessionPartRefineResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
	PartType   uint32   `protobuf:"varint,3,opt,name=PartType,proto3" json:"PartType"`
	CurLevel   uint32   `protobuf:"varint,4,opt,name=CurLevel,proto3" json:"CurLevel"`
}

type ProfessionPartRefineTupoRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
	PartType   uint32   `protobuf:"varint,3,opt,name=PartType,proto3" json:"PartType"`
	CurLevel   uint32   `protobuf:"varint,4,opt,name=CurLevel,proto3" json:"CurLevel"`
}

type ProfessionPartRefineTupoResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
	PartType   uint32   `protobuf:"varint,3,opt,name=PartType,proto3" json:"PartType"`
	CurLevel   uint32   `protobuf:"varint,4,opt,name=CurLevel,proto3" json:"CurLevel"`
}

type ProfessionPartResetRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
	PartType   uint32   `protobuf:"varint,3,opt,name=PartType,proto3" json:"PartType"`
}

type ProfessionPartResetResponse struct {
	PacketHead *IPacket                      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32                        `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
	PartList   []*PBPlayerSystemProfPartInfo `protobuf:"bytes,3,rep,name=PartList,proto3" json:"PartList"`
}

type ProfessionPeakLevelRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
	CurLevel   uint32   `protobuf:"varint,3,opt,name=CurLevel,proto3" json:"CurLevel"`
}

type ProfessionPeakLevelResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
	CurLevel   uint32   `protobuf:"varint,3,opt,name=CurLevel,proto3" json:"CurLevel"`
}

type ProfessionPeakRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
}

type ProfessionPeakResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`
	PeakLevel  uint32   `protobuf:"varint,3,opt,name=PeakLevel,proto3" json:"PeakLevel"`
}

type ProtocolNameNotify struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Name       string   `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name"`
	ProtocolId uint32   `protobuf:"varint,3,opt,name=ProtocolId,proto3" json:"ProtocolId"`
}

type RankData struct {
	RankType   uint32   `protobuf:"varint,1,opt,name=RankType,proto3" json:"RankType"`
	RegionID   uint32   `protobuf:"varint,2,opt,name=RegionID,proto3" json:"RegionID"`
	CreateTime uint64   `protobuf:"varint,3,opt,name=CreateTime,proto3" json:"CreateTime"`
	HasRewards []uint64 `protobuf:"varint,7,rep,packed,name=HasRewards,proto3" json:"HasRewards"`
}

type RankInfo struct {
	Rank    uint32               `protobuf:"varint,1,opt,name=Rank,proto3" json:"Rank"`
	Display *PBPlayerBaseDisplay `protobuf:"bytes,2,opt,name=Display,proto3" json:"Display"`
	Value   uint64               `protobuf:"varint,3,opt,name=Value,proto3" json:"Value"`
}

type RankRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RankType   uint32   `protobuf:"varint,2,opt,name=RankType,proto3" json:"RankType"`
	Begin      uint32   `protobuf:"varint,3,opt,name=Begin,proto3" json:"Begin"`
	End        uint32   `protobuf:"varint,4,opt,name=End,proto3" json:"End"`
	CreateTime uint64   `protobuf:"varint,5,opt,name=CreateTime,proto3" json:"CreateTime"`
}

type RankResponse struct {
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RankType   uint32      `protobuf:"varint,2,opt,name=RankType,proto3" json:"RankType"`
	Begin      uint32      `protobuf:"varint,3,opt,name=Begin,proto3" json:"Begin"`
	End        uint32      `protobuf:"varint,4,opt,name=End,proto3" json:"End"`
	RankList   []*RankInfo `protobuf:"bytes,5,rep,name=RankList,proto3" json:"RankList"`
	SelfInfo   *RankInfo   `protobuf:"bytes,6,opt,name=SelfInfo,proto3" json:"SelfInfo"`
	TotalRank  int64       `protobuf:"varint,7,opt,name=TotalRank,proto3" json:"TotalRank"`
}

type RankRewardRequest struct {
	PacketHead *IPacket     `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RankType   uint32       `protobuf:"varint,2,opt,name=RankType,proto3" json:"RankType"`
	Notify     bool         `protobuf:"varint,3,opt,name=Notify,proto3" json:"Notify"`
	Doing      EmDoingType  `protobuf:"varint,4,opt,name=Doing,proto3,enum=common.EmDoingType" json:"Doing"`
	Rewards    []*PBAddItem `protobuf:"bytes,5,rep,name=Rewards,proto3" json:"Rewards"`
	CreateTime uint64       `protobuf:"varint,6,opt,name=CreateTime,proto3" json:"CreateTime"`
}

type RankRewardResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type RankUpdateNotify struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RankType   uint32   `protobuf:"varint,2,opt,name=RankType,proto3" json:"RankType"`
	CreateTime uint64   `protobuf:"varint,3,opt,name=CreateTime,proto3" json:"CreateTime"`
	Member     string   `protobuf:"bytes,7,opt,name=Member,proto3" json:"Member"`
	Score      float64  `protobuf:"fixed64,8,opt,name=Score,proto3" json:"Score"`
}

type RouteInfo struct {
	Gate   uint32 `protobuf:"varint,1,opt,name=Gate,proto3" json:"Gate"`
	Game   uint32 `protobuf:"varint,2,opt,name=Game,proto3" json:"Game"`
	Gm     uint32 `protobuf:"varint,3,opt,name=Gm,proto3" json:"Gm"`
	Db     uint32 `protobuf:"varint,4,opt,name=Db,proto3" json:"Db"`
	Dip    uint32 `protobuf:"varint,5,opt,name=Dip,proto3" json:"Dip"`
	Record uint32 `protobuf:"varint,6,opt,name=Record,proto3" json:"Record"`
}

type RpcHead struct {
	SrcClusterId   uint32     `protobuf:"varint,1,opt,name=SrcClusterId,proto3" json:"SrcClusterId"`
	ClusterId      uint32     `protobuf:"varint,2,opt,name=ClusterId,proto3" json:"ClusterId"`
	DestServerType SERVICE    `protobuf:"varint,3,opt,name=DestServerType,proto3,enum=common.SERVICE" json:"DestServerType"`
	RouteType      uint32     `protobuf:"varint,4,opt,name=RouteType,proto3" json:"RouteType"`
	RegionID       uint32     `protobuf:"varint,5,opt,name=RegionID,proto3" json:"RegionID"`
	Id             uint64     `protobuf:"varint,6,opt,name=Id,proto3" json:"Id"`
	SocketId       uint32     `protobuf:"varint,7,opt,name=SocketId,proto3" json:"SocketId"`
	SendType       SEND       `protobuf:"varint,8,opt,name=SendType,proto3,enum=rpc3.SEND" json:"SendType"`
	SeqId          uint32     `protobuf:"varint,9,opt,name=SeqId,proto3" json:"SeqId"`
	ActorName      string     `protobuf:"bytes,10,opt,name=ActorName,proto3" json:"ActorName"`
	FuncName       string     `protobuf:"bytes,11,opt,name=FuncName,proto3" json:"FuncName"`
	Reply          string     `protobuf:"bytes,12,opt,name=Reply,proto3" json:"Reply"`
	Route          *RouteInfo `protobuf:"bytes,13,opt,name=Route,proto3" json:"Route"`
}

type RpcPacket struct {
	RpcHead *RpcHead `protobuf:"bytes,1,opt,name=RpcHead,proto3" json:"RpcHead"`
	ArgLen  uint32   `protobuf:"varint,2,opt,name=ArgLen,proto3" json:"ArgLen"`
	RpcBody []byte   `protobuf:"bytes,3,opt,name=RpcBody,proto3" json:"RpcBody"`
}

type SetClientRequest struct {
	PacketHead *IPacket      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ClientData *PBClientData `protobuf:"bytes,2,opt,name=ClientData,proto3" json:"ClientData"`
}

type SetClientResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type SetPlayerDataRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	DataType   int32    `protobuf:"varint,2,opt,name=DataType,proto3" json:"DataType"`
	JsonData   string   `protobuf:"bytes,3,opt,name=JsonData,proto3" json:"JsonData"`
}

type SetPlayerDataResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type SevenDayActivePrizeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
}

type SevenDayActivePrizeResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	PrizeValue uint32   `protobuf:"varint,3,opt,name=PrizeValue,proto3" json:"PrizeValue"`
}

type SevenDayBuyGiftRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	GiftId     uint32   `protobuf:"varint,3,opt,name=GiftId,proto3" json:"GiftId"`
}

type SevenDayBuyGiftResponse struct {
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32    `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	Value      uint32    `protobuf:"varint,3,opt,name=Value,proto3" json:"Value"`
	GiftInfo   *PBU32U32 `protobuf:"bytes,4,opt,name=GiftInfo,proto3" json:"GiftInfo"`
}

type SevenDayGiftNotify struct {
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	GiftInfo   *PBU32U32 `protobuf:"bytes,2,opt,name=GiftInfo,proto3" json:"GiftInfo"`
	Value      uint32    `protobuf:"varint,3,opt,name=Value,proto3" json:"Value"`
}

type SevenDayListNotify struct {
	PacketHead *IPacket          `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	AddList    []*PBSevenDayInfo `protobuf:"bytes,2,rep,name=AddList,proto3" json:"AddList"`
	Delist     []uint32          `protobuf:"varint,3,rep,packed,name=Delist,proto3" json:"Delist"`
}

type SevenDayTaskPrizeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	TaskId     uint32   `protobuf:"varint,3,opt,name=TaskId,proto3" json:"TaskId"`
}

type SevenDayTaskPrizeResponse struct {
	PacketHead    *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id            uint32           `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	Value         uint32           `protobuf:"varint,3,opt,name=Value,proto3" json:"Value"`
	TaskStageInfo *PBTaskStageInfo `protobuf:"bytes,4,opt,name=TaskStageInfo,proto3" json:"TaskStageInfo"`
}

type ShopBuyRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopType   uint32   `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"`
	GoodsID    uint32   `protobuf:"varint,3,opt,name=GoodsID,proto3" json:"GoodsID"`
	AdvertType uint32   `protobuf:"varint,4,opt,name=AdvertType,proto3" json:"AdvertType"`
}

type ShopBuyResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopType   uint32   `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"`
}

type ShopExchangeGoodNotify struct {
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopType   uint32    `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"`
	GoodInfo   *PBU32U32 `protobuf:"bytes,3,opt,name=GoodInfo,proto3" json:"GoodInfo"`
}

type ShopExchangeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopType   uint32   `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"`
	GoodsID    uint32   `protobuf:"varint,3,opt,name=GoodsID,proto3" json:"GoodsID"`
}

type ShopExchangeResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopType   uint32   `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"`
	GoodsID    uint32   `protobuf:"varint,3,opt,name=GoodsID,proto3" json:"GoodsID"`
	BuyTimes   uint32   `protobuf:"varint,4,opt,name=BuyTimes,proto3" json:"BuyTimes"`
}

type ShopListNotify struct {
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopList   []*PBU32U64 `protobuf:"bytes,2,rep,name=ShopList,proto3" json:"ShopList"`
}

type ShopOpenRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopType   uint32   `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"`
}

type ShopOpenResponse struct {
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopType   uint32           `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"`
	GoodList   []*PBShopGoodCfg `protobuf:"bytes,3,rep,name=GoodList,proto3" json:"GoodList"`
}

type ShopRedNotify struct {
	PacketHead  *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopRedList []*PBU32U32 `protobuf:"bytes,2,rep,name=ShopRedList,proto3" json:"ShopRedList"`
}

type ShopRefreshRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopType   uint32   `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"`
	IsFree     bool     `protobuf:"varint,3,opt,name=IsFree,proto3" json:"IsFree"`
}

type ShopRefreshResponse struct {
	PacketHead  *IPacket           `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopType    uint32             `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"`
	RefreshInfo *PBShopRefreshInfo `protobuf:"bytes,3,opt,name=RefreshInfo,proto3" json:"RefreshInfo"`
}

type ShopRefreshTimeNotify struct {
	PacketHead  *IPacket           `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopType    uint32             `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"`
	RefreshInfo *PBShopRefreshInfo `protobuf:"bytes,3,opt,name=RefreshInfo,proto3" json:"RefreshInfo"`
}

type ShopUpdateNotify struct {
	PacketHead *IPacket            `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopType   uint32              `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"`
	Shop       *PBPlayerSystemShop `protobuf:"bytes,3,opt,name=Shop,proto3" json:"Shop"`
}

type ShopUpdateOneGoodsNotify struct {
	PacketHead *IPacket        `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopType   uint32          `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"`
	ShopGood   *PBShopGoodInfo `protobuf:"bytes,3,opt,name=ShopGood,proto3" json:"ShopGood"`
}

type StubMailBox struct {
	Id        uint64 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`
	LeaseId   uint64 `protobuf:"varint,2,opt,name=LeaseId,proto3" json:"LeaseId"`
	StubType  STUB   `protobuf:"varint,3,opt,name=StubType,proto3,enum=rpc3.STUB" json:"StubType"`
	ClusterId uint32 `protobuf:"varint,5,opt,name=ClusterId,proto3" json:"ClusterId"`
}

type SystemOpenNotify struct {
	PacketHead    *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SystemOpenIds []uint32 `protobuf:"varint,2,rep,packed,name=SystemOpenIds,proto3" json:"SystemOpenIds"`
}

type SystemOpenPrizeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
}

type SystemOpenPrizeResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
}

type SystemPropNotify struct {
	PacketHead     *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SyetemPropType EmSyetemPropType `protobuf:"varint,2,opt,name=SyetemPropType,proto3,enum=common.EmSyetemPropType" json:"SyetemPropType"`
	PropInfoList   []*PBPropInfo    `protobuf:"bytes,3,rep,name=PropInfoList,proto3" json:"PropInfoList"`
}

type UpdateRouteNotify struct {
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Users      []*UserRouteInfo `protobuf:"bytes,2,rep,name=Users,proto3" json:"Users"`
}

type UserRouteInfo struct {
	UID       uint64     `protobuf:"varint,1,opt,name=UID,proto3" json:"UID"`
	RouteInfo *RouteInfo `protobuf:"bytes,2,opt,name=RouteInfo,proto3" json:"RouteInfo"`
}

type WorldBossBattleBeginRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type WorldBossBattleBeginResponse struct {
	PacketHead      *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	DailyEnterCount uint32   `protobuf:"varint,2,opt,name=DailyEnterCount,proto3" json:"DailyEnterCount"`
}

type WorldBossBattleEndRequest struct {
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Battle     *BattleInfo `protobuf:"bytes,2,opt,name=Battle,proto3" json:"Battle"`
	IsFinish   uint32      `protobuf:"varint,3,opt,name=IsFinish,proto3" json:"IsFinish"`
}

type WorldBossBattleEndResponse struct {
	PacketHead       *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	DailyMaxDamage   uint64           `protobuf:"varint,2,opt,name=DailyMaxDamage,proto3" json:"DailyMaxDamage"`
	DailyTotalDamage uint64           `protobuf:"varint,3,opt,name=DailyTotalDamage,proto3" json:"DailyTotalDamage"`
	DailyPrizeCount  uint32           `protobuf:"varint,4,opt,name=DailyPrizeCount,proto3" json:"DailyPrizeCount"`
	ItemInfo         []*PBAddItemData `protobuf:"bytes,5,rep,name=ItemInfo,proto3" json:"ItemInfo"`
	DailyEnterCount  uint32           `protobuf:"varint,6,opt,name=DailyEnterCount,proto3" json:"DailyEnterCount"`
}

type WorldBossBuyCountRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type WorldBossBuyCountResponse struct {
	PacketHead    *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	DailyBuyCount uint32   `protobuf:"varint,2,opt,name=DailyBuyCount,proto3" json:"DailyBuyCount"`
}

type WorldBossNotify struct {
	PacketHead *IPacket                 `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	WorldBoss  *PBPlayerSystemWorldBoss `protobuf:"bytes,2,opt,name=WorldBoss,proto3" json:"WorldBoss"`
}

type WorldBossRecordRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	AccountId  uint64   `protobuf:"varint,2,opt,name=AccountId,proto3" json:"AccountId"`
	ServerId   uint32   `protobuf:"varint,3,opt,name=ServerId,proto3" json:"ServerId"`
}

type WorldBossRecordResponse struct {
	PacketHead *IPacket            `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RecordInfo *PBPlayerBattleData `protobuf:"bytes,2,opt,name=RecordInfo,proto3" json:"RecordInfo"`
}

type WorldBossStagePrizeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type WorldBossStagePrizeResponse struct {
	PacketHead        *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	DailyPrizeStageId uint32   `protobuf:"varint,2,opt,name=DailyPrizeStageId,proto3" json:"DailyPrizeStageId"`
}

type WorldBossSweepRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}

type WorldBossSweepResponse struct {
	PacketHead      *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	DailyEnterCount uint32   `protobuf:"varint,2,opt,name=DailyEnterCount,proto3" json:"DailyEnterCount"`
	DailyPrizeCount uint32   `protobuf:"varint,3,opt,name=DailyPrizeCount,proto3" json:"DailyPrizeCount"`
}

type Z_C_ENTITY struct {
	PacketHead *IPacket             `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	EntityInfo []*Z_C_ENTITY_Entity `protobuf:"bytes,2,rep,name=EntityInfo,proto3" json:"EntityInfo"`
}

type Z_C_ENTITY_Entity struct {
	Id    uint64                       `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`
	Data  *Z_C_ENTITY_Entity_DataMask  `protobuf:"bytes,2,opt,name=Data,proto3" json:"Data"`
	Move  *Z_C_ENTITY_Entity_MoveMask  `protobuf:"bytes,3,opt,name=Move,proto3" json:"Move"`
	Stats *Z_C_ENTITY_Entity_StatsMask `protobuf:"bytes,4,opt,name=Stats,proto3" json:"Stats"`
}

type Z_C_ENTITY_Entity_DataMask struct {
	Type       int32                                     `protobuf:"varint,1,opt,name=Type,proto3" json:"Type"`
	RemoveFlag bool                                      `protobuf:"varint,2,opt,name=RemoveFlag,proto3" json:"RemoveFlag"`
	NpcData    *Z_C_ENTITY_Entity_DataMask_NpcDataMask   `protobuf:"bytes,3,opt,name=NpcData,proto3" json:"NpcData"`
	SpellData  *Z_C_ENTITY_Entity_DataMask_SpellDataMask `protobuf:"bytes,4,opt,name=SpellData,proto3" json:"SpellData"`
}

type Z_C_ENTITY_Entity_DataMask_NpcDataMask struct {
	DataId int32 `protobuf:"varint,1,opt,name=DataId,proto3" json:"DataId"`
}

type Z_C_ENTITY_Entity_DataMask_SpellDataMask struct {
	DataId   int32 `protobuf:"varint,1,opt,name=DataId,proto3" json:"DataId"`
	LifeTime int32 `protobuf:"varint,2,opt,name=LifeTime,proto3" json:"LifeTime"`
	FlySpeed int32 `protobuf:"varint,3,opt,name=FlySpeed,proto3" json:"FlySpeed"`
}

type Z_C_ENTITY_Entity_MoveMask struct {
	Pos             *Point3F `protobuf:"bytes,1,opt,name=Pos,proto3" json:"Pos"`
	Rotation        float32  `protobuf:"fixed32,2,opt,name=Rotation,proto3" json:"Rotation"`
	ForceUpdateFlag bool     `protobuf:"varint,3,opt,name=ForceUpdateFlag,proto3" json:"ForceUpdateFlag"`
}

type Z_C_ENTITY_Entity_StatsMask struct {
	HP                int32 `protobuf:"varint,1,opt,name=HP,proto3" json:"HP"`
	MP                int32 `protobuf:"varint,2,opt,name=MP,proto3" json:"MP"`
	MaxHP             int32 `protobuf:"varint,3,opt,name=MaxHP,proto3" json:"MaxHP"`
	MaxMP             int32 `protobuf:"varint,4,opt,name=MaxMP,proto3" json:"MaxMP"`
	PhyDamage         int32 `protobuf:"varint,5,opt,name=PhyDamage,proto3" json:"PhyDamage"`
	PhyDefence        int32 `protobuf:"varint,6,opt,name=PhyDefence,proto3" json:"PhyDefence"`
	SplDamage         int32 `protobuf:"varint,7,opt,name=SplDamage,proto3" json:"SplDamage"`
	SplDefence        int32 `protobuf:"varint,8,opt,name=SplDefence,proto3" json:"SplDefence"`
	Heal              int32 `protobuf:"varint,9,opt,name=Heal,proto3" json:"Heal"`
	CriticalTimes     int32 `protobuf:"varint,10,opt,name=CriticalTimes,proto3" json:"CriticalTimes"`
	Critical          int32 `protobuf:"varint,11,opt,name=Critical,proto3" json:"Critical"`
	AntiCriticalTimes int32 `protobuf:"varint,12,opt,name=AntiCriticalTimes,proto3" json:"AntiCriticalTimes"`
	AntiCritical      int32 `protobuf:"varint,13,opt,name=AntiCritical,proto3" json:"AntiCritical"`
	Dodge             int32 `protobuf:"varint,14,opt,name=Dodge,proto3" json:"Dodge"`
	Hit               int32 `protobuf:"varint,15,opt,name=Hit,proto3" json:"Hit"`
}

type Z_C_LoginMap struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Id         uint64   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	Pos        *Point3F `protobuf:"bytes,3,opt,name=Pos,proto3" json:"Pos"`
	Rotation   float32  `protobuf:"fixed32,4,opt,name=Rotation,proto3" json:"Rotation"`
}
type CHAT int32

const (
	CHAT_MSG_TYPE_WORLD   CHAT = 0
	CHAT_MSG_TYPE_PRIVATE CHAT = 1
	CHAT_MSG_TYPE_ORG     CHAT = 2
	CHAT_MSG_TYPE_COUNT   CHAT = 3
)

type EmBagType int32

const (
	EmBagType_BagType_Equip    EmBagType = 0
	EmBagType_BagType_Item     EmBagType = 1
	EmBagType_BagType_PetPiece EmBagType = 2
	EmBagType_BagType_Special  EmBagType = 3
	EmBagType_BagType_GodEquip EmBagType = 4
)

type EmBattleType int32

const (
	EmBattleType_EBT_None   EmBattleType = 0
	EmBattleType_EBT_Normal EmBattleType = 1
	EmBattleType_EBT_Tower  EmBattleType = 2
	EmBattleType_EBT_Hook   EmBattleType = 3
)

type EmDoingType int32

const (
	EmDoingType_EDT_Gm                       EmDoingType = 0
	EmDoingType_EDT_Other                    EmDoingType = 1
	EmDoingType_EDT_Client                   EmDoingType = 2
	EmDoingType_EDT_ItemUse                  EmDoingType = 3
	EmDoingType_EDT_GiftCode                 EmDoingType = 4
	EmDoingType_EDT_ProfessionLevel          EmDoingType = 5
	EmDoingType_EDT_ProfessionPeakLevel      EmDoingType = 6
	EmDoingType_EDT_ProfessionPartLevel      EmDoingType = 7
	EmDoingType_EDT_EquipSplit               EmDoingType = 8
	EmDoingType_EDT_HeroAwaken               EmDoingType = 9
	EmDoingType_EDT_HeroRebirth              EmDoingType = 10
	EmDoingType_EDT_BoxScoreExchange         EmDoingType = 11
	EmDoingType_EDT_Battle                   EmDoingType = 12
	EmDoingType_EDT_BoxConsume               EmDoingType = 13
	EmDoingType_EDT_BoxOpen                  EmDoingType = 14
	EmDoingType_EDT_Task                     EmDoingType = 15
	EmDoingType_EDT_System                   EmDoingType = 16
	EmDoingType_EDT_Login                    EmDoingType = 17
	EmDoingType_EDT_ChangePlayerName         EmDoingType = 18
	EmDoingType_EDT_Mail                     EmDoingType = 19
	EmDoingType_EDT_BlackShop                EmDoingType = 20
	EmDoingType_EDT_Advert                   EmDoingType = 21
	EmDoingType_EDT_CrystalRobotUpgrade      EmDoingType = 22
	EmDoingType_EDT_CrystalRedefine          EmDoingType = 23
	EmDoingType_EDT_CrystalBookUpgrade       EmDoingType = 24
	EmDoingType_EDT_Draw                     EmDoingType = 25
	EmDoingType_EDT_Charge                   EmDoingType = 26
	EmDoingType_EDT_BattleHook               EmDoingType = 27
	EmDoingType_EDT_Offline                  EmDoingType = 28
	EmDoingType_EDT_ProfessionPartRefine     EmDoingType = 29
	EmDoingType_EDT_ProfessionPartRefineTupo EmDoingType = 30
	EmDoingType_EDT_HeroBook                 EmDoingType = 31
	EmDoingType_EDT_StarSource               EmDoingType = 32
	EmDoingType_EDT_SevenDay                 EmDoingType = 33
	EmDoingType_EDT_Shop                     EmDoingType = 34
	EmDoingType_EDT_DailyTask                EmDoingType = 35
	EmDoingType_EDT_Reset                    EmDoingType = 36
	EmDoingType_EDT_Entry                    EmDoingType = 37
	EmDoingType_EDT_BattleNormal             EmDoingType = 38
	EmDoingType_EDT_StarSourceDraw           EmDoingType = 39
	EmDoingType_EDT_RankReward               EmDoingType = 40
	EmDoingType_EDT_WorldBoss                EmDoingType = 41
	EmDoingType_EDT_Championship             EmDoingType = 42
	EmDoingType_EDT_FirstCharge              EmDoingType = 43
	EmDoingType_EDT_BP                       EmDoingType = 44
	EmDoingType_EDT_ChargeCard               EmDoingType = 45
	EmDoingType_EDT_ChargeGift               EmDoingType = 46
	EmDoingType_EDT_GrowRoad                 EmDoingType = 47
	EmDoingType_EDT_HookTech                 EmDoingType = 48
	EmDoingType_EDT_CrystalUpgrade           EmDoingType = 49
	EmDoingType_EDT_Adventure                EmDoingType = 50
	EmDoingType_EDT_ItemBuy                  EmDoingType = 51
	EmDoingType_EDT_Activity                 EmDoingType = 52
	EmDoingType_EDT_AdvertEject              EmDoingType = 53
	EmDoingType_EDT_CrystalGenerate          EmDoingType = 54
	EmDoingType_EDT_BattleHookPassReward     EmDoingType = 55
)

type EmGiftCodeType int32

const (
	EmGiftCodeType_GAT_Common EmGiftCodeType = 0
	EmGiftCodeType_GAT_Code   EmGiftCodeType = 1
	EmGiftCodeType_GAT_Week   EmGiftCodeType = 2
	EmGiftCodeType_GAT_Month  EmGiftCodeType = 3
)

type EmGmFuncType int32

const (
	EmGmFuncType_GFT_AddItem  EmGmFuncType = 0
	EmGmFuncType_GFT_AddEquip EmGmFuncType = 1
	EmGmFuncType_GFT_AddHero  EmGmFuncType = 2
	EmGmFuncType_GFT_NB       EmGmFuncType = 3
	EmGmFuncType_GFT_Relogin  EmGmFuncType = 4
	EmGmFuncType_GFT_Charge   EmGmFuncType = 5
)

type EmGmParamType int32

const (
	EmGmParamType_GPT_None   EmGmParamType = 0
	EmGmParamType_GPT_GmOpen EmGmParamType = 1
)

type EmItemExpendType int32

const (
	EmItemExpendType_EIET_None       EmItemExpendType = 0
	EmItemExpendType_EIET_Cash       EmItemExpendType = 1
	EmItemExpendType_EIET_Gold       EmItemExpendType = 2
	EmItemExpendType_EIET_SplitScore EmItemExpendType = 25
	EmItemExpendType_EIET_Max        EmItemExpendType = 100
)

type EmMailState int32

const (
	EmMailState_NoRead      EmMailState = 0
	EmMailState_ReadRecieve EmMailState = 1
)

type EmPlatType int32

const (
	EmPlatType_Local EmPlatType = 0
)

type EmPlayerOfflineType int32

const (
	EmPlayerOfflineType_EPOT_Mail EmPlayerOfflineType = 0
	EmPlayerOfflineType_EPOT_Item EmPlayerOfflineType = 1
)

type EmShopType int32

const (
	EmShopType_EST_None      EmShopType = 0
	EmShopType_EST_BlackShop EmShopType = 1
)

type EmSyetemPropType int32

const (
	EmSyetemPropType_SPT_HeroBook EmSyetemPropType = 0
)

type EmTaskState int32

const (
	EmTaskState_ETS_Ing    EmTaskState = 0
	EmTaskState_ETS_Finish EmTaskState = 1
	EmTaskState_ETS_Award  EmTaskState = 2
)

type LoginState int32

const (
	LoginState_None    LoginState = 0
	LoginState_Init    LoginState = 1
	LoginState_SetName LoginState = 2
)

type MAIL int32

const (
	MAIL_M_PlayerMgr  MAIL = 0
	MAIL_M_AccountMgr MAIL = 1
)

type PlayerDataType int32

const (
	PlayerDataType_Crystal            PlayerDataType = 0
	PlayerDataType_Base               PlayerDataType = 1
	PlayerDataType_System             PlayerDataType = 2
	PlayerDataType_Bag                PlayerDataType = 3
	PlayerDataType_Equipment          PlayerDataType = 4
	PlayerDataType_Client             PlayerDataType = 5
	PlayerDataType_Hero               PlayerDataType = 6
	PlayerDataType_Mail               PlayerDataType = 7
	PlayerDataType_Max                PlayerDataType = 8
	PlayerDataType_SystemCommon       PlayerDataType = 10
	PlayerDataType_SystemChat         PlayerDataType = 11
	PlayerDataType_SystemProfession   PlayerDataType = 12
	PlayerDataType_SystemBox          PlayerDataType = 13
	PlayerDataType_SystemBattle       PlayerDataType = 14
	PlayerDataType_SystemTask         PlayerDataType = 15
	PlayerDataType_SystemShop         PlayerDataType = 16
	PlayerDataType_SystemDraw         PlayerDataType = 17
	PlayerDataType_SystemCharge       PlayerDataType = 18
	PlayerDataType_SystemGene         PlayerDataType = 19
	PlayerDataType_SystemOffline      PlayerDataType = 20
	PlayerDataType_SystemHookTech     PlayerDataType = 21
	PlayerDataType_SystemSevenDay     PlayerDataType = 22
	PlayerDataType_SystemWorldBoss    PlayerDataType = 23
	PlayerDataType_SystemChampionship PlayerDataType = 24
	PlayerDataType_SystemActivity     PlayerDataType = 25
	PlayerDataType_SystemRepair       PlayerDataType = 26
	PlayerDataType_SystemMax          PlayerDataType = 27
)

type Protocol_Player int32

const (
	Protocol_Player_P_C2S_Protocol_Player         Protocol_Player = 0
	Protocol_Player_P_C2S_Protocol_Common         Protocol_Player = 1
	Protocol_Player_P_C2S_Protocol_Copymap        Protocol_Player = 2
	Protocol_Player_P_C2S_Protocol_Pet            Protocol_Player = 3
	Protocol_Player_P_C2S_Protocol_Item           Protocol_Player = 4
	Protocol_Player_P_C2S_Protocol_Fight          Protocol_Player = 5
	Protocol_Player_P_C2S_Protocol_Task           Protocol_Player = 6
	Protocol_Player_P_C2S_Protocol_Mail           Protocol_Player = 7
	Protocol_Player_P_C2S_Protocol_TopList        Protocol_Player = 8
	Protocol_Player_P_C2S_Protocol_Challenge      Protocol_Player = 9
	Protocol_Player_P_C2S_Protocol_Faction        Protocol_Player = 10
	Protocol_Player_P_C2S_Protocol_Team           Protocol_Player = 11
	Protocol_Player_P_C2S_Protocol_Call           Protocol_Player = 12
	Protocol_Player_P_C2S_Protocol_Sail           Protocol_Player = 13
	Protocol_Player_P_C2S_Protocol_Hook           Protocol_Player = 14
	Protocol_Player_P_C2S_Protocol_Artifact       Protocol_Player = 15
	Protocol_Player_P_C2S_Protocol_Shop           Protocol_Player = 16
	Protocol_Player_P_C2S_Protocol_Train          Protocol_Player = 17
	Protocol_Player_P_C2S_Protocol_Achieve        Protocol_Player = 18
	Protocol_Player_P_C2S_Protocol_Expedition     Protocol_Player = 19
	Protocol_Player_P_C2S_Protocol_Shape          Protocol_Player = 20
	Protocol_Player_P_C2S_Protocol_Temple         Protocol_Player = 21
	Protocol_Player_P_C2S_Protocol_Friend         Protocol_Player = 22
	Protocol_Player_P_C2S_Protocol_Element        Protocol_Player = 23
	Protocol_Player_P_C2S_Protocol_Risk           Protocol_Player = 24
	Protocol_Player_P_C2S_Protocol_Dan            Protocol_Player = 25
	Protocol_Player_P_C2S_Protocol_Ladder         Protocol_Player = 26
	Protocol_Player_P_C2S_Protocol_Champion       Protocol_Player = 27
	Protocol_Player_P_C2S_Protocol_Holy           Protocol_Player = 28
	Protocol_Player_P_C2S_Protocol_Video          Protocol_Player = 29
	Protocol_Player_P_C2S_Protocol_Privilege      Protocol_Player = 30
	Protocol_Player_P_C2S_Protocol_Weal           Protocol_Player = 31
	Protocol_Player_P_C2S_Protocol_Activity       Protocol_Player = 32
	Protocol_Player_P_C2S_Protocol_Platform       Protocol_Player = 33
	Protocol_Player_P_C2S_Protocol_Talk           Protocol_Player = 34
	Protocol_Player_P_C2S_Protocol_Treasure       Protocol_Player = 35
	Protocol_Player_P_C2S_Protocol_HeavenDungeon  Protocol_Player = 36
	Protocol_Player_P_C2S_Protocol_CrossChallenge Protocol_Player = 37
	Protocol_Player_P_C2S_Protocol_Tablet         Protocol_Player = 38
	Protocol_Player_P_C2S_Protocol_WeekChampion   Protocol_Player = 39
	Protocol_Player_P_C2S_Protocol_TeamCampaign   Protocol_Player = 40
	Protocol_Player_P_C2S_Protocol_Operate        Protocol_Player = 255
)

type SEND int32

const (
	SEND_POINT      SEND = 0
	SEND_BOARD_CAST SEND = 1
)

type SERVICE int32

const (
	SERVICE_NONE   SERVICE = 0
	SERVICE_CLIENT SERVICE = 1
	SERVICE_GATE   SERVICE = 2
	SERVICE_GM     SERVICE = 3
	SERVICE_GAME   SERVICE = 4
	SERVICE_DB     SERVICE = 5
	SERVICE_Dip    SERVICE = 6
	SERVICE_Record SERVICE = 7
)

type STUB int32

const (
	STUB_Master          STUB = 0
	STUB_DbPlayerMgr     STUB = 1
	STUB_PlayerMgr       STUB = 2
	STUB_ChatChannelMgr  STUB = 4
	STUB_DbChatMgr       STUB = 5
	STUB_AccountMgr      STUB = 6
	STUB_BattleRecordMgr STUB = 7
	STUB_RankMgr         STUB = 8
	STUB_GlobalMgr       STUB = 9
	STUB_END             STUB = 10
)
