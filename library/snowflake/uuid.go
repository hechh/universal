package snowflake

import (
	"time"
	"universal/common/pb"
	"universal/library/uerror"

	"github.com/sony/sonyflake"
)

var (
	uuidGen *sonyflake.Sonyflake
)

func Init(nn *pb.Node) error {
	if nn.Type > 255 || nn.Id > 255 {
		return uerror.New(1, -1, "SnowFlow初始化错误")
	}

	uuidGen = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		MachineID: func() (uint16, error) {
			MachineID := uint16(nn.Type)<<8 | uint16(nn.Id)
			return MachineID, nil
		},
		CheckMachineID: nil,
	})
	return nil
}

func GenUUID() (uint64, error) {
	return uuidGen.NextID()
}
