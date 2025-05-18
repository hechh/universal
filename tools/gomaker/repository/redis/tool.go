package redis

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"forevernine.com/planet/server/tool/gomaker/internal/base"
	"forevernine.com/planet/server/tool/gomaker/internal/manager"
)

type RedisMockRule struct {
	PackageName string // 包名
	PBName      string // 协议名
	Desc        string // 注释
	Args        string // 参数说明
	IsHash      bool   // 是否为hash
}

func genTool(rule, path string, buf *bytes.Buffer) {
	if len(path) <= 0 {
		path = base.ROOT
	}

	rules := []*RedisMockRule{}
	for _, item := range manager.GetRules(RuleTypeString) {
		val, ok := item.(*RedisAttr)
		if !ok || val == nil {
			panic(fmt.Sprintf("data is empty or type is error, data: %v", val))
		}
		rules = append(rules, &RedisMockRule{
			PackageName: val.Package,
			PBName:      val.Name,
			Args:        val.Key.Values.Arg(),
			Desc:        val.Desc,
		})
	}
	for _, item := range manager.GetRules(RuleTypeHash) {
		val, ok := item.(*RedisAttr)
		if !ok || val == nil {
			panic(fmt.Sprintf("data is empty or type is error, data: %v", val))
		}
		rules = append(rules, &RedisMockRule{
			PackageName: val.Package,
			PBName:      val.Name,
			Args:        val.Key.Values.Join(val.Field.Values).Arg(),
			Desc:        val.Desc, //val.Desc,
			IsHash:      true,
		})
	}
	sort.Slice(rules, func(i, j int) bool {
		return strings.Compare(rules[i].PBName, rules[j].PBName) <= 0
	})
	manager.Execute(Action, "tool.tpl", buf, rules)
	base.GenGo(buf, fmt.Sprintf(toolFile, path), true)
	buf.Reset()
}
