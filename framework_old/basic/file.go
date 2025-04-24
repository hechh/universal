package basic

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
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

// 遍历目录所有文件
func Glob(dir, pattern, filter string, recursive bool) (rets []string, err error) {
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
	// 过滤无用
	if len(filter) > 0 {
		fre, err := regexp.Compile(filter)
		if err != nil {
			return nil, err
		}
		j := -1
		for _, val := range rets {
			if fre.MatchString(val) {
				continue
			}
			j++
			rets[j] = val
		}
		rets = rets[:j+1]
	}
	return
}
