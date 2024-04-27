package network

import "encoding/binary"

//websocket包结构
// bodySize(4B) | md5(4B) | body

type WebSocketHeader []byte

const (
	WebSocketHeaderSize  = 8
	WebSocketMaxBodySize = 1024 * 1024
	WebSocketMaxLimit    = WebSocketHeaderSize + WebSocketMaxBodySize
)

func (d WebSocketHeader) GetSize() uint32 {
	return binary.LittleEndian.Uint32(d)
}

func (d WebSocketHeader) SetSize(size uint32) {
	binary.LittleEndian.PutUint32(d, size)
}

func (d WebSocketHeader) GetCrc32() uint32 {
	return binary.LittleEndian.Uint32(d[4:])
}

func (d WebSocketHeader) SetCrc32(crc uint32) {
	binary.LittleEndian.PutUint32(d[4:], crc)
}
