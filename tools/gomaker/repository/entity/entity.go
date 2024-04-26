package entity

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"universal/framework/fbasic"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/base"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/typespec"
)

type EntityAttr struct {
	typespec.BaseFunc
	PkgName string
	Name    string
	Field   *typespec.Field
	Struct  *typespec.Struct
}

func Init() {
	manager.Register(&base.Action{
		Name:  domain.ENTITY,
		Help:  "生成entity模板文件",
		Param: `"entity:{XxxEntity}|map:field@pbname"`,
		Gen:   Gen,
	})
}

func parseParams(params string) (*EntityAttr, error) {
	ret := &EntityAttr{}
	strs := strings.Split(params, "|")
	for _, str := range strs {
		pos := strings.Index(str, ":")
		switch str[:pos] {
		case "entity":
			ret.Name = str[pos+1:]
		case "map":
			ipos := strings.Index(str, "@")
			ret.Struct = manager.GetStruct(str[ipos+1:])
			ret.Field = ret.Struct.Fields[str[pos+1:ipos]]
		}
	}
	return ret, nil
}

func Gen(action string, dst string, params string) error {
	attr, err := parseParams(params)
	if err != nil {
		return err
	}
	// 生成文档
	if !strings.HasSuffix(dst, ".go") {
		dst += fmt.Sprintf("/%s.go", attr.Name)
	}
	attr.PkgName = base.GetFilePathBase(dst)
	// 生成包头
	buf := bytes.NewBuffer(nil)
	if err := manager.GenPackage(dst, buf); err != nil {
		return err
	}
	// 模版
	if tpl := manager.GetTpl(action); tpl == nil {
		return fbasic.NewUError(1, -1, fmt.Sprintf("The action of %s is not supported", action))
	} else {
		// 生成文件
		if err := tpl.ExecuteTemplate(buf, action+".tpl", attr); err != nil {
			return fbasic.NewUError(1, -1, err)
		}
	}
	// 格式化
	result, err := format.Source(buf.Bytes())
	if err != nil {
		//ioutil.WriteFile("./gen.go", buf.Bytes(), os.FileMode(0644))
		return fbasic.NewUError(1, -1, err)
	}
	if err := os.MkdirAll(filepath.Dir(dst), os.FileMode(0777)); err != nil {
		return fbasic.NewUError(1, -1, err)
	}
	if err := ioutil.WriteFile(dst, result, os.FileMode(0666)); err != nil {
		return fbasic.NewUError(1, -1, err)
	}
	return nil
}
