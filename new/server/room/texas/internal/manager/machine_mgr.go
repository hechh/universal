package manager

import (
	"poker_server/common/pb"
	"poker_server/server/room/texas/domain"
)

var (
	stateMgr = make(map[pb.TexasGameState]domain.IState)
)

func RegisterState(val pb.TexasGameState, state domain.IState) {
	if _, ok := stateMgr[val]; ok {
		panic("StateMgr: 注册状态失败，状态已存在")
	}
	stateMgr[val] = state
}

func GetState(val pb.TexasGameState) domain.IState {
	return stateMgr[val]
}
