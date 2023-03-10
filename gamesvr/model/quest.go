package model

import (
	"shared/protobuf/pb"
	"shared/utility/servertime"
	"time"
)

type QuestState int32

const (
	QuestStateWaiting     QuestState = 1
	QuestStateProgressing QuestState = 2
	QuestStateComplete    QuestState = 3
	QuestStateRecieve     QuestState = 4
)

type QuestModule int32

type Quest struct {
	Id           int32      `json:"id"`
	State        QuestState `json:"state"`
	Progress     int32      `json:"progress"`
	CompleteTime int64      `json:"complete_time"`
}

type QuestActivity struct {
	Module       QuestModule    `json:"module"`
	Point        int32          `json:"point"`
	RewardRecord map[int32]bool `json:"reward_record"`
}

type QuestPack struct {
	Quests     map[int32]*Quest               `json:"quests"`
	Activities map[QuestModule]*QuestActivity `json:"activities"`
	*DailyRefreshChecker
}

func NewQuest(id int32) *Quest {
	return &Quest{
		Id:           id,
		State:        QuestStateWaiting,
		Progress:     0,
		CompleteTime: 0,
	}
}

func NewQuestActivity(module QuestModule) *QuestActivity {
	return &QuestActivity{
		Module:       module,
		Point:        0,
		RewardRecord: map[int32]bool{},
	}
}

func NewQuestPack() *QuestPack {
	return &QuestPack{
		Quests:              map[int32]*Quest{},
		Activities:          map[QuestModule]*QuestActivity{},
		DailyRefreshChecker: NewDailyRefreshChecker(),
	}
}

//----------------------------------------
//QuestPack
//----------------------------------------
func (qp *QuestPack) GetAllQuest() map[int32]*Quest {
	return qp.Quests
}

func (qp *QuestPack) GetQuest(id int32) (*Quest, bool) {
	quest, ok := qp.Quests[id]
	return quest, ok
}

func (qp *QuestPack) IsQuestRecieved(questId int32) bool {
	quest, ok := qp.GetQuest(questId)
	if !ok {
		return false
	}

	return quest.IsRecieved()
}

func (qp *QuestPack) IsQuestCompleted(questId int32) bool {
	quest, ok := qp.GetQuest(questId)
	if !ok {
		return false
	}

	return quest.IsCompleted()
}

func (qp *QuestPack) AddQuest(questId int32) *Quest {
	quest := NewQuest(questId)

	qp.Quests[questId] = quest

	return quest
}

func (qp *QuestPack) RemoveQuest(questId int32) {
	delete(qp.Quests, questId)
}

func (qp *QuestPack) GetActivityData(module QuestModule) *QuestActivity {
	activityData, ok := qp.Activities[module]
	if !ok {
		activityData = NewQuestActivity(module)
		qp.Activities[module] = activityData
	}

	return activityData
}

func (qp *QuestPack) GetLastWeeklyRefreshTime() int64 {
	lastDailyTime := time.Unix(qp.GetLastDailyRefreshTime(), 0)
	return WeekRefreshTime(lastDailyTime).Unix()
}

func (qp *QuestPack) VOQuestInfo() *pb.VOQuestInfo {
	questsData := make([]*pb.VOQuest, 0)
	for _, quest := range qp.Quests {
		questsData = append(questsData, quest.VOQuest())
	}

	questActivitiesData := make([]*pb.VOQuestActivity, 0, len(qp.Activities))
	for _, questActivity := range qp.Activities {
		questActivitiesData = append(questActivitiesData, questActivity.VOQuestActivity())
	}

	return &pb.VOQuestInfo{
		Quests:                questsData,
		QuestActivities:       questActivitiesData,
		NextDailyRefreshTime:  qp.GetLastDailyRefreshTime() + servertime.SecondPerDay,
		NextWeeklyRefreshTime: qp.GetLastWeeklyRefreshTime() + servertime.SecondPerWeek,
	}
}

//----------------------------------------
//QuestActivity
//----------------------------------------
func (qa *QuestActivity) AddPoint(point int32) {
	qa.Point = qa.Point + point
}

func (qa *QuestActivity) GetPoint() int32 {
	return qa.Point
}

func (qa *QuestActivity) IsRewardReceived(rewardId int32) bool {
	_, ok := qa.RewardRecord[rewardId]
	return ok
}

func (qa *QuestActivity) RecordReward(rewardId int32) {
	qa.RewardRecord[rewardId] = true
}

func (qa *QuestActivity) Reset() {
	qa.Point = 0
	qa.RewardRecord = map[int32]bool{}
}

func (qa *QuestActivity) VOQuestActivity() *pb.VOQuestActivity {
	rewardRecords := make([]int32, 0, len(qa.RewardRecord))

	for id := range qa.RewardRecord {
		rewardRecords = append(rewardRecords, id)
	}

	return &pb.VOQuestActivity{
		Module:       int32(qa.Module),
		Point:        qa.Point,
		RewardRecord: rewardRecords,
	}
}

//----------------------------------------
//Quest
//----------------------------------------
func (q *Quest) GetId() int32 {
	return q.Id
}

func (q *Quest) GetProgress() int32 {
	return q.Progress
}

func (q *Quest) UpdateProgress(newProgress int32) {
	q.Progress = newProgress
}

func (q *Quest) Restart() {
	q.Progress = 0
	q.State = QuestStateProgressing
}

func (q *Quest) Complete() {
	q.State = QuestStateComplete
	q.CompleteTime = servertime.Now().Unix()
}

func (q *Quest) Recieve() {
	q.State = QuestStateRecieve
}

func (q *Quest) IsProgressing() bool {
	return q.State == QuestStateProgressing
}

func (q *Quest) IsCompleted() bool {
	return q.State == QuestStateComplete
}

func (q *Quest) IsRecieved() bool {
	return q.State == QuestStateRecieve
}

func (q *Quest) IsCounting() bool {
	return q.State == QuestStateProgressing || q.State == QuestStateComplete
}

func (q *Quest) VOQuest() *pb.VOQuest {
	return &pb.VOQuest{
		Id:           q.Id,
		State:        int32(q.State),
		Progress:     q.Progress,
		CompleteTime: q.CompleteTime,
	}
}
