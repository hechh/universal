package config

import (
	"fmt"
	"io/ioutil"
	"path"
	"universal/common/pb"
	"universal/common/yaml"
	"universal/library/uerror"
)

var (
	configureDir string
	fileMgr      = make(map[string]func(string) error)
)

func Register(sheet string, f func(string) error) {
	if _, ok := fileMgr[sheet]; ok {
		panic(fmt.Sprintf("%s已经注册过了", sheet))
	}
	fileMgr[sheet] = f
}

func Init(cfg *yaml.CommonConfig) error {
	configureDir = cfg.ConfigPath
	for sheet, f := range fileMgr {
		fileName := sheet + ".conf"
		// 加载整个文件
		buf, err := ioutil.ReadFile(path.Join(configureDir, fileName))
		if err != nil {
			return uerror.New(1, pb.ErrorCode_OPEN_FILE_FAILED, fileName)
		}
		if err := f(string(buf)); err != nil {
			return uerror.New(1, pb.ErrorCode_PARSE_FAILED, "加载%s配置错误： %v", fileName, err)
		}
	}
	return nil
}
