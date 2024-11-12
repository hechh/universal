package parse

import (
	"fmt"
	"stock/internal/base"
	"stock/internal/typespec"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

func ParseFilter(filename string, f func(string, ...string)) error {
	// 打开文件
	fp, err := excelize.OpenFile(filename)
	if err != nil {
		return base.NewUError(1, -1, "%s: %v", filename, err)
	}
	defer fp.Close()
	// 解析转换表
	if rows, err := fp.GetRows("转换表"); err == nil {
		for _, vals := range rows {
			if len(vals[0]) > 0 && len(vals[1]) > 0 {
				//fmt.Println(pos, "=====>", val, val[:pos], val[pos+1:])
				f(vals[0], strings.Split(vals[1], ",")...)
			}
		}
	}
	// 解析过滤表
	if rows, err := fp.GetRows("过滤表"); err == nil {
		for _, vals := range rows {
			for _, val := range vals {
				for _, item := range strings.Split(val, ",") {
					f(item)
				}
			}
		}
	}
	return nil
}

func ParseStock(filename string, f func(*typespec.Stock)) error {
	// 打开文件
	fp, err := excelize.OpenFile(filename)
	if err != nil {
		return base.NewUError(1, -1, "%s: %v", filename, err)
	}
	defer fp.Close()
	// 读取所有数据
	for _, sheet := range fp.GetSheetList() {
		rows, err := fp.GetRows(sheet)
		if len(rows) <= 0 {
			continue
		}
		if err != nil {
			return base.NewUError(1, -1, "%s(%s): %v", filename, sheet, err)
		}
		// 解析下标
		tmps := map[string]int{}
		for i, val := range rows[1] {
			if len(val) > 0 {
				tmps[val] = i
			}
		}
		// 解析数据
		for _, vals := range rows[2:] {
			f(&typespec.Stock{
				Code:      fmt.Sprintf("%06s", vals[tmps["Code"]]),
				Name:      vals[tmps["Name"]],
				Reason:    vals[tmps["Reason"]],
				Themes:    strings.Split(vals[tmps["Reason"]], "+"),
				BeginTime: vals[tmps["BeginTime"]],
				EndTime:   vals[tmps["EndTime"]],
				BeginUnix: timeParse(sheet, vals[tmps["BeginTime"]]),
				EndUnix:   timeParse(sheet, vals[tmps["EndTime"]]),
				Continue:  cast.ToInt(vals[tmps["Continue"]]),
				Total:     vals[tmps["Total"]],
				Time: &typespec.TimeInfo{
					Date:  cast.ToInt64(sheet),
					Year:  cast.ToInt32(sheet[:4]),
					Month: cast.ToInt32(sheet[4:6]),
					Day:   cast.ToInt32(sheet[6:]),
				},
			})
		}
	}
	return nil
}

func timeParse(date, tt string) int64 {
	ti, _ := time.Parse("20060102 15:04:05", fmt.Sprintf("%s %s", date, tt))
	return ti.Unix()
}
