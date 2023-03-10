package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/csv/entry"
	"shared/protobuf/pb"
	"shared/utility/coordinate"
	"shared/utility/errors"
)

type Graveyard struct {
	*DailyRefreshChecker
	Incr                     int64                                `json:"incr"`                         // 建筑uid递增
	Buildings                map[int64]*common.UserGraveyardBuild `json:"buildings"`                    // 所有建筑（包括在背包中）
	MainTowerUid             int64                                `json:"-"`                            // 主堡uid
	SendRequestCount         int32                                `json:"send_request_num"`             // 公会互助今日发送互助请求次数
	PlotRecord               *GraveyardPlotRecord                 `json:"plot_record"`                  // 小人对话领奖次数 ，每天随机3个整点时间，每到一个整点后可领一次奖
	DailyGuildGoldByHelp     int32                                `json:"daily_guild_gold_by_help"`     // 互助玩法今日获取代币数
	DailyAddActivationByHelp int32                                `json:"daily_add_activation_by_help"` // 互助玩法今日增加贡献值

}

func NewGraveyard() *Graveyard {
	return &Graveyard{
		DailyRefreshChecker: NewDailyRefreshChecker(),
		Incr:                0,
		Buildings:           map[int64]*common.UserGraveyardBuild{},
		SendRequestCount:    0,
		PlotRecord:          NewGraveyardPlotRecord(),
	}
}

type GraveyardPlotRecord struct {
	RewardHours []int32 `json:"reward_hours"` // 领奖随机区间
	RewardNum   int32   `json:"reward_num"`   // 今日已领奖次数
}

func NewGraveyardPlotRecord() *GraveyardPlotRecord {
	return &GraveyardPlotRecord{
		RewardHours: manager.CSV.GraveyardEntry.RandomRewardHours(),
		RewardNum:   0,
	}
}

// GetRewardHours 返回RewardHours slice 去掉已领奖的
func (g *GraveyardPlotRecord) GetRewardHours() []int32 {
	return g.RewardHours[g.RewardNum:]
}

func (g *GraveyardPlotRecord) GetRewardNum() int32 {
	return g.RewardNum

}

func (g *GraveyardPlotRecord) SetRewardNum(newVal int32) {
	g.RewardNum = newVal
}

func (g *GraveyardPlotRecord) IncrRewardNum(addVal int32) int32 {
	incr := g.GetRewardNum() + addVal
	g.SetRewardNum(incr)
	return incr
}

func (g *Graveyard) GetMainTower() (*common.UserGraveyardBuild, bool) {
	if g.MainTowerUid == 0 {
		for uid, build := range g.Buildings {
			if build.BuildId == entry.MainTowerBuildId {
				g.MainTowerUid = uid
				break
			}
		}
	}
	mainTower, ok := g.Buildings[g.MainTowerUid]
	return mainTower, ok
}

func (g *Graveyard) GetMainTowerLevel() int32 {

	mainTower, ok := g.GetMainTower()
	if ok {
		return mainTower.FetchLevel()
	}
	return 0
}

func (g *Graveyard) CreateBuild(build *common.UserGraveyardBuild) int64 {
	g.Incr++

	g.Buildings[g.Incr] = build
	return g.Incr
}

func (g *Graveyard) GetBuildsByBuildId(buildIds ...int32) map[int64]*common.UserGraveyardBuild {
	builds := map[int64]*common.UserGraveyardBuild{}
	for _, buildId := range buildIds {
		for uid, build := range g.Buildings {
			//

			if build.BuildId == buildId {
				builds[uid] = build
			}
		}
	}

	return builds
}

func (g *Graveyard) GetBuildsCountByBuildId(buildId int32) int32 {
	return int32(len(g.GetBuildsByBuildId(buildId)))
}

func (g *Graveyard) GetLocatedBuildsList() map[int32][]*coordinate.Position {

	return g.GetLocatedBuildsListExcept(nil)
}

type buildUIdList []int64

func (b *buildUIdList) contains(buildUid int64) bool {
	for _, uid := range *b {
		if uid == buildUid {
			return true
		}
	}
	return false
}

