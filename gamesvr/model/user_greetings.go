package model

import (
	"context"
	"gamesvr/manager"
	"shared/protobuf/pb"
	"shared/statistic/logreason"
	"shared/utility/errors"
	"shared/utility/servertime"
)

func (u *User) GetGreetings(ctx context.Context, guildGreetings []*pb.VOGreetings) ([]*pb.VOGreetings, error) {
	sequence := make([]*Greeting, 0, len(guildGreetings))

	for _, guildGreeting := range guildGreetings {
		if guildGreeting.Timestamp <= servertime.Now().Unix() {
			continue
		}
		key := NewGreetingKey(guildGreeting.Id, guildGreeting.GreetingType)
		greeting := NewGreeting(guildGreeting.Star, guildGreeting.Count, guildGreeting.Timestamp, *key)
		sequence = append(sequence, greeting)
	}

	u.Guild.GuildGreetings.Sequence = sequence

	return u.Guild.GuildGreetings.VOGreetings(), nil
}

func (u *User) GreetingsRewards(ctx context.Context) (int64, error) {

	var lastTimestamp int64

	// fmt.Println("==============length", len(u.Guild.GuildGreetings.Sequence))

	for _, greeting := range u.Guild.GuildGreetings.Sequence {
		if lastTimestamp < greeting.Timestamp {
			lastTimestamp = greeting.Timestamp
		}
		if greeting.Timestamp <= servertime.Now().Unix() {
			continue
		}

		// fmt.Println("============greeting: ", *greeting)
		// 根据星级读取csv的奖励
		dropId, err := manager.CSV.GreetingsEntry.GetDropId(greeting.Star)
		if err != nil {
			return 0, errors.WrapTrace(err)
		}

		err = manager.Global.Greetings.SetUserGreetingCount(ctx, u.GetUserId(), greeting.Id, greeting.GreetingType, 1)
		if err != nil {
			return 0, errors.WrapTrace(err)
		}

		reason := logreason.NewReason(logreason.GuildGreating)
		_, err = u.AddRewardsByDropId(dropId, reason)
		if err != nil {
			return 0, errors.WrapTrace(err)
		}
		// 领取过奖励之后，将其时间戳置0，就不会出现重复领取的情况了
		greeting.Timestamp = 0
	}

	return lastTimestamp, nil
}

func (u *User) AddGreetings(gid int32, gType int32, star int32) {
	key := NewGreetingKey(gid, gType)
	greetings := NewGreeting(star, 0, 0, *key)
	u.Guild.GuildGreetings.AddSendGreetings(greetings)
}

func (u *User) ClearGreetings() {
	u.Guild.GuildGreetings.SendQueue = []*Greeting{}
}
