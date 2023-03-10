package model

import (
	"testing"
)

func TestGachaRecords_Gacha(t *testing.T) {
	gachaDrop, err := TestUser.UserGachaDrop(1, false)
	if err != nil {
		t.Errorf("failed, %s", err.Error())
		return
	}

	t.Logf("gachaDrop: %v", gachaDrop)
	t.Logf("resourceResult:%+v", TestUser.VOResourceResult())
}
