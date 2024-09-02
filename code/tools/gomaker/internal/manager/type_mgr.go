package manager

import (
	"fmt"
	"sort"
	"strings"
	"universal/tools/gomaker/internal/typespec"
)

var (
	alias = make(map[string]*typespec.Alias)
	types = make(map[string]*typespec.Type)
)

func GetOrAddType(tt *typespec.Type) *typespec.Type {
	key := fmt.Sprintf("%s.%s", tt.Selector, tt.Name)
	if val, ok := types[key]; !ok {
		types[key] = tt
	} else {
		if len(tt.Doc) > 0 {
			val.Doc = tt.Doc
		}
		if tt.Kind > 0 {
			val.Kind = tt.Kind
		}
	}
	return types[key]
}

// -----------自定义别名---------------
func AddAlias(vv *typespec.Alias) {
	if vv != nil {
		alias[fmt.Sprintf("%s.%s", vv.Type.Selector, vv.Type.Name)] = vv
	}
}

func GetAliasList() (rets []*typespec.Alias) {
	for _, val := range alias {
		rets = append(rets, val)
	}
	sort.Slice(rets, func(i, j int) bool {
		return strings.Compare(rets[i].Type.Name, rets[j].Type.Name) < 0
	})
	return
}
