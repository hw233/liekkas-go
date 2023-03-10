package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/csv/static"
	"shared/statistic/bilog"
	"shared/statistic/logreason"
	"shared/utility/errors"
	"shared/utility/glog"
	"time"
)

type QuestCache struct {
	ProgressingQuests map[int32]map[int32]bool
	Inited            bool
}

func NewQuestCache() *QuestCache {
	return &QuestCache{
		ProgressingQuests: map[int32]map[int32]bool{},
		Inited:            false,
	}
}

func (u *User) InitQuest() {
	u.QuestCache = NewQuestCache()

	achievementFinishCount := int32(0)

	for questId, quest := range u.QuestPack.GetAllQuest() {
		if quest.IsCounting() {
			u.registerProgressingQuest(questId)
		} else if quest.IsRecieved() {
			questCfg, err := manager.CSV.Quest.Get(questId)
			if err != nil {
				continue
			}

			if questCfg.Module == static.TaskModuleAchievement {
				achievementFinishCount = achievementFinishCount + 1
			}
		}
	}

	u.Info.SetAchievementFinishCount(achievementFinishCount)

	u.QuestCache.Inited = true
}

func (u *User) TryAcceptNewQuests() {
	for questId, questCfg := range manager.CSV.Quest.GetAll() {
		if isAutoAccpetQuest(questCfg.Module) {
			u.TryAcceptQuest(questId)
		}
	}
}

func (u *User) TryAcceptQuest(questId int32) {
	_, ok := u.QuestPack.GetQuest(questId)
	if ok {
		return
	}

	if u.checkQuestUnlockCondition(questId) {
		u.AcceptQuest(questId)
	}
}

func (u *User) TryAcceptQuestByGroup(groupId int32) {
	questIds := manager.CSV.Quest.GetQuestIdsByGroup(groupId)
	for _, questId := range questIds {
		u.TryAcceptQuest(questId)
	}
}

func (u *User) AcceptQuest(questId int32) *Quest {
	quest, ok := u.QuestPack.GetQuest(questId)
	if !ok {
		quest = u.QuestPack.AddQuest(questId)
	}
	quest.Restart()

	u.registerProgressingQuest(questId)
	u.questInitProgress(quest)

	u.AddQuestUpdateNotify([]int32{questId})

	u.BIQuest(quest, u.CreateBICommon(bilog.EventNameQuest))

	return quest
}

func (u *User) RemoveQuest(questId int32) {
	_, ok := u.QuestPack.GetQuest(questId)
	if !ok {
		return
	}

	u.unregisterProgressQuest(questId)
	u.QuestPack.RemoveQuest(questId)
}

func (u *User) RemoveQuestByGroup(groupId int32) {
	questIds := manager.CSV.Quest.GetQuestIdsByGroup(groupId)
	for _, questId := range questIds {
		u.RemoveQuest(questId)
	}
}

func (u *User) TriggerQuestUpdate(conditionType int32, params ...interface{}) {
	questIds := u.getProgressingQuestIds(conditionType)
	for questId := range questIds {
		if u.tryUpdateQuestProgress(questId, conditionType, params...) {
			u.AddQuestUpdateNotify([]int32{questId})
		}
	}
}

func (u *User) RecieveQuest(questId int32) error {
	quest, ok := u.QuestPack.GetQuest(questId)
	if !ok {
		return errors.Swrapf(common.ErrQuestNotFound, questId)
	}

	if !quest.IsCompleted() {
		return errors.Swrapf(common.ErrQuestProgressNotArrival, questId)

	}

	questCfg, err := manager.CSV.Quest.Get(questId)
	if err != nil {
		return err

	}

	quest.Recieve()

	u.unregisterProgressQuest(quest.Id)

	biCommon := u.CreateBICommon(bilog.EventNameQuest)

	if questCfg.DropId > 0 {
		reson := logreason.NewReason(logreason.QuestComplete, logreason.AddRelateLog(biCommon.GetLogId()))
		_, err = u.AddRewardsByDropId(questCfg.DropId, reson)
		if err != nil {
			return err
		}
	}

	u.BIQuest(quest, biCommon)

	questActivityData := u.QuestPack.GetActivityData(QuestModule(questCfg.Module))
	questActivityData.AddPoint(questCfg.AddPoint)

	u.Info.AddAchievementFinishCount(1)

	return nil
}

