package secure

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"forevernine.com/planet/server/tool/gomaker/domain"
	"forevernine.com/planet/server/tool/gomaker/internal/manager"

	"forevernine.com/planet/server/tool/gomaker/internal/base"
)

const (
	Action               = "secure"
	RuleTypeSecure       = "secure"
	RuleTypeSecureConfig = "//@gomaker:pbconfig:secure"
	RuleTypeRange        = "//@gomaker:pbconfig:range"
	goFile               = "%s/common/pbconfig/repository/config/%s/%s/security.gen.go"
)

func Init() {
	manager.RegisterAction(Action, RuleTypeSecure, RuleTypeSecureConfig, RuleTypeRange)
	manager.RegisterParser(RuleTypeRange, parseRange)
	manager.RegisterParser(RuleTypeSecureConfig, parseCfg)
	manager.RegisterCreator(RuleTypeSecure, genSafe)
	manager.RegisterCreator(RuleTypeSecureConfig, genCfg)
}

func parseRange(pbname, rule string) interface{} {
	arrs, comment, _ := base.RuleSplit(rule)
	rules := []*base.Index{}
	for i := 0; i < len(arrs)-1; i++ {
		rules = append(rules, base.ParseIndex(comment[arrs[i]+1:arrs[i+1]]))
	}
	return rules
}

func genSafe(rule, path string, buf *bytes.Buffer) {
	if len(path) <= 0 {
		path = base.ROOT
	}
	buf.WriteString("package pbclass\n\n")

	isGen := false
	manager.WalkConfig(func(item *domain.AstStruct) {
		manager.Execute(Action, "security.tpl", buf, item)
		isGen = true
	})
	if isGen {
		base.GenGo(buf, filepath.Join(path, "common", "pbclass", "security.gen.go"), true)
	}
}

func genCfg(rule, path string, buf *bytes.Buffer) {
	if len(path) <= 0 {
		path = base.ROOT
	}

	for _, item := range manager.GetRules(rule) {
		vals, ok := item.(*Attribute)
		if !ok || vals == nil {
			panic(fmt.Sprintf("data is empty or type is error, data: %v", vals))
		}

		pbname := strings.TrimSuffix(vals.PBName, "S")
		if iii, ok := manager.GetRule(RuleTypeRange, pbname).([]*base.Index); ok && iii != nil {
			vals.Range = iii
		}
		vals.Ast = manager.GetAstStruct(pbname)

		manager.Execute(Action, "pbconfigS.tpl", buf, vals)
		base.GenGo(buf, fmt.Sprintf(goFile, path, vals.ModuleName, vals.PackageName), true)
		buf.Reset()
	}
}

type Attribute struct {
	ModuleName  string            // 模块名称
	PackageName string            // 包名
	SheetName   string            // 表明
	PBName      string            // pb名字
	IsList      bool              // 是否为List结构配置
	IsStruct    bool              // 是否为struct结构配置
	IsCheck     bool              // 是否需要配置检测
	Priority    string            // 配置加载顺序
	Map         []*base.Index     // map结构配置
	Group       []*base.Index     // group配置结构
	Max         []*base.Index     // 极大值
	Min         []*base.Index     // 极小值
	Sum         []*base.Index     // 求和
	Range       []*base.Index     // 范围搜索
	Ast         *domain.AstStruct // pbname结构
}

func parseCfg(pbname, rule string) interface{} {
	arrs, comment, _ := base.RuleSplit(rule)
	sheet := strings.TrimSuffix(pbname, "ConfigAry")
	tmp := &Attribute{
		ModuleName:  comment[arrs[0]+1 : arrs[1]],
		PackageName: base.ToUnderline(sheet),
		SheetName:   sheet,
		IsCheck:     true,
		Priority:    "0",
		PBName:      comment[arrs[1]+1:arrs[2]] + "S",
	}

	sumfilter := map[string]struct{}{}
	for i := 2; i < len(arrs)-1; i++ {
		val := base.ParseIndex(comment[arrs[i]+1 : arrs[i+1]])
		switch val.Field {
		case "check":
			tmp.IsCheck = !(strings.ToLower(val.Values[0].Name) == "false")
		case "priority":
			tmp.Priority = val.Values[0].Name
		case "struct":
			tmp.IsStruct = true
		case "list":
			tmp.IsList = true
		case "map":
			tmp.Map = append(tmp.Map, val)
		case "group":
			tmp.Group = append(tmp.Group, val)
		case "max":
			tmp.Max = append(tmp.Max, val)
		case "min":
			tmp.Min = append(tmp.Min, val)
		case "sum":
			sumfilter[val.Field] = struct{}{}
			tmp.Sum = append(tmp.Sum, val)
		}
		val.Field = base.FirstToBig(val.Field)
	}
	return tmp
}
