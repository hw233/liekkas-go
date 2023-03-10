package method

import (
	"context"
	"gm/manager"
	"gm/param"
	"shared/common"
	"shared/protobuf/pb"
)

func (p *HttpPostHandler) AddAnnouncement(ctx context.Context, param *param.AddAnnouncement) error {
	id, err := manager.Global.GenAnnouncementId(ctx)
	if err != nil {
		return err
	}

	announcement := common.NewAnnouncement(id, param.Type, param.Module, param.Title, param.Content, param.Image,
		param.Tag, param.StartTime, param.EndTime, param.ShowStartTime, param.ShowEndTime, param.Priority)

	err = common.DBAddAnnouncement(ctx, announcement, manager.DB)
	if err != nil {
		return err
	}

	pushByModule(ctx, param.Module)

	return nil
}

func (p *HttpPostHandler) UpdateAnnouncement(ctx context.Context, param *param.UpdateAnnouncement) error {
	announcement := common.NewAnnouncement(param.Id, param.Type, param.Module, param.Title, param.Content, param.Image,
		param.Tag, param.StartTime, param.EndTime, param.ShowStartTime, param.ShowEndTime, param.Priority)

	err := common.DBUpdateAnnouncement(ctx, announcement, manager.DB)
	if err != nil {
		return err
	}

	pushByModule(ctx, param.Module)

	return nil
}

func (p *HttpPostHandler) FetchAnnouncement(ctx context.Context, param *param.FetchAnnouncement) ([]*common.Announcement, error) {
	announcements, err := common.DBLoadAnnouncement(ctx, param.Module, 0, manager.DB)
	if err != nil {
		return nil, err
	}

	return announcements, nil
}

func (p *HttpPostHandler) DeleteAnnouncement(ctx context.Context, param *param.DeleteAnnouncement) error {
	err := common.DBDeleteAnnouncement(ctx, param.Id, manager.DB)

	pushByModule(ctx, common.AnnouncementModuleGame)
	pushByModule(ctx, common.AnnouncementModuleForeplay)

	return err
}

func (p *HttpPostHandler) AddBanner(ctx context.Context, param *param.AddBanner) error {
	id, err := manager.Global.GenBannerId(ctx)
	if err != nil {
		return err
	}

	banner := common.NewBanner(id, param.Type, param.Module, param.Image, param.Jump, param.StartTime,
		param.EndTime, param.ShowStartTime, param.ShowEndTime, param.Priority)

	err = common.DBAddBanner(ctx, banner, manager.DB)
	if err != nil {
		return err
	}

	pushByModule(ctx, param.Module)

	return nil
}

func (p *HttpPostHandler) UpdateBanner(ctx context.Context, param *param.UpdateBanner) error {
	Banner := common.NewBanner(param.Id, param.Type, param.Module, param.Image, param.Jump,
		param.StartTime, param.EndTime, param.ShowStartTime, param.ShowEndTime, param.Priority)

	err := common.DBUpdateBanner(ctx, Banner, manager.DB)
	if err != nil {
		return err
	}

	pushByModule(ctx, param.Module)

	return nil
}

func (p *HttpPostHandler) FetchBanner(ctx context.Context, param *param.FetchBanner) ([]*common.Banner, error) {
	banners, err := common.DBLoadBanner(ctx, param.Module, 0, manager.DB)
	if err != nil {
		return nil, err
	}

	return banners, nil
}

func (p *HttpPostHandler) DeleteBanner(ctx context.Context, param *param.DeleteBanner) error {
	err := common.DBDeleteBanner(ctx, param.Id, manager.DB)

	pushByModule(ctx, common.AnnouncementModuleGame)
	pushByModule(ctx, common.AnnouncementModuleForeplay)

	return err
}

func (p *HttpPostHandler) AddCaution(ctx context.Context, param *param.AddCaution) error {
	id, err := manager.Global.GenCautionId(ctx)
	if err != nil {
		return err
	}

	cauton := common.NewCaution(id, param.Content, param.StartTime, param.EndTime)

	err = common.DBAddCaution(ctx, cauton, manager.DB)
	if err != nil {
		return err
	}

	pushByModule(ctx, common.AnnouncementModuleForeplay)

	return nil
}

func (p *HttpPostHandler) UpdateCaution(ctx context.Context, param *param.UpdateCaution) error {
	Caution := common.NewCaution(param.Id, param.Content, param.StartTime, param.EndTime)

	err := common.DBUpdateCaution(ctx, Caution, manager.DB)
	if err != nil {
		return err
	}

	pushByModule(ctx, common.AnnouncementModuleForeplay)

	return nil
}

func (p *HttpGetHandler) FetchCaution(ctx context.Context) ([]*common.Caution, error) {
	cautions, err := common.DBLoadCaution(ctx, 0, manager.DB)
	if err != nil {
		return nil, err
	}

	return cautions, nil
}

func (p *HttpPostHandler) DeleteCaution(ctx context.Context, param *param.DeleteCaution) error {
	err := common.DBDeleteCaution(ctx, param.Id, manager.DB)

	pushByModule(ctx, common.AnnouncementModuleForeplay)

	return err
}

func pushByModule(ctx context.Context, module int32) error {
	var err error
	switch module {
	case common.AnnouncementModuleGame:
		_, err = manager.RPCGameClient.ReloadAnnouncement(
			ctx,
			&pb.ReloadAnnouncementReq{},
		)
	case common.AnnouncementModuleForeplay:
		_, err = manager.RPCForeplayClient.ReloadAnnouncement(
			ctx,
			&pb.ReloadAnnouncementReq{},
		)
	}

	return err
}
