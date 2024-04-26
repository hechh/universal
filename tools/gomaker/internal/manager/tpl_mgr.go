package manager

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"universal/framework/fbasic"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/base"
)

var (
	tpls = make(map[string]*template.Template)
)

func GetTpl(action string) *template.Template {
	return tpls[action]
}

// ./templates目录路径
func InitTpl(root string) {
	// 便利所有模版文件
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if path == root {
			return nil
		}
		if info.IsDir() {
			tpls[filepath.Base(path)] = template.Must(template.ParseGlob(path + "/*.tpl")).Funcs(template.FuncMap{"html": func(s string) string { return s }})
			return nil
		}
		return nil
	})
}

func GenPackage(dst string, buf *bytes.Buffer) error {
	pkg := GetTpl(domain.PACKAGE)
	if pkg == nil {
		return fbasic.NewUError(2, -1, fmt.Sprintf("The tpl of %s.tpl is not supported", domain.PACKAGE))
	}
	if err := pkg.ExecuteTemplate(buf, domain.PACKAGE+".tpl", base.GetFilePathBase(dst)); err != nil {
		return fbasic.NewUError(2, -1, err)
	}
	return nil
}
