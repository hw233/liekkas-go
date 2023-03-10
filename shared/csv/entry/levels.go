package entry

import (
	"shared/common"
	"shared/csv/static"
	"shared/utility/errors"
	"shared/utility/transfer"
	"sync"
)

type Level struct {
	Id                int32
	ChapterId         int32
	UnlockConditions  *common.Conditions `rule:"conditions"`
	TeamExp           int32
	CharacterExp      int32
	Cost              *common.Rewards `ignore:"true"`
	BattleID          int32
	AchievementsDrops [][]int32 `ignore:"true"`
	TargetDrops       []int32
	FirstDrop         []int32
	NormalDrop        []int32
	SweepDrop         []int32
	YggdrasilDrop     []int32
	ChallengeTimes    int32
	RefreshType       int32
	AchievementsCount int32 `ignore:"true"`
	TargetCount       int32 `ignore:"true"`

	ChapterStory            int32 // 战前vn
	BattleStory             int32 // 战中vn
	WinStory                int32 // 战后vn
}

type LevelsEntry struct {
	sync.RWMutex

	Levels map[int32]Level
}

func NewLevelsEntry() *LevelsEntry {
	return &LevelsEntry{
		Levels: map[int32]Level{},
	}
}

func (le *LevelsEntry) Check(config *Config) error {
	return nil
}

func (le *LevelsEntry) Reload(config *Config) error {
	le.Lock()
	defer le.Unlock()

	for _, levelCSV := range config.CfgExploreChapterLevelConfig.GetAllData() {
		level := &Level{
			AchievementsCount: 0,
			TargetCount:       0,
		}

		err := transfer.Transfer(levelCSV, level)
		if err != nil {
			return errors.WrapTrace(err)
		}

		rewards := common.NewRewards()
		rewards.AddReward(common.NewReward(static.CommonResourceTypeEnergy, levelCSV.PowerPay))
		level.Cost = rewards

		level.AchievementsDrops = make([][]int32, 0)
		if levelCSV.Achievement1ID > 0 {
			level.AchievementsDrops = append(level.AchievementsDrops, levelCSV.Achievement1DropId)
			level.AchievementsCount = level.AchievementsCount + 1
		}
		if levelCSV.Achievement2ID > 0 {
			level.AchievementsDrops = append(level.AchievementsDrops, levelCSV.Achievement2DropId)
			level.AchievementsCount = level.AchievementsCount + 1
		}

		if levelCSV.Target1ID > 0 {
			level.TargetCount = level.TargetCount + 1
		}

		if levelCSV.Target2ID > 0 {
			level.TargetCount = level.TargetCount + 1
		}

		if levelCSV.Target3ID > 0 {
			level.TargetCount = level.TargetCount + 1
		}

		le.Levels[level.Id] = *level
	}

	return nil
}

func (le *LevelsEntry) GetLevel(id int32) (*Level, error) {
	levelCfg, ok := le.Levels[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrLevelConfigNotFound, id)
	}

	return &levelCfg, nil
}
