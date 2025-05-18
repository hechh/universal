package xlsx

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"forevernine.com/planet/server/tool/gomaker/domain"
	"forevernine.com/planet/server/tool/gomaker/internal/base"
	"forevernine.com/planet/server/tool/gomaker/internal/manager"
	"github.com/xuri/excelize/v2"
)

const (
	Action       = "xlsx"
	RuleTypeXlsx = "//@gomaker:xlsx"
)

func Init() {
	manager.RegisterAction(Action, RuleTypeXlsx)
	manager.RegisterParser(RuleTypeXlsx, parse)
	manager.RegisterCreator(RuleTypeXlsx, gen)
}

//@gomaker:xlsx|sheet@pbname,...
type Attribute struct {
	XlsxName string
	Sheets   base.Params
}

func parse(pbname, comment string) interface{} {
	return &Attribute{
		XlsxName: pbname,
		Sheets:   base.ParseParams(comment[strings.Index(comment, "|")+1:]),
	}
}

func gen(rule, path string, buf *bytes.Buffer) {
	for xlsxName, v := range manager.GetRules(rule) {
		vals, ok := v.(*Attribute)
		if !ok || vals == nil {
			panic(fmt.Sprintf("data is empty or type is error, data: %v", v))
		}
		handle(path, xlsxName, vals.Sheets)
	}
}

func handle(path, xlsxName string, params base.Params) {
	f := excelize.NewFile()
	for _, item := range params {
		var st *domain.AstStruct
		if st = manager.GetAstStruct(item.Type); st == nil {
			panic(fmt.Sprintf("%s is not found in proto file", item.Type))
		}

		f.NewSheet(item.Name)
		rows := parseXlsxRow(st)
		f.SetSheetRow(item.Name, "A2", &rows)
	}
	f.DeleteSheet("sheet1")
	filename := filepath.Join(path, fmt.Sprintf("output/xlsx/%s.xlsx", xlsxName))
	if err := os.MkdirAll(filepath.Dir(filename), os.FileMode(0777)); err != nil {
		panic(err)
	}
	if err := f.SaveAs(filename); err != nil {
		panic(err)
	}
}

func parseXlsxRow(elem *domain.AstStruct) (result []interface{}) {
	for _, item := range elem.Idents {
		// 防止循环解析相同类型
		if elem.Type.Name == item.Type.Name {
			continue
		}

		if item.Type.Token&domain.STRUCT == 0 {
			result = append(result, item.Name)
		} else {
			for _, val := range parseXlsxRow(manager.GetAstStruct(item.Type.Name)) {
				result = append(result, fmt.Sprintf("%s.%s", item.Name, val.(string)))
			}
		}
	}
	for _, item := range elem.Arrays {
		// 防止循环解析相同类型
		if elem.Type.Name == item.Type.Name {
			continue
		}

		if item.Type.Token&domain.STRUCT == 0 {
			result = append(result, fmt.Sprintf("%s#1", item.Name))
		} else {
			for _, val := range parseXlsxRow(manager.GetAstStruct(item.Type.Name)) {
				result = append(result, fmt.Sprintf("%s#1.%s", item.Name, val.(string)))
			}
		}
	}
	return
}
