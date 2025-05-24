package test

import (
	"reflect"
	"sync"
	"testing"
	"universal/framework/internal/cluster"
	"universal/library/encode"
)

var (
	rType = reflect.TypeOf(&cluster.Cluster{})
	m, _  = rType.MethodByName("Del")
)

func TestGob(t *testing.T) {
	rType := reflect.TypeOf(&cluster.Cluster{})
	m, _ := rType.MethodByName("Del")

	for i := uint32(1); i < 100; i++ {
		buf := encode.Encode(i, i)
		rets, err := encode.Decode(buf, m, 1)
		t.Log(err, "===>", rets[1].Interface(), rets[2].Interface())
	}
}

func BenchmarkGob(b *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for i := 0; i < b.N; i++ {
			buf := encode.Encode(uint32(i), uint32(i))
			if _, err := encode.Decode(buf, m, 1); err != nil {
				b.Fatal(err)
				return
			}
		}
		b.Log(b.N)
		wg.Done()
	}()
	for i := 0; i < b.N; i++ {
		buf := encode.Encode(uint32(i), uint32(i))
		if _, err := encode.Decode(buf, m, 1); err != nil {
			b.Fatal(err)
			return
		}
	}
	b.Log(b.N)
	wg.Wait()
}
