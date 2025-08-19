package message

import (
	"poker_server/common/pb"
	"poker_server/framework/request"
)

func Init() {
	initActor()
	initCmd()
}

func initActor() {
	request.RegisterActor(pb.NodeType_NodeTypeGate, "Player", 0)        // uid
	request.RegisterActor(pb.NodeType_NodeTypeGate, "GatePlayerMgr", 1) // 唯一值

	request.RegisterActor(pb.NodeType_NodeTypeDb, "DbGeneratorMgr", 2)          // 随意值
	request.RegisterActor(pb.NodeType_NodeTypeDb, "PlayerDataMgr", 3)           // 随意值
	request.RegisterActor(pb.NodeType_NodeTypeDb, "ReportDataMgr", 4)           // 随意值
	request.RegisterActor(pb.NodeType_NodeTypeDb, "RoomInfoMgr", 5)             // 随意值
	request.RegisterActor(pb.NodeType_NodeTypeDb, "DbRummyRoomMgr", 6)          // 随意值
	request.RegisterActor(pb.NodeType_NodeTypeDb, "RummyExtSettleMatchPool", 7) // 随意值
	request.RegisterActor(pb.NodeType_NodeTypeDb, "RummySettlePool", 8)         // 随意值
	request.RegisterActor(pb.NodeType_NodeTypeDb, "RummySettleMatchPool", 9)    // 随意值
	request.RegisterActor(pb.NodeType_NodeTypeDb, "DbTexasRoomMgr", 10)         // 随意值
	request.RegisterActor(pb.NodeType_NodeTypeDb, "UserInfoMgr", 11)            // 随意值

	request.RegisterActor(pb.NodeType_NodeTypeBuilder, "BuilderTexasGenerator", 12) // 随意值
	request.RegisterActor(pb.NodeType_NodeTypeBuilder, "BuilderRummyGenerator", 13) // 随意值

	request.RegisterActor(pb.NodeType_NodeTypeGame, "Player", 0)     // uid
	request.RegisterActor(pb.NodeType_NodeTypeGame, "PropMgr", 14)   // 随意值
	request.RegisterActor(pb.NodeType_NodeTypeGame, "PlayerMgr", 15) // 随意值

	request.RegisterActor(pb.NodeType_NodeTypeMatch, "MatchRummyRoomMgr", 16) // 随意值
	request.RegisterActor(pb.NodeType_NodeTypeMatch, "MatchRummyRoom", 17)    // 随意值
	request.RegisterActor(pb.NodeType_NodeTypeMatch, "MatchTexasRoomMgr", 18) // 随意值
	request.RegisterActor(pb.NodeType_NodeTypeMatch, "MatchTexasRoom", 19)    // 随意值
	request.RegisterActor(pb.NodeType_NodeTypeMatch, "SngRoomMgr", 20)        // 随意值

	request.RegisterActor(pb.NodeType_NodeTypeRoom, "TexasGameMgr", 21) // 随意值
	request.RegisterActor(pb.NodeType_NodeTypeRoom, "TexasGame", 22)    // roomId
	request.RegisterActor(pb.NodeType_NodeTypeRoom, "SngTexasGame", 23) // roomId
	request.RegisterActor(pb.NodeType_NodeTypeRoom, "RummyGameMgr", 24) // 随意值
	request.RegisterActor(pb.NodeType_NodeTypeRoom, "RummyGame", 25)    // roomId
}

