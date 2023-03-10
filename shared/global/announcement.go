package global

import (
	"context"
	"shared/utility/global"

	"github.com/go-redis/redis/v8"
)

var (
	KeyAnnouncementId = "announcement:annc:id"
	KeyBannerId       = "announcement:banner:id"
	KeyCautionId      = "announcement:caution:id"
)

type Announcement struct {
	anncId   *global.IncrID
	bannerId *global.IncrID
	// foreplayanncId   *global.IncrID
	// foreplayBannerId *global.IncrID
	cautionId *global.IncrID
	getter    *global.String
}

func NewAnnouncement(client *redis.Client) *Announcement {
	return &Announcement{
		anncId:   global.NewIncrID(client),
		bannerId: global.NewIncrID(client),
		// foreplayanncId:   global.NewIncrID(client),
		// foreplayBannerId: global.NewIncrID(client),
		cautionId: global.NewIncrID(client),
		getter:    global.NewString(client),
	}
}

func (a *Announcement) GenAnnouncementId(ctx context.Context) (int64, error) {
	return a.anncId.GenID(ctx, KeyAnnouncementId)
}

func (a *Announcement) LastAnnouncementId(ctx context.Context) (int64, error) {
	return a.getter.GetInt64(ctx, KeyAnnouncementId)
}

func (a *Announcement) GenBannerId(ctx context.Context) (int64, error) {
	return a.bannerId.GenID(ctx, KeyBannerId)
}

func (a *Announcement) LastBannerId(ctx context.Context) (int64, error) {
	return a.getter.GetInt64(ctx, KeyBannerId)
}

// func (a *Announcement) GenForeplayAnnouncementId(ctx context.Context) (int64, error) {
// 	return a.foreplayanncId.GenID(ctx, KeyForeplayAnnouncementId)
// }

// func (a *Announcement) LastForeplayAnnouncementId(ctx context.Context) (int64, error) {
// 	return a.getter.GetInt64(ctx, KeyForeplayAnnouncementId)
// }

// func (a *Announcement) GenForeplayBannerId(ctx context.Context) (int64, error) {
// 	return a.foreplayBannerId.GenID(ctx, KeyForeplayBannerId)
// }

// func (a *Announcement) LastForeplayBannerId(ctx context.Context) (int64, error) {
// 	return a.getter.GetInt64(ctx, KeyForeplayBannerId)
// }

func (a *Announcement) GenCautionId(ctx context.Context) (int64, error) {
	return a.cautionId.GenID(ctx, KeyCautionId)
}

func (a *Announcement) LastCautionId(ctx context.Context) (int64, error) {
	return a.getter.GetInt64(ctx, KeyCautionId)
}
