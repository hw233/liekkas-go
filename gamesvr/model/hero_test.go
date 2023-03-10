package model

import (
	"context"
	"fmt"
	"shared/common"
	"testing"
)

func TestUnlockHero(t *testing.T) {
	user := NewUser(1001)
	user.InitForCreate(context.Background())

	heroId := int32(1)

	hero, err := user.UnlockHero(heroId)
	if err != nil {
		t.Errorf("failed, %s", err.Error())
		return
	}

	fmt.Printf("hero: %v\n", hero)

	user.addReward(common.NewReward(6, 100), nil)

	err = user.HeroLevelUp()
	if err != nil {
		t.Errorf("failed, %s", err.Error())
		return
	}
	fmt.Printf("user.HeroPack.Level: %v\n", user.HeroPack.Level)

	err = user.HeroSkillLevelUpgrade(heroId, 110)
	if err != nil {
		t.Errorf("failed, %s", err.Error())
		return
	}
	fmt.Printf("hero: %v\n", hero)
}
