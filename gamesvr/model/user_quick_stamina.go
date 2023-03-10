package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/statistic/logreason"
	"shared/utility/errors"
)

// 购买体力
func (u *User) QuickPurchaseStamina() error {

	// 1. 获取体力表数据
	record := u.Info.StaminaRecord.Record
	maxTimes := manager.CSV.QuickPurchaseStamina.GetMaxTimes()
	if record >= maxTimes {
		return errors.Swrapf(common.ErrQuickPurchaseStaminaExceedMaxTimes)
	}
	stamina, err := manager.CSV.QuickPurchaseStamina.GetStaminaData(record)
	if err != nil {
		return err
	}

	// 3. 生成新奖励
	costs := common.NewRewards()
	costs.AddReward(common.NewReward(static.CommonResourceTypeDiamondGift, stamina.Cost))

	gain := common.NewRewards()
	gain.AddReward(common.NewReward(static.CommonResourceTypeEnergy, stamina.Count))

	// 4. 检查钻石是否足够
	err = u.CheckRewardsEnough(costs)
	if err != nil {
		return err
	}

	// 5. 扣钻石，加体力
	reason := logreason.NewReason(logreason.QuickBuyStamina)
	_, err = u.addRewards(gain, reason)
	if err != nil {
		return err
	}
	err = u.CostRewards(costs, reason)
	if err != nil {
		return err
	}

	record += 1
	u.Info.StaminaRecord.Record = record
	return nil
}

func (u *User) CheckStaminaUpdate() {
	if u.Info.StaminaRecord == nil {
		u.Info.StaminaRecord = NewStaminaRecord()
	}
	u.Info.StaminaRecord.UpdateRecord()
}

func (u *User) GetStaminaInfo() (*pb.VOStaminaInfo, error) {

	record := u.Info.StaminaRecord.Record
	stamina, err := manager.CSV.QuickPurchaseStamina.GetStaminaData(record)
	if err != nil {
		return nil, err
	}

	return &pb.VOStaminaInfo{
		Cnt:   record,
		Price: stamina.Cost,
	}, nil
}
