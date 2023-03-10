package static

// sheetFileName: cfg_yggdrasil_dispatch【派遣】.xlsx
const (
	YggdrasilConditionTypeYggDispatchLevel        = 1 // 世界探索派遣角色等级，参数:[要求的角色数量:等级]
	YggdrasilConditionTypeYggDispatchCharaCamp    = 2 // 世界探索派遣角色阵营，参数:[要求的角色数量:阵营id]
	YggdrasilConditionTypeYggDispatchCharaCareer  = 3 // 世界探索派遣角色职业，参数:[要求的角色数量:职业id]
	YggdrasilConditionTypeYggDispatchCharaRarity  = 4 // 世界探索派遣角色稀有度，参数:[要求的角色数量:稀有度]
	YggdrasilConditionTypeYggDispatchCharaStar    = 5 // 世界探索派遣角色星级，参数:[要求的角色数量:星级]
	YggdrasilConditionTypeYggDispatchPower        = 6 // 世界探索派遣角色总战力，参数:[战力]
	YggdrasilConditionTypeYggDispatchSpecificChar = 7 // 世界探索派遣指定角色，参数:[角色id:等级:阶级:星级]
	YggdrasilConditionTypeYggDispatchCamp         = 8 // 世界探索派遣所有角色统一阵营，参数:[阵营id]
	YggdrasilConditionTypeYggDispatchCareer       = 9 // 世界探索派遣所有角色统一职业，参数:[职业id]
)
