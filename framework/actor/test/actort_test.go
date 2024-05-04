package test

import (
	"fmt"
	"sync"
	"testing"
	"universal/common/pb"
	"universal/framework/actor/internal/manager"
	"universal/framework/fbasic"

	"github.com/spf13/cast"
)

var (
	wg = sync.WaitGroup{}
)

func TestMain(m *testing.M) {
	manager.SetActorHandle(handle)
	m.Run()
}

func handle(ctx *fbasic.Context, buf []byte) func() {
	return func() {
		fmt.Println(ctx.PacketHead, "----->", string(buf), ctx.GetValue("test"))
		wg.Done()
	}
}

func Test_Send(t *testing.T) {
	for i := 0; i <= 100; i++ {
		wg.Add(1)
		Head := &pb.PacketHead{
			SendType:       pb.SendType_POINT,
			ApiCode:        2,
			UID:            100100600 + uint64(i),
			SrcClusterType: pb.ClusterType_GATE,
			DstClusterType: pb.ClusterType_GAME,
		}
		key := cast.ToString(Head.UID)
		act := manager.GetIActor(key)
		act.SetObject("test", "this is a test")
		manager.Send(key, &pb.Packet{Head: Head})
	}
	wg.Wait()
}
