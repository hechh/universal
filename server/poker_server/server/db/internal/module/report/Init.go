package report

import (
	"poker_server/common/dao"
	"poker_server/common/dao/domain"
	"poker_server/common/pb"
)

func init() {
	dao.RegisterMysqlTable(domain.MYSQL_DB_PLAYER_DATA,
		&pb.TexasPlayerFlowReport{},
		//&pb.TexasPlayerReport{},
		&pb.TexasRoomReport{},
		&pb.TexasGameReport{})
}