func initCmd() {
	request.RegisterCmd(pb.NodeType_NodeTypeGate, pb.CMD_GATE_LOGIN_REQUEST, "Player.Login") // 登录请求
	request.RegisterCmd(pb.NodeType_NodeTypeGate, pb.CMD_KICK_PLAYER_NOTIFY, "Player.Kick")  // 剔除玩家通知

	request.RegisterCmd(pb.NodeType_NodeTypeGame, pb.CMD_GATE_HEART_REQUEST, "Player.HeartRequest")         // 心跳请求
	request.RegisterCmd(pb.NodeType_NodeTypeGame, pb.CMD_GET_BAG_REQ, "Player.GetBagReq")                   // 查询背包
	request.RegisterCmd(pb.NodeType_NodeTypeGame, pb.CMD_TEXAS_JOIN_ROOM_REQ, "PlayerMgr.TexasJoinRoomReq") // 德州扑克加入房间请求
	request.RegisterCmd(pb.NodeType_NodeTypeGame, pb.CMD_TEXAS_QUIT_ROOM_REQ, "PlayerMgr.TexasQuitRoomReq") // 离开房间
	request.RegisterCmd(pb.NodeType_NodeTypeGame, pb.CMD_TEXAS_BUY_IN_REQ, "PlayerMgr.TexasBuyInReq")       // 买入请求
	request.RegisterCmd(pb.NodeType_NodeTypeGame, pb.CMD_TEXAS_CHANGE_ROOM_REQ, "PlayerMgr.TexasChangeReq") // 德州换房间

	request.RegisterCmd(pb.NodeType_NodeTypeGame, pb.CMD_RUMMY_JOIN_ROOM_REQ, "PlayerMgr.RummyJoinRoomReq")       // Rummy加入房间请求
	request.RegisterCmd(pb.NodeType_NodeTypeGame, pb.CMD_RUMMY_QUIT_ROOM_REQ, "PlayerMgr.RummyQuitRoomReq")       // Rummy退出房间请求
	request.RegisterCmd(pb.NodeType_NodeTypeGame, pb.CMD_RUMMY_CHANGE_ROOM_REQ, "PlayerMgr.RummyChangeRoomReq")   // Rummy更换房间请求
	request.RegisterCmd(pb.NodeType_NodeTypeGame, pb.CMD_RUMMY_MATCH_REQ, "PlayerMgr.RummyMatchReq")              // Rummy分支玩法开启匹配
	request.RegisterCmd(pb.NodeType_NodeTypeGame, pb.CMD_RUMMY_CANCEL_MATCH_REQ, "PlayerMgr.RummyCancelMatchReq") // Rummy分支玩法退出匹配
	request.RegisterCmd(pb.NodeType_NodeTypeGame, pb.CMD_RUMMY_GIVE_UP_REQ, "PlayerMgr.RummyGiveUpReq")           // Rummy分支玩法退出房间

	request.RegisterCmd(pb.NodeType_NodeTypeMatch, pb.CMD_TEXAS_ROOM_LIST_REQ, "MatchTexasRoom.RoomListReq") // 德州房间列表请求
	request.RegisterCmd(pb.NodeType_NodeTypeMatch, pb.CMD_RUMMY_ROOM_LIST_REQ, "MatchRummyRoom.RoomListReq") // 获取房间列表
	request.RegisterCmd(pb.NodeType_NodeTypeMatch, pb.CMD_SNG_JOIN_ROOM_REQ, "SngRoomMgr.SngJoinRoomReq")    // sng加入房间

	request.RegisterCmd(pb.NodeType_NodeTypeDb, pb.CMD_TEXAS_GAME_REPORT_REQ, "ReportDataMgr.GameReportReq")          // 查询德州牌局记录
	request.RegisterCmd(pb.NodeType_NodeTypeDb, pb.CMD_RUMMY_MATCH_SELECT_REQ, "RummySettleMatchPool.Select")         // 游戏内查询玩家输赢记录
	request.RegisterCmd(pb.NodeType_NodeTypeDb, pb.CMD_GET_USER_INFO_REQ, "UserInfoMgr.Query")                        // 查询缓存玩家信息
	request.RegisterCmd(pb.NodeType_NodeTypeDb, pb.CMD_GET_TEXAS_GAME_REPORT_REQ, "ReportDataMgr.GetTexasGameReport") // 查询牌局信息

	request.RegisterCmd(pb.NodeType_NodeTypeRoom, pb.CMD_TEXAS_SIT_DOWN_REQ, "TexasGameMgr.SitDownReq") // 坐下请求
	request.RegisterCmd(pb.NodeType_NodeTypeRoom, pb.CMD_TEXAS_STAND_UP_REQ, "TexasGameMgr.StandUpReq") // 站起请求
	request.RegisterCmd(pb.NodeType_NodeTypeRoom, pb.CMD_TEXAS_DO_BET_REQ, "TexasGameMgr.DoBetReq")     // 下注请求

	request.RegisterCmd(pb.NodeType_NodeTypeRoom, pb.CMD_RUMMY_SAVE_CARD_GROUP_REQ, "RummyGame.SaveCardGroup")  // Rummy玩家保存手牌
	request.RegisterCmd(pb.NodeType_NodeTypeRoom, pb.CMD_RUMMY_OPR_CARD_REQ, "RummyGame.OprCardReq")            // Rummy玩家操作
	request.RegisterCmd(pb.NodeType_NodeTypeRoom, pb.CMD_RUMMY_FIX_CARD_REQ, "RummyGame.FixCard")               // Rummy玩家胡牌操作
	request.RegisterCmd(pb.NodeType_NodeTypeRoom, pb.CMD_RUMMY_READY_ROOM_REQ, "RummyGame.ReadyRoomReq")        // Rummy玩家准备请求
	request.RegisterCmd(pb.NodeType_NodeTypeRoom, pb.CMD_RUMMY_GET_OUT_CARDS_REQ, "RummyGame.GetRummyOutCards") // Rummy玩家请求出牌列表
	request.RegisterCmd(pb.NodeType_NodeTypeRoom, pb.CMD_TEXAS_STATISTICS_REQ, "TexasGameMgr.StatisticsReq")    // 查询房间统计信息
	request.RegisterCmd(pb.NodeType_NodeTypeRoom, pb.CMD_SNG_RANK_REQ, "SngTexasGame.SngRankReq")               // 查询sng排行榜

}
