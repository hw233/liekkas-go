package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/statistic/logreason"
	"shared/utility/errors"
	"shared/utility/glog"
	"shared/utility/rand"
	"shared/utility/servertime"
	"time"
)

func (u *User) StoreDailyRefresh(refreshTime int64) {
	for sid, store := range u.StoreInfo.Stores {
		storeData, err := manager.CSV.Store.GetStoreGeneral(sid)
		if err != nil {
			glog.Errorf("StoreDailyRefresh GetStoreGeneral err:%+v", err)
			return
		}
		if store.checkUpdate() {
			err = u.UpdateStore(sid, store, storeData.StoreType)
			if err != nil {
				glog.Errorf("StoreDailyRefresh UpdateStore err:%+v", err)
				return
			}
			store.UpdateTimes = 0
			u.StoreInfo.Stores[sid] = store
		}
	}
	u.AddStoreNotify()
}

// 消耗资源强制刷新商店
func (u *User) ForceToUpdate(storeID int32) error {
	if storeID == static.StoreNameGachaStore {
		err := u.CheckActionUnlock(static.ActionIdTypeGachadrop)
		if err != nil {
			return err
		}
	} else {
		err := u.CheckActionUnlock(static.ActionIdTypeStoreunlock)
		if err != nil {
			return err
		}
	}
	if !u.checkStoreUnlockCondition(storeID) {
		return errors.Swrapf(common.ErrStoreUnlockConditionNotSatisfied, storeID)
	}

	store, ok := u.StoreInfo.Stores[storeID]
	if !ok {
		return errors.Swrapf(common.ErrStoreWrongStoreIDForStoreInfo, storeID)
	}

	storeData, err := manager.CSV.Store.GetStoreGeneral(storeID)
	if err != nil {
		return err
	}

	updateRule, err := manager.CSV.Store.GetUpdateRule(storeData.UpdateRule)
	if err != nil {
		return err
	}

	// 如果Times为0，就代表不可强制刷新
	if store.UpdateTimes >= updateRule.Times {
		return errors.Swrapf(common.ErrStoreExceedMaxUpdateTimes, storeID)
	}
	costs := common.NewRewards()
	cost := common.NewReward(updateRule.Currency, updateRule.Cnt[store.UpdateTimes])
	costs.AddReward(cost)

	err = u.CheckRewardsEnough(costs)
	if err != nil {
		return err
	}

	err = u.UpdateStore(storeID, store, storeData.StoreType)
	if err != nil {
		return err
	}

	reason := logreason.NewReason(logreason.RefreshShop)
	err = u.CostRewards(costs, reason)
	if err != nil {
		return err
	}
	store.UpdateTimes += 1

	return nil
}

func (u *User) CheckGoodsInfo(storeID int32) error {
	if storeID == static.StoreNameGachaStore {
		err := u.CheckActionUnlock(static.ActionIdTypeGachadrop)
		if err != nil {
			return err
		}
	} else {
		err := u.CheckActionUnlock(static.ActionIdTypeStoreunlock)
		if err != nil {
			return err
		}
	}

	storeData, err := manager.CSV.Store.GetStoreGeneral(storeID)
	if err != nil {
		return err
	}

	// if !u.checkStoreUnlockCondition(storeID) {
	// 	// return errors.Swrapf(common.ErrStoreUnlockConditionNotSatisfied, storeID)
	// 	return nil
	// }

	// var store *Store
	// 1.查看种类商店是否存在，若不存在则建立种类商店
	store, ok := u.StoreInfo.Stores[storeID]
	if !ok {

		updateRule, err := manager.CSV.Store.GetUpdateRule(storeData.UpdateRule)
		if err != nil {
			return err
		}
		store = NewStore(updateRule.Period)
		err = u.UpdateStore(storeID, store, storeData.StoreType)
		if err != nil {
			return err
		}
		store.UpdateTimes = 0
	}

	u.StoreInfo.Stores[storeID] = store

	return nil
}

// 根据概率计算此栏位应该是哪个商品,返回其id
func calProbability(ids []int32) (int32, error) {

	probs := make([]int32, 0, len(ids))

	for _, id := range ids {
		goods, err := manager.CSV.Store.GetGoods(id)
		if err != nil {
			return 0, err
		}
		probs = append(probs, goods.Probability)
	}

	if len(probs) < 1 {
		return 0, errors.Swrapf(common.ErrStoreNoGoodsInCellToGenerate)
	}

	index := rand.SinglePerm(probs)

	return ids[index], nil
}

