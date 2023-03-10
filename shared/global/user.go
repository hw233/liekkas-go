package global

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-redis/redis/v8"

	"shared/utility/global"
	"shared/utility/key"
)

const (
	KeyNicknameIndex = "NickNameIndex" // 重复名字redisKey

	UserCacheExpire = 6 * time.Hour // 玩家缓存过期时间

	UserCacheWithGuild  = 1 << iota // 公会信息
	UserCacheWithOnline             // 在线信息
)

type User struct {
	*UserCacheHandler
	*RepeatUserNameSelector
}

func NewUser(db *sql.DB, client *redis.Client) *User {
	return &User{
		UserCacheHandler:       NewUserCacheHandler(db, client),
		RepeatUserNameSelector: NewRepeatUserNameSelector(client),
	}
}

// 取玩家缓存数据
type UserCacheHandler struct {
	handler *global.CacheHandler
	*UserStatus
	*GuildUser
	*NickNameSuffix
}

func NewUserCacheHandler(db *sql.DB, client *redis.Client) *UserCacheHandler {
	return &UserCacheHandler{
		handler: global.NewCacheHandler(client, KeyUserCache, func() global.Cache {
			return &UserCache{
				ImplementedCache: &global.ImplementedCache{
					DB:     db,
					Client: client,
				},
			}
		}, UserCacheExpire),
		UserStatus:     NewUserStatus(client),
		GuildUser:      NewGuildUser(client),
		NickNameSuffix: NewNickNameSuffix(client),
	}
}

func (h *UserCacheHandler) GetUserCache(ctx context.Context, userID int64) (*UserCache, error) {
	caches, err := h.handler.GetCaches(ctx, []int64{userID})
	if err != nil {
		return nil, err
	}

	return caches[userID].(*UserCache), err
}

func (h *UserCacheHandler) GetUserCaches(ctx context.Context, userIDs []int64) (map[int64]*UserCache, error) {
	caches, err := h.handler.GetCaches(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	ret := make(map[int64]*UserCache, len(caches))

	for k, v := range caches {
		ret[k] = v.(*UserCache)
	}

	return ret, err
}

// EX: GetUserCachesExtension(crx, userIDs, UserCacheWithGuild|UserCacheWithOnline)
func (h *UserCacheHandler) GetUserCachesExtension(ctx context.Context, userIDs []int64, extension int) (map[int64]*UserCache, error) {
	caches, err := h.handler.GetCaches(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	ret := make(map[int64]*UserCache, len(caches))

	for userID, v := range caches {
		userCache := v.(*UserCache)

		if extension&UserCacheWithGuild != 0 {
			// 存在第三方干预玩家公会，如通过审批和踢出等操作，所以公会从redis取比较准确
			userGuild, err := h.GetUserGuildData(ctx, userID)
			if err == nil {
				userCache.GuildID = userGuild.GuildID
				userCache.GuildName = userGuild.GuildName
			}

		}

		if extension&UserCacheWithOnline != 0 {
			// 玩家在线状态实时性比较高，所以从redis取
			lastLoginTime, err := h.UserLastLoginTime(ctx, userID)
			if err == nil {
				if lastLoginTime != 0 {
					userCache.LastLoginTime = lastLoginTime
					userCache.OnlineStatus = 1
				}
			}
		}

		ret[userID] = userCache
	}

	return ret, err
}

// 玩家状态
type UserStatus struct {
	key    string
	client *redis.Client
}

func NewUserStatus(client *redis.Client) *UserStatus {
	return &UserStatus{
		key:    KeyUserStatus,
		client: client,
	}
}

func (g *UserStatus) UserOnline(ctx context.Context, userID int64, lastLoginTime int64, expire time.Duration) error {
	err := g.client.Set(ctx, key.MakeRedisKey(g.key, userID), lastLoginTime, expire).Err()
	if err != nil {
		return err
	}

	return nil
}

func (g *UserStatus) UserOffline(ctx context.Context, userID int64) error {
	err := g.client.Del(ctx, key.MakeRedisKey(g.key, userID)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (g *UserStatus) UserLastLoginTime(ctx context.Context, userID int64) (int64, error) {
	ret, err := g.client.Get(ctx, key.MakeRedisKey(g.key, userID)).Int64()
	if err == redis.Nil {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return ret, nil
}

type RepeatUserNameSelector struct {
	client *redis.Client
}

func NewRepeatUserNameSelector(client *redis.Client) *RepeatUserNameSelector {
	return &RepeatUserNameSelector{
		client: client,
	}
}

func (r *RepeatUserNameSelector) FetchNicknameRepeatNum(ctx context.Context, nickname string) (int32, error) {
	count, err := r.client.HGet(ctx, KeyNicknameIndex, nickname).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return int32(count), nil
}
func (r *RepeatUserNameSelector) IncrNicknameRepeatNum(ctx context.Context, nickname string) (int32, error) {
	count, err := r.client.HIncrBy(ctx, KeyNicknameIndex, nickname, 1).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return int32(count), nil
}

// NickNameSuffix 生成初始昵称后缀 //TODO: 临时规则
type NickNameSuffix struct {
	incr *global.IncrID
}

func NewNickNameSuffix(client *redis.Client) *NickNameSuffix {
	return &NickNameSuffix{
		incr: global.NewIncrID(client),
	}
}

func (n *NickNameSuffix) GenNickNameSuffix(ctx context.Context) (int64, error) {
	return n.incr.GenID(ctx, "nickName:prefix")
}
