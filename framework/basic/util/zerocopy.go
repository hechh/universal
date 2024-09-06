package util

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"reflect"
	"unsafe"
)

func StringToBytes(str string) []byte {
	if len(str) == 0 {
		return nil
	}
	s := *(*reflect.StringHeader)(unsafe.Pointer(&str))
	b := &reflect.SliceHeader{Data: s.Data, Len: s.Len, Cap: s.Len}
	return *(*[]byte)(unsafe.Pointer(b))
}

func BytesToString(bts []byte) string {
	if len(bts) == 0 {
		return ""
	}
	b := *(*reflect.SliceHeader)(unsafe.Pointer(&bts))
	s := &reflect.StringHeader{Data: b.Data, Len: b.Len}
	return *(*string)(unsafe.Pointer(s))
}

// 整形转换成字节
func IntToBytes(val int) []byte {
	tmp := uint32(val)
	buff := make([]byte, 4)
	binary.BigEndian.PutUint32(buff, tmp)
	return buff
}

// 字节转换成整形
func BytesToInt(data []byte) int {
	buff := make([]byte, 4)
	copy(buff, data)
	tmp := int32(binary.BigEndian.Uint32(buff))
	return int(tmp)
}

func MD5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func GetCrc(name string) uint32 {
	return crc32.ChecksumIEEE([]byte(name))
}

func GetString(str, def string) string {
	if len(str) > 0 {
		return str
	}
	return def
}
