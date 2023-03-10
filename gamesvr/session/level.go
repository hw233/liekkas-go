package session

import (
	"context"
	"gamesvr/manager"
	"gamesvr/model"
	"shared/protobuf/pb"
)

func (s *Session) LevelInfo(ctx context.Context, req *pb.C2SLevelInfo) (*pb.S2CLevelInfo, error) {
	return &pb.S2CLevelInfo{
		LevelInfo: s.User.LevelsInfo.VOLevelInfo(),
	}, nil
}

func (s *Session) LevelBattleStart(ctx context.Context, req *pb.C2SLevelBattleStart) (*pb.S2CLevelBattleStart, error) {
	seed, detail, err := s.User.StartLevel(req.LevelId, req.BattleFormation, req.LevelSystemType, req.LevelSystemParams)
	if err != nil {
		return nil, err
	}

	return &pb.S2CLevelBattleStart{
		LevelId:  req.LevelId,
		Seed:     seed,
		Attacker: detail,
		Defender: nil,
	}, nil
}

func (s *Session) LevelBattleEnd(ctx context.Context, req *pb.C2SLevelBattleEnd) (*pb.S2CLevelBattleEnd, error) {
	//todo get level id by battleId
	levelId := s.User.LevelsInfo.GetCurLevelId()
	battleData := &model.BattleVerifyData{
		IsWin:                req.IsWin,
		BattleResult:         req.BattleResult,
		BattleInput:          req.BattleInput,
		BattleOutput:         req.BattleOutput,
		CompleteTargets:      req.Targets,
		CompleteAchievements: req.Achievements,
		Statistic:            req.BattleStatis,
		BattleCharacters:     req.BattleCharacters,
	}

	levelResult, err := s.User.LevelEnd(ctx, levelId, battleData)
	if err != nil {
		return nil, err
	}

	level := s.User.LevelsInfo.GetOrCreateLevel(levelId)

	resp := &pb.S2CLevelBattleEnd{
		IsWin:        req.IsWin,
		BattleResult: req.BattleResult,
		VerifyResult: 1,
		Level:        level.VOLevel(),
		ResouceResults: &pb.VOLevelResouceResults{
			FirstPassResult:   levelResult.RewardRewsult.FirstPassResult,
			PassResult:        levelResult.RewardRewsult.PassResult,
			TargetResult:      levelResult.RewardRewsult.TargetResult,
			AchievementResult: levelResult.RewardRewsult.AchievementResult,
			EntityChange:      s.VOYggdrasilEntityChange(ctx),
		},
		Charas: make([]*pb.VOUserCharacter, 0, len(levelResult.Charas)),
	}

	for _, charaId := range levelResult.Charas {
		chara, _ := s.User.CharacterPack.Get(charaId)
		resp.Charas = append(resp.Charas, chara.VOUserCharacter())
	}

	if levelResult.ChapterId > 0 {
		chapter, _ := s.User.ChapterInfo.GetChapter(levelResult.ChapterId)
		resp.Chapter = chapter.VOChapter()
	}

	return resp, nil
}

func (s *Session) LevelSweep(ctx context.Context, req *pb.C2SLevelSweep) (*pb.S2CLevelSweep, error) {
	err := s.User.SweepLevel(req.LevelId, req.Times)
	if err != nil {
		return nil, err
	}

	level, _ := s.User.LevelsInfo.GetLevel(req.LevelId)
	return &pb.S2CLevelSweep{
		Level:          level.VOLevel(),
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

//----------------------------------------
//Notify
//----------------------------------------
func (s *Session) TryPushLevelNotify() {
	notify := s.User.PopLevelNofity()
	if notify == nil {
		return
	}

	s.push(manager.CSV.Protocol.Pushes.LevelNotify, notify)
}
