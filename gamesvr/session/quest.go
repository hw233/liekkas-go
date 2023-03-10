package session

import (
	"context"
	"gamesvr/manager"
	"gamesvr/model"
	"shared/protobuf/pb"
)

func (s *Session) GetQuestInfo(ctx context.Context, req *pb.C2SGetQuestInfo) (*pb.S2CGetQuestInfo, error) {
	return &pb.S2CGetQuestInfo{QuestInfo: s.User.VOQuestInfo()}, nil
}

func (s *Session) RecieveQuest(ctx context.Context, req *pb.C2SRecieveQuest) (*pb.S2CRecieveQuest, error) {
	err := s.User.RecieveQuest(req.QuestId)
	if err != nil {
		return nil, err
	}

	return &pb.S2CRecieveQuest{
		QuestId:        req.QuestId,
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

func (s *Session) RecieveQuestActivityReward(ctx context.Context,
	req *pb.C2SRecieveQuestActivityReward) (*pb.S2CRecieveQuestActivityReward, error) {

	err := s.User.RecieveQuestActivityReward(req.ActivityId)
	if err != nil {
		return nil, err
	}

	questActivityCfg, _ := manager.CSV.Quest.GetQuestActivity(req.ActivityId)
	activityData := s.User.QuestPack.GetActivityData(model.QuestModule(questActivityCfg.Module))

	return &pb.S2CRecieveQuestActivityReward{
		QuestActivity:  activityData.VOQuestActivity(),
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

//----------------------------------------
//Notify
//----------------------------------------
func (s *Session) TryPushQuestUpdateNotify() {
	notify := s.User.PopQuestUpdateNotify()
	if notify == nil {
		return
	}

	s.push(manager.CSV.Protocol.Pushes.QuestUpdateNotify, notify)
}
