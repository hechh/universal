/*
* 本代码由pbtool工具生成，请勿手动修改
 */

package pb

/*
import (
	"github.com/golang/protobuf/proto"
)
*/

func (d *LoginReq) SetToken(v string) {
	d.Token = v
}

func (d *LoginRsp) SetHead(v *RspHead) {
	d.Head = v
}

func (d *LogoutReq) SetUid(v uint64) {
	d.Uid = v
}

func (d *LogoutRsp) SetHead(v *RspHead) {
	d.Head = v
}

func (d *HeartReq) SetUtc(v int64) {
	d.Utc = v
}
func (d *HeartReq) SetBeginTime(v int64) {
	d.BeginTime = v
}

func (d *HeartRsp) SetHead(v *RspHead) {
	d.Head = v
}
func (d *HeartRsp) SetUtc(v int64) {
	d.Utc = v
}
func (d *HeartRsp) SetBeginTime(v int64) {
	d.BeginTime = v
}
func (d *HeartRsp) SetEndTime(v int64) {
	d.EndTime = v
}

func (d *Node) SetName(v string) {
	d.Name = v
}
func (d *Node) SetType(v NodeType) {
	d.Type = v
}
func (d *Node) SetId(v int32) {
	d.Id = v
}
func (d *Node) SetIp(v string) {
	d.Ip = v
}
func (d *Node) SetPort(v int32) {
	d.Port = v
}

func (d *Router) SetGate(v int32) {
	d.Gate = v
}
func (d *Router) SetGame(v int32) {
	d.Game = v
}
func (d *Router) SetDb(v int32) {
	d.Db = v
}
func (d *Router) SetBuild(v int32) {
	d.Build = v
}
func (d *Router) SetRoom(v int32) {
	d.Room = v
}
func (d *Router) SetMatch(v int32) {
	d.Match = v
}
func (d *Router) SetGm(v int32) {
	d.Gm = v
}

func (d *NodeRouter) SetNodeType(v NodeType) {
	d.NodeType = v
}
func (d *NodeRouter) SetNodeId(v int32) {
	d.NodeId = v
}
func (d *NodeRouter) SetRouterId(v uint64) {
	d.RouterId = v
}
func (d *NodeRouter) SetActorFunc(v uint32) {
	d.ActorFunc = v
}
func (d *NodeRouter) SetActorId(v uint64) {
	d.ActorId = v
}
func (d *NodeRouter) SetRouter(v *Router) {
	d.Router = v
}

func (d *Head) SetSendType(v SendType) {
	d.SendType = v
}
func (d *Head) SetSrc(v *NodeRouter) {
	d.Src = v
}
func (d *Head) SetDst(v *NodeRouter) {
	d.Dst = v
}
func (d *Head) SetUid(v uint64) {
	d.Uid = v
}
func (d *Head) SetSeq(v uint32) {
	d.Seq = v
}
func (d *Head) SetCmd(v uint32) {
	d.Cmd = v
}
func (d *Head) SetReference(v uint32) {
	d.Reference = v
}
func (d *Head) SetReply(v string) {
	d.Reply = v
}
func (d *Head) SetActorName(v string) {
	d.ActorName = v
}
func (d *Head) SetFuncName(v string) {
	d.FuncName = v
}
func (d *Head) SetActorId(v uint64) {
	d.ActorId = v
}

func (d *Packet) SetHead(v *Head) {
	d.Head = v
}
func (d *Packet) SetBody(v []byte) {
	d.Body = v
}

func (d *RspHead) SetCode(v int32) {
	d.Code = v
}
func (d *RspHead) SetErrMsg(v string) {
	d.ErrMsg = v
}

/*
var (
	factorys = make(map[string]func() proto.Message)
)

func init() {
	factorys["LoginReq"] = func() proto.Message { return &LoginReq{} }
	factorys["LoginRsp"] = func() proto.Message { return &LoginRsp{} }
	factorys["LogoutReq"] = func() proto.Message { return &LogoutReq{} }
	factorys["LogoutRsp"] = func() proto.Message { return &LogoutRsp{} }
	factorys["HeartReq"] = func() proto.Message { return &HeartReq{} }
	factorys["HeartRsp"] = func() proto.Message { return &HeartRsp{} }
	factorys["Node"] = func() proto.Message { return &Node{} }
	factorys["Router"] = func() proto.Message { return &Router{} }
	factorys["NodeRouter"] = func() proto.Message { return &NodeRouter{} }
	factorys["Head"] = func() proto.Message { return &Head{} }
	factorys["Packet"] = func() proto.Message { return &Packet{} }
	factorys["RspHead"] = func() proto.Message { return &RspHead{} }
}
*/
