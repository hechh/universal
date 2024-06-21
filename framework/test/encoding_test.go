package test

import (
	"testing"
	"universal/common/pb"
	"universal/framework/library/encoding"
)

func TestMain(m *testing.M) {
	encoding.RegisterProto(&pb.RpcHead{})
	m.Run()
}

func TestEncoding(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		// true测试
		if buf, err := encoding.Encode(true); err != nil {
			t.Log("bool", err)
		} else {
			results, num, _ := encoding.Decode(buf)
			if vv, ok := results.(bool); !ok || vv != true {
				t.Log(true, vv, num)
			}
			t.Log("----pass---->", true, "num: ", num)
		}

		// false测试
		if buf, err := encoding.Encode(false); err != nil {
			t.Log("bool", err)
		} else {
			result, num, _ := encoding.Decode(buf)
			if vv, ok := result.(bool); !ok || vv != false {
				t.Log(false, vv, num)
			}
			t.Log("----pass---->", false, "num: ", num)
		}
	})
	t.Run("int8", func(t *testing.T) {
		value := int8(-127)
		if buf, err := encoding.Encode(value); err != nil {
			t.Log("int8", err)
		} else {
			result, num, _ := encoding.Decode(buf)
			if vv, ok := result.(int8); !ok || vv != value {
				t.Log(value, vv, num, "---->", buf)
			}
			t.Log("----pass---->", value, "num: ", num)
		}
		value = int8(127)
		if buf, err := encoding.Encode(value); err != nil {
			t.Log("int8", err)
		} else {
			result, num, _ := encoding.Decode(buf)
			if vv, ok := result.(int8); !ok || vv != value {
				t.Log(value, vv, num, "---->", buf)
			}
			t.Log("----pass---->", value, "num: ", num)
		}
	})
	t.Run("int16", func(t *testing.T) {
		value := int16(32767)
		if buf, err := encoding.Encode(value); err != nil {
			t.Log("int16", err)
		} else {
			result, num, _ := encoding.Decode(buf)
			if vv, ok := result.(int16); !ok || vv != value {
				t.Log(value, vv, num, "---->", buf)
			}
			t.Log("----pass---->", value, "num: ", num)
		}
		value = int16(127)
		if buf, err := encoding.Encode(value); err != nil {
			t.Log("int16", err)
		} else {
			result, num, _ := encoding.Decode(buf)
			if vv, ok := result.(int16); !ok || vv != value {
				t.Log(value, vv, num, "---->", buf)
			}
			t.Log("----pass---->", value, "num: ", num)
		}
	})
	t.Run("int32", func(t *testing.T) {
		value := int32(2147483647)
		if buf, err := encoding.Encode(value); err != nil {
			t.Log("int16", err)
		} else {
			result, num, _ := encoding.Decode(buf)
			if vv, ok := result.(int32); !ok || vv != value {
				t.Log(value, vv, num, "---->", buf)
			}
			t.Log("----pass---->", value, "num: ", num)
		}
		value = int32(127)
		if buf, err := encoding.Encode(value); err != nil {
			t.Log("int16", err)
		} else {
			result, num, _ := encoding.Decode(buf)
			if vv, ok := result.(int32); !ok || vv != value {
				t.Log(value, vv, num, "---->", buf)
			}
			t.Log("----pass---->", value, "num: ", num)
		}
	})
	t.Run("int64", func(t *testing.T) {
		value := int64(0x7fffffffffffffff)
		if buf, err := encoding.Encode(value); err != nil {
			t.Log("int16", err)
		} else {
			result, num, _ := encoding.Decode(buf)
			if vv, ok := result.(int64); !ok || vv != value {
				t.Log(value, vv, num, "---->", buf)
			}
			t.Log("----pass---->", value, "num: ", num)
		}
		value = int64(127)
		if buf, err := encoding.Encode(value); err != nil {
			t.Log("int16", err)
		} else {
			result, num, _ := encoding.Decode(buf)
			if vv, ok := result.(int64); !ok || vv != value {
				t.Log(value, vv, num, "---->", buf)
			}
			t.Log("----pass---->", value, "num: ", num)
		}
	})
	t.Run("[]byte", func(t *testing.T) {
		value := []byte{1, 2, 3, 4, 5, 6, 7}
		if buf, err := encoding.Encode(value); err != nil {
			t.Log("int16", err)
			return
		} else {
			result, num, _ := encoding.Decode(buf)
			if vv, ok := result.([]byte); !ok {
				t.Log(value, vv, num, "---->", buf)
				return
			} else {
				for i, elem := range value {
					if elem != vv[i] {
						t.Log(value, vv, num, "---->", buf)
						return
					}
				}
			}
			t.Log("----pass---->", value, "num: ", num)
		}

		value = []byte{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		}
		if buf, err := encoding.Encode(value); err != nil {
			t.Log("int16", err)
			return
		} else {
			result, num, _ := encoding.Decode(buf)
			if vv, ok := result.([]byte); !ok {
				t.Log(value, vv, num, "---->", buf)
				return
			} else {
				for i, elem := range value {
					if elem != vv[i] {
						t.Log(value, vv, num, "---->", buf)
						return
					}
				}
			}
			t.Log("----pass---->", value, "num: ", num)
		}
	})
	t.Run("string", func(t *testing.T) {
		value := "this is a"
		if buf, err := encoding.Encode(value); err != nil {
			t.Log("int16", err)
		} else {
			result, num, _ := encoding.Decode(buf)
			if vv, ok := result.(string); !ok || len(value) != len(vv) {
				t.Log(value, vv, num, "---->", buf)
				return
			} else {
				for i, elem := range value {
					if elem != rune(vv[i]) {
						t.Log(value, vv, num, "---->", buf)
						return
					}
				}
			}
			t.Log("----pass---->", value, "num: ", num)
		}

		value = "asdfaksdfja;ksjdf;laksjdf;akjdf;lkajsd;fkja;skdfj;alksjdfaksjdfa;lkjdf;akjdf;alksjdf;kajsd;fkja;sdfj;askdjf;aksjdf;a"
		if buf, err := encoding.Encode(value); err != nil {
			t.Log("int16", err)
			return
		} else {
			result, num, _ := encoding.Decode(buf)
			if vv, ok := result.(string); !ok || len(value) != len(vv) {
				t.Log(value, vv, num, "---->", buf)
			} else {
				for i, elem := range value {
					if elem != rune(vv[i]) {
						t.Log(value, vv, num, "---->", buf)
						return
					}
				}
			}
			t.Log("----pass---->", value, "num: ", num)
		}
	})
	t.Run("proto", func(t *testing.T) {
		value := &pb.RpcHead{Code: 123234123, ErrMsg: "this"}
		if buf, err := encoding.Encode(value); err != nil {
			t.Log("int16", err)
		} else {
			result, num, _ := encoding.Decode(buf)
			if vv, ok := result.(*pb.RpcHead); !ok || vv == nil {
				t.Log(value, vv, num, "---->", buf)
				return
			} else {
				if value.Code != vv.Code || value.ErrMsg != vv.ErrMsg {
					t.Log(value, vv, num, "---->", buf)
					return
				}
			}
			t.Log("----pass---->", value, "num: ", num)
		}

		value = &pb.RpcHead{Code: 123234123, ErrMsg: "this is a testsfeasdfasdfasdfasdfasdfadfadfadfadfadfadfasdfasdfasdfasdfasdfasdfasdf"}
		if buf, err := encoding.Encode(value); err != nil {
			t.Log("int16", err)
		} else {
			result, num, _ := encoding.Decode(buf)
			if vv, ok := result.(*pb.RpcHead); !ok || vv == nil {
				t.Log(value, vv, num, "---->", buf)
				return
			} else {
				if value.Code != vv.Code || value.ErrMsg != vv.ErrMsg {
					t.Log(value, vv, num, "---->", buf)
					return
				}
			}
			t.Log("----pass---->", value, "num: ", num)
		}
	})
}
