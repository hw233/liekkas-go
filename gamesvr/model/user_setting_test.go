package model

import (
	"context"
	"gamesvr/manager"
	"testing"
)

func TestUser_VOVisitingCard(t *testing.T) {
	vo := TestUser.VOVisitingCard()
	t.Log(vo)
	userCache, err := manager.Global.GetUserCache(context.Background(), 20220107181855)
	t.Log(err)

	vo = NewOthersUserVisitingCard(userCache).VOVisitingCard()
	t.Log(vo)

}
