package static

// sheetFileName: const_condition_type【玩家条件类型】.xlsx
const (
	ConditionTypeUserLevel                 = 1    // 玩家等级, 参数: [等级]
	ConditionTypePassLevel                 = 2    // 通过关卡, 参数: [关卡id]
	ConditionTypeFinishChapter             = 3    // 完成章节, 参数: [章节id]
	ConditionTypeFinishQuest               = 4    // 完成任务, 参数: [任务id]
	ConditionTypeExplorePoint              = 5    // 主线进度, 参数: [map_object_type:id]
	ConditionTypeYggCompleteSubtask1       = 10   // 完成除了最后一条子任务，参数:[subtaskId]
	ConditionTypeYggCompleteSubtask2       = 11   // 完成最后一条子任务，参数:[subtaskId]
	ConditionTypeYggAcceptSubtask          = 12   // 领取世界探索子任务，参数:[subtaskId]
	ConditionTypeGuide                     = 13   // 完成某一引导，参数:[guideId]
	ConditionTypeYggAreaPrestige           = 14   // 世界探索区域声望达到，参数:[areaId:prestige]
	ConditionTypeGraveyardBuildLv          = 15   // 模拟经营某个建筑达到某等级(没揭幕是0级)，参数:[buildId:lv]
	ConditionTypeYggCompleteTask           = 16   // 完成世界探索-某个任务，参数:[taskId]
	ConditionTypeYggMoveCount              = 17   // 世界探索移动x步，参数:[移动步数]
	ConditionTypeCharaLevel                = 2001 // 角色等级, 参数: [角色id: 等级]
	ConditionTypeCharaStar                 = 2002 // 角色星级, 参数: [角色id: 星级]
	ConditionTypeHeroLevel                 = 3001 // 至尊等级, 参数: [等级]
	ConditionTypeHeroSkillLevel            = 3002 // 至尊技能等级, 参数: [至尊id: 技能id: 等级]
	ConditionTypeHeroSkillItemUsed         = 3003 // 至尊技能点使用, 参数: [至尊id: 使用量]
	ConditionTypeHasItem                   = 4001 // 拥有道具, 参数: [道具id: 数量]
	ConditionTypeApNoMoreThan              = 5001 // 世界探索精力不超过，参数：[AP]
	ConditionTypeYggPosBuildable           = 5002 // 世界探索地块可建造，参数：[Type](0，不可建造，1-可建造)
	ConditionTypeYggPosNoInitiativeMonster = 5003 // 世界探索多少格内无主动怪，参数：[Radius]
)
