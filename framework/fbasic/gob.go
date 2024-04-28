package fbasic

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

type DecodeTypes []reflect.Type

func (d DecodeTypes) DecodeValues(buf []byte, ps []reflect.Value) {
	decode := gob.NewDecoder(bytes.NewReader(buf))
	for i := 1; i < len(d); i++ {
		ps[i-1] = reflect.New(d[i]).Elem()
		decode.DecodeValue(ps[i-1])
	}
}

func (d DecodeTypes) Decode(buf []byte, ps []interface{}) {
	decode := gob.NewDecoder(bytes.NewReader(buf))
	for i := 1; i < len(d); i++ {
		ps[i-1] = reflect.New(d[i]).Elem().Interface()
		decode.Decode(ps[i-1])
	}
}

type EncodeValues []reflect.Value

func (d EncodeValues) Encode() []byte {
	bb := GetBuffer()
	defer PutBuffer(bb)
	enc := gob.NewEncoder(bb)
	for _, param := range d {
		enc.EncodeValue(param)
	}
	return bb.Bytes()
}

type EncodeAnys []interface{}

func (d EncodeAnys) Encode() []byte {
	bb := GetBuffer()
	defer PutBuffer(bb)
	enc := gob.NewEncoder(bb)
	for _, param := range d {
		enc.Encode(param)
	}
	return bb.Bytes()
}
