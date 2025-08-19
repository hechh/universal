package http

import (
	"encoding/json"
	"fmt"
	"poker_server/common/pb"
	"testing"
)

func TestChargeTransDecode(t *testing.T) {
	buf := []byte(`{"resp_msg": {"code": 6005,"message": ""},"resp_data": []}`)
	rsp := &pb.HttpTransferInRsp{}
	json.Unmarshal(buf, rsp)
	fmt.Printf("%v \n", rsp)
}
