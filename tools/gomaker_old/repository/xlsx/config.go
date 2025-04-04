package xlsx

import (
	"bytes"
	"fmt"
	"go/format"
	"hego/framework/uerror"
	"hego/tools/gomaker/internal/manager"
	"hego/tools/gomaker/internal/util"
	"path/filepath"
	"strings"
	"text/template"
)

func ConfigGen(dst string, tpls *template.Template, extra ...string) error {
	if strings.HasSuffix(dst, ".go") {
		dst = filepath.Dir(dst)
	}

	buf := bytes.NewBuffer(nil)
	tpl := template.Must(template.New("config").Parse(configTpl))
	for _, ms := range manager.GetMessageList() {
		// 过滤不需要生成的代码
		if !ms.IsList && !ms.IsStruct && len(ms.Group) <= 0 && len(ms.Map) <= 0 {
			continue
		}
		buf.Reset()
		tpl.Execute(buf, ms)
		result, err := format.Source(buf.Bytes())
		if err != nil {
			return uerror.NewUError(1, -1, "%v", err)
		}
		filename := fmt.Sprintf("%s.gen.go", ms.GetPkg())
		if err := util.SaveFile(filepath.Join(dst, ms.GetPkg(), filename), result); err != nil {
			return err
		}
	}
	return nil
}

const (
	configTpl = `
{{$groups := .Group}}
{{$maps := .Map}}
{{$gIndexs := .GetGIndexs}}
{{$mIndexs := .GetMIndexs}}
{{$pbname := .Config}}
{{$st := .Struct}}

package {{.GetPkg}}

import (
	"sync/atomic"
	"hego/common/config/internal/manager"
	"hego/common/pb"
	"hego/framework/uerror"
	"github.com/golang/protobuf/proto"
)

var (
	data atomic.Value
)

{{range $i, $name := $mIndexs}} type {{$name}} struct { 
	{{range $val := index $maps $i}} {{$val.Name}} {{$val.Type.GetType ""}} 
	{{end}} 
} 
{{end}}

{{range $i,$name := $gIndexs}} type {{$name}} struct { 
	{{range $val := index $groups $i}} {{$val.Name}} {{$val.Type.GetType ""}} 
	{{end}} 
} 
{{end}}

type {{$pbname}}Data struct {
	{{if .IsStruct}} cfg *{{$st.Type.GetType ""}}  {{end}} 
	{{if .IsList}} cfgs []*{{$st.Type.GetType ""}} {{end}}
	{{range $name := $mIndexs}} {{$name}} map[{{$name}}]*{{$st.Type.GetType ""}}
	{{end}} {{range $name := $gIndexs}} {{$name}} map[{{$name}}][]*{{$st.Type.GetType ""}}
	{{end}}
}

func init(){
	manager.Register("{{$pbname}}", load)
}

func load(buf []byte) error {
	ary := &{{$st.Type.GetType ""}}Ary{}
	if err := proto.Unmarshal(buf, ary); err != nil {
		return uerror.NewUError(1, -1, "加载{{$pbname}}配置失败: %v", err)	
	}

	dd := &{{$pbname}}Data{
		{{if .IsStruct}} cfg: ary.Ary[0], {{end}}
		{{range $name := $mIndexs}} {{$name}}: make(map[{{$name}}]*{{$st.Type.GetType ""}}),
		{{end}} {{range $name := $gIndexs}} {{$name}}: make(map[{{$name}}][]*{{$st.Type.GetType ""}}),
		{{end}}
	}
	for _, item := range ary.Ary {
		{{if .IsList}} dd.cfgs = append(dd.cfgs, item) 
		{{end}} {{range $i, $name := $mIndexs}} {{$args := index $maps $i}} dd.{{$name}}[{{$name}}{ {{$args.GetParams "item"}} }] = item
		{{end}} {{range $i, $name := $gIndexs}} {{$args := index $groups $i}} dd.{{$name}}[{{$name}}{ {{$args.GetParams "item"}} }] = append(dd.{{$name}}[{{$name}}{ {{$args.GetParams "item"}} }], item)
		{{end}}
	}
	data.Store(dd)
	return nil
}

func getObj() *{{$pbname}}Data {
	if obj, ok := data.Load().(*{{$pbname}}Data); ok && obj != nil {
		return obj
	}
	return nil
}

{{if .IsStruct}}
func Get() *{{$st.Type.GetType ""}} {
	return getObj().cfg
}
{{end}}

{{if .IsList}}
func Gets() (rets []*{{$st.Type.GetType ""}}) {
	list := getObj().cfgs
	rets = make([]*{{$st.Type.GetType ""}}, len(list))
	copy(rets, list)
	return 
}
{{end}}

{{range $i, $name := $mIndexs}} 
{{$args := index $maps $i}}
func GetBy{{$args.GetName}}({{$args.GetArgs ""}}) *{{$st.Type.GetType ""}} {
	return getObj().{{$name}}[{{$name}}{ {{$args.GetParams ""}} }]
}
{{end}}

{{range $i, $name := $gIndexs}} 
{{$args := index $groups $i}}
func GetsBy{{$args.GetName}}({{$args.GetArgs ""}}) (rets []*{{$st.Type.GetType ""}}) {
	list := getObj().{{$name}}[{{$name}}{ {{$args.GetParams ""}} }]
	if len(list) > 0 {
		rets = make([]*{{$st.Type.GetType ""}}, len(list))
		copy(rets, list)
	}
	return
}
{{end}}

`
)
