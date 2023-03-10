package model

import (
	"shared/utility/glog"
	"time"
)

func (u *User) RegisterDailyRefreshers() {
	u.DailyRefreshers = []*DailyRefresher{
		NewDailyRefresher(u.QuestPack, u.QuestDailyRefresh),
		NewDailyRefresher(u.LevelsInfo, u.LevelDailyRefresh),
		NewDailyRefresher(u.TowerInfo, u.TowerDailyRefresh),
		NewDailyRefresher(u.Yggdrasil, u.YggdrasilDailyRefresh),
		NewDailyRefresher(u.Graveyard, u.GraveyardDailyRefresh),
		NewDailyRefresher(u.GachaRecords, u.GachaDailyRefresh),
		NewDailyRefresher(u.Info, u.DailyRefresh),
		NewDailyRefresher(u.StoreInfo, u.StoreDailyRefresh),
		NewDailyRefresher(u.Mercenary, u.MercenaryDailyRefresh),
	}
}

func (u *User) TryDailyRefresh(now int64) {
	defer func() {
		if err := recover(); err != nil {
			glog.Errorf("error: %s\n", err)
		}
	}()

	now_tm := time.Unix(now, 0)
	curRefreshTime := DailyRefreshTime(now_tm).Unix()
	if curRefreshTime <= u.Info.LastDailyRefreshTime {
		return
	}

	for _, dailyRefresher := range u.DailyRefreshers {
		dailyRefresher.TryRefresh(now)
	}

	u.AddDailyRefreshNotify()
	u.Info.LastDailyRefreshTime = curRefreshTime
}
