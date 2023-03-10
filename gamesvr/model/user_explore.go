package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/statistic/logreason"
	"shared/utility/errors"
	"shared/utility/glog"
	"shared/utility/servertime"
)

func (u *User) EnterChapterMap(chapterId int32) error {
	err := u.CheckChapterUnlock(chapterId)
	if err != nil {
		return err
	}

	_, ok := u.ExploreInfo.GetMap(chapterId)
	if !ok {
		chapterCfg, _ := manager.CSV.ChapterEntry.GetChapter(chapterId)
		u.ExploreInfo.AddMap(chapterId, chapterCfg.InitPos)
	}

	u.ExploreInfo.SetCurChapter(chapterId)

	return nil
}

func (u *User) ExploreUpdatePosition(chapterId int32, pos *common.Vec2) error {
	if u.ExploreInfo.GetCurChapter() != chapterId {
		return errors.Swrapf(common.ErrNotInChapter, chapterId)
	}

	exploreMap, _ := u.ExploreInfo.GetMap(chapterId)
	exploreMap.SetCurPosition(pos)

	return nil
}

func (u *User) ExploreNPCInteraction(npcId, option int32) error {
	npcCfg, err := manager.CSV.ExploreEntry.GetExploreNPC(npcId)
	if err != nil {
		return err
	}

	if int(option) > len(npcCfg.OptionDrops) || option < 1 {
		return errors.Swrapf(common.ErrExploreNPCOptionInvalid, npcId, option)
	}

	eventId := npcCfg.EventPointId
	err = u.checkExploreEventPointAvailable(eventId)
	if err != nil {
		return err
	}

	cost := npcCfg.OptionCost[option-1]
	err = u.CheckRewardsEnough(cost)
	if err != nil {
		return err
	}

	eventPoint, ok := u.ExploreInfo.GetEventPoint(eventId)
	if !ok {
		eventPoint = u.ExploreInfo.AddEventPoint(eventId)
	}
	eventPoint.RecordInteration()

	reason := logreason.NewReason(logreason.ExploreNPC)
	u.AddRewardsByDropId(npcCfg.OptionDrops[option-1], reason)

	return nil
}

func (u *User) ExploreRewardPointInteraction(rewardPointId int32) error {
	rewardPointCfg, err := manager.CSV.ExploreEntry.GetExploreRewardPoint(rewardPointId)
	if err != nil {
		return err
	}

	eventId := rewardPointCfg.EventPointId
	err = u.checkExploreEventPointAvailable(eventId)
	if err != nil {
		return err
	}

	eventPoint, ok := u.ExploreInfo.GetEventPoint(eventId)
	if !ok {
		eventPoint = u.ExploreInfo.AddEventPoint(eventId)
	}
	eventPoint.RecordInteration()

	reason := logreason.NewReason(logreason.ExploreRewardPoint)
	u.AddRewardsByDropId(rewardPointCfg.GatherDrop, reason)

	return nil
}

func (u *User) ExploreUnlockFog(fogId int32) error {
	fogCfg, err := manager.CSV.ExploreEntry.GetFog(fogId)
	if err != nil {
		return err
	}

	err = u.CheckUserConditions(fogCfg.UnlockConditions)
	if err != nil {
		return err
	}

	u.ExploreInfo.UnlockFog(fogId)

	return nil
}

func (u *User) ExploreStartCollectResource(id int32) error {
	resourceCfg, err := manager.CSV.ExploreEntry.GetResource(id)
	if err != nil {
		return err
	}

	resourcePoint, ok := u.ExploreInfo.GetResourcePoint(id)
	if !ok {
		resourcePoint = u.ExploreInfo.AddResourcePoint(id)
		if resourceCfg.Monstermap <= 0 {
			resourcePoint.ClearMonster()
		}
	}

	if resourcePoint.IsMonsterExist() {
		return errors.Swrapf(common.ErrLevelNotPassed, resourceCfg.Monstermap)
	}

	if resourcePoint.IsCollecting() {
		return errors.Swrapf(common.ErrExploreResourceIsCollecting, id)
	}

	if resourcePoint.GetCollectTimes() >= resourceCfg.RefreshLimit {
		return errors.Swrapf(common.ErrExploreResourceCollectTimesLimit, id)
	}

	resourcePoint.StartCollect()

	return nil
}

