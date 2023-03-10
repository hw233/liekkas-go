package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/servertime"
)

const (
	UpdateForNone  int32 = 0
	UpdateForDay   int32 = 1
	UpdateForWeek  int32 = 2
	UpdateForMonth int32 = 3
)

type StoreInfo struct {
	*DailyRefreshChecker
	Stores    map[int32]*Store    `json:"stores"`
	SubStores map[int32]*SubStore `json:"subStores"`
}

type Store struct {
	*GoodsUpdate `json:"update_info"`
	SubStores    []int32 `json:"subStore_ids"`
}

type GoodsUpdate struct {
	UpdateTimes    int32 `json:"update_times"`
	UpdatePeriod   int32 `json:"update_period"`
	LastTimeUpdate int64 `json:"last_time_update"`
}

type SubStore struct {
	ID int32 `json:"subStore_id"`
	//GoodsRecord map[int32]int32
	Cells []*Cell `json:"cells"`
}

type Cell struct {
	GoodsID      int32 `json:"goods_id"`
	NumOfRemains int32 `json:"num_of_remains"`
}

func NewStoreInfo() *StoreInfo {
	return &StoreInfo{
		DailyRefreshChecker: NewDailyRefreshChecker(),
		Stores:              map[int32]*Store{},
		SubStores:           map[int32]*SubStore{},
	}
}

func NewStore(period int32) *Store {

	return &Store{
		GoodsUpdate: NewGoodsUpdate(period),
		SubStores:   []int32{},
	}
}

func NewGoodsUpdate(period int32) *GoodsUpdate {
	return &GoodsUpdate{
		UpdateTimes:    0,
		UpdatePeriod:   period,
		LastTimeUpdate: 0,
	}
}

func NewSubStore(id int32) *SubStore {
	return &SubStore{
		ID:    id,
		Cells: []*Cell{},
	}

}

//----------------------------------------
//Store
//----------------------------------------

// 返回值为真，代表商店中的数据需要更新
// func (s *Store) checkForUpdate() bool {
// 	return s.checkUpdate()
// }

func (s *Store) checkForSubStoreID(id int32) bool {
	for _, subID := range s.SubStores {
		if subID == id {
			return true
		}
	}

	return false
}

//----------------------------------------
//GoodsUpdate
//----------------------------------------
func (g *GoodsUpdate) checkUpdate() bool {
	if g.LastTimeUpdate == 0 { // 代表是新建的商店
		return true
	}
	switch g.UpdatePeriod {
	case UpdateForNone:
		return false
	case UpdateForDay:
		return g.updateForDay()
	case UpdateForWeek:
		return g.updateForWeek()
	case UpdateForMonth:
		return g.updateForMonth()
	}

	return false
}

func (g *GoodsUpdate) updateForWeek() bool {
	tm := servertime.Now()
	return g.LastTimeUpdate < WeekRefreshTime(tm).Unix()
}

func (g *GoodsUpdate) updateForDay() bool {
	tm := servertime.Now()
	return g.LastTimeUpdate < DailyRefreshTime(tm).Unix()
}

func (g *GoodsUpdate) updateForMonth() bool {
	tm := servertime.Now()
	return g.LastTimeUpdate < MonthRefreshTime(tm).Unix()
}

//----------------------------------------
//Subtore
//----------------------------------------

func (s *SubStore) VOSubStoreInfo() *pb.VOSubStoreInfo {

	cellInfos := make([]*pb.VOCellInfo, 0, len(s.Cells))
	for _, cell := range s.Cells {
		cellInfo := &pb.VOCellInfo{
			GoodsId:    cell.GoodsID,
			RemainsNum: cell.NumOfRemains,
		}
		cellInfos = append(cellInfos, cellInfo)
	}
	return &pb.VOSubStoreInfo{
		SubStoreId: s.ID,
		CellInfos:  cellInfos,
	}
}

func (s *SubStore) checkForNum(index, goodsID, num int32) error {

	if int(index) >= len(s.Cells) {
		return errors.Swrapf(common.ErrStoreIndexOutOfRangeForInfo, s.ID, index)
	}

	if goodsID != s.Cells[index].GoodsID {
		return errors.Swrapf(common.ErrStoreWrongGoodsIDInCellForInfo, s.ID, goodsID)
	}

	if s.Cells[index].NumOfRemains != -1 && num > s.Cells[index].NumOfRemains {
		return errors.Swrapf(common.ErrStoreNumOfRemainsNotEnough, s.ID, goodsID)
	}

	return nil
}

// 检查子商店所有商品是否已经售空，适用于阶段式商店调用.返回值为true，代表售空
func (s *SubStore) checkForSoldOut() bool {
	for _, cell := range s.Cells {
		if cell.NumOfRemains != 0 {
			return false
		}
	}

	return true
}

// 刷新子商店的所有商品
func (s *SubStore) updateGoods(ids []int32) error {

	record := make([]*Cell, 0, 20)
	for _, id := range ids {
		cell := &Cell{}
		goods, err := manager.CSV.Store.GetGoods(id)
		if err != nil {
			return err
		}
		cell.GoodsID = id
		cell.NumOfRemains = goods.Times
		record = append(record, cell)
	}

	s.Cells = record

	return nil
}

// 刷新单个商品，适用于购买行为之后, 返回-1代表商品数量是无限的
func (s *SubStore) updateAfterPurchase(index, goodsID int32, num int32) (int32, error) {
	if int(index) >= len(s.Cells) {
		return 0, errors.Swrapf(common.ErrStoreIndexOutOfRangeForInfo, s.ID, index)
	}

	cnt := s.Cells[index].NumOfRemains

	if cnt != -1 { // 如果不是无限数量
		cnt -= num
		s.Cells[index].NumOfRemains = cnt
	}

	return cnt, nil
}
