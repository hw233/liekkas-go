package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/utility/errors"
	"shared/utility/lua_state"

	lua "github.com/yuin/gopher-lua"
)

// 对外接口，计算user的总战力
func (u *User) GetUserPower() int32 {
	var result int32
	for _, character := range *u.CharacterPack {
		result += character.Power
	}

	for _, hero := range u.HeroPack.Heros {
		result += hero.Power
	}
	u.Info.Power = result
	return result
}

func (u *User) CalUserCombatPower() (int32, error) {

	manager.LuaState.Lock()
	defer manager.LuaState.Unlock()

	power, err := calUserCombatPower(manager.LuaState, u)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}
	return power, nil
}

// 对外接口，计算character的战力
func (u *User) CalCharacterCombatPower(charac *Character) (int32, error) {
	manager.LuaState.Lock()
	defer manager.LuaState.Unlock()

	characCombatData, err := GetCharacterCombatData(manager.LuaState, charac, u)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}

	characPower, err := manager.LuaState.CalCharacterCombatPower(characCombatData)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}
	// fmt.Println("=========characPower: ", characPower)
	return characPower, nil
}

func (u *User) CalHeroCombatPower(hero *Hero) (int32, error) {
	manager.LuaState.Lock()
	defer manager.LuaState.Unlock()

	heroCombatData, err := GetHeroCombatData(manager.LuaState, hero, u)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}

	heroPower, err := manager.LuaState.CalHeroCombatPower(heroCombatData)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}

	return heroPower, nil
}

func calUserCombatPower(L *lua_state.LuaMatcher, u *User) (int32, error) {
	var result int32

	// 逐个计算角色的战力
	for _, charac := range *u.CharacterPack {
		// fmt.Printf("========characId: %v, eid: %v, wid: %d\n", charac.ID, charac.Equipments, charac.WorldItem)
		characCombatData, err := GetCharacterCombatData(L, charac, u)
		if err != nil {
			return 0, errors.WrapTrace(err)
		}

		characPower, err := L.CalCharacterCombatPower(characCombatData)
		if err != nil {
			return 0, errors.WrapTrace(err)
		}
		charac.Power = characPower
		// fmt.Println("charaId: ", charac.ID, "power: ", characPower)
		result += characPower
	}

	for _, hero := range u.HeroPack.Heros {
		heroCombatData, err := GetHeroCombatData(L, hero, u)
		if err != nil {
			return 0, errors.WrapTrace(err)
		}

		heroPower, err := L.CalHeroCombatPower(heroCombatData)
		if err != nil {
			return 0, errors.WrapTrace(err)
		}

		// fmt.Println("heroPower: ", heroPower)

		result += heroPower
	}

	return result, nil
}

func GetHeroCombatData(L *lua_state.LuaMatcher, h *Hero, u *User) (lua.LValue, error) {

	// 计算出至尊所有侍从的info
	characInfos := L.NewTable()

	for _, attend := range h.Attendants {
		charac, err := u.CharacterPack.Get(attend.CharaId)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		characInfo, err := GetCharacterCombatData(L, charac, u)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		characInfos.Append(characInfo)
	}

	heroInfo := setHeroData(L, h, characInfos, u.HeroPack.GetLevel())

	return heroInfo, nil
}

func GetCharacterCombatData(L *lua_state.LuaMatcher, charac *Character, u *User) (lua.LValue, error) {
	if charac.ID == 0 {
		return nil, errors.Swrapf(common.ErrCharacterNotFound, charac.ID)
	}

	equips := make([]*common.Equipment, 0, CharacterEquipmentCounts)
	for _, eid := range charac.Equipments {
		if eid != 0 {
			equip, err := u.EquipmentPack.Get(eid)
			if err != nil {
				return nil, errors.WrapTrace(err)
			}
			equips = append(equips, equip)
		}
	}
	worldItem := &common.WorldItem{}
	if charac.WorldItem != 0 {
		wItem, err := u.WorldItemPack.Get(charac.WorldItem)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		worldItem = wItem
	}

	characInfo := setCharacterData(L, charac, equips, worldItem, charac.Level)

	return characInfo, nil
}

// 把Charater转换成LValue
func getCharacterInfo(L *lua_state.LuaMatcher, chara *Character) lua.LValue {

	if chara.ID == 0 {
		return lua.LNil
	}

	charaLua := L.NewTable()
	charaLua.RawSetString("characterId", lua.LNumber(chara.ID))
	charaLua.RawSetString("level", lua.LNumber(chara.Level))
	charaLua.RawSetString("star", lua.LNumber(chara.Star))
	charaLua.RawSetString("stage", lua.LNumber(chara.Stage))
	charaLua.RawSetString("createAt", lua.LNumber(chara.CTime))
	charaLua.RawSetString("heroId", lua.LNumber(chara.HeroId))
	charaLua.RawSetString("heroPos", lua.LNumber(chara.HeroSlot))
	charaLua.RawSetString("power", lua.LNumber(chara.Power))
	charaLua.RawSetString("canYggdrasilTime", lua.LNumber(chara.CanYggdrasilTime))
	charaLua.RawSetString("rarity", lua.LNumber(chara.Rarity))

	skillList := L.NewTable()

	for skillId, level := range chara.Skills {
		sk := L.NewTable()
		sk.RawSetString("skillId", lua.LNumber(skillId))
		sk.RawSetString("level", lua.LNumber(level))
		skillList.Append(sk)
	}

	charaLua.RawSetString("skills", skillList)
	return charaLua
}