func (u *User) Purchase(storeID, subStoreID, index, goodsID, currency, num int32) (int32, error) {

	if storeID == static.StoreNameGachaStore {
		err := u.CheckActionUnlock(static.ActionIdTypeGachadrop)
		if err != nil {
			return 0, err
		}
	} else {
		err := u.CheckActionUnlock(static.ActionIdTypeStoreunlock)
		if err != nil {
			return 0, err
		}
	}

	if !u.checkSubStoreUnlockCondition(subStoreID) {
		return 0, errors.Swrapf(common.ErrStoreSubStoreUnlockConditionNotSatisfied)
	}

	if index < 0 || index >= 19 {
		return 0, errors.Swrapf(common.ErrStoreIndexOutOfRangeForInfo)
	}
	// 检查要访问的栏位是否已经解锁
	if !u.checkCellUnlockCondition(subStoreID, int32(index)) {
		return 0, errors.Swrapf(common.ErrStoreCellUnlockConditionNotSatisfied, index)
	}

	// 检查商品是否已经解锁
	if !u.checkGoodsUnlockCondition(goodsID) {
		return 0, errors.Swrapf(common.ErrStoreGoodsUnlockConditionNotSatisfied, goodsID)
	}

	if !u.checkSubStoreUnlockCondition(subStoreID) {
		return 0, errors.Swrapf(common.ErrStoreSubStoreUnlockConditionNotSatisfied, subStoreID)
	}

	store, ok := u.StoreInfo.Stores[storeID]
	if !ok {
		return 0, errors.Swrapf(common.ErrStoreWrongStoreIDForStoreInfo, storeID)
	}

	storeData, err := manager.CSV.Store.GetStoreGeneral(storeID)
	if err != nil {
		return 0, err
	}

	if !store.checkForSubStoreID(subStoreID) {
		return 0, errors.Swrapf(common.ErrStoreWrongSubStoreIDInStoreForInfo, subStoreID)
	}

	// 对于阶段式商店，需要做额外检查。目前检查规则为底阶段商店售罄才会开启高阶段商店
	if int(storeData.StoreType) == 2 {
		err = u.checkStageStoreForSoldOut(store, subStoreID)
		if err != nil {
			return 0, err
		}
	}

	subStore, ok := u.StoreInfo.SubStores[subStoreID]
	if !ok {
		return 0, errors.Swrapf(common.ErrStoreWrongSubStoreIDForInfo, subStoreID)
	}

	err = subStore.checkForNum(index, goodsID, num)
	if err != nil {
		return 0, err
	}

	// 计算货币的消耗以及资源的增加
	err = u.calRewards(goodsID, currency, num)
	if err != nil {
		return 0, err
	}

	// 购买之后，更新物品的数量信息
	cnt, err := subStore.updateAfterPurchase(index, goodsID, num)
	if err != nil {
		return 0, err
	}

	return cnt, nil
}

func (u *User) QuickPurchase(goodsID int32) error {

	goods, err := manager.CSV.Store.GetGoods(goodsID)
	if err != nil {
		return err
	}

	// 检查商品是否已经解锁，在此只是一个简单的检查
	if !u.checkGoodsUnlockCondition(goodsID) {
		return errors.Swrapf(common.ErrStoreGoodsUnlockConditionNotSatisfied, goodsID)
	}

	if int(goods.Times) != -1 || int(goods.Probability) != 10000 {
		return errors.Swrapf(common.ErrStoreWrongGoodsIDOfQuickPurchaseForData, goodsID)
	}

	err = u.calRewards(goodsID, static.CommonResourceTypeDiamondGift, 1)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) calRewards(goodsID int32, currency int32, num int32) error {
	goods, err := manager.CSV.Store.GetGoods(goodsID)
	if err != nil {
		return err
	}

	rewards := goods.Gain

	costs := common.NewRewards()

	currencyCorrect := false

	for i, supportCurrency := range goods.Currencies {
		if currency == supportCurrency {
			currencyCorrect = true
			cost := common.NewReward(currency, goods.Price[i])
			costs.AddReward(cost)
			costs = costs.Multiple(num)
			break
		}
	}

	if !currencyCorrect {
		return errors.Swrapf(common.ErrStoreCurrencyNotMatchDuringPurchase, goodsID)
	}

	// 加道具，减钱
	err = u.CheckRewardsEnough(costs)
	if err != nil {
		return err
	}

	reason := logreason.NewReason(logreason.Shop)
	err = u.CostRewards(costs, reason)
	if err != nil {
		return err
	}
	rewards = rewards.Multiple(num)

	_, err = u.addRewards(rewards, reason)
	if err != nil {
		return err
	}
	return nil
}

// 应该是对整个种类商店的刷新，所以对于阶段式商店，应该是全部刷新
func (u *User) UpdateStore(sid int32, store *Store, storeType int32) error {

	if storeType == 1 {
		subIDs, err := u.SelectSubStore(sid)
		if err != nil {
			return err
		}

		uselessSubStore := store.SubStores
		for _, id := range uselessSubStore { // 替换式商店只会显示一个子商店，随着高阶的子商店解锁，之前的低阶子商店可以删掉
			delete(u.StoreInfo.SubStores, id)
		}

		store.SubStores = subIDs
	} else {
		subIDs, err := u.SelectSubStores(sid)
		if err != nil {
			return err
		}
		store.SubStores = subIDs
	}

	for _, subID := range store.SubStores {

		err := u.UpdateSubStore(subID)
		if err != nil {
			return err
		}
	}

	store.LastTimeUpdate = servertime.Now().Unix()

	return nil
}

