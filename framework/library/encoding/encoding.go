package encoding

import (
	"fmt"
	"math"
)

/*
* --------------内网跨服务调用api的编码规范--------------
* 00 --int8,int16,int32,int64,uint8,uint16,uint32,uint64,bool
*	00	--表示bool类型
		000x|	--表示true或false
*	01	--字节编码
*	10	--variant编码
		x	--是否有符号
*			000|	--uint8
*			001|	--int8
*			010|	--uint16
*			011|	--int16
*			100|	--uint32
*			101|	--int32
*			110|	--uint64
*			111|	--int64
*	11	--表示float类型
		0000 	--float32
		0001 	--float64
* 01 --[]byte
*	x	--表示下一个字节是否为长度字节。0表示5bit表示长度。1表示有13bit长度(8191)
* 10 --string
*	x	--表示下一个字节是否为长度字节。0表示5bit表示长度。1表示有13bit长度(8191)
* 11 --proto
*	x	--表示下一个字节是否为长度字节。0表示5bit表示长度。1表示有13bit长度(8191)
*/

type WireType uint64

const (
	WireTypeVariant WireType = 0
	WireType16Bit   WireType = 1
	WireType32Bit   WireType = 2
	WireType64Bit   WireType = 3
)

type DataType uint64

const (
	DataTypeBool    DataType = 1
	DataTypeInt8    DataType = 2
	DataTypeUint8   DataType = 3
	DataTypeInt16   DataType = 4
	DataTypeUint16  DataType = 5
	DataTypeInt32   DataType = 6
	DataTypeUint32  DataType = 7
	DataTypeInt64   DataType = 8
	DataTypeUint64  DataType = 9
	DataTypeFloat32 DataType = 10 // 0x01
	DataTypeFloat64 DataType = 11 // 0x0b
	DataTypeBytes   DataType = 12 // 0x0c
	DataTypeString  DataType = 13 // 0x0d
	DataTypeProto   DataType = 14 // 0x0e
)

// Zigzag编码
func ZigzagInt16Encode(n int16) uint16 {
	return uint16((n << 1) ^ (n >> 15))
}

func ZigzagInt32Encode(n int32) uint32 {
	return uint32((n << 1) ^ (n >> 31))
}

func ZigzagInt64Encode(n int64) uint64 {
	return uint64((n << 1) ^ (n >> 63))
}

// Zigzag解码
func ZigzagInt16Decode(n uint16) int16 {
	return int16((n >> 1) ^ -(n & 1))
}

func ZigzagInt32Decode(n uint32) int32 {
	return int32((n >> 1) ^ -(n & 1))
}

func ZigzagInt64Decode(n uint64) int64 {
	return int64((n >> 1) ^ -(n & 1))
}

// Variant编码
func VariantUint64Encode(value uint64) (buf []byte) {
	for {
		b := byte(value & 0x7f)
		value >>= 7
		if value != 0 {
			b |= 0x80
		}
		buf = append(buf, b)
		if value == 0 {
			break
		}
	}
	return
}

// Variant解码
func VariantUint64Decode(data []byte) (value uint64) {
	var shift uint
	for _, b := range data {
		value |= (uint64(b&0x7f) << shift)
		shift += 7
		if b&0x80 == 0 {
			break
		}
	}
	return
}

func BytesToUint64(buf []byte) (val uint64) {
	for i, item := range buf {
		val |= (uint64(item) << (i * 8))
	}
	return
}

func Uint64ToBytes(val uint64) (buf []byte) {
	for ; val != 0; val >>= 8 {
		buf = append(buf, byte(val&0xff))
	}
	return
}

func Uint64Encode(dataType DataType, wireType WireType, val uint64) (buf []byte) {
	buf = append(buf, byte(((dataType&0x0f)<<4)|DataType(wireType&0x0f)))
	switch wireType {
	case WireTypeVariant:
		buf = append(buf, VariantUint64Encode(val)...)
	default:
		buf = append(buf, Uint64ToBytes(val)...)
	}
	return
}

func Uint64Decode(wireType WireType, buf []byte) (val uint64) {
	switch wireType {
	case WireTypeVariant:
		val = VariantUint64Decode(buf)
	default:
		val = BytesToUint64(buf)
	}
	return
}

