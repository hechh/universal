package manager

import (
	"html/template"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	tpls = make(map[string]*template.Template)
)

func init() {
	// 获取绝对地址
	_, filename, _, _ := runtime.Caller(0)
	datapath := filepath.Join(path.Dir(filename), "../../templates/")
	absPath, _ := filepath.Abs(datapath)
	// 便利所有模版文件
	filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || path == absPath {
			return nil
		}
		path = filepath.Dir(path)
		action := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
		tpls[action] = template.Must(template.ParseGlob(path + "/*.tpl"))
		return nil
	})

}