func (u *User) GetStoreInfo(storeID int32) (*pb.VOStoreInfo, error) {
	if !u.checkStoreUnlockCondition(storeID) {
		// return errors.Swrapf(common.ErrStoreUnlockConditionNotSatisfied, storeID)
		return nil, nil
	}
	if storeID == static.StoreNameGachaStore {
		err := u.CheckActionUnlock(static.ActionIdTypeGachadrop)
		if err != nil {
			return nil, err
		}
	} else {
		err := u.CheckActionUnlock(static.ActionIdTypeStoreunlock)
		if err != nil {
			return nil, err
		}
	}

	store, ok := u.StoreInfo.Stores[storeID]
	if !ok {
		return nil, errors.Swrapf(common.ErrStoreWrongStoreIDForStoreInfo, storeID)
	}

	storeData, err := manager.CSV.Store.GetStoreGeneral(storeID)
	if err != nil {
		return nil, err
	}

	updateRule, err := manager.CSV.Store.GetUpdateRule(storeData.UpdateRule)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	subIDs := store.SubStores
	if len(subIDs) == 0 {
		err = u.UpdateStore(storeID, store, storeData.StoreType)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		if len(subIDs) == 0 {
			return nil, nil
		}
	}
	subInfos := make([]*pb.VOSubStoreInfo, 0, len(subIDs))

	for _, subID := range subIDs {
		subStore, ok := u.StoreInfo.SubStores[subID]
		if !ok {
			return nil, errors.Swrapf(common.ErrStoreWrongSubStoreIDForInfo, subID)
		}

		subInfos = append(subInfos, subStore.VOSubStoreInfo())
	}

	return &pb.VOStoreInfo{
		StoreID:       storeID,
		StoreType:     storeData.StoreType,
		UpdateTimes:   store.UpdateTimes,
		RestTime:      calRestTime(updateRule.Period),
		SubStoreInfos: subInfos,
	}, nil
}

func calRestTime(period int32) int64 {

	var restTime int64
	now := servertime.Now()
	switch period {
	case 1:
		// todo 时间修改，减少时间对象和时间戳之间的转换
		tomorrowRefreshTime := DailyRefreshTime(now).AddDate(0, 0, 1)
		restTime = int64(tomorrowRefreshTime.Sub(now).Seconds()) // 明天的刷新时间距离现在的秒数
	case 2:
		// todo 时间修改，减少时间对象和时间戳之间的转换
		nextWeekRefreshTime := WeekRefreshTime(now.AddDate(0, 0, 7))
		restTime = int64(nextWeekRefreshTime.Sub(now).Seconds()) // 下一周的刷新时间距离现在的秒数
	case 3:
		nowDate := servertime.Now()
		// currentMonthFirstDay := nowDate.AddDate(0, 0, -nowDate.Day()+1) // 本月的第一天
		currentMonthFirstDay := time.Date(nowDate.Year(), nowDate.Month(), 1, 0, 0, 0, 0, nowDate.Location())
		nextMonthFirstDay := currentMonthFirstDay.AddDate(0, 1, 0)               // 下个月的第一天
		restTime = int64(MonthRefreshTime(nextMonthFirstDay).Sub(now).Seconds()) // 下一个月的刷新时间距离现在的秒数
	default:
		restTime = -1
	}

	return restTime
}

// entry 对应的地方加入商品的检查，子商店中的商品id需要存在于商品表中
func (u *User) UpdateSubStore(subID int32) error {

	subData, err := manager.CSV.Store.GetSubStore(subID)
	if err != nil {
		return err
	}

	allChosen := make([]int32, 0, 20)

	for _, cell := range subData.Cells {

		// 获取当前栏位目前已经解锁的所有商品ID
		goodsIDs := make([]int32, 0, len(cell.Goods))
		for _, goodsID := range cell.Goods {
			if u.checkGoodsUnlockCondition(goodsID) {
				goodsIDs = append(goodsIDs, goodsID)
			}
		}

		if len(goodsIDs) < 1 {
			continue
		}

		chosen, err := calProbability(goodsIDs)
		if err != nil {
			return err
		}
		allChosen = append(allChosen, chosen)
	}

	// 把此时的商品数据记录下来
	subStore := NewSubStore(subID)
	err = subStore.updateGoods(allChosen)
	if err != nil {
		return err
	}

	u.StoreInfo.SubStores[subID] = subStore

	return nil
}
