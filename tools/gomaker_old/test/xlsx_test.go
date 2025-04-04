package test

import (
	"hego/tools/gomaker/internal/parser"
	"io/ioutil"
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

func TestProto(t *testing.T) {
	filename := "../../../configure/proto/packet.proto"
	buf, _ := ioutil.ReadFile(filename)
	pp := &parser.PbParser{}
	pp.ParseFile(buf)
}
