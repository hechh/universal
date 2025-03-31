package base

type Type struct {
	Name    string // 类型名称
	TypeOf  uint32 // 1表示内置类型，2表示枚举，3表示结构体
	ValueOf uint32 // 1 表示单值，2:表示数组，3:表示map，4 表示group
}

type Field struct {
	Name     string // 字段名
	Type     *Type  // 字段类型
	Desc     string // 字段标签
	Position int    // 字段索引
}

type Struct struct {
	Name      string              // 结构体名称
	List      []*Field            // 字段类型
	Converts  map[string][]*Field // 转换表
	SheetName string              // 表明
	FileName  string              // 文件名
}

// 生成表
type Config struct {
	Name      string   // 表名称
	List      []*Field // 表列表
	SheetName string   // 表明
	FileName  string   // 文件名
}

// 枚举类型定义
type Enum struct {
	Name      string             // 枚举名称
	Values    map[string]*EValue // 枚举值
	SheetName string             // 表明
	FileName  string             // 文件名
}

type EValue struct {
	Name  string // 枚举值名称
	Value uint32 // 枚举值
	Desc  string // 枚举值描述
}

type Value struct {
	TypeOf    uint32
	Type      string // 类型名称
	Name      string // 枚举值名称
	Value     uint32 // 枚举值
	Desc      string // 枚举值描述
	SheetName string // 表明
	FileName  string // 文件名
}
