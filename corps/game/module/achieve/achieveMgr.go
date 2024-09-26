package achieve

//成就累计数据
import (
	"corps/base/cfgEnum"
	"corps/common"
	"corps/common/cfgData"
	"corps/common/serverCommon"
	"corps/framework/plog"
	"corps/pb"
)

type AchieveBase struct {
	mapData    map[uint32]map[common.AchieveKey]uint32 //成就数据
	mapAchieve map[cfgEnum.EAchieveSystemType]IAchieve
}

func NewAchieveBase() *AchieveBase {
	return &AchieveBase{
		mapData:    make(map[uint32]map[common.AchieveKey]uint32),
		mapAchieve: make(map[cfgEnum.EAchieveSystemType]IAchieve),
	}
}

func (this *AchieveBase) RegisterAchieve(emType cfgEnum.EAchieveSystemType, pAchieve IAchieve) {
	this.mapAchieve[emType] = pAchieve
}

// 加载
func (this *AchieveBase) LoadAchieveBase(info *pb.PBAchieveInfo) {
	if _, ok := this.mapData[info.AchieveType]; !ok {
		this.mapData[info.AchieveType] = make(map[common.AchieveKey]uint32)
	}
	pAchieveKey := common.AchieveKey{
		AchieveType: info.AchieveType,
	}
	if len(info.Params) > 0 {
		pAchieveKey.Param1 = info.Params[0]
	}
	if len(info.Params) > 1 {
		pAchieveKey.Param2 = info.Params[1]
	}

	this.mapData[info.AchieveType][pAchieveKey] = info.Value
}

// 存储
func (this *AchieveBase) SaveAchieveBase(info *[]*pb.PBAchieveInfo) {
	for key, list := range this.mapData {
		for skey, value := range list {
			*info = append(*info, &pb.PBAchieveInfo{
				AchieveType: key,
				Params:      []uint32{skey.Param1, skey.Param2},
				Value:       value,
			})
		}
	}
}

// 上阵英雄
func (this *AchieveBase) TriggerAchieveGameFightList(emAchieveType cfgEnum.AchieveType, mapQuality map[uint32]uint32) {
	this.mapData[uint32(emAchieveType)] = make(map[common.AchieveKey]uint32)
	for key, value := range mapQuality {
		pAchieveKey := common.AchieveKey{
			AchieveType: uint32(emAchieveType),
			Param1:      key,
		}

		this.mapData[uint32(emAchieveType)][pAchieveKey] = value

		//触发
		for _, pAchieve := range this.mapAchieve {
			pAchieve.Trigger(&pAchieveKey, value)
		}
	}

	//触发
	for _, pAchieve := range this.mapAchieve {
		pAchieve.Trigger(&common.AchieveKey{
			AchieveType: uint32(emAchieveType),
			Param1:      0,
		}, 1)
	}
}

// 触发成就
func (this *AchieveBase) DeleteAchieveBase(emAchieveType cfgEnum.AchieveType, params ...uint32) {
	if _, ok := this.mapData[uint32(emAchieveType)]; !ok {
		return
	}

	pAchieveKey := common.AchieveKey{
		AchieveType: uint32(emAchieveType),
	}
	if len(params) > 0 {
		pAchieveKey.Param1 = params[0]
	}
	if len(params) > 1 {
		pAchieveKey.Param2 = params[1]
	}

	delete(this.mapData[uint32(emAchieveType)], pAchieveKey)
}

