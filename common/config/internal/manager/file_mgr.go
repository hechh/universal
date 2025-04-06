package manager

import (
	"hego/common/config/internal/parse"
	"hego/framework/plog"
	"time"
)

var (
	configPath string                           // 配置文件路径
	fileExt    string                           // 配置文件后缀
	files      = make(map[string]*parse.Parser) // 所有配置解析器
)

func Register(name string, f func([]byte) error) {
	if _, ok := files[name]; ok {
		panic("配置文件名重复: " + name)
	}
	files[name] = parse.NewParser(name, f)
}

func Init(dir, ext string, ttl time.Duration) error {
	configPath = dir
	fileExt = ext
	// 加载所有配置
	for _, par := range files {
		if err := par.Load(dir, ext); err != nil {
			return err
		}
	}

	// 定时检查
	go check(ttl)
	return nil
}

func check(ttl time.Duration) {
	tt := time.NewTicker(ttl)
	defer tt.Stop()

	for {
		<-tt.C
		for name, par := range files {
			if !par.IsChange(configPath, fileExt) {
				continue
			}

			// 重新加载配置
			if err := par.Load(configPath, fileExt); err != nil {
				plog.Fatal("%s/%s.%s加载失败", configPath, name, fileExt)
			}
		}
	}
}
