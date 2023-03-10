package entry

import (
	"shared/common"
	"shared/utility/errors"
	"shared/utility/servertime"
	"shared/utility/transfer"
	"sync"
)

const CfgNameEverydayEnergyReceive = "cfg_everyday_energy_receive"

type UserSettingEntry struct {
	sync.RWMutex
	dailyRewards []DailyReward
}

type DailyReward struct {
	Id       int32
	StartSec int32 `src:"StarReceiveTime" rule:"dailyTimeToSec"` // 时间格式转换成秒，如16:30
	EndSec   int32 `src:"EndReceiveTime" rule:"dailyTimeToSec"`  // 时间格式转换成秒，如16:30
	DropID   int32
}

func NewUserSettingEntry() *UserSettingEntry {
	return &UserSettingEntry{
		dailyRewards: nil,
	}
}

func (u *UserSettingEntry) Check(config *Config) error {
	return nil
}

func (u *UserSettingEntry) Reload(config *Config) error {
	u.Lock()
	defer u.Unlock()
	dailyRewards := make([]DailyReward, 0, 3)
	for i := 1; i <= 3; i++ {
		v, ok := config.CfgEverydayEnergyReceiveConfig.Find(int32(i))
		if !ok {
			return errors.WrapTrace(errors.Swrapf(common.ErrNotFoundInCSV, CfgNameEverydayEnergyReceive, i))
		}
		dailyReward := &DailyReward{}
		err := transfer.Transfer(v, dailyReward)
		if err != nil {
			return errors.WrapTrace(err)
		}
		dailyRewards = append(dailyRewards, *dailyReward)
	}

	u.dailyRewards = dailyRewards
	return nil
}

func (u *UserSettingEntry) GetDailyRewardByIndex(index int) (*DailyReward, error) {
	u.RLock()
	defer u.RUnlock()
	if index > 2 || index < 0 {
		return nil, errors.WrapTrace(common.ErrNotFoundInCSV)
	}
	dailyReward := u.dailyRewards[index]
	hour, min, sec := servertime.Now().Clock()
	totalSec := int32(hour*servertime.SecondPerHour + min*servertime.SecondPerMinute + sec)
	if totalSec < dailyReward.StartSec || totalSec > dailyReward.EndSec {
		return nil, errors.WrapTrace(common.ErrParamError)
	}
	return &dailyReward, nil
}
