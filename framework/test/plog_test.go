package test

import (
	"encoding/binary"
	"hash/crc32"
	"testing"
	"time"
	"universal/framework/library/plog"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestError(t *testing.T) {
	plog.Trace("[aaadfjaskdjf;alksdjf;alskjdg;alkdj;alksjdf;alksdjf;lkasdjf;alksdjf;alkjsdf;lkajsdf;klajsdf]")
	plog.Debug("[aaadfjaskdjf;alksdjf;alskjdg;alkdj;alksjdf;alksdjf;lkasdjf;alksdjf;alkjsdf;lkajsdf;klajsdf]")
	plog.Warn("[aaadfjaskdjf;alksdjf;alskjdg;alkdj;alksjdf;alksdjf;lkasdjf;alksdjf;alkjsdf;lkajsdf;klajsdf]")
	plog.Info("[aaadfjaskdjf;alksdjf;alskjdg;alkdj;alksjdf;alksdjf;lkasdjf;alksdjf;alkjsdf;lkajsdf;klajsdf]")
	plog.Error("[aaadfjaskdjf;alksdjf;alskjdg;alkdj;alksjdf;alksdjf;lkasdjf;alksdjf;alkjsdf;lkajsdf;klajsdf]")
	plog.Fatal("[aaadfjaskdjf;alksdjf;alskjdg;alkdj;alksjdf;alksdjf;lkasdjf;alksdjf;alkjsdf;lkajsdf;klajsdf]")
	time.Sleep(1 * time.Second)
}

func BenchmarkPlog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		plog.Fatal("[aaadfjaskdjf;alksdjf;alskjdg;alkdj;alksjdf;alksdjf;lkasdjf;alksdjf;alkjsdf;lkajsdf;klajsdf] %d", i)
	}
	b.Log(b.N)
}

// websocket包结构
// bodySize(4B) | md5(4B) | body
type Frame struct{}

func (d *Frame) GetHeadSize() int {
	return 10
}

func (d *Frame) GetBodySize(head []byte) int {
	return int(binary.LittleEndian.Uint32(head))
}

func (d *Frame) Check(head []byte, body []byte) bool {
	oldCrc := binary.LittleEndian.Uint32(head[4:])
	crc := crc32.ChecksumIEEE(body)
	return crc == oldCrc
}

func (d *Frame) Build(frame []byte, body []byte) []byte {
	// 设置包头
	binary.LittleEndian.PutUint32(frame, uint32(len(body)))
	binary.LittleEndian.PutUint32(frame[4:], crc32.ChecksumIEEE(body))
	// 拷贝
	headSize := d.GetHeadSize()
	copy(frame[headSize:], body)
	return frame
}
