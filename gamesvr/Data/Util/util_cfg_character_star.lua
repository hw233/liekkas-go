local cfg_character_star = require "Data/Config/cfg_character_star"
local util = {}
local cache = {}
-- cache[角色id][star] = 配置表

for k, cfg in pairs(cfg_character_star) do
    if not cache[cfg.charID] then
        cache[cfg.charID] = {}
    end
    cache[cfg.charID][cfg.star] = cfg
end

--- 根据[角色id] [角色星级] 返回 [配置]
---
--- 如果星级空则返回角色所有星级的配置
function util.get(characterId, star)
    if nil == cache[characterId] then
        print("cfg_character_star中缺少配置,角色id:" .. characterId)
    end
    if nil == star then
        return cache[characterId]
    end
    return cache[characterId][star]
end

return util
