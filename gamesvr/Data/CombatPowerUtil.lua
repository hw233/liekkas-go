if CS then
    --客户端
    LuaState = CS.GameApplication.LuaState.Instance
    -- 预加载脚本
    LuaState:LoadScript("Data/Config/cfg_combat_power")
    LuaState:LoadScript("Data/Config/cfg_passsive_attributes")
    LuaState:LoadScript("Data/Config/cfg_character")
    LuaState:LoadScript("Data/Config/cfg_character_star")
    LuaState:LoadScript("Data/Config/cfg_character_stage")
    LuaState:LoadScript("Data/Config/cfg_character_data")
    LuaState:LoadScript("Data/Config/cfg_equip_attributes_data")
    LuaState:LoadScript("Data/Config/cfg_world_item_attributes_data")
    LuaState:LoadScript("Data/Config/cfg_equip_data")
    LuaState:LoadScript("Data/Config/cfg_character_skill")
    LuaState:LoadScript("Data/Config/const_skill_type")
    LuaState:LoadScript("Data/Config/cfg_hero")
    LuaState:LoadScript("Data/Config/cfg_hero_data")
    LuaState:LoadScript("Data/Config/cfg_hero_skill")
    LuaState:LoadScript("Data/Config/cfg_equip_skill")
    LuaState:LoadScript("Data/Config/cfg_world_item_data")
    LuaState:LoadScript("Data/Util/util_cfg_character_star")
    LuaState:LoadScript("Data/Util/util_cfg_character_stage")
    LuaState:LoadScript("Data/Util/util_cfg_character_data")
    LuaState:LoadScript("Data/Util/util_cfg_equip_attributes_data")
    LuaState:LoadScript("Data/Util/util_cfg_world_item_attributes_data")
    LuaState:LoadScript("Data/Util/util_cfg_equip_data")
    LuaState:LoadScript("Data/Util/util_cfg_character_skill")
    LuaState:LoadScript("Data/Util/util_cfg_hero_data")
    LuaState:LoadScript("Data/Util/util_cfg_hero_skill")
    LuaState:LoadScript("Data/Util/util_cfg_world_item_data")
else
    -- 服务器
end

-- 四舍五入保留整数部分
local function Round(n)
    return math.floor(n + 0.5)
end

-- 保留三位小数
local function RoundToThree(n)
    if n < 0 then
        return math.ceil(n * 1000) * 0.001
    end
    return math.floor(n * 1000) * 0.001
end

