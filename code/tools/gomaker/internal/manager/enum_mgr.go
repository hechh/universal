package manager

import (
	"fmt"
	"sort"
	"strings"
	"universal/tools/gomaker/internal/typespec"
)

var (
	enums = make(map[string]*typespec.Enum)
)

// -------枚举类型---------
func AddValue(vv *typespec.Value) {
	if vv == nil {
		return
	}

	name := fmt.Sprintf("%s.%s", vv.Type.Selector, vv.Type.Name)
	if eval, ok := enums[name]; !ok {
		enums[name] = typespec.NewEnum(vv.Type).AddValue(vv)
	} else {
		eval.AddValue(vv)
	}
}

func AddEnum(vv *typespec.Enum) {
	if vv == nil {
		return
	}

	name := fmt.Sprintf("%s.%s", vv.Type.Selector, vv.Type.Name)
	if eval, ok := enums[name]; !ok {
		enums[name] = vv
	} else {
		for _, item := range vv.Values {
			eval.AddValue(item)
		}
	}
}

func GetEnumList() (rets []*typespec.Enum) {
	for _, val := range enums {
		rets = append(rets, val)
	}
	sort.Slice(rets, func(i, j int) bool {
		return strings.Compare(rets[i].Type.Name, rets[j].Type.Name) < 0
	})
	return
}