func (u *User) RecieveQuestActivityReward(rewardId int32) error {
	questActivityCfg, err := manager.CSV.Quest.GetQuestActivity(rewardId)
	if err != nil {
		return err
	}

	module := QuestModule(questActivityCfg.Module)
	questActivityData := u.QuestPack.GetActivityData(module)

	if questActivityData.Point < questActivityCfg.Condition {
		errors.Swrapf(common.ErrQuestActivityProgressNotArrival, rewardId)
	}

	if questActivityData.IsRewardReceived(rewardId) {
		return errors.Swrapf(common.ErrQuestActivityRewardGot, rewardId)
	}

	reason := logreason.NewReason(logreason.QuestActivityReward)
	_, err = u.AddRewardsByDropId(questActivityCfg.DropId, reason)
	if err != nil {
		return err
	}

	questActivityData.RecordReward(rewardId)

	return nil
}

func (u *User) QuestDailyRefresh(refreshTime int64) {
	u.RefreshQuests(static.TaskModuleDaily)

	lastWeekRefreshTime := WeekRefreshTime(time.Unix(u.QuestPack.LastDailyRefreshTime, 0))
	weekRefreshTime := WeekRefreshTime(time.Unix(refreshTime, 0))
	if weekRefreshTime.After(lastWeekRefreshTime) {
		u.RefreshQuests(static.TaskModuleWeekly)
	}

	u.TriggerQuestUpdate(static.TaskTypeDailyLogin)
}

func (u *User) RefreshQuests(module QuestModule) {
	questIds := manager.CSV.Quest.GetQuestIdsByModule(int32(module))

	for _, questId := range questIds {
		_, ok := u.QuestPack.GetQuest(questId)
		if ok {
			u.RemoveQuest(questId)
		}

		u.TryAcceptQuest(questId)
	}

	questActivityData := u.QuestPack.GetActivityData(module)
	questActivityData.Reset()
}

