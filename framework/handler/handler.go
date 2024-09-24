package handler

import (
	"reflect"
	"strings"
	"universal/common/pb"

	"github.com/golang/protobuf/proto"
)

var (
	apis = make(map[uint32]*ApiInfo)
)

type IRspHead interface {
	GetRspHead() *pb.RspHead
}

type Handler func(ctx *Context, req proto.Message, rsp proto.Message) error

type ApiInfo struct {
	apiID   uint32
	reqname string
	rspname string
	req     reflect.Type
	rsp     reflect.Type
	fun     Handler
}

func RegisterApi(apiID uint32, f Handler, req proto.Message, rsp proto.Message) {
	reqType := reflect.TypeOf(req).Elem()
	rspType := reflect.TypeOf(rsp).Elem()
	apis[apiID] = &ApiInfo{
		apiID:   apiID,
		reqname: GetProtoName(req),
		rspname: GetProtoName(rsp),
		req:     reqType,
		rsp:     rspType,
		fun:     f,
	}
}

func Get(apiID uint32) *ApiInfo {
	return apis[apiID]
}

func Walk(f func(api *ApiInfo) bool) {
	for _, api := range apis {
		if !f(api) {
			break
		}
	}
}

func BuildRspHead(id uint64, dst pb.SERVER, arrParam ...uint32) *pb.RspHead {
	code := uint32(0)
	if len(arrParam) > 0 {
		code = arrParam[0]
	}
	return &pb.RspHead{
		Stx:           uint32(0x27),
		Ckx:           uint32(0x72),
		DstServerType: dst,
		UID:           id,
		Code:          code,
	}
}

func GetProtoName(val proto.Message) string {
	sType := proto.MessageName(val)
	index := strings.Index(sType, ".")
	if index != -1 {
		sType = sType[index+1:]
	}
	return sType
}

func (d *ApiInfo) GetReqName() string {
	return d.reqname
}

func (d *ApiInfo) GetRspName() string {
	return d.rspname
}

func (d *ApiInfo) GetApiID() uint32 {
	return d.apiID
}

func (d *ApiInfo) NewRequest() proto.Message {
	return reflect.New(d.req).Interface().(proto.Message)
}

func (d *ApiInfo) NewResponse() proto.Message {
	return reflect.New(d.rsp).Interface().(proto.Message)
}
