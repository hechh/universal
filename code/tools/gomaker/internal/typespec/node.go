package typespec

import (
	"fmt"
	"go/token"
	"strings"

	"github.com/spf13/cast"
)

type SheetNode struct {
	Sheet string
	Rows  [][]string
}

func (d *SheetNode) Pos() token.Pos { return 0 }
func (d *SheetNode) End() token.Pos { return 0 }

// 枚举节点
type EnumNode struct {
	Sheet string
	Type  string
	Name  string
	Value uint32
	Doc   string
}

func NewEnumNode(sheet string, ss []string) *EnumNode {
	return &EnumNode{
		Sheet: sheet,
		Type:  ss[2],
		Name:  fmt.Sprintf("%s_%s", ss[2], ss[3]),
		Value: cast.ToUint32(ss[4]),
		Doc:   ss[1],
	}
}

func (d *EnumNode) Pos() token.Pos { return 0 }
func (d *EnumNode) End() token.Pos { return 0 }

// 枚举节点
type InnerEnumNode struct {
	Sheet string
	Type  string
	Name  string
	Value uint32
	Doc   string
}

func NewInnerEnumNode(sheet string, ss []string) *EnumNode {
	return &EnumNode{
		Sheet: sheet,
		Type:  ss[2],
		Name:  fmt.Sprintf("%s_%s", ss[2], ss[3]),
		Value: cast.ToUint32(ss[4]),
		Doc:   ss[1],
	}
}

func (d *InnerEnumNode) Pos() token.Pos { return 0 }
func (d *InnerEnumNode) End() token.Pos { return 0 }

// 代对节点
type ProxyNode struct {
	Sheet     string // sheet名
	IsCreator bool   // 是否需要生成服务器配置代码
	English   string // 英文名
	Name      string // 中文名
}

func NewProxyNode(sheet string, ss []string) *ProxyNode {
	return &ProxyNode{
		Sheet:     sheet,
		IsCreator: strings.Contains(strings.ToLower(ss[0]), "s"),
		English:   ss[2],
		Name:      fmt.Sprintf("@%s", ss[1]),
	}
}

func (d *ProxyNode) Pos() token.Pos { return 0 }
func (d *ProxyNode) End() token.Pos { return 0 }

// 配置节点
type FieldNode struct {
	IsProxy  bool
	Index    int
	Type     string
	Name     string
	Original string
	Doc      string
}

type StructNode struct {
	Sheet  string
	Fields []*FieldNode
}

func NewStructNode(sheet string, defines, docs []string) (ret *StructNode) {
	ret = &StructNode{Sheet: sheet}
	for i, name := range defines {
		if len(docs) < len(defines) {
			docs = append(docs, "")
		}
		isProxy := strings.HasPrefix(name, "$")
		original := strings.TrimPrefix(name, "$")
		if pos := strings.Index(original, "_"); pos > 0 {
			ret.Fields = append(ret.Fields, &FieldNode{
				IsProxy:  isProxy,
				Original: original,
				Index:    i,
				Doc:      docs[i],
				Type:     strings.ToLower(original[:pos]),
				Name:     name[pos+1:],
			})
		}
	}
	return
}

func (d *StructNode) Pos() token.Pos { return 0 }
func (d *StructNode) End() token.Pos { return 0 }

// table值节点
type ValueNode struct {
	Sheet  string
	Values [][]string
}

func (d *ValueNode) Pos() token.Pos { return 0 }
func (d *ValueNode) End() token.Pos { return 0 }
