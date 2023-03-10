package model

import "shared/protobuf/pb"

type Tower struct {
	Id             int32 `json:"id"`
	CurStage       int32 `json:"cur_stage"`
	TodayGoUpTimes int32 `json:"tody_go_up_times"`
}

type TowerInfo struct {
	Towers map[int32]*Tower `json:"towers"`
	*DailyRefreshChecker
}

func NewTower(id int32) *Tower {
	return &Tower{
		Id:             id,
		CurStage:       0,
		TodayGoUpTimes: 0,
	}
}

func NewTowerInfo() *TowerInfo {
	return &TowerInfo{
		Towers:              map[int32]*Tower{},
		DailyRefreshChecker: NewDailyRefreshChecker(),
	}
}

//----------------------------------------
//TowerInfo
//----------------------------------------
func (ti *TowerInfo) AddTower(id int32) *Tower {
	tower, ok := ti.Towers[id]
	if ok {
		return tower
	}

	tower = NewTower(id)
	ti.Towers[id] = tower

	return tower
}

func (ti *TowerInfo) GetTower(id int32) (*Tower, bool) {
	tower, ok := ti.Towers[id]
	return tower, ok
}

func (ti *TowerInfo) GetOrCreateTower(id int32) *Tower {
	tower, ok := ti.GetTower(id)
	if !ok {
		tower = ti.AddTower(id)
	}

	return tower
}

func (ti *TowerInfo) ResetDaily() {
	for _, tower := range ti.Towers {
		tower.ResetDaily()
	}
}

func (ti *TowerInfo) VOTowerInfo() []*pb.VOTower {
	data := make([]*pb.VOTower, 0, len(ti.Towers))

	for _, tower := range ti.Towers {
		data = append(data, tower.VOTower())
	}

	return data
}

//----------------------------------------
//Tower
//----------------------------------------
func (t *Tower) RecordGoUp() {
	t.CurStage = t.CurStage + 1
	t.TodayGoUpTimes = t.TodayGoUpTimes + 1
}

func (t *Tower) GetCurStage() int32 {
	return t.CurStage
}

func (t *Tower) GetTodayGoUpTimes() int32 {
	return t.TodayGoUpTimes
}

func (t *Tower) ResetDaily() {
	t.TodayGoUpTimes = 0
}

func (t *Tower) VOTower() *pb.VOTower {
	return &pb.VOTower{
		Id:             t.Id,
		CurStage:       t.CurStage,
		TodayGoUpTimes: t.TodayGoUpTimes,
	}
}
