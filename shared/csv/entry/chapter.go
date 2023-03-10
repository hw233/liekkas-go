package entry

import (
	"shared/common"
	"shared/utility/errors"
	"shared/utility/transfer"
	"sync"
)

type Chapter struct {
	Id              int32
	ChapterType     int32
	UnlockCondition *common.Conditions `rule:"conditions"`
	InitPos         *common.Vec2       `rule:"floatsToVec2"`
}

type ChapterReward struct {
	Id        int32
	ChapterId int32
	Number    int32
	DropId    int32
}

type ChapterEntry struct {
	sync.RWMutex

	Chapters       map[int32]*Chapter
	ChapterRewards map[int32]*ChapterReward

	TypeToChapters   map[int32][]int32
	ChapterToRewards map[int32][]int32
}

func NewChapterEntry() *ChapterEntry {
	return &ChapterEntry{
		Chapters:         map[int32]*Chapter{},
		ChapterRewards:   map[int32]*ChapterReward{},
		TypeToChapters:   map[int32][]int32{},
		ChapterToRewards: map[int32][]int32{},
	}
}

func (ce *ChapterEntry) Check(config *Config) error {
	return nil
}

func (ce *ChapterEntry) Reload(config *Config) error {
	ce.Lock()
	defer ce.Unlock()

	typeToChapters := map[int32][]int32{}
	chapterToRewards := map[int32][]int32{}

	for _, chapterCSV := range config.CfgExploreChapterConfig.GetAllData() {
		chapter := &Chapter{}

		err := transfer.Transfer(chapterCSV, chapter)
		if err != nil {
			return errors.WrapTrace(err)
		}

		ce.Chapters[chapter.Id] = chapter

		_, ok := typeToChapters[chapter.ChapterType]
		if !ok {
			typeToChapters[chapter.ChapterType] = []int32{}
		}
		typeToChapters[chapter.ChapterType] = append(typeToChapters[chapter.ChapterType], chapter.Id)
	}

	for _, chapterRewardCSV := range config.CfgExploreChapterRewardConfig.GetAllData() {
		chapterReward := &ChapterReward{}

		err := transfer.Transfer(chapterRewardCSV, chapterReward)
		if err != nil {
			return errors.WrapTrace(err)
		}

		chapterRewardId := chapterReward.Id
		chapterId := chapterReward.ChapterId
		ce.ChapterRewards[chapterRewardId] = chapterReward
		_, ok := chapterToRewards[chapterId]
		if !ok {
			chapterToRewards[chapterId] = []int32{}
		}
		chapterToRewards[chapterId] = append(chapterToRewards[chapterId], chapterRewardId)
	}

	ce.TypeToChapters = typeToChapters
	ce.ChapterToRewards = chapterToRewards

	return nil
}

func (ce *ChapterEntry) GetChapter(id int32) (*Chapter, error) {
	chapter, ok := ce.Chapters[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrChapterConfigNotFound, id)
	}

	return chapter, nil
}

func (ce *ChapterEntry) GetChapterReward(id int32) (*ChapterReward, error) {
	chapterReward, ok := ce.ChapterRewards[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrChapterRewardConfigNotFound, id)
	}

	return chapterReward, nil
}

func (ce *ChapterEntry) GetChapterIdsByType(chapterType int32) []int32 {
	chapters, ok := ce.TypeToChapters[chapterType]
	if !ok {
		chapters = []int32{}
	}

	return chapters
}

func (ce *ChapterEntry) GetChapterRewardIds(chapterId int32) []int32 {
	ids, ok := ce.ChapterToRewards[chapterId]
	if !ok {
		ids = []int32{}
	}

	return ids
}
