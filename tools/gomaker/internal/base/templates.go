package base

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"
	"universal/framework/fbasic"
)

type Templates struct {
	tpls map[string]*template.Template
}

func NewTemplates(curPath string) *Templates {
	tmps := make(map[string]*template.Template)
	if len(curPath) <= 0 {
		return &Templates{tpls: tmps}
	}
	// 便利所有模版文件
	funcs := template.FuncMap{"html": func(s string) string { return s }}
	filepath.Walk(curPath, func(path string, info os.FileInfo, err error) error {
		if path != curPath && info.IsDir() {
			tmps[filepath.Base(path)] = template.Must(template.ParseGlob(path + "/*.tpl")).Funcs(funcs)
		}
		return nil
	})
	return &Templates{tpls: tmps}
}

func (d *Templates) GenPackage(buf *bytes.Buffer, data interface{}) error {
	action := "package"
	val, ok := d.tpls[action]
	if !ok {
		return fbasic.NewUError(1, -1, action)
	}
	if err := val.ExecuteTemplate(buf, action+".tpl", data); err != nil {
		return fbasic.NewUError(1, -1, err)
	}
	return nil
}

func (d *Templates) Execute(action string, tplfile string, buf *bytes.Buffer, data interface{}) error {
	val, ok := d.tpls[action]
	if !ok {
		return fbasic.NewUError(1, -1, action)
	}
	if err := val.ExecuteTemplate(buf, tplfile, data); err != nil {
		return fbasic.NewUError(1, -1, err)
	}
	return nil
}
