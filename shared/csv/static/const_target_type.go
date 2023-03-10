package static

// sheetFileName: cfg_targeting.xlsx
const (
	TargetTypeOwn               = 1  // 自己
	TargetTypeFriend            = 2  // 友军，即包含友方单位，和全体友方召唤物。不包括自己和英雄
	TargetTypeHostile           = 3  // 敌军，即包含敌方单位，和全体敌方召唤物。不包括英雄
	TargetTypeSkilltarget       = 4  // skill筛选出的一级目标（有可能是空，因为目标离开范围）
	TargetTypeAitarget          = 5  // 当前对象ai的攻击目标
	TargetTypeFriendsummon      = 6  // 友方召唤
	TargetTypeHostilesummon     = 7  // 敌方召唤
	TargetTypeFriendhero        = 8  // 友方英雄
	TargetTypeHostilehero       = 9  // 敌方英雄
	TargetTypeFriendall         = 10 // 全体友方，包括自己
	TargetTypeIgnore            = 11 // 忽略目标，特效使用
	TargetTypeAitargetorhostile = 12 // 如果ai目标存在则使用ai目标，否则重新寻找目标
	TargetTypeAll               = 13 // 全单位，包括自己
	TargetTypeOthers            = 14 // 全单位，不包括自己
)
