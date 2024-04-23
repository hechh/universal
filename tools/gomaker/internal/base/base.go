package base

import (
	"go/ast"
	"universal/tools/gomaker/domain"
)

func AddToken(d *domain.Type, val int32) {
	d.Token <<= 4
	d.Token |= val & 0x0f
}

func ParseStruct(fields []*ast.Field, st *domain.Struct) {
	for _, field := range fields {
		// 解析字段类型
		typ := &domain.Type{}
		ParseType(field.Type, typ)

		// 保存字段
		st.Fields[field.Names[0].Name] = &domain.Field{
			Name:    field.Names[0].Name,
			Comment: ParseComment(field.Comment),
			Type:    typ,
		}
	}
}

func ParseComment(g *ast.CommentGroup) string {
	if g == nil && len(g.List) <= 0 {
		return ""
	}
	return g.List[0].Text
}

func ParseType(val ast.Expr, typ *domain.Type) {
	switch v := val.(type) {
	case *ast.MapType:
		AddToken(typ, domain.MAP)
		ParseType(v.Key, typ)
		typ.Value = new(domain.Type)
		ParseType(v.Value, typ.Value)
	case *ast.SelectorExpr:
		typ.PkgName = v.X.(*ast.Ident).Name
		typ.Name = v.Sel.Name
	case *ast.StarExpr:
		AddToken(typ, domain.POINTER)
		ParseType(v.X, typ)
	case *ast.ArrayType:
		AddToken(typ, domain.ARRAY)
		ParseType(v.Elt, typ)
	case *ast.Ident:
		typ.Name = v.Name
	}
}
