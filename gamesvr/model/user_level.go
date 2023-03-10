package model

import (
	"context"
	"gamesvr/manager"
	"math"
	"shared/common"
	"shared/csv/entry"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/statistic/bilog"
	"shared/statistic/logreason"
	"shared/utility/errors"
	"shared/utility/glog"
	"shared/utility/rand"
	"shared/utility/slice"
)

type BattleVerifyData struct {
	IsWin                bool
	BattleResult         int32
	BattleInput          string
	BattleOutput         string
	CompleteTargets      []int32
	CompleteAchievements []int32
	Statistic            *pb.VOBattleStatis
	BattleCharacters     []*pb.VOBattleCharacter // 战斗同步血量
}

type LevelRewardResult struct {
	FirstPassResult   *pb.VOResourceResult
	PassResult        *pb.VOResourceResult
	TargetResult      *pb.VOResourceResult
	AchievementResult *pb.VOResourceResult
	TotalRewards      *common.Rewards
	Cost              *common.Rewards
}

type LevelResult struct {
	RewardRewsult *LevelRewardResult
	Charas        []int32
	ChapterId     int32
}

func NewLevelResult() *LevelResult {
	return &LevelResult{
		RewardRewsult: &LevelRewardResult{},
		Charas:        []int32{},
		ChapterId:     0,
	}
}

type BattleFormationPosition struct {
	Position int32
	CharaId  int32
}

func NewBattleFormationPosition(position, charaId int32) *BattleFormationPosition {
	return &BattleFormationPosition{
		Position: position,
		CharaId:  charaId,
	}
}

func (u *User) StartLevel(levelId int32, formation *pb.VOBattleFormation,
	systemType int32, systemParams []int64) (int32, *pb.VOBattleFighterDetail, error) {
	err := u.checkLevelStart(levelId, formation, systemType, systemParams)
	if err != nil {
		return 0, nil, err
	}

	seed := int32(0)

	levelCfg, _ := manager.CSV.LevelsEntry.GetLevel(levelId)

	var detail *pb.VOBattleFighterDetail = nil

	if levelCfg.BattleID > 0 {
		seed = rand.RangeInt32(1, math.MaxInt32)
		detail, err = u.genBattleDetailData(formation)
		if err != nil {
			return seed, detail, err
		}
	}

	u.LevelsInfo.SetCurLevel(levelId, formation, systemType, systemParams, seed)

	return seed, detail, nil
}

func (u *User) LevelEnd(ctx context.Context, levelId int32, battleData *BattleVerifyData) (*LevelResult, error) {
	curLevelInfo := u.LevelsInfo.GetCurLevel()
	systemType := curLevelInfo.SystemType

	switch systemType {
	case static.BattleTypeYgg, static.BattleTypeChallengeAltar:

		err := u.Yggdrasil.TravelInfo.ContainCharacters(battleData.BattleCharacters)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
	}

	curLevelId := u.LevelsInfo.GetCurLevelId()

	levelCfg, err := manager.CSV.LevelsEntry.GetLevel(levelId)
	if err != nil {
		return nil, err
	}

	if levelCfg.BattleID > 0 {
		if curLevelId == 0 {
			return nil, errors.Swrapf(common.ErrLevelNotStarted, levelId)
		}

		if curLevelId != levelId {
			return nil, errors.Swrapf(common.ErrOtherLevelStarted, curLevelId)
		}

		// todo check battle result
		for _, index := range battleData.CompleteTargets {
			if index <= 0 || index > 3 {
				return nil, errors.Swrapf(common.ErrLevelInvalidTarget, levelId)
			}
		}

		for _, index := range battleData.CompleteAchievements {
			if index <= 0 || int(index) > len(levelCfg.AchievementsDrops) {
				return nil, errors.Swrapf(common.ErrLevelInvalidAchievement, levelId)
			}
		}
	}

	biCommon := u.CreateBICommon(bilog.EventNameLevel)

	var levelResult *LevelResult
	if battleData.BattleResult > 0 {
		switch battleData.BattleResult {
		case static.BattleEndTypeWin:
			passReason := logreason.NewReason(logreason.PassLevel, logreason.AddRelateLog(biCommon.GetLogId()))
			levelResult, err = u.onLevelWin(ctx, levelId, battleData, passReason)
			if err != nil {
				return nil, err
			}
		case static.BattleEndTypeLose,
			static.BattleEndTypeGiveUp:
			err := u.onLevelFailed(ctx, u.LevelsInfo.GetCurLevel(), battleData.BattleCharacters, battleData.BattleResult)
			if err != nil {
				return nil, errors.WrapTrace(err)
			}
			levelResult = NewLevelResult()
		default:
			return nil, errors.Swrapf(common.ErrBattleInvalidEndType, battleData.BattleResult)
		}

	} else {
		//remove in the feature
		if battleData.IsWin {
			passReason := logreason.NewReason(logreason.PassLevel, logreason.AddRelateLog(biCommon.GetLogId()))
			levelResult, err = u.onLevelWin(ctx, levelId, battleData, passReason)
			if err != nil {
				return nil, err
			}
		} else {
			err := u.onLevelFailed(ctx, u.LevelsInfo.GetCurLevel(), battleData.BattleCharacters, static.BattleEndTypeLose)
			if err != nil {
				return nil, errors.WrapTrace(err)
			}
			levelResult = NewLevelResult()
		}

	}

	score := int32(len(battleData.CompleteTargets))
	rewardResult := levelResult.RewardRewsult
	isFirst := rewardResult.FirstPassResult != nil
	u.BILevels(levelId, systemType, battleData.IsWin, isFirst, score,
		curLevelInfo.Formation, rewardResult.Cost, rewardResult.TotalRewards, battleData, biCommon)

	u.TriggerQuestUpdate(static.TaskTypeLevelEndByType, systemType, int32(1))

	heroId := u.LevelsInfo.GetCurLevelHeroId()
	if heroId > 0 {
		u.TriggerQuestUpdate(static.TaskTypeLevelEndHero, heroId, int32(1))
	}

	u.LevelsInfo.ClearCurLevel()

	return levelResult, nil
}

