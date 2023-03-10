package model

import (
	"context"
	"guild/manager"
	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/rand"
	"shared/utility/servertime"
	"time"
)

func (g *Guild) GuildTaskRefresh() {
	weekRefreshTime := servertime.WeekOffsetZeroTime(servertime.Now(), manager.CSV.GlobalEntry.DailyRefreshTimeOffset())

	big := manager.CSV.GlobalEntry.GuildTaskRefreshWeekDay[1] - 1
	small := manager.CSV.GlobalEntry.GuildTaskRefreshWeekDay[0] - 1
	biggerRefreshTime := weekRefreshTime.Add(24 * time.Hour * time.Duration(big))
	smallerRefreshTime := weekRefreshTime.Add(24 * time.Hour * time.Duration(small))
	// 如果当前的时间已经超过了本周的周五刷新时间
	if servertime.Now().After(biggerRefreshTime) {
		// 如果上次刷新的时间在周五之前
		if time.Unix(g.TaskLastRefreshTime, 0).Before(biggerRefreshTime) {
			g.RefreshTask()
		}
		// 如果当前的时间已经超过了本周的周二刷新时间
	} else if servertime.Now().After(smallerRefreshTime) {
		if time.Unix(g.TaskLastRefreshTime, 0).Before(smallerRefreshTime) {
			g.RefreshTask()
		}
	}
}

func (g *Guild) RefreshTask() {
	normal, separate := manager.CSV.Guild.GetAllTaskIds()
	// fmt.Println("===============normal: ", normal, "=======separate: ", separate)
	g.Tasks = make([]*common.GuildTask, 0, GuildNormalTaskNum+GuildSeparateTaskNum)
	g.GenerateTask(normal, GuildNormalTaskNum)
	g.GenerateTask(separate, GuildSeparateTaskNum)
	g.TaskLastRefreshTime = servertime.Now().Unix()

}

func (g *Guild) GenerateTask(taskIds []int32, num int) {
	if num > len(taskIds) {
		num = len(taskIds)
	}
	probs := make([]int32, 0, len(taskIds))
	for i := 0; i < len(taskIds); i++ {
		probs = append(probs, 1)
	}
	taskIndex := rand.UniquePerm(num, probs)
	// fmt.Println("=====afterrand", taskIndex)

	for _, index := range taskIndex {
		g.Tasks = append(g.Tasks, common.NewGuildTask(taskIds[index]))
	}
}

// ---------------------公会任务---------------------

func (g *Guild) VOGuildTaskList(ctx context.Context, userId int64) []*pb.VOGuildTask {
	result := make([]*pb.VOGuildTask, 0, len(g.Tasks))
	// g.RefreshTask()
	for _, task := range g.Tasks {
		// fmt.Println("====================task", *task)
		for _, member := range g.Members {
			if member.UserID == userId {
				isReceived, ok := member.TaskRewardsReceived[task.ID]
				if !ok {
					isReceived = false
				}
				result = append(result, task.VOGuildTask(isReceived))
				// fmt.Println("===========task:", *task)
			}
		}
	}
	return result
}
func (g *Guild) ReceiveTaskRewards(taskId int32, userId int64) (int32, error) {
	var activation int32
	for _, task := range g.Tasks {
		if task.ID == taskId {
			taskConfig, err := manager.CSV.Guild.GetTaskConfig(taskId)
			if err != nil {
				return 0, errors.WrapTrace(err)
			}
			activation = taskConfig.Activation
			if !task.Count.Enough(taskConfig.Count) {
				return 0, errors.Swrapf(common.ErrGuildTaskNotFinish, taskId, task.Count.Value(), taskConfig.Count)
			}
			for _, member := range g.Members {
				if member.UserID == userId {
					isReceived, ok := member.TaskRewardsReceived[task.ID]
					if !ok {
						isReceived = false
					}
					if isReceived {
						return 0, errors.Swrapf(common.ErrGuildTaskRewardsHasReceived, taskId)
					}
					g.AddGuildMemberActivation(member, activation)
					isReceived = true
					member.TaskRewardsReceived[task.ID] = isReceived
				}
			}
		}
	}
	return activation, nil
}

func (g *Guild) GuildTaskAddProgress(taskId, value int32) error {
	for _, task := range g.Tasks {
		if task.ID == taskId {
			taskConfig, err := manager.CSV.Guild.GetTaskConfig(taskId)
			if err != nil {
				return errors.WrapTrace(err)
			}
			task.Count.Plus(value)
			if task.Count.Enough(taskConfig.Count) {
				task.Count.SetValue(taskConfig.Count)
			}
		}
	}
	return nil
}
