import (
	"corps/base"
	"corps/base/cfgEnum"
	"corps/common/cfgData"
	"corps/pb"
	"corps/report"
	"corps/rpc"
	"corps/server/packet"
	"corps/server/serverCommon"

	"github.com/golang/protobuf/proto"
)

{{$name := .}}

type {{$name}} struct {
	PlayerFun
}

//初始化
func (this *{{$name}}) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
}

// 客户端数据
func (this *{{$name}}) SaveDataToClient(pbData *pb.PBPlayerData) {
}

// 设置玩家数据, web管理后台
func (this *{{$name}}) SetUserTypeInfo(message proto.Message) bool {
    return false
}

// 加载数据(非system类型数据)
func (this *{{$name}}) Load(pData []byte){

}

// 存储数据(非system类型数据)
func (this *{{$name}}) Save(bNow bool){
}

//新系统
func (this *{{$name}}) NewPlayer() {
}

// 加载系统数据(system类型数据)
func (this *{{$name}}) LoadSystem(pbSystem *pb.PBPlayerSystem) {
}

// 存储数据 返回存储标志(system类型数据)
func (this *{{$name}}) SaveSystem(pbSystem *pb.PBPlayerSystem) bool {
    return false
}

// ----------以上被playerMgr调用----------
//初始化
func (this *{{$name}}) UpdateCommon(common *FunCommon) {
}

//是否存储数据
func (this *{{$name}}) IsSave() bool {
   return false 
}

//加载完成
func (this *{{$name}}) LoadComplete() {
}

//加载系统数据DB数据 数据初始化用
func (this *{{$name}}) LoadPlayerDBFinish()  {
}

//心跳包
func (this *{{$name}}) Heat()  {
}

//是否跨天
func (this *{{$name}}) PassDay(isDay, isWeek, isMonth bool) {
}

//保存
func (this *{{$name}}) UpdateSave(bSave bool) {
}
