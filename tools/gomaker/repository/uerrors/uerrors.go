package uerrors

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"universal/framework/fbasic"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/base"
	"universal/tools/gomaker/internal/manager"
)

func Gen(action string, dst string, params string) error {
	en := manager.GetEnum("ErrorCode")
	if en == nil {
		return fbasic.NewUError(1, -1, fmt.Sprintf("The enum of ErrorCode is not found in typespec"))
	}
	// 生成文档
	if !strings.HasSuffix(dst, ".go") {
		dst += "/uerrors.gen.go"
	}
	// 生成包头
	buf := bytes.NewBuffer(nil)
	if err := manager.GenPackage(dst, buf); err != nil {
		return err
	}
	// 模版
	if tpl := manager.GetTpl(action); tpl == nil {
		return fbasic.NewUError(1, -1, fmt.Sprintf("The action of %s is not supported", action))
	} else {
		// 生成文件
		if err := tpl.ExecuteTemplate(buf, action+".tpl", en); err != nil {
			return fbasic.NewUError(1, -1, err)
		}
	}
	// 格式化
	result, err := format.Source(buf.Bytes())
	if err != nil {
		return fbasic.NewUError(1, -1, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), os.FileMode(0777)); err != nil {
		return fbasic.NewUError(1, -1, err)
	}
	if err := ioutil.WriteFile(dst, result, os.FileMode(0666)); err != nil {
		return fbasic.NewUError(1, -1, err)
	}
	return nil
}

func Init() {
	manager.Register(&base.Action{
		Name: domain.UERRORS,
		Help: "ErrorCode生成UError错误码",
		Gen:  Gen,
	})
}
