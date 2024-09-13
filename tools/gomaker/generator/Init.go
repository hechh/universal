package generator

import "universal/tools/gomaker/internal/manager"

func Init() {
	manager.Register("client", "生成client代码", HttpKitGenerator, OmitEmptyGenerator, ProtoGenerator)
	manager.Register("pb", "xlsx转pb结构", enumGenerator, tableGenerator)
}
