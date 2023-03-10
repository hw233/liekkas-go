package model

import (
	"context"
	"gamesvr/manager"
	"shared/common"
	"shared/csv/entry"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/statistic/logreason"
	"shared/utility/errors"
	"shared/utility/strings"
)

const VisitingCardPosCount = 4

// DailyGiftDailyRefresh 每日刷新
func (u *User) DailyGiftDailyRefresh(refreshTime int64) {
	// 日常奖励领取情况
	u.Info.DailyGiftInfoRewarded = []bool{false, false, false}
}
func (u *User) SetKanban(characterId int32) error {
	_, err := u.CharacterPack.Get(characterId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	u.Info.Kanban = characterId

	return nil
}

func (u *User) SetAvatar(avatar int32) error {
	if !u.ItemPack.Enough(avatar, 1) {
		return errors.WrapTrace(common.ErrParamError)
	}
	item, ok := manager.CSV.Item.GetItem(avatar)
	if !ok {
		return errors.WrapTrace(errors.Swrapf(common.ErrNotFoundInCSV, entry.CfgItemData, avatar))
	}
	if item.ItemType != static.ItemTypeCharacterAvatar && item.ItemType != static.ItemTypeCharacterAvatarAuto {
		return errors.WrapTrace(common.ErrParamError)

	}
	u.Info.Avatar = avatar
	return nil
}

func (u *User) ChangeNickname(ctx context.Context, nickname string) error {

	if len(nickname) == 0 {
		return errors.WrapTrace(common.ErrParamError)
	}
	// 长度判断
	if strings.StringDisplayLen(nickname) > 8 {
		return errors.WrapTrace(common.ErrParamError)
	}
	if nickname == u.Name {
		return errors.WrapTrace(common.ErrUserNicknameNotChange)
	}
	//todo:sdk敏感词检测

	// 第一次改名不消耗
	if !u.Info.IsFirstRename {
		// 消耗

		rewards := manager.CSV.TeamLevelCache.GetChangeNameConsume()

		err := u.CheckRewardsEnough(rewards)
		if err != nil {
			return errors.WrapTrace(err)
		}
		reason := logreason.NewReason(logreason.ChangeNickName)
		err = u.CostRewards(rewards, reason)
		if err != nil {
			return errors.WrapTrace(err)
		}
	}

	return u.changeNickname(ctx, nickname)
}

func (u *User) changeNickname(ctx context.Context, nickname string) error {

	index, err := manager.Global.FetchNicknameRepeatNum(ctx, nickname)
	if err != nil {
		return errors.WrapTrace(err)
	}
	if index > 9999 {
		return errors.WrapTrace(common.ErrUserNicknameRepeatedToMany)
	}

	index, err = manager.Global.IncrNicknameRepeatNum(ctx, nickname)
	if err != nil {
		return errors.WrapTrace(err)
	}

	u.Name = nickname
	u.Info.NameIndex = index
	u.Info.IsFirstRename = false
	return nil
}

func (u *User) ChangeSignature(signature string) error {
	// 长度判断
	if strings.StringDisplayLen(signature) > 50 {
		return errors.WrapTrace(common.ErrParamError)
	}
	//todo:sdk敏感词检测
	u.Info.Signature = signature
	return nil

}

func (u *User) VisitingCardSetCharacter(list []*pb.VOCharacterPos) error {
	if len(list) > VisitingCardPosCount {
		return errors.WrapTrace(common.ErrParamError)
	}
	//参数检查
	characterMap := map[int32]int32{}
	for _, pos := range list {
		// 同pos重复设置
		_, ok := characterMap[pos.Position]
		if ok {
			return errors.WrapTrace(common.ErrParamError)
		}
		// 同chara重复设置
		for _, charaId := range characterMap {
			if charaId == pos.CharacterId {
				return errors.WrapTrace(common.ErrParamError)
			}
		}

		if pos.Position > VisitingCardPosCount || pos.Position < 1 {
			return errors.WrapTrace(common.ErrParamError)
		}
		_, err := u.CharacterPack.Get(pos.CharacterId)
		if err != nil {
			return errors.WrapTrace(err)
		}
		characterMap[pos.Position] = pos.CharacterId
	}

	u.Info.CardShow.Characters = characterMap
	u.Info.CardShow.CharactersSet = true
	return nil
}

func (u *User) VisitingCardSetWorldItem(worldItemUid int64) error {
	if worldItemUid == 0 {
		u.Info.CardShow.WorldItemId = 0
		u.Info.CardShow.WorldItemUId = 0
		return nil
	}
	worldItem, err := u.WorldItemPack.Get(worldItemUid)
	if err != nil {
		return errors.WrapTrace(err)
	}
	u.Info.CardShow.WorldItemId = worldItem.WID
	u.Info.CardShow.WorldItemUId = worldItemUid
	return nil
}

func (u *User) GetDailyGiftReward(index int) error {
	dailyReward, err := manager.CSV.UserSetting.GetDailyRewardByIndex(index)
	if err != nil {
		return errors.WrapTrace(err)
	}
	if u.Info.DailyGiftInfoRewarded[index] {
		return errors.WrapTrace(common.ErrParamError)
	}

	reason := logreason.NewReason(logreason.DailyReward)
	_, err = u.AddRewardsByDropId(dailyReward.DropID, reason)
	if err != nil {
		return errors.WrapTrace(err)
	}

	u.Info.DailyGiftInfoRewarded[index] = true
	return nil
}

func (u *User) SetFrame(frame int32) error {
	if !u.ItemPack.Enough(frame, 1) {
		return errors.WrapTrace(common.ErrParamError)
	}
	item, ok := manager.CSV.Item.GetItem(frame)
	if !ok {
		return errors.WrapTrace(errors.Swrapf(common.ErrNotFoundInCSV, entry.CfgItemData, frame))
	}
	if item.ItemType != static.ItemTypeCharacterFrame {
		return errors.WrapTrace(common.ErrParamError)

	}
	u.Info.Frame = frame
	return nil
}
