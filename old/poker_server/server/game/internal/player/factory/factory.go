package factory

import (
	"poker_server/common/pb"
	"poker_server/server/game/internal/player/domain"
	"poker_server/server/game/internal/player/playerfun"
)

var (
	LoadList = []pb.PlayerDataType{
		pb.PlayerDataType_PLAYER_DATA_BASE,
		pb.PlayerDataType_PLAYER_DATA_BAG,
	}
	FUNCS = make(map[pb.PlayerDataType]func(*playerfun.PlayerFun) domain.IPlayerFun)
)

func init() {
	FUNCS[pb.PlayerDataType_PLAYER_DATA_BASE] = playerfun.NewPlayerBaseFun
	FUNCS[pb.PlayerDataType_PLAYER_DATA_BAG] = playerfun.NewPlayerBagFun
}
