package base

import (
	"math/rand"
	"poker_server/common/pb"
	"time"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

func Intn(val int) int {
	return rand.Intn(val)
}

func Int32n(val int32) int32 {
	return rand.Int31n(val)
}

func GetPlayerState(usr *pb.TexasPlayerData) (rets []pb.TexasPlayerState) {
	for i := pb.TexasPlayerState_TPS_JOIN_ROOM; i <= pb.TexasPlayerState_TPS_QUIT_ROOM; i++ {
		if usr.PlayerState&(1<<(i-1)) > 0 {
			rets = append(rets, i)
		}
	}
	return rets
}

func SetPlayerState(usr *pb.TexasPlayerData, states ...pb.TexasPlayerState) {
	usr.PlayerState = 0
	AddPlayerState(usr, states...)
}

func AddPlayerState(usr *pb.TexasPlayerData, states ...pb.TexasPlayerState) {
	for _, s := range states {
		usr.PlayerState |= (1 << (s - 1))
	}
}

func DelPlayerState(usr *pb.TexasPlayerData, states ...pb.TexasPlayerState) {
	for _, s := range states {
		usr.PlayerState &= ^(1 << (s - 1))
	}
}

func HasPlayerState(usr *pb.TexasPlayerData, states ...pb.TexasPlayerState) bool {
	for _, s := range states {
		if IsPlayerState(usr, s) {
			return true
		}
	}
	return false
}

func IsPlayerState(usr *pb.TexasPlayerData, state pb.TexasPlayerState) bool {
	return usr.PlayerState&(1<<(state-1)) > 0
}

func ToRet(code pb.ErrorCode, msg string) *pb.RspHead {
	return &pb.RspHead{
		Code: int32(code),
		Msg:  msg,
	}
}
