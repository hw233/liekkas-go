package session

import (
	"context"
	"shared/common"
	"shared/csv/static"
	"shared/protobuf/pb"
)

// 侍从玩家所有角色
func (s *Session) GetCharacterList(ctx context.Context, req *pb.C2SGetCharacterList) (*pb.S2CGetCharacterList, error) {
	return &pb.S2CGetCharacterList{Characters: s.User.VOUserCharacter()}, nil
}

// 侍从升级（使用经验药水）
func (s *Session) CharacterLvUp(ctx context.Context, req *pb.C2SCharacterLvUp) (*pb.S2CCharacterLvUp, error) {
	// check param
	if len(req.Costs) == 0 {
		return nil, common.ErrParamError
	}
	rewards, err := common.ParseFromVOConsume(req.Costs)
	if err != nil {
		return nil, err
	}
	// check consume
	err = s.User.CheckRewardsEnough(rewards)
	if err != nil {
		return nil, err
	}

	character, err := s.User.CharacterLvUp(req.CharacterId, rewards)
	if err != nil {
		return nil, err
	}
	return &pb.S2CCharacterLvUp{
		Result: character.VOUserCharacter(),
		// 在session层返回VOResourceResult
		ResourceResult: s.VOResourceResult(),
	}, nil
}

// 侍从升星
func (s *Session) CharacterStarUp(ctx context.Context, req *pb.C2SCharacterStarUp) (*pb.S2CCharacterStarUp, error) {
	chara, voRewardResult, err := s.User.CharacterStarUp(req.CharacterId)

	if err != nil {
		return nil, err
	}

	if chara.Rarity == static.RaritySsr {
		s.User.AddGreetings(chara.ID, 1, chara.GetStar())
	}

	ret := &pb.S2CCharacterStarUp{
		Result:         chara.VOUserCharacter(),
		ResourceResult: voRewardResult,
	}

	return ret, nil
}

// 侍从升阶
func (s *Session) CharacterStageUp(ctx context.Context, req *pb.C2SCharacterStageUp) (*pb.S2CCharacterStageUp, error) {
	chara, voRewardResult, err := s.User.CharacterStageUp(req.CharacterId)

	if err != nil {
		return nil, err
	}

	ret := &pb.S2CCharacterStageUp{
		Result:         chara.VOUserCharacter(),
		ResourceResult: voRewardResult,
	}

	return ret, nil
}

// 侍从等级提升
func (s *Session) CharacterSkillLvUp(ctx context.Context, req *pb.C2SCharacterSkillLvUp) (*pb.S2CCharacterSkillLvUp, error) {
	// check param
	if req.LvUpAmount <= 0 {
		return nil, common.ErrParamError
	}

	character, err := s.User.CharacterSkillLvUp(req.CharacterId, req.SkillId, req.LvUpAmount)
	if err != nil {
		return nil, err
	}

	return &pb.S2CCharacterSkillLvUp{
		Character: character.VOUserCharacter(),
		// 在session层返回VOResourceResult
		ResourceResult: s.VOResourceResult(),
	}, nil
}

// 穿装备
func (s *Session) CharacterWear(ctx context.Context, req *pb.C2SCharacterWear) (*pb.S2CCharacterWear, error) {
	// check param
	equipments, worldItems, err := s.User.CharacterWear(req.CharacterID, req.EquipmentUIDs, req.WorldItemUIDs)
	if err != nil {
		return nil, err
	}

	voEquipments := make([]*pb.VOUserEquipment, 0, len(equipments))
	for _, equipment := range equipments {
		voEquipments = append(voEquipments, equipment.VOUserEquipment())
	}

	voWorldItems := make([]*pb.VOUserWorldItem, 0, len(worldItems))
	for _, worldItem := range worldItems {
		voWorldItems = append(voWorldItems, worldItem.VOUserWorldItem())
	}

	return &pb.S2CCharacterWear{
		Equipments: voEquipments,
		WorldItems: voWorldItems,
	}, nil
}

// 脱装备
func (s *Session) CharacterUndress(ctx context.Context, req *pb.C2SCharacterUndress) (*pb.S2CCharacterUndress, error) {
	equipments, worldItems, err := s.User.CharacterUndress(req.CharacterID, req.EquipmentUIDs, req.WorldItemUIDs)
	if err != nil {
		return nil, err
	}

	voEquipments := make([]*pb.VOUserEquipment, 0, len(equipments))
	for _, equipment := range equipments {
		voEquipments = append(voEquipments, equipment.VOUserEquipment())
	}

	voWorldItems := make([]*pb.VOUserWorldItem, 0, len(worldItems))
	for _, worldItem := range worldItems {
		voWorldItems = append(voWorldItems, worldItem.VOUserWorldItem())
	}

	return &pb.S2CCharacterUndress{
		Equipments: voEquipments,
		WorldItems: voWorldItems,
	}, nil
}
