package parse

import (
	"go/ast"
	"path/filepath"
	"strings"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/typespec"
)

type CfgParser struct {
	filename string            // 文件名字
	tables   map[string]string // 中文--->英文
}

func NewCfgParser() *CfgParser {
	return &CfgParser{
		tables: make(map[string]string),
	}
}

func (d *CfgParser) SetFile(filename string) {
	d.filename = filepath.Base(filename)
}

func (d *CfgParser) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *typespec.EnumNode:
	case *typespec.SheetNode:
		if n.Sheet == "代对表" {
			for _, vals := range n.Rows {
				for _, val := range vals {
					ss := strings.Split(val, ":")
					switch len(ss) {
					case 3:
						d.Visit(typespec.NewProxyNode(n.Sheet, ss))
					}
				}
			}
		} else if strings.HasPrefix(n.Sheet, "@") {
			d.Visit(typespec.NewStructNode(n.Sheet, n.Rows[0], n.Rows[1]))
		}
		return nil
	case *typespec.ProxyNode:
		if n.IsCreator {
			d.tables[n.Name] = n.English
		}
	case *typespec.StructNode:
		if name, ok := d.tables[n.Sheet]; ok {
			// 解析struct结构
			item := typespec.NewStruct(d.filename, manager.GetOrAddType(&typespec.Type{domain.KIND_STRUCT, "pb", name, ""}))
			for _, ff := range n.Fields {
				item.Add(manager.ParseRule(ff))
			}
			manager.AddStruct(item)
		}
	}
	return d
}
