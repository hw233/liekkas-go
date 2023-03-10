package model

import (
	"encoding/json"
	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/glog"
	"shared/utility/number"
)

type ManualInfo map[int32]bool

func NewManualInfo() *ManualInfo {
	return (*ManualInfo)(&map[int32]bool{})
}

func (m *ManualInfo) Active(manualId int32) {
	// 已获得图鉴
	if _, ok := (*m)[manualId]; ok {
		return
	}
	// 获得图鉴并设置未领奖状态
	(*m)[manualId] = false

}

func (m *ManualInfo) IsActive(manualId int32) bool {
	_, ok := (*m)[manualId]
	return ok
}

func (m *ManualInfo) VOManualInfo() []*pb.VOManualInfo {
	vos := make([]*pb.VOManualInfo, 0, len(*m))
	for manualId, rewarded := range *m {
		vos = append(vos, &pb.VOManualInfo{
			ManualId: manualId,
			Rewarded: rewarded,
		})
	}
	return vos
}

func (m *ManualInfo) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *ManualInfo) Unmarshal(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	err := json.Unmarshal(data, &m)
	if err != nil {
		glog.Errorf("json.Unmarshal error: %v", err)
		return err
	}

	return nil
}

func (m *ManualInfo) IsRewarded(manualId int32) (bool, error) {
	v, ok := (*m)[manualId]
	if !ok {
		return false, errors.WrapTrace(common.ErrManualNotGet)
	}
	return v, nil
}

func (m *ManualInfo) SetRewarded(manualIds *number.NonRepeatableArr) {
	for _, manualId := range manualIds.Values() {
		(*m)[manualId] = true

	}
}
