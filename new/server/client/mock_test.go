package main

import (
	"poker_server/framework/mock"
	"testing"
	"time"
)

func TestMock(t *testing.T) {
	client, err := mock.NewMockClient("../../env/hch/local.yaml", 1, 1000001)
	if err != nil {
		t.Fatalf("MockClient初始化失败: %v", err)
		return
	}

	time.Sleep(20 * time.Second)
	client.Close()
}
