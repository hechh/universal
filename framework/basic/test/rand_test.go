package test

import (
	"hego/common/pb"
	"hego/framework/basic"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"
)

func TestRand(t *testing.T) {
	t.Log(basic.RangeInt63n(1, 2))
	t.Log(rand.Int())
}

func TestWalk(t *testing.T) {
	root := "../"
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if root != path {
				return filepath.SkipDir
			}
		}
		t.Log(path, info.Name(), info.IsDir())
		return nil
	})
}

func TestMessage(t *testing.T) {
	vv := &pb.Head{}

	t.Log("========>", proto.MessageName(vv))

	tt := proto.MessageType("pb.Head")
	t.Log(tt == nil, tt)
}
