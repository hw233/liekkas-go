package global

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestGetUserCache(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Username: "root",
		Password: "",
		Addr:     ":6379",
	})

	db, err := sql.Open("mysql", "root:123456@tcp(10.24.12.30:3306)/overlord_user_go_dev?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil {
		t.Errorf("open mysql error: %v", err)
		return
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	err = db.Ping()
	if err != nil {
		t.Errorf("init mysql error: %v", err)
		return
	}

	ch := NewUserCacheHandler(db, redisClient)

	c, err := ch.GetUserCaches(context.Background(), []int64{20210804101717, 20210913212912, 20211214154146})
	if err != nil {
		log.Printf("get cache error: %v", err)
		return
	}

	fmt.Printf("ret: %+v \n", c[20210804101717])
	fmt.Printf("ret: %+v \n", c[20210913212912])
	fmt.Printf("ret: %+v \n", c[20211214154146])
}
