package domain

const (
	// 类类型
	TypeOfConfig = 4
	TypeOfStruct = 3
	TypeOfEnum   = 2
	TypeOfBase   = 1

	// 值类型
	ValueOfBase  = 1
	ValueOfList  = 2
	ValueOfMap   = 3
	ValueOfGroup = 4
)

var (
	PkgName = ""
)

/*
@config|sheet:结构名|map:字段名[,字段名]:别名|group:字段名[,字段名]:别名
@struct|sheet:结构名
@enum|sheet
E|道具类型-金币|PropertType|Coin|1
*/
