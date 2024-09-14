package generator

import "universal/tools/gomaker_new/internal/manager"

func Init() {
	manager.Register("proto", "xlsx转pb结构", enumGenerator, tableGenerator)
}
