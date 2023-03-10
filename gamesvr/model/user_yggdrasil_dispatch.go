package model

import (
	"context"
	"gamesvr/manager"
	"shared/common"
	"shared/csv/entry"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/statistic/logreason"
	"shared/utility/errors"
	"shared/utility/rand"
)

// todo 功能解锁

func (u *User) YggdrasilDispatchDailyRefresh(refreshTime int64) {
	u.Yggdrasil.Dispatch.YggDispatchDailyRefresh()
}

// todo 临时的派遣任务
func (u *User) YggdrasilTempDispatchInfo() ([]*pb.VOYggDispatchTaskState, []*pb.VOYggDispatchTaskState, error) {

	tasksCsv := manager.CSV.Yggdrasil.GetAllDispatchTask()
	if len(tasksCsv) == 0 {
		return nil, nil, errors.New("no taskcsv")
	}
	dailyTasks := []*pb.VOYggDispatchTaskState{}
	guildTasks := []*pb.VOYggDispatchTaskState{}

	// todo entry要检查日常派遣的type
	for _, taskCsv := range tasksCsv {
		if taskCsv.Type == static.YggdrasilDispatchTypeGuild {
			u.Yggdrasil.Dispatch.AddNewGuildTask(taskCsv.Id, taskCsv.CloseTime)
		}
		//  else {
		// 	u.Yggdrasil.Dispatch.AddNewGuildTask(taskCSV.Id, taskCSV.CloseTime)
		// }
	}

	for _, guildTask := range u.Yggdrasil.Dispatch.GuildTasks {
		if len(guildTasks) >= 10 {
			break
		}
		guildTasks = append(guildTasks, guildTask.VOYggDispatchTaskState())
	}

	return dailyTasks, guildTasks, nil
}

func (u *User) YggdrasilDispatchGetInfo() ([]*pb.VOYggDispatchTaskState, error) {
	var generateTasks []int32 // 需要新生成的任务
	var taskInfos []*pb.VOYggDispatchTaskState

	// 筛选出所有已经解锁派遣任务的区域id
	for _, area := range u.Yggdrasil.Areas.Areas {
		// fmt.Printf("----------areaId: %d, Prestige: %d \n", area.AreaId, area.Prestige)
		areaCSV, err := manager.CSV.Yggdrasil.GetArea(area.AreaId)
		if err != nil {
			return nil, err
		}
		if areaCSV.DispatchUnlock <= area.FetchPrestige(u) { // 解锁
			// 查看这个区域每个星级需要生成多少个任务
			for star, num := range areaCSV.DailyStar {
				// 根据每个区域解锁的日常任务的星级和数量去对应生成任务
				taskIDs, err := u.YggDispatchGenerateDailyTasks(areaCSV.Id, star, num)
				if err != nil {
					return nil, errors.WrapTrace(err)
				}
				// fmt.Printf("===================areaId: %d, Prestige: %d, taskId: %v", area.AreaId, area.Prestige, taskIDs)
				generateTasks = append(generateTasks, taskIDs...)
			}
		}
	}

	for _, newTask := range generateTasks {
		u.Yggdrasil.Dispatch.AddNewDailyTask(newTask)
	}

	// 更新所有的状态
	for _, taskState := range u.Yggdrasil.Dispatch.DailyTasks {
		taskState.YggDispatchUpdateTask()
		if !taskState.CheckStateAfterRewards() {
			taskInfos = append(taskInfos, taskState.VOYggDispatchTaskState())
		}
	}

	// lua.MaxArrayIndex = 4000
	// proto, err := CompileLua("E:/code/go/src/overlord-backend-go-pro/gamesvr/model/CombatPowerUtil.lua")
	// if err != nil {
	// 	fmt.Println(common.ErrCompileLuaFileFailed)
	// }

	// L := NewLuaMatcher()
	// defer L.Close()

	// err = L.DoCompiledFile(proto)
	// if err != nil {
	// 	fmt.Println(errors.WrapTrace(err))
	// }
	// characId := 1001

	// _, err = u.addCharacter(int32(characId), 1, u)

	// // // todo 计算角色战斗力
	// userPower, err := CalUserCombatPower(u)
	// if err != nil {
	// 	return nil, errors.WrapTrace(err)
	// }

	// fmt.Printf("===============taskId: , TEAMPOWER: %d==================\n", userPower)

	return taskInfos, nil
}