-- 人物基础
local function chara_basic(chara)
    local util_cfg_character_data = require "Data/Util/util_cfg_character_data"
    local util_cfg_character_star = require "Data/Util/util_cfg_character_star"
    local util_cfg_character_stage = require "Data/Util/util_cfg_character_stage"
    local id = chara.characterId
    local level = chara.level
    local star = chara.star
    local stage = chara.stage
    local cfg_chara = util_cfg_character_data.get(id)
    local cfg_star = util_cfg_character_star.get(id, star)
    local cfg_stage = util_cfg_character_stage.get(id, stage)
    local basic = {}
    local level_ratio = level - 1
    basic.hpMax = level_ratio * cfg_star.hpRatio + cfg_chara.hpMax + cfg_stage.hpMax + cfg_star.hpMax
    basic.phyAtk = level_ratio * cfg_star.phyAtkRatio + cfg_chara.phyAtk + cfg_stage.phyAtk + cfg_star.phyAtk
    basic.magAtk = level_ratio * cfg_star.magAtkRatio + cfg_chara.magAtk + cfg_stage.magAtk + cfg_star.magAtk
    basic.phyDfs = level_ratio * cfg_star.phyDfsRatio + cfg_chara.phyDfs + cfg_stage.phyDfs + cfg_star.phyDfs
    basic.magDfs = level_ratio * cfg_star.magDfsRatio + cfg_chara.magDfs + cfg_stage.magDfs + cfg_star.magDfs
    basic.critAtkRatio = cfg_chara.critAtkRatio + cfg_stage.critAtkRatio
    basic.critDfsRatio = cfg_chara.critDfsRatio + cfg_stage.critDfsRatio
    basic.critAtkValue = cfg_chara.critAtkValue + cfg_stage.critAtkValue
    basic.critDfsValue = cfg_chara.critDfsValue + cfg_stage.critDfsValue
    basic.hitRateValue = cfg_chara.hitRateValue + cfg_stage.hitRateValue + cfg_chara.hitRateBasic
    basic.evadeValue = cfg_chara.evadeValue + cfg_stage.evadeValue
    basic.normalAtkUp = cfg_stage.normalAtkUp
    basic.normalDfsUp = cfg_stage.normalDfsUp
    basic.skillAtkUp = cfg_stage.skillAtkUp
    basic.skillDfsUp = cfg_stage.skillDfsUp
    basic.ultraAtkUp = cfg_stage.ultraAtkUp
    basic.ultraDfsUp = cfg_stage.ultraDfsUp
    basic.phyAtkUp = cfg_stage.skillPhyAtkUp
    basic.phyDfsUp = cfg_stage.skillPhyDfsUp
    basic.magAtkUp = cfg_stage.skillMagAtkUp
    basic.magDfsUp = cfg_stage.skillMagDfsUp
    basic.cureUp = 0
    basic.healUp = 0
    basic.finalDmg = 0
    basic.phyPen = 0
    basic.magPen = 0
    return basic
end

-- 对应配置表const_equip_rand_attributer,若有变动，需要更新
local equip_attrs_key = {
    "hpMaxPercent",
    "phyAtkPercent",
    "magAtkPercent",
    "phyDfsPercent",
    "magDfsPercent",
    "critAtkRatio",
    "critDfsRatio",
    "critAtkValue",
    "critDfsValue",
    "healUp",
    "phyDamAdd",
    "phyDamReduce",
    "magDamAdd",
    "magDamReduce",
    "phyPen",
    "magPen"
}

-- 装备词条
local function equip_attrs(attrs)
    attrs = attrs or {}
    local addition = {}
    for _, key in pairs(equip_attrs_key) do
        addition[key] = 0
    end
    for _, v in pairs(attrs) do
        local key = equip_attrs_key[v.attr]
        addition[key] = v.value
    end
    return addition
end

