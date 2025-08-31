package parse

import (
	"go/ast"
	"go/token"
	"universal/library/util"
	"universal/tool/pbtool/domain"
)

type TypeParser struct {
	kinds map[string]int32
}

func NewTypeParser() *TypeParser {
	ret := &TypeParser{kinds: make(map[string]int32)}
	ret.kinds["uint32"] = domain.KIND_IDENT
	ret.kinds["uint64"] = domain.KIND_IDENT
	ret.kinds["int"] = domain.KIND_IDENT
	ret.kinds["int32"] = domain.KIND_IDENT
	ret.kinds["int64"] = domain.KIND_IDENT
	ret.kinds["bool"] = domain.KIND_IDENT
	ret.kinds["float32"] = domain.KIND_IDENT
	ret.kinds["float64"] = domain.KIND_IDENT
	ret.kinds["string"] = domain.KIND_IDENT
	ret.kinds["[]byte"] = domain.KIND_IDENT
	return ret
}

func (p *TypeParser) GetKind(name string) int32 {
	return p.kinds[name]
}

func (p *TypeParser) Visit(n ast.Node) ast.Visitor {
	switch vv := n.(type) {
	case *ast.File:
		return p
	case *ast.GenDecl:
		return util.Or[*TypeParser](vv.Tok != token.TYPE, nil, p)
	case *ast.TypeSpec:
		switch nn := vv.Type.(type) {
		case *ast.StructType:
			p.kinds[vv.Name.Name] = domain.KIND_STRUCT
		case *ast.Ident:
			if nn.Name == "int32" {
				p.kinds[vv.Name.Name] = domain.KIND_ENUM
			} else {
				p.kinds[vv.Name.Name] = domain.KIND_IDENT
			}
		}
	}
	return nil
}
