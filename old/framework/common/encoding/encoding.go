package encoding

import (
	"fmt"
	"hash/crc32"
	"math"
	"reflect"
	"strings"
	"unsafe"

	"github.com/golang/protobuf/proto"
)

/*
* --------------内网跨服务调用api的编码规范--------------
* 00 --int8,int16,int32,int64,uint8,uint16,uint32,uint64,bool
*	00	--表示无编码
*		0x00(0000)|	--false
*		0x01(0001)|	--true
*		0x02(0010)|	--uint8
*		0x03(0011)|	--int8
*	01	--字节编码
*		0000|	--uint16	<0x10(0001 0000)>
*		0001|	--int16		<0x11(0001 0001)>
*		0010|	--uint32	<0x12(0001 0010)>
*		0011|	--int32		<0x13(0001 0011)>
*		0100|	--uint64	<0x14(0001 0100)>
*		0101|	--int64		<0x15(0001 0101)>
*		0110| 	--float32	<0x16(0011 0110)>
*		0111| 	--float64	<0x17(0011 0111)>
*	10	--variant编码
*		0000|	--uint16	<0x20(0010 0000)>
*		0001|	--int16		<0x21(0010 0001)>
*		0010|	--uint32	<0x22(0010 0010)>
*		0011|	--int32		<0x23(0010 0011)>
*		0100|	--uint64	<0x24(0010 0100)>
*		0101|	--int64		<0x25(0010 0101)>
*		0110| 	--float32	<0x26(0010 0110)>
*		0111| 	--float64	<0x27(0010 0111)>
* 01 --[]byte,string
*	010x xxxx				--0表示5bit表示长度(32)		<0x40|0xHH>
*	011x xxxx|xxxx xxxx		--1表示有13bit长度(8191) 	<0x60|0xHH>
* 10 --string
*	100x xxxx				--0表示5bit表示长度(32)		<0x80|0xHH>
*	101x xxxx|xxxx xxxx		--1表示有13bit长度(8191) 	<0xa0|0xHH>
* 11 --proto
*	110x xxxx|--crc32--|				--0表示5bit表示长度(32)		<0xC0|0xHH>
*	11xx xxxx|xxxx xxxx|--crc32--|		--0表示13bit表示长度(8191)	<0xE0|0xHH>
 */

const (
	DATA_TYPE_MASK  = 3 << 6
	DataTypeIdent   = 0x00
	DataTypeBytes   = 0x40
	DataTypeString  = 0x80
	DataTypeProto   = 0xC0
	WIRE_TYPE_MASK  = 3 << 4
	WireTypeNone    = 0x00
	WireTypeByte    = 0x10
	WireTypeVariant = 0x20
	SIZE_FLAG_MASK  = (1 << 5)
	DATA_SIZE_MASK  = 1<<5 - 1
)

var (
	protos = make(map[uint32]reflect.Type)
)

func RegisterProto(aa proto.Message) {
	name := getProtoName(aa)
	crc := crc32.ChecksumIEEE([]byte(name))
	protos[crc] = reflect.TypeOf(aa).Elem()
}

func getProtoName(packet proto.Message) string {
	sType := proto.MessageName(packet)
	index := strings.Index(sType, ".")
	if index != -1 {
		sType = sType[index+1:]
	}
	return sType
}

func BytesToUint64(buf []byte, pos int, num int) (val uint64) {
	for i := 0; i < num; i++ {
		val |= (uint64(buf[i+pos]) << (i * 8))
	}
	return
}

func Uint64ToBytes(val uint64) (buf []byte) {
	for ; val != 0; val >>= 8 {
		buf = append(buf, byte(val&0xff))
	}
	return
}

// Zigzag编码
func ZigzagInt64Encode(n int64) uint64 {
	return uint64((n << 1) ^ (n >> 63))
}

// Zigzag解码
func ZigzagInt64Decode(n uint64) int64 {
	return int64((n >> 1) ^ -(n & 1))
}

