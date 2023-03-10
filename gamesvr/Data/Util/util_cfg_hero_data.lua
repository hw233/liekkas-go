local cfg_hero_data = require "Data/Config/cfg_hero_data"
local util = {}
local cache = {}
-- 缓存结构
-- cache[角色id][技能id][技能level] = 配置表

for k, cfg in pairs(cfg_hero_data) do
    if not cache[cfg.heroID] then
        cache[cfg.heroID] = {}
    end
    cache[cfg.heroID][cfg.heroLv] = cfg
end

--- 通过[英雄id] [英雄lv] 返回 英雄配置
function util.get(heroID, heroLv)
    if nil == cache[heroID] then
        return nil
    end
    return cache[heroID][heroLv]
end

return util