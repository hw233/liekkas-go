package model

import (
	"shared/protobuf/pb"
)

const (
	ChapterTypeCampain = 1 //主线
	ChapterTypeChaos   = 2 //混沌
	ChapterTypeElite   = 3 //历练
)

type Chapter struct {
	Id           int32          `json:"id"`
	Score        int32          `json:"score"`
	RewardRecord map[int32]bool `json:"reward_record"`
}

type ChapterInfo struct {
	Chapters map[int32]*Chapter `json:"chapters"`
}

func NewChapter(id int32) *Chapter {
	chapter := &Chapter{
		Id:           id,
		Score:        0,
		RewardRecord: map[int32]bool{},
	}

	return chapter
}

func NewChapterInfo() *ChapterInfo {
	return &ChapterInfo{
		Chapters: map[int32]*Chapter{},
	}
}

//----------------------------------------
// ChapterInfo
//----------------------------------------
func (ci *ChapterInfo) AddChapter(id int32) *Chapter {
	chapter, ok := ci.Chapters[id]
	if ok {
		return chapter
	}

	chapter = NewChapter(id)

	ci.Chapters[id] = chapter

	return chapter
}

func (ci *ChapterInfo) GetChapter(id int32) (*Chapter, bool) {
	chapter, ok := ci.Chapters[id]
	return chapter, ok
}

func (ci *ChapterInfo) GetOrCreateChapter(id int32) *Chapter {
	chapter, ok := ci.GetChapter(id)
	if !ok {
		chapter = ci.AddChapter(id)
	}

	return chapter
}

func (c *ChapterInfo) VOChapterInfo() *pb.VOChapterInfo {
	chapters := make([]*pb.VOChapter, 0, len(c.Chapters))

	for _, chapter := range c.Chapters {
		chapters = append(chapters, chapter.VOChapter())
	}

	return &pb.VOChapterInfo{
		Chapters: chapters,
	}
}

//----------------------------------------
// Chapter
//----------------------------------------
func (c *Chapter) AddScore(score int32) {
	c.Score = c.Score + score
}

func (c *Chapter) GetScore() int32 {
	return c.Score
}

func (c *Chapter) IsReceiveReward(rewardId int32) bool {
	received, ok := c.RewardRecord[rewardId]
	if !ok {
		return false
	}

	return received
}

func (c *Chapter) RecordReward(rewardId int32) {
	c.RewardRecord[rewardId] = true
}

func (c *Chapter) VOChapter() *pb.VOChapter {
	chapterData := &pb.VOChapter{
		ChapterId:    c.Id,
		Score:        c.Score,
		RewardRecord: make([]int32, 0, len(c.RewardRecord)),
	}

	for index := range c.RewardRecord {
		chapterData.RewardRecord = append(chapterData.RewardRecord, index)
	}

	return chapterData
}
