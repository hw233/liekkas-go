local cfg_world_item_data = require "Data/Config/cfg_world_item_data"
local util = {}

function util.get(id)
    return cfg_world_item_data[id]
end

function util.get_attribute_id(id, level, star)
    local cfg_world_item = cfg_world_item_data[id]
    if nil == cfg_world_item.passiveAttriID then
        return 0
    end
    -- 13011:1:10,13012:2:1,13013:2:2,13014:2:3,13015:2:4
    local passiveIDs = {}
    for passiveID in string.gmatch(cfg_world_item.passiveAttriID, "([^,]+)") do
        passiveIDs[#passiveIDs + 1] = passiveID
    end
    -- 倒序，从后往前匹配，找到则返回
    for i = #passiveIDs, 1, -1 do
        local tmp = {}
        for v in string.gmatch(passiveIDs[i], "([^:]+)") do
            tmp[#tmp + 1] = tonumber(v)
        end
        local id = tmp[1]
        -- 类型
        local type = tmp[2]
        -- 比较值
        local value = tmp[3]
        if type == 1 then
            if level >= value then
                return id
            end
        end

        if type == 2 then
            if star >= value then
                return id
            end
        end
    end
    return 0
end

function util.get_passive_id(id, level, star)
    local cfg_world_item = cfg_world_item_data[id]
    if nil == cfg_world_item.passiveSkillID or "" == cfg_world_item.passiveSkillID then
        return 0
    end
    -- 70011:1:10,70012:2:1,70013:2:2,70014:2:3
    local passiveIDs = {}
    for passiveID in string.gmatch(cfg_world_item.passiveSkillID, "([^,]+)") do
        passiveIDs[#passiveIDs + 1] = passiveID
    end
    -- 倒序，从后往前匹配，找到则返回
    for i = #passiveIDs, 1, -1 do
        local tmp = {}
        for v in string.gmatch(passiveIDs[i], "([^:]+)") do
            tmp[#tmp + 1] = tonumber(v)
        end
        local id = tmp[1]
        local type = tmp[2]
        local value = tmp[3]
        if type == 1 then
            if level >= value then
                return id
            end
        end

        if type == 2 then
            if star >= value then
                return id
            end
        end
    end
    return 0
end

return util
