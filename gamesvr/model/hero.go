package model

import (
	"shared/protobuf/pb"
)

type HeroAttendant struct {
	CharaId      int32 `json:"chara_id"`
	LastCalcAttr int32 `json:"last_calc_attr"`
}

type Hero struct {
	ID            int32                    `json:"id"`
	Skills        map[int32]int32          `json:"skills"`     //id -> 等级
	Attendants    map[int32]*HeroAttendant `json:"attendants"` //侍从, 序号 -> 角色id
	SkillItemUsed int32                    `json:"skill_item_used"`
	LastCalcAttr  int32                    `json:"last_calc_attr"`
	Power         int32                    `json:"power"`
}

type HeroPack struct {
	Heros  map[int32]*Hero `json:"heros"`
	Level  int32           `json:"level"`
	Exp    int32           `json:"exp"`
	NewOne bool            `json:"new_one"`
}

func NewHeroPack() *HeroPack {
	return &HeroPack{
		Heros:  map[int32]*Hero{},
		Level:  1,
		Exp:    0,
		NewOne: true,
	}
}

func NewHero(heroId int32) *Hero {
	hero := &Hero{
		ID:            heroId,
		Skills:        map[int32]int32{},
		Attendants:    map[int32]*HeroAttendant{},
		SkillItemUsed: 0,
		LastCalcAttr:  0,
	}
	return hero
}

//----------------------------------------
//HeroPack
//----------------------------------------
func (hp *HeroPack) GetHero(id int32) (*Hero, bool) {
	hero, ok := hp.Heros[id]

	return hero, ok
}

func (hp *HeroPack) AddHero(id int32) *Hero {
	hero, ok := hp.GetHero(id)
	if ok {
		return hero
	}

	hero = NewHero(id)

	hp.Heros[id] = hero

	return hero
}

func (hp *HeroPack) GetLevel() int32 {
	return hp.Level
}
func (hp *HeroPack) SetLevel(level int32) {
	hp.Level = level
}
func (hp *HeroPack) GetExp() int32 {
	return hp.Exp
}

func (hp *HeroPack) AddExp(exp int32) {
	hp.Exp = hp.Exp + exp
}

func (hp *HeroPack) SetNewOneFlag(newOne bool) {
	hp.NewOne = newOne
}

func (hp *HeroPack) VOHeroInfo() *pb.VOHeroInfo {
	data := &pb.VOHeroInfo{
		Heros:  make([]*pb.VOHero, 0, len(hp.Heros)),
		Level:  hp.Level,
		Exp:    hp.Exp,
		NewOne: hp.NewOne,
	}

	for _, hero := range hp.Heros {
		data.Heros = append(data.Heros, hero.VOHero())
	}

	return data
}

//----------------------------------------
//Hero
//----------------------------------------
func (h *Hero) SetSkillLevel(skillId, level int32) {
	h.Skills[skillId] = level
}

func (h *Hero) GetSkillLevel(skillId int32) int32 {
	level, ok := h.Skills[skillId]
	if !ok {
		return 0
	}

	return level
}

func (h *Hero) RecordSkillItemUsed(useCount int32) {
	h.SkillItemUsed = h.SkillItemUsed + useCount
}

func (h *Hero) GetSkillItemUsed() int32 {
	return h.SkillItemUsed
}

func (h *Hero) SetAttendant(slot, charaId int32) {
	h.Attendants[slot] = &HeroAttendant{
		CharaId:      charaId,
		LastCalcAttr: 0,
	}
}

func (h *Hero) RemoveAttendant(slot int32) {
	delete(h.Attendants, slot)
}

func (h *Hero) GetAttendantCharaId(slot int32) int32 {
	attendant, ok := h.Attendants[slot]
	if !ok {
		return 0
	}

	return attendant.CharaId
}

func (h *Hero) UpdateLastAttr(heroAttr int32, attrs map[int32]int32) {
	h.LastCalcAttr = heroAttr
	for slot, attr := range attrs {
		attendant, ok := h.Attendants[slot]
		if !ok {
			continue
		}

		attendant.LastCalcAttr = attr
	}
}

func (h *Hero) GetPower() int32 {
	return h.Power
}

func (h *Hero) SetPower(power int32) {
	h.Power = power
}

func (h *Hero) VOHero() *pb.VOHero {
	data := &pb.VOHero{
		HeroId:        h.ID,
		Skills:        make([]*pb.VOHeroSkill, 0, len(h.Skills)),
		Attendants:    make([]*pb.VOHeroAttendant, 0, len(h.Attendants)),
		LastCalcAttr:  h.LastCalcAttr,
		SkillItemUsed: h.SkillItemUsed,
	}

	for skillId, level := range h.Skills {
		voHeroSkill := &pb.VOHeroSkill{
			SkillId: skillId,
			Level:   level,
		}

		data.Skills = append(data.Skills, voHeroSkill)
	}

	for slot, attendant := range h.Attendants {
		voHeroAttendant := &pb.VOHeroAttendant{
			Slot:         slot,
			CharaId:      attendant.CharaId,
			LastCalcAttr: attendant.LastCalcAttr,
		}

		data.Attendants = append(data.Attendants, voHeroAttendant)
	}

	return data
}
