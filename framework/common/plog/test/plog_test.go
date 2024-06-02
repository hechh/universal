package test

import (
	"fmt"
	"path"
	"runtime"
	"testing"
)

func Print() {
	pc, file, line, _ := runtime.Caller(0)
	funcName := runtime.FuncForPC(pc).Name()
	fmt.Println(file, line, path.Base(funcName))
}

func TestLog(t *testing.T) {
	Print()
}
