package global

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
)

func TestGreetings(t *testing.T) {

	GreetingsGlobal := NewGreetings(redis.NewClient(&redis.Options{
		Username: "root",
		Password: "",
		Addr:     ":6379",
	}))

	ctx := context.Background()

	err := GreetingsGlobal.SetUserGreetingCount(ctx, 1, 1001, 1, 1)
	if err != nil {
		t.Fatal(err)
	}

	record, err := GreetingsGlobal.GetUserGreetingCount(ctx, 1, 1001, 1)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(record)

}
