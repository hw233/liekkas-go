package session

import (
	"context"
	"gamesvr/manager"
	"shared/protobuf/pb"
)

func (s *Session) TowerInfo(ctx context.Context, req *pb.C2STowerInfo) (*pb.S2CTowerInfo, error) {
	return &pb.S2CTowerInfo{
		Towers: s.User.TowerInfo.VOTowerInfo(),
	}, nil
}

//----------------------------------------
//Notify
//----------------------------------------
func (s *Session) TryPushTowerUpdateNotify() {
	notify := s.User.PopTowerUpdateNotify()
	if notify == nil {
		return
	}

	s.push(manager.CSV.Protocol.Pushes.TowerUpdateNotify, notify)
}
