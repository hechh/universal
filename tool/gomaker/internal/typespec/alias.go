package typespec

import "fmt"

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
