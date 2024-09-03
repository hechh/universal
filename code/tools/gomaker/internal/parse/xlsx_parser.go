package parse

import (
	"fmt"
	"go/ast"
	"path/filepath"
	"strings"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/typespec"

	"github.com/spf13/cast"
)

type XlsxParser struct {
	filename string                     // 文件名字
	enums    map[string]*typespec.Value // 中文--->所有代对
	tables   map[string]string          // 中文--->英文
	fields   []*typespec.FieldNode      // 字段信息
}

func NewXlsxParser() *XlsxParser {
	return &XlsxParser{
		enums:  make(map[string]*typespec.Value),
		tables: make(map[string]string),
	}
}

func (d *XlsxParser) SetFile(name string) {
	d.filename = filepath.Base(name)
}

func (d *XlsxParser) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *typespec.EnumNode:
		ss := strings.Split(strings.TrimSpace(n.Value), ":")
		item := &typespec.Value{
			Name: fmt.Sprintf("%s_%s", ss[2], ss[3]),
			Type: manager.GetOrAddType(&typespec.Type{
				Kind:    domain.KIND_ENUM,
				PkgName: "pb",
				Name:    ss[2],
			}),
			Value: cast.ToInt32(ss[4]),
			Doc:   ss[1],
		}
		d.enums[item.Doc] = item
		manager.AddValue(item, d.filename)

	case *typespec.ProxyNode:
		if n.IsCreator {
			ss := strings.Split(strings.TrimSpace(n.Value), ":")
			d.tables[fmt.Sprintf("@%s", ss[1])] = ss[2]
		}

	case *typespec.TableNode:
		name, ok := d.tables[n.SheetName]
		if !ok {
			return nil
		}
		d.fields = n.Fields
		tmp := typespec.NewStruct(manager.GetOrAddType(&typespec.Type{
			Kind:    domain.KIND_STRUCT,
			PkgName: "pb",
			Name:    name,
		}), d.filename)
		for _, ff := range n.Fields {
			tmp.Add(manager.ParseRule(ff))
		}
		manager.AddStruct(tmp)

	case *typespec.ValueNode:
		if sheetName, ok := d.tables[n.SheetName]; ok {
			name := fmt.Sprintf("%s.%s.json", strings.TrimSuffix(d.filename, filepath.Ext(d.filename)), sheetName)
			tmps := map[string]interface{}{}
			for _, ff := range d.fields {
				tmps[ff.Name] = manager.CastRule(ff, d.enums, n.Values[ff.Index])
			}
			manager.AddJson(name, tmps)
		}
	}
	return d
}
