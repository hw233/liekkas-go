package model

import (
	"context"
	"guild/manager"
	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/servertime"
	"strconv"
)

const (
	GreetingTypeCharacter = 1
	GreetingTypeWorldItem = 2
)

const (
	GreetingsMailTemplateId = 8
)

type GuildGreetings struct {
	Greetings map[int64]*MemberRecord `json:"greetings"`
}

type MemberRecord struct {
	MemberGreetings []*Greeting `json:"member_greetings"`
}

type Greeting struct {
	*GreetingKey
	Star      int32 `json:"star"`
	Timestamp int64 `json:"timestamp"`
}

type GreetingKey struct {
	Id           int32 `json:"id"`
	GreetingType int32 `json:"greeting_type"`
}

func NewGreetings() *GuildGreetings {
	return &GuildGreetings{
		Greetings: map[int64]*MemberRecord{},
	}
}

func NewGreetingKey(id, gType int32) *GreetingKey {
	return &GreetingKey{
		Id:           id,
		GreetingType: gType,
	}
}

func NewGreeting(star int32, timestamp int64, key GreetingKey) *Greeting {
	return &Greeting{
		GreetingKey: &key,
		Star:        star,
		Timestamp:   timestamp,
	}
}

func (g *Greeting) VOGuildGreeting(count int32) *pb.VOGreetings {
	return &pb.VOGreetings{
		GreetingType: g.GreetingType,
		Id:           g.Id,
		Star:         g.Star,
		Count:        count,
		Timestamp:    g.Timestamp,
	}
}

func NewMemBerRecord() *MemberRecord {
	return &MemberRecord{
		MemberGreetings: []*Greeting{},
	}
}

// 删除已经领取的问候奖励
func (mr *MemberRecord) UpdateGreetings(lastTimestamp int64) {

	var lastGreetingIndex int
	// 从后往前找到第一个时间戳比lastTimestamp小的元素，这个元素以及之前的元素都需要被删除
	for i := len(mr.MemberGreetings) - 1; i >= 0; i-- {
		if mr.MemberGreetings[i].Timestamp <= lastTimestamp {
			lastGreetingIndex = i + 1
			break
		}
	}

	// fmt.Println("lastGreetingIndex: ", lastGreetingIndex)
	// fmt.Println("membergreetings len: ", len(mr.MemberGreetings))

	newGreetings := make([]*Greeting, len(mr.MemberGreetings)-lastGreetingIndex)
	// fmt.Println("newGreetings len: ", len(mr.MemberGreetings)-lastGreetingIndex)
	copy(newGreetings, mr.MemberGreetings[lastGreetingIndex:])
	mr.MemberGreetings = newGreetings
}

func (g *Guild) GetMembers() []int64 {
	members := make([]int64, 0, len(g.Members))
	for _, m := range g.Members {
		members = append(members, m.UserID)
	}
	return members
}

// -----------------guild_Greetings

func (g *Guild) AddGreetings(ctx context.Context, userId int64, voGreetings []*pb.VOGreetings) error {
	// 判断玩家是否属于此公会
	err := g.CheckIsMember(userId)
	if err != nil {
		return errors.WrapTrace(err)
	}

	for _, greetings := range voGreetings {
		key := NewGreetingKey(greetings.Id, greetings.GreetingType)
		greeting := NewGreeting(greetings.Star, servertime.Now().Unix()+servertime.SecondPerDay*int64(manager.CSV.GlobalEntry.GreetingsAliveDuration), *key)

		// 获取当前公会所有成员，分别加入新的问候
		for _, userId := range g.GetMembers() {
			member, ok := g.GuildGreetings.Greetings[userId]
			if !ok {
				member = NewMemBerRecord()
			}
			member.MemberGreetings = append(member.MemberGreetings, greeting)
			g.GuildGreetings.Greetings[userId] = member
		}
	}

	return nil
}

func (g *Guild) MemberGreetings(ctx context.Context, userId int64) ([]*pb.VOGreetings, error) {
	// 判断玩家是否属于此公会
	err := g.CheckIsMember(userId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	member, ok := g.GuildGreetings.Greetings[userId]
	if !ok {
		// 如果是公会新成员，就新建问候记录
		member = NewMemBerRecord()
		g.GuildGreetings.Greetings[userId] = member
	}

	greetingsCount := map[GreetingKey]int32{}
	voGuildGreetings := make([]*pb.VOGreetings, 0, len(member.MemberGreetings))

	// 遍历目前所有收到的问候，将不超过限定数量的返回给user
	for _, greeting := range member.MemberGreetings {

		if greeting.Timestamp <= servertime.Now().Unix() {
			continue
		}

		key := NewGreetingKey(greeting.Id, greeting.GreetingType)

		record, err := manager.Global.Greetings.GetUserGreetingCount(ctx, userId, greeting.Id, greeting.GreetingType)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}

		// for test: var record int32 = 1

		count, ok := greetingsCount[*key]
		if !ok {
			count = 0
		}

		// 超过上限的直接跳过
		if greeting.GreetingType == GreetingTypeCharacter && record+count+1 > manager.CSV.GlobalEntry.CharacterGreetingsLimit {
			continue
		} else if record+count+1 > manager.CSV.GlobalEntry.WorldItemGreetingsLimit {
			continue
		}

		count++
		greeting := NewGreeting(greeting.Star, greeting.Timestamp, *key)
		greetingsCount[*key] = count

		voGuildGreetings = append(voGuildGreetings, greeting.VOGuildGreeting(count))
	}

	return voGuildGreetings, nil
}

