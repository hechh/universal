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
	for i := int32(-1000); i < 1000; i++ {
		// 编码
		buf := encoding.IntegerEncode(i)
		// 解码
		val := encoding.IntegerDecode(buf)
		if vv, ok := val.(int32); !ok || vv != i {
			t.Log(i, ok, val)
		}
	}
	for i := int16(-1000); i < 1000; i++ {
		// 编码
		buf := encoding.IntegerEncode(i)
		// 解码
		val := encoding.IntegerDecode(buf)
		if vv, ok := val.(int16); !ok || vv != i {
			t.Log(i, ok, val)
		}
	}
	for i := int8(-100); i < 100; i++ {
		// 编码
		buf := encoding.IntegerEncode(i)
		// 解码
		val := encoding.IntegerDecode(buf)
		if vv, ok := val.(int8); !ok || vv != i {
			t.Log(i, ok, val)
		}
	}
}