// 把Equipment转换成LValue
func getEquitmentInfo(L *lua_state.LuaMatcher, equips []*common.Equipment) lua.LValue {
	if len(equips) == 0 {
		return lua.LNil
	}
	equipsLua := L.NewTable()

	for _, equip := range equips {
		equipLua := L.NewTable()

		equipLua.RawSetString("equipmentUId", lua.LNumber(equip.ID))
		equipLua.RawSetString("equipmentId", lua.LNumber(equip.EID))
		equipLua.RawSetString("characterId", lua.LNumber(equip.CID))
		equipLua.RawSetString("exp", lua.LNumber(equip.EXP.Value()))
		equipLua.RawSetString("stage", lua.LNumber(equip.Stage.Value()))
		equipLua.RawSetString("level", lua.LNumber(equip.Level.Value()))
		equipLua.RawSetString("characterId", lua.LNumber(equip.CID))
		equipLua.RawSetString("lockFlag", lua.LBool(equip.IsLocked))
		equipLua.RawSetString("camp", lua.LNumber(int32(equip.Camp)))
		equipLua.RawSetString("campRandom", lua.LNumber(0))
		equipLua.RawSetString("createAt", lua.LNumber(equip.CTime))
		equipLua.RawSetString("recastCamp", lua.LNumber(int32(equip.RecastCamp)))

		attrs := L.NewTable()
		for i, ea := range equip.Attrs {
			attr := L.NewTable()
			attr.RawSetString("attr", lua.LNumber(int32(ea.Attr)))
			attr.RawSetString("value", lua.LNumber(ea.Value))
			attrs.Insert(i, attr)
		}
		equipLua.RawSetString("attrs", attrs)

		equipsLua.Append(equipLua)
	}

	return equipsLua
}

// 把WorldItem转为LValue
func getWorldItemInfo(L *lua_state.LuaMatcher, worldItem *common.WorldItem) lua.LValue {
	if worldItem.ID == 0 {
		return lua.LNil
	}

	worldItemLua := L.NewTable()

	worldItemLua.RawSetString("worldItemUId", lua.LNumber(worldItem.ID))
	worldItemLua.RawSetString("worldItemId", lua.LNumber(worldItem.WID))
	worldItemLua.RawSetString("characterId", lua.LNumber(worldItem.CID))
	worldItemLua.RawSetString("exp", lua.LNumber(worldItem.EXP.Value()))
	worldItemLua.RawSetString("star", lua.LNumber(worldItem.Stage.Value()))
	worldItemLua.RawSetString("level", lua.LNumber(worldItem.Level.Value()))
	worldItemLua.RawSetString("lockFlag", lua.LBool(worldItem.IsLock))
	worldItemLua.RawSetString("createAt", lua.LNumber(worldItem.CTime))

	return worldItemLua

}

// 把至尊结构体转换为LValue
func getHeroInfo(L *lua_state.LuaMatcher, h *Hero) lua.LValue {
	heroLua := L.NewTable()

	heroLua.RawSetString("heroId", lua.LNumber(h.ID))

	skillList := L.NewTable()
	for skillId, level := range h.Skills {
		sk := L.NewTable()
		sk.RawSetString("skillId", lua.LNumber(skillId))
		sk.RawSetString("level", lua.LNumber(level))
		skillList.Append(sk)
	}

	heroLua.RawSetString("skills", skillList)

	attendList := L.NewTable()
	for slot, a := range h.Attendants {
		attend := L.NewTable()
		attend.RawSetString("slot", lua.LNumber(slot))
		attend.RawSetString("charaId", lua.LNumber(a.CharaId))
		attend.RawSetString("lastCalcAttr", lua.LNumber(a.LastCalcAttr))
		attendList.Append(attend)
	}

	heroLua.RawSetString("attendants", attendList)

	heroLua.RawSetString("skillItemUsed", lua.LNumber(h.SkillItemUsed))
	heroLua.RawSetString("lastCalcAttr", lua.LNumber(h.LastCalcAttr))

	return heroLua
}

// 把角色，装备以及世界道具转换成一个LValue
func setCharacterData(L *lua_state.LuaMatcher, c *Character, e []*common.Equipment, w *common.WorldItem, level int32) lua.LValue {

	cLua := getCharacterInfo(L, c)
	eLua := getEquitmentInfo(L, e)
	wLua := getWorldItemInfo(L, w)

	data := L.NewTable()

	data.RawSetString("Character", cLua)
	data.RawSetString("Equips", eLua)
	data.RawSetString("Item", wLua)
	data.RawSetString("Level", lua.LNumber(level))
	return data
}

func setHeroData(L *lua_state.LuaMatcher, h *Hero, charaCombatData lua.LValue, level int32) lua.LValue {
	hLua := getHeroInfo(L, h)

	data := L.NewTable()

	data.RawSetString("Hero", hLua)
	data.RawSetString("Pledges", charaCombatData)
	data.RawSetString("Level", lua.LNumber(level))

	return data
}
