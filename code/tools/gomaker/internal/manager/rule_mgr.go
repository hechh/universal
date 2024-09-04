package manager

import (
	"fmt"
	"strings"
	"time"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/typespec"

	"github.com/spf13/cast"
)

var (
	rules = make(map[string]*RuleInfo)
)

type CastFunc func(map[string]*typespec.Value, string) interface{}
type ParseFunc func(*typespec.FieldNode) *typespec.Field

type RuleInfo struct {
	castFunc  CastFunc
	parseFunc ParseFunc
}

func ParseRule(field *typespec.FieldNode) *typespec.Field {
	if rr, ok := rules[field.Type]; ok {
		return rr.parseFunc(field)
	}
	panic(fmt.Errorf("[%s]不支持", field.Type))
}

func CastRule(field *typespec.FieldNode, alls map[string]*typespec.Value, val string) interface{} {
	return rules[field.Type].castFunc(alls, val)
}

func init() {
	rules["t"] = &RuleInfo{castT, parseT}
	rules["b"] = &RuleInfo{castB, parseB}
	rules["s"] = &RuleInfo{castS, parseS}
	rules["i"] = &RuleInfo{castI, parseI}
	rules["il"] = &RuleInfo{castIL, parseIL}
	rules["ill"] = &RuleInfo{castILL, parseILL}
	rules["in"] = &RuleInfo{castIN, parseIN}
	rules["inl"] = &RuleInfo{castINL, parseINL}
	rules["inll"] = &RuleInfo{castINLL, parseINLL}
	rules["f"] = &RuleInfo{castF, parseF}
	rules["fl"] = &RuleInfo{castFL, parseFL}
	rules["fll"] = &RuleInfo{castFLL, parseFLL}
}

// t
func parseT(field *typespec.FieldNode) *typespec.Field {
	return &typespec.Field{
		Type:  GetOrAddType(&typespec.Type{Kind: domain.KIND_IDENT, Name: "uint64"}),
		Name:  field.Name,
		Index: field.Index,
		Doc:   field.Doc,
	}
}

func castT(alls map[string]*typespec.Value, val string) interface{} {
	t, _ := time.Parse("2006-01-02 15:04:05", val)
	return uint64(t.Unix())
}

// b
func parseB(field *typespec.FieldNode) *typespec.Field {
	return &typespec.Field{
		Type:  GetOrAddType(&typespec.Type{Kind: domain.KIND_IDENT, Name: "bool"}),
		Name:  field.Name,
		Index: field.Index,
		Doc:   field.Doc,
	}
}

func castB(alls map[string]*typespec.Value, val string) interface{} {
	return cast.ToBool(val)
}

// s
func parseS(field *typespec.FieldNode) *typespec.Field {
	return &typespec.Field{
		Type:  GetOrAddType(&typespec.Type{Kind: domain.KIND_IDENT, Name: "string"}),
		Name:  field.Name,
		Index: field.Index,
		Doc:   field.Doc,
	}
}

func castS(alls map[string]*typespec.Value, val string) interface{} {
	return val
}

// i
func parseI(field *typespec.FieldNode) *typespec.Field {
	return &typespec.Field{
		Type:  GetOrAddType(&typespec.Type{Kind: domain.KIND_IDENT, Name: "uint32"}),
		Name:  field.Name,
		Index: field.Index,
		Doc:   field.Doc,
	}
}

func castI(alls map[string]*typespec.Value, val string) interface{} {
	if value, ok := alls[val]; ok {
		return uint32(value.Value)
	}
	return cast.ToUint32(val)
}

func parseIL(field *typespec.FieldNode) *typespec.Field {
	return &typespec.Field{
		Type:  GetOrAddType(&typespec.Type{Kind: domain.KIND_IDENT, Name: "uint32"}),
		Name:  field.Name,
		Index: field.Index,
		Doc:   field.Doc,
		Token: []uint32{domain.TOKEN_ARRAY},
	}
}

func castIL(alls map[string]*typespec.Value, val string) interface{} {
	rets := []interface{}{}
	for _, ss := range strings.Split(strings.ReplaceAll(val, "|", ","), ",") {
		rets = append(rets, castI(alls, ss))
	}
	return rets
}

