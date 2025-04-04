package manager

import (
	"hego/common/config/domain"
	"hego/common/config/internal/parse"
	"hego/framework/basic"
	"hego/framework/plog"
	"runtime/debug"
	"time"
)

var (
	path  string                           // 配置文件路径
	files = make(map[string]*parse.Parser) // 所有配置解析器
)

func Init(dir string, ttl time.Duration) error {
	path = dir
	// 加载所有配置
	for _, par := range files {
		if err := par.Load(dir); err != nil {
			return err
		}
	}
	// 定时检查
	tt := time.NewTicker(ttl)
	basic.SafeGo(func(err interface{}) {
		plog.Fatal("error: %v, stack: %s", err, string(debug.Stack()))
	}, func() {
		for {
			<-tt.C
			check()
		}
	})
	return nil
}

func Register(name string, cfgs ...domain.LoadFunc) {
	val, ok := files[name]
	if !ok {
		files[name] = parse.NewParser(name, cfgs...)
		return
	}
	val.Register(cfgs...)
}

// 定时检测配置变更协程
func check() {
	// 检测变更
	for name, par := range files {
		if !par.Check(path) {
			continue
		}

		// 重新加载配置
		if err := par.Load(path); err != nil {
			plog.Fatal("%s/%s.bytes加载失败", path, name)
		}
	}
}
