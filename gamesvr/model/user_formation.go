package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/slice"
)

func (u *User) SetFormation(id, battleType int32, voFormation *pb.VOBattleFormation) error {
	formationType := BattleTypeToFormationType(int(battleType))

	switch formationType {
	case BattleFormationTypeExplore,
		BattleFormationTypeElite,
		BattleFormationTypeYGG,
		BattleFormationTypeChaos:
		if id > 1 {
			return errors.Swrapf(common.ErrInvalidFormation, formationType, id)
		}

	case BattleFormationTypeTower:
		_, err := manager.CSV.Tower.GetTower(id)
		if err != nil {
			return errors.Swrapf(common.ErrInvalidFormation, formationType, id)
		}
	}

	err := u.checkBattleFormation(battleType, 0, voFormation)
	if err != nil {
		return err
	}

	u.FormationInfo.SetFormation(id, int32(formationType), voFormation)

	return nil
}

func (u *User) checkBattleFormation(battleType, battleId int32, battleFormation *pb.VOBattleFormation) error {
	charas := map[int32]int32{}
	positions := map[int32]int32{}

	if battleFormation == nil {
		return errors.Swrapf(common.ErrInvalidFormationS)
	}

	if battleFormation.BattleCharacters != nil {
		for _, battleChara := range battleFormation.BattleCharacters {
			charaId := battleChara.CharacterId
			chara, err := u.CharacterPack.Get(charaId)
			if err != nil {
				return err
			}

			for _, equipmentId := range chara.Equipments {
				if equipmentId > 0 {
					_, err := u.EquipmentPack.Get(equipmentId)
					if err != nil {
						return err
					}
				}
			}

			if chara.WorldItem > 0 {
				_, err = u.WorldItemPack.Get(chara.WorldItem)
				if err != nil {
					return err
				}
			}

			_, ok := charas[charaId]
			if ok {
				errors.Swrapf(common.ErrBattleDuplicateChara, charaId)
			}

			position := battleChara.Position
			_, ok = positions[position]
			if ok {
				errors.Swrapf(common.ErrBattleDuplicatePosition, position)
			}

			charas[charaId] = charaId
			positions[position] = position
		}
	}

	if battleFormation.Npcs != nil {
		extraNpcs := []int32{}
		if battleType == static.BattleTypeYgg {
			yggNpcs := u.Yggdrasil.GetBattleNpc()
			for _, yggNpc := range yggNpcs {
				extraNpcs = append(extraNpcs, yggNpc.NpcId)
			}
		}

		for _, npc := range battleFormation.Npcs {
			charaId, err := u.checkBattleNpc(npc.NpcId, npc.NpcType, battleId, battleType, extraNpcs)
			if err != nil {
				return err
			}

			_, ok := charas[charaId]
			if ok {
				errors.Swrapf(common.ErrBattleDuplicateChara, charaId)
			}

			position := npc.Position
			_, ok = positions[position]
			if ok {
				errors.Swrapf(common.ErrBattleDuplicatePosition, position)
			}

			charas[charaId] = charaId
			positions[position] = position
		}
	}

	battleHero := battleFormation.BattleHero
	if battleHero != nil {
		heroId := battleHero.HeroId
		if heroId > 0 {
			hero, ok := u.HeroPack.GetHero(heroId)
			if !ok {
				return errors.Swrapf(common.ErrHeroNotFound, heroId)
			}

			for _, skillId := range battleHero.UseSkillIds {
				if skillId > 0 && hero.GetSkillLevel(skillId) <= 0 {
					return errors.Swrapf(common.ErrHeroHasNotSkill, heroId, skillId)
				}
			}
		}
	}

	return nil
}

func (u *User) checkBattleNpc(npcId, npcType, battleId, systemType int32, extraNPC []int32) (int32, error) {
	var charaId int32 = 0
	switch npcType {
	case static.BattleNpcTypeSystem:
		npcCfg, err := manager.CSV.Battle.GetBattleNPC(npcId)
		if err != nil {
			return charaId, err
		}
		charaId = npcCfg.CharaID

		if slice.SliceInt32HasEle(extraNPC, npcId) {
			return charaId, nil
		}

		battleCfg, err := manager.CSV.Battle.GetBattleConfig(battleId)
		if err != nil {
			return charaId, err
		}

		if !slice.SliceInt32HasEle(battleCfg.Npc, npcId) &&
			!slice.SliceInt32HasEle(battleCfg.ControlNpc, npcId) {
			return charaId, errors.Swrapf(common.ErrBattleInvalidNPC, npcId)
		}

	case static.BattleNpcTypePlayer:
		err := u.MercenaryCheckUserCount(systemType, npcId)
		if err != nil {
			return npcId, err
		}

		charaId = npcId
	default:
		return charaId, errors.Swrapf(common.ErrBattleInvalidNPC, npcId)
	}

	return charaId, nil
}