func parseILL(field *typespec.FieldNode) *typespec.Field {
	return &typespec.Field{
		Type:  GetOrAddType(&typespec.Type{Kind: domain.KIND_IDENT, Name: "uint32"}),
		Name:  field.Name,
		Index: field.Index,
		Doc:   field.Doc,
		Token: []uint32{domain.TOKEN_ARRAY, domain.TOKEN_ARRAY},
	}
}

func castILL(alls map[string]*typespec.Value, val string) interface{} {
	rets := []interface{}{}
	for _, sval := range strings.Split(val, "#") {
		rets = append(rets, castIL(alls, sval))
	}
	return rets
}

// in
func parseIN(field *typespec.FieldNode) *typespec.Field {
	return &typespec.Field{
		Type:  GetOrAddType(&typespec.Type{Kind: domain.KIND_IDENT, Name: "int32"}),
		Name:  field.Name,
		Index: field.Index,
		Doc:   field.Doc,
	}
}

func castIN(alls map[string]*typespec.Value, val string) interface{} {
	if value, ok := alls[val]; ok {
		return int32(value.Value)
	}
	return cast.ToInt32(val)
}

func parseINL(field *typespec.FieldNode) *typespec.Field {
	return &typespec.Field{
		Type:  GetOrAddType(&typespec.Type{Kind: domain.KIND_IDENT, Name: "int32"}),
		Name:  field.Name,
		Index: field.Index,
		Doc:   field.Doc,
		Token: []uint32{domain.TOKEN_ARRAY},
	}
}

func castINL(alls map[string]*typespec.Value, val string) interface{} {
	rets := []interface{}{}
	for _, sval := range strings.Split(strings.ReplaceAll(val, "|", ","), ",") {
		rets = append(rets, castIN(alls, sval))
	}
	return rets
}

func parseINLL(field *typespec.FieldNode) *typespec.Field {
	return &typespec.Field{
		Type:  GetOrAddType(&typespec.Type{Kind: domain.KIND_IDENT, Name: "int32"}),
		Name:  field.Name,
		Index: field.Index,
		Doc:   field.Doc,
		Token: []uint32{domain.TOKEN_ARRAY, domain.TOKEN_ARRAY},
	}
}

func castINLL(alls map[string]*typespec.Value, val string) interface{} {
	rets := []interface{}{}
	for _, sval := range strings.Split(val, "#") {
		rets = append(rets, castINL(alls, sval))
	}
	return rets
}

// f
func parseF(field *typespec.FieldNode) *typespec.Field {
	return &typespec.Field{
		Type:  GetOrAddType(&typespec.Type{Kind: domain.KIND_IDENT, Name: "float32"}),
		Name:  field.Name,
		Index: field.Index,
		Doc:   field.Doc,
	}
}

func castF(alls map[string]*typespec.Value, val string) interface{} {
	if value, ok := alls[val]; ok {
		return float32(value.Value)
	}
	return cast.ToFloat32(val)
}

// fl
func parseFL(field *typespec.FieldNode) *typespec.Field {
	return &typespec.Field{
		Type:  GetOrAddType(&typespec.Type{Kind: domain.KIND_IDENT, Name: "float32"}),
		Name:  field.Name,
		Index: field.Index,
		Doc:   field.Doc,
		Token: []uint32{domain.TOKEN_ARRAY},
	}
}

func castFL(alls map[string]*typespec.Value, val string) interface{} {
	rets := []interface{}{}
	for _, ss := range strings.Split(strings.ReplaceAll(val, "|", ","), ",") {
		rets = append(rets, castF(alls, ss))
	}
	return rets
}

// fll
func parseFLL(field *typespec.FieldNode) *typespec.Field {
	return &typespec.Field{
		Type:  GetOrAddType(&typespec.Type{Kind: domain.KIND_IDENT, Name: "float32"}),
		Name:  field.Name,
		Index: field.Index,
		Doc:   field.Doc,
		Token: []uint32{domain.TOKEN_ARRAY, domain.TOKEN_ARRAY},
	}
}

func castFLL(alls map[string]*typespec.Value, val string) interface{} {
	rets := []interface{}{}
	for _, sval := range strings.Split(val, "#") {
		rets = append(rets, castFL(alls, sval))
	}
	return rets
}
