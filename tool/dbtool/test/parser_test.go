package test

import (
	"bytes"
	"testing"
	"universal/library/util"
	"universal/tool/dbtool/domain"
	"universal/tool/dbtool/internal/parse"
	"universal/tool/dbtool/service"
)

func Test_Parser(t *testing.T) {
	domain.PbPath = "../../../common/pb/"
	domain.RedisPath = "../../../common/redis/repository/"

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
