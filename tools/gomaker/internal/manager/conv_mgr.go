package manager

import (
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/typespec"

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

func AddConvType(typ, proto string) {
	if val, ok := convs[typ]; ok {
		val.ProtoType = proto
	} else {
		convs[typ] = &ConvInfo{TableType: typ, ProtoType: proto}
	}
}

func AddConvFunc(typ string, f domain.ConvFunc) {
	if val, ok := convs[typ]; ok {
		val.conv = f
	} else {
		convs[typ] = &ConvInfo{TableType: typ, conv: f}
	}
}

func IsValidType(typ string) bool {
	val, ok := convs[typ]
	return ok && len(val.ProtoType) > 0
}

func GetProtoType(typ string) string {
	return convs[typ].ProtoType
}

func AddEnumValue(item *typespec.Value) {
	evals[item.Doc] = item
}

func ToEValue(str string) int32 {
	if val, ok := evals[str]; ok {
		return val.Value
	}
	return cast.ToInt32(str)
}

func DefaultEnumConv(str string) interface{} {
	return cast.ToInt32(str)
}
