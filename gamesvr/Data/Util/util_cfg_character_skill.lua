local cfg_character_skill = require "Data/Config/cfg_character_skill"
local const_skill_type = require "Data/Config/const_skill_type"
local util = {}
local cache = {}
-- 缓存结构
-- cache[角色id][技能id][技能level] = 配置表

for k, cfg in pairs(cfg_character_skill) do
    if not cache[cfg.roleID] then
        cache[cfg.roleID] = {}
    end
    if not cache[cfg.roleID][cfg.skillNumber] then
        cache[cfg.roleID][cfg.skillNumber] = {}
    end
    cache[cfg.roleID][cfg.skillNumber][cfg.skillLevel] = cfg
end

--- 根据[人物id] [技能id] [技能level] 返回 [技能配置] 
---
--- 如果skillId空则返回角色所有技能
---
--- 如果skillLevel空则返回角色技能所有等级的技能
function util.get(roleID, skillId, skillLevel)
    if nil == skillId then
        return cache[roleID]
    end

    if nil == skillLevel then
        return cache[roleID][skillId]
    end

    return cache[roleID][skillId][skillLevel]
end

--- 获取角色技能等级的所有技能
function util.get_role_level_cfgs(roleID, skillLevel)
    local t = {}
    local tmp = util.get(roleID)
    for k, skills in pairs(tmp) do
        t[#t+1] = skills[skillLevel]
    end 
    return t
end

--- 获取角色基础被动技能,注意和isPassive区分
function util.get_role_base_passive_cfgs(roleID)
    local t = {}
    local tmp = util.get(roleID)
    for k, skills in pairs(tmp) do
        for v, skill in pairs(skills) do
            if skill.skillType == const_skill_type.PassiveSkill then
                t[#t+1] = skill
            end
        end
    end 
    return t
end

-- --- 获取角色基础被动技能的最终解锁的配置
-- function util.get_role_base_passive_unlock_cfg(VoUserCharacter)
--     local cfgs_passive = util.get_role_base_passive_cfgs(VoUserCharacter.characterId)
--     local cfg_passive = nil
--     for k, cfg_skill in pairs(cfgs_passive) do
--         if util.is_unlock(cfg_skill, VoUserCharacter) then
--             if cfg_passive == nil or cfg_passive.skillLevel < cfg_skill.skillLevel then
--                 cfg_passive = cfg_skill
--             end
--         end
--     end
--     return cfg_passive
-- end

-- --- 根据服务器人物数据判断当前技能是否解锁
-- function util.is_unlock(cfg_skill, VoUserCharacter)
--     local level = VoUserCharacter.level
--     local stage = VoUserCharacter.stage
--     local star = VoUserCharacter.star

--     if cfg_skill.skillUnlock == 1 then
--         return (level >= cfg_skill.unlockParam)
--     end

--     if cfg_skill.skillUnlock == 2 then
--         return (stage >= cfg_skill.unlockParam)
--     end

--     if cfg_skill.skillUnlock == 3 then
--         return (star >= cfg_skill.unlockParam)
--     end

--     if cfg_skill.skillUnlock == 4 then
--         return (star >= cfg_skill.unlockParam)
--     end

--     return true
-- end

return util
