package entry

import (
	"shared/common"
	"shared/utility/coordinate"
	"shared/utility/errors"
	"shared/utility/servertime"
	"testing"
	"time"
)

func TestGraveyardEntry(t *testing.T) {
	area, err := CSV.GraveyardEntry.CalUnlockArea(1, nil)
	if err != nil {
		t.Errorf("%+v", errors.Format(err))

	}

	t.Logf("unlockarea :%+v", area)
}

func TestUnlockArea_Minus(t *testing.T) {
	m := map[int32][]*coordinate.Position{}
	m[11] = []*coordinate.Position{{-6, -12}}
	m[3] = []*coordinate.Position{{-8, -21}}

	area, err := CSV.GraveyardEntry.CalUnlockArea(2, m)
	if err != nil {
		t.Errorf("%+v", errors.Format(err))
		return
	}
	//building, _ := CSV.GraveyardEntry.GetBuildById(11)
	//err = area.MinusBuildingArea(building.BuildingArea, &common.Position{
	//	X: -7,
	//	Y: -12,
	//})
	//if err != nil {
	//	t.Errorf("%+v", errors.Format(err))
	//	return
	//}
	//// ----
	//building, _ = CSV.GraveyardEntry.GetBuildById(3)
	//err = area.MinusBuildingArea(building.BuildingArea, &common.Position{
	//	X: 2,
	//	Y: -9,
	//})
	//if err != nil {
	//	t.Errorf("%+v", errors.Format(err))
	//	return
	//}

	// ----
	building, _ := CSV.GraveyardEntry.GetBuildById(2)
	err = area.MinusBuildingArea(building.BuildingArea, &coordinate.Position{
		X: 5,
		Y: -7,
	})
	if err != nil {
		t.Errorf("%+v", errors.Format(err))
		return
	}
}

func TestGraveyardEntry_GetLvUp(t *testing.T) {
	up, err := CSV.GraveyardEntry.GetLvUp(2, 1)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	t.Logf("up :%+v", up)

}

func TestGraveyardEntry_GetStageUp(t *testing.T) {
	up, err := CSV.GraveyardEntry.GetStageUp(3, 5)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	t.Logf("up :%+v", up)
}

func TestGraveyardEntry_GetUnitCount(t *testing.T) {
	t.Logf("count :%d", CSV.GraveyardEntry.GetUnitCount(1, 1, UnitTypeCharacter))
	t.Logf("count :%d", CSV.GraveyardEntry.GetUnitCount(1, 1, UnitTypePopulation))
}

func TestGraveyardEntry_RefreshOutPut2(t *testing.T) {
	build := &common.UserGraveyardBuild{
		UserGraveyardBuildBase:  common.NewUserGraveyardBuildBase(2, 0, 0, nil),
		UserGraveyardTransition: nil,
		UserGraveyardProduce:    common.NewUserGraveyardProduce(),
	}

	build.ProduceStartAt = servertime.Now().Unix() - 45

	t.Logf("produce before:%+v", build.UserGraveyardProduce.Productions)

	CSV.RefreshOutPut(build, nil, false)
	build.ProduceStartAt = build.ProduceStartAt - 45
	CSV.RefreshOutPut(build, nil, false)

	//todo:
	t.Logf("produce after:%+v", build.UserGraveyardProduce.VOGraveyardProduce(nil))
}

func TestGraveyardEntry_RefreshOutPut(t *testing.T) {
	build := &common.UserGraveyardBuild{
		UserGraveyardBuildBase:  common.NewUserGraveyardBuildBase(2, 2, 0, nil),
		UserGraveyardTransition: nil,
		UserGraveyardProduce:    common.NewUserGraveyardProduce(),
	}

	build.ProduceStartAt = servertime.Now().Unix() - 60*60*24
	config, err := CSV.GraveyardEntry.GetBuffConfig(2)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	build.AddBuff(config.Id, config.Type, config.TypeContent, servertime.Now().Unix()-60*60*12)

	build.AddBuff(config.Id, config.Type, config.TypeContent, servertime.Now().Unix()-60*60*2)
	build.AccRecords.Records[0].AccAt = servertime.Now().Unix() - 60*60*2 + 1
	t.Logf("produce before:%+v", build.UserGraveyardProduce.Productions)
	characterStarMap := map[int32]int32{}

	characterStarMap[1001] = 2
	characterStarMap[1002] = 2

	CSV.RefreshOutPut(build, characterStarMap, false)
	//todo:
	t.Logf("produce after:%+v", build.UserGraveyardProduce.Productions)
}

