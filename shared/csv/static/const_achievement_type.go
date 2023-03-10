package static

// sheetFileName: cfg_chapter_level_achievement.xlsx
const (
	AchievementTypeTime        = 1  // (X)秒内通关时间
	AchievementTypeAlive       = 2  // 全员存活
	AchievementTypeDead        = 3  // 阵亡不超过(X)人
	AchievementTypeCharacter   = 4  // 上阵特定角色(X)
	AchievementTypeCareer      = 5  // 上阵(X)个职业(Y)
	AchievementTypeNoCareer    = 6  // 不上阵(X)个职业(Y)，最小只能填1
	AchievementTypeUltra       = 7  // 使用(X)次奥义
	AchievementTypeHeroSkill   = 8  // 使用(X)次英雄技能
	AchievementTypeUltraId     = 9  // 使用(X)次侍从(Y)的奥义
	AchievementTypeHeroSkillId = 10 // 使用(X)次英雄技能(Y)
	AchievementTypeNoUltra     = 11 // 不使用奥义
	AchievementTypeNoHeroSkill = 12 // 不使用英雄技能
	AchievementTypeRarity      = 13 // 上阵(X)个品质(Y)的角色
	AchievementTypeNoRarity    = 14 // 不上阵(X)个品质(Y)角色，最小只能填1
	AchievementTypeHp          = 15 // 任何角色的血量不低于(X)%通关
	AchievementTypeAnd         = 16 // 达成(X)成就、(Y)成就与(Z)成就。不可嵌套类型16和17
	AchievementTypeOr          = 17 // 达成(X)成就、(Y)成就或(Z)成就。不可嵌套类型16和17
	AchievementTypeCamp        = 18 // 上阵阵营限制，参数填写阵营ID,数量：使用(X)个阵营(Y)通关，达成该条件。
	AchievementTypeNum         = 19 // 上阵人数限制，参数填写人数(X)，通关阵容人数≤(X)时，达成该条件。
)
