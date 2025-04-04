package basic

import (
	"os"
	"path/filepath"
	"regexp"
)

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
