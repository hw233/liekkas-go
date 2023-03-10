package entry

import (
	"shared/common"
	"shared/utility/errors"
	"shared/utility/transfer"
	"sync"
	"time"
)

type SignInData struct {
	Id     int32
	Type   int32
	Year   int32
	Month  int32
	Start  int32
	DayCnt int32
	DropID []int32
}

type SignInEntry struct {
	sync.RWMutex
	SignInData map[int32]SignInData
}

func NewSignEntry() *SignInEntry {
	return &SignInEntry{}
}

func (s *SignInEntry) Check(config *Config) error {
	return nil
}

func (s *SignInEntry) Reload(config *Config) error {
	s.Lock()
	defer s.Unlock()

	signinData := map[int32]SignInData{}
	for _, signinCsv := range config.CfgSigninDataConfig.GetAllData() {
		monthDataCfg := &SignInData{}
		err := transfer.Transfer(signinCsv, monthDataCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		if int32(len(monthDataCfg.DropID)) != (monthDataCfg.DayCnt) {
			return errors.Swrapf(common.ErrSignInDropIDMismatchDayCnt, signinCsv.Id)
		}
		signinData[signinCsv.Id] = *monthDataCfg
	}

	s.SignInData = signinData

	return nil

}

// 根据签到ID和签到记录来寻找对应的dropID
func (s *SignInEntry) GetDropID(id int32, index int32) (int32, error) {
	s.RLock()
	defer s.RUnlock()

	signinData, ok := s.SignInData[id]
	if !ok {
		return 0, errors.Swrapf(common.ErrSignInWrongIDForSignInData, id)
	}

	if index >= int32(len(signinData.DropID)) {
		return 0, errors.Swrapf(common.ErrSignInIndexOutOfDropID, id)
	}

	return signinData.DropID[index], nil
}

// 遍历所有签到ID，ID对应一个时间段(一个月，或者一个活动时间段), 返回所有时间段内包含所给时间点的ID
func (s *SignInEntry) GetID(timestamp int64, offset int64) map[int32]int32 {
	s.RLock()
	defer s.RUnlock()

	result := map[int32]int32{}

	for _, signinData := range s.SignInData {
		//fmt.Printf("signinentry--------->id: %v\n", signinData.Id)
		date := time.Date(int(signinData.Year), time.Month(signinData.Month), int(signinData.Start), 0, 0, 0, 0, time.Local)
		//fmt.Printf("signinentry--------------->Date: %v year, %v month, %v, day.\n", date.Year(), date.Month(), date.Day())
		startTime := date.Unix() + offset // offset代表系统刷新时间相对于0点0分0秒的偏移
		endTime := startTime + int64(signinData.DayCnt)*24*60*60
		if timestamp >= startTime && timestamp < endTime {
			result[signinData.Id] = signinData.Type
		}
	}

	return result

}

func (s *SignInEntry) GetAll() *map[int32]SignInData {
	s.RLock()
	defer s.RUnlock()

	return &s.SignInData
}
