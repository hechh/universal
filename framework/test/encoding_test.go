package test

import (
	"testing"
	"universal/framework/library/encoding"
)

func TestEncoding(t *testing.T) {
	for i := int64(-1000); i < 1000; i++ {
		// 编码
		buf := encoding.IntegerEncode(i)
		// 解码
		val := encoding.IntegerDecode(buf)
		if vv, ok := val.(int64); !ok || vv != i {
			t.Log(i, ok, val)
		}
	}
}
