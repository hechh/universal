package domain

const (
	KIND_IDENT    = 0
	KIND_ENUM     = 1
	KIND_ALIAS    = 2
	KIND_STRUCT   = 3
	TOKEN_NONE    = 0
	TOKEN_POINTER = 1
	TOKEN_ARRAY   = 2
	TOKEN_MAP     = 3
)

var (
	PbPath string
)

// 获取类型
type IFactory interface {
	GetKind(string) int32 // 获取类型
}

// 类型接口
type IType interface {
	Push(int32)             // 添加token
	GetKind() int32         // 获取类型
	GetName() string        // 获取类型名字
	GetPkg() string         // 获取包名
	FullName(string) string // 获取完整类型名
}

// 字段接口
type IAttribute interface {
	IType
	GetName() string // 获取字段名
}

// 类型接口
type IClass interface {
	IType
	Add(IAttribute)
	Get(string) IAttribute
	GetAll() []IAttribute
}