-- 装备基础
local function equip_basic(equip, camp)
    local attrs = equip_attrs(equip.attrs)
    local util_cfg_equip_data = require "Data/Util/util_cfg_equip_data"
    local util_cfg_equip_attributes_data = require "Data/Util/util_cfg_equip_attributes_data"
    local cfg_passsive_attributes = require "Data/Config/cfg_passsive_attributes"
    local id = equip.equipmentId
    local stage = equip.stage
    local level = equip.level
    local cfg_equip = util_cfg_equip_data.get(id)
    local cfg_basic = util_cfg_equip_attributes_data.get(id, 0)
    local cfg_extra = util_cfg_equip_attributes_data.get_extra(id, stage)
    local attribute_id = util_cfg_equip_data.get_attribute_id(id, level, stage)
    local cfg_passive = cfg_passsive_attributes[attribute_id]
    local level_ratio = level - 1
    local basic = {}
    local camp_addition = 0
    if equip.camp == camp then
        camp_addition = cfg_equip.campAddition
    end
    local percent = 1 + camp_addition * 0.001
    -- 装备属性数值加成
    basic.hpMax =
        (level_ratio * cfg_extra.hpExtraRatio + cfg_basic.hpMaxExtra + cfg_extra.hpMaxExtra) * percent +
        cfg_passive.hpMaxExtra
    basic.phyAtk =
        (level_ratio * cfg_extra.phyAtkExtraRatio + cfg_basic.phyAtkExtra + cfg_extra.phyAtkExtra) * percent +
        cfg_passive.phyAtkExtra
    basic.magAtk =
        (level_ratio * cfg_extra.magAtkExtraRatio + cfg_basic.magAtkExtra + cfg_extra.magAtkExtra) * percent +
        cfg_passive.magAtkExtra
    basic.phyDfs =
        (level_ratio * cfg_extra.phyDfsExtraRatio + cfg_basic.phyDfsExtra + cfg_extra.phyDfsExtra) * percent +
        cfg_passive.phyDfsExtra
    basic.magDfs =
        (level_ratio * cfg_extra.magDfsExtraRatio + cfg_basic.magDfsExtra + cfg_extra.magDfsExtra) * percent +
        cfg_passive.magDfsExtra
    basic.critAtkRatio = cfg_basic.critAtkRatio + cfg_extra.critAtkRatio + cfg_passive.critAtkRatio + attrs.critAtkRatio
    basic.critDfsRatio = cfg_basic.critDfsRatio + cfg_extra.critDfsRatio + cfg_passive.critDfsRatio + attrs.critDfsRatio
    basic.critAtkValue = cfg_basic.critAtkValue + cfg_extra.critAtkValue + cfg_passive.critAtkValue + attrs.critAtkValue
    basic.critDfsValue = cfg_basic.critDfsValue + cfg_extra.critDfsValue + cfg_passive.critDfsValue + attrs.critDfsValue
    basic.hitRateValue = cfg_basic.hitRateValue + cfg_extra.hitRateValue + cfg_passive.hitRateValue
    basic.evadeValue = cfg_basic.evadeValue + cfg_extra.evadeValue + cfg_passive.evadeValue
    basic.normalAtkUp = cfg_basic.normalAtkUp + cfg_extra.normalAtkUp + cfg_passive.normalAtkUp
    basic.normalDfsUp = cfg_basic.normalDfsUp + cfg_extra.normalDfsUp + cfg_passive.normalDfsUp
    basic.skillAtkUp = cfg_basic.skillAtkUp + cfg_extra.skillAtkUp + cfg_passive.skillAtkUp
    basic.skillDfsUp = cfg_basic.skillDfsUp + cfg_extra.skillDfsUp + cfg_passive.skillDfsUp
    basic.ultraAtkUp = cfg_basic.ultraAtkUp + cfg_extra.ultraAtkUp + cfg_passive.ultraAtkUp
    basic.ultraDfsUp = cfg_basic.ultraDfsUp + cfg_extra.ultraDfsUp + cfg_passive.ultraDfsUp
    basic.phyAtkUp = cfg_basic.phyDamAdd + cfg_extra.phyDamAdd + cfg_passive.skillPhyAtkUp + attrs.phyDamAdd
    basic.phyDfsUp = cfg_basic.phyDamReduce + cfg_extra.phyDamReduce + cfg_passive.skillPhyDfsUp + attrs.phyDamReduce
    basic.magAtkUp = cfg_basic.magDamAdd + cfg_extra.magDamAdd + cfg_passive.skillMagAtkUp + attrs.magDamAdd
    basic.magDfsUp = cfg_basic.magDamReduce + cfg_extra.magDamReduce + cfg_passive.skillMagDfsUp + attrs.magDamReduce
    basic.cureUp = cfg_basic.cureUp + cfg_extra.cureUp + cfg_passive.cureUp
    basic.healUp = cfg_basic.healUp + cfg_extra.healUp + cfg_passive.healUp + attrs.healUp
    basic.finalDmg = cfg_passive.finalDmg
    basic.phyPen = cfg_basic.phyPen + cfg_extra.phyPen + attrs.phyPen
    basic.magPen = cfg_basic.magPen + cfg_extra.magPen + attrs.magPen
    -- 装备属性百分比加成（相对于角色自身）
    local percent = {}
    percent.hpMaxPercent = cfg_extra.hpMaxPercent + cfg_passive.hpMaxPercent + attrs.hpMaxPercent
    percent.phyAtkPercent = cfg_extra.phyAtkPercent + cfg_passive.phyAtkPercent + attrs.phyAtkPercent
    percent.magAtkPercent = cfg_extra.magAtkPercent + cfg_passive.magAtkPercent + attrs.magAtkPercent
    percent.phyDfsPercent = cfg_extra.phyDfsPercent + cfg_passive.phyDfsPercent + attrs.phyDfsPercent
    percent.magDfsPercent = cfg_extra.magDfsPercent + cfg_passive.magDfsPercent + attrs.magDfsPercent
    return basic, percent
