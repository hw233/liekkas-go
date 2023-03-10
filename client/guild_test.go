package main

import (
	"fmt"
	"testing"
	"time"

	"shared/protobuf/pb"
)

const i = 1100

// 创建公会
func TestGuildCreate(t *testing.T) {
	client := Pool.NewClient()

	_, err := client.Login(&pb.C2SLogin{
		UserId: i,
	})
	if err != nil {
		return
	}

	_, err = client.GuildCreate(&pb.C2SGuildCreate{
		Name:      fmt.Sprintf("Guild%d", i),
		Icon:      &pb.VOGuildIcon{},
		JoinModel: 0,
	})
	if err != nil {
		t.Logf("GuildCreate error: %v", err)
		return
	}
}

func TestGuildInfo(t *testing.T) {
	client := Pool.NewClient()

	_, err := client.Login(&pb.C2SLogin{
		UserId: i,
	})
	if err != nil {
		return
	}

	_, err = client.GuildInfo(&pb.C2SGuildInfo{})
	if err != nil {
		return
	}
}

// 创建公会+申请+自动批准入会
func TestGuildCreateAndApply(t *testing.T) {
	client := Pool.NewClient()
	client2 := Pool.NewClient()

	_, err := client.Login(&pb.C2SLogin{
		UserId: i,
	})
	if err != nil {
		return
	}

	guildCreateRet, err := client.GuildCreate(&pb.C2SGuildCreate{
		Name:      fmt.Sprintf("Guild%d", i),
		Icon:      &pb.VOGuildIcon{},
		JoinModel: 0,
	})
	if err != nil {
		t.Logf("GuildCreate error: %v", err)
		return
	}

	time.Sleep(time.Second)

	_, err = client2.Login(&pb.C2SLogin{
		UserId: i + 1,
	})
	if err != nil {
		return
	}

	_, err = client2.GuildApply(&pb.C2SGuildApply{
		// GuildID: 16,
		GuildID: guildCreateRet.GuildInfo.GuildID,
	})
	if err != nil {
		return
	}
}

// 创建公会+申请+手动批准入会
func TestGuildCreateAndApply2(t *testing.T) {
	client := Pool.NewClient()
	client2 := Pool.NewClient()

	_, err := client.Login(&pb.C2SLogin{
		UserId: i,
	})
	if err != nil {
		return
	}

	guildCreateRet, err := client.GuildCreate(&pb.C2SGuildCreate{
		Name:      fmt.Sprintf("Guild%d", i),
		Icon:      &pb.VOGuildIcon{},
		JoinModel: 1,
	})
	if err != nil {
		t.Logf("GuildCreate error: %v", err)
		return
	}

	_, err = client.GuildInfo(&pb.C2SGuildInfo{})
	if err != nil {
		return
	}

	time.Sleep(time.Second)

	_, err = client2.Login(&pb.C2SLogin{
		UserId: i + 1,
	})
	if err != nil {
		return
	}

	_, err = client2.GuildApply(&pb.C2SGuildApply{
		// GuildID: 16,
		GuildID: guildCreateRet.GuildInfo.GuildID,
	})
	if err != nil {
		return
	}

	time.Sleep(time.Second)

	_, err = client.GuildHandleApplied(&pb.C2SGuildHandleApplied{
		Approve: []int64{i + 1},
	})
	if err != nil {
		return
	}

	_, err = client2.GuildInfo(&pb.C2SGuildInfo{})
	if err != nil {
		return
	}
}

// 创建公会+解散公会
func TestGuildCreateAndDissolve(t *testing.T) {
	client := Pool.NewClient()

	_, err := client.Login(&pb.C2SLogin{
		UserId: i,
	})
	if err != nil {
		return
	}

	_, err = client.GuildCreate(&pb.C2SGuildCreate{
		Name:      fmt.Sprintf("Guild%d", i),
		Icon:      &pb.VOGuildIcon{},
		JoinModel: 1,
	})
	if err != nil {
		return
	}

	_, err = client.GuildInfo(&pb.C2SGuildInfo{})
	if err != nil {
		return
	}

	_, err = client.GuildDissolve(&pb.C2SGuildDissolve{})
	if err != nil {
		return
	}
}

// 创建公会+申请+自动批准入会+升职
func TestGuildCreateAndApplyAndPromotion(t *testing.T) {
	client := Pool.NewClient()
	client2 := Pool.NewClient()

	_, err := client.Login(&pb.C2SLogin{
		UserId: i,
	})
	if err != nil {
		return
	}

	guildCreateRet, err := client.GuildCreate(&pb.C2SGuildCreate{
		Name:      fmt.Sprintf("Guild%d", i),
		Icon:      &pb.VOGuildIcon{},
		JoinModel: 0,
	})
	if err != nil {
		t.Logf("GuildCreate error: %v", err)
		return
	}

	time.Sleep(time.Second)

	_, err = client2.Login(&pb.C2SLogin{
		UserId: i + 1,
	})
	if err != nil {
		return
	}

	_, err = client2.GuildApply(&pb.C2SGuildApply{
		// GuildID: 16,
		GuildID: guildCreateRet.GuildInfo.GuildID,
	})
	if err != nil {
		return
	}

	_, err = client.GuildPromotion(&pb.C2SGuildPromotion{
		UserID: i + 1,
	})
	if err != nil {
		return
	}

	_, err = client.GuildPromotion(&pb.C2SGuildPromotion{
		UserID: i + 1,
	})
	if err != nil {
		return
	}
}
