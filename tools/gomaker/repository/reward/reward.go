package reward

import (
	"bytes"
	"fmt"
	"strings"

	"forevernine.com/planet/server/tool/gomaker/internal/base"

	"forevernine.com/planet/server/tool/gomaker/internal/manager"
)

const (
	Action               = "reward"
	RuleTypeRewardString = "//@gomaker:reward:string"
	RuleTypeRewardHash   = "//@gomaker:reward:hash"
	goFile               = "%s/%s/%s.go"
)

func Init() {
	manager.RegisterAction(Action, RuleTypeRewardString, RuleTypeRewardHash)
	manager.RegisterParser(RuleTypeRewardString, func(pbname, comment string) interface{} {
		return parse(false, pbname, comment)
	})
	manager.RegisterParser(RuleTypeRewardHash, func(pbname, comment string) interface{} {
		return parse(true, pbname, comment)
	})

	manager.RegisterCreator(RuleTypeRewardString, gen)
	manager.RegisterCreator(RuleTypeRewardHash, gen)
}

type Attribute struct {
	IsHash        bool
	Package       string
	DaoPackage    string
	DataHis       string
	Data          string
	ActivityType  string
	PropertyType  string
	PropertyTypes []string
}

// @gomaker:reward:string|data|ActivityType|PropertyType,...
// @gomaker:reward:hash|data|ActivityType|PropertyType,...
func parse(flag bool, pbname, comment string) interface{} {
	arrs, comment, _ := base.RuleSplit(comment)
	name := base.ToUnderline(pbname)
	ret := &Attribute{IsHash: flag, Package: name + "_reward", DaoPackage: name, DataHis: pbname}
	switch len(arrs) {
	case 2:
		ret.PropertyTypes = strings.Split(comment[arrs[0]+1:arrs[1]], ",")
	case 3:
		if strings.Contains(comment, "ActivityType") {
			ret.ActivityType = comment[arrs[0]+1 : arrs[1]]
		} else {
			ret.Data = comment[arrs[0]+1 : arrs[1]]
		}
		ret.PropertyTypes = strings.Split(comment[arrs[1]+1:arrs[2]], ",")
	case 4:
		ret.Data = comment[arrs[0]+1 : arrs[1]]
		ret.ActivityType = comment[arrs[1]+1 : arrs[2]]
		ret.PropertyTypes = strings.Split(comment[arrs[2]+1:arrs[3]], ",")
	}
	return ret
}

func gen(rule, path string, buf *bytes.Buffer) {
	if len(path) <= 0 {
		path = base.ROOT
	}

	for _, v := range manager.GetRules(rule) {
		val, ok := v.(*Attribute)
		if !ok || val == nil {
			panic(fmt.Sprintf("data is empty or type is error, data: %v", v))
		}
		manager.Execute(Action, "reward_api.tpl", buf, val)
		base.GenGo(buf, fmt.Sprintf(goFile, path, val.Package, "reward_api"), true)
		buf.Reset()

		manager.Execute(Action, "reward_entity.tpl", buf, val)
		base.GenGo(buf, fmt.Sprintf(goFile, path, val.Package, "reward_entity"), true)
		buf.Reset()

		for _, pp := range val.PropertyTypes {
			val.PropertyType = pp
			manager.Execute(Action, "reward_property.tpl", buf, val)
			base.GenGo(buf, fmt.Sprintf(goFile, path, val.Package, base.ToUnderline(pp)), true)
			buf.Reset()
		}
	}
}