func (u *User) YggDispatchGuildGetInfo() ([]*pb.VOYggDispatchTaskState, error) {

	var taskInfos []*pb.VOYggDispatchTaskState
	if u.Yggdrasil.Dispatch.IfUpdateGuild {
		// 筛选出已经解锁了派遣任务的区域中，id最大的区域
		var generateTasks []int32
		var mostAdvancedArea int32
		var starTask []int32
		var numRange []int32

		for _, area := range u.Yggdrasil.Areas.Areas {
			areaCSV, err := manager.CSV.Yggdrasil.GetArea(area.AreaId)
			if err != nil {
				return nil, err
			}
			if areaCSV.DispatchUnlock <= area.FetchPrestige(u) { // 解锁
				if areaCSV.Id > mostAdvancedArea {
					mostAdvancedArea = areaCSV.Id
					starTask = areaCSV.GuildStar
					numRange = areaCSV.GuildNum
				}
			}
		}

		// 如果有已解锁公会派遣的区域可用
		if mostAdvancedArea > 0 {
			for _, star := range starTask {
				tasks, err := u.YggDispatchGenerateGuildTasks(mostAdvancedArea, star, numRange)
				if err != nil {
					return nil, errors.WrapTrace(err)
				}
				generateTasks = append(generateTasks, tasks...)
			}
		}

		for _, newTask := range generateTasks {
			taskCSV, err := manager.CSV.Yggdrasil.GetDispatchTask(newTask)
			if err != nil {
				return nil, errors.WrapTrace(err)
			}
			u.Yggdrasil.Dispatch.AddNewGuildTask(newTask, int64(taskCSV.CloseTime))
		}

		u.Yggdrasil.Dispatch.IfUpdateGuild = false
	}

	// 更新所有公会任务的状态
	u.Yggdrasil.Dispatch.YggDispatchUpdateGuildStates()
	for _, task := range u.Yggdrasil.Dispatch.GuildTasks {
		taskInfos = append(taskInfos, task.VOYggDispatchTaskState())
	}

	return taskInfos, nil
}

// 公会任务生成
func (u *User) YggDispatchGenerateGuildTasks(areaId int32, star int32, numRange []int32) ([]int32, error) {
	var taskPool []int32
	var result []int32

	taskCSVs, err := manager.CSV.Yggdrasil.GetSpecificTasks(areaId, static.YggdrasilDispatchTypeGuild, star)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	for _, taskCSV := range taskCSVs {
		task, ok := u.Yggdrasil.Dispatch.GuildTasks[taskCSV.Id]
		if ok { // 公会任务不会刷新只会增加和关闭
			delete(taskCSVs, task.Id)
		} else {
			taskPool = append(taskPool, taskCSV.Id)
		}
	}

	num := rand.RangeInt32(numRange[0], numRange[1]) // 随机出本次要增加的任务数量

	if len(taskPool) == 0 {
		return result, nil
	}

	// 从taskPool中随机抽取num个任务,若不足num个，则会全部抽取
	simpleShuffle := rand.Perm(len(taskPool))
	for i, index := range simpleShuffle {
		if i >= int(num) {
			break
		}
		result = append(result, taskPool[index])
	}

	return result, nil
}

// 根据区域id和星级新增任务，当前已存在的对应区域和星级的任务不参与新增
func (u *User) YggDispatchGenerateDailyTasks(areaId, star, num int32) ([]int32, error) {
	var taskPool []int32
	var result []int32

	// 根据区域和星级获取所有的任务
	taskCSVs, err := manager.CSV.Yggdrasil.GetSpecificTasks(areaId, static.YggdrasilDispatchTypeDaily, star)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	for _, taskCSV := range taskCSVs {
		_, ok := u.Yggdrasil.Dispatch.DailyTasks[taskCSV.Id]
		if ok { // 如果特定区域特定星级的任务已在map中，则直接占用新增名额
			num -= 1
		} else {
			taskPool = append(taskPool, taskCSV.Id) // 否则将任务加入任务池等待抽取
		}
	}

	// 从taskPool中随机抽取num个任务
	simpleShuffle := rand.Perm(len(taskPool))
	for i, index := range simpleShuffle {
		if i >= int(num) {
			break
		}
		result = append(result, taskPool[index])
	}

	return result, nil
}

