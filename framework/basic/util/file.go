package util

import (
	"os"
	"path"
)

func NewFile(fileName string) (fb *os.File, err error) {
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

func SameFile(fb *os.File, name string) bool {
	st2, err := os.Stat(name)
	if err != nil {
		return false
	}
	st1, _ := fb.Stat()
	return os.SameFile(st1, st2)
}
