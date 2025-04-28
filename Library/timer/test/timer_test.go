package test

import (
	"encoding/binary"
	"sync/atomic"
	"testing"
	"time"
	"universal/common/pb"
	"universal/library/mlog"
	"universal/library/timer"

	"github.com/golang/protobuf/proto"
)

var count int64

func Print() {
	atomic.AddInt64(&count, 1)
}

func TestTimer(t *testing.T) {

	tt := timer.NewTimer(4, 7, 4, mlog.Fatal)

	tmps := map[uint64]*uint64{}
	for i := uint64(1); i <= 100000; i++ {
		tmps[i] = &i
		tt.AddTaskFun(tmps[i], Print, 1*time.Second, 5)
	}

	time.Sleep(6 * time.Second)
	t.Log(atomic.LoadInt64(&count))
}

/*
goos: windows
goarch: amd64
pkg: universal/library/timer/test
cpu: Intel(R) Core(TM) i5-14600KF
BenchmarkPb

	d:\project\src\universal\library\timer\test\timer_test.go:42: 1
	d:\project\src\universal\library\timer\test\timer_test.go:42: 100
	d:\project\src\universal\library\timer\test\timer_test.go:42: 10000
	d:\project\src\universal\library\timer\test\timer_test.go:42: 1000000
	d:\project\src\universal\library\timer\test\timer_test.go:42: 3487502

BenchmarkPb-20    	 3487502	       347.8 ns/op	     454 B/op	       7 allocs/op
*/
func BenchmarkPb(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pack := &pb.Packet{Head: &pb.Head{RegionID: uint64(i), UID: uint64(i), SocketID: uint32(i), ApiID: uint64(i)}, Body: []byte("1234567890")}
		buf, _ := proto.Marshal(pack)
		pack2 := &pb.Packet{}
		proto.Unmarshal(buf, pack2)
	}
	b.Log(b.N)
}

func Benchmark(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pack := &CSPacket{Header: CSPacketHeader{Version: uint16(i), Uid: uint64(i), Cmd: uint32(i), Seq: uint32(i)}, Body: []byte("1234567890")}
		buf := pack.Header.ToBytes()
		pack2 := &CSPacket{Header: CSPacketHeader{}}
		pack2.Header.From(buf)
	}
	b.Log(b.N)
}

type CSPacket struct {
	Header CSPacketHeader
	Body   []byte
}

// 注意这里的排列是考虑了内存对齐的情况，调整时请注意。
type CSPacketHeader struct {
	Version  uint16
	PassCode uint16
	Seq      uint32

	Uid uint64

	AppVersion uint32
	Cmd        uint32

	BodyLen uint32
}

func ByteLenOfCSPacketHeader() int {
	return 28
}

func ByteLenOfCSPacketBody(header []byte) int {
	return int(binary.BigEndian.Uint32(header[ByteLenOfCSPacketHeader()-4:]))
}

func (h *CSPacketHeader) From(b []byte) {
	pos := 0
	h.Version = binary.BigEndian.Uint16(b[pos:])
	pos += 2
	h.PassCode = binary.BigEndian.Uint16(b[pos:])
	pos += 2
	h.Seq = binary.BigEndian.Uint32(b[pos:])
	pos += 4
	h.Uid = binary.BigEndian.Uint64(b[pos:])
	pos += 8
	h.AppVersion = binary.BigEndian.Uint32(b[pos:])
	pos += 4
	h.Cmd = binary.BigEndian.Uint32(b[pos:])
	pos += 4
	h.BodyLen = binary.BigEndian.Uint32(b[pos:])
	pos += 4
}

func (h *CSPacketHeader) To(b []byte) {
	pos := uintptr(0)
	binary.BigEndian.PutUint16(b[pos:], h.Version)
	pos += 2
	binary.BigEndian.PutUint16(b[pos:], h.PassCode)
	pos += 2
	binary.BigEndian.PutUint32(b[pos:], h.Seq)
	pos += 4
	binary.BigEndian.PutUint64(b[pos:], h.Uid)
	pos += 8
	binary.BigEndian.PutUint32(b[pos:], h.AppVersion)
	pos += 4
	binary.BigEndian.PutUint32(b[pos:], h.Cmd)
	pos += 4
	binary.BigEndian.PutUint32(b[pos:], h.BodyLen)
	pos += 4
}

func (h *CSPacketHeader) ToBytes() []byte {
	bytes := make([]byte, ByteLenOfCSPacketHeader())
	h.To(bytes)
	return bytes
}

func (h *CSPacketHeader) Size() int {
	return ByteLenOfCSPacketHeader()
}
