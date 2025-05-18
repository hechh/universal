package config

import (
	"fmt"
	"io/ioutil"
	"path"
	"poker_server/common/yaml"
	"poker_server/framework/library/uerror"
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

func InitLocal(cfg *yaml.ConfigureConfig) error {
	configureDir = cfg.LocalPath
	for sheet, f := range fileMgr {
		fileName := sheet + ".json"
		// 加载整个文件
		buf, err := ioutil.ReadFile(path.Join(configureDir, fileName))
		if err != nil {
			return uerror.New(1, -1, err.Error())
		}
		if err := f(string(buf)); err != nil {
			return uerror.New(1, -1, "加载%s配置错误： %v", fileName, err)
		}
	}
	return nil
}

/*
// 初始化配置中心
func InitNet(client config_client.IConfigClient, group string) error {
	for sheet, f := range fileMgr {
		fileName := sheet + ".conf"
		content, err := client.GetConfig(vo.ConfigParam{DataId: fileName, Group: group})
		if err != nil || content == "" {
			return uerror.New(1, -1, "nacos.GetConfig(%s): %v", fileName, err)
		}

		err = client.ListenConfig(vo.ConfigParam{
			DataId: fileName,
			Group:  group,
			OnChange: func(namespace, group, dataId, data string) {
				logger.Infof("gameconf changed !! ** update file: [group: %v , dataId: %v ] **", group, dataId)
				if err := f(data); err != nil {
					logger.Errorf("gameconf changed !! ** update file: [group: %v , dataId: %v ] **", group, dataId)
				}
			},
		})
		if err != nil {
			return uerror.New(1, -1, "nacos.ListenConfig(%s): %v", fileName, err)
		}

		if err := f(content); err != nil {
			return uerror.New(1, -1, "加载配置错误： %v", err)
		}
	}
	return nil
}
*/
