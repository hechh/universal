package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"universal/framework/basic/uerror"
)

func Search(dir string, pattern string) (files []string, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		patterns, err := filepath.Glob(filepath.Join(path, pattern))
		if err != nil {
			return uerror.NewUError(1, -1, "%v", err)
		}
		if len(patterns) > 0 {
			files = append(files, patterns...)
		}
		return nil
	})
	return
}

func SaveJson(dst string, jsons []map[string]interface{}) error {
	result, err := json.MarshalIndent(&jsons, "", "  ")
	if err != nil {
		return err
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
