package service

import (
	"bytes"
	"hego/Library/file"
	"hego/Library/uerror"
	"hego/tools/cfgtool/internal/manager"
	"hego/tools/cfgtool/internal/parser"

	"github.com/jhump/protoreflect/dynamic"
)

func GenData(dataPath string, buf *bytes.Buffer) error {
	if err := manager.ParseProto(); err != nil {
		return err
	}

	for _, st := range manager.GetConfigMap() {
		RootMsgDesc := manager.GetMessageDescriptor(st.FileName, st.Name+"Ary")
		if RootMsgDesc == nil {
			return uerror.New(1, -1, "find message desc nil: %s", st.Name)
		}
		ItemMsgDesc := manager.GetMessageDescriptor(st.FileName, st.Name)
		if ItemMsgDesc == nil {
			return uerror.New(1, -1, "find item message desc nil: %s", st.Name)
		}

		// 反射new一个对象s
		RootMsg := dynamic.NewMessage(RootMsgDesc)
		if RootMsg == nil {
			return uerror.New(1, -1, "create proto message nil")
		}

		// 加载xlsx数据
		tab := manager.GetTable(st.FileName, st.Sheet)
		for _, vals := range tab.Rows[3:] {
			item := dynamic.NewMessage(ItemMsgDesc)
			for _, field := range st.FieldList {
				item.SetFieldByName(field.Name, parser.ConvertField(field, vals...))
			}
			RootMsg.AddRepeatedFieldByName("Ary", item)
		}

		// 保存数据
		buf, err := RootMsg.Marshal()
		if err != nil {
			return err
		}
		if err := file.Save(dataPath, st.Name+".data", buf); err != nil {
			return err
		}
	}
	return nil
}
