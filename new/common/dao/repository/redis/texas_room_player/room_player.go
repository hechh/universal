package texas_room_player

import (
	"fmt"
	"poker_server/common/dao"

	"github.com/spf13/cast"
)

func GetKey() string {
	return "texas_room"
}

func GetField(uid uint64) string {
	return fmt.Sprintf("texas_room_player_%d", uid)
}

func HSet(uid uint64, roomId uint64) error {
	return dao.HSet("poker", GetKey(), GetField(uid), roomId)
}

func HGet(uid uint64) (uint64, error) {
	val, err := dao.HGet("poker", GetKey(), GetField(uid))
	if err != nil {
		return 0, err
	}
	return cast.ToUint64(val), nil
}

func HDel(uid uint64) error {
	return dao.HDel("poker", GetKey(), GetField(uid))
}
