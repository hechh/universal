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
	"universal/framework/basic"
	"universal/framework/uerror"
)

func Panic(err error) {
	if err == nil {
		return
	}
	if uerr, ok := err.(*uerror.UError); ok {
		panic(uerr.String())
	} else {
		panic(err)
	}
}

// 解析go文件
func ParseFiles(v ast.Visitor, files ...string) error {
	fset := token.NewFileSet()
	for _, filename := range files {
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

func SaveFile(filename string, buf *bytes.Buffer) error {
	// 创建目录
	if err := os.MkdirAll(filepath.Dir(filename), os.FileMode(0777)); err != nil {
		return uerror.NewUError(1, -1, "filename: %s, error: %v", filename, err)
	}

	// 写入文件
	if err := ioutil.WriteFile(filename, buf.Bytes(), os.FileMode(0666)); err != nil {
		return uerror.NewUError(1, -1, "filename: %s, error: %v", filename, err)
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

// 打开所有tpl模板文件
func OpenTemplate(dir string, pattern, filter string, recursive bool) (*template.Template, error) {
	files, err := basic.Glob(dir, pattern, filter, recursive)
	if err != nil {
		return nil, uerror.NewUError(1, -1, "%v", err)
	}
	if len(files) > 0 {
		return template.Must(template.ParseFiles(files...)), nil
	}
	return nil, nil
}
