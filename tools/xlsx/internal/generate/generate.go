package generate

import (
	"bytes"
	"fmt"
	"hego/Library/file"
	"hego/Library/util"
	"hego/tools/xlsx/domain"
	"hego/tools/xlsx/internal/base"
)

func gindex(i int) string {
	return fmt.Sprintf("group%d", i)
}

func mindex(i int) string {
	return fmt.Sprintf("map%d", i)
}

func mapKey(i int) string {
	return fmt.Sprintf("mapData%d", i)
}

func groupKey(i int) string {
	return fmt.Sprintf("groupData%d", i)
}

func Generate(codePath string, cfg *base.Config, buf *bytes.Buffer) error {
	sts := NewStInfoMgr("cfg", cfg)
	sts.Push(sts.DataName, "listData", fmt.Sprintf("[]*%s", sts.CType))
	for i, vals := range cfg.MapList {
		ktype, index := mapKey(i), mindex(i)
		var args, values, refs []string
		for _, field := range vals {
			args = append(args, fmt.Sprintf("%s %s", field.Name, field.Type.GetType(sts.PkgName)))
			values = append(values, field.Name)
			refs = append(refs, fmt.Sprintf("item.%s", field.Name))
			if len(vals) > 1 {
				sts.Push(index, field.Name, field.Type.GetType(sts.PkgName))
			}
		}
		sts.AddFun(&FunInfo{
			Pos:     i,
			ValueOf: domain.VALUE_OF_MAP,
			Name:    ktype,
			Index:   util.Ifelse(len(vals) == 1, vals[0].Type.GetType(sts.PkgName), index),
			Args:    args,
			Values:  values,
			Refs:    refs,
		})
		if len(vals) == 1 {
			sts.Add(sts.DataName, ktype, fmt.Sprintf("map[%s]*%s", vals[0].Type.GetType(sts.PkgName), sts.CType))
		} else {
			sts.Add(sts.DataName, ktype, fmt.Sprintf("map[%s]*%s", index, sts.CType))
		}
	}
	for i, vals := range cfg.GroupList {
		ktype, index := groupKey(i), gindex(i)
		var args, values, refs []string
		for _, field := range vals {
			args = append(args, fmt.Sprintf("%s %s", field.Name, field.Type.GetType(sts.PkgName)))
			values = append(values, field.Name)
			refs = append(refs, fmt.Sprintf("item.%s", field.Name))
			if len(vals) > 1 {
				sts.Push(index, field.Name, field.Type.GetType(sts.PkgName))
			}
		}
		sts.AddFun(&FunInfo{
			Pos:     i,
			ValueOf: domain.VALUE_OF_GROUP,
			Name:    ktype,
			Index:   util.Ifelse(len(vals) == 1, vals[0].Type.GetType(sts.PkgName), index),
			Args:    args,
			Values:  values,
			Refs:    refs,
		})
		if len(vals) == 1 {
			sts.Add(sts.DataName, ktype, fmt.Sprintf("map[%s][]*%s", vals[0].Type.GetType(sts.PkgName), sts.CType))
		} else {
			sts.Add(sts.DataName, ktype, fmt.Sprintf("map[%s][]*%s", index, sts.CType))
		}
	}

	buf.WriteString(sts.Package())
	buf.WriteString(sts.Define())
	buf.WriteString(sts.Func())
	buf.WriteString(sts.Parse())
	return file.SaveGo(codePath, fmt.Sprintf("%s.gen.go", sts.DataName), buf.Bytes())
}
