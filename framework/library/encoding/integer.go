package encoding

func Uint16Encode(value uint16) (buf []byte) {
	if value>>14 != 0 {
		buf = append(buf, byte(((DataTypeUint16<<4)&0xf0)|(DataType(WireTypeVariant)&0x0f)))
		buf = append(buf, byte(value&0xff), byte((value>>8)&0xff))
	} else {
		buf = append(buf, byte(((DataTypeUint16<<4)&0xf0)|(DataType(WireType16Bit)&0x0f)))
		buf = append(buf, VariantUint64Encode(uint64(value))...)
	}
	return
}

func Int16Encode(val int16) (buf []byte) {
	var value uint16
	if value <= 0 {
		value = ZigzagInt16Encode(val)
	} else {
		value = uint16(val)
	}

	// 编码
	if value>>14 != 0 {
		buf = append(buf, byte(((DataTypeInt16<<4)&0xf0)|(DataType(WireTypeVariant)&0x0f)))
		buf = append(buf, byte(value&0xff), byte((value>>8)&0xff))
	} else {
		buf = append(buf, byte(((DataTypeInt16<<4)&0xf0)|(DataType(WireType16Bit)&0x0f)))
		buf = append(buf, VariantUint64Encode(uint64(value))...)
	}
	return
}

func Uint32Encode(value uint32) (buf []byte) {
	if value>>28 != 0 {
		buf = append(buf, byte(((DataTypeUint32<<4)&0xf0)|(DataType(WireTypeVariant)&0x0f)))
		buf = append(buf, byte(value&0xff), byte((value>>8)&0xff), byte((value>>16)&0xff), byte((value>>24)&0xff))
	} else {
		buf = append(buf, byte(((DataTypeUint32<<4)&0xf0)|(DataType(WireType32Bit)&0x0f)))
		buf = append(buf, VariantUint64Encode(uint64(value))...)
	}
	return
}

func Int32Encode(val int32) (buf []byte) {
	var value uint32
	if value <= 0 {
		value = ZigzagInt32Encode(val)
	} else {
		value = uint32(val)
	}

	// 编码
	if value>>28 != 0 {
		buf = append(buf, byte(((DataTypeInt32<<4)&0xf0)|(DataType(WireTypeVariant)&0x0f)))
		buf = append(buf, byte(value&0xff), byte((value>>8)&0xff), byte((value>>16)&0xff), byte((value>>24)&0xff))
	} else {
		buf = append(buf, byte(((DataTypeInt32<<4)&0xf0)|(DataType(WireType32Bit)&0x0f)))
		buf = append(buf, VariantUint64Encode(uint64(value))...)
	}
	return
}

func Uint64Encode(value uint64) (buf []byte) {
	if value>>56 != 0 {
		buf = append(buf, byte(((DataTypeUint64<<4)&0xf0)|(DataType(WireTypeVariant)&0x0f)))
		buf = append(buf, byte(value&0xff), byte((value>>8)&0xff), byte((value>>16)&0xff), byte((value>>24)&0xff), byte((value>>32)&0xff), byte((value>>40)&0xff), byte((value>>48)&0xff), byte((value>>56)&0xff))
	} else {
		buf = append(buf, byte(((DataTypeUint64<<4)&0xf0)|(DataType(WireType64Bit)&0x0f)))
		buf = append(buf, VariantUint64Encode(uint64(value))...)
	}
	return
}

func Int64Encode(val int64) (buf []byte) {
	var value uint64
	if value <= 0 {
		value = ZigzagInt64Encode(val)
	} else {
		value = uint64(val)
	}

	// 编码
	if value>>56 != 0 {
		buf = append(buf, byte(((DataTypeInt64<<4)&0xf0)|(DataType(WireTypeVariant)&0x0f)))
		buf = append(buf, byte(value&0xff), byte((value>>8)&0xff), byte((value>>16)&0xff), byte((value>>24)&0xff), byte((value>>32)&0xff), byte((value>>40)&0xff), byte((value>>48)&0xff), byte((value>>56)&0xff))
	} else {
		buf = append(buf, byte(((DataTypeInt64<<4)&0xf0)|(DataType(WireType64Bit)&0x0f)))
		buf = append(buf, VariantUint64Encode(uint64(value))...)
	}
	return
}
