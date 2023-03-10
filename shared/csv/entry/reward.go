package entry

import (
	"sync"

	"shared/common"
	"shared/csv/static"
	"shared/utility/errors"
)

const (
	RewardCfgItemDataConfig = "cfg_item_data_config"
)

type Reward struct {
	sync.RWMutex

	rewardsType map[int32]int32

	// parameter
	autoRandomParam       map[int32][]int32
	autoRandomRemoveParam map[int32][]int32
}

func NewReward() *Reward {
	return &Reward{
		rewardsType:           map[int32]int32{},
		autoRandomParam:       map[int32][]int32{},
		autoRandomRemoveParam: map[int32][]int32{},
	}
}

func (r *Reward) Check(config *Config) error {
	return nil
}

func (r *Reward) Reload(config *Config) error {
	r.Lock()
	defer r.Unlock()

	rewardsType := map[int32]int32{}
	autoRandom := map[int32][]int32{}
	autoRandomRemove := map[int32][]int32{}

	for id, reward := range config.CfgItemDataConfig.GetAllData() {
		rewardsType[id] = reward.ItemType
		switch reward.ItemType {
		case static.ItemTypeCurrency: // 一般道具（金币，钻石，体力等）
		case static.ItemTypeGacha: // 抽卡道具
		case static.ItemTypeCharExp: // 角色经验道具
		case static.ItemTypeCharSkill: // 角色进阶/技能道具
		case static.ItemTypeCharPiece: // 角色碎片
		case static.ItemTypeWorldItemExp: // 神器经验道具
		case static.ItemTypeWorldItemBreak: // 神器突破道具
		case static.ItemTypeEquipExp: // 装备经验道具
		case static.ItemTypeEquipBreak: // 装备突破道具
		case static.ItemTypeGraveyardAccelerate: // 模拟经营加速道具
		case static.ItemTypeGraveyardGetProduct: // 模拟经营使用道具直接获得产出
		case static.ItemTypeGraveyardProduceBuff: // 模拟经营生产buff道具
		case static.ItemTypeEnergyItem: // 体力药水
		case static.ItemTypeCommon: // 一般道具
		case static.ItemTypeGiftSelectOne: // N选1礼包
		case static.ItemTypeGiftRandomDrop: // 随机普通礼包
		case static.ItemTypeAutoRandom: // 获得后不进背包，自动使用 且 使用后执行一般掉落逻辑的道具-即按照原掉落规则纯随机
			var dropIDs []int32
			dropIDs = append(dropIDs, reward.UseParam...)

			autoRandom[id] = dropIDs
		case static.ItemTypeAutoRandomRemove: // 获得后不进背包，自动使用 且 使用后执行奖池掉落逻辑的道具
			var dropIDs []int32
			dropIDs = append(dropIDs, reward.UseParam...)

			autoRandomRemove[id] = dropIDs
		case static.ItemTypeYggItemSpecial: // 世界探索特殊道具（魔石，不占用世界探索格子，不能存到仓库）
		case static.ItemTypeEquipment: // 装备
		case static.ItemTypeWorldItem: // 世界级道具
		case static.ItemTypeCharacter: // 角色
		}
	}

	r.rewardsType = rewardsType
	common.SetRewardsType(rewardsType)
	r.autoRandomParam = autoRandom
	r.autoRandomRemoveParam = autoRandomRemove

	return nil
}

func (r *Reward) RewardsType() map[int32]int32 {
	r.RLock()
	defer r.RUnlock()

	return r.rewardsType
}

func (r *Reward) AutoRandomParam(id int32) ([]int32, error) {
	r.RLock()
	defer r.RUnlock()

	param, ok := r.autoRandomParam[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, RewardCfgItemDataConfig, id)
	}

	return param, nil
}

func (r *Reward) AutoRandomRemoveParam(id int32) ([]int32, error) {
	r.RLock()
	defer r.RUnlock()

	param, ok := r.autoRandomRemoveParam[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, RewardCfgItemDataConfig, id)
	}

	return param, nil
}

// case csv.ConstItem.Common: // 一般道具（金币，钻石，体力等）
// 		switch reward.ID {
// 		case csv.ConstCommonResource.TeamExp: // 账号经验
// 			u.Info.addExp(reward.Num)
// 		case csv.ConstCommonResource.Money: // 金币
// 			u.Info.Gold.Plus(reward.Num)
// 		case csv.ConstCommonResource.DiamondGift: // 免费钻石
// 			u.Info.DiamondGift.Plus(reward.Num)
// 		case csv.ConstCommonResource.Energy: // 体力
// 			u.Info.Energy.Plus(reward.Num)
// 		case csv.ConstCommonResource.DiamondCash: // 付费钻石
// 			u.Info.DiamondCash.Plus(reward.Num)
// 		}
//
// 	case csv.ConstItem.Gacha: // 抽卡道具
// 	case csv.ConstItem.CharExp: // 角色经验道具
// 	case csv.ConstItem.CharSkill: // 角色进阶/技能道具
// 	case csv.ConstItem.CharPiece: // 角色碎片
// 	case csv.ConstItem.WorldItemExp: // 神器经验道具
// 	case csv.ConstItem.WorldItemBreak: // 神器突破道具
// 	case csv.ConstItem.EquipExp: // 装备经验道具
// 	case csv.ConstItem.EquipBreak: // 装备突破道具
// 	case csv.ConstItem.GraveyardAccelerate: // 模拟经营加速道具
// 	case csv.ConstItem.GiftSelectOne: // N选1礼包
// 	case csv.ConstItem.GiftRandomDrop: // 随机普通礼包
// 	case csv.ConstItem.AutoRandom: // 获得后不进背包，自动使用 且 使用后执行一般掉落逻辑的道具-即按照原掉落规则纯随机
// 	case csv.ConstItem.AutoRandomRemove: // 获得后不进背包，自动使用 且 使用后执行奖池掉落逻辑的道具
