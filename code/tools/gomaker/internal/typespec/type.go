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

func (d *Type) GetDoc() string {
	if len(d.Doc) <= 0 {
		return ""
	}
	return fmt.Sprintf("// %s", d.Doc)
}
