package domain

const (
	INIT    = 0
	BASE    = 0x01 // 基础类型
	ENUM    = 0x02 // 枚举类型
	STRUCT  = 0x04 // 自定义类型
	POINTER = 0x08 // 指针类型
	ARRAY   = 0x10 // 数据类型
)

type AstType struct {
	Token int32
	Name  string
}

type AstField struct {
	Type *AstType
	Name string
}

type AstMapField struct {
	KType *AstType
	VType *AstType
	Name  string
}

type AstValue struct {
	Type   *AstType
	Name   string
	Value  int32
	StrVal string
}

type AstEnum struct {
	Type   *AstType
	Values map[string]*AstValue // field = value
}

type AstStruct struct {
	Type   *AstType
	Idents []*AstField
	Arrays []*AstField
	Maps   []*AstMapField
}

func (d *AstType) IsPointer() bool {
	return d.Token&POINTER > 0
}
func (d *AstType) IsBase() bool {
	return d.Token&BASE > 0
}
func (d *AstType) IsEnum() bool {
	return d.Token&ENUM > 0
}
func (d *AstType) IsStruct() bool {
	return d.Token&STRUCT > 0
}
func (d *AstType) IsArray() bool {
	return d.Token&ARRAY > 0
}

func (d *AstType) GetType() string {
	ret := ""
	if d.Token&ARRAY > 0 {
		ret += "[]"
	}
	if d.Token&POINTER > 0 {
		ret += "*"
	}
	if d.Token&ENUM > 0 || d.Token&STRUCT > 0 {
		ret += "pb."
	}
	return ret + d.Name
}

func (d *AstType) IsReward() bool {
	return (d.Name == "Reward") && (d.Token&STRUCT > 0)
}

// ----- safe -----
func (d *AstField) GetNameS() string {
	return "_" + d.Name
}

func (d *AstMapField) GetNameS() string {
	return "_" + d.Name
}

func (d *AstType) GetTypeS() string {
	ret := ""
	if d.Token&ARRAY > 0 {
		ret += "[]"
	}
	if d.Token&POINTER > 0 {
		ret += "*"
	}
	if d.Token&STRUCT > 0 {
		return ret + d.Name + "S"
	}
	return ret + d.Name
}
