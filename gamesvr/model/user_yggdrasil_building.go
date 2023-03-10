package model

import (
	"context"
	"gamesvr/manager"
	"shared/common"
	"shared/csv/entry"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/statistic/logreason"
	"shared/utility/coordinate"
	"shared/utility/errors"
	"shared/utility/servertime"
)

const (
	EntityTypeBuilding   = "building"
	EntityTypeObject     = "object"
	EntityTypeMessage    = "message"
	EntityTypeCoBuilding = "co_building"
	PortalTypeCity       = 1
	PortalTypeCobuild    = 2
)

// 建造的时候是否只需要判断距离玩家当前区域所有所在城市的中心半径即可
// appendbuild的时候，按照build的id来算还是按照类型

func (u *User) YggdrasilBuildCreate(ctx context.Context, buildID int32) (*YggBuild, error) {

	// 1. 获得对应建筑的csv数据，检查解锁
	buildCsv, err := manager.CSV.Yggdrasil.GetYggdrasilBuilding(buildID)
	if err != nil {
		return nil, err
	}

	if buildCsv.BuildingType == static.YggBuildingTypeStepladder {
		err = u.Yggdrasil.Entities.CheckElevatorCanBuild(*u.Yggdrasil.TravelPos)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
	}
	err = u.CheckUserConditions(buildCsv.UnlockConditions)
	if err != nil {
		return nil, err
	}

	// 2. 检查是否可以建造，根据建造范围限制
	isBlank := u.Yggdrasil.Entities.IsBlankPos(*u.Yggdrasil.TravelPos)
	if !isBlank {
		return nil, errors.Swrapf(common.ErrYggdrasilPosAlreadyHasEntity, u.Yggdrasil.TravelPos.X, u.Yggdrasil.TravelPos.Y)
	}

	if !u.Yggdrasil.Entities.BuildablePos(*u.Yggdrasil.TravelPos) {
		return nil, errors.Swrapf(common.ErrYggdrasilPosTypeForbidBuild, u.Yggdrasil.TravelPos.X, u.Yggdrasil.TravelPos.Y)
	}

	// 检查主城的不可建造半径
	areaId, err := manager.CSV.Yggdrasil.GetPosAreaId(*u.Yggdrasil.TravelPos)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	area, err := manager.CSV.Yggdrasil.GetArea(areaId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	for _, cityId := range area.AreaCityIds {
		city, err := manager.CSV.Yggdrasil.GetYggCityById(cityId)
		if err != nil {
			return nil, err
		}
		dis := coordinate.CubeDistance(*u.Yggdrasil.TravelPos, *city.CityCenterPos)
		if dis <= city.CityBanR {
			return nil, errors.Swrapf(common.ErrYggdrasilPosInsideCityBanRadius, u.Yggdrasil.TravelPos.X, u.Yggdrasil.TravelPos.Y)
		}
	}

	// 检查周围有无同type建筑
	err = u.Yggdrasil.Entities.CheckSameTypeBuildingAround(*u.Yggdrasil.TravelPos, buildCsv.BuildingType, buildCsv.BuildingR)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 3. 检查钱是否足够
	buildingCosts, ok := buildCsv.SingleCost[areaId]
	if !ok {
		return nil, errors.Swrapf(common.ErrYggdrasilBuildCSVNoFindAreaCost, buildCsv.Id, areaId)
	}
	// buildingCosts := buildCsv.SingleCosts
	err = u.CheckRewardsEnough(buildingCosts)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 4. 如果可以建造，则new一个建筑, 并把usecount置0
	newBuild, err := NewYggBuild(ctx, buildID, *u.Yggdrasil.TravelPos)
	if err != nil {
		return nil, err
	}

	// 5. 更新数据库数据
	err = u.Yggdrasil.AppendOwnBuild(ctx, u.ID, newBuild)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	reason := logreason.NewReason(logreason.YggBuild)
	err = u.CostRewards(buildingCosts, reason)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 任务监听
	u.Yggdrasil.Task.ProcessNum(ctx, u, static.YggdrasilSubTaskTypeTypeBuild, 1, buildID)
	// 6. 返回值
	return newBuild, nil
}

func (u *User) YggdrasilBuildDestroy(ctx context.Context, uid int64) (*YggBuild, error) {

	build, ok := u.Yggdrasil.Entities.Builds[uid]
	if !ok {
		return nil, errors.Swrapf(common.ErrYggdrasilWrongUIDForEntity, EntityTypeBuilding, uid)
	}

	err := u.Yggdrasil.RemoveBuild(ctx, u.ID, build, true)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// 减少亲密度
	if !build.IsOwn() {
		//todo:配表
		err := u.IntimacyChange(ctx, build.MatchUserId, -1)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
	}

	return build, nil
}

func (u *User) YggdrasilBuildBeforeUse() (*YggBuild, *entry.YggdrasilBuilding, error) {

	// 1. 获得建筑
	build, ok := u.Yggdrasil.Entities.FindBuildingsByPos(*u.Yggdrasil.TravelPos)
	if !ok {
		return nil, nil, errors.Swrapf(common.ErrYggdrasilPosHasNoEntityToUse, EntityTypeBuilding, u.Yggdrasil.TravelPos.X, u.Yggdrasil.TravelPos.Y)
	}

	// 2. 获得对应建筑的csv数据，检查
	buildCsv, err := manager.CSV.Yggdrasil.GetYggdrasilBuilding(build.BuildId)
	if err != nil {
		return nil, nil, err
	}

	// if buildCsv.BuildingType != static.YggBuildingTypeApSpring {
	// 	return nil, nil, errors.Swrapf(common.ErrYggdrasilBuildNoMatchProtocol)
	// }

	err = u.CheckUserConditions(buildCsv.UnlockConditions)
	if err != nil {
		return nil, nil, err
	}

	if build.UseCount >= buildCsv.UsingTimes {
		return nil, nil, errors.Swrapf(common.ErrYggdrasilUseCountOutLimit, build.UseCount)
	}

	// 3. 检查钱是否足够
	usingCosts := buildCsv.UsingCost
	err = u.CheckRewardsEnough(usingCosts)
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}

	return build, buildCsv, nil
}

func (u *User) YggdrasilBuildAfterUse(ctx context.Context, build *YggBuild, buildCsv *entry.YggdrasilBuilding) (*YggBuild, error) {
	err := build.AfterUse(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	u.Yggdrasil.Entities.Builds[build.Uid] = build

	reason := logreason.NewReason(logreason.YggUsingBuilding)
	err = u.CostRewards(buildCsv.UsingCost, reason)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// 增加亲密度
	if !build.IsOwn() {
		//todo:配表
		err := u.IntimacyChange(ctx, build.MatchUserId, 1)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
	}
	u.Guild.AddTaskItem(static.GuildTaskYggbuild, 1)
	return build, nil
}

func (u *User) YggdrasilBuildAddAp(ctx context.Context) (*pb.VOYggdrasilBuild, int32, error) {
	build, buildCsv, err := u.YggdrasilBuildBeforeUse()
	if err != nil {
		return nil, 0, errors.WrapTrace(err)
	}

	// 4. 使用活力之泉，人物精力增加，可使用次数减一， 减钱
	u.Yggdrasil.TravelInfo.TravelAp += buildCsv.UsingParam[0]

	buildAfterUse, err := u.YggdrasilBuildAfterUse(ctx, build, buildCsv)
	if err != nil {
		return nil, 0, errors.WrapTrace(err)
	}

	// 6. 返回UseBuild{VOYggdrasilBuild, resourceResult}, ap
	voYggdrasilBuild, err := buildAfterUse.VOYggdrasilBuild(ctx)
	if err != nil {
		return nil, 0, errors.WrapTrace(err)
	}
	return voYggdrasilBuild, u.Yggdrasil.TravelInfo.TravelAp, nil
}

func (u *User) YggdrasilBuildAddHp(ctx context.Context) (*pb.VOYggdrasilBuild, []*pb.VOYggdrasilCharacter, error) {

	build, buildCsv, err := u.YggdrasilBuildBeforeUse()
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}

	// 4. 使用生命祭坛， 角色血量增加，可使用次数减一，减钱
	u.Yggdrasil.TravelInfo.AddHpForAll(buildCsv.UsingParam[0])
	buildAfterUse, err := u.YggdrasilBuildAfterUse(ctx, build, buildCsv)
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}
	// 6. 返回值
	voYggdrasilBuild, err := buildAfterUse.VOYggdrasilBuild(ctx)
	if err != nil {
		return nil, nil, errors.WrapTrace(err)
	}
	return voYggdrasilBuild, u.Yggdrasil.TravelInfo.VOYggdrasilCharacterHp(), nil
}

func (u *User) YggdrasilBuildUsePeeing(ctx context.Context) (*pb.VOYggdrasilBuild, error) {

	build, buildCsv, err := u.YggdrasilBuildBeforeUse()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 4. 使用凝视魔塔， 解除伪装，可使用次数减一，减钱
	radius := buildCsv.UsingParam[0]
	positions := coordinate.CubeRange(*u.Yggdrasil.TravelPos, radius)
	for _, p := range positions {
		b, err := u.Yggdrasil.Entities.FindObjectByPosAndType(p, static.YggObjectTypeMysticmagic)
		if err != nil {
			continue
		}
		// 如果类型为9，则更新object为下一个状态
		u.Yggdrasil.ObjectChangeToNextState(ctx, u, b)
	}

	buildAfterUse, err := u.YggdrasilBuildAfterUse(ctx, build, buildCsv)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 6. 返回值
	voYggdrasilBuild, err := buildAfterUse.VOYggdrasilBuild(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return voYggdrasilBuild, nil
}

func (u *User) YggdrasilBuildUseTransPort(ctx context.Context, goodsIdList []int64) (*pb.VOYggdrasilBuild, error) {

	build, buildCsv, err := u.YggdrasilBuildBeforeUse()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 4. 使用，把物品都加入对应玩家的仓库，可使用次数减一，减钱
	// 根据goodsIdList找到对应goods，调用TakeBackPackGood函数
	if len(goodsIdList) > int(buildCsv.UsingParam[0]) {
		return nil, errors.Swrapf(common.ErrYggdrasilBuildUseOutOfLimit, build.BuildId)
	}
	packGoods := make([]*YggPackGoods, 0, len(goodsIdList))
	for _, goodsId := range goodsIdList {
		goods, ok := u.Yggdrasil.Pack.Get(goodsId)
		if !ok {
			return nil, errors.Swrapf(common.ErrYggdrasilWrongGoodsIDForPack, goodsId)
		}
		packGoods = append(packGoods, goods)
	}

	err = u.Yggdrasil.TakeBackPackGood(ctx, u, packGoods...)
	if err != nil {
		return nil, err
	}

	buildAfterUse, err := u.YggdrasilBuildAfterUse(ctx, build, buildCsv)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 6. 返回值
	return buildAfterUse.VOYggdrasilBuild(ctx)
}

//// todo
//func (u *User) YggdrasilBuildUseStepladder() (*pb.VOYggdrasilBuild, error) {
//	build, buildCsv, err := u.YggdrasilBuildBeforeUse()
//	if err != nil {
//		return nil, errors.WrapTrace(err)
//	}
//
//	// 4. 使用，玩家的次元层数变化，可使用次数减一，减钱
//
//	buildAfterUse, err := u.YggdrasilBuildAfterUse(build, buildCsv)
//	if err != nil {
//		return nil, errors.WrapTrace(err)
//	}
//
//	// 6. 返回值
//	return buildAfterUse.VOYggdrasilBuild(), nil
//}

func (u *User) YggdrasilBuildUseGetBuff() error {

	return nil
}

func (u *User) YggdrasilMessageCreate(ctx context.Context, m string) (*YggMessage, error) {
	// 1. 检查当前坐标是否可以留言
	//isBlank := u.Yggdrasil.Entities.IsBlankPos(*u.Yggdrasil.TravelPos)
	//if !isBlank {
	//	return nil, errors.Swrapf(common.ErrYggdrasilPosHasNoEntityToUse, EntityTypeBuilding, u.Yggdrasil.TravelPos.X, u.Yggdrasil.TravelPos.Y)
	//}
	// 留言可以和其它entity重合 但留言和留言不能重合
	if _, ok := u.Yggdrasil.Entities.FindMessageByPos(*u.Yggdrasil.TravelPos); ok {
		return nil, errors.WrapTrace(errors.Swrapf(common.ErrYggdrasilMessageRepeated, u.Yggdrasil.TravelPos.X, u.Yggdrasil.TravelPos.Y))
	}

	areaId, err := manager.CSV.Yggdrasil.GetPosAreaId(*u.Yggdrasil.TravelPos)
	if err != nil {
		return nil, err
	}
	area, err := manager.CSV.Yggdrasil.GetArea(areaId)
	if err != nil {
		return nil, err
	}
	for _, cityId := range area.AreaCityIds {
		city, err := manager.CSV.Yggdrasil.GetYggCityById(cityId)
		if err != nil {
			return nil, err
		}
		dis := coordinate.CubeDistance(*u.Yggdrasil.TravelPos, *city.CityCenterPos)
		if dis <= city.CityBanR {
			return nil, errors.Swrapf(common.ErrYggdrasilPosInsideCityBanRadius, u.Yggdrasil.TravelPos.X, u.Yggdrasil.TravelPos.Y)
		}
	}

	msg, err := NewYggMessage(ctx, m, *u.Yggdrasil.TravelPos)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 2. 添加留言
	err = u.Yggdrasil.AppendOwnMessage(ctx, u.ID, msg)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// 3. 返回值
	return msg, nil
}

// 目前是需要走到对应坐标点上才能修改
func (u *User) YggdrasilMessageUpdate(ctx context.Context, m string) (*pb.VOYggdrasilMessage, error) {

	msg, ok := u.Yggdrasil.Entities.FindMessageByPos(*u.Yggdrasil.TravelPos)
	if !ok {
		return nil, errors.Swrapf(common.ErrYggdrasilPosHasNoEntityToUse, EntityTypeMessage, u.Yggdrasil.TravelPos.X, u.Yggdrasil.TravelPos.Y)
	}
	err := u.Yggdrasil.UpdateMessage(ctx, u.ID, msg, m)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return msg.VOYggdrasilMessage(), nil
}

func (u *User) YggdrasilMessageDestroy(ctx context.Context, uid int64) (*YggMessage, error) {

	msg, ok := u.Yggdrasil.Entities.Messages[uid]
	if !ok {
		return nil, errors.Swrapf(common.ErrYggdrasilWrongUIDForEntity, EntityTypeMessage, uid)
	}

	err := u.Yggdrasil.RemoveMessage(ctx, u.ID, msg)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return msg, nil
}

// ----------------传送门-------------------
// todo 匹配功能开启后才解锁联合建筑?
// todo 再过一下GetInfo

type RPCGuildCoBuildImproveFunc func(context.Context, int32) (*pb.GuildCoBuildImproveResp, error)
type RPCGuildCoBuildGetInfoFunc func(context.Context, int32) (*pb.GuildCoBuildGetInfoResp, error)

func (u *User) YggdrasilTransferPortalGetBuild() (*YggCoBuild, error) {
	coBuild, ok := u.Yggdrasil.Entities.FindCoBuildingsByPos(*u.Yggdrasil.TravelPos)
	if !ok {
		return nil, errors.Swrapf(common.ErrYggdrasilPosHasNoEntityToUse, EntityTypeCoBuilding, u.Yggdrasil.TravelPos.X, u.Yggdrasil.TravelPos.Y)
	}
	return coBuild, nil
}

func (u *User) YggdrasilCoBuildGetInfo() (*pb.VOYggdrasilCoBuild, error) {
	coBuild, err := u.YggdrasilTransferPortalGetBuild()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	u.Yggdrasil.Entities.CoBuilds[coBuild.Uid] = coBuild

	return coBuild.VOYggdrasilCoBuild(u.Yggdrasil, coBuild.Progress, []int64{u.GetUserId()}, coBuild.TotalUseCount), nil
}

func (u *User) YggdrasilTransferPortalBuild(ctx context.Context, rpcGuild RPCGuildCoBuildImproveFunc) (*pb.VOYggdrasilCoBuild, bool, error) {

	coBuild, err := u.YggdrasilTransferPortalGetBuild()
	if err != nil {
		return nil, false, errors.WrapTrace(err)
	}

	if coBuild.CheckIsActive() {
		return nil, false, errors.Swrapf(common.ErrYggdrasilTransferPortalAlreadyActivated, coBuild.BuildId)
	}

	coBuildCsv, err := manager.CSV.Yggdrasil.GetYggCoBuilding(coBuild.BuildId)
	if err != nil {
		return nil, false, errors.WrapTrace(err)
	}

	buildCosts := coBuildCsv.BuildCost
	err = u.CheckRewardsEnough(buildCosts)
	if err != nil {
		return nil, false, errors.WrapTrace(err)
	}

	var progress int32
	var contributorList []int64
	buildNotExecuted := false

	// 个人进度无论如何是否有公会都会+1
	progress = coBuild.PersonalProgressImprove()

	if u.Guild.HasJoinedGuild() { // 已经加入公会
		ret, err := rpcGuild(ctx, coBuild.BuildId)
		if err != nil {
			return nil, false, errors.WrapTrace(err)
		}
		// 这个变量的名称改一下，应该是completed，表示在当前玩家的请求生效前此公会建筑是否已经建造完成
		if ret.BuildNotExcuted {
			buildNotExecuted = true
		}
		progress = ret.CoBuild.Progress
		contributorList = ret.CoBuild.ContributorList
	} else {
		if coBuild.CheckIsCompleted() {
			buildNotExecuted = true
		}
		contributorList = append(contributorList, u.ID)
	}

	reason := logreason.NewReason(logreason.YggBuildTranportPortal)
	err = u.CostRewards(buildCosts, reason)
	if err != nil {
		return nil, false, errors.WrapTrace(err)
	}

	// 5. 返回值, 返回值中可传送地点列表为空
	return coBuild.VOYggdrasilCoBuild(u.Yggdrasil, progress, contributorList, coBuild.TotalUseCount), buildNotExecuted, nil
}

func (u *User) YggdrasilTransferPortalActivate(ctx context.Context, rpcGuild RPCGuildCoBuildGetInfoFunc) (*pb.VOYggdrasilCoBuild, error) {

	coBuild, err := u.YggdrasilTransferPortalGetBuild()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	if coBuild.CheckIsActive() {
		return nil, errors.Swrapf(common.ErrYggdrasilTransferPortalAlreadyActivated, coBuild.BuildId)
	}

	var contributorList []int64

	if u.Guild.HasJoinedGuild() { // 已经加入了公会
		ret, err := rpcGuild(ctx, coBuild.BuildId)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		if ret.CoBuild.Progress < YggdrasilTransferPortalMaxProgress {
			return nil, errors.Swrapf(common.ErrYggdrasilCoBuildProgressNotEnoughForActivation, coBuild.BuildId)
		}
		coBuild.TotalUseCount = int32(ret.CoBuild.TotalUseCount)
		contributorList = ret.CoBuild.ContributorList
	} else {
		if !coBuild.CheckIsCompleted() {
			return nil, errors.Swrapf(common.ErrYggdrasilCoBuildProgressNotEnoughForActivation, coBuild.BuildId)
		}
		contributorList = append(contributorList, u.ID)
	}

	coBuild.SetActive()

	u.Yggdrasil.Entities.CoBuilds[coBuild.Uid] = coBuild

	return coBuild.VOYggdrasilCoBuild(u.Yggdrasil, coBuild.Progress, contributorList, coBuild.TotalUseCount), nil
}

func (u *User) YggdrasilTransferPortalUse(ctx context.Context, target YggPortalLocation) error {

	// 确定传送门的逻辑是站在建筑上还是附近也可

	if !u.Yggdrasil.CheckPortalActive(target) {
		return errors.Swrapf(common.ErrYggdrasilPortalTargetNotInList, target.LocType, target.LocId)
	}

	// 修改坐标
	if target.LocType == PortalTypeCity { // 城市
		// TODO:消耗

		// 回城
		err := u.Yggdrasil.TransferToCity(ctx, u, target.LocId)
		if err != nil {
			return errors.WrapTrace(err)
		}
	} else if target.LocType == PortalTypeCobuild { // 传送门

		coBuild, err := u.YggdrasilTransferPortalGetBuild()
		if err != nil {
			return errors.WrapTrace(err)
		}
		coBuildCsv, err := manager.CSV.Yggdrasil.GetYggCoBuilding(coBuild.BuildId)
		if err != nil {
			return errors.WrapTrace(err)
		}
		// 消耗
		err = u.CheckRewardsEnough(coBuildCsv.UsingCost)
		if err != nil {
			return errors.WrapTrace(err)
		}

		reason := logreason.NewReason(logreason.YggUseTranportPortal)
		err = u.CostRewards(coBuildCsv.UsingCost, reason)
		if err != nil {
			return errors.WrapTrace(err)
		}
		// 坐标改变
		_, _, err = u.Yggdrasil.ChangePos(ctx, u, *coBuild.Position)
		if err != nil {
			return errors.WrapTrace(err)
		}
		// 更新建筑数据
		coBuild.TotalUseCount += 1
		coBuild.UpdateAt = servertime.Now().Unix()
		u.Yggdrasil.Entities.CoBuilds[coBuild.Uid] = coBuild

	} else {
		return errors.Swrapf(common.ErrYggdrasilTransferPortalWrongTypeForPortalLocation, target.LocType)
	}

	return nil
}