func (u *User) YggdrasilDispatchDailyTaskInfo(taskId int32) ([]int32, error) {

	_, ok := u.Yggdrasil.Dispatch.DailyTasks[taskId]
	if !ok {
		return nil, errors.Swrapf(common.ErrYggdrasilDispatchTaskNotFound, taskId)
	}

	var characIds []int32

	for _, task := range u.Yggdrasil.Dispatch.DailyTasks {
		characIds = append(characIds, task.CharacterId...)
	}

	return characIds, nil
}

func (u *User) YggdrasilDispatchGuildTaskInfo(taskId int32) ([]int32, int32, error) {
	_, ok := u.Yggdrasil.Dispatch.GuildTasks[taskId]
	if !ok {
		return nil, 0, errors.Swrapf(common.ErrYggdrasilDispatchTaskNotFound, taskId)
	}

	var characIds []int32

	for _, task := range u.Yggdrasil.Dispatch.GuildTasks {
		characIds = append(characIds, task.CharacterId...)
	}

	// 获得csv
	taskCsv, err := manager.CSV.Yggdrasil.GetDispatchTask(taskId)
	if err != nil {
		return nil, 0, errors.WrapTrace(err)
	}

	return characIds, taskCsv.GuildCharacId, nil
}

func (u *User) YggdrasilDispatchDailyTaskBegin(taskId int32, userCharacIds []int32) (*pb.VOYggDispatchTaskState, error) {
	task, ok := u.Yggdrasil.Dispatch.DailyTasks[taskId]
	if !ok {
		return nil, errors.Swrapf(common.ErrYggdrasilDispatchTaskNotFound, taskId)
	}
	if !task.CheckStateReady() {
		return nil, errors.Swrapf(common.ErrYggdrasilDispatchTaskStateNotReadyForMission, taskId, task.TaskState)
	}

	// 计算必要条件和额外条件的解锁
	taskCSV, err := manager.CSV.Yggdrasil.GetDispatchTask(task.Id)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	if taskCSV.TeamSize != int32(len(userCharacIds)) {
		return nil, errors.Swrapf(common.ErrYggDispatchWrongTeamSizeForDispatch, task.Id)
	}

	// 如果必要条件不满足，报错
	if u.CheckYggDispatchConditions(taskCSV.NecessaryConditions, userCharacIds) != nil {
		return nil, errors.Swrapf(common.ErrYggdrasilDispatchNecessaryConditionNotSatisfied, taskId)
	}

	// 计算额外条件是否满足，满足则进行记录奖励的index或者在此就记录dropID
	for i, condition := range *taskCSV.ExtraConditions {
		if u.CheckYggDispatchCondition(&condition, userCharacIds) == nil {
			task.ExtraRewards = append(task.ExtraRewards, int32(i))
		}
	}

	// lua.MaxArrayIndex = 4000
	// proto, err := CompileLua("E:/code/go/src/overlord/gamesvr/model/CombatPowerUtil.lua")
	// if err != nil {
	// 	fmt.Println(common.ErrCompileLuaFileFailed)
	// }

	// L := NewLuaMatcher()
	// defer L.Close()

	// err = L.DoCompiledFile(proto)
	// if err != nil {
	// 	fmt.Println(errors.WrapTrace(err))
	// }

	// var teamPower int32

	// 在map中把这些角色状态改变，如果没找到这个角色，报错
	// 侍从策划确定了只会增加不会减少
	for _, characId := range userCharacIds {
		_, err := u.CharacterPack.Get(characId)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		// if stateId != 0 {
		// 	return nil, errors.Swrapf(common.ErrYggdrasilDispatchCharacterIsOnMission, characId)
		// }
		if u.Yggdrasil.Dispatch.FindCharacterDailyOnMission(characId) {
			return nil, errors.Swrapf(common.ErrYggdrasilDispatchCharacterIsOnMission, characId)
		}
		// u.Yggdrasil.Dispatch.CharacStateForDaily[characId] = taskId
	}

	// 修改任务的状态
	task.StateReadyToGo(taskCSV.TimeCost, userCharacIds)

	u.Yggdrasil.Dispatch.DailyTasks[taskId] = task

	// 返回任务的状态信息
	return task.VOYggDispatchTaskState(), nil
}

