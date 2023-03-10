package model

import (
	"gamesvr/manager"
	"math"
	"shared/common"
	"shared/csv/entry"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/statistic/bilog"
	"shared/statistic/logreason"
	"shared/utility/coordinate"
	"shared/utility/errors"
	"shared/utility/glog"
	"shared/utility/servertime"
	"strconv"
	"time"
)

// GraveyardDailyRefresh 每日刷新
func (u *User) GraveyardDailyRefresh(refreshTime int64) {
	// 小人领奖次数置0
	u.Graveyard.PlotRecord.SetRewardNum(0)
	// 小人领奖时间重置
	u.Graveyard.PlotRecord.RewardHours = manager.CSV.GraveyardEntry.RandomRewardHours()
	// 重置公会互助相关
	u.Graveyard.SendRequestCount = 0
	u.Graveyard.DailyGuildGoldByHelp = 0
	u.Graveyard.DailyAddActivationByHelp = 0
}

func (u *User) GraveyardGetInfo() *pb.S2CGraveyardGetInfo {

	return &pb.S2CGraveyardGetInfo{
		Builds:           u.Graveyard.VOAll(u),
		SendRequestCount: u.Graveyard.SendRequestCount,
		RewardHours:      u.GraveyardGetRewardHours(),
	}
}

func (u *User) GraveyardGetRewardHours() []int32 {
	return u.Graveyard.PlotRecord.GetRewardHours()

}

func (u *User) GraveyardBuildCreate(buildId int32, position *coordinate.Position) (*pb.VOGraveyardBuild, error) {
	building, err := manager.CSV.GraveyardEntry.GetBuildById(buildId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// 主堡不能再建
	if building.Type == static.GraveyardBuildTypeMaintower {
		return nil, common.ErrGraveyardTypeCannotBuild
	}

	mainTowerLevel := u.Graveyard.GetMainTowerLevel()
	// 超过最多可建造个数
	buildCount := u.Graveyard.GetBuildsCountByBuildId(buildId)

	// 区域不可建造
	unlockArea, err := manager.CSV.GraveyardEntry.CalUnlockArea(mainTowerLevel, u.Graveyard.GetLocatedBuildsList())
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	err = unlockArea.MinusBuildingArea(building.BuildingArea, position)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	biCommon := u.CreateBICommon(bilog.EventNameGraveyard)

	var newBuild *common.UserGraveyardBuild
	switch building.Type {
	case static.GraveyardBuildTypeContinuous, static.GraveyardBuildTypeConsumeProduceItem:
		if buildCount-1 >= manager.CSV.GraveyardEntry.CanBuildCount(mainTowerLevel, buildId) {
			return nil, common.ErrGraveyardNumLimitCannotBuild
		}
		err := u.CheckUserConditions(building.UnlockCondition)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		// 产出型建造
		// 建造所需时间
		lvUp, err := manager.CSV.GraveyardEntry.GetCreate(buildId)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		// 消耗
		consume := lvUp.LevelUpConsume
		err = u.CheckRewardsEnough(consume)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		reason := logreason.NewReason(logreason.GraveyardBuild, logreason.AddRelateLog(biCommon.GetLogId()))
		err = u.CostRewards(consume, reason)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}

		newBuild = common.NewUserGraveyardBuild(buildId, position, lvUp.LevelUpTime)
		// 新建筑拷贝buff
		for _, tmp := range u.Graveyard.GetBuildsByBuildId(buildId) {
			newBuild.UserGraveyardProduceBuffs = tmp.Clone()
			newBuild.InvalidateBuff()
			break
		}

	case static.GraveyardBuildTypeDecoration:
		if buildCount >= manager.CSV.GraveyardEntry.CanBuildCount(mainTowerLevel, buildId) {
			return nil, common.ErrGraveyardNumLimitCannotBuild
		}
		// 装饰物
		newBuild = common.NewUserGraveyardAdornment(buildId, position)
	default:
		return nil, common.ErrGraveyardTypeCannotBuild
	}
	// 放入地图
	buildUid := u.Graveyard.CreateBuild(newBuild)

	u.BIGraveyard(bilog.GraveyardOpBuild, mainTowerLevel, newBuild, 0, biCommon)

	if building.Type == static.GraveyardBuildTypeDecoration {
		u.TriggerQuestUpdate(static.TaskTypeGraveyarDecorationSetting, 1)
	}

	return u.VOGraveyardBuild(buildUid, newBuild), nil
}

