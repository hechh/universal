package base

import (
	"path/filepath"
)

func GetAbsPath(src string, root string) string {
	if len(src) <= 0 {
		return src
	}
	if !filepath.IsAbs(src) {
		return filepath.Join(root, src)
	}
	return src
}

func GetPathDefault(dst string, defaultEnv string) string {
	if len(dst) > 0 {
		return dst
	}
	return defaultEnv
}

func GetFilePathBase(dst string) string {
	ext := filepath.Ext(dst)
	if len(ext) <= 0 {
		return filepath.Base(dst)
	}
	return filepath.Base(filepath.Dir(dst))
}
