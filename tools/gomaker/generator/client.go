package generator

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/util"
)

// 生成代码
func HttpKitGenerator(dst string, param string, tpls *template.Template) error {
	if strings.HasSuffix(dst, ".go") {
		dst = filepath.Dir(dst)
	}
	// 生成包头
	buf := bytes.NewBufferString(fmt.Sprintf("package %s\n", filepath.Base(dst)))
	// 包引用
	buf.WriteString("\nimport (\"net/http\")\n")
	// 注册
	arrs := []string{}
	for _, st := range manager.GetStructList() {
		name := st.Type.Name
		if strings.HasSuffix(name, "Request") || strings.HasSuffix(name, "Response") || strings.HasSuffix(name, "Notify") {
			arrs = append(arrs, fmt.Sprintf("http.HandleFunc(\"/api/%s\", handle)", name))
		}
	}
	buf.WriteString(fmt.Sprintf("func init(){\n%s\n}\n", strings.Join(arrs, "\n")))
	// 生成文件
	return util.SaveGo(filepath.Join(dst, "init.gen.go"), buf)
}

func OmitEmptyGenerator(dst string, param string, tpls *template.Template) error {
	if strings.HasSuffix(dst, ".go") {
		dst = filepath.Dir(dst)
	}
	// 生成包头
	buf := bytes.NewBufferString(fmt.Sprintf("package %s\n", filepath.Base(dst)))
	// 生成结构体数据
	arrs := []string{}
	for _, tmp := range manager.GetStructList() {
		arrs = append(arrs, strings.ReplaceAll(tmp.Clone().Format(), ",omitempty", ""))
	}
	for _, tmp := range manager.GetEnumList() {
		arrs = append(arrs, tmp.Format())
	}
	buf.WriteString(strings.Join(arrs, "\n"))
	// 生成文件
	return util.SaveGo(filepath.Join(dst, "pb.gen.go"), buf)
}

func ProtoGenerator(dst string, param string, tpls *template.Template) error {
	if strings.HasSuffix(dst, ".go") {
		dst = filepath.Dir(dst)
	}
	// 生成包头
	buf := bytes.NewBufferString(fmt.Sprintf("package %s\n", filepath.Base(dst)))
	// 生成注册信息
	arrs := []string{}
	for _, st := range manager.GetStructList() {
		name := st.Type.Name
		if strings.HasSuffix(name, "Request") {
			arrs = append(arrs, fmt.Sprintf("registerJson(\"%s\", &%s{})", name, name))
		}
	}
	buf.WriteString(fmt.Sprintf("func init(){\n%s\n}", strings.Join(arrs, "\n")))
	// 生成文件
	return util.SaveGo(filepath.Join(dst, "json.gen.go"), buf)
}
