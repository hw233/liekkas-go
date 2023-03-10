package model

import "time"

type TimeRefreshChecker interface {
	CheckRefresh(checkTime int64) bool
	UpdateRefreshTime(checkTime int64)
}

type DailyRefreshChecker struct {
	LastDailyRefreshTime int64 `json:"last_daily_refresh_time"`
}

func NewDailyRefreshChecker() *DailyRefreshChecker {
	return &DailyRefreshChecker{
		LastDailyRefreshTime: 0,
	}
}

func (drc *DailyRefreshChecker) CheckRefresh(checkTime int64) bool {
	if drc.LastDailyRefreshTime >= checkTime {
		return false
	}

	return true
}

func (drc *DailyRefreshChecker) UpdateRefreshTime(checkTime int64) {
	drc.LastDailyRefreshTime = checkTime
}

func (drc *DailyRefreshChecker) GetLastDailyRefreshTime() int64 {
	return drc.LastDailyRefreshTime
}

type DailyRefresher struct {
	Checker     TimeRefreshChecker
	RefreshFunc func(int64)
}

func NewDailyRefresher(checker TimeRefreshChecker, refreshFunc func(int64)) *DailyRefresher {
	return &DailyRefresher{
		Checker:     checker,
		RefreshFunc: refreshFunc,
	}
}

func (dr *DailyRefresher) TryRefresh(now int64) {
	checkTime := DailyRefreshTime(time.Unix(now, 0)).Unix()
	if dr.Checker.CheckRefresh(checkTime) {
		dr.RefreshFunc(checkTime)
		dr.Checker.UpdateRefreshTime(checkTime)
	}
}
