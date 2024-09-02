package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"universal/framework/basic/uerror"
	"universal/tools/cfgtool/domain"
)

// 保存json文件
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

func Glob(dir, suffix string, recursive bool) (files []string, err error) {
	if !recursive {
		files, err = filepath.Glob(filepath.Join(dir, "*.go"))
	} else {
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				return nil
			}
			patterns, err := filepath.Glob(filepath.Join(path, suffix))
			if err != nil {
				return uerror.NewUError(1, -1, "%v", err)
			}
			if len(patterns) > 0 {
				files = append(files, patterns...)
			}
			return nil
		})
	}
	return
}

func NewFileInfo(filename string, alls map[string]*domain.Enum) *domain.FileType {
	ext := filepath.Ext(filename)
	name := filepath.Base(filename)
	return &domain.FileType{
		Name:   strings.TrimSuffix(name, ext),
		Enums:  make(map[string][]*domain.Enum),
		Alls:   alls,
		Tables: make(map[string]*domain.Table),
	}
}
