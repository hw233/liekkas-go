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
}

func NewAnnouncementInfo() *AnnouncementInfo {
	return &AnnouncementInfo{
		Announcements: []*common.Announcement{},
		Banners:       []*common.Banner{},
	}
}

func (ai *AnnouncementInfo) LoadAnnouncements(ctx context.Context) error {
	nowTs := servertime.Now().Unix()
	announcements, err := common.DBLoadAnnouncement(ctx, common.AnnouncementModuleGame, nowTs, DB)
	if err != nil {
		return err
	}

	banners, err := common.DBLoadBanner(ctx, common.AnnouncementModuleGame, nowTs, DB)
	if err != nil {
		return err
	}

	ai.Announcements = announcements
	ai.Banners = banners

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

func AnnouncementOnTick() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	Announcements.LoadAnnouncements(ctx)
}

func AnnouncementInit() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err := Announcements.LoadAnnouncements(ctx)
	if err != nil {
		return err
	}

	Timer.ScheduleFunc(time.Minute, AnnouncementOnTick)

	return nil
}