// todo 传入参数还有一个公会角色的id, 只会有一个公会角色
func (u *User) YggdrasilDispatchGuildTaskBegin(taskId int32, userCharacIds []int32, guildCustomIds []int32) (*pb.VOYggDispatchTaskState, error) {
	task, ok := u.Yggdrasil.Dispatch.GuildTasks[taskId]
	if !ok {
		return nil, errors.Swrapf(common.ErrYggdrasilDispatchTaskNotFound, taskId)
	}
	if !task.CheckStateReady() {
		return nil, errors.Swrapf(common.ErrYggdrasilDispatchTaskStateNotReadyForMission, taskId, task.TaskState)
	}

	// todo 检查是否有重复角色 .......公会角色和玩家自身角色可以重复

	// 计算必要条件和额外条件的解锁
	taskCSV, err := manager.CSV.Yggdrasil.GetDispatchTask(task.Id)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	if taskCSV.TeamSize != int32(len(userCharacIds)+len(guildCustomIds)) {
		return nil, errors.Swrapf(common.ErrYggDispatchWrongTeamSizeForDispatch, task.Id)
	}

	// 如果必要条件不满足，报错
	if u.CheckYggDispatchConditions(taskCSV.NecessaryConditions, userCharacIds) != nil {
		return nil, errors.Swrapf(common.ErrYggdrasilDispatchNecessaryConditionNotSatisfied, taskId)
	}

	// 计算额外条件是否满足，满足则进行记录奖励的index或者在此就记录dropID
	for i, condition := range *taskCSV.ExtraConditions {
		if u.CheckYggDispatchCondition(&condition, userCharacIds) == nil {
			task.ExtraRewards = append(task.ExtraRewards, int32(i))
		}
	}

	for _, characId := range userCharacIds {
		_, err := u.CharacterPack.Get(characId)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		if u.Yggdrasil.Dispatch.FindCharacterDailyOnMission(characId) {
			return nil, errors.Swrapf(common.ErrYggdrasilDispatchCharacterIsOnMission, characId)
		}
	}

	// 修改任务的状态
	task.StateReadyToGo(taskCSV.TimeCost, userCharacIds)

	u.Yggdrasil.Dispatch.GuildTasks[taskId] = task
	return nil, nil
}

