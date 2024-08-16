package generator

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
	"universal/framework/uerror"
	"universal/tool/gomaker/domain"
	"universal/tool/gomaker/internal/manager"
	"universal/tool/gomaker/internal/typespec"
	"universal/tool/gomaker/internal/util"
)

// 生成代码
func HttpKitGenerator(dst string, param string, tpls *template.Template) error {
	if !strings.HasSuffix(dst, ".go") {
		return fmt.Errorf("未定义生成go文件")
	}
	// 生成包头
	buf := bytes.NewBuffer(nil)
	dir := filepath.Dir(dst)
	if err := tpls.ExecuteTemplate(buf, domain.PACKAGE, filepath.Base(dir)); err != nil {
		return uerror.NewUError(1, -1, "%v", err)
	}
	// 获取数据
	arrs := []*typespec.Struct{}
	for _, st := range manager.GetStruct() {
		name := st.Type.Name
		if strings.HasSuffix(name, "Request") || strings.HasSuffix(name, "Response") || strings.HasSuffix(name, "Notify") {
			arrs = append(arrs, st)
		}
	}
	// 模板生成
	if err := tpls.ExecuteTemplate(buf, domain.HTTPKIT, arrs); err != nil {
		return uerror.NewUError(1, -1, "%v", err)
	}
	// 生成文件
	return util.SaveGo(dst, buf)
}

func OmitEmptyGenerator(dst string, param string, tpls *template.Template) error {
	if !strings.HasSuffix(dst, ".go") {
		return fmt.Errorf("未定义生成go文件")
	}
	// 生成包头
	buf := bytes.NewBuffer(nil)
	dir := filepath.Dir(dst)
	if err := tpls.ExecuteTemplate(buf, domain.PACKAGE, filepath.Base(dir)); err != nil {
		return uerror.NewUError(1, -1, "%v", err)
	}
	// 获取struct数据
	arrs := []*typespec.Struct{}
	for _, st := range manager.GetStruct() {
		tmp := st.Clone()
		for _, ff := range tmp.List {
			ff.Tag = strings.ReplaceAll(ff.Tag, ",omitempty", "")
		}
		arrs = append(arrs, tmp)
	}
	// 模板生成
	if err := tpls.ExecuteTemplate(buf, domain.PBCLASS, arrs); err != nil {
		return uerror.NewUError(1, -1, "%v", err)
	}
	// 获取枚举数据
	tmps := []*typespec.Enum{}
	for _, st := range manager.GetEnum() {
		tmps = append(tmps, st)
	}
	// 模板生成
	if err := tpls.ExecuteTemplate(buf, domain.PBCLASS, tmps); err != nil {
		return uerror.NewUError(1, -1, "%v", err)
	}
	// 生成文件
	return util.SaveGo(dst, buf)
}

func ProtoGenerator(dst string, param string, tpls *template.Template) error {
	if !strings.HasSuffix(dst, ".go") {
		return fmt.Errorf("未定义生成go文件")
	}
	// 生成包头
	buf := bytes.NewBuffer(nil)
	dir := filepath.Dir(dst)
	if err := tpls.ExecuteTemplate(buf, domain.PACKAGE, filepath.Base(dir)); err != nil {
		return uerror.NewUError(1, -1, "%v", err)
	}
	// 获取数据
	arrs := []*typespec.Struct{}
	for _, st := range manager.GetStruct() {
		name := st.Type.Name
		if strings.HasSuffix(name, "Request") || strings.HasSuffix(name, "Response") || strings.HasSuffix(name, "Notify") {
			arrs = append(arrs, st)
		}
	}
	// 模板生成
	if err := tpls.ExecuteTemplate(buf, domain.PROTO, arrs); err != nil {
		return uerror.NewUError(1, -1, "%v", err)
	}
	// 生成文件
	return util.SaveGo(dst, buf)
}
