package model

import (
	"context"
	"fmt"
	"shared/csv/static"
	"shared/utility/servertime"
	"testing"
)

func TestQuestParamsConvert(t *testing.T) {
	paramsObj := struct {
		I32 int32
		I64 int64
	}{}

	params := make([]interface{}, 0, 2)
	params = append(params, 32)
	params = append(params, 64)

	err := questConvertParams(&paramsObj, params)
	if err != nil {
		t.Errorf("failed, %s", err.Error())
	}

	fmt.Printf("paramsObj: %+v\n", paramsObj)
}

func TestQuestTrigger(t *testing.T) {
	user := NewUser(100001)
	user.InitForCreate(context.Background())
	user.OnOnline(servertime.Now().Unix())

	battleData := &BattleVerifyData{
		IsWin: true,
	}

	curLevelCache := &LevelCacheInfo{
		LevelId:      10101,
		SystemType:   static.BattleTypeTower,
		SystemParams: []int64{1, 1},
	}

	user.LevelsInfo.CurLevel = *curLevelCache
	user.LevelEnd(context.Background(), 10101, battleData)

	fmt.Printf("user.Notifies.QuestUpdateNotify: %+v\n", user.Notifies.QuestUpdateNotify)
}
