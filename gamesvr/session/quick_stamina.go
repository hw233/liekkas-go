package session

import (
	"context"
	"shared/protobuf/pb"
)

func (s *Session) GetStaminaInfo(ctx context.Context, req *pb.C2SGetStaminaInfo) (*pb.S2CGetStaminaInfo, error) {
	s.User.CheckStaminaUpdate()

	stamina, err := s.User.GetStaminaInfo()
	if err != nil {
		return nil, err
	}

	return &pb.S2CGetStaminaInfo{
		StaminaInfo: stamina,
	}, nil
}

func (s *Session) QuickPurchaseStamina(ctx context.Context, req *pb.C2SQuickPurchaseStamina) (*pb.S2CQuickPurchaseStamina, error) {

	s.User.CheckStaminaUpdate()

	err := s.User.QuickPurchaseStamina()
	if err != nil {
		return nil, err
	}

	stamina, err := s.User.GetStaminaInfo()
	if err != nil {
		return nil, err
	}

	return &pb.S2CQuickPurchaseStamina{
		StaminaInfo:    stamina,
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}
