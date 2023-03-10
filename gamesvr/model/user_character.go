package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/csv/entry"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/statistic/bilog"
	"shared/statistic/logreason"
	"shared/utility/errors"
	"shared/utility/glog"
)

func (u *User) CharacterStarUp(charaId int32) (*Character, *pb.VOResourceResult, error) {
	err := u.CheckActionUnlock(static.ActionIdTypeCharacterstarup)
	if err != nil {
		return nil, nil, err
	}

	chara, err := u.CharacterPack.Get(charaId)
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}

	targetStar := chara.GetStar() + 1

	cost, err := u.CharacterPack.CalcCharacterStarUpCost(charaId, targetStar)
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}

	err = u.CheckRewardsEnough(&cost)
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}

	reason := logreason.NewReason(logreason.CharacterStarUp)
	err = u.CostRewards(&cost, reason)
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}

	chara.SetStar(targetStar)

	cfg, err := manager.CSV.Character.Star(charaId, targetStar)
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}
	if cfg.RarityUp > 0 {
		chara.SetRare(cfg.RarityUp)
	}

	u.TriggerQuestUpdate(static.TaskTypeCharacterStar, charaId, chara.GetStar())
	u.QuestCheckUnlock(static.ConditionTypeCharaStar)

	u.ActiveCharacterSkill(charaId)
	u.skillAutoLvUp(charaId)

	charaPower, err := u.CalCharacterCombatPower(chara)
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}
	chara.Power = charaPower

	u.BICharaOp(chara, bilog.CharaOpStar)

	return chara, u.VOResourceResult(), nil
}

func (u *User) CharacterStageUp(charaId int32) (*Character, *pb.VOResourceResult, error) {

	err := u.CheckActionUnlock(static.ActionIdTypeCharacterstageup)
	if err != nil {
		return nil, nil, err
	}

	chara, err := u.CharacterPack.Get(charaId)
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}

	if !u.characterCanBeStageUp(chara) {
		return nil, nil, common.ErrCharacterCanNotBeStageUp
	}

	targetStage := chara.GetStage() + 1

	cost, err := u.CharacterPack.CalcCharacterStageUpCost(charaId, targetStage)
	if err != nil {
		return nil, nil, err
	}

	err = u.CheckRewardsEnough(&cost)
	if err != nil {
		return nil, nil, err
	}

	reason := logreason.NewReason(logreason.CharacterStageUp)
	err = u.CostRewards(&cost, reason)
	if err != nil {
		return nil, nil, err
	}

	chara.SetStage(targetStage)

	u.TriggerQuestUpdate(static.TaskTypeCharacterStageCount, chara.GetRare(), targetStage, targetStage-1)
	u.TriggerQuestUpdate(static.TaskTypeCharacterStage, charaId, targetStage)

	u.ActiveCharacterSkill(charaId)
	u.skillAutoLvUp(charaId)
	charaPower, err := u.CalCharacterCombatPower(chara)
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}
	chara.Power = charaPower
	u.BICharaOp(chara, bilog.CharaOpStage)

	return chara, u.VOResourceResult(), nil
}

func (u *User) characterCanBeStageUp(chara *Character) bool {
	cfg, err := manager.CSV.Character.Stage(chara.GetId(), chara.GetStage()+1)

	if err != nil {
		return false
	}

	if cfg.LevelLimit > u.Info.GetLevel() {
		return false
	}

	return true
}

