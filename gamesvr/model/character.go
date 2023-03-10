package model

import (
	"encoding/json"
	"time"

	"gamesvr/manager"
	"shared/common"
	"shared/csv/entry"
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/glog"
)

const (
	CharacterEquipmentCounts = 4 // 角色可穿装备数量
	CharacterWorldItemCounts = 1 // 角色可穿世界级道具数量
)

type CharacterPack map[int32]*Character

func NewCharacterPack() *CharacterPack {
	return (*CharacterPack)(&(map[int32]*Character{}))
}

func (cp *CharacterPack) Contains(id int32) bool {
	_, ok := (*cp)[id]
	return ok
}

func (cp *CharacterPack) NewCharacter(id int32) *Character {
	character := NewCharacter(id)
	(*cp)[character.ID] = character
	return character
}

func (cp *CharacterPack) Put(character *Character) {
	(*cp)[character.ID] = character
}

func (cp *CharacterPack) Get(id int32) (*Character, error) {
	character, ok := (*cp)[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrCharacterNotFound, id)
	}
	return character, nil
}

func (cp *CharacterPack) VOUserCharacter() []*pb.VOUserCharacter {
	vos := make([]*pb.VOUserCharacter, 0, len(*cp))

	for _, v := range *cp {

		vos = append(vos, v.VOUserCharacter())
	}

	return vos
}

type Character struct {
	ID               int32           `json:"id"`
	Exp              int32           `json:"exp"`
	Level            int32           `json:"level"`
	Star             int32           `json:"star"`
	Stage            int32           `json:"stage"`
	CTime            int64           `json:"c_time"`
	Skills           map[int32]int32 `json:"skills"`
	HeroId           int32           `json:"hero_id"`
	HeroSlot         int32           `json:"hero_slot"`
	Power            int32           `json:"power"`
	CanYggdrasilTime int64           `json:"can_yggdrasil_time"`
	Rarity           int32           `json:"rarity"`
	Equipments       []int64         `json:"equipments"`
	WorldItem        int64           `json:"world_item"`
}

func (c *Character) VOUserCharacter() *pb.VOUserCharacter {
	skills := make([]*pb.VOCharacterSkill, 0, len(c.Skills))
	for skillId, level := range c.Skills {
		skills = append(skills, &pb.VOCharacterSkill{
			SkillId: skillId,
			Level:   level,
		})
	}

	return &pb.VOUserCharacter{
		CharacterId:      c.ID,
		Exp:              c.Exp,
		Level:            c.Level,
		Star:             c.Star,
		Stage:            c.Stage,
		CreateAt:         c.CTime,
		Skills:           skills,
		HeroId:           c.HeroId,
		HeroPos:          c.HeroSlot,
		Power:            c.Power,
		CanYggdrasilTime: c.CanYggdrasilTime,
		Rarity:           c.Rarity,
		// Equipments:      c.Equipments,
		// WorldItem:        c.WorldItem,
	}
}

func NewCharacter(id int32) *Character {
	character := &Character{
		ID:               id,
		Exp:              0,
		Level:            1,
		Star:             1,
		Stage:            0,
		CTime:            time.Now().Unix(),
		Skills:           map[int32]int32{},
		HeroId:           0,
		HeroSlot:         0,
		Power:            0,
		CanYggdrasilTime: time.Now().Unix(),
		Rarity:           0,
		Equipments:       make([]int64, CharacterEquipmentCounts, CharacterEquipmentCounts),
		WorldItem:        0,
	}

	return character
}

func (c *Character) EquipEquipment(equipmentUID int64, part int32) int64 {
	i := int(part - 1)

	if i >= len(c.Equipments) {
		return 0
	}

	oldEquipmentUID := c.Equipments[i]

	c.Equipments[i] = equipmentUID

	return oldEquipmentUID
}

func (c *Character) StripEquipment(equipmentUID int64) bool {
	for i, v := range c.Equipments {
		if v == equipmentUID {
			c.Equipments[i] = 0
			return true
		}
	}

	return false
}

func (c *Character) GetId() int32 {
	return c.ID
}

func (c *Character) FetchSkillLevelCanLevelUp(skillNum int32) (int32, error) {
	lv, ok := c.Skills[skillNum]
	// 未激活
	if !ok {
		return 0, errors.Swrapf(common.ErrCharacterSkillLevelUnlock, c.ID, skillNum)
	}
	if lv >= entry.CharacterSkillsLevelUpCount {
		return 0, common.ErrCharacterSkillCannotLevelUp
	}

	return lv, nil
}

// func (c *Character) addExp(addNum int32) {
//	if addNum <= 0 {
//		return
//	}
//
//	nowLevel := c.Level
//	nowExp := c.Exp + addNum
//	isMaxLv := false
//	for ; ; nowLevel++ {
//		// 升到本级所需经验值
//		exp, ok := manager.CSV.Character.GetExpToNextLevel(nowLevel)
//		if !ok {
//			// 升到最高级了
//			isMaxLv = true
//			break
//		}
//		if exp > nowExp {
//			// 当前经验不可升级
//			break
//		}
//		nowExp = nowExp - exp
//	}
//	// 升级
//	c.SetLevel(nowLevel)
//	if isMaxLv {
//		nowExp = 0
//	}
//	c.SetExp(nowExp)
//
// }

func (cp *CharacterPack) Marshal() ([]byte, error) {
	return json.Marshal(cp)
}

func (cp *CharacterPack) Unmarshal(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	err := json.Unmarshal(data, &cp)
	if err != nil {
		glog.Errorf("json.Unmarshal error: %v", err)
		return err
	}

	return nil
}

func (cp *CharacterPack) CalcCharacterStarUpCost(charaId, star int32) (common.Rewards, error) {
	cfg, err := manager.CSV.Character.Star(charaId, star)
	if err != nil {
		return nil, err
	}

	return *cfg.Cost, nil
}

func (cp *CharacterPack) CalcCharacterStageUpCost(charaId, stage int32) (common.Rewards, error) {
	cfg, err := manager.CSV.Character.Stage(charaId, stage)
	if err != nil {
		return nil, err
	}

	return *cfg.Cost, nil
}

func (c *Character) GetStar() int32 {
	return c.Star
}

func (c *Character) SetStar(star int32) {
	c.Star = star
}

func (c *Character) GetStage() int32 {
	return c.Stage
}

func (c *Character) SetStage(stage int32) {
	c.Stage = stage

}

func (c *Character) GetLevel() int32 {
	return c.Level
}

func (c *Character) SetLevel(level int32, user *User) {
	c.Level = level

}

func (c *Character) GetExp() int32 {
	return c.Exp
}

func (c *Character) SetExp(exp int32) {
	c.Exp = exp

}

func (c *Character) GetMaxLevel(u *User) int32 {
	characterMaxLv, err := manager.CSV.TeamLevelCache.GetCharacterMaxLv(u.Info.Level.Value())
	if err != nil {
		return 0
	}

	return characterMaxLv
}

func (c *Character) SetRare(rare int32) {
	c.Rarity = rare
}

func (c *Character) GetRare() int32 {
	return c.Rarity
}

func (c *Character) SetHero(heroId, slot int32) {
	c.HeroId = heroId
	c.HeroSlot = slot
}

func (c *Character) GetHero() (int32, int32) {
	return c.HeroId, c.HeroSlot
}
