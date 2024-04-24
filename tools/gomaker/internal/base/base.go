package base

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetAbsPath(src string, root string) (string, error) {
	if len(src) <= 0 {
		return "", fmt.Errorf("Relative path is empty")
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
