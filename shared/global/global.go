package global

import (
	"database/sql"

	"shared/utility/global"

	"github.com/go-redis/redis/v8"
)

const (
	// key
	KeyUserCache  = "user:cache"  // string
	KeyUserStatus = "user:status" // string
	KeyUserID     = "user:id"     // string

	KeyGuildUser = "guild:user" // hash 存储玩家的公会数据
	KeyGuildName = "guild:name" // hash 通过公会名称查询公会id
	KeyGuildID   = "guild:id"   // string 生成公会id
	KeyGuildSet  = "guild:set"  // set  人数未满的公会
)

type Global struct {
	*global.Global
	*Guild
	*User
	*Intimacy
	*Yggdrasil
	*Greetings
	*UserMercenary
	*WhiteList
	*Announcement
	*ServerSwitch
}

func NewGlobal(db *sql.DB, client *redis.Client) *Global {
	return &Global{
		Global:        global.NewGlobal(client),
		Guild:         NewGuild(client),
		User:          NewUser(db, client),
		Intimacy:      NewIntimacy(client),
		Yggdrasil:     NewYggdrasil(client),
		Greetings:     NewGreetings(client),
		UserMercenary: NewUserMercenary(client),
		WhiteList:     NewWhiteList(client),
		Announcement:  NewAnnouncement(client),
		ServerSwitch:  NewServerSwitch(client),
	}
}
