package typespec

import (
	"fmt"
	"strings"
	"universal/tools/gomaker/internal/base"

	"github.com/xuri/excelize/v2"
)

type Params []*Field

type Sheet struct {
	Rule     string   // 规则
	Sheet    string   // 表明
	Config   string   // 表明
	Class    string   // 分类
	Struct   *Struct  // 结构
	Group    []Params // 类型数据
	Map      []Params // 类型数据
	IsList   bool     // 是否为list数据
	IsStruct bool     // 是否为单个数据
	fp       *excelize.File
}

func NewSheet(r string, fp *excelize.File) *Sheet {
	return &Sheet{Rule: r, fp: fp}
}

func SHEET(r, class, sheet, cfg string, fp *excelize.File) *Sheet {
	return &Sheet{Rule: r, Class: class, Sheet: sheet, Config: cfg, fp: fp}
}

func (d *Sheet) GetPkg() string {
	return base.ToUnderline(d.Config)
}

func (d *Sheet) GetRows() ([][]string, error) {
	return d.fp.GetRows(d.Sheet)
}

func (d *Sheet) GetGIndexs() (rets []string) {
	for i := range d.Group {
		rets = append(rets, fmt.Sprintf("group%d", i))
	}
	return
}

func (d *Sheet) GetMIndexs() (rets []string) {
	for i := range d.Map {
		rets = append(rets, fmt.Sprintf("map%d", i))
	}
	return
}

func (d Params) GetArg(pkg string) string {
	rets := []string{}
	for _, item := range d {
		rets = append(rets, fmt.Sprintf("%s %s", item.Name, item.Type.GetType(pkg)))
	}
	return strings.Join(rets, ",")
}

func (d Params) GetParam(str string) string {
	rets := []string{}
	for _, item := range d {
		if len(str) > 0 {
			rets = append(rets, fmt.Sprintf("%s.%s", str, item.Name))
		} else {
			rets = append(rets, item.Name)
		}
	}
	return strings.Join(rets, ",")
}

func (d Params) GetName() string {
	rets := []string{}
	for _, item := range d {
		rets = append(rets, item.Name)
	}
	return strings.Join(rets, "")
}
