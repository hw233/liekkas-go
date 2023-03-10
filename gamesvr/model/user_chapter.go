package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/statistic/logreason"
	"shared/utility/errors"
)

func (u *User) ReceiveChapterReward(rewardId int32) error {
	chapterRewardCfg, err := manager.CSV.ChapterEntry.GetChapterReward(rewardId)
	if err != nil {
		return err
	}

	chapterId := chapterRewardCfg.ChapterId
	chapter, ok := u.ChapterInfo.GetChapter(chapterId)
	if !ok {
		return errors.Swrapf(common.ErrChapterScoreNotArrival, chapterId)
	}

	if chapter.GetScore() < chapterRewardCfg.Number {
		return errors.Swrapf(common.ErrChapterScoreNotArrival, chapterId)
	}

	if chapter.IsReceiveReward(rewardId) {
		return errors.Swrapf(common.ErrChapterRewardReceived, chapterId, rewardId)
	}

	chapter.RecordReward(rewardId)

	reason := logreason.NewReason(logreason.ChapterReward)
	u.AddRewardsByDropId(chapterRewardCfg.DropId, reason)

	return nil
}

func (u *User) CheckChapterUnlock(chapterId int32) error {
	chapterCfg, err := manager.CSV.ChapterEntry.GetChapter(chapterId)
	if err != nil {
		return err
	}

	err = u.CheckUserConditions(chapterCfg.UnlockCondition)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) CheckChapterRewardNoticeByType(chapterType int32) int32 {
	count := int32(0)

	chapterIds := manager.CSV.ChapterEntry.GetChapterIdsByType(chapterType)
	for _, chapterId := range chapterIds {
		chapter, ok := u.ChapterInfo.GetChapter(chapterId)
		if !ok {
			continue
		}

		rewardIds := manager.CSV.ChapterEntry.GetChapterRewardIds(chapterId)
		for _, rewardId := range rewardIds {
			if chapter.IsReceiveReward(rewardId) {
				continue
			}

			rewardCfg, err := manager.CSV.ChapterEntry.GetChapterReward(rewardId)
			if err != nil {
				continue
			}

			if chapter.GetScore() >= rewardCfg.Number {
				count = count + 1
			}
		}
	}

	return count
}
