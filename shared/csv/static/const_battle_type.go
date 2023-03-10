package static

// sheetFileName: cfg_explore_chapter_level.xlsx
const (
	BattleTypeExplore         = 1  // 主线探索副本, 参数[]
	BattleTypeExploreResource = 2  // 主线探索资源点, 参数[资源点id]
	BattleTypeExploreMonster  = 3  // 主线探索怪物, 参数[怪物表id]
	BattleTypeExploreElite    = 4  // 历练, 参数[]
	BattleTypeTower           = 5  // 塔, 参数[塔id, 层数]
	BattleTypeYgg             = 6  // 世界探索
	BattleTypeTeam            = 7  // 组队战
	BattleTypeArena           = 8  // 竞技场
	BattleTypeFormation       = 9  // 编队
	BattleTypeGate            = 10 // 混沌之门
	BattleTypeChallengeAltar  = 11 // 挑战祭坛
)
