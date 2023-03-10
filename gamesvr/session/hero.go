package session

import (
	"context"
	"shared/csv/static"
	"shared/protobuf/pb"
)

func (s *Session) HeroInfo(ctx context.Context, req *pb.C2SHeroInfo) (*pb.S2CHeroInfo, error) {
	return &pb.S2CHeroInfo{
		HeroInfo: s.User.HeroPack.VOHeroInfo(),
	}, nil
}

func (s *Session) HeroLevelUp(ctx context.Context, req *pb.C2SHeroLevelUp) (*pb.S2CHeroLevelUp, error) {
	err := s.User.HeroLevelUp()
	if err != nil {
		return nil, err
	}

	return &pb.S2CHeroLevelUp{
		Level: s.User.HeroPack.GetLevel(),
	}, nil
}

func (s *Session) UnlockHero(ctx context.Context, req *pb.C2SUnlockHero) (*pb.S2CUnlockHero, error) {
	err := s.User.CheckActionUnlock(static.ActionIdTypeHerogrowthunlock)
	if err != nil {
		return nil, err
	}

	hero, err := s.User.UnlockHero(req.HeroId)
	if err != nil {
		return nil, err
	}

	return &pb.S2CUnlockHero{
		Hero:           hero.VOHero(),
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

func (s *Session) HeroSkillUpgrade(ctx context.Context, req *pb.C2SHeroSkillUpgrade) (*pb.S2CHeroSkillUpgrade, error) {
	err := s.User.HeroSkillLevelUpgrade(req.HeroId, req.SkillId)
	if err != nil {
		return nil, err
	}

	hero, _ := s.User.HeroPack.GetHero(req.HeroId)

	return &pb.S2CHeroSkillUpgrade{
		Hero:           hero.VOHero(),
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

func (s *Session) AddHeroAttendant(ctx context.Context, req *pb.C2SAddHeroAttendant) (*pb.S2CAddHeroAttendant, error) {
	oldCharaId, oldHeroId, err := s.User.AddHeroAttendant(req.HeroId, req.Slot, req.CharaId)
	if err != nil {
		return nil, err
	}

	resp := &pb.S2CAddHeroAttendant{
		Heros:      make([]*pb.VOHero, 0),
		Characters: make([]*pb.VOUserCharacter, 0),
	}

	hero, _ := s.User.HeroPack.GetHero(req.HeroId)
	chara, _ := s.User.CharacterPack.Get(req.CharaId)

	resp.Heros = append(resp.Heros, hero.VOHero())
	resp.Characters = append(resp.Characters, chara.VOUserCharacter())

	if oldCharaId > 0 {
		oldChara, _ := s.User.CharacterPack.Get(oldCharaId)
		resp.Characters = append(resp.Characters, oldChara.VOUserCharacter())
	}

	if oldHeroId > 0 {
		oldHero, _ := s.User.HeroPack.GetHero(oldHeroId)
		resp.Heros = append(resp.Heros, oldHero.VOHero())
	}

	return resp, nil
}

func (s *Session) RemoveHeroAttendant(ctx context.Context, req *pb.C2SRemoveHeroAttendant) (*pb.S2CRemoveHeroAttendant, error) {
	oldCharaId, err := s.User.RemoveHeroAttendant(req.HeroId, req.Slot)
	if err != nil {
		return nil, err
	}

	hero, _ := s.User.HeroPack.GetHero(req.HeroId)
	resp := &pb.S2CRemoveHeroAttendant{
		Hero: hero.VOHero(),
	}

	if oldCharaId > 0 {
		oldChara, _ := s.User.CharacterPack.Get(oldCharaId)
		resp.Character = oldChara.VOUserCharacter()
	}

	return resp, nil
}

func (s *Session) UpdateHeroLastCalcAttr(ctx context.Context, req *pb.C2SUpdateHeroLastCalcAttr) (*pb.S2CUpdateHeroLastCalcAttr, error) {
	attrs := make(map[int32]int32, len(req.Charas))
	for _, chara := range req.Charas {
		attrs[chara.Slot] = chara.LastCalcAttr
	}

	err := s.User.UpdateHeroLastCalcAttr(req.HeroId, req.HeroAttr, attrs)
	if err != nil {
		return nil, err
	}

	hero, _ := s.User.HeroPack.GetHero(req.HeroId)
	return &pb.S2CUpdateHeroLastCalcAttr{
		Hero: hero.VOHero(),
	}, nil
}

func (s *Session) HeroSetNewOneFlag(ctx context.Context, req *pb.C2SHeroSetNewOneFlag) (*pb.S2CHeroSetNewOneFlag, error) {
	err := s.User.SetHeroNewOneFlag(req.NewOne)
	if err != nil {
		return nil, err
	}

	return &pb.S2CHeroSetNewOneFlag{
		NewOne: req.NewOne,
	}, nil
}
