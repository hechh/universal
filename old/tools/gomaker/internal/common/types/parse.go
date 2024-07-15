package types

import (
	"go/ast"
	"sort"
	"strings"
	"unicode"
	"universal/tools/gomaker/domain"

	"github.com/spf13/cast"
)

func ParseType(val ast.Expr, typ *Type) {
	switch v := val.(type) {
	case *ast.MapType:
		typ.Token |= domain.MAP
		typ.Key, typ.Value = &Type{}, &Type{}
		ParseType(v.Key, typ.Key)
		ParseType(v.Value, typ.Value)
	case *ast.StarExpr:
		typ.Token |= domain.POINTER
		ParseType(v.X, typ)
	case *ast.ArrayType:
		typ.Token |= domain.ARRAY
		ParseType(v.Elt, typ)
	case *ast.SelectorExpr:
		typ.PkgName, typ.Name = v.X.(*ast.Ident).Name, v.Sel.Name
	case *ast.Ident:
		typ.Name = v.Name
	}
}

func ParseComment(g *ast.CommentGroup) string {
	if g == nil || len(g.List) <= 0 {
		return ""
	}
	return g.List[0].Text
}

func ParseTag(tag *ast.BasicLit) string {
	if tag != nil {
		return tag.Value
	}
	return ""
}

func ParseAlias(pkg, name, doc string, node *ast.TypeSpec) *Alias {
	subType := &Type{}
	ParseType(node.Type, subType)

	return &Alias{
		Type:      &Type{PkgName: pkg, Name: name, Token: domain.ALIAS},
		Doc:       doc,
		Comment:   ParseComment(node.Comment),
		Reference: subType,
	}
}

func ParseEnum(pkg, doc string, specs ...ast.Spec) *Enum {
	item := &Enum{
		Type:   &Type{PkgName: pkg, Token: domain.ENUM},
		Doc:    doc,
		Fields: make(map[string]*Value),
	}
	for _, node := range specs {
		vv, ok := node.(*ast.ValueSpec)
		if !ok || vv == nil || vv.Type == nil {
			continue
		}
		// 解析枚举类型
		item.Type.Name = vv.Type.(*ast.Ident).Name
		// 保存字段
		val := &Value{
			Name:    vv.Names[0].Name,
			Comment: ParseComment(vv.Comment),
			Value:   cast.ToInt32(vv.Values[0].(*ast.BasicLit).Value),
		}
		item.Fields[vv.Names[0].Name] = val
		item.List = append(item.List, val)
	}
	sort.Slice(item.List, func(i, j int) bool {
		return item.List[i].Value < item.List[j].Value
	})
	return item

}

func ParseStruct(pkg, name, doc string, Fields []*ast.Field) *Struct {
	item := &Struct{
		Type:   &Type{PkgName: pkg, Name: name, Token: domain.STRUCT},
		Doc:    doc,
		Fields: make(map[string]*Field),
		List:   make([]*Field, 0),
	}
	for _, field := range Fields {
		// 过滤不对外开放接
		if unicode.IsLower(rune(field.Names[0].Name[0])) {
			continue
		}
		// 解析类型
		fieldType := &Type{}
		ParseType(field.Type, fieldType)
		// 存储
		item.Fields[field.Names[0].Name] = &Field{
			Type:    fieldType,
			Name:    field.Names[0].Name,
			Comment: ParseComment(field.Comment),
			Tag:     ParseTag(field.Tag),
		}
		item.List = append(item.List, item.Fields[field.Names[0].Name])
	}
	sort.Slice(item.List, func(i, j int) bool {
		return strings.Compare(item.List[i].Name, item.List[j].Name) < 0
	})
	return item
}
