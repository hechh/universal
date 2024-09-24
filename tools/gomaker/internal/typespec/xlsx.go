package typespec

import (
	"github.com/xuri/excelize/v2"
)

type Sheet struct {
	Rule     string     // 规则
	Sheet    string     // 表明
	Config   string     // 表明
	Class    string     // 分类
	Group    [][]*Field // 类型数据
	Map      [][]*Field // 类型数据
	IsList   bool       // 是否为list数据
	IsStruct bool       // 是否为单个数据
	fp       *excelize.File
}

func (d *Sheet) GetRows() ([][]string, error) {
	return d.fp.GetRows(d.Sheet)
}

func NewSheet(r string, fp *excelize.File) *Sheet {
	return &Sheet{Rule: r, fp: fp}
}

func SHEET(r, class, sheet, cfg string, fp *excelize.File) *Sheet {
	return &Sheet{Rule: r, Class: class, Sheet: sheet, Config: cfg, fp: fp}
}
