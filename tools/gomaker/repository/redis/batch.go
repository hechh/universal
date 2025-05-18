package redis

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"forevernine.com/planet/server/tool/gomaker/internal/base"
	"forevernine.com/planet/server/tool/gomaker/internal/manager"
)

type PBNameList struct {
	PBName string
	Field  *base.Index
}

type ActivityBatchRule struct {
	UniqueID string      // 唯一ID // 数据库名称+key类型
	DBName   string      // redis数据库名
	Key      *base.Index // redis的key
	List     []*PBNameList
}

func genBatch(rule, path string, buf *bytes.Buffer) {
	if len(path) <= 0 {
		path = base.ROOT
	}

	// 格式处理
	tmp := []*ActivityBatchRule{}
	filter := map[string]*ActivityBatchRule{}
	for _, v := range manager.GetRules(rule) {
		val, ok := v.(*RedisAttr)
		if !ok || val == nil || !val.IsHash() {
			continue
		}
		if val.Key.Count() != 1 || val.Field.Count() > 0 || val.Key.Values[0].UniqueID() != "uid@string" {
			continue
		}

		id := fmt.Sprintf("%s%s", val.DbName, val.Key.Field)
		if item, ok := filter[id]; !ok {
			filter[id] = &ActivityBatchRule{
				UniqueID: fmt.Sprintf("%s%s", val.DbName, val.Key.Field),
				DBName:   val.DbName,
				Key:      val.Key,
				List:     []*PBNameList{{PBName: val.Name, Field: val.Field}},
			}
			tmp = append(tmp, filter[id])
		} else {
			item.List = append(item.List, &PBNameList{PBName: val.Name, Field: val.Field})
		}
	}

	// 排序
	sort.Slice(tmp, func(i, j int) bool {
		return strings.Compare(tmp[i].UniqueID, tmp[j].UniqueID) >= 0
	})
	for _, item := range tmp {
		sort.Slice(item.List, func(i, j int) bool {
			return strings.Compare(item.List[i].PBName, item.List[j].PBName) >= 0
		})
	}
	manager.Execute(Action, "batchRedis.tpl", buf, tmp)
	base.GenGo(buf, fmt.Sprintf(goFile, path, "activity_batch"), true)
	buf.Reset()
}