func Encode(vv interface{}) (buf []byte, err error) {
	switch value := vv.(type) {
	case string:
		/*
			s := *(*reflect.StringHeader)(unsafe.Pointer(&value))
			b := &reflect.SliceHeader{Data: s.Data, Len: s.Len, Cap: s.Len}
			bb := *(*[]byte)(unsafe.Pointer(b))
		*/
	case []byte:
		ll := uint64(len(buf))
		if limit := uint64(0x1<<12 - 1); ll > limit {
			err = fmt.Errorf("[]byte out of length limit(%d)", limit)
			return
		}
		buf = append(buf, byte(uint64(DataTypeBytes&0x0f)<<4|ll>>8), byte(ll&0xff))
		buf = append(buf, value...)
	case bool:
		if value {
			buf = append(buf, byte(DataTypeBool<<4+0x01))
		} else {
			buf = append(buf, byte(DataTypeBool<<4))
		}
	case uint8:
		buf = append(buf, byte(DataTypeUint8<<4), byte(value))
	case int8:
		buf = append(buf, byte(DataTypeInt8<<4), byte(value))
	case float32:
		val := math.Float32bits(value)
		if val>>28 != 0 {
			buf = Uint64Encode(DataTypeFloat32, WireType32Bit, uint64(val))
		} else {
			buf = Uint64Encode(DataTypeFloat32, WireTypeVariant, uint64(val))
		}
	case float64:
		val := math.Float64bits(value)
		if val>>56 != 0 {
			buf = Uint64Encode(DataTypeFloat64, WireType64Bit, val)
		} else {
			buf = Uint64Encode(DataTypeFloat64, WireTypeVariant, val)
		}
	case int16:
		val := uint16((value << 1) ^ (value >> 15))
		if val>>14 != 0 {
			buf = Uint64Encode(DataTypeInt16, WireType16Bit, uint64(val))
		} else {
			buf = Uint64Encode(DataTypeInt16, WireTypeVariant, uint64(val))
		}
	case int32:
		val := uint32((value << 1) ^ (value >> 31))
		if val>>28 != 0 {
			buf = Uint64Encode(DataTypeInt32, WireType32Bit, uint64(val))
		} else {
			buf = Uint64Encode(DataTypeInt32, WireTypeVariant, uint64(val))
		}
	case int64:
		val := uint64((value << 1) ^ (value >> 63))
		if val>>56 != 0 {
			buf = Uint64Encode(DataTypeInt64, WireType64Bit, val)
		} else {
			buf = Uint64Encode(DataTypeInt64, WireTypeVariant, val)
		}
	case uint16:
		if value>>14 != 0 {
			buf = Uint64Encode(DataTypeUint16, WireType16Bit, uint64(value))
		} else {
			buf = Uint64Encode(DataTypeUint16, WireTypeVariant, uint64(value))
		}
	case uint32:
		if value>>28 != 0 {
			buf = Uint64Encode(DataTypeUint32, WireType32Bit, uint64(value))
		} else {
			buf = Uint64Encode(DataTypeUint32, WireTypeVariant, uint64(value))
		}
	case uint64:
		if value>>56 != 0 {
			buf = Uint64Encode(DataTypeUint64, WireType64Bit, value)
		} else {
			buf = Uint64Encode(DataTypeUint64, WireTypeVariant, value)
		}
	}
	return
}

func Decode(buf []byte) interface{} {
	wireType := WireType(buf[0] & 0x0f)
	switch DataType((buf[0] >> 4) & 0x0f) {
	case DataTypeBytes:
		ll := (buf[0]&0x0f)<<8 | buf[1]
		return buf[2 : ll+2]
	case DataTypeBool:
		return buf[0]&0x01 == 1
	case DataTypeInt8:
		return int8(buf[1])
	case DataTypeUint8:
		return uint8(buf[1])
	case DataTypeFloat32:
		return math.Float32frombits(uint32(Uint64Decode(wireType, buf[1:])))
	case DataTypeFloat64:
		return math.Float64frombits(Uint64Decode(wireType, buf[1:]))
	case DataTypeInt16:
		val := uint16(Uint64Decode(wireType, buf[1:]))
		return int16((val >> 1) ^ -(val & 1))
	case DataTypeInt32:
		val := uint32(Uint64Decode(wireType, buf[1:]))
		return int32((val >> 1) ^ -(val & 1))
	case DataTypeInt64:
		val := Uint64Decode(wireType, buf[1:])
		return int64((val >> 1) ^ -(val & 1))
	case DataTypeUint16:
		return uint16(Uint64Decode(wireType, buf[1:]))
	case DataTypeUint32:
		return uint32(Uint64Decode(wireType, buf[1:]))
	case DataTypeUint64:
		return Uint64Decode(wireType, buf[1:])
	}
	return nil
}
