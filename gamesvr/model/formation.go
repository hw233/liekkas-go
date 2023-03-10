package model

import (
	"shared/csv/static"
	"shared/protobuf/pb"
)

const (
	BattleFormationTypeUnkown  = 0
	BattleFormationTypeExplore = 1
	BattleFormationTypeElite   = 4
	BattleFormationTypeTower   = 5
	BattleFormationTypeYGG     = 6
	BattleFormationTypeChaos   = 10
)

func BattleTypeToFormationType(battleType int) int {
	switch battleType {
	case static.BattleTypeExplore,
		static.BattleTypeExploreResource,
		static.BattleTypeExploreMonster:
		return BattleFormationTypeExplore
	case static.BattleTypeExploreElite:
		return BattleFormationTypeElite
	case static.BattleTypeTower:
		return BattleFormationTypeTower
	case static.BattleTypeYgg:
		return BattleFormationTypeYGG
	case static.BattleTypeGate:
		return BattleFormationTypeChaos
	default:
		return BattleFormationTypeUnkown
	}
}

type CharaFormationPos struct {
	CharaId  int32 `json:"chara_id"`
	Position int32 `json:"position"`
}

func (cfp *CharaFormationPos) ParseFromVOBattleChara(voData *pb.VOBattleCharacter) {
	cfp.CharaId = voData.CharacterId
	cfp.Position = voData.Position
}

func (cfp *CharaFormationPos) VOBattleChara() *pb.VOBattleCharacter {
	return &pb.VOBattleCharacter{
		CharacterId: cfp.CharaId,
		Position:    cfp.Position,
		HpPercent:   10000,
	}
}

type NpcFormationPos struct {
	NpcId    int32 `json:"npc_id"`
	NpcType  int32 `json:"npc_type"`
	Position int32 `json:"position"`
}

func (nfp *NpcFormationPos) ParseFromVOBattleNPC(voData *pb.VOBattleNPC) {
	nfp.NpcId = voData.NpcId
	nfp.NpcType = voData.NpcType
	nfp.Position = voData.Position
}

func (nfp *NpcFormationPos) VOBattleNPC() *pb.VOBattleNPC {
	return &pb.VOBattleNPC{
		NpcId:     nfp.NpcId,
		NpcType:   nfp.NpcType,
		Position:  nfp.Position,
		HpPercent: 10000,
	}
}

type HeroFormationInfo struct {
	HeroId     int32   `json:"hero_id"`
	UsedSkills []int32 `json:"used_skills"`
}

func (hfi *HeroFormationInfo) ParseFromVOBattleHero(voData *pb.VOBattleHero) {
	hfi.HeroId = voData.HeroId
	hfi.UsedSkills = make([]int32, 0, len(voData.UseSkillIds))
	hfi.UsedSkills = append(hfi.UsedSkills, voData.UseSkillIds...)
}

func (hfi *HeroFormationInfo) VOBattleHero() *pb.VOBattleHero {
	battleHero := &pb.VOBattleHero{
		HeroId:      hfi.HeroId,
		UseSkillIds: make([]int32, 0, len(hfi.UsedSkills)),
	}

	battleHero.UseSkillIds = append(battleHero.UseSkillIds, hfi.UsedSkills...)

	return battleHero
}

type BattleSequence struct {
	UnitOperations []*BattleUnitOperation `json:"unit_operations"`
}

type BattleUnitOperation struct {
	UnitId    int32 `json:"unit_id"`
	UnitType  int32 `json:"unit_type"`
	UnitSkill int32 `json:"unit_skill"`
}

func (bs *BattleSequence) ParseFromVOBattleSequence(voData *pb.VOBattleSequence) {
	if voData.Units == nil {
		bs.UnitOperations = []*BattleUnitOperation{}
		return
	}

	bs.UnitOperations = make([]*BattleUnitOperation, 0, len(voData.Units))
	for _, voUnitOp := range voData.Units {
		unitOp := &BattleUnitOperation{
			UnitId:    voUnitOp.UnitId,
			UnitType:  voUnitOp.UnitType,
			UnitSkill: voUnitOp.UnitSkillId,
		}

		bs.UnitOperations = append(bs.UnitOperations, unitOp)
	}
}

func (bs *BattleSequence) VOBattleSequence() *pb.VOBattleSequence {
	voData := &pb.VOBattleSequence{
		Units: make([]*pb.VOBattleUnit, 0, len(bs.UnitOperations)),
	}

	for _, unitOp := range bs.UnitOperations {
		voOp := &pb.VOBattleUnit{
			UnitId:      unitOp.UnitId,
			UnitType:    unitOp.UnitType,
			UnitSkillId: unitOp.UnitSkill,
		}

		voData.Units = append(voData.Units, voOp)
	}

	return voData
}

