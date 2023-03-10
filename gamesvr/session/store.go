package session

import (
	"context"
	"gamesvr/manager"
	"shared/protobuf/pb"
)

func (s *Session) StoreGetGoods(ctx context.Context, req *pb.C2SStoreGetGoods) (*pb.S2CStoreGetGoods, error) {

	storeIDs := req.StoreIds

	infos := []*pb.VOStoreInfo{}

	for _, storeID := range storeIDs {
		err := s.User.CheckGoodsInfo(storeID)
		if err != nil {
			return nil, err
		}

		storeInfo, err := s.User.GetStoreInfo(storeID)
		if err != nil {
			return nil, err
		}
		if storeInfo != nil {
			infos = append(infos, storeInfo)
		}
	}

	return &pb.S2CStoreGetGoods{
		StoreInfos: infos,
	}, nil
}

func (s *Session) StorePurchase(ctx context.Context, req *pb.C2SStorePurchase) (*pb.S2CStorePurchase, error) {

	cnt, err := s.User.Purchase(req.StoreId, req.SubStoreId, req.CellIndex, req.GoodsId, req.Currency, req.Num)
	if err != nil {
		return nil, err
	}

	return &pb.S2CStorePurchase{
		StoreId:        req.StoreId,
		SubStoreId:     req.SubStoreId,
		CellIndex:      req.CellIndex,
		GoodsId:        req.GoodsId,
		Num:            cnt,
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

func (s *Session) StoreUpdate(ctx context.Context, req *pb.C2SStoreUpdate) (*pb.S2CStoreUpdate, error) {

	err := s.User.ForceToUpdate(req.StoreId)
	if err != nil {
		return nil, err
	}

	storeInfo, err := s.User.GetStoreInfo(req.StoreId)
	if err != nil {
		return nil, err
	}

	return &pb.S2CStoreUpdate{
		StoreInfos: storeInfo,
	}, nil
}

func (s *Session) StoreQuickPurchase(ctx context.Context, req *pb.C2SStoreQuickPurchase) (*pb.S2CStoreQuickPurchase, error) {

	err := s.User.QuickPurchase(req.GoodsID)
	if err != nil {
		return nil, err
	}
	return &pb.S2CStoreQuickPurchase{
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

// ---------notify----------------
func (s *Session) TryPushStoreNotify() {
	notify := s.PopStoreNotify()
	if notify == nil {
		return
	}
	s.push(manager.CSV.Protocol.Pushes.StoreDailyRefresh, notify)
}
