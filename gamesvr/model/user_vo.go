package model

import (
	"context"
	"shared/protobuf/pb"
	"sort"
)

func (u *User) VOUserInfo() *pb.VOUserInfo {
	return &pb.VOUserInfo{
		UserId:    u.ID,
		Nickname:  u.Name,
		NameIndex: u.Info.NameIndex,
		TeamLevel: u.Info.Level.Value(),
		TeamExp:   u.Info.Exp.Value(),
		Avatar:    u.Info.Avatar,
		Frame:     u.Info.Frame,
		GuildId:   u.Guild.GuildID,
		GuildName: u.Guild.GuildName,
	}
}

func (u *User) VOUserInfoSimple() *pb.VOUserInfoSimple {
	return &pb.VOUserInfoSimple{
		UserId:    u.ID,
		Nickname:  u.Name,
		TeamLevel: u.Info.Level.Value(),
		Avatar:    u.Info.Avatar,
		NameIndex: u.Info.NameIndex,
	}
}

func (u *User) VOUserResource() *pb.VOUserResource {
	return &pb.VOUserResource{
		Money:           int64(u.Info.Gold.Value()),
		DiamondGift:     u.Info.DiamondGift.Value(),
		DiamondCash:     u.Info.DiamondCash.Value(),
		Energy:          u.Info.Energy.Value(),
		EnergyRefreshAt: u.Info.Energy.Last(),
	}
}

func (u *User) VOUserEquipment() []*pb.VOUserEquipment {
	return u.EquipmentPack.VOUserEquipment()
}

func (u *User) VOItemInfo() []*pb.VOItemInfo {
	return u.ItemPack.VOItemInfo()
}

func (u *User) VOUserCharacter() []*pb.VOUserCharacter {
	return u.CharacterPack.VOUserCharacter()
}

func (u *User) VOUserWorldItem() []*pb.VOUserWorldItem {
	return u.WorldItemPack.VOUserWorldItem()
}

func (u *User) VOManualInfo() []*pb.VOManualInfo {
	return u.ManualInfo.VOManualInfo()
}

func (u *User) VOQuestInfo() *pb.VOQuestInfo {
	return u.QuestPack.VOQuestInfo()
}

func (u *User) VOResourceResult() *pb.VOResourceResult {
	ret := u.RewardsResult.VOResourceResult()
	// 清除数据
	u.RewardsResult.Clear()
	return ret
}

func (u *User) VOMergedResourceResult() *pb.VOResourceResult {
	ret := u.RewardsResult.VOMergedResourceResult()
	// 清除数据
	u.RewardsResult.Clear()
	return ret
}

func (u *User) VOYggdrasilResourceResult() *pb.VOYggdrasilResourceResult {
	ret := u.Yggdrasil.VOYggdrasilResourceResult()
	return ret
}

func (u *User) VOYggdrasilEntityChange(ctx context.Context) *pb.VOYggdrasilEntityChange {
	ret := u.Yggdrasil.VOYggdrasilEntityChange(ctx)
	return ret
}

func (u *User) VOBattleCharacter(charaId int32) *pb.VOBattleCharacterProp {
	chara, _ := u.CharacterPack.Get(charaId)
	voBattleChara := &pb.VOBattleCharacterProp{
		Character:  chara.VOUserCharacter(),
		Equipments: make([]*pb.VOUserEquipment, 0, len(chara.Equipments)),
	}

	for _, equipmentId := range chara.Equipments {
		if equipmentId > 0 {
			equip, _ := u.EquipmentPack.Get(equipmentId)
			voBattleChara.Equipments = append(voBattleChara.Equipments, equip.VOUserEquipment())
		}
	}

	if chara.WorldItem > 0 {
		worldItem, _ := u.WorldItemPack.Get(chara.WorldItem)
		voBattleChara.WorldItem = worldItem.VOUserWorldItem()
	}

	return voBattleChara
}

// VOVisitingCard 玩家名片信息
func (u *User) VOVisitingCard() *pb.VOVisitingCard {
	u.RefreshCardShowCharacterCache()
	return NewOwnUserVisitingCard(u.ID, u.Name, u.Info, u.Guild).VOVisitingCard()

}

// RefreshCardShowCharacterCache 更新名片中角色展示的数据
func (u *User) RefreshCardShowCharacterCache() {
	characters := map[int32]*pb.VOUserCharacter{}
	if u.Info.CardShow.CharactersSet {
		for i := 1; i <= VisitingCardPosCount; i++ {
			pos := int32(i)
			characterId, ok := u.Info.CardShow.Characters[pos]
			if !ok {
				characters[pos] = nil
				continue
			}
			character, err := u.CharacterPack.Get(characterId)
			if err != nil {
				characters[pos] = nil
				continue
			}
			characters[pos] = character.VOUserCharacter()

		}

	} else {

		// 排序取前四个
		totalCharacters := make([]*Character, 0, len(*u.CharacterPack))
		for _, character := range *u.CharacterPack {
			totalCharacters = append(totalCharacters, character)
		}
		/**
		稀有度高>稀有度低
		等级大>等级小
		星级高>星级低
		阶段高>阶段低
		id小>大
		*/
		sort.Slice(totalCharacters, func(i, j int) bool {
			characteri := totalCharacters[i]
			characterj := totalCharacters[j]
			if characteri.Rarity != characterj.Rarity {
				return characterj.Rarity < characteri.Rarity
			}
			if characteri.Level != characterj.Level {
				return characterj.Level < characteri.Level
			}
			if characteri.Star != characterj.Star {
				return characterj.Star < characteri.Star
			}
			if characteri.Stage != characterj.Stage {
				return characterj.Stage < characteri.Stage
			}
			if characteri.ID != characterj.ID {
				return characteri.ID < characterj.ID
			}
			return true
		})
		for i := 1; i <= VisitingCardPosCount; i++ {
			pos := int32(i)

			if i > len(totalCharacters) {
				characters[pos] = nil
			} else {
				characters[pos] = totalCharacters[i-1].VOUserCharacter()
			}
		}
	}
	u.Info.CardShow.CharactersCache = characters
}
