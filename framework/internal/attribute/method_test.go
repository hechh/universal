package attribute

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"universal/common/pb"

	"github.com/golang/protobuf/proto"
)

type Player struct {
}

func (p *Player) Login(h *pb.Head, req proto.Message, rsp proto.Message) error {
	fmt.Println(req, "---->", rsp)
	return nil
}

func parseName(rr reflect.Type) string {
	name := rr.String()
	if index := strings.Index(name, "."); index > -1 {
		name = name[index+1:]
	}
	return name
}

func TestHandler(t *testing.T) {
	pl := &Player{}
	vv := reflect.TypeOf(pl.Login)
	t.Log("=======>", vv.String())
	t.Log("------>", vv.NumIn(), parseName(vv.In(1)))
}
