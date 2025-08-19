package encode

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"sync"

	"github.com/golang/protobuf/proto"
)

type GobEncoder struct {
	buf *bytes.Buffer
	enc *gob.Encoder
}

type GobDecoder struct {
	buf *bytes.Buffer
	dec *gob.Decoder
}

var (
	encPool = sync.Pool{
		New: func() interface{} {
			buf := bytes.NewBuffer(make([]byte, 0, 1024))
			enc := gob.NewEncoder(buf)
			return &GobEncoder{buf: buf, enc: enc}
		},
	}
	decPool = sync.Pool{
		New: func() interface{} {
			buf := bytes.NewBuffer(make([]byte, 0, 1024))
			dec := gob.NewDecoder(buf)
			return &GobDecoder{buf: buf, dec: dec}
		},
	}
)

// 编码
func Encode(args ...interface{}) ([]byte, error) {
	item := encPool.Get().(*GobEncoder)
	defer encPool.Put(item)
	item.buf.Reset()
	for _, arg := range args {
		if err := item.enc.Encode(arg); err != nil {
			return nil, err
		}
	}
	rets := make([]byte, item.buf.Len())
	copy(rets, item.buf.Bytes())
	return rets, nil
}

// 解码
func Decode(data []byte, mfun reflect.Method, rets []reflect.Value, pos int) (err error) {
	item := decPool.Get().(*GobDecoder)
	defer decPool.Put(item)
	item.buf.Reset()
	item.buf.Write(data)
	for i := pos; i < mfun.Type.NumIn(); i++ {
		val := reflect.New(mfun.Type.In(i))
		if err = item.dec.DecodeValue(val); err != nil {
			return err
		}
		rets[i] = val.Elem()
	}
	return
}

func Marshal(args ...interface{}) ([]byte, error) {
	if len(args) == 1 {
		switch vv := args[0].(type) {
		case []byte:
			return vv, nil
		case proto.Message:
			buf, err := proto.Marshal(vv)
			if err != nil {
				return nil, err
			}
			return buf, nil
		}
	}
	return Encode(args...)
}
