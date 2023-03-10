package model

import (
	"context"
	"gamesvr/manager"
	"shared/common"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/glog"
	"shared/utility/servertime"
	"time"
)

func (u *User) MercenaryDailyRefresh(refreshTime int64) {
	u.Mercenary.DailyRefresh()
}

// 佣兵已经使用次数+1
func (u *User) AddMercenaryUseCount(characterId int32, systemType int32) error {
	if systemType != static.BattleTypeChallengeAltar && systemType != static.BattleTypeTower {
		return errors.Swrapf(common.ErrWrongSystemType, systemType)
	}
	var has bool
	for _, mercenary := range u.Mercenary.Mercenaries {
		if mercenary.Character.ID == characterId {
			has = true
			count, ok := mercenary.UseCount[systemType]
			if !ok {
				count = 0
			}
			if count >= manager.CSV.GlobalEntry.MercenaryUseLimit {
				return errors.Swrapf(common.ErrMercenaryExceedUseLimit, characterId, systemType, count)
			}
			count += 1
			mercenary.UseCount[systemType] = count
		}
	}
	if !has {
		return errors.Swrapf(common.ErrMercenaryNotFound, characterId)
	}

	return nil
}

// 获取佣兵数据
func (u *User) GetMercenary(characterId int32) (*Mercenary, error) {

	var has bool
	var result *Mercenary
	for _, mercenary := range u.Mercenary.Mercenaries {
		if mercenary.Character.ID == characterId {
			has = true
			result = mercenary
		}
	}
	if !has {
		return nil, errors.Swrapf(common.ErrMercenaryNotFound, characterId)
	}

	return result, nil
}

func (u *User) UpdateMercenaryBorrow(ctx context.Context) {
	characters := manager.Global.CharacterListGetAndClear(ctx, u.GetUserId())
	for _, character := range characters {
		if len(u.Mercenary.Mercenaries) < 3 {
			u.MercenaryAdd(character)

		} else {
			glog.Errorf("MercenaryCharacterList MercenaryAdd exceed limit")
			continue
		}
	}
}

func (u *User) MercenaryCheckUserCount(systemType, characterId int32) error {
	var has bool
	for _, mercenary := range u.Mercenary.Mercenaries {
		if mercenary.Character.ID == characterId {
			has = true
			useCount, ok := mercenary.UseCount[systemType]
			if !ok {
				useCount = 0
			}
			if useCount >= manager.CSV.GlobalEntry.MercenaryUseLimit {
				return errors.Swrapf(common.ErrMercenaryExceedUseLimit, characterId, systemType, useCount)
			}
		}
	}
	if !has {
		return errors.Swrapf(common.ErrMercenaryNotFound, characterId)
	}
	return nil
}

func (u *User) GetMercenaryList(ctx context.Context, friends []int64, guildElites []int64) (*pb.S2CGetMercenaryList, error) {
	// 获取佣兵列表前首先刷新一次，get最新的已借到的佣兵缓存
	u.UpdateMercenaryBorrow(ctx)
	// 可以借的佣兵
	mercenaryUsers := map[int64]*MercenaryUser{}
	// 根据id去缓存拿name
	caches, err := manager.Global.GetUserCaches(ctx, guildElites)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	for _, cache := range caches {
		if cache.ID == u.GetUserId() {
			continue
		}
		mercenaryUser := NewMercenaryUser(cache.ID, cache.Name, MercenaryRelationGuild)

		values, err := manager.Global.GetMercenaryCharacter(ctx, cache.ID)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		for _, v := range values {
			// 如果已经借到这个佣兵，就不需要显示了
			// if u.Mercenary.CheckMercenaryOwn(v.Id) {
			// 	continue
			// }
			// character := MercenaryCharacter{}

			// 查看角色是否已经借出
			isAvaliable, err := manager.Global.CheckMercenaryAvaliable(ctx, cache.ID, v.Id)
			if err != nil {
				return nil, errors.WrapTrace(err)
			}
			if !isAvaliable {
				continue
			}
			// 获取这个角色的申请人的人数
			count, err := manager.Global.MercenaryApplyGet(ctx, cache.ID, v.Id)
			if err != nil {
				return nil, errors.WrapTrace(err)
			}
			character := NewMercenaryCharacter(v.Id, v.Level, v.Power, count, v.Star)
			mercenaryUser.Characters[v.Id] = character
		}
		if len(mercenaryUser.Characters) <= 0 {
			continue
		}

		mercenaryUsers[mercenaryUser.Uid] = mercenaryUser
	}

	// 遍历处于正在申请和撤销申请状态的申请，改变状态
	for _, sent := range u.Mercenary.MercenarySentApply {
		userMercenary, ok := mercenaryUsers[sent.SentTo]
		if !ok {
			sent.IfDelete = true
			continue
		}
		character, ok := userMercenary.Characters[sent.CharacterId]
		if !ok {
			sent.IfDelete = true
			continue
		}
		// 如果角色已经被借到，这个申请也要被删除
		if u.Mercenary.CheckMercenaryOwn(sent.CharacterId) {
			sent.IfDelete = true
		}
		if sent.IsCanceled {
			if servertime.Now().Unix() > sent.ExpireTime {
				sent.IfDelete = true
				// delete(userMercenary.Characters, sent.CharacterId)
				// continue
			} else {
				character.Status = MercenaryStatusCancelApplied // 正在撤销
				character.ExpireTime = sent.ExpireTime
			}

		} else {
			character.Status = MercenaryStatusApplied // 正在申请
		}
		mercenaryUsers[sent.SentTo].Characters[sent.CharacterId] = character
	}

	// 删除已经发送的申请中不合规的
	u.Mercenary.UpdateMercenarySend()

	voMercenaries := make([]*pb.VOMercenaries, 0, len(mercenaryUsers))
	for _, m := range mercenaryUsers {
		voMercenaries = append(voMercenaries, m.VOMercenaries())
	}

	info := make([]*pb.VOMercenarySimple, 0, 3)
	for _, character := range u.Mercenary.Mercenaries {
		// fmt.Println("===========mercenary: ", character.ID)
		if character.Character == nil {
			glog.Errorf("GetMercenaryList mercenary no data , characterId: %d", character.ID)
			continue
		}
		// fmt.Println("===========mercenary:data  ")
		c := &pb.VOMercenarySimple{
			CharacterId: character.Character.ID,
			NameFrom:    character.Owner,
			Level:       character.Character.Level,
			Relation:    int32(character.Relation),
		}
		info = append(info, c)
	}

	return &pb.S2CGetMercenaryList{
		Info:           info,
		ExpireTime:     u.Mercenary.ExpireTime,
		UserCharacters: voMercenaries,
	}, nil
}

