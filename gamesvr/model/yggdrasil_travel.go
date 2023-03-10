package model

import (
	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/errors"
)

// YggdrasilInTravelInfo 当前旅行信息
type YggdrasilInTravelInfo struct {
	TravelAp     int32           `json:"explore_ap"`    // 本次探索剩余ap
	CharactersHp map[int32]int32 `json:"characters_hp"` // 携带角色的血量
}

func NewYggdrasilInTravelInfo() *YggdrasilInTravelInfo {
	return &YggdrasilInTravelInfo{
		TravelAp:     0,
		CharactersHp: map[int32]int32{},
	}

}

// AddHpForAll 给所有角色回血,num代表万分比
func (y *YggdrasilInTravelInfo) AddHpForAll(num int32) {
	for id, hp := range y.CharactersHp {
		if hp <= 0 {
			continue
		}
		if hp <= (10000 - num) {
			hp += num
		} else {
			hp = 10000
		}
		y.CharactersHp[id] = hp
	}
}

func (y *YggdrasilInTravelInfo) SetAllDead() {
	for k := range y.CharactersHp {
		y.CharactersHp[k] = 0
	}
}

func (y *YggdrasilInTravelInfo) AllDead() bool {
	for _, hp := range y.CharactersHp {
		if hp > 0 {
			return false
		}
	}
	return true
}

func (y *YggdrasilInTravelInfo) ContainCharacters(characters []*pb.VOBattleCharacter) error {
	for _, character := range characters {
		hp, ok := y.CharactersHp[character.CharacterId]
		if !ok {
			return errors.WrapTrace(errors.Swrapf(common.ErrYggdrasilCharacterNotExist, character.CharacterId))
		}
		if hp <= 0 {
			return errors.WrapTrace(errors.Swrapf(common.ErrYggdrasilCharacterNotExist, character.CharacterId))
		}
	}
	return nil
}
func (y *YggdrasilInTravelInfo) GetCharacterHp(characterId int32) (int32, error) {
	hp, ok := y.CharactersHp[characterId]
	if !ok {
		return 0, errors.WrapTrace(errors.Swrapf(common.ErrYggdrasilCharacterNotExist, characterId))
	}
	return hp, nil
}

func (y *YggdrasilInTravelInfo) SetCharacterHp(characterId, hp int32) {

	_, ok := y.CharactersHp[characterId]
	if ok {
		y.CharactersHp[characterId] = hp
	}
}

func (y *YggdrasilInTravelInfo) VOYggdrasilCharacterHp() []*pb.VOYggdrasilCharacter {
	Characters := make([]*pb.VOYggdrasilCharacter, 0, len(y.CharactersHp))

	for charId, hpPercent := range y.CharactersHp {
		Characters = append(Characters, &pb.VOYggdrasilCharacter{CharacterId: charId, HpPercent: hpPercent})
	}
	return Characters
}

func (y *YggdrasilInTravelInfo) VOYggdrasilTravelInfo() *pb.VOYggdrasilTravelInfo {
	return &pb.VOYggdrasilTravelInfo{
		Ap:         y.TravelAp,
		Characters: y.VOYggdrasilCharacterHp(),
		HelpCount:  0, //todo:建造帮助
	}

}
