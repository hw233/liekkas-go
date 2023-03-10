package model

import (
	"shared/protobuf/pb"
)

const YggdrasilTransferPortalMaxProgress = 20

// todo 如果某人离开公会，需要遍历公会中的建筑，将其从列表中清楚
type YggCoBuild struct {
	Id              int32
	Progress        int32
	TotalUseCount   int64
	ContributorList []int64
}

func NewYggCoBuild(id int32) *YggCoBuild {
	// 建筑进度每次增加百分之五，最多会有20个贡献者
	contributors := make([]int64, 0, YggdrasilTransferPortalMaxProgress)

	return &YggCoBuild{
		Id:              id,
		Progress:        0,
		TotalUseCount:   0,
		ContributorList: contributors,
	}
}

// todo 增加亲密度(默契值)
func (y *YggCoBuild) AddUseCnt(uid int64) {
	y.TotalUseCount += 1
}

func (y *YggCoBuild) AddProgress(uid int64) bool {
	if y.Progress >= YggdrasilTransferPortalMaxProgress {
		return true
	}
	y.Progress += 1
	if !y.FindContributor(uid) {
		y.ContributorList = append(y.ContributorList, uid)
	}

	return false
}

func (y *YggCoBuild) FindContributor(uid int64) bool {
	for _, p := range y.ContributorList {
		if uid == p {
			return true
		}
	}

	return false
}

func (y *YggCoBuild) VOGuildYggdrailCoBuild() *pb.VOGuildYggdrasilCoBuild {
	return &pb.VOGuildYggdrasilCoBuild{
		BuildID:         y.Id,
		Progress:        y.Progress,
		TotalUseCount:   y.TotalUseCount,
		ContributorList: y.ContributorList,
	}
}
