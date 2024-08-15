package util

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"text/template"
	"universal/framework/uerror"
	"universal/tool/gomaker/domain"
)

// 获取绝对值
func GetAbsPath(cwd, pp string) string {
	if len(cwd) <= 0 || filepath.IsAbs(pp) {
		return filepath.Clean(pp)
	}
	return filepath.Clean(filepath.Join(cwd, pp))
}

func OpenTemplate(tpl string) (map[string]*template.Template, error) {
	ret := make(map[string]*template.Template)
	if len(tpl) <= 0 {
		return ret, nil
	}
	// 遍历所有目录
	err := filepath.Walk(tpl, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		// 模式匹配
		pattern, err := filepath.Glob(filepath.Join(path, "*.tpl"))
		if err != nil {
			return uerror.NewUError(1, -1, "%v", err)
		}
		if len(pattern) > 0 {
			ret[filepath.Base(path)] = template.Must(template.ParseFiles(pattern...))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// 添加默认包名模板
	ret[domain.PACKAGE] = template.Must(template.New(domain.PACKAGE).Parse("package {{.}}"))
	return ret, nil
}

// 解析文件
func ParseDir(v ast.Visitor, fset *token.FileSet, src string) error {
	if len(src) <= 0 {
		return nil
	}
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		pattern, err := filepath.Glob(filepath.Join(path, "*.go"))
		if err != nil {
			return uerror.NewUError(1, -1, "%v", err)
		}
		for _, filename := range pattern {
			if err := ParseFile(v, fset, filename); err != nil {
				return err
			}
		}
		return nil
	})
}

// 解析单个文件
func ParseFile(v ast.Visitor, fset *token.FileSet, filename string) error {
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return uerror.NewUError(1, -1, "%v", err)
	}
	ast.Walk(v, f)
	return nil
}
