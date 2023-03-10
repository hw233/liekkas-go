local cfg_character_data = require "Data/Config/cfg_character_data"
local util = {}
local cache = {}

--- 通过[角色id] 返回 角色配置
function util.get(characterId)
    return cfg_character_data[characterId]
end

return util

