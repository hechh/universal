package manager

import (
	"strings"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/typespec"

	"github.com/spf13/cast"
)

var (
	rules = make(map[string]*RuleInfo)
	jsons = make(map[string][]map[string]interface{})
)

type CastFunc func(map[string]*typespec.Value, string) interface{}
type ParseFunc func(*typespec.FieldNode) *typespec.Field

type RuleInfo struct {
	castFunc  CastFunc
	parseFunc ParseFunc
}

func AddJson(name string, data map[string]interface{}) {
	if _, ok := jsons[name]; !ok {
		jsons[name] = make([]map[string]interface{}, 0)
	}
	jsons[name] = append(jsons[name], data)
}

func ParseRule(field *typespec.FieldNode) *typespec.Field {
	return rules[field.Type].parseFunc(field)
}

func CastRule(field *typespec.FieldNode, alls map[string]*typespec.Value, val string) interface{} {
	return rules[field.Type].castFunc(alls, val)
}

func init() {
	rules["b"] = &RuleInfo{castB, parseB}
	rules["s"] = &RuleInfo{castS, parseS}
	rules["i"] = &RuleInfo{castI, parseI}
	rules["il"] = &RuleInfo{castIL, parseIL}
	rules["ill"] = &RuleInfo{castILL, parseILL}
	rules["f"] = &RuleInfo{castF, parseF}
	rules["fl"] = &RuleInfo{castFL, parseFL}
	rules["fll"] = &RuleInfo{castFLL, parseFLL}
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

// il
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

// ill
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
	for _, sval := range strings.Split(strings.ReplaceAll(val, "|", ","), "#") {
		tmps := []interface{}{}
		for _, ss := range strings.Split(sval, ",") {
			tmps = append(tmps, castI(alls, ss))
		}
		rets = append(rets, tmps)
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
	for _, sval := range strings.Split(strings.ReplaceAll(val, "|", ","), "#") {
		tmps := []interface{}{}
		for _, ss := range strings.Split(sval, ",") {
			tmps = append(tmps, castF(alls, ss))
		}
		rets = append(rets, tmps)
	}
	return rets
}
