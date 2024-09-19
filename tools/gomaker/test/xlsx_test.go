package test

import (
	"io/ioutil"
	"testing"
	"universal/tools/gomaker/internal/parser"

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
	//t.Log(string(buf))
	pp := parser.NewPbParser(buf)
	pp.Parse()
}
