package manager

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"forevernine.com/planet/server/tool/gomaker/internal/base"
)

var (
	_templs = make(map[string]*template.Template)
)

func Execute(action, name string, buf *bytes.Buffer, data interface{}) {
	val, ok := _templs[action]
	if !ok {
		panic(fmt.Sprintf("%s is not found", action))
	}
	if len(name) <= 0 {
		base.Must(val.Execute(buf, data))
		return
	}
	if !strings.Contains(name, ".tpl") {
		name = name + ".tpl"
	}
	if tt := val.Lookup(name); tt != nil {
		base.Must(tt.Execute(buf, data))
		return
	}
	panic(fmt.Sprintf("%s(%s) is not found", action, name))
}

func init() {
	root := base.AbsPath("../../templates/")
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() || path == root {
			return nil
		}
		_templs[info.Name()] = template.Must(template.ParseGlob(path + "/*.tpl"))
		return nil
	})
}
