package test

import (
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"universal/framework/basic"
)

func TestRand(t *testing.T) {
	t.Log(basic.RangeInt63n(1, 2))
	t.Log(rand.Int())
}

func TestWalk(t *testing.T) {
	root := "../"
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if root != path {
				return filepath.SkipDir
			}
		}
		t.Log(path, info.Name(), info.IsDir())
		return nil
	})
}
