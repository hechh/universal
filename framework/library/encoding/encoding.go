package encoding

// Variant编码
func VariantEncode(value uint64) []byte {
	var buf []byte
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
	return buf
}

// Variant解码
func VariantDecode(data []byte) uint {
	var value uint
	var shift uint

	for _, b := range data {
		value |= uint(b&0x7f) << shift
		shift += 7
		if b&0x80 == 0 {
			break
		}
	}

	return value
}

// Zigzag编码
func ZigzagEncode(n int64) uint64 {
	return uint64((n << 1) ^ (n >> 63))
}

// Zigzag解码
func ZigzagDecode(n uint64) int {
	return int((n >> 1) ^ -(n & 1))
}
