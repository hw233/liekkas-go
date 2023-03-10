local cfg_hero_skill = require "Data/Config/cfg_hero_skill"
local util = {}
local cache = {}

for k, cfg in pairs(cfg_hero_skill) do
    if not cache[cfg.skillID] then
        cache[cfg.skillID] = {}
    end
    cache[cfg.skillID][cfg.skillLevel] = cfg
end

function util.get(skillId)
    return cfg_hero_skill[skillId]
end

--- 根据技能的基础id和等级获取最终技能
function util.get_level_cfg(baseId, level)
    return cache[baseId][level]
end

return util