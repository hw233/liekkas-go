package static

// sheetFileName: cfg_fx.xlsx
const (
	FxTypeFx          = 1 // 通常特效
	FxTypeQueuebuff   = 2 // 同种多列式展示方式，与2互斥
	FxTypeOverlaybuff = 3 // 同种叠层展示方式，与1互斥
	FxTypeIconbuff    = 4 // 血条下方的Buff，图标展示
	FxTypeLightning   = 5 // 闪电
	FxTypeField       = 6 // 场地技能特效，在场地技能的索敌位置播放
)