func (u *User) GraveyardRelocation(relocationMap map[int64]*coordinate.Position) (*pb.S2CGraveyardRelocation, error) {

	buildings := map[int64]*common.UserGraveyardBuild{}
	var buildsNeedRelocation []int64
	// 检查uid是否是已拥有建筑
	for uid, position := range relocationMap {
		build, err := u.Graveyard.FindByUid(uid)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		// 只有装饰物可放入背包
		if position == nil {
			building, err := manager.CSV.GraveyardEntry.GetBuildById(build.BuildId)
			if err != nil {
				return nil, errors.WrapTrace(err)
			}
			if building.Type != static.GraveyardBuildTypeDecoration {
				return nil, errors.WrapTrace(common.ErrParamError)
			}
		}
		buildings[uid] = build
		buildsNeedRelocation = append(buildsNeedRelocation, uid)
	}

	// 判断是否能放入区域
	area, err := manager.CSV.GraveyardEntry.CalUnlockArea(u.Graveyard.GetMainTowerLevel(), u.Graveyard.GetLocatedBuildsListExcept(buildsNeedRelocation))
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	for uid, build := range buildings {
		position := relocationMap[uid]
		if position == nil {
			continue
		}
		buildConfig, err := manager.CSV.GraveyardEntry.GetBuildById(build.BuildId)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		err = area.MinusBuildingArea(buildConfig.BuildingArea, position)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
	}

	now := servertime.Now().Unix()
	// 重新设置位置
	for uid, build := range buildings {
		oldPosition := build.Position
		build.Position = relocationMap[uid]
		if build.Position == nil {
			build.InBagTime = now
		}

		if oldPosition == nil {
			buildConfig, err := manager.CSV.GraveyardEntry.GetBuildById(build.BuildId)
			if err != nil {
				continue
			}
			if buildConfig.Type == static.GraveyardBuildTypeDecoration {
				u.TriggerQuestUpdate(static.TaskTypeGraveyarDecorationSetting, 1)
			}
		}
	}
	// 组增量的返回数据
	return &pb.S2CGraveyardRelocation{
		Builds: u.Graveyard.VOPartial(buildings, u),
	}, nil

}

