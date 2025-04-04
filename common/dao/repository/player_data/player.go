package player_data

/*
import (
	"hego/common/dao/internal/manager"
	"hego/common/pb"
	"hego/framework/basic/util"

	"github.com/astaxie/beego/orm"

	"github.com/golang/protobuf/proto"
)

type PlayerData struct {
	AccountID     uint64 `orm:"column(AccountID);pk"`
	PlayerLevel   uint64 `orm:"column(PlayerLevel)"`
	PlayerName    string `orm:"column(PlayerName)"`
	CreateTime    uint64 `orm:"column(CreateTime)"`
	LastDailyTime uint64 `orm:"column(LastDailyTime)"`
	BaseData      string `orm:"column(BaseData);type(mediumblob)"`
	SystemData    string `orm:"column(SystemData);type(mediumblob)"`
	BagData       string `orm:"column(BagData);type(mediumblob)"`
	EquipmentData string `orm:"column(EquipmentData);type(mediumblob)"`
	ClientData    string `orm:"column(ClientData);type(mediumblob)"`
	MailData      string `orm:"column(MailData);type(mediumblob)"`
	CrystalData   string `orm:"column(CrystalData);type(mediumblob)"`
	HeroData      string `orm:"column(HeroData);type(mediumblob)"`
}

func init() {
	orm.RegisterModel(new(PlayerData))
}

func (d *PlayerData) TableName() string {
	return "t_player"
}

func Get(dbname string, uid uint64) (*pb.PBPlayerData, error) {
	cli, err := manager.NewOrmer(dbname)
	if err != nil {
		return nil, err
	}
	data := &PlayerData{AccountID: uid}
	if err := cli.Read(data); err != nil {
		return nil, err
	}
	ret := &pb.PBPlayerData{
		Base:      new(pb.PBPlayerBase),
		System:    new(pb.PBPlayerSystem),
		Bag:       new(pb.PBPlayerBag),
		Equipment: new(pb.PBPlayerEquipment),
		Client:    new(pb.PBPlayerClientData),
		Hero:      new(pb.PBPlayerHero),
		Mail:      new(pb.PBPlayerMail),
		Crystal:   new(pb.PBPlayerCrystal),
	}
	if err := proto.Unmarshal(util.StringToBytes(data.BaseData), ret.Base); err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(util.StringToBytes(data.SystemData), ret.System); err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(util.StringToBytes(data.BagData), ret.Bag); err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(util.StringToBytes(data.EquipmentData), ret.Equipment); err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(util.StringToBytes(data.ClientData), ret.Client); err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(util.StringToBytes(data.HeroData), ret.Hero); err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(util.StringToBytes(data.MailData), ret.Mail); err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(util.StringToBytes(data.CrystalData), ret.Crystal); err != nil {
		return nil, err
	}
	return ret, nil
}
*/
