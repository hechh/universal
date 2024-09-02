package util

import (
	"os"
	"path/filepath"
	"universal/framework/basic/uerror"
)

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
