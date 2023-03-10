package static

// sheetFileName: cfg_yggdrasil_sub_task.xlsx
const (
	YggdrasilSubTaskTypeTypeVn                   = 1  // 对话任务
	YggdrasilSubTaskTypeTypeChapter              = 2  // 完成对应ID的战斗关卡
	YggdrasilSubTaskTypeTypeMonster              = 3  // 击杀X个指定怪物
	YggdrasilSubTaskTypeTypeDeliverItem          = 4  // 交付物品
	YggdrasilSubTaskTypeTypeConvoy               = 5  // 护送
	YggdrasilSubTaskTypeTypeLeadWay              = 6  // 带路
	YggdrasilSubTaskTypeTypeBuild                = 7  // 建造建筑
	YggdrasilSubTaskTypeTypeHelpBuild            = 8  // 帮助建造建筑
	YggdrasilSubTaskTypeTypeMove                 = 9  // 移动
	YggdrasilSubTaskTypeTypeMultiTask            = 10 // 完成多个其他任务ID
	YggdrasilSubTaskTypeTypeObjectStateChange    = 11 // 与对应object交互
	YggdrasilSubTaskTypeTypeCity                 = 12 // 在某个主城内
	YggdrasilSubTaskTypeTypeOwn                  = 13 // 拥有对应物品
	YggdrasilSubTaskTypeTypeDeliverItemSelectOne = 14 // 多类道具交付其中1个
	YggdrasilSubTaskTypeTypeMultiSubTask         = 15 // 子任务同时有多个子任务
)
