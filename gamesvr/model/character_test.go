package model

import (
	"testing"
)

const (
	testCharId = 1001
)

func TestCharaStarUp(t *testing.T) {
	user := NewUser(1001)
	user.addCharacter(testCharId, 1, user, nil)

	chara, _ := user.CharacterPack.Get(testCharId)

	cost, err := user.CharacterPack.CalcCharacterStarUpCost(testCharId, chara.GetStar()+1)
	if err != nil {
		t.Errorf("failed, %s", err.Error())
	}

	user.addRewards(&cost, nil)

	_, vo, err := user.CharacterStarUp(testCharId)
	if err != nil {
		t.Errorf("failed, %s", err.Error())
	}
	t.Logf("VOResourceResult: %v\n", vo)
}

func TestCharaStageUp(t *testing.T) {
	user := NewUser(1001)
	user.addCharacter(testCharId, 1, user, nil)
	chara, _ := user.CharacterPack.Get(testCharId)

	cost, err := user.CharacterPack.CalcCharacterStageUpCost(testCharId, chara.GetStage()+1)
	if err != nil {
		t.Errorf("failed, %s", err.Error())
	}

	user.addRewards(&cost, nil)

	_, vo, err := user.CharacterStageUp(testCharId)
	if err != nil {
		t.Errorf("failed, %s", err.Error())
	}
	t.Logf("VOResourceResult: %v\n", vo)
}

func TestUser_CharacterSkillLvUp(t *testing.T) {
	user := NewUser(1001)
	user.addCharacter(1011, 1, user, nil)
	user.CharacterSkillLvUp(1011, 101101, 1)
}
