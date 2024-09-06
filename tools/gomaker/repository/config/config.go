package config

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"universal/framework/basic/uerror"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/util"
)

func XlsxGenerator(dst string, param string, tpls *template.Template) error {
	if !strings.HasSuffix(dst, ".go") {
		return fmt.Errorf("未定义生成的go文件")
	}
	// 生成包头
	buf := bytes.NewBuffer(nil)
	dir := filepath.Dir(dst)
	if err := tpls.ExecuteTemplate(buf, domain.PACKAGE, filepath.Base(dir)); err != nil {
		return uerror.NewUError(1, -1, "dst: %s, param: %s, error: %v", dst, param, err)
	}

	// 获取数据
	arrs := manager.GetStructList()
	sort.Slice(arrs, func(i, j int) bool {
		return strings.Compare(arrs[i].Type.Name, arrs[j].Type.Name) < 0
	})

	// 模板生成
	if err := tpls.ExecuteTemplate(buf, domain.XLSX, arrs); err != nil {
		return uerror.NewUError(1, -1, "%v", err)
	}
	// 生成文件
	return util.SaveGo(dst, buf)
}