func (u *User) onLevelWin(ctx context.Context, levelId int32, battleData *BattleVerifyData,
	passReason *logreason.Reason) (*LevelResult, error) {
	LevelResult := NewLevelResult()

	levelCfg, _ := manager.CSV.LevelsEntry.GetLevel(levelId)

	// commonData := u.CreateUserCommonDIData("stage_flow", servertime.Now())

	reason := logreason.NewReason(logreason.PassLevel)
	u.CostRewards(levelCfg.Cost, reason)

	curLevelInfo := u.LevelsInfo.GetCurLevel()

	charaIds := []int32{}
	if curLevelInfo.Formation != nil {
		for _, formationChara := range curLevelInfo.Formation.BattleCharacters {
			charaIds = append(charaIds, formationChara.CharacterId)
		}

		statistic := battleData.Statistic
		if statistic != nil {
			heroStatis := statistic.Hero
			if heroStatis != nil && heroStatis.HeroId == curLevelInfo.Formation.BattleHero.HeroId {
				u.TriggerQuestUpdate(static.TaskTypeBattleHeroDamage, heroStatis.HeroId, heroStatis.Damage)
				u.TriggerQuestUpdate(static.TaskTypeBattleHeroCure, heroStatis.HeroId, heroStatis.Cure)
				u.TriggerQuestUpdate(static.TaskTypeBattleHeroHurt, heroStatis.HeroId, heroStatis.Hurt)

				for _, castRecord := range heroStatis.SkillCasts {
					if castRecord.CastTimes > 0 {
						u.TriggerQuestUpdate(static.TaskTypeBattleHeroSkillCast,
							heroStatis.HeroId, castRecord.SkillId, castRecord.CastTimes)
					}
				}
			}

			for _, charaStatis := range statistic.Charas {
				if !slice.SliceInt32HasEle(charaIds, charaStatis.CharaId) {
					continue
				}

				u.TriggerQuestUpdate(static.TaskTypeBattleCharaDamage, charaStatis.CharaId, charaStatis.Damage)
				u.TriggerQuestUpdate(static.TaskTypeBattleCharaCure, charaStatis.CharaId, charaStatis.Cure)
				u.TriggerQuestUpdate(static.TaskTypeBattleCharaHurt, charaStatis.CharaId, charaStatis.Hurt)

				for _, castRecord := range charaStatis.SkillRecord {
					if castRecord.CastTimes > 0 {
						u.TriggerQuestUpdate(static.TaskTypeBattleCharaSkillCast,
							charaStatis.CharaId, castRecord.SkillNum, castRecord.CastTimes)
					}
				}
			}
		}
	}

	u.TriggerQuestUpdate(static.TaskTypeLevelCharacter, levelId, charaIds)

	levelRewardResult, err := u.PassLevel(ctx, levelCfg, battleData.CompleteTargets, battleData.CompleteAchievements,
		curLevelInfo.SystemType, curLevelInfo.SystemParams, charaIds, battleData.BattleCharacters, passReason)
	if err != nil {
		return nil, err
	}

	LevelResult.RewardRewsult = levelRewardResult
	LevelResult.Charas = charaIds
	LevelResult.RewardRewsult.Cost = levelCfg.Cost

	return LevelResult, nil
}

