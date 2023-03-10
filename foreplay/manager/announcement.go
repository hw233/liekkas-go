package manager

import (
	"context"
	"shared/common"
	"shared/utility/servertime"
	"time"
)

type AnnouncementInfo struct {
	Announcements []*common.Announcement
	Banners       []*common.Banner
	Cautions      []*common.Caution
}

func NewAnnouncementInfo() *AnnouncementInfo {
	return &AnnouncementInfo{
		Announcements: []*common.Announcement{},
		Banners:       []*common.Banner{},
		Cautions:      []*common.Caution{},
	}
}

func (ai *AnnouncementInfo) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err := ai.LoadAnnouncements(ctx)
	if err != nil {
		return err
	}

	Timer.ScheduleFunc(time.Minute, ai.OnTick)

	return nil
}

func (ai *AnnouncementInfo) LoadAnnouncements(ctx context.Context) error {
	nowTs := servertime.Now().Unix()
	announcements, err := common.DBLoadAnnouncement(ctx, common.AnnouncementModuleForeplay, nowTs, DB)
	if err != nil {
		return err
	}

	banners, err := common.DBLoadBanner(ctx, common.AnnouncementModuleForeplay, nowTs, DB)
	if err != nil {
		return err
	}

	cautions, err := common.DBLoadCaution(ctx, nowTs, DB)
	if err != nil {
		return err
	}

	ai.Announcements = announcements
	ai.Banners = banners
	ai.Cautions = cautions

	return nil
}

func (ai *AnnouncementInfo) GetAnnouncements() ([]*common.Announcement, []*common.Banner) {
	announcenments := make([]*common.Announcement, 0, len(ai.Announcements))
	banners := make([]*common.Banner, 0, len(ai.Banners))

	nowTs := servertime.Now().Unix()

	for _, annc := range ai.Announcements {
		if nowTs > annc.StartTime && annc.EndTime >= nowTs {
			announcenments = append(announcenments, annc)
		}
	}

	for _, banner := range ai.Banners {
		if nowTs > banner.StartTime && banner.EndTime >= nowTs {
			banners = append(banners, banner)
		}
	}

	return announcenments, banners
}

func (ai *AnnouncementInfo) GetCautions() []*common.Caution {
	cautions := make([]*common.Caution, 0, len(ai.Cautions))

	nowTs := servertime.Now().Unix()

	for _, caution := range ai.Cautions {
		if nowTs > caution.StartTime && caution.EndTime >= nowTs {
			cautions = append(cautions, caution)
		}
	}

	return cautions
}

func (ai *AnnouncementInfo) OnTick() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	ai.LoadAnnouncements(ctx)
}