// Variant编码
func VariantEncode(value uint64) (buf []byte) {
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
func VariantDecode(data []byte) (value uint64, num int) {
	var shift uint
	for _, b := range data {
		num++
		value |= (uint64(b&0x7f) << shift)
		shift += 7
		if b&0x80 == 0 {
			break
		}
	}
	return
}

func Encode(vv interface{}) (buf []byte, err error) {
	switch value := vv.(type) {
	case bool:
		if !value {
			buf = append(buf, 0x00)
		} else {
			buf = append(buf, 0x01)
		}
	case uint8:
		buf = append(buf, 0x02, value)
	case int8:
		buf = append(buf, 0x03, byte(value))
	case uint16:
		if value>>14 != 0 {
			buf = append(buf, 0x10)
			buf = append(buf, Uint64ToBytes(uint64(value))...)
		} else {
			buf = append(buf, 0x20)
			buf = append(buf, VariantEncode(uint64(value))...)
		}
	case int16:
		if value>>14 != 0 {
			buf = append(buf, 0x11)
			buf = append(buf, Uint64ToBytes(uint64(value))...)
		} else {
			val := uint16((value << 1) ^ (value >> 15))
			buf = append(buf, 0x21)
			buf = append(buf, VariantEncode(uint64(val))...)
		}
	case uint32:
		if value>>28 != 0 {
			buf = append(buf, 0x12)
			buf = append(buf, Uint64ToBytes(uint64(value))...)
		} else {
			buf = append(buf, 0x22)
			buf = append(buf, VariantEncode(uint64(value))...)
		}
	case int32:
		if value>>28 != 0 {
			buf = append(buf, 0x13)
			buf = append(buf, Uint64ToBytes(uint64(value))...)
		} else {
			val := uint32((value << 1) ^ (value >> 31))
			buf = append(buf, 0x23)
			buf = append(buf, VariantEncode(uint64(val))...)
		}
	case uint64:
		if value>>56 != 0 {
			buf = append(buf, 0x14)
			buf = append(buf, Uint64ToBytes(value)...)
		} else {
			buf = append(buf, 0x24)
			buf = append(buf, VariantEncode(value)...)
		}
	case int64:
		if value>>56 != 0 {
			buf = append(buf, 0x15)
			buf = append(buf, Uint64ToBytes(uint64(value))...)
		} else {
			val := uint64((value << 1) ^ (value >> 63))
			buf = append(buf, 0x25)
			buf = append(buf, VariantEncode(val)...)
		}
	case float32:
		if val := math.Float32bits(value); val>>28 != 0 {
			buf = append(buf, 0x16)
			buf = append(buf, Uint64ToBytes(uint64(val))...)
		} else {
			buf = append(buf, 0x26)
			buf = append(buf, VariantEncode(uint64(val))...)
		}
	case float64:
		if val := math.Float64bits(value); val>>56 != 0 {
			buf = append(buf, 0x17)
			buf = append(buf, Uint64ToBytes(val)...)
		} else {
			buf = append(buf, 0x27)
			buf = append(buf, VariantEncode(val)...)
		}
	case []byte:
		if ll := uint64(len(value)); ll > 0x1fff {
			return nil, fmt.Errorf("[]byte out of length limit(%d)", 0x1fff)
		} else if ll <= 0x1f {
			buf = append(buf, 0x40|uint8(ll))
		} else {
			buf = append(buf, 0x60|uint8(ll>>8), uint8(ll))
		}
		buf = append(buf, value...)
	case string:
		if ll := uint64(len(value)); ll > 0x1fff {
			return nil, fmt.Errorf("string out of length limit(%d)", 0x1fff)
		} else if ll <= 0x1f {
			buf = append(buf, 0x80|uint8(ll))
		} else {
			buf = append(buf, 0xa0|uint8(ll>>8), uint8(ll))
		}
		buf = append(buf, value...)
	case proto.Message:
		bb, err := proto.Marshal(value)
		if err != nil {
			return nil, err
		}
		if ll := uint64(len(bb)); ll > 0x1fff {
			return nil, fmt.Errorf("proto.Marshal out of length limit(%d)", 0x1fff)
		} else if ll <= 0x1f {
			buf = append(buf, 0xC0|uint8(ll))
		} else {
			buf = append(buf, 0xE0|uint8(ll>>8), uint8(ll))
		}
		crcVal := crc32.ChecksumIEEE([]byte(getProtoName(value)))
		buf = append(buf, Uint64ToBytes(uint64(crcVal))...)
		buf = append(buf, bb...)
	}
	return
}

func Decode(buf []byte) (ret interface{}, shift int, err error) {
	switch buf[0] & DATA_TYPE_MASK {
	case DataTypeIdent:
		switch buf[0] & WIRE_TYPE_MASK {
		case WireTypeNone:
			switch buf[0] & 0x0f {
			case 0x00:
				ret = false
				shift = 1
			case 0x01:
				ret = true
				shift = 1
			case 0x02:
				ret = uint8(buf[1])
				shift = 2
			case 0x03:
				ret = int8(buf[1])
				shift = 2
			}
		case WireTypeByte:
			switch buf[0] & 0x0f {
			case 0x00:
				ret = uint16(BytesToUint64(buf, 1, 2))
				shift = 3
			case 0x01:
				ret = int16(BytesToUint64(buf, 1, 2))
				shift = 3
			case 0x02:
				ret = uint32(BytesToUint64(buf, 1, 4))
				shift = 5
			case 0x03:
				ret = int32(BytesToUint64(buf, 1, 4))
				shift = 5
			case 0x04:
				ret = BytesToUint64(buf, 1, 8)
				shift = 9
			case 0x05:
				ret = int64(BytesToUint64(buf, 1, 8))
				shift = 9
			case 0x06:
				val := uint32(BytesToUint64(buf, 1, 4))
				ret = math.Float32frombits(val)
				shift = 5
			case 0x07:
				val := BytesToUint64(buf, 1, 8)
				ret = math.Float64frombits(val)
				shift = 9
			}
		case WireTypeVariant:
			switch buf[0] & 0x0f {
			case 0x00:
				val, num := VariantDecode(buf[1:])
				ret = uint16(val)
				shift = num + 1
			case 0x01:
				val, num := VariantDecode(buf[1:])
				ret = int16((val >> 1) ^ -(val & 1))
				shift = num + 1
			case 0x02:
				val, num := VariantDecode(buf[1:])
				ret = uint32(val)
				shift = num + 1
			case 0x03:
				val, num := VariantDecode(buf[1:])
				ret = int32((uint32(val) >> 1) ^ -(uint32(val) & 1))
				shift = num + 1
			case 0x04:
				val, num := VariantDecode(buf[1:])
				ret = val
				shift = num + 1
			case 0x05:
				val, num := VariantDecode(buf[1:])
				ret = int64((val >> 1) ^ -(val & 1))
				shift = num + 1
			case 0x06:
				val, num := VariantDecode(buf[1:])
				ret = math.Float32frombits(uint32(val))
				shift = num + 1
			case 0x07:
				val, num := VariantDecode(buf[1:])
				ret = math.Float64frombits(val)
				shift = num + 1
			}
		}
	case DataTypeBytes:
		shift = 1
		ll := int(buf[0] & DATA_SIZE_MASK)
		if SIZE_FLAG_MASK&buf[0] != 0 {
			ll = (ll << 8) | int(buf[1])
			shift++
		}
		result := make([]byte, ll)
		copy(result, buf[shift:shift+ll])
		ret = result
		shift += ll
	case DataTypeString:
		shift = 1
		ll := int(buf[0] & DATA_SIZE_MASK)
		if SIZE_FLAG_MASK&buf[0] != 0 {
			ll = (ll << 8) | int(buf[1])
			shift++
		}
		result := make([]byte, ll)
		copy(result, buf[shift:shift+ll])
		b := *(*reflect.SliceHeader)(unsafe.Pointer(&result))
		s := &reflect.StringHeader{Data: b.Data, Len: b.Len}
		ret = *(*string)(unsafe.Pointer(s))
		shift += ll

	case DataTypeProto:
		shift = 5
		ll := int(buf[0] & DATA_SIZE_MASK)
		if SIZE_FLAG_MASK&buf[0] != 0 {
			ll = (ll << 8) | int(buf[1])
			shift++
		}
		// 获取crc
		crc := uint32(BytesToUint64(buf, shift-4, 4))
		tt, ok := protos[crc]
		if !ok {
			err = fmt.Errorf("crc32(%d) not supported", crc)
			return
		}
		// 创建对象
		result := reflect.New(tt).Interface().(proto.Message)
		if err = proto.Unmarshal(buf[shift:shift+ll], result); err != nil {
			return
		}
		ret = result
		shift += ll
	}
	return
}
