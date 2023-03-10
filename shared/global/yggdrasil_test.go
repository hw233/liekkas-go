package global

import (
	"context"
	"github.com/go-redis/redis/v8"
	"testing"
)

func TestYggGoodsTakenBack_FetchGoodsTakenBack(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Username: "root",
		Password: "",
		Addr:     ":6379",
	})

	Yggdrasil := NewYggdrasil(redisClient)
	back, err := Yggdrasil.FetchGoodsTakenBack(context.Background(), 1)
	t.Log(back)
	t.Log(err)

}
