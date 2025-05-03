package actor

import (
	"reflect"
	"sync"
	"universal/framework/define"
)

type ActorGroup struct {
	methods map[string]reflect.Method
	mutex   sync.Mutex
	actors  map[int64]define.IActor
}
