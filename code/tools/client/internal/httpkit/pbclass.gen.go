package httpkit

type AchieveTaskInfoNotify struct {
	PacketHead *IPacket           `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SystemType uint32             `protobuf:"varint,2,opt,name=SystemType,proto3" json:"SystemType"` // 成就系统类型
	TaskList   []*PBTaskStageInfo `protobuf:"bytes,3,rep,name=TaskList,proto3" json:"TaskList"`      // 更新的任务进度信息
}
type ActivityDataNewNotify struct {
	ActivityType uint32                  `protobuf:"varint,2,opt,name=ActivityType,proto3" json:"ActivityType"` // 活动类型
	Info         *PBPlayerSystemActivity `protobuf:"bytes,3,opt,name=Info,proto3" json:"Info"`                  // 新增活动
	PacketHead   *IPacket                `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ActivityFreePrizeRequest struct {
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"` // ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ActivityFreePrizeResponse struct {
	Id                 uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`                                 // ID
	NextDailyPrizeTime uint64   `protobuf:"varint,4,opt,name=NextDailyPrizeTime,proto3" json:"NextDailyPrizeTime"` // 下次奖励时间
	PacketHead         *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ActivityListNotify struct {
	ActivityList []*PBPlayerActivityInfo `protobuf:"bytes,2,rep,name=ActivityList,proto3" json:"ActivityList"`   // 活动列表
	DelIdList    []uint32                `protobuf:"varint,3,rep,packed,name=DelIdList,proto3" json:"DelIdList"` // 活动结束的列表
	PacketHead   *IPacket                `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ActivityOpenRequest struct {
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"` // ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ActivityOpenResponse struct {
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`            // ID
	JsonData   []string `protobuf:"bytes,3,rep,name=JsonData,proto3" json:"JsonData"` // json数据
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ActivityRedNotify struct {
	IdList     []uint32 `protobuf:"varint,2,rep,packed,name=IdList,proto3" json:"IdList"` // 活动类型 活动ID列表
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type AdventureRewardRequest struct {
	CfgID      uint32   `protobuf:"varint,3,opt,name=CfgID,proto3" json:"CfgID"` // 奖励配置ID
	ID         uint32   `protobuf:"varint,2,opt,name=ID,proto3" json:"ID"`       // 活动ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type AdventureRewardResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type AdvertNotify struct {
	AdvertList []*PBAdvertInfo `protobuf:"bytes,2,rep,name=AdvertList,proto3" json:"AdvertList"` // 广告信息
	PacketHead *IPacket        `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"` // 包头
}
type AdvertRequest struct {
	AdvestType uint32   `protobuf:"varint,2,opt,name=AdvestType,proto3" json:"AdvestType"` // 广告类型
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`  // 包头
}
type AdvertResponse struct {
	AdvestInfo *PBAdvertInfo `protobuf:"bytes,2,opt,name=AdvestInfo,proto3" json:"AdvestInfo"` // 广告数据
	PacketHead *IPacket      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"` // 包头
}
type AllPlayerInfoNotify struct {
	Mark       uint32        `protobuf:"varint,2,opt,name=Mark,proto3" json:"Mark"`            // 位标记 PlayerDataType 0表示发送完成
	PacketHead *IPacket      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"` // 包头
	PlayerData *PBPlayerData `protobuf:"bytes,3,opt,name=PlayerData,proto3" json:"PlayerData"` // 所有玩家数据
}
type AvatarFrameNotify struct {
	Frames     []*PBAvatarFrame `protobuf:"bytes,2,rep,name=Frames,proto3" json:"Frames"`
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type AvatarNotify struct {
	Avatars    []*PBAvatar `protobuf:"bytes,2,rep,name=Avatars,proto3" json:"Avatars"`
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type AwardMailRequest struct {
	MailId     uint32   `protobuf:"varint,2,opt,name=MailId,proto3" json:"MailId"` // 邮件ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type AwardMailResponse struct {
	Mail       *PBMail  `protobuf:"bytes,2,opt,name=Mail,proto3" json:"Mail"` // 邮件数据
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type BPAcitiveNotify struct {
	BPType     uint32   `protobuf:"varint,2,opt,name=BPType,proto3" json:"BPType"`         // BP类型
	ChargeTime uint64   `protobuf:"varint,4,opt,name=ChargeTime,proto3" json:"ChargeTime"` // 充值时间
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	StageId    uint32   `protobuf:"varint,3,opt,name=StageId,proto3" json:"StageId"` // BP类型
}
type BPNewNotify struct {
	BPInfo     *PBBPInfo `protobuf:"bytes,2,opt,name=BPInfo,proto3" json:"BPInfo"` // 新BP
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type BPNewStageNotify struct {
	BPType     uint32           `protobuf:"varint,2,opt,name=BPType,proto3" json:"BPType"`          // BP类型
	DelList    []uint32         `protobuf:"varint,4,rep,packed,name=DelList,proto3" json:"DelList"` // 删除的期数
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	StageList  []*PBBPStageInfo `protobuf:"bytes,3,rep,name=StageList,proto3" json:"StageList"` // 新BP期数
}
type BPPrizeRequest struct {
	BPType     uint32   `protobuf:"varint,2,opt,name=BPType,proto3" json:"BPType"` // BP类型
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	StageId    uint32   `protobuf:"varint,3,opt,name=StageId,proto3" json:"StageId"` // BP期数
}
type BPPrizeResponse struct {
	BPType        uint32   `protobuf:"varint,2,opt,name=BPType,proto3" json:"BPType"`               // BP类型
	ExtralPrizeId uint32   `protobuf:"varint,5,opt,name=ExtralPrizeId,proto3" json:"ExtralPrizeId"` // 额外领奖ID
	NoramlPrizeId uint32   `protobuf:"varint,4,opt,name=NoramlPrizeId,proto3" json:"NoramlPrizeId"` // 普通领奖ID
	PacketHead    *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	StageId       uint32   `protobuf:"varint,3,opt,name=StageId,proto3" json:"StageId"` // BP期数
}
type BPValueNotify struct {
	BPType     uint32   `protobuf:"varint,2,opt,name=BPType,proto3" json:"BPType"` // BP类型
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Value      uint32   `protobuf:"varint,3,opt,name=Value,proto3" json:"Value"` // BP类型
}
type BattleBeginRequest struct {
	BattleType EmBattleType `protobuf:"varint,2,opt,name=BattleType,proto3,enum=common.EmBattleType" json:"BattleType"` // 战斗类型 1地图，2爬塔，3挂机
	MapId      uint32       `protobuf:"varint,3,opt,name=MapId,proto3" json:"MapId"`                                    // 地图ID
	PacketHead *IPacket     `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Params     []uint32     `protobuf:"varint,5,rep,packed,name=Params,proto3" json:"Params"` // 参数
	StageId    uint32       `protobuf:"varint,4,opt,name=StageId,proto3" json:"StageId"`      // 关卡ID
}
type BattleBeginResponse struct {
	BattleType EmBattleType `protobuf:"varint,2,opt,name=BattleType,proto3,enum=common.EmBattleType" json:"BattleType"` // 战斗类型 1地图，2爬塔，3挂机
	FightCount uint32       `protobuf:"varint,5,opt,name=FightCount,proto3" json:"FightCount"`                          // 挑战次数
	MapId      uint32       `protobuf:"varint,3,opt,name=MapId,proto3" json:"MapId"`                                    // 地图ID
	PacketHead *IPacket     `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Params     []uint32     `protobuf:"varint,6,rep,packed,name=Params,proto3" json:"Params"` // 参数
	StageId    uint32       `protobuf:"varint,4,opt,name=StageId,proto3" json:"StageId"`      // 关卡ID
}
type BattleEndRequest struct {
	Battle     *BattleInfo `protobuf:"bytes,2,opt,name=Battle,proto3" json:"Battle"`
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type BattleEndResponse struct {
	BattleType EmBattleType     `protobuf:"varint,2,opt,name=BattleType,proto3,enum=common.EmBattleType" json:"BattleType"` // 战斗类型 1地图，2爬塔，3挂机
	ItemInfo   []*PBAddItemData `protobuf:"bytes,5,rep,name=ItemInfo,proto3" json:"ItemInfo"`                               // 道具信息
	MapId      uint32           `protobuf:"varint,3,opt,name=MapId,proto3" json:"MapId"`                                    // 地图ID
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	StageId    uint32           `protobuf:"varint,4,opt,name=StageId,proto3" json:"StageId"` // 关卡ID
}
type BattleFunBuyRequest struct {
	BattleFunType uint32   `protobuf:"varint,2,opt,name=BattleFunType,proto3" json:"BattleFunType"` // 刷新类型 1是复活  2是卡牌刷新
	PacketHead    *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type BattleFunBuyResponse struct {
	BattleFunType uint32   `protobuf:"varint,2,opt,name=BattleFunType,proto3" json:"BattleFunType"` // 刷新类型 1是复活  2是卡牌刷新
	PacketHead    *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type BattleInfo struct {
	BattleType  EmBattleType             `protobuf:"varint,1,opt,name=BattleType,proto3,enum=common.EmBattleType" json:"BattleType"` // 战斗类型 1地图，2爬塔，3挂机
	ClientData  *PBBattleClientData      `protobuf:"bytes,8,opt,name=ClientData,proto3" json:"ClientData"`                           // 战场内部客户端数据
	IsSucc      uint32                   `protobuf:"varint,4,opt,name=IsSucc,proto3" json:"IsSucc"`                                  // 结果 1胜利
	MapId       uint32                   `protobuf:"varint,2,opt,name=MapId,proto3" json:"MapId"`                                    // 地图ID
	MonsterInfo []*BattleKillMonsterInfo `protobuf:"bytes,7,rep,name=MonsterInfo,proto3" json:"MonsterInfo"`                         // 击杀数据
	StageId     uint32                   `protobuf:"varint,3,opt,name=StageId,proto3" json:"StageId"`                                // 关卡ID
	StageRate   uint32                   `protobuf:"varint,6,opt,name=StageRate,proto3" json:"StageRate"`                            // 通关进度百分比
	TotalDamage uint64                   `protobuf:"varint,9,opt,name=TotalDamage,proto3" json:"TotalDamage"`                        // 总伤害
	UseTime     uint32                   `protobuf:"varint,5,opt,name=UseTime,proto3" json:"UseTime"`                                // 使用时间秒
}
type BattleKillMonsterInfo struct {
	KillCount   uint32 `protobuf:"varint,2,opt,name=KillCount,proto3" json:"KillCount"`     // 击杀数量
	MaxCount    uint32 `protobuf:"varint,3,opt,name=MaxCount,proto3" json:"MaxCount"`       // 最大数量
	MonsterType uint32 `protobuf:"varint,1,opt,name=MonsterType,proto3" json:"MonsterType"` // 怪物类型 小怪 精英怪 boss
}
type BattleMapNotify struct {
	BattleType EmBattleType     `protobuf:"varint,2,opt,name=BattleType,proto3,enum=common.EmBattleType" json:"BattleType"` // 战斗类型 1地图，2爬塔，3挂机
	MapInfo    *PBBattleMapInfo `protobuf:"bytes,3,opt,name=MapInfo,proto3" json:"MapInfo"`                                 // 战场信息
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type BattleRecordRequest struct {
	BattleType EmBattleType `protobuf:"varint,2,opt,name=BattleType,proto3,enum=common.EmBattleType" json:"BattleType"` // 战斗类型 1地图，2爬塔，3挂机
	MapId      uint32       `protobuf:"varint,3,opt,name=MapId,proto3" json:"MapId"`                                    // 地图ID
	PacketHead *IPacket     `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	StageId    uint32       `protobuf:"varint,4,opt,name=StageId,proto3" json:"StageId"` // 关卡ID
}
type BattleRecordResponse struct {
	BattleType EmBattleType          `protobuf:"varint,2,opt,name=BattleType,proto3,enum=common.EmBattleType" json:"BattleType"` // 战斗类型 1地图，2爬塔，3挂机
	MapId      uint32                `protobuf:"varint,3,opt,name=MapId,proto3" json:"MapId"`                                    // 地图ID
	PacketHead *IPacket              `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RecordList []*PBPlayerBattleData `protobuf:"bytes,5,rep,name=RecordList,proto3" json:"RecordList"` // 记录数据
	StageId    uint32                `protobuf:"varint,4,opt,name=StageId,proto3" json:"StageId"`      // 关卡ID
}
type BattleReliveNotify struct {
	PacketHead *IPacket        `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Relive     *PBBattleRelive `protobuf:"bytes,2,opt,name=Relive,proto3" json:"Relive"`
}
type BattleReliveRequest struct {
	AdvertType uint32       `protobuf:"varint,4,opt,name=AdvertType,proto3" json:"AdvertType"`                          // 广告类型
	BattleType EmBattleType `protobuf:"varint,2,opt,name=BattleType,proto3,enum=common.EmBattleType" json:"BattleType"` // 战斗类型 1地图，2爬塔，3挂机
	HeroId     uint32       `protobuf:"varint,3,opt,name=HeroId,proto3" json:"HeroId"`                                  // 英雄ID
	PacketHead *IPacket     `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type BattleReliveResponse struct {
	BattleType EmBattleType    `protobuf:"varint,2,opt,name=BattleType,proto3,enum=common.EmBattleType" json:"BattleType"` // 战斗类型 1地图，2爬塔，3挂机
	HeroId     uint32          `protobuf:"varint,3,opt,name=HeroId,proto3" json:"HeroId"`                                  // 英雄ID
	PacketHead *IPacket        `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Relive     *PBBattleRelive `protobuf:"bytes,4,opt,name=Relive,proto3" json:"Relive"` // 复活数据
}
type BattleScheduleSaveRequest struct {
	BattleSchedule *PBBattleSchedule `protobuf:"bytes,2,opt,name=BattleSchedule,proto3" json:"BattleSchedule"`
	PacketHead     *IPacket          `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type BattleScheduleSaveResponse struct {
	BattleSchedule *PBBattleSchedule `protobuf:"bytes,2,opt,name=BattleSchedule,proto3" json:"BattleSchedule"`
	PacketHead     *IPacket          `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type BookCollectionCoinRequest struct {
	CrystalID  uint32   `protobuf:"varint,2,opt,name=CrystalID,proto3" json:"CrystalID"` // 晶核ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type BookCollectionCoinResponse struct {
	Coin       uint32   `protobuf:"varint,2,opt,name=Coin,proto3" json:"Coin"`   // 收藏币数量
	Level      uint32   `protobuf:"varint,3,opt,name=Level,proto3" json:"Level"` // 当前图鉴系统等级
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type BookStageRewardRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type BookStageRewardResponse struct {
	BookInfo   *PBCrystalBook `protobuf:"bytes,2,opt,name=BookInfo,proto3" json:"BookInfo"` // 图鉴数据
	PacketHead *IPacket       `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type BoxOpenRequest struct {
	ItemID     uint32   `protobuf:"varint,2,opt,name=ItemID,proto3" json:"ItemID"`
	ItemNum    uint32   `protobuf:"varint,3,opt,name=ItemNum,proto3" json:"ItemNum"`
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type BoxOpenResponse struct {
	ItemInfo   []*PBAddItemData `protobuf:"bytes,3,rep,name=ItemInfo,proto3" json:"ItemInfo"` // 道具信息
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Score      uint32           `protobuf:"varint,2,opt,name=Score,proto3" json:"Score"` // 返回积分
}
type BoxProgressRewardRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type BoxProgressRewardResponse struct {
	Level      uint32   `protobuf:"varint,3,opt,name=Level,proto3" json:"Level"`         // 当前挡位
	NeedScore  uint32   `protobuf:"varint,2,opt,name=NeedScore,proto3" json:"NeedScore"` // 所需积分
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Recycle    uint32   `protobuf:"varint,5,opt,name=Recycle,proto3" json:"Recycle"` // 轮数
	Score      uint32   `protobuf:"varint,4,opt,name=Score,proto3" json:"Score"`     // 当前拥有积分
}
type BroadcastNotify struct {
	BroadcastType uint32   `protobuf:"varint,3,opt,name=BroadcastType,proto3" json:"BroadcastType"` // 广播类型
	Channel       uint32   `protobuf:"varint,2,opt,name=Channel,proto3" json:"Channel"`             // 广播频道
	Content       string   `protobuf:"bytes,4,opt,name=Content,proto3" json:"Content"`              // 广播类容
	Extends       []byte   `protobuf:"bytes,5,opt,name=Extends,proto3" json:"Extends"`              // 扩展字段
	PacketHead    *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`        // 包头
}
type C_Z_LoginCopyMap struct {
	DataId     int32    `protobuf:"varint,2,opt,name=DataId,proto3" json:"DataId"`
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type C_Z_Move struct {
	Move       *C_Z_Move_Move `protobuf:"bytes,2,opt,name=move,proto3" json:"move"`
	PacketHead *IPacket       `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type C_Z_Move_Move struct {
	Mode   int32                 `protobuf:"varint,1,opt,name=Mode,proto3" json:"Mode"`
	Normal *C_Z_Move_Move_Normal `protobuf:"bytes,2,opt,name=normal,proto3" json:"normal"`
}
type C_Z_Move_Move_Normal struct {
	Duration float32  `protobuf:"fixed32,3,opt,name=Duration,proto3" json:"Duration"`
	Pos      *Point3F `protobuf:"bytes,1,opt,name=Pos,proto3" json:"Pos"`
	Yaw      float32  `protobuf:"fixed32,2,opt,name=Yaw,proto3" json:"Yaw"`
}
type C_Z_Skill struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SkillId    int32    `protobuf:"varint,2,opt,name=SkillId,proto3" json:"SkillId"`
	TargetId   int64    `protobuf:"varint,3,opt,name=TargetId,proto3" json:"TargetId"`
}
type ChampionshipInfoRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"` // 包头
}
type ChampionshipInfoResponse struct {
	List       []*ChampionshipRankInfo `protobuf:"bytes,2,rep,name=List,proto3" json:"List"`             // 第一名信息
	PacketHead *IPacket                `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"` // 包头
}
type ChampionshipNotify struct {
	CreateTime uint64                  `protobuf:"varint,2,opt,name=CreateTime,proto3" json:"CreateTime"` // 锦标赛活动开启时间
	Expire     uint64                  `protobuf:"varint,3,opt,name=Expire,proto3" json:"Expire"`         // 活动有效时长
	List       []*ChampionshipTimeInfo `protobuf:"bytes,4,rep,name=List,proto3" json:"List"`              // 第一名信息
	PacketHead *IPacket                `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ChampionshipRankInfo struct {
	First    *RankInfo `protobuf:"bytes,2,opt,name=First,proto3" json:"First"`        // 第一名信息
	RankType uint32    `protobuf:"varint,1,opt,name=RankType,proto3" json:"RankType"` // 排行榜类型
}
type ChampionshipTaskRewardRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RankType   uint32   `protobuf:"varint,2,opt,name=RankType,proto3" json:"RankType"`
	TaskID     uint32   `protobuf:"varint,3,opt,name=TaskID,proto3" json:"TaskID"`
}
type ChampionshipTaskRewardResponse struct {
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Task       *PBTaskStageInfo `protobuf:"bytes,2,opt,name=Task,proto3" json:"Task"` // 主线任务
}
type ChampionshipTimeInfo struct {
	Active   uint64 `protobuf:"varint,3,opt,name=Active,proto3" json:"Active"`     // 活跃时长
	Interval uint64 `protobuf:"varint,2,opt,name=Interval,proto3" json:"Interval"` // 间隔多长时间开启
	RankType uint32 `protobuf:"varint,1,opt,name=RankType,proto3" json:"RankType"` // 排行榜类型
	Reward   uint64 `protobuf:"varint,4,opt,name=Reward,proto3" json:"Reward"`     // 结算时长（领奖时长）
	Show     uint64 `protobuf:"varint,5,opt,name=Show,proto3" json:"Show"`         // 展示时长
}
type ChangeAvatarFrameRequest struct {
	FrameID    uint32   `protobuf:"varint,2,opt,name=FrameID,proto3" json:"FrameID"` // 切换到的头像框ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ChangeAvatarFrameResponse struct {
	FrameID    uint32   `protobuf:"varint,2,opt,name=FrameID,proto3" json:"FrameID"` // 切换到的头像框ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ChangeAvatarRequest struct {
	AvatarID   uint32   `protobuf:"varint,2,opt,name=AvatarID,proto3" json:"AvatarID"` // 切换到的头像ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ChangeAvatarResponse struct {
	AvatarID   uint32   `protobuf:"varint,2,opt,name=AvatarID,proto3" json:"AvatarID"` // 切换到的头像ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ChangePlayerNameRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PlayerName string   `protobuf:"bytes,2,opt,name=PlayerName,proto3" json:"PlayerName"` // 玩家名称
}
type ChangePlayerNameResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PlayerName string   `protobuf:"bytes,2,opt,name=PlayerName,proto3" json:"PlayerName"` // 玩家名称
}
type ChargeCardNewNotify struct {
	CardInfo   *PBChargeCard `protobuf:"bytes,2,opt,name=CardInfo,proto3" json:"CardInfo"` // 充值卡
	PacketHead *IPacket      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ChargeCardPrizeRequest struct {
	CardType   uint32   `protobuf:"varint,2,opt,name=CardType,proto3" json:"CardType"` // 充值卡类型
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ChargeCardPrizeResponse struct {
	CardType   uint32   `protobuf:"varint,2,opt,name=CardType,proto3" json:"CardType"` // 充值卡类型
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PrizeTime  uint64   `protobuf:"varint,3,opt,name=PrizeTime,proto3" json:"PrizeTime"` // 领奖时间
}
type ChargeCardUpdateNotify struct {
	CardInfo   []*PBChargeCard `protobuf:"bytes,2,rep,name=CardInfo,proto3" json:"CardInfo"`       // 充值卡
	DelList    []uint32        `protobuf:"varint,3,rep,packed,name=DelList,proto3" json:"DelList"` // 过期删除的卡
	PacketHead *IPacket        `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ChargeGiftBuyNotify struct {
	BuyInfo    *PBU32U32 `protobuf:"bytes,3,opt,name=BuyInfo,proto3" json:"BuyInfo"` // 礼包ID|数量
	Id         uint32    `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`          // 活动ID
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ChargeNotify struct {
	Charge     *PBCharge `protobuf:"bytes,2,opt,name=Charge,proto3" json:"Charge"` // 充值数据
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ChargeOrderRequest struct {
	IsNeigou   bool     `protobuf:"varint,3,opt,name=IsNeigou,proto3" json:"IsNeigou"` // 是否是内购
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProductId  uint32   `protobuf:"varint,2,opt,name=ProductId,proto3" json:"ProductId"` // 商品ID
}
type ChargeOrderResponse struct {
	BingchuanOrder *PBChargeBingchuanOrder `protobuf:"bytes,2,opt,name=BingchuanOrder,proto3" json:"BingchuanOrder"` // 冰川订单
	PacketHead     *IPacket                `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ChargeQueryNotify struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProductId  uint32   `protobuf:"varint,2,opt,name=ProductId,proto3" json:"ProductId"` // 充值数据
}
type ChargeQueryRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ChargeQueryResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProductIds []uint32 `protobuf:"varint,2,rep,packed,name=ProductIds,proto3" json:"ProductIds"` // 商品ID
}
type ClientJsonNotify struct {
	JsonList   []*PBJsonInfo `protobuf:"bytes,2,rep,name=JsonList,proto3" json:"JsonList"`     // json内容
	PacketHead *IPacket      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"` // 包头
}
type ClusterInfo struct {
	ClusterID  uint32  `protobuf:"varint,7,opt,name=ClusterID,proto3" json:"ClusterID"`
	CreateTime uint64  `protobuf:"varint,8,opt,name=CreateTime,proto3" json:"CreateTime"` // 加入时间
	Ip         string  `protobuf:"bytes,2,opt,name=Ip,proto3" json:"Ip"`
	Port       int32   `protobuf:"varint,3,opt,name=Port,proto3" json:"Port"`
	SocketId   uint32  `protobuf:"varint,5,opt,name=SocketId,proto3" json:"SocketId"`
	Type       SERVICE `protobuf:"varint,1,opt,name=Type,proto3,enum=common.SERVICE" json:"Type"`
	Version    uint32  `protobuf:"varint,6,opt,name=Version,proto3" json:"Version"` // 版本号
	Weight     int32   `protobuf:"varint,4,opt,name=Weight,proto3" json:"Weight"`
}
type CommonNotify struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"` // 包头
}
type CommonPrizeNotify struct {
	DoingType  EmDoingType      `protobuf:"varint,2,opt,name=DoingType,proto3,enum=common.EmDoingType" json:"DoingType"` // 原因
	ItemInfo   []*PBAddItemData `protobuf:"bytes,3,rep,name=ItemInfo,proto3" json:"ItemInfo"`                            // 道具信息
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`                        // 包头
}
type CrystalNotify struct {
	CrystalInfo []*PBCrystal `protobuf:"bytes,2,rep,name=CrystalInfo,proto3" json:"CrystalInfo"` // 变更之后的晶核信息
	PacketHead  *IPacket     `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type CrystalRedefineRequest struct {
	CrystalID  uint32   `protobuf:"varint,2,opt,name=CrystalID,proto3" json:"CrystalID"` // 要改造的晶核ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type CrystalRedefineResponse struct {
	CrystalID  uint32   `protobuf:"varint,2,opt,name=CrystalID,proto3" json:"CrystalID"` // 要改造的晶核ID
	CurStar    uint32   `protobuf:"varint,3,opt,name=CurStar,proto3" json:"CurStar"`     // 升级之后的星级
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type CrystalRobotBattleRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RobotID    uint32   `protobuf:"varint,2,opt,name=RobotID,proto3" json:"RobotID"` // 机器人ID
}
type CrystalRobotBattleResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RobotID    uint32   `protobuf:"varint,2,opt,name=RobotID,proto3" json:"RobotID"` // 机器人ID
}
type CrystalRobotNotify struct {
	PacketHead *IPacket        `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RobotInfo  *PBCrystalRobot `protobuf:"bytes,2,opt,name=RobotInfo,proto3" json:"RobotInfo"` // 变更之后的机器人信息
}
type CrystalRobotUpgradeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RobotID    uint32   `protobuf:"varint,2,opt,name=RobotID,proto3" json:"RobotID"` // 机器人ID
}
type CrystalRobotUpgradeResponse struct {
	CurLevel   uint32   `protobuf:"varint,3,opt,name=CurLevel,proto3" json:"CurLevel"` // 升级之后的当前等级
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RobotID    uint32   `protobuf:"varint,2,opt,name=RobotID,proto3" json:"RobotID"` // 机器人ID
}
type CrystalUpgradeRequest struct {
	CrystalID  uint32   `protobuf:"varint,2,opt,name=CrystalID,proto3" json:"CrystalID"` // 要改造的晶核ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type CrystalUpgradeResponse struct {
	CrystalID  uint32   `protobuf:"varint,2,opt,name=CrystalID,proto3" json:"CrystalID"` // 要改造的晶核ID
	CurLevel   uint32   `protobuf:"varint,3,opt,name=CurLevel,proto3" json:"CurLevel"`   // 升级之后的等级
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type DailyTasFinishResponse struct {
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Score      uint32           `protobuf:"varint,3,opt,name=Score,proto3" json:"Score"` // 活跃值
	Task       *PBTaskStageInfo `protobuf:"bytes,2,opt,name=Task,proto3" json:"Task"`    // 主线任务
}
type DailyTaskFinishRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	TaskId     uint32   `protobuf:"varint,2,opt,name=TaskId,proto3" json:"TaskId"` // 任务ID
}
type DailyTaskNotify struct {
	DailyTask  *PBDailyTask `protobuf:"bytes,2,opt,name=DailyTask,proto3" json:"DailyTask"` // 每日任务
	PacketHead *IPacket     `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type DailyTaskScorePrizeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type DailyTaskScorePrizeResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PrizeScore uint32   `protobuf:"varint,2,opt,name=PrizeScore,proto3" json:"PrizeScore"` // 领取的活跃值
}
type DeleteMailRequest struct {
	MailId     uint32   `protobuf:"varint,2,opt,name=MailId,proto3" json:"MailId"` // 邮件ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type DeleteMailResponse struct {
	MailId     uint32   `protobuf:"varint,2,opt,name=MailId,proto3" json:"MailId"` // 邮件ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type DrawNotify struct {
	DelDrawList []uint32      `protobuf:"varint,3,rep,packed,name=DelDrawList,proto3" json:"DelDrawList"` // 删除的抽奖
	DrawList    []*PBDrawInfo `protobuf:"bytes,2,rep,name=DrawList,proto3" json:"DrawList"`               // 新增的抽奖
	PacketHead  *IPacket      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type DrawPrizeInfoRequest struct {
	DrawId     uint32   `protobuf:"varint,2,opt,name=DrawId,proto3" json:"DrawId"` // 抽奖id
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type DrawPrizeInfoResponse struct {
	DrawId     uint32             `protobuf:"varint,2,opt,name=DrawId,proto3" json:"DrawId"` // 抽奖id
	PacketHead *IPacket           `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PrizeList  []*PBDrawPrizeInfo `protobuf:"bytes,3,rep,name=PrizeList,proto3" json:"PrizeList"` // 奖励信息
}
type DrawRequest struct {
	AdvertType   uint32   `protobuf:"varint,5,opt,name=AdvertType,proto3" json:"AdvertType"`     // 广告类型
	DrawCount    uint32   `protobuf:"varint,3,opt,name=DrawCount,proto3" json:"DrawCount"`       // 抽奖次数
	DrawId       uint32   `protobuf:"varint,2,opt,name=DrawId,proto3" json:"DrawId"`             // 抽奖id
	IsUseReplace bool     `protobuf:"varint,4,opt,name=IsUseReplace,proto3" json:"IsUseReplace"` // 是否用替换道具
	PacketHead   *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type DrawResponse struct {
	DrawInfo   *PBDrawInfo `protobuf:"bytes,2,opt,name=DrawInfo,proto3" json:"DrawInfo"` // 抽奖信息
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type DrawScorePrizeRequest struct {
	DrawId     uint32   `protobuf:"varint,2,opt,name=DrawId,proto3" json:"DrawId"` // 抽奖id
	Id         uint32   `protobuf:"varint,3,opt,name=Id,proto3" json:"Id"`         // 抽奖积分配置配置的ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type DrawScorePrizeResponse struct {
	DrawId     uint32   `protobuf:"varint,2,opt,name=DrawId,proto3" json:"DrawId"` // 抽奖id
	Id         uint32   `protobuf:"varint,3,opt,name=Id,proto3" json:"Id"`         // 抽奖积分配置配置的ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type EntryCondition struct {
	CfgID      uint32 `protobuf:"varint,1,opt,name=CfgID,proto3" json:"CfgID"`           // 配置id
	Process    uint32 `protobuf:"varint,2,opt,name=Process,proto3" json:"Process"`       // 条件完成进度
	Times      uint32 `protobuf:"varint,3,opt,name=Times,proto3" json:"Times"`           // 触发次数
	UpdateTime uint64 `protobuf:"varint,4,opt,name=UpdateTime,proto3" json:"UpdateTime"` // 刷新时间
}
type EntryConditionNotify struct {
	Condition  *EntryCondition `protobuf:"bytes,2,opt,name=Condition,proto3" json:"Condition"`
	PacketHead *IPacket        `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type EntryEffect struct {
	List       []*EntryEffectData `protobuf:"bytes,3,rep,name=List,proto3" json:"List"`              // 效果参数
	ParamsType uint32             `protobuf:"varint,1,opt,name=ParamsType,proto3" json:"ParamsType"` // 参数类型
	Type       uint32             `protobuf:"varint,2,opt,name=Type,proto3" json:"Type"`             // 效果类型
}
type EntryEffectData struct {
	Object uint32              `protobuf:"varint,1,opt,name=Object,proto3" json:"Object"` // 生效对象
	Values []*EntryEffectValue `protobuf:"bytes,2,rep,name=Values,proto3" json:"Values"`  // 效果参数
}
type EntryEffectNotify struct {
	Effect     *EntryEffect `protobuf:"bytes,2,opt,name=Effect,proto3" json:"Effect"`
	PacketHead *IPacket     `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type EntryEffectValue struct {
	List []uint32 `protobuf:"varint,1,rep,packed,name=List,proto3" json:"List"` // 参数
}
type EntryTriggerRequest struct {
	EntryType  uint32   `protobuf:"varint,2,opt,name=EntryType,proto3" json:"EntryType"` // 触发的词条类型
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Params     []uint32 `protobuf:"varint,4,rep,packed,name=Params,proto3" json:"Params"` // 扩展参数（例如品质、星级、等）
	Times      uint32   `protobuf:"varint,3,opt,name=Times,proto3" json:"Times"`          // 触发次数
}
type EntryTriggerResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type EntryUnlockRequest struct {
	PacketHead     *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PassiveSkillID uint32   `protobuf:"varint,2,opt,name=PassiveSkillID,proto3" json:"PassiveSkillID"` // 被动技能id
}
type EntryUnlockResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type EquipmentAutoSplitRequest struct {
	PacketHead  *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	QualityList []uint32 `protobuf:"varint,2,rep,packed,name=QualityList,proto3" json:"QualityList"` // 品质列表
}
type EquipmentAutoSplitResponse struct {
	PacketHead  *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	QualityList []uint32 `protobuf:"varint,2,rep,packed,name=QualityList,proto3" json:"QualityList"` // 品质列表
}
type EquipmentBuyPosRequest struct {
	CurPosBuyCount uint32   `protobuf:"varint,2,opt,name=CurPosBuyCount,proto3" json:"CurPosBuyCount"` // 当前购买次数
	PacketHead     *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type EquipmentBuyPosResponse struct {
	PacketHead  *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PosBuyCount uint32   `protobuf:"varint,2,opt,name=PosBuyCount,proto3" json:"PosBuyCount"` // 当前购买次数
}
type EquipmentLockRequest struct {
	IsLock     bool     `protobuf:"varint,3,opt,name=IsLock,proto3" json:"IsLock"` // 是否加锁
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Sn         uint32   `protobuf:"varint,2,opt,name=Sn,proto3" json:"Sn"` // 装备Sn
}
type EquipmentLockResponse struct {
	IsLock     bool     `protobuf:"varint,3,opt,name=IsLock,proto3" json:"IsLock"` // 是否加锁
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Sn         uint32   `protobuf:"varint,2,opt,name=Sn,proto3" json:"Sn"` // 装备Sn
}
type EquipmentNotify struct {
	Equipment  []*PBEquipment `protobuf:"bytes,2,rep,name=Equipment,proto3" json:"Equipment"` // 装备数据
	IsHook     bool           `protobuf:"varint,3,opt,name=IsHook,proto3" json:"IsHook"`      // 是否挂机背包
	PacketHead *IPacket       `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type EquipmentSplitRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SnList     []uint32 `protobuf:"varint,2,rep,packed,name=SnList,proto3" json:"SnList"` // 装备sn列表
}
type EquipmentSplitResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SnList     []uint32 `protobuf:"varint,2,rep,packed,name=SnList,proto3" json:"SnList"` // 装备sn列表
}
type EquipmentSplitScoreNotify struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SplitScore uint32   `protobuf:"varint,2,opt,name=SplitScore,proto3" json:"SplitScore"` // 当前进度积分
}
type FirstChargeNotify struct {
	FirstChargeList []*PBFirstCharge `protobuf:"bytes,2,rep,name=FirstChargeList,proto3" json:"FirstChargeList"` // 新增的首冲数据
	PacketHead      *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type FirstChargePrizeRequest struct {
	FirstChargeId uint32   `protobuf:"varint,2,opt,name=FirstChargeId,proto3" json:"FirstChargeId"` // 首冲类型
	PacketHead    *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type FirstChargePrizeResponse struct {
	FirstChargeId uint32   `protobuf:"varint,2,opt,name=FirstChargeId,proto3" json:"FirstChargeId"` // 首冲类型
	PacketHead    *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PrizeDay      uint32   `protobuf:"varint,3,opt,name=PrizeDay,proto3" json:"PrizeDay"` // 领取的最新奖励天数
}
type GeneCardActiveInfo struct {
	CardID   uint32 `protobuf:"varint,2,opt,name=CardID,proto3" json:"CardID"`     // 卡牌ID
	IsActive bool   `protobuf:"varint,1,opt,name=IsActive,proto3" json:"IsActive"` // 是否激活
}
type GeneCardActiveRequest struct {
	Actives        []*GeneCardActiveInfo `protobuf:"bytes,3,rep,name=Actives,proto3" json:"Actives"` // 激活列表
	PacketHead     *IPacket              `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RobotPositions []uint32              `protobuf:"varint,4,rep,packed,name=RobotPositions,proto3" json:"RobotPositions"` // 重置机器人位置列表
	SchemeID       uint32                `protobuf:"varint,2,opt,name=SchemeID,proto3" json:"SchemeID"`                    // 方案ID
}
type GeneCardActiveResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type GeneChangeNameRequest struct {
	Name       string   `protobuf:"bytes,3,opt,name=Name,proto3" json:"Name"` // 改成的名字
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SchemeID   uint32   `protobuf:"varint,2,opt,name=SchemeID,proto3" json:"SchemeID"` // 基因方案ID
}
type GeneChangeNameResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type GeneRobotActiveRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Position   uint32   `protobuf:"varint,3,opt,name=Position,proto3" json:"Position"` // 位置
	RobotID    uint32   `protobuf:"varint,4,opt,name=RobotID,proto3" json:"RobotID"`   // 机器人id
	SchemeID   uint32   `protobuf:"varint,2,opt,name=SchemeID,proto3" json:"SchemeID"` // 方案ID
}
type GeneRobotActiveResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type GeneRobotCardActiveRequest struct {
	Actives    []*GeneCardActiveInfo `protobuf:"bytes,4,rep,name=Actives,proto3" json:"Actives"` // 激活列表
	PacketHead *IPacket              `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Position   uint32                `protobuf:"varint,3,opt,name=Position,proto3" json:"Position"` // 位置
	SchemeID   uint32                `protobuf:"varint,2,opt,name=SchemeID,proto3" json:"SchemeID"` // 方案ID
}
type GeneRobotCardActiveResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type GeneSchemeChangeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SchemeID   uint32   `protobuf:"varint,2,opt,name=SchemeID,proto3" json:"SchemeID"` // 基因方案ID
}
type GeneSchemeChangeResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type GeneSchemeResetRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SchemeID   uint32   `protobuf:"varint,2,opt,name=SchemeID,proto3" json:"SchemeID"` // 基因方案ID
}
type GeneSchemeResetResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type GetPlayerDataRequest struct {
	DataType   int32    `protobuf:"varint,2,opt,name=DataType,proto3" json:"DataType"`
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type GetPlayerDataResponse struct {
	PacketHead *IPacket      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PlayerData *PBPlayerData `protobuf:"bytes,2,opt,name=PlayerData,proto3" json:"PlayerData"`
}
type GiftCodeRequest struct {
	Acode      string   `protobuf:"bytes,2,opt,name=Acode,proto3" json:"Acode"` // 兑换码
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type GiftCodeResponse struct {
	Acode      string   `protobuf:"bytes,2,opt,name=Acode,proto3" json:"Acode"` // 兑换码
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type GmFuncRequest struct {
	GmType     uint32   `protobuf:"varint,2,opt,name=GmType,proto3" json:"GmType"` // gm类型
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Param      []string `protobuf:"bytes,3,rep,name=Param,proto3" json:"Param"` // 参数数据
}
type GmFuncResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type GrowRoadTaskPrizeRequest struct {
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"` // 活动ID growroad表的ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	TaskId     uint32   `protobuf:"varint,3,opt,name=TaskId,proto3" json:"TaskId"` // 任务ID
}
type GrowRoadTaskPrizeResponse struct {
	Id            uint32           `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"` // 活动ID growroad表的ID
	PacketHead    *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	TaskStageInfo *PBTaskStageInfo `protobuf:"bytes,4,opt,name=TaskStageInfo,proto3" json:"TaskStageInfo"` // 任务数据
}
type HeardPacket struct {
}
type HeartbeatRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"` // 包头
	Time       uint64   `protobuf:"varint,2,opt,name=Time,proto3" json:"Time"`            // 当前时间戳
}
type HeartbeatResponse struct {
	CurTime    uint64   `protobuf:"varint,4,opt,name=CurTime,proto3" json:"CurTime"`      // 当前时间戳
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"` // 包头
	RecvTime   uint64   `protobuf:"varint,3,opt,name=RecvTime,proto3" json:"RecvTime"`    // 收到时间戳
	SendTime   uint64   `protobuf:"varint,2,opt,name=SendTime,proto3" json:"SendTime"`    // 发送时间戳
}
type HeroAutoUpStarRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type HeroAutoUpStarResponse struct {
	DelSnList  []uint32    `protobuf:"varint,3,rep,packed,name=DelSnList,proto3" json:"DelSnList"` // 消耗的英雄
	HeroList   []*PBU32U32 `protobuf:"bytes,2,rep,name=HeroList,proto3" json:"HeroList"`           // 升星的英雄列表 key:sn value:星级
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type HeroAwakenLevelRequest struct {
	CurLevel   uint32   `protobuf:"varint,3,opt,name=CurLevel,proto3" json:"CurLevel"` // 当前等级
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Sn         uint32   `protobuf:"varint,2,opt,name=Sn,proto3" json:"Sn"` // 觉醒的sn
}
type HeroAwakenLevelResponse struct {
	CurLevel   uint32   `protobuf:"varint,3,opt,name=CurLevel,proto3" json:"CurLevel"` // 当前等级
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Sn         uint32   `protobuf:"varint,2,opt,name=Sn,proto3" json:"Sn"` // 觉醒的sn
}
type HeroBattleStarChangeNotify struct {
	Heros      []*HeroBattleStarInfo `protobuf:"bytes,2,rep,name=Heros,proto3" json:"Heros"`
	PacketHead *IPacket              `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type HeroBattleStarChangeRequest struct {
	Heros      []*HeroBattleStarInfo `protobuf:"bytes,2,rep,name=Heros,proto3" json:"Heros"`
	PacketHead *IPacket              `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type HeroBattleStarChangeResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type HeroBattleStarInfo struct {
	HeroID uint32 `protobuf:"varint,1,opt,name=HeroID,proto3" json:"HeroID"` // 英雄id
	Total  uint32 `protobuf:"varint,2,opt,name=Total,proto3" json:"Total"`   // 当前总量
}
type HeroBookActiveRequest struct {
	HeroId     uint32   `protobuf:"varint,2,opt,name=HeroId,proto3" json:"HeroId"` // 英雄ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type HeroBookActiveResponse struct {
	HeroBook   *PBHeroBook `protobuf:"bytes,2,opt,name=HeroBook,proto3" json:"HeroBook"` // 英雄图鉴
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type HeroBookNotify struct {
	HeroBook   *PBHeroBook `protobuf:"bytes,2,opt,name=HeroBook,proto3" json:"HeroBook"` // 英雄图鉴 0星表示需要激活
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type HeroBookUpStarRequest struct {
	HeroId     uint32   `protobuf:"varint,2,opt,name=HeroId,proto3" json:"HeroId"` // 英雄ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type HeroBookUpStarResponse struct {
	HeroBook   *PBHeroBook `protobuf:"bytes,2,opt,name=HeroBook,proto3" json:"HeroBook"` // 英雄图鉴
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type HeroFightPowerNotify struct {
	FightPower uint32   `protobuf:"varint,2,opt,name=FightPower,proto3" json:"FightPower"` // 战斗力
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type HeroGameHeroListNotify struct {
	PacketHead *IPacket        `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Team       *PBHeroTeamList `protobuf:"bytes,2,opt,name=Team,proto3" json:"Team"` // 上阵列表
}
type HeroGameHeroListRequest struct {
	PacketHead *IPacket        `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Team       *PBHeroTeamList `protobuf:"bytes,2,opt,name=Team,proto3" json:"Team"` // 上阵列表
}
type HeroGameHeroListResponse struct {
	PacketHead *IPacket        `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Team       *PBHeroTeamList `protobuf:"bytes,2,opt,name=Team,proto3" json:"Team"` // 上阵列表
}
type HeroNewStarNotify struct {
	Info       []*PBHero `protobuf:"bytes,2,rep,name=Info,proto3" json:"Info"` // 英雄数据
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type HeroNotify struct {
	Info       []*PBHero `protobuf:"bytes,2,rep,name=Info,proto3" json:"Info"` // 英雄数据
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type HeroRebirthRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Sn         uint32   `protobuf:"varint,2,opt,name=Sn,proto3" json:"Sn"` // 英雄sn
}
type HeroRebirthResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Sn         uint32   `protobuf:"varint,2,opt,name=Sn,proto3" json:"Sn"`     // 升星的sn
	Star       uint32   `protobuf:"varint,3,opt,name=Star,proto3" json:"Star"` // 星级
}
type HeroUpStarRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Sn         uint32   `protobuf:"varint,2,opt,name=Sn,proto3" json:"Sn"`                      // 升星的sn
	UseSnList  []uint32 `protobuf:"varint,3,rep,packed,name=UseSnList,proto3" json:"UseSnList"` // 消耗的英雄
}
type HeroUpStarResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Sn         uint32   `protobuf:"varint,2,opt,name=Sn,proto3" json:"Sn"`                      // 升星的sn
	Star       uint32   `protobuf:"varint,3,opt,name=Star,proto3" json:"Star"`                  // 星级
	UseSnList  []uint32 `protobuf:"varint,4,rep,packed,name=UseSnList,proto3" json:"UseSnList"` // 消耗的英雄
}
type HookBattleAutoMapRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type HookBattleAutoMapResponse struct {
	AutoMap    bool     `protobuf:"varint,2,opt,name=AutoMap,proto3" json:"AutoMap"` // 是否自动推关
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type HookBattleLootRequest struct {
	MonsterInfo []*BattleKillMonsterInfo `protobuf:"bytes,2,rep,name=MonsterInfo,proto3" json:"MonsterInfo"` // 击杀数据
	PacketHead  *IPacket                 `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type HookBattleLootResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type HookEquipmentAwardRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Sn         uint32   `protobuf:"varint,2,opt,name=Sn,proto3" json:"Sn"` // 装备Sn
}
type HookEquipmentAwardResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Sn         uint32   `protobuf:"varint,2,opt,name=Sn,proto3" json:"Sn"` // 装备Sn
}
type HookTechLevelNotify struct {
	HookTechList []*PBHookTech `protobuf:"bytes,2,rep,name=HookTechList,proto3" json:"HookTechList"` // 挂机科技列表
	PacketHead   *IPacket      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type HookTechLevelRequest struct {
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"` // Id
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type HookTechLevelResponse struct {
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`               // Id
	LevelTime  uint64   `protobuf:"varint,3,opt,name=LevelTime,proto3" json:"LevelTime"` // 升级到达时间
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type HookTechSpeedRequest struct {
	AdvertType uint32   `protobuf:"varint,3,opt,name=AdvertType,proto3" json:"AdvertType"` // 广告类型 0无广告
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`                 // Id
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type HookTechSpeedResponse struct {
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`               // Id
	Level      uint32   `protobuf:"varint,3,opt,name=Level,proto3" json:"Level"`         // 等级
	LevelTime  uint64   `protobuf:"varint,4,opt,name=LevelTime,proto3" json:"LevelTime"` // 升级到达时间
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type IPacket struct {
	Ckx            uint32  `protobuf:"varint,2,opt,name=ckx,proto3" json:"ckx"`
	Code           uint32  `protobuf:"varint,6,opt,name=code,proto3" json:"code"`                                         // 错误码 ErrorCode
	DestServerType SERVICE `protobuf:"varint,3,opt,name=destServerType,proto3,enum=common.SERVICE" json:"destServerType"` // 发送者服务器类型
	Id             uint64  `protobuf:"varint,4,opt,name=id,proto3" json:"id"`                                             // 目标ID
	Seqid          uint32  `protobuf:"varint,5,opt,name=seqid,proto3" json:"seqid"`                                       // 序列号
	Stx            uint32  `protobuf:"varint,1,opt,name=stx,proto3" json:"stx"`
}
type ItemBuyRequest struct {
	Count      uint32   `protobuf:"varint,3,opt,name=Count,proto3" json:"Count"` // 个数
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`       // 道具Id
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ItemBuyResponse struct {
	Count      uint32   `protobuf:"varint,3,opt,name=Count,proto3" json:"Count"` // 使用个数
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`       // 道具Id
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ItemSelectRequest struct {
	Id         uint32      `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"` // 道具Id
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SelectList []*PBU32U32 `protobuf:"bytes,3,rep,name=SelectList,proto3" json:"SelectList"` // key：ItemUseShowResponse的id  value:数量
}
type ItemSelectResponse struct {
	Id         uint32      `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"` // 道具Id
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SelectList []*PBU32U32 `protobuf:"bytes,3,rep,name=SelectList,proto3" json:"SelectList"` // key：ItemUseShowResponse的id  value:数量
}
type ItemUpdateNotify struct {
	ItemList   []*PBItem `protobuf:"bytes,2,rep,name=ItemList,proto3" json:"ItemList"` // 更新的道具
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ItemUseRequest struct {
	Count      uint32   `protobuf:"varint,3,opt,name=Count,proto3" json:"Count"` // 使用个数
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`       // 道具Id
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ItemUseResponse struct {
	Count      uint32   `protobuf:"varint,3,opt,name=Count,proto3" json:"Count"` // 使用个数
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`       // 道具Id
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ItemUseShowInfo struct {
	Id   uint32     `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"` // selectId
	Item *PBAddItem `protobuf:"bytes,2,opt,name=Item,proto3" json:"Item"`
}
type ItemUseShowRequest struct {
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"` // 道具Id
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type ItemUseShowResponse struct {
	Id         uint32             `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`            // 道具ID
	ItemList   []*ItemUseShowInfo `protobuf:"bytes,3,rep,name=ItemList,proto3" json:"ItemList"` // 道具展示数据
	PacketHead *IPacket           `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type LoginRequest struct {
	AccountName string   `protobuf:"bytes,2,opt,name=AccountName,proto3" json:"AccountName"` // 账号名
	PacketHead  *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`   // 包头
	TokenKey    string   `protobuf:"bytes,3,opt,name=TokenKey,proto3" json:"TokenKey"`       // 秘钥 账号名称+accountid的+token md5
}
type LoginResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"` // 包头
	Time       uint64   `protobuf:"varint,2,opt,name=Time,proto3" json:"Time"`            // 当前时间戳
}
type MailBox struct {
	ClusterId uint32 `protobuf:"varint,5,opt,name=ClusterId,proto3" json:"ClusterId"` // 集群id
	Id        uint64 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`
	LeaseId   uint64 `protobuf:"varint,2,opt,name=LeaseId,proto3" json:"LeaseId"`
	MailType  MAIL   `protobuf:"varint,3,opt,name=MailType,proto3,enum=rpc3.MAIL" json:"MailType"`
}
type MainTaskFinishRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type MainTaskFinishResponse struct {
	MainTask   *PBTaskStageInfo `protobuf:"bytes,2,opt,name=MainTask,proto3" json:"MainTask"` // 主线任务
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type NewMailNotify struct {
	Mail       *PBMail  `protobuf:"bytes,2,opt,name=Mail,proto3" json:"Mail"` // 邮件数据
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type NormalBattlePrizeRequest struct {
	MapId      uint32   `protobuf:"varint,2,opt,name=MapId,proto3" json:"MapId"` // 地图ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PrizeStage uint32   `protobuf:"varint,4,opt,name=PrizeStage,proto3" json:"PrizeStage"` // 阶段ID
	StageId    uint32   `protobuf:"varint,3,opt,name=StageId,proto3" json:"StageId"`       // 关卡ID
}
type NormalBattlePrizeResponse struct {
	PacketHead   *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PrizeMapId   uint32   `protobuf:"varint,2,opt,name=PrizeMapId,proto3" json:"PrizeMapId"`        // 领取的地图ID
	PrizeStage   []uint32 `protobuf:"varint,4,rep,packed,name=PrizeStage,proto3" json:"PrizeStage"` // 领取的奖励进度
	PrizeStageId uint32   `protobuf:"varint,3,opt,name=PrizeStageId,proto3" json:"PrizeStageId"`    // 领取的关卡ID
}
type NoticeNotify struct {
	IsNew      bool      `protobuf:"varint,2,opt,name=IsNew,proto3" json:"IsNew"`  // 是否新的
	Notice     *PBNotice `protobuf:"bytes,3,opt,name=Notice,proto3" json:"Notice"` // 公告内容
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type NoticeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type NoticeResponse struct {
	NoticeList []*PBNotice `protobuf:"bytes,2,rep,name=NoticeList,proto3" json:"NoticeList"` // 公告列表
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type OfflineIncomeRewardRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type OfflineIncomeRewardResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type OnekeyAwardMailRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type OnekeyAwardMailResponse struct {
	Mails      []*PBMail `protobuf:"bytes,2,rep,name=Mails,proto3" json:"Mails"` // 邮件数据
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type OnekeyDeleteMailRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type OnekeyDeleteMailResponse struct {
	MailIds    []uint32 `protobuf:"varint,2,rep,packed,name=MailIds,proto3" json:"MailIds"` // 邮件ID组
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type OpenBossRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type OpenBossResponse struct {
	PacketHead    *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	WorldBossRank *PBU32U64 `protobuf:"bytes,2,opt,name=WorldBossRank,proto3" json:"WorldBossRank"` // 世界boss排名数据 key:名次 value:分数
}
type OpenServerGiftBuyNotify struct {
	BuyInfo    *PBU32U32 `protobuf:"bytes,4,opt,name=BuyInfo,proto3" json:"BuyInfo"` // id|数量
	GiftId     uint32    `protobuf:"varint,3,opt,name=GiftId,proto3" json:"GiftId"`  // 礼包ID
	Id         uint32    `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`          // sID
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type OpenServerGiftNewNotify struct {
	GiftInfo   *PBOpenServerGiftInfo `protobuf:"bytes,3,opt,name=GiftInfo,proto3" json:"GiftInfo"` // 礼包ID|数量
	PacketHead *IPacket              `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	SId        uint32                `protobuf:"varint,2,opt,name=SId,proto3" json:"SId"` // 活动ID
}
type PBAchieveInfo struct {
	AchieveType uint32   `protobuf:"varint,1,opt,name=AchieveType,proto3" json:"AchieveType"` // 成就类型
	Params      []uint32 `protobuf:"varint,2,rep,packed,name=Params,proto3" json:"Params"`    // 参数
	Value       uint32   `protobuf:"varint,3,opt,name=Value,proto3" json:"Value"`             // 当前值
}
type PBActivityAdventure struct {
	BeginTime    uint64   `protobuf:"varint,2,opt,name=BeginTime,proto3" json:"BeginTime"`            // 开始时间
	EndTime      uint64   `protobuf:"varint,3,opt,name=EndTime,proto3" json:"EndTime"`                // 结束时间
	Id           uint32   `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`                          // 活动ID
	PrizeIdList  []uint32 `protobuf:"varint,5,rep,packed,name=PrizeIdList,proto3" json:"PrizeIdList"` // 领取的奖励配置ID
	RegisterTime uint64   `protobuf:"varint,4,opt,name=RegisterTime,proto3" json:"RegisterTime"`      // 注册时间
}
type PBActivityChargeGift struct {
	BeginTime uint64      `protobuf:"varint,2,opt,name=BeginTime,proto3" json:"BeginTime"` // 开始时间
	BuyList   []*PBU32U32 `protobuf:"bytes,4,rep,name=BuyList,proto3" json:"BuyList"`      // 已列表 礼包ID|数量
	EndTime   uint64      `protobuf:"varint,3,opt,name=EndTime,proto3" json:"EndTime"`     // 结束时间
	Id        uint32      `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`               // 活动ID
}
type PBActivityGrowRoadInfo struct {
	BeginTime uint64             `protobuf:"varint,2,opt,name=BeginTime,proto3" json:"BeginTime"` // 开始时间
	EndTime   uint64             `protobuf:"varint,3,opt,name=EndTime,proto3" json:"EndTime"`     // 结束时间
	Id        uint32             `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`               // 活动ID
	TaskList  []*PBTaskStageInfo `protobuf:"bytes,4,rep,name=TaskList,proto3" json:"TaskList"`    // 任务列表
}
type PBActivityOpenServerGift struct {
	BeginTime          uint64                  `protobuf:"varint,2,opt,name=BeginTime,proto3" json:"BeginTime"`                   // 开始时间
	EndTime            uint64                  `protobuf:"varint,3,opt,name=EndTime,proto3" json:"EndTime"`                       // 结束时间
	GiftList           []*PBOpenServerGiftInfo `protobuf:"bytes,4,rep,name=GiftList,proto3" json:"GiftList"`                      // 礼包列表
	Id                 uint32                  `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`                                 // 活动SID
	NextDailyPrizeTime uint64                  `protobuf:"varint,5,opt,name=NextDailyPrizeTime,proto3" json:"NextDailyPrizeTime"` // 下次领取时间
}
type PBAddItem struct {
	Count  int64    `protobuf:"varint,3,opt,name=Count,proto3" json:"Count"`          // 数量
	Id     uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`                // 道具ID
	Kind   uint32   `protobuf:"varint,1,opt,name=Kind,proto3" json:"Kind"`            // 道具类型 EItemKindType
	Params []uint32 `protobuf:"varint,4,rep,packed,name=Params,proto3" json:"Params"` // 参数 英雄:星级  装备：品质，星级，sn
}
type PBAddItemData struct {
	Count     int64        `protobuf:"varint,3,opt,name=Count,proto3" json:"Count"`                                 // 物品数量
	DoingType EmDoingType  `protobuf:"varint,4,opt,name=DoingType,proto3,enum=common.EmDoingType" json:"DoingType"` // 操作类型
	Equipment *PBEquipment `protobuf:"bytes,6,opt,name=Equipment,proto3" json:"Equipment"`                          // 装备数据
	Id        uint32       `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`                                       // 物品ID
	Kind      uint32       `protobuf:"varint,1,opt,name=Kind,proto3" json:"Kind"`                                   // 类型
	Params    []uint32     `protobuf:"varint,5,rep,packed,name=Params,proto3" json:"Params"`                        // 物品参数
}
type PBAdvertInfo struct {
	DailyCount uint32 `protobuf:"varint,2,opt,name=DailyCount,proto3" json:"DailyCount"` // 每日观看次数
	Type       uint32 `protobuf:"varint,1,opt,name=Type,proto3" json:"Type"`             // 广告类型
}
type PBAllChatMsgInfo struct {
	Msg []*PBChatMsgInfo `protobuf:"bytes,1,rep,name=Msg,proto3" json:"Msg"` // 单条聊天信息
}
type PBAvatar struct {
	AvatarID uint32 `protobuf:"varint,1,opt,name=AvatarID,proto3" json:"AvatarID"` // 头像ID
	Type     uint32 `protobuf:"varint,2,opt,name=Type,proto3" json:"Type"`         // 获取方式(EmDointType)
}
type PBAvatarFrame struct {
	FrameID uint32 `protobuf:"varint,1,opt,name=FrameID,proto3" json:"FrameID"` // 头像框ID
	Type    uint32 `protobuf:"varint,2,opt,name=Type,proto3" json:"Type"`       // 获取方式(EmDointType)
}
type PBBPInfo struct {
	BPType    uint32           `protobuf:"varint,1,opt,name=BPType,proto3" json:"BPType"`      // BP类型
	MaxStage  uint32           `protobuf:"varint,4,opt,name=MaxStage,proto3" json:"MaxStage"`  // 历史最大期
	StageList []*PBBPStageInfo `protobuf:"bytes,3,rep,name=StageList,proto3" json:"StageList"` // 期数数据
	Value     uint32           `protobuf:"varint,2,opt,name=Value,proto3" json:"Value"`        // 当前值
}
type PBBPStageInfo struct {
	BeginTime     uint64 `protobuf:"varint,2,opt,name=BeginTime,proto3" json:"BeginTime"`         // 结束时间
	ChargeTime    uint64 `protobuf:"varint,6,opt,name=ChargeTime,proto3" json:"ChargeTime"`       // 充值时间
	EndTime       uint64 `protobuf:"varint,3,opt,name=EndTime,proto3" json:"EndTime"`             // 结束时间
	ExtralPrizeId uint32 `protobuf:"varint,5,opt,name=ExtralPrizeId,proto3" json:"ExtralPrizeId"` // 额外领奖ID
	NoramlPrizeId uint32 `protobuf:"varint,4,opt,name=NoramlPrizeId,proto3" json:"NoramlPrizeId"` // 普通领奖ID
	StageId       uint32 `protobuf:"varint,1,opt,name=StageId,proto3" json:"StageId"`             // 期数
}
type PBBattleClientData struct {
	CryscalRobotId  []uint32    `protobuf:"varint,4,rep,packed,name=CryscalRobotId,proto3" json:"CryscalRobotId"` // 使徒ID
	DropBoxCount    uint32      `protobuf:"varint,5,opt,name=DropBoxCount,proto3" json:"DropBoxCount"`            // 空投次数
	HeroBattleLevel []*PBU32U32 `protobuf:"bytes,1,rep,name=HeroBattleLevel,proto3" json:"HeroBattleLevel"`       // 战场内部英雄等级数据
	LeaderId        uint32      `protobuf:"varint,3,opt,name=LeaderId,proto3" json:"LeaderId"`                    // 队长ID
	SelectCard      []uint32    `protobuf:"varint,2,rep,packed,name=SelectCard,proto3" json:"SelectCard"`         // 战场内部选卡数据
}
type PBBattleHookInfo struct {
	AutoMap        bool             `protobuf:"varint,4,opt,name=AutoMap,proto3" json:"AutoMap"`               // 是否自动推关
	BeginLootTime  uint64           `protobuf:"varint,5,opt,name=BeginLootTime,proto3" json:"BeginLootTime"`   // 开始掉落时间
	CurMapId       uint32           `protobuf:"varint,2,opt,name=CurMapId,proto3" json:"CurMapId"`             // 当前地图ID
	CurStageId     uint32           `protobuf:"varint,3,opt,name=CurStageId,proto3" json:"CurStageId"`         // 当前关卡ID
	MapInfo        *PBBattleMapInfo `protobuf:"bytes,1,opt,name=MapInfo,proto3" json:"MapInfo"`                // 战场信息
	TotalLootCount uint32           `protobuf:"varint,6,opt,name=TotalLootCount,proto3" json:"TotalLootCount"` // 累计掉落件数
}
type PBBattleMapInfo struct {
	FightCount   uint32 `protobuf:"varint,6,opt,name=FightCount,proto3" json:"FightCount"`     // 挑战次数 begin就算一次
	IsSuceess    uint32 `protobuf:"varint,7,opt,name=IsSuceess,proto3" json:"IsSuceess"`       // 是否通关 1通关
	MapId        uint32 `protobuf:"varint,1,opt,name=MapId,proto3" json:"MapId"`               // 地图ID
	RebirthCount uint32 `protobuf:"varint,8,opt,name=RebirthCount,proto3" json:"RebirthCount"` // 总复活次数
	StageId      uint32 `protobuf:"varint,2,opt,name=StageId,proto3" json:"StageId"`           // 关卡ID
	StageRate    uint32 `protobuf:"varint,4,opt,name=StageRate,proto3" json:"StageRate"`       // 通关进度万分比
	Time         uint64 `protobuf:"varint,3,opt,name=Time,proto3" json:"Time"`                 // 时间戳
	UseTime      uint32 `protobuf:"varint,5,opt,name=UseTime,proto3" json:"UseTime"`           // 使用时间秒
}
type PBBattleNormalInfo struct {
	MapInfo      *PBBattleMapInfo `protobuf:"bytes,1,opt,name=MapInfo,proto3" json:"MapInfo"`               // 战场信息
	PrizeMapId   uint32           `protobuf:"varint,2,opt,name=PrizeMapId,proto3" json:"PrizeMapId"`        // 领取的地图ID
	PrizeStage   []uint32         `protobuf:"varint,4,rep,packed,name=PrizeStage,proto3" json:"PrizeStage"` // 领取的奖励进度
	PrizeStageId uint32           `protobuf:"varint,3,opt,name=PrizeStageId,proto3" json:"PrizeStageId"`    // 领取的关卡ID
}
type PBBattleRelive struct {
	AdvestReliveCount uint32 `protobuf:"varint,1,opt,name=AdvestReliveCount,proto3" json:"AdvestReliveCount"` // 广告复活次数
	ShareReliveCount  uint32 `protobuf:"varint,2,opt,name=ShareReliveCount,proto3" json:"ShareReliveCount"`   // 分享复活次数
}
type PBBattleSchedule struct {
	BattleType   uint32                   `protobuf:"varint,1,opt,name=BattleType,proto3" json:"BattleType"`     // 战斗类型
	ClientData   *PBBattleClientData      `protobuf:"bytes,5,opt,name=ClientData,proto3" json:"ClientData"`      // 战场客户端数据
	MonsterInfo  []*BattleKillMonsterInfo `protobuf:"bytes,6,rep,name=MonsterInfo,proto3" json:"MonsterInfo"`    // 击杀数据
	RebirthCount uint32                   `protobuf:"varint,4,opt,name=RebirthCount,proto3" json:"RebirthCount"` // 复活次数
	StageRate    uint32                   `protobuf:"varint,2,opt,name=StageRate,proto3" json:"StageRate"`       // 通关进度万分比
	UseTime      uint32                   `protobuf:"varint,3,opt,name=UseTime,proto3" json:"UseTime"`           // 使用时间秒
}
type PBBlackShop struct {
	Items           []*PBShopGoodInfo  `protobuf:"bytes,2,rep,name=Items,proto3" json:"Items"`                      // 橱窗列表(刷新生成)
	NextRefreshTime uint64             `protobuf:"varint,1,opt,name=NextRefreshTime,proto3" json:"NextRefreshTime"` // 下一次刷新时间点
	RefreshInfo     *PBShopRefreshInfo `protobuf:"bytes,3,opt,name=RefreshInfo,proto3" json:"RefreshInfo"`          // 刷新数据
}
type PBBoxInfo struct {
	ItemID    uint32 `protobuf:"varint,1,opt,name=ItemID,proto3" json:"ItemID"`
	OpenTimes uint32 `protobuf:"varint,2,opt,name=OpenTimes,proto3" json:"OpenTimes"` // 开宝箱次数
}
type PBCharge struct {
	DailyCharge uint32 `protobuf:"varint,3,opt,name=DailyCharge,proto3" json:"DailyCharge"` // 每日充值数据
	MonthCharge uint32 `protobuf:"varint,5,opt,name=MonthCharge,proto3" json:"MonthCharge"` // 每月充值数据
	OrderId     uint32 `protobuf:"varint,1,opt,name=OrderId,proto3" json:"OrderId"`         // 流水ID
	TotalCharge uint32 `protobuf:"varint,2,opt,name=TotalCharge,proto3" json:"TotalCharge"` // 累计充值
	WeekCharge  uint32 `protobuf:"varint,4,opt,name=WeekCharge,proto3" json:"WeekCharge"`   // 每周充值数据
}
type PBChargeBingchuanOrder struct {
	ActorId          string `protobuf:"bytes,6,opt,name=actorId,proto3" json:"actorId"`                   // 角色ID
	CurrencyType     string `protobuf:"bytes,7,opt,name=currencyType,proto3" json:"currencyType"`         // 货币类型
	DeveloperPayload string `protobuf:"bytes,8,opt,name=developerPayload,proto3" json:"developerPayload"` // 透传参数
	OrderItem        string `protobuf:"bytes,1,opt,name=OrderItem,proto3" json:"OrderItem"`               // 订单明细
	OrderNo          string `protobuf:"bytes,2,opt,name=OrderNo,proto3" json:"OrderNo"`                   // 游戏订单号
	OrderSign        string `protobuf:"bytes,5,opt,name=orderSign,proto3" json:"orderSign"`               // 签名
	PayNum           string `protobuf:"bytes,3,opt,name=payNum,proto3" json:"payNum"`                     // 订单明细
	UserId           string `protobuf:"bytes,4,opt,name=userId,proto3" json:"userId"`                     // 冰川用户ID
}
type PBChargeCard struct {
	BeginTime uint64 `protobuf:"varint,2,opt,name=BeginTime,proto3" json:"BeginTime"` // 开始时间
	CardType  uint32 `protobuf:"varint,1,opt,name=CardType,proto3" json:"CardType"`   // 月卡类型
	EndTime   uint64 `protobuf:"varint,3,opt,name=EndTime,proto3" json:"EndTime"`     // 结束时间
	PrizeTime uint64 `protobuf:"varint,4,opt,name=PrizeTime,proto3" json:"PrizeTime"` // 上次领奖时间
}
type PBChatMsgInfo struct {
	Display *PBPlayerBaseDisplay `protobuf:"bytes,2,opt,name=Display,proto3" json:"Display"` // 发送者数据
	Index   uint64               `protobuf:"varint,1,opt,name=Index,proto3" json:"Index"`    // 消息序号
	Msg     string               `protobuf:"bytes,3,opt,name=Msg,proto3" json:"Msg"`         // 发送的消息
	Time    uint64               `protobuf:"varint,4,opt,name=Time,proto3" json:"Time"`      // 发送的时间
}
type PBClientData struct {
	Data string `protobuf:"bytes,2,opt,name=Data,proto3" json:"Data"` // 数据
	Type string `protobuf:"bytes,1,opt,name=Type,proto3" json:"Type"` // 类型
}
type PBCrystal struct {
	CrystalID       uint32   `protobuf:"varint,1,opt,name=CrystalID,proto3" json:"CrystalID"`                    // 晶核ID
	Element         uint32   `protobuf:"varint,2,opt,name=Element,proto3" json:"Element"`                        // 元素属性
	Level           uint32   `protobuf:"varint,7,opt,name=Level,proto3" json:"Level"`                            // 晶核等级
	PassiveSkillIds []uint32 `protobuf:"varint,6,rep,packed,name=PassiveSkillIds,proto3" json:"PassiveSkillIds"` // 解锁被动技能ID
	Quality         uint32   `protobuf:"varint,3,opt,name=Quality,proto3" json:"Quality"`                        // 品质属性
	RewardCoinTimes uint32   `protobuf:"varint,5,opt,name=RewardCoinTimes,proto3" json:"RewardCoinTimes"`        // 升星对应的领奖次数(每一次升星，对应一次领取奖励，例如收藏币)
	Star            uint32   `protobuf:"varint,4,opt,name=Star,proto3" json:"Star"`                              // 星星数量
}
type PBCrystalBook struct {
	Coin          uint32 `protobuf:"varint,1,opt,name=Coin,proto3" json:"Coin"`                   // 总收藏币数量
	FinishedStage uint32 `protobuf:"varint,3,opt,name=FinishedStage,proto3" json:"FinishedStage"` // 已经领取的等级
	Stage         uint32 `protobuf:"varint,2,opt,name=Stage,proto3" json:"Stage"`                 // 当前图鉴系统等级
}
type PBCrystalProp struct {
	Key   uint32 `protobuf:"varint,1,opt,name=Key,proto3" json:"Key"`     // 主键
	Value uint32 `protobuf:"varint,2,opt,name=Value,proto3" json:"Value"` // 值
}
type PBCrystalRobot struct {
	Crystals         []uint32 `protobuf:"varint,7,rep,packed,name=Crystals,proto3" json:"Crystals"`             // 装备的晶核
	FinishStage      uint32   `protobuf:"varint,3,opt,name=FinishStage,proto3" json:"FinishStage"`              // 完成等级
	RobotID          uint32   `protobuf:"varint,1,opt,name=RobotID,proto3" json:"RobotID"`                      // 机器人ID
	RoleSkillID      uint32   `protobuf:"varint,4,opt,name=RoleSkillID,proto3" json:"RoleSkillID"`              // 机器人的共鸣技能ID
	RoleSkillPercent uint32   `protobuf:"varint,5,opt,name=RoleSkillPercent,proto3" json:"RoleSkillPercent"`    // 技能参数提升百分比
	Stage            uint32   `protobuf:"varint,2,opt,name=Stage,proto3" json:"Stage"`                          // 当前等级
	UnlockLinkages   []uint32 `protobuf:"varint,6,rep,packed,name=UnlockLinkages,proto3" json:"UnlockLinkages"` // 已经解锁的共鸣技能
}
type PBDailyTask struct {
	PrizeScore uint32             `protobuf:"varint,3,opt,name=PrizeScore,proto3" json:"PrizeScore"` // 领取的活跃值
	Score      uint32             `protobuf:"varint,2,opt,name=Score,proto3" json:"Score"`           // 活跃值
	TaskList   []*PBTaskStageInfo `protobuf:"bytes,1,rep,name=TaskList,proto3" json:"TaskList"`      // 任务列表
}
type PBDrawInfo struct {
	AdvertNextTime uint64   `protobuf:"varint,5,opt,name=AdvertNextTime,proto3" json:"AdvertNextTime"` // 下次广告购买时间
	BeginTime      uint64   `protobuf:"varint,6,opt,name=BeginTime,proto3" json:"BeginTime"`           // 开启时间戳
	DrawCount      uint32   `protobuf:"varint,2,opt,name=DrawCount,proto3" json:"DrawCount"`           // 抽奖次数
	DrawId         uint32   `protobuf:"varint,1,opt,name=DrawId,proto3" json:"DrawId"`                 // 抽奖次数
	EndTime        uint64   `protobuf:"varint,7,opt,name=EndTime,proto3" json:"EndTime"`               // 结束时间戳
	Guar2Count     uint32   `protobuf:"varint,4,opt,name=Guar2Count,proto3" json:"Guar2Count"`         // 保底2次数
	Guar3Count     uint32   `protobuf:"varint,8,opt,name=Guar3Count,proto3" json:"Guar3Count"`         // 保底3次数
	GuarCount      uint32   `protobuf:"varint,3,opt,name=GuarCount,proto3" json:"GuarCount"`           // 保底次数
	ScorePrize     []uint32 `protobuf:"varint,9,rep,packed,name=ScorePrize,proto3" json:"ScorePrize"`  // 积分奖励进度
}
type PBDrawPrizeInfo struct {
	ItemList []*PBAddItem `protobuf:"bytes,3,rep,name=ItemList,proto3" json:"ItemList"` // 道具列表
	Name     string       `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name"`         // 名称
	Rate     uint32       `protobuf:"varint,2,opt,name=Rate,proto3" json:"Rate"`        // 概率万分比
}
type PBEquipment struct {
	EquipProfession uint32             `protobuf:"varint,9,opt,name=EquipProfession,proto3" json:"EquipProfession"` // 穿戴职业
	Id              uint32             `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`                           // 装备ID
	IsLock          bool               `protobuf:"varint,10,opt,name=IsLock,proto3" json:"IsLock"`                  // 是否锁定
	LinkPropList    []*PBEquipmentProp `protobuf:"bytes,8,rep,name=LinkPropList,proto3" json:"LinkPropList"`        // 共鸣词条
	MainProp        *PBEquipmentProp   `protobuf:"bytes,5,opt,name=MainProp,proto3" json:"MainProp"`                // 主词条
	MinorPropList   []*PBEquipmentProp `protobuf:"bytes,6,rep,name=MinorPropList,proto3" json:"MinorPropList"`      // 次词条
	Quality         uint32             `protobuf:"varint,3,opt,name=Quality,proto3" json:"Quality"`                 // 品质
	Sn              uint32             `protobuf:"varint,1,opt,name=Sn,proto3" json:"Sn"`                           // 装备唯一
	Star            uint32             `protobuf:"varint,4,opt,name=Star,proto3" json:"Star"`                       // 星级
	VicePropList    []*PBEquipmentProp `protobuf:"bytes,7,rep,name=VicePropList,proto3" json:"VicePropList"`        // 副词条
}
type PBEquipmentProp struct {
	PropId uint32 `protobuf:"varint,1,opt,name=PropId,proto3" json:"PropId"` // 属性ID
	Score  uint32 `protobuf:"varint,3,opt,name=Score,proto3" json:"Score"`   // 评分 对应装备评分表的区间
	Value  uint32 `protobuf:"varint,2,opt,name=Value,proto3" json:"Value"`   // 属性值
}
type PBFirstCharge struct {
	ActiveTime    uint64 `protobuf:"varint,2,opt,name=ActiveTime,proto3" json:"ActiveTime"`       // 激活时间
	FirstChargeId uint32 `protobuf:"varint,1,opt,name=FirstChargeId,proto3" json:"FirstChargeId"` // 首冲类型
	PrizeDay      uint32 `protobuf:"varint,3,opt,name=PrizeDay,proto3" json:"PrizeDay"`           // 领取的最新奖励天数
}
type PBGeneRobot struct {
	Position uint32       `protobuf:"varint,2,opt,name=Position,proto3" json:"Position"` // 位置
	RobotID  uint32       `protobuf:"varint,1,opt,name=RobotID,proto3" json:"RobotID"`   // 机器人id
	Tags     []*PBGeneTag `protobuf:"bytes,3,rep,name=Tags,proto3" json:"Tags"`          // 机器人激活的强化卡
}
type PBGeneScheme struct {
	Name     string         `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name"`          // 方案名字(为空，表示使用默认值)
	Robots   []*PBGeneRobot `protobuf:"bytes,4,rep,name=Robots,proto3" json:"Robots"`      // 共鸣使徒
	SchemeID uint32         `protobuf:"varint,1,opt,name=SchemeID,proto3" json:"SchemeID"` // 方案id
	Tags     []*PBGeneTag   `protobuf:"bytes,3,rep,name=Tags,proto3" json:"Tags"`          // 基因触发器
}
type PBGeneTag struct {
	Cards []uint32 `protobuf:"varint,2,rep,packed,name=Cards,proto3" json:"Cards"` // 激活的卡牌
	TagID uint32   `protobuf:"varint,1,opt,name=TagID,proto3" json:"TagID"`        // 标签ID
}
type PBHero struct {
	AwakenLevel uint32 `protobuf:"varint,4,opt,name=AwakenLevel,proto3" json:"AwakenLevel"` // 英雄技能觉醒等级
	BattleStar  uint32 `protobuf:"varint,5,opt,name=BattleStar,proto3" json:"BattleStar"`   // 上阵星星数量
	Id          uint32 `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`                   // 英雄id
	Sn          uint32 `protobuf:"varint,1,opt,name=Sn,proto3" json:"Sn"`                   // 英雄唯一
	Star        uint32 `protobuf:"varint,3,opt,name=Star,proto3" json:"Star"`               // 英雄星级
}
type PBHeroBook struct {
	Id      uint32 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`           // 英雄id
	MaxStar uint32 `protobuf:"varint,3,opt,name=MaxStar,proto3" json:"MaxStar"` // 最大星级
	Star    uint32 `protobuf:"varint,2,opt,name=Star,proto3" json:"Star"`       // 星级 0表示需要激活
}
type PBHeroTeamList struct {
	HeroSn   []uint32 `protobuf:"varint,2,rep,packed,name=HeroSn,proto3" json:"HeroSn"` // 英雄ID
	TeamType uint32   `protobuf:"varint,1,opt,name=TeamType,proto3" json:"TeamType"`    // 编队类型
}
type PBHookTech struct {
	Id        uint32 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`               // Id
	Level     uint32 `protobuf:"varint,2,opt,name=Level,proto3" json:"Level"`         // 等级
	LevelTime uint64 `protobuf:"varint,3,opt,name=LevelTime,proto3" json:"LevelTime"` // 结束升级时间
}
type PBItem struct {
	Count int64  `protobuf:"varint,2,opt,name=Count,proto3" json:"Count"` // 物品数量
	Id    uint32 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`       // 物品ID
}
type PBJsonInfo struct {
	JsonData string `protobuf:"bytes,2,opt,name=JsonData,proto3" json:"JsonData"` // json数据
	JsonName string `protobuf:"bytes,1,opt,name=JsonName,proto3" json:"JsonName"` // json名称
}
type PBMail struct {
	AwardTime  uint64           `protobuf:"varint,6,opt,name=AwardTime,proto3" json:"AwardTime"`                 // 过期时间
	Content    string           `protobuf:"bytes,8,opt,name=Content,proto3" json:"Content"`                      // 正文
	ExpireTime uint64           `protobuf:"varint,5,opt,name=ExpireTime,proto3" json:"ExpireTime"`               // 过期时间
	Id         uint32           `protobuf:"varint,3,opt,name=Id,proto3" json:"Id"`                               // 编号
	Item       []*PBAddItemData `protobuf:"bytes,10,rep,name=Item,proto3" json:"Item"`                           // 道具
	Receiver   uint64           `protobuf:"varint,2,opt,name=Receiver,proto3" json:"Receiver"`                   // 接收者玩家ID
	SendTime   uint64           `protobuf:"varint,4,opt,name=SendTime,proto3" json:"SendTime"`                   // 发送时间
	Sender     string           `protobuf:"bytes,1,opt,name=Sender,proto3" json:"Sender"`                        // 发送者玩家ID
	State      EmMailState      `protobuf:"varint,9,opt,name=State,proto3,enum=common.EmMailState" json:"State"` // 已读标记 (0代表未读,其他已读) EmMailType
	Title      string           `protobuf:"bytes,7,opt,name=Title,proto3" json:"Title"`                          // 标题
}
type PBNotice struct {
	BeginTime uint64 `protobuf:"varint,4,opt,name=BeginTime,proto3" json:"BeginTime"` // 开始时间戳
	Content   string `protobuf:"bytes,3,opt,name=Content,proto3" json:"Content"`      // 正文
	EndTime   uint64 `protobuf:"varint,5,opt,name=EndTime,proto3" json:"EndTime"`     // 结束时间戳
	Id        uint32 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`               // 公告ID 越大越优先
	Sender    string `protobuf:"bytes,6,opt,name=Sender,proto3" json:"Sender"`        // 发件人
	Title     string `protobuf:"bytes,2,opt,name=Title,proto3" json:"Title"`          // 标题
}
type PBOfflineData struct {
	DoingType   EmDoingType         `protobuf:"varint,4,opt,name=DoingType,proto3,enum=common.EmDoingType" json:"DoingType"`             // 发放原因
	Item        []*PBAddItemData    `protobuf:"bytes,3,rep,name=Item,proto3" json:"Item"`                                                // 道具
	Mail        *PBMail             `protobuf:"bytes,2,opt,name=Mail,proto3" json:"Mail"`                                                // 邮件
	Notify      bool                `protobuf:"varint,5,opt,name=Notify,proto3" json:"Notify"`                                           // 是否发送恭喜获得
	OfflineType EmPlayerOfflineType `protobuf:"varint,1,opt,name=OfflineType,proto3,enum=common.EmPlayerOfflineType" json:"OfflineType"` // 离线数据类型
}
type PBOpenServerGiftInfo struct {
	BeginTime uint64      `protobuf:"varint,2,opt,name=BeginTime,proto3" json:"BeginTime"` // 开始时间
	EndTime   uint64      `protobuf:"varint,3,opt,name=EndTime,proto3" json:"EndTime"`     // 结束时间
	GiftId    uint32      `protobuf:"varint,1,opt,name=GiftId,proto3" json:"GiftId"`       // 礼包ID
	StageList []*PBU32U32 `protobuf:"bytes,4,rep,name=StageList,proto3" json:"StageList"`  // 档次列表 ID|数量
}
type PBPlayerActivityInfo struct {
	ActivityId uint32 `protobuf:"varint,1,opt,name=ActivityId,proto3" json:"ActivityId"` // 活动ID
	BeginTime  uint64 `protobuf:"varint,2,opt,name=BeginTime,proto3" json:"BeginTime"`   // 开始时间
	EndTime    uint64 `protobuf:"varint,3,opt,name=EndTime,proto3" json:"EndTime"`       // 结束时间
}
type PBPlayerBag struct {
	DailyBuyItem []*PBU32U32 `protobuf:"bytes,2,rep,name=DailyBuyItem,proto3" json:"DailyBuyItem"` // 日购买道具
	ItemList     []*PBItem   `protobuf:"bytes,1,rep,name=ItemList,proto3" json:"ItemList"`         // 道具背包
}
type PBPlayerBase struct {
	AccountName       string               `protobuf:"bytes,10,opt,name=AccountName,proto3" json:"AccountName"`                      // 账号名称（openid）
	CreateTime        uint64               `protobuf:"varint,2,opt,name=CreateTime,proto3" json:"CreateTime"`                        // 创建时间
	Display           *PBPlayerBaseDisplay `protobuf:"bytes,1,opt,name=Display,proto3" json:"Display"`                               // 展示数据
	LastDailyTime     uint64               `protobuf:"varint,5,opt,name=LastDailyTime,proto3" json:"LastDailyTime"`                  // 上一次跨天时间戳
	LastModifyTime    uint64               `protobuf:"varint,6,opt,name=LastModifyTime,proto3" json:"LastModifyTime"`                // 上一次修改名字时间戳(扩展字段)
	LoginState        LoginState           `protobuf:"varint,3,opt,name=LoginState,proto3,enum=common.LoginState" json:"LoginState"` // 登录状态
	NewPlayerTypeList []uint32             `protobuf:"varint,7,rep,packed,name=NewPlayerTypeList,proto3" json:"NewPlayerTypeList"`   // 是否new过的系统
	PlatSystemType    uint32               `protobuf:"varint,9,opt,name=PlatSystemType,proto3" json:"PlatSystemType"`                // 平台系统类型 安卓/ios/海外安卓/海外ios
	PlatType          uint32               `protobuf:"varint,8,opt,name=PlatType,proto3" json:"PlatType"`                            // 平台类型 本地0,冰川1
	SeverStartTime    uint64               `protobuf:"varint,11,opt,name=SeverStartTime,proto3" json:"SeverStartTime"`               // 开服时间
}
type PBPlayerBaseDisplay struct {
	AccountId     uint64 `protobuf:"varint,1,opt,name=AccountId,proto3" json:"AccountId"`         // 账号ID
	AvatarFrameID uint32 `protobuf:"varint,6,opt,name=AvatarFrameID,proto3" json:"AvatarFrameID"` // 头像框ID
	AvatarID      uint32 `protobuf:"varint,5,opt,name=AvatarID,proto3" json:"AvatarID"`           // 头像ID
	PlayerLevel   uint32 `protobuf:"varint,3,opt,name=PlayerLevel,proto3" json:"PlayerLevel"`     // 角色等级
	PlayerName    string `protobuf:"bytes,2,opt,name=PlayerName,proto3" json:"PlayerName"`        // 角色名称
	SeverId       uint32 `protobuf:"varint,7,opt,name=SeverId,proto3" json:"SeverId"`             // 内部区服ID
	VipLevel      uint32 `protobuf:"varint,4,opt,name=vipLevel,proto3" json:"vipLevel"`           // vip等级
}
type PBPlayerBattleData struct {
	ClientData *PBBattleClientData  `protobuf:"bytes,6,opt,name=ClientData,proto3" json:"ClientData"`  // 战场内部数据
	Display    *PBPlayerBaseDisplay `protobuf:"bytes,1,opt,name=Display,proto3" json:"Display"`        // 基本数据
	FightPower uint32               `protobuf:"varint,5,opt,name=FightPower,proto3" json:"FightPower"` // 战斗力
	HeroList   []*PBHero            `protobuf:"bytes,2,rep,name=HeroList,proto3" json:"HeroList"`      // 英雄列表
	Time       uint64               `protobuf:"varint,3,opt,name=Time,proto3" json:"Time"`             // 通关时间戳
	UseTime    uint32               `protobuf:"varint,4,opt,name=UseTime,proto3" json:"UseTime"`       // 使用时间
}
type PBPlayerClientData struct {
	ClientDataList []*PBClientData `protobuf:"bytes,1,rep,name=ClientDataList,proto3" json:"ClientDataList"` // 前端数据
}
type PBPlayerCrystal struct {
	Book       *PBCrystalBook    `protobuf:"bytes,1,opt,name=Book,proto3" json:"Book"`             // 图鉴系统
	Conditions []*EntryCondition `protobuf:"bytes,4,rep,name=Conditions,proto3" json:"Conditions"` // 词条条件
	Crystals   []*PBCrystal      `protobuf:"bytes,3,rep,name=Crystals,proto3" json:"Crystals"`     // 获取的晶核
	Effects    []*EntryEffect    `protobuf:"bytes,5,rep,name=Effects,proto3" json:"Effects"`       // 效果数据
	Robots     []*PBCrystalRobot `protobuf:"bytes,2,rep,name=Robots,proto3" json:"Robots"`         // 已经解锁的机器人
}
type PBPlayerData struct {
	Bag       *PBPlayerBag        `protobuf:"bytes,3,opt,name=Bag,proto3" json:"Bag"`             // 玩家背包
	Base      *PBPlayerBase       `protobuf:"bytes,1,opt,name=Base,proto3" json:"Base"`           // 角色基础数据
	Client    *PBPlayerClientData `protobuf:"bytes,5,opt,name=Client,proto3" json:"Client"`       // 角色前端数据
	Crystal   *PBPlayerCrystal    `protobuf:"bytes,8,opt,name=Crystal,proto3" json:"Crystal"`     // 晶核系统
	Equipment *PBPlayerEquipment  `protobuf:"bytes,4,opt,name=Equipment,proto3" json:"Equipment"` // 玩家装备
	Hero      *PBPlayerHero       `protobuf:"bytes,6,opt,name=Hero,proto3" json:"Hero"`           // 伙伴数据
	Mail      *PBPlayerMail       `protobuf:"bytes,7,opt,name=Mail,proto3" json:"Mail"`           // 邮件数据
	System    *PBPlayerSystem     `protobuf:"bytes,2,opt,name=System,proto3" json:"System"`       // 角色功能数据
}
type PBPlayerEquipment struct {
	AutoSplitQuality  []uint32       `protobuf:"varint,3,rep,packed,name=AutoSplitQuality,proto3" json:"AutoSplitQuality"` // 自动分解品质
	EquipmentList     []*PBEquipment `protobuf:"bytes,2,rep,name=equipmentList,proto3" json:"equipmentList"`               // 装备列表
	HookEquipmentList []*PBEquipment `protobuf:"bytes,7,rep,name=HookEquipmentList,proto3" json:"HookEquipmentList"`       // 挂机装备列表
	OrderId           uint32         `protobuf:"varint,1,opt,name=OrderId,proto3" json:"OrderId"`                          // 序列id
	PosBuyCount       uint32         `protobuf:"varint,4,opt,name=PosBuyCount,proto3" json:"PosBuyCount"`                  // 购买格子次数
	SplitAddBoxCount  uint32         `protobuf:"varint,6,opt,name=SplitAddBoxCount,proto3" json:"SplitAddBoxCount"`        // 分解增加宝箱个数
	SplitEquipCount   uint32         `protobuf:"varint,8,opt,name=SplitEquipCount,proto3" json:"SplitEquipCount"`          // 分解增加宝箱个数
	SplitScore        uint32         `protobuf:"varint,5,opt,name=SplitScore,proto3" json:"SplitScore"`                    // 分解积分
	TotalSplitScore   uint64         `protobuf:"varint,9,opt,name=TotalSplitScore,proto3" json:"TotalSplitScore"`          // 累计积分
}
type PBPlayerGiftCode struct {
	Acode string `protobuf:"bytes,1,opt,name=Acode,proto3" json:"Acode"`  // 兑换码
	Count uint32 `protobuf:"varint,2,opt,name=Count,proto3" json:"Count"` // 次数
	Time  uint64 `protobuf:"varint,3,opt,name=Time,proto3" json:"Time"`   // 兑换时间
}
type PBPlayerHero struct {
	FightPower           uint32            `protobuf:"varint,4,opt,name=FightPower,proto3" json:"FightPower"`                     // 战斗力
	HeroBookList         []*PBHeroBook     `protobuf:"bytes,6,rep,name=HeroBookList,proto3" json:"HeroBookList"`                  // 英雄图鉴列表
	HeroList             []*PBHero         `protobuf:"bytes,2,rep,name=HeroList,proto3" json:"HeroList"`                          // 英雄列表
	MaxHistoryFightPower uint32            `protobuf:"varint,7,opt,name=MaxHistoryFightPower,proto3" json:"MaxHistoryFightPower"` // 历史最大战斗力
	OrderId              uint32            `protobuf:"varint,1,opt,name=OrderId,proto3" json:"OrderId"`                           // 自增ID
	TeamList             []*PBHeroTeamList `protobuf:"bytes,3,rep,name=TeamList,proto3" json:"TeamList"`                          // 编队列表
	UpStarCount          []*PBU32U32       `protobuf:"bytes,5,rep,name=UpStarCount,proto3" json:"UpStarCount"`                    // 升星次数
}
type PBPlayerMail struct {
	AllOrderId uint32    `protobuf:"varint,3,opt,name=AllOrderId,proto3" json:"AllOrderId"` // 全服邮件领取序号
	MailList   []*PBMail `protobuf:"bytes,1,rep,name=MailList,proto3" json:"MailList"`      // 邮件数据
	OrderId    uint32    `protobuf:"varint,2,opt,name=OrderId,proto3" json:"OrderId"`       // 邮件序号
}
type PBPlayerSystem struct {
	Activity     *PBPlayerSystemActivity     `protobuf:"bytes,15,opt,name=Activity,proto3" json:"Activity"`         // 活动
	Battle       *PBPlayerSystemBattle       `protobuf:"bytes,4,opt,name=Battle,proto3" json:"Battle"`              // 战斗
	Box          *PBPlayerSystemBox          `protobuf:"bytes,5,opt,name=Box,proto3" json:"Box"`                    // 宝箱系统
	Championship *PBPlayerSystemChampionship `protobuf:"bytes,14,opt,name=Championship,proto3" json:"Championship"` // 锦标赛
	Charge       *PBPlayerSystemCharge       `protobuf:"bytes,8,opt,name=Charge,proto3" json:"Charge"`              // 充值
	Common       *PBPlayerSystemCommon       `protobuf:"bytes,1,opt,name=Common,proto3" json:"Common"`              // 用户通用数据
	Draw         *PBPlayerSystemDraw         `protobuf:"bytes,7,opt,name=Draw,proto3" json:"Draw"`                  // 抽奖
	Gene         *PBPlayerSystemGene         `protobuf:"bytes,9,opt,name=Gene,proto3" json:"Gene"`                  // 基因系统
	HookTech     *PBPlayerSystemHookTech     `protobuf:"bytes,11,opt,name=HookTech,proto3" json:"HookTech"`         // 挂机科技系统
	Offline      *PBPlayerSystemOffline      `protobuf:"bytes,10,opt,name=Offline,proto3" json:"Offline"`           // 离线系统
	Prof         *PBPlayerSystemProfession   `protobuf:"bytes,3,opt,name=Prof,proto3" json:"Prof"`                  // 职业
	SevenDay     *PBPlayerSystemSevenDay     `protobuf:"bytes,12,opt,name=SevenDay,proto3" json:"SevenDay"`         // 七天活动
	Shop         *PBPlayerSystemShop         `protobuf:"bytes,6,opt,name=Shop,proto3" json:"Shop"`                  // 商店
	Task         *PBPlayerSystemTask         `protobuf:"bytes,2,opt,name=Task,proto3" json:"Task"`                  // 任务信息
	WorldBoss    *PBPlayerSystemWorldBoss    `protobuf:"bytes,13,opt,name=WorldBoss,proto3" json:"WorldBoss"`       // 世界boss
}
type PBPlayerSystemActivity struct {
	ActivityList       []*PBPlayerActivityInfo     `protobuf:"bytes,1,rep,name=ActivityList,proto3" json:"ActivityList"`             // 活动列表
	AdventureList      []*PBActivityAdventure      `protobuf:"bytes,4,rep,name=AdventureList,proto3" json:"AdventureList"`           // 冒险奖励列表
	GiftList           []*PBActivityChargeGift     `protobuf:"bytes,3,rep,name=GiftList,proto3" json:"GiftList"`                     // 直购礼包
	GrowRoadList       []*PBActivityGrowRoadInfo   `protobuf:"bytes,2,rep,name=GrowRoadList,proto3" json:"GrowRoadList"`             // 成长之路列表
	OpenServerGiftList []*PBActivityOpenServerGift `protobuf:"bytes,5,rep,name=OpenServerGiftList,proto3" json:"OpenServerGiftList"` // 开服特惠礼包
}
type PBPlayerSystemBattle struct {
	BattleHook    *PBBattleHookInfo   `protobuf:"bytes,2,opt,name=BattleHook,proto3" json:"BattleHook"`       // 挂机关卡数据
	BattleNormal  *PBBattleNormalInfo `protobuf:"bytes,1,opt,name=BattleNormal,proto3" json:"BattleNormal"`   // 精英关卡数据
	Battlechedule *PBBattleSchedule   `protobuf:"bytes,3,opt,name=Battlechedule,proto3" json:"Battlechedule"` // 战场进度
	Relive        *PBBattleRelive     `protobuf:"bytes,4,opt,name=Relive,proto3" json:"Relive"`               // 复活次数
}
type PBPlayerSystemBox struct {
	BoxScore     uint32       `protobuf:"varint,1,opt,name=BoxScore,proto3" json:"BoxScore"`         // 宝箱积分
	Boxs         []*PBBoxInfo `protobuf:"bytes,4,rep,name=Boxs,proto3" json:"Boxs"`                  // 当前未开宝箱类型
	CurrentLevel uint32       `protobuf:"varint,2,opt,name=CurrentLevel,proto3" json:"CurrentLevel"` // 当前等级
	RecycleTimes uint32       `protobuf:"varint,3,opt,name=RecycleTimes,proto3" json:"RecycleTimes"` // 宝箱积分循环次数
}
type PBPlayerSystemChampionship struct {
	Battle          *PBTaskStageInfo `protobuf:"bytes,6,opt,name=Battle,proto3" json:"Battle"`                    // 关卡排行榜任务
	BattleHasReward uint32           `protobuf:"varint,2,opt,name=BattleHasReward,proto3" json:"BattleHasReward"` // 是否已经领奖
	Damage          *PBTaskStageInfo `protobuf:"bytes,7,opt,name=Damage,proto3" json:"Damage"`                    // 试炼排行榜任务
	DamageHasReward uint32           `protobuf:"varint,3,opt,name=DamageHasReward,proto3" json:"DamageHasReward"` // 是否已经领奖
	Level           *PBTaskStageInfo `protobuf:"bytes,5,opt,name=Level,proto3" json:"Level"`                      // 等级排行榜任务
	LevelHasReward  uint32           `protobuf:"varint,1,opt,name=LevelHasReward,proto3" json:"LevelHasReward"`   // 0表示没加入,1表示加入,2表示领完奖
	Power           *PBTaskStageInfo `protobuf:"bytes,8,opt,name=Power,proto3" json:"Power"`                      // 战力排行榜任务
	PowerHasReward  uint32           `protobuf:"varint,4,opt,name=PowerHasReward,proto3" json:"PowerHasReward"`   // 是否已经领奖
}
type PBPlayerSystemCharge struct {
	BPList          []*PBBPInfo      `protobuf:"bytes,3,rep,name=BPList,proto3" json:"BPList"`                   // bp
	CardList        []*PBChargeCard  `protobuf:"bytes,4,rep,name=CardList,proto3" json:"CardList"`               // 充值卡
	Charge          *PBCharge        `protobuf:"bytes,1,opt,name=Charge,proto3" json:"Charge"`                   // 充值数据
	FirstChargeList []*PBFirstCharge `protobuf:"bytes,2,rep,name=FirstChargeList,proto3" json:"FirstChargeList"` // 首冲数据
}
type PBPlayerSystemCommon struct {
	AdvertList          []*PBAdvertInfo     `protobuf:"bytes,7,rep,name=AdvertList,proto3" json:"AdvertList"`                           // 广告信息
	AvatarFrames        []*PBAvatarFrame    `protobuf:"bytes,5,rep,name=AvatarFrames,proto3" json:"AvatarFrames"`                       // 头像框
	Avatars             []*PBAvatar         `protobuf:"bytes,4,rep,name=Avatars,proto3" json:"Avatars"`                                 // 头像
	GiftCode            []*PBPlayerGiftCode `protobuf:"bytes,2,rep,name=GiftCode,proto3" json:"GiftCode"`                               // 兑换码
	LastChatTime        uint64              `protobuf:"varint,1,opt,name=LastChatTime,proto3" json:"LastChatTime"`                      // 上次聊天时间
	MaxNoticeId         uint32              `protobuf:"varint,6,opt,name=MaxNoticeId,proto3" json:"MaxNoticeId"`                        // 最大的公告ID
	SystemOpenIds       []uint32            `protobuf:"varint,3,rep,packed,name=SystemOpenIds,proto3" json:"SystemOpenIds"`             // 系统开关列表
	SystemOpenPrizeList []uint32            `protobuf:"varint,8,rep,packed,name=SystemOpenPrizeList,proto3" json:"SystemOpenPrizeList"` // 系统开关领取列表
}
type PBPlayerSystemDraw struct {
	DrawList []*PBDrawInfo `protobuf:"bytes,1,rep,name=DrawList,proto3" json:"DrawList"` // 抽奖信息
}
type PBPlayerSystemGene struct {
	SchemeID uint32          `protobuf:"varint,1,opt,name=SchemeID,proto3" json:"SchemeID"` // 当前基因方案
	Schemes  []*PBGeneScheme `protobuf:"bytes,2,rep,name=Schemes,proto3" json:"Schemes"`    // 基因方案
}
type PBPlayerSystemHookTech struct {
	HookTechList []*PBHookTech `protobuf:"bytes,1,rep,name=HookTechList,proto3" json:"HookTechList"` // 挂机科技列表
}
type PBPlayerSystemOffline struct {
	AddEquipmentBag     uint32           `protobuf:"varint,6,opt,name=AddEquipmentBag,proto3" json:"AddEquipmentBag"`         // 加入背包数量
	IncomTime           uint32           `protobuf:"varint,3,opt,name=IncomTime,proto3" json:"IncomTime"`                     // 离线收益时长秒
	LoginTime           uint64           `protobuf:"varint,1,opt,name=LoginTime,proto3" json:"LoginTime"`                     // 上次登录时间
	LogoutTime          uint64           `protobuf:"varint,2,opt,name=LogoutTime,proto3" json:"LogoutTime"`                   // 登出时间
	MaxIncomTime        uint32           `protobuf:"varint,8,opt,name=MaxIncomTime,proto3" json:"MaxIncomTime"`               // 最大离线收益时长秒
	Rewards             []*PBAddItemData `protobuf:"bytes,4,rep,name=Rewards,proto3" json:"Rewards"`                          // 离线收益
	SplitEquipmentScore uint64           `protobuf:"varint,7,opt,name=SplitEquipmentScore,proto3" json:"SplitEquipmentScore"` // 分解装备积分
	TotalEquipment      uint32           `protobuf:"varint,5,opt,name=TotalEquipment,proto3" json:"TotalEquipment"`           // 装备数量
}
type PBPlayerSystemProfInfo struct {
	Grade        uint32                        `protobuf:"varint,3,opt,name=Grade,proto3" json:"Grade"`               // 突破等级
	LastLinkStar uint32                        `protobuf:"varint,6,opt,name=LastLinkStar,proto3" json:"LastLinkStar"` // 上一次套装品阶
	Level        uint32                        `protobuf:"varint,2,opt,name=Level,proto3" json:"Level"`               // 职业等级
	PartList     []*PBPlayerSystemProfPartInfo `protobuf:"bytes,5,rep,name=PartList,proto3" json:"PartList"`          // 职业装备信息
	PeakLevel    uint32                        `protobuf:"varint,4,opt,name=PeakLevel,proto3" json:"PeakLevel"`       // 巅峰等级
	ProfType     uint32                        `protobuf:"varint,1,opt,name=ProfType,proto3" json:"ProfType"`         // 职业类型  EmProfessionType
}
type PBPlayerSystemProfPartInfo struct {
	EquipSn    uint32 `protobuf:"varint,3,opt,name=EquipSn,proto3" json:"EquipSn"`       // 装备Sn 0表示无
	Level      uint32 `protobuf:"varint,2,opt,name=Level,proto3" json:"Level"`           // 强化等级
	Part       uint32 `protobuf:"varint,1,opt,name=Part,proto3" json:"Part"`             // 部位ID EmEquipPartType
	Refine     uint32 `protobuf:"varint,4,opt,name=Refine,proto3" json:"Refine"`         // 精炼等级
	RefineTupo uint32 `protobuf:"varint,5,opt,name=RefineTupo,proto3" json:"RefineTupo"` // 精炼突破等级
}
type PBPlayerSystemProfession struct {
	ProfList []*PBPlayerSystemProfInfo `protobuf:"bytes,1,rep,name=ProfList,proto3" json:"ProfList"` // 职业信息
}
type PBPlayerSystemSevenDay struct {
	SevenDayList []*PBSevenDayInfo `protobuf:"bytes,1,rep,name=SevenDayList,proto3" json:"SevenDayList"` // 活动列表
}
type PBPlayerSystemShop struct {
	BlackShop *PBBlackShop  `protobuf:"bytes,1,opt,name=BlackShop,proto3" json:"BlackShop"` // 黑市商店
	ShopList  []*PBShopInfo `protobuf:"bytes,2,rep,name=ShopList,proto3" json:"ShopList"`   // 商店列表
}
type PBPlayerSystemTask struct {
	AchieveList []*PBAchieveInfo `protobuf:"bytes,3,rep,name=AchieveList,proto3" json:"AchieveList"` // 成就
	DailyTask   *PBDailyTask     `protobuf:"bytes,2,opt,name=DailyTask,proto3" json:"DailyTask"`     // 每日任务
	MainTask    *PBTaskStageInfo `protobuf:"bytes,1,opt,name=MainTask,proto3" json:"MainTask"`       // 主线任务
}
type PBPlayerSystemWorldBoss struct {
	BossId            uint32 `protobuf:"varint,1,opt,name=BossId,proto3" json:"BossId"`                       // bossid
	DailyBuyCount     uint32 `protobuf:"varint,5,opt,name=DailyBuyCount,proto3" json:"DailyBuyCount"`         // 购买次数
	DailyEnterCount   uint32 `protobuf:"varint,6,opt,name=DailyEnterCount,proto3" json:"DailyEnterCount"`     // 进入次数
	DailyMaxDamage    uint64 `protobuf:"varint,3,opt,name=DailyMaxDamage,proto3" json:"DailyMaxDamage"`       // 最大伤害值
	DailyPrizeCount   uint32 `protobuf:"varint,7,opt,name=DailyPrizeCount,proto3" json:"DailyPrizeCount"`     // 结算次数
	DailyPrizeStageId uint32 `protobuf:"varint,2,opt,name=DailyPrizeStageId,proto3" json:"DailyPrizeStageId"` // 领取的进度ID
	DailyTotalDamage  uint64 `protobuf:"varint,4,opt,name=DailyTotalDamage,proto3" json:"DailyTotalDamage"`   // 每日累计伤害值
	MaxDamage         uint64 `protobuf:"varint,8,opt,name=MaxDamage,proto3" json:"MaxDamage"`                 // 历史最大伤害值
}
type PBPropInfo struct {
	PropId uint32 `protobuf:"varint,1,opt,name=PropId,proto3" json:"PropId"` // 属性ID
	Value  uint32 `protobuf:"varint,2,opt,name=Value,proto3" json:"Value"`   // 属性值
}
type PBSevenDayInfo struct {
	BeginTime  uint64             `protobuf:"varint,2,opt,name=BeginTime,proto3" json:"BeginTime"`   // 活动开始时间
	EndTime    uint64             `protobuf:"varint,3,opt,name=EndTime,proto3" json:"EndTime"`       // 活动结束时间
	GiftList   []*PBU32U32        `protobuf:"bytes,7,rep,name=GiftList,proto3" json:"GiftList"`      // 礼包列表 礼包ID|数量
	Id         uint32             `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`                 // 活动ID
	PrizeValue uint32             `protobuf:"varint,5,opt,name=PrizeValue,proto3" json:"PrizeValue"` // 领取的活跃值
	TaskList   []*PBTaskStageInfo `protobuf:"bytes,6,rep,name=TaskList,proto3" json:"TaskList"`      // 任务进度列表
	Value      uint32             `protobuf:"varint,4,opt,name=Value,proto3" json:"Value"`           // 活跃值
}
type PBShopGoodCfg struct {
	AddItem     []*PBAddItemData `protobuf:"bytes,6,rep,name=AddItem,proto3" json:"AddItem"`         // 商品道具
	BuyTimes    uint32           `protobuf:"varint,2,opt,name=BuyTimes,proto3" json:"BuyTimes"`      // 购买次数  钻石商店大于0不展示首充双倍
	Discount    uint32           `protobuf:"varint,4,opt,name=Discount,proto3" json:"Discount"`      // 折扣率
	GoodsID     uint32           `protobuf:"varint,1,opt,name=GoodsID,proto3" json:"GoodsID"`        // 商品ID
	MaxTimes    uint32           `protobuf:"varint,3,opt,name=MaxTimes,proto3" json:"MaxTimes"`      // 最大购买次数
	NeedItem    *PBAddItem       `protobuf:"bytes,5,opt,name=NeedItem,proto3" json:"NeedItem"`       // 需要道具
	Price       uint32           `protobuf:"varint,9,opt,name=Price,proto3" json:"Price"`            // 人民币价格分
	ProductId   uint32           `protobuf:"varint,7,opt,name=ProductId,proto3" json:"ProductId"`    // 充值商品ID
	ProductName string           `protobuf:"bytes,8,opt,name=ProductName,proto3" json:"ProductName"` // 商品名称
	SortTag     uint32           `protobuf:"varint,11,opt,name=SortTag,proto3" json:"SortTag"`       // 排序优先级(从大到小)
	ValueTips   string           `protobuf:"bytes,10,opt,name=ValueTips,proto3" json:"ValueTips"`    // 价值Tips
}
type PBShopGoodInfo struct {
	BuyTimes  uint32       `protobuf:"varint,3,opt,name=BuyTimes,proto3" json:"BuyTimes"`   // 购买次数
	Discount  uint32       `protobuf:"varint,2,opt,name=Discount,proto3" json:"Discount"`   // 折扣率
	Equipment *PBEquipment `protobuf:"bytes,5,opt,name=Equipment,proto3" json:"Equipment"`  // 装备信息
	FreeTimes uint32       `protobuf:"varint,4,opt,name=FreeTimes,proto3" json:"FreeTimes"` // 可获得的免费次数
	GoodsID   uint32       `protobuf:"varint,1,opt,name=GoodsID,proto3" json:"GoodsID"`     // 商品ID
}
type PBShopInfo struct {
	Items           []*PBU32U32 `protobuf:"bytes,3,rep,name=Items,proto3" json:"Items"`                      // 商品购买数据
	NextRefreshTime uint64      `protobuf:"varint,2,opt,name=NextRefreshTime,proto3" json:"NextRefreshTime"` // 下一次刷新时间点
	ShopType        uint32      `protobuf:"varint,1,opt,name=ShopType,proto3" json:"ShopType"`               // 商店类型
}
type PBShopRefreshInfo struct {
	DailyBuyCount       uint32 `protobuf:"varint,1,opt,name=DailyBuyCount,proto3" json:"DailyBuyCount"`             // 每日购买次数
	DailyFreeMaxCount   uint32 `protobuf:"varint,4,opt,name=DailyFreeMaxCount,proto3" json:"DailyFreeMaxCount"`     // 每日免费最大次数
	DailyFreeUseCount   uint32 `protobuf:"varint,2,opt,name=DailyFreeUseCount,proto3" json:"DailyFreeUseCount"`     // 每日免费刷新使用次数
	NextFreeRefreshTime uint64 `protobuf:"varint,3,opt,name=NextFreeRefreshTime,proto3" json:"NextFreeRefreshTime"` // 免费刷新时间 0表示下次刷新
}
type PBStringInt64 struct {
	Key   string `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key"`      // 名称
	Value int64  `protobuf:"varint,2,opt,name=value,proto3" json:"value"` // 值
}
type PBTaskStageInfo struct {
	Id       uint32      `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`                               // 任务ID
	MaxValue uint32      `protobuf:"varint,3,opt,name=MaxValue,proto3" json:"MaxValue"`                   // 最大值
	State    EmTaskState `protobuf:"varint,4,opt,name=State,proto3,enum=common.EmTaskState" json:"State"` // 领取状态 0未达到 1完成 2已经领取
	Value    uint32      `protobuf:"varint,2,opt,name=Value,proto3" json:"Value"`                         // 进度值
}
type PBU32U32 struct {
	Key   uint32 `protobuf:"varint,1,opt,name=Key,proto3" json:"Key"`     // 主键
	Value uint32 `protobuf:"varint,2,opt,name=Value,proto3" json:"Value"` // 值
}
type PBU32U64 struct {
	Key   uint32 `protobuf:"varint,1,opt,name=Key,proto3" json:"Key"`     // 主键
	Value uint64 `protobuf:"varint,2,opt,name=Value,proto3" json:"Value"` // 值
}
type Packet struct {
	Buff      []byte     `protobuf:"bytes,3,opt,name=Buff,proto3" json:"Buff"`           // proto压缩数据
	Id        uint32     `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`              // socketid
	Reply     string     `protobuf:"bytes,2,opt,name=Reply,proto3" json:"Reply"`         // call sessionid
	RpcPacket *RpcPacket `protobuf:"bytes,4,opt,name=RpcPacket,proto3" json:"RpcPacket"` // rpc packet
}
type PageOpenRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PageType   uint32   `protobuf:"varint,2,opt,name=PageType,proto3" json:"PageType"` // 页面类型
}
type PageOpenResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type PassTimeNotify struct {
	CurTime    uint64   `protobuf:"varint,5,opt,name=CurTime,proto3" json:"CurTime"`      // 当前时间戳
	IsDay      bool     `protobuf:"varint,2,opt,name=IsDay,proto3" json:"IsDay"`          // 是否跨天
	IsMonth    bool     `protobuf:"varint,4,opt,name=IsMonth,proto3" json:"IsMonth"`      // 是否跨月
	IsWeek     bool     `protobuf:"varint,3,opt,name=IsWeek,proto3" json:"IsWeek"`        // 是否跨周
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"` // 包头
}
type PlayerData struct {
	Ok         bool   `protobuf:"varint,4,opt,name=Ok,proto3" json:"Ok"`
	PlayerGold int32  `protobuf:"varint,3,opt,name=PlayerGold,proto3" json:"PlayerGold"`
	PlayerID   uint64 `protobuf:"varint,1,opt,name=PlayerID,proto3" json:"PlayerID"`
	PlayerName string `protobuf:"bytes,2,opt,name=PlayerName,proto3" json:"PlayerName"`
}
type PlayerUpdateKvNotify struct {
	ListInfo   []*PBStringInt64 `protobuf:"bytes,2,rep,name=ListInfo,proto3" json:"ListInfo"`     // 多个数据返回
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"` // 包头
}
type Point3F struct {
	X float32 `protobuf:"fixed32,1,opt,name=X,proto3" json:"X"`
	Y float32 `protobuf:"fixed32,2,opt,name=Y,proto3" json:"Y"`
	Z float32 `protobuf:"fixed32,3,opt,name=Z,proto3" json:"Z"`
}
type ProfessionEquipRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PartType   uint32   `protobuf:"varint,3,opt,name=PartType,proto3" json:"PartType"` // 部位类型
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
	Sn         uint32   `protobuf:"varint,4,opt,name=Sn,proto3" json:"Sn"`             // 装备sn 0表示脱
}
type ProfessionEquipResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PartType   uint32   `protobuf:"varint,3,opt,name=PartType,proto3" json:"PartType"` // 部位类型
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
	Sn         uint32   `protobuf:"varint,4,opt,name=Sn,proto3" json:"Sn"`             // 当前部位ID
}
type ProfessionGradeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
}
type ProfessionGradeResponse struct {
	Grade      uint32   `protobuf:"varint,3,opt,name=Grade,proto3" json:"Grade"` // 突破等级
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
}
type ProfessionLevelRequest struct {
	AddLevel   uint32   `protobuf:"varint,4,opt,name=AddLevel,proto3" json:"AddLevel"` // 新增等级
	CurLevel   uint32   `protobuf:"varint,3,opt,name=CurLevel,proto3" json:"CurLevel"` // 当前等级
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
}
type ProfessionLevelResponse struct {
	CurLevel   uint32   `protobuf:"varint,3,opt,name=CurLevel,proto3" json:"CurLevel"` // 新等级
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
}
type ProfessionOnekeyUnEquipRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
}
type ProfessionOnekeyUnEquipResponse struct {
	EquipSnList []uint32 `protobuf:"varint,3,rep,packed,name=EquipSnList,proto3" json:"EquipSnList"` // 装备列表
	PacketHead  *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType    uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
}
type ProfessionPartLevelRequest struct {
	CurLevel   uint32   `protobuf:"varint,4,opt,name=CurLevel,proto3" json:"CurLevel"` // 当前等级
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PartType   uint32   `protobuf:"varint,3,opt,name=PartType,proto3" json:"PartType"` // 部位类型
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
}
type ProfessionPartLevelResponse struct {
	CurLevel   uint32   `protobuf:"varint,4,opt,name=CurLevel,proto3" json:"CurLevel"` // 当前等级
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PartType   uint32   `protobuf:"varint,3,opt,name=PartType,proto3" json:"PartType"` // 部位类型
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
}
type ProfessionPartOnekeyLevelRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
}
type ProfessionPartOnekeyLevelResponse struct {
	PacketHead *IPacket                      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PartList   []*PBPlayerSystemProfPartInfo `protobuf:"bytes,3,rep,name=PartList,proto3" json:"PartList"`  // 职业部位信息
	ProfType   uint32                        `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
}
type ProfessionPartOnekeyRefineRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
}
type ProfessionPartOnekeyRefineResponse struct {
	PacketHead *IPacket                      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PartList   []*PBPlayerSystemProfPartInfo `protobuf:"bytes,3,rep,name=PartList,proto3" json:"PartList"`  // 职业部位信息
	ProfType   uint32                        `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
}
type ProfessionPartRefineRequest struct {
	CurLevel   uint32   `protobuf:"varint,4,opt,name=CurLevel,proto3" json:"CurLevel"` // 当前等级
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PartType   uint32   `protobuf:"varint,3,opt,name=PartType,proto3" json:"PartType"` // 部位类型
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
}
type ProfessionPartRefineResponse struct {
	CurLevel   uint32   `protobuf:"varint,4,opt,name=CurLevel,proto3" json:"CurLevel"` // 当前等级
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PartType   uint32   `protobuf:"varint,3,opt,name=PartType,proto3" json:"PartType"` // 部位类型
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
}
type ProfessionPartRefineTupoRequest struct {
	CurLevel   uint32   `protobuf:"varint,4,opt,name=CurLevel,proto3" json:"CurLevel"` // 当前等级
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PartType   uint32   `protobuf:"varint,3,opt,name=PartType,proto3" json:"PartType"` // 部位类型
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
}
type ProfessionPartRefineTupoResponse struct {
	CurLevel   uint32   `protobuf:"varint,4,opt,name=CurLevel,proto3" json:"CurLevel"` // 当前等级
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PartType   uint32   `protobuf:"varint,3,opt,name=PartType,proto3" json:"PartType"` // 部位类型
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
}
type ProfessionPartResetRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PartType   uint32   `protobuf:"varint,3,opt,name=PartType,proto3" json:"PartType"` // 部位类型 999表示任意
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
}
type ProfessionPartResetResponse struct {
	PacketHead *IPacket                      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PartList   []*PBPlayerSystemProfPartInfo `protobuf:"bytes,3,rep,name=PartList,proto3" json:"PartList"`  // 职业装备信息
	ProfType   uint32                        `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
}
type ProfessionPeakLevelRequest struct {
	CurLevel   uint32   `protobuf:"varint,3,opt,name=CurLevel,proto3" json:"CurLevel"` // 当前巅峰等级
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
}
type ProfessionPeakLevelResponse struct {
	CurLevel   uint32   `protobuf:"varint,3,opt,name=CurLevel,proto3" json:"CurLevel"` // 新巅峰等级
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
}
type ProfessionPeakRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"` // 职业类型
}
type ProfessionPeakResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PeakLevel  uint32   `protobuf:"varint,3,opt,name=PeakLevel,proto3" json:"PeakLevel"` // 新巅峰等级
	ProfType   uint32   `protobuf:"varint,2,opt,name=ProfType,proto3" json:"ProfType"`   // 职业类型
}
type ProtocolNameNotify struct {
	Name       string   `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name"`              // 名称
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`  // 包头
	ProtocolId uint32   `protobuf:"varint,3,opt,name=ProtocolId,proto3" json:"ProtocolId"` // 协议ID
}
type RankData struct {
	CreateTime uint64   `protobuf:"varint,3,opt,name=CreateTime,proto3" json:"CreateTime"`        // 开启时间点
	HasRewards []uint64 `protobuf:"varint,7,rep,packed,name=HasRewards,proto3" json:"HasRewards"` // 已经领取奖励的玩家
	RankType   uint32   `protobuf:"varint,1,opt,name=RankType,proto3" json:"RankType"`            // 排行榜类型
	RegionID   uint32   `protobuf:"varint,2,opt,name=RegionID,proto3" json:"RegionID"`            // 分区ID
}
type RankInfo struct {
	Display *PBPlayerBaseDisplay `protobuf:"bytes,2,opt,name=Display,proto3" json:"Display"` // 基本数据
	Rank    uint32               `protobuf:"varint,1,opt,name=Rank,proto3" json:"Rank"`      // 名次 1开始
	Value   uint64               `protobuf:"varint,3,opt,name=Value,proto3" json:"Value"`    // 数据 通关进度：useTime:4位 rate：5位 stageid:3位 mapid:3位
}
type RankRequest struct {
	Begin      uint32   `protobuf:"varint,3,opt,name=Begin,proto3" json:"Begin"`           // 开始名次
	CreateTime uint64   `protobuf:"varint,5,opt,name=CreateTime,proto3" json:"CreateTime"` // game发送gm服务携带的字段信息
	End        uint32   `protobuf:"varint,4,opt,name=End,proto3" json:"End"`               // 结束名次
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RankType   uint32   `protobuf:"varint,2,opt,name=RankType,proto3" json:"RankType"` // 排行类型
}
type RankResponse struct {
	Begin      uint32      `protobuf:"varint,3,opt,name=Begin,proto3" json:"Begin"` // 开始名次 1开始
	End        uint32      `protobuf:"varint,4,opt,name=End,proto3" json:"End"`     // 结束名次
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RankList   []*RankInfo `protobuf:"bytes,5,rep,name=RankList,proto3" json:"RankList"`    // 玩家数据
	RankType   uint32      `protobuf:"varint,2,opt,name=RankType,proto3" json:"RankType"`   // 排行类型
	SelfInfo   *RankInfo   `protobuf:"bytes,6,opt,name=SelfInfo,proto3" json:"SelfInfo"`    // 自己的数据
	TotalRank  int64       `protobuf:"varint,7,opt,name=TotalRank,proto3" json:"TotalRank"` // 总数据
}
type RankRewardRequest struct {
	CreateTime uint64       `protobuf:"varint,6,opt,name=CreateTime,proto3" json:"CreateTime"`               // 时间点(内部服务节点转发使用)
	Doing      EmDoingType  `protobuf:"varint,4,opt,name=Doing,proto3,enum=common.EmDoingType" json:"Doing"` // 发放奖励原因(内部服务节点转发使用)
	Notify     bool         `protobuf:"varint,3,opt,name=Notify,proto3" json:"Notify"`                       // 是否需要发送恭喜获得(内部服务节点转发使用)
	PacketHead *IPacket     `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RankType   uint32       `protobuf:"varint,2,opt,name=RankType,proto3" json:"RankType"` // 排行类型
	Rewards    []*PBAddItem `protobuf:"bytes,5,rep,name=Rewards,proto3" json:"Rewards"`    // 发放的奖励(内部服务节点转发使用)
}
type RankRewardResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type RankUpdateNotify struct {
	CreateTime uint64   `protobuf:"varint,3,opt,name=CreateTime,proto3" json:"CreateTime"` // 开启时间点
	Member     string   `protobuf:"bytes,7,opt,name=Member,proto3" json:"Member"`          // 待更新成员
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RankType   uint32   `protobuf:"varint,2,opt,name=RankType,proto3" json:"RankType"` // 排行类型
	Score      float64  `protobuf:"fixed64,8,opt,name=Score,proto3" json:"Score"`      // 积分
}
type RouteInfo struct {
	Db     uint32 `protobuf:"varint,4,opt,name=Db,proto3" json:"Db"`
	Dip    uint32 `protobuf:"varint,5,opt,name=Dip,proto3" json:"Dip"`
	Game   uint32 `protobuf:"varint,2,opt,name=Game,proto3" json:"Game"`
	Gate   uint32 `protobuf:"varint,1,opt,name=Gate,proto3" json:"Gate"`
	Gm     uint32 `protobuf:"varint,3,opt,name=Gm,proto3" json:"Gm"`
	Record uint32 `protobuf:"varint,6,opt,name=Record,proto3" json:"Record"`
}
type RpcHead struct {
	ActorName      string     `protobuf:"bytes,10,opt,name=ActorName,proto3" json:"ActorName"`                               // actor名称
	ClusterId      uint32     `protobuf:"varint,2,opt,name=ClusterId,proto3" json:"ClusterId"`                               // 目标集群id 大于0优先
	DestServerType SERVICE    `protobuf:"varint,3,opt,name=DestServerType,proto3,enum=common.SERVICE" json:"DestServerType"` // 目标集群类型
	FuncName       string     `protobuf:"bytes,11,opt,name=FuncName,proto3" json:"FuncName"`                                 // actor函数名称
	Id             uint64     `protobuf:"varint,6,opt,name=Id,proto3" json:"Id"`                                             // 目标ID
	RegionID       uint32     `protobuf:"varint,5,opt,name=RegionID,proto3" json:"RegionID"`                                 // 区服ID
	Reply          string     `protobuf:"bytes,12,opt,name=Reply,proto3" json:"Reply"`                                       // call sessionid
	Route          *RouteInfo `protobuf:"bytes,13,opt,name=Route,proto3" json:"Route"`                                       // 路由表
	RouteType      uint32     `protobuf:"varint,4,opt,name=RouteType,proto3" json:"RouteType"`                               // 路由类型
	SendType       SEND       `protobuf:"varint,8,opt,name=SendType,proto3,enum=rpc3.SEND" json:"SendType"`                  // 消息类型
	SeqId          uint32     `protobuf:"varint,9,opt,name=SeqId,proto3" json:"SeqId"`                                       // 序列号
	SocketId       uint32     `protobuf:"varint,7,opt,name=SocketId,proto3" json:"SocketId"`                                 // 客户端连接id
	SrcClusterId   uint32     `protobuf:"varint,1,opt,name=SrcClusterId,proto3" json:"SrcClusterId"`                         // 源集群id
}
type RpcPacket struct {
	ArgLen  uint32   `protobuf:"varint,2,opt,name=ArgLen,proto3" json:"ArgLen"`  // 可变参数个数
	RpcBody []byte   `protobuf:"bytes,3,opt,name=RpcBody,proto3" json:"RpcBody"` // 存储可变参数 pb
	RpcHead *RpcHead `protobuf:"bytes,1,opt,name=RpcHead,proto3" json:"RpcHead"`
}
type SetClientRequest struct {
	ClientData *PBClientData `protobuf:"bytes,2,opt,name=ClientData,proto3" json:"ClientData"`
	PacketHead *IPacket      `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type SetClientResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type SetPlayerDataRequest struct {
	DataType   int32    `protobuf:"varint,2,opt,name=DataType,proto3" json:"DataType"`
	JsonData   string   `protobuf:"bytes,3,opt,name=JsonData,proto3" json:"JsonData"`
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type SetPlayerDataResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type SevenDayActivePrizeRequest struct {
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"` // 活动ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type SevenDayActivePrizeResponse struct {
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"` // 活动ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PrizeValue uint32   `protobuf:"varint,3,opt,name=PrizeValue,proto3" json:"PrizeValue"` // 当前领取的活跃值
}
type SevenDayBuyGiftRequest struct {
	GiftId     uint32   `protobuf:"varint,3,opt,name=GiftId,proto3" json:"GiftId"` // 礼包ID
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`         // 活动ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type SevenDayBuyGiftResponse struct {
	GiftInfo   *PBU32U32 `protobuf:"bytes,4,opt,name=GiftInfo,proto3" json:"GiftInfo"` // 礼包数据
	Id         uint32    `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`            // 活动ID
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Value      uint32    `protobuf:"varint,3,opt,name=Value,proto3" json:"Value"` // 活跃值
}
type SevenDayGiftNotify struct {
	GiftInfo   *PBU32U32 `protobuf:"bytes,2,opt,name=GiftInfo,proto3" json:"GiftInfo"` // 礼包数据
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Value      uint32    `protobuf:"varint,3,opt,name=Value,proto3" json:"Value"` // 活跃值
}
type SevenDayListNotify struct {
	AddList    []*PBSevenDayInfo `protobuf:"bytes,2,rep,name=AddList,proto3" json:"AddList"`       // 增加的活动
	Delist     []uint32          `protobuf:"varint,3,rep,packed,name=Delist,proto3" json:"Delist"` // 删除的活动
	PacketHead *IPacket          `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type SevenDayTaskPrizeRequest struct {
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"` // 活动ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	TaskId     uint32   `protobuf:"varint,3,opt,name=TaskId,proto3" json:"TaskId"` // 任务ID
}
type SevenDayTaskPrizeResponse struct {
	Id            uint32           `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"` // 活动ID
	PacketHead    *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	TaskStageInfo *PBTaskStageInfo `protobuf:"bytes,4,opt,name=TaskStageInfo,proto3" json:"TaskStageInfo"` // 任务数据
	Value         uint32           `protobuf:"varint,3,opt,name=Value,proto3" json:"Value"`                // 活跃值
}
type ShopBuyRequest struct {
	AdvertType uint32   `protobuf:"varint,4,opt,name=AdvertType,proto3" json:"AdvertType"` // 广告类型 0无广告
	GoodsID    uint32   `protobuf:"varint,3,opt,name=GoodsID,proto3" json:"GoodsID"`       // 商品ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopType   uint32   `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"` // 商店类型
}
type ShopBuyResponse struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopType   uint32   `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"` // 商店类型
}
type ShopExchangeGoodNotify struct {
	GoodInfo   *PBU32U32 `protobuf:"bytes,3,opt,name=GoodInfo,proto3" json:"GoodInfo"` // 商品数据 key：商品ID  Value:购买次数
	PacketHead *IPacket  `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopType   uint32    `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"` // 商店类型
}
type ShopExchangeRequest struct {
	GoodsID    uint32   `protobuf:"varint,3,opt,name=GoodsID,proto3" json:"GoodsID"` // 商品ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopType   uint32   `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"` // 商店类型
}
type ShopExchangeResponse struct {
	BuyTimes   uint32   `protobuf:"varint,4,opt,name=BuyTimes,proto3" json:"BuyTimes"` // 购买次数
	GoodsID    uint32   `protobuf:"varint,3,opt,name=GoodsID,proto3" json:"GoodsID"`   // 商品ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopType   uint32   `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"` // 商店类型
}
type ShopListNotify struct {
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopList   []*PBU32U64 `protobuf:"bytes,2,rep,name=ShopList,proto3" json:"ShopList"` // 商店列表 key:商店类型 value:下次刷新时间
}
type ShopOpenRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopType   uint32   `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"` // 商店类型
}
type ShopOpenResponse struct {
	GoodList   []*PBShopGoodCfg `protobuf:"bytes,3,rep,name=GoodList,proto3" json:"GoodList"` // 商品信息
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopType   uint32           `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"` // 商店类型
}
type ShopRefreshRequest struct {
	IsFree     bool     `protobuf:"varint,3,opt,name=IsFree,proto3" json:"IsFree"` // 是否免费
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopType   uint32   `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"` // 商店类型
}
type ShopRefreshResponse struct {
	PacketHead  *IPacket           `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RefreshInfo *PBShopRefreshInfo `protobuf:"bytes,3,opt,name=RefreshInfo,proto3" json:"RefreshInfo"` // 刷新数据
	ShopType    uint32             `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"`      // 商店类型
}
type ShopRefreshTimeNotify struct {
	PacketHead  *IPacket           `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RefreshInfo *PBShopRefreshInfo `protobuf:"bytes,3,opt,name=RefreshInfo,proto3" json:"RefreshInfo"` // 刷新数据
	ShopType    uint32             `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"`      // 商店类型
}
type ShopUpdateNotify struct {
	PacketHead *IPacket            `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Shop       *PBPlayerSystemShop `protobuf:"bytes,3,opt,name=Shop,proto3" json:"Shop"`          // 商店类型
	ShopType   uint32              `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"` // 商店类型
}
type ShopUpdateOneGoodsNotify struct {
	PacketHead *IPacket        `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ShopGood   *PBShopGoodInfo `protobuf:"bytes,3,opt,name=ShopGood,proto3" json:"ShopGood"`  // 商品数据
	ShopType   uint32          `protobuf:"varint,2,opt,name=ShopType,proto3" json:"ShopType"` // 商店类型
}
type StubMailBox struct {
	ClusterId uint32 `protobuf:"varint,5,opt,name=ClusterId,proto3" json:"ClusterId"` // 集群id
	Id        uint64 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`
	LeaseId   uint64 `protobuf:"varint,2,opt,name=LeaseId,proto3" json:"LeaseId"`
	StubType  STUB   `protobuf:"varint,3,opt,name=StubType,proto3,enum=rpc3.STUB" json:"StubType"`
}
type SystemOpenNotify struct {
	PacketHead    *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`               // 包头
	SystemOpenIds []uint32 `protobuf:"varint,2,rep,packed,name=SystemOpenIds,proto3" json:"SystemOpenIds"` // 系统开启ID
}
type SystemOpenPrizeRequest struct {
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"` // 开关ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type SystemOpenPrizeResponse struct {
	Id         uint32   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"` // 开关ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type SystemPropNotify struct {
	PacketHead     *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	PropInfoList   []*PBPropInfo    `protobuf:"bytes,3,rep,name=PropInfoList,proto3" json:"PropInfoList"`                                   // 属性数据
	SyetemPropType EmSyetemPropType `protobuf:"varint,2,opt,name=SyetemPropType,proto3,enum=common.EmSyetemPropType" json:"SyetemPropType"` // 系统类型
}
type UpdateRouteNotify struct {
	PacketHead *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	Users      []*UserRouteInfo `protobuf:"bytes,2,rep,name=Users,proto3" json:"Users"`
}
type UserRouteInfo struct {
	RouteInfo *RouteInfo `protobuf:"bytes,2,opt,name=RouteInfo,proto3" json:"RouteInfo"`
	UID       uint64     `protobuf:"varint,1,opt,name=UID,proto3" json:"UID"`
}
type WorldBossBattleBeginRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type WorldBossBattleBeginResponse struct {
	DailyEnterCount uint32   `protobuf:"varint,2,opt,name=DailyEnterCount,proto3" json:"DailyEnterCount"` // 使用免费次数 +1 最大三次
	PacketHead      *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type WorldBossBattleEndRequest struct {
	Battle     *BattleInfo `protobuf:"bytes,2,opt,name=Battle,proto3" json:"Battle"`
	IsFinish   uint32      `protobuf:"varint,3,opt,name=IsFinish,proto3" json:"IsFinish"` // 是否最终上报
	PacketHead *IPacket    `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type WorldBossBattleEndResponse struct {
	DailyEnterCount  uint32           `protobuf:"varint,6,opt,name=DailyEnterCount,proto3" json:"DailyEnterCount"`   // 免费次数 -1
	DailyMaxDamage   uint64           `protobuf:"varint,2,opt,name=DailyMaxDamage,proto3" json:"DailyMaxDamage"`     // 最大伤害值
	DailyPrizeCount  uint32           `protobuf:"varint,4,opt,name=DailyPrizeCount,proto3" json:"DailyPrizeCount"`   // 结算次数
	DailyTotalDamage uint64           `protobuf:"varint,3,opt,name=DailyTotalDamage,proto3" json:"DailyTotalDamage"` // 每日累计伤害值
	ItemInfo         []*PBAddItemData `protobuf:"bytes,5,rep,name=ItemInfo,proto3" json:"ItemInfo"`                  // 道具信息
	PacketHead       *IPacket         `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type WorldBossBuyCountRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type WorldBossBuyCountResponse struct {
	DailyBuyCount uint32   `protobuf:"varint,2,opt,name=DailyBuyCount,proto3" json:"DailyBuyCount"` // 购买次数
	PacketHead    *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type WorldBossNotify struct {
	PacketHead *IPacket                 `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	WorldBoss  *PBPlayerSystemWorldBoss `protobuf:"bytes,2,opt,name=WorldBoss,proto3" json:"WorldBoss"`
}
type WorldBossRecordRequest struct {
	AccountId  uint64   `protobuf:"varint,2,opt,name=AccountId,proto3" json:"AccountId"` // 玩家ID
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	ServerId   uint32   `protobuf:"varint,3,opt,name=ServerId,proto3" json:"ServerId"` // 服务器ID
}
type WorldBossRecordResponse struct {
	PacketHead *IPacket            `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
	RecordInfo *PBPlayerBattleData `protobuf:"bytes,2,opt,name=RecordInfo,proto3" json:"RecordInfo"` // 记录数据
}
type WorldBossStagePrizeRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type WorldBossStagePrizeResponse struct {
	DailyPrizeStageId uint32   `protobuf:"varint,2,opt,name=DailyPrizeStageId,proto3" json:"DailyPrizeStageId"` // 阶段ID
	PacketHead        *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type WorldBossSweepRequest struct {
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type WorldBossSweepResponse struct {
	DailyEnterCount uint32   `protobuf:"varint,2,opt,name=DailyEnterCount,proto3" json:"DailyEnterCount"` // 进入次数
	DailyPrizeCount uint32   `protobuf:"varint,3,opt,name=DailyPrizeCount,proto3" json:"DailyPrizeCount"` // 奖励次数
	PacketHead      *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type Z_C_ENTITY struct {
	EntityInfo []*Z_C_ENTITY_Entity `protobuf:"bytes,2,rep,name=EntityInfo,proto3" json:"EntityInfo"`
	PacketHead *IPacket             `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
}
type Z_C_ENTITY_Entity struct {
	Data  *Z_C_ENTITY_Entity_DataMask  `protobuf:"bytes,2,opt,name=Data,proto3" json:"Data"`   // 初始化信息
	Id    uint64                       `protobuf:"varint,1,opt,name=Id,proto3" json:"Id"`      // 唯一id
	Move  *Z_C_ENTITY_Entity_MoveMask  `protobuf:"bytes,3,opt,name=Move,proto3" json:"Move"`   // 移动
	Stats *Z_C_ENTITY_Entity_StatsMask `protobuf:"bytes,4,opt,name=Stats,proto3" json:"Stats"` // 基础属性
}
type Z_C_ENTITY_Entity_DataMask struct {
	NpcData    *Z_C_ENTITY_Entity_DataMask_NpcDataMask   `protobuf:"bytes,3,opt,name=NpcData,proto3" json:"NpcData"`        // npc初始信息
	RemoveFlag bool                                      `protobuf:"varint,2,opt,name=RemoveFlag,proto3" json:"RemoveFlag"` // 销毁状态
	SpellData  *Z_C_ENTITY_Entity_DataMask_SpellDataMask `protobuf:"bytes,4,opt,name=SpellData,proto3" json:"SpellData"`    // effect初始信息
	Type       int32                                     `protobuf:"varint,1,opt,name=Type,proto3" json:"Type"`             // 类型，player(9)，npc(5)
}
type Z_C_ENTITY_Entity_DataMask_NpcDataMask struct {
	DataId int32 `protobuf:"varint,1,opt,name=DataId,proto3" json:"DataId"` // 模板id
}
type Z_C_ENTITY_Entity_DataMask_SpellDataMask struct {
	DataId   int32 `protobuf:"varint,1,opt,name=DataId,proto3" json:"DataId"`     // 特效id
	FlySpeed int32 `protobuf:"varint,3,opt,name=FlySpeed,proto3" json:"FlySpeed"` // 飞行速度
	LifeTime int32 `protobuf:"varint,2,opt,name=LifeTime,proto3" json:"LifeTime"` // 存在事件
}
type Z_C_ENTITY_Entity_MoveMask struct {
	ForceUpdateFlag bool     `protobuf:"varint,3,opt,name=ForceUpdateFlag,proto3" json:"ForceUpdateFlag"` // 强制拉回
	Pos             *Point3F `protobuf:"bytes,1,opt,name=Pos,proto3" json:"Pos"`
	Rotation        float32  `protobuf:"fixed32,2,opt,name=Rotation,proto3" json:"Rotation"`
}
type Z_C_ENTITY_Entity_StatsMask struct {
	AntiCritical      int32 `protobuf:"varint,13,opt,name=AntiCritical,proto3" json:"AntiCritical"`           // 防暴击全局百分比几率
	AntiCriticalTimes int32 `protobuf:"varint,12,opt,name=AntiCriticalTimes,proto3" json:"AntiCriticalTimes"` // 抗暴击全局百分比伤害倍数
	Critical          int32 `protobuf:"varint,11,opt,name=Critical,proto3" json:"Critical"`                   // 暴击全局百分比几率
	CriticalTimes     int32 `protobuf:"varint,10,opt,name=CriticalTimes,proto3" json:"CriticalTimes"`         // 暴击全局百分比伤害倍数
	Dodge             int32 `protobuf:"varint,14,opt,name=Dodge,proto3" json:"Dodge"`                         // 闪避全局百分比几率
	HP                int32 `protobuf:"varint,1,opt,name=HP,proto3" json:"HP"`                                // 生命
	Heal              int32 `protobuf:"varint,9,opt,name=Heal,proto3" json:"Heal"`                            // 治疗
	Hit               int32 `protobuf:"varint,15,opt,name=Hit,proto3" json:"Hit"`                             // 命中全局百分比几率
	MP                int32 `protobuf:"varint,2,opt,name=MP,proto3" json:"MP"`                                // 真气
	MaxHP             int32 `protobuf:"varint,3,opt,name=MaxHP,proto3" json:"MaxHP"`                          // 最大生命
	MaxMP             int32 `protobuf:"varint,4,opt,name=MaxMP,proto3" json:"MaxMP"`                          // 最大真气
	PhyDamage         int32 `protobuf:"varint,5,opt,name=PhyDamage,proto3" json:"PhyDamage"`                  // 物理伤害
	PhyDefence        int32 `protobuf:"varint,6,opt,name=PhyDefence,proto3" json:"PhyDefence"`                // 物理防御
	SplDamage         int32 `protobuf:"varint,7,opt,name=SplDamage,proto3" json:"SplDamage"`                  // 法术伤害
	SplDefence        int32 `protobuf:"varint,8,opt,name=SplDefence,proto3" json:"SplDefence"`                // 法术防御
}
type Z_C_LoginMap struct {
	Id         uint64   `protobuf:"varint,2,opt,name=Id,proto3" json:"Id"`
	PacketHead *IPacket `protobuf:"bytes,1,opt,name=PacketHead,proto3" json:"PacketHead"`
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
	EmBagType_BagType_Equip    EmBagType = 0 // 装备背包类型
	EmBagType_BagType_Item     EmBagType = 1 // 道具背包类型
	EmBagType_BagType_PetPiece EmBagType = 2 // 碎片背包类型
	EmBagType_BagType_Special  EmBagType = 3 // 特殊背包类型
	EmBagType_BagType_GodEquip EmBagType = 4 // 神装背包类型
)

type EmBattleType int32

const (
	EmBattleType_EBT_None   EmBattleType = 0 // 默认
	EmBattleType_EBT_Normal EmBattleType = 1 // 精英关卡
	EmBattleType_EBT_Tower  EmBattleType = 2 // 爬塔
	EmBattleType_EBT_Hook   EmBattleType = 3 // 挂机
)

type EmDoingType int32

const (
	EmDoingType_EDT_Gm                       EmDoingType = 0  // gm后台添加
	EmDoingType_EDT_Other                    EmDoingType = 1  // 其它
	EmDoingType_EDT_Client                   EmDoingType = 2  // 客户端请求
	EmDoingType_EDT_ItemUse                  EmDoingType = 3  // 道具使用
	EmDoingType_EDT_GiftCode                 EmDoingType = 4  // 兑换码
	EmDoingType_EDT_ProfessionLevel          EmDoingType = 5  // 职业升级
	EmDoingType_EDT_ProfessionPeakLevel      EmDoingType = 6  // 职业巅峰升级
	EmDoingType_EDT_ProfessionPartLevel      EmDoingType = 7  // 职业部位升级
	EmDoingType_EDT_EquipSplit               EmDoingType = 8  // 装备分解
	EmDoingType_EDT_HeroAwaken               EmDoingType = 9  // 英雄技能觉醒
	EmDoingType_EDT_HeroRebirth              EmDoingType = 10 // 英雄重生
	EmDoingType_EDT_BoxScoreExchange         EmDoingType = 11 // 宝箱积分兑换
	EmDoingType_EDT_Battle                   EmDoingType = 12 // 战斗
	EmDoingType_EDT_BoxConsume               EmDoingType = 13 // 宝箱消耗
	EmDoingType_EDT_BoxOpen                  EmDoingType = 14 // 开宝箱
	EmDoingType_EDT_Task                     EmDoingType = 15 // 任务
	EmDoingType_EDT_System                   EmDoingType = 16 // 系统增加
	EmDoingType_EDT_Login                    EmDoingType = 17 // 登录
	EmDoingType_EDT_ChangePlayerName         EmDoingType = 18 // 修改名字
	EmDoingType_EDT_Mail                     EmDoingType = 19 // 邮件
	EmDoingType_EDT_BlackShop                EmDoingType = 20 // 黑市商店
	EmDoingType_EDT_Advert                   EmDoingType = 21 // 广告
	EmDoingType_EDT_CrystalRobotUpgrade      EmDoingType = 22 // 晶核系统机器人升级
	EmDoingType_EDT_CrystalRedefine          EmDoingType = 23 // 晶核改造消耗
	EmDoingType_EDT_CrystalBookUpgrade       EmDoingType = 24 // 晶核图鉴系统升级
	EmDoingType_EDT_Draw                     EmDoingType = 25 // 抽奖
	EmDoingType_EDT_Charge                   EmDoingType = 26 // 充值
	EmDoingType_EDT_BattleHook               EmDoingType = 27 // 挂机关卡
	EmDoingType_EDT_Offline                  EmDoingType = 28 // 离线收益
	EmDoingType_EDT_ProfessionPartRefine     EmDoingType = 29 // 职业部位精炼
	EmDoingType_EDT_ProfessionPartRefineTupo EmDoingType = 30 // 职业部位精炼突破
	EmDoingType_EDT_HeroBook                 EmDoingType = 31 // 英雄图鉴
	EmDoingType_EDT_StarSource               EmDoingType = 32 // 挂机科技
	EmDoingType_EDT_SevenDay                 EmDoingType = 33 // 七天登录
	EmDoingType_EDT_Shop                     EmDoingType = 34 // 商店
	EmDoingType_EDT_DailyTask                EmDoingType = 35 // 每日任务
	EmDoingType_EDT_Reset                    EmDoingType = 36 // 部位重置
	EmDoingType_EDT_Entry                    EmDoingType = 37 // 词条
	EmDoingType_EDT_BattleNormal             EmDoingType = 38 // 精英官咖
	EmDoingType_EDT_StarSourceDraw           EmDoingType = 39 // 星源抽奖
	EmDoingType_EDT_RankReward               EmDoingType = 40 // 排行榜奖励
	EmDoingType_EDT_WorldBoss                EmDoingType = 41 // 世界boss
	EmDoingType_EDT_Championship             EmDoingType = 42 // 锦标赛
	EmDoingType_EDT_FirstCharge              EmDoingType = 43 // 首冲
	EmDoingType_EDT_BP                       EmDoingType = 44 // bp
	EmDoingType_EDT_ChargeCard               EmDoingType = 45 // 充值卡
	EmDoingType_EDT_ChargeGift               EmDoingType = 46 // 直购礼包
	EmDoingType_EDT_GrowRoad                 EmDoingType = 47 // 成长之路
	EmDoingType_EDT_HookTech                 EmDoingType = 48 // 星源抽奖
	EmDoingType_EDT_CrystalUpgrade           EmDoingType = 49 // 晶核升级
	EmDoingType_EDT_Adventure                EmDoingType = 50 // 冒险奖励
	EmDoingType_EDT_ItemBuy                  EmDoingType = 51 // 道具购买
	EmDoingType_EDT_Activity                 EmDoingType = 52 // 活动
)

type EmGiftCodeType int32

const (
	EmGiftCodeType_GAT_Common EmGiftCodeType = 0 // 通用礼包
	EmGiftCodeType_GAT_Code   EmGiftCodeType = 1 // 兑换码礼包
	EmGiftCodeType_GAT_Week   EmGiftCodeType = 2 // 周礼包
	EmGiftCodeType_GAT_Month  EmGiftCodeType = 3 // 月礼包
)

type EmGmFuncType int32

const (
	EmGmFuncType_GFT_AddItem  EmGmFuncType = 0 // 增加道具 参数：类型 ID 数量
	EmGmFuncType_GFT_AddEquip EmGmFuncType = 1 // 增加装备 参数：装备id 品质 星级 数量
	EmGmFuncType_GFT_AddHero  EmGmFuncType = 2 // 增加英雄 参数：英雄id 星级 数量
	EmGmFuncType_GFT_NB       EmGmFuncType = 3 // 牛逼
	EmGmFuncType_GFT_Relogin  EmGmFuncType = 4 // 重登陆
	EmGmFuncType_GFT_Charge   EmGmFuncType = 5 // 直冲
)

type EmGmParamType int32

const (
	EmGmParamType_GPT_None   EmGmParamType = 0 // gm是否开启
	EmGmParamType_GPT_GmOpen EmGmParamType = 1 // gm是否开启
)

type EmItemExpendType int32

const (
	EmItemExpendType_EIET_None EmItemExpendType = 0   // 无
	EmItemExpendType_EIET_Cash EmItemExpendType = 1   // 元宝
	EmItemExpendType_EIET_Gold EmItemExpendType = 2   // 金币
	EmItemExpendType_EIET_Max  EmItemExpendType = 100 // 最大值
)

type EmMailState int32

const (
	EmMailState_NoRead      EmMailState = 0 // 未读
	EmMailState_ReadRecieve EmMailState = 1 // 已读已领取
)

type EmPlatType int32

const (
	EmPlatType_Local EmPlatType = 0 // 本地
)

type EmPlayerOfflineType int32

const (
	EmPlayerOfflineType_EPOT_Mail EmPlayerOfflineType = 0 // 邮件
	EmPlayerOfflineType_EPOT_Item EmPlayerOfflineType = 1 // 道具
)

type EmShopType int32

const (
	EmShopType_EST_None      EmShopType = 0
	EmShopType_EST_BlackShop EmShopType = 1 // 黑市商店
)

type EmSyetemPropType int32

const (
	EmSyetemPropType_SPT_HeroBook EmSyetemPropType = 0 // 图鉴
)

type EmTaskState int32

const (
	EmTaskState_ETS_Ing    EmTaskState = 0 // 进行中
	EmTaskState_ETS_Finish EmTaskState = 1 // 已完成
	EmTaskState_ETS_Award  EmTaskState = 2 // 已领取
)

type LoginState int32

const (
	LoginState_None    LoginState = 0 // 未初始化
	LoginState_Init    LoginState = 1 // 初始化
	LoginState_SetName LoginState = 2 // 取名完成
)

type MAIL int32

const (
	MAIL_M_PlayerMgr  MAIL = 0 // 游戏player
	MAIL_M_AccountMgr MAIL = 1 // 网关player
)

type PlayerDataType int32

const (
	PlayerDataType_Crystal            PlayerDataType = 0  // 晶核系统
	PlayerDataType_Base               PlayerDataType = 1  // 角色基础数据
	PlayerDataType_System             PlayerDataType = 2  // 角色功能数据
	PlayerDataType_Bag                PlayerDataType = 3  // 玩家背包
	PlayerDataType_Equipment          PlayerDataType = 4  // 玩家装备
	PlayerDataType_Client             PlayerDataType = 5  // 角色前端数据
	PlayerDataType_Hero               PlayerDataType = 6  // 伙伴数据
	PlayerDataType_Mail               PlayerDataType = 7  // 邮件数据
	PlayerDataType_Max                PlayerDataType = 8  // 最大值
	PlayerDataType_SystemCommon       PlayerDataType = 10 // 功能通用系统
	PlayerDataType_SystemChat         PlayerDataType = 11 // 功能聊天系统
	PlayerDataType_SystemProfession   PlayerDataType = 12 // 功能职业系统
	PlayerDataType_SystemBox          PlayerDataType = 13 // 宝箱系统
	PlayerDataType_SystemBattle       PlayerDataType = 14 // 功能战斗系统
	PlayerDataType_SystemTask         PlayerDataType = 15 // 功能任务系统
	PlayerDataType_SystemShop         PlayerDataType = 16 // 黑市商店
	PlayerDataType_SystemDraw         PlayerDataType = 17 // 抽奖系统
	PlayerDataType_SystemCharge       PlayerDataType = 18 // 充值系统
	PlayerDataType_SystemGene         PlayerDataType = 19 // 基因系统
	PlayerDataType_SystemOffline      PlayerDataType = 20 // 离线系统
	PlayerDataType_SystemHookTech     PlayerDataType = 21 // 挂机科技
	PlayerDataType_SystemSevenDay     PlayerDataType = 22 // 七天系统
	PlayerDataType_SystemWorldBoss    PlayerDataType = 23 // 世界buss
	PlayerDataType_SystemChampionship PlayerDataType = 24 // 锦标赛
	PlayerDataType_SystemActivity     PlayerDataType = 25 // 活动
	PlayerDataType_SystemMax          PlayerDataType = 26 // 功能最大
)

type Protocol_Player int32

const (
	Protocol_Player_P_C2S_Protocol_Player         Protocol_Player = 0   // 玩家操作请求	_emC2S_Player_Protocol
	Protocol_Player_P_C2S_Protocol_Common         Protocol_Player = 1   // 通用功能模块 _emC2S_Common_Protocol
	Protocol_Player_P_C2S_Protocol_Copymap        Protocol_Player = 2   // 副本模块 _emC2S_Copymap_Protocol
	Protocol_Player_P_C2S_Protocol_Pet            Protocol_Player = 3   // 伙伴模块 _emC2S_Pet_Protocol
	Protocol_Player_P_C2S_Protocol_Item           Protocol_Player = 4   // 道具模块 _emC2S_Item_Protocol
	Protocol_Player_P_C2S_Protocol_Fight          Protocol_Player = 5   // 战斗模块 _emC2S_Fight_Protocol
	Protocol_Player_P_C2S_Protocol_Task           Protocol_Player = 6   // 任务模块 _emC2S_Task_Protocol
	Protocol_Player_P_C2S_Protocol_Mail           Protocol_Player = 7   // 邮件系统 _emC2S_Mail_Protocol
	Protocol_Player_P_C2S_Protocol_TopList        Protocol_Player = 8   // 排行榜系统 _emC2S_TopList_Protocol
	Protocol_Player_P_C2S_Protocol_Challenge      Protocol_Player = 9   // 竞技场	_emC2S_Challenge_Protocol
	Protocol_Player_P_C2S_Protocol_Faction        Protocol_Player = 10  // 帮派相关操作 _emC2S_Faction_Protocol
	Protocol_Player_P_C2S_Protocol_Team           Protocol_Player = 11  // 帮派相关操作 _emC2S_Team_Protocol
	Protocol_Player_P_C2S_Protocol_Call           Protocol_Player = 12  // 召唤系统 _emC2S_Call_Protocol
	Protocol_Player_P_C2S_Protocol_Sail           Protocol_Player = 13  // 远航系统 _emC2S_Sail_Protocol
	Protocol_Player_P_C2S_Protocol_Hook           Protocol_Player = 14  // 挂机系统 _emC2S_Hook_Protocol
	Protocol_Player_P_C2S_Protocol_Artifact       Protocol_Player = 15  // 神器系统 _emC2S_Artifact_Protocol
	Protocol_Player_P_C2S_Protocol_Shop           Protocol_Player = 16  // 商店系统 _emC2S_Shop_Protocol
	Protocol_Player_P_C2S_Protocol_Train          Protocol_Player = 17  // 试炼系统 _emC2S_Train_Protocol
	Protocol_Player_P_C2S_Protocol_Achieve        Protocol_Player = 18  // 成就系统 _emC2S_Achieve_Protocol
	Protocol_Player_P_C2S_Protocol_Expedition     Protocol_Player = 19  // 远征系统 _emC2S_Expedition_Protocol
	Protocol_Player_P_C2S_Protocol_Shape          Protocol_Player = 20  // 外显系统 _emC2S_Shape_Protocol
	Protocol_Player_P_C2S_Protocol_Temple         Protocol_Player = 21  // 神殿系统 _emC2S_Temple_Protocol
	Protocol_Player_P_C2S_Protocol_Friend         Protocol_Player = 22  // 好友系统 _emC2S_Friend_Protocol
	Protocol_Player_P_C2S_Protocol_Element        Protocol_Player = 23  // 元素系统 _emC2S_Element_Protocol
	Protocol_Player_P_C2S_Protocol_Risk           Protocol_Player = 24  // 冒险系统 _emC2S_Risk_Protocol
	Protocol_Player_P_C2S_Protocol_Dan            Protocol_Player = 25  // 超凡段位系统 _emC2S_Dan_Protocol
	Protocol_Player_P_C2S_Protocol_Ladder         Protocol_Player = 26  // 跨服天梯系统 _emC2S_Ladder_Protocol
	Protocol_Player_P_C2S_Protocol_Champion       Protocol_Player = 27  // 冠军赛系统 _emC2S_Champion_Protocol
	Protocol_Player_P_C2S_Protocol_Holy           Protocol_Player = 28  // 圣物系统 _emC2S_Holy_Protocol
	Protocol_Player_P_C2S_Protocol_Video          Protocol_Player = 29  // 录像系统 _emC2S_Video_Protocol
	Protocol_Player_P_C2S_Protocol_Privilege      Protocol_Player = 30  // 特权系统 _emC2S_Privilege_Protocol
	Protocol_Player_P_C2S_Protocol_Weal           Protocol_Player = 31  // 福利系统 _emC2S_Weal_Protocol
	Protocol_Player_P_C2S_Protocol_Activity       Protocol_Player = 32  // 活动系统 _emC2S_Activity_Protocol
	Protocol_Player_P_C2S_Protocol_Platform       Protocol_Player = 33  // 平台系统 _emC2S_Platform_Protocol
	Protocol_Player_P_C2S_Protocol_Talk           Protocol_Player = 34  // 聊天系统 _emC2S_Talk_Protocol
	Protocol_Player_P_C2S_Protocol_Treasure       Protocol_Player = 35  // 探宝系统 _emC2S_Treasure_Protocol
	Protocol_Player_P_C2S_Protocol_HeavenDungeon  Protocol_Player = 36  // 天界副本系统 _emC2S_HeavenDungeon_Protocol
	Protocol_Player_P_C2S_Protocol_CrossChallenge Protocol_Player = 37  // 跨服竞技场 _emC2S_CrossChallenge_Protocol
	Protocol_Player_P_C2S_Protocol_Tablet         Protocol_Player = 38  // 晶碑 _emC2S_Tablet_Protocol
	Protocol_Player_P_C2S_Protocol_WeekChampion   Protocol_Player = 39  // 周冠军赛 _emC2S_WeekChampion_Protocol
	Protocol_Player_P_C2S_Protocol_TeamCampaign   Protocol_Player = 40  // 组队征战 _emC2S_TeamCampaign_Protocol
	Protocol_Player_P_C2S_Protocol_Operate        Protocol_Player = 255 // 网络层相关操作 _emC2S_Operate_Protocol
)

type SEND int32

const (
	SEND_POINT      SEND = 0 // 指定集群id
	SEND_BOARD_CAST SEND = 1 // 广播
)

type SERVICE int32

const (
	SERVICE_NONE   SERVICE = 0
	SERVICE_CLIENT SERVICE = 1 // 客户端
	SERVICE_GATE   SERVICE = 2 // 网关,转发服务
	SERVICE_GM     SERVICE = 3 // gamemgr
	SERVICE_GAME   SERVICE = 4 // game
	SERVICE_DB     SERVICE = 5 // db
	SERVICE_Dip    SERVICE = 6 // 后台服务
	SERVICE_Record SERVICE = 7 // 日志服务
)

type STUB int32

const (
	STUB_Master          STUB = 0 // master
	STUB_DbPlayerMgr     STUB = 1 // db玩家数据
	STUB_PlayerMgr       STUB = 2 // game玩家数据
	STUB_ChatChannelMgr  STUB = 4 // gm聊天数据
	STUB_DbChatMgr       STUB = 5 // gm聊天数据
	STUB_AccountMgr      STUB = 6 // 登录
	STUB_BattleRecordMgr STUB = 7 // 战斗记录
	STUB_RankMgr         STUB = 8 // 排行榜
	STUB_GlobalMgr       STUB = 9 // 全局数据
	STUB_END             STUB = 10
)
