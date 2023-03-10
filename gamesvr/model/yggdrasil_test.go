package model

import (
	"context"
	"fmt"
	"gamesvr/manager"
	"math"
	"shared/common"
	"shared/csv/base"
	"shared/statistic/logreason"
	"shared/utility/coordinate"
	"testing"
)

func TestYggdrasil_InitAreaInitObj(t *testing.T) {
	TestUser.Yggdrasil.InitAreaInitObj(context.Background(), 1)

}

func TestYggdrasil_GetBlockInfoByEnter(t *testing.T) {

	TestUser.Yggdrasil.initTravelInfo(nil, 40)
	vo := TestUser.Yggdrasil.GetBlockInfoByEnter(context.Background())
	t.Log(vo)
}

func TestYggEntityTree_Append(t *testing.T) {

}

func TestYggdrasil_AddRewards(t *testing.T) {
	rewards := common.NewRewards()
	rewards.AddReward(common.NewReward(2, 10000))
	rewards.AddReward(common.NewReward(1111001, 30))
	rewards.AddReward(common.NewReward(1183001, 2))
	rewards.AddReward(common.NewReward(171000, 1))
	err := TestUser.Yggdrasil.AddRewards(context.Background(), TestUser, rewards, 1, logreason.EmptyReason())
	t.Log(err)
	t.Log(TestUser.VOYggdrasilResourceResult())
	t.Log(TestUser.VOResourceResult())

}

func TestYggdrasil_TakeBackPackGood(t *testing.T) {
	TestYggdrasil_AddRewards(t)
	err := TestUser.Yggdrasil.TakeBackPackGood(context.Background(), TestUser, TestUser.Yggdrasil.Pack.Values()...)
	t.Log(err)
	t.Log(TestUser.VOResourceResult())

}

func TestYggdrasil_CostRewards(t *testing.T) {
	TestYggdrasil_AddRewards(t)
	rewards := common.NewRewards()
	rewards.AddReward(common.NewReward(2, 10000))

	rewards.AddReward(common.NewReward(1183001, 2))
	rewards.AddReward(common.NewReward(171000, 1))
	TestUser.Yggdrasil.CostRewards(TestUser, rewards, logreason.EmptyReason())
	t.Log(TestUser.VOYggdrasilResourceResult())
	t.Log(TestUser.VOResourceResult())

}
func TestYggdrasil_Task(t *testing.T) {
	config, _ := manager.CSV.Yggdrasil.GetTaskConfig(8000)
	task, err := TestUser.Yggdrasil.Task.AcceptTask(context.Background(), TestUser, config)
	t.Log(task)
	t.Log(err)
	t.Log(TestUser.VOYggdrasilEntityChange(context.Background()))
	t.Log(TestUser.VOYggdrasilResourceResult())
}

func TestYggdrasilTask_Process(t *testing.T) {
	TestYggdrasil_Task(t)
	TestUser.Yggdrasil.Task.ProcessPos(context.Background(), TestUser, 9, *coordinate.NewPosition(-8, -21))
	t.Log(TestUser.Notifies.YggPush)
}
func TestYggdrasilTask_ChooseNext(t *testing.T) {
	TestYggdrasilTask_Process(t)
	task, err := TestUser.YggdrasilChooseNext(context.Background(), 1010102)
	t.Log(task)
	t.Log(err)
	t.Log(TestUser.VOYggdrasilEntityChange(context.Background()))
	t.Log(TestUser.VOYggdrasilResourceResult())
	t.Log(TestUser.VOResourceResult())

}

func TestYggBuildCreate(t *testing.T) {
	TestUser.Yggdrasil.InitAreaInitObj(context.Background(), 1)
	vo, err := TestUser.YggdrasilBuildCreate(context.Background(), 1)
	if err != nil {
		fmt.Println(err)
	}
	t.Log(vo)
}

func TestYggBuildC(t *testing.T) {
	TestUser.Yggdrasil.InitAreaInitObj(context.Background(), 1)
	vo, err := TestUser.YggdrasilBuildCreate(context.Background(), 1)
	if err != nil {
		fmt.Println(err)
	}
	t.Log(vo)
}

