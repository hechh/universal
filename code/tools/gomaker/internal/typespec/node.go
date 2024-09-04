package typespec

import (
	"go/token"
	"strings"
)

// 代对节点
type ProxyNode struct {
	SheetName string
	IsCreator bool
	Value     string
}

func (d *ProxyNode) Pos() token.Pos { return 0 }
func (d *ProxyNode) End() token.Pos { return 0 }

// 枚举节点
type EnumNode struct {
	SheetName string
	Value     string
}

func (d *EnumNode) Pos() token.Pos { return 0 }
func (d *EnumNode) End() token.Pos { return 0 }

// 配置节点
type FieldNode struct {
	IsProxy  bool
	Index    int
	Type     string
	Name     string
	Original string
	Doc      string
}

type TableNode struct {
	SheetName string
	Fields    []*FieldNode
}

func (d *TableNode) Pos() token.Pos { return 0 }
func (d *TableNode) End() token.Pos { return 0 }

// table值节点
type ValueNode struct {
	SheetName string
	Values    []string
}

func (d *ValueNode) Pos() token.Pos { return 0 }
func (d *ValueNode) End() token.Pos { return 0 }

func NewTableNode(sheet string, defines, docs []string) *TableNode {
	docFunc := func(i int) string {
		if len(docs) <= i {
			return ""
		}
		return docs[i]
	}
	tmps := []*FieldNode{}
	for pos, name := range defines {
		if index := strings.Index(name, "_"); index > 0 {
			tmps = append(tmps, &FieldNode{
				IsProxy:  strings.HasPrefix(name, "$"),
				Original: strings.TrimPrefix(name, "$"),
				Index:    pos,
				Type:     strings.ToLower(strings.TrimSpace(strings.TrimPrefix(name[:index], "$"))),
				Name:     name[index+1:],
				Doc:      docFunc(pos),
			})
		}
	}
	/*
		buf, _ := json.Marshal(&tmps)
		fmt.Println(sheet, "===>", string(buf))
	*/
	return &TableNode{SheetName: sheet, Fields: tmps}
}
