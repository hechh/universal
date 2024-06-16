package fbasic

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"sync"
)

var (
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(nil)
		},
	}
)

func GetBuffer() *bytes.Buffer {
	if obj, ok := bufferPool.Get().(*bytes.Buffer); ok && obj != nil {
		return obj
	}
	return bytes.NewBuffer(nil)
}

func PutBuffer(obj *bytes.Buffer) {
	if obj != nil {
		obj.Reset()
		bufferPool.Put(obj)
	}
}

func DecodeValue(buf []byte, typs []reflect.Type, pos int) (ps []reflect.Value) {
	ps = make([]reflect.Value, len(typs))
	decode := gob.NewDecoder(bytes.NewReader(buf))
	for i := pos; i < len(typs); i++ {
		ps[i] = reflect.New(typs[i]).Elem()
		decode.DecodeValue(ps[i])
	}
	return
}

func DecodeAny(buf []byte, typs []reflect.Type, pos int) (ps []interface{}) {
	ps = make([]interface{}, len(typs))
	decode := gob.NewDecoder(bytes.NewReader(buf))
	for i := pos; i < len(typs); i++ {
		vv := reflect.New(typs[i]).Elem()
		decode.DecodeValue(vv)
		ps[i] = vv.Interface()
	}
	return
}

func EncodeAny(items ...interface{}) []byte {
	bb := GetBuffer()
	defer PutBuffer(bb)
	enc := gob.NewEncoder(bb)
	for _, param := range items {
		enc.Encode(param)
	}
	return bb.Bytes()
}

func EncodeValue(items ...reflect.Value) []byte {
	bb := GetBuffer()
	defer PutBuffer(bb)
	enc := gob.NewEncoder(bb)
	for _, param := range items {
		enc.EncodeValue(param)
	}
	return bb.Bytes()
}