func (u *User) CharacterLvUp(characterId int32, rewards *common.Rewards) (*Character, error) {

	err := u.CheckActionUnlock(static.ActionIdTypeCharacterlvup)
	if err != nil {
		return nil, err
	}

	character, err := u.CharacterPack.Get(characterId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	items, err := manager.CSV.Item.GetSortedCharExpItem(rewards)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	characterMaxLv, err := manager.CSV.TeamLevelCache.GetCharacterMaxLv(u.Info.Level.Value())
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// 计算经验缺口
	totalExp := manager.CSV.Character.GetNeedExp(character.Exp, characterMaxLv)

	consume, addExp, err := manager.CSV.Item.CalRealConsume(items, totalExp)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// 消耗道具
	reason := logreason.NewReason(logreason.CharacterLevelUp)
	err = u.CostRewards(consume, reason)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	u.CharacterAddExp(characterId, addExp)

	u.TriggerQuestUpdate(static.TaskTypeCharacterLevelStrengthenTimes, character.GetRare())

	charaPower, err := u.CalCharacterCombatPower(character)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	character.Power = charaPower

	return character, nil
}

func (u *User) CharacterAddExp(characterId int32, addExp int32) {
	character, err := u.CharacterPack.Get(characterId)
	if err != nil {
		return
	}

	oldLevel := character.GetLevel()

	// 加经验
	u.AddExpCommon(addExp, character, manager.CSV.Character)
	u.ActiveCharacterSkill(characterId)
	u.skillAutoLvUp(characterId)

	u.BICharaOp(character, bilog.CharaOpExp)

	if character.GetLevel() > oldLevel {
		u.TriggerQuestUpdate(static.TaskTypeCharacterLevelUpCount, character, oldLevel)
		u.TriggerQuestUpdate(static.TaskTypeCharacterLevelCount, character, oldLevel)
		u.TriggerQuestUpdate(static.TaskTypeCharacterLevel, characterId, character.GetLevel())
		u.QuestCheckUnlock(static.ConditionTypeCharaLevel)
		u.BICharaOp(character, bilog.CharaOpLevel)
	}
}

func (u *User) CharacterSkillLvUp(characterId int32, skillNum int32, amount int32) (*Character, error) {

	err := u.CheckActionUnlock(static.ActionIdTypeCharacterskillup)
	if err != nil {
		return nil, err
	}

	character, err := u.CharacterPack.Get(characterId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	lv, err := character.FetchSkillLevelCanLevelUp(skillNum)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// skill配置
	_, err = manager.CSV.Character.FindCanUpgradeSkill(characterId, skillNum)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	toLv := lv + amount

	// 消耗
	consume := common.NewRewards()
	for i := lv + 1; i <= toLv; i++ {
		skillLevel, err := manager.CSV.Character.SkillLevel(skillNum, i)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		// 自动升星的技能不可手动升级
		if !skillLevel.CanSkillLevelUp() {
			return nil, common.ErrCharacterSkillCannotLevelUp
		}

		// check 是否解锁
		if !skillLevel.CanSkillUnlock(character.Level, character.GetStar(), character.GetStage()) {
			return nil, errors.Swrapf(common.ErrCharacterSkillLevelUnlock, character.ID)
		}
		consume.Append(skillLevel.Cost)
	}

	// 消耗
	reason := logreason.NewReason(logreason.CharacterSkillLevelUp)
	err = u.CostRewards(consume, reason)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	u.ManualUpdateSkillLevel(characterId, skillNum, toLv)
	u.ActiveCharacterSkill(characterId)
	u.skillAutoLvUp(characterId)

	charaPower, err := u.CalCharacterCombatPower(character)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	character.Power = charaPower

	return character, nil
}

func (u *User) ManualUpdateSkillLevel(characterId, skillNum, lv int32) {
	c, err := u.CharacterPack.Get(characterId)
	if err != nil {
		return
	}
	u.SetCharacterSkillLv(characterId, skillNum, lv)
	// 有关联技能的技能不可手动升级，会随着关联技能升级而自动升级
	associated, ok := manager.CSV.Character.GetSkillsAssociated(c.ID, skillNum)
	if ok {
		u.SetCharacterSkillLv(characterId, associated, lv)
	}

}

// 角色初始获得，升星，升阶 会触发技能解锁
func (u *User) ActiveCharacterSkill(characterId int32) {
	c, err := u.CharacterPack.Get(characterId)
	if err != nil {
		return
	}
	skillMap, err := manager.CSV.Character.AllSkill(c.ID)
	if err != nil {
		glog.Errorf("AllSkill is nil characterId:%v", characterId)
		return
	}

	for skillNum := range skillMap {
		_, ok := c.Skills[skillNum]
		if ok {
			// 已激活
			continue
		}
		// 取一级的配置
		skillLevel, err := manager.CSV.Character.SkillLevel(skillNum, entry.CharacterSkillInitLevel)
		if err != nil {
			glog.Errorf("SkillLevel is nil skillNum:%v,level:%v", skillNum, entry.CharacterSkillInitLevel)
			return
		}
		// 判断是否激活技能
		if skillLevel.CanSkillUnlock(c.Level, c.Star, c.Stage) {
			u.SetCharacterSkillLv(characterId, skillNum, entry.CharacterSkillInitLevel)
		}

	}
}

func (u *User) skillAutoLvUp(characterId int32) {
	c, err := u.CharacterPack.Get(characterId)
	if err != nil {
		return
	}
	skillMap, err := manager.CSV.Character.AllSkill(c.ID)
	if err != nil {
		// log
		return
	}
	for skillNum := range skillMap {
		lv, ok := c.Skills[skillNum]
		if !ok {
			// 未激活 取1级
			lv = entry.CharacterSkillInitLevel
		} else {
			// 已激活 取下一级
			lv++
		}
		skillLevel, err := manager.CSV.Character.SkillLevel(skillNum, lv)
		if err != nil {
			continue
		}
		ok = skillLevel.CanSkillAutoLvUp(c.GetStar())
		if ok {
			u.SetCharacterSkillLv(characterId, skillNum, lv)
		}

	}
}

func (u *User) SetCharacterSkillLv(characterId, skillNum, lv int32) {
	character, err := u.CharacterPack.Get(characterId)
	if err != nil {
		return
	}
	oldLv, ok := character.Skills[skillNum]
	if !ok {
		oldLv = 0
	}
	character.Skills[skillNum] = lv

	skillCfg, _ := manager.CSV.Character.GetSkillByNum(characterId, skillNum)
	u.TriggerQuestUpdate(static.TaskTypeCharacterSkillLevelCount, character.GetRare(), skillCfg.SkillType, lv, oldLv)
	u.TriggerQuestUpdate(static.TaskTypeCharacterSkillLevel, characterId, skillNum, lv)
}

func (u *User) CharacterWear(characterID int32, equipmentUIDs, worldItemUIDs []int64) ([]*common.Equipment, []*common.WorldItem, error) {
	err := u.CheckActionUnlock(static.ActionIdTypeCharacterequipup)
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}

	character, err := u.CharacterPack.Get(characterID)
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}

	var equipments []*common.Equipment
	var worldItems []*common.WorldItem

	if len(equipmentUIDs) > 0 {
		// 最多只有四件装备
		if len(equipmentUIDs) > CharacterEquipmentCounts {
			return nil, nil, common.ErrCharacterWearCountInvalid
		}

		equipments, err = u.EquipmentPack.BatchGet(equipmentUIDs)
		if err != nil {
			return nil, nil, errors.WrapTrace(err)
		}

		mPart := map[int64]int32{}
		mPartCheck := map[int32]bool{}

		// 检查穿戴装备是否有重复部位
		for _, equipment := range equipments {
			part, err := manager.CSV.Equipment.Part(equipment.EID)
			if err != nil {
				return nil, nil, errors.WrapTrace(err)
			}

			if mPartCheck[part] {
				return nil, nil, common.ErrCharacterWearRepeatPart
			}

			mPart[equipment.ID] = part
			mPartCheck[part] = true
		}

		for _, equipment := range equipments {
			// 检查职业
			career, err := manager.CSV.Character.Career(character.ID)
			if err != nil {
				return nil, nil, errors.WrapTrace(err)
			}

			err = manager.CSV.Equipment.CheckCareer(equipment.EID, career)
			if err != nil {
				return nil, nil, errors.WrapTrace(err)
			}

			equipment.CID = characterID
			oldEquipmentUID := character.EquipEquipment(equipment.ID, mPart[equipment.ID])

			u.BIEquipmentOp(equipment, bilog.EquipmentOpWear, logreason.EmptyReason())

			// 清空替换下来装备的CID
			if oldEquipmentUID != 0 {
				oldEquipment, err := u.EquipmentPack.Get(oldEquipmentUID)
				if err != nil {
					return nil, nil, errors.WrapTrace(err)
				}

				oldEquipment.CID = 0
				equipments = append(equipments, oldEquipment)
				u.BIEquipmentOp(oldEquipment, bilog.EquipmentOpTakeOff, logreason.EmptyReason())
			}
		}
	}

	if len(worldItemUIDs) > 0 {
		if len(worldItemUIDs) > CharacterWorldItemCounts {
			return nil, nil, common.ErrCharacterWearCountInvalid
		}

		worldItem, err := u.WorldItemPack.Get(worldItemUIDs[0])
		if err != nil {
			return nil, nil, errors.WrapTrace(err)
		}

		// 清除旧数据
		if character.WorldItem != 0 {
			oldWorldItem, err := u.WorldItemPack.Get(character.WorldItem)
			if err != nil {
				return nil, nil, errors.WrapTrace(err)
			}

			oldWorldItem.CID = 0
			worldItems = append(worldItems, oldWorldItem)
			u.BIWorldItemOp(worldItem, bilog.WorldItemOpTakeOff, logreason.EmptyReason())
		}

		// 更新
		worldItem.CID = characterID
		character.WorldItem = worldItem.ID
		worldItems = append(worldItems, worldItem)

		u.BIWorldItemOp(worldItem, bilog.WorldItemOpWear, logreason.EmptyReason())

	}
	charaPower, err := u.CalCharacterCombatPower(character)
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}
	character.Power = charaPower

	return equipments, worldItems, nil
}

func (u *User) CharacterUndress(characterID int32, equipmentUIDs, worldItemUIDs []int64) ([]*common.Equipment, []*common.WorldItem, error) {
	err := u.CheckActionUnlock(static.ActionIdTypeCharacterequipup)
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}

	character, err := u.CharacterPack.Get(characterID)
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}

	var equipments []*common.Equipment
	var worldItems []*common.WorldItem

	if len(equipmentUIDs) > 0 {
		targetEquipments, err := u.EquipmentPack.BatchGet(equipmentUIDs)
		if err != nil {
			return nil, nil, errors.WrapTrace(err)
		}

		for _, equipment := range targetEquipments {
			if character.StripEquipment(equipment.ID) {
				// 确实脱下来是这个角色的才能重置
				equipment.CID = 0
				equipments = append(equipments, equipment)
			}
		}
	}

	if len(worldItemUIDs) > 0 {
		worldItem, err := u.WorldItemPack.Get(worldItemUIDs[len(worldItemUIDs)-1])
		if err != nil {
			return nil, nil, errors.WrapTrace(err)
		}

		if worldItem.ID == character.WorldItem {
			worldItem.CID = 0
			character.WorldItem = 0
			worldItems = append(worldItems, worldItem)
		}
	}
	charaPower, err := u.CalCharacterCombatPower(character)
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}
	character.Power = charaPower
	return equipments, worldItems, nil
}
