package playerFun

import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common/cfgData"
	"corps/common/serverCommon"
	"corps/framework/common/uerror"
	"corps/pb"

	"github.com/golang/protobuf/proto"
)

const (
	DefaultSchemeCount = 3
)

// ----gomaker生成的模板-------
type PlayerSystemGeneFun struct {
	PlayerFun
	entity *PlayerSystemGeneSchemeEntity
}

func (this *PlayerSystemGeneFun) GetEntity() *PlayerSystemGeneSchemeEntity {
	return this.entity
}

// --------------------通用接口实现------------------------------
// 初始化
func (this *PlayerSystemGeneFun) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
}

// 新系统
func (this *PlayerSystemGeneFun) NewPlayer() {
	// 获取有多少方案
	count := cfgData.GetCfgGlobalConfig(cfgEnum.GlobalConfig_GLOBAL_GENE_SCHEME_COUNT)
	if count <= 0 {
		count = DefaultSchemeCount
	}
	// 初始化
	items := &pb.PBPlayerSystemGene{SchemeID: 1}
	for i := uint32(1); i <= count; i++ {
		items.Schemes = append(items.Schemes, &pb.PBGeneScheme{SchemeID: i})
	}
	this.entity = NewPlayerSystemGeneSchemeEntity(items)
	// 保存数据
	this.UpdateSave(true)
}
func (this *PlayerSystemGeneFun) LoadPlayerDBFinish() {
	if this.entity == nil {
		this.NewPlayer()
		return
	}
}

// 加载系统数据(system类型数据)
func (this *PlayerSystemGeneFun) LoadSystem(pbSystem *pb.PBPlayerSystem) {
	if pbSystem.Gene == nil {
		this.NewPlayer()
		return
	}
	// 加载数据
	this.entity = NewPlayerSystemGeneSchemeEntity(pbSystem.Gene)
	this.UpdateSave(false)
}

// 存储数据 返回存储标志(system类型数据)
func (this *PlayerSystemGeneFun) SaveSystem(pbSystem *pb.PBPlayerSystem) bool {
	pbSystem.Gene = this.entity.ToProto()
	return true
}

// 客户端数据
func (this *PlayerSystemGeneFun) SaveDataToClient(pbData *pb.PBPlayerData) {
	if pbData.System == nil {
		pbData.System = &pb.PBPlayerSystem{}
	}
	pbData.System.Gene = this.entity.ToProto()
}
func (this *PlayerSystemGeneFun) GetProtoPtr() proto.Message {
	return &pb.PBPlayerSystemGene{}
}

// 设置玩家数据, web管理后台
func (this *PlayerSystemGeneFun) SetUserTypeInfo(pbData proto.Message) bool {
	if pbData == nil {
		return false
	}
	pbSystem, ok := pbData.(*pb.PBPlayerSystemGene)
	if !ok || pbSystem == nil {
		return false
	}
	// 加载数据
	this.entity = NewPlayerSystemGeneSchemeEntity(pbSystem)
	return true
}

// --------------------交互接口实现------------------------------
// 切换基因方案
func (this *PlayerSystemGeneFun) GeneSchemeChangeRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.GeneSchemeChangeRequest)
	scheme := this.entity.GetScheme(req.SchemeID)
	if scheme == nil {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_GeneSchemeNotExist, "head: %v, req: %v", head, req)
	}
	this.entity.ChangeScheme(req.SchemeID)
	this.UpdateSave(true)
	return nil
}

// 重置整个基因方案
func (this *PlayerSystemGeneFun) GeneSchemeResetRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.GeneSchemeResetRequest)
	//rsp := response.(*pb.GeneSchemeResetResponse)
	scheme := this.entity.GetScheme(req.SchemeID)
	if scheme == nil {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_GeneSchemeNotExist, "head: %v, req: %v", head, req)
	}
	// 去掉所有卡牌和机器人的激活状态
	scheme.Reset()
	// 保存数据
	this.UpdateSave(true)
	return nil
}

// 修改基因方案
func (this *PlayerSystemGeneFun) GeneChangeNameRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.GeneChangeNameRequest)
	//rsp := response.(*pb.GeneChangeNameResponse)
	scheme := this.entity.GetScheme(req.SchemeID)
	if scheme == nil {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_GeneSchemeNotExist, "head: %v, req: %v", head, req)
	}

	//检查屏蔽字库
	uCode := serverCommon.CheckMaskWord(&serverCommon.CheckMaskWordRequest{
		Content:    req.Name,
		SenderID:   int(this.AccountId),
		SenderName: this.getPlayerBaseFun().GetDisplay().PlayerName,
		SendTime:   int(base.GetNow())})

	if uCode != cfgEnum.ErrorCode_Success {
		return uerror.NewUErrorf(1, uCode, "head: %v, req: %v", head, req)
	}

	// 设置名字
	scheme.SetName(req.Name)
	// 保存数据
	this.UpdateSave(true)
	return nil
}

