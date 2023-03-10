package session

import (
	"context"
	"gamesvr/manager"
	"shared/protobuf/pb"
)

func (s *Session) ChapterInfo(ctx context.Context, req *pb.C2SChapterInfo) (*pb.S2CChapterInfo, error) {
	return &pb.S2CChapterInfo{
		ChapterInfo: s.User.ChapterInfo.VOChapterInfo(),
	}, nil
}

func (s *Session) ReceiveChapterReward(ctx context.Context, req *pb.C2SReceiveChapterReward) (*pb.S2CReceiveChapterReward, error) {
	err := s.User.ReceiveChapterReward(req.RewardId)
	if err != nil {
		return nil, err
	}

	rewardCfg, _ := manager.CSV.ChapterEntry.GetChapterReward(req.RewardId)
	chapter, _ := s.User.ChapterInfo.GetChapter(rewardCfg.ChapterId)

	return &pb.S2CReceiveChapterReward{
		Chapter:        chapter.VOChapter(),
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}
