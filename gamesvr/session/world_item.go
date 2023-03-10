package session

import (
	"context"

	"shared/common"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/utility/errors"
)

// 世界级道具列表
func (s *Session) WorldItemList(ctx context.Context, req *pb.C2SWorldItemList) (*pb.S2CWorldItemList, error) {
	return &pb.S2CWorldItemList{
		WorldItems: s.User.VOUserWorldItem(),
		RecastID:   s.User.WorldItemPack.RecastID,
		RecastCamp: int32(s.User.WorldItemPack.RecastCamp),
	}, nil
}

// 世界级道具强化
func (s *Session) WorldItemStrengthen(ctx context.Context, req *pb.C2SWorldItemStrengthen) (*pb.S2CWorldItemStrengthen, error) {
	if len(req.Items) != len(common.WorldItemEXPItems) {
		return nil, common.ErrParamError
	}

	equipment, err := s.User.WorldItemStrengthen(req.Id, req.Items, req.Materials)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CWorldItemStrengthen{
		WorldItem:      equipment.VOUserWorldItem(),
		ResourceResult: s.VOResourceResult(),
		Materials:      req.Materials,
	}, nil
}

// 世界级道具进阶
func (s *Session) WorldItemAdvance(ctx context.Context, req *pb.C2SWorldItemAdvance) (*pb.S2CWorldItemAdvance, error) {
	// 两种强化方式只能选一种
	if len(req.Materials) > 0 && req.ItemID != 0 {
		return nil, common.ErrParamError
	}

	equipment, err := s.User.WorldItemAdvance(req.Id, req.Materials, req.ItemID)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	if equipment.Rarity == static.RaritySsr {
		s.User.AddGreetings(equipment.WID, 2, equipment.Stage.Value())
	}

	return &pb.S2CWorldItemAdvance{
		WorldItem:      equipment.VOUserWorldItem(),
		ResourceResult: s.VOResourceResult(),
		Materials:      req.Materials,
	}, nil
}

// 世界级道具加锁解锁
func (s *Session) WorldItemLock(ctx context.Context, req *pb.C2SWorldItemLock) (*pb.S2CWorldItemLock, error) {
	err := s.User.WorldItemLock(req.Id, req.IsLock)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CWorldItemLock{
		Id:     req.Id,
		IsLock: req.IsLock,
	}, nil
}
