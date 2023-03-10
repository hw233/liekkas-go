package session

import (
	"context"

	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/errors"
)

// 装备列表
func (s *Session) EquipmentList(ctx context.Context, req *pb.C2SEquipmentList) (*pb.S2CEquipmentList, error) {
	return &pb.S2CEquipmentList{
		Equipments: s.User.VOUserEquipment(),
	}, nil
}

// 装备强化
func (s *Session) EquipmentStrengthen(ctx context.Context, req *pb.C2SEquipmentStrengthen) (*pb.S2CEquipmentStrengthen, error) {
	if len(req.Items) != len(common.EquipmentEXPItems) {
		return nil, common.ErrParamError
	}

	equipment, err := s.User.EquipmentStrengthen(req.Id, req.Items, req.Materials)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CEquipmentStrengthen{
		Equipment:      equipment.VOUserEquipment(),
		ResourceResult: s.VOResourceResult(),
		Materials:      req.Materials,
	}, nil
}

// 装备进阶
func (s *Session) EquipmentAdvance(ctx context.Context, req *pb.C2SEquipmentAdvance) (*pb.S2CEquipmentAdvance, error) {
	equipment, err := s.User.EquipmentAdvance(req.Id, req.Materials)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CEquipmentAdvance{
		Equipment:      equipment.VOUserEquipment(),
		ResourceResult: s.VOResourceResult(),
		Materials:      req.Materials,
	}, nil
}

// 装备阵营重铸
func (s *Session) EquipmentRecastCamp(ctx context.Context, req *pb.C2SEquipmentRecastCamp) (*pb.S2CEquipmentRecastCamp, error) {
	camp, err := s.User.EquipmentRecastCamp(req.Id)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CEquipmentRecastCamp{
		Camp:           int32(camp),
		ResourceResult: s.VOResourceResult(),
	}, nil
}

// 装备阵营重铸确认
func (s *Session) EquipmentConfirmRecastCamp(ctx context.Context, req *pb.C2SEquipmentConfirmRecastCamp) (*pb.S2CEquipmentConfirmRecastCamp, error) {
	equipment, err := s.User.EquipmentConfirmRecastCamp(req.Id, req.Confirm)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CEquipmentConfirmRecastCamp{
		Equipment: equipment.VOUserEquipment(),
	}, nil
}

// 装备加锁解锁
func (s *Session) EquipmentLock(ctx context.Context, req *pb.C2SEquipmentLock) (*pb.S2CEquipmentLock, error) {
	err := s.User.EquipmentLock(req.Id, req.IsLock)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CEquipmentLock{
		Id:     req.Id,
		IsLock: req.IsLock,
	}, nil
}
