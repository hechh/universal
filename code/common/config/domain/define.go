package domain

type IConfig interface {
	LoadFile([]byte) error // 加载配置数据
}
