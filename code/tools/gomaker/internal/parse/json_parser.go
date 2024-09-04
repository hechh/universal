package parse

import (
	"fmt"
	"go/ast"
	"path/filepath"
	"strings"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/typespec"
	"universal/tools/gomaker/internal/util"
)

type JsonParser struct {
	dst      string                     // 输出目录
	filename string                     // 文件名字
	enums    map[string]*typespec.Value // 中文--->所有代对
	tables   map[string]string          // 中文--->英文
	fields   []*typespec.FieldNode      // 字段信息
}

func NewJsonParser(dst string) *JsonParser {
	return &JsonParser{
		dst:    dst,
		enums:  make(map[string]*typespec.Value),
		tables: make(map[string]string),
	}
}

func (d *JsonParser) SetFile(filename string) {
	d.filename = filepath.Base(filename)
}

func (d *JsonParser) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *typespec.EnumNode:
		item := &typespec.Value{
			Type:  manager.GetOrAddType(&typespec.Type{domain.KIND_ENUM, "pb", n.Type, ""}),
			Name:  n.Name,
			Value: int32(n.Value),
			Doc:   n.Doc,
		}
		d.enums[item.Doc] = item
		manager.AddValue(d.filename, item)
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
		if sheetName, ok := d.tables[n.Sheet]; ok {
			name := fmt.Sprintf("%s.%s.json", strings.TrimSuffix(d.filename, filepath.Ext(d.filename)), sheetName)
			result := []map[string]interface{}{}
			for _, vals := range n.Values {
				result = append(result, d.toMap(vals))
			}
			util.SaveJson(filepath.Join(d.dst, name), result)
		}
	}
	return d
}

func (d *JsonParser) toMap(vals []string) map[string]interface{} {
	tmp := map[string]interface{}{}
	for _, ff := range d.fields {
		if ff.Index < len(vals) {
			tmp[ff.Original] = manager.CastRule(ff, d.enums, vals[ff.Index])
		}
	}
	return tmp
}
