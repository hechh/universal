package generator

import (
	"bytes"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"universal/tools/gomaker_new/internal/manager"
	"universal/tools/gomaker_new/internal/util"
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
		buf.WriteRune('\n')
		buf.WriteString(st.Proto())
		buf.WriteRune('\n')
	}
	// 生成文件
	return util.SaveFile(filepath.Join(dst, "enum.gen.proto"), buf.Bytes())
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
		buf.WriteRune('\n')
		buf.WriteString(st.Proto())
		buf.WriteRune('\n')
	}
	// 生成文件
	return util.SaveFile(filepath.Join(dst, "table.gen.proto"), buf.Bytes())
}
