package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"universal/common/yaml"
	"universal/library/mlog"
	"universal/library/uerror"
)

var (
	fileMgr = make(map[string]func(string) error)
	watcher *Watcher
)

func Register(sheet string, f func(string) error) {
	if _, ok := fileMgr[sheet]; ok {
		panic(fmt.Sprintf("%s已经注册过了", sheet))
	}
	fileMgr[sheet] = f
}

func Init(cfg *yaml.EtcdConfig, ccfg *yaml.DataConfig) (err error) {
	var w *Watcher
	if w, err = NewWatcher(cfg, ccfg); err != nil {
		return err
	}
	tmps := make(map[string]struct{})
	if !ccfg.IsRemote {
		mlog.Infof("加载本地磁盘中的配置")
		if err := LoadAndUpload(ccfg.Path, w, tmps); err != nil {
			return err
		}
	} else {
		mlog.Infof("加载etcd中的配置")
		if err := w.Download(tmps); err != nil {
			return err
		}
	}
	mlog.Infof("监听etcd中的配置")
	return w.Watch(tmps)
}

func LoadAndUpload(cpath string, w *Watcher, tmps map[string]struct{}) error {
	for sheet, f := range fileMgr {
		filename := path.Join(cpath, sheet+".conf")
		if _, err := os.Stat(filename); err != nil {
			if os.IsNotExist(err) {
				return err
			}
		}
		buf, err := ioutil.ReadFile(filename)
		if err != nil {
			return uerror.New(1, -1, "加载%s配置错误%v", sheet, err)
		}
		if err := f(string(buf)); err != nil {
			return uerror.New(1, -1, "加载%s配置错误： %v", sheet, err)
		}
		if err := w.Upload(sheet, buf); err != nil {
			return err
		}
		tmps[sheet] = struct{}{}
	}
	return nil
}

func Close() error {
	if watcher != nil {
		return watcher.Close()
	}
	return nil
}
