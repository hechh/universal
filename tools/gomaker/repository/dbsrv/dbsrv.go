package dbsrv

import (
	"bytes"
	"sort"

	"forevernine.com/planet/server/tool/gomaker/internal/manager"

	"forevernine.com/planet/server/tool/gomaker/domain"
	"forevernine.com/planet/server/tool/gomaker/internal/base"
)

const (
	Action        = "dbsrv"
	RuleTypeDbsrv = "dbsrv"
	goFile        = "srv/dbsrv/internal/biz/proc/internal/dbsrv.go"
)

func Init() {
	manager.RegisterAction(Action, RuleTypeDbsrv)
	manager.RegisterCreator(RuleTypeDbsrv, genDbsrv)
}

type Field struct {
	Name  string
	Type  string
	Value string
}

func (d *Field) GetName() string {
	return d.Name
}

func (d *Field) GetHanldeFuncName() string {
	return base.FirstToBig(d.Name)
}

func (d *Field) GetType() string {
	return d.Type
}

func (d *Field) GetValue() string {
	return d.Value
}

type Attribute struct {
	Package string
	Sets    []*Field
}

func (d *Attribute) Push(elem *Field) {
	d.Sets = append(d.Sets, elem)
}

func genDbsrv(rule, path string, buf *bytes.Buffer) {
	if len(path) <= 0 {
		path = base.ROOT
	}

	attrs := &Attribute{Package: "internal"}
	manager.WalkAstEnum("routeKeyType", func(item *domain.AstValue) {
		attrs.Sets = append(attrs.Sets, &Field{
			Name:  item.Name,
			Type:  item.Type.Name,
			Value: item.StrVal,
		})
	})
	sort.Slice(attrs.Sets, func(i, j int) bool {
		return attrs.Sets[i].Name < attrs.Sets[j].Name
	})

	manager.Execute(Action, "", buf, attrs)
	base.GenGo(buf, path+"/"+goFile, true)
}
