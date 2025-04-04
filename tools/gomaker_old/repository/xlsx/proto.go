package xlsx

import (
	"bytes"
	"fmt"
	"hego/tools/gomaker/domain"
	"hego/tools/gomaker/internal/manager"
	"hego/tools/gomaker/internal/util"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

const (
	protoHeader = `
syntax = "proto3";
package pb;
option  go_package = "../../common/pb";
`
)

func getDoc(doc string, isnewline bool) string {
	fmtstr := "// %s"
	if isnewline {
		fmtstr = "// %s\n"
	}
	if len(doc) > 0 {
		return fmt.Sprintf(fmtstr, doc)
	}
	return doc
}

// 生成enum.gen.proto
func EnumGen(dst string, tpls *template.Template, extra ...string) error {
	if strings.HasSuffix(dst, ".go") {
		dst = filepath.Dir(dst)
	}
	// 生成message信息
	class := map[string][]string{}
	for _, st := range manager.GetEnumList() {
		// 排序
		sort.Slice(st.List, func(i, j int) bool { return st.List[i].Value < st.List[j].Value })
		tmps := []string{}
		for _, val := range st.List {
			tmps = append(tmps, fmt.Sprintf("\t%s = %d; %s", val.Name, val.Value, getDoc(val.Doc, false)))
		}
		if _, ok := class[st.Type.Class]; !ok {
			class[st.Type.Class] = []string{}
		}
		class[st.Type.Class] = append(class[st.Type.Class], fmt.Sprintf("%senum %s {\n%s\n}", getDoc(st.Doc, true), st.Type.Name, strings.Join(tmps, "\n")))
	}
	// 生成包头
	buf := bytes.NewBuffer(nil)
	for cl, list := range class {
		buf.Reset()
		buf.WriteString(protoHeader)
		buf.WriteString("\n\n")
		buf.WriteString(strings.Join(list, "\n\n"))
		// 生成文件
		if err := util.SaveFile(filepath.Join(dst, fmt.Sprintf("%s.gen.proto", cl)), buf.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

func TableGen(dst string, tpls *template.Template, extra ...string) error {
	if strings.HasSuffix(dst, ".go") {
		dst = filepath.Dir(dst)
	}
	// 生成message信息
	class := map[string][]string{}
	refs := map[string]struct{}{}
	for _, st := range manager.GetStructList() {
		if len(st.Type.Class) <= 0 {
			continue
		}
		tmps := []string{}
		for _, ff := range st.List {
			// 获取枚举类型分类，用于import
			ttt := manager.GetTypeReference(ff.Type)
			if ttt.Kind != domain.KindTypeIdent {
				refs[ttt.Class] = struct{}{}
			}
			if len(ff.Token) > 0 {
				tmps = append(tmps, fmt.Sprintf("\trepeated %s %s = %d; %s", ff.Type.Name, ff.Name, ff.Index+1, getDoc(ff.Doc, false)))
			} else {
				tmps = append(tmps, fmt.Sprintf("\t%s %s = %d; %s", ff.Type.Name, ff.Name, ff.Index+1, getDoc(ff.Doc, false)))
			}
		}
		if _, ok := class[st.Type.Class]; !ok {
			class[st.Type.Class] = []string{}
		}
		class[st.Type.Class] = append(class[st.Type.Class], fmt.Sprintf("%smessage %s {\n%s\n}", getDoc(st.Doc, true), st.Type.Name, strings.Join(tmps, "\n")))
		class[st.Type.Class] = append(class[st.Type.Class], fmt.Sprintf("message %sAry {\n\trepeated %s Ary = 1;\n}", st.Type.Name, st.Type.Name))
	}
	// 生成包头
	buf := bytes.NewBuffer(nil)
	for cl, list := range class {
		buf.Reset()
		buf.WriteString(protoHeader)
		for name := range refs {
			buf.WriteString(fmt.Sprintf("import \"%s.gen.proto\";\n", name))
		}
		buf.WriteString("\n\n")
		buf.WriteString(strings.Join(list, "\n\n"))
		// 生成文件
		if err := util.SaveFile(filepath.Join(dst, fmt.Sprintf("%s.gen.proto", cl)), buf.Bytes()); err != nil {
			return err
		}
	}
	return nil
}
