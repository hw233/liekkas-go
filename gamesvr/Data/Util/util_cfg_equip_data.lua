local cfg_equip_data = require "Data/Config/cfg_equip_data"
local util = {}

function util.get(equipId)
    return cfg_equip_data[equipId]
end

-- 客户端EquipModel.AttributeId
function util.get_attribute_id(equipId, equipLevel, equipStage)
    local cfg_equip = cfg_equip_data[equipId]
    if nil == cfg_equip.attributeID then
        return 0
    end
    -- 13011:1:10,13012:2:1,13013:2:2,13014:2:3,13015:2:4
    local passiveIDs = {}
    for passiveID in string.gmatch(cfg_equip.attributeID, "([^,]+)") do
        passiveIDs[#passiveIDs+1] = passiveID
    end
    -- 倒序，从后往前匹配，找到则返回
    for i=#passiveIDs, 1, -1 do
        local tmp = {}
        for v in string.gmatch(passiveIDs[i], "([^:]+)") do
            tmp[#tmp+1] = tonumber(v)
        end
        local id = tmp[1]
        local type = tmp[2]
        local level = tmp[3]
        if type == 1 then
            if equipLevel >= level then
                return id
            end
        end

        if type == 2 then
            if equipStage >= level then
                return id
            end
        end
    end
    return 0
end

function util.get_passive_id(equipId, equipLevel, equipStage)
    local cfg_equip = cfg_equip_data[equipId]
    if nil == cfg_equip.passiveID or '' ==  cfg_equip.passiveID then
        return 0
    end
    -- 70011:1:10,70012:2:1,70013:2:2,70014:2:3
    local passiveIDs = {}
    for passiveID in string.gmatch(cfg_equip.passiveID, "([^,]+)") do
        passiveIDs[#passiveIDs+1] = passiveID
    end
    -- 倒序，从后往前匹配，找到则返回
    for i=#passiveIDs, 1, -1 do
        local tmp = {}
        for v in string.gmatch(passiveIDs[i], "([^:]+)") do
            tmp[#tmp+1] = tonumber(v)
        end
        local id = tmp[1]
        local type = tmp[2]
        local level = tmp[3]
        if type == 1 then
            if equipLevel >= level then
                return id
            end
        end

        if type == 2 then
            if equipStage >= level then
                return id
            end
        end
    end
    return 0
end

return util