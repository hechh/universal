package domain

const (
	INIT    = 0
	BASE    = 0x01
	ENUM    = 0x02
	STRUCT  = 0x04
	POINTER = 0x08
	ARRAY   = 0x10
	MAP     = 0x20
)

type Type struct {
	Token   int32
	Name    string
	PkgName string
}

type Field struct {
	Name string
	Type *Type
}

type MapField struct {
	Name  string
	KType *Type
	VType *Type
}

type Struct struct {
	Type   *Type
	Maps   []*MapField
	Others []*Field
}

type EnumValue struct {
	Name    string
	Value   int32
	Comment string
}

type Enum struct {
	Type   *Type
	Fields map[string]*EnumValue
}
