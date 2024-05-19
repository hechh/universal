package uerrors

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"universal/framework/common/uerror"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/base"
	"universal/tools/gomaker/internal/maker"
	"universal/tools/gomaker/internal/manager"
)

func Gen(cmdLine *domain.CmdLine, tpls *base.Templates) error {
	en := manager.GetEnum("ErrorCode")
	if en == nil {
		return uerror.NewUError(1, -1, fmt.Sprintf("The enum of ErrorCode is not found in typespec"))
	}
	// 生成文档
	dstFile := cmdLine.Dst
	if !strings.HasSuffix(dstFile, ".go") {
		dstFile += "/uerrors.gen.go"
	}
	// 生成包头
	buf := bytes.NewBuffer(nil)
	if err := tpls.GenPackage(buf, base.GetPkgName(dstFile)); err != nil {
		return err
	}
	// 模版
	if err := tpls.Execute(cmdLine.Action, cmdLine.Action+".tpl", buf, en); err != nil {
		return err
	}
	// 格式化
	result, err := format.Source(buf.Bytes())
	if err != nil {
		return uerror.NewUError(1, -1, err)
	}
	if err := os.MkdirAll(filepath.Dir(dstFile), os.FileMode(0777)); err != nil {
		return uerror.NewUError(1, -1, err)
	}
	if err := ioutil.WriteFile(dstFile, result, os.FileMode(0666)); err != nil {
		return uerror.NewUError(1, -1, err)
	}
	return nil
}

func Init() {
	manager.Register("uerrors", maker.NewBaseMaker(Gen, "", "ErrorCode生成UError错误码"))
}
