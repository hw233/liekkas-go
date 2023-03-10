package entry

import (
	"shared/common"
	"shared/utility/errors"
	"shared/utility/transfer"
	"sync"
)

type SurveyInfo struct {
	Id              int32
	SurveyType      int32
	StartDate       int32
	DayCnt          int32
	UnlockCondition *common.Conditions `rule:"conditions"`
	QuestionCnt     int32
	MaxLen          int32
	DropID          int32
}

type SurveyEntry struct {
	sync.RWMutex

	SurveyInfos map[int32]SurveyInfo
}

func NewSurveyEntry() *SurveyEntry {
	return &SurveyEntry{}
}

func (s *SurveyEntry) Check(config *Config) error {
	return nil
}

func (s *SurveyEntry) Reload(config *Config) error {
	s.Lock()
	defer s.Unlock()

	surveyInfos := map[int32]SurveyInfo{}
	for _, surveyCsv := range config.CfgSurveyInfoConfig.GetAllData() {
		surveyCfg := &SurveyInfo{}
		err := transfer.Transfer(surveyCsv, surveyCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		if surveyCfg.SurveyType > 3 || surveyCfg.SurveyType < 1 {
			return errors.Swrapf(common.ErrSurveyTypeBeyondLimit, surveyCsv.Id)
		}
		surveyInfos[surveyCsv.Id] = *surveyCfg
	}

	s.SurveyInfos = surveyInfos

	return nil
}

// 根据survey的id找到对应的info
func (s *SurveyEntry) GetInfo(id int32) (*SurveyInfo, error) {
	s.RLock()
	defer s.RUnlock()

	surveyInfo, ok := s.SurveyInfos[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrSurveyWrongIDForSurveyInfos, id)
	}

	return &surveyInfo, nil

}

func (s *SurveyEntry) GetAll() *map[int32]SurveyInfo {
	s.RLock()
	defer s.RUnlock()

	return &s.SurveyInfos
}
