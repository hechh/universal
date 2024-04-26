package playerFun

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
	"universal/tools/gomaker/internal/typespec"
)

type PlayerFunAttr struct {
	typespec.BaseFunc
	Name    string
	ReqList []string
}

// name, pbname
func Gen(action string, dst string, params string) error {
	attr := parseParams(params)
	// 生成文档
	if !strings.HasSuffix(dst, ".go") {
		dst += fmt.Sprintf("/%s.go", attr.Name)
	}
	// 生成包头
	buf := bytes.NewBuffer(nil)
	if err := manager.GenPackage(dst, buf); err != nil {
		return err
	}
	// 模版
	if tpl := manager.GetTpl(action); tpl == nil {
		return fbasic.NewUError(2, -1, fmt.Sprintf("The action of %s is not supported", action))
	} else {
		// 生成文件
		if err := tpl.ExecuteTemplate(buf, action+".tpl", attr); err != nil {
			return fbasic.NewUError(2, -1, err)
		}
	}
	// 格式化
	result, err := format.Source(buf.Bytes())
	if err != nil {
		//ioutil.WriteFile("./gen.go", buf.Bytes(), os.FileMode(0644))
		return fbasic.NewUError(2, -1, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), os.FileMode(0777)); err != nil {
		return fbasic.NewUError(2, -1, err)
	}
	if err := ioutil.WriteFile(dst, result, os.FileMode(0666)); err != nil {
		return fbasic.NewUError(2, -1, err)
	}
	return nil
}

func Init() {
	manager.Register(&base.Action{
		Name:  domain.PLAYER_FUN,
		Help:  "生成playerFun模板文件",
		Param: `"fun:{PlayerXxxxFun}|req:{req1},{req2},..."`,
		Gen:   Gen,
	})
}

func parseParams(params string) *PlayerFunAttr {
	ret := &PlayerFunAttr{}
	strs := strings.Split(params, "|")
	for _, str := range strs {
		pos := strings.Index(str, ":")
		switch str[:pos+1] {
		case "fun":
			ret.Name = str[pos:]
		case "req":
			ret.ReqList = strings.Split(str[:pos], ",")
		}
	}
	return ret
}