type Formation struct {
	Id              int32                `json:"id"`
	FormationType   int32                `json:"formation_type"`
	HeroInfo        *HeroFormationInfo   `json:"hero_info"`
	CharaPositions  []*CharaFormationPos `json:"chara_positions"`
	NpcPositions    []*NpcFormationPos   `json:"npc_positions"`
	BattleSequences []*BattleSequence    `json:"battle_sequences"`
}

func (f *Formation) ParseFromVOBattleFormation(voData *pb.VOBattleFormation) {
	battleHero := voData.BattleHero
	if battleHero != nil {
		heroInfo := &HeroFormationInfo{}
		heroInfo.ParseFromVOBattleHero(battleHero)
		f.HeroInfo = heroInfo
	} else {
		f.HeroInfo = nil
	}

	if voData.BattleCharacters != nil {
		f.CharaPositions = make([]*CharaFormationPos, 0, len(voData.BattleCharacters))
		for _, battleChara := range voData.BattleCharacters {
			charaPos := &CharaFormationPos{}
			charaPos.ParseFromVOBattleChara(battleChara)
			f.CharaPositions = append(f.CharaPositions, charaPos)
		}
	} else {
		f.CharaPositions = []*CharaFormationPos{}
	}

	if voData.Npcs != nil {
		f.NpcPositions = make([]*NpcFormationPos, 0, len(voData.Npcs))
		for _, battleNpc := range voData.Npcs {
			npcPos := &NpcFormationPos{}
			npcPos.ParseFromVOBattleNPC(battleNpc)
			f.NpcPositions = append(f.NpcPositions, npcPos)
		}
	} else {
		f.NpcPositions = []*NpcFormationPos{}
	}

	if voData.BattleSequences != nil {
		f.BattleSequences = make([]*BattleSequence, 0, len(voData.BattleSequences))
		for _, voSeq := range voData.BattleSequences {
			seq := &BattleSequence{}
			seq.ParseFromVOBattleSequence(voSeq)
			f.BattleSequences = append(f.BattleSequences, seq)
		}
	}
}

func (f *Formation) VOFormation() *pb.VOFormation {
	voFormation := &pb.VOFormation{
		Id:            f.Id,
		FormationType: f.FormationType,
	}
	battleFormation := &pb.VOBattleFormation{
		BattleCharacters: make([]*pb.VOBattleCharacter, 0, len(f.CharaPositions)),
		Npcs:             make([]*pb.VOBattleNPC, 0, len(f.NpcPositions)),
		BattleSequences:  make([]*pb.VOBattleSequence, 0, len(f.BattleSequences)),
	}
	if f.HeroInfo != nil {
		battleFormation.BattleHero = f.HeroInfo.VOBattleHero()
	}

	for _, charaPos := range f.CharaPositions {
		battleFormation.BattleCharacters = append(battleFormation.BattleCharacters, charaPos.VOBattleChara())
	}

	for _, npcPos := range f.NpcPositions {
		battleFormation.Npcs = append(battleFormation.Npcs, npcPos.VOBattleNPC())
	}

	for _, battleSequence := range f.BattleSequences {
		battleFormation.BattleSequences = append(battleFormation.BattleSequences, battleSequence.VOBattleSequence())
	}

	voFormation.BattleFormation = battleFormation

	return voFormation
}

type FormationInfo struct {
	Formations map[int32]map[int32]*Formation `json:"formations"`
}

func NewFormationInfo() *FormationInfo {
	return &FormationInfo{
		Formations: map[int32]map[int32]*Formation{},
	}
}

func (fi *FormationInfo) SetFormation(id, formationType int32, voFormation *pb.VOBattleFormation) {
	formations, ok := fi.Formations[formationType]
	if !ok {
		formations = map[int32]*Formation{}
		fi.Formations[formationType] = formations
	}

	formation := &Formation{
		Id:            id,
		FormationType: formationType,
	}
	formation.ParseFromVOBattleFormation(voFormation)
	formations[id] = formation
}

func (fi *FormationInfo) GetFormation(id, formationType int32) (*Formation, bool) {
	formations, ok := fi.Formations[formationType]
	if !ok {
		return nil, false
	}

	formation, ok := formations[id]

	return formation, ok
}

func (fi *FormationInfo) VOFormations() []*pb.VOFormation {
	voData := []*pb.VOFormation{}
	for _, formations := range fi.Formations {
		for _, formation := range formations {
			voData = append(voData, formation.VOFormation())
		}
	}

	return voData
}