// 激活触发器卡牌
func (this *PlayerSystemGeneFun) GeneCardActiveRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.GeneCardActiveRequest)
	//rsp := response.(*pb.GeneCardActiveResponse)
	scheme := this.entity.GetScheme(req.SchemeID)
	if scheme == nil {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_GeneSchemeNotExist, "head: %v, req: %v", head, req)
	}
	if err := scheme.ActiveTagCard(req.Actives); err != nil {
		return err
	}
	scheme.ResetRobot(req.RobotPositions)

	//成就统计
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_GeneActive, 1)
	// 保存数据
	this.UpdateSave(true)
	return nil
}

// 激活机器人
func (this *PlayerSystemGeneFun) GeneRobotActiveRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.GeneRobotActiveRequest)
	//rsp := response.(*pb.GeneRobotActiveResponse)
	scheme := this.entity.GetScheme(req.SchemeID)
	if scheme == nil {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_GeneSchemeNotExist, "head: %v, req: %v", head, req)
	}

	if err := scheme.ActiveRobot(req.Position, req.RobotID); err != nil {
		return err
	}

	//成就统计
	this.getPlayerSystemTaskFun().TriggerAchieve(head, cfgEnum.AchieveType_GeneActive, 1)

	this.UpdateSave(true)
	return nil
}

// 激活机器人卡牌
func (this *PlayerSystemGeneFun) GeneRobotCardActiveRequest(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.GeneRobotCardActiveRequest)
	//rsp := response.(*pb.GeneRobotCardActiveResponse)
	scheme := this.entity.GetScheme(req.SchemeID)
	if scheme == nil {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_GeneSchemeNotExist, "head: %v, req: %v", head, req)
	}
	if err := scheme.ActiveRobotTagCard(req.Position, req.Actives); err != nil {
		return err
	}
	this.UpdateSave(true)
	return nil
}

// ------------------------------------------Entity封装---------------------------------------------------
type PlayerSystemGeneSchemeEntity struct {
	schemeID uint32
	schemes  map[uint32]*GeneSchemeEntity
}

func NewPlayerSystemGeneSchemeEntity(data *pb.PBPlayerSystemGene) *PlayerSystemGeneSchemeEntity {
	ret := &PlayerSystemGeneSchemeEntity{
		schemeID: data.SchemeID,
		schemes:  make(map[uint32]*GeneSchemeEntity),
	}
	for _, item := range data.Schemes {
		ret.schemes[item.SchemeID] = NewGeneSchemeEntity(item)
	}
	return ret
}

func (d *PlayerSystemGeneSchemeEntity) GetCount() int {
	return len(d.schemes)
}

func (d *PlayerSystemGeneSchemeEntity) GetScheme(id uint32) *GeneSchemeEntity {
	return d.schemes[id]
}

func (d *PlayerSystemGeneSchemeEntity) ChangeScheme(id uint32) {
	d.schemeID = id
}

func (d *PlayerSystemGeneSchemeEntity) ToProto() *pb.PBPlayerSystemGene {
	item := &pb.PBPlayerSystemGene{SchemeID: d.schemeID}
	for _, scheme := range d.schemes {
		item.Schemes = append(item.Schemes, scheme.ToProto())
	}
	return item
}

// -------------------------具体方案---------------------------
type GeneSchemeEntity struct {
	schemeID uint32
	name     string
	tags     GeneTagMap
	robots   map[uint32]*GeneRobotEntity
}

func NewGeneSchemeEntity(data *pb.PBGeneScheme) *GeneSchemeEntity {
	ret := &GeneSchemeEntity{
		schemeID: data.SchemeID,
		name:     data.Name,
		tags:     make(map[uint32]*GeneTagEntity),
		robots:   make(map[uint32]*GeneRobotEntity),
	}
	for _, item := range data.Tags {
		ret.tags[item.TagID] = NewGeneTagEntity(item)
	}
	for _, item := range data.Robots {
		ret.robots[item.Position] = NewGeneRobotEntity(item)
	}
	return ret
}

func (d *GeneSchemeEntity) ToProto() *pb.PBGeneScheme {
	return &pb.PBGeneScheme{
		SchemeID: d.schemeID,
		Name:     d.name,
		Tags:     d.tags.ToProto(),
		Robots: func() (rets []*pb.PBGeneRobot) {
			for _, item := range d.robots {
				rets = append(rets, item.ToProto())
			}
			return
		}(),
	}
}

func (d *GeneSchemeEntity) GetName() string {
	return d.name
}

func (d *GeneSchemeEntity) GetTags() map[uint32]*GeneTagEntity {
	return d.tags
}

func (d *GeneSchemeEntity) GetRobots() map[uint32]*GeneRobotEntity {
	return d.robots
}

