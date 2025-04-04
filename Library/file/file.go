package file

import (
	"go/format"
	"hego/Library/uerror"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func NewOrOpen(fileName string) (fb *os.File, err error) {
	// 判断路径是否存在
	pp := path.Dir(fileName)
	if _, err = os.Stat(pp); os.IsNotExist(err) {
		if err := os.MkdirAll(pp, os.FileMode(0755)); err != nil {
			return nil, err
		}
	}
	// 创建文件
	if fb, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644); err != nil {
		return nil, err
	}
	return
}

func IsSame(fb *os.File, name string) bool {
	st2, err := os.Stat(name)
	if err != nil {
		return false
	}
	st1, _ := fb.Stat()
	return os.SameFile(st1, st2)
}

func Save(path, filename string, buf []byte) error {
	fileName := filepath.Join(path, filename)
	// 创建目录
	if err := os.MkdirAll(filepath.Dir(fileName), os.FileMode(0777)); err != nil {
		return uerror.New(1, -1, "filename: %s, error: %v", fileName, err)
	}

	// 写入文件
	if err := ioutil.WriteFile(fileName, buf, os.FileMode(0666)); err != nil {
		return uerror.New(1, -1, "filename: %s, error: %v", fileName, err)
	}
	return nil
}

func SaveGo(path, filename string, buf []byte) error {
	result, err := format.Source(buf)
	if err != nil {
		return uerror.New(1, -1, "格式化失败: %v", err)
	}
	return Save(path, filename, result)
}