// GetLocatedBuildsListExcept 获得有坐标的建筑列表
func (g *Graveyard) GetLocatedBuildsListExcept(buildUIdList buildUIdList) map[int32][]*coordinate.Position {
	buildPositionListMap := map[int32][]*coordinate.Position{}
	for uid, build := range g.Buildings {
		//
		if !build.InBag() && !buildUIdList.contains(uid) {
			positionList, ok := buildPositionListMap[build.BuildId]
			if !ok {
				positionList = []*coordinate.Position{}
			}

			positionList = append(positionList, build.Position)
			buildPositionListMap[build.BuildId] = positionList

		}
	}

	return buildPositionListMap
}

// VOPartial 增量数据
func (g *Graveyard) VOPartial(buildings map[int64]*common.UserGraveyardBuild, u *User) []*pb.VOGraveyardBuild {
	var builds []*pb.VOGraveyardBuild
	for uid, build := range buildings {
		builds = append(builds, u.VOGraveyardBuild(uid, build))
	}
	return builds
}

// VOAll 全部数据
func (g *Graveyard) VOAll(u *User) []*pb.VOGraveyardBuild {
	return g.VOPartial(g.Buildings, u)
}

func (g *Graveyard) initForCreate() {
	g.CreateBuild(common.NewUserGraveyardMainTower(entry.MainTowerBuildId, manager.CSV.GraveyardEntry.MainTowerInitPosition()))
}

func (g *Graveyard) FindByUid(buildUid int64) (*common.UserGraveyardBuild, error) {
	build, ok := g.Buildings[buildUid]
	if !ok {
		return nil, common.ErrGraveyardBuildNotExist
	}
	return build, nil
}

// CheckLvUpLimitAndConsume 升级条件判断
func (g *Graveyard) CheckLvUpLimitAndConsume(lvUp *entry.BuildingLvUp, u *User) error {
	if lvUp.MainTowerLevelLimit > g.GetMainTowerLevel() {
		return errors.WrapTrace(common.ErrGraveyardMainTowerLvLimit)
	}

	err := u.CheckUserConditions(lvUp.UnlockConditions)
	if err != nil {
		return errors.WrapTrace(err)
	}

	err = g.checkBuildingLevelLimits(lvUp.BuildingLevelLimits)
	if err != nil {
		return errors.WrapTrace(err)
	}
	return u.CheckRewardsEnough(lvUp.LevelUpConsume)

}

// 判断指定buildId,等级大于指定lv的build数量
func (g *Graveyard) checkBuildingLevelLimits(buildingLevelLimits *entry.BuildingLevelLimits) error {
	if buildingLevelLimits == nil {
		return nil
	}
	for _, limit := range buildingLevelLimits.Limits {
		if g.getBuildCount(limit.BuildId, limit.Lv) < limit.NeedNum {
			return errors.WrapTrace(common.ErrGraveyardBuildCountLimit)

		}
	}
	return nil

}

// 获得指定buildId,等级大于指定lv的build数量
func (g *Graveyard) getBuildCount(buildId, lv int32) int32 {
	var i int32
	for _, build := range g.Buildings {
		if build.BuildId == buildId && build.FetchLevel() >= lv {
			i++
		}
	}
	return i
}

// CheckStageUpLimitAndConsume 升阶条件判断
func (g *Graveyard) CheckStageUpLimitAndConsume(build *common.UserGraveyardBuild, stageUp *entry.BuildingStageUp, u *User) error {
	if stageUp.MainTowerLevelLimit > g.GetMainTowerLevel() {
		return errors.WrapTrace(common.ErrGraveyardMainTowerLvLimit)
	}
	err := g.checkBuildingLevelLimits(stageUp.BuildingLevelLimits)
	if err != nil {
		return errors.WrapTrace(err)
	}
	if build.TotalProduceNum < stageUp.ProduceCountLimit {
		return errors.WrapTrace(common.ErrGraveyardProduceTimeLimit)
	}

	return u.CheckRewardsEnough(stageUp.StageUpConsume)

}

func (g *Graveyard) CheckBuildLv(buildId, lv int32) error {
	builds := g.GetBuildsByBuildId(buildId)
	if len(builds) == 0 {
		return errors.WrapTrace(common.ErrParamError)
	}

	for _, build := range builds {
		if build.FetchLevel() >= lv {
			return nil
		}
	}
	return errors.WrapTrace(common.ErrParamError)

}
