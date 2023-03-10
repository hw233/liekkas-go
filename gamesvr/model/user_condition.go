package model

import (
	"shared/common"
	"shared/csv/static"
	"shared/utility/errors"
)

func (u *User) CheckUserCondition(condition *common.Condition) error {
	switch condition.ConditionType {
	case static.ConditionTypeUserLevel:
		if u.Info.GetLevel() < condition.Params[0] {
			return common.ErrUserLevelNotArrival
		}
	case static.ConditionTypePassLevel:
		levelId := condition.Params[0]
		if !u.LevelsInfo.IsLevelPassed(levelId) {
			return errors.Swrapf(common.ErrLevelNotPassed, levelId)
		}

	case static.ConditionTypeFinishQuest:
		questId := condition.Params[0]
		if !u.QuestPack.IsQuestCompleted(questId) {
			return errors.Swrapf(common.ErrQuestNotComplete, questId)
		}

	case static.ConditionTypeExplorePoint:
		err := u.checkExplorePointCondition(condition)
		if err != nil {
			return err
		}
	case static.ConditionTypeHeroLevel:
		level := condition.Params[0]
		if level > u.HeroPack.GetLevel() {
			return errors.Swrapf(common.ErrHeroLevelNotArrival)
		}
	case static.ConditionTypeHeroSkillLevel:
		heroId := condition.Params[0]
		skillId := condition.Params[1]
		skillLevel := condition.Params[2]

		hero, ok := u.HeroPack.GetHero(heroId)
		if !ok {
			return errors.Swrapf(common.ErrHeroNotFound, heroId)
		}

		if skillLevel > hero.GetSkillLevel(skillId) {
			return errors.Swrapf(common.ErrHeroSkillLevelNotArrival, heroId, skillId)
		}

	case static.ConditionTypeHeroSkillItemUsed:
		heroId := condition.Params[0]
		itemUse := condition.Params[1]

		hero, ok := u.HeroPack.GetHero(heroId)
		if !ok {
			return errors.Swrapf(common.ErrHeroNotFound, heroId)
		}

		if hero.GetSkillItemUsed() < itemUse {
			return errors.Swrapf(common.ErrHeroSkillItemUsedNotArrival, heroId)
		}

	case static.ConditionTypeCharaLevel:
		charaId := condition.Params[0]
		level := condition.Params[1]
		chara, err := u.CharacterPack.Get(charaId)
		if err != nil {
			return err
		}

		if level > chara.GetLevel() {
			return errors.Swrapf(common.ErrCharacterLevelNotArrival, charaId)
		}

	case static.ConditionTypeCharaStar:
		charaId := condition.Params[0]
		star := condition.Params[1]
		chara, err := u.CharacterPack.Get(charaId)
		if err != nil {
			return err
		}

		if star > chara.GetStar() {
			return errors.Swrapf(common.ErrCharacterStarNotArrival, charaId)
		}
	case static.ConditionTypeGuide:
		guideId := condition.Params[0]
		if !u.Info.PassedGuideIds.Contains(guideId) {
			return errors.Swrapf(common.ErrGuideNotPass, guideId)

		}

	case static.ConditionTypeYggAreaPrestige:
		areaId := condition.Params[0]
		prestige := condition.Params[1]
		return u.CheckPrestigeEnough(areaId, prestige)
	case static.ConditionTypeGraveyardBuildLv:
		buildId := condition.Params[0]
		lv := condition.Params[1]
		return u.Graveyard.CheckBuildLv(buildId, lv)
	case static.ConditionTypeYggCompleteTask:
		taskId := condition.Params[0]
		if !u.Yggdrasil.Task.CompleteTaskIds.Contains(taskId) {
			return errors.WrapTrace(common.ErrParamError)
		}

	case static.ConditionTypeYggMoveCount:
		moveCount := condition.Params[0]
		return u.Yggdrasil.SpecialStatics.IsMoveCountReach(moveCount)
	case static.ConditionTypeHasItem:
		itemId := condition.Params[0]
		count := condition.Params[1]
		if !u.ItemPack.Enough(itemId, count) {
			return errors.Swrapf(common.ErrItemNotEnough, itemId, count)
		}
	case static.ConditionTypeApNoMoreThan:
		ap := condition.Params[0]
		if u.Yggdrasil.TravelInfo.TravelAp > ap {
			return errors.Swrapf(common.ErrYggdrasilApNoMoreThan, ap, u.Yggdrasil.TravelInfo.TravelAp)
		}
	}

	return nil
}

func (u *User) CheckUserConditions(conditions *common.Conditions) error {
	if conditions == nil {
		return nil
	}

	for _, condition := range *conditions {
		err := u.CheckUserCondition(&condition)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *User) checkExplorePointCondition(condition *common.Condition) error {
	mapObjType := condition.Params[0]
	switch mapObjType {
	case static.ExploreMapObjectTypeNormalLevel:
		levelId := condition.Params[1]
		if !u.LevelsInfo.IsLevelPassed(levelId) {
			return errors.Swrapf(common.ErrLevelNotPassed, levelId)
		}
	case static.ExploreMapObjectTypeFog:
		fogId := condition.Params[1]
		if !u.ExploreInfo.IsFogUnlocked(fogId) {
			return errors.Swrapf(common.ErrExploreFogLocked, fogId)
		}

	case static.ExploreMapObjectTypeResource:
		resourceId := condition.Params[1]
		resourcePoint, ok := u.ExploreInfo.GetResourcePoint(resourceId)
		if !ok {
			return errors.Swrapf(common.ErrExploreResourceNotCollected, resourceId)
		}

		if resourcePoint.GetCollectTimes() <= 0 {
			return errors.Swrapf(common.ErrExploreResourceNotCollected, resourceId)
		}

	case static.ExploreMapObjectTypeEvent:
		eventId := condition.Params[1]
		_, ok := u.ExploreInfo.GetEventPoint(eventId)
		if !ok {
			return errors.Swrapf(common.ErrExploreNotInteracted, eventId)
		}
	}

	return nil
}

func (u *User) CheckUserCompoundConditions(compoundConditions *common.CompoundConditions) error {
	err := u.CheckUserConditions(compoundConditions.And)
	if err != nil {
		return errors.WrapTrace(err)
	}
	for _, or := range compoundConditions.Or {
		err = u.CheckUserConditions(or)
		if err == nil {
			break
		}
	}
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}
