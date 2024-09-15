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
func EnumGen(dst string, tpls *template.Template, extra ...string) error {
	if strings.HasSuffix(dst, ".go") {
		dst = filepath.Dir(dst)
	}
	// 生成message信息
	list := []string{}
	for _, st := range manager.GetEnumList() {
		// 排序
		sort.Slice(st.List, func(i, j int) bool { return st.List[i].Value < st.List[j].Value })
		tmps := []string{}
		for _, val := range st.List {
			if len(val.Doc) > 0 {
				tmps = append(tmps, fmt.Sprintf("\t %s = %d; // %s", val.Name, val.Value, val.Doc))
			} else {
				tmps = append(tmps, fmt.Sprintf("\t %s = %d;", val.Name, val.Value))
			}
		}
		if len(st.Doc) > 0 {
			list = append(list, fmt.Sprintf("// %s\nenum %s {\n%s\n}", st.Doc, st.Type.Name, strings.Join(tmps, "\n")))
		} else {
			list = append(list, fmt.Sprintf("enum %s {\n%s\n}", st.Type.Name, strings.Join(tmps, "\n")))
		}
	}
	// 生成包头
	buf := bytes.NewBufferString(`syntax = "proto3";`)
	buf.WriteByte('\n')
	buf.WriteString(`package pb;`)
	buf.WriteByte('\n')
	buf.WriteString(`option go_package = "../../common/pb";`)
	buf.WriteString("\n\n")
	buf.WriteString(strings.Join(list, "\n\n"))
	// 生成文件
	return util.SaveFile(filepath.Join(dst, "enum.gen.proto"), buf.Bytes())
}

func TableGen(dst string, tpls *template.Template, extra ...string) error {
	if strings.HasSuffix(dst, ".go") {
		dst = filepath.Dir(dst)
	}
	// 生成message信息
	list := []string{}
	for _, st := range manager.GetStructList() {
		tmps := []string{}
		for _, ff := range st.List {
			if len(ff.Token) > 0 {
				if len(ff.Doc) > 0 {
					tmps = append(tmps, fmt.Sprintf("\t repeated %s %s = %d; // %s", ff.Type.Name, ff.Name, ff.Index+1, ff.Doc))
				} else {
					tmps = append(tmps, fmt.Sprintf("\t repeated %s %s = %d;", ff.Type.Name, ff.Name, ff.Index+1))
				}
			} else {
				if len(ff.Doc) > 0 {
					tmps = append(tmps, fmt.Sprintf("\t %s %s = %d; // %s", ff.Type.Name, ff.Name, ff.Index+1, ff.Doc))
				} else {
					tmps = append(tmps, fmt.Sprintf("\t %s %s = %d;", ff.Type.Name, ff.Name, ff.Index+1))
				}
			}
		}
		if len(st.Doc) > 0 {
			list = append(list, fmt.Sprintf("// %s \nmessage %s {\n%s\n}", st.Doc, st.Type.Name, strings.Join(tmps, "\n")))
		} else {
			list = append(list, fmt.Sprintf("message %s {\n%s\n}", st.Type.Name, strings.Join(tmps, "\n")))
		}
		list = append(list, fmt.Sprintf("message %sAry {\nrepeated %s Ary = 1;\n}", st.Type.Name, st.Type.Name))
	}
	// 生成包头
	buf := bytes.NewBufferString(`syntax = "proto3";`)
	buf.WriteByte('\n')
	buf.WriteString(`package pb;`)
	buf.WriteByte('\n')
	buf.WriteString(`import "enum.gen.proto";`)
	buf.WriteByte('\n')
	buf.WriteString(`option go_package = "../../common/pb";`)
	buf.WriteString("\n\n")
	buf.WriteString(strings.Join(list, "\n\n"))
	// 生成文件
	return util.SaveFile(filepath.Join(dst, "table.gen.proto"), buf.Bytes())
}
