package entry

import (
	"shared/utility/coordinate"
	"testing"
)

func TestYggdrasilEntry_Reload(t *testing.T) {
	position1 := *coordinate.NewPosition(1, 1)
	position2 := *coordinate.NewPosition(1, 1)
	t.Log(position1 == position2)

}

func TestYggdrasilEntry_GetPosAreaId(t *testing.T) {
	areaId, err := CSV.Yggdrasil.GetPosAreaId(*coordinate.NewPosition(-8, -24))
	t.Log(err)
	t.Log(areaId)
	areaId, err = CSV.Yggdrasil.GetPosAreaId(*coordinate.NewPosition(-113, 58))
	t.Log(err)

	t.Log(areaId)
}
func TestYggdrasilEntry_GetPosInitHeightType(t *testing.T) {
	H, T := CSV.Yggdrasil.GetPosInitHeightAndType(*coordinate.NewPosition(-8, -24))
	t.Log(H)

	t.Log(T)
}

func TestYggdrasilEntry_Check(t *testing.T) {
	_, posType := CSV.Yggdrasil.GetPosInitHeightType()
	t.Log(posType[*coordinate.NewPosition(-7, -26)])
}
func TestYggdrasilEntry_GetMostCost(t *testing.T) {
	cost, err := CSV.Yggdrasil.GetMostCost(*coordinate.NewPosition(12, 0))
	t.Log(cost)
	t.Log(err)

}

func TestYggdrasilEntry_Area(t *testing.T) {
	area1, _ := CSV.Yggdrasil.GetArea(1)
	t.Log(area1.Area)
	t.Log(area1.Area.Count())
	area2, _ := CSV.Yggdrasil.GetArea(2)
	cut := area1.Area.CutArea(area2.Area)
	minus := area1.Area.MinusArea(area2.Area)

	t.Log(cut.Count())
	t.Log(minus.Count())

}

func TestYggdrasilEntry_GetYggBagAllCount(t *testing.T) {
	count := CSV.Yggdrasil.GetYggBagAllCount(59)
	t.Log(count)
	teamCount := CSV.Yggdrasil.GetYggEditTeamCount(59)
	t.Log(teamCount)

}

func TestYggdrasilEntry_GetYggInitPos(t *testing.T) {

	pos := CSV.Yggdrasil.GetYggInitPos()
	i := &pos
	i.X = 23
	t.Log(CSV.Yggdrasil.GetYggInitPos())
	count := CSV.Yggdrasil.GetYggMailMaxCount()
	count = 1
	t.Log(count)

	t.Log(CSV.Yggdrasil.GetYggMailMaxCount())

}

func TestYggdrasilEntry_GetExploreProcessIndex(t *testing.T) {

	index, err := CSV.Yggdrasil.GetExploreProcessIndex(1, 38)
	t.Log(index)
	t.Log(err)

	area, err := CSV.Yggdrasil.GetArea(1)
	dropId := area.ExploredProgressDrop[index]
	t.Log(dropId)

}

func TestYggdrasilEntry_IsCityEntrance(t *testing.T) {
	cityId := CSV.Yggdrasil.IsCityEntrance(coordinate.NewPosition(-35, -6))
	t.Log(cityId)

}

func TestYggdrasilEntry_GetArea(t *testing.T) {
	id, err := CSV.Yggdrasil.GetPosAreaId(*coordinate.NewPosition(11, -26))
	t.Log(id)
	t.Log(err)
}
func TestYggdrasilEntry_GetClosestSafePos(t *testing.T) {
	p, err := CSV.Yggdrasil.GetClosestSafePos(*coordinate.NewPosition(11, -26))
	t.Log(p)
	t.Log(err)
}
