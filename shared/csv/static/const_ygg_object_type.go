package static

// sheetFileName: cfg_yggdrasil_object_state.xlsx
const (
	YggObjectTypeChest             = 1  // 宝箱
	YggObjectTypePassivermonster   = 2  // 被动怪物
	YggObjectTypeInitiativemonster = 3  // 主动怪物
	YggObjectTypeNpc               = 4  // npc
	YggObjectTypeBuilding          = 5  // 建筑
	YggObjectTypeMagictable        = 6  // 魔台,有vn
	YggObjectTypeEffect            = 7  // 表现类、无Vnor交互参数
	YggObjectTypeMysticmagic       = 9  // 伪装魔法，驱散后显示宝箱
	YggObjectTypeTaskmonster       = 10 // 护送npc任务怪
	YggObjectTypeFollownpc         = 11 // 跟随npc【ID不能与其他object重复】
	YggObjectTypeBuildingnpc       = 12 // 建筑类npc【类似魔台的弹出形式，但是没有选择按钮，纯展示用】
	YggObjectTypeFollowbattlenpc   = 13 // 跟随助战npc，objectParam对应cfg_battle_npc表的id
	YggObjectTypeChallengealtar    = 14 // 挑战祭台
)
