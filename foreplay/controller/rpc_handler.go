package controller

import (
	"context"
	"foreplay/manager"
	"shared/protobuf/pb"
)

type ForeplayHandler struct {
	*pb.UnimplementedForeplayServer
}

func NewForeplayHandler() (*ForeplayHandler, error) {
	return &ForeplayHandler{
		UnimplementedForeplayServer: &pb.UnimplementedForeplayServer{},
	}, nil
}

func (fh *ForeplayHandler) ReloadAnnouncement(ctx context.Context, req *pb.ReloadAnnouncementReq) (*pb.ReloadAnnouncementResp, error) {
	err := manager.Announcements.LoadAnnouncements(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ReloadAnnouncementResp{}, nil
}

func (fh *ForeplayHandler) ReloadMaintain(ctx context.Context, req *pb.ReloadMaintainReq) (*pb.ReloadMaintainResp, error) {
	err := manager.Server.LoadServerInfo(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ReloadMaintainResp{}, nil
}
