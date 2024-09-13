package typespec

import "fmt"

type Type struct {
	Kind    uint32 // 基础类型
	PkgName string // 引用的包名
	Name    string // 字段名称
	Doc     string // 注释
}

func (d *Type) GetName(pkg string) string {
	if len(d.PkgName) <= 0 || d.PkgName == pkg {
		return d.Name
	}
	return fmt.Sprintf("%s.%s", d.PkgName, d.Name)
}

func (d *Type) GetPkgType() string {
	if len(d.PkgName) <= 0 {
		return d.Name
	}
	return fmt.Sprintf("%s.%s", d.PkgName, d.Name)
}

func (d *Type) GetDoc() string {
	if len(d.Doc) <= 0 {
		return ""
	}
	return fmt.Sprintf("// %s", d.Doc)
}

type Alias struct {
	Token    []uint32 // 数据类型
	Type     *Type    // 引用类型
	RealType *Type    // 真实类型
	Doc      string   // 注释
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