// 领取奖励之后更新
func (g *Guild) UpdateMember(ctx context.Context, userId int64, lastTimestamp int64) error {
	// 判断玩家是否属于此公会
	err := g.CheckIsMember(userId)
	if err != nil {
		return errors.WrapTrace(err)
	}

	member, ok := g.GuildGreetings.Greetings[userId]
	if !ok {
		return errors.Swrapf(common.ErrGuildNotMember, g.ID, userId)
	}

	member.UpdateGreetings(lastTimestamp)

	g.GuildGreetings.Greetings[userId] = member

	return nil
}

func (g *Guild) GuildGreetingsLeave(ctx context.Context, userIds []int64) error {
	for _, userId := range userIds {
		delete(g.GuildGreetings.Greetings, userId)
	}
	return nil
}

// 玩家离开公会的时候，无论主动或者被动，调用此接口通过邮件发送奖励
func (g *Guild) SendGreetingsByMails(ctx context.Context, userIds []int64) error {

	for _, userId := range userIds {
		err := g.SendGreetingsToMember(ctx, userId)
		if err != nil {
			return errors.WrapTrace(err)
		}
	}
	return nil
}

func (g *Guild) SendGreetingsToMember(ctx context.Context, userId int64) error {

	// 判断玩家是否属于此公会
	err := g.CheckIsMember(userId)
	if err != nil {
		return errors.WrapTrace(err)
	}

	member, ok := g.GuildGreetings.Greetings[userId]
	if !ok {
		member = NewMemBerRecord()
		g.GuildGreetings.Greetings[userId] = member
	}

	greetingsCount := map[GreetingKey]int32{}
	allRewards := common.NewRewards()

	// 遍历目前所有收到的问候，去除超时的以及超过上限的，将其他的都合并奖励
	for _, greeting := range member.MemberGreetings {
		if greeting.Timestamp <= servertime.Now().Unix() {
			continue
		}

		key := NewGreetingKey(greeting.Id, greeting.GreetingType)
		record, err := manager.Global.Greetings.GetUserGreetingCount(ctx, userId, greeting.Id, greeting.GreetingType)
		if err != nil {
			return errors.WrapTrace(err)
		}
		count, ok := greetingsCount[*key]
		if !ok {
			count = 0
		}

		// 超过上限的直接跳过
		if greeting.GreetingType == GreetingTypeCharacter && record+count+1 > manager.CSV.GlobalEntry.CharacterGreetingsLimit {
			continue
		} else if record+count+1 > manager.CSV.GlobalEntry.WorldItemGreetingsLimit {
			continue
		}

		dropId, err := manager.CSV.GreetingsEntry.GetDropId(greeting.Star)
		if err != nil {
			return errors.WrapTrace(err)
		}
		rewards, err := manager.CSV.Drop.DropRewards(dropId)
		if err != nil {
			return errors.WrapTrace(err)
		}
		allRewards.AddRewards(rewards)
		count++
		greetingsCount[*key] = count

	}
	// 统计一共多少个问候奖励
	var allCount int32
	for _, count := range greetingsCount {
		allCount += count
	}
	if allCount == 0 {
		return nil
	}
	countStr := strconv.Itoa(int(allCount))

	req := &pb.SendPersonalMailReq{
		TemplateId:  GreetingsMailTemplateId,
		Title:       "",
		TitleArgs:   []string{},
		Content:     "",
		ContentArgs: []string{g.Name, countStr},
		Attachment:  allRewards.MergeVOResource(),
		Sender:      "",
		StartTime:   servertime.Now().Unix(),
		EndTime:     0,
		ExpireTime:  0,
		Users:       []int64{userId},
	}

	// 不需要附上id
	_, err = manager.RPCMailClient.SendPersonalMail(ctx, req)
	if err != nil {
		return errors.WrapTrace(err)
	}

	for key, count := range greetingsCount {
		err := manager.Global.Greetings.SetUserGreetingCount(ctx, userId, key.Id, key.GreetingType, count)
		if err != nil {
			return errors.WrapTrace(err)
		}
	}

	return nil
}
