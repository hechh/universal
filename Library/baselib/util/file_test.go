package util

import (
	"os"
	"testing"
)

func TestFile(t *testing.T) {
	err := os.MkdirAll("./test", os.FileMode(0755))
	t.Log(err)

	err = os.MkdirAll("./test", os.FileMode(0755))
	t.Log(err)
}
