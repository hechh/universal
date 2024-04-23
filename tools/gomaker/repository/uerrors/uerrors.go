package uerrors

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
	"universal/tools/gomaker/internal/manager"
)

func Gen() error {
	en := manager.GetEnum("ErrorCode")
	if en == nil {
		return nil
	}
	// 模版
	tpl := template.Must(template.ParseFiles("./templates/uerrors/uerrors.tpl"))
	buf := bytes.NewBuffer(nil)
	if err := tpl.Execute(buf, en); err != nil {
		return err
	}
	// 格式化
	result, err := format.Source(buf.Bytes())
	if err != nil {
		//ioutil.WriteFile("./gen.go", buf.Bytes(), os.FileMode(0644))
		return err
	}
	// 生成文档
	genFile := "../../common/uerrors/uerrors.gen.go"
	if err := os.MkdirAll(filepath.Dir(genFile), os.FileMode(0777)); err != nil {
		return err
	}
	if err := ioutil.WriteFile(genFile, result, os.FileMode(0666)); err != nil {
		return err
	}
	return nil
}
