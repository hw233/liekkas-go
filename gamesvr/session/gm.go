package session

import (
	"context"
	"reflect"
	"strings"
	"time"

	"gamesvr/manager"
	"shared/common"
	"shared/csv/static"
	"shared/statistic/logreason"
	"shared/utility/errors"
	"shared/utility/glog"
	"shared/utility/param"
	"shared/utility/servertime"
)

func (s *Session) isGMCode(code string) bool {
	return strings.HasPrefix(code, "#GM:")
}

func (s *Session) handelGMCode(ctx context.Context, code string) error {
	// 解析命令码
	// "#GM:addItem(1,2);addItem(1,3)"
	gmCodes := strings.Split(strings.TrimPrefix(code, "#GM:"), ";")
	// ["addItem(1,2)", "addItem(1,3)"]
	gmTasks := make([]common.GMTask, 0, len(gmCodes))
	for _, gmCode := range gmCodes {
		splitGmCode := strings.Split(gmCode, "(")
		if len(splitGmCode) != 2 {
			return errors.New("code fmt error")
		}

		// ["addItem", "1,2)"]

		// 解析参数["1", "2"]
		params := strings.Split(strings.TrimSuffix(splitGmCode[1], ")"), ",")

		glog.Infof("cmd: %v", "gm"+strings.Title(splitGmCode[0]))

		// 搜索函数名：gmAddItem
		method, ok := reflect.TypeOf(s).MethodByName("GM" + strings.Title(splitGmCode[0]))
		if !ok {
			return errors.New("not found gm")
		}

		gmTasks = append(gmTasks, *common.NewGMTask(method.Func, s, ctx, param.NewParam(params)))
	}

	// 按顺序执行任务
	for _, gmTask := range gmTasks {
		err := gmTask.Do()
		if err != nil {
			return errors.WrapTrace(err)
		}
	}

	return nil
}

func (s *Session) GMConfigTask(ctx context.Context, param *param.Param) error {
	if param.Len() != 1 {
		return common.ErrParamError
	}

	id, err := param.GetInt32(0)
	if err != nil {
		return errors.WrapTrace(err)
	}

	code, err := manager.CSV.GM.GetGMCode(id)
	if err != nil {
		return errors.WrapTrace(err)
	}

	if s.isGMCode(code) {
		return s.handelGMCode(ctx, code)
	}

	return errors.New("gm fmt error")
}

func (s *Session) GMAddItem(ctx context.Context, param *param.Param) error {
	if param.Len() != 2 {
		return common.ErrParamError
	}

	id, err := param.GetInt32(0)
	if err != nil {
		return errors.WrapTrace(err)
	}

	count, err := param.GetInt32(1)
	if err != nil {
		return errors.WrapTrace(err)
	}

	rewards := common.NewRewards()
	rewards.AddReward(common.NewReward(id, count))

	reason := logreason.NewReason(logreason.GMAddReward)
	_, err = s.AddRewards(rewards, reason)
	if err != nil {
		return errors.WrapTrace(err)
	}

	return nil
}

func (s *Session) GMAddEquipmentEXP(ctx context.Context, param *param.Param) error {
	if param.Len() != 2 {
		return common.ErrParamError
	}

	id, err := param.GetInt64(0)
	if err != nil {
		return errors.WrapTrace(err)
	}

	count, err := param.GetInt32(1)
	if err != nil {
		return errors.WrapTrace(err)
	}

	equipment, err := s.User.EquipmentPack.Get(id)
	if err != nil {
		return errors.WrapTrace(err)
	}

	equipment.EXP.Plus(count)

	err = manager.CSV.Equipment.SyncLevelAndUnlockAttr(equipment)
	if err != nil {
		return errors.WrapTrace(err)
	}

	return nil
}

func (s *Session) GMAddMail(ctx context.Context, params *param.Param) error {
	if params.Len() != 2 {
		return common.ErrParamError
	}

	templateId, err := params.GetInt32(0)
	if err != nil {
		return err
	}

	templateCfg, err := manager.CSV.Mail.GetTemplate(int32(templateId))
	if err != nil {
		return err
	}

	var rewards *common.Rewards = nil
	dropId, err := params.GetInt32(1)
	if err != nil {
		return err
	}

	if dropId != 0 {
		rewards, err = manager.CSV.Drop.DropRewards(dropId)
		if err != nil {
			return err
		}
	}

	sendTime := servertime.Now().Unix()
	expireTime := sendTime + int64(templateCfg.ExpireDays)*servertime.SecondPerDay

	s.User.AddMail(templateId, "", []string{}, "", []string{},
		rewards, "", sendTime, expireTime)

	return nil
}

