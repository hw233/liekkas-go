package model

import (
	"context"
	"encoding/json"
	"gamesvr/manager"
	"math"
	"shared/common"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/statistic/logreason"
	"shared/utility/coordinate"
	"shared/utility/errors"
	"shared/utility/glog"
	"shared/utility/number"
)

const (
	searchBlackMaxLoop = 10000

	YggdrasilStateQuit   int32 = 0
	YggdrasilStateTravel int32 = 1
)

type Yggdrasil struct {
	*DailyRefreshChecker

	Status           int32                    `json:"status"`
	TravelTime       int32                    `json:"travel_time"`
	CityId           int32                    `json:"city_id"`             // 所在在城市
	CanTravelCityIds *number.NonRepeatableArr `json:"can_travel_city_ids"` // 可以传送的城市

	WorldId         int32                `json:"world_id"`
	TravelPos       *coordinate.Position `json:"travel_pos"`
	TravelPosBefore *coordinate.Position `json:"travel_pos_before"` // 上一步移动的位置

	TravelAreaId int32 `json:"travel_area_id"` // 探索所在区域

	TravelInfo *YggdrasilInTravelInfo `json:"travel_info"`
	Pack       *YggPack               `json:"pack"`      // 背包
	TaskPack   *YggTaskPack           `json:"task_pack"` // 任务背包

	Task     *YggdrasilTask       `json:"task"`
	Dispatch *YggdrasilDispatches `json:"yggdrasil_dispatch"`

	Areas      *YggAreas            `json:"areas"`
	UnlockArea *common.Area         `json:"unlock_area"` // 照亮的区域
	Entities   *YggEntities         `json:"entities"`    // 地图上东西（包括Builds,Marks,DiscardGoods,Objects）
	TrackMark  *coordinate.Position `json:"track_mark"`  // 追踪标记

	MailBox *YggdrasilMailBox `json:"mail_box"`

	DailyMonsters *YggdrasilDailyMonsters `json:"daily_monster"`

	SpecialStatics *YggdrasilSpecialStatics `json:"special_statics"`
	InitObjects    *YggdrasilInitObjects    `json:"init_objects"` //记录AreaPos已经刷的object

	EntityChange *YggdrasilEntityChange `json:"-"` // 同步
	MatchUserIds map[int64]struct{}     `json:"-"` // 匹配到的用户id

}

func NewYggdrasil() *Yggdrasil {
	initPos := manager.CSV.Yggdrasil.GetYggInitPos()
	ygg := &Yggdrasil{
		DailyRefreshChecker: NewDailyRefreshChecker(),
		Status:              YggdrasilStateQuit,
		TravelTime:          0,
		CityId:              0,
		CanTravelCityIds:    number.NewNonRepeatableArr(),

		WorldId:         1,
		TravelPos:       &initPos,
		TravelPosBefore: &initPos,
		TravelAreaId:    0,

		TravelInfo: NewYggdrasilInTravelInfo(),
		Pack:       NewYggPack(),
		TaskPack:   NewYggTaskPack(),
		Dispatch:   NewYggDispatches(),

		Task:  NewYggdrasilTask(),
		Areas: NewYggAreas(),

		UnlockArea: common.NewArea(),
		Entities:   NewYggEntities(),
		TrackMark:  nil,

		MailBox: NewYggdrasilMailBox(),

		DailyMonsters: NewYggdrasilDailyMonsters(),

		SpecialStatics: NewYggdrasilSpecialStatics(),

		EntityChange: NewYggdrasilEntityChange(),
		MatchUserIds: map[int64]struct{}{},
		InitObjects:  NewYggdrasilInitObjects(),
	}

	ygg.CanTravelCityIds.Append(manager.CSV.Yggdrasil.GetYggInitCity())
	return ygg

}

func (y *Yggdrasil) init(ctx context.Context) {
	areaId, err := manager.CSV.Yggdrasil.GetPosAreaId(*y.TravelPos)
	if err != nil {
		glog.Errorf("Yggdrasil init GetPosAreaId err:%+v", err)
	}
	y.TravelAreaId = areaId

	// 加载所有entity到block中（必须先于）LightAround
	y.Entities.init(y.MatchUserIds)

	y.Areas.init(y)

	y.Pack.init(y.MatchUserIds)

	// 点亮周围
	y.LightAround(ctx)

	// 初始化日常怪物
	y.DailyMonsters.init(ctx, y)

	// 刷新增cfg_yggdrasil_area_pos的obj
	y.RefreshAreaPos(ctx)
	// 清空资源变动
	y.EntityChange.clear()

}

func (y *Yggdrasil) RefreshAreaPos(ctx context.Context) {
	for _, area := range y.Areas.Areas {
		y.InitAreaInitObj(ctx, area.AreaId)
	}
}

// InitAreaInitObj 初始化地图上obj
func (y *Yggdrasil) InitAreaInitObj(ctx context.Context, areaId int32) {
	areaConfig, err := manager.CSV.Yggdrasil.GetArea(areaId)
	if err != nil {
		glog.Errorf("InitAreaInitObj GetArea error:%+v", err)
		return
	}
	for _, obj := range areaConfig.InitObjs {
		if y.InitObjects.Contains(*obj.Position) {
			continue
		}
		objectConfig, err := manager.CSV.Yggdrasil.GetYggdrasilObjectConfig(obj.ObjectId)
		if err != nil {
			glog.Errorf("InitAreaInitObj GetYggdrasilObjectConfig error:%+v", err)
			continue

		}
		object, err := NewYggObject(ctx, *obj.Position, obj.ObjectId, objectConfig.DefaultState)
		if err != nil {
			glog.Errorf("InitAreaInitObj NewYggObject error:%+v", err)
			continue

		}
		y.Entities.AppendObject(object)
		y.InitObjects.Append(*obj.Position)

	}

}
func (y *Yggdrasil) InitAreaCoBuilds(ctx context.Context, areaId int32) {
	coBuilds := manager.CSV.Yggdrasil.GetAllYggCoBuild(areaId)
	for _, build := range coBuilds {
		newBuild, err := NewYggCoBuild(ctx, *build.Pos)
		if err != nil {
			glog.Errorf("InitAreaCoBuilds NewCoBuild error:%+v", err)
			return
		}
		y.Entities.AppendCoBuild(newBuild)
	}
}

// AppendOwnBuild 添加自建的建筑
func (y *Yggdrasil) AppendOwnBuild(ctx context.Context, userId int64, build *YggBuild) error {

	area := y.Areas.getByCreate(ctx, y, build.AreaId)
	areaCsv, err := manager.CSV.Yggdrasil.GetArea(build.AreaId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	// 检查当前种类的建筑是否超出最大建筑数量的限制

	buildCfg, err := manager.CSV.Yggdrasil.GetYggdrasilBuilding(build.BuildId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	if area.BuildCreateCount[build.BuildId] > areaCsv.GetMaxBuildBuildCount(buildCfg.BuildingType) {
		return errors.Swrapf(common.ErrYggdrasilBuildCountOutOfLimit, buildCfg.BuildingType)
	}
	err = y.Entities.AppendBuild(ctx, userId, build)
	if err != nil {
		return errors.WrapTrace(err)
	}
	area.BuildCreateCount[build.BuildId] = area.BuildCreateCount[build.BuildId] + 1
	return nil
}

// RemoveBuild 移除建筑
func (y *Yggdrasil) RemoveBuild(ctx context.Context, userId int64, build *YggBuild, updateBuildCreateCount bool) error {
	area := y.Areas.getByCreate(ctx, y, build.AreaId)
	err := y.Entities.RemoveBuild(ctx, userId, build)
	if err != nil {
		return errors.WrapTrace(err)
	}
	if updateBuildCreateCount && build.IsOwn() {
		area.BuildCreateCount[build.BuildId] -= 1
	}
	return nil
}

// AppendOwnMessage 添加留言
func (y *Yggdrasil) AppendOwnMessage(ctx context.Context, userId int64, msg *YggMessage) error {
	area := y.Areas.getByCreate(ctx, y, msg.AreaId)
	areaCsv, err := manager.CSV.Yggdrasil.GetArea(msg.AreaId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	if area.MessageCreateCount > areaCsv.MessageMaxCreateCount {
		return errors.Swrapf(common.ErrYggdrasilMessageCountOutOfLimit)
	}
	err = y.Entities.AppendMessage(ctx, userId, msg)
	if err != nil {
		return errors.WrapTrace(err)
	}
	area.MessageCreateCount += 1
	return nil
}

func (y *Yggdrasil) UpdateMessage(ctx context.Context, userId int64, msg *YggMessage, m string) error {
	return y.Entities.UpdateMessage(ctx, userId, msg, m)
}

func (y *Yggdrasil) RemoveMessage(ctx context.Context, userId int64, msg *YggMessage) error {
	area := y.Areas.getByCreate(ctx, y, msg.AreaId)
	err := y.Entities.RemoveMessage(ctx, userId, msg)
	if err != nil {
		return errors.WrapTrace(err)
	}
	if msg.IsOwn() {
		area.MessageCreateCount -= 1
	}
	return nil
}

func (y *Yggdrasil) RecoverFromBefore(now, before *YggObject) {
	y.Entities.RemoveObject(before)
	now.State = before.State
	now.X = before.X
	now.Y = before.Y
	y.Entities.AppendObject(now)
	y.EntityChange.AppendObject(now)
}

func (y *Yggdrasil) RemoveObject(object *YggObject) {
	y.Entities.RemoveObject(object)
	y.EntityChange.AppendRemovedObjectUid(object.Uid)

}

// AppendObject 添加互动物
func (y *Yggdrasil) AppendObject(ctx context.Context, object *YggObject) {

	y.Areas.getByCreate(ctx, y, object.AreaId)
	y.Entities.AppendObject(object)
	y.EntityChange.AppendObject(object)

}

// AppendMark 添加标记
func (y *Yggdrasil) AppendMark(ctx context.Context, mark *YggMark) {

	y.Areas.getByCreate(ctx, y, mark.AreaId)
	y.Entities.AppendMark(mark)
}

// RemoveMark 移除标记
func (y *Yggdrasil) RemoveMark(mark *YggMark) {
	y.Entities.RemoveMark(mark)
}

func (y *Yggdrasil) VOYggdrasilUser(u *User) *pb.VOYggdrasilUser {
	return &pb.VOYggdrasilUser{
		Status:          y.Status,
		Ap:              u.Info.Ap.Value(),
		ApRefreshAt:     u.Info.Ap.Last(),
		CityId:          y.CityId,
		TodayTravelTime: y.TravelTime,
	}
}

func (y *Yggdrasil) VOYggdrasilBlockAndArea(ctx context.Context, corners []coordinate.Position) *pb.VOYggdrasilBlockAndArea {
	var BlockList []*pb.VOYggdrasilBlock
	var AreaList []*pb.VOYggdrasilArea
	// block所包含的区域id
	areaList := number.NewNonRepeatableArr()

	for _, corner := range corners {
		block := y.Entities.Blocks.getByCreate(corner)
		areaList.Append(block.Areas.Values()...)
	}
	// 因为此时block可能没有初始化 ，先尝试组area 数据，初始化area时候会初始化block，
	for _, areaId := range areaList.Values() {
		area := y.Areas.getByCreate(ctx, y, areaId)
		AreaList = append(AreaList, area.VOYggdrasilArea())
	}
	for _, corner := range corners {
		block := y.Entities.Blocks.getByCreate(corner)
		BlockList = append(BlockList, block.VOYggdrasilBlock(ctx, y.UnlockArea))
	}
	return &pb.VOYggdrasilBlockAndArea{
		BlockList: BlockList,
		AreaList:  AreaList,
	}
}
func (y *Yggdrasil) VOYggdrasilTravelInfo() *pb.VOYggdrasilTravelInfo {
	if y.Status == YggdrasilStateQuit {
		return nil
	}
	return y.TravelInfo.VOYggdrasilTravelInfo()

}

func (y *Yggdrasil) VOYggdrasilTaskTotalInfo() *pb.VOYggdrasilTaskTotalInfo {
	return y.Task.VOYggdrasilTaskTotalInfo()
}

// GetBlockInfoByEnter 进入时候获得地块信息
func (y *Yggdrasil) GetBlockInfoByEnter(ctx context.Context) *pb.VOYggdrasilBlockAndArea {
	corners := calNineCorner(*y.TravelPos)
	return y.VOYggdrasilBlockAndArea(ctx, corners)
}

func (y *Yggdrasil) initTravelInfo(characters []int32, ap int32) error {

	y.Status = YggdrasilStateTravel
	y.TravelTime++
	CharactersHp := map[int32]int32{}

	for _, character := range characters {
		CharactersHp[character] = 10000
	}
	y.TravelInfo.TravelAp = ap
	y.TravelInfo.CharactersHp = CharactersHp
	// 初始化建筑使用次数
	for _, build := range y.Entities.Builds {
		build.UseCount = 0
	}
	return nil
}

func (y *Yggdrasil) VOYggdrasilPosition() *pb.VOYggdrasilPosition {
	if y.CityId != 0 {
		return nil
	}

	return &pb.VOYggdrasilPosition{
		Position: y.TravelPos.VOPosition(),
		WorldId:  y.WorldId,
	}
}

// 获取一个点 周围九宫格block corner
func calNineCorner(position coordinate.Position) []coordinate.Position {
	h, v := manager.CSV.Yggdrasil.GetYggBlockLengthAndWidth()

	center := calCorner(position)
	var result []coordinate.Position
	result = append(result, center)
	result = append(result, *coordinate.NewPosition(center.X, center.Y+v))
	result = append(result, *coordinate.NewPosition(center.X, center.Y-v))

	result = append(result, *coordinate.NewPosition(center.X-h, center.Y))
	result = append(result, *coordinate.NewPosition(center.X-h, center.Y+v))
	result = append(result, *coordinate.NewPosition(center.X-h, center.Y-v))

	result = append(result, *coordinate.NewPosition(center.X+h, center.Y))
	result = append(result, *coordinate.NewPosition(center.X+h, center.Y+v))
	result = append(result, *coordinate.NewPosition(center.X+h, center.Y-v))

	return result
}

func calCorner(position coordinate.Position) coordinate.Position {
	h, v := manager.CSV.Yggdrasil.GetYggBlockLengthAndWidth()
	return *coordinate.NewPosition(int32(math.Floor(float64(position.X)/float64(h)))*h, int32(math.Floor(float64(position.Y)/float64(v)))*v)
}

// 移动
func (y *Yggdrasil) moveTo(ctx context.Context, u *User, target coordinate.Position) ([]*pb.VOPosition, []*pb.VOExploredPosCountUpdate, error) {
	from := *y.TravelPos
	//objects, ok := y.Entities.FindObjectsByPos(from)
	//if ok {
	//	for _, object := range objects {
	//		objectType, err := manager.CSV.Yggdrasil.GetYggdrasilObjectType(object.ObjectId, object.State)
	//		if err != nil {
	//			return nil, nil, errors.WrapTrace(err)
	//		}
	//		// 主动怪物和被动怪物
	//		if objectType == static.YggObjectTypePassivermonster || objectType == static.YggObjectTypeInitiativemonster {
	//			return nil, nil, errors.WrapTrace(common.ErrYggdrasilCannotMoveMonster)
	//		}
	//	}
	//}

	moveCost, err := y.canMoveTo(from, target)
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}

	if moveCost > y.TravelInfo.TravelAp {
		return nil, nil, errors.WrapTrace(common.ErrYggdrasilApNotEnough)
	}
	y.TravelInfo.TravelAp -= moveCost
	y.SpecialStatics.TriggerMoveCount()
	return y.ChangePos(ctx, u, target)

}

func (y *Yggdrasil) IfElevator(pos coordinate.Position) bool {
	build, ok := y.Entities.FindBuildingsByPos(pos)
	if !ok {
		return false
	}
	buildingCfg, err := manager.CSV.Yggdrasil.GetYggdrasilBuilding(build.BuildId)
	if err != nil {
		return false
	}
	//todo: 次元电梯type
	if buildingCfg.BuildingType != static.YggBuildingTypeStepladder {
		return false
	}
	return true
}

// 返回移动移动消耗
func (y *Yggdrasil) canMoveTo(from, to coordinate.Position) (int32, error) {
	distance := coordinate.CubeDistance(from, to)
	if distance != 1 {
		return 0, errors.WrapTrace(common.ErrYggdrasilCannotMoveDistanceIllegal)
	}
	// 判断是否可行走区域
	err := y.Entities.WalkablePos(to)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}

	if y.IfElevator(from) || y.IfElevator(to) {
		return 0, nil
	}
	// 有地势差时处理
	if y.TerrainDiff(from, to) {
		return 0, errors.WrapTrace(common.ErrYggdrasilCannotMoveTerrainDiff)

	}
	moveCost, err := manager.CSV.Yggdrasil.GetMostCost(to)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}
	return moveCost, nil
}

// TerrainDiff 是否有地势差
func (y *Yggdrasil) TerrainDiff(from, to coordinate.Position) bool {
	// 判断地形
	fromHeight := y.Entities.GetPosHeight(from)
	toHeight := y.Entities.GetPosHeight(to)
	if fromHeight < toHeight {
		return true
	}
	if fromHeight > toHeight {
		if fromHeight-toHeight > 1 {
			return true
		}
	}
	return false
}

func (y *Yggdrasil) DeletePackGoods(goodsIds ...int64) {
	for _, goodsId := range goodsIds {
		y.Pack.Delete(goodsId)
		y.EntityChange.AppendDeletePackGoods(goodsId)
	}

}

func (y *Yggdrasil) PickUp(discardGoods *YggDiscardGoods) error {

	packGoods, err := NewYggPackGoodsByDiscard(discardGoods)
	if err != nil {
		return errors.WrapTrace(err)
	}
	y.AddPackGoods(packGoods)
	return nil
}
func (y *Yggdrasil) AddPackGoods(packGoods *YggPackGoods) {

	y.Pack.Add(packGoods)
	y.EntityChange.AppendPackGoods(packGoods)
}

// AppendDiscardGoods 添加丢弃物
func (y *Yggdrasil) AppendDiscardGoods(ctx context.Context, userId int64, discardGoods *YggDiscardGoods) error {

	y.Areas.getByCreate(ctx, y, discardGoods.AreaId)
	err := y.Entities.AppendDiscardGoods(ctx, userId, discardGoods)
	if err != nil {
		return errors.WrapTrace(err)
	}
	y.EntityChange.AppendDiscardGoods(discardGoods)
	return nil
}

func (y *Yggdrasil) RemoveDiscardGoods(ctx context.Context, userId int64, discardGoods *YggDiscardGoods) error {
	err := y.Entities.RemoveDiscardGoods(ctx, userId, discardGoods)
	if err != nil {
		return errors.WrapTrace(err)
	}
	y.EntityChange.AppendDeleteDiscardGoods(discardGoods.Uid)
	return nil
}

func (y *Yggdrasil) ObjectChangeToNextState(ctx context.Context, u *User, object *YggObject) {
	state, err := manager.CSV.Yggdrasil.GetObjectState(object.ObjectId, object.State)
	if err != nil {
		glog.Errorf("ObjectChangeToNextState GetObjectState err:%+v", err)
		return
	}
	// 记录obj原始state和坐标，任务放弃或者失败的时候重置成操作前的状态
	y.Task.AppendTemporalObject(object)
	if state.NextState > 0 {
		y.Entities.SetObjectNextState(object, state.NextState)
		y.EntityChange.AppendObject(object)
		y.Task.ProcessNum(ctx, u, static.YggdrasilSubTaskTypeTypeObjectStateChange, state.NextState, object.ObjectId)

	} else {
		y.RemoveObject(object)
	}

}

func (y *Yggdrasil) ObjectMove(ctx context.Context, u *User, obj *YggObject, to coordinate.Position) error {
	areaId, err := manager.CSV.Yggdrasil.GetPosAreaId(to)
	if err != nil {
		return errors.WrapTrace(err)
	}

	// 记录obj原始state和坐标，任务放弃或者失败的时候重置成操作前的状态
	y.Task.AppendTemporalObject(obj)

	y.RemoveObject(obj)
	obj.Position = &to
	obj.AreaId = areaId
	y.AppendObject(ctx, obj)
	y.Task.ProcessPos(ctx, u, static.YggdrasilSubTaskTypeTypeConvoy, to, obj.ObjectId)
	y.Task.ProcessPos(ctx, u, static.YggdrasilSubTaskTypeTypeLeadWay, to, obj.ObjectId)
	// object移动不走 EntityChange
	y.EntityChange.clear()
	return nil
}

