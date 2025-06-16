package main

import (
	"context"
	"flag"
	"hash/crc32"
	"os"
	"path"
	"strings"
	"time"
	"universal/common/yaml"
	"universal/library/util"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	var config, dataPath string
	flag.StringVar(&dataPath, "data", "./data", "配置文件目录")
	flag.StringVar(&config, "config", "config.yaml", "配置文件目录")
	flag.Parse()
	cfg, err := yaml.NewConfig(config)
	if err != nil {
		panic(err)
	}
	client, err := clientv3.New(clientv3.Config{
		Endpoints:            cfg.Etcd.Endpoints,
		DialTimeout:          5 * time.Second,
		DialKeepAliveTime:    30 * time.Second,
		DialKeepAliveTimeout: 3 * time.Second,
		MaxCallSendMsgSize:   1024 * 1024 * 1024,
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()
	rsp, err := client.Get(context.Background(), cfg.Common.ConfigTopic, clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	tmps := map[string]uint32{}
	for _, kv := range rsp.Kvs {
		tmps[string(kv.Key)] = crc32.ChecksumIEEE(kv.Value)
	}
	files, err := util.Glob(dataPath, ".*\\.conf", true)
	if err != nil {
		panic(err)
	}
	for _, filename := range files {
		buf, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		sheet := path.Base(strings.TrimSuffix(filename, ".conf"))
		val := crc32.ChecksumIEEE(buf)
		old, ok := tmps[sheet]
		if ok && old == val {
			continue
		}
		if _, err := client.Put(context.Background(), path.Join(cfg.Common.ConfigTopic, sheet), string(buf)); err != nil {
			panic(err)
		}
	}
}
