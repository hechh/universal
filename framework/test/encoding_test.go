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
	t.Run("int8", func(t *testing.T) {
		value := int8(-1)
		buf := encoding.Encode(value)
		val := encoding.Decode(buf)
		if vv, ok := val.(int8); !ok || vv != value {
			t.Log(value, ok, val)
		}

		value02 := uint8(0xff)
		buf02 := encoding.Encode(value02)
		val02 := encoding.Decode(buf02)
		if vv, ok := val02.(uint8); !ok || vv != value02 {
			t.Log(value02, ok, val02)
		}
	})
	t.Run("int16", func(t *testing.T) {
		value := int16(-1)
		buf := encoding.Encode(value)
		val := encoding.Decode(buf)
		if vv, ok := val.(int16); !ok || vv != value {
			t.Log(value, ok, val)
		}
		value02 := int16(0x09)
		buf02 := encoding.Encode(value02)
		val02 := encoding.Decode(buf02)
		if vv, ok := val02.(int16); !ok || vv != value02 {
			t.Log(value02, ok, val02)
		}

		value03 := uint16(12)
		buf03 := encoding.Encode(value03)
		val03 := encoding.Decode(buf03)
		if vv, ok := val03.(uint16); !ok || vv != value03 {
			t.Log(value03, ok, val03)
		}

		value04 := uint16(0x09)
		buf04 := encoding.Encode(value04)
		val04 := encoding.Decode(buf04)
		if vv, ok := val04.(uint16); !ok || vv != value04 {
			t.Log(value04, ok, val04)
		}
	})
	t.Run("int32", func(t *testing.T) {
		value := int32(-1)
		buf := encoding.Encode(value)
		val := encoding.Decode(buf)
		if vv, ok := val.(int32); !ok || vv != value {
			t.Log(value, ok, val)
		}
		value02 := int32(0x09)
		buf02 := encoding.Encode(value02)
		val02 := encoding.Decode(buf02)
		if vv, ok := val02.(int32); !ok || vv != value02 {
			t.Log(value02, ok, val02)
		}

		value03 := uint32(12)
		buf03 := encoding.Encode(value03)
		val03 := encoding.Decode(buf03)
		if vv, ok := val03.(uint32); !ok || vv != value03 {
			t.Log(value03, ok, val03)
		}

		value04 := uint32(0x09)
		buf04 := encoding.Encode(value04)
		val04 := encoding.Decode(buf04)
		if vv, ok := val04.(uint32); !ok || vv != value04 {
			t.Log(value04, ok, val04)
		}
	})
	t.Run("int64", func(t *testing.T) {
		value := int64(-1)
		buf := encoding.Encode(value)
		val := encoding.Decode(buf)
		if vv, ok := val.(int64); !ok || vv != value {
			t.Log(value, ok, val)
		}
		value02 := int64(0xfffffffffffffff)
		buf02 := encoding.Encode(value02)
		val02 := encoding.Decode(buf02)
		if vv, ok := val02.(int64); !ok || vv != value02 {
			t.Log(value02, ok, val02)
		}

		value03 := uint64(0xffffffffffffffff)
		buf03 := encoding.Encode(value03)
		val03 := encoding.Decode(buf03)
		if vv, ok := val03.(uint64); !ok || vv != value03 {
			t.Log(value03, ok, val03)
		}

		value04 := uint64(0x09)
		buf04 := encoding.Encode(value04)
		val04 := encoding.Decode(buf04)
		if vv, ok := val04.(uint64); !ok || vv != value04 {
			t.Log(value04, ok, val04)
		}
	})

}