func (s *Session) GMPassLevel(ctx context.Context, params *param.Param) error {
	if params.Len() < 1 {
		return common.ErrParamError
	}

	levelId, err := params.GetInt32(0)
	if err != nil {
		return err
	}

	levelCfg, err := manager.CSV.LevelsEntry.GetLevel(levelId)
	if err != nil {
		return err
	}

	passAchievement := []int32{}
	for i := int32(1); i <= levelCfg.AchievementsCount; i++ {
		passAchievement = append(passAchievement, i)
	}

	passTarget := []int32{}
	for i := int32(1); i <= levelCfg.TargetCount; i++ {
		passTarget = append(passTarget, i)
	}

	s.User.PassLevel(ctx, levelCfg, passTarget, passAchievement, 0,
		[]int64{}, []int32{}, nil, logreason.EmptyReason())

	return nil
}

func (s *Session) GMResetLevel(ctx context.Context, params *param.Param) error {
	if params.Len() < 1 {
		return common.ErrParamError
	}

	levelId, err := params.GetInt32(0)
	if err != nil {
		return err
	}

	level, ok := s.User.LevelsInfo.GetLevel(levelId)

	if !ok {
		return errors.Swrapf(common.ErrLevelNotPassed, levelId)
	}

	level.Reset()

	return nil
}

func (s *Session) GMPassTower(ctx context.Context, params *param.Param) error {
	if params.Len() < 2 {
		return common.ErrParamError
	}

	towerId, err := params.GetInt32(0)
	if err != nil {
		return err
	}

	towerStage, err := params.GetInt32(1)
	if err != nil {
		return err
	}

	towerStageCfg, err := manager.CSV.Tower.GetTowerStage(towerId, towerStage)
	if err != nil {
		return err
	}

	levelCfg, _ := manager.CSV.LevelsEntry.GetLevel(towerStageCfg.LevelId)

	passAchievement := []int32{}
	for i := int32(1); i <= levelCfg.AchievementsCount; i++ {
		passAchievement = append(passAchievement, i)
	}

	passTarget := []int32{}
	for i := int32(1); i <= levelCfg.TargetCount; i++ {
		passTarget = append(passTarget, i)
	}

	s.User.PassLevel(ctx, levelCfg, passTarget, passAchievement, static.BattleTypeTower,
		[]int64{int64(towerId), int64(towerStage)}, []int32{}, nil, logreason.EmptyReason())

	return nil
}

func (s *Session) GMSetTime(ctx context.Context, params *param.Param) error {
	if params.Len() < 1 {
		return common.ErrParamError
	}

	timeStr, err := params.GetString(0)
	if err != nil {
		return err
	}

	timeStr = strings.Trim(timeStr, "\"")
	t, err := servertime.ParseTime(timeStr)
	if err != nil {
		return err
	}

	timeOffset := t - servertime.OriginNow().Unix()
	servertime.SetTimeOffset(timeOffset)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	manager.Global.String.Set(ctx, servertime.TimeOffsetRedisName, timeOffset)

	return nil
}

