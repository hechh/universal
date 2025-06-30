package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"universal/common/yaml"
	"universal/library/uerror"
)

var (
	fileMgr = make(map[string]func([]byte) error)
	watcher *Watcher
)

func Register(sheet string, f func([]byte) error) {
	if _, ok := fileMgr[sheet]; ok {
		panic(fmt.Sprintf("%s已经注册过了", sheet))
	}
	fileMgr[sheet] = f
}

func Init(cfg *yaml.TableConfig) (err error) {
	tmps := make(map[string]struct{})
	if err := InitConfig(cfg.Path, tmps); err != nil {
		return err
	}
	if cfg.IsRemote {
		if watcher, err = NewWatcher(cfg); err != nil {
			return
		}
		if err := watcher.Load(tmps); err != nil {
			return err
		}
		watcher.Watch(tmps)
	}
	return
}

func InitConfig(cpath string, tmps map[string]struct{}) error {
	for sheet, f := range fileMgr {
		filename := path.Join(cpath, sheet+".conf")
		if _, err := os.Stat(filename); err != nil {
			if os.IsNotExist(err) {
				continue
			}
		}
		buf, err := ioutil.ReadFile(filename)
		if err != nil {
			return uerror.N(1, -1, sheet)
		}
		if err := f(buf); err != nil {
			return uerror.N(1, -1, "加载%s配置错误： %v", sheet, err)
		}
		if tmps != nil {
			tmps[sheet] = struct{}{}
		}
	}
	return nil
}

func Close() error {
	if watcher != nil {
		return watcher.Close()
	}
	return nil
}
