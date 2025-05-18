package texas_room_list

import (
	"fmt"
	"poker_server/common/dao"
	"strings"

	"github.com/spf13/cast"
)

func GetUniqueId(gameType, coinType int32) uint64 {
	// 生成唯一ID
	return uint64(gameType)<<32 | uint64(coinType)
}

func GetKey(id uint64) string {
	return fmt.Sprintf("texas_room_list_%d", id)
}

// 加载房间列表
func Get(id uint64) (rets []uint64, err error) {
	// 从redis中读取房间列表
	val, err := dao.Get("poker", GetKey(id))
	if err != nil {
		return nil, err
	}
	for _, vv := range strings.Split(val, ",") {
		rets = append(rets, cast.ToUint64(vv))
	}
	return
}

func Set(id uint64, ids ...uint64) error {
	if len(ids) <= 0 {
		return nil
	}

	strs := []string{}
	for _, v := range ids {
		strs = append(strs, cast.ToString(v))
	}
	str := strings.Join(strs, ",")
	return dao.Set("poker", GetKey(id), str)
}
