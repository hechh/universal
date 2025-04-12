package service

import (
	"bytes"
	"hego/Library/file"
	"hego/Library/uerror"
	"hego/tools/cfgtool/domain"
	"hego/tools/cfgtool/internal/base"
	"hego/tools/cfgtool/internal/manager"
	"strings"
)

func GenData(dataPath string, buf *bytes.Buffer) error {
	for _, cfg := range manager.GetConfigMap() {
		// 反射new一个对象
		ary := manager.NewProto(cfg.FileName, cfg.Name+"Ary")
		if ary == nil {
			return uerror.New(1, -1, "new %sAry is nil", cfg.Name)
		}

		// 加载xlsx数据
		tab := manager.GetTable(cfg.FileName, cfg.Sheet)
		for _, vals := range tab.Rows[3:] {
			// 反射new一个对象
			item := manager.NewProto(cfg.FileName, cfg.Name)
			if item == nil {
				return uerror.New(1, -1, "new %s is nil", cfg.Name)
			}
			ary.AddRepeatedFieldByName("Ary", configValue(cfg, vals))
		}

		// 保存数据
		buf, err := ary.Marshal()
		if err != nil {
			return err
		}
		if err := file.Save(dataPath, cfg.Name+".data", buf); err != nil {
			return err
		}
	}
	return nil
}

func configValue(f *base.Config, vals []string) interface{} {
	item := manager.NewProto(f.FileName, f.Name)
	for i, field := range f.FieldList {
		if field.Position >= len(vals) {
			item.SetFieldByName(field.Name, nil)
			continue
		}

		switch field.Type.TypeOf {
		case domain.TypeOfBase, domain.TypeOfEnum:
			item.SetFieldByName(field.Name, fieldValue(field, vals[i]))
		case domain.TypeOfStruct:
			st := manager.GetStruct(field.Type.Name)
			item.SetFieldByName(field.Name, structValue(st, vals[i]))
		}
	}
	return nil
}

func structValue(f *base.Struct, val string) interface{} {
	strs := strings.Split(val, "|")
	if len(f.Converts[strs[0]]) >= len(strs) {
		return uerror.New(1, -1, "struct %s field %s not enough", f.Name, strs[0])
	}

	item := manager.NewProto(f.FileName, f.Name)
	for i, field := range f.Converts[strs[0]] {
		item.SetFieldByName(field.Name, fieldValue(field, strs[i]))
	}
	return item
}

func fieldValue(f *base.Field, val string) interface{} {
	switch f.Type.ValueOf {
	case domain.ValueOfBase:
		return f.ConvFunc(val)
	case domain.ValueOfList:
		return f.Convert(strings.Split(val, ",")...)
	}
	return nil
}