func (u *User) ExploreFinishResourceCollect(id int32) error {
	resourcePoint, ok := u.ExploreInfo.GetResourcePoint(id)
	if !ok {
		return errors.Swrapf(common.ErrExploreResourceNotCollecting, id)
	}

	if !resourcePoint.IsCollecting() {
		return errors.Swrapf(common.ErrExploreResourceNotCollecting, id)
	}

	resourceCfg, err := manager.CSV.ExploreEntry.GetResource(id)
	if err != nil {
		return err
	}

	if servertime.Now().Unix() < resourcePoint.GetCollectStartTime()+resourceCfg.Time {
		return errors.Swrapf(common.ErrExploreResourceIsCollecting, id)
	}

	resourcePoint.FinishCollect()

	reason := logreason.NewReason(logreason.ExploreResource)
	u.AddRewardsByDropId(resourceCfg.Drop, reason)

	return nil
}

func (u *User) ExploreTransportGateEnter(id int32) error {
	tpCfg, err := manager.CSV.ExploreEntry.GetTransportGate(id)
	if err != nil {
		return err
	}

	transportGate, ok := u.ExploreInfo.GetTransportGate(id)
	if ok {
		if tpCfg.CompleteLimit != 0 && transportGate.GetUseTimes() >= tpCfg.CompleteLimit {
			return errors.Swrapf(common.ErrExploreTPGateUseTimesLimited, id)
		}
	}

	err = u.CheckUserConditions(tpCfg.UnlockConditions)
	if err != nil {
		return err
	}

	if !ok {
		transportGate = u.ExploreInfo.AddTransportGate(id)
	}

	transportGate.RecordUsed()

	return nil
}

func (u *User) ExploreTransportGateDestroy(id int32) error {
	_, err := manager.CSV.ExploreEntry.GetTransportGate(id)
	if err != nil {
		return err
	}

	transportGate, ok := u.ExploreInfo.GetTransportGate(id)
	if !ok {
		transportGate = u.ExploreInfo.AddTransportGate(id)
	}

	transportGate.SetDestroyShowed()

	return nil
}

func (u *User) CheckExploreMonsterLevel(levelId, monsterId int32) error {
	monsterCfg, err := manager.CSV.ExploreEntry.GetExploreMonster(monsterId)
	if err != nil {
		return err
	}

	if monsterCfg.Monstermap != levelId {
		return errors.Swrapf(common.ErrExploreMonsterHasNotLevel, monsterId, levelId)
	}

	err = u.checkExploreEventPointAvailable(monsterCfg.EventPointId)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) CheckExploreResourcePointLevel(levelId, resourceId int32) error {
	resourceCfg, err := manager.CSV.ExploreEntry.GetResource(resourceId)
	if err != nil {
		return err
	}
	if resourceCfg.Monstermap != levelId {
		return errors.Swrapf(common.ErrExploreResourceHasNotLevel, resourceId, levelId)
	}

	resourcePoint, ok := u.ExploreInfo.GetResourcePoint(resourceId)
	if ok && !resourcePoint.IsMonsterExist() {
		return errors.Swrapf(common.ErrExploreResourceLevelPassed, resourceId, levelId)
	}

	return nil
}

func (u *User) OnExploreMonsterLevelPass(levelId, monsterId int32) {
	monsterCfg, err := manager.CSV.ExploreEntry.GetExploreMonster(monsterId)
	if err != nil {
		glog.Error(err)
		return
	}

	eventId := monsterCfg.EventPointId
	eventPoint, ok := u.ExploreInfo.GetEventPoint(eventId)
	if !ok {
		eventPoint = u.ExploreInfo.AddEventPoint(eventId)
	}
	eventPoint.RecordInteration()

	u.AddExploreEventNotify(eventId)
}

func (u *User) OnExploreResourceLevelPass(levelId, resourceId int32) {
	resourcePoint, ok := u.ExploreInfo.GetResourcePoint(resourceId)
	if !ok {
		resourcePoint = u.ExploreInfo.AddResourcePoint(resourceId)
	}

	resourcePoint.ClearMonster()

	u.AddExploreResourceNotify(resourceId)
}

func (u *User) checkExploreEventPointAvailable(eventId int32) error {
	objCfg, err := manager.CSV.ExploreEntry.GetExploreEventPoint(eventId)
	if err != nil {
		return err
	}

	err = u.CheckUserConditions(objCfg.UnlockConditions)
	if err != nil {
		return err
	}

	eventPoint, ok := u.ExploreInfo.GetEventPoint(eventId)
	if !ok {
		return nil
	}

	if eventPoint.GetInteractTimes() >= objCfg.RefreshLimit {
		return errors.Swrapf(common.ErrExploreInteracted, eventId)
	}

	if eventPoint.GetInteractTime()+objCfg.RefreshTime > servertime.Now().Unix() {
		return errors.Swrapf(common.ErrExploreInteracted, eventId)
	}

	return nil
}
