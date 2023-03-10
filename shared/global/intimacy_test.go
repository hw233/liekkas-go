package global

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
)

func Test_intimacy(t *testing.T) {

	intimacyGlobal := NewIntimacy(redis.NewClient(&redis.Options{
		Username: "root",
		Password: "",
		Addr:     ":6379",
	}))
	ctx := context.Background()
	var guild, userId1, userId2 int64 = 1, 1, 2
	intimacy, err := intimacyGlobal.ChangeIntimacy(ctx, guild, userId1, userId2, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(intimacy)
	intimacy, err = intimacyGlobal.ChangeIntimacy(ctx, guild, userId2, userId1, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(intimacy)

	intimacy, err = intimacyGlobal.GetIntimacy(ctx, guild, userId2, userId1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(intimacy)

	err = intimacyGlobal.ClearIntimacy(ctx, guild, userId1, userId2)
	if err != nil {
		t.Fatal(err)
	}

	intimacy, err = intimacyGlobal.GetIntimacy(ctx, guild, userId2, userId1)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(intimacy)
	err = intimacyGlobal.ClearGuildIntimacy(ctx, guild)
	if err != nil {
		t.Fatal(err)
	}
	intimacy, err = intimacyGlobal.GetIntimacy(ctx, guild, userId2, userId1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(intimacy)
}

func Test_intimacy2(t *testing.T) {

	client := redis.NewClient(&redis.Options{
		Username: "root",
		Password: "",
		Addr:     ":6379",
	})
	ctx := context.Background()
	client.HSet(ctx, "test", "1", 2)
}

func TestIntimacy_GetGuildIntimacyMap(t *testing.T) {
	intimacyGlobal := NewIntimacy(redis.NewClient(&redis.Options{
		Username: "root",
		Password: "",
		Addr:     ":6379",
	}))
	ctx := context.Background()

	intimacyMap, err := intimacyGlobal.GetGuildIntimacyMap(ctx, 1, 2)
	t.Log(intimacyMap)
	t.Log(err)
}
