package fbasic

import (
	"bytes"
	"encoding/gob"
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

func ValueToDecode(buf []byte, typs []reflect.Type, pos int) (ps []reflect.Value) {
	ps = make([]reflect.Value, len(typs))
	decode := gob.NewDecoder(bytes.NewReader(buf))
	for i := pos; i < len(typs); i++ {
		ps[i] = reflect.New(typs[i]).Elem()
		decode.DecodeValue(ps[i])
	}
	return
}

func AnyToDecode(buf []byte, typs []reflect.Type, pos int) (ps []interface{}) {
	ps = make([]interface{}, len(typs))
	decode := gob.NewDecoder(bytes.NewReader(buf))
	for i := pos; i < len(typs); i++ {
		vv := reflect.New(typs[i]).Elem()
		decode.DecodeValue(vv)
		ps[i] = vv.Interface()
	}
	return
}

func AnyToEncode(items ...interface{}) []byte {
	bb := GetBuffer()
	defer PutBuffer(bb)
	enc := gob.NewEncoder(bb)
	for _, param := range items {
		enc.Encode(param)
	}
	return bb.Bytes()
}

func ValueToEncode(items ...reflect.Value) []byte {
	bb := GetBuffer()
	defer PutBuffer(bb)
	enc := gob.NewEncoder(bb)
	for _, param := range items {
		enc.EncodeValue(param)
	}
	return bb.Bytes()
}
