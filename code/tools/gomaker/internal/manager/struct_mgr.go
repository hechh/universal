package manager

import (
	"fmt"
	"sort"
	"strings"
	"universal/tools/gomaker/internal/typespec"
)

var (
	structs = make(map[string]*typespec.Struct)
)

func AddStruct(vv *typespec.Struct) {
	if vv != nil {
		structs[fmt.Sprintf("%s.%s", vv.Type.Selector, vv.Type.Name)] = vv
	}
}

func GetStructList() (rets []*typespec.Struct) {
	for _, val := range structs {
		rets = append(rets, val)
	}
	sort.Slice(rets, func(i, j int) bool {
		return strings.Compare(rets[i].Type.Name, rets[j].Type.Name) < 0
	})
	return
}
