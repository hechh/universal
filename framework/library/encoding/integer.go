package encoding

/*
* |dddd 0000|--1byte--| 		int8、uint8类型：无编码类型。
* |dddd xx00|--<=2byte--|		int16、uint16类型：当val>>14 等于0时, 采用variant编码，否则采用16bit编码
* |dddd xx00|--<=4byte--|		int32、uint32类型：当val>>28 等于0时，采用variant编码，否则采用32bit编码
* |dddd xx00|--<=8byte--|		int64、uint64类型：当val>>56 等于0时，采用variant编码，否则采用64bit编码
 */

func IntegerEncode(vv interface{}) (buf []byte) {
	switch value := vv.(type) {
	case int16:
		val := uint16((value << 1) ^ (value >> 15))
		if val>>14 != 0 {
			buf = append(buf, byte(((DataTypeInt16<<4)&0xf0)|(DataType(WireTypeVariant)&0x0f)))
			buf = append(buf, byte(val&0xff), byte((val>>8)&0xff))
		} else {
			buf = append(buf, byte(((DataTypeInt16<<4)&0xf0)|(DataType(WireType16Bit)&0x0f)))
			buf = append(buf, VariantUint64Encode(uint64(val))...)
		}
	case int32:
		val := uint32((value << 1) ^ (value >> 31))
		if val>>28 != 0 {
			buf = append(buf, byte(((DataTypeInt32<<4)&0xf0)|(DataType(WireTypeVariant)&0x0f)))
			buf = append(buf, byte(val&0xff), byte((val>>8)&0xff), byte((val>>16)&0xff), byte((val>>24)&0xff))
		} else {
			buf = append(buf, byte(((DataTypeInt32<<4)&0xf0)|(DataType(WireType32Bit)&0x0f)))
			buf = append(buf, VariantUint64Encode(uint64(val))...)
		}
	case int64:
		val := uint64((value << 1) ^ (value >> 63))
		if val>>56 != 0 {
			buf = append(buf, byte(((DataTypeInt64<<4)&0xf0)|(DataType(WireTypeVariant)&0x0f)))
			buf = append(buf, byte(val&0xff), byte((val>>8)&0xff), byte((val>>16)&0xff), byte((val>>24)&0xff), byte((val>>32)&0xff), byte((val>>40)&0xff), byte((val>>48)&0xff), byte((val>>56)&0xff))
		} else {
			buf = append(buf, byte(((DataTypeInt64<<4)&0xf0)|(DataType(WireType64Bit)&0x0f)))
			buf = append(buf, VariantUint64Encode(uint64(val))...)
		}
	case uint16:
		if value>>14 != 0 {
			buf = append(buf, byte(((DataTypeUint16<<4)&0xf0)|(DataType(WireTypeVariant)&0x0f)))
			buf = append(buf, byte(value&0xff), byte((value>>8)&0xff))
		} else {
			buf = append(buf, byte(((DataTypeUint16<<4)&0xf0)|(DataType(WireType16Bit)&0x0f)))
			buf = append(buf, VariantUint64Encode(uint64(value))...)
		}
	case uint32:
		if value>>28 != 0 {
			buf = append(buf, byte(((DataTypeUint32<<4)&0xf0)|(DataType(WireTypeVariant)&0x0f)))
			buf = append(buf, byte(value&0xff), byte((value>>8)&0xff), byte((value>>16)&0xff), byte((value>>24)&0xff))
		} else {
			buf = append(buf, byte(((DataTypeUint32<<4)&0xf0)|(DataType(WireType32Bit)&0x0f)))
			buf = append(buf, VariantUint64Encode(uint64(value))...)
		}
	case uint64:
		if value>>56 != 0 {
			buf = append(buf, byte(((DataTypeUint64<<4)&0xf0)|(DataType(WireTypeVariant)&0x0f)))
			buf = append(buf, byte(value&0xff), byte((value>>8)&0xff), byte((value>>16)&0xff), byte((value>>24)&0xff), byte((value>>32)&0xff), byte((value>>40)&0xff), byte((value>>48)&0xff), byte((value>>56)&0xff))
		} else {
			buf = append(buf, byte(((DataTypeUint64<<4)&0xf0)|(DataType(WireType64Bit)&0x0f)))
			buf = append(buf, VariantUint64Encode(uint64(value))...)
		}
	}
	return
}

func IntegerDecode(buf []byte) interface{} {
	wireType := WireType((buf[0] >> 2) & 0x01)
	switch DataType(buf[0] >> 4) {
	case DataTypeInt16:
		if wireType == WireTypeVariant {
			val := uint16(VariantUint64Decode(buf[1:]))
			return int16((val >> 1) ^ -(val & 1))
		}
		return int16(buf[1] | (buf[2] << 8))
	case DataTypeInt32:
		if wireType == WireTypeVariant {
			val := uint32(VariantUint64Decode(buf[1:]))
			return int32((val >> 1) ^ -(val & 1))
		}
		return int32(buf[1] | (buf[2] << 8) | (buf[3] << 16) | (buf[4] << 24))
	case DataTypeInt64:
		if wireType == WireTypeVariant {
			val := uint64(VariantUint64Decode(buf[1:]))
			return int64((val >> 1) ^ -(val & 1))
		}
		return int64(buf[1] | (buf[2] << 8) | (buf[3] << 16) | (buf[4] << 24) | (buf[5] << 32) | (buf[6] << 40) | (buf[7] << 56))
	case DataTypeUint16:
		if wireType == WireTypeVariant {
			return uint16(VariantUint64Decode(buf[1:]))
		}
		return uint16(buf[1] | (buf[2] << 8))
	case DataTypeUint32:
		if wireType == WireTypeVariant {
			return uint32(VariantUint64Decode(buf[1:]))
		}
		return uint32(buf[1] | (buf[2] << 8) | (buf[3] << 16) | (buf[4] << 24))
	case DataTypeUint64:
		if wireType == WireTypeVariant {
			return uint64(VariantUint64Decode(buf[1:]))
		}
		return uint64(buf[1] | (buf[2] << 8) | (buf[3] << 16) | (buf[4] << 24) | (buf[5] << 32) | (buf[6] << 40) | (buf[7] << 56))
	}
	return nil
}
