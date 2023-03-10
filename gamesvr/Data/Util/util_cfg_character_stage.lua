local cfg_character_stage = require "Data/Config/cfg_character_stage"
local util = {}
local cache = {}
-- cache[角色id][stage] = 配置表

for k, cfg in pairs(cfg_character_stage) do
    if not cache[cfg.charId] then
        cache[cfg.charId] = {}
    end
    cache[cfg.charId][cfg.stage] = cfg
end

--- 根据[角色id] [阶级stage] 返回 [配置]
---
--- 如果星级空则返回角色所有阶级的配置
function util.get(characterId, stage)
    if nil == cache[characterId] then
        print("cfg_character_stage中缺少配置,角色id:" .. characterId)
    end
    if nil == stage then
        return cache[characterId]
    end
    return cache[characterId][stage]
end

return util