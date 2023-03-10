package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/csv/static"
	"shared/statistic/bilog"
	"shared/statistic/logreason"
	"shared/utility/errors"
)

func (u *User) UnlockHero(heroId int32) (*Hero, error) {
	_, ok := u.HeroPack.GetHero(heroId)
	if ok {
		return nil, errors.Swrapf(common.ErrHeroUnlocked, heroId)
	}

	heroCfg, err := manager.CSV.Hero.GetHero(heroId)
	if err != nil {
		return nil, err
	}

	if !heroCfg.Show {
		return nil, errors.Swrapf(common.ErrHeroNotFound, heroId)
	}

	err = u.CheckRewardsEnough(heroCfg.UnlockCost)
	if err != nil {
		return nil, err
	}

	reason := logreason.NewReason(logreason.HeroUnlock)
	u.CostRewards(heroCfg.UnlockCost, reason)

	hero := u.HeroPack.AddHero(heroId)

	u.BIHeroOp(hero, u.HeroPack.GetLevel(), bilog.HeroOpUnlock)

	return hero, nil
}

func (u *User) HeroLevelUp() error {

	err := u.CheckActionUnlock(static.ActionIdTypeHerogrowthunlock)
	if err != nil {
		return err
	}

	nextLevel := u.HeroPack.GetLevel() + 1

	maxLevel, err := manager.CSV.TeamLevelCache.GetHeroMaxLevel(u.Info.GetLevel())
	if err != nil {
		return err
	}

	if maxLevel < nextLevel {
		return common.ErrHeroLevelMax
	}

	nextLevelCfg, err := manager.CSV.Hero.GetLevel(nextLevel)
	if err != nil {
		return err
	}

	exp := u.HeroPack.GetExp()
	if exp < nextLevelCfg.Exp {
		return common.ErrHeroExpNotEnough
	}

	u.HeroPack.SetLevel(nextLevel)

	u.TriggerQuestUpdate(static.TaskTypeHeroLevel, nextLevel)
	u.QuestCheckUnlock(static.ConditionTypeHeroLevel)

	for _, hero := range u.HeroPack.Heros {
		power, err := u.CalHeroCombatPower(hero)
		if err != nil {
			return errors.WrapTrace(err)
		}
		hero.Power = power
	}

	u.BIHeroOp(NewHero(0), nextLevel, bilog.HeroOpLevelUp)

	return nil
}

func (u *User) HeroSkillLevelUpgrade(heroId, skillId int32) error {

	err := u.CheckActionUnlock(static.ActionIdTypeHerogrowthunlock)
	if err != nil {
		return err
	}

	hero, ok := u.HeroPack.GetHero(heroId)
	if !ok {
		return errors.Swrapf(common.ErrHeroNotFound, heroId)
	}

	skillLevel := hero.GetSkillLevel(skillId)

	nextSkillLevel := skillLevel + 1
	skillCfg, err := manager.CSV.Hero.GetSkill(heroId, skillId, nextSkillLevel)
	if err != nil {
		return err
	}

	cost := skillCfg.CostItems
	err = u.CheckRewardsEnough(cost)
	if err != nil {
		return err
	}

	reason := logreason.NewReason(logreason.HeroSkillLevelUp)
	u.CostRewards(cost, reason)

	mergedCost := cost.MergeValue()
	heroCfg, _ := manager.CSV.Hero.GetHero(heroId)
	for _, reward := range mergedCost {
		if reward.ID == heroCfg.SkillItem {
			hero.RecordSkillItemUsed(reward.Num)
		}
	}

	hero.SetSkillLevel(skillId, nextSkillLevel)

	u.QuestCheckUnlock(static.ConditionTypeHeroSkillLevel)

	if skillLevel == 0 {
		u.TriggerQuestUpdate(static.TaskTypeHeroSkillUnlockCount, heroId)
	}

	isLevelMax, _ := manager.CSV.Hero.IsSkillLevelMax(heroId, skillId, nextSkillLevel)
	if isLevelMax {
		u.TriggerQuestUpdate(static.TaskTypeHeroSkillLevelMaxCount, heroId)
	}

	power, err := u.CalHeroCombatPower(hero)
	if err != nil {
		return errors.WrapTrace(err)
	}
	hero.Power = power

	u.BIHeroOp(hero, u.HeroPack.GetLevel(), bilog.HeroOpSkillLevelUp)

	return nil
}

func (u *User) AddHeroAttendant(heroId, slot, charaId int32) (int32, int32, error) {

	err := u.CheckActionUnlock(static.ActionIdTypeHeroattendantunlock)
	if err != nil {
		return 0, 0, err
	}

	hero, ok := u.HeroPack.GetHero(heroId)
	if !ok {
		return 0, 0, errors.Swrapf(common.ErrHeroNotFound, heroId)
	}

	conditions, err := manager.CSV.Hero.GetAttendantUnlockConditions(heroId, slot)
	if err != nil {
		return 0, 0, err
	}

	err = u.CheckUserConditions(conditions)
	if err != nil {
		return 0, 0, err
	}

	chara, err := u.CharacterPack.Get(charaId)
	if err != nil {
		return 0, 0, err
	}

	charaHeroId, charaHeroSlot := chara.GetHero()
	if charaHeroId > 0 {
		oldHero, _ := u.HeroPack.GetHero(charaHeroId)
		u.removeHeroAttendant(oldHero, charaHeroSlot)
	}
	chara.SetHero(heroId, slot)

	oldCharaId := hero.GetAttendantCharaId(slot)
	if oldCharaId > 0 {
		u.removeHeroAttendant(hero, slot)
	}
	hero.SetAttendant(slot, charaId)

	power, err := u.CalHeroCombatPower(hero)
	if err != nil {
		return 0, 0, errors.WrapTrace(err)
	}
	hero.Power = power

	u.BIHeroOp(hero, u.HeroPack.GetLevel(), bilog.HeroOpAddChara)

	return oldCharaId, charaHeroId, nil
}

func (u *User) RemoveHeroAttendant(heroId, slot int32) (int32, error) {

	err := u.CheckActionUnlock(static.ActionIdTypeHerogrowthunlock)
	if err != nil {
		return 0, err
	}

	hero, ok := u.HeroPack.GetHero(heroId)
	if !ok {
		return 0, errors.Swrapf(common.ErrHeroNotFound, heroId)
	}

	removeId := u.removeHeroAttendant(hero, slot)

	power, err := u.CalHeroCombatPower(hero)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}
	hero.Power = power

	u.BIHeroOp(hero, u.HeroPack.GetLevel(), bilog.HeroOpAddChara)

	return removeId, nil
}

func (u *User) UpdateHeroLastCalcAttr(heroId int32, heroAttr int32, attrs map[int32]int32) error {

	err := u.CheckActionUnlock(static.ActionIdTypeHerogrowthunlock)
	if err != nil {
		return err
	}

	hero, ok := u.HeroPack.GetHero(heroId)
	if !ok {
		return errors.Swrapf(common.ErrHeroNotFound, heroId)
	}

	hero.UpdateLastAttr(heroAttr, attrs)

	power, err := u.CalHeroCombatPower(hero)
	if err != nil {
		return errors.WrapTrace(err)
	}
	hero.Power = power

	return nil
}

func (u *User) SetHeroNewOneFlag(newOne bool) error {
	err := u.CheckActionUnlock(static.ActionIdTypeHerogrowthunlock)
	if err != nil {
		return err
	}

	u.HeroPack.SetNewOneFlag(newOne)

	return nil
}

func (u *User) removeHeroAttendant(hero *Hero, slot int32) int32 {

	charaId := hero.GetAttendantCharaId(slot)
	if charaId <= 0 {
		return 0
	}

	chara, _ := u.CharacterPack.Get(charaId)

	chara.SetHero(0, 0)
	hero.RemoveAttendant(slot)

	return charaId
}
