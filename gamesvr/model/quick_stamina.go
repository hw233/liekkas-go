package model

import (
	"shared/utility/servertime"
	"time"
)

type StaminaRecord struct {
	Record         int32 `json:"stamina_record"`
	LastTimeUpdate int64 `json:"last_time_update"`
}

func NewStaminaRecord() *StaminaRecord {
	return &StaminaRecord{
		Record:         0,
		LastTimeUpdate: 0,
	}
}

func (s *StaminaRecord) UpdateRecord() {
	now := servertime.Now()
	if s.checkForUpdate(now) {
		s.Record = 0
		s.LastTimeUpdate = now.Unix()
	}
}

func (s *StaminaRecord) checkForUpdate(t time.Time) bool {
	return s.LastTimeUpdate < DailyRefreshTime(t).Unix()
}
