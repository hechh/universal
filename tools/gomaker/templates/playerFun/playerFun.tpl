import (
	"context"
	"corps/base"
	"corps/common/uerror"
	"corps/pb"
	"corps/server/packet"

	"github.com/golang/protobuf/proto"
)

{{$name := .Name}}

// ----gomaker生成的模板-------
type {{$name}} struct {
	PlayerFun
}

// --------------------通用接口实现------------------------------
//初始化
func (this *{{$name}}) Init(pbType pb.PlayerDataType, common *FunCommon) {
	this.PlayerFun.Init(pbType, common)
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

// 客户端数据
func (this *{{$name}}) SaveDataToClient(pbData *pb.PBPlayerData) {
}

// 设置玩家数据, web管理后台
func (this *{{$name}}) SetUserTypeInfo(message proto.Message) bool {
    return false
}

// --------------------交互接口实现------------------------------
{{range $req := .ReqList}} {{$rsp := $.Join ($.TrimSuffix $req "Request") "Response"}}
func (this *{{$name}}) {{$req}}(head *pb.RpcHead, request, response proto.Message) error {
	//req := request.(*pb.{{$req}})
	//rsp := response.(*pb.{{$rsp}})
	// to
	return nil
}
{{end}}

/*
{{range $req := .ReqList}} {{$rsp := $.Join ($.TrimSuffix $req "Request") "Response"}}
func (this *Player) {{$req}}(ctx context.Context, req *pb.{{$req}}) {
	head := this.GetRpcHead(ctx)
	rsp := &pb.{{$rsp}}{PacketHead: &pb.IPacket{}}
	err := this.get{{$name}}().{{$req}}(head, req, rsp)
	if err != nil {
		base.Debugf(uerror.GetCode(err), "error: %v", err)
	}
	packet.SendToClient(head, rsp, uerror.GetCode(err))
}
{{end}}
*/


/*
{{range $req := .ReqList}} {{$rsp := $.Join ($.TrimSuffix $req "Request") "Response"}}
RegisterPacket(&pb.{{$req}}{}, "game{{html "<"}}-Player.{{$req}}", &pb.{{$rsp}}{}, WithCmd(true))
{{end}}
*/