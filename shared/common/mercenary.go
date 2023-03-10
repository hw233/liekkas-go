package common

import "shared/protobuf/pb"

type CharacterData struct {
	ID         int32           `json:"id"`
	Level      int32           `json:"level"`
	Star       int32           `json:"star"`
	Stage      int32           `json:"stage"`
	Skills     map[int32]int32 `json:"skills"`
	Power      int32           `json:"power"`
	Rarity     int32           `json:"rarity"`
	Equipments []*Equipment    `json:"equipments"`
	WorldItem  *WorldItem      `json:"world_item"`
}

func NewCharacterData(id, level, star, stage, power, rarity int32, skills map[int32]int32, equips []*Equipment, w *WorldItem) *CharacterData {
	return &CharacterData{
		ID:         id,
		Level:      level,
		Star:       star,
		Stage:      stage,
		Skills:     skills,
		Power:      power,
		Rarity:     rarity,
		Equipments: equips,
		WorldItem:  w,
	}
}

// ----------redis----------------
// 单个角色需要存储的数据
type MercenaryCharacter struct {
	Id    int32 `json:"id"`
	Level int32 `json:"level"`
	Star  int32 `json:"star"`
	Power int32 `json:"power"`
}

func NewMercenaryCharacter(id, level, star, power int32) *MercenaryCharacter {
	return &MercenaryCharacter{
		Id:    id,
		Level: level,
		Star:  star,
		Power: power,
	}
}

type MercenaryData struct {
	Uid      int64          `json:"uid"`
	Name     string         `json:"name"`
	Relation int8           `json:"relation"`
	Data     *CharacterData `json:"data"`
}

func NewMercenaryData(uid int64, name string, relation int8, c *CharacterData) *MercenaryData {
	return &MercenaryData{
		Uid:      uid,
		Name:     name,
		Relation: relation,
		Data:     c,
	}
}

type MercenarySend struct {
	IsCancel    bool   `json:"is_cancel"`
	CharacterId int32  `json:"character_id"`
	Uid         int64  `json:"uid"`
	SendTime    int64  `json:"send_time"`
	Name        string `json:"name"`
}

func NewMercenarySend(uid int64, name string, characterId int32) *MercenarySend {
	return &MercenarySend{
		IsCancel:    false,
		CharacterId: characterId,
		Uid:         uid,
		Name:        name,
	}
}

type MercenarySendRecord struct {
	UserTo      int64 `json:"user_to"`
	CharacterId int32 `json:"character_id"`
}

func NewMercenarySendRecord(uid int64, cid int32) *MercenarySendRecord {
	return &MercenarySendRecord{
		UserTo:      uid,
		CharacterId: cid,
	}
}

func (c *CharacterData) VOCharacter() *pb.VOUserCharacter {
	skills := make([]*pb.VOCharacterSkill, 0, len(c.Skills))
	for skillId, level := range c.Skills {
		skills = append(skills, &pb.VOCharacterSkill{
			SkillId: skillId,
			Level:   level,
		})
	}

	return &pb.VOUserCharacter{
		CharacterId: c.ID,
		Level:       c.Level,
		Star:        c.Star,
		Stage:       c.Stage,
		Skills:      skills,
		Power:       c.Power,
		Rarity:      c.Rarity,
	}
}
