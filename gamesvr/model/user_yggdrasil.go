package model

import (
	"context"
	"gamesvr/manager"
	"math"
	"shared/common"
	"shared/csv/entry"
	"shared/csv/static"
	"shared/global"
	"shared/protobuf/pb"
	"shared/statistic/logreason"
	"shared/utility/coordinate"
	"shared/utility/errors"
	"shared/utility/servertime"
)

var (
	objectHandles = map[int32]objectHandle{}
)

type objectHandle func(ctx context.Context, u *User, object *YggObject) error

func init() {
	objectHandles[static.YggObjectTypeChest] = HandleChestObject
}

// YggdrasilDailyRefresh 每日刷新
func (u *User) YggdrasilDailyRefresh(refreshTime int64) {
	// 每日探险次数清零
	u.Yggdrasil.TravelTime = 0
	// 刷日常怪物
	u.Yggdrasil.DailyMonsters.DailyRefresh(context.Background(), u.Yggdrasil)

	// 日常派遣刷新
	u.Yggdrasil.Dispatch.YggDispatchDailyRefresh()
}

func (u *User) YggdrasilGetBlockInfo(ctx context.Context, positions []coordinate.Position) (*pb.S2CYggdrasilGetBlockInfo, error) {
	h, v := manager.CSV.Yggdrasil.GetYggBlockLengthAndWidth()

	for _, position := range positions {
		if position.X%h != 0 || position.Y%v != 0 {
			return nil, errors.WrapTrace(common.ErrParamError)
		}
	}
	return &pb.S2CYggdrasilGetBlockInfo{
		BlockAndArea: u.Yggdrasil.VOYggdrasilBlockAndArea(ctx, positions),
	}, nil
}

func (u *User) YggdrasilExploreStart(characters []int32) error {
	if u.Yggdrasil.Status == YggdrasilStateTravel {

		return errors.WrapTrace(common.ErrYggdrasilInTravel)
	}
	now := servertime.Now().Unix()
	// check 角色能否出行
	for _, characterId := range characters {
		chara, err := u.CharacterPack.Get(characterId)
		if err != nil {
			return errors.WrapTrace(err)
		}
		if chara.CanYggdrasilTime > now {
			return errors.WrapTrace(common.ErrYggdrasilCharacterCannotCarry)
		}
	}

	if len(characters) == 0 {
		return errors.WrapTrace(common.ErrParamError)
	}

	if int32(len(characters)) > manager.CSV.Yggdrasil.GetYggEditTeamCount(u.Info.Level.Value()) {
		return errors.WrapTrace(common.ErrParamError)
	}
	// check ap
	levelConfig, ok := manager.CSV.TeamLevelCache.GetByLv(u.Info.Level.Value())
	if !ok {
		return errors.Swrapf(common.ErrNotFoundInCSV, entry.CfgTeamLevelConfig, u.Info.Level.Value())

	}

	// check 今日可探索次数
	if u.Yggdrasil.TravelTime >= manager.CSV.Yggdrasil.GetYggDailyTravelTime() {
		return errors.WrapTrace(common.ErrYggdrasilNoTravelTime)
	}
	// 消耗ap
	consume := common.NewRewards()
	consume.AddReward(common.NewReward(static.CommonResourceTypeAp, levelConfig.ExploreAp))
	err := u.CheckRewardsEnough(consume)
	if err != nil {
		return errors.WrapTrace(err)
	}

	reason := logreason.NewReason(logreason.YggStartExplore)
	err = u.CostRewards(consume, reason)
	if err != nil {
		return errors.WrapTrace(err)
	}
	u.Guild.AddTaskItem(static.GuildTaskYgg, 1)
	// 初始化本次探索信息
	return u.Yggdrasil.initTravelInfo(characters, levelConfig.ExploreAp)
}

func (u *User) YggdrasilExploreMove(ctx context.Context, position coordinate.Position) ([]*pb.VOPosition, []*pb.VOExploredPosCountUpdate, error) {
	if u.Yggdrasil.Status == YggdrasilStateQuit {

		return nil, nil, errors.WrapTrace(common.ErrYggdrasilNotInTravel)
	}

	return u.Yggdrasil.moveTo(ctx, u, position)
}

