package event

import (
	"github.com/go-redis/redis/v8"
	"log"
	"shared/utility/global"
	"shared/utility/glog"
	"shared/utility/param"
	"testing"
)

func TestEventLooper_Get(t *testing.T) {

	glog.InitLog()
	Redis := &redis.Options{
		Username: "root",
		Password: "",
		Addr:     ":6379",
	}
	RedisClient := redis.NewClient(Redis)
	eventLooper := NewEventQueue(RedisClient, global.NewGlobal(RedisClient))
	UserEventHandler.Register(1, 1, func(Param *param.Param) error {
		log.Println(Param.GetString(0))
		return nil
	})

	eventLooper.Push(1, NewEvent(1, "test1"))
	eventLooper.Push(1, NewEvent(1, "test2"))

	eventLooper.ExecuteEventsInQueue(1)

}