func (u *User) GraveyardBuildLvUp(buildUId int64) (*pb.VOGraveyardBuild, error) {
	build, err := u.Graveyard.FindByUid(buildUId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	err = manager.CSV.GraveyardEntry.CheckInNormalState(build)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 获得升级的配置
	toLv := build.FetchLevel() + 1
	buildingLvUp, err := manager.CSV.GraveyardEntry.GetLvUp(build.BuildId, toLv)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// 检查
	err = u.Graveyard.CheckLvUpLimitAndConsume(buildingLvUp, u)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	biCommon := u.CreateBICommon(bilog.EventNameGraveyard)

	reason := logreason.NewReason(logreason.GraveyardBuildingLevelUp, logreason.AddRelateLog(biCommon.GetLogId()))
	err = u.CostRewards(buildingLvUp.LevelUpConsume, reason)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// 进入升级状态
	build.LvUpTransition(toLv, buildingLvUp.LevelUpTime)

	mainTowerLevel := u.Graveyard.GetMainTowerLevel()
	u.BIGraveyard(bilog.GraveyardOpLevelUp, mainTowerLevel, build, 0, biCommon)

	return u.VOGraveyardBuild(buildUId, build), nil

}

func (u *User) GraveyardBuildStageUp(buildUId int64) (*pb.VOGraveyardBuild, error) {
	build, err := u.Graveyard.FindByUid(buildUId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	err = manager.CSV.GraveyardEntry.CheckInNormalState(build)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// 获得升阶的配置
	toStage := build.FetchStage() + 1
	buildingStageUp, err := manager.CSV.GraveyardEntry.GetStageUp(build.BuildId, toStage)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	err = u.Graveyard.CheckStageUpLimitAndConsume(build, buildingStageUp, u)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	biCommon := u.CreateBICommon(bilog.EventNameGraveyard)
	reason := logreason.NewReason(logreason.GraveyardBuildingStageUp, logreason.AddRelateLog(biCommon.GetLogId()))
	err = u.CostRewards(buildingStageUp.StageUpConsume, reason)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// 进入升阶状态
	build.StageUpTransition(toStage, buildingStageUp.StageUpTime)

	mainTowerLevel := u.Graveyard.GetMainTowerLevel()
	u.BIGraveyard(bilog.GraveyardOpLevelUp, mainTowerLevel, build, 0, biCommon)

	return u.VOGraveyardBuild(buildUId, build), nil

}

func (u *User) GraveyardOpenCurtain(buildUId int64) (*pb.VOGraveyardBuild, error) {
	build, err := u.Graveyard.FindByUid(buildUId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	if build.UserGraveyardTransition == nil {
		return nil, common.ErrGraveyardBuildNotInTransaction
	}
	if build.UserGraveyardTransition.EndAt > servertime.Now().Unix() {
		return nil, common.ErrGraveyardBuildTransactionTimeNotEnough
	}

	biCommon := u.CreateBICommon(bilog.EventNameGraveyard)

	// 结算产出
	u.RefreshOutPut(build, false)
	// 发揭幕奖励
	drop := manager.CSV.GraveyardEntry.GetCurtainDrop(build)
	if drop > 0 {
		curtainReasonEnum := logreason.GraveyardCreateCurtain
		switch build.UserGraveyardTransition.CurtainType {
		case common.CurtainTypeCreate:
			curtainReasonEnum = logreason.GraveyardCreateCurtain
		case common.CurtainTypeLvUp:
			curtainReasonEnum = logreason.GraveyardLevelUpCurtain
		case common.CurtainTypeStageUp:
			curtainReasonEnum = logreason.GraveyardStageUpCurtain
		}
		reason := logreason.NewReason(curtainReasonEnum, logreason.AddRelateLog(biCommon.GetLogId()))
		_, err = u.AddRewardsByDropId(drop, reason)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
	}

	transition := build.UserGraveyardTransition

	// 建筑从建造揭幕开始计算产出
	build.StartProduce()
	// 结束进度条
	build.EndTransition()

	opType := bilog.GraveyardOpBuild

	switch transition.CurtainType {
	case common.CurtainTypeLvUp:
		newLevel := build.FetchLevel()
		u.TriggerQuestUpdate(static.TaskTypeGraveyarBuildingLevelCount, build.BuildId, newLevel-1, newLevel)

		opType = bilog.GraveyardOpLevelUp

	case common.CurtainTypeStageUp:
		newStage := build.FetchStage()
		u.TriggerQuestUpdate(static.TaskTypeGraveyarBuildingStageCount, build.BuildId, newStage-1, newStage)

		opType = bilog.GraveyardOpStageUp

	case common.CurtainTypeCreate:
		opType = bilog.GraveyardOpBuild
		u.TriggerQuestUpdate(static.TaskTypeGraveyarBuildingLevelCount, build.BuildId, 0, 1)
		u.TriggerQuestUpdate(static.TaskTypeGraveyarBuildingStageCount, build.BuildId, 0, 1)

	}
	mainTowerLevel := u.Graveyard.GetMainTowerLevel()
	u.BIGraveyard(opType, mainTowerLevel, build, 0, biCommon)

	return u.VOGraveyardBuild(buildUId, build), nil
}

func (u *User) GraveyardProduceStart(buildUId int64, produceNum int32) ([]*pb.VOGraveyardBuild, error) {
	build, err := u.Graveyard.FindByUid(buildUId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	consume, err := manager.CSV.GraveyardEntry.CheckProduceNum(build, produceNum)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// 消耗
	err = u.CheckRewardsEnough(consume)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	biCommon := u.CreateBICommon(bilog.EventNameGraveyard)
	reason := logreason.NewReason(logreason.GraveyardProduce, logreason.AddRelateLog(biCommon.LogId))
	err = u.CostRewards(consume, reason)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	build.StartConsumeProduce(produceNum)
	var buildVos []*pb.VOGraveyardBuild
	// 同id建筑共享次数类buff
	for tmpUid, tmpBuild := range u.Graveyard.GetBuildsByBuildId(build.BuildId) {
		if tmpUid != buildUId {
			if tmpBuild.CountBuff != nil {
				tmpBuild.CountBuff.RestCount--
				//若同id建筑存在多个，且次数型buff剩余1次：此情况下首次进行生产的建筑享有buff加成，生产进度条上显示icon；其他同id建筑界面隐藏buff剩余信息，且进行生产不再享有buff加成且无icon显示
				tmpBuild.InvalidateCountBuff()
			}
		}

		buildVos = append(buildVos, u.VOGraveyardBuild(tmpUid, tmpBuild))
	}

	mainTowerLevel := u.Graveyard.GetMainTowerLevel()
	u.BIGraveyard(bilog.GraveyardOpStartProduce, mainTowerLevel, build, 0, biCommon)

	u.Guild.AddTaskItem(static.GuildTaskEquipment, produceNum)
	return buildVos, nil
}

func (u *User) GraveyardProductionGet(buildUIdList []int64) ([]*pb.VOGraveyardBuild, error) {
	// 传入的uid如果不能收获也不报错直接返回
	var voList []*pb.VOGraveyardBuild

	// 传入的uid如果不能收获也不报错直接返回
	for _, buildUId := range buildUIdList {
		build, err := u.Graveyard.FindByUid(buildUId)
		if err != nil {
			continue
		}

		// 刷新产出
		u.RefreshOutPut(build, true)
		if manager.CSV.GraveyardEntry.CanProductionGet(build) {
			biCommon := u.CreateBICommon(bilog.EventNameGraveyard)

			glog.Debugf("GraveyardProductionGet: productions: %v", build.Productions)
			reason := logreason.NewReason(logreason.GraveyardProduction)
			realRewards, err := u.addRewards(build.Productions, reason)
			if err != nil {
				glog.Errorf(" GraveyardProductionGet addRewards error:%+v", errors.Format(err))
			}

			u.triggerGraveyardProductQuest(realRewards)

			mainTowerLevel := u.Graveyard.GetMainTowerLevel()
			buildType := manager.CSV.GraveyardEntry.GetBuildType(build.BuildId)
			if buildType == static.GraveyardBuildTypeContinuous {
				// 重置产出
				build.ResetContinuousProduce()
				u.BIGraveyard(bilog.GraveyardOpGetProduction, mainTowerLevel, build, 0, biCommon)
			} else if buildType == static.GraveyardBuildTypeConsumeProduceItem {
				build.ResetConsumeProduce()

				u.BIGraveyard(bilog.GraveyardOpGetProduction, mainTowerLevel, build, 0, biCommon)
			}

			voList = append(voList, u.VOGraveyardBuild(buildUId, build))

		}

		glog.Debugf("GraveyardProductionGet: rewards: %v", u.RewardsResult)
	}

	if len(voList) > 0 {
		u.TriggerQuestUpdate(static.TaskTypeGraveyarProductCollectTimes, int32(1))
	}

	return voList, nil
}

func (u *User) GraveyardRefreshBuildInfo(buildUIdList []int64) ([]*pb.VOGraveyardBuild, error) {
	var voList []*pb.VOGraveyardBuild

	//传入的uid如果不能刷新也不报错直接返回
	for _, buildUId := range buildUIdList {
		build, err := u.Graveyard.FindByUid(buildUId)
		if err != nil {
			continue
		}
		voList = append(voList, u.VOGraveyardBuild(buildUId, build))
	}

	return voList, nil
}

func (u *User) GraveyardCharacterDispatch(buildUId int64, characters *common.CharacterPositions) ([]*pb.VOGraveyardBuild, error) {
	build, err := u.Graveyard.FindByUid(buildUId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查角色是否存在
	for cid := range characters.Positions {
		_, err := u.CharacterPack.Get(cid)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
	}
	// 入驻角色替换
	characterPositionMap, err := manager.CSV.GraveyardEntry.CalCharacterPositionMap(u.Graveyard.Buildings, buildUId, build, characters)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	var voList []*pb.VOGraveyardBuild

	newUpCharaCount := 0
	for uid, charactersPosition := range characterPositionMap {
		tmpBuild, err := u.Graveyard.FindByUid(uid)
		if err != nil {
			return nil, err
		}

		oldCharas, ok := tmpBuild.FetchCharacters()
		if !ok {
			oldCharas = &common.CharacterPositions{Positions: map[int32]int32{}}
		}
		// 结算产出
		u.RefreshOutPut(tmpBuild, false)
		tmpBuild.SetCharacters(charactersPosition)

		for cid := range charactersPosition.Positions {
			_, ok := oldCharas.Positions[cid]
			if !ok {
				newUpCharaCount = newUpCharaCount + 1
			}
		}

		voList = append(voList, u.VOGraveyardBuild(uid, tmpBuild))
	}

	biCommon := u.CreateBICommon(bilog.EventNameGraveyard)
	mainTowerLevel := u.Graveyard.GetMainTowerLevel()
	u.BIGraveyard(bilog.GraveyardOpGetProduction, mainTowerLevel, build, 0, biCommon)

	u.TriggerQuestUpdate(static.TaskTypeGraveyarCharacterSetting, newUpCharaCount)

	return voList, nil
}

func (u *User) GraveyardPopulationSet(populationMap map[int64]int32) ([]*pb.VOGraveyardBuild, error) {

	// 参数合法性检查
	for uid, population := range populationMap {
		byUid, err := u.Graveyard.FindByUid(uid)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		err = manager.CSV.GraveyardEntry.CanPopulationSet(byUid.BuildId, byUid.FetchLevel(), population)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
	}

	var voList []*pb.VOGraveyardBuild
	for uid, population := range populationMap {
		byUid, _ := u.Graveyard.FindByUid(uid)

		// 结算产出
		u.RefreshOutPut(byUid, false)
		byUid.PopulationSet(population)
		voList = append(voList, u.VOGraveyardBuild(uid, byUid))

	}

	return voList, nil

}

// RefreshOutPut 用于刷新建筑产出，更新 build.UserGraveyardProduce.Productions
func (u *User) RefreshOutPut(build *common.UserGraveyardBuild, manualFinish bool) {
	manager.CSV.RefreshOutPut(build, u.GetBuildCharacterStar(build), manualFinish)

}

// GetBuildCharacterStar 获得该建筑入驻角色星级
func (u *User) GetBuildCharacterStar(build *common.UserGraveyardBuild) map[int32]int32 {
	characterStarMap := map[int32]int32{}
	if build.UserGraveyardProduce == nil {
		return characterStarMap
	}
	for _, i := range build.Characters.GetCharacters() {
		character, err := u.CharacterPack.Get(i)
		if err != nil {
			glog.Errorf("RefreshOutPut CharacterPack.ExecuteEventsInLoop error:%+v", errors.Format(err))
			continue
		}
		characterStarMap[i] = character.GetStar()
	}
	return characterStarMap
}

func (u *User) GraveyardAccelerate(buildUId int64, rewards *common.Rewards) (*pb.VOGraveyardBuild, error) {
	build, err := u.Graveyard.FindByUid(buildUId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	mergeValue := rewards.MergeValue()

	// 检查道具参数
	items := map[*entry.GraveyardAccItem]int32{}
	for _, reward := range mergeValue {
		graveyardAcc, ok := manager.CSV.Item.GetGraveyardAcc(reward.ID)
		if !ok {
			return nil, common.ErrParamError
		}
		items[graveyardAcc] = reward.Num
	}
	if len(items) == 0 {
		return nil, common.ErrParamError
	}
	// 加速
	consume, err := manager.CSV.Acc(build, items, u.GetBuildCharacterStar(build))
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	biCommon := u.CreateBICommon(bilog.EventNameGraveyard)

	// 消耗
	reason := logreason.NewReason(logreason.GraveyardAccelerate, logreason.AddRelateLog(biCommon.GetLogId()))
	err = u.CostRewards(consume, reason)

	// 生产加速时间任务监听
	var accSec int32
	for _, reward := range consume.MergeValue() {
		graveyardAcc, ok := manager.CSV.Item.GetGraveyardAcc(reward.ID)
		if !ok {
			continue
		}
		if graveyardAcc.Acc.AccType != static.GraveyardAccelerateTypeReduceProduceTime {
			continue
		}
		accSec += reward.Num * graveyardAcc.Acc.Sec
	}
	// 生产加速时间任务监听
	if accSec > 0 {
		u.TriggerQuestUpdate(static.TaskTypeGraveyarAccelerateTimes, build.BuildId, int32(1))
		u.TriggerQuestUpdate(static.TaskTypeGraveyarAccelerateTime, build.BuildId, int32(accSec/servertime.SecondPerMinute))

		mainTowerLevel := u.Graveyard.GetMainTowerLevel()
		u.BIGraveyard(bilog.GraveyardOpAccelerate, mainTowerLevel, build, 0, biCommon)

	}

	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return u.VOGraveyardBuild(buildUId, build), nil
}

func (u *User) VOGraveyardBuild(buildUId int64, build *common.UserGraveyardBuild) *pb.VOGraveyardBuild {
	u.RefreshOutPut(build, false)
	return common.VOGraveyardBuild(buildUId, build)
}

func (u *User) GraveyardReceiveHelp(helpType int, buildUid int64, sec int32, helpAt int64) {
	build, err := u.Graveyard.FindByUid(buildUid)
	if err != nil {
		return
	}
	switch helpType {
	case common.GraveyardHelpTypeBuild, common.GraveyardHelpTypeLvUp, common.GraveyardHelpTypeStageUp:
		if build.UserGraveyardTransition == nil {
			return
		}
		if int32(helpType) != build.CurtainType {
			return
		}
		build.BuildAcc(sec)
	case common.GraveyardHelpTypeProduct:
		if manager.CSV.GraveyardEntry.GetBuildType(build.BuildId) != static.GraveyardBuildTypeConsumeProduceItem {
			return
		}
		// 不在生产中
		if build.CurrProduceNum == 0 {
			return
		}
		build.AccRecords.AddRecord(sec, helpAt)
	default:
		return
	}
	build.IncreaseReceiveHelpCount()
	// 推送S2CGraveyardReceiveHelp
	u.AddGraveyardPush(&pb.S2CGraveyardReceiveHelp{
		Build: u.VOGraveyardBuild(buildUid, build),
	})

}

func (u *User) GenVOAddGraveyardRequest(buildUid int64, build *common.UserGraveyardBuild) (*pb.VOAddGraveyardRequest, error) {

	if u.Graveyard.SendRequestCount >= manager.CSV.GraveyardEntry.GetGraveyardDailyHelpSendCount() {
		return nil, errors.WrapTrace(common.ErrGraveyardNoRemainedHelpCount)
	}
	if build.RequestState == common.RequestStateDuring {
		return nil, errors.WrapTrace(common.ErrGraveyardBuildInHelpNow)
	}
	var helpType int
	var expireAt int64
	var totalSec int32
	if build.UserGraveyardTransition != nil {
		helpType = int(build.CurtainType)
		expireAt = build.EndAt
		totalSec = build.TransitionNeedSec
	} else {
		buildingType := manager.CSV.GraveyardEntry.GetBuildType(build.BuildId)
		if buildingType != static.GraveyardBuildTypeConsumeProduceItem {
			return nil, errors.WrapTrace(common.ErrParamError)
		}
		// 非可持续建筑不在生产中
		if build.CurrProduceNum == 0 {
			return nil, errors.WrapTrace(common.ErrParamError)
		}
		produceTime, err := manager.CSV.GraveyardEntry.GetProduceTime(build, u.GetBuildCharacterStar(build))
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		// 已生产完成
		totalSec = int32(math.Floor(float64(build.CurrProduceNum) * produceTime))
		produceSecs := build.CalProduceSecs()
		if totalSec <= produceSecs {
			return nil, errors.WrapTrace(err)
		}

		expireAt = servertime.Now().Unix() + int64(totalSec-produceSecs)

		helpType = common.GraveyardHelpTypeProduct

	}

	return &pb.VOAddGraveyardRequest{
		BuildId:  build.BuildId,
		BuildUid: buildUid,
		BuildLv:  build.FetchLevel(),
		UserId:   u.ID,
		HelpType: int32(helpType),
		TotalSec: totalSec,
		ExpireAt: expireAt,
	}, nil

}

func (u *User) setBuildLvByGm(params []string) (*pb.VOGraveyardBuild, error) {
	if len(params) != 2 {
		return nil, errors.WrapTrace(common.ErrParamError)
	}

	uid, err := strconv.ParseInt(params[0], 10, 64)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	lv, err := strconv.ParseInt(params[1], 10, 32)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	build, err := u.Graveyard.FindByUid(uid)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	build.Lv = int32(lv)

	return u.VOGraveyardBuild(uid, build), nil
}

func (u *User) GraveyardReceivePlotReward() ([]int32, error) {
	hours := u.Graveyard.PlotRecord.GetRewardHours()
	if len(hours) == 0 {
		return nil, errors.WrapTrace(common.ErrGraveyardPlotRewardNumMax)
	}
	now := servertime.Now()
	start := servertime.DayOffsetZeroTime(now, time.Hour*time.Duration(hours[0]))
	end := TodayRefreshTime().AddDate(0, 0, 1)
	// 判断是否在可领取时间内
	if now.Before(start) || now.After(end) {
		return nil, errors.WrapTrace(common.ErrGraveyardCannotPlotRewardNow)
	}

	u.Graveyard.PlotRecord.IncrRewardNum(1)
	reason := logreason.NewReason(logreason.GraveyardPlot)
	_, err := u.AddRewardsByDropId(manager.CSV.GraveyardEntry.GraveyardPlotRewardDropId(), reason)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return u.Graveyard.PlotRecord.GetRewardHours(), nil
}

func (u *User) GraveyardUseBuff(itemId, useNum int32) ([]*pb.VOGraveyardBuild, error) {
	consume := common.NewRewards()
	consume.AddReward(common.NewReward(itemId, useNum))

	// 检查消耗
	err := u.CheckRewardsEnough(consume)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	item, ok := manager.CSV.Item.GetGraveyardBuffItem(itemId)
	if !ok {
		return nil, errors.WrapTrace(common.ErrParamError)
	}

	var buffs []*entry.ProduceBuff
	buildsExist := false
	// 判断能否使用buff
	for _, buffId := range item.BuffIds {
		config, err := manager.CSV.GraveyardEntry.GetBuffConfig(buffId)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		builds := u.Graveyard.GetBuildsByBuildId(config.BuildId)
		if len(builds) > 0 {
			buildsExist = true
		}
		for _, build := range builds {
			duringBuff, ok := build.DuringBuff()
			if ok && duringBuff.GetBuffId() != config.Id {
				// 不能使用不同buff
				return nil, common.ErrGraveyardBuffInUse
			}
		}
		buffs = append(buffs, config)
	}
	if !buildsExist {
		return nil, common.ErrGraveyardCannotUseBuffBuildsNotExist
	}

	// 消耗道具
	reason := logreason.NewReason(logreason.GraveyardUseBuff)
	err = u.CostRewards(consume, reason)
	var vos []*pb.VOGraveyardBuild
	for ; useNum > 0; useNum-- {
		//使用buff
		for _, buff := range buffs {
			for uid, build := range u.Graveyard.GetBuildsByBuildId(buff.BuildId) {
				build.AddBuff(buff.Id, buff.Type, buff.TypeContent, servertime.Now().Unix())
				vos = append(vos, u.VOGraveyardBuild(uid, build))
			}
		}
	}
	return vos, nil

}

func (u *User) triggerGraveyardProductQuest(product *common.Rewards) {
	for _, rewards := range *product {
		for _, reward := range rewards {
			u.TriggerQuestUpdate(static.TaskTypeGraveyarProductRewardType, reward.Type, reward.Num)
			u.TriggerQuestUpdate(static.TaskTypeGraveyarProductItem, reward.ID, reward.Num)

			if reward.Type == static.ItemTypeEquipment {
				part, err := manager.CSV.Equipment.Part(reward.ID)
				if err != nil {
					glog.Errorf(err.Error())
				}

				u.TriggerQuestUpdate(static.TaskTypeGraveyarProductEquipmentPart, part, reward.Num)
			}
		}
	}
}
