package test

import (
	"hash/crc32"
	"testing"
)

func TestNode(t *testing.T) {
	t.Log(crc32.ChecksumIEEE([]byte("adfasdf")))
}
