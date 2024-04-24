package base

import (
	"fmt"
	"os"
	"path/filepath"
	"universal/framework/basic"
)

func GetAbsPath(src string, root string) (string, error) {
	if len(src) <= 0 {
		return "", basic.NewUError(2, -1, fmt.Sprintf("Relative path is empty"))
	}
	if !filepath.IsAbs(src) {
		return filepath.Join(root, src), nil
	}
	return src, nil
}

func GetPathDefault(dst string, defaultEnv string) string {
	if len(dst) > 0 {
		return dst
	}
	return os.Getenv(defaultEnv)
}

func GetFilePathBase(dst string) string {
	ext := filepath.Ext(dst)
	if len(ext) <= 0 {
		return filepath.Base(dst)
	}
	return filepath.Base(filepath.Dir(dst))
}
