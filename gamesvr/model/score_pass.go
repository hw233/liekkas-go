package model

import "shared/protobuf/pb"

const (
	ScorePassStateWaiting     = 1
	ScorePassStateProgressing = 2
	ScorePassStateEnd         = 3
)

type ScorePassSeason struct {
	Id        int32                     `json:"id"`
	Phases    map[int32]*ScorePassPhase `json:"phases"`
	State     int32                     `json:"state"`
	StartTime int64                     `json:"start_time"`
}

type ScorePassPhase struct {
	Id           int32          `json:"id"`
	Groups       map[int32]bool `json:"groups"`
	RewardRecord map[int32]bool `json:"reward_record"`
	StartTime    int64          `json:"start_time"`
}

type ScorePassInfo struct {
	Seasons map[int32]*ScorePassSeason `json:"seasons"`
}

func NewScorePassInfo() *ScorePassInfo {
	return &ScorePassInfo{
		Seasons: map[int32]*ScorePassSeason{},
	}
}

//ScorePassInfo
func (spi *ScorePassInfo) StartSeason(id int32, startTime int64) *ScorePassSeason {
	season := &ScorePassSeason{
		Id:        id,
		Phases:    map[int32]*ScorePassPhase{},
		State:     ScorePassStateProgressing,
		StartTime: startTime,
	}

	spi.Seasons[id] = season

	return season
}

func (spi *ScorePassInfo) GetSeason(id int32) (*ScorePassSeason, bool) {
	season, ok := spi.Seasons[id]

	return season, ok
}

func (spi *ScorePassInfo) VOScorePassInfo() *pb.VOScorePassInfo {
	data := &pb.VOScorePassInfo{
		Seasons: make([]*pb.VOScorePassSeason, 0),
	}

	for _, season := range spi.Seasons {
		if season.IsProgressing() {
			data.Seasons = append(data.Seasons, season.VOScorePassSeason())
		}
	}

	return data
}

//ScorePassSeason
func (sps *ScorePassSeason) StartPhase(id int32, startTime int64) *ScorePassPhase {
	phase := &ScorePassPhase{
		Id:           id,
		Groups:       map[int32]bool{},
		RewardRecord: map[int32]bool{},
		StartTime:    startTime,
	}

	sps.Phases[id] = phase
	return phase
}

func (sps *ScorePassSeason) GetPhase(id int32) (*ScorePassPhase, bool) {
	phase, ok := sps.Phases[id]
	return phase, ok
}

func (sps *ScorePassSeason) End() {
	sps.State = ScorePassStateEnd
}

func (sps *ScorePassSeason) IsProgressing() bool {
	return sps.State == ScorePassStateProgressing
}

func (sps *ScorePassSeason) IsEnd() bool {
	return sps.State == ScorePassStateEnd
}

func (sps *ScorePassSeason) GetStartTime() int64 {
	return sps.StartTime
}

func (sps *ScorePassSeason) VOScorePassSeason() *pb.VOScorePassSeason {
	data := &pb.VOScorePassSeason{
		Id:     sps.Id,
		Phases: make([]*pb.VOScorePassPhase, 0, len(sps.Phases)),
		IsEnd:  sps.IsEnd(),
	}

	for _, phase := range sps.Phases {
		data.Phases = append(data.Phases, phase.VOScorePassPhase())
	}

	return data
}

//ScorePassPhase
func (spp *ScorePassPhase) StartGroup(groupId int32) {
	spp.Groups[groupId] = true
}

func (spp *ScorePassPhase) IsGroupStarted(groupId int32) bool {
	_, ok := spp.Groups[groupId]
	return ok
}

func (spp *ScorePassPhase) RecordReward(rewardId int32) {
	spp.RewardRecord[rewardId] = true
}

func (spp *ScorePassPhase) IsRewardReceived(rewardId int32) bool {
	_, ok := spp.RewardRecord[rewardId]
	return ok
}

func (spp *ScorePassPhase) GetStartTime() int64 {
	return spp.StartTime
}

func (spp *ScorePassPhase) VOScorePassPhase() *pb.VOScorePassPhase {
	data := &pb.VOScorePassPhase{
		Id:           spp.Id,
		RewardRecord: make([]int32, 0, len(spp.RewardRecord)),
		Groups:       make([]int32, 0, len(spp.Groups)),
	}

	for rewardId := range spp.RewardRecord {
		data.RewardRecord = append(data.RewardRecord, rewardId)
	}

	for groupId := range spp.Groups {
		data.Groups = append(data.Groups, groupId)
	}

	return data
}
