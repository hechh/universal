package manager

import (
	"bytes"
	"strings"

	"forevernine.com/planet/server/tool/gomaker/internal/base"
)

var (
	_parsers = make(map[string]parseFunc) // rule --- parse
	_gens    = make(map[string][]genFunc) // rule --- gen
	_actions = make(map[string][]string)  // module --- rules
)

type parseFunc func(pbname string, comment string) interface{}
type genFunc func(rule, path string, buf *bytes.Buffer)

func RegisterAction(module string, rules ...string) {
	if _, ok := _actions[module]; !ok {
		_actions[module] = rules
	} else {
		_actions[module] = append(_actions[module], rules...)
	}
}

func RegisterParser(rule string, p parseFunc) {
	_parsers[rule] = p
}

func RegisterCreator(rule string, p ...genFunc) {
	_gens[rule] = append(_gens[rule], p...)
}

func ParseRule(pbname string, doc string) {
	if pos := strings.Index(doc, "|"); pos > 0 {
		rule := doc[:pos]
		if _, ok := _parsers[rule]; !ok {
			return
		}
		if _, ok := _rules[rule]; !ok {
			_rules[rule] = make(map[string]interface{})
		}
		_rules[rule][pbname] = _parsers[rule](pbname, doc[pos+1:])
	}
}

func AddConfigAry(pbname string) {
	if strings.HasSuffix(pbname, "ConfigAry") {
		_pbconfigs[pbname] = struct{}{}
	}
}

func GenCode(rule, path string, buf *bytes.Buffer) {
	for _, gen := range _gens[rule] {
		buf.Reset()
		gen(rule, path, buf)
	}
}

func GetRuleTypes() (rets []string) {
	for rule := range _gens {
		rets = append(rets, rule)
	}
	return
}

func GetRuleType(action string) (rets []string) {
	base.StringSplit(action, ',', func(m string) {
		if vals, ok := _actions[m]; ok {
			rets = append(rets, vals...)
		}
	})
	return rets
}
