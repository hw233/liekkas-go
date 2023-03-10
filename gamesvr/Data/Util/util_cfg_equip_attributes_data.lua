local cfg_equip_attributes_data = require "Data/Config/cfg_equip_attributes_data"
local util = {}
local cache = {}

for k, cfg in pairs(cfg_equip_attributes_data) do
    if not cache[cfg.equipID] then
        cache[cfg.equipID] = {}
    end
    cache[cfg.equipID][cfg.stage] = cfg
end

function util.get(equipId, stage)
    if nil == stage then
        return cache[equipId]
    end
    return cache[equipId][stage]
end

function util.get_extra(equipId, stage)
    if nil == stage then
        return cache[equipId]
    end
    local extra = cache[equipId][stage]
    if 0 == stage then
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
