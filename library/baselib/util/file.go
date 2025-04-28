package util

import (
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// os.O_CREATE|os.O_APPEND|os.O_RDWR
func NewOrOpenFile(fileName string, flag int) (fb *os.File, err error) {
	// 判断路径是否存在
	pp := path.Dir(fileName)
	if err := os.MkdirAll(pp, os.FileMode(0755)); err != nil {
		return nil, err
	}
	// 创建文件
	if fb, err = os.OpenFile(fileName, flag, os.FileMode(0644)); err != nil {
		return nil, err
	}
	return
}

// 文件是否相同
func IsSameFile(fb *os.File, name string) bool {
	st2, _ := os.Stat(name)
	st1, _ := fb.Stat()
	return os.SameFile(st1, st2)
}

// 保存文件
func SaveFile(path, filename string, buf []byte) (err error) {
	fileName := filepath.Join(path, filename)
	// 是否为go文件
	if strings.HasSuffix(fileName, ".go") {
		if buf, err = format.Source(buf); err != nil {
			return
		}
	}
	// 创建目录
	if err = os.MkdirAll(filepath.Dir(fileName), os.FileMode(0777)); err != nil {
		return err
	}
	// 写入文件
	return ioutil.WriteFile(fileName, buf, os.FileMode(0666))
}

// 遍历目录所有文件
func Glob(dir, pattern string, recursive bool) (rets []string, err error) {
	pre, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// 不深度迭代
		if !recursive && info.IsDir() && dir != path {
			return filepath.SkipDir
		}
		// 过滤目录
		if info.IsDir() {
			return nil
		}
		// 是否配置
		if pre.MatchString(path) {
			rets = append(rets, path)
		}
		return nil
	})
	return
}
