package parse

import (
	"go/ast"
	"path/filepath"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/typespec"
)

type CfgParser struct {
	filename string                     // 文件名字
	enums    map[string]*typespec.Value // 中文--->所有代对
	tables   map[string]string          // 中文--->英文
	fields   []*typespec.FieldNode      // 字段信息
}

func NewCfgParser() *CfgParser {
	return &CfgParser{
		enums:  make(map[string]*typespec.Value),
		tables: make(map[string]string),
	}
}

func (d *CfgParser) SetFile(filename string) {
	d.filename = filepath.Base(filename)
}

func (d *CfgParser) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *typespec.EnumNode:
		return nil
	case *typespec.ProxyNode:
		if n.IsCreator {
			d.tables[n.Name] = n.English
		}
	case *typespec.TableNode:
		if name, ok := d.tables[n.Sheet]; ok {
			// 解析struct结构
			item := typespec.NewStruct(d.filename, manager.GetOrAddType(&typespec.Type{domain.KIND_STRUCT, "pb", name, ""}))
			for _, ff := range n.Fields {
				item.Add(manager.ParseRule(ff))
			}
			manager.AddStruct(item)

			// 临时存储
			d.fields = n.Fields
		}
	case *typespec.ValueNode:
		return nil
	}
	return d
}

func (d *CfgParser) toMap(vals []string) map[string]interface{} {
	tmp := map[string]interface{}{}
	for _, ff := range d.fields {
		if ff.Index < len(vals) {
			tmp[ff.Original] = manager.CastRule(ff, d.enums, vals[ff.Index])
		}
	}
	return tmp
}
