package manager

import (
	"universal/tools/cfgtool/domain"
	"universal/tools/cfgtool/internal/typespec"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
)

var (
	enums   = make(map[string]*typespec.Enum)
	structs = make(map[string]*typespec.Struct)
	configs = make(map[string]*typespec.Config)
	protos  = make(map[string]*typespec.Proto)
	files   = make(map[string]string)
	descMap = make(map[string]*desc.FileDescriptor)
)

func NewProto(fileName, name string) *dynamic.Message {
	val, ok := descMap[fileName+".proto"]
	if !ok {
		return nil
	}
	typeOf := val.FindMessage("universal." + name)
	if typeOf == nil {
		return nil
	}
	return dynamic.NewMessage(typeOf)
}

func AddProtoFile(fileList []string, tmps map[string]string) error {
	files = tmps
	paser := protoparse.Parser{Accessor: protoparse.FileContentsFromMap(tmps)}
	descs, err := paser.ParseFiles(fileList...)
	if err != nil {
		return nil
	}
	for i := range fileList {
		descMap[fileList[i]] = descs[i]
	}
	return nil
}

func GetProtoFiles() map[string]string {
	return files
}

func GetTypeOf(name string) int {
	if _, ok := enums[name]; ok {
		return domain.TypeOfEnum
	}
	if _, ok := structs[name]; ok {
		return domain.TypeOfStruct
	}
	if _, ok := configs[name]; ok {
		return domain.TypeOfConfig
	}
	return domain.TypeOfBase
}

func GetEnum(name string) *typespec.Enum {
	return enums[name]
}

func GetStruct(name string) *typespec.Struct {
	return structs[name]
}

func GetConfig(name string) *typespec.Config {
	return configs[name]
}

func GetProto(file string) *typespec.Proto {
	return protos[file]
}

func GetEnums() map[string]*typespec.Enum {
	return enums
}

func GetStructs() map[string]*typespec.Struct {
	return structs
}

func GetConfigs() map[string]*typespec.Config {
	return configs
}

func GetProtos() map[string]*typespec.Proto {
	return protos
}

func GetOrNewEnum(name string) *typespec.Enum {
	if val, ok := enums[name]; ok {
		return val
	}
	val := &typespec.Enum{
		Name:   name,
		Values: make(map[string]*typespec.Value),
	}
	enums[name] = val
	return val
}

func GetOrNewStruct(name string) *typespec.Struct {
	if val, ok := structs[name]; ok {
		return val
	}
	val := &typespec.Struct{
		Name:     name,
		Fields:   make(map[string]*typespec.Field),
		Converts: make(map[string][]*typespec.Field),
	}
	structs[name] = val
	return val
}

func GetOrNewConfig(name string) *typespec.Config {
	if val, ok := configs[name]; ok {
		return val
	}
	val := &typespec.Config{
		Name:   name,
		Fields: make(map[string]*typespec.Field),
		Indexs: make(map[string]*typespec.Index),
	}
	configs[name] = val
	return val
}

func GetOrNewProto(file string) *typespec.Proto {
	if val, ok := protos[file]; ok {
		return val
	}
	val := &typespec.Proto{
		FileName:   file,
		References: make(map[string]struct{}),
	}
	protos[file] = val
	return val
}
