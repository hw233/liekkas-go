package model

import (
	"gamesvr/manager"
	"shared/utility/servertime"
	"time"
)

func DailyRefreshTime(t time.Time) time.Time {
	return servertime.DayOffsetZeroTime(t, manager.CSV.GlobalEntry.DailyRefreshTimeOffset())
}

func TodayRefreshTime() time.Time {
	return DailyRefreshTime(servertime.Now())
}

func WeekRefreshTime(t time.Time) time.Time {
	return servertime.WeekOffsetZeroTime(t, manager.CSV.GlobalEntry.DailyRefreshTimeOffset())
}

func ThisWeekRefreshTime() time.Time {
	return WeekRefreshTime(servertime.Now())
}

func MonthRefreshTime(t time.Time) time.Time {
	return servertime.MonthOffsetZeroTime(t, manager.CSV.GlobalEntry.DailyRefreshTimeOffset())
}

func ThisMonthRefreshTime() time.Time {
	return MonthRefreshTime(servertime.Now())
}