func TestYggBuildAddAp(t *testing.T) {

	TestUser.Yggdrasil.initTravelInfo(nil, 40)
	_, err := TestUser.YggdrasilBuildCreate(context.Background(), 1)
	if err != nil {
		fmt.Println(err)
	}

	vo, ap, err := TestUser.YggdrasilBuildAddAp(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	t.Log(vo, ap)
}

func TestYggBuildPeeing(t *testing.T) {
	_, err := TestUser.YggdrasilBuildCreate(context.Background(), 3)
	if err != nil {
		fmt.Println(err)
	}
	_, err = TestUser.YggdrasilBuildUsePeeing(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	t.Log(TestUser.VOYggdrasilEntityChange(context.Background()))
}

func TestYggMessage(t *testing.T) {
	vo, err := TestUser.YggdrasilMessageCreate(context.Background(), "hello world")
	if err != nil {
		fmt.Errorf(err.Error())
	}
	t.Log(vo)

	vo_new, err := TestUser.YggdrasilMessageUpdate(context.Background(), "hello again")
	if err != nil {
		fmt.Errorf(err.Error())
	}

	t.Log(vo_new)

	_, err = TestUser.YggdrasilMessageDestroy(context.Background(), vo_new.MessageUid)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	t.Log(*TestUser.Yggdrasil.TravelPos)

	isBlank := TestUser.Yggdrasil.Entities.IsBlankPos(*TestUser.Yggdrasil.TravelPos)

	t.Log(isBlank)

	v, ok := TestUser.Yggdrasil.Entities.EntityFinder.Entities[*TestUser.Yggdrasil.TravelPos]

	fmt.Println(*v, ok)
}

func TestNewYggdrasilMail(t *testing.T) {
	pageNum := 40
	mailBox := NewYggdrasilMailBox()
	for i := 0; i < 1; i++ {
		mailBox.AddOne(context.Background(), 0, fmt.Sprint(i), common.NewRewards())
	}
	t.Log(len(mailBox.PagingSearch(10, pageNum)))
	t.Log(len(mailBox.PagingSearch(1, pageNum)))

	for i := 0; i < 1000; i++ {
		mailBox.AddOne(context.Background(), 0, fmt.Sprint(i), common.NewRewards())
	}
	var offset int64 = math.MaxInt64

	for {
		search := mailBox.PagingSearch(offset, pageNum)
		t.Log(len(search))

		if len(search) < 40 {
			break
		}
		offset = search[len(search)-1].Uid

	}

}

func TestYggdrasil_LightAround(t *testing.T) {
	TestUser.Yggdrasil.TravelPos = coordinate.NewPosition(1111, 13)
	around, updates := TestUser.Yggdrasil.LightAround(context.Background())
	t.Log(around)
	t.Log(updates)
}

func TestUser_YggdrasilMatch(t *testing.T) {
	// todo：测试用
	matchUserIds, err := manager.Global.SMembers(context.Background(), "matchUserIds")

	if err != nil {
		t.Log(err)
		return
	}

	otherMembers := make([]int64, 0, len(matchUserIds)-1)
	for _, str := range matchUserIds {
		id, ok := base.String2Int64(str)
		if !ok {
			continue
		}
		if id == TestUser.GetUserId() {
			continue
		}
		otherMembers = append(otherMembers, id)
	}

	area := TestUser.Yggdrasil.Areas.getByCreate(context.Background(), TestUser.Yggdrasil, 1)
	area.IsTaskDone = true
	err = TestUser.YggdrasilMatch(context.Background(), otherMembers)
	t.Log(err)
}

func BenchmarkYggdrasil_YggInit(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmpUser := NewUser(1001)
		tmpUser.InitForCreate(context.Background())
		tmpUser.Init(context.Background())
		tmpUser.RewardsResult.Clear()
		TestUser.Yggdrasil.GetBlockInfoByEnter(context.Background())

	}
	b.ReportAllocs()

}

func TestYggdrasilTask_InitEnv(t *testing.T) {
	yggdrasilTask := NewYggdrasilTask()
	subtask, err := manager.CSV.Yggdrasil.GetSubTaskConfig(1010202)
	t.Log(err)

	env, err := subtask.Envs.RandomEnv()
	t.Log(err)

	err = yggdrasilTask.InitEnv(context.Background(), TestUser, subtask.Id, env)
	t.Log(err)
}

func TestYggdrasilDispatch_GetInfo(t *testing.T) {
	ctx := context.Background()
	TestUser.Yggdrasil.InitAreaInitObj(ctx, 1)
	TestUser.Yggdrasil.Areas.getByCreate(ctx, TestUser.Yggdrasil, 2)

	taskInfoBefore, err := TestUser.YggdrasilDispatchGetInfo()
	if err != nil {
		t.Log(err)
	}
	t.Log("DispatchTask: ", taskInfoBefore)

	taskInfoAfter, err := TestUser.YggdrasilDispatchGetInfo()
	if err != nil {
		t.Log(err)
	}
	t.Log("DispatchTask: ", taskInfoAfter)

}

func TestYggdrasil_GetTotalIntimacy(t *testing.T) {
	intimacy, err := TestUser.GetTotalIntimacy(context.Background())
	t.Log(intimacy)
	t.Log(err)

}

func TestYggdrasil_GuildDispatchGenerate(t *testing.T) {
	taskIds, err := TestUser.YggDispatchGenerateGuildTasks(1, 2, []int32{2, 3})
	t.Log(taskIds)
	t.Log(err)

	taskIds, err = TestUser.YggDispatchGenerateGuildTasks(5, 5, []int32{2, 3})
	t.Log(taskIds)
	t.Log(err)

	taskIds, err = TestUser.YggDispatchGenerateGuildTasks(5, 5, []int32{4, 5})
	t.Log(taskIds)
	t.Log(err)
}

func TestYggdrasilTask_InitSubTask(t *testing.T) {
	task, err := TestUser.Yggdrasil.Task.InitSubTask(context.Background(), TestUser, 2010101)
	t.Log(task)
	t.Log(err)
}

func TestYggdrasilTask_InitSubTask2(t *testing.T) {
	config, _ := manager.CSV.Yggdrasil.GetTaskConfig(10101)
	task, err := TestUser.Yggdrasil.Task.AcceptTask(context.Background(), TestUser, config)
	t.Log(task)
	t.Log(err)
}

func TestYggdrasilTask_InitSubTask3(t *testing.T) {
	config, _ := manager.CSV.Yggdrasil.GetTaskConfig(10202)
	ctx := context.Background()
	_, err := TestUser.Yggdrasil.Task.InitSubTask(ctx, TestUser, config.NextSubTaskId)
	if err != nil {
		t.Log(err)
	}
	info := TestUser.Yggdrasil.Task.TaskInProgress[10202]

	err = TestUser.Yggdrasil.Task.ForceCompleteSubTask(ctx, TestUser, info.Base, info)
	if err != nil {
		t.Log(err)
	}
	t.Log(TestUser.VOYggdrasilEntityChange(ctx))
	TestUser.Yggdrasil.Task.ChooseNext(ctx, TestUser, 1020202)
	err = TestUser.Yggdrasil.Task.ForceCompleteSubTask(ctx, TestUser, info.Base, info)
	if err != nil {
		t.Log(err)
	}
	t.Log(TestUser.VOYggdrasilEntityChange(ctx))

	t.Log(err)
}

func TestYggdrasilTask_InitSubTask4(t *testing.T) {
	ctx := context.Background()
	_, err := TestUser.Yggdrasil.Task.InitSubTask(ctx, TestUser, 1010301)
	if err != nil {
		t.Log(err)
	}
	TestUser.Yggdrasil.Task.ForceCompleteTask(ctx, TestUser, 20111)
	t.Log(TestUser.VOYggdrasilEntityChange(ctx))
	TestUser.Yggdrasil.Task.ForceCompleteTask(ctx, TestUser, 20112)

	t.Log(err)
}

func TestYggdrasil_EnvRemove(t *testing.T) {

	TestUser.Yggdrasil.Entities.CreateEnvTerrain(1, 1, 2, TestUser.Yggdrasil)
	TestUser.Yggdrasil.Entities.CreateEnvTerrain(1, 1, 1, TestUser.Yggdrasil)
	TestUser.Yggdrasil.Entities.CreateEnvTerrain(1, 1, 3, TestUser.Yggdrasil)
	t.Log(TestUser.Yggdrasil.Entities.EnvTerrains)

	TestUser.Yggdrasil.Entities.RemoveEnvTerrainDeleteAt(1, TestUser.Yggdrasil)
	t.Log(TestUser.Yggdrasil.Entities.EnvTerrains)

}