// 返回值为0，代表发送成功；返回值为1，代表此玩家已经并非你的朋友或者公会精英；返回值为2，代表该角色已经借出
func (u *User) SendApply(ctx context.Context, userId int64, characterId int32, friends, elites []int64) (int32, error) {
	u.UpdateMercenaryBorrow(ctx)
	// todo 输出错误
	if userId == u.GetUserId() {
		return 0, nil
	}
	if u.Mercenary.CheckMercenaryOwn(characterId) {
		return 0, errors.Swrapf(common.ErrMercenaryAlreadyHad, characterId)
	}
	if u.Mercenary.SendApplyCountWithoutCancel() >= 5 {
		return 0, errors.Swrapf(common.ErrMercenaryApplyExceedLimit)
	}
	// 判断关系
	var isFriend, isElite bool
	for _, friend := range friends {
		if friend == userId {
			isFriend = true
		}
	}
	for _, elite := range elites {
		if elite == userId {
			isElite = true
		}
	}
	if !isElite && !isFriend {
		return 1, nil
	}

	// redis 查看是否已经借出
	IsAvailable, err := manager.Global.UserMercenary.CheckMercenaryAvaliable(ctx, userId, characterId)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}

	if !IsAvailable {
		return 2, nil
	}

	// 添加已经发送申请到本地
	err = u.Mercenary.AddSendApply(userId, characterId)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}

	sendTime := servertime.Now().Unix()
	nextWeekRefresh := ThisWeekRefreshTime().Add(time.Hour * 24 * 7)
	sendApply := common.NewMercenarySend(u.GetUserId(), u.GetUserName(), characterId)
	sendApply.SendTime = sendTime

	// 发送申请
	err = manager.Global.MercenaryApplySend(ctx, userId, sendApply, nextWeekRefresh)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}

	// 添加已发送申请到缓存
	// record := common.NewMercenarySendRecord(userId, characterId)
	// err = manager.Global.UserMercenary.MercenarySendRecordAdd(ctx, u.GetUserId(), record, nextWeekRefresh)
	// if err != nil {
	// 	return 0, errors.WrapTrace(err)
	// }

	// 申请人数+1
	err = manager.Global.UserMercenary.MercenaryApplyAdd(ctx, userId, characterId, 1, nextWeekRefresh)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}

	return 0, nil
}

