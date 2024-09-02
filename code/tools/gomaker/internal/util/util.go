package util

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
	"universal/framework/basic/uerror"
	"universal/tools/gomaker/domain"
)

// 获取绝对值
func GetAbsPath(cwd, pp string) string {
	if len(cwd) <= 0 || filepath.IsAbs(pp) {
		return filepath.Clean(pp)
	}
	return filepath.Clean(filepath.Join(cwd, pp))
}

// 保存文件
func SaveGo(dst string, buf *bytes.Buffer) error {
	// 格式化数据
	result, err := format.Source(buf.Bytes())
	if err != nil {
		ioutil.WriteFile("./error.gen.go", buf.Bytes(), os.FileMode(0666))
		return uerror.NewUError(1, -1, "%v", err)
	}
	// 创建目录
	if err := os.MkdirAll(filepath.Dir(dst), os.FileMode(0777)); err != nil {
		return uerror.NewUError(1, -1, "%v", err)
	}
	// 写入文件
	if err := ioutil.WriteFile(dst, result, os.FileMode(0666)); err != nil {
		return uerror.NewUError(1, -1, "%v", err)
	}
	return nil
}

func Glob(dir, suffix string, recursive bool) (files []string, err error) {
	if !recursive {
		files, err = filepath.Glob(filepath.Join(dir, suffix))
	} else {
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				return nil
			}
			patterns, err := filepath.Glob(filepath.Join(path, suffix))
			if err != nil {
				return uerror.NewUError(1, -1, "%v", err)
			}
			if len(patterns) > 0 {
				files = append(files, patterns...)
			}
			return nil
		})
	}
	return
}

func OpenTemplate(tpl string) (*template.Template, error) {
	files, err := Glob(tpl, "*.tpl", true)
	if err != nil {
		return nil, err
	}
	result := template.Must(template.ParseFiles(files...))
	result.New(domain.PACKAGE).Parse("package {{.}}")
	return result, nil
}

// 解析整个目录
func ParseDir(v ast.Visitor, fset *token.FileSet, src string) error {
	files, err := Glob(src, "*.go", true)
	if err != nil {
		return err
	}
	for _, filename := range files {
		if err := ParseFile(v, fset, filename); err != nil {
			return err
		}
	}
	return nil
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
