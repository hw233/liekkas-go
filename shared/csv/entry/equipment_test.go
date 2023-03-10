package entry

import (
	"testing"

	"shared/common"
	"shared/utility/errors"
)

func TestEquipmentStrengthenEXP(t *testing.T) {
	target, err := CSV.Equipment.NewEquipment(1, 1174001)
	if err != nil {
		t.Errorf("%+v", errors.Format(err))
	}

	material1, err := CSV.Equipment.NewEquipment(1, 1174001)
	if err != nil {
		t.Errorf("%+v", errors.Format(err))
	}

	material1.EXP.SetValue(4238080)

	err = CSV.Equipment.SyncLevelAndUnlockAttr(material1)
	if err != nil {
		t.Errorf("%+v", errors.Format(err))
	}

	// material2, err := CSV.Equipment.NewEquipment(1, 1173001)
	// if err != nil {
	// 	t.Errorf("%+v", errors.Format(err))
	// }

	var materials = []*common.Equipment{material1}

	addEXP, err := CSV.Equipment.StrengthenEXP(target, 0, materials, 60)
	if err != nil {
		t.Errorf("%+v", errors.Format(err))
	}

	costGolds, err := CSV.Equipment.StrengthenGoldCost(target, 0, materials)

	target.EXP.Plus(addEXP)

	err = CSV.Equipment.SyncLevelAndUnlockAttr(target)
	if err != nil {
		t.Errorf("%+v", errors.Format(err))
	}

	t.Logf("addEXP: %d, costGolds: %d , targetEXP: %d, targetLevel: %d", addEXP, costGolds, target.EXP, target.Level)
}

func TestEquipmentSyncLevelAndUnlockAttr(t *testing.T) {
	equipment, err := CSV.Equipment.NewEquipment(1, 1172001)
	if err != nil {
		t.Errorf("%+v", errors.Format(err))
	}

	t.Logf("equipment attr before: %+v", equipment.Attrs)

	equipment.EXP.Plus(155000)
	err = CSV.Equipment.SyncLevelAndUnlockAttr(equipment)
	if err != nil {
		t.Errorf("%+v", errors.Format(err))
	}

	t.Logf("equipment level: %d, attr after: %+v", equipment.Level, equipment.Attrs)
}

func TestEquipmentRecastCamp(t *testing.T) {
	camp, err := CSV.Equipment.RecastCamp(common.NewEquipment(1, 1172001))
	if err != nil {
		t.Errorf("%+v", errors.Format(err))
	}

	t.Logf("camp: %v", camp)
}

func TestEquipmentCheckAdvanceMaterials(t *testing.T) {
	target, err := CSV.Equipment.NewEquipment(1, 1174001)
	if err != nil {
		t.Errorf("%+v", errors.Format(err))
	}

	material1, err := CSV.Equipment.NewEquipment(1, 1174001)
	if err != nil {
		t.Errorf("%+v", errors.Format(err))
	}

	var materials = []*common.Equipment{material1}

	cost, err := CSV.Equipment.AdvanceCostGold(target)
	if err != nil {
		t.Errorf("AdvanceCostGold: %+v", errors.Format(err))
		return
	}

	err = CSV.Equipment.CheckAdvanceMaterials(target, materials)
	if err != nil {
		t.Errorf("CheckAdvanceMaterials: %+v", errors.Format(err))
		return
	}

	err = CSV.Equipment.CheckStageUpToLimit(target)
	if err != nil {
		t.Errorf("CheckStageUpToLimit: %+v", errors.Format(err))
		return
	}
	t.Logf("cost: %v", cost)
}
