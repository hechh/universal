package player

import (
	"encoding/binary"
)

type Frame struct{}

func (d *Frame) GetHeadSize() int {
	return 4
}

func (d *Frame) GetBodySize(head []byte) int {
	return int(binary.BigEndian.Uint32(head))
}

func (d *Frame) Check(head []byte, body []byte) bool {
	return true
}

func (d *Frame) Build(frame []byte, body []byte) []byte {
	// 设置包头
	binary.BigEndian.PutUint32(frame, uint32(len(body)))
	// 拷贝
	headSize := d.GetHeadSize()
	copy(frame[headSize:], body)
	return frame
}
