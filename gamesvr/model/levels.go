package model

import (
	"shared/protobuf/pb"
)

const (
	LevelTargetCount = 3
)

type Level struct {
	Id                   int32          `json:"id"`
	TodayPassTimes       int32          `json:"today_pass_times"`
	TotalPassTimes       int32          `json:"total_pass_times"`
	CompleteTargets      map[int32]bool `json:"complete_targets"`
	CompleteAchievements map[int32]bool `json:"complete_achievements"`
}

type LevelCacheInfo struct {
	LevelId      int32                 `json:"level_id"`
	Formation    *pb.VOBattleFormation `json:"formation"`
	SystemType   int32                 `json:"system_type"`
	SystemParams []int64               `json:"system_params"`
	Seed         int32                 `json:"seed"`
}

type LevelsInfo struct {
	Levels   map[int32]*Level `json:"levels"`
	CurLevel LevelCacheInfo   `json:"cur_level"`
	*DailyRefreshChecker
}

func NewLevelsInfo() *LevelsInfo {
	return &LevelsInfo{
		Levels: map[int32]*Level{},
		CurLevel: LevelCacheInfo{
			LevelId:    0,
			Formation:  nil,
			SystemType: 0,
		},
		DailyRefreshChecker: NewDailyRefreshChecker(),
	}
}

func (li *LevelsInfo) GetLevel(levelId int32) (*Level, bool) {
	level, ok := li.Levels[levelId]
	return level, ok
}

func (li *LevelsInfo) GetOrCreateLevel(levelId int32) *Level {
	level, ok := li.GetLevel(levelId)
	if !ok {
		level = NewLevel(levelId)
		li.Levels[levelId] = level
	}

	return level
}

func (li *LevelsInfo) GetCurLevelId() int32 {
	return li.CurLevel.LevelId
}

func (li *LevelsInfo) GetCurLevel() LevelCacheInfo {
	return li.CurLevel
}

func (li *LevelsInfo) ClearCurLevel() {
	li.CurLevel.LevelId = 0
	li.CurLevel.Formation = nil
	li.CurLevel.SystemType = 0
	li.CurLevel.SystemParams = nil
	li.CurLevel.Seed = 0
}

func (li *LevelsInfo) SetCurLevel(levelId int32, formation *pb.VOBattleFormation,
	systemType int32, systemParams []int64, seed int32) {
	li.CurLevel.LevelId = levelId
	li.CurLevel.Formation = formation
	li.CurLevel.SystemType = systemType
	li.CurLevel.SystemParams = systemParams
	li.CurLevel.Seed = seed
}

func (li *LevelsInfo) GetCurLevelHeroId() int32 {
	formation := li.CurLevel.Formation
	if formation == nil {
		return 0
	}

	if formation.BattleHero == nil {
		return 0
	}

	return formation.BattleHero.HeroId
}

func (li *LevelsInfo) LevelPass(levelId int32, targets, achievements []int32) ([]int32, []int32) {
	level := li.GetOrCreateLevel(levelId)

	return level.Pass(targets, achievements)
}

func (li *LevelsInfo) IsLevelPassed(levelId int32) bool {
	level, ok := li.GetLevel(levelId)
	if !ok {
		return false
	}

	return level.IsPassed()
}

func (li *LevelsInfo) VOLevelInfo() *pb.VOLevelInfo {
	levelInfo := &pb.VOLevelInfo{
		Levels:   make([]*pb.VOLevel, 0, len(li.Levels)),
		CurLevel: li.CurLevel.LevelId,
	}

	for _, level := range li.Levels {
		levelInfo.Levels = append(levelInfo.Levels, level.VOLevel())
	}

	return levelInfo
}

func NewLevel(levelId int32) *Level {
	return &Level{
		Id:                   levelId,
		TodayPassTimes:       0,
		TotalPassTimes:       0,
		CompleteTargets:      map[int32]bool{},
		CompleteAchievements: map[int32]bool{},
	}
}

func (l *Level) Pass(targets, achievements []int32) ([]int32, []int32) {
	l.TodayPassTimes = l.TodayPassTimes + 1
	l.TotalPassTimes = l.TotalPassTimes + 1

	newTargets := make([]int32, 0)
	for _, index := range targets {
		_, ok := l.CompleteTargets[index]
		if !ok {
			l.CompleteTargets[index] = true
			newTargets = append(newTargets, index)
		}
	}

	newAchievements := make([]int32, 0)
	for _, index := range achievements {
		_, ok := l.CompleteAchievements[index]
		if !ok {
			l.CompleteAchievements[index] = true
			newAchievements = append(newAchievements, index)
		}
	}

	return newTargets, newAchievements
}

func (l *Level) Sweep(times int32) {
	l.TodayPassTimes = l.TodayPassTimes + times
	l.TotalPassTimes = l.TotalPassTimes + times
}

func (l *Level) GetStar() int32 {
	return int32(len(l.CompleteTargets))
}

func (l *Level) IsFirstPass() bool {
	return l.TotalPassTimes == 1
}

func (l *Level) IsPassed() bool {
	return l.TotalPassTimes >= 1
}

func (l *Level) GetTodayPass() int32 {
	return l.TodayPassTimes
}

func (l *Level) GetTotalPass() int32 {
	return l.TotalPassTimes
}

func (l *Level) ResetDailyTime() {
	l.TodayPassTimes = 0
}

func (l *Level) Reset() {
	l.TotalPassTimes = 0
	l.TodayPassTimes = 0
	l.CompleteTargets = map[int32]bool{}
	l.CompleteAchievements = map[int32]bool{}
}

func (l *Level) IsAllTargetCleared() bool {
	return len(l.CompleteTargets) >= LevelTargetCount
}

func (l *Level) VOLevel() *pb.VOLevel {
	levelData := &pb.VOLevel{
		LevelId:              l.Id,
		TodayPassTimes:       l.TodayPassTimes,
		TotalPassTimes:       l.TotalPassTimes,
		CompleteTargets:      make([]int32, 0, len(l.CompleteTargets)),
		CompleteAchievements: make([]int32, 0, len(l.CompleteAchievements)),
	}

	for index := range l.CompleteTargets {
		levelData.CompleteTargets = append(levelData.CompleteTargets, index)
	}

	for index := range l.CompleteAchievements {
		levelData.CompleteAchievements = append(levelData.CompleteAchievements, index)
	}

	return levelData
}
