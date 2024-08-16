package handler

import (
	"hash/crc32"
	"reflect"
	"strings"
	"universal/common/pb"

	"github.com/golang/protobuf/proto"
)

var (
	apis = make(map[uint32]*ApiInfo)
)

type Handler func(ctx *Context, req proto.Message, rsp proto.Message) error

func BuildPacketHead(id uint64, dst pb.SERVICE, arrParam ...uint32) *pb.IPacket {
	code := uint32(0)
	if len(arrParam) > 0 {
		code = arrParam[0]
	}
	return &pb.IPacket{
		Stx:            uint32(0x27),
		Ckx:            uint32(0x72),
		DestServerType: dst,
		Id:             id,
		Code:           code,
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

func Register(f Handler, req proto.Message, rsp proto.Message) {
	reqType := reflect.TypeOf(req).Elem()
	rspType := reflect.TypeOf(rsp).Elem()

	item := &ApiInfo{
		reqname: GetProtoName(req),
		rspname: GetProtoName(rsp),
		req:     reqType,
		rsp:     rspType,
		fun:     f,
	}
	apis[item.GetReqCrc()] = item
	apis[item.GetRspCrc()] = item
}

func GetCrc(name string) uint32 {
	return crc32.ChecksumIEEE([]byte(name))
}

type ApiInfo struct {
	reqname string
	rspname string
	req     reflect.Type
	rsp     reflect.Type
	fun     Handler
}

func (d *ApiInfo) GetReqCrc() uint32 {
	return GetCrc(d.reqname)
}

func (d *ApiInfo) GetRspCrc() uint32 {
	return GetCrc(d.rspname)
}

func (d *ApiInfo) NewRequest() proto.Message {
	return reflect.New(d.req).Interface().(proto.Message)
}

func (d *ApiInfo) NewResponse() proto.Message {
	return reflect.New(d.rsp).Interface().(proto.Message)
}
