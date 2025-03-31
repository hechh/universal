package base

import (
	"universal/framework/uerror"

	"github.com/xuri/excelize/v2"
)

type Table struct {
	TypeOf    uint32
	SheetName string
	FileName  string
	fp        *excelize.File
	rows      *excelize.Rows
}

func (t *Table) SetFp(fp *excelize.File) {
	t.fp = fp
}

func (t *Table) ScanRows(count int) (rets [][]string, err error) {
	if t.rows == nil {
		t.rows, err = t.fp.Rows(t.SheetName)
		if err != nil {
			return
		}
	}
	defer t.rows.Close()
	for i := 0; i < count && t.rows.Next(); i++ {
		row, err := t.rows.Columns()
		if err != nil {
			return nil, uerror.NewUError(1, -1, "获取行失败: %v", err)
		}
		rets = append(rets, row)
	}
	return
}

func (t *Table) GetRows() ([][]string, error) {
	return t.fp.GetRows(t.SheetName)
}

func (t *Table) GetCols() ([][]string, error) {
	return t.fp.GetCols(t.SheetName)
}
