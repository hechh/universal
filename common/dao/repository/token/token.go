package token

import (
	"fmt"
	"universal/common/dao/domain"
	"universal/common/dao/internal/manager"
)

func GetLoginToken(uid uint64) (key string, err error) {
	cli := manager.GetRedisByUID(uid)
	if cli != nil {
		err = fmt.Errorf("RedisClient(%d) not found", uid)
		return
	}
	key, err = cli.Get(fmt.Sprintf("%s_%d", domain.ERK_LoginToken, uid))
	return
}
