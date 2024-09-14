package test

import (
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestXlsx(t *testing.T) {
	fb, err := excelize.OpenFile("define.xlsx")
	if err != nil {
		t.Log(err)
		return
	}

	vals, err := fb.GetRows("testset")
	vv, ok := err.(excelize.ErrSheetNotExist)
	t.Log(vv, ok, len(vals))
}
