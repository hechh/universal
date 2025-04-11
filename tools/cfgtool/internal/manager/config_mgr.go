package manager

import (
	"hego/tools/cfgtool/internal/base"
	"sort"
	"strings"
)

var (
	configMgr = make(map[string]*base.Config)
)

func GetOrNewConfig(file, sheet, name string) *base.Config {
	if val, ok := configMgr[name]; ok {
		return val
	}
	configMgr[name] = &base.Config{
		Name:     name,
		FileName: file,
		Sheet:    sheet,
		Fields:   make(map[string]*base.Field),
		Indexs:   make(map[int][]*base.Index),
	}
	return configMgr[name]
}

func GetConfigList() (rets []*base.Config) {
	for _, val := range configMgr {
		rets = append(rets, val)
	}
	sort.Slice(rets, func(i, j int) bool {
		return strings.Compare(rets[i].Name, rets[j].Name) <= 0
	})
	return
}

func GetConfigMap() map[string]*base.Config {
	return configMgr
}