// 设置基因方案名字
func (d *GeneSchemeEntity) SetName(name string) {
	d.name = name
}

// 重置基因方案
func (d *GeneSchemeEntity) Reset() {
	d.tags = make(map[uint32]*GeneTagEntity)
	d.robots = make(map[uint32]*GeneRobotEntity)
}

// 激活卡牌
func (d *GeneSchemeEntity) ActiveTagCard(list []*pb.GeneCardActiveInfo) error {
	return d.tags.ActiveCard(list)
}

// 激活机器人卡牌
func (d *GeneSchemeEntity) ActiveRobotTagCard(pos uint32, list []*pb.GeneCardActiveInfo) error {
	robot, ok := d.robots[pos]
	if !ok {
		return uerror.NewUErrorf(1, cfgEnum.ErrorCode_GeneRobotNotActive, "robot not active in position(%d)", pos)
	}
	return robot.tags.ActiveCard(list)
}

// 激活机器人
func (d *GeneSchemeEntity) ActiveRobot(position uint32, robotID uint32) error {
	if robot, ok := d.robots[position]; !ok {
		d.robots[position] = NewGeneRobotEntity(&pb.PBGeneRobot{RobotID: robotID, Position: position})
	} else {
		// 重复激活
		if robot.robotID == robotID {
			return uerror.NewUErrorf(1, cfgEnum.ErrorCode_GeneRobotRepeatActive, "Robot(%d,%d) repeat active", robotID, position)
		}
		d.robots[position] = NewGeneRobotEntity(&pb.PBGeneRobot{RobotID: robotID, Position: position})
	}

	return nil
}

// 重置机器人
func (d *GeneSchemeEntity) ResetRobot(positions []uint32) {
	for _, pos := range positions {
		delete(d.robots, pos)
	}
}

// ---------------------------------基因使徒--------------------------------
type GeneTagMap map[uint32]*GeneTagEntity

type GeneRobotEntity struct {
	robotID  uint32
	position uint32
	tags     GeneTagMap
}

func NewGeneRobotEntity(data *pb.PBGeneRobot) *GeneRobotEntity {
	ret := &GeneRobotEntity{
		robotID:  data.RobotID,
		position: data.Position,
		tags:     make(map[uint32]*GeneTagEntity),
	}
	for _, item := range data.Tags {
		ret.tags[item.TagID] = NewGeneTagEntity(item)
	}
	return ret
}

func (d *GeneRobotEntity) GetTags() map[uint32]*GeneTagEntity {
	return d.tags
}

func (d *GeneRobotEntity) GetRobotID() uint32 {
	return d.robotID
}

func (d *GeneRobotEntity) ToProto() *pb.PBGeneRobot {
	return &pb.PBGeneRobot{
		RobotID:  d.robotID,
		Position: d.position,
		Tags:     GeneTagMap(d.tags).ToProto(),
	}
}

func (tags GeneTagMap) ToProto() (rets []*pb.PBGeneTag) {
	for _, item := range tags {
		rets = append(rets, item.ToProto())
	}
	return
}

// ---------------------------------基因触发器--------------------------------
type GeneTagEntity struct {
	id    uint32
	cards map[uint32]struct{}
}

func NewGeneTagEntity(data *pb.PBGeneTag) *GeneTagEntity {
	ret := &GeneTagEntity{
		id:    data.TagID,
		cards: make(map[uint32]struct{}),
	}
	for _, id := range data.Cards {
		ret.cards[id] = struct{}{}
	}
	return ret
}

func (d *GeneTagEntity) GetCards() map[uint32]struct{} {
	return d.cards
}

func (d *GeneTagEntity) SetActive(active bool, cardID uint32) {
	if active {
		d.cards[cardID] = struct{}{}
	} else {
		delete(d.cards, cardID)
	}
}

// 激活标签模块中的卡牌
func (tags GeneTagMap) ActiveCard(list []*pb.GeneCardActiveInfo) error {
	for _, item := range list {
		rogueCfg := cfgData.GetCfgRoguelikeCard(item.CardID)
		if rogueCfg == nil {
			return uerror.NewUErrorf(1, cfgData.GetRoguelikeCardErrorCode(item.CardID), "RogueLikeCard(%d) not found", item.CardID)
		}
		// 判断该标签是否存在
		tag, ok := tags[rogueCfg.Prof]
		if !ok {
			tag = NewGeneTagEntity(&pb.PBGeneTag{TagID: rogueCfg.Prof})
			tags[rogueCfg.Prof] = tag
		}
		// 修改激活状态
		tag.SetActive(item.IsActive, item.CardID)
	}
	return nil
}

func (d *GeneTagEntity) ToProto() *pb.PBGeneTag {
	return &pb.PBGeneTag{
		TagID: d.id,
		Cards: func() (rets []uint32) {
			for cardid := range d.cards {
				rets = append(rets, cardid)
			}
			return
		}(),
	}
}
