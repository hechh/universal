package test

import (
	"encoding/json"

	"google.golang.org/protobuf/proto"
)

type TestReportClient struct{}

func (d *TestReportClient) Write(buf []byte) (int, error) {
	//base.Debug(0, string(buf))
	return 0, nil
}

func (d *TestReportClient) Close() {
	//base.Debug(0, "close TestReportClient")
}

func ToString(db interface{}) string {
	buf, _ := json.Marshal(db)
	return string(buf)
}

func Copy(dst, src proto.Message) proto.Message {
	buf, _ := proto.Marshal(src)
	proto.Unmarshal(buf, dst)
	return dst
}
