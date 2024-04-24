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

// ----gomaker生成模板-------
{{$name := .Name}}

type {{$name}} struct {
	PlayerFun
}

// --------------------通用接口实现------------------------------
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

// --------------------交互接口实现------------------------------
{{range $req := .ReqList}} {{$rsp := $.Join ($.TrimSuffix $req "Request") "Response"}}
func (this *{{$name}}) {{$req}}(head *pb.RpcHead, request, response proto.Message) error {
	req := request.(*pb.{{$req}})
	rsp := response.(*pb.{{$rsp}})

	return nil
}
{{end}}