func (u *User) onLevelFailed(ctx context.Context, levelCacheInfo LevelCacheInfo, characters []*pb.VOBattleCharacter, battleResult int32) error {

	switch levelCacheInfo.SystemType {
	case static.BattleTypeYgg:
		objectUid := levelCacheInfo.SystemParams[0]
		err := u.onYggdrasilLevelFail(ctx, objectUid, characters, battleResult)
		if err != nil {
			return errors.WrapTrace(err)
		}
	case static.BattleTypeChallengeAltar:
		objectUid := levelCacheInfo.SystemParams[0]
		err := u.onChallengeAltarFail(ctx, objectUid, characters, battleResult)
		if err != nil {
			return errors.WrapTrace(err)
		}
	}
	return nil
}

func (u *User) PassLevel(ctx context.Context, levelCfg *entry.Level, targets, achievements []int32,
	systemType int32, systemParams []int64, charaIds []int32,
	charactersAfterBattle []*pb.VOBattleCharacter, passReason *logreason.Reason) (*LevelRewardResult, error) {
	levelReswardResult := &LevelRewardResult{
		TotalRewards: common.NewRewards(),
	}

	levelId := levelCfg.Id

	//index start from 1
	passAchievement := []int32{}
	for _, achievementIdx := range achievements {
		if achievementIdx > 0 && achievementIdx <= levelCfg.AchievementsCount {
			passAchievement = append(passAchievement, achievementIdx)
		}
	}

	passTarget := []int32{}
	for _, targetIdx := range targets {
		if targetIdx > 0 && targetIdx <= levelCfg.TargetCount {
			passTarget = append(passTarget, targetIdx)
		}
	}

	newTargets, newAchievements := u.LevelsInfo.LevelPass(levelId, passTarget, passAchievement)
	for range newTargets {
		reason := logreason.NewReason(logreason.LevelTarget, logreason.AddRelateLog(passReason.RelateLog()))
		realRewards, err := u.AddRewardsByDropIds(levelCfg.TargetDrops, reason)
		if err != nil {
			glog.Errorf("level target drop err: %+v\n", err)
			return nil, err
		}

		levelReswardResult.TotalRewards.Append(realRewards)
	}
	levelReswardResult.TargetResult = u.VOResourceResult()

	for _, index := range newAchievements {
		reason := logreason.NewReason(logreason.LevelAchievement, logreason.AddRelateLog(passReason.RelateLog()))
		realRewards, err := u.AddRewardsByDropIds(levelCfg.AchievementsDrops[index-1], reason)
		if err != nil {
			glog.Errorf("level achievements drop err: %+v\n", err)
			return nil, err
		}

		levelReswardResult.TotalRewards.Append(realRewards)
	}

	levelReswardResult.AchievementResult = u.VOResourceResult()

	level, _ := u.LevelsInfo.GetLevel(levelId)
	if level.IsFirstPass() {
		realRewards, err := u.AddRewardsByDropIds(levelCfg.FirstDrop, passReason)
		if err != nil {
			glog.Errorf("level first pass drop err: %+v\n", err)
			return nil, err
		}

		levelReswardResult.TotalRewards.Append(realRewards)

		levelReswardResult.FirstPassResult = u.VOResourceResult()
	}

	realRewards, err := u.AddRewardsByDropIds(levelCfg.NormalDrop, passReason)
	if err != nil {
		glog.Errorf("level pass drop err: %+v\n", err)
	}

	levelReswardResult.TotalRewards.Append(realRewards)

	expReward := common.NewReward(static.CommonResourceTypeTeamExp, levelCfg.TeamExp)
	realRewards, _, _ = u.addReward(expReward, passReason)
	levelReswardResult.TotalRewards.Append(realRewards)

	for _, charaId := range charaIds {
		u.CharacterAddExp(charaId, levelCfg.CharacterExp)
	}

	switch systemType {
	case static.BattleTypeExplore:
		if level.IsFirstPass() {
			u.Info.SetLatestExploreLevel(levelId)
		}
	case static.BattleTypeExploreMonster:
		monsterId := systemParams[0]
		u.OnExploreMonsterLevelPass(levelId, int32(monsterId))

	case static.BattleTypeExploreResource:
		resourceId := systemParams[0]
		u.OnExploreResourceLevelPass(levelId, int32(resourceId))

	case static.BattleTypeTower:
		towerId := systemParams[0]
		stage := systemParams[1]
		u.onTowerLevelPass(int32(towerId), int32(stage))
	case static.BattleTypeYgg:
		objectUid := systemParams[0]
		err = u.onYggdrasilLevelPass(ctx, objectUid, levelCfg, charactersAfterBattle)
		if err != nil {
			glog.Errorf("PassLevel onYggdrasilLevelPass err: %+v", err)
			return nil, err
		}
	case static.BattleTypeChallengeAltar:
		objectUid := systemParams[0]
		err = u.onChallengeAltarPass(ctx, objectUid, levelCfg, charactersAfterBattle)
		if err != nil {
			glog.Errorf("PassLevel onChallengeAltarPass err: %+v", err)
			return nil, err
		}

	}

	if levelCfg.ChapterId > 0 {
		addScore := int32(len(newTargets))
		chapter := u.ChapterInfo.GetOrCreateChapter(levelCfg.ChapterId)
		chapter.AddScore(addScore)
		u.Info.SetExploredLevelVns(levelId)
		chapterCfg, err := manager.CSV.ChapterEntry.GetChapter(levelCfg.ChapterId)
		if err != nil {
			glog.Error(err.Error())
		} else {
			u.TriggerQuestUpdate(static.TaskTypeChapterScore, chapterCfg.ChapterType, addScore)
		}
	}

	levelReswardResult.PassResult = u.VOResourceResult()

	targetsMap := map[int32]bool{}
	for _, index := range newTargets {
		targetsMap[index] = true
	}
	score := int32(len(targets))
	u.TriggerQuestUpdate(static.TaskTypeLevelPass, levelId, score)
	u.TriggerQuestUpdate(static.TaskTypeLevelPassByType, systemType, int32(1))

	u.QuestCheckUnlock(static.ConditionTypePassLevel, levelId)

	return levelReswardResult, nil
}

