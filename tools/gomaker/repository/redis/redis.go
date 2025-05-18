package redis

import (
	"strings"

	"forevernine.com/planet/server/tool/gomaker/domain"
	"forevernine.com/planet/server/tool/gomaker/internal/base"
	"forevernine.com/planet/server/tool/gomaker/internal/manager"
)

const (
	StringRedis    = 1
	HashRedis      = 2
	ZSetRedis      = 3
	Action         = "redis"
	RuleTypeCMD    = "//@gomaker:redis:cmd"
	RuleTypeString = "//@gomaker:redis:string"
	RuleTypeHash   = "//@gomaker:redis:hash"
	RuleTypeTool   = "dao_tool"
	goFile         = "%s/common/dao/repository/redis/%s/redis_api.gen.go"
	toolFile       = "%s/common/dao/repository/tools/init_tool.gen.go"
)

func Init() {
	manager.RegisterAction(Action, RuleTypeCMD, RuleTypeString, RuleTypeHash, RuleTypeTool)
	manager.RegisterParser(RuleTypeCMD, parseCmd)
	manager.RegisterParser(RuleTypeString, parseString)
	manager.RegisterParser(RuleTypeHash, parseHash)

	manager.RegisterCreator(RuleTypeString, genString)
	manager.RegisterCreator(RuleTypeHash, genHash, genBatch)
	manager.RegisterCreator(RuleTypeTool, genTool)
}

type RedisAttr struct {
	Type    int32             // redis类型，string,hash,list,set,zset
	Package string            // 包命
	DbName  string            // 数据库名字
	Name    string            // pbname
	Desc    string            // 注释
	IsCache bool              // 是否需要ctx
	Key     *base.Index       // redis的key
	Field   *base.Index       // redis的field
	UUID    *base.Index       // 唯一搜索字段
	Ast     *domain.AstStruct // proto信息
	SubAst  *domain.AstStruct // proto信息
}

func (d *RedisAttr) IsString() bool {
	return d.Type == StringRedis
}

func (d *RedisAttr) IsHash() bool {
	return d.Type == HashRedis
}

// @gomaker:redis:cmd|GET,SET,MGET,MSET,...
func parseCmd(pbname, comment string) interface{} {
	return strings.Split(comment, ",")
}
