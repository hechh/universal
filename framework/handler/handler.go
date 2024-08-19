package handler

import (
	"hash/crc32"
	"reflect"
	"strings"
	"universal/common/pb"
	"universal/framework/util"

	"github.com/golang/protobuf/proto"
)

var (
	apis = make(map[uint32]*ApiInfo)
)

type Handler func(ctx *Context, req proto.Message, rsp proto.Message) error

type IHead interface {
	GetPacketHead() *pb.IPacket
}

type ApiInfo struct {
	reqname string
	rspname string
	req     reflect.Type
	rsp     reflect.Type
	fun     Handler
}

func (d *ApiInfo) GetReqName() string {
	return d.reqname
}

func (d *ApiInfo) GetRspName() string {
	return d.rspname
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

func Walk(f func(api *ApiInfo) bool) {
	for _, api := range apis {
		if !f(api) {
			break
		}
	}
}

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

func Get(crc uint32) *ApiInfo {
	return apis[crc]
}

func GetByName(name string) *ApiInfo {
	return apis[GetCrc(name)]
}

func GetCrc(name string) uint32 {
	return crc32.ChecksumIEEE([]byte(name))
}

func Encode(packet proto.Message) []byte {
	crc := GetCrc(GetProtoName(packet))
	buff, _ := proto.Marshal(packet)
	data := append(util.IntToBytes(int(crc)), buff...)
	return data
}

func Decode(buff []byte) (uint32, []byte) {
	packetId := uint32(util.BytesToInt(buff[0:4]))
	return packetId, buff[4:]
}
