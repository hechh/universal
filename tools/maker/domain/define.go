package domain

var (
	DefaultPkg = "pb"
	// 基础数据类型, go类型映射到proto类型
	BasicTypes = map[string]bool{
		"uint32":  true,
		"uint64":  true,
		"int":     true,
		"int32":   true,
		"int64":   true,
		"bool":    true,
		"float32": true,
		"float64": true,
		"string":  true,
		"[]byte":  true,
	}
)

// 配置表规则
const (
	RuleTypeEnum       = "E:"
	GomakerTypeEnum    = "@gomaker:enum"
	GomakerTypeMessage = "@gomaker:message"
	GenTable           = "生成表"
	DefaultEnumClass   = "other"
	SYNTAX             = "syntax"
	PACKAGE            = "package"
	OPTION             = "option"
	IMPORT             = "import"
	MESSAGE            = "message"
	ENUM               = "enum"
	KindTypeIdent      = 0
	KindTypeEnum       = 1
	KindTypeAlias      = 2
	KindTypeStruct     = 3
	TokenTypeNone      = 0
	TokenTypePointer   = 1
	TokenTypeArray     = 2
	TokenTypeMap       = 3
	SourceTypeXlsx     = 1
	SourceTypeProto    = 2
	SourceTypeGo       = 3
)

// 代码生成
type GenFunc func(dst string, extra ...string) error

// xlsx转bytes
type ConvFunc func(string) interface{}
