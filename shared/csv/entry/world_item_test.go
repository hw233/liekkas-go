package entry

import (
	"testing"

	"shared/common"
	"shared/utility/errors"
)

func TestWorldItemStrengthenEXP(t *testing.T) {
	target, err := CSV.WorldItem.NewWorldItem(1, 10001)
	if err != nil {
		t.Errorf("%+v", errors.Format(err))
	}

	material1, err := CSV.WorldItem.NewWorldItem(1, 10002)
	if err != nil {
		t.Errorf("%+v", errors.Format(err))
	}

	material1.EXP.SetValue(1000)

	material2, err := CSV.WorldItem.NewWorldItem(1, 10003)
	if err != nil {
		t.Errorf("%+v", errors.Format(err))
	}

	var materials = []*common.WorldItem{material1, material2}

	addEXP, err := CSV.WorldItem.StrengthenEXP(target, 0, materials, 10)
	if err != nil {
		t.Errorf("%+v", errors.Format(err))
	}

	t.Logf("addEXP: %d", addEXP)
}

func TestWorldItemAdvance(t *testing.T) {
	target, err := CSV.WorldItem.NewWorldItem(1, 10001)
	if err != nil {
		t.Errorf("%+v", errors.Format(err))
	}

	material1, err := CSV.WorldItem.NewWorldItem(1, 10001)
	if err != nil {
		t.Errorf("%+v", errors.Format(err))
	}

	material1.Stage.Plus(9)

	var materials = []*common.WorldItem{material1}

	err = CSV.WorldItem.CheckStageUpToLimit(target)
	if err != nil {
		t.Errorf("%+v", errors.Format(err))
	}

	err = CSV.WorldItem.CheckAdvanceMaterials(target, materials, 0)
	if err != nil {
		t.Errorf("%+v", errors.Format(err))
	}

	addStage, err := CSV.WorldItem.CalAddStage(target, materials, 0)
	if err != nil {
		t.Errorf("%+v", errors.Format(err))
	}

	target.Stage.Plus(addStage)

	t.Logf("materialStage: %v, addStage: %d,  retStage: %v", material1.Stage, addStage, target.Stage)
}

func TestWorldItemCalculateEXP(t *testing.T) {
	exp := CSV.Item.CalculateTotalWorldItemEXP([]int32{1, 0, 0, 0})

	t.Logf("exp: %d", exp)
}
