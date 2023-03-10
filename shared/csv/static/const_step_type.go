package static

// sheetFileName: cfg_guide_step【新手引导步骤配置表】.xlsx
const (
	StepTypeDialog      = 1  // 对话框
	StepTypeClick       = 2  // 点击
	StepTypeArrow       = 3  // 箭头
	StepTypeDelay       = 4  // 延时
	StepTypeHidedialog  = 5  // 隐藏对话框
	StepTypeHidearrow   = 6  // 隐藏箭头
	StepTypeComplete    = 7  // 完成引导模块
	StepTypeWaittrigger = 8  // 等待触发下一步引导
	StepTypeCoding      = 9  // 特殊逻辑
	StepTypeHideprefab  = 10 //
	StepTypeShowprefab  = 11 //
	StepTypeHidemask    = 12 // 隐藏遮罩
	StepTypeHint        = 13 // 弱引导组件
)
