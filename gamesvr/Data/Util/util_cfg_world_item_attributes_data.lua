local cfg_world_item_attributes_data = require "Data/Config/cfg_world_item_attributes_data"
local util = {}
local cache = {}

for k, cfg in pairs(cfg_world_item_attributes_data) do
    if not cache[cfg.worldItemId] then
        cache[cfg.worldItemId] = {}
    end
    cache[cfg.worldItemId][cfg.star] = cfg
end

function util.get(worldItemId, star)
    if nil == star then
        return cache[worldItemId]
    end
    return cache[worldItemId][star]
end

function util.get_extra(worldItemId, star)
    if nil == star then
        return cache[worldItemId]
    end
    local extra = cache[worldItemId][star]
    if 0 == star then
        local t = {}
        for key, _ in pairs(extra) do
            t[key] = 0
        end
        -- 0阶extra保留基础属性的等级加成
        t.hpExtraRatio = extra.hpExtraRatio
        t.phyAtkExtraRatio = extra.phyAtkExtraRatio
        t.magAtkExtraRatio = extra.magAtkExtraRatio
        t.phyDfsExtraRatio = extra.phyDfsExtraRatio
        t.magDfsExtraRatio = extra.magDfsExtraRatio
        extra = t
    end
    return extra
end

return util