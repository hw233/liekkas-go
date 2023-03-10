package model

import (
	"context"
	"shared/protobuf/pb"
	"shared/utility/coordinate"
	"shared/utility/errors"
)

const (
	YggdrasilCoBuildStateUnComplete = 0
	YggdrasilCoBuildStateCompleted  = 1
	YggdrasilCoBuildStateActivated  = 2
)

const YggdrasilTransferPortalMaxProgress = 20

type YggCoBuild struct {
	*YggEntity    `json:"entity"`
	BuildId       int32 `json:"build_id"`
	Progress      int32 `json:"personal_progress"`
	State         int32 `json:"build_state"`
	TotalUseCount int32 `json:"total_use_count"` // 使用int32还是int64
}

type YggPortalLocation struct {
	LocType int32 `json:"location_type"`
	LocId   int32 `json:"location_id"`
}

func NewYggCoBuild(ctx context.Context, pos coordinate.Position) (*YggCoBuild, error) {

	entity, err := NewYggEntity(ctx, pos)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &YggCoBuild{
		YggEntity:     entity,
		BuildId:       0,
		Progress:      0,
		State:         YggdrasilCoBuildStateUnComplete,
		TotalUseCount: 0,
	}, nil
}

func NewYggPortalLocation(locType, locId int32) *YggPortalLocation {
	return &YggPortalLocation{
		LocType: locType,
		LocId:   locId,
	}
}

func (ypl *YggPortalLocation) VOYggdrasilPortalLocation() *pb.VOYggdrasilPortalLocation {
	return &pb.VOYggdrasilPortalLocation{
		LocationType: ypl.LocType,
		LocationId:   ypl.LocId,
	}
}

// 返回值为true，代表已经激活；返回值为false，代表未激活
func (ycb *YggCoBuild) CheckIsActive() bool {
	return ycb.State >= YggdrasilCoBuildStateActivated
}

func (ycb *YggCoBuild) SetActive() {
	// 如果公会的完成了，个人的还没完成，那么被激活后由于会永久保留建筑，所以个人的需要同步进度至20
	ycb.Progress = YggdrasilTransferPortalMaxProgress
	ycb.State = YggdrasilCoBuildStateActivated
}

func (ycb *YggCoBuild) GetPersonalProgress() int32 {
	return ycb.Progress
}

func (ycb *YggCoBuild) CheckIsCompleted() bool {
	return ycb.State == YggdrasilCoBuildStateCompleted
}

// 个人进度+1，并且返回个人进度的值
func (ycb *YggCoBuild) PersonalProgressImprove() int32 {
	ycb.Progress = (ycb.Progress + 1) % (YggdrasilTransferPortalMaxProgress + 1)
	if ycb.Progress >= YggdrasilTransferPortalMaxProgress {
		ycb.State = YggdrasilCoBuildStateCompleted
	}
	return ycb.GetPersonalProgress()
}

func (ycb *YggCoBuild) VOYggdrasilCoBuildBase(progress int32, contributorList []int64, cnt int32) *pb.VOYggdrasilCoBuildBase {
	return &pb.VOYggdrasilCoBuildBase{
		Position:             ycb.VOPosition(),
		IsActivated:          ycb.State >= 2,
		BuildId:              ycb.BuildId,
		Progress:             progress,
		TotalUseCount:        cnt,
		OriContributorIdList: contributorList,
	}
}

func (ycb *YggCoBuild) VOYggdrasilCoBuild(yggdrasil *Yggdrasil, progress int32, contributorList []int64, cnt int32) *pb.VOYggdrasilCoBuild {
	return &pb.VOYggdrasilCoBuild{
		CoBuildBase:     ycb.VOYggdrasilCoBuildBase(progress, contributorList, cnt),
		PortalLocations: yggdrasil.VOPortalLocations(),
	}
}
