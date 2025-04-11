package manager

import (
	"bytes"
	"hego/Library/uerror"
	"hego/tools/cfgtool/internal/base"
	"hego/tools/xlsx/domain"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/spf13/cast"
)

var (
	convertMgr = make(map[string]*base.Convert)
	protoMgr   = make(map[string]string)
	protoList  = []string{}
	descMap    = make(map[string]*desc.FileDescriptor)
)

func GetConvFunc(name string) func(string) interface{} {
	if val, ok := convertMgr[name]; ok {
		return val.ConvFunc
	}

	// 默认枚举转换函数
	if item, ok := enumMgr[name]; ok {
		return func(str string) interface{} {
			if vv, ok := item.Values[str]; ok {
				return cast.ToInt32(vv)
			}
			return cast.ToInt32(str)
		}
	}
	return nil
}

func GetConvType(name string) string {
	if val, ok := convertMgr[name]; ok {
		return val.Name
	}
	return name
}

func AddProto(filename string, buf *bytes.Buffer) {
	protoMgr[filename] = buf.String()
	protoList = append(protoList, filename)
}

func GetProtoList() []string {
	return protoList
}

func GetProtoMap() map[string]string {
	return protoMgr
}

func ParseProto() error {
	paser := protoparse.Parser{Accessor: protoparse.FileContentsFromMap(protoMgr)}
	descs, err := paser.ParseFiles(protoList...)
	if err != nil {
		return uerror.New(1, -1, "parse proto file error: %s", err.Error())
	}
	for i := range protoList {
		descMap[protoList[i]] = descs[i]
	}
	return nil
}

func GetMessageDescriptor(fileName, name string) *desc.MessageDescriptor {
	if val, ok := descMap[fileName]; ok {
		return val.FindMessage(domain.PkgName + "." + name)
	}
	return nil
}
