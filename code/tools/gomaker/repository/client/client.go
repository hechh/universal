package client

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
	"universal/framework/basic/uerror"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/typespec"
	"universal/tools/gomaker/internal/util"
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
	for _, st := range manager.GetStructList() {
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
	for _, tmp := range manager.GetStructList() {
		j := -1
		for _, ff := range tmp.List {
			if unicode.IsLower(rune(ff.Name[0])) {
				continue
			}
			j++
			tmp.List[j] = ff
			ff.Tag = strings.ReplaceAll(ff.Tag, ",omitempty", "")
		}
		tmp.List = tmp.List[:j+1]
		arrs = append(arrs, tmp)
	}
	// 模板生成
	if err := tpls.ExecuteTemplate(buf, domain.PBCLASS, arrs); err != nil {
		return uerror.NewUError(1, -1, "%v", err)
	}
	// 获取枚举数据 + 模板生成
	if err := tpls.ExecuteTemplate(buf, domain.PBCLASS, manager.GetEnumList()); err != nil {
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
	for _, st := range manager.GetStructList() {
		name := st.Type.Name
		if strings.HasSuffix(name, "Request") {
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
