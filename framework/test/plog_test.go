package test

import (
	"testing"
	"time"
	"universal/framework/common/plog"
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
