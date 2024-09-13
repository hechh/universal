package manager

import (
	"universal/tools/gomaker/domain"

	"github.com/spf13/cast"
)

type ConvInfo struct {
	TableType string
	ProtoType string
	conv      domain.ConvFunc
}

var (
	convs = make(map[string]*ConvInfo)
)

func AddConv(typ, proto string, f domain.ConvFunc) {
	val, ok := convs[typ]
	if !ok {
		convs[typ] = &ConvInfo{TableType: typ, ProtoType: proto, conv: f}
		return
	}
	if len(proto) > 0 {
		val.ProtoType = proto
	}
	if f != nil {
		val.conv = f
	}
}

func IsValidType(typ string) bool {
	val, ok := convs[typ]
	return ok && len(val.ProtoType) > 0
}

func GetProtoType(typ string) string {
	return convs[typ].ProtoType
}

func DefaultEnumConv(str string) interface{} {
	if val, ok := evals[str]; ok {
		return val.Value
	}
	return cast.ToInt32(str)
}
