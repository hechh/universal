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

// 解析go文件
func ParseFiles(v domain.IParser, files ...string) error {
	fset := token.NewFileSet()
	for _, filename := range files {
		v.SetFile(filename)
		// 解析语法树
		fs, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			return uerror.NewUError(1, -1, "filename: %v, error: %v", filename, err)
		}

		// 遍历语法树
		ast.Walk(v, fs)
	}
	return nil
}

// 保存文件
func SaveGo(filename string, buf *bytes.Buffer) error {
	// 格式化数据
	result, err := format.Source(buf.Bytes())
	if err != nil {
		ioutil.WriteFile("./error.gen.go", buf.Bytes(), os.FileMode(0666))
		return uerror.NewUError(1, -1, "filename: %s, error: %v", filename, err)
	}

	// 创建目录
	if err := os.MkdirAll(filepath.Dir(filename), os.FileMode(0777)); err != nil {
		return uerror.NewUError(1, -1, "filename: %s, error: %v", filename, err)
	}

	// 写入文件
	if err := ioutil.WriteFile(filename, result, os.FileMode(0666)); err != nil {
		return uerror.NewUError(1, -1, "filename: %s, error: %v", filename, err)
	}
	return nil
}

// 遍历目录所有文件
func Glob(dir, pattern string, recursive bool) (files []string, err error) {
	if !recursive {
		files, err = filepath.Glob(filepath.Join(dir, pattern))
	} else {
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				results, err := filepath.Glob(filepath.Join(path, pattern))
				if err != nil {
					return uerror.NewUError(1, -1, "dir: %s, pattern: %s, error: %v", dir, pattern, err)
				}
				if len(results) > 0 {
					files = append(files, results...)
				}
			}
			return nil
		})
	}
	return
}

// 打开所有tpl模板文件
func OpenTemplate(dir string, pattern string, recursive bool) (*template.Template, error) {
	files, err := Glob(dir, pattern, recursive)
	if err != nil {
		return nil, err
	}
	result := template.Must(template.ParseFiles(files...))
	result.New(domain.PACKAGE).Parse("package {{.}}")
	return result, nil
}

// 获取绝对值
func GetAbsPath(cwd, pp string) string {
	if len(cwd) <= 0 || filepath.IsAbs(pp) {
		return filepath.Clean(pp)
	}
	return filepath.Clean(filepath.Join(cwd, pp))
}
