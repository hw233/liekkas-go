package model

import (
	"context"
	"fmt"
	"shared/protobuf/pb"
	"testing"
)

func TestLevelPass(t *testing.T) {
	user := NewUser(1001)
	user.InitForCreate(context.Background())

	levelId := int32(501001)
	_, _, err := user.StartLevel(levelId, &pb.VOBattleFormation{}, 0, []int64{})
	if err != nil {
		t.Errorf("failed, %s", err.Error())
		return
	}

	fmt.Printf("user.LevelsInfo.GetCurLevel(): %v\n", user.LevelsInfo.GetCurLevel())

	battleData := &BattleVerifyData{
		IsWin:                true,
		BattleInput:          "",
		BattleOutput:         "",
		CompleteTargets:      []int32{0, 1},
		CompleteAchievements: []int32{0, 1},
	}
	resouceResult, err := user.LevelEnd(context.Background(), levelId, battleData)
	if err != nil {
		t.Errorf("failed, %s", err.Error())
		return
	}

	level, _ := user.LevelsInfo.GetLevel(levelId)
	fmt.Printf("level: %+v\n", level)

	fmt.Printf("user.LevelsInfo.GetCurLevel(): %+v\n", user.LevelsInfo.GetCurLevel())

	fmt.Printf("resouceResult: %+v\n", resouceResult)
}
