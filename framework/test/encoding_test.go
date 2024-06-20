package test

import (
	"testing"
	"universal/framework/library/encoding"
)

func TestEncoding(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		buf := encoding.Encode(true)
		val := encoding.Decode(buf)
		if vv, ok := val.(bool); !ok || vv != true {
			t.Log(true, ok, vv)
		}
	})
	t.Run("false", func(t *testing.T) {
		buf := encoding.Encode(false)
		val := encoding.Decode(buf)
		if vv, ok := val.(bool); !ok || vv != false {
			t.Log(false, ok, vv)
		}
	})
	t.Run("float", func(t *testing.T) {
		for i := float32(-134.123); i < 1000; i++ {
			// 编码
			buf := encoding.Encode(i)
			// 解码
			val := encoding.Decode(buf)
			if vv, ok := val.(float32); !ok || vv != i {
				t.Log("float32", i, ok, val)
			}
		}
		for i := float64(-12134.123234); i < 1000; i++ {
			// 编码
			buf := encoding.Encode(i)
			// 解码
			val := encoding.Decode(buf)
			if vv, ok := val.(float64); !ok || vv != i {
				t.Log("float64", i, ok, val)
			}
		}
	})
	t.Run("uinteger", func(t *testing.T) {
		for i := uint64(0); i < 1000; i++ {
			// 编码
			buf := encoding.Encode(i)
			// 解码
			val := encoding.Decode(buf)
			if vv, ok := val.(uint64); !ok || vv != i {
				t.Log(i, ok, val)
			}
		}
		for i := uint32(0); i < 1000; i++ {
			// 编码
			buf := encoding.Encode(i)
			// 解码
			val := encoding.Decode(buf)
			if vv, ok := val.(uint32); !ok || vv != i {
				t.Log(i, ok, val)
			}
		}
		for i := uint16(0); i < 1000; i++ {
			// 编码
			buf := encoding.Encode(i)
			// 解码
			val := encoding.Decode(buf)
			if vv, ok := val.(uint16); !ok || vv != i {
				t.Log(i, ok, val)
			}
		}
		for i := uint8(0); i < 225; i++ {
			// 编码
			buf := encoding.Encode(i)
			// 解码
			val := encoding.Decode(buf)
			if vv, ok := val.(uint8); !ok || vv != i {
				t.Log(i, ok, val)
			}
		}
	})
	t.Run("integer", func(t *testing.T) {
		for i := int64(-1000); i < 1000; i++ {
			// 编码
			buf := encoding.Encode(i)
			// 解码
			val := encoding.Decode(buf)
			if vv, ok := val.(int64); !ok || vv != i {
				t.Log(i, ok, val)
			}
		}
		for i := int32(-1000); i < 1000; i++ {
			// 编码
			buf := encoding.Encode(i)
			// 解码
			val := encoding.Decode(buf)
			if vv, ok := val.(int32); !ok || vv != i {
				t.Log(i, ok, val)
			}
		}
		for i := int16(-1000); i < 1000; i++ {
			// 编码
			buf := encoding.Encode(i)
			// 解码
			val := encoding.Decode(buf)
			if vv, ok := val.(int16); !ok || vv != i {
				t.Log(i, ok, val)
			}
		}
		for i := int8(-100); i < 100; i++ {
			// 编码
			buf := encoding.Encode(i)
			// 解码
			val := encoding.Decode(buf)
			if vv, ok := val.(int8); !ok || vv != i {
				t.Log(i, ok, val)
			}
		}
	})
}
