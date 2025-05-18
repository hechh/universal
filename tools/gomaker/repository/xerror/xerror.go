package xerror

import (
	"bytes"
	"path/filepath"
	"sort"
	"strings"

	"forevernine.com/planet/server/tool/gomaker/domain"
	"forevernine.com/planet/server/tool/gomaker/internal/base"
	"forevernine.com/planet/server/tool/gomaker/internal/manager"
	"github.com/spf13/cast"
)

const (
	Rule   = "xerror"
	Action = "xerror"
	goFile = "common/xerrors/errors_auto.gen.go"
)

func Init() {
	manager.RegisterAction(Action, Rule)
	manager.RegisterCreator(Rule, creator)
}

type Attribute struct {
	FName  string
	ErrMsg string
	Code   int32
}

func creator(rule, path string, buf *bytes.Buffer) {
	if len(path) <= 0 {
		path = base.ROOT
	}
	data01 := []*Attribute{}
	manager.WalkAstEnum("BASE_ERROR_CODE", func(item *domain.AstValue) {
		data01 = append(data01, &Attribute{
			FName:  getFName(item.Name),
			ErrMsg: getErrMsg(item.Name),
			Code:   cast.ToInt32(item.Value),
		})
	})
	sort.Slice(data01, func(i, j int) bool {
		return data01[i].Code < data01[j].Code
	})
	data02 := []*Attribute{}
	manager.WalkAstEnum("DB_ERROR_CODE", func(item *domain.AstValue) {
		data02 = append(data02, &Attribute{
			FName:  getFName(item.Name),
			ErrMsg: getErrMsg(item.Name),
			Code:   cast.ToInt32(item.Value),
		})
	})
	sort.Slice(data02, func(i, j int) bool {
		return data02[i].Code < data02[j].Code
	})
	data01 = append(data01, data02...)
	data03 := []*Attribute{}
	manager.WalkAstEnum("ERROR", func(item *domain.AstValue) {
		data03 = append(data03, &Attribute{
			FName:  getFName(item.Name),
			ErrMsg: getErrMsg(item.Name),
			Code:   cast.ToInt32(item.Value),
		})
	})
	sort.Slice(data03, func(i, j int) bool {
		return data03[i].Code < data03[j].Code
	})
	data01 = append(data01, data03...)

	manager.Execute(Action, "", buf, data01)
	base.GenGo(buf, filepath.Join(path, goFile), true)
}

func getFName(a string) string {
	x := strings.Replace(a, "ERR_CODE_", "ERR_", 1)
	result := []string{}
	for _, a := range strings.Split(x, "_") {
		result = append(result, base.ToBigHump(a))
	}
	return strings.Join(result, "")
}

func getErrMsg(x string) string {
	a := strings.TrimPrefix(x, "ERR_")
	a = strings.Replace(a, "_", " ", -1)
	return strings.ToLower(a)
}
