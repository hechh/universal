package test

import (
	"testing"
	"time"
)

func TestStock(t *testing.T) {
	layout := "20060102 15:04:05"
	tt, err := time.Parse(layout, "20241111 9:12:56")
	if err != nil {
		panic(err)
	}
	t.Log(tt.Unix())
}
