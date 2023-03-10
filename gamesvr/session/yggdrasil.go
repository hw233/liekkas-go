package session

import (
	"context"
	"gamesvr/manager"
	"gamesvr/model"
	"shared/common"
	"shared/csv/base"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/utility/coordinate"
	"shared/utility/errors"
	"shared/utility/global"
	"shared/utility/glog"
	"shared/utility/servertime"
)

// -----------------------------推送--------------------------------------------

func (s *Session) TryPushYgg() {
	for _, push := range s.User.Notifies.YggPush {
		cmd, err := manager.CSV.Protocol.GetCmdByProtoName(push)
		if err != nil {
			glog.Errorf("TryPushYgg  GetCmdByProtoName err:%+v", err)
			continue
		}
		err = s.push(cmd, push)
		if err != nil {
			glog.Errorf("TryPushYgg push err:%+v", err)
			continue
		}
	}
	s.User.Notifies.YggPush = nil

}

// -----------------------------出行，移动相关--------------------------------------------

// YggdrasilGetMain 世界探索-进入玩法主界面
func (s *Session) YggdrasilGetMain(ctx context.Context, req *pb.C2SYggdrasilGetMain) (*pb.S2CYggdrasilGetMain, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	totalIntimacy, err := s.GetTotalIntimacy(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	//todo:测试
	manager.Global.SAdd(ctx, "matchUserIds", s.ID)
	// 重新进界面的处理
	s.Yggdrasil.ReEnter(ctx, s.User)

	return &pb.S2CYggdrasilGetMain{
		User:            s.Yggdrasil.VOYggdrasilUser(s.User),
		Position:        s.Yggdrasil.VOYggdrasilPosition(),
		TravelInfo:      s.Yggdrasil.VOYggdrasilTravelInfo(),
		VoTaskTotalInfo: s.Yggdrasil.VOYggdrasilTaskTotalInfo(),
		BlockAndArea:    s.Yggdrasil.GetBlockInfoByEnter(ctx),
		PackInfo:        s.Yggdrasil.VOYggdrasilPackGoods(),
		TaskPackInfo:    s.Yggdrasil.VOTaskPackInfo(),
		TrackMark:       s.Yggdrasil.VOTrackMarkPosition(),
		MarkTotalCount:  s.Yggdrasil.GetMarkTotalCount(),
		TotalIntimacy:   totalIntimacy,
		EnvTerrain:      s.Yggdrasil.Entities.GetEnvTerrain(),
	}, nil
}
func (s *Session) YggdrasilGetAllArea(ctx context.Context, req *pb.C2SYggdrasilGetAllArea) (*pb.S2CYggdrasilGetAllArea, error) {
	return &pb.S2CYggdrasilGetAllArea{
		Areas: s.Yggdrasil.Areas.VOYggdrasilAreas(),
	}, nil
}

// YggdrasilGetBlockInfo 世界探索-获得块信息
func (s *Session) YggdrasilGetBlockInfo(ctx context.Context, req *pb.C2SYggdrasilGetBlockInfo) (*pb.S2CYggdrasilGetBlockInfo, error) {
	// 参数检查
	if len(req.Positions) == 0 {
		return nil, errors.WrapTrace(common.ErrParamError)
	}
	var posList []coordinate.Position
	for _, position := range req.Positions {
		if position == nil {
			return nil, errors.WrapTrace(common.ErrParamError)

		}
		posList = append(posList, *coordinate.NewPositionFromVo(position))

	}

	return s.User.YggdrasilGetBlockInfo(ctx, posList)
}

// YggdrasilExploreStart 世界探索-开始探索
func (s *Session) YggdrasilExploreStart(ctx context.Context, req *pb.C2SYggdrasilExploreStart) (*pb.S2CYggdrasilExploreStart, error) {
	if len(req.ExploreCharacters) == 0 {
		return nil, errors.WrapTrace(common.ErrParamError)
	}
	err := s.User.YggdrasilExploreStart(req.ExploreCharacters)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// 匹配
	err = s.YggdrasilMatch(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CYggdrasilExploreStart{
		User:         s.Yggdrasil.VOYggdrasilUser(s.User),
		BlockAndArea: s.Yggdrasil.GetBlockInfoByEnter(ctx),
		TravelInfo:   s.Yggdrasil.VOYggdrasilTravelInfo(),
	}, nil
}

func (s *Session) YggdrasilMatch(ctx context.Context) error {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return errors.WrapTrace(err)
	}
	//// 未加入公会，不匹配
	//if s.Guild.GuildID == 0 {
	//	return nil
	//}
	//memberRsp, err := s.RPCGuildMembers(ctx)
	//if err != nil {
	//	return errors.WrapTrace(err)
	//}
	//// 其他成员
	//otherMembers := make([]int64, 0, len(memberRsp.UserIds)-1)
	//for i, id := range memberRsp.UserIds {
	//	if id == s.GetUserId() {
	//		otherMembers = append(memberRsp.UserIds[:i], memberRsp.UserIds[i+1:]...)
	//      break;
	//	}
	//}
	// todo：测试用
	matchUserIds, err := manager.Global.SMembers(ctx, "matchUserIds")
	if err == global.ErrNil {
		return nil
	}
	if err != nil {
		return errors.WrapTrace(err)
	}

	otherMembers := make([]int64, 0, len(matchUserIds)-1)
	for _, str := range matchUserIds {
		id, ok := base.String2Int64(str)
		if !ok {
			continue
		}
		if id == s.GetUserId() {
			continue
		}
		otherMembers = append(otherMembers, id)
	}
	return s.User.YggdrasilMatch(ctx, otherMembers)

}

// YggdrasilExploreMove 世界探索-移动
func (s *Session) YggdrasilExploreMove(ctx context.Context, req *pb.C2SYggdrasilExploreMove) (*pb.S2CYggdrasilExploreMove, error) {

	if req.Position == nil {
		return nil, errors.WrapTrace(common.ErrParamError)
	}

	vos, updates, err := s.User.YggdrasilExploreMove(ctx, *coordinate.NewPositionFromVo(req.Position))
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CYggdrasilExploreMove{
		Position:       s.Yggdrasil.VOYggdrasilPosition(),
		Ap:             s.Yggdrasil.TravelInfo.TravelAp,
		UnlockPosList:  vos,
		PosCountUpdate: updates,
	}, nil
}

// YggdrasilExploreQuit 世界探索-体力耗完休息
func (s *Session) YggdrasilExploreQuit(ctx context.Context, req *pb.C2SYggdrasilExploreQuit) (*pb.S2CYggdrasilExploreQuit, error) {

	characters, err := s.User.YggdrasilExploreQuit(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CYggdrasilExploreQuit{
		User:              s.Yggdrasil.VOYggdrasilUser(s.User),
		ResourceResult:    s.User.VOResourceResult(),
		ExploreCharacters: characters,
	}, nil
}

// YggdrasilExploreReturnCity 世界探索-回城
func (s *Session) YggdrasilExploreReturnCity(ctx context.Context, req *pb.C2SYggdrasilExploreReturnCity) (*pb.S2CYggdrasilExploreReturnCity, error) {
	err := s.User.YggdrasilExploreReturnCity(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CYggdrasilExploreReturnCity{
		ReturnCity:     s.User.Yggdrasil.VOYggdrasilReturnCity(),
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

// YggdrasilExploreLeaveCity 世界探索-从城市出发
func (s *Session) YggdrasilExploreLeaveCity(ctx context.Context, req *pb.C2SYggdrasilExploreLeaveCity) (*pb.S2CYggdrasilExploreLeaveCity, error) {
	err := s.User.YggdrasilExploreLeaveCity(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CYggdrasilExploreLeaveCity{
		User:         s.Yggdrasil.VOYggdrasilUser(s.User),
		Position:     s.Yggdrasil.VOYggdrasilPosition(),
		TravelInfo:   s.Yggdrasil.VOYggdrasilTravelInfo(),
		BlockAndArea: s.Yggdrasil.GetBlockInfoByEnter(ctx),
	}, nil
}

// -----------------------------背包，丢弃相关--------------------------------------------

// YggdrasilGoodsDiscard 世界探索-丢弃物品
func (s *Session) YggdrasilGoodsDiscard(ctx context.Context, req *pb.C2SYggdrasilGoodsDiscard) (*pb.S2CYggdrasilGoodsDiscard, error) {
	err := s.User.YggdrasilGoodsDiscard(ctx, false, req.PackGoodsId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CYggdrasilGoodsDiscard{
		YggdrasilResourceResult: s.User.VOYggdrasilResourceResult(),
	}, nil
}

// YggdrasilGoodsPickUp 世界探索-物品捡起
func (s *Session) YggdrasilGoodsPickUp(ctx context.Context, req *pb.C2SYggdrasilGoodsPickUp) (*pb.S2CYggdrasilGoodsPickUp, error) {
	err := s.User.YggdrasilGoodsPickUp(ctx, req.ReplacedGoodId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CYggdrasilGoodsPickUp{
		YggdrasilResourceResult: s.User.VOYggdrasilResourceResult(),
	}, nil
}

// -----------------------------互动物相关--------------------------------------------

// YggdrasilObjectHandle 操作object (宝箱 魔台等)
func (s *Session) YggdrasilObjectHandle(ctx context.Context, req *pb.C2SYggdrasilObjectHandle) (*pb.S2CYggdrasilObjectHandle, error) {
	err := s.User.YggdrasilObjectHandle(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CYggdrasilObjectHandle{
		EntityChange:   s.VOYggdrasilEntityChange(ctx),
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

// YggdrasilObjectMove 世界探索-object移动
func (s *Session) YggdrasilObjectMove(ctx context.Context, req *pb.C2SYggdrasilObjectMove) (*pb.S2CYggdrasilObjectMove, error) {
	if len(req.ObjectIdList) == 0 {
		return nil, errors.WrapTrace(common.ErrParamError)
	}

	posMap := map[int64]*coordinate.Position{}

	for _, vo := range req.ObjectIdList {
		if vo == nil {
			return nil, errors.WrapTrace(common.ErrParamError)
		}
		if vo.GetPosition() == nil {
			return nil, errors.WrapTrace(common.ErrParamError)
		}
		posMap[vo.ObjectUid] = coordinate.NewPositionFromVo(vo.GetPosition())
	}

	voList, err := s.User.YggdrasilObjectMove(ctx, posMap)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CYggdrasilObjectMove{
		ObjectList: voList,
	}, nil
}

// YggdrasilQueryPosition 世界探索-查询object位置
func (s *Session) YggdrasilQueryPosition(ctx context.Context, req *pb.C2SYggdrasilQueryPosition) (*pb.S2CYggdrasilQueryPosition, error) {
	vos, err := s.User.YggdrasilQueryPosition(ctx, req.ObjectId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CYggdrasilQueryPosition{
		Objects: vos,
	}, nil
}

// YggdrasilAreaProgressReward 世界探索-领取探索进度奖励
func (s *Session) YggdrasilAreaProgressReward(ctx context.Context, req *pb.C2SYggdrasilAreaProgressReward) (*pb.S2CYggdrasilAreaProgressReward, error) {
	vo, err := s.User.YggdrasilAreaProgressReward(ctx, req.AreaId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CYggdrasilAreaProgressReward{
		Area:           vo,
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

// -----------------------------任务相关--------------------------------------------

// YggdrasilAcceptTask 世界探索-领任务
func (s *Session) YggdrasilAcceptTask(ctx context.Context, req *pb.C2SYggdrasilAcceptTask) (*pb.S2CYggdrasilAcceptTask, error) {
	vo, err := s.User.YggdrasilAcceptTask(ctx, req.TaskId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CYggdrasilAcceptTask{
		TaskInfo:     vo,
		EntityChange: s.User.VOYggdrasilEntityChange(ctx),
	}, nil
}

// YggdrasilCompleteTask 世界探索-完成任务
func (s *Session) YggdrasilCompleteTask(ctx context.Context, req *pb.C2SYggdrasilCompleteTask) (*pb.S2CYggdrasilCompleteTask, error) {
	err := s.User.YggdrasilCompleteTask(ctx, req.TaskId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CYggdrasilCompleteTask{
		CompletedTasks: s.Yggdrasil.Task.CompleteTaskIds.Values(),
		EntityChange:   s.User.VOYggdrasilEntityChange(ctx),
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

// YggdrasilAbandonTask 世界探索-放弃任务
func (s *Session) YggdrasilAbandonTask(ctx context.Context, req *pb.C2SYggdrasilAbandonTask) (*pb.S2CYggdrasilAbandonTask, error) {
	err := s.User.YggdrasilAbandonTask(ctx, req.TaskId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CYggdrasilAbandonTask{
		TaskId:       req.TaskId,
		EntityChange: s.User.VOYggdrasilEntityChange(ctx),
	}, nil
}

// YggdrasilDeliverTaskGoods 世界探索-交付任务物品
func (s *Session) YggdrasilDeliverTaskGoods(ctx context.Context, req *pb.C2SYggdrasilDeliverTaskGoods) (*pb.S2CYggdrasilDeliverTaskGoods, error) {
	if req.Resources == nil {
		return nil, errors.WrapTrace(common.ErrParamError)
	}
	err := s.User.YggdrasilDeliverTaskGoods(ctx, req.SubTaskId, req.Resources)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CYggdrasilDeliverTaskGoods{
		ResourceResult:  s.User.VOResourceResult(),
		YggdrasilResult: s.User.VOYggdrasilResourceResult(),
	}, nil
}

// YggdrasilSetTrackTask 世界探索-设置追踪任务
func (s *Session) YggdrasilSetTrackTask(ctx context.Context, req *pb.C2SYggdrasilSetTrackTask) (*pb.S2CYggdrasilSetTrackTask, error) {
	err := s.User.YggdrasilSetTrackTask(req.TrackTask)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CYggdrasilSetTrackTask{
		TrackTask: req.TrackTask,
	}, nil
}

// YggdrasilChooseNext 世界探索-选择下一个子任务
func (s *Session) YggdrasilChooseNext(ctx context.Context, req *pb.C2SYggdrasilChooseNext) (*pb.S2CYggdrasilChooseNext, error) {
	vo, err := s.User.YggdrasilChooseNext(ctx, req.SubTaskId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CYggdrasilChooseNext{
		TaskInfo:       vo,
		EntityChange:   s.User.VOYggdrasilEntityChange(ctx),
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

// -----------------------------建筑和标注相关--------------------------------------------

// YggdrasilBuildCreate 世界探索-建筑建造
func (s *Session) YggdrasilBuildCreate(ctx context.Context, req *pb.C2SYggdrasilBuildCreate) (*pb.S2CYggdrasilBuildCreate, error) {

	build, err := s.User.YggdrasilBuildCreate(ctx, req.BuildId)
	if err != nil {
		return nil, err
	}

	vo, err := build.VOYggdrasilBuild(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.S2CYggdrasilBuildCreate{
		NewBuild:       vo,
		ResourceResult: s.User.VOResourceResult(),
		Area:           s.User.Yggdrasil.Areas.VOYggdrasilArea(build.AreaId),
	}, nil
}

// YggdrasilBuildDestroy 世界探索-建筑拆除
func (s *Session) YggdrasilBuildDestroy(ctx context.Context, req *pb.C2SYggdrasilBuildDestroy) (*pb.S2CYggdrasilBuildDestroy, error) {
	build, err := s.User.YggdrasilBuildDestroy(ctx, req.BuildUid)
	if err != nil {
		return nil, err
	}

	// todo 是否还需要返回uid呢,需要跟客户端商量
	return &pb.S2CYggdrasilBuildDestroy{
		DestroyBuildUid: req.BuildUid,
		Area:            s.User.Yggdrasil.Areas.VOYggdrasilArea(build.AreaId),
	}, nil
}

// YggdrasilBuildAddAp 世界探索-建筑使用精力泉水
func (s *Session) YggdrasilBuildAddAp(ctx context.Context, req *pb.C2SYggdrasilBuildAddAp) (*pb.S2CYggdrasilBuildAddAp, error) {
	build, ap, err := s.User.YggdrasilBuildAddAp(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.S2CYggdrasilBuildAddAp{
		VoUseBuild: &pb.VOYggdrasilUseBuild{
			Build:          build,
			ResourceResult: s.User.VOResourceResult(),
		},
		Ap: ap,
	}, nil
}

// YggdrasilBuildAddHp 世界探索-建筑使用血量泉水
func (s *Session) YggdrasilBuildAddHp(ctx context.Context, req *pb.C2SYggdrasilBuildAddHp) (*pb.S2CYggdrasilBuildAddHp, error) {
	build, characs, err := s.User.YggdrasilBuildAddHp(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.S2CYggdrasilBuildAddHp{
		VoUseBuild: &pb.VOYggdrasilUseBuild{
			Build:          build,
			ResourceResult: s.User.VOResourceResult(),
		},
		ExploreCharacters: characs,
	}, nil
}

// YggdrasilBuildUsePeeing 世界探索-建筑使用魔法窥视站
func (s *Session) YggdrasilBuildUsePeeing(ctx context.Context, req *pb.C2SYggdrasilBuildUsePeeing) (*pb.S2CYggdrasilBuildUsePeeing, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	build, err := s.User.YggdrasilBuildUsePeeing(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.S2CYggdrasilBuildUsePeeing{
		VoUseBuild: &pb.VOYggdrasilUseBuild{
			Build:          build,
			ResourceResult: s.User.VOResourceResult(),
		},
		EntityChange: s.VOYggdrasilEntityChange(ctx),
	}, nil
}

// YggdrasilBuildUseTransPort 世界探索-建筑使用魔法传送阵
func (s *Session) YggdrasilBuildUseTransPort(ctx context.Context, req *pb.C2SYggdrasilBuildUseTransPort) (*pb.S2CYggdrasilBuildUseTransPort, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	build, err := s.User.YggdrasilBuildUseTransPort(ctx, req.GoodsIdList)
	if err != nil {
		return nil, err
	}

	return &pb.S2CYggdrasilBuildUseTransPort{
		VoUseBuild:        &pb.VOYggdrasilUseBuild{Build: build, ResourceResult: s.User.VOResourceResult()},
		YggResourceResult: s.VOYggdrasilResourceResult(),
	}, nil
}

// YggdrasilBuildUseGetBuff 世界探索-建筑使用魔法灵息
func (s *Session) YggdrasilBuildUseGetBuff(ctx context.Context, req *pb.C2SYggdrasilBuildUseGetBuff) (*pb.S2CYggdrasilBuildUseGetBuff, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CYggdrasilBuildUseGetBuff{}, nil
}

// // YggdrasilBuildUseGetBuff 世界探索-建筑使用次元电梯
// func (s *Session) YggdrasilBuildUseStepladder(ctx context.Context, req *pb.C2SYggdrasilBuildUseStepladder) (*pb.S2CYggdrasilBuildStepladder, error) {
// 	return &pb.S2CYggdrasilBuildUseStepladder{}, nil
// }

// YggdrasilMarkCreate 世界探索-标注创建
func (s *Session) YggdrasilMarkCreate(ctx context.Context, req *pb.C2SYggdrasilMarkCreate) (*pb.S2CYggdrasilMarkCreate, error) {
	if req.Pos == nil {
		return nil, errors.WrapTrace(common.ErrParamError)
	}

	vo, err := s.User.YggdrasilMarkCreate(ctx, req.MarkId, *coordinate.NewPositionFromVo(req.Pos))
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CYggdrasilMarkCreate{
		NewMark:        vo,
		MarkTotalCount: s.Yggdrasil.GetMarkTotalCount(),
	}, nil
}

// YggdrasilMarkDestroy 世界探索-标注销毁
func (s *Session) YggdrasilMarkDestroy(ctx context.Context, req *pb.C2SYggdrasilMarkDestroy) (*pb.S2CYggdrasilMarkDestroy, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	err = s.User.YggdrasilMarkDestroy(req.MarkUid)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CYggdrasilMarkDestroy{
		DeleteMarkUid:  req.MarkUid,
		MarkTotalCount: s.Yggdrasil.GetMarkTotalCount(),
	}, nil
}

// YggdrasilTrackMark 世界探索-设置追踪标注
func (s *Session) YggdrasilTrackMark(ctx context.Context, req *pb.C2SYggdrasilTrackMark) (*pb.S2CYggdrasilTrackMark, error) {
	var pos *coordinate.Position
	if req.TrackMark == nil {
		pos = nil
	} else {
		pos = coordinate.NewPositionFromVo(req.TrackMark)
	}
	s.User.YggdrasilTrackMark(pos)
	return &pb.S2CYggdrasilTrackMark{
		TrackMark: req.TrackMark,
	}, nil
}

func (s *Session) YggdrasilMessageCreate(ctx context.Context, req *pb.C2SYggdrasilMessageCreate) (*pb.S2CYggdrasilMessageCreate, error) {
	msg, err := s.User.YggdrasilMessageCreate(ctx, req.Comment)
	if err != nil {
		return nil, err
	}

	return &pb.S2CYggdrasilMessageCreate{
		Msg:  msg.VOYggdrasilMessage(),
		Area: s.User.Yggdrasil.Areas.VOYggdrasilArea(msg.AreaId),
	}, nil
}

func (s *Session) YggdrasilMessageUpdate(ctx context.Context, req *pb.C2SYggdrasilMessageUpdate) (*pb.S2CYggdrasilMessageUpdate, error) {
	msg, err := s.User.YggdrasilMessageUpdate(ctx, req.Comment)
	if err != nil {
		return nil, err
	}

	return &pb.S2CYggdrasilMessageUpdate{
		Msg: msg,
	}, nil
}

func (s *Session) YggdrasilMessageDestroy(ctx context.Context, req *pb.C2SYggdrasilMessageDestroy) (*pb.S2CYggdrasilMessageDestroy, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	msg, err := s.User.YggdrasilMessageDestroy(ctx, req.MessageUid)
	if err != nil {
		return nil, err
	}

	return &pb.S2CYggdrasilMessageDestroy{
		MessageUid: msg.Uid,
		Area:       s.User.Yggdrasil.Areas.VOYggdrasilArea(msg.AreaId),
	}, nil
}

// --------------------联合建筑----------------------
func (s *Session) YggdrasilCoBuildGetInfo(ctx context.Context, req *pb.C2SYggdrasilCoBuildGetInfo) (*pb.S2CYggdrasilCoBuildGetInfo, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	voBuild, err := s.User.YggdrasilCoBuildGetInfo()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 如果加入了公会，就用公会中对应建筑的数据
	if s.Guild.HasJoinedGuild() { // 已经加入了公会
		ret, err := s.RPCGuildCoBuildGetInfo(ctx, voBuild.CoBuildBase.BuildId)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		if !voBuild.CoBuildBase.IsActivated { // 未激活，就使用公会的进度；如果激活了，就永久激活
			voBuild.CoBuildBase.Progress = ret.CoBuild.Progress
		}

		voBuild.CoBuildBase.OriContributorIdList = ret.CoBuild.ContributorList
		voBuild.CoBuildBase.TotalUseCount = int32(ret.CoBuild.TotalUseCount)
	}

	return &pb.S2CYggdrasilCoBuildGetInfo{
		Build: voBuild,
	}, nil
}
func (s *Session) YggdrasilTransferPortalActivate(ctx context.Context, req *pb.C2SYggTransferPortalActivate) (*pb.S2CYggTransferPortalActivate, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	voBuild, err := s.User.YggdrasilTransferPortalActivate(ctx, s.RPCGuildCoBuildGetInfo)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CYggTransferPortalActivate{
		VoCoBuild: voBuild,
	}, nil
}

func (s *Session) YggdrasilTransferPortalUse(ctx context.Context, req *pb.C2SYggTransferPortalUse) (*pb.S2CYggTransferPortalUse, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	err = s.User.YggdrasilTransferPortalUse(ctx, *model.NewYggPortalLocation(req.Location.LocationType, req.Location.LocationId))
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CYggTransferPortalUse{
		ReturnCity:   s.User.Yggdrasil.VOYggdrasilReturnCity(),
		Pos:          s.User.Yggdrasil.VOYggdrasilPosition(),
		Result:       s.VOResourceResult(),
		BlockAndArea: s.Yggdrasil.GetBlockInfoByEnter(ctx),
	}, nil
}

func (s *Session) YggdrasilTransferPortalBuild(ctx context.Context, req *pb.C2SYggTransferPortalBuild) (*pb.S2CYggTransferPortalBuild, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	voBuild, buildNotExecuted, err := s.User.YggdrasilTransferPortalBuild(ctx, s.RPCGuildCoBuildImprove)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CYggTransferPortalBuild{
		VoUseCoBuild:     &pb.VOYggdrasilUseCoBuild{Build: voBuild, ResourceResult: s.VOResourceResult()},
		BuildNotExecuted: buildNotExecuted,
	}, nil
}

// -----------------------------异界邮箱--------------------------------------------
func (s *Session) YggdrasilMailGetByPage(ctx context.Context, req *pb.C2SYggdrasilMailGetByPage) (*pb.S2CYggdrasilMailGetByPage, error) {
	if req.Num <= 0 {
		return nil, errors.WrapTrace(common.ErrParamError)
	}
	return &pb.S2CYggdrasilMailGetByPage{
		Mails: s.User.YggdrasilMailGetByPage(req.Offset, req.Num),
	}, nil
}

func (s *Session) YggdrasilMailReceiveOne(ctx context.Context, req *pb.C2SYggdrasilMailReceiveOne) (*pb.S2CYggdrasilMailReceiveOne, error) {
	err := s.User.YggdrasilMailReceiveOne(ctx, req.MailUid)
	if err != nil {
		return nil, err
	}

	return &pb.S2CYggdrasilMailReceiveOne{
		DeletedMailUid: req.MailUid,
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

func (s *Session) YggdrasilMailReceiveAll(ctx context.Context, req *pb.C2SYggdrasilMailReceiveAll) (*pb.S2CYggdrasilMailReceiveAll, error) {
	deletedUids, err := s.User.YggdrasilMailReceiveAll(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.S2CYggdrasilMailReceiveAll{
		DeletedMailUids: deletedUids,
		ResourceResult:  s.User.VOResourceResult(),
	}, nil
}

// ----------派遣------------
func (s *Session) YggDispatchGetInfo(ctx context.Context, req *pb.C2SYggDispatchGetInfo) (*pb.S2CYggDispatchGetInfo, error) {

	dailyTasks, err := s.User.YggdrasilDispatchGetInfo()
	if err != nil {
		return nil, err
	}

	voGuilds, err := s.User.YggDispatchGuildGetInfo()
	if err != nil {
		return nil, err
	}

	dailyInfo := pb.VOYggDispatchInfo{
		DispatchType: static.YggdrasilDispatchTypeDaily,
		UpdateTime:   -1,
		TaskStates:   dailyTasks,
	}

	now := servertime.Now()
	nextDailyRefreshTime := model.DailyRefreshTime(now.AddDate(0, 0, 1))
	guildInfo := pb.VOYggDispatchInfo{
		DispatchType: static.YggdrasilDispatchTypeGuild,
		UpdateTime:   int64(nextDailyRefreshTime.Sub(now).Seconds()),
		TaskStates:   voGuilds,
	}

	return &pb.S2CYggDispatchGetInfo{
		Dispatches: []*pb.VOYggDispatchInfo{&dailyInfo, &guildInfo},
	}, nil
}

func (s *Session) YggDispatchGetDailyTaskInfo(ctx context.Context, req *pb.C2SYggDispatchGetDailyTaskInfo) (*pb.S2CYggDispatchGetDailyTaskInfo, error) {

	charaIds, err := s.User.YggdrasilDispatchDailyTaskInfo(req.TaskId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CYggDispatchGetDailyTaskInfo{
		TaskId:             req.TaskId,
		CharacIdInDispatch: charaIds,
	}, nil
}

func (s *Session) YggDispatchGetGuildTaskInfo(ctx context.Context, req *pb.C2SYggDispatchGetGuildTaskInfo) (*pb.S2CYggDispatchGetGuildTaskInfo, error) {
	// if !s.Guild.HasJoinedGuild(ctx, s.ID) { // 如果没加入公会
	// 	return nil, errors.Swrapf(common.ErrGuildHasNotJoin)
	// }

	// characIds, _, err := s.User.YggdrasilDispatchGuildTaskInfo(req.TaskId)
	// if err != nil {
	// 	return nil, errors.WrapTrace(err)
	// }
	// // s.User.Guild.Join(1, "Guild0010")

	// // todo 修改协议
	// // guildCharac, err := s.RPCGuildGetDispatchCharac(ctx, guildCharacId)
	// // if err != nil {
	// // 	return nil, errors.WrapTrace(err)
	// // }

	// // fmt.Println("---------", guildCharac)

	// return &pb.S2CYggDispatchGetGuildTaskInfo{
	// 	TaskId:             req.TaskId,
	// 	CharacIdInDispatch: characIds,
	// 	GuildCharacs:       []*pb.VOYggDispatchGuildCharacter{},
	// }, nil
	return nil, nil
}

func (s *Session) YggDailyDispatchBegin(ctx context.Context, req *pb.C2SYggDailyDispatchBegin) (*pb.S2CYggDailyDispatchBegin, error) {
	taskState, err := s.User.YggdrasilDispatchDailyTaskBegin(req.TaskId, req.UserCharacIds)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CYggDailyDispatchBegin{
		TaskInfo: taskState,
	}, nil
}

// 协议需要修改，客户端需要把选中的公会角色的全部信息发过来
func (s *Session) YggGuildDispatchBegin(ctx context.Context, req *pb.C2SYggGuildDispatchBegin) (*pb.S2CYggGuildDispatchBegin, error) {
	// voTaskState, err := s.User.YggdrasilDispatchGuildTaskBegin(req.TaskId, req.UserCharacIds, req.GuildCharacter)
	// if err != nil {
	// 	return nil, errors.WrapTrace(err)
	// }

	// return &pb.S2CYggGuildDispatchBegin{
	// 	TaskInfo: voTaskState,
	// }, nil
	return nil, nil
}

func (s *Session) YggDispatchReceiveRewards(ctx context.Context, req *pb.C2SYggDispatchReceiveRewards) (*pb.S2CYggDispatchReceiveRewards, error) {
	resourceResult, err := s.User.YggdrasilDispatchReward(req.TaskId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CYggDispatchReceiveRewards{
		TaskId:  req.TaskId,
		Rewards: resourceResult,
	}, nil
}
func (s *Session) YggDispatchCancel(ctx context.Context, req *pb.C2SYggDispatchCancel) (*pb.S2CYggDispatchCancel, error) {
	voTaskState, err := s.User.YggDispatchCancel(ctx, req.TaskId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CYggDispatchCancel{
		TaskInfo: voTaskState,
	}, nil
}

// YggSpecialStatics 世界探索-ygg特殊的统计（用于新手引导）
func (s *Session) YggSpecialStatics(ctx context.Context, req *pb.C2SYggSpecialStatics) (*pb.S2CYggSpecialStatics, error) {
	return &pb.S2CYggSpecialStatics{
		Statics: s.Yggdrasil.VOYggSpecialStatics(),
	}, nil
}

func (s *Session) UserQuerySimpleInfo(ctx context.Context, req *pb.C2SUserQuerySimpleInfo) (*pb.S2CUserQuerySimpleInfo, error) {
	vos, err := s.User.QuerySimpleInfo(ctx, req.UserIdList)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CUserQuerySimpleInfo{
		SimpleInfoList: vos,
	}, nil
}

func (s *Session) YggBattleGiveUp(ctx context.Context, req *pb.C2SYggBattleGiveUp) (*pb.S2CYggBattleGiveUp, error) {
	err := s.User.YggBattleGiveUp(ctx, req.ObjectUid)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CYggBattleGiveUp{
		Position:   s.Yggdrasil.VOYggdrasilPosition(),
		TravelInfo: s.Yggdrasil.VOYggdrasilTravelInfo(),
	}, nil
}

func (s *Session) YggMonsterInitPos(ctx context.Context, req *pb.C2SYggMonsterInitPos) (*pb.S2CYggMonsterInitPos, error) {
	vo, err := s.User.YggMonsterInitPos(req.ObjectUid)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CYggMonsterInitPos{
		Position:  vo,
		ObjectUid: req.ObjectUid,
	}, nil
}

func (s *Session) YggMonsterBackInitPos(ctx context.Context, req *pb.C2SYggMonsterBackInitPos) (*pb.S2CYggMonsterBackInitPos, error) {
	vo, err := s.User.YggMonsterBackInitPos(ctx, req.ObjectUid)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CYggMonsterBackInitPos{
		Monster: vo,
	}, nil
}
