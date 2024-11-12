package manager

import (
	"fmt"
	"sort"
	"stock/internal/parse"
	"stock/internal/typespec"
	"strings"

	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

var (
	dates   = []int64{}                                  // 时间序列
	datas   = make(map[int64]map[string]*typespec.Stock) // 涨停数据
	hots    = make(map[int64][]*typespec.Hotspot)        // 热点数据
	convs   = make(map[string]string)                    // 重定义题材
	filters = make(map[string]struct{})                  // 被过滤题材
)

func Write(fp *excelize.File) {
	// 根据时间遍历
	for _, date := range dates {
		sheet := cast.ToString(date)
		fp.NewSheet(sheet)
		// 写入xlsx表
		rowNum := 1
		for _, hot := range hots[date] {
			//fmt.Println(date, "------>", hot)
			old := rowNum
			for _, con := range hot.List {
				rows := getRow(hot.Name, con, hot.Members[con]...)
				cellName, _ := excelize.CoordinatesToCellName(1, rowNum)
				fp.SetSheetRow(sheet, cellName, &rows)
				rowNum++
			}
			top, _ := excelize.CoordinatesToCellName(1, old)
			button, _ := excelize.CoordinatesToCellName(1, rowNum-1)
			fp.MergeCell(sheet, top, button)
		}
	}
}

func getRow(name string, con int, list ...*typespec.Stock) (ret []interface{}) {
	tmps := []string{}
	for _, st := range list {
		tmps = append(tmps, st.Name)
	}
	ret = append(ret, name, fmt.Sprintf("%d连板", con), strings.Join(tmps, ","))
	return
}

func Analyse() {
	for date, kv := range datas {
		dates = append(dates, date)
		list := []*typespec.Hotspot{}
		tmps := map[string]*typespec.Hotspot{}
		for _, st := range kv {
			// fmt.Println(date, "------>", st)
			// 题材分类
			for _, th := range st.Themes {
				hot, ok := tmps[th]
				if !ok {
					hot = &typespec.Hotspot{
						Name:        th,
						MaxContinue: st.Continue,
						Members:     make(map[int][]*typespec.Stock),
					}
					tmps[th] = hot
					list = append(list, hot)
				}
				// 设置最大连板
				if hot.MaxContinue > st.Continue {
					hot.MaxContinue = st.Continue
				}
				if _, ok := hot.Members[st.Continue]; !ok {
					hot.List = append(hot.List, st.Continue)
				}
				hot.Members[st.Continue] = append(hot.Members[st.Continue], st)
			}
		}
		// 排序
		sort.Slice(list, func(i, j int) bool { return list[i].MaxContinue > list[j].MaxContinue })
		for _, tt := range list {
			//fmt.Println("--------->", tt)
			sort.Slice(tt.List, func(i, j int) bool { return tt.List[i] > tt.List[j] })
		}
		hots[date] = list
	}
	// 时间序列
	sort.Slice(dates, func(i, j int) bool { return dates[i] > dates[j] })
}

func ParseStocks(files ...string) error {
	for _, filename := range files {
		err := parse.ParseStock(filename, func(item *typespec.Stock) {
			if _, ok := datas[item.Time.Date]; !ok {
				datas[item.Time.Date] = make(map[string]*typespec.Stock)
			}
			// 过滤无用概念
			j := -1
			for _, val := range item.Themes {
				if _, ok := filters[val]; ok {
					continue
				}
				j++
				if nval, ok := convs[val]; ok {
					item.Themes[j] = nval
				} else {
					item.Themes[j] = val
				}
			}
			item.Themes = item.Themes[:j+1]
			// 保存数据
			datas[item.Time.Date][item.Code] = item
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func ParseFilters(files ...string) error {
	for _, filename := range files {
		err := parse.ParseFilter(filename, func(key string, vals ...string) {
			if len(vals) > 0 {
				for _, kk := range vals {
					convs[kk] = key
				}
			} else {
				filters[key] = struct{}{}
			}
		})
		if err != nil {
			return err
		}
	}
	return nil
}
