package model

import (
	"context"
	"gamesvr/manager"
	"shared/common"
	"shared/protobuf/pb"
	"shared/statistic/logreason"
	"shared/utility/errors"
	"shared/utility/servertime"
	"strings"
	"time"

	"strconv"
)

const (
	GMTypeAddReward           = 1 //添加物品, args: [id, count]
	GMTypeGraveyardBuildLvSet = 2 // 设置模拟经营建筑等级， args:[buildUid,lv]

	//添加gm邮件
	//	args: [templateId, dropId(附件，可以不填), title(空字符直接读表), titleArgs(逗号间隔),
	//			content(空字符直接读表), contentArgs(逗号间隔)]
	GMTypeAddMail = 3

	GMTypePassLevel  = 4 //通过副本(可以解锁其他系统), args: [levelId, systemType, param, param, ...]
	GMTypePassLevels = 5 //通过多个副本, args: [levelId, levelId, ...]
	GMTypeSetTime    = 6 //设置时间, args:["2021-12-06 03:59:00"]
)

const (
	addRewardParamCount = 2
	addMailParamCount   = 2
)

func (u *User) ExecuteGM(ctx context.Context, gmType int32, params []string) (*pb.S2CGM, error) {
	ret := &pb.S2CGM{}
	var err error

	switch gmType {
	case GMTypeAddReward:
		err = u.gmAddReward(params)
	case GMTypeGraveyardBuildLvSet:
		voBuild, err := u.setBuildLvByGm(params)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		ret.VoBuild = voBuild
	case GMTypeAddMail:
		err = u.gmAddMail(params)
	case GMTypePassLevel:
		err = u.gmPassLevel(ctx, params)
	case GMTypePassLevels:
		levels, err := u.gmPassLevels(ctx, params)
		if err == nil {
			ret.Levels = []*pb.VOLevel{}
			for _, level := range levels {
				ret.Levels = append(ret.Levels, level.VOLevel())
			}
		}
	case GMTypeSetTime:
		err = u.gmSetTime(params)
	default:
		err = common.ErrParamError
	}

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (u *User) gmAddReward(params []string) error {
	if len(params) != addRewardParamCount {
		return common.ErrParamError
	}

	id, err := strconv.ParseInt(params[0], 10, 32)
	if err != nil {
		return err
	}

	count, err := strconv.ParseInt(params[1], 10, 32)
	if err != nil {
		return err
	}

	rewards := common.NewRewards()
	rewards.AddReward(common.NewReward(int32(id), int32(count)))

	reason := logreason.NewReason(logreason.GMAddReward)
	_, err = u.addRewards(rewards, reason)
	if err != nil {
		return errors.WrapTrace(err)
	}

	return nil
}

func (u *User) gmAddMail(params []string) error {
	if len(params) < addMailParamCount {
		return common.ErrParamError
	}

	templateId, err := strconv.ParseInt(params[0], 10, 32)
	if err != nil {
		return err
	}

	dropId64 := int64(0)
	if params[1] != "" {
		dropId64, err = strconv.ParseInt(params[1], 10, 32)
		if err != nil {
			return err
		}
	}

	title := ""
	if len(params) >= 3 {
		title = params[2]
	}

	titleArgs := []string{}
	if len(params) >= 4 && params[5] != "" {
		titleArgs = strings.Split(params[3], ",")
	}

	content := ""
	if len(params) >= 5 {
		content = params[4]
	}

	contentArgs := []string{}
	if len(params) >= 6 && params[5] != "" {
		contentArgs = strings.Split(params[5], ",")
	}

	var rewards *common.Rewards = nil

	dropId := int32(dropId64)
	if dropId != 0 {
		rewards, err = manager.CSV.Drop.DropRewards(dropId)
		if err != nil {
			return err
		}
	}

	templateCfg, err := manager.CSV.Mail.GetTemplate(int32(templateId))
	if err != nil {
		return err
	}

	sendTime := servertime.Now().Unix()
	expireTime := sendTime + int64(templateCfg.ExpireDays)*servertime.SecondPerDay

	u.AddMail(int32(templateId), title, titleArgs, content, contentArgs,
		rewards, "", sendTime, expireTime)

	return nil
}

func (u *User) gmPassLevel(ctx context.Context, params []string) error {
	if len(params) < 1 {
		return common.ErrParamError
	}

	levelId64, err := strconv.ParseInt(params[0], 10, 32)
	if err != nil {
		return err
	}

	levelId := int32(levelId64)

	levelCfg, err := manager.CSV.LevelsEntry.GetLevel(levelId)
	if err != nil {
		return err
	}

	systemType := int32(0)
	var systemParams []int64
	if len(params) > 2 {
		systemType64, err := strconv.ParseInt(params[1], 10, 32)
		if err != nil {
			return err
		}

		systemType = int32(systemType64)

		for i := 2; i < len(params); i++ {
			param, err := strconv.ParseInt(params[i], 10, 32)
			if err != nil {
				return err
			}

			systemParams = append(systemParams, param)
		}
	}

	u.PassLevel(ctx, levelCfg, []int32{0, 1, 2}, []int32{0, 1, 2}, systemType, systemParams,
		[]int32{}, nil, logreason.EmptyReason())

	return nil
}

func (u *User) gmPassLevels(ctx context.Context, params []string) ([]*Level, error) {
	if len(params) < 1 {
		return nil, common.ErrParamError
	}

	levels := []*Level{}

	for _, levelIdStr := range params {
		levelId64, err := strconv.ParseInt(levelIdStr, 10, 32)
		if err != nil {
			return nil, err
		}

		levelId := int32(levelId64)
		levelCfg, err := manager.CSV.LevelsEntry.GetLevel(levelId)
		if err != nil {
			return nil, err
		}

		u.PassLevel(ctx, levelCfg, []int32{}, []int32{}, 0, []int64{}, []int32{},
			nil, logreason.EmptyReason())

		level := u.LevelsInfo.GetOrCreateLevel(levelId)
		levels = append(levels, level)
	}

	return levels, nil
}

func (u *User) gmSetTime(params []string) error {
	if len(params) < 1 {
		return common.ErrParamError
	}

	timeStr := params[0]

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
