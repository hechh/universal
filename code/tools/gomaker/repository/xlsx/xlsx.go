package xlsx

import (
	"path/filepath"
	"text/template"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/util"
)

// 生成代码
func JsonGenerator(dst string, param string, tpls *template.Template) error {
	for name, data := range manager.GetJsons() {
		if err := util.SaveJson(filepath.Join(dst, name), data); err != nil {
			return err
		}
	}
	return nil
}