func (u *User) SweepLevel(levelId, times int32) error {
	level, ok := u.LevelsInfo.GetLevel(levelId)
	if !ok {
		return errors.Swrapf(common.ErrLevelNotPassed, levelId)
	}

	if !level.IsAllTargetCleared() {
		return errors.Swrapf(common.ErrLevelTargetNotCompleteAll, levelId)
	}

	levelCfg, err := manager.CSV.LevelsEntry.GetLevel(levelId)
	if err != nil {
		return err
	}

	if !checkLevelPassTimes(level, levelCfg, times) {
		return errors.Swrapf(common.ErrLevelPassTimesLimited, levelId)
	}

	cost := levelCfg.Cost.Multiple(times)
	err = u.CheckRewardsEnough(cost)
	if err != nil {
		return err
	}

	reason := logreason.NewReason(logreason.SweepLevel)
	for i := int32(0); i < times; i++ {
		u.AddRewardsByDropIds(levelCfg.NormalDrop, reason)
		u.AddRewardsByDropIds(levelCfg.SweepDrop, reason)
		expReward := common.NewReward(static.CommonResourceTypeTeamExp, levelCfg.TeamExp)
		u.addReward(expReward, reason)
	}

	u.CostRewards(cost, reason)
	level.Sweep(times)

	u.TriggerQuestUpdate(static.TaskTypeLevelPassByType, int32(static.BattleTypeExploreElite), times)
	u.TriggerQuestUpdate(static.TaskTypeLevelEndByType, int32(static.BattleTypeExploreElite), times)

	return nil
}

func (u *User) LevelDailyRefresh(refreshTime int64) {
	for _, level := range u.LevelsInfo.Levels {
		level.ResetDailyTime()
		levelCfg, err := manager.CSV.LevelsEntry.GetLevel(level.Id)
		if err != nil {
			glog.Errorf("level daily refresh error: %+v", err)
			continue
		}

		if levelCfg.RefreshType == static.LevelRefreshTypeDaily {
			u.AddLevelNotify([]int32{level.Id})
		}
	}
}