func (s *Session) GMResetSetTime(ctx context.Context, params *param.Param) error {
	servertime.SetTimeOffset(0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	manager.Global.String.Set(ctx, servertime.TimeOffsetRedisName, 0)

	return nil
}

func (s *Session) GMPassGuide(ctx context.Context, params *param.Param) error {
	if params.Len() != 1 {
		return errors.WrapTrace(common.ErrParamError)
	}

	guideId, err := params.GetInt32(0)
	if err != nil {
		return errors.WrapTrace(err)
	}
	config, err := manager.CSV.Guide.GetConfigById(guideId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	before := manager.CSV.Guide.GetGuideBefore(config.GuideOrder)

	reason := logreason.NewReason(logreason.GMPassGuide)
	for _, guideConfig := range before {
		if s.User.Info.PassedGuideIds.Contains(guideConfig.Id) {
			continue
		}
		if config.DropId > 0 {
			_, err = s.User.AddRewardsByDropId(config.DropId, reason)
			if err != nil {
				return errors.WrapTrace(err)
			}
		}
		s.User.Info.PassedGuideIds.Append(guideConfig.Id)

	}
	return nil
}

func (s *Session) GMAddYggdrasilMail(ctx context.Context, params *param.Param) error {
	if params.Len() != 2 {
		return errors.WrapTrace(common.ErrParamError)
	}

	num, err := params.GetInt32(0)
	if err != nil {
		return errors.WrapTrace(err)
	}
	dropId, err := params.GetInt32(1)
	if err != nil {
		return errors.WrapTrace(err)
	}
	for i := 0; i < int(num); i++ {
		rewards, err := manager.CSV.Drop.DropRewards(dropId)
		if err != nil {
			return errors.WrapTrace(err)
		}
		err = s.User.Yggdrasil.MailBox.AddOne(ctx, s.ID, s.Name, rewards)
		if err != nil {
			return errors.WrapTrace(err)
		}
	}

	return nil
}

func (s *Session) GMUnlockCity(ctx context.Context, params *param.Param) error {

	cityId, err := params.GetInt32(0)
	if err != nil {
		return errors.WrapTrace(err)
	}

	_, err = manager.CSV.Yggdrasil.GetYggCityById(cityId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	s.Yggdrasil.CanTravelCityIds.Append(cityId)
	return nil

}

func (s *Session) GMCompleteTracedSubTask(ctx context.Context, params *param.Param) error {
	taskId := s.Yggdrasil.Task.FetchTrackTaskId()
	if taskId == 0 {
		return nil
	}
	info := s.Yggdrasil.Task.TaskInProgress[taskId]

	for _, multi := range info.Multi {
		err := s.Yggdrasil.Task.ForceCompleteSubTask(ctx, s.User, multi, info)
		if err != nil {
			return errors.WrapTrace(err)
		}
	}
	err := s.Yggdrasil.Task.ForceCompleteSubTask(ctx, s.User, info.Base, info)
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}
func (s *Session) GMAcceptTask(ctx context.Context, params *param.Param) error {
	taskId, err := params.GetInt32(0)
	if err != nil {
		return errors.WrapTrace(err)
	}
	config, err := manager.CSV.Yggdrasil.GetTaskConfig(taskId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	_, err = s.Yggdrasil.Task.InitSubTask(ctx, s.User, config.NextSubTaskId)
	if err != nil {
		return err
	}
	return nil

}
func (s *Session) GMCompleteTask(ctx context.Context, params *param.Param) error {
	taskId, err := params.GetInt32(0)
	if err != nil {
		return errors.WrapTrace(err)
	}
	return s.Yggdrasil.Task.ForceCompleteTask(ctx, s.User, taskId)
}

func (s *Session) GMAddTravelAp(ctx context.Context, params *param.Param) error {
	add, err := params.GetInt32(0)
	if err != nil {
		return errors.WrapTrace(err)
	}
	s.Yggdrasil.TravelInfo.TravelAp += add

	return nil
}

func (s *Session) GMAddItemInYgg(ctx context.Context, params *param.Param) error {
	if params.Len() != 2 {
		return common.ErrParamError
	}
	id, err := params.GetInt32(0)
	if err != nil {
		return errors.WrapTrace(err)
	}

	count, err := params.GetInt32(1)
	if err != nil {
		return errors.WrapTrace(err)
	}

	rewards := common.NewRewards()
	rewards.AddReward(common.NewReward(id, count))

	reason := logreason.NewReason(logreason.GMAddReward)
	err = s.Yggdrasil.AddRewards(ctx, s.User, rewards, 0, reason)
	if err != nil {
		return errors.WrapTrace(err)
	}

	return nil
}

func (s *Session) GMYggLightAll(ctx context.Context, params *param.Param) error {

	for _, config := range manager.CSV.Yggdrasil.GetAllAreaConfig() {
		for _, p := range config.Area.Points() {
			s.Yggdrasil.UnlockArea.AppendPoint(p)
		}
	}

	return nil
}
func (s *Session) GMYggClearMatchPool(ctx context.Context, params *param.Param) error {
	err := manager.Global.SDel(ctx, "matchUserIds")
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}

func (s *Session) GMYggResetExploreTimes(ctx context.Context, params *param.Param) error {
	s.Yggdrasil.TravelTime = 0
	return nil
}

func (s *Session) GMYggResetCharacterTime(ctx context.Context, params *param.Param) error {
	for _, character := range *s.CharacterPack {
		character.CanYggdrasilTime = servertime.Now().Unix()
	}
	return nil
}
