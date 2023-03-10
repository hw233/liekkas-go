package model

import (
	"shared/common"
	"shared/statistic/logreason"
	"testing"
)

func TestGiftRandom(t *testing.T) {
	user := NewUser(1001)
	err := GiftRandom(user, 20015, 10, 0, logreason.EmptyReason())
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	t.Log(user.VOResourceResult())
}

func TestGiftSelectOne(t *testing.T) {
	user := NewUser(1001)

	err := GiftSelectOne(user, 20014, 10, 0, logreason.EmptyReason())
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	t.Log(user.VOResourceResult())
}
func TestDrop(t *testing.T) {
	user := NewUser(1001)
	rewards := common.NewRewards()
	rewards.AddReward(common.NewReward(1600211, 17))
	rewards.AddReward(common.NewReward(1600221, 3))

	_, err := user.addRewards(rewards, logreason.EmptyReason())
	if err != nil {
		t.Errorf("%+v", err)
		return
	}
	result := user.VOResourceResult()
	t.Log(result)
	t.Log(len(result.Rewards))

	t.Log(len(result.Equipments))
}
