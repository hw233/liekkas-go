package servertime

import (
	"sync"
	"time"
)

const (
	SecondPerMinute = 60
	SecondPerHour   = SecondPerMinute * 60
	SecondPerDay    = SecondPerHour * 24
	SecondPerWeek   = SecondPerDay * 7
	TimeFormat      = "2006-01-02 15:04:05"

	TimeOffsetRedisName = "timeoffset"
)

var (
	refreshOffset time.Duration = 5 * time.Hour // 刷新时间（hour）
	timeOffset    time.Duration = 0             //时间偏移
)

var timeOffsetLock sync.RWMutex

func SetDailyRefreshHour(hour int32) {
	refreshOffset = time.Hour * time.Duration(hour)
}

func SetTimeOffset(offset int64) {
	timeOffsetLock.Lock()
	timeOffset = time.Second * time.Duration(offset)
	timeOffsetLock.Unlock()
}

func Now() time.Time {
	return time.Now().Add(timeOffset)
}
func OriginNow() time.Time {
	return time.Now()
}

func DayZeroTime(t time.Time) time.Time {
	year, month, day := t.Date()
	zeroTm := time.Date(year, month, day, 0, 0, 0, 0, time.Local)

	return zeroTm
}

func DayOffsetZeroTime(t time.Time, offset time.Duration) time.Time {
	zeroTime := DayZeroTime(t)

	offsetZero := zeroTime.Add(offset)
	if offsetZero.After(t) {
		offsetZero = offsetZero.AddDate(0, 0, -1)
	}

	return offsetZero
}

func WeekZeroTime(t time.Time) time.Time {
	weekday := t.Weekday()

	dayZeroTime := DayZeroTime(t)

	dayBack := (int(weekday) + 6) % 7
	weekZeroTime := dayZeroTime.AddDate(0, 0, -dayBack)

	return weekZeroTime
}

func WeekOffsetZeroTime(t time.Time, offset time.Duration) time.Time {
	weekZeroTime := WeekZeroTime(t)
	weekOffsetZeroTime := weekZeroTime.Add(offset)

	if weekOffsetZeroTime.After(t) {
		weekOffsetZeroTime = weekOffsetZeroTime.AddDate(0, 0, -7)
	}

	return weekOffsetZeroTime
}

func MonthZeroTime(t time.Time) time.Time {
	year, month, _ := t.Date()
	zeroTime := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)

	return zeroTime
}

func MonthOffsetZeroTime(t time.Time, offset time.Duration) time.Time {
	monthZeroTime := MonthZeroTime(t)
	monthOffsetZeroTime := monthZeroTime.Add(offset)

	if monthOffsetZeroTime.After(t) {
		monthOffsetZeroTime = monthOffsetZeroTime.AddDate(0, -1, 0)
	}

	return monthOffsetZeroTime
}

func GetHourWhen(timestamp int64) int32 {
	return int32(time.Unix(timestamp, 0).Hour())
}

func ParseTime(timeStr string) (int64, error) {
	timeObj, err := time.ParseInLocation(TimeFormat, timeStr, time.Local)

	if err != nil {
		return 0, err
	}

	return timeObj.Unix(), nil
}
