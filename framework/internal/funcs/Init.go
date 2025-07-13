package funcs

import (
	"reflect"
	"sync/atomic"
	"universal/common/pb"
	"universal/library/uerror"
	"universal/library/util"

	"github.com/golang/protobuf/proto"
)

var (
	apis    = make(map[uint32]*Method)
	sendRsp func(*pb.Head, proto.Message) error
	args    = util.ArrayPool[reflect.Value](6)
)

func Init(f func(*pb.Head, proto.Message) error) {
	sendRsp = f
}

func AddReference(head *pb.Head) {
	if rr, ok := apis[head.Dst.ActorFunc]; ok {
		if rr.mask == (HEAD_FLAG | REQ_FLAG | RSP_FLAG) {
			atomic.AddUint32(&head.Reference, 1)
		}
	}
}

func ParseNodeRouter(head *pb.Head, ffs ...string) error {
	var ok bool
	var rr *Method
	if head.Dst.ActorFunc > 0 {
		rr, ok = apis[head.Dst.ActorFunc]
	} else if len(ffs) > 0 {
		rr, ok = apis[util.GetCrc32(ffs[0])]
	}
	if !ok {
		return uerror.New(1, -1, "请求接口不存在%v", head.Dst)
	}
	actorType := head.Dst.ActorId & 0xFF
	if rr.act.GetActorType() != uint32(actorType) {
		return uerror.New(1, -1, "ActorType错误%v", head.Dst)
	}
	head.ActorName = rr.act.GetActorName()
	head.FuncName = rr.Name
	head.ActorId = (head.Dst.ActorId >> 8)
	return nil
}

func register(mm *Method) *Method {
	apis[util.GetCrc32(mm.act.GetActorName()+"."+mm.Name)] = mm
	return mm
}

func get(size int) []reflect.Value {
	rets := args.Get().([]reflect.Value)
	return rets[:size]
}

func put(rets []reflect.Value) {
	args.Put(rets)
}

func toRspHead(err error) *pb.RspHead {
	switch vv := err.(type) {
	case nil:
		return nil
	case *uerror.UError:
		return &pb.RspHead{Code: vv.GetCode(), ErrMsg: vv.GetMsg()}
	case error:
		return &pb.RspHead{Code: -1, ErrMsg: vv.Error()}
	default:
		return nil
	}
}
