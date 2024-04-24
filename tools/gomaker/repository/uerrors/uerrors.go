package uerrors

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"universal/tools/gomaker/internal/manager"
)

func Gen(action string, dst string) error {
	en := manager.GetEnum("ErrorCode")
	if en == nil {
		return fmt.Errorf("The enum of ErrorCode is not found in typespec")
	}
	// 模版
	tpl := manager.GetTpl(action)
	if tpl == nil {
		return fmt.Errorf("The action of %s is not supported", action)
	}
	// 生成文件
	buf := bytes.NewBuffer(nil)
	if err := tpl.ExecuteTemplate(buf, "uerrors.tpl", en); err != nil {
		return err
	}
	// 格式化
	result, err := format.Source(buf.Bytes())
	if err != nil {
		//ioutil.WriteFile("./gen.go", buf.Bytes(), os.FileMode(0644))
		return err
	}
	// 生成文档
	if !strings.HasSuffix(dst, ".go") {
		dst += "/uerrors.gen.go"
	}
	if err := os.MkdirAll(filepath.Dir(dst), os.FileMode(0777)); err != nil {
		return err
	}
	if err := ioutil.WriteFile(dst, result, os.FileMode(0666)); err != nil {
		return err
	}
	return nil
}
