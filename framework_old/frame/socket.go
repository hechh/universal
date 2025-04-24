package frame

type IFrame interface {
	GetHeadSize() int       // 包头大小
	GetBodySize([]byte) int // 获取包体大小
}