// 触发成就
func (this *AchieveBase) AddAchieveBase(emAchieveType cfgEnum.AchieveType, uAdd uint32, params ...uint32) bool {
	if _, ok := this.mapData[uint32(emAchieveType)]; !ok {
		this.mapData[uint32(emAchieveType)] = make(map[common.AchieveKey]uint32)
	}

	pAchieveKey := common.AchieveKey{
		AchieveType: uint32(emAchieveType),
	}
	if len(params) > 0 {
		pAchieveKey.Param1 = params[0]
	}
	if len(params) > 1 {
		pAchieveKey.Param2 = params[1]
	}

	cfgAchieveType := cfgData.GetCfgAchieveTypeConfig(uint32(emAchieveType))
	if cfgAchieveType == nil {
		plog.Info("(this *PlayerSystemTaskFun) TriggerAchieve cfg error", emAchieveType, uAdd, params)
		return false
	}

	bChange := false
	if cfgAchieveType.CompareType == uint32(cfgEnum.ECompareType_Equal) {
		this.mapData[uint32(emAchieveType)][pAchieveKey] += uAdd
		bChange = true
	} else {
		//单参数直接替换
		if len(params) == 1 {
			this.mapData[uint32(emAchieveType)][pAchieveKey] += uAdd
			bChange = true
		} else {
			for pTmpKey, _ := range this.mapData[uint32(emAchieveType)] {
				if cfgAchieveType.Param1Type == uint32(cfgEnum.ECompareType_Equal) {
					if pTmpKey.Param1 != pAchieveKey.Param1 {
						continue
					}

					//更新最大值
					if cfgAchieveType.Param2Type == uint32(cfgEnum.ECompareType_NotSmall) {
						if pAchieveKey.Param2 <= pTmpKey.Param2 {
							continue
						}
					}
				} else if cfgAchieveType.Param1Type == uint32(cfgEnum.ECompareType_NotSmall) {

					if pAchieveKey.Param1 <= pTmpKey.Param1 {
						continue
					}

				}

				//需要删除map，重新增加一个
				delete(this.mapData[uint32(emAchieveType)], pTmpKey)
				break
			}
			this.mapData[uint32(emAchieveType)][pAchieveKey] = uAdd
			bChange = true
		}

	}

	if cfgAchieveType.IsSet > 0 {
		this.mapData[uint32(emAchieveType)][pAchieveKey] = uAdd
		bChange = true
	}

	//触发
	for _, pAchieve := range this.mapAchieve {
		pAchieve.Trigger(&pAchieveKey, uAdd)
	}

	return bChange
}

// 获取成就值
func (this *AchieveBase) GetAchieveValue(uAchieveType uint32, params ...uint32) uint32 {
	arrParams := make([]uint32, 0)
	for _, v := range params {
		arrParams = append(arrParams, v)
	}
	if len(arrParams) > 2 {
		if uAchieveType == uint32(cfgEnum.AchieveType_BattleMap) {
			arrParams[1] = serverCommon.MAKE_BATTLE_MAP(params[1], params[2])
		}
	}

	if _, ok := this.mapData[uAchieveType]; !ok {
		return 0
	}

	cfgAchieveType := cfgData.GetCfgAchieveTypeConfig(uAchieveType)
	if cfgAchieveType == nil {
		plog.Info("(this *PlayerSystemTaskFun) TriggerAchieve cfg error", uAchieveType, params)
		return 0
	}

	if cfgAchieveType.CompareType == uint32(cfgEnum.ECompareType_Equal) {
		pAchieveKey := common.AchieveKey{
			AchieveType: uAchieveType,
		}
		if len(arrParams) > 0 {
			pAchieveKey.Param1 = arrParams[0]
		}
		if len(arrParams) > 1 {
			pAchieveKey.Param2 = arrParams[1]
		}

		if _, ok := this.mapData[uAchieveType][pAchieveKey]; !ok {
			return 0
		}

		return this.mapData[uAchieveType][pAchieveKey]
	}

	//遍历获取
	uValue := uint32(0)
	for pAchieveType, value := range this.mapData[uAchieveType] {
		//判断
		if cfgAchieveType.Param1Type == uint32(cfgEnum.ECompareType_Equal) {
			if arrParams[0] != pAchieveType.Param1 {
				continue
			}

			if cfgAchieveType.Param2Type == uint32(cfgEnum.ECompareType_NotSmall) {
				if len(arrParams) < 1 {
					plog.Error("(this *AchieveBase) GetAchieveValue param is not 1 uAchieveType:%d", uAchieveType)
					continue
				}
				if arrParams[1] <= pAchieveType.Param2 {
					return value
				}
			} else if cfgAchieveType.Param2Type == uint32(cfgEnum.ECompareType_Equal) {
				if len(arrParams) < 1 {
					plog.Error("(this *AchieveBase) GetAchieveValue param is not 1 uAchieveType:%d", uAchieveType)
					continue
				}
				if arrParams[1] == pAchieveType.Param2 {
					return value
				}
			}
			return 0

		} else {
			if arrParams[0] > pAchieveType.Param1 {
				continue
			}
			uValue += value
		}
	}

	return uValue
}
