package manager

import (
	"time"
	"universal/common/config/domain"
	"universal/framework/basic/util"
	"universal/framework/plog"
)

var (
	path  string                     // 配置文件路径
	files = make(map[string]*Parser) // 所有配置解析器
)

func Init(dir string) error {
	path = dir
	// 加载所有配置
	for _, par := range files {
		if err := par.Load(dir); err != nil {
			return err
		}
	}
	// 定时检查
	initCheck()
	return nil
}

func Register(name string, cfgs ...domain.IConfig) {
	val, ok := files[name]
	if !ok {
		files[name] = NewParser(name, cfgs...)
		return
	}
	val.Register(cfgs...)
}

// 定时检测配置变更协程
func initCheck() {
	tt := time.NewTicker(1 * time.Second)
	util.SafeGo(nil, func() {
		for {
			<-tt.C
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
	})
}
