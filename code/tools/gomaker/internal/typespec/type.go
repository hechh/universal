package typespec

import "fmt"

// 数据类型
type Type struct {
	Kind     uint32 // 基础类型
	Selector string // 引用的包名
	Name     string // 字段名称
	Doc      string // 注释
}

func (d *Type) Clone() *Type {
	return &Type{d.Kind, d.Selector, d.Name, d.Doc}
}

func (d *Type) GetName(pkg string) string {
	if len(d.Selector) <= 0 || d.Selector == pkg {
		return d.Name
	}
	return fmt.Sprintf("%s.%s", d.Selector, d.Name)
}

func (d *Type) GetDoc() string {
	if len(d.Doc) <= 0 {
		return ""
	}
	return fmt.Sprintf("// %s", d.Doc)
}

type Alias struct {
	Token    []uint32 // 数据类型
	Type     *Type    // 真实类型
	RealType *Type    // 引用类型
	Comment  string   // 注释
}

func (d *Alias) Format() string {
	return fmt.Sprintf("%s\ntype %s %s", d.Type.GetDoc(), d.Type.Name, d.RealType.GetName(d.Type.Selector))
}

func (d *Alias) Clone() *Alias {
	tmps := make([]uint32, len(d.Token))
	copy(tmps, d.Token)
	return &Alias{tmps, d.Type.Clone(), d.RealType.Clone(), d.Comment}
}