// AddRewardsByDropId 通过dropId掉落
func (y *Yggdrasil) AddRewardsByDropId(ctx context.Context, user *User, dropId int32, subtaskId int32,
	reason *logreason.Reason) error {
	rewards, err := manager.CSV.Drop.DropRewards(dropId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	return y.AddRewards(ctx, user, rewards, subtaskId, reason)

}

func (y *Yggdrasil) AddRewardsByDropIds(ctx context.Context, user *User, dropIds []int32, subtaskId int32,
	reason *logreason.Reason) error {
	for _, dropId := range dropIds {
		err := y.AddRewardsByDropId(ctx, user, dropId, subtaskId, reason)
		if err != nil {
			return err
		}
	}

	return nil
}

// 奖励分割成三个 1.发到玩家背包 2.发到世界探索内 3.发到世界探索任务背包
func (y *Yggdrasil) separateReward(rewards *common.Rewards) (*common.Rewards, *common.Rewards, *common.Rewards) {

	// 发到玩家背包
	rewardsInUser := common.NewRewards()
	// 发到世界探索内
	rewardsInYggdrasil := common.NewRewards()
	// 发到世界探索任务背包
	rewardsInTaskPack := common.NewRewards()

	for _, reward := range rewards.MergeValue() {
		if reward.Type == static.ItemTypeYggItemSpecial {
			rewardsInTaskPack.AddReward(&reward)
		} else {
			// 如果在city中，那么应该除了任务道具其他都发到玩家背包
			if y.CityId > 0 {
				rewardsInUser.AddReward(&reward)

			} else {
				if reward.Type != static.ItemTypeCurrency &&
					reward.Type != static.ItemTypeCharacterAvatarAuto &&
					reward.Type != static.ItemTypeCharacterAvatar &&
					reward.Type != static.ItemTypeCharacterFrame &&
					reward.Type != static.ItemTypeYggPrestige &&
					reward.Type != static.ItemTypeCharacter {
					rewardsInYggdrasil.AddReward(&reward)

				} else {
					rewardsInUser.AddReward(&reward)

				}
			}
		}
	}
	return rewardsInUser, rewardsInYggdrasil, rewardsInTaskPack

}

// AddRewards 发奖
func (y *Yggdrasil) AddRewards(ctx context.Context, user *User, rewards *common.Rewards, subtaskId int32,
	reason *logreason.Reason) error {

	rewardsInUser, rewardsInYggdrasil, rewardsInTaskPack := y.separateReward(rewards)
	err := y.AddRewardsInYgg(ctx, user, rewardsInYggdrasil)
	if err != nil {
		return errors.WrapTrace(err)
	}

	_, err = user.addRewards(rewardsInUser, reason)
	if err != nil {
		return errors.WrapTrace(err)
	}
	y.AddRewardsInTaskPack(ctx, user, rewardsInTaskPack, subtaskId)
	return nil
}

func (y *Yggdrasil) AddRewardsInYgg(ctx context.Context, user *User, rewards *common.Rewards) error {

	var r int32 = 0
	var cubeRing []coordinate.Position
	mergedRewards := rewards.MergeValue()
	for _, reward := range mergedRewards {
		for j := 0; j < int(reward.Num); j++ {

			if !y.Pack.IsFull(user.Info.Level.Value()) {
				byReward, err := NewYggPackGoods(ctx, reward.ID)
				if err != nil {
					return errors.WrapTrace(err)
				}
				y.AddPackGoods(byReward)
			} else {
				// 掉落地上
				for i := 0; i < searchBlackMaxLoop; i++ {

					if len(cubeRing) == 0 {
						cubeRing = append(cubeRing, coordinate.CubeRing(*y.TravelPos, r)...)
						r++
					} else {
						position := cubeRing[0]
						cubeRing = cubeRing[1:]
						if y.Entities.IsBlankPos(position) {
							byReward, err := NewYggDiscardGoodsByReward(ctx, position, reward.ID)
							if err != nil {
								return errors.WrapTrace(err)
							}
							err = y.AppendDiscardGoods(ctx, user.ID, byReward)
							if err != nil {
								return errors.WrapTrace(err)
							}
							break
						}
					}
				}

			}
		}
		user.RewardsResult.AddReward(&reward)
	}

	return nil
}

func (y *Yggdrasil) AddRewardsInTaskPack(ctx context.Context, user *User, rewards *common.Rewards, subtaskId int32) {
	for _, reward := range rewards.MergeValue() {
		num := y.TaskPack.Add(reward.ID, reward.Num, subtaskId)
		y.EntityChange.TaskItemChange(reward.ID, num)
		// 任务监听
		y.Task.ProcessNum(ctx, user, static.YggdrasilSubTaskTypeTypeOwn, num, reward.ID)
		user.RewardsResult.AddReward(&reward)
	}

}

// MergePack 合并三个背包 1.玩家背包 2.世界探索内背包(不包含匹配的) 3.世界探索任务背包
func (y *Yggdrasil) MergePack(u *User) *ItemPack {

	ret := NewItemPack()
	for i, eventNumber := range *u.ItemPack {
		ret.Add(i, eventNumber.Value())
	}
	for _, goods := range *u.Yggdrasil.Pack.GetAllOwnGoods() {
		ret.Add(goods.ItemId, YggdrasilPackUnitCount)
	}

	for i, eventNumber := range *u.Yggdrasil.TaskPack.merge() {
		ret.Add(i, eventNumber.Value())

	}
	return ret

}

func (y *Yggdrasil) CheckRewardsEnough(u *User, rewards *common.Rewards) error {
	pack := y.MergePack(u)
	for _, reward := range rewards.Value() {
		if !pack.Enough(reward.ID, reward.Num) {
			return errors.Swrapf(common.ErrItemNotEnough, reward.ID, reward.Num)
		}
	}
	return nil
}

func (y *Yggdrasil) CostRewards(u *User, rewards *common.Rewards, reason *logreason.Reason) {
	// 任务道具交付的时候 1.玩家背包 2.世界探索内背包 3.世界探索任务背包 ，这三个背包的消耗顺序是啥，是321
	rewards = y.CostItemPack(u, rewards, reason)
	rewards = y.CostYggdrasilPack(rewards)
	rewards = y.CostYggdrasilTaskPack(rewards)

}

// CostItemPack 消耗背包物品 返回还应该再消耗
func (y *Yggdrasil) CostItemPack(u *User, rewards *common.Rewards, reason *logreason.Reason) *common.Rewards {
	ret := common.NewRewards()
	itemPack := u.ItemPack
	for _, reward := range rewards.MergeValue() {

		enough := itemPack.Enough(reward.ID, reward.Num)
		if enough {
			u.minusItem(reward.ID, reward.Num, reason)
		} else {
			eventNumber, ok := (*itemPack)[reward.ID]
			if ok {
				value := eventNumber.Value()
				u.minusItem(reward.ID, value, reason)
				ret.AddReward(common.NewReward(reward.ID, reward.Num-value))

			} else {
				ret.AddReward(&reward)
			}
		}
	}
	return ret
}

// CostYggdrasilPack 消耗背包 返回还应该再消耗
func (y *Yggdrasil) CostYggdrasilPack(rewards *common.Rewards) *common.Rewards {
	ret := common.NewRewards()

	ownGoods := y.Pack.GetAllOwnGoods()
	for _, reward := range rewards.MergeValue() {
		restReward, needDelete := ownGoods.CostReward(reward)
		y.DeletePackGoods(needDelete...)
		if restReward.Num > 0 {
			ret.AddReward(restReward)
		}
	}
	return ret
}

// CostYggdrasilTaskPack 消耗任务背包 返回还应该再消耗 ，不要在里面写任务监听
func (y *Yggdrasil) CostYggdrasilTaskPack(rewards *common.Rewards) *common.Rewards {
	ret := common.NewRewards()

	for _, reward := range rewards.MergeValue() {

		restReward, restNum := y.TaskPack.CostReward(reward)
		if restReward.Num != reward.Num {
			// 任务背包有变化
			y.EntityChange.TaskItemChange(reward.ID, restNum)
		}
		if restReward.Num > 0 {
			ret.AddReward(restReward)
		}
	}
	return ret
}

func (y *Yggdrasil) VOYggdrasilResourceResult() *pb.VOYggdrasilResourceResult {
	return y.EntityChange.VOYggdrasilResourceResult()
}

func (y *Yggdrasil) FindObjectById(objectId int32) ([]*YggObject, error) {
	obj, ok := y.Entities.FindObjectById(objectId)
	if !ok {
		return nil, errors.WrapTrace(errors.Swrapf(common.ErrYggdrasilObjectNotFoundWithParam, objectId))
	}
	return obj, nil
}

func (y *Yggdrasil) Marshal() ([]byte, error) {
	return json.Marshal(y)
}

func (y *Yggdrasil) Unmarshal(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	err := json.Unmarshal(data, &y)
	if err != nil {
		glog.Errorf("json.Unmarshal error: %v", err)
		return err
	}

	return nil
}

func (y *Yggdrasil) VOYggdrasilPackGoods() []*pb.VOYggdrasilPackGoods {
	// 背包
	var PackInfo []*pb.VOYggdrasilPackGoods
	for _, goods := range *y.Pack {
		PackInfo = append(PackInfo, goods.VOYggdrasilPackGoods())
	}
	return PackInfo
}

func (y *Yggdrasil) VOTaskPackInfo() []*pb.VOResource {
	return y.TaskPack.VOTaskPackInfo()

}

func (y *Yggdrasil) TransferToCity(ctx context.Context, u *User, cityId int32) error {
	_, err := manager.CSV.Yggdrasil.GetYggCityById(cityId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	err = u.Yggdrasil.TakeBackPackGood(ctx, u, u.Yggdrasil.Pack.Values()...)
	if err != nil {
		return errors.WrapTrace(err)
	}
	u.Yggdrasil.CanTravelCityIds.Append(cityId)
	u.Yggdrasil.CityId = cityId
	// 任务监听
	u.Yggdrasil.Task.ProcessNum(ctx, u, static.YggdrasilSubTaskTypeTypeCity, cityId)
	return nil
}

// TakeBackPackGood 带回背包物品
func (y *Yggdrasil) TakeBackPackGood(ctx context.Context, user *User, packGoods ...*YggPackGoods) error {
	rewards := common.NewRewards()
	for _, good := range packGoods {
		//  判断是否已被别人带回
		takenBack, err := manager.Global.FetchGoodsTakenBack(ctx, good.Uid)
		if err != nil {
			return errors.WrapTrace(err)
		}
		// 没带回的话带回
		if !takenBack {
			// 自己的物品直接发奖
			if good.IsOwn() {
				rewards.AddReward(common.NewReward(good.ItemId, YggdrasilPackUnitCount))
				err = manager.Global.SetGoodsTakenBack(ctx, good.Uid)
				if err != nil {
					return errors.WrapTrace(err)
				}
			} else {
				// 带回别人的物品
				manager.EventQueue.Push(ctx, good.MatchUserId, common.NewYggdrasilMailEvent(user.Name, good.Uid))
				// 增加亲密度
				//todo:配表
				err := user.IntimacyChange(ctx, good.MatchUserId, 1)
				if err != nil {
					return errors.WrapTrace(err)
				}

			}
		}
		y.DeletePackGoods(good.Uid)

	}

	reason := logreason.NewReason(logreason.YggReturnToCity)
	_, err := user.addRewards(rewards, reason)
	if err != nil {
		return errors.WrapTrace(err)
	}

	return nil
}

//RecoverEnv env暂时删除的entity恢复
func (y *Yggdrasil) RecoverEnv(ctx context.Context, userId int64, subTaskId int32) error {
	builds := y.Entities.EnvRemovedBuilds[subTaskId]
	for _, build := range builds {
		err := y.Entities.AppendBuild(ctx, userId, build)
		if err != nil {
			return errors.WrapTrace(err)
		}
		y.EntityChange.AppendBuild(build)
	}
	objects := y.Entities.EnvRemovedObjects[subTaskId]
	for _, object := range objects {
		y.AppendObject(ctx, object)

	}

	discardGoods := y.Entities.EnvRemovedDiscardGoods[subTaskId]
	for _, goods := range discardGoods {
		err := y.AppendDiscardGoods(ctx, userId, goods)
		if err != nil {
			return errors.WrapTrace(err)
		}
	}
	return nil

}

func (y *Yggdrasil) VOYggdrasilEntityChange(ctx context.Context) *pb.VOYggdrasilEntityChange {

	return y.EntityChange.VOYggdrasilEntityChange(ctx, y)

}

func (y *Yggdrasil) ProgressReward(ctx context.Context, u *User, areaId int32) (*pb.VOYggdrasilArea, error) {
	areaConfig, err := manager.CSV.Yggdrasil.GetArea(areaId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	area := y.Areas.getByCreate(ctx, y, areaId)
	index, err := manager.CSV.Yggdrasil.GetExploreProcessIndex(areaId, area.ExploredPosCount*100/areaConfig.PosCount)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	if area.ExploredProgressRewardIndex >= index {
		return nil, errors.WrapTrace(common.ErrYggdrasilProgressRewardBefore)
	}

	reason := logreason.NewReason(logreason.YggProgress)
	for i := area.ExploredProgressRewardIndex + 1; i <= index; i++ {
		_, err := u.AddRewardsByDropId(areaConfig.ExploredProgressDrop[i], reason)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
	}
	// 设为已领奖
	area.ExploredProgressRewardIndex = index
	return area.VOYggdrasilArea(), nil
}

type LightPosList []coordinate.Position

func NewLightPosList() *LightPosList {
	return (*LightPosList)(&[]coordinate.Position{})
}

func (l *LightPosList) Append(centerHeight int32, new coordinate.Position, yggdrasil *Yggdrasil) {
	for _, position := range *l {
		if position == new {
			return
		}
	}
	height := yggdrasil.Entities.GetPosHeight(new)
	// 7.照亮区域低处不能照高处，反之可以
	if !yggdrasil.UnlockArea.Contains(new) && centerHeight >= height {
		*l = append(*l, new)
	}

}

// LightAround 点亮周围
func (y *Yggdrasil) LightAround(ctx context.Context) ([]*pb.VOPosition, []*pb.VOExploredPosCountUpdate) {

	centerHeight := y.Entities.GetPosHeight(*y.TravelPos)
	newLightPositions := NewLightPosList()
	for _, p := range coordinate.CubeRange(*y.TravelPos, manager.CSV.Yggdrasil.GetYggLightRadius()) {
		// 如果包含城市点亮区域则照亮整个城市
		cityId, ok := manager.CSV.Yggdrasil.CityLightPos(p)
		if ok {
			city, err := manager.CSV.Yggdrasil.GetYggCityById(cityId)
			if err != nil {
				glog.Errorf("LightAround GetYggCityById,err:%+v", err)
				continue
			}
			for _, cityP := range coordinate.CubeRange(*city.CityCenterPos, city.CityRadius) {
				newLightPositions.Append(centerHeight, cityP, y)
			}
		} else {
			newLightPositions.Append(centerHeight, p, y)

		}

	}
	var vos []*pb.VOPosition

	if len(*newLightPositions) > 0 {
		// 照亮的点如果是新的area则初始化新area
		for _, p := range *newLightPositions {
			areaId, err := manager.CSV.Yggdrasil.GetPosAreaId(p)
			if err != nil {
				continue
			}
			yggArea := y.Areas.getByCreate(ctx, y, areaId)
			// 区域探索进度更新

			yggArea.ExploredPosCount++

			y.UnlockArea.AppendPoint(p)
			vos = append(vos, p.VOPosition())
		}
	}

	return vos, y.GetExploredPosCountUpdate(ctx, newLightPositions)
}

// CheckPrestigeEnough 声望是否达到
func (u *User) CheckPrestigeEnough(areaId, val int32) error {

	area, err := manager.CSV.Yggdrasil.GetArea(areaId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	if !u.ItemPack.Enough(area.PrestigeItemID, val) {
		return errors.WrapTrace(common.ErrYggdrasilPrestigeNotEnough)
	}
	return nil
}

func (y *Yggdrasil) VOTrackMarkPosition() *pb.VOPosition {
	if y.TrackMark == nil {
		return nil
	}
	return y.TrackMark.VOPosition()
}

func (y *Yggdrasil) GetMarkTotalCount() int32 {
	return int32(len(y.Entities.Marks))
}

func (y *Yggdrasil) TryTakeBackGoods(ctx context.Context, userId int64, fromUserName string, goodsUid int64) error {
	glog.Debugf("TryTakeBackGoods fromUserId:%d,fromName:%s,goodsUid:%d", userId, fromUserName, goodsUid)
	goods, ok := y.FindGoodByUid(goodsUid)
	if !ok {
		// 已经被自己带回
		glog.Debugf("TryTakeBackGoods takenBackByOwn,goodsUid:%d", goodsUid)
		return nil
	}
	takenBack, err := manager.Global.FetchGoodsTakenBack(ctx, goodsUid)
	if err != nil {
		return errors.WrapTrace(err)
	}
	if takenBack {
		// 已经被他人带回
		glog.Debugf("TryTakeBackGoods takenBackByOthers,goodsUid:%d", goodsUid)
		return nil
	}
	if y.MailBox.IsMailBoxFull() {
		// 邮箱满了
		glog.Debugf("TryTakeBackGoods MailBoxFull,goodsUid:%d", goodsUid)
		return nil
	}

	rewards := common.NewRewards()
	rewards.AddReward(common.NewReward(goods.ItemId, YggdrasilPackUnitCount))
	err = y.MailBox.AddOne(ctx, userId, fromUserName, rewards)
	if err != nil {
		return errors.WrapTrace(err)
	}
	err = manager.Global.SetGoodsTakenBack(ctx, goodsUid)
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}

func (y *Yggdrasil) FindGoodByUid(uid int64) (*YggGoods, bool) {
	discardGoods, ok := y.Entities.DiscardGoods[uid]
	if ok {
		return discardGoods.YggGoods, ok
	}
	packGoods, ok := y.Pack.Get(uid)
	if ok {
		return packGoods.YggGoods, ok
	}
	return nil, false
}
func (y *Yggdrasil) GetActivePortals() []YggPortalLocation {
	var locations []YggPortalLocation
	for _, build := range y.Entities.CoBuilds {
		buildCsv, err := manager.CSV.Yggdrasil.GetYggCoBuilding(build.BuildId)
		if err != nil {
			glog.Errorf("GetActivePortals GetYggCoBuilding err:%+v", err)
			continue
		}
		if buildCsv.BuildingType == static.YggCobuildingTypeTransferPortal && build.CheckIsActive() {

			locations = append(locations, *NewYggPortalLocation(PortalTypeCobuild, build.BuildId))
		}
	}
	for _, cityID := range y.CanTravelCityIds.Values() {
		locations = append(locations, *NewYggPortalLocation(PortalTypeCity, cityID))
	}
	return locations
}

func (y *Yggdrasil) CheckPortalActive(target YggPortalLocation) bool {

	for _, location := range y.GetActivePortals() {
		if target == location {
			return true
		}
	}

	return false
}

func (y *Yggdrasil) VOPortalLocations() []*pb.VOYggdrasilPortalLocation {
	locations := y.GetActivePortals()
	vos := make([]*pb.VOYggdrasilPortalLocation, 0, len(locations))
	for _, location := range locations {
		vos = append(vos, location.VOYggdrasilPortalLocation())
	}
	return vos
}
func (y *Yggdrasil) VOYggdrasilReturnCity() *pb.VOYggdrasilReturnCity {
	if y.CityId == 0 {
		return nil
	}
	return &pb.VOYggdrasilReturnCity{
		CityId:                  y.CityId,
		YggdrasilResourceResult: y.VOYggdrasilResourceResult(),
		PortalLocations:         y.VOPortalLocations(),
	}
}

func (y *Yggdrasil) ChangePos(ctx context.Context, u *User, position coordinate.Position) ([]*pb.VOPosition, []*pb.VOExploredPosCountUpdate, error) {
	// 更新区域
	areaId, err := manager.CSV.Yggdrasil.GetPosAreaId(position)
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}
	y.TravelPosBefore = coordinate.NewPosition(y.TravelPos.X, y.TravelPos.Y)
	y.TravelPos = &position
	y.TravelAreaId = areaId

	y.Task.ProcessPos(ctx, u, static.YggdrasilSubTaskTypeTypeMove, *y.TravelPos)
	// 点亮周围
	around, updates := y.LightAround(ctx)
	return around, updates, nil
}

// EnvRemove env暂时移除该位置上的entity
func (y *Yggdrasil) EnvRemove(ctx context.Context, userId int64, position coordinate.Position, subTaskId int32, clearAllTypeObjects bool) error {
	// 建筑暂时移除
	building, ok := y.Entities.FindBuildingsByPos(position)
	if ok {
		//  暂时移除的build 该区域build建造数量不变（因为会恢复）
		err := y.RemoveBuild(ctx, userId, building, false)
		if err != nil {
			return errors.WrapTrace(err)
		}
		m, ok := y.Entities.EnvRemovedBuilds[subTaskId]
		if !ok {
			m = map[int64]*YggBuild{}
			y.Entities.EnvRemovedBuilds[subTaskId] = m
		}
		m[building.Uid] = building
		y.EntityChange.AppendRemovedBuildingUid(building.Uid)
	}
	// obj暂时移除
	objects, ok := y.Entities.FindObjectsByPos(position)
	if ok {
		for _, object := range objects {
			if clearAllTypeObjects {
				y.RemoveObject(object)
				m, ok := y.Entities.EnvRemovedObjects[subTaskId]
				if !ok {
					m = map[int64]*YggObject{}
					y.Entities.EnvRemovedObjects[subTaskId] = m
				}
				m[object.Uid] = object

			} else {
				//清理地块只清理怪物
				_, err := checkObjectState(object, static.YggObjectTypePassivermonster, static.YggObjectTypeInitiativemonster)
				if err == nil {
					y.RemoveObject(object)
					m, ok := y.Entities.EnvRemovedObjects[subTaskId]
					if !ok {
						m = map[int64]*YggObject{}
						y.Entities.EnvRemovedObjects[subTaskId] = m
					}
					m[object.Uid] = object

				}
			}
		}
	}

	// 丢弃物暂时移除
	discardGoods, ok := y.Entities.FindDiscardGoodsByPos(position)
	if ok {
		err := y.RemoveDiscardGoods(ctx, userId, discardGoods)
		if err != nil {
			return errors.WrapTrace(err)
		}

		m, ok := y.Entities.EnvRemovedDiscardGoods[subTaskId]
		if !ok {
			m = map[int64]*YggDiscardGoods{}
			y.Entities.EnvRemovedDiscardGoods[subTaskId] = m
		}
		m[discardGoods.Uid] = discardGoods
	}
	//  留言 不处理
	return nil
}

func (y *Yggdrasil) GetExploredPosCount(areaId int32) int32 {
	area, ok := y.Areas.get(areaId)
	if !ok {
		return 0
	}
	return area.ExploredPosCount
}

func (y *Yggdrasil) GetExploredPosCountUpdate(ctx context.Context, posList *LightPosList) []*pb.VOExploredPosCountUpdate {
	m := map[int32]int32{}
	for _, position := range *posList {
		areaId, err := manager.CSV.Yggdrasil.GetPosAreaId(position)
		if err != nil {
			continue
		}
		area := y.Areas.getByCreate(ctx, y, areaId)
		m[area.AreaId] = area.ExploredPosCount
	}
	vos := make([]*pb.VOExploredPosCountUpdate, 0, len(m))
	for k, v := range m {
		vos = append(vos, &pb.VOExploredPosCountUpdate{
			AreaId:           k,
			ExploredPosCount: v,
		})
	}
	return vos
}

func (y *Yggdrasil) VOYggSpecialStatics() *pb.VOYggSpecialStatics {
	return y.SpecialStatics.VOYggSpecialStatics()
}

func (y *Yggdrasil) GetBattleNpc() []*YggdrasilBattleNpc {
	ret := make([]*YggdrasilBattleNpc, 0, len(y.Entities.NpcHp))
	for npcId, hp := range y.Entities.NpcHp {
		ret = append(ret, NewYggdrasilBattleNpc(npcId, hp))
	}
	return ret
}

// ReEnter 重新进界面的处理
func (y *Yggdrasil) ReEnter(ctx context.Context, user *User) {
	y.EntityChange.Clear()
	// 重置跟随npc
	for _, object := range y.Entities.Objects {
		state, err := manager.CSV.Yggdrasil.GetObjectState(object.ObjectId, object.State)
		if err != nil {
			glog.Errorf("GetObjectState err :%v", err)
			continue
		}
		if state.ObjectType == static.YggObjectTypeFollownpc || state.ObjectType == static.YggObjectTypeFollowbattlenpc {
			if *y.TravelPos != *object.Position {
				err = y.ObjectMove(ctx, user, object, *y.TravelPos)
				if err != nil {
					glog.Errorf("ObjectMove err :%v", err)
					continue
				}
			}
		}
	}
}