func (u *User) CancelApply(ctx context.Context, userId int64, characterId int32) (int64, int32, error) {
	u.UpdateMercenaryBorrow(ctx)

	if u.Mercenary.CheckMercenaryOwn(characterId) {
		return 0, 1, nil
	}

	var ok bool
	var coolDown int64
	for _, sent := range u.Mercenary.MercenarySentApply {
		if sent.SentTo == userId && sent.CharacterId == characterId {
			ok = true
			sent.IsCanceled = true
			sent.ExpireTime = servertime.Now().Add(10 * time.Minute).Unix()
			coolDown = sent.ExpireTime

			sendApply := common.NewMercenarySend(u.GetUserId(), u.GetUserName(), characterId)
			sendApply.IsCancel = true
			nextWeekRefresh := ThisWeekRefreshTime().Add(time.Hour * 24 * 7)
			// 发送取消申请
			err := manager.Global.MercenaryApplySend(ctx, userId, sendApply, nextWeekRefresh)
			if err != nil {
				return 0, 0, errors.WrapTrace(err)
			}

			// record := common.NewMercenarySendRecord(userId, characterId)
			// err = manager.Global.UserMercenary.MercenarySendRecordDelete(ctx, u.GetUserId(), record)
			// if err != nil {
			// 	return errors.WrapTrace(err)
			// }

			// 申请人数-1
			err = manager.Global.UserMercenary.MercenaryApplyAdd(ctx, userId, characterId, -1, nextWeekRefresh)
			if err != nil {
				return 0, 0, errors.WrapTrace(err)
			}
			break
		}
	}
	if !ok {
		return 0, 0, errors.Swrapf(common.ErrMercenarySendApplyNotFound, characterId, userId)
	}
	return coolDown, 0, nil
}

// 去缓存里面拿最新的申请数据
func (u *User) UpdateMercenaryReceivedApply(ctx context.Context, friends, members []int64) error {
	message, err := manager.Global.MercenaryApplyReceive(ctx, u.GetUserId())
	if err != nil {
		return errors.WrapTrace(err)
	}

	friendsSet := make(map[int64]struct{})
	membersSet := make(map[int64]struct{})
	for _, f := range friends {
		friendsSet[f] = struct{}{}
	}
	for _, e := range members {
		membersSet[e] = struct{}{}
	}

	for _, m := range message {
		if m.IsCancel {
			// 是有序的
			u.Mercenary.DeleteSpecificApply(m.CharacterId, m.Uid)
			continue
		}
		var isFriend, isMembers int32
		_, ok := friendsSet[m.Uid]
		if ok {
			isFriend = 1
		}
		_, ok = membersSet[m.Uid]
		if ok {
			isMembers = 2
		}
		relation := isFriend + isMembers
		if relation == 0 { // 如果此时非公会或者好友关系，则忽略此申请
			continue
		}
		apply := NewMercenaryApply(m.CharacterId, relation, m.Uid, m.Name, m.SendTime)
		subMap, ok := u.Mercenary.MercenaryReceivedApply[m.CharacterId]
		if !ok {
			subMap = make(map[int64]*MercenaryApply)
		}
		subMap[m.Uid] = apply
		u.Mercenary.MercenaryReceivedApply[m.CharacterId] = subMap
	}
	return nil
}

// 拉去所有申请，并返回VO
func (u *User) GetMercenaryManagement(ctx context.Context, friends, members []int64) ([]*pb.VOMercenaryApply, error) {
	err := u.UpdateMercenaryReceivedApply(ctx, friends, members)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	voApply := make([]*pb.VOMercenaryApply, 0, len(u.Mercenary.MercenaryReceivedApply))

	for _, apply := range u.Mercenary.MercenaryReceivedApply {
		for _, value := range apply {
			voApply = append(voApply, value.VOMercenaryApply())
		}
	}

	return voApply, nil
}

