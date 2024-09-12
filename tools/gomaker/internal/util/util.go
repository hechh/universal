package util

import (
	"bytes"
	"encoding/json"
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
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/typespec"

	"github.com/xuri/excelize/v2"
)

// 保存文件
func SaveJson(filename string, data []map[string]interface{}) error {
	// 创建目录
	if err := os.MkdirAll(filepath.Dir(filename), os.FileMode(0777)); err != nil {
		return uerror.NewUError(1, -1, "filename: %s, error: %v", filename, err)
	}

	// 写入文件
	buf, _ := json.MarshalIndent(&data, "", "	")
	if err := ioutil.WriteFile(filename, buf, os.FileMode(0666)); err != nil {
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

// 获取绝对值
func GetAbsPath(cwd, pp string) string {
	if len(cwd) <= 0 || filepath.IsAbs(pp) {
		return filepath.Clean(pp)
	}
	return filepath.Clean(filepath.Join(cwd, pp))
}

// 打开所有tpl模板文件
func OpenTemplate(dir string, pattern string, recursive bool) (*template.Template, error) {
	files, err := basic.Glob(dir, pattern, "", recursive)
	if err != nil {
		return nil, err
	}
	result := template.Must(template.ParseFiles(files...))
	result.New(domain.PACKAGE).Parse("package {{.}}")
	return result, nil
}

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

// 读取xlsx文件
func ReadXlsx(filename string) (map[string][][]string, []string, error) {
	fb, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, nil, uerror.NewUError(1, -1, "filename: %s, error: %v", filename, err)
	}
	defer fb.Close()
	result := make(map[string][][]string)
	sheets := fb.GetSheetList()
	for _, sheet := range sheets {
		if values, err := fb.GetRows(sheet); err != nil {
			return nil, nil, err
		} else {
			result[sheet] = values
		}
	}
	return result, sheets, nil
}

// 解析xlsx
func ParseXlsxs(v domain.IParser, f func([]string) []string, files ...string) error {
	if f == nil {
		f = func(s []string) []string { return s }
	}
	for _, filename := range files {
		v.SetFile(filename)
		// 加载数据
		rows, sheets, err := ReadXlsx(filename)
		if err != nil {
			return err
		}
		// 解析数据
		for _, sheet := range f(sheets) {
			v.Visit(&typespec.SheetNode{Sheet: sheet, Rows: rows[sheet]})
		}
	}
	return nil
}

// 需要规定xlsx文件解析顺序
func SortXlsx(files []string) []string {
	j := -1
	for i, filename := range files {
		if filepath.Base(filename) == "define.xlsx" {
			j++
			files[i] = files[j]
			files[j] = filename
			break
		}
	}
	return files
}

// 需要规定sheet解析顺序
func SortSheet(list []string) []string {
	j := -1
	for i, name := range list {
		if name == "代对表" {
			j++
			list[i] = list[j]
			list[j] = name
			break
		}
	}
	return list
}