func (u *User) YggdrasilDispatchReward(taskId int32) (*pb.VOYggDispatchResourceResult, error) {
	taskCSV, err := manager.CSV.Yggdrasil.GetDispatchTask(taskId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	voResult := &pb.VOYggDispatchResourceResult{}
	if taskCSV.Type == static.YggdrasilDispatchTypeDaily {
		voResult, err = u.YggdrasilDispatchDailyRewards(taskCSV)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
	} else {
		// todo 公会
	}

	return voResult, nil
}

func (u *User) YggdrasilDispatchDailyRewards(taskCSV *entry.CfgYggdrasilDispatch) (*pb.VOYggDispatchResourceResult, error) {
	task, ok := u.Yggdrasil.Dispatch.DailyTasks[taskCSV.Id]
	if !ok {
		return nil, errors.Swrapf(common.ErrYggdrasilDispatchTaskNotFound, taskCSV.Id)
	}

	if !task.CheckStateMissionComplete() {
		return nil, errors.Swrapf(common.ErrYggdrasilDispatchTaskStateNotReadyForReceiveRewards, task.Id)
	}

	reason := logreason.NewReason(logreason.YggDailyDispatch)
	_, err := u.AddRewardsByDropId(taskCSV.BaseRewards, reason)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	baseResult := u.VOResourceResult()

	extraDropIds := make([]int32, 0, len(task.ExtraRewards))
	for _, index := range task.ExtraRewards {
		extraDropIds = append(extraDropIds, taskCSV.ExtraRewards[index])
	}

	_, err = u.AddRewardsByDropIds(extraDropIds, reason)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	extraResult := u.VOResourceResult()

	task.NextState()

	return &pb.VOYggDispatchResourceResult{
		BasicResult: baseResult,
		ExtraResult: extraResult,
	}, nil
}

func (u *User) YggdrasilDispatchGuildRewards(taskCSV *entry.CfgYggdrasilDispatch) (*pb.VOYggDispatchResourceResult, error) {
	task, ok := u.Yggdrasil.Dispatch.GuildTasks[taskCSV.Id]
	if !ok {
		return nil, errors.Swrapf(common.ErrYggdrasilDispatchTaskNotFound, taskCSV.Id)
	}

	if !task.CheckStateMissionComplete() {
		return nil, errors.Swrapf(common.ErrYggdrasilDispatchTaskStateNotReadyForReceiveRewards, task.Id)
	}

	reason := logreason.NewReason(logreason.YggGuildDispatch)
	_, err := u.AddRewardsByDropId(taskCSV.BaseRewards, reason)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	baseResult := u.VOResourceResult()

	extraDropIds := make([]int32, 0, len(taskCSV.ExtraRewards))
	for _, index := range task.ExtraRewards {
		extraDropIds = append(extraDropIds, taskCSV.ExtraRewards[index])
	}

	_, err = u.AddRewardsByDropIds(extraDropIds, reason)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	extraResult := u.VOResourceResult()

	task.NextState()

	return &pb.VOYggDispatchResourceResult{
		BasicResult: baseResult,
		ExtraResult: extraResult,
	}, nil
}

func (u *User) YggDispatchCancel(ctx context.Context, taskId int32) (*pb.VOYggDispatchTaskState, error) {
	taskCSV, err := manager.CSV.Yggdrasil.GetDispatchTask(taskId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	voTaskInfo, err := u.Yggdrasil.Dispatch.CancelDispatch(taskId, taskCSV.Type)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return voTaskInfo, nil
}

// ------------------派遣条件解锁-----------------

// 如果满足解锁条件，返回nil
func (u *User) CheckYggDispatchConditions(conditions *common.Conditions, userCharacIds []int32) error {
	if conditions == nil {
		return nil
	}

	for _, condition := range *conditions {
		err := u.CheckYggDispatchCondition(&condition, userCharacIds)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *User) CheckYggDispatchCondition(condition *common.Condition, userCharacIds []int32) error {
	switch condition.ConditionType {
	case static.YggdrasilConditionTypeYggDispatchLevel:
		targetNum := condition.Params[0]
		level := condition.Params[1]

		var num int32
		for _, charaId := range userCharacIds {
			chara, err := u.CharacterPack.Get(charaId)
			if err != nil {
				return err
			}

			if chara.GetLevel() >= level {
				num += 1
			}
		}

		if targetNum > num {
			return errors.Swrapf(common.ErrYggdrasilDispatchCharacterLevelNotArrival)
		}

	case static.YggdrasilConditionTypeYggDispatchCharaCamp:
		targetNum := condition.Params[0]
		campId := condition.Params[1]

		var num int32
		for _, charaId := range userCharacIds {
			chara, err := u.CharacterPack.Get(charaId)
			if err != nil {
				return errors.WrapTrace(err)
			}

			camp, err := manager.CSV.Character.Camp(chara.ID)
			if err != nil {
				return errors.WrapTrace(err)
			}

			if camp == campId {
				num += 1
			}
		}

		if targetNum > num {
			return errors.Swrapf(common.ErrYggdrasilDispatchCharacterCampNotArrival)
		}

	case static.YggdrasilConditionTypeYggDispatchCharaCareer:
		targetNum := condition.Params[0]
		careerId := condition.Params[1]

		var num int32
		for _, charaId := range userCharacIds {
			chara, err := u.CharacterPack.Get(charaId)
			if err != nil {
				return errors.WrapTrace(err)
			}

			career, err := manager.CSV.Character.Career(chara.ID)
			if err != nil {
				return errors.WrapTrace(err)
			}
			if career == careerId {
				num += 1
			}
		}

		if targetNum > num {
			return errors.Swrapf(common.ErrYggdrasilDispatchCharacterCareerNotArrival)
		}

	case static.YggdrasilConditionTypeYggDispatchCharaRarity:
		targetNum := condition.Params[0]
		targetRarity := condition.Params[1]

		var num int32
		for _, charaId := range userCharacIds {
			chara, err := u.CharacterPack.Get(charaId)
			if err != nil {
				return errors.WrapTrace(err)
			}
			rarity := chara.GetRare()

			if rarity >= targetRarity { // 达到目标稀有度，数量加1
				num += 1
			}
		}

		if targetNum > num { // 如果数量并未达到目标数量，则并未达成条件
			return errors.Swrapf(common.ErrYggdrasilDispatchCharacterRarityNotArrival)
		}
	case static.YggdrasilConditionTypeYggDispatchCharaStar:
		targetNum := condition.Params[0]
		targetStar := condition.Params[1]

		var num int32
		for _, charaId := range userCharacIds {
			chara, err := u.CharacterPack.Get(charaId)
			if err != nil {
				return errors.WrapTrace(err)
			}

			if chara.GetStar() >= targetStar { // 达到目标稀有度，数量加1
				num += 1
			}
		}

		if targetNum > num { // 如果数量并未达到目标数量，则并未达成条件
			return errors.Swrapf(common.ErrYggdrasilDispatchCharacterStarNotArrival)
		}
	case static.YggdrasilConditionTypeYggDispatchPower:
		// todo 后续战力为实时的话，这里不需要计算，直接取角色战力相加即可
		powerTarget := condition.Params[0]
		var powerSum int32
		for _, charaId := range userCharacIds {
			chara, err := u.CharacterPack.Get(charaId)
			if err != nil {
				return errors.WrapTrace(err)
			}
			powerSum += chara.Power
		}
		if powerTarget > powerSum {
			return errors.Swrapf(common.ErrYggdrasilDispatchCharacterPowerNotArrival)
		}

	case static.YggdrasilConditionTypeYggDispatchSpecificChar:
		charaId := condition.Params[0]
		level := condition.Params[1]
		stage := condition.Params[2]
		star := condition.Params[3]

		ifContain := false
		for _, cid := range userCharacIds {
			if cid == charaId {
				ifContain = true
			}
		}

		if !ifContain {
			return errors.Swrapf(common.ErrYggdrasilDispatchSpecificCharacterNotArrival)
		}
		chara, err := u.CharacterPack.Get(charaId)
		if err != nil {
			return errors.WrapTrace(err)
		}

		if level > chara.GetLevel() || stage > chara.GetStage() || star > chara.GetStar() {
			return errors.Swrapf(common.ErrYggdrasilDispatchSpecificCharacterNotArrival)
		}

	case static.YggdrasilConditionTypeYggDispatchCamp:
		campId := condition.Params[0]
		for _, charaId := range userCharacIds {
			chara, err := u.CharacterPack.Get(charaId)
			if err != nil {
				return err
			}

			camp, err := manager.CSV.Character.Camp(chara.ID)
			if err != nil {
				return errors.WrapTrace(err)
			}

			if camp != campId {
				return errors.Swrapf(common.ErrYggdrasilDispatchCampNotArrival)
			}
		}

	case static.YggdrasilConditionTypeYggDispatchCareer:
		careerId := condition.Params[0]
		for _, charaId := range userCharacIds {
			chara, err := u.CharacterPack.Get(charaId)
			if err != nil {
				return err
			}

			career, err := manager.CSV.Character.Career(chara.ID)
			if err != nil {
				return errors.WrapTrace(err)
			}

			if career != careerId {
				return errors.Swrapf(common.ErrYggdrasilDispatchCareerNotArrival)
			}
		}
	}

	return nil
}
