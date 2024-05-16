package network

import (
	"encoding/binary"
	"universal/common/pb"
	"universal/framework/fbasic"
)

// websocket包结构
// bodySize(4B) | md5(4B) | body
type SocketFrame struct {
}

func (d *SocketFrame) GetHeadSize() int {
	return 10
}

func (d *SocketFrame) GetBodySize(head []byte) int {
	return int(binary.LittleEndian.Uint32(head))
}

func (d *SocketFrame) Check(head []byte, body []byte) error {
	oldCrc := binary.LittleEndian.Uint32(head[4:])
	crc := fbasic.GetCrc32(body)
	if crc != oldCrc {
		return fbasic.NewUError(1, pb.ErrorCode_SocketFrameCheck, oldCrc, crc)
	}
	return nil
}

func (d *SocketFrame) Build(frame []byte, body []byte) []byte {
	// 设置包头
	binary.LittleEndian.PutUint32(frame, uint32(len(body)))
	binary.LittleEndian.PutUint32(frame[4:], fbasic.GetCrc32(body))
	// 拷贝
	headSize := d.GetHeadSize()
	copy(frame[headSize:], body)
	return frame
}
