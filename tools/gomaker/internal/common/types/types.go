package types

type Type struct {
	Token   int32  // 类型
	PkgName string // 引用类型所在的包
	Name    string // 引用的类型名称
	Key     *Type  // map特殊处理（key类型）
	Value   *Type  // map特殊处理（value类型）
}

type Field struct {
	Type    *Type  // 类型
	Name    string // 字段名字
	Comment string // 注释
	Tag     string // 标签
}

type Struct struct {
	Type   *Type             // 类型
	Doc    string            // 注释
	Fields map[string]*Field // 解析的字段
	List   []*Field          // 排序
}

type Value struct {
	Name    string
	Value   int32
	Comment string
}

type Enum struct {
	Type   *Type             // 类型
	Doc    string            // 注释规则
	Fields map[string]*Value // 解析的字段
	List   []*Value          // 排序队列
}

type Alias struct {
	Type      *Type  // 别名类型
	Doc       string // 规则注释
	Comment   string // 字段注释
	Reference *Type  // 引用类型
}
