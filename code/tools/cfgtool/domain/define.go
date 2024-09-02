package domain

type Enum struct {
	Type  string // 枚举类型
	Name  string // 枚举名字
	Value int32  // 枚举值
	Doc   string // 枚举注释
}

type Field struct {
	IsEnum   bool   // 是否为枚举
	Index    int    // 下标
	Type     string // 类型
	Name     string // 字段名字
	Original string // 原始值
	Doc      string // 注释
}

type Table struct {
	IsServer bool     // 是否生成后台配置
	Name     string   // 英文名
	Fields   []*Field // 字段数据
}

type FileType struct {
	Name   string             // 文件名
	Enums  map[string][]*Enum // 枚举类型 -- 枚举值
	Alls   map[string]*Enum   // 所有枚举类型
	Tables map[string]*Table  // 配置表 中文名-表
}
