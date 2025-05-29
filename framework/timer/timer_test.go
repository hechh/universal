package timer

import (
	"fmt"
	"testing"
	"time"
)

func Print() {
	fmt.Println("hello world")
}

func TestTimer(t *testing.T) {
	timer := NewTimer(4, 7, 4)

	taskId := uint64(123)
	if err := timer.Register(&taskId, Print, 17*time.Millisecond, 2); err != nil {
		t.Fatalf("Register failed: %v", err)
		return
	}
	time.Sleep(100 * time.Millisecond)
}

// 1748483395234   109,280,212,202
// 				   109,280,212,201

// 197 19BB 4EAF
// 197 19BB 4E9F
