package service

import (
	"bytes"
	"universal/library/util"
	"universal/tools/cfgtool/domain"
	"universal/tools/cfgtool/internal/manager"
	"universal/tools/cfgtool/internal/templ"
)

func GenProto(buf *bytes.Buffer) error {
	files := []string{}
	tmps := map[string]string{}
	for _, item := range manager.GetProtos() {
		buf.Reset()
		if err := templ.ProtoTpl.Execute(buf, item); err != nil {
			return err
		}

		filename := item.FileName + ".proto"
		data := buf.Bytes()
		if len(domain.ProtoPath) > 0 {
			if err := util.Save(domain.ProtoPath, filename, data); err != nil {
				return err
			}
		}
		files = append(files, filename)
		tmps[filename] = string(data)
	}

	manager.AddProtoFile(files, tmps)
	return nil
}