// 处理之前拉去一下申请列表防止有更新
// 返回值为0， 代表成功，返回值为1，代表此人已经并非朋友或者公会成员; 返回值为2，代表对方已经借到别人的角色或者角色已满;返回值为3，代表对方已经撤销申请
func (u *User) HandleApply(ctx context.Context, ifIgnore bool, userId int64, characterId int32, friends, members []int64) (int32, error) {
	err := u.UpdateMercenaryReceivedApply(ctx, friends, members)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}
	if ifIgnore {
		u.Mercenary.DeleteSpecificApply(characterId, userId)
		return 0, nil
	}
	if u.Mercenary.CheckMercenaryRecord(characterId) {
		return 4, nil
	}
	apply, ok := u.Mercenary.MercenaryReceivedApply[characterId][userId]
	if !ok {
		return 3, nil
	}

	// 判断关系
	var isFriend, isMembers int32
	for _, friend := range friends {
		if friend == userId {
			isFriend = 1
		}
	}
	for _, member := range members {
		if member == userId {
			isMembers = 2
		}
	}
	relation := isMembers + isFriend
	if relation == 0 { // 已经不是公会或者朋友关系了
		u.Mercenary.DeleteSpecificApply(characterId, userId)
		return 1, nil
	}

	// 将过期时间设置为下周的刷新时间
	nextWeekRefresh := WeekRefreshTime(servertime.Now()).Add(time.Hour * 24 * 7)

	// 检查对方是否已经拥有此id的佣兵
	isSuccess, err := manager.Global.PutMercenaryOwn(ctx, userId, characterId, nextWeekRefresh)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}
	if !isSuccess {
		u.Mercenary.DeleteSpecificApply(characterId, userId)
		return 2, nil
	}

	// 将当前角色设置为已经借出
	err = manager.Global.UserMercenary.ChangeMercenaryAvaliable(ctx, u.GetUserId(), characterId, nextWeekRefresh)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}
	// 发送数据
	charac, err := u.CharacterPack.Get(characterId)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}
	equips := make([]*common.Equipment, 0, 4)
	for _, eid := range charac.Equipments {
		if eid > 0 {
			equip, err := u.EquipmentPack.Get(eid)
			if err != nil {
				return 0, errors.WrapTrace(err)
			}
			equips = append(equips, equip)
		}
	}
	worldItem := &common.WorldItem{}
	if charac.WorldItem > 0 {
		target, err := u.WorldItemPack.Get(charac.WorldItem)
		if err != nil {
			return 0, errors.WrapTrace(err)
		}
		worldItem = target
	}

	characData := common.NewCharacterData(charac.ID, charac.Level, charac.Star, charac.Stage, charac.Power, charac.Rarity, charac.Skills, equips, worldItem)
	err = manager.Global.CharacterList.CharacterListPush(ctx, u.GetUserId(), u.GetUserName(), userId, int8(relation), characData, nextWeekRefresh)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}
	record := NewMercenaryRecord(characterId, relation, userId, apply.ApplicantName, charac.Star, charac.Level, charac.Power)
	u.Mercenary.AddMercenaryRecord(record)
	// 将有关的其他申请删除
	u.Mercenary.DeleteApply(characterId)

	return 0, nil
}

func (u *User) GetMercenaryRecord(ctx context.Context) []*pb.VOMercenaryRecord {

	result := make([]*pb.VOMercenaryRecord, 0, len(u.Mercenary.MercenaryRecords))

	// 遍历所有的record, 发送id
	for _, record := range u.Mercenary.MercenaryRecords {
		result = append(result, record.VOMercenaryRecord())
	}

	return result
}

func (u *User) MercenaryAdd(c *common.MercenaryData) {
	for _, mercenary := range u.Mercenary.Mercenaries {
		if c.Uid == mercenary.ID && c.Data.ID == mercenary.Character.ID {
			glog.Errorf("MercenaryCharacterList MercenaryAdd err: characterData repeated")
			return
		}
	}
	character := NewMercenary(c.Uid, c.Name, c.Relation, c.Data)
	u.Mercenary.Mercenaries = append(u.Mercenary.Mercenaries, character)
}

// 给前端返回佣兵所有数据
func (u *User) GetMercenaryCharacter(ctx context.Context, systemType int32) ([]*pb.VOMercenaryDetail, error) {

	details, err := u.GetMercenaryCharacterData(ctx, systemType)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return details, nil
}

func (u *User) GetMercenaryCharacterData(ctx context.Context, systemType int32) ([]*pb.VOMercenaryDetail, error) {
	details := make([]*pb.VOMercenaryDetail, 0, 3)
	for _, mercenary := range u.Mercenary.Mercenaries {

		useCount, ok := mercenary.UseCount[systemType]
		if !ok {
			useCount = 0
		}
		if useCount >= manager.CSV.GlobalEntry.MercenaryUseLimit {
			continue
		}
		if mercenary == nil {
			continue
		}
		voCharacter := &pb.VOBattleCharacterProp{}

		newCharacter := NewCharacter(mercenary.Character.ID)
		newCharacter.Level = mercenary.Character.Level
		newCharacter.Stage = mercenary.Character.Stage
		newCharacter.Star = mercenary.Character.Star
		newCharacter.Skills = mercenary.Character.Skills
		newCharacter.Power = mercenary.Character.Power
		newCharacter.Rarity = mercenary.Character.Rarity

		voCharacter.Character = newCharacter.VOUserCharacter()

		for _, equip := range mercenary.Equipments {
			voCharacter.Equipments = append(voCharacter.Equipments, equip.VOUserEquipment())
		}
		if mercenary.WorldItem != nil && mercenary.WorldItem.ID > 0 {
			voCharacter.WorldItem = mercenary.WorldItem.VOUserWorldItem()
		}

		voDetail := &pb.VOMercenaryDetail{
			SystemType: systemType,
			UseCount:   useCount,
			Character:  voCharacter,
		}
		details = append(details, voDetail)
	}

	return details, nil
}