end

-- 世界道具基础
local function world_item_basic(item, camp)
    local util_cfg_world_item_data = require "Data/Util/util_cfg_world_item_data"
    local util_cfg_world_item_attributes_data = require "Data/Util/util_cfg_world_item_attributes_data"
    local cfg_passsive_attributes = require "Data/Config/cfg_passsive_attributes"
    local id = item.worldItemId
    local level = item.level
    local star = item.star
    local cfg_basic = util_cfg_world_item_attributes_data.get(id, 0)
    local cfg_extra = util_cfg_world_item_attributes_data.get_extra(id, star)
    local attribute_id = util_cfg_world_item_data.get_attribute_id(id, level, star)
    local cfg_passive = cfg_passsive_attributes[attribute_id]
    local level_ratio = level - 1
    local basic = {}
    -- 世界道具属性数值加成
    basic.hpMax =
        level_ratio * cfg_extra.hpExtraRatio + cfg_basic.hpMaxExtra + cfg_extra.hpMaxExtra + cfg_passive.hpMaxExtra
    basic.phyAtk =
        level_ratio * cfg_extra.phyAtkExtraRatio + cfg_basic.phyAtkExtra + cfg_extra.phyAtkExtra +
        cfg_passive.phyAtkExtra
    basic.magAtk =
        level_ratio * cfg_extra.magAtkExtraRatio + cfg_basic.magAtkExtra + cfg_extra.magAtkExtra +
        cfg_passive.magAtkExtra
    basic.phyDfs =
        level_ratio * cfg_extra.phyDfsExtraRatio + cfg_basic.phyDfsExtra + cfg_extra.phyDfsExtra +
        cfg_passive.phyDfsExtra
    basic.magDfs =
        level_ratio * cfg_extra.magDfsExtraRatio + cfg_basic.magDfsExtra + cfg_extra.magDfsExtra +
        cfg_passive.magDfsExtra
    basic.critAtkRatio = cfg_basic.critAtkRatio + cfg_extra.critAtkRatio + cfg_passive.critAtkRatio
    basic.critDfsRatio = cfg_basic.critDfsRatio + cfg_extra.critDfsRatio + cfg_passive.critDfsRatio
    basic.critAtkValue = cfg_basic.critAtkValue + cfg_extra.critAtkValue + cfg_passive.critAtkValue
    basic.critDfsValue = cfg_basic.critDfsValue + cfg_extra.critDfsValue + cfg_passive.critDfsValue
    basic.hitRateValue = cfg_basic.hitRateValue + cfg_extra.hitRateValue + cfg_passive.hitRateValue
    basic.evadeValue = cfg_basic.evadeValue + cfg_extra.evadeValue + cfg_passive.evadeValue
    basic.normalAtkUp = cfg_basic.normalAtkUp + cfg_extra.normalAtkUp + cfg_passive.normalAtkUp
    basic.normalDfsUp = cfg_basic.normalDfsUp + cfg_extra.normalDfsUp + cfg_passive.normalDfsUp
    basic.skillAtkUp = cfg_basic.skillAtkUp + cfg_extra.skillAtkUp + cfg_passive.skillAtkUp
    basic.skillDfsUp = cfg_basic.skillDfsUp + cfg_extra.skillDfsUp + cfg_passive.skillDfsUp
    basic.ultraAtkUp = cfg_basic.ultraAtkUp + cfg_extra.ultraAtkUp + cfg_passive.ultraAtkUp
    basic.ultraDfsUp = cfg_basic.ultraDfsUp + cfg_extra.ultraDfsUp + cfg_passive.ultraDfsUp
    basic.phyAtkUp = cfg_basic.skillPhyAtkUp + cfg_extra.skillPhyAtkUp + cfg_passive.skillPhyAtkUp
    basic.phyDfsUp = cfg_basic.skillPhyDfsUp + cfg_extra.skillPhyDfsUp + cfg_passive.skillPhyDfsUp
    basic.magAtkUp = cfg_basic.skillMagAtkUp + cfg_extra.skillMagAtkUp + cfg_passive.skillMagAtkUp
    basic.magDfsUp = cfg_basic.skillMagDfsUp + cfg_extra.skillMagDfsUp + cfg_passive.skillMagDfsUp
    basic.cureUp = cfg_basic.cureUp + cfg_extra.cureUp + cfg_passive.cureUp
    basic.healUp = cfg_basic.healUp + cfg_extra.healUp + cfg_passive.healUp
    basic.finalDmg = 0
    basic.phyPen = cfg_basic.phyPen + cfg_extra.phyPen
    basic.magPen = cfg_basic.magPen + cfg_extra.magPen
    -- 世界道具属性百分比加成（相对于角色自身）
    local percent = {}
    percent.hpMaxPercent = cfg_extra.hpMaxPercent + cfg_passive.hpMaxPercent
    percent.phyAtkPercent = cfg_extra.phyAtkPercent + cfg_passive.phyAtkPercent
    percent.magAtkPercent = cfg_extra.magAtkPercent + cfg_passive.magAtkPercent
    percent.phyDfsPercent = cfg_extra.phyDfsPercent + cfg_passive.phyDfsPercent
    percent.magDfsPercent = cfg_extra.magDfsPercent + cfg_passive.magDfsPercent
    return basic, percent
