package playerFun

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"universal/framework/basic"
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
	strs := strings.Split(params, ",")
	attr := &PlayerFunAttr{
		Name: strs[0],
		ReqList: func() (ret []string) {
			if len(strs) > 1 {
				ret = strs[1:]
			}
			return
		}(),
	}
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
		return basic.NewUError(2, -1, fmt.Sprintf("The action of %s is not supported", action))
	} else {
		// 生成文件
		if err := tpl.ExecuteTemplate(buf, action+".tpl", attr); err != nil {
			return basic.NewUError(2, -1, err)
		}
	}
	// 格式化
	result, err := format.Source(buf.Bytes())
	if err != nil {
		//ioutil.WriteFile("./gen.go", buf.Bytes(), os.FileMode(0644))
		return basic.NewUError(2, -1, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), os.FileMode(0777)); err != nil {
		return basic.NewUError(2, -1, err)
	}
	if err := ioutil.WriteFile(dst, result, os.FileMode(0666)); err != nil {
		return basic.NewUError(2, -1, err)
	}
	return nil
}
