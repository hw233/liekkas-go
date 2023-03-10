package model

import (
	"context"
	"fmt"
	"gamesvr/manager"
	"shared/csv/static"
	"testing"
	"time"
)

func TestExploreEnterChapter(t *testing.T) {
	chapterId := int32(1)
	err := TestUser.EnterChapterMap(chapterId)
	if err != nil {
		t.Errorf("failed, %s", err.Error())
		return
	}

	fmt.Printf("user.ExploreInfo: %v\n", TestUser.ExploreInfo)
}

func TestExploreChapterReward(t *testing.T) {
	chapterId := int32(1)
	err := TestUser.EnterChapterMap(chapterId)
	if err != nil {
		t.Errorf("failed, %s", err.Error())
		return
	}

	chapter, _ := TestUser.ChapterInfo.GetChapter(chapterId)
	chapter.AddScore(100)

	TestUser.RewardsResult.Clear()

	rewardId := int32(1)
	err = TestUser.ReceiveChapterReward(rewardId)
	if err != nil {
		t.Errorf("failed, %s", err.Error())
		return
	}

	fmt.Printf("chapter: %v\n", chapter)
	fmt.Printf("TestUser.VOResourceResult(): %v\n", TestUser.VOResourceResult())
}

func TestExploreNPC(t *testing.T) {
	// npcId := int32(1)
	// levelId := int32(10207)

	// TestUser.onLevelEnd(levelId, true, []int32{}, []int32{})
	// TestUser.RewardsResult.Clear()

	// err := TestUser.ExploreNPCInteraction(npcId, 1)
	// if err != nil {
	// 	t.Errorf("failed, %s", err.Error())
	// 	return
	// }

	// rewardResult := TestUser.VOResourceResult()
	// fmt.Printf("rewardResult: %v\n", rewardResult)

	// for _, eep := range TestUser.ExploreInfo.EventPoints {
	// 	fmt.Printf("eep: %v\n", eep)
	// }
}

func TestExploreRewardPoint(t *testing.T) {
	rewardPointId := int32(1)

	err := TestUser.ExploreRewardPointInteraction(rewardPointId)
	if err != nil {
		t.Errorf("failed, %s", err.Error())
		return
	}

	rewardResult := TestUser.VOResourceResult()
	fmt.Printf("rewardResult: %v\n", rewardResult)

	for _, eep := range TestUser.ExploreInfo.EventPoints {
		fmt.Printf("eep: %v\n", eep)
	}
}
func TestExploreMonster(t *testing.T) {
	monsterId := int32(1)

	monsterCfg, err := manager.CSV.ExploreEntry.GetExploreMonster(monsterId)
	if err != nil {
		t.Errorf("failed, %s", err.Error())
		return
	}

	levelId := monsterCfg.Monstermap
	_, _, err = TestUser.StartLevel(levelId, nil,
		static.BattleTypeExploreMonster, []int64{int64(monsterId)})
	if err != nil {
		t.Errorf("failed, %s", err.Error())
		return
	}

	battleData := &BattleVerifyData{
		IsWin:                true,
		BattleInput:          "",
		BattleOutput:         "",
		CompleteTargets:      []int32{0, 1},
		CompleteAchievements: []int32{0, 1},
	}
	_, err = TestUser.LevelEnd(context.Background(), levelId, battleData)
	if err != nil {
		t.Errorf("failed, %s", err.Error())
		return
	}

	for _, eep := range TestUser.ExploreInfo.EventPoints {
		fmt.Printf("eep: %v\n", eep)
	}
}

func TestExploreResource(t *testing.T) {
	resourceId := int32(1)

	resourceCfg, _ := manager.CSV.ExploreEntry.GetResource(resourceId)
	levelId := resourceCfg.Monstermap
	_, _, err := TestUser.StartLevel(levelId, nil,
		static.BattleTypeExploreResource, []int64{})
	if err != nil {
		t.Errorf("failed, %s", err.Error())
		return
	}

	battleData := &BattleVerifyData{
		IsWin:                true,
		BattleInput:          "",
		BattleOutput:         "",
		CompleteTargets:      []int32{1, 2},
		CompleteAchievements: []int32{1, 2},
	}
	_, err = TestUser.LevelEnd(context.Background(), levelId, battleData)
	if err != nil {
		t.Errorf("failed, %s", err.Error())
		return
	}
	TestUser.RewardsResult.Clear()

	err = TestUser.ExploreStartCollectResource(resourceId)
	if err != nil {
		t.Errorf("failed, %s", err.Error())
		return
	}

	resourcePoint, _ := TestUser.ExploreInfo.GetResourcePoint(resourceId)
	fmt.Printf("resourcePoint: %v\n", resourcePoint)

	time.Sleep(time.Second * time.Duration(resourceCfg.Time))

	err = TestUser.ExploreFinishResourceCollect(resourceId)
	if err != nil {
		t.Errorf("failed, %s", err.Error())
		return
	}

	fmt.Printf("resourcePoint: %v\n", resourcePoint)

	fmt.Printf("TestUser.VOResourceResult(): %v\n", TestUser.VOResourceResult())
}