func TestGraveyardEntry_buildAcc(t *testing.T) {
	build := &common.UserGraveyardBuild{
		UserGraveyardBuildBase:  common.NewUserGraveyardBuildBase(3, 1, 1, nil),
		UserGraveyardTransition: common.NewUserGraveyardTransition(1000, common.CurtainTypeCreate),
		UserGraveyardProduce:    common.NewUserGraveyardProduce(),
	}
	items := map[*GraveyardAccItem]int32{}
	acc, ok := CSV.Item.GetGraveyardAcc(10100001)
	if !ok {
		return
	}
	items[acc] = 1000
	consume, err := CSV.Acc(build, items, nil)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	t.Logf("consume:%v", consume)
	t.Logf("Transition:%+v", build.VOGraveyardTransition())

}

func TestGraveyardEntry_productAcc(t *testing.T) {
	build := &common.UserGraveyardBuild{
		// 经验药水
		UserGraveyardBuildBase:  common.NewUserGraveyardBuildBase(3, 1, 1, nil),
		UserGraveyardTransition: nil,
		UserGraveyardProduce:    common.NewUserGraveyardProduce(),
	}
	items := map[*GraveyardAccItem]int32{}
	acc, ok := CSV.Item.GetGraveyardAcc(10200001)
	if !ok {
		return
	}
	items[acc] = 1000
	consume, err := CSV.Acc(build, items, nil)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	t.Logf("consume:%v", consume)
	CSV.RefreshOutPut(build, nil, false)
	t.Logf("produce after:%+v", build.UserGraveyardProduce.Productions)

}

func TestGraveyardEntry_reduceProduceTimeAcc(t *testing.T) {
	build := &common.UserGraveyardBuild{
		// 武器
		UserGraveyardBuildBase:  common.NewUserGraveyardBuildBase(6, 1, 0, nil),
		UserGraveyardTransition: nil,
		UserGraveyardProduce:    common.NewUserGraveyardProduce(),
	}
	build.StartConsumeProduce(3)
	//build.PopulationCount=10
	items := map[*GraveyardAccItem]int32{}
	acc, ok := CSV.Item.GetGraveyardAcc(1611001)
	if !ok {
		return
	}
	items[acc] = 1
	consume, err := CSV.Acc(build, items, nil)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	t.Logf("consume:%v", consume)
	CSV.RefreshOutPut(build, nil, true)
	t.Logf("produce after:%+v", build.UserGraveyardProduce.VOGraveyardProduce(build))
	CSV.RefreshOutPut(build, nil, true)
	t.Logf("produce after:%+v", build.UserGraveyardProduce.VOGraveyardProduce(build))

}

func TestGraveyardEntry_RewardsWithLowUp(t *testing.T) {
	rewards := common.NewRewards()
	rewards.AddReward(common.NewReward(1, 100))
	rewards.AddReward(common.NewReward(1, 200))
	result := CSV.GraveyardEntry.RewardsWithLowUp(rewards, 781, 1222)
	t.Logf("result:%v", result)

}

func TestGraveyardEntry_CreateBuild(t *testing.T) {
	build := common.NewUserGraveyardBuild(2, nil, 3)
	t.Logf("Transition before:%+v", build.VOGraveyardTransition())
	//time.Sleep(3 * time.Second)
	t.Logf("Transition after:%+v", build.VOGraveyardTransition())

	t.Logf("CanBuildCount :%+v", CSV.GraveyardEntry.CanBuildCount(1, 2))

}

func TestGraveyardEntry_NewTimeSegments(t *testing.T) {
	now := servertime.Now().Unix()
	produce := common.NewUserGraveyardProduce()
	produce.ProduceStartAt = now - 100
	produce.AccRecords.AddRecord(10, now-10)
	//produce.AccRecords.AddRecord(10,now-10)
	segments := NewTimeSegments(produce)
	t.Logf("segments:%+v", segments)
	t.Logf("now:%+v", servertime.Now())

	t.Logf("tickTime:%+v", time.Unix(segments.getTickTime(90), 0))

	t.Logf("tickTime:%+v", time.Unix(segments.getTickTime(100), 0))
	t.Logf("tickTime:%+v", time.Unix(segments.getTickTime(110), 0))

}

func TestGraveyardEntry_RandomRewardHours(t *testing.T) {
	for i := 0; i < 1000; i++ {
		hours := CSV.GraveyardEntry.RandomRewardHours()
		t.Logf("RandomRewardHours:%+v", hours)
		now := servertime.GetHourWhen(servertime.Now().Unix())
		t.Logf("NowHour:%+v", now)

		num := CSV.GraveyardEntry.GetPlotRewardNum(hours, now)
		t.Logf("GetPlotRewardNum:%+v", num)

	}

}
