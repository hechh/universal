package main

import (
	"bytes"
	"flag"
	"universal/library/util"
	"universal/tool/dbtool/domain"
	"universal/tool/dbtool/internal/parse"
	"universal/tool/dbtool/service"
)

// dbtool工具用于生成dao/repository/redis/目录下的代码
func main() {
	flag.StringVar(&domain.PbPath, "pb", "", "proto文件目录")
	flag.StringVar(&domain.RedisPath, "redis", "", "redis文件目录")
	flag.Parse()

	if len(domain.PbPath) <= 0 {
		panic("proto文件目录不能为空")
	}

	// 加载所有文件
	files, err := util.Glob(domain.PbPath, ".*\\.pb\\.go", true)
	if err != nil {
		panic(err)
	}

	// 解析所有文件
	if err := parse.ParseFiles(&parse.Parser{}, files...); err != nil {
		panic(err)
	}

	// 生成代码
	buf := bytes.NewBuffer(nil)
	service.GenString(buf)
	service.GenHash(buf)
}
