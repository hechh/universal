package xlsx

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/util"
)

const (
	configHeader = `
import (
	"sync/atomic"
	"universal/common/config/internal/manager"
	"universal/common/pb"
	"universal/framework/uerror"
	"github.com/golang/protobuf/proto"
)

var (
	data atomic.Value
)	
	`
)

func ConfigGen(dst string, tpls *template.Template, extra ...string) error {
	if strings.HasSuffix(dst, ".go") {
		dst = filepath.Dir(dst)
	}

	buf := bytes.NewBuffer(nil)
	for _, ms := range manager.GetMessageList() {
		// 过滤不需要生成的代码
		if !ms.IsList && !ms.IsStruct && len(ms.Group) <= 0 && len(ms.Map) <= 0 {
			continue
		}
		pkg := util.ToUnderline(ms.Config)
		buf.Reset()
		buf.WriteString(fmt.Sprintf("package %s\n", pkg))
		buf.WriteString(configHeader)
		buf.WriteString(fmt.Sprintf("func init(){\n manager.Register(\"%s\", load%s)}\n", pkg, ms.Config))

	}

	return nil
}

const (
	configTpl = `
package {{.GetPkg}}

import (
	"sync/atomic"
	"universal/common/config/internal/manager"
	"universal/common/pb"
	"universal/framework/uerror"
	"github.com/golang/protobuf/proto"
)

var (
	data atomic.Value
)

type {{.Config}}Parser struct {
	{{if .IsStruct}} cfg *pb.{{.Config}} 
	{{end}} {{if .IsList}} cfgs []*pb.{{.Config}} 
	{{end}}
}

func init(){
	manager.Register("{{.Config}}", load{{.Config}})
}


	
	`
)