func (u *User) checkLevelStart(levelId int32, formation *pb.VOBattleFormation,
	systemType int32, systemParams []int64) error {
	levelCfg, err := manager.CSV.LevelsEntry.GetLevel(levelId)
	if err != nil {
		return err
	}

	if levelCfg.ChapterId > 0 {
		err = u.CheckChapterUnlock(levelCfg.ChapterId)
		if err != nil {
			return err
		}
	}

	level, ok := u.LevelsInfo.GetLevel(levelId)
	if ok {
		if !checkLevelPassTimes(level, levelCfg, 1) {
			return errors.Swrapf(common.ErrLevelPassTimesLimited, levelId)
		}
	}

	err = u.CheckRewardsEnough(levelCfg.Cost)
	if err != nil {
		return err
	}

	err = u.CheckUserConditions(levelCfg.UnlockConditions)
	if err != nil {
		return err
	}

	charas := []int32{}
	battleId := levelCfg.BattleID
	if battleId > 0 {
		err := u.checkBattleFormation(systemType, battleId, formation)
		if err != nil {
			return err
		}
	}

	switch systemType {
	case static.BattleTypeExplore:
		chapterId := levelCfg.ChapterId
		if chapterId <= 0 {
			return errors.Swrapf(common.ErrNotCampainLevel, levelId)
		}

		chapterCfg, err := manager.CSV.ChapterEntry.GetChapter(chapterId)
		if err != nil {
			return errors.Swrapf(common.ErrNotCampainLevel, levelId)
		}

		if chapterCfg.ChapterType != ChapterTypeCampain {
			return errors.Swrapf(common.ErrNotCampainLevel, levelId)
		}

	case static.BattleTypeExploreElite:
		chapterId := levelCfg.ChapterId
		if chapterId <= 0 {
			return errors.Swrapf(common.ErrNotExploreEliteLevel, levelId)
		}

		chapterCfg, err := manager.CSV.ChapterEntry.GetChapter(chapterId)
		if err != nil {
			return errors.Swrapf(common.ErrNotExploreEliteLevel, levelId)
		}

		if chapterCfg.ChapterType != ChapterTypeElite {
			return errors.Swrapf(common.ErrNotExploreEliteLevel, levelId)
		}
	case static.BattleTypeGate:
		chapterId := levelCfg.ChapterId
		if chapterId <= 0 {
			return errors.Swrapf(common.ErrNotCampainLevel, levelId)
		}

		chapterCfg, err := manager.CSV.ChapterEntry.GetChapter(chapterId)
		if err != nil {
			return errors.Swrapf(common.ErrNotCampainLevel, levelId)
		}

		if chapterCfg.ChapterType != ChapterTypeChaos {
			return errors.Swrapf(common.ErrLevelInvalidSystemParam, levelId)
		}

	case static.BattleTypeExploreMonster:
		if len(systemParams) < 1 {
			return errors.Swrapf(common.ErrLevelInvalidSystemParam, levelId)
		}
		monsterId := systemParams[0]
		err := u.CheckExploreMonsterLevel(levelId, int32(monsterId))
		if err != nil {
			return err
		}

	case static.BattleTypeExploreResource:
		if len(systemParams) < 1 {
			return errors.Swrapf(common.ErrLevelInvalidSystemParam, levelId)
		}
		resouceId := systemParams[0]
		err := u.CheckExploreResourcePointLevel(levelId, int32(resouceId))
		if err != nil {
			return err
		}

	case static.BattleTypeTower:
		err = u.CheckActionUnlock(static.ActionIdTypeTowerunlock)
		if err != nil {
			return err
		}
		if len(systemParams) < 2 {
			return errors.Swrapf(common.ErrLevelInvalidSystemParam, levelId)
		}

		towerId := systemParams[0]
		towerStage := systemParams[1]
		err := u.checkTowerLevel(int32(towerId), int32(towerStage), levelId, charas)
		if err != nil {
			return err
		}

	case static.BattleTypeYgg:
		if len(systemParams) != 1 {
			return errors.Swrapf(common.ErrLevelInvalidSystemParam, levelId)
		}
		objectUid := systemParams[0]
		err := u.checkYggdrasilBattle(objectUid, levelId, formation.BattleCharacters, formation.Npcs)
		if err != nil {
			return errors.WrapTrace(err)
		}
	case static.BattleTypeChallengeAltar:
		if len(systemParams) != 1 {
			return errors.Swrapf(common.ErrLevelInvalidSystemParam, levelId)
		}
		objectUid := systemParams[0]
		err := u.checkChallengeAltarBattle(objectUid, levelId, formation.BattleCharacters, formation.Npcs)
		if err != nil {
			return errors.WrapTrace(err)
		}

	}

	return nil
}

