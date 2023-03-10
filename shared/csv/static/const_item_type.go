package static

// sheetFileName: cfg_item_data【道具】.xlsx
const (
	ItemTypeCurrency              = 1   // 货币（金币，钻石，体力）
	ItemTypeGacha                 = 2   // 抽卡道具
	ItemTypeCharExp               = 3   // 角色经验道具
	ItemTypeCharSkill             = 4   // 角色进阶/技能道具
	ItemTypeCharPiece             = 5   // 角色碎片
	ItemTypeWorldItemExp          = 6   // 神器经验道具
	ItemTypeWorldItemBreak        = 7   // 神器突破道具
	ItemTypeEquipExp              = 8   // 装备经验道具
	ItemTypeEquipBreak            = 9   // 装备突破道具
	ItemTypeGraveyardAccelerate   = 10  // 模拟经营加速道具
	ItemTypeGraveyardGetProduct   = 11  // 模拟经营使用道具直接获得产出
	ItemTypeGraveyardProduceBuff  = 12  // 模拟经营生产buff道具
	ItemTypeEnergyItem            = 13  // 体力药水
	ItemTypeCommon                = 14  // 一般道具
	ItemTypeCharacterAdaptiveGift = 15  // 根据新手池子抽出的角色自适应的道具礼包
	ItemTypeCharacterAvatarAuto   = 16  // 获得角色时自动获取对应头像
	ItemTypeCharacterAvatar       = 17  // 头像
	ItemTypeCharacterFrame        = 18  // 头像框
	ItemTypeHeroUnlockItem        = 19  // 至尊解锁道具
	ItemTypeGiftSelectOne         = 101 // N选1礼包
	ItemTypeGiftRandomDrop        = 102 // 随机普通礼包
	ItemTypeAutoRandom            = 103 // 获得后不进背包，自动使用 且 使用后执行一般掉落逻辑的道具-即按照原掉落规则纯随机
	ItemTypeAutoRandomRemove      = 104 // 获得后不进背包，自动使用 且 使用后执行奖池掉落逻辑的道具
	ItemTypeYggItemSpecial        = 105 // 世界探索特殊道具（魔石，不占用世界探索格子，不能存到仓库）
	ItemTypeYggPrestige           = 106 // 世界探索声望
	ItemTypeEquipment             = 301 // 装备
	ItemTypeWorldItem             = 401 // 世界级道具
	ItemTypeCharacter             = 201 // 角色
)
