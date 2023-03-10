package model

import (
	"shared/common"
	"shared/protobuf/pb"
)

type RewardsResult struct {
	user *User
	*common.Rewards
	items      map[int32]int32
	characters []*Character
	equipments []*common.Equipment
	worldItems []*common.WorldItem
}

func NewRewardsResult(user *User) *RewardsResult {
	return &RewardsResult{
		user:    user,
		Rewards: common.NewRewards(),
	}
}

func (r *RewardsResult) Clear() {

	r.Rewards = common.NewRewards()

	if r.items != nil {
		r.items = nil
	}

	if r.characters != nil {
		r.characters = nil
	}

	if r.equipments != nil {
		r.equipments = nil
	}

	if r.worldItems != nil {
		r.worldItems = nil
	}
}

func (r *RewardsResult) AddItem(id, num int32) {
	if r.items == nil {
		r.items = map[int32]int32{}
	}

	r.items[id] = num
}

func (r *RewardsResult) AddCharacter(character *Character) {
	if r.characters == nil {
		r.characters = []*Character{}
	}

	r.characters = append(r.characters, character)
}

func (r *RewardsResult) AddEquipment(equipment *common.Equipment) {
	if r.equipments == nil {
		r.equipments = []*common.Equipment{}
	}

	r.equipments = append(r.equipments, equipment)
}

func (r *RewardsResult) AddWorldItem(worldItem *common.WorldItem) {
	if r.worldItems == nil {
		r.worldItems = []*common.WorldItem{}
	}

	r.worldItems = append(r.worldItems, worldItem)
}

func (r *RewardsResult) Append(ret *RewardsResult) {
	if ret == nil {
		return
	}

	r.Rewards.Append(ret.Rewards)

	if ret.items != nil {
		if r.items == nil {
			r.items = map[int32]int32{}
		}

		for k, v := range ret.items {
			r.items[k] = v
		}
	}

	if ret.characters != nil {
		if r.characters == nil {
			r.characters = []*Character{}
		}

		r.characters = append(r.characters, ret.characters...)
	}

	if ret.equipments != nil {
		if r.equipments == nil {
			r.equipments = []*common.Equipment{}
		}

		r.equipments = append(r.equipments, ret.equipments...)
	}

	if ret.worldItems != nil {
		if r.worldItems == nil {
			r.worldItems = []*common.WorldItem{}
		}

		r.worldItems = append(r.worldItems, ret.worldItems...)
	}
}

// 转为proto 用于和前端同步数据
func (r *RewardsResult) VOResourceResult() *pb.VOResourceResult {
	result := r.voResourceResultWithoutRewards()

	// 资源变动信息
	result.Rewards = r.Rewards.VOResourceMultiple()

	return result
}

// 转为proto 用于和前端同步数据
func (r *RewardsResult) VOMergedResourceResult() *pb.VOResourceResult {
	result := r.voResourceResultWithoutRewards()

	// 资源变动信息
	result.Rewards = r.Rewards.VOMergedResourceMultiple()

	return result
}

func (r *RewardsResult) voResourceResultWithoutRewards() *pb.VOResourceResult {
	result := &pb.VOResourceResult{}

	result.CommonResource = &pb.VOCommonResource{
		Gold:            r.user.Info.Gold.Value(),
		DiamondGift:     r.user.Info.DiamondGift.Value(),
		DiamondCash:     r.user.Info.DiamondCash.Value(),
		Energy:          r.user.Info.Energy.Value(),
		EnergyRefreshAt: r.user.Info.Energy.Last(),
		Exp:             r.user.Info.Exp.Value(),
		Level:           r.user.Info.Level.Value(),
		HeroExp:         r.user.HeroPack.Exp,
		GuildGold:       r.user.Guild.GuildGold.Value(),
	}

	// item变动
	if r.items != nil {
		Items := make([]*pb.VOItemInfo, 0, len(r.items))
		for id, num := range r.items {
			Items = append(Items, &pb.VOItemInfo{
				ItemId: id,
				Amount: num,
			})
		}
		result.Items = Items
	}

	// 新获得角色
	if r.characters != nil {
		Characters := make([]*pb.VOUserCharacter, 0, len(r.characters))
		for _, character := range r.characters {
			Characters = append(Characters, character.VOUserCharacter())
		}
		result.Characters = Characters
	}

	// 新获得装备
	if r.equipments != nil {
		Equipments := make([]*pb.VOUserEquipment, 0, len(r.equipments))
		for _, equipment := range r.equipments {
			Equipments = append(Equipments, equipment.VOUserEquipment())
		}
		result.Equipments = Equipments
	}

	// 新获得世界级道具
	if r.worldItems != nil {
		WorldItems := make([]*pb.VOUserWorldItem, 0, len(r.worldItems))
		for _, worldItem := range r.worldItems {
			WorldItems = append(WorldItems, worldItem.VOUserWorldItem())
		}
		result.WorldItems = WorldItems
	}

	return result
}
