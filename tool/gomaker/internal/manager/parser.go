package manager

import (
	"go/ast"
	"go/token"
	"universal/tool/gomaker/domain"

	"github.com/spf13/cast"
)

var (
	enums   = make(map[string]*domain.Enum)
	alias   = make(map[string]*domain.Alias)
	structs = make(map[string]*domain.Struct)
)

type TypeParser struct {
	pkgName string
}

func getDoc(n *ast.CommentGroup) string {
	if n == nil || len(n.List) <= 0 {
		return ""
	}
	return n.List[0].Text
}

func (d *TypeParser) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.File:
		d.pkgName = n.Name.Name
		return d
	case *ast.GenDecl:
		if n.Tok != token.CONST && n.Tok != token.TYPE {
			return nil
		}
		// 遍历所有节点
		for _, node := range n.Specs {
			switch vv := node.(type) {
			case *ast.ValueSpec:
				d.addEnum(getDoc(n.Doc), vv)
			case *ast.TypeSpec:
				/*
					switch vvv := vv.Type.(type) {
					case *ast.StructType:
							item := &domain.Struct{
								Type: &domain.Type{
									Token:   domain.STRUCT,
									PkgName: d.pkgName,
									Name:    vv.Name.Name,
								},
								Doc:    getDoc(n.Doc),
								Fields: make(map[string]*domain.Field),
							}
							for _, field := range vvv.Fields.List {
								name := field.Names[0].Name
								// 过滤不对外开放接
								if unicode.IsLower(rune(name[0])) {
									continue
								}
								ff := &domain.Field{
									Type: &domain.Type{PkgName: d.pkgName, },
									Name: name,
								}
							}
					default:
					}
				*/
			}
		}
	}
	return nil
}

func (d *TypeParser) addEnum(doc string, vv *ast.ValueSpec) {
	val := &domain.Value{
		Name:    vv.Names[0].Name,
		Value:   cast.ToInt32(vv.Values[0].(*ast.BasicLit).Value),
		Comment: getDoc(vv.Comment),
	}
	// 获取enum类型名称
	enumType := vv.Type.(*ast.Ident).Name
	vvv, ok := enums[enumType]
	if !ok {
		vvv = &domain.Enum{
			Type: &domain.Type{
				Token:   domain.ENUM,
				PkgName: d.pkgName,
				Name:    enumType,
			},
			Doc:    doc,
			Fields: make(map[string]*domain.Value),
		}
		enums[enumType] = vvv
	}
	// 保存
	vvv.Fields[val.Name] = val
	vvv.List = append(enums[enumType].List, val)
}
