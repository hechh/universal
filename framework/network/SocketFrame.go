package network

import (
	"encoding/binary"
)

//websocket包结构
// bodySize(4B) | md5(4B) | body

type SocketFrame []byte
type SocketFrameHeader []byte

const (
	SocketFrameHeaderSize       = 8           // 包头大小
	SocketFrameBodySizeMaxLimit = 1024 * 1024 // 包体大小
	// 最大包长
	SocketFrameSizeMaxLimit = SocketFrameHeaderSize + SocketFrameBodySizeMaxLimit
)

// 包头
func (d SocketFrameHeader) GetSize() uint32 {
	return binary.LittleEndian.Uint32(d)
}

func (d SocketFrameHeader) SetSize(size uint32) {
	binary.LittleEndian.PutUint32(d, size)
}

func (d SocketFrameHeader) GetCrc32() uint32 {
	return binary.LittleEndian.Uint32(d[4:])
}

func (d SocketFrameHeader) SetCrc32(crc uint32) {
	binary.LittleEndian.PutUint32(d[4:], crc)
}
