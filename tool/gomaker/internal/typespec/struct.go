package typespec

import (
	"fmt"
	"strings"
	"universal/tool/gomaker_old/domain"
)

// 字段结构
type Field struct {
	Token   []uint32 // 数据类型
	Name    string   // 字段名字
	Type    *Type    // 类型
	Tag     string   // 标签
	Comment string   // 注释
}

type Struct struct {
	Type   *Type             // 类型
	Fields map[string]*Field // 字段
	List   []*Field          // 排序队列
}

func (d *Struct) Format() string {
	vals := []string{}
	for _, val := range d.List {
		strToken := ""
		for _, val := range val.Token {
			switch val {
			case domain.POINTER:
				strToken += "*"
			case domain.ARRAY:
				strToken += "[]"
			}
		}
		if len(val.Comment) > 0 {
			vals = append(vals, fmt.Sprintf("%s %s%s %s // %s", val.Name, strToken, val.Type.Format(d.Type.Selector), val.Tag, val.Comment))
		} else {
			vals = append(vals, fmt.Sprintf("%s %s%s %s", val.Name, strToken, val.Type.Format(d.Type.Selector), val.Tag))
		}
	}
	return fmt.Sprintf("type %s struct {\n%s\n}", d.Type.Name, strings.Join(vals, "\n"))
}
