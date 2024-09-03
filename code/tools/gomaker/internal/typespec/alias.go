package typespec

import "fmt"

type Alias struct {
	Token    []uint32 // 数据类型
	Type     *Type    // 引用类型
	RealType *Type    // 真实类型
	FileName string   // 文件名
	Doc      string   // 注释
}

func (d *Alias) GetType() string {
	return d.Type.GetName(d.Type.PkgName)
}

func (d *Alias) GetDoc() string {
	if len(d.Doc) <= 0 {
		return ""
	}
	return fmt.Sprintf("// %s", d.Doc)
}

func (d *Alias) Format() string {
	return fmt.Sprintf("%s\ntype %s %s", d.Type.GetDoc(), d.Type.Name, d.RealType.GetName(d.Type.PkgName))
}
