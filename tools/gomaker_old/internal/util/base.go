package util

import (
	"go/ast"
	"go/parser"
	"go/token"
	"hego/framework/basic"
	"hego/framework/uerror"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
	"unicode"
)

func Panic(err interface{}) {
	switch vv := err.(type) {
	case nil:
	case *uerror.UError:
		panic(vv.String())
	default:
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

func SaveFile(filename string, buf []byte) error {
	// 创建目录
	if err := os.MkdirAll(filepath.Dir(filename), os.FileMode(0777)); err != nil {
		return uerror.NewUError(1, -1, "filename: %s, error: %v", filename, err)
	}

	// 写入文件
	if err := ioutil.WriteFile(filename, buf, os.FileMode(0666)); err != nil {
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

// 大小写转成下划线
func ToUnderline(word string) string {
	result := []byte{}
	for i, ch := range word {
		if !unicode.IsUpper(ch) {
			result = append(result, byte(ch))
		} else {
			if i == 0 {
				result = append(result, byte(unicode.ToLower(ch)))
			} else {
				result = append(result, '_')
				result = append(result, byte(unicode.ToLower(ch)))
			}
		}
	}
	return basic.BytesToString(result)
}
