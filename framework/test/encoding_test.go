package test

import (
	"fmt"
	"testing"
	"universal/framework/library/encoding"
)

func TestEncoding(t *testing.T) {
	// Variant编码与解码示例
	value := uint64(16512)
	encoded := encoding.VariantEncode(value)
	fmt.Println("Variant encoded value:", encoded)
	decoded := encoding.VariantDecode(encoded)
	fmt.Println("Variant decoded value:", decoded)

	// Zigzag编码与解码示例
	n := int64(-453456421)
	zigzagEncoded := encoding.ZigzagEncode(n)
	fmt.Println("Zigzag encoded value:", zigzagEncoded)
	zigzagDecoded := encoding.ZigzagDecode(zigzagEncoded)
	fmt.Println("Zigzag decoded value:", zigzagDecoded)
}
