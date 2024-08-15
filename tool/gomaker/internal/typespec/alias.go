package typespec

import "fmt"

type Alias struct {
	Token    []uint32 // 数据类型
	Type     *Type    // 真实类型
	RealType *Type    // 引用类型
	Comment  string   // 注释
}

func (d *Alias) Format() string {
	return fmt.Sprintf("// %s\ntype %s %s", d.Type.Doc, d.Type.Name, d.RealType.Format(d.Type.Selector))
}
