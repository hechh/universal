package room_router_data

import (
	"fmt"
	"poker_server/common/dao"
)

func GetKey(uid uint64) string {
	return fmt.Sprintf("user_desk:%d", uid)
}

func GetField(roomId uint64) string {
	return fmt.Sprintf("%d", roomId)
}

func HSet(uid uint64, addr string) error {
	return dao.HSet("poker", GetKey(uid), GetField(uid), addr)
}

func HDel(uid uint64) error {
	return dao.HDel("poker", GetKey(uid), GetField(uid))
}
