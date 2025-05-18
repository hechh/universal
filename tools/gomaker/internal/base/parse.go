package base

import (
	"go/ast"

	"forevernine.com/planet/server/tool/gomaker/domain"
	"github.com/spf13/cast"
)

func ParseAstStruct(item *ast.TypeSpec) *domain.AstStruct {
	val, ok := item.Type.(*ast.StructType)
	if !ok || val == nil {
		return nil
	}

	ret := &domain.AstStruct{
		Type: &domain.AstType{
			Name:  item.Name.Name,
			Token: domain.STRUCT,
		},
	}
	parseStructType(ret, val)
	return ret
}

func parseStructType(ret *domain.AstStruct, node *ast.StructType) {
	for _, item := range node.Fields.List {
		switch v := item.Type.(type) {
		case *ast.Ident, *ast.StarExpr:
			ret.Idents = append(ret.Idents, &domain.AstField{
				Type: getAstType(v),
				Name: item.Names[0].Name,
			})
		case *ast.ArrayType:
			ret.Arrays = append(ret.Arrays, &domain.AstField{
				Type: getAstType(v),
				Name: item.Names[0].Name,
			})
		case *ast.MapType:
			ret.Maps = append(ret.Maps, &domain.AstMapField{
				Name:  item.Names[0].Name,
				KType: getAstType(v.Key),
				VType: getAstType(v.Value),
			})
		}
	}
}

func ParseAstEnum(nn *ast.GenDecl) *domain.AstEnum {
	ret := &domain.AstEnum{Values: make(map[string]*domain.AstValue)}
	for _, spec := range nn.Specs {
		item, ok := spec.(*ast.ValueSpec)
		if !ok || item == nil {
			continue
		}

		if vval, ok := item.Values[0].(*ast.BasicLit); ok {
			astVal := &domain.AstValue{
				Type:   getAstType(item.Type),
				Name:   item.Names[0].Name,
				Value:  cast.ToInt32(vval.Value),
				StrVal: vval.Value,
			}
			ret.Values[astVal.Name] = astVal
			ret.Type = astVal.Type
		}
	}
	if len(ret.Values) <= 0 {
		return nil
	}
	return ret
}

func getAstType(node ast.Node) *domain.AstType {
	item := &domain.AstType{}
	parseType(node, item)
	return item
}

func parseType(n ast.Node, item *domain.AstType) {
	switch v := n.(type) {
	case *ast.ArrayType:
		item.Token |= domain.ARRAY
		parseType(v.Elt, item)
	case *ast.StarExpr:
		item.Token |= domain.POINTER
		parseType(v.X, item)
	case *ast.Ident:
		item.Name = v.Name
	}
}