func checkLevelPassTimes(level *Level, levelCfg *entry.Level, times int32) bool {
	switch levelCfg.RefreshType {
	case static.LevelRefreshTypeInfinity:

	case static.LevelRefreshTypeOnce:
		if level.GetTotalPass()+times > levelCfg.ChallengeTimes {
			return false
		}

	case static.LevelRefreshTypeDaily:
		if level.GetTodayPass()+times > levelCfg.ChallengeTimes {
			return false
		}
	}

	return true
}

func (u *User) genBattleDetailData(formation *pb.VOBattleFormation) (*pb.VOBattleFighterDetail, error) {
	battleCharaProps := make([]*pb.VOBattleCharacterProp, 0, len(formation.BattleCharacters)+len(formation.Npcs))
	for _, battleChara := range formation.BattleCharacters {
		charaProp := u.VOBattleCharacter(battleChara.CharacterId)
		battleCharaProps = append(battleCharaProps, charaProp)
	}

	for _, npc := range formation.Npcs {
		charaProp, err := u.genBattleNPCDetailData(npc.NpcId, npc.NpcType)
		if err != nil {
			return nil, err
		}
		battleCharaProps = append(battleCharaProps, charaProp)
	}

	detail := &pb.VOBattleFighterDetail{
		UserInfo:        u.VOUserInfoSimple(),
		Power:           0,
		BattleFormation: formation,
		Characters:      battleCharaProps,
	}

	battleHero := formation.GetBattleHero()
	heroId := battleHero.GetHeroId()
	hero, ok := u.HeroPack.GetHero(heroId)
	if ok {
		detail.Hero = &pb.VOBattleHeroProp{
			Hero:      hero.VOHero(),
			HeroLevel: u.HeroPack.GetLevel(),
		}
	}

	return detail, nil
}

func (u *User) genBattleNPCDetailData(npcId, npcType int32) (*pb.VOBattleCharacterProp, error) {
	var data *pb.VOBattleCharacterProp = &pb.VOBattleCharacterProp{}

	switch npcType {
	case static.BattleNpcTypeSystem:
		npcCfg, err := manager.CSV.Battle.GetBattleNPC(npcId)
		if err != nil {
			return nil, err
		}

		chara := &pb.VOUserCharacter{
			CharacterId:      npcCfg.CharaID,
			Exp:              0,
			Level:            npcCfg.Level,
			Star:             npcCfg.Star,
			Stage:            npcCfg.Stage,
			CreateAt:         0,
			Skills:           make([]*pb.VOCharacterSkill, 0, len(npcCfg.SkillLv)),
			HeroId:           0,
			HeroPos:          0,
			Power:            0,
			CanYggdrasilTime: 0,
			Rarity:           0,
		}

		for skillId, level := range npcCfg.SkillLv {
			chara.Skills = append(chara.Skills, &pb.VOCharacterSkill{SkillId: skillId, Level: level})
		}

		data.Character = chara
		data.Equipments = make([]*pb.VOUserEquipment, 0, len(npcCfg.Equipments))
		for _, equipment := range npcCfg.Equipments {
			data.Equipments = append(data.Equipments, equipment.VOUserEquipment())
		}

		if npcCfg.WorldItem > 0 {
			worldItem := common.NewWorldItem(0, npcCfg.WorldItem)
			data.WorldItem = worldItem.VOUserWorldItem()
		}

	case static.BattleNpcTypePlayer:
		mercenary, _ := u.GetMercenary(npcId)
		data.Character = mercenary.Character.VOUserCharacter()
		data.Equipments = make([]*pb.VOUserEquipment, 0, len(mercenary.Equipments))
		for _, equipment := range mercenary.Equipments {
			data.Equipments = append(data.Equipments, equipment.VOUserEquipment())
		}

		if mercenary.WorldItem != nil {
			data.WorldItem = mercenary.WorldItem.VOUserWorldItem()
		}
	}

	return data, nil
}
