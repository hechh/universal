package parser

import (
	"hego/tools/cfgtool/internal/base"
	"hego/tools/xlsx/domain"
	"hego/tools/xlsx/internal/manager"
	"strings"

	"github.com/jhump/protoreflect/dynamic"
)

func ConvertField(f *base.Field, vals ...string) interface{} {
	if f.Position >= len(vals) {
		return nil
	}
	switch f.Type.TypeOf {
	case domain.TypeOfBase:
		switch f.Type.ValueOf {
		case domain.ValueOfBase:
			return f.ConvFunc(vals[f.Position])
		case domain.ValueOfList:
			return f.Convert(strings.Split(vals[f.Position], ",")...)
		}
	case domain.TypeOfEnum:
		switch f.Type.ValueOf {
		case domain.ValueOfBase:
			return f.ConvFunc(vals[f.Position])
		case domain.ValueOfList:
			return f.Convert(strings.Split(vals[f.Position], ",")...)
		}
	case domain.TypeOfStruct:
		st := manager.GetStruct(f.Type.Name)
		itemDesc := manager.GetMessageDescriptor(st.FileName, st.Name)
		newItem := dynamic.NewMessage(itemDesc)
		switch f.Type.ValueOf {
		case domain.ValueOfBase:
			for _, field := range st.FieldList {
				newItem.SetFieldByName(field.Name)
			}
		case domain.ValueOfList:
		}
	}
	return nil
}
