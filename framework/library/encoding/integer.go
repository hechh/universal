package encoding

import "math"

/*
* |0001 0000| 					bool类型(1 byte): 无编码类型，第一位bit表示bool值
* |dddd 0000|--1byte--| 		int8、uint8类型：无编码类型。
* |dddd 00xx|--<=2byte--|		int16、uint16类型：当val>>14 等于0时, 采用variant编码，否则采用16bit编码
* |dddd 00xx|--<=4byte--|		int32、uint32类型：当val>>28 等于0时，采用variant编码，否则采用32bit编码
* |dddd 00xx|--<=8byte--|		int64、uint64类型：当val>>56 等于0时，采用variant编码，否则采用64bit编码
* |1010 00xx|--4byte--|			float32类型：固定编码类型，直接采用32bit编码
* |1011 00xx|--8byte--|			float64类型：固定编码类型，直接采用64bit编码
 */

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

func Encode(vv interface{}) (buf []byte) {
	switch value := vv.(type) {
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