end

-- 角色最终属性计算
function chara_addition(chara, equips, item)
    local cfg_chara_info = require "Data/Config/cfg_character"
    local chara_info = cfg_chara_info[chara.characterId] or {}
    local camp = chara_info.camp
    equips = equips or {}
    local addition = {}
    -- 基础属性
    local chara_basic = chara_basic(chara)
    for k, v in pairs(chara_basic) do
        addition[k] = v
    end
    -- 计算所有装备的基础属性和基础百分比加成
    for _, equip in pairs(equips) do
        local equip_basic, equip_percent = equip_basic(equip, camp)
        -- 累加装备属性
        for k, v in pairs(equip_basic) do
            addition[k] = addition[k] + v
        end
        -- 累加百分比加成
        addition.hpMax = addition.hpMax + chara_basic.hpMax * equip_percent.hpMaxPercent * 0.0001
        addition.phyAtk = addition.phyAtk + chara_basic.phyAtk * equip_percent.phyAtkPercent * 0.0001
        addition.magAtk = addition.magAtk + chara_basic.magAtk * equip_percent.magAtkPercent * 0.0001
        addition.phyDfs = addition.phyDfs + chara_basic.phyDfs * equip_percent.phyDfsPercent * 0.0001
        addition.magDfs = addition.magDfs + chara_basic.magDfs * equip_percent.magDfsPercent * 0.0001
    end
    -- 计算世界道具的基础属性和百分比加成
    if item then
        local world_item_basic, world_item_percent = world_item_basic(item, camp)
        -- 累加世界道具属性
        for k, v in pairs(world_item_basic) do
            addition[k] = addition[k] + v
        end
        -- 累加百分比加成
        addition.hpMax = addition.hpMax + chara_basic.hpMax * world_item_percent.hpMaxPercent * 0.0001
        addition.phyAtk = addition.phyAtk + chara_basic.phyAtk * world_item_percent.phyAtkPercent * 0.0001
        addition.magAtk = addition.magAtk + chara_basic.magAtk * world_item_percent.magAtkPercent * 0.0001
        addition.phyDfs = addition.phyDfs + chara_basic.phyDfs * world_item_percent.phyDfsPercent * 0.0001
        addition.magDfs = addition.magDfs + chara_basic.magDfs * world_item_percent.magDfsPercent * 0.0001
    end
    return addition
