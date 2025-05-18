package module

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"forevernine.com/planet/server/tool/gomaker/domain"
	"forevernine.com/planet/server/tool/gomaker/internal/manager"

	"forevernine.com/planet/server/tool/gomaker/internal/base"
)

const (
	Action         = "module"
	RuleTypeModule = "//@gomaker:module"
	apiGo          = "%s/%s/%s_api.gen.go"
	eventGo        = "%s/%s/%s_event.gen.go"
	initGo         = "%s/%s/init.gen.go"
	actionGo       = "%s/%s/action.gen.go"
	triggerGo      = "%s/%s/trigger.gen.go"
	actProto       = "deps/proto/activity.proto"
	cmdProto       = "deps/proto/cmd.proto"
)

// @gomaker:module|activityType|req1,req2
func Init() {
	manager.RegisterAction(Action, RuleTypeModule)
	manager.RegisterParser(RuleTypeModule, parse)
	manager.RegisterCreator(RuleTypeModule, activity, cmd, api)
}

type ApiInfo struct {
	PackageName string
	FuncName    string
	Req         string
	Rsp         string
	CmdReq      string
	CmdRsp      string
}

type EventInfo struct {
	PackageName string
	FuncName    string
	Event       string
	CmdEvent    string
}

type Attribute struct {
	PackageName  string
	ActivityType string
	ApiList      []*ApiInfo
	EventList    []*EventInfo
}

func parse(pbname string, comment string) interface{} {
	rules := &Attribute{
		PackageName: strings.ToLower(pbname),
		ActivityType: func() string {
			if strings.Contains(comment, "ActivityType") {
				return comment[:strings.Index(comment, "|")]
			}
			return ""
		}(),
	}
	base.StringSplit(comment[strings.LastIndex(comment, "|")+1:], ',', func(req string) {
		if strings.HasSuffix(req, "Req") {
			cmdReq := base.ToCmd(req)
			rules.ApiList = append(rules.ApiList, &ApiInfo{
				PackageName: rules.PackageName,
				FuncName:    strings.TrimPrefix(strings.TrimSuffix(req, "Req"), "Game"),
				Req:         req,
				Rsp:         strings.Replace(req, "Req", "Rsp", 1),
				CmdReq:      cmdReq,
				CmdRsp:      strings.Replace(cmdReq, "_REQ", "_RSP", 1),
			})
		} else {
			cmdEvent := base.ToEvent(req)
			rules.EventList = append(rules.EventList, &EventInfo{
				PackageName: rules.PackageName,
				FuncName:    strings.TrimPrefix(strings.TrimSuffix(req, "Event"), "Game"),
				Event:       req,
				CmdEvent:    cmdEvent,
			})
		}
	})
	return rules
}

func activity(rule, path string, buf *bytes.Buffer) {
	path = base.ROOT
	isGen := false
	max := int32(0)
	rules := []*domain.AstValue{}
	filter := manager.GetAstEnum("ActivityType")

	manager.WalkAstEnum("ActivityType", func(item *domain.AstValue) {
		rules = append(rules, item)
		if max < item.Value {
			max = item.Value
		}
	})
	for _, item := range manager.GetRules(rule) {
		val, ok := item.(*Attribute)
		if !ok || val == nil {
			panic(fmt.Sprintf("data is empty or type is error, data: %v", val))
		}
		if len(val.ActivityType) <= 0 {
			continue
		}
		if _, ok := filter.Values[val.ActivityType]; ok {
			log.Printf("%s is already exist", val.ActivityType)
			continue
		}

		max++
		rules = append(rules, &domain.AstValue{Name: val.ActivityType, Value: max})
		isGen = true
	}
	if isGen {
		sort.Slice(rules, func(i, j int) bool {
			return rules[i].Value < rules[j].Value
		})
		manager.Execute(Action, "activity.tpl", buf, rules)
		base.GenProto(buf, filepath.Join(path, actProto), true)
	}
}

func cmd(rule, path string, buf *bytes.Buffer) {
	path = base.ROOT
	nums := map[int32]struct{}{}
	rules := []*domain.AstValue{}
	manager.WalkAstEnum("CMD", func(item *domain.AstValue) {
		rules = append(rules, item)
		nums[item.Value] = struct{}{}
	})

	isGen := false
	j := int32(9999)
	filter := manager.GetAstEnum("CMD")
	for _, item := range manager.GetRules(rule) {
		val, ok := item.(*Attribute)
		if !ok || val == nil {
			panic(fmt.Sprintf("data is empty or type is error, data: %v", val))
		}
		for _, vval := range val.ApiList {
			if _, ok := filter.Values[vval.CmdReq]; ok {
				log.Printf("%s is registered", vval.CmdReq)
				continue
			}
			for ok01, ok02 := true, true; ok01 || ok02; j += 2 {
				_, ok01 = nums[j]
				_, ok02 = nums[j+1]
			}
			rules = append(rules, &domain.AstValue{Name: vval.CmdReq, Value: j - 2}, &domain.AstValue{Name: vval.CmdRsp, Value: j - 1})
			log.Printf("%s: %d, %s: %d", vval.CmdReq, j, vval.CmdRsp, j+1)
			isGen = true
		}
		for _, vval := range val.EventList {
			if _, ok := filter.Values[vval.CmdEvent]; ok {
				log.Printf("%s is registered", vval.CmdEvent)
				continue
			}
			for ok01, ok02 := true, true; ok01 || ok02; j += 2 {
				_, ok01 = nums[j]
				_, ok02 = nums[j+1]
			}
			rules = append(rules, &domain.AstValue{Name: vval.CmdEvent, Value: j - 2})
			log.Printf("%s: %d", vval.CmdEvent, j-2)
			isGen = true
		}
	}
	if isGen {
		sort.Slice(rules, func(i, j int) bool {
			return rules[i].Value < rules[j].Value
		})
		manager.Execute(Action, "cmd.tpl", buf, rules)
		base.GenProto(buf, filepath.Join(path, cmdProto), true)
	}
}

func api(rule, path string, buf *bytes.Buffer) {
	if len(path) <= 0 {
		path = base.ROOT
	}

	for _, data := range manager.GetRules(rule) {
		vals, ok := data.(*Attribute)
		if !ok || vals == nil {
			panic(fmt.Sprintf("data is empty or type is error, data: %v", vals))
		}
		manager.Execute(Action, "init.tpl", buf, vals)
		base.GenGo(buf, fmt.Sprintf(initGo, path, vals.PackageName), true)
		buf.Reset()

		manager.Execute(Action, "action.tpl", buf, vals)
		base.GenGo(buf, fmt.Sprintf(actionGo, path, vals.PackageName), false)
		buf.Reset()

		manager.Execute(Action, "trigger.tpl", buf, vals)
		base.GenGo(buf, fmt.Sprintf(triggerGo, path, vals.PackageName), false)

		for _, api := range vals.ApiList {
			buf.Reset()
			manager.Execute(Action, "api.tpl", buf, api)
			base.GenGo(buf, fmt.Sprintf(apiGo, path, vals.PackageName, api.FuncName), false)
		}
		for _, api := range vals.EventList {
			buf.Reset()
			manager.Execute(Action, "event.tpl", buf, api)
			base.GenGo(buf, fmt.Sprintf(eventGo, path, vals.PackageName, api.FuncName), false)
		}
	}
}
