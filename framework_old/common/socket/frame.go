package socket

import (
	"encoding/binary"
	"universal/framework/common/fbasic"
)

// websocket包结构
// bodySize(4B) | md5(4B) | body
type Frame struct {
}

func (d *Frame) GetHeadSize() int {
	return 10
}

func (d *Frame) GetBodySize(head []byte) int {
	return int(binary.LittleEndian.Uint32(head))
}

func (d *Frame) Check(head []byte, body []byte) bool {
	oldCrc := binary.LittleEndian.Uint32(head[4:])
	crc := fbasic.GetCrc32(string(body))
	return crc == oldCrc
}

func (d *Frame) Build(frame []byte, body []byte) []byte {
	// 设置包头
	binary.LittleEndian.PutUint32(frame, uint32(len(body)))
	binary.LittleEndian.PutUint32(frame[4:], fbasic.GetCrc32(string(body)))
	// 拷贝
	headSize := d.GetHeadSize()
	copy(frame[headSize:], body)
	return frame
}
