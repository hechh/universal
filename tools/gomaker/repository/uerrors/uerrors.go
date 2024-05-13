package uerrors

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"universal/framework/fbasic"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/base"
	"universal/tools/gomaker/internal/manager"
	parse "universal/tools/gomaker/internal/parser"
)

func Gen(cwd string, cmdLine *domain.CmdLine, tpls map[string]*template.Template) error {
	en := manager.GetEnum("ErrorCode")
	if en == nil {
		return fbasic.NewUError(1, -1, fmt.Sprintf("The enum of ErrorCode is not found in typespec"))
	}
	// 生成文档
	if !strings.HasSuffix(cmdLine.Dst, ".go") {
		cmdLine.Dst += "/uerrors.gen.go"
	}
	// 生成包头
	buf := bytes.NewBuffer(nil)
	if err := base.MapTpl(tpls).GenPackage(cmdLine.Dst, buf); err != nil {
		return err
	}
	// 模版
	if tpl, ok := tpls[cmdLine.Action]; !ok {
		return fbasic.NewUError(1, -1, fmt.Sprintf("The action of %s is not supported", cmdLine.Action))
	} else {
		// 生成文件
		if err := tpl.ExecuteTemplate(buf, cmdLine.Action+".tpl", en); err != nil {
			return fbasic.NewUError(1, -1, err)
		}
	}
	// 格式化
	result, err := format.Source(buf.Bytes())
	if err != nil {
		return fbasic.NewUError(1, -1, err)
	}
	if err := os.MkdirAll(filepath.Dir(cmdLine.Dst), os.FileMode(0777)); err != nil {
		return fbasic.NewUError(1, -1, err)
	}
	if err := ioutil.WriteFile(cmdLine.Dst, result, os.FileMode(0666)); err != nil {
		return fbasic.NewUError(1, -1, err)
	}
	return nil
}

func Init() {
	manager.Register(parse.NewBaseParser(Gen, domain.UERRORS, "", "ErrorCode生成UError错误码"))
}
