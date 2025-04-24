package mlog

import (
	"testing"
	"time"
	"universal/library/builder"
)

func TestMain(m *testing.M) {
	Init("./log", "test", LOG_DEBUG)
	m.Run()
}

func TestBuilder(t *testing.T) {
	ww := builder.NewFileBuilder(1024 * 1024)
	ww.SetWriter("./log/file.test")
	for i := 0; i < 10000; i++ {
		ww.Write([]byte("adsfakls;dfja;lksdfj;aklsdjf;aklsdjf;aklsdfjal;dks"))
	}
	ww.Flush()
	ww.Close()
}

func TestWriter(t *testing.T) {
	ww := NewWriter("./log", "test", 1024)
	ww.Write(NewFormat(1, LOG_DEBUG, "test %d", 1))
	ww.Write(NewFormat(1, LOG_DEBUG, "test %d", 1))
	time.Sleep(6 * time.Second)
	if err := ww.Close(); err != nil {
		t.Log(err)
	}
}

// go test -benchmem -run=BenchmarkError -bench -count=1 -v -cpuprofile=cpu.prof
func BenchmarkError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Error("--%d-------askfja;lksdjf;alkdjf;alkjdf;alsjkdf-------", i)
	}
	b.Log(b.N)
}
