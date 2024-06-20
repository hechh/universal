package test

import (
	"testing"
	"universal/framework/library/encoding"
)

func TestEncoding(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		// true测试
		if buf, err := encoding.Encode(true); err != nil {
			t.Log("bool", err)
		} else {
			results, num := encoding.Decode(buf)
			if vv, ok := results.(bool); !ok || vv != true {
				t.Log(true, vv, num)
			}
		}

		// false测试
		if buf, err := encoding.Encode(false); err != nil {
			t.Log("bool", err)
		} else {
			result, num := encoding.Decode(buf)
			if vv, ok := result.(bool); !ok || vv != false {
				t.Log(false, vv, num)
			}
		}
	})
	t.Run("int16", func(t *testing.T) {
		value := int16(32767)
		if buf, err := encoding.Encode(value); err != nil {
			t.Log("int16", err)
		} else {
			result, num := encoding.Decode(buf)
			if vv, ok := result.(int16); !ok || vv != value {
				t.Log(value, vv, num, "---->", buf)
			}
		}
		value = int16(127)
		if buf, err := encoding.Encode(value); err != nil {
			t.Log("int16", err)
		} else {
			result, num := encoding.Decode(buf)
			if vv, ok := result.(int16); !ok || vv != value {
				t.Log(value, vv, num, "---->", buf)
			}
		}
	})
}
