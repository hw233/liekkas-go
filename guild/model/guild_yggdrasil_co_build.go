package model

import (
	"shared/common"
	"shared/utility/errors"
)

// 返回值为false，代表执行成功，进度+1；返回值为true，代表进度已达上限，执行失败
func (g *Guild) YggCoBuildImprove(buildID int32, uid int64) (*YggCoBuild, bool, error) {

	coBuild, ok := g.YggCoBuilds[buildID]
	if !ok {
		return nil, false, errors.Swrapf(common.ErrGuildYggWrongIDForCoBuild, buildID)
	}

	return coBuild, coBuild.AddProgress(uid), nil
}

func (g *Guild) YggCoBuildUse(buildID int32, uid int64) (*YggCoBuild, error) {
	coBuild, ok := g.YggCoBuilds[buildID]
	if !ok {
		return nil, errors.Swrapf(common.ErrGuildYggWrongIDForCoBuild, buildID)
	}

	coBuild.AddUseCnt(uid)

	return coBuild, nil
}

func (g *Guild) YggCoBuildGetInfo(buildID int32) *YggCoBuild {
	coBuild, ok := g.YggCoBuilds[buildID]
	if !ok {
		coBuild = NewYggCoBuild(buildID)
	}
	return coBuild
}
