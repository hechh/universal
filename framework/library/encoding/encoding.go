package encoding

/*
* --------------内网跨服务调用api的编码规范--------------
* |0001 0000| 					bool类型(1 byte): 无编码类型，第一位bit表示bool值
* |dddd 0000|--1byte--| 		int8、uint8类型：无编码类型。
* |dddd 0000|--<=2byte--|		int16、uint16类型：当val>>14 等于0时, 采用variant编码，否则采用16bit编码
* |dddd xx00|--<=4byte--|		int32、uint32类型：当val>>28 等于0时，采用variant编码，否则采用32bit编码
* |dddd xx00|--<=8byte--|		int64、uint64类型：当val>>56 等于0时，采用variant编码，否则采用64bit编码
* |1010 0000|--4byte--|			float32类型：固定编码类型，直接采用32bit编码
* |1011 0000|--8byte--|			float64类型：固定编码类型，直接采用64bit编码
* |1100 xxxx|xxxx xxxx|--变长--| []byte类型：无编码类型，12bit的长度，最大长度限制4095
* |1101 xxxx|xxxx xxxx|--变长--| string类型：无编码类型，12bit的长度，最大长度限制4095
* |1110 xxxx|xxxx xxxx|--变长--| proto类型：无编码类型，12bit的长度，最大长度限制4095
 */

type WireType uint8

const (
	WireTypeVariant WireType = 0
	WireType16Bit   WireType = 1
	WireType32Bit   WireType = 2
	WireType64Bit   WireType = 3
)

type DataType uint8

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
