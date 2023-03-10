package global

import (
	"context"

	"github.com/go-redis/redis/v8"

	"shared/utility/global"
)

type UserID struct {
	incrID *global.IncrID
}

func NewUserID(client *redis.Client) *UserID {
	return &UserID{
		incrID: global.NewIncrID(client),
	}
}

func (u *UserID) GenID(ctx context.Context) (int64, error) {
	return u.incrID.GenID(ctx, KeyUserID)
}

func (u *UserID) SetInitID(ctx context.Context, initID int64) error {
	return u.incrID.SetInitID(ctx, KeyUserID, initID)
}
