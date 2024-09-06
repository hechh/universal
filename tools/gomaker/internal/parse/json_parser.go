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

func (d *JsonParser) parseEnum(n *typespec.SheetNode) {
	for _, vals := range n.Rows {
		for _, val := range vals {
			ss := strings.Split(val, ":")
			if len(ss) != 4 {
				continue
			}

			if d.Visit(typespec.NewEnumNode(n.Sheet, ss)) == nil {
				return
			}
		}
	}
}

func (d *JsonParser) parseProxy(n *typespec.SheetNode) {
	for _, vals := range n.Rows {
		for _, val := range vals {
			ss := strings.Split(val, ":")
			switch len(ss) {
			case 4:
				d.Visit(typespec.NewInnerEnumNode(n.Sheet, ss))
			case 3:
				d.Visit(typespec.NewProxyNode(n.Sheet, ss))
			}
		}
	}
}

func (d *JsonParser) parseStruct(n *typespec.SheetNode) {
	if d.Visit(typespec.NewStructNode(n.Sheet, n.Rows[0], n.Rows[1])) == nil {
		return
	}
	d.Visit(&typespec.ValueNode{n.Sheet, n.Rows[2:]})
}

func (d *JsonParser) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *typespec.SheetNode:
		if d.filename == "define.xlsx" {
			d.parseEnum(n)
		} else if n.Sheet == "代对表" {
			d.parseProxy(n)
		} else if strings.HasPrefix(n.Sheet, "@") {
			d.parseStruct(n)
		}
		return nil
	case *typespec.EnumNode:
		tt := &typespec.Type{domain.KIND_ENUM, "pb", n.Type, ""}
		item := typespec.NewValue(manager.GetOrAddType(tt), n.Name, int32(n.Value), n.Doc)
		d.enums[item.Doc] = item
		manager.AddEnumValue(d.filename, item)
	case *typespec.InnerEnumNode:
		tt := &typespec.Type{domain.KIND_ENUM, "pb", n.Type, ""}
		item := typespec.NewValue(manager.GetOrAddType(tt), n.Name, int32(n.Value), n.Doc)
		d.enums[item.Doc] = item
		manager.AddEnumValue(d.filename, item)
	case *typespec.ProxyNode:
		if n.IsCreator {
			d.tables[n.Name] = n.English
		}
	case *typespec.StructNode:
		if name, ok := d.tables[n.Sheet]; ok {
			item := typespec.NewStruct(d.filename, manager.GetOrAddType(&typespec.Type{domain.KIND_STRUCT, "pb", name, ""}))
			for _, ff := range n.Fields {
				item.Add(manager.ParseRule(ff))
			}
			manager.AddStruct(item)
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
