package session

import (
	"context"
	"gamesvr/manager"
	"shared/common"
	"shared/protobuf/pb"
)

func (s *Session) ExploreInfo(ctx context.Context, req *pb.C2SExploreInfo) (*pb.S2CExploreInfo, error) {
	return &pb.S2CExploreInfo{
		ExploreInfo: s.User.ExploreInfo.VOExploreInfo(),
	}, nil
}

func (s *Session) EnterChapterMap(ctx context.Context, req *pb.C2SEnterChapterMap) (*pb.S2CEnterChapterMap, error) {
	err := s.User.EnterChapterMap(req.ChapterId)
	if err != nil {
		return nil, err
	}

	exploreMap, _ := s.User.ExploreInfo.GetMap(req.ChapterId)

	return &pb.S2CEnterChapterMap{
		ExploreMap: exploreMap.VOExploreMap(),
	}, nil
}

func (s *Session) ExploreUpdatePosition(ctx context.Context, req *pb.C2SExploreUpdatePosition) (*pb.S2CExploreUpdatePosition, error) {
	if req.Pos == nil {
		return nil, common.ErrParamError
	}
	pos := common.NewVec2(float64(req.Pos.PosX), float64(req.Pos.PosY))
	err := s.User.ExploreUpdatePosition(req.ChapterId, pos)
	if err != nil {
		return nil, err
	}

	return &pb.S2CExploreUpdatePosition{
		ChapterId: req.ChapterId,
	}, nil
}

func (s *Session) ExploreNPCInteract(ctx context.Context, req *pb.C2SExploreNPCInteract) (*pb.S2CExploreNPCInteract, error) {
	npcId := req.NpcId

	err := s.User.ExploreNPCInteraction(npcId, req.Option)
	if err != nil {
		return nil, err
	}

	npcCfg, _ := manager.CSV.ExploreEntry.GetExploreNPC(npcId)
	exploreEvent, _ := s.User.ExploreInfo.GetEventPoint(npcCfg.EventPointId)

	return &pb.S2CExploreNPCInteract{
		NpcId:          npcId,
		ExploreEvent:   exploreEvent.VOExploreEvent(),
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

func (s *Session) ExploreGather(ctx context.Context, req *pb.C2SExploreGather) (*pb.S2CExploreGather, error) {
	rewardId := req.GatherId

	err := s.User.ExploreRewardPointInteraction(rewardId)
	if err != nil {
		return nil, err
	}

	rewardCfg, _ := manager.CSV.ExploreEntry.GetExploreRewardPoint(rewardId)
	exploreEvent, _ := s.User.ExploreInfo.GetEventPoint(rewardCfg.EventPointId)
	return &pb.S2CExploreGather{
		GatherId:       rewardId,
		ExploreEvent:   exploreEvent.VOExploreEvent(),
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

func (s *Session) UnlockFog(ctx context.Context, req *pb.C2SUnlockFog) (*pb.S2CUnlockFog, error) {
	err := s.User.ExploreUnlockFog(req.FogId)
	if err != nil {
		return nil, err
	}

	return &pb.S2CUnlockFog{
		FogId: req.FogId,
	}, nil
}

func (s *Session) ExploreStartCollectResource(ctx context.Context, req *pb.C2SExploreStartCollectResource) (*pb.S2CExploreStartCollectResource, error) {
	err := s.User.ExploreStartCollectResource(req.ResourceId)
	if err != nil {
		return nil, err
	}

	resource, _ := s.User.ExploreInfo.GetResourcePoint(req.ResourceId)
	return &pb.S2CExploreStartCollectResource{
		Resource: resource.VOExploreResource(),
	}, nil
}

func (s *Session) ExploreFinishCollectResource(ctx context.Context, req *pb.C2SExploreFinishCollectResource) (*pb.S2CExploreFinishCollectResource, error) {
	err := s.User.ExploreFinishResourceCollect(req.ResourceId)
	if err != nil {
		return nil, err
	}

	resource, _ := s.User.ExploreInfo.GetResourcePoint(req.ResourceId)
	return &pb.S2CExploreFinishCollectResource{
		Resource:       resource.VOExploreResource(),
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

func (s *Session) ExploreEnterTPGate(ctx context.Context, req *pb.C2SExploreEnterTPGate) (*pb.S2CExploreEnterTPGate, error) {
	err := s.User.ExploreTransportGateEnter(req.Id)
	if err != nil {
		return nil, err
	}

	gate, _ := s.User.ExploreInfo.GetTransportGate(req.Id)

	return &pb.S2CExploreEnterTPGate{
		Gate: gate.VOExploreTransportGate(),
	}, nil
}

func (s *Session) ExploreDestroyTPGate(ctx context.Context, req *pb.C2SExploreDestroyTPGate) (*pb.S2CExploreDestroyTPGate, error) {
	err := s.User.ExploreTransportGateDestroy(req.Id)
	if err != nil {
		return nil, err
	}

	gate, _ := s.User.ExploreInfo.GetTransportGate(req.Id)

	return &pb.S2CExploreDestroyTPGate{
		Gate: gate.VOExploreTransportGate(),
	}, nil
}

//----------------------------------------
//Notify
//----------------------------------------
func (s *Session) TryPushExplore() {
	eventNotify := s.User.PopExploreEventNotify()
	if eventNotify == nil {
		return
	}

	s.push(manager.CSV.Protocol.Pushes.ExploreEventNotify, eventNotify)

	resourceNotify := s.User.PopExploreResourceNotify()
	if resourceNotify == nil {
		return
	}

	s.push(manager.CSV.Protocol.Pushes.ExploreResourceNotify, resourceNotify)
}
