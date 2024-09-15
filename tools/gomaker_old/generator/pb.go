package generator

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/util"
)

// 生成enum.gen.proto
func enumGenerator(dst string, tpls *template.Template, extra ...string) error {
	if strings.HasSuffix(dst, ".go") {
		dst = filepath.Dir(dst)
	}
	// 生成包头
	buf := bytes.NewBufferString(`
syntax = "proto3";
package pb;
option go_package = "../../common/pb";
	`)
	// 注册
	for _, st := range manager.GetEnumList() {
		sort.Slice(st.List, func(i, j int) bool { return st.List[i].Value < st.List[j].Value })
		arrs := []string{}
		for _, val := range st.List {
			arrs = append(arrs, fmt.Sprintf("\t%s = %d; // %s", val.Name, val.Value, val.Doc))
		}
		buf.WriteString(fmt.Sprintf("\nenum %s {\n%s\n}\n", st.Type.Name, strings.Join(arrs, "\n")))
	}
	// 生成文件
	return util.SaveFile(filepath.Join(dst, "enum.gen.proto"), buf)
}

func tableGenerator(dst string, tpls *template.Template, extra ...string) error {
	if strings.HasSuffix(dst, ".go") {
		dst = filepath.Dir(dst)
	}
	// 生成包头
	buf := bytes.NewBufferString(`
syntax = "proto3";
package pb;
import "enum.gen.proto";
option go_package = "../../common/pb";
	`)
	// 注册
	for _, st := range manager.GetStructList() {
		arrs := []string{}
		for i, val := range st.List {
			arrs = append(arrs, fmt.Sprintf("\t%s %s = %d; // %s", manager.GetProtoType(val.Type.GetPkgType()), val.Name, i+1, val.Doc))
		}
		buf.WriteString(fmt.Sprintf("\nmessage %s {\n%s\n}\n", st.Type.Name, strings.Join(arrs, "\n")))
		buf.WriteString(fmt.Sprintf("\nmessage %sAry {\n repeated %s Ary = 1;\n}\n", st.Type.Name, st.Type.Name))
	}
	// 生成文件
	return util.SaveFile(filepath.Join(dst, "table.gen.proto"), buf)
}
