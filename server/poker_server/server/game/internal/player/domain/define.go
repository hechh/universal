package domain

import "poker_server/common/pb"

type IPlayerFun interface {
	Load(*pb.PlayerData) error // 加载数据
	Save(*pb.PlayerData) error // 保存数据
	LoadComplate() error       // 加载完成
	Change()                   // 数据变更通知
}
