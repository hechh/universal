package typespec

import (
	"fmt"
	"go/ast"
	"sort"

	"github.com/spf13/cast"
)

type Value struct {
	Name    string
	Value   int32
	Comment string
}

type Enum struct {
	BaseFunc
	PkgName string            // 所在包名
	Name    string            // 引用的类型名称
	Doc     string            // 注释规则
	Fields  map[string]*Value // 解析的字段
	List    []*Value          // 排序队列
}

func NewEnum(pkg, doc string, specs []ast.Spec) *Enum {
	item := &Enum{
		PkgName: pkg,
		Doc:     doc,
		Fields:  make(map[string]*Value),
	}
	for _, node := range specs {
		// 判断是否有效
		vv, ok := node.(*ast.ValueSpec)
		if !ok || vv == nil || vv.Type == nil {
			continue
		}
		// 解析枚举类型
		item.Name = vv.Type.(*ast.Ident).Name
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

func (d *Enum) GetTypeString(pkg string) string {
	if len(d.PkgName) > 0 && len(pkg) > 0 && pkg != d.PkgName {
		return d.PkgName + "." + d.Name
	}
	return d.Name
}

func (d *Enum) String() string {
	str := ""
	for _, v := range d.Fields {
		str += fmt.Sprintf("\t %s %s = %d %s\n", v.Name, d.Name, v.Value, v.Comment)
	}
	return fmt.Sprintf("type %s int32\nconst (\n %s )", d.Name, str)
}