func (u *User) YggdrasilExploreQuit(ctx context.Context) ([]*pb.VOUserCharacter, error) {

	// 旅行体力为0才可以
	if u.Yggdrasil.TravelInfo.TravelAp != 0 {
		return nil, errors.WrapTrace(common.ErrParamError)
	}

	if u.Yggdrasil.Status == YggdrasilStateQuit {

		return nil, errors.WrapTrace(common.ErrYggdrasilNotInTravel)
	}
	return u.doYggdrasilExploreQuit(ctx)
}

func (u *User) doYggdrasilExploreQuit(ctx context.Context) ([]*pb.VOUserCharacter, error) {
	now := servertime.Now().Unix()
	// 本次携带角色进入休息状态
	var characters []*pb.VOUserCharacter
	for characterId := range u.Yggdrasil.TravelInfo.CharactersHp {
		chara, err := u.CharacterPack.Get(characterId)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		chara.CanYggdrasilTime = now + manager.CSV.Yggdrasil.GetYggCharacterRestSec()
		characters = append(characters, chara.VOUserCharacter())
	}
	// 主动怪物回到起始点
	for _, object := range u.Yggdrasil.Entities.Objects {
		objectType, err := manager.CSV.Yggdrasil.GetYggdrasilObjectType(object.ObjectId, object.State)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		if objectType == static.YggObjectTypeInitiativemonster {
			err := u.Yggdrasil.ObjectMove(ctx, u, object, *object.OrgPos)
			if err != nil {
				return nil, errors.WrapTrace(err)
			}
		}
	}

	u.Yggdrasil.Status = YggdrasilStateQuit
	return characters, nil
}
func (u *User) YggdrasilExploreReturnCity(ctx context.Context) error {
	if u.Yggdrasil.Status == YggdrasilStateQuit {
		return errors.WrapTrace(common.ErrYggdrasilNotInTravel)
	}
	if u.Yggdrasil.CityId != 0 {
		return errors.WrapTrace(common.ErrYggdrasilInCityNow)
	}
	cityId := manager.CSV.Yggdrasil.IsCityEntrance(u.Yggdrasil.TravelPos)
	if cityId == 0 {
		return errors.WrapTrace(common.ErrYggdrasilCannotReturnCityThisPos)
	}
	err := u.Yggdrasil.TransferToCity(ctx, u, cityId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}

func (u *User) YggdrasilExploreLeaveCity(ctx context.Context) error {
	if u.Yggdrasil.CityId == 0 {
		return errors.WrapTrace(common.ErrYggdrasilNotInCityNow)
	}

	areaCity, err := manager.CSV.Yggdrasil.GetYggCityById(u.Yggdrasil.CityId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	_, _, err = u.Yggdrasil.ChangePos(ctx, u, *areaCity.CityExitPos)
	if err != nil {
		return errors.WrapTrace(err)
	}
	u.Yggdrasil.CityId = 0
	return nil
}

func (u *User) YggdrasilGoodsDiscard(ctx context.Context, replace bool, goodsId int64) error {
	if u.Yggdrasil.Status == YggdrasilStateQuit {
		return errors.WrapTrace(common.ErrYggdrasilNotInTravel)
	}
	goods, ok := u.Yggdrasil.Pack.Get(goodsId)
	if !ok {
		return errors.WrapTrace(common.ErrYggdrasilPackGoodsNotFound)
	}
	if !replace {
		if !u.Yggdrasil.Entities.IsBlankPos(*u.Yggdrasil.TravelPos) {
			return errors.WrapTrace(common.ErrYggdrasilOverlapping)
		}
	}

	u.Yggdrasil.DeletePackGoods(goodsId)
	yggDiscardGoods, err := NewYggDiscardGoods(*u.Yggdrasil.TravelPos, goods)
	if err != nil {
		return errors.WrapTrace(err)
	}
	return u.Yggdrasil.AppendDiscardGoods(ctx, u.ID, yggDiscardGoods)

}

func (u *User) YggdrasilGoodsPickUp(ctx context.Context, replaceGoodsId int64) error {

	if u.Yggdrasil.Status == YggdrasilStateQuit {
		return errors.WrapTrace(common.ErrYggdrasilNotInTravel)
	}

	discard, ok := u.Yggdrasil.Entities.FindDiscardGoodsByPos(*u.Yggdrasil.TravelPos)
	if !ok {
		return errors.WrapTrace(common.ErrYggdrasilDiscardGoodsNotFound)
	}

	//走替换流程
	if replaceGoodsId > 0 {
		err := u.YggdrasilGoodsDiscard(ctx, true, replaceGoodsId)
		if err != nil {
			return errors.WrapTrace(err)
		}
	} else {
		if u.Yggdrasil.Pack.IsFull(u.Info.Level.Value()) {
			return errors.WrapTrace(common.ErrYggdrasilBagIsFull)
		}
	}

	err := u.Yggdrasil.PickUp(discard)
	if err != nil {
		return errors.WrapTrace(err)
	}
	return u.Yggdrasil.RemoveDiscardGoods(ctx, u.ID, discard)

}

func (u *User) YggdrasilObjectHandle(ctx context.Context) error {
	if u.Yggdrasil.Status == YggdrasilStateQuit {
		return errors.WrapTrace(common.ErrYggdrasilNotInTravel)
	}

	object, err := u.Yggdrasil.Entities.FindObjectByPosAndType(*u.Yggdrasil.TravelPos, static.YggObjectTypeMagictable, static.YggObjectTypeEffect, static.YggObjectTypeChest)
	if err != nil {
		return errors.WrapTrace(err)
	}

	state, err := manager.CSV.Yggdrasil.GetObjectState(object.ObjectId, object.State)
	if err != nil {
		return errors.WrapTrace(err)
	}
	if state.SubTaskID > 0 {
		// 判断是否在subtask
		if _, _, ok := u.Yggdrasil.Task.YggSubTaskProgressInfoInProcess(state.SubTaskID); !ok {
			return errors.WrapTrace(common.ErrYggdrasilTaskNotInProgress)
		}
	}

	handle, ok := objectHandles[state.ObjectType]
	if ok {
		err := handle(ctx, u, object)
		if err != nil {
			return errors.WrapTrace(err)
		}
	}
	u.Yggdrasil.ObjectChangeToNextState(ctx, u, object)

	return nil

}

func HandleChestObject(ctx context.Context, u *User, object *YggObject) error {
	state, err := manager.CSV.Yggdrasil.GetObjectState(object.ObjectId, object.State)
	if err != nil {
		return errors.WrapTrace(err)
	}

	reason := logreason.NewReason(logreason.YggHandleChest)
	err = u.Yggdrasil.AddRewardsByDropId(ctx, u, state.ObjectParam, 0, reason)
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}

func (u *User) YggdrasilObjectMove(ctx context.Context, posMap map[int64]*coordinate.Position) ([]*pb.VOYggdrasilObject, error) {
	m := map[*YggObject]coordinate.Position{}
	for uid, position := range posMap {
		obj, ok := u.Yggdrasil.Entities.FindObjectByUid(uid)
		if !ok {
			return nil, errors.WrapTrace(common.ErrYggdrasilObjectNotFound)
		}
		_, err := checkObjectState(obj, static.YggObjectTypeInitiativemonster, static.YggObjectTypeNpc,
			static.YggObjectTypeFollownpc, static.YggObjectTypeFollowbattlenpc)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		m[obj] = *position
	}
	var voList []*pb.VOYggdrasilObject
	for obj, to := range m {
		err := u.Yggdrasil.ObjectMove(ctx, u, obj, to)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}

		voList = append(voList, obj.VOYggdrasilObject())
	}

	return voList, nil
}

func (u *User) YggdrasilQueryPosition(ctx context.Context, objectId int32) ([]*pb.VOYggdrasilObject, error) {

	// 初始化obj所在的area
	for _, areaId := range manager.CSV.Yggdrasil.GetObjectInitArea(objectId) {
		u.Yggdrasil.Areas.getByCreate(ctx, u.Yggdrasil, areaId)
	}
	objects, err := u.Yggdrasil.FindObjectById(objectId)

	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	var vos []*pb.VOYggdrasilObject
	for _, object := range objects {
		vos = append(vos, object.VOYggdrasilObject())
	}
	return vos, nil
}

func (u *User) checkYggdrasilBattle(objectUid int64, levelId int32, characters []*pb.VOBattleCharacter, npcs []*pb.VOBattleNPC) error {
	if u.Yggdrasil.Status == YggdrasilStateQuit {
		return errors.WrapTrace(common.ErrYggdrasilNotInTravel)
	}
	object, ok := u.Yggdrasil.Entities.FindObjectByUid(objectUid)
	if !ok {
		return errors.WrapTrace(common.ErrYggdrasilObjectNotFound)
	}
	// todo:主动怪物位置有时候会和玩家不一致
	//objPos := *object.Position
	//travelPos := *u.Yggdrasil.TravelPos
	//if objPos != travelPos {
	//	return errors.WrapTrace(common.ErrYggdrasilObjectNotFound)
	//}
	state, err := checkObjectState(object, static.YggObjectTypePassivermonster, static.YggObjectTypeInitiativemonster)
	if err != nil {
		return errors.WrapTrace(err)
	}
	if state.ObjectParam != levelId {
		return errors.WrapTrace(common.ErrParamError)
	}
	// 血量为0的不可出战
	for _, character := range characters {
		hp, err := u.Yggdrasil.TravelInfo.GetCharacterHp(character.CharacterId)
		if err != nil {
			return errors.WrapTrace(err)
		}
		if hp <= 0 {
			return errors.WrapTrace(errors.Swrapf(common.ErrYggdrasilCharacterHpErr, character.CharacterId, hp))
		}
		// 设置战斗血量
		character.HpPercent = hp
	}

	return nil
}

func (u *User) onYggdrasilLevelFail(ctx context.Context, objectUid int64, charactersAfterBattle []*pb.VOBattleCharacter, result int32) error {
	object, ok := u.Yggdrasil.Entities.FindObjectByUid(objectUid)
	if !ok {
		return errors.WrapTrace(common.ErrYggdrasilObjectNotFound)
	}
	state, err := checkObjectState(object, static.YggObjectTypePassivermonster, static.YggObjectTypeInitiativemonster)
	if err != nil {
		return errors.WrapTrace(err)
	}

	if state.ObjectType == static.YggObjectTypePassivermonster {
		// 被动怪物

		if result == static.BattleEndTypeGiveUp {
			// 战斗中放弃
			return u.YggBattleGiveUp(ctx, objectUid)
		} else {
			// 战斗失败
			// 设置血量
			u.SetYggCharacterHp(ctx, charactersAfterBattle)
			return u.YggBattleGiveUp(ctx, objectUid)

		}

	} else if state.ObjectType == static.YggObjectTypeInitiativemonster {
		// 主动怪物
		return u.SetAllDeadAndReturnSafePos(ctx)

	}

	return nil

}

func (u *User) onYggdrasilLevelPass(ctx context.Context, objectUid int64, levelCfg *entry.Level, charactersAfterBattle []*pb.VOBattleCharacter) error {
	object, ok := u.Yggdrasil.Entities.FindObjectByUid(objectUid)
	if !ok {
		return errors.WrapTrace(common.ErrYggdrasilObjectNotFound)
	}
	_, err := checkObjectState(object, static.YggObjectTypePassivermonster, static.YggObjectTypeInitiativemonster)
	if err != nil {
		return errors.WrapTrace(err)
	}
	// 设置血量
	u.SetYggCharacterHp(ctx, charactersAfterBattle)

	reason := logreason.NewReason(logreason.YggLevelPass)
	err = u.Yggdrasil.AddRewardsByDropIds(ctx, u, levelCfg.YggdrasilDrop, 0, reason)
	if err != nil {
		return errors.WrapTrace(err)
	}
	u.Yggdrasil.ObjectChangeToNextState(ctx, u, object)

	u.Yggdrasil.Task.ProcessNum(ctx, u, static.YggdrasilSubTaskTypeTypeChapter, 1, levelCfg.Id)
	u.Yggdrasil.Task.ProcessNum(ctx, u, static.YggdrasilSubTaskTypeTypeMonster, 1, object.ObjectId)

	return nil
}

func (u *User) checkChallengeAltarBattle(objectUid int64, levelId int32, characters []*pb.VOBattleCharacter, npcs []*pb.VOBattleNPC) error {
	err := u.CheckActionUnlock(static.ActionIdTypeChallengealtarunlock)
	if err != nil {
		return errors.WrapTrace(err)
	}
	if u.Yggdrasil.Status == YggdrasilStateQuit {
		return errors.WrapTrace(common.ErrYggdrasilNotInTravel)
	}
	object, ok := u.Yggdrasil.Entities.FindObjectByUid(objectUid)
	if !ok {
		return errors.WrapTrace(common.ErrYggdrasilObjectNotFound)
	}
	objPos := *object.Position
	travelPos := *u.Yggdrasil.TravelPos
	if objPos != travelPos {
		return errors.WrapTrace(common.ErrYggdrasilObjectNotFound)
	}
	state, err := checkObjectState(object, static.YggObjectTypeChallengealtar)
	if err != nil {
		return errors.WrapTrace(err)
	}
	if state.ObjectParam != levelId {
		return errors.WrapTrace(common.ErrParamError)
	}
	// 血量为0的不可出战
	for _, character := range characters {
		hp, err := u.Yggdrasil.TravelInfo.GetCharacterHp(character.CharacterId)
		if err != nil {
			return errors.WrapTrace(err)
		}
		if hp <= 0 {
			return errors.WrapTrace(errors.Swrapf(common.ErrYggdrasilCharacterHpErr, character.CharacterId, hp))
		}
		// 设置战斗血量
		character.HpPercent = hp
	}

	return nil
}

func (u *User) onChallengeAltarFail(ctx context.Context, objectUid int64, charactersAfterBattle []*pb.VOBattleCharacter, result int32) error {
	object, ok := u.Yggdrasil.Entities.FindObjectByUid(objectUid)
	if !ok {
		return errors.WrapTrace(common.ErrYggdrasilObjectNotFound)
	}
	_, err := checkObjectState(object, static.YggObjectTypeChallengealtar)
	if err != nil {
		return errors.WrapTrace(err)
	}

	if result == static.BattleEndTypeGiveUp {
		// 战斗中放弃
		return u.YggBattleGiveUp(ctx, objectUid)
	} else {
		// 战斗失败
		// 设置血量
		u.SetYggCharacterHp(ctx, charactersAfterBattle)

	}
	return u.YggBattleGiveUp(ctx, objectUid)

}

func (u *User) onChallengeAltarPass(ctx context.Context, objectUid int64, levelCfg *entry.Level, charactersAfterBattle []*pb.VOBattleCharacter) error {
	object, ok := u.Yggdrasil.Entities.FindObjectByUid(objectUid)
	if !ok {
		return errors.WrapTrace(common.ErrYggdrasilObjectNotFound)
	}
	_, err := checkObjectState(object, static.YggObjectTypeChallengealtar)
	if err != nil {
		return errors.WrapTrace(err)
	}
	// 设置血量
	u.SetYggCharacterHp(ctx, charactersAfterBattle)

	reason := logreason.NewReason(logreason.YggLevelPass)
	err = u.Yggdrasil.AddRewardsByDropIds(ctx, u, levelCfg.YggdrasilDrop, 0, reason)
	if err != nil {
		return errors.WrapTrace(err)
	}
	u.Yggdrasil.ObjectChangeToNextState(ctx, u, object)

	u.Yggdrasil.Task.ProcessNum(ctx, u, static.YggdrasilSubTaskTypeTypeChapter, 1, levelCfg.Id)
	u.Yggdrasil.Task.ProcessNum(ctx, u, static.YggdrasilSubTaskTypeTypeMonster, 1, object.ObjectId)

	return nil
}

func (u *User) YggdrasilAcceptTask(ctx context.Context, taskId int32) (*pb.VOYggdrasilTaskInfo, error) {
	config, err := manager.CSV.Yggdrasil.GetTaskConfig(taskId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	err = u.CheckUserConditions(config.UnlockCondition)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	info, err := u.Yggdrasil.Task.AcceptTask(ctx, u, config)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return info.VOYggdrasilTaskInfo(), nil

}

func (u *User) YggdrasilSetTrackTask(taskId int32) error {
	return u.Yggdrasil.Task.SetTrackTask(taskId)
}

func (u *User) YggdrasilCompleteTask(ctx context.Context, taskId int32) error {
	config, err := manager.CSV.Yggdrasil.GetTaskConfig(taskId)
	if err != nil {
		return errors.WrapTrace(err)
	}

	// 完成任务
	err = u.Yggdrasil.Task.CompleteTask(ctx, u, taskId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	// 发奖
	if config.DropId > 0 {
		reason := logreason.NewReason(logreason.YggTaskComplete)
		err := u.Yggdrasil.AddRewardsByDropId(ctx, u, config.DropId, 0, reason)
		if err != nil {
			return errors.WrapTrace(err)
		}
	}

	return nil
}

func (u *User) YggdrasilChooseNext(ctx context.Context, subTaskId int32) (*pb.VOYggdrasilTaskInfo, error) {
	info, err := u.Yggdrasil.Task.ChooseNext(ctx, u, subTaskId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return info.VOYggdrasilTaskInfo(), nil
}

//todo: complete_event 删掉    env可以跨多个subtask 添加删除
//todo: env中有的特效类永久保存

// 任务背包每个格子只有一个吗, 累积
// 休息时候会带背包里东西回去吗 ，不会，回城市或者用建筑
// 带回去的东西在“异界邮箱”里面吗 ，别人带回去的进“异界邮箱”，自己带的直接发奖，如果邮箱满了就不匹配（我的掉落物不进入其他人的匹配）

func (u *User) YggdrasilAbandonTask(ctx context.Context, taskId int32) error {
	return u.Yggdrasil.Task.TaskAbandon(ctx, u, taskId)
}

func (u *User) YggdrasilDeliverTaskGoods(ctx context.Context, subTaskId int32, resources []*pb.VOResource) error {
	return u.Yggdrasil.Task.DeliverTaskGoods(ctx, u, subTaskId, resources)
}

func (u *User) YggdrasilAreaProgressReward(ctx context.Context, areaId int32) (*pb.VOYggdrasilArea, error) {

	return u.Yggdrasil.ProgressReward(ctx, u, areaId)

}

func (u *User) YggdrasilMailGetByPage(offset int64, num int32) []*pb.VOYggdrasilMail {
	search := u.Yggdrasil.MailBox.PagingSearch(offset, int(num))

	vos := make([]*pb.VOYggdrasilMail, 0, len(search))

	for _, mail := range search {
		vos = append(vos, mail.VOYggdrasilMail())
	}
	return vos
}

func (u *User) YggdrasilMailReceiveOne(ctx context.Context, uid int64) error {
	yggdrasilMail, ok := u.Yggdrasil.MailBox.Get(uid)
	if !ok {
		return errors.WrapTrace(common.ErrYggdrasilMailNotExist)
	}

	reason := logreason.NewReason(logreason.YggMail)
	_, err := u.addRewards(yggdrasilMail.Attachment, reason)
	if err != nil {
		return errors.WrapTrace(err)
	}
	err = u.Yggdrasil.MailBox.Delete(ctx, u.ID, uid)
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}

func (u *User) YggdrasilMailReceiveAll(ctx context.Context) ([]int64, error) {
	all := u.Yggdrasil.MailBox.PagingSearch(math.MaxInt64, manager.CSV.Yggdrasil.GetMailReceiveMaxCount())
	if len(all) == 0 {
		return nil, errors.WrapTrace(common.ErrYggdrasilMailNotExist)
	}
	var ret []int64
	for _, mail := range all {
		err := u.YggdrasilMailReceiveOne(ctx, mail.Uid)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		ret = append(ret, mail.Uid)
	}
	return ret, nil
}

func (u *User) YggdrasilMarkCreate(ctx context.Context, markId int32, pos coordinate.Position) (*pb.VOYggdrasilMark, error) {
	_, err := manager.CSV.Yggdrasil.GetMark(markId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// mark上限
	if u.Yggdrasil.GetMarkTotalCount() >= manager.CSV.Yggdrasil.GetYggMarkTotalCount() {
		return nil, errors.WrapTrace(common.ErrYggdrasilMarkCountLimit)
	}
	mark, err := NewYggMark(ctx, markId, pos)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	u.Yggdrasil.AppendMark(ctx, mark)
	return mark.VOYggdrasilMark(), nil
}

func (u *User) YggdrasilMarkDestroy(markUId int64) error {
	mark, ok := u.Yggdrasil.Entities.FindMarkByUid(markUId)
	if !ok {
		return errors.WrapTrace(common.ErrYggdrasilMarkNotExist)
	}
	u.Yggdrasil.RemoveMark(mark)
	return nil
}

func (u *User) YggdrasilTrackMark(pos *coordinate.Position) {
	u.Yggdrasil.TrackMark = pos
}

func (u *User) YggdrasilMatch(ctx context.Context, otherMembers []int64) error {

	matchPool, err := GetYggdrasilInRedis(ctx, otherMembers)
	if err != nil {
		return errors.WrapTrace(err)
	}

	u.Yggdrasil.Entities.Match(ctx, u.Yggdrasil, u.Guild.GuildID, u.GetUserId(), matchPool)
	return nil
}

func GetYggdrasilInRedis(ctx context.Context, userIds []int64) (*MatchPool, error) {
	matchPool := NewMatchPool()
	for _, userId := range userIds {
		entities, err := LoadUserMatchEntities(ctx, userId)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		(*matchPool)[userId] = entities
	}

	return matchPool, nil
}

func (u *User) GetTotalIntimacy(ctx context.Context) (int32, error) {
	guildID := u.Guild.GuildID

	return getTotalIntimacy(ctx, u.GetUserId(), guildID)
}

func getTotalIntimacy(ctx context.Context, userId, guildID int64) (int32, error) {
	if guildID == 0 {
		return 0, nil
	}
	intimacyMap, err := manager.Global.GetGuildIntimacyMap(ctx, guildID, userId)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}
	var ret int32
	for _, i := range intimacyMap {
		ret += i
	}
	return ret, nil
}

// IntimacyChange 亲密度变化
func (u *User) IntimacyChange(ctx context.Context, other int64, changeVal int32) error {
	caches, err := manager.Global.GetUserCachesExtension(ctx, []int64{other}, global.UserCacheWithOnline|global.UserCacheWithGuild)
	if err != nil {
		return errors.WrapTrace(err)
	}
	otherCache, ok := caches[other]
	if !ok {
		return nil
	}
	guildID := u.Guild.GuildID
	// 不在公会
	if guildID == 0 {
		return nil
	}
	// 不在同一公会
	if guildID != otherCache.GuildID {
		return nil
	}
	intimacy, err := manager.Global.ChangeIntimacy(ctx, guildID, u.ID, other, changeVal)
	if err != nil {
		return errors.WrapTrace(err)
	}

	totalIntimacy, err := getTotalIntimacy(ctx, guildID, u.ID)
	if err != nil {
		return errors.WrapTrace(err)
	}
	u.AddYggPush(&pb.S2CYggdrasilIntimacyChange{
		UserId:        other,
		IntimacyValue: intimacy,
		TotalIntimacy: totalIntimacy,
	})

	// 其他玩家在线时才推
	if otherCache.OnlineStatus == 1 {
		othersTotalIntimacy, err := getTotalIntimacy(ctx, guildID, other)
		if err != nil {
			return errors.WrapTrace(err)
		}
		manager.EventQueue.Push(ctx, other, common.NewYggdrasilIntimacyChangeEvent(u.ID, intimacy, othersTotalIntimacy))
	}

	return nil
}

func (u *User) QuerySimpleInfo(ctx context.Context, userIds []int64) ([]*pb.VOUserInfoSimple, error) {
	if len(userIds) == 0 {
		return nil, errors.WrapTrace(common.ErrParamError)
	}
	for _, userId := range userIds {
		_, ok := u.Yggdrasil.MatchUserIds[userId]
		if !ok {
			return nil, errors.WrapTrace(common.ErrNoPermissionError)
		}
	}

	caches, err := manager.Global.GetUserCaches(ctx, userIds)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	vos := make([]*pb.VOUserInfoSimple, 0, len(caches))
	for _, cache := range caches {
		vos = append(vos, cache.VOUserInfoSimple())
	}
	return vos, nil
}

func (u *User) YggBattleGiveUp(ctx context.Context, uid int64) error {
	if u.Yggdrasil.Status == YggdrasilStateQuit {
		return errors.WrapTrace(common.ErrYggdrasilNotInTravel)
	}
	object, ok := u.Yggdrasil.Entities.FindObjectByUid(uid)
	if !ok {
		return errors.WrapTrace(common.ErrYggdrasilObjectNotFound)
	}
	if *(object.Position) != *(u.Yggdrasil.TravelPos) {
		return errors.WrapTrace(common.ErrYggdrasilObjectNotFound)
	}
	state, err := checkObjectState(object, static.YggObjectTypePassivermonster, static.YggObjectTypeInitiativemonster, static.YggObjectTypeChallengealtar)
	if err != nil {
		return errors.WrapTrace(err)
	}

	if state.ObjectType == static.YggObjectTypePassivermonster || state.ObjectType == static.YggObjectTypeChallengealtar {

		// 被动怪物和挑战祭坛
		pos := *u.Yggdrasil.TravelPos
		if coordinate.CubeDistance(*u.Yggdrasil.TravelPosBefore, *u.Yggdrasil.TravelPos) == 1 {
			pos = *u.Yggdrasil.TravelPosBefore
		} else {
			ring := coordinate.CubeRing(*u.Yggdrasil.TravelPos, 1)
			for _, position := range ring {
				if _, err := u.Yggdrasil.canMoveTo(*u.Yggdrasil.TravelPos, position); err != nil {
					pos = position
					break
				}
			}
		}
		_, _, err = u.Yggdrasil.ChangePos(ctx, u, pos)
		if err != nil {
			return errors.WrapTrace(err)
		}

	} else if state.ObjectType == static.YggObjectTypeInitiativemonster {
		err := u.SetAllDeadAndReturnSafePos(ctx)
		if err != nil {
			return errors.WrapTrace(err)
		}
	}
	return nil
}

func (u *User) SetAllDeadAndReturnSafePos(ctx context.Context) error {
	pos, err := manager.CSV.Yggdrasil.GetClosestSafePos(*u.Yggdrasil.TravelPos)
	if err != nil {
		return errors.WrapTrace(err)
	}
	_, _, err = u.Yggdrasil.ChangePos(ctx, u, *pos)
	if err != nil {
		return errors.WrapTrace(err)
	}
	u.Yggdrasil.TravelInfo.SetAllDead()
	u.Yggdrasil.TravelInfo.TravelAp = 0
	return nil
}

func (u *User) YggMonsterInitPos(objectUid int64) (*pb.VOPosition, error) {
	//if u.Yggdrasil.Status == YggdrasilStateQuit {
	//	return nil, errors.WrapTrace(common.ErrYggdrasilNotInTravel)
	//}
	object, ok := u.Yggdrasil.Entities.FindObjectByUid(objectUid)
	if !ok {
		return nil, errors.WrapTrace(common.ErrYggdrasilObjectNotFound)
	}
	_, err := checkObjectState(object, static.YggObjectTypeInitiativemonster)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return object.OrgPos.VOPosition(), nil
}

func (u *User) YggMonsterBackInitPos(ctx context.Context, objectUid int64) (*pb.VOYggdrasilObject, error) {
	//if u.Yggdrasil.Status == YggdrasilStateQuit {
	//	return nil, errors.WrapTrace(common.ErrYggdrasilNotInTravel)
	//}
	object, ok := u.Yggdrasil.Entities.FindObjectByUid(objectUid)
	if !ok {
		return nil, errors.WrapTrace(common.ErrYggdrasilObjectNotFound)
	}
	_, err := checkObjectState(object, static.YggObjectTypeInitiativemonster)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	err = u.Yggdrasil.ObjectMove(ctx, u, object, *object.OrgPos)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return object.VOYggdrasilObject(), nil
}

func (u *User) SetYggCharacterHp(ctx context.Context, charactersAfterBattle []*pb.VOBattleCharacter) {

	for _, character := range charactersAfterBattle {
		u.Yggdrasil.TravelInfo.SetCharacterHp(character.CharacterId, character.HpPercent)

	}
	if u.Yggdrasil.TravelInfo.AllDead() {
		u.Yggdrasil.TravelInfo.TravelAp = 0
	}
}
