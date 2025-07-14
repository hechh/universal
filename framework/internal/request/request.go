package request

import (
	"hash/crc32"
	"reflect"
	"strings"
	"universal/common/pb"
	"universal/library/uerror"
	"universal/library/util"

	"github.com/golang/protobuf/proto"
)

var (
	self         *pb.Node
	apis         = make(map[uint32]IRequest)
	names        = make(map[string]uint32)
	args         = util.ArrayPool[reflect.Value](6)
	SendResponse func(*pb.Head, proto.Message) error
)

type IRequest interface {
	GetActorType() uint32
	GetActorName() string
	GetFuncName() string
}

func Init(nn *pb.Node, f func(*pb.Head, proto.Message) error) {
	self = nn
	SendResponse = f
}

func GetReflectValueArray(size int) []reflect.Value {
	rets := args.Get().([]reflect.Value)
	return rets[:size]
}

func PutReflectValueArray(rets []reflect.Value) {
	args.Put(rets)
}

func Register(mm IRequest) {
	apis[GetCrc32(mm.GetActorName(), mm.GetFuncName())] = mm
}

func GetCrc32(strs ...string) uint32 {
	actorFunc := strings.Join(strs, ".")
	if _, ok := names[actorFunc]; !ok {
		names[actorFunc] = crc32.ChecksumIEEE(util.StringToBytes(actorFunc))
	}
	return names[actorFunc]
}

func NewNodeRouter(actorFunc string, id uint64) *pb.NodeRouter {
	actId := GetCrc32(actorFunc)
	rr, ok := apis[actId]
	if !ok {
		return nil
	}
	return &pb.NodeRouter{
		NodeType:  self.Type,
		NodeId:    self.Id,
		ActorFunc: actId,
		ActorId:   id<<8 | uint64(rr.GetActorType()&0xFF),
	}
}

func Parse(head *pb.Head, ffs ...string) error {
	var ok bool
	var rr IRequest
	if head.Dst.ActorFunc > 0 {
		rr, ok = apis[head.Dst.ActorFunc]
	} else if len(ffs) > 0 {
		rr, ok = apis[GetCrc32(ffs[0])]
	}
	if !ok {
		return uerror.New(1, -1, "请求接口不存在%v", head.Dst)
	}
	actorType := uint32(head.Dst.ActorId & 0xFF)
	if rr.GetActorType() != actorType {
		return uerror.New(1, -1, "ActorType错误%v", head.Dst)
	}
	head.ActorName = rr.GetActorName()
	head.FuncName = rr.GetFuncName()
	head.ActorId = head.Dst.ActorId >> 8
	return nil
}
