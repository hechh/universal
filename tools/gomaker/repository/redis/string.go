package redis

import (
	"bytes"
	"fmt"

	"forevernine.com/planet/server/tool/gomaker/internal/base"
	"forevernine.com/planet/server/tool/gomaker/internal/manager"
)

// @gomaker:redis:string|数据库|key变量名:组成key的参数1,参数2,...|#注释
func parseString(pbname, comment string) interface{} {
	arrs, comment, desc := base.RuleSplit(comment)
	return &RedisAttr{
		Type:    StringRedis,
		Package: base.ToUnderline(pbname),
		DbName:  comment[arrs[0]+1 : arrs[1]],
		Name:    pbname,
		Desc:    desc,
		Key:     base.ParseIndex(comment[arrs[1]+1 : arrs[2]]),
		Ast:     manager.GetAstStruct(pbname),
	}
}

func genString(rule, path string, buf *bytes.Buffer) {
	if len(path) <= 0 {
		path = base.ROOT
	}
	for _, v := range manager.GetRules(rule) {
		val, ok := v.(*RedisAttr)
		if !ok || val == nil {
			panic(fmt.Sprintf("data is empty or type is error, data: %v", v))
		}
		cmds, ok := manager.GetRule(RuleTypeCMD, val.Name).([]string)
		if ok && len(cmds) > 0 {
			for _, cmd := range cmds {
				if cmd == "CACHE" {
					val.IsCache = true
					break
				}
			}
		}
		manager.Execute(Action, "stringRedis.tpl", buf, val)
		for _, cmd := range cmds {
			manager.Execute(Action, cmd, buf, val)
		}
		base.GenGo(buf, fmt.Sprintf(goFile, path, val.Package), true)
		buf.Reset()
	}
}
