package entity

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"universal/common/pb"
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
		Param: "{XxxEntity}:field@pbname|...",
		Gen:   Gen,
	})
}

func parseParams(params string) (rets []*EntityAttr, err error) {
	for _, str := range strings.Split(params, "|") {
		pos1 := strings.Index(str, ":")
		pos2 := strings.Index(str, "@")
		st := manager.GetStruct(str[pos2+1:])
		if st == nil {
			return nil, fbasic.NewUError(1, pb.ErrorCode_NotFound, fmt.Sprintf("Struct not found, %s", str))
		}
		field := st.Fields[str[pos1+1:pos2]]
		if st == nil {
			return nil, fbasic.NewUError(1, pb.ErrorCode_NotFound, fmt.Sprintf("Field not found, %s", str))
		}

		rets = append(rets, &EntityAttr{
			Name:   str[:pos1],
			Struct: st,
			Field:  field,
		})
	}
	return
}

func Gen(action string, dst string, params string) error {
	attrs, err := parseParams(params)
	if err != nil {
		return err
	}
	// 生成文档
	if !strings.HasSuffix(dst, ".go") {
		dst += fmt.Sprintf("/%s.go", attrs[0].Name)
	}
	// 生成包头
	buf := bytes.NewBuffer(nil)
	if err := manager.GenPackage(dst, buf); err != nil {
		return err
	}
	// 循环生成
	pkgName := base.GetFilePathBase(dst)
	for _, attr := range attrs {
		attr.PkgName = pkgName
		// 模版
		if tpl := manager.GetTpl(action); tpl == nil {
			return fbasic.NewUError(1, -1, fmt.Sprintf("The action of %s is not supported", action))
		} else {
			// 生成文件
			if err := tpl.ExecuteTemplate(buf, action+".tpl", attr); err != nil {
				return fbasic.NewUError(1, -1, err)
			}
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
