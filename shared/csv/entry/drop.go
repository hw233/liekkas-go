package entry

import (
	"shared/utility/glog"
	"sync"

	"shared/common"
	"shared/utility/errors"
	"shared/utility/rand"
)

const (
	dropModeMutex    = 0 // 互斥掉落（在一个掉落组中，不可重复掉落同一条掉落数据）
	dropModeRandom   = 1 // 随机掉落（在一个掉落组中，可重复掉落同一条掉落数据）
	dropTypeAbsolute = 1 // 绝对概率（道具是否掉落只受自己本身的概率影响，填写掉落千分比）
	dropTypeRelative = 2 // 相对概率（道具相对概率数值/同组道具相对概率数值之和）
)

const (
	CfgDropDataConfig  = "cfg_drop_data"
	CfgDropGroupConfig = "cfg_drop_group"
)

type Drop struct {
	sync.RWMutex

	drops map[int32]DropInfo
}

type DropInfo struct {
	Mode    int32
	Type    int32
	Count   int
	Details []DropDetail
}

type DropDetail struct {
	Prob               int32
	*common.RandReward `rule:"randReward" src:"DropItem,DropNumber[0],DropNumber[1]"`
}

func NewDrop() *Drop {
	return &Drop{
		drops: map[int32]DropInfo{},
	}
}

func (d *Drop) Check(config *Config) error {
	for _, vGroup := range config.CfgDropGroupConfig.GetAllData() {
		if len(vGroup.DropNumber) != 2 {
			return common.ErrCSVFormatInvalid
		}
	}

	return nil
}

func (d *Drop) Reload(config *Config) error {
	d.Lock()
	defer d.Unlock()

	drops := map[int32]DropInfo{}

	details := map[int32][]DropDetail{}

	for _, vGroup := range config.CfgDropGroupConfig.GetAllData() {
		if len(vGroup.DropNumber) != 2 {
			return errors.WrapTrace(errors.Swrapf(common.ErrCSVFormatInvalid, CfgDropGroupConfig, vGroup.Id))
		}

		details[vGroup.DropGroup] = append(details[vGroup.DropGroup], DropDetail{
			Prob:       vGroup.DropPr,
			RandReward: common.NewRandReward(vGroup.DropItem, vGroup.DropNumber[0], vGroup.DropNumber[1]),
		})
	}

	for _, vDrop := range config.CfgDropDataConfig.GetAllData() {
		detail, ok := details[vDrop.Id]
		if !ok {
			return errors.WrapTrace(errors.Swrapf(common.ErrCSVFormatInvalid, CfgDropDataConfig, vDrop.Id))
		}

		drops[vDrop.Id] = DropInfo{
			Mode:    vDrop.DropMode,
			Type:    vDrop.DropType,
			Count:   int(vDrop.DropCount),
			Details: detail,
		}
	}

	d.drops = drops

	return nil
}

func (d *Drop) DropRewards(dropID int32) (*common.Rewards, error) {
	d.RLock()
	defer d.RUnlock()

	drop, ok := d.drops[dropID]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgDropDataConfig, dropID)
	}

	rewards := common.NewRewards()
	probs := make([]int32, 0, len(drop.Details))
	for _, detail := range drop.Details {
		probs = append(probs, detail.Prob)
	}

	glog.Infof("drop dropID: %v, type: %v mode: %v, count: %v, prob: %v, rewards: %v", dropID, drop.Type, drop.Mode, drop.Count, probs, drop.Details)

	switch drop.Type {
	case dropTypeAbsolute:
		for _, i := range rand.TryPerm(probs) {
			reward := drop.Details[i].NewReward()
			if reward.Num == 0 {
				continue
			}
			rewards.AddReward(reward)
		}
	case dropTypeRelative:
		switch drop.Mode {
		case dropModeMutex:
			for _, i := range rand.UniquePerm(drop.Count, probs) {
				reward := drop.Details[i].NewReward()
				if reward.Num == 0 {
					continue
				}
				rewards.AddReward(reward)
			}
		case dropModeRandom:
			for _, i := range rand.RepeatPerm(drop.Count, probs) {
				reward := drop.Details[i].NewReward()
				if reward.Num == 0 {
					continue
				}
				rewards.AddReward(reward)
			}
		default:
			return nil, errors.Swrapf(common.ErrCSVFormatInvalid, CfgDropDataConfig, dropID)
		}
	default:
		return nil, errors.Swrapf(common.ErrCSVFormatInvalid, CfgDropDataConfig, dropID)
	}

	return rewards, nil
}
