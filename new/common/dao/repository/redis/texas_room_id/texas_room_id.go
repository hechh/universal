package texas_room_id

import (
	"poker_server/common/dao"

	"github.com/spf13/cast"
)

const (
	ROOM_ID_BEGIN = 100000
	ROOM_ID_MAX   = 1000000
)

func GetKey() string {
	return "texas_room_id"
}

func Incr() (uint64, error) {
	val, err := dao.IncrBy("poker", GetKey(), 1)
	if err != nil {
		return 0, err
	}
	val = (val + ROOM_ID_BEGIN) % ROOM_ID_MAX
	return uint64(val), nil
}

func Get() (uint64, error) {
	val, err := dao.Get("poker", GetKey())
	if err != nil {
		return 0, err
	}
	id := cast.ToUint64(val)
	id = (id + ROOM_ID_BEGIN) % ROOM_ID_MAX
	return id, nil
}
