package session

import (
	"context"
	"gamesvr/manager"
	"shared/common"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/statistic/logreason"
	"shared/utility/coordinate"
	"shared/utility/errors"
	"shared/utility/glog"
)

//TryPushGraveyard 因为Graveyard公会互助，别人触发自己的推送，所以只放在timmerpush
func (s *Session) TryPushGraveyard() {
	for _, push := range s.User.Notifies.GraveyardPush {
		cmd, err := manager.CSV.Protocol.GetCmdByProtoName(push)
		if err != nil {
			glog.Errorf("TryPushGraveyard  GetCmdByProtoName err:%+v", err)
			continue
		}
		err = s.push(cmd, push)
		if err != nil {
			glog.Errorf("TryPushGraveyard push err:%+v", err)
			continue
		}
	}
	s.User.Notifies.GraveyardPush = nil

}

// GraveyardGetInfo 进玩法
func (s *Session) GraveyardGetInfo(ctx context.Context, req *pb.C2SGraveyardGetInfo) (*pb.S2CGraveyardGetInfo, error) {
	return s.User.GraveyardGetInfo(), nil
}

//--------------------------建造相关

// GraveyardBuildCreate 建造
func (s *Session) GraveyardBuildCreate(ctx context.Context, req *pb.C2SGraveyardBuildCreate) (*pb.S2CGraveyardBuildCreate, error) {
	// 检查模拟经营是否解锁
	err := s.User.CheckActionUnlock(static.ActionIdTypeTycoonunlock)
	if err != nil {
		return nil, err
	}

	buildId := req.GetBuildId()
	if buildId <= 0 {
		return nil, errors.WrapTrace(common.ErrParamError)
	}
	if req.Position == nil {
		return nil, errors.WrapTrace(common.ErrParamError)
	}
	position := coordinate.NewPositionFromVo(req.GetPosition())
	vo, err := s.User.GraveyardBuildCreate(buildId, position)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CGraveyardBuildCreate{
		Build:          vo,
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

// GraveyardRelocation 移动或者收到背包
func (s *Session) GraveyardRelocation(ctx context.Context, req *pb.C2SGraveyardRelocation) (*pb.S2CGraveyardRelocation, error) {
	// 检查模拟经营是否解锁
	err := s.User.CheckActionUnlock(static.ActionIdTypeTycoonunlock)
	if err != nil {
		return nil, err
	}

	// uid和要移动到的position
	relocationMap := map[int64]*coordinate.Position{}
	// 参数检查
	for _, relocation := range req.Relocations {
		_, ok := relocationMap[relocation.BuildUid]
		if ok {
			return nil, errors.WrapTrace(common.ErrParamError)
		}
		if relocation.Position == nil {
			relocationMap[relocation.BuildUid] = nil
		} else {
			position := coordinate.NewPositionFromVo(relocation.Position)
			relocationMap[relocation.BuildUid] = position
		}
	}
	if len(relocationMap) == 0 {
		return nil, errors.WrapTrace(common.ErrParamError)
	}

	return s.User.GraveyardRelocation(relocationMap)
}

// GraveyardBuildLvUp 升级
func (s *Session) GraveyardBuildLvUp(ctx context.Context, req *pb.C2SGraveyardBuildLvUp) (*pb.S2CGraveyardBuildLvUp, error) {
	// 检查模拟经营是否解锁
	err := s.User.CheckActionUnlock(static.ActionIdTypeTycoonunlock)
	if err != nil {
		return nil, err
	}

	vo, err := s.User.GraveyardBuildLvUp(req.BuildUid)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CGraveyardBuildLvUp{
		Build:          vo,
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

// GraveyardBuildStageUp 升阶
func (s *Session) GraveyardBuildStageUp(ctx context.Context, req *pb.C2SGraveyardBuildStageUp) (*pb.S2CGraveyardBuildStageUp, error) {
	// 检查模拟经营是否解锁
	err := s.User.CheckActionUnlock(static.ActionIdTypeTycoonunlock)
	if err != nil {
		return nil, err
	}

	vo, err := s.User.GraveyardBuildStageUp(req.BuildUid)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CGraveyardBuildStageUp{
		Build:          vo,
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

// GraveyardOpenCurtain 建筑揭幕
func (s *Session) GraveyardOpenCurtain(ctx context.Context, req *pb.C2SGraveyardOpenCurtain) (*pb.S2CGraveyardOpenCurtain, error) {
	// 检查模拟经营是否解锁
	err := s.User.CheckActionUnlock(static.ActionIdTypeTycoonunlock)
	if err != nil {
		return nil, err
	}

	vo, err := s.User.GraveyardOpenCurtain(req.BuildUid)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CGraveyardOpenCurtain{
		Build:          vo,
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

//--------------------------产出相关

// GraveyardProduceStart 非持续性建筑开始产出
func (s *Session) GraveyardProduceStart(ctx context.Context, req *pb.C2SGraveyardProduceStart) (*pb.S2CGraveyardProduceStart, error) {
	// 检查模拟经营是否解锁
	err := s.User.CheckActionUnlock(static.ActionIdTypeTycoonunlock)
	if err != nil {
		return nil, err
	}

	vos, err := s.User.GraveyardProduceStart(req.BuildUid, req.ProduceNum)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CGraveyardProduceStart{
		Builds:         vos,
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

// GraveyardProductionGet 获得产出
func (s *Session) GraveyardProductionGet(ctx context.Context, req *pb.C2SGraveyardProductionGet) (*pb.S2CGraveyardProductionGet, error) {
	// 检查模拟经营是否解锁
	err := s.User.CheckActionUnlock(static.ActionIdTypeTycoonunlock)
	if err != nil {
		return nil, err
	}

	voList, err := s.User.GraveyardProductionGet(req.BuildUidList)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CGraveyardProductionGet{
		Builds:         voList,
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

// GraveyardRefreshBuildInfo 刷新产出，不收获
func (s *Session) GraveyardRefreshBuildInfo(ctx context.Context, req *pb.C2SGraveyardRefreshBuildInfo) (*pb.S2CGraveyardRefreshBuildInfo, error) {
	// 检查模拟经营是否解锁
	err := s.User.CheckActionUnlock(static.ActionIdTypeTycoonunlock)
	if err != nil {
		return nil, err
	}

	voList, err := s.User.GraveyardRefreshBuildInfo(req.BuildUidList)
	if err != nil {
		return nil, err
	}
	return &pb.S2CGraveyardRefreshBuildInfo{
		BuildInfoList: voList,
	}, nil
}

// GraveyardCharacterDispatch 派遣侍从
func (s *Session) GraveyardCharacterDispatch(ctx context.Context, req *pb.C2SGraveyardCharacterDispatch) (*pb.S2CGraveyardCharacterDispatch, error) {
	// 检查模拟经营是否解锁
	err := s.User.CheckActionUnlock(static.ActionIdTypeTycoonunlock)
	if err != nil {
		return nil, err
	}

	//参数检查
	characterMap := map[int32]int32{}
	for _, character := range req.Characters {
		_, ok := characterMap[character.Position]
		if ok {
			return nil, errors.WrapTrace(common.ErrParamError)
		}
		characterMap[character.CharacterId] = character.Position
	}
	voList, err := s.User.GraveyardCharacterDispatch(req.BuildUid, common.NewCharacterPositions(characterMap))
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CGraveyardCharacterDispatch{
		Builds: voList,
	}, nil
}

// GraveyardPopulationSet 设置人口
func (s *Session) GraveyardPopulationSet(ctx context.Context, req *pb.C2SGraveyardPopulationSet) (*pb.S2CGraveyardPopulationSet, error) {
	// 检查模拟经营是否解锁
	err := s.User.CheckActionUnlock(static.ActionIdTypeTycoonunlock)
	if err != nil {
		return nil, err
	}

	//参数检查
	populationMap := map[int64]int32{}
	for _, population := range req.Populations {

		if population.GetCurrPopulationCount() < 0 {
			return nil, errors.WrapTrace(common.ErrParamError)
		}
		_, ok := populationMap[population.BuildUid]
		if ok {
			return nil, errors.WrapTrace(common.ErrParamError)
		}

		populationMap[population.BuildUid] = population.GetCurrPopulationCount()
	}
	if len(populationMap) == 0 {
		return nil, errors.WrapTrace(common.ErrParamError)
	}

	voList, err := s.User.GraveyardPopulationSet(populationMap)
	if err != nil {
		return nil, err
	}
	return &pb.S2CGraveyardPopulationSet{
		BuildInfos: voList,
	}, nil
}

// GraveyardAccelerate 加速
func (s *Session) GraveyardAccelerate(ctx context.Context, req *pb.C2SGraveyardAccelerate) (*pb.S2CGraveyardAccelerate, error) {
	// 检查模拟经营是否解锁
	err := s.User.CheckActionUnlock(static.ActionIdTypeTycoonunlock)
	if err != nil {
		return nil, err
	}
	// check param
	if len(req.Consumes) == 0 {
		return nil, common.ErrParamError
	}
	rewards, err := common.ParseFromVOConsume(req.Consumes)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// check consume
	err = s.User.CheckRewardsEnough(rewards)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	vo, err := s.User.GraveyardAccelerate(req.BuildUid, rewards)

	return &pb.S2CGraveyardAccelerate{
		VoBuild:        vo,
		ResourceResult: s.User.VOResourceResult(),
	}, nil

}

// GraveyardHelp 帮助他人减cd
func (s *Session) GraveyardHelp(ctx context.Context, req *pb.C2SGraveyardHelp) (*pb.S2CGraveyardHelp, error) {
	// 检查模拟经营是否解锁
	err := s.User.CheckActionUnlock(static.ActionIdTypeTycoonunlock)
	if err != nil {
		return nil, err
	}
	err = s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	if s.Guild.GuildID == 0 {
		return nil, errors.Swrapf(common.ErrGuildNotFound)
	}

	resp, err := s.RPCGuildHelpRequestsHandle(ctx)
	if err != nil {
		return nil, err
	}

	if resp.GetDailyAddGold() > s.User.Graveyard.DailyGuildGoldByHelp {
		rewards := common.NewRewards()
		rewards.AddReward(common.NewReward(static.CommonResourceTypeGuildMoney, resp.GetDailyAddGold()-s.User.Graveyard.DailyGuildGoldByHelp))
		_, err = s.User.AddRewards(rewards, logreason.NewReason(logreason.GraveyardHelp))
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
	}
	s.User.Graveyard.DailyGuildGoldByHelp = resp.GetDailyAddGold()
	s.User.Graveyard.DailyAddActivationByHelp = resp.GetDailyAddActivation()

	s.User.Guild.AddTaskItem(static.GuildTaskHelp, 1)
	return &pb.S2CGraveyardHelp{
		HelpCount:      resp.HelpCount,
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

// GraveyardSendHelpRequest 请求减cd
func (s *Session) GraveyardSendHelpRequest(ctx context.Context, req *pb.C2SGraveyardSendHelpRequest) (*pb.S2CGraveyardSendHelpRequest, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	if s.Guild.GuildID == 0 {
		return nil, errors.Swrapf(common.ErrGuildNotFound)
	}
	// 检查模拟经营是否解锁
	err = s.User.CheckActionUnlock(static.ActionIdTypeTycoonunlock)
	if err != nil {
		return nil, err
	}
	build, err := s.Graveyard.FindByUid(req.GetBuildUid())
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	requestAdd, err := s.User.GenVOAddGraveyardRequest(req.GetBuildUid(), build)
	if err != nil {
		return nil, err
	}
	resp, err := s.RPCGuildAddGraveyardRequest(ctx, requestAdd)
	if err != nil {
		return nil, err
	}
	build.SetRequestId(resp.RequestUid)
	s.Graveyard.SendRequestCount += 1
	return &pb.S2CGraveyardSendHelpRequest{
		Build:            s.VOGraveyardBuild(req.BuildUid, build),
		SendRequestCount: s.Graveyard.SendRequestCount,
	}, nil
}

// GraveyardGetHelpRequests 获得请求帮助列表
func (s *Session) GraveyardGetHelpRequests(ctx context.Context, req *pb.C2SGraveyardGetHelpRequests) (*pb.S2CGraveyardGetHelpRequests, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	if s.Guild.GuildID == 0 {
		return nil, errors.Swrapf(common.ErrGuildNotFound)
	}
	resp, err := s.RPCGuildGetGraveyardRequests(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.S2CGraveyardGetHelpRequests{
		Requests:               resp.Requests,
		TodayHelpAddActivation: s.Graveyard.DailyAddActivationByHelp,
	}, nil
}

// GraveyardReceivePlotReward 领取小人对话奖励
func (s *Session) GraveyardReceivePlotReward(ctx context.Context, req *pb.C2SGraveyardReceivePlotReward) (*pb.S2CGraveyardReceivePlotReward, error) {
	// 检查模拟经营是否解锁
	err := s.User.CheckActionUnlock(static.ActionIdTypeTycoonunlock)
	if err != nil {
		return nil, err
	}

	rewardHours, err := s.User.GraveyardReceivePlotReward()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CGraveyardReceivePlotReward{
		RewardHours:    rewardHours,
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

func (s *Session) GraveyardGetRewardHours(ctx context.Context, req *pb.C2SGraveyardGetRewardHours) (*pb.S2CGraveyardGetRewardHours, error) {
	//// 检查模拟经营是否解锁
	//err := s.User.CheckActionUnlock(static.ActionIdTypeTycoonunlock)
	//if err != nil {
	//	return nil, err
	//}

	return &pb.S2CGraveyardGetRewardHours{
		RewardHours: s.User.GraveyardGetRewardHours(),
	}, nil
}
func (s *Session) GraveyardUseBuff(ctx context.Context, req *pb.C2SGraveyardUseBuff) (*pb.S2CGraveyardUseBuff, error) {
	// 检查模拟经营是否解锁
	err := s.User.CheckActionUnlock(static.ActionIdTypeTycoonunlock)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	vos, err := s.User.GraveyardUseBuff(req.ItemId, req.Amount)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CGraveyardUseBuff{
		ResourceResult: s.User.VOResourceResult(),
		BuildInfoList:  vos,
	}, nil

}
