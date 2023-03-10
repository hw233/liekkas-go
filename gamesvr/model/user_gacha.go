package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/statistic/bilog"
	"shared/statistic/logreason"
	"shared/utility/errors"
	"shared/utility/servertime"
)

// 每日刷新
func (u *User) GachaDailyRefresh(refreshTime int64) {
	// 每个池子每日抽卡次数清零
	for _, record := range u.GachaRecords.PoolRecords {
		record.TodayGachaCount = 0
	}
	// 清除过期的抽卡记录
	u.GachaRecords.removeRecordIfOutOfDate()

}

func (u *User) GetGachaList() *pb.S2CGetGachaList {
	// 清理过期的池子
	u.GachaRecords.removePoolsIfOutOfDate()
	var result []*pb.VOGachaInfo
	pools := manager.CSV.Gacha.GetCurrentPools(u.GachaRecords.NewPlayerDrew)
	for _, pool := range pools {
		if u.CheckUserConditions(pool.UnlockCondition) == nil {
			result = append(result, u.GachaRecords.VOGachaInfo(pool))
		}

	}
	return &pb.S2CGetGachaList{
		GachaList: result,
	}
}

func (u *User) UserGachaDrop(poolId int32, single bool) ([]*pb.VOGachaInfo, error) {
	err := u.CheckActionUnlock(static.ActionIdTypeGachadrop)
	if err != nil {
		return nil, err
	}

	pool, err := manager.CSV.Gacha.GetGachaPool(poolId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	if !pool.CheckPoolInTime(servertime.Now().Unix()) {
		return nil, errors.WrapTrace(common.ErrParamError)
	}
	err = u.CheckUserConditions(pool.UnlockCondition)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	if pool.Type == static.GachaPoolTypeNewPlayer {
		if single == true {
			// 新手池子只能10连
			return nil, errors.WrapTrace(common.ErrParamError)
		}
		if u.GachaRecords.NewPlayerDrew {
			// 新手池子已经抽过了
			return nil, errors.WrapTrace(common.ErrParamError)
		}
	}

	biCommon := u.CreateBICommon(bilog.EventNameGacha)
	var num int32 = 1
	var reasonEnum int32 = logreason.GachaOnce
	if !single {
		num = manager.CSV.Gacha.GetGachaMultipleValue()
		reasonEnum = logreason.GachaTen
	}
	reason := logreason.NewReason(reasonEnum, logreason.AddRelateLog(biCommon.GetLogId()))

	poolRecord, typeRecord := u.GachaRecords.GetByPoolIdCreate(pool)
	if poolRecord.TodayGachaCount >= pool.DailyLimit {
		// 超过每日限制
		return nil, errors.WrapTrace(common.ErrParamError)
	}
	// 抽卡消耗
	multiple := manager.CSV.Gacha.GachaConsume(pool.SingleConsume, num)
	err = u.CheckRewardsEnough(multiple)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	err = u.CostRewards(multiple, reason)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	gachaRewards, err := u.GachaRecords.Gacha(u.ID, pool, num, poolRecord, typeRecord)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	rewards, err := u.addRewardsForMerge(gachaRewardsToRewards(gachaRewards).Shuffle(), false, reason)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	var vos []*pb.VOGachaInfo
	vos = append(vos, u.GachaRecords.VOGachaInfo(*pool))
	pools := manager.CSV.Gacha.GetCurrentPools(u.GachaRecords.NewPlayerDrew)
	// 返回同类型的池子 因为同类型池子回共享概率增幅
	for _, tmpPool := range pools {
		if tmpPool.Type != pool.Type {
			continue
		}
		if tmpPool.Id == pool.Id {
			continue
		}
		if u.CheckUserConditions(pool.UnlockCondition) == nil {
			vos = append(vos, u.GachaRecords.VOGachaInfo(tmpPool))
		}

	}

	u.GachaRecords.RecordTotalTimes(num)

	u.TriggerQuestUpdate(static.TaskTypeGachaTimes, pool.Type, num)
	u.TriggerQuestUpdate(static.TaskTypeGachaRare, gachaRewards)
	u.TriggerQuestUpdate(static.TaskTypeGachaSame, gachaRewards)
	u.TriggerQuestUpdate(static.TaskTypeGachaTimesSsr, gachaRewards)
	u.TriggerQuestUpdate(static.TaskTypeGachaItemTimes, gachaRewards)

	u.BIGacha(poolId, pool.Type, !single, multiple, rewards, u.GachaRecords.GetTotalTimes(), biCommon)

	return vos, nil

}

func gachaRewardsToRewards(gachaRewards []*common.GachaReward) *common.Rewards {
	ret := common.NewRewards()
	for _, reward := range gachaRewards {
		ret.AddReward(reward.Reward)
	}
	return ret
}

func (u *User) UserGachaRecords(charaOrWorldItem bool, offset int64, num int32) []*pb.VOGachaResultRecord {

	var records *common.GachaResultRecords
	if charaOrWorldItem {
		records = u.GachaRecords.CharacterRecord
	} else {
		records = u.GachaRecords.WorldItemRecord

	}
	result := records.PagingSearch(offset, int(num))
	var vos []*pb.VOGachaResultRecord
	// 只显示n个月
	time := servertime.Now().AddDate(0, -manager.CSV.Gacha.GetGachaRecordShowMonth(), 0).Unix()
	for _, record := range result {
		if record.CreateAt > time {
			vos = append(vos, record.VOGachaResultRecord())
		}

	}
	return vos
}
