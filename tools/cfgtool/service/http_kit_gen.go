package service

import (
	"bytes"
	"path/filepath"
	"universal/tools/cfgtool/domain"
	"universal/tools/cfgtool/internal/base"
	"universal/tools/cfgtool/internal/manager"
	"universal/tools/cfgtool/internal/templ"
)

type HttpKit struct {
	Pkg  string
	Data map[uint32]string
}

func GenHttpKit(buf *bytes.Buffer) error {
	if len(domain.ClientPath) <= 0 {
		return nil
	}

	item := &HttpKit{
		Pkg:  filepath.Base(domain.ClientPath),
		Data: manager.GetCmdMap(),
	}

	buf.Reset()
	if err := templ.HttpKitTpl.Execute(buf, item); err != nil {
		return err
	}
	// 保存代码
	if err := base.SaveGo(domain.ClientPath, "Init.gen.go", buf.Bytes()); err != nil {
		return err
	}
	return nil
}
