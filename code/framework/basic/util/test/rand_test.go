package test

import (
	"math/rand"
	"testing"
	"universal/framework/basic/util"
)

func TestRand(t *testing.T) {
	t.Log(util.RangeInt63n(1, 2))
	t.Log(rand.Int())
}
