package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"universal/tool/gomaker/domain"
)

func GenPBClass(action, dst string, tpls map[string]*template.Template) error {
	if strings.HasSuffix(dst, ".go") {
		return fmt.Errorf("生成文件必须是.go")
	}
	// 判断模板文件是否存在
	tpl, ok := tpls[action]
	if !ok {
		return fmt.Errorf("%s模板文件不存在", action)
	}
	// 生成包头
	buf := bytes.NewBuffer(nil)
	if err := tpls[domain.PACKAGE].ExecuteTemplate(buf, domain.PACKAGE, filepath.Base(dst)); err != nil {
		return err
	}
	// 生成文件
	return nil
}