func (u *User) QuestCheckUnlock(conditionType int32, params ...interface{}) {
	var err error

	switch conditionType {
	case static.ConditionTypeUserLevel:
		paramsObj := struct {
			OldLevel int32
			NewLevel int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		for level := paramsObj.OldLevel + 1; level <= paramsObj.NewLevel; level++ {
			questIds := manager.CSV.Quest.GetPassLevelUnlocks(level)
			for _, questId := range questIds {
				u.TryAcceptQuest(questId)
			}
		}
	case static.ConditionTypePassLevel:
		paramsObj := struct {
			LevelId int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		questIds := manager.CSV.Quest.GetPassLevelUnlocks(paramsObj.LevelId)
		for _, questId := range questIds {
			u.TryAcceptQuest(questId)
		}

	case static.ConditionTypeFinishQuest:
		paramsObj := struct {
			QuestId int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		questCfg, err := manager.CSV.Quest.Get(paramsObj.QuestId)
		if err != nil {
			break
		}

		for i := 0; i < len(questCfg.NextQuests); i++ {
			nextQuestId := questCfg.NextQuests[i]
			u.TryAcceptQuest(nextQuestId)
		}

	default:
		questIds := manager.CSV.Quest.GetCommonConditionUnlocks(conditionType)
		for _, questId := range questIds {
			u.TryAcceptQuest(questId)
		}
	}

	if err != nil {
		glog.Errorf("check quest unlock failure, err: %+v\n", err)
		return
	}

}

func (u *User) tryUpdateQuestProgress(questId, conditionType int32, params ...interface{}) bool {
	quest, ok := u.QuestPack.GetQuest(questId)
	if !ok {
		return false
	}

	if !quest.IsCounting() {
		glog.Errorf("quest [%d] update error: quest not couting\n", questId)
		return false
	}

	questCfg, err := manager.CSV.Quest.Get(questId)
	if err != nil {
		//todo error log
		glog.Errorf("quest update error: %v\n", err)
		return false
	}

	if questCfg.ConditionType != conditionType {
		return false
	}

	progress, err := u.calcNewQuestProgress(questCfg, quest.Progress, params...)
	if err != nil {
		glog.Errorf("quest [%d] update error: %+v", questId, err)
		return false
	}

	if progress == quest.GetProgress() {
		return false
	}

	u.updateQuestProgress(quest, progress)

	return true
}

func (u *User) updateQuestProgress(quest *Quest, newProgress int32) {
	quest.UpdateProgress(newProgress)

	questId := quest.GetId()

	questCfg, _ := manager.CSV.Quest.Get(questId)
	if quest.GetProgress() >= questCfg.TargetCount {
		quest.Complete()
		u.onQuestCompleted(questId)
	}
}

func (u *User) registerProgressingQuest(questId int32) {
	questCfg, err := manager.CSV.Quest.Get(questId)
	if err != nil {
		glog.Error(err.Error())
		return
	}

	conditionType := questCfg.ConditionType

	questSet, ok := u.QuestCache.ProgressingQuests[conditionType]
	if !ok {
		questSet = make(map[int32]bool)
		u.QuestCache.ProgressingQuests[conditionType] = questSet
	}

	questSet[questId] = true
}

func (u *User) unregisterProgressQuest(questId int32) {
	questCfg, err := manager.CSV.Quest.Get(questId)
	if err != nil {
		glog.Error(err.Error())
		return
	}

	conditionType := questCfg.ConditionType

	questSet, ok := u.QuestCache.ProgressingQuests[conditionType]
	if !ok {
		return
	}

	delete(questSet, questId)
}

func (u *User) getProgressingQuestIds(conditionType int32) map[int32]bool {
	questIds, ok := u.QuestCache.ProgressingQuests[conditionType]
	if !ok {
		return make(map[int32]bool)
	}

	return questIds
}

func (u *User) onQuestCompleted(questId int32) {
	quest, _ := u.QuestPack.GetQuest(questId)
	u.BIQuest(quest, u.CreateBICommon(bilog.EventNameQuest))

	u.QuestCheckUnlock(static.ConditionTypeFinishQuest, questId)

	if manager.CSV.Quest.IsAutoRecieve(questId) {
		u.RecieveQuest(questId)
	}
}

func (u *User) questInitProgress(quest *Quest) {
	questCfg, err := manager.CSV.Quest.Get(quest.Id)
	if err != nil {
		glog.Errorf("quest [%d] init failure, err: %+v", quest.Id, err)
		return
	}

	var progress int32 = 0

	if len(questCfg.PreQuests) > 0 {
		for _, preQuestId := range questCfg.PreQuests {
			preQuestCfg, err := manager.CSV.Quest.Get(preQuestId)
			if err != nil {
				glog.Errorf("quest [%d] init failure, err: %+v", quest.Id, err)
				return
			}

			if questConditionCheckSame(questCfg.ConditionParams, preQuestCfg.ConditionParams) {
				preQuest, ok := u.QuestPack.GetQuest(preQuestId)
				if !ok {
					glog.Errorf("quest [%d] init failure, pre quest [%d] not found", quest.Id, preQuestId)
					return
				}

				if preQuest.Progress > progress {
					progress = preQuest.Progress
				}
			}
		}
	} else {
		switch questCfg.ConditionType {
		case static.TaskTypeAccountLevel:
			progress, err = u.calcNewQuestProgress(questCfg, progress, u.Info.GetLevel())

		case static.TaskTypeCharacterLevelCount:
			for _, chara := range *u.CharacterPack {
				progress, err = u.calcNewQuestProgress(questCfg, progress, chara, int32(0))
				if err != nil {
					break
				}
			}

		case static.TaskTypeCharacterSkillLevelCount:
			for charaId, chara := range *u.CharacterPack {
				for skillNum, skillLevel := range chara.Skills {
					skillCfg, _ := manager.CSV.Character.GetSkillByNum(charaId, skillNum)
					progress, err = u.calcNewQuestProgress(questCfg, progress, chara.GetRare(),
						skillCfg.SkillType, skillLevel, int32(0))
					if err != nil {
						break
					}
				}
			}

		case static.TaskTypeCharacterStageCount:
			for _, chara := range *u.CharacterPack {
				progress, err = u.calcNewQuestProgress(questCfg, progress, chara.GetRare(), chara.GetStage(), int32(0))
				if err != nil {
					break
				}
			}

		case static.TaskTypeCharacterCampCount:
			for charaId, chara := range *u.CharacterPack {
				camp, err := manager.CSV.Character.Camp(charaId)
				if err != nil {
					break
				}
				progress, err = u.calcNewQuestProgress(questCfg, progress, chara.GetRare(), camp)
				if err != nil {
					break
				}
			}

		case static.TaskTypeCharacterLevel:
			charaId := questCfg.ConditionParams[0]
			if charaId > 0 {
				chara, getErr := u.CharacterPack.Get(charaId)
				if getErr != nil {
					break
				}
				progress = chara.GetLevel()

			} else {
				for charaId, chara := range *u.CharacterPack {
					progress, err = u.calcNewQuestProgress(questCfg, progress, charaId, chara.GetLevel())
					if err != nil {
						break
					}
				}
			}

		case static.TaskTypeCharacterSkillLevel:
			cfgCharaId := questCfg.ConditionParams[0]
			cfgSkillNum := questCfg.ConditionParams[1]
			if cfgCharaId > 0 {
				chara, getErr := u.CharacterPack.Get(cfgCharaId)
				if getErr != nil {
					break
				}

				for skillNum, skillLevel := range chara.Skills {
					if cfgSkillNum == 0 || cfgSkillNum == skillNum {
						progress, err = u.calcNewQuestProgress(questCfg, progress, cfgCharaId,
							skillNum, skillLevel)
						if err != nil {
							break
						}
					}
				}
			} else {
				for charaId, chara := range *u.CharacterPack {
					for skillNum, skillLevel := range chara.Skills {
						progress, err = u.calcNewQuestProgress(questCfg, progress, charaId,
							skillNum, skillLevel)
						if err != nil {
							break
						}
					}
				}
			}

		case static.TaskTypeCharacterStage:
			charaId := questCfg.ConditionParams[0]
			if charaId > 0 {
				chara, getErr := u.CharacterPack.Get(charaId)
				if getErr != nil {
					break
				}
				progress = chara.GetStage()

			} else {
				for charaId, chara := range *u.CharacterPack {
					progress, err = u.calcNewQuestProgress(questCfg, progress, charaId, chara.GetStage())
					if err != nil {
						break
					}
				}
			}

		case static.TaskTypeCharacterStar:
			charaId := questCfg.ConditionParams[0]
			if charaId > 0 {
				chara, getErr := u.CharacterPack.Get(charaId)
				if getErr != nil {
					break
				}
				progress = chara.GetStar()

			} else {
				for charaId, chara := range *u.CharacterPack {
					progress, err = u.calcNewQuestProgress(questCfg, progress, charaId, chara.GetStar())
					if err != nil {
						break
					}
				}
			}

		case static.TaskTypeHasCharacters:
			for charaId := range *u.CharacterPack {
				progress, err = u.calcNewQuestProgress(questCfg, progress, charaId)
				if err != nil {
					break
				}
			}

		case static.TaskTypeEquipmentLevelCount:
			for _, equip := range u.EquipmentPack.Equipments {
				progress, err = u.calcNewQuestProgress(questCfg, progress, equip.Rarity, 0, equip.Level.Value())
			}

		case static.TaskTypeHeroLevel:
			progress, err = u.calcNewQuestProgress(questCfg, progress, u.HeroPack.GetLevel())

		case static.TaskTypeHeroSkillUnlockCount:
			cfgHeroId := questCfg.ConditionParams[0]
			for heroId, hero := range u.HeroPack.Heros {
				if cfgHeroId == 0 || cfgHeroId == heroId {
					progress = int32(len(hero.Skills))
				}
			}

		case static.TaskTypeHeroSkillLevelMaxCount:
			cfgHeroId := questCfg.ConditionParams[0]
			for heroId, hero := range u.HeroPack.Heros {
				if cfgHeroId == 0 || cfgHeroId == heroId {
					for skillId, skillLevel := range hero.Skills {
						var isMax bool = false
						isMax, err = manager.CSV.Hero.IsSkillLevelMax(heroId, skillId, skillLevel)
						if err != nil {
							break
						}

						if isMax {
							progress = progress + 1
						}
					}
				}
			}

		case static.TaskTypeGraveyarBuildingLevelCount:
			buildId := questCfg.ConditionParams[0]
			level := questCfg.ConditionParams[1]

			for _, build := range u.Graveyard.Buildings {
				if (build.BuildId == buildId || buildId == 0) && build.FetchLevel() >= level {
					progress = progress + 1
				}
			}

		case static.TaskTypeGraveyarBuildingStageCount:
			buildId := questCfg.ConditionParams[0]
			stage := questCfg.ConditionParams[1]

			for _, build := range u.Graveyard.Buildings {
				if (build.BuildId == buildId || buildId == 0) && build.FetchStage() >= stage {
					progress = progress + 1
				}
			}

		case static.TaskTypeWorldLevelCount:
			for _, worldItem := range u.WorldItemPack.WorldItems {
				progress, err = u.calcNewQuestProgress(questCfg, progress, worldItem.Rarity, 0, worldItem.Level.Value())
				if err != nil {
					break
				}
			}

		case static.TaskTypeWorldStarCount:
			for _, worldItem := range u.WorldItemPack.WorldItems {
				progress, err = u.calcNewQuestProgress(questCfg, progress, worldItem.Rarity, 0, worldItem.Stage.Value())
				if err != nil {
					break
				}
			}

		case static.TaskTypeTowerStagePassed:
			towerId := questCfg.ConditionParams[0]

			if towerId == 0 {
				for _, tower := range u.TowerInfo.Towers {
					stage := tower.GetCurStage()
					if stage > progress {
						progress = stage
					}
				}
			} else {
				tower, ok := u.TowerInfo.GetTower(towerId)
				if !ok {
					break
				}

				progress = tower.GetCurStage()
			}

		case static.TaskTypeChapterScore:
			chapterType := questCfg.ConditionParams[0]

			if chapterType == 0 {
				for _, chapter := range u.ChapterInfo.Chapters {
					progress = progress + chapter.GetScore()
				}
			} else {
				chapterIds := manager.CSV.ChapterEntry.GetChapterIdsByType(chapterType)
				for _, chapterId := range chapterIds {
					chapter, ok := u.ChapterInfo.GetChapter(chapterId)
					if ok {
						progress = progress + chapter.GetScore()
					}
				}
			}

		case static.TaskTypeManualCollect:
			for manualId, _ := range *u.ManualInfo {
				manualVersion, cfgErr := manager.CSV.Manual.GetManualVersion(manualId)
				if err != nil {
					glog.Error(cfgErr.Error())
					continue
				}

				if questConditionCheckConfigContains(questCfg.ConditionParams, []int32{manualVersion}) {
					progress = progress + 1
				}
			}

		}

		if err != nil {
			glog.Errorf("quest [%d] init failure, err: %+v", quest.Id, err)
			return
		}
	}

	u.updateQuestProgress(quest, progress)
}

func isAutoAccpetQuest(questModule int32) bool {
	switch questModule {
	case static.TaskModuleDaily,
		static.TaskModuleWeekly,
		static.TaskModuleAchievement,
		static.TaskModuleHero:
		return true
	default:
		return false
	}
}
