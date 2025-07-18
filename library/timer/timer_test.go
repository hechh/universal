package timer

import (
	"fmt"
	"testing"
	"time"
	"universal/library/mlog"
)

func TestTimer(t *testing.T) {
	timer := NewTimer(4)
	taskId := uint64(123)
	for i := 0; i < 2; i++ {
		if err := timer.Register(&taskId, func() { fmt.Println("-->", i, time.Now().Unix()) }, 1*time.Second, -1); err != nil {
			mlog.Errorf("Register failed: %v", err)
			return
		}
	}
	time.Sleep(4 * time.Second)
	//select {}
}
