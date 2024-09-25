package manager

import (
	"fmt"
	"universal/tools/gomaker/internal/typespec"
)

var (
	menums   = []*typespec.Sheet{}
	messages = []*typespec.Sheet{}
)

func GetMEnumsReference() *[]*typespec.Sheet {
	return &menums
}

func GetMessageReference() *[]*typespec.Sheet {
	return &messages
}

func GetMEnumList() []*typespec.Sheet {
	return menums
}

func GetMessageList() []*typespec.Sheet {
	return messages
}

// 修复sheet结构中的类型数据
func InitSheet() {
	for _, sh := range messages {
		sh.Struct = GetStruct(fmt.Sprintf("pb.%s", sh.Config))
		// 延后解析group
		for _, gg := range sh.Group {
			for i, ff := range gg {
				gg[i] = sh.Struct.Fields[ff.Name]
			}
		}
		// 延后解析map
		for _, mm := range sh.Map {
			for i, ff := range mm {
				mm[i] = sh.Struct.Fields[ff.Name]
			}
		}
	}
}
