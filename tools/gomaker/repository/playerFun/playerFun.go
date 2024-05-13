package playerFun

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
	"universal/tools/gomaker/internal/typespec"
)

type PlayerFunAttr struct {
	typespec.BaseFunc
	Name    string
	ReqList []string
}

func parseParams(params string) *PlayerFunAttr {
	strs := strings.Split(params, ":")
	return &PlayerFunAttr{Name: strs[0], ReqList: strings.Split(strs[1], ",")}
}

// name, pbname
func Gen(cwd string, cmdLine *domain.CmdLine, tpls map[string]*template.Template) error {
	attr := parseParams(cmdLine.Param)
	// 生成文档
	if !strings.HasSuffix(cmdLine.Dst, ".go") {
		cmdLine.Dst += fmt.Sprintf("/%s.go", attr.Name)
	}
	// 生成包头
	buf := bytes.NewBuffer(nil)
	if err := base.MapTpl(tpls).GenPackage(cmdLine.Dst, buf); err != nil {
		return err
	}
	// 模版
	if tpl, ok := tpls[cmdLine.Action]; !ok {
		return fbasic.NewUError(1, -1, fmt.Sprintf("The cmdline.action of %s is not supported", cmdLine.Action))
	} else {
		// 生成文件
		if err := tpl.ExecuteTemplate(buf, cmdLine.Action+".tpl", attr); err != nil {
			return fbasic.NewUError(1, -1, err)
		}
	}
	// 格式化
	result, err := format.Source(buf.Bytes())
	if err != nil {
		//ioutil.WriteFile("./gen.go", buf.Bytes(), os.FileMode(0644))
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
	manager.Register(parse.NewBaseParser(Gen, domain.PLAYER_FUN, "{PlayerXxxxFun}:{req1},{req2},...", "生成playerFun模板文件"))
}