end

-- 角色技能相关
local function chara_skill_basic(chara)
    local util_cfg_character_skill = require "Data/Util/util_cfg_character_skill"
    local skills = {}
    -- 基础技能
    local cfgs_skill = util_cfg_character_skill.get_role_level_cfgs(chara.characterId, 1)
    for _, cfg_skill in pairs(cfgs_skill) do
        skills[cfg_skill.skillNumber] = {}
        skills[cfg_skill.skillNumber].level = 1
        skills[cfg_skill.skillNumber].cfg = cfg_skill
    end
    -- 升级后技能
    if chara.skills then
        for _, skill in pairs(chara.skills) do
            skills[skill.skillId].level = skill.level
            skills[skill.skillId].unlock = true
        end
    end
    return skills
end

-- 装备技能相关
local function equips_skill_basic(equips)
    local util_cfg_equip_data = require "Data/Util/util_cfg_equip_data"
    equips = equips or {}
    local skills = {}
    for _, equip in pairs(equips) do
        local passiveId = util_cfg_equip_data.get_passive_id(equip.equipmentId, equip.level, equip.stage)
        skills[#skills + 1] = passiveId
    end
    return skills
end

-- 世界物品技能相关
local function item_skill_basic(item)
    if nil == item then
        return nil, nil
    end
    local util_cfg_world_item_data = require "Data/Util/util_cfg_world_item_data"
    local campLimit = util_cfg_world_item_data.get(item.worldItemId).equipCampLmt or {}
    local passiveId = util_cfg_world_item_data.get_passive_id(item.worldItemId, item.level, item.star)
    return passiveId, campLimit
end

-- 获取角色战力
function GetCharacterCombatPower(data)
    local chara = data.Character
    local equips = data.Equips or {}
    local item = data.Item
    local level = data.Level
    local util_cfg_character_skill = require "Data/Util/util_cfg_character_skill"
    local cfg_combat_power = require "Data/Config/cfg_combat_power"
    local cfg_equip_skill = require "Data/Config/cfg_equip_skill"
    local cfg_power = cfg_combat_power[chara.characterId]
    local chara_addition = chara_addition(chara, equips, item)
    local level_ratio = RoundToThree(level / 60)
    -- 属性战力
    local property_power = 0
    property_power = property_power + chara_addition.hpMax * cfg_power.hpMaxFinal
    property_power = property_power + chara_addition.phyAtk * cfg_power.phyAtkFinal
    property_power = property_power + chara_addition.magAtk * cfg_power.magAtkFinal
    property_power = property_power + chara_addition.phyDfs * cfg_power.phyDfsFinal
    property_power = property_power + chara_addition.magDfs * cfg_power.magDfsFinal
    property_power = property_power + chara_addition.critAtkRatio * cfg_power.critAtkRatio * level_ratio
    property_power = property_power + chara_addition.critDfsRatio * cfg_power.critDfsRatio * level_ratio
    property_power = property_power + chara_addition.critAtkValue * cfg_power.critAtkValue * level_ratio
    property_power = property_power + chara_addition.critDfsValue * cfg_power.critDfsValue * level_ratio
    property_power = property_power + (chara_addition.hitRateValue - 10000) * cfg_power.hitRateValue * level_ratio --基础命中不算战力,此处10000是基础命中
    property_power = property_power + chara_addition.evadeValue * cfg_power.evadeValue * level_ratio
    property_power = property_power + chara_addition.normalAtkUp * cfg_power.normalAtkUp * level_ratio
    property_power = property_power + chara_addition.normalDfsUp * cfg_power.normalDfsUp * level_ratio
    property_power = property_power + chara_addition.skillAtkUp * cfg_power.skillAtkUp * level_ratio
    property_power = property_power + chara_addition.skillDfsUp * cfg_power.skillDfsUp * level_ratio
    property_power = property_power + chara_addition.ultraAtkUp * cfg_power.ultraAtkUp * level_ratio
    property_power = property_power + chara_addition.ultraDfsUp * cfg_power.ultraDfsUp * level_ratio
    property_power = property_power + chara_addition.phyAtkUp * cfg_power.skillPhyAtkUp * level_ratio
    property_power = property_power + chara_addition.phyDfsUp * cfg_power.skillPhyDfsUp * level_ratio
    property_power = property_power + chara_addition.magAtkUp * cfg_power.skillMagAtkUp * level_ratio
    property_power = property_power + chara_addition.magDfsUp * cfg_power.skillMagDfsUp * level_ratio
    property_power = property_power + chara_addition.cureUp * cfg_power.cureUp * level_ratio
    property_power = property_power + chara_addition.healUp * cfg_power.healUp * level_ratio
    property_power = property_power + chara_addition.finalDmg * cfg_power.finalDmg
    property_power = property_power + chara_addition.phyPen * cfg_power.phyPen * level_ratio
    property_power = property_power + chara_addition.magPen * cfg_power.magPen * level_ratio
    -- 技能战力
    local skill_power = 0
    -- 角色技能战力
    local skills = chara_skill_basic(chara)
    for id, skill in pairs(skills) do
        if skill.unlock then
            local cfg_skill = util_cfg_character_skill.get(chara.characterId, id, skill.level)
            if cfg_skill.combatPower then
                skill_power =
                    skill_power + cfg_skill.combatPower[1] +
                    RoundToThree(property_power * cfg_skill.combatPower[2] / 10000)
            end
        end
    end
    -- 装备技能战力
    local equips_skill = equips_skill_basic(equips)
    for _, skillId in pairs(equips_skill) do
        local cfg_skill = cfg_equip_skill[skillId]
        if cfg_skill and cfg_skill.combatPower then
            skill_power =
                skill_power + cfg_skill.combatPower[1] + RoundToThree(property_power * cfg_skill.combatPower[2] / 10000)
        end
    end
    -- 世界物品技能战力
    local item_skill, campLimit = item_skill_basic(item)
    if item_skill and item_skill > 0 then
        -- 世界道具也使用装备配置表？策划是这样说的
        local cfg_skill = cfg_equip_skill[item_skill]
        if cfg_skill and cfg_skill.combatPower then
            skill_power =
                skill_power + cfg_skill.combatPower[1] + RoundToThree(property_power * cfg_skill.combatPower[2] / 10000)
        end
    end
    return Round(property_power + skill_power)
end

-- 至尊基础属性
local function hero_basic(hero, level)
    local util_cfg_hero_data = require "Data/Util/util_cfg_hero_data"
    local id = hero.heroId
    local cfg_hero_data = util_cfg_hero_data.get(id, level)
    local addition = {}
    addition.hpMax = cfg_hero_data.hpMax
    addition.phyAtk = cfg_hero_data.phyAtk
    addition.magAtk = cfg_hero_data.magAtk
    addition.phyDfs = cfg_hero_data.phyDfs
    addition.magDfs = cfg_hero_data.magDfs
    addition.critAtkRatio = cfg_hero_data.critAtkRatio
    addition.critAtkValue = cfg_hero_data.critAtkValue
    return addition
end

-- 缔结誓约角色
local function pledge_basic(pledge, hero)
    local cfg_hero = require "Data/Config/cfg_hero"
    local hero_id = hero.heroId
    local chara = pledge.Character
    local equips = pledge.Equips
    local item = pledge.Item
    local chara_addition = chara_addition(chara, equips, item)
    local chara_rarity = chara.rarity
    local cfg_hero = cfg_hero[hero_id]
    local addition = {}
    local percent = cfg_hero.attrRates[chara_rarity]
    local key = {"hpMax", "phyAtk", "magAtk", "phyDfs", "magDfs"}
    key = key[cfg_hero.mainAttr]
    addition[key] = chara_addition[key]
    return addition, percent, chara_rarity -- 加成，百分比，稀有度
end

-- 骨王技能相关
function hero_skill_basic(hero)
    local cfg_hero = require "Data/Config/cfg_hero"
    local util_cfg_hero_skill = require "Data/Util/util_cfg_hero_skill"
    local id = hero.heroId
    local skills = {}
    local passive_percent_skill = nil
    -- 基础技能
    local ids = cfg_hero[id].skills
    for i = 1, #ids do
        local cfg_skill = util_cfg_hero_skill.get_level_cfg(ids[i], 1)
        skills[cfg_skill.skillID] = {}
        skills[cfg_skill.skillID].level = 1
        skills[cfg_skill.skillID].cfg = cfg_skill
        if cfg_skill.invokeType == 4 then
            passive_percent_skill = skills[cfg_skill.skillID]
        end
    end
    -- 服务器技能数据
    if hero.skills then
        for k, skill in pairs(hero.skills) do
            skills[skill.skillId].level = skill.level
            -- 解锁技能才会出现再hero.skills
            skills[skill.skillId].unlock = true
            skills[skill.skillId].cfg = util_cfg_hero_skill.get_level_cfg(skill.skillId, skill.level)
        end
    end
    return skills, passive_percent_skill
end

-- 至尊最终属性加成
function hero_addition(hero, pledges, level, passive_percent_skill)
    pledges = pledges or {}
    local addition = {}
    -- 基础属性
    local hero_basic = hero_basic(hero, level)
    local passive_percent = 0
    for k, v in pairs(hero_basic) do
        addition[k] = v
    end

    -- 计算所有缔结誓约的基础属性
    for _, pledge in pairs(pledges) do
        local pledge_basic, percent, chara_rarity = pledge_basic(pledge, hero)
        -- 至尊属性加成变动比较大，被动加成修改缔结加成的percent
        if passive_percent_skill and passive_percent_skill.unlock then
            passive_percent = passive_percent_skill.cfg.attrRates[chara_rarity] or passive_percent
        end
        -- pledge_basic虽然是table,但其中只有一条数据
        for k, v in pairs(pledge_basic) do
            addition[k] = addition[k] + Round(v * (percent + passive_percent) / 10000)
        end
    end
    return addition
end

-- 获取至尊战力
function GetHeroCombatPower(data)
    local hero = data.Hero
    local level = data.Level
    local pledges = data.Pledges
    local cfg_combat_power = require "Data/Config/cfg_combat_power"
    -- 注意skills中包含passive_percent_skill
    local skills, passive_percent_skill = hero_skill_basic(hero, level)
    local hero_addition = hero_addition(hero, pledges, level, passive_percent_skill)
    local cfg_power = cfg_combat_power[hero.heroId]
    -- 属性战力
    local property_power = 0
    property_power = property_power + hero_addition.hpMax * cfg_power.hpMaxFinal
    property_power = property_power + hero_addition.phyAtk * cfg_power.phyAtkFinal
    property_power = property_power + hero_addition.magAtk * cfg_power.magAtkFinal
    property_power = property_power + hero_addition.phyDfs * cfg_power.phyDfsFinal
    property_power = property_power + hero_addition.magDfs * cfg_power.magDfsFinal
    -- 技能战力
    local skill_power = 0
    for id, skill in pairs(skills) do
        if skill.unlock then
            local cfg_skill = skill.cfg
            if cfg_skill.combatPower then
                local cfg_power_1 = cfg_skill.combatPower[1] or 0
                local cfg_power_2 = cfg_skill.combatPower[2] or 0
                skill_power = skill_power + cfg_power_1 + (property_power * cfg_power_2 / 10000)
            end
        end
    end
    return Round(property_power + skill_power)
end
