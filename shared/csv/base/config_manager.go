package base

import "shared/utility/glog"

type ConfigManager struct {
	CfgActionUnlockConfig             *CfgActionUnlockConfig
	CfgActivityConfig                 *CfgActivityConfig
	CfgActivityFuncConfig             *CfgActivityFuncConfig
	CfgBattleLevelConfig              *CfgBattleLevelConfig
	CfgBattleNpcConfig                *CfgBattleNpcConfig
	CfgBattlePositionConfig           *CfgBattlePositionConfig
	CfgBossInforConfig                *CfgBossInforConfig
	CfgBuffConfig                     *CfgBuffConfig
	CfgBulletConfig                   *CfgBulletConfig
	CfgChapterLevelAchievementConfig  *CfgChapterLevelAchievementConfig
	CfgCharacterConfig                *CfgCharacterConfig
	CfgCharacterBattleConfig          *CfgCharacterBattleConfig
	CfgCharacterCampConfig            *CfgCharacterCampConfig
	CfgCharacterCareerConfig          *CfgCharacterCareerConfig
	CfgCharacterDataConfig            *CfgCharacterDataConfig
	CfgCharacterFeatureConfig         *CfgCharacterFeatureConfig
	CfgCharacterLevelConfig           *CfgCharacterLevelConfig
	CfgCharacterRarityConfig          *CfgCharacterRarityConfig
	CfgCharacterSkillConfig           *CfgCharacterSkillConfig
	CfgCharacterSkinConfig            *CfgCharacterSkinConfig
	CfgCharacterStageConfig           *CfgCharacterStageConfig
	CfgCharacterStarConfig            *CfgCharacterStarConfig
	CfgCombatPowerConfig              *CfgCombatPowerConfig
	CfgCombatPowerAdaptConfig         *CfgCombatPowerAdaptConfig
	CfgDropDataConfig                 *CfgDropDataConfig
	CfgDropGroupConfig                *CfgDropGroupConfig
	CfgEquipLevelUpDataConfig         *CfgEquipLevelUpDataConfig
	CfgEquipAdvanceConfig             *CfgEquipAdvanceConfig
	CfgEquipAttributesDataConfig      *CfgEquipAttributesDataConfig
	CfgEquipBreakConfig               *CfgEquipBreakConfig
	CfgEquipDataConfig                *CfgEquipDataConfig
	CfgEquipRandAttributesConfig      *CfgEquipRandAttributesConfig
	CfgEquipRandAttributesValueConfig *CfgEquipRandAttributesValueConfig
	CfgEquipSkillConfig               *CfgEquipSkillConfig
	CfgEverydayEnergyReceiveConfig    *CfgEverydayEnergyReceiveConfig
	CfgExploreChapterConfig           *CfgExploreChapterConfig
	CfgExploreChapterLevelConfig      *CfgExploreChapterLevelConfig
	CfgExploreChapterLevelAdaptConfig *CfgExploreChapterLevelAdaptConfig
	CfgExploreChapterRewardConfig     *CfgExploreChapterRewardConfig
	CfgExploreEventConfig             *CfgExploreEventConfig
	CfgExploreFogConfig               *CfgExploreFogConfig
	CfgExploreGatherConfig            *CfgExploreGatherConfig
	CfgExploreMapCoordinateConfig     *CfgExploreMapCoordinateConfig
	CfgExploreMonsterConfig           *CfgExploreMonsterConfig
	CfgExploreNpcConfig               *CfgExploreNpcConfig
	CfgExploreResourceConfig          *CfgExploreResourceConfig
	CfgFieldSkillConfig               *CfgFieldSkillConfig
	CfgFxConfig                       *CfgFxConfig
	CfgGachaDropConfig                *CfgGachaDropConfig
	CfgGachaGlobalConfig              *CfgGachaGlobalConfig
	CfgGachaPoolConfig                *CfgGachaPoolConfig
	CfgGlobalStringConfig             *CfgGlobalStringConfig
	CfgGraveyardBuildConfig           *CfgGraveyardBuildConfig
	CfgGraveyardBuildLevelConfig      *CfgGraveyardBuildLevelConfig
	CfgGraveyardBuildStageConfig      *CfgGraveyardBuildStageConfig
	CfgGraveyardMainTowerLevelConfig  *CfgGraveyardMainTowerLevelConfig
	CfgGraveyardProduceBuffConfig     *CfgGraveyardProduceBuffConfig
	CfgGuideModuleConfig              *CfgGuideModuleConfig
	CfgGuideStepConfig                *CfgGuideStepConfig
	CfgGuildGreetingsConfig           *CfgGuildGreetingsConfig
	CfgGuildLevelConfig               *CfgGuildLevelConfig
	CfgGuildTaskConfig                *CfgGuildTaskConfig
	CfgHeroConfig                     *CfgHeroConfig
	CfgHeroCampSkillFilterConfig      *CfgHeroCampSkillFilterConfig
	CfgHeroDataConfig                 *CfgHeroDataConfig
	CfgHeroExtraDataConfig            *CfgHeroExtraDataConfig
	CfgHeroLevelConfig                *CfgHeroLevelConfig
	CfgHeroSkillConfig                *CfgHeroSkillConfig
	CfgHeroTaskDataConfig             *CfgHeroTaskDataConfig
	CfgItemDataConfig                 *CfgItemDataConfig
	CfgMailTemplateConfig             *CfgMailTemplateConfig
	CfgManualConfig                   *CfgManualConfig
	CfgManualCampConfig               *CfgManualCampConfig
	CfgManualCharacterVoiceConfig     *CfgManualCharacterVoiceConfig
	CfgManualGroupConfig              *CfgManualGroupConfig
	CfgManualLevelAttributeDataConfig *CfgManualLevelAttributeDataConfig
	CfgMaskWordConfig                 *CfgMaskWordConfig
	CfgMonsterBattleConfig            *CfgMonsterBattleConfig
	CfgMonsterDataConfig              *CfgMonsterDataConfig
	CfgNamesNpcConfig                 *CfgNamesNpcConfig
	CfgPassiveEffectConfig            *CfgPassiveEffectConfig
	CfgPasssiveAttributesConfig       *CfgPasssiveAttributesConfig
	CfgQuickPurchaseStaminaConfig     *CfgQuickPurchaseStaminaConfig
	CfgRarityConfig                   *CfgRarityConfig
	CfgScorePassGroupConfig           *CfgScorePassGroupConfig
	CfgScorePassPhaseConfig           *CfgScorePassPhaseConfig
	CfgScorePassRewardConfig          *CfgScorePassRewardConfig
	CfgScorePassTaskConfig            *CfgScorePassTaskConfig
	CfgSigninDataConfig               *CfgSigninDataConfig
	CfgSkillConfig                    *CfgSkillConfig
	CfgSkillEffectConfig              *CfgSkillEffectConfig
	CfgSkillEventConfig               *CfgSkillEventConfig
	CfgSoundConfig                    *CfgSoundConfig
	CfgStoreGeneralConfig             *CfgStoreGeneralConfig
	CfgStoreGoodsConfig               *CfgStoreGoodsConfig
	CfgStoreSubstoreConfig            *CfgStoreSubstoreConfig
	CfgStoreUpdateConfig              *CfgStoreUpdateConfig
	CfgSurveyInfoConfig               *CfgSurveyInfoConfig
	CfgSurveyQuestionConfig           *CfgSurveyQuestionConfig
	CfgTargetingConfig                *CfgTargetingConfig
	CfgTaskDataConfig                 *CfgTaskDataConfig
	CfgTaskRewardConfig               *CfgTaskRewardConfig
	CfgTeamLevelConfig                *CfgTeamLevelConfig
	CfgTowerConfig                    *CfgTowerConfig
	CfgTowerLevelConfig               *CfgTowerLevelConfig
	CfgTowerStageConfig               *CfgTowerStageConfig
	CfgTpGateConfig                   *CfgTpGateConfig
	CfgUltraSkillConfig               *CfgUltraSkillConfig
	CfgUserAvatarConfig               *CfgUserAvatarConfig
	CfgVnConfig                       *CfgVnConfig
	CfgWorldItemAdvanceConfig         *CfgWorldItemAdvanceConfig
	CfgWorldItemAttributesDataConfig  *CfgWorldItemAttributesDataConfig
	CfgWorldItemDataConfig            *CfgWorldItemDataConfig
	CfgWorldItemLevelUpDataConfig     *CfgWorldItemLevelUpDataConfig
	CfgYggdrasilAreaConfig            *CfgYggdrasilAreaConfig
	CfgYggdrasilAreaCityConfig        *CfgYggdrasilAreaCityConfig
	CfgYggdrasilAreaPosConfig         *CfgYggdrasilAreaPosConfig
	CfgYggdrasilBuildingConfig        *CfgYggdrasilBuildingConfig
	CfgYggdrasilCobuildingConfig      *CfgYggdrasilCobuildingConfig
	CfgYggdrasilDailyMonsterConfig    *CfgYggdrasilDailyMonsterConfig
	CfgYggdrasilDispatchConfig        *CfgYggdrasilDispatchConfig
	CfgYggdrasilMarkConfig            *CfgYggdrasilMarkConfig
	CfgYggdrasilObjectConfig          *CfgYggdrasilObjectConfig
	CfgYggdrasilObjectStateConfig     *CfgYggdrasilObjectStateConfig
	CfgYggdrasilSimpleTalkConfig      *CfgYggdrasilSimpleTalkConfig
	CfgYggdrasilSubTaskConfig         *CfgYggdrasilSubTaskConfig
	CfgYggdrasilSubTaskEnvConfig      *CfgYggdrasilSubTaskEnvConfig
	CfgYggdrasilTaskConfig            *CfgYggdrasilTaskConfig
	CfgYggdrasilWorldConfig           *CfgYggdrasilWorldConfig
	GmConfig                          *GmConfig
	ProtocolConfig                    *ProtocolConfig
}

func (m *ConfigManager) Init() {
	m.CfgActionUnlockConfig = NewCfgActionUnlockConfig()
	m.CfgActivityConfig = NewCfgActivityConfig()
	m.CfgActivityFuncConfig = NewCfgActivityFuncConfig()
	m.CfgBattleLevelConfig = NewCfgBattleLevelConfig()
	m.CfgBattleNpcConfig = NewCfgBattleNpcConfig()
	m.CfgBattlePositionConfig = NewCfgBattlePositionConfig()
	m.CfgBossInforConfig = NewCfgBossInforConfig()
	m.CfgBuffConfig = NewCfgBuffConfig()
	m.CfgBulletConfig = NewCfgBulletConfig()
	m.CfgChapterLevelAchievementConfig = NewCfgChapterLevelAchievementConfig()
	m.CfgCharacterConfig = NewCfgCharacterConfig()
	m.CfgCharacterBattleConfig = NewCfgCharacterBattleConfig()
	m.CfgCharacterCampConfig = NewCfgCharacterCampConfig()
	m.CfgCharacterCareerConfig = NewCfgCharacterCareerConfig()
	m.CfgCharacterDataConfig = NewCfgCharacterDataConfig()
	m.CfgCharacterFeatureConfig = NewCfgCharacterFeatureConfig()
	m.CfgCharacterLevelConfig = NewCfgCharacterLevelConfig()
	m.CfgCharacterRarityConfig = NewCfgCharacterRarityConfig()
	m.CfgCharacterSkillConfig = NewCfgCharacterSkillConfig()
	m.CfgCharacterSkinConfig = NewCfgCharacterSkinConfig()
	m.CfgCharacterStageConfig = NewCfgCharacterStageConfig()
	m.CfgCharacterStarConfig = NewCfgCharacterStarConfig()
	m.CfgCombatPowerConfig = NewCfgCombatPowerConfig()
	m.CfgCombatPowerAdaptConfig = NewCfgCombatPowerAdaptConfig()
	m.CfgDropDataConfig = NewCfgDropDataConfig()
	m.CfgDropGroupConfig = NewCfgDropGroupConfig()
	m.CfgEquipLevelUpDataConfig = NewCfgEquipLevelUpDataConfig()
	m.CfgEquipAdvanceConfig = NewCfgEquipAdvanceConfig()
	m.CfgEquipAttributesDataConfig = NewCfgEquipAttributesDataConfig()
	m.CfgEquipBreakConfig = NewCfgEquipBreakConfig()
	m.CfgEquipDataConfig = NewCfgEquipDataConfig()
	m.CfgEquipRandAttributesConfig = NewCfgEquipRandAttributesConfig()
	m.CfgEquipRandAttributesValueConfig = NewCfgEquipRandAttributesValueConfig()
	m.CfgEquipSkillConfig = NewCfgEquipSkillConfig()
	m.CfgEverydayEnergyReceiveConfig = NewCfgEverydayEnergyReceiveConfig()
	m.CfgExploreChapterConfig = NewCfgExploreChapterConfig()
	m.CfgExploreChapterLevelConfig = NewCfgExploreChapterLevelConfig()
	m.CfgExploreChapterLevelAdaptConfig = NewCfgExploreChapterLevelAdaptConfig()
	m.CfgExploreChapterRewardConfig = NewCfgExploreChapterRewardConfig()
	m.CfgExploreEventConfig = NewCfgExploreEventConfig()
	m.CfgExploreFogConfig = NewCfgExploreFogConfig()
	m.CfgExploreGatherConfig = NewCfgExploreGatherConfig()
	m.CfgExploreMapCoordinateConfig = NewCfgExploreMapCoordinateConfig()
	m.CfgExploreMonsterConfig = NewCfgExploreMonsterConfig()
	m.CfgExploreNpcConfig = NewCfgExploreNpcConfig()
	m.CfgExploreResourceConfig = NewCfgExploreResourceConfig()
	m.CfgFieldSkillConfig = NewCfgFieldSkillConfig()
	m.CfgFxConfig = NewCfgFxConfig()
	m.CfgGachaDropConfig = NewCfgGachaDropConfig()
	m.CfgGachaGlobalConfig = NewCfgGachaGlobalConfig()
	m.CfgGachaPoolConfig = NewCfgGachaPoolConfig()
	m.CfgGlobalStringConfig = NewCfgGlobalStringConfig()
	m.CfgGraveyardBuildConfig = NewCfgGraveyardBuildConfig()
	m.CfgGraveyardBuildLevelConfig = NewCfgGraveyardBuildLevelConfig()
	m.CfgGraveyardBuildStageConfig = NewCfgGraveyardBuildStageConfig()
	m.CfgGraveyardMainTowerLevelConfig = NewCfgGraveyardMainTowerLevelConfig()
	m.CfgGraveyardProduceBuffConfig = NewCfgGraveyardProduceBuffConfig()
	m.CfgGuideModuleConfig = NewCfgGuideModuleConfig()
	m.CfgGuideStepConfig = NewCfgGuideStepConfig()
	m.CfgGuildGreetingsConfig = NewCfgGuildGreetingsConfig()
	m.CfgGuildLevelConfig = NewCfgGuildLevelConfig()
	m.CfgGuildTaskConfig = NewCfgGuildTaskConfig()
	m.CfgHeroConfig = NewCfgHeroConfig()
	m.CfgHeroCampSkillFilterConfig = NewCfgHeroCampSkillFilterConfig()
	m.CfgHeroDataConfig = NewCfgHeroDataConfig()
	m.CfgHeroExtraDataConfig = NewCfgHeroExtraDataConfig()
	m.CfgHeroLevelConfig = NewCfgHeroLevelConfig()
	m.CfgHeroSkillConfig = NewCfgHeroSkillConfig()
	m.CfgHeroTaskDataConfig = NewCfgHeroTaskDataConfig()
	m.CfgItemDataConfig = NewCfgItemDataConfig()
	m.CfgMailTemplateConfig = NewCfgMailTemplateConfig()
	m.CfgManualConfig = NewCfgManualConfig()
	m.CfgManualCampConfig = NewCfgManualCampConfig()
	m.CfgManualCharacterVoiceConfig = NewCfgManualCharacterVoiceConfig()
	m.CfgManualGroupConfig = NewCfgManualGroupConfig()
	m.CfgManualLevelAttributeDataConfig = NewCfgManualLevelAttributeDataConfig()
	m.CfgMaskWordConfig = NewCfgMaskWordConfig()
	m.CfgMonsterBattleConfig = NewCfgMonsterBattleConfig()
	m.CfgMonsterDataConfig = NewCfgMonsterDataConfig()
	m.CfgNamesNpcConfig = NewCfgNamesNpcConfig()
	m.CfgPassiveEffectConfig = NewCfgPassiveEffectConfig()
	m.CfgPasssiveAttributesConfig = NewCfgPasssiveAttributesConfig()
	m.CfgQuickPurchaseStaminaConfig = NewCfgQuickPurchaseStaminaConfig()
	m.CfgRarityConfig = NewCfgRarityConfig()
	m.CfgScorePassGroupConfig = NewCfgScorePassGroupConfig()
	m.CfgScorePassPhaseConfig = NewCfgScorePassPhaseConfig()
	m.CfgScorePassRewardConfig = NewCfgScorePassRewardConfig()
	m.CfgScorePassTaskConfig = NewCfgScorePassTaskConfig()
	m.CfgSigninDataConfig = NewCfgSigninDataConfig()
	m.CfgSkillConfig = NewCfgSkillConfig()
	m.CfgSkillEffectConfig = NewCfgSkillEffectConfig()
	m.CfgSkillEventConfig = NewCfgSkillEventConfig()
	m.CfgSoundConfig = NewCfgSoundConfig()
	m.CfgStoreGeneralConfig = NewCfgStoreGeneralConfig()
	m.CfgStoreGoodsConfig = NewCfgStoreGoodsConfig()
	m.CfgStoreSubstoreConfig = NewCfgStoreSubstoreConfig()
	m.CfgStoreUpdateConfig = NewCfgStoreUpdateConfig()
	m.CfgSurveyInfoConfig = NewCfgSurveyInfoConfig()
	m.CfgSurveyQuestionConfig = NewCfgSurveyQuestionConfig()
	m.CfgTargetingConfig = NewCfgTargetingConfig()
	m.CfgTaskDataConfig = NewCfgTaskDataConfig()
	m.CfgTaskRewardConfig = NewCfgTaskRewardConfig()
	m.CfgTeamLevelConfig = NewCfgTeamLevelConfig()
	m.CfgTowerConfig = NewCfgTowerConfig()
	m.CfgTowerLevelConfig = NewCfgTowerLevelConfig()
	m.CfgTowerStageConfig = NewCfgTowerStageConfig()
	m.CfgTpGateConfig = NewCfgTpGateConfig()
	m.CfgUltraSkillConfig = NewCfgUltraSkillConfig()
	m.CfgUserAvatarConfig = NewCfgUserAvatarConfig()
	m.CfgVnConfig = NewCfgVnConfig()
	m.CfgWorldItemAdvanceConfig = NewCfgWorldItemAdvanceConfig()
	m.CfgWorldItemAttributesDataConfig = NewCfgWorldItemAttributesDataConfig()
	m.CfgWorldItemDataConfig = NewCfgWorldItemDataConfig()
	m.CfgWorldItemLevelUpDataConfig = NewCfgWorldItemLevelUpDataConfig()
	m.CfgYggdrasilAreaConfig = NewCfgYggdrasilAreaConfig()
	m.CfgYggdrasilAreaCityConfig = NewCfgYggdrasilAreaCityConfig()
	m.CfgYggdrasilAreaPosConfig = NewCfgYggdrasilAreaPosConfig()
	m.CfgYggdrasilBuildingConfig = NewCfgYggdrasilBuildingConfig()
	m.CfgYggdrasilCobuildingConfig = NewCfgYggdrasilCobuildingConfig()
	m.CfgYggdrasilDailyMonsterConfig = NewCfgYggdrasilDailyMonsterConfig()
	m.CfgYggdrasilDispatchConfig = NewCfgYggdrasilDispatchConfig()
	m.CfgYggdrasilMarkConfig = NewCfgYggdrasilMarkConfig()
	m.CfgYggdrasilObjectConfig = NewCfgYggdrasilObjectConfig()
	m.CfgYggdrasilObjectStateConfig = NewCfgYggdrasilObjectStateConfig()
	m.CfgYggdrasilSimpleTalkConfig = NewCfgYggdrasilSimpleTalkConfig()
	m.CfgYggdrasilSubTaskConfig = NewCfgYggdrasilSubTaskConfig()
	m.CfgYggdrasilSubTaskEnvConfig = NewCfgYggdrasilSubTaskEnvConfig()
	m.CfgYggdrasilTaskConfig = NewCfgYggdrasilTaskConfig()
	m.CfgYggdrasilWorldConfig = NewCfgYggdrasilWorldConfig()
	m.GmConfig = NewGmConfig()
	m.ProtocolConfig = NewProtocolConfig()
}

// set CfgActionUnlock interface
func (m *ConfigManager) SetCfgActionUnlockConfig(config *CfgActionUnlockConfig) {
	m.CfgActionUnlockConfig = config
}

// reload config CfgActionUnlock interface
func (m *ConfigManager) ReloadCfgActionUnlockConfig() *CfgActionUnlockConfig {
	reloadConfig := NewCfgActionUnlockConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_action_unlock.csv") {
		glog.Info("Load ./shared/csv/data/cfg_action_unlock.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgActivity interface
func (m *ConfigManager) SetCfgActivityConfig(config *CfgActivityConfig) {
	m.CfgActivityConfig = config
}

// reload config CfgActivity interface
func (m *ConfigManager) ReloadCfgActivityConfig() *CfgActivityConfig {
	reloadConfig := NewCfgActivityConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_activity.csv") {
		glog.Info("Load ./shared/csv/data/cfg_activity.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgActivityFunc interface
func (m *ConfigManager) SetCfgActivityFuncConfig(config *CfgActivityFuncConfig) {
	m.CfgActivityFuncConfig = config
}

// reload config CfgActivityFunc interface
func (m *ConfigManager) ReloadCfgActivityFuncConfig() *CfgActivityFuncConfig {
	reloadConfig := NewCfgActivityFuncConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_activity_func.csv") {
		glog.Info("Load ./shared/csv/data/cfg_activity_func.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgBattleLevel interface
func (m *ConfigManager) SetCfgBattleLevelConfig(config *CfgBattleLevelConfig) {
	m.CfgBattleLevelConfig = config
}

// reload config CfgBattleLevel interface
func (m *ConfigManager) ReloadCfgBattleLevelConfig() *CfgBattleLevelConfig {
	reloadConfig := NewCfgBattleLevelConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_battle_level.csv") {
		glog.Info("Load ./shared/csv/data/cfg_battle_level.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgBattleNpc interface
func (m *ConfigManager) SetCfgBattleNpcConfig(config *CfgBattleNpcConfig) {
	m.CfgBattleNpcConfig = config
}

// reload config CfgBattleNpc interface
func (m *ConfigManager) ReloadCfgBattleNpcConfig() *CfgBattleNpcConfig {
	reloadConfig := NewCfgBattleNpcConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_battle_npc.csv") {
		glog.Info("Load ./shared/csv/data/cfg_battle_npc.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgBattlePosition interface
func (m *ConfigManager) SetCfgBattlePositionConfig(config *CfgBattlePositionConfig) {
	m.CfgBattlePositionConfig = config
}

// reload config CfgBattlePosition interface
func (m *ConfigManager) ReloadCfgBattlePositionConfig() *CfgBattlePositionConfig {
	reloadConfig := NewCfgBattlePositionConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_battle_position.csv") {
		glog.Info("Load ./shared/csv/data/cfg_battle_position.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgBossInfor interface
func (m *ConfigManager) SetCfgBossInforConfig(config *CfgBossInforConfig) {
	m.CfgBossInforConfig = config
}

// reload config CfgBossInfor interface
func (m *ConfigManager) ReloadCfgBossInforConfig() *CfgBossInforConfig {
	reloadConfig := NewCfgBossInforConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_boss_infor.csv") {
		glog.Info("Load ./shared/csv/data/cfg_boss_infor.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgBuff interface
func (m *ConfigManager) SetCfgBuffConfig(config *CfgBuffConfig) {
	m.CfgBuffConfig = config
}

// reload config CfgBuff interface
func (m *ConfigManager) ReloadCfgBuffConfig() *CfgBuffConfig {
	reloadConfig := NewCfgBuffConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_buff.csv") {
		glog.Info("Load ./shared/csv/data/cfg_buff.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgBullet interface
func (m *ConfigManager) SetCfgBulletConfig(config *CfgBulletConfig) {
	m.CfgBulletConfig = config
}

// reload config CfgBullet interface
func (m *ConfigManager) ReloadCfgBulletConfig() *CfgBulletConfig {
	reloadConfig := NewCfgBulletConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_bullet.csv") {
		glog.Info("Load ./shared/csv/data/cfg_bullet.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgChapterLevelAchievement interface
func (m *ConfigManager) SetCfgChapterLevelAchievementConfig(config *CfgChapterLevelAchievementConfig) {
	m.CfgChapterLevelAchievementConfig = config
}

// reload config CfgChapterLevelAchievement interface
func (m *ConfigManager) ReloadCfgChapterLevelAchievementConfig() *CfgChapterLevelAchievementConfig {
	reloadConfig := NewCfgChapterLevelAchievementConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_chapter_level_achievement.csv") {
		glog.Info("Load ./shared/csv/data/cfg_chapter_level_achievement.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgCharacter interface
func (m *ConfigManager) SetCfgCharacterConfig(config *CfgCharacterConfig) {
	m.CfgCharacterConfig = config
}

// reload config CfgCharacter interface
func (m *ConfigManager) ReloadCfgCharacterConfig() *CfgCharacterConfig {
	reloadConfig := NewCfgCharacterConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_character.csv") {
		glog.Info("Load ./shared/csv/data/cfg_character.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgCharacterBattle interface
func (m *ConfigManager) SetCfgCharacterBattleConfig(config *CfgCharacterBattleConfig) {
	m.CfgCharacterBattleConfig = config
}

// reload config CfgCharacterBattle interface
func (m *ConfigManager) ReloadCfgCharacterBattleConfig() *CfgCharacterBattleConfig {
	reloadConfig := NewCfgCharacterBattleConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_character_battle.csv") {
		glog.Info("Load ./shared/csv/data/cfg_character_battle.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgCharacterCamp interface
func (m *ConfigManager) SetCfgCharacterCampConfig(config *CfgCharacterCampConfig) {
	m.CfgCharacterCampConfig = config
}

// reload config CfgCharacterCamp interface
func (m *ConfigManager) ReloadCfgCharacterCampConfig() *CfgCharacterCampConfig {
	reloadConfig := NewCfgCharacterCampConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_character_camp.csv") {
		glog.Info("Load ./shared/csv/data/cfg_character_camp.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgCharacterCareer interface
func (m *ConfigManager) SetCfgCharacterCareerConfig(config *CfgCharacterCareerConfig) {
	m.CfgCharacterCareerConfig = config
}

// reload config CfgCharacterCareer interface
func (m *ConfigManager) ReloadCfgCharacterCareerConfig() *CfgCharacterCareerConfig {
	reloadConfig := NewCfgCharacterCareerConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_character_career.csv") {
		glog.Info("Load ./shared/csv/data/cfg_character_career.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgCharacterData interface
func (m *ConfigManager) SetCfgCharacterDataConfig(config *CfgCharacterDataConfig) {
	m.CfgCharacterDataConfig = config
}

// reload config CfgCharacterData interface
func (m *ConfigManager) ReloadCfgCharacterDataConfig() *CfgCharacterDataConfig {
	reloadConfig := NewCfgCharacterDataConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_character_data.csv") {
		glog.Info("Load ./shared/csv/data/cfg_character_data.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgCharacterFeature interface
func (m *ConfigManager) SetCfgCharacterFeatureConfig(config *CfgCharacterFeatureConfig) {
	m.CfgCharacterFeatureConfig = config
}

// reload config CfgCharacterFeature interface
func (m *ConfigManager) ReloadCfgCharacterFeatureConfig() *CfgCharacterFeatureConfig {
	reloadConfig := NewCfgCharacterFeatureConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_character_feature.csv") {
		glog.Info("Load ./shared/csv/data/cfg_character_feature.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgCharacterLevel interface
func (m *ConfigManager) SetCfgCharacterLevelConfig(config *CfgCharacterLevelConfig) {
	m.CfgCharacterLevelConfig = config
}

// reload config CfgCharacterLevel interface
func (m *ConfigManager) ReloadCfgCharacterLevelConfig() *CfgCharacterLevelConfig {
	reloadConfig := NewCfgCharacterLevelConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_character_level.csv") {
		glog.Info("Load ./shared/csv/data/cfg_character_level.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgCharacterRarity interface
func (m *ConfigManager) SetCfgCharacterRarityConfig(config *CfgCharacterRarityConfig) {
	m.CfgCharacterRarityConfig = config
}

// reload config CfgCharacterRarity interface
func (m *ConfigManager) ReloadCfgCharacterRarityConfig() *CfgCharacterRarityConfig {
	reloadConfig := NewCfgCharacterRarityConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_character_rarity.csv") {
		glog.Info("Load ./shared/csv/data/cfg_character_rarity.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgCharacterSkill interface
func (m *ConfigManager) SetCfgCharacterSkillConfig(config *CfgCharacterSkillConfig) {
	m.CfgCharacterSkillConfig = config
}

// reload config CfgCharacterSkill interface
func (m *ConfigManager) ReloadCfgCharacterSkillConfig() *CfgCharacterSkillConfig {
	reloadConfig := NewCfgCharacterSkillConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_character_skill.csv") {
		glog.Info("Load ./shared/csv/data/cfg_character_skill.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgCharacterSkin interface
func (m *ConfigManager) SetCfgCharacterSkinConfig(config *CfgCharacterSkinConfig) {
	m.CfgCharacterSkinConfig = config
}

// reload config CfgCharacterSkin interface
func (m *ConfigManager) ReloadCfgCharacterSkinConfig() *CfgCharacterSkinConfig {
	reloadConfig := NewCfgCharacterSkinConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_character_skin.csv") {
		glog.Info("Load ./shared/csv/data/cfg_character_skin.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgCharacterStage interface
func (m *ConfigManager) SetCfgCharacterStageConfig(config *CfgCharacterStageConfig) {
	m.CfgCharacterStageConfig = config
}

// reload config CfgCharacterStage interface
func (m *ConfigManager) ReloadCfgCharacterStageConfig() *CfgCharacterStageConfig {
	reloadConfig := NewCfgCharacterStageConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_character_stage.csv") {
		glog.Info("Load ./shared/csv/data/cfg_character_stage.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgCharacterStar interface
func (m *ConfigManager) SetCfgCharacterStarConfig(config *CfgCharacterStarConfig) {
	m.CfgCharacterStarConfig = config
}

// reload config CfgCharacterStar interface
func (m *ConfigManager) ReloadCfgCharacterStarConfig() *CfgCharacterStarConfig {
	reloadConfig := NewCfgCharacterStarConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_character_star.csv") {
		glog.Info("Load ./shared/csv/data/cfg_character_star.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgCombatPower interface
func (m *ConfigManager) SetCfgCombatPowerConfig(config *CfgCombatPowerConfig) {
	m.CfgCombatPowerConfig = config
}

// reload config CfgCombatPower interface
func (m *ConfigManager) ReloadCfgCombatPowerConfig() *CfgCombatPowerConfig {
	reloadConfig := NewCfgCombatPowerConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_combat_power.csv") {
		glog.Info("Load ./shared/csv/data/cfg_combat_power.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgCombatPowerAdapt interface
func (m *ConfigManager) SetCfgCombatPowerAdaptConfig(config *CfgCombatPowerAdaptConfig) {
	m.CfgCombatPowerAdaptConfig = config
}

// reload config CfgCombatPowerAdapt interface
func (m *ConfigManager) ReloadCfgCombatPowerAdaptConfig() *CfgCombatPowerAdaptConfig {
	reloadConfig := NewCfgCombatPowerAdaptConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_combat_power_adapt.csv") {
		glog.Info("Load ./shared/csv/data/cfg_combat_power_adapt.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgDropData interface
func (m *ConfigManager) SetCfgDropDataConfig(config *CfgDropDataConfig) {
	m.CfgDropDataConfig = config
}

// reload config CfgDropData interface
func (m *ConfigManager) ReloadCfgDropDataConfig() *CfgDropDataConfig {
	reloadConfig := NewCfgDropDataConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_drop_data.csv") {
		glog.Info("Load ./shared/csv/data/cfg_drop_data.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgDropGroup interface
func (m *ConfigManager) SetCfgDropGroupConfig(config *CfgDropGroupConfig) {
	m.CfgDropGroupConfig = config
}

// reload config CfgDropGroup interface
func (m *ConfigManager) ReloadCfgDropGroupConfig() *CfgDropGroupConfig {
	reloadConfig := NewCfgDropGroupConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_drop_group.csv") {
		glog.Info("Load ./shared/csv/data/cfg_drop_group.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgEquipLevelUpData interface
func (m *ConfigManager) SetCfgEquipLevelUpDataConfig(config *CfgEquipLevelUpDataConfig) {
	m.CfgEquipLevelUpDataConfig = config
}

// reload config CfgEquipLevelUpData interface
func (m *ConfigManager) ReloadCfgEquipLevelUpDataConfig() *CfgEquipLevelUpDataConfig {
	reloadConfig := NewCfgEquipLevelUpDataConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_equipLevelUp_data.csv") {
		glog.Info("Load ./shared/csv/data/cfg_equipLevelUp_data.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgEquipAdvance interface
func (m *ConfigManager) SetCfgEquipAdvanceConfig(config *CfgEquipAdvanceConfig) {
	m.CfgEquipAdvanceConfig = config
}

// reload config CfgEquipAdvance interface
func (m *ConfigManager) ReloadCfgEquipAdvanceConfig() *CfgEquipAdvanceConfig {
	reloadConfig := NewCfgEquipAdvanceConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_equip_advance.csv") {
		glog.Info("Load ./shared/csv/data/cfg_equip_advance.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgEquipAttributesData interface
func (m *ConfigManager) SetCfgEquipAttributesDataConfig(config *CfgEquipAttributesDataConfig) {
	m.CfgEquipAttributesDataConfig = config
}

// reload config CfgEquipAttributesData interface
func (m *ConfigManager) ReloadCfgEquipAttributesDataConfig() *CfgEquipAttributesDataConfig {
	reloadConfig := NewCfgEquipAttributesDataConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_equip_attributes_data.csv") {
		glog.Info("Load ./shared/csv/data/cfg_equip_attributes_data.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgEquipBreak interface
func (m *ConfigManager) SetCfgEquipBreakConfig(config *CfgEquipBreakConfig) {
	m.CfgEquipBreakConfig = config
}

// reload config CfgEquipBreak interface
func (m *ConfigManager) ReloadCfgEquipBreakConfig() *CfgEquipBreakConfig {
	reloadConfig := NewCfgEquipBreakConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_equip_break.csv") {
		glog.Info("Load ./shared/csv/data/cfg_equip_break.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgEquipData interface
func (m *ConfigManager) SetCfgEquipDataConfig(config *CfgEquipDataConfig) {
	m.CfgEquipDataConfig = config
}

// reload config CfgEquipData interface
func (m *ConfigManager) ReloadCfgEquipDataConfig() *CfgEquipDataConfig {
	reloadConfig := NewCfgEquipDataConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_equip_data.csv") {
		glog.Info("Load ./shared/csv/data/cfg_equip_data.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgEquipRandAttributes interface
func (m *ConfigManager) SetCfgEquipRandAttributesConfig(config *CfgEquipRandAttributesConfig) {
	m.CfgEquipRandAttributesConfig = config
}

// reload config CfgEquipRandAttributes interface
func (m *ConfigManager) ReloadCfgEquipRandAttributesConfig() *CfgEquipRandAttributesConfig {
	reloadConfig := NewCfgEquipRandAttributesConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_equip_rand_attributes.csv") {
		glog.Info("Load ./shared/csv/data/cfg_equip_rand_attributes.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgEquipRandAttributesValue interface
func (m *ConfigManager) SetCfgEquipRandAttributesValueConfig(config *CfgEquipRandAttributesValueConfig) {
	m.CfgEquipRandAttributesValueConfig = config
}

// reload config CfgEquipRandAttributesValue interface
func (m *ConfigManager) ReloadCfgEquipRandAttributesValueConfig() *CfgEquipRandAttributesValueConfig {
	reloadConfig := NewCfgEquipRandAttributesValueConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_equip_rand_attributes_value.csv") {
		glog.Info("Load ./shared/csv/data/cfg_equip_rand_attributes_value.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgEquipSkill interface
func (m *ConfigManager) SetCfgEquipSkillConfig(config *CfgEquipSkillConfig) {
	m.CfgEquipSkillConfig = config
}

// reload config CfgEquipSkill interface
func (m *ConfigManager) ReloadCfgEquipSkillConfig() *CfgEquipSkillConfig {
	reloadConfig := NewCfgEquipSkillConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_equip_skill.csv") {
		glog.Info("Load ./shared/csv/data/cfg_equip_skill.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgEverydayEnergyReceive interface
func (m *ConfigManager) SetCfgEverydayEnergyReceiveConfig(config *CfgEverydayEnergyReceiveConfig) {
	m.CfgEverydayEnergyReceiveConfig = config
}

// reload config CfgEverydayEnergyReceive interface
func (m *ConfigManager) ReloadCfgEverydayEnergyReceiveConfig() *CfgEverydayEnergyReceiveConfig {
	reloadConfig := NewCfgEverydayEnergyReceiveConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_everyday_energy_receive.csv") {
		glog.Info("Load ./shared/csv/data/cfg_everyday_energy_receive.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgExploreChapter interface
func (m *ConfigManager) SetCfgExploreChapterConfig(config *CfgExploreChapterConfig) {
	m.CfgExploreChapterConfig = config
}

// reload config CfgExploreChapter interface
func (m *ConfigManager) ReloadCfgExploreChapterConfig() *CfgExploreChapterConfig {
	reloadConfig := NewCfgExploreChapterConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_explore_chapter.csv") {
		glog.Info("Load ./shared/csv/data/cfg_explore_chapter.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgExploreChapterLevel interface
func (m *ConfigManager) SetCfgExploreChapterLevelConfig(config *CfgExploreChapterLevelConfig) {
	m.CfgExploreChapterLevelConfig = config
}

// reload config CfgExploreChapterLevel interface
func (m *ConfigManager) ReloadCfgExploreChapterLevelConfig() *CfgExploreChapterLevelConfig {
	reloadConfig := NewCfgExploreChapterLevelConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_explore_chapter_level.csv") {
		glog.Info("Load ./shared/csv/data/cfg_explore_chapter_level.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgExploreChapterLevelAdapt interface
func (m *ConfigManager) SetCfgExploreChapterLevelAdaptConfig(config *CfgExploreChapterLevelAdaptConfig) {
	m.CfgExploreChapterLevelAdaptConfig = config
}

// reload config CfgExploreChapterLevelAdapt interface
func (m *ConfigManager) ReloadCfgExploreChapterLevelAdaptConfig() *CfgExploreChapterLevelAdaptConfig {
	reloadConfig := NewCfgExploreChapterLevelAdaptConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_explore_chapter_level_adapt.csv") {
		glog.Info("Load ./shared/csv/data/cfg_explore_chapter_level_adapt.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgExploreChapterReward interface
func (m *ConfigManager) SetCfgExploreChapterRewardConfig(config *CfgExploreChapterRewardConfig) {
	m.CfgExploreChapterRewardConfig = config
}

// reload config CfgExploreChapterReward interface
func (m *ConfigManager) ReloadCfgExploreChapterRewardConfig() *CfgExploreChapterRewardConfig {
	reloadConfig := NewCfgExploreChapterRewardConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_explore_chapter_reward.csv") {
		glog.Info("Load ./shared/csv/data/cfg_explore_chapter_reward.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgExploreEvent interface
func (m *ConfigManager) SetCfgExploreEventConfig(config *CfgExploreEventConfig) {
	m.CfgExploreEventConfig = config
}

// reload config CfgExploreEvent interface
func (m *ConfigManager) ReloadCfgExploreEventConfig() *CfgExploreEventConfig {
	reloadConfig := NewCfgExploreEventConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_explore_event.csv") {
		glog.Info("Load ./shared/csv/data/cfg_explore_event.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgExploreFog interface
func (m *ConfigManager) SetCfgExploreFogConfig(config *CfgExploreFogConfig) {
	m.CfgExploreFogConfig = config
}

// reload config CfgExploreFog interface
func (m *ConfigManager) ReloadCfgExploreFogConfig() *CfgExploreFogConfig {
	reloadConfig := NewCfgExploreFogConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_explore_fog.csv") {
		glog.Info("Load ./shared/csv/data/cfg_explore_fog.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgExploreGather interface
func (m *ConfigManager) SetCfgExploreGatherConfig(config *CfgExploreGatherConfig) {
	m.CfgExploreGatherConfig = config
}

// reload config CfgExploreGather interface
func (m *ConfigManager) ReloadCfgExploreGatherConfig() *CfgExploreGatherConfig {
	reloadConfig := NewCfgExploreGatherConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_explore_gather.csv") {
		glog.Info("Load ./shared/csv/data/cfg_explore_gather.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgExploreMapCoordinate interface
func (m *ConfigManager) SetCfgExploreMapCoordinateConfig(config *CfgExploreMapCoordinateConfig) {
	m.CfgExploreMapCoordinateConfig = config
}

// reload config CfgExploreMapCoordinate interface
func (m *ConfigManager) ReloadCfgExploreMapCoordinateConfig() *CfgExploreMapCoordinateConfig {
	reloadConfig := NewCfgExploreMapCoordinateConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_explore_map_coordinate.csv") {
		glog.Info("Load ./shared/csv/data/cfg_explore_map_coordinate.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgExploreMonster interface
func (m *ConfigManager) SetCfgExploreMonsterConfig(config *CfgExploreMonsterConfig) {
	m.CfgExploreMonsterConfig = config
}

// reload config CfgExploreMonster interface
func (m *ConfigManager) ReloadCfgExploreMonsterConfig() *CfgExploreMonsterConfig {
	reloadConfig := NewCfgExploreMonsterConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_explore_monster.csv") {
		glog.Info("Load ./shared/csv/data/cfg_explore_monster.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgExploreNpc interface
func (m *ConfigManager) SetCfgExploreNpcConfig(config *CfgExploreNpcConfig) {
	m.CfgExploreNpcConfig = config
}

// reload config CfgExploreNpc interface
func (m *ConfigManager) ReloadCfgExploreNpcConfig() *CfgExploreNpcConfig {
	reloadConfig := NewCfgExploreNpcConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_explore_npc.csv") {
		glog.Info("Load ./shared/csv/data/cfg_explore_npc.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgExploreResource interface
func (m *ConfigManager) SetCfgExploreResourceConfig(config *CfgExploreResourceConfig) {
	m.CfgExploreResourceConfig = config
}

// reload config CfgExploreResource interface
func (m *ConfigManager) ReloadCfgExploreResourceConfig() *CfgExploreResourceConfig {
	reloadConfig := NewCfgExploreResourceConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_explore_resource.csv") {
		glog.Info("Load ./shared/csv/data/cfg_explore_resource.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgFieldSkill interface
func (m *ConfigManager) SetCfgFieldSkillConfig(config *CfgFieldSkillConfig) {
	m.CfgFieldSkillConfig = config
}

// reload config CfgFieldSkill interface
func (m *ConfigManager) ReloadCfgFieldSkillConfig() *CfgFieldSkillConfig {
	reloadConfig := NewCfgFieldSkillConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_field_skill.csv") {
		glog.Info("Load ./shared/csv/data/cfg_field_skill.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgFx interface
func (m *ConfigManager) SetCfgFxConfig(config *CfgFxConfig) {
	m.CfgFxConfig = config
}

// reload config CfgFx interface
func (m *ConfigManager) ReloadCfgFxConfig() *CfgFxConfig {
	reloadConfig := NewCfgFxConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_fx.csv") {
		glog.Info("Load ./shared/csv/data/cfg_fx.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgGachaDrop interface
func (m *ConfigManager) SetCfgGachaDropConfig(config *CfgGachaDropConfig) {
	m.CfgGachaDropConfig = config
}

// reload config CfgGachaDrop interface
func (m *ConfigManager) ReloadCfgGachaDropConfig() *CfgGachaDropConfig {
	reloadConfig := NewCfgGachaDropConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_gacha_drop.csv") {
		glog.Info("Load ./shared/csv/data/cfg_gacha_drop.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgGachaGlobal interface
func (m *ConfigManager) SetCfgGachaGlobalConfig(config *CfgGachaGlobalConfig) {
	m.CfgGachaGlobalConfig = config
}

// reload config CfgGachaGlobal interface
func (m *ConfigManager) ReloadCfgGachaGlobalConfig() *CfgGachaGlobalConfig {
	reloadConfig := NewCfgGachaGlobalConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_gacha_global.csv") {
		glog.Info("Load ./shared/csv/data/cfg_gacha_global.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgGachaPool interface
func (m *ConfigManager) SetCfgGachaPoolConfig(config *CfgGachaPoolConfig) {
	m.CfgGachaPoolConfig = config
}

// reload config CfgGachaPool interface
func (m *ConfigManager) ReloadCfgGachaPoolConfig() *CfgGachaPoolConfig {
	reloadConfig := NewCfgGachaPoolConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_gacha_pool.csv") {
		glog.Info("Load ./shared/csv/data/cfg_gacha_pool.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgGlobalString interface
func (m *ConfigManager) SetCfgGlobalStringConfig(config *CfgGlobalStringConfig) {
	m.CfgGlobalStringConfig = config
}

// reload config CfgGlobalString interface
func (m *ConfigManager) ReloadCfgGlobalStringConfig() *CfgGlobalStringConfig {
	reloadConfig := NewCfgGlobalStringConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_global_string.csv") {
		glog.Info("Load ./shared/csv/data/cfg_global_string.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgGraveyardBuild interface
func (m *ConfigManager) SetCfgGraveyardBuildConfig(config *CfgGraveyardBuildConfig) {
	m.CfgGraveyardBuildConfig = config
}

// reload config CfgGraveyardBuild interface
func (m *ConfigManager) ReloadCfgGraveyardBuildConfig() *CfgGraveyardBuildConfig {
	reloadConfig := NewCfgGraveyardBuildConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_graveyard_build.csv") {
		glog.Info("Load ./shared/csv/data/cfg_graveyard_build.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgGraveyardBuildLevel interface
func (m *ConfigManager) SetCfgGraveyardBuildLevelConfig(config *CfgGraveyardBuildLevelConfig) {
	m.CfgGraveyardBuildLevelConfig = config
}

// reload config CfgGraveyardBuildLevel interface
func (m *ConfigManager) ReloadCfgGraveyardBuildLevelConfig() *CfgGraveyardBuildLevelConfig {
	reloadConfig := NewCfgGraveyardBuildLevelConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_graveyard_build_level.csv") {
		glog.Info("Load ./shared/csv/data/cfg_graveyard_build_level.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgGraveyardBuildStage interface
func (m *ConfigManager) SetCfgGraveyardBuildStageConfig(config *CfgGraveyardBuildStageConfig) {
	m.CfgGraveyardBuildStageConfig = config
}

// reload config CfgGraveyardBuildStage interface
func (m *ConfigManager) ReloadCfgGraveyardBuildStageConfig() *CfgGraveyardBuildStageConfig {
	reloadConfig := NewCfgGraveyardBuildStageConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_graveyard_build_stage.csv") {
		glog.Info("Load ./shared/csv/data/cfg_graveyard_build_stage.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgGraveyardMainTowerLevel interface
func (m *ConfigManager) SetCfgGraveyardMainTowerLevelConfig(config *CfgGraveyardMainTowerLevelConfig) {
	m.CfgGraveyardMainTowerLevelConfig = config
}

// reload config CfgGraveyardMainTowerLevel interface
func (m *ConfigManager) ReloadCfgGraveyardMainTowerLevelConfig() *CfgGraveyardMainTowerLevelConfig {
	reloadConfig := NewCfgGraveyardMainTowerLevelConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_graveyard_main_tower_level.csv") {
		glog.Info("Load ./shared/csv/data/cfg_graveyard_main_tower_level.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgGraveyardProduceBuff interface
func (m *ConfigManager) SetCfgGraveyardProduceBuffConfig(config *CfgGraveyardProduceBuffConfig) {
	m.CfgGraveyardProduceBuffConfig = config
}

// reload config CfgGraveyardProduceBuff interface
func (m *ConfigManager) ReloadCfgGraveyardProduceBuffConfig() *CfgGraveyardProduceBuffConfig {
	reloadConfig := NewCfgGraveyardProduceBuffConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_graveyard_produce_buff.csv") {
		glog.Info("Load ./shared/csv/data/cfg_graveyard_produce_buff.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgGuideModule interface
func (m *ConfigManager) SetCfgGuideModuleConfig(config *CfgGuideModuleConfig) {
	m.CfgGuideModuleConfig = config
}

// reload config CfgGuideModule interface
func (m *ConfigManager) ReloadCfgGuideModuleConfig() *CfgGuideModuleConfig {
	reloadConfig := NewCfgGuideModuleConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_guide_module.csv") {
		glog.Info("Load ./shared/csv/data/cfg_guide_module.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgGuideStep interface
func (m *ConfigManager) SetCfgGuideStepConfig(config *CfgGuideStepConfig) {
	m.CfgGuideStepConfig = config
}

// reload config CfgGuideStep interface
func (m *ConfigManager) ReloadCfgGuideStepConfig() *CfgGuideStepConfig {
	reloadConfig := NewCfgGuideStepConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_guide_step.csv") {
		glog.Info("Load ./shared/csv/data/cfg_guide_step.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgGuildGreetings interface
func (m *ConfigManager) SetCfgGuildGreetingsConfig(config *CfgGuildGreetingsConfig) {
	m.CfgGuildGreetingsConfig = config
}

// reload config CfgGuildGreetings interface
func (m *ConfigManager) ReloadCfgGuildGreetingsConfig() *CfgGuildGreetingsConfig {
	reloadConfig := NewCfgGuildGreetingsConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_guild_greetings.csv") {
		glog.Info("Load ./shared/csv/data/cfg_guild_greetings.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgGuildLevel interface
func (m *ConfigManager) SetCfgGuildLevelConfig(config *CfgGuildLevelConfig) {
	m.CfgGuildLevelConfig = config
}

// reload config CfgGuildLevel interface
func (m *ConfigManager) ReloadCfgGuildLevelConfig() *CfgGuildLevelConfig {
	reloadConfig := NewCfgGuildLevelConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_guild_level.csv") {
		glog.Info("Load ./shared/csv/data/cfg_guild_level.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgGuildTask interface
func (m *ConfigManager) SetCfgGuildTaskConfig(config *CfgGuildTaskConfig) {
	m.CfgGuildTaskConfig = config
}

// reload config CfgGuildTask interface
func (m *ConfigManager) ReloadCfgGuildTaskConfig() *CfgGuildTaskConfig {
	reloadConfig := NewCfgGuildTaskConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_guild_task.csv") {
		glog.Info("Load ./shared/csv/data/cfg_guild_task.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgHero interface
func (m *ConfigManager) SetCfgHeroConfig(config *CfgHeroConfig) {
	m.CfgHeroConfig = config
}

// reload config CfgHero interface
func (m *ConfigManager) ReloadCfgHeroConfig() *CfgHeroConfig {
	reloadConfig := NewCfgHeroConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_hero.csv") {
		glog.Info("Load ./shared/csv/data/cfg_hero.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgHeroCampSkillFilter interface
func (m *ConfigManager) SetCfgHeroCampSkillFilterConfig(config *CfgHeroCampSkillFilterConfig) {
	m.CfgHeroCampSkillFilterConfig = config
}

// reload config CfgHeroCampSkillFilter interface
func (m *ConfigManager) ReloadCfgHeroCampSkillFilterConfig() *CfgHeroCampSkillFilterConfig {
	reloadConfig := NewCfgHeroCampSkillFilterConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_hero_camp_skill_filter.csv") {
		glog.Info("Load ./shared/csv/data/cfg_hero_camp_skill_filter.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgHeroData interface
func (m *ConfigManager) SetCfgHeroDataConfig(config *CfgHeroDataConfig) {
	m.CfgHeroDataConfig = config
}

// reload config CfgHeroData interface
func (m *ConfigManager) ReloadCfgHeroDataConfig() *CfgHeroDataConfig {
	reloadConfig := NewCfgHeroDataConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_hero_data.csv") {
		glog.Info("Load ./shared/csv/data/cfg_hero_data.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgHeroExtraData interface
func (m *ConfigManager) SetCfgHeroExtraDataConfig(config *CfgHeroExtraDataConfig) {
	m.CfgHeroExtraDataConfig = config
}

// reload config CfgHeroExtraData interface
func (m *ConfigManager) ReloadCfgHeroExtraDataConfig() *CfgHeroExtraDataConfig {
	reloadConfig := NewCfgHeroExtraDataConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_hero_extra_data.csv") {
		glog.Info("Load ./shared/csv/data/cfg_hero_extra_data.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgHeroLevel interface
func (m *ConfigManager) SetCfgHeroLevelConfig(config *CfgHeroLevelConfig) {
	m.CfgHeroLevelConfig = config
}

// reload config CfgHeroLevel interface
func (m *ConfigManager) ReloadCfgHeroLevelConfig() *CfgHeroLevelConfig {
	reloadConfig := NewCfgHeroLevelConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_hero_level.csv") {
		glog.Info("Load ./shared/csv/data/cfg_hero_level.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgHeroSkill interface
func (m *ConfigManager) SetCfgHeroSkillConfig(config *CfgHeroSkillConfig) {
	m.CfgHeroSkillConfig = config
}

// reload config CfgHeroSkill interface
func (m *ConfigManager) ReloadCfgHeroSkillConfig() *CfgHeroSkillConfig {
	reloadConfig := NewCfgHeroSkillConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_hero_skill.csv") {
		glog.Info("Load ./shared/csv/data/cfg_hero_skill.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgHeroTaskData interface
func (m *ConfigManager) SetCfgHeroTaskDataConfig(config *CfgHeroTaskDataConfig) {
	m.CfgHeroTaskDataConfig = config
}

// reload config CfgHeroTaskData interface
func (m *ConfigManager) ReloadCfgHeroTaskDataConfig() *CfgHeroTaskDataConfig {
	reloadConfig := NewCfgHeroTaskDataConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_hero_task_data.csv") {
		glog.Info("Load ./shared/csv/data/cfg_hero_task_data.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgItemData interface
func (m *ConfigManager) SetCfgItemDataConfig(config *CfgItemDataConfig) {
	m.CfgItemDataConfig = config
}

// reload config CfgItemData interface
func (m *ConfigManager) ReloadCfgItemDataConfig() *CfgItemDataConfig {
	reloadConfig := NewCfgItemDataConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_item_data.csv") {
		glog.Info("Load ./shared/csv/data/cfg_item_data.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgMailTemplate interface
func (m *ConfigManager) SetCfgMailTemplateConfig(config *CfgMailTemplateConfig) {
	m.CfgMailTemplateConfig = config
}

// reload config CfgMailTemplate interface
func (m *ConfigManager) ReloadCfgMailTemplateConfig() *CfgMailTemplateConfig {
	reloadConfig := NewCfgMailTemplateConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_mail_template.csv") {
		glog.Info("Load ./shared/csv/data/cfg_mail_template.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgManual interface
func (m *ConfigManager) SetCfgManualConfig(config *CfgManualConfig) {
	m.CfgManualConfig = config
}

// reload config CfgManual interface
func (m *ConfigManager) ReloadCfgManualConfig() *CfgManualConfig {
	reloadConfig := NewCfgManualConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_manual.csv") {
		glog.Info("Load ./shared/csv/data/cfg_manual.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgManualCamp interface
func (m *ConfigManager) SetCfgManualCampConfig(config *CfgManualCampConfig) {
	m.CfgManualCampConfig = config
}

// reload config CfgManualCamp interface
func (m *ConfigManager) ReloadCfgManualCampConfig() *CfgManualCampConfig {
	reloadConfig := NewCfgManualCampConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_manual_camp.csv") {
		glog.Info("Load ./shared/csv/data/cfg_manual_camp.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgManualCharacterVoice interface
func (m *ConfigManager) SetCfgManualCharacterVoiceConfig(config *CfgManualCharacterVoiceConfig) {
	m.CfgManualCharacterVoiceConfig = config
}

// reload config CfgManualCharacterVoice interface
func (m *ConfigManager) ReloadCfgManualCharacterVoiceConfig() *CfgManualCharacterVoiceConfig {
	reloadConfig := NewCfgManualCharacterVoiceConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_manual_character_voice.csv") {
		glog.Info("Load ./shared/csv/data/cfg_manual_character_voice.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgManualGroup interface
func (m *ConfigManager) SetCfgManualGroupConfig(config *CfgManualGroupConfig) {
	m.CfgManualGroupConfig = config
}

// reload config CfgManualGroup interface
func (m *ConfigManager) ReloadCfgManualGroupConfig() *CfgManualGroupConfig {
	reloadConfig := NewCfgManualGroupConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_manual_group.csv") {
		glog.Info("Load ./shared/csv/data/cfg_manual_group.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgManualLevelAttributeData interface
func (m *ConfigManager) SetCfgManualLevelAttributeDataConfig(config *CfgManualLevelAttributeDataConfig) {
	m.CfgManualLevelAttributeDataConfig = config
}

// reload config CfgManualLevelAttributeData interface
func (m *ConfigManager) ReloadCfgManualLevelAttributeDataConfig() *CfgManualLevelAttributeDataConfig {
	reloadConfig := NewCfgManualLevelAttributeDataConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_manual_level_attribute_data.csv") {
		glog.Info("Load ./shared/csv/data/cfg_manual_level_attribute_data.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgMaskWord interface
func (m *ConfigManager) SetCfgMaskWordConfig(config *CfgMaskWordConfig) {
	m.CfgMaskWordConfig = config
}

// reload config CfgMaskWord interface
func (m *ConfigManager) ReloadCfgMaskWordConfig() *CfgMaskWordConfig {
	reloadConfig := NewCfgMaskWordConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_mask_word.csv") {
		glog.Info("Load ./shared/csv/data/cfg_mask_word.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgMonsterBattle interface
func (m *ConfigManager) SetCfgMonsterBattleConfig(config *CfgMonsterBattleConfig) {
	m.CfgMonsterBattleConfig = config
}

// reload config CfgMonsterBattle interface
func (m *ConfigManager) ReloadCfgMonsterBattleConfig() *CfgMonsterBattleConfig {
	reloadConfig := NewCfgMonsterBattleConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_monster_battle.csv") {
		glog.Info("Load ./shared/csv/data/cfg_monster_battle.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgMonsterData interface
func (m *ConfigManager) SetCfgMonsterDataConfig(config *CfgMonsterDataConfig) {
	m.CfgMonsterDataConfig = config
}

// reload config CfgMonsterData interface
func (m *ConfigManager) ReloadCfgMonsterDataConfig() *CfgMonsterDataConfig {
	reloadConfig := NewCfgMonsterDataConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_monster_data.csv") {
		glog.Info("Load ./shared/csv/data/cfg_monster_data.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgNamesNpc interface
func (m *ConfigManager) SetCfgNamesNpcConfig(config *CfgNamesNpcConfig) {
	m.CfgNamesNpcConfig = config
}

// reload config CfgNamesNpc interface
func (m *ConfigManager) ReloadCfgNamesNpcConfig() *CfgNamesNpcConfig {
	reloadConfig := NewCfgNamesNpcConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_names_npc.csv") {
		glog.Info("Load ./shared/csv/data/cfg_names_npc.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgPassiveEffect interface
func (m *ConfigManager) SetCfgPassiveEffectConfig(config *CfgPassiveEffectConfig) {
	m.CfgPassiveEffectConfig = config
}

// reload config CfgPassiveEffect interface
func (m *ConfigManager) ReloadCfgPassiveEffectConfig() *CfgPassiveEffectConfig {
	reloadConfig := NewCfgPassiveEffectConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_passive_effect.csv") {
		glog.Info("Load ./shared/csv/data/cfg_passive_effect.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgPasssiveAttributes interface
func (m *ConfigManager) SetCfgPasssiveAttributesConfig(config *CfgPasssiveAttributesConfig) {
	m.CfgPasssiveAttributesConfig = config
}

// reload config CfgPasssiveAttributes interface
func (m *ConfigManager) ReloadCfgPasssiveAttributesConfig() *CfgPasssiveAttributesConfig {
	reloadConfig := NewCfgPasssiveAttributesConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_passsive_attributes.csv") {
		glog.Info("Load ./shared/csv/data/cfg_passsive_attributes.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgQuickPurchaseStamina interface
func (m *ConfigManager) SetCfgQuickPurchaseStaminaConfig(config *CfgQuickPurchaseStaminaConfig) {
	m.CfgQuickPurchaseStaminaConfig = config
}

// reload config CfgQuickPurchaseStamina interface
func (m *ConfigManager) ReloadCfgQuickPurchaseStaminaConfig() *CfgQuickPurchaseStaminaConfig {
	reloadConfig := NewCfgQuickPurchaseStaminaConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_quick_purchase_stamina.csv") {
		glog.Info("Load ./shared/csv/data/cfg_quick_purchase_stamina.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgRarity interface
func (m *ConfigManager) SetCfgRarityConfig(config *CfgRarityConfig) {
	m.CfgRarityConfig = config
}

// reload config CfgRarity interface
func (m *ConfigManager) ReloadCfgRarityConfig() *CfgRarityConfig {
	reloadConfig := NewCfgRarityConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_rarity.csv") {
		glog.Info("Load ./shared/csv/data/cfg_rarity.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgScorePassGroup interface
func (m *ConfigManager) SetCfgScorePassGroupConfig(config *CfgScorePassGroupConfig) {
	m.CfgScorePassGroupConfig = config
}

// reload config CfgScorePassGroup interface
func (m *ConfigManager) ReloadCfgScorePassGroupConfig() *CfgScorePassGroupConfig {
	reloadConfig := NewCfgScorePassGroupConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_score_pass_group.csv") {
		glog.Info("Load ./shared/csv/data/cfg_score_pass_group.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgScorePassPhase interface
func (m *ConfigManager) SetCfgScorePassPhaseConfig(config *CfgScorePassPhaseConfig) {
	m.CfgScorePassPhaseConfig = config
}

// reload config CfgScorePassPhase interface
func (m *ConfigManager) ReloadCfgScorePassPhaseConfig() *CfgScorePassPhaseConfig {
	reloadConfig := NewCfgScorePassPhaseConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_score_pass_phase.csv") {
		glog.Info("Load ./shared/csv/data/cfg_score_pass_phase.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgScorePassReward interface
func (m *ConfigManager) SetCfgScorePassRewardConfig(config *CfgScorePassRewardConfig) {
	m.CfgScorePassRewardConfig = config
}

// reload config CfgScorePassReward interface
func (m *ConfigManager) ReloadCfgScorePassRewardConfig() *CfgScorePassRewardConfig {
	reloadConfig := NewCfgScorePassRewardConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_score_pass_reward.csv") {
		glog.Info("Load ./shared/csv/data/cfg_score_pass_reward.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgScorePassTask interface
func (m *ConfigManager) SetCfgScorePassTaskConfig(config *CfgScorePassTaskConfig) {
	m.CfgScorePassTaskConfig = config
}

// reload config CfgScorePassTask interface
func (m *ConfigManager) ReloadCfgScorePassTaskConfig() *CfgScorePassTaskConfig {
	reloadConfig := NewCfgScorePassTaskConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_score_pass_task.csv") {
		glog.Info("Load ./shared/csv/data/cfg_score_pass_task.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgSigninData interface
func (m *ConfigManager) SetCfgSigninDataConfig(config *CfgSigninDataConfig) {
	m.CfgSigninDataConfig = config
}

// reload config CfgSigninData interface
func (m *ConfigManager) ReloadCfgSigninDataConfig() *CfgSigninDataConfig {
	reloadConfig := NewCfgSigninDataConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_signin_data.csv") {
		glog.Info("Load ./shared/csv/data/cfg_signin_data.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgSkill interface
func (m *ConfigManager) SetCfgSkillConfig(config *CfgSkillConfig) {
	m.CfgSkillConfig = config
}

// reload config CfgSkill interface
func (m *ConfigManager) ReloadCfgSkillConfig() *CfgSkillConfig {
	reloadConfig := NewCfgSkillConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_skill.csv") {
		glog.Info("Load ./shared/csv/data/cfg_skill.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgSkillEffect interface
func (m *ConfigManager) SetCfgSkillEffectConfig(config *CfgSkillEffectConfig) {
	m.CfgSkillEffectConfig = config
}

// reload config CfgSkillEffect interface
func (m *ConfigManager) ReloadCfgSkillEffectConfig() *CfgSkillEffectConfig {
	reloadConfig := NewCfgSkillEffectConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_skill_effect.csv") {
		glog.Info("Load ./shared/csv/data/cfg_skill_effect.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgSkillEvent interface
func (m *ConfigManager) SetCfgSkillEventConfig(config *CfgSkillEventConfig) {
	m.CfgSkillEventConfig = config
}

// reload config CfgSkillEvent interface
func (m *ConfigManager) ReloadCfgSkillEventConfig() *CfgSkillEventConfig {
	reloadConfig := NewCfgSkillEventConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_skill_event.csv") {
		glog.Info("Load ./shared/csv/data/cfg_skill_event.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgSound interface
func (m *ConfigManager) SetCfgSoundConfig(config *CfgSoundConfig) {
	m.CfgSoundConfig = config
}

// reload config CfgSound interface
func (m *ConfigManager) ReloadCfgSoundConfig() *CfgSoundConfig {
	reloadConfig := NewCfgSoundConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_sound.csv") {
		glog.Info("Load ./shared/csv/data/cfg_sound.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgStoreGeneral interface
func (m *ConfigManager) SetCfgStoreGeneralConfig(config *CfgStoreGeneralConfig) {
	m.CfgStoreGeneralConfig = config
}

// reload config CfgStoreGeneral interface
func (m *ConfigManager) ReloadCfgStoreGeneralConfig() *CfgStoreGeneralConfig {
	reloadConfig := NewCfgStoreGeneralConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_store_general.csv") {
		glog.Info("Load ./shared/csv/data/cfg_store_general.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgStoreGoods interface
func (m *ConfigManager) SetCfgStoreGoodsConfig(config *CfgStoreGoodsConfig) {
	m.CfgStoreGoodsConfig = config
}

// reload config CfgStoreGoods interface
func (m *ConfigManager) ReloadCfgStoreGoodsConfig() *CfgStoreGoodsConfig {
	reloadConfig := NewCfgStoreGoodsConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_store_goods.csv") {
		glog.Info("Load ./shared/csv/data/cfg_store_goods.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgStoreSubstore interface
func (m *ConfigManager) SetCfgStoreSubstoreConfig(config *CfgStoreSubstoreConfig) {
	m.CfgStoreSubstoreConfig = config
}

// reload config CfgStoreSubstore interface
func (m *ConfigManager) ReloadCfgStoreSubstoreConfig() *CfgStoreSubstoreConfig {
	reloadConfig := NewCfgStoreSubstoreConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_store_substore.csv") {
		glog.Info("Load ./shared/csv/data/cfg_store_substore.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgStoreUpdate interface
func (m *ConfigManager) SetCfgStoreUpdateConfig(config *CfgStoreUpdateConfig) {
	m.CfgStoreUpdateConfig = config
}

// reload config CfgStoreUpdate interface
func (m *ConfigManager) ReloadCfgStoreUpdateConfig() *CfgStoreUpdateConfig {
	reloadConfig := NewCfgStoreUpdateConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_store_update.csv") {
		glog.Info("Load ./shared/csv/data/cfg_store_update.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgSurveyInfo interface
func (m *ConfigManager) SetCfgSurveyInfoConfig(config *CfgSurveyInfoConfig) {
	m.CfgSurveyInfoConfig = config
}

// reload config CfgSurveyInfo interface
func (m *ConfigManager) ReloadCfgSurveyInfoConfig() *CfgSurveyInfoConfig {
	reloadConfig := NewCfgSurveyInfoConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_survey_info.csv") {
		glog.Info("Load ./shared/csv/data/cfg_survey_info.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgSurveyQuestion interface
func (m *ConfigManager) SetCfgSurveyQuestionConfig(config *CfgSurveyQuestionConfig) {
	m.CfgSurveyQuestionConfig = config
}

// reload config CfgSurveyQuestion interface
func (m *ConfigManager) ReloadCfgSurveyQuestionConfig() *CfgSurveyQuestionConfig {
	reloadConfig := NewCfgSurveyQuestionConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_survey_question.csv") {
		glog.Info("Load ./shared/csv/data/cfg_survey_question.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgTargeting interface
func (m *ConfigManager) SetCfgTargetingConfig(config *CfgTargetingConfig) {
	m.CfgTargetingConfig = config
}

// reload config CfgTargeting interface
func (m *ConfigManager) ReloadCfgTargetingConfig() *CfgTargetingConfig {
	reloadConfig := NewCfgTargetingConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_targeting.csv") {
		glog.Info("Load ./shared/csv/data/cfg_targeting.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgTaskData interface
func (m *ConfigManager) SetCfgTaskDataConfig(config *CfgTaskDataConfig) {
	m.CfgTaskDataConfig = config
}

// reload config CfgTaskData interface
func (m *ConfigManager) ReloadCfgTaskDataConfig() *CfgTaskDataConfig {
	reloadConfig := NewCfgTaskDataConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_task_data.csv") {
		glog.Info("Load ./shared/csv/data/cfg_task_data.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgTaskReward interface
func (m *ConfigManager) SetCfgTaskRewardConfig(config *CfgTaskRewardConfig) {
	m.CfgTaskRewardConfig = config
}

// reload config CfgTaskReward interface
func (m *ConfigManager) ReloadCfgTaskRewardConfig() *CfgTaskRewardConfig {
	reloadConfig := NewCfgTaskRewardConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_task_reward.csv") {
		glog.Info("Load ./shared/csv/data/cfg_task_reward.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgTeamLevel interface
func (m *ConfigManager) SetCfgTeamLevelConfig(config *CfgTeamLevelConfig) {
	m.CfgTeamLevelConfig = config
}

// reload config CfgTeamLevel interface
func (m *ConfigManager) ReloadCfgTeamLevelConfig() *CfgTeamLevelConfig {
	reloadConfig := NewCfgTeamLevelConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_team_level.csv") {
		glog.Info("Load ./shared/csv/data/cfg_team_level.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgTower interface
func (m *ConfigManager) SetCfgTowerConfig(config *CfgTowerConfig) {
	m.CfgTowerConfig = config
}

// reload config CfgTower interface
func (m *ConfigManager) ReloadCfgTowerConfig() *CfgTowerConfig {
	reloadConfig := NewCfgTowerConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_tower.csv") {
		glog.Info("Load ./shared/csv/data/cfg_tower.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgTowerLevel interface
func (m *ConfigManager) SetCfgTowerLevelConfig(config *CfgTowerLevelConfig) {
	m.CfgTowerLevelConfig = config
}

// reload config CfgTowerLevel interface
func (m *ConfigManager) ReloadCfgTowerLevelConfig() *CfgTowerLevelConfig {
	reloadConfig := NewCfgTowerLevelConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_tower_level.csv") {
		glog.Info("Load ./shared/csv/data/cfg_tower_level.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgTowerStage interface
func (m *ConfigManager) SetCfgTowerStageConfig(config *CfgTowerStageConfig) {
	m.CfgTowerStageConfig = config
}

// reload config CfgTowerStage interface
func (m *ConfigManager) ReloadCfgTowerStageConfig() *CfgTowerStageConfig {
	reloadConfig := NewCfgTowerStageConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_tower_stage.csv") {
		glog.Info("Load ./shared/csv/data/cfg_tower_stage.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgTpGate interface
func (m *ConfigManager) SetCfgTpGateConfig(config *CfgTpGateConfig) {
	m.CfgTpGateConfig = config
}

// reload config CfgTpGate interface
func (m *ConfigManager) ReloadCfgTpGateConfig() *CfgTpGateConfig {
	reloadConfig := NewCfgTpGateConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_tp_gate.csv") {
		glog.Info("Load ./shared/csv/data/cfg_tp_gate.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgUltraSkill interface
func (m *ConfigManager) SetCfgUltraSkillConfig(config *CfgUltraSkillConfig) {
	m.CfgUltraSkillConfig = config
}

// reload config CfgUltraSkill interface
func (m *ConfigManager) ReloadCfgUltraSkillConfig() *CfgUltraSkillConfig {
	reloadConfig := NewCfgUltraSkillConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_ultra_skill.csv") {
		glog.Info("Load ./shared/csv/data/cfg_ultra_skill.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgUserAvatar interface
func (m *ConfigManager) SetCfgUserAvatarConfig(config *CfgUserAvatarConfig) {
	m.CfgUserAvatarConfig = config
}

// reload config CfgUserAvatar interface
func (m *ConfigManager) ReloadCfgUserAvatarConfig() *CfgUserAvatarConfig {
	reloadConfig := NewCfgUserAvatarConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_user_avatar.csv") {
		glog.Info("Load ./shared/csv/data/cfg_user_avatar.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgVn interface
func (m *ConfigManager) SetCfgVnConfig(config *CfgVnConfig) {
	m.CfgVnConfig = config
}

// reload config CfgVn interface
func (m *ConfigManager) ReloadCfgVnConfig() *CfgVnConfig {
	reloadConfig := NewCfgVnConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_vn.csv") {
		glog.Info("Load ./shared/csv/data/cfg_vn.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgWorldItemAdvance interface
func (m *ConfigManager) SetCfgWorldItemAdvanceConfig(config *CfgWorldItemAdvanceConfig) {
	m.CfgWorldItemAdvanceConfig = config
}

// reload config CfgWorldItemAdvance interface
func (m *ConfigManager) ReloadCfgWorldItemAdvanceConfig() *CfgWorldItemAdvanceConfig {
	reloadConfig := NewCfgWorldItemAdvanceConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_world_item_advance.csv") {
		glog.Info("Load ./shared/csv/data/cfg_world_item_advance.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgWorldItemAttributesData interface
func (m *ConfigManager) SetCfgWorldItemAttributesDataConfig(config *CfgWorldItemAttributesDataConfig) {
	m.CfgWorldItemAttributesDataConfig = config
}

// reload config CfgWorldItemAttributesData interface
func (m *ConfigManager) ReloadCfgWorldItemAttributesDataConfig() *CfgWorldItemAttributesDataConfig {
	reloadConfig := NewCfgWorldItemAttributesDataConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_world_item_attributes_data.csv") {
		glog.Info("Load ./shared/csv/data/cfg_world_item_attributes_data.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgWorldItemData interface
func (m *ConfigManager) SetCfgWorldItemDataConfig(config *CfgWorldItemDataConfig) {
	m.CfgWorldItemDataConfig = config
}

// reload config CfgWorldItemData interface
func (m *ConfigManager) ReloadCfgWorldItemDataConfig() *CfgWorldItemDataConfig {
	reloadConfig := NewCfgWorldItemDataConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_world_item_data.csv") {
		glog.Info("Load ./shared/csv/data/cfg_world_item_data.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgWorldItemLevelUpData interface
func (m *ConfigManager) SetCfgWorldItemLevelUpDataConfig(config *CfgWorldItemLevelUpDataConfig) {
	m.CfgWorldItemLevelUpDataConfig = config
}

// reload config CfgWorldItemLevelUpData interface
func (m *ConfigManager) ReloadCfgWorldItemLevelUpDataConfig() *CfgWorldItemLevelUpDataConfig {
	reloadConfig := NewCfgWorldItemLevelUpDataConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_world_item_level_up_data.csv") {
		glog.Info("Load ./shared/csv/data/cfg_world_item_level_up_data.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgYggdrasilArea interface
func (m *ConfigManager) SetCfgYggdrasilAreaConfig(config *CfgYggdrasilAreaConfig) {
	m.CfgYggdrasilAreaConfig = config
}

// reload config CfgYggdrasilArea interface
func (m *ConfigManager) ReloadCfgYggdrasilAreaConfig() *CfgYggdrasilAreaConfig {
	reloadConfig := NewCfgYggdrasilAreaConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_yggdrasil_area.csv") {
		glog.Info("Load ./shared/csv/data/cfg_yggdrasil_area.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgYggdrasilAreaCity interface
func (m *ConfigManager) SetCfgYggdrasilAreaCityConfig(config *CfgYggdrasilAreaCityConfig) {
	m.CfgYggdrasilAreaCityConfig = config
}

// reload config CfgYggdrasilAreaCity interface
func (m *ConfigManager) ReloadCfgYggdrasilAreaCityConfig() *CfgYggdrasilAreaCityConfig {
	reloadConfig := NewCfgYggdrasilAreaCityConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_yggdrasil_area_city.csv") {
		glog.Info("Load ./shared/csv/data/cfg_yggdrasil_area_city.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgYggdrasilAreaPos interface
func (m *ConfigManager) SetCfgYggdrasilAreaPosConfig(config *CfgYggdrasilAreaPosConfig) {
	m.CfgYggdrasilAreaPosConfig = config
}

// reload config CfgYggdrasilAreaPos interface
func (m *ConfigManager) ReloadCfgYggdrasilAreaPosConfig() *CfgYggdrasilAreaPosConfig {
	reloadConfig := NewCfgYggdrasilAreaPosConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_yggdrasil_area_pos.csv") {
		glog.Info("Load ./shared/csv/data/cfg_yggdrasil_area_pos.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgYggdrasilBuilding interface
func (m *ConfigManager) SetCfgYggdrasilBuildingConfig(config *CfgYggdrasilBuildingConfig) {
	m.CfgYggdrasilBuildingConfig = config
}

// reload config CfgYggdrasilBuilding interface
func (m *ConfigManager) ReloadCfgYggdrasilBuildingConfig() *CfgYggdrasilBuildingConfig {
	reloadConfig := NewCfgYggdrasilBuildingConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_yggdrasil_building.csv") {
		glog.Info("Load ./shared/csv/data/cfg_yggdrasil_building.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgYggdrasilCobuilding interface
func (m *ConfigManager) SetCfgYggdrasilCobuildingConfig(config *CfgYggdrasilCobuildingConfig) {
	m.CfgYggdrasilCobuildingConfig = config
}

// reload config CfgYggdrasilCobuilding interface
func (m *ConfigManager) ReloadCfgYggdrasilCobuildingConfig() *CfgYggdrasilCobuildingConfig {
	reloadConfig := NewCfgYggdrasilCobuildingConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_yggdrasil_cobuilding.csv") {
		glog.Info("Load ./shared/csv/data/cfg_yggdrasil_cobuilding.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgYggdrasilDailyMonster interface
func (m *ConfigManager) SetCfgYggdrasilDailyMonsterConfig(config *CfgYggdrasilDailyMonsterConfig) {
	m.CfgYggdrasilDailyMonsterConfig = config
}

// reload config CfgYggdrasilDailyMonster interface
func (m *ConfigManager) ReloadCfgYggdrasilDailyMonsterConfig() *CfgYggdrasilDailyMonsterConfig {
	reloadConfig := NewCfgYggdrasilDailyMonsterConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_yggdrasil_daily_monster.csv") {
		glog.Info("Load ./shared/csv/data/cfg_yggdrasil_daily_monster.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgYggdrasilDispatch interface
func (m *ConfigManager) SetCfgYggdrasilDispatchConfig(config *CfgYggdrasilDispatchConfig) {
	m.CfgYggdrasilDispatchConfig = config
}

// reload config CfgYggdrasilDispatch interface
func (m *ConfigManager) ReloadCfgYggdrasilDispatchConfig() *CfgYggdrasilDispatchConfig {
	reloadConfig := NewCfgYggdrasilDispatchConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_yggdrasil_dispatch.csv") {
		glog.Info("Load ./shared/csv/data/cfg_yggdrasil_dispatch.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgYggdrasilMark interface
func (m *ConfigManager) SetCfgYggdrasilMarkConfig(config *CfgYggdrasilMarkConfig) {
	m.CfgYggdrasilMarkConfig = config
}

// reload config CfgYggdrasilMark interface
func (m *ConfigManager) ReloadCfgYggdrasilMarkConfig() *CfgYggdrasilMarkConfig {
	reloadConfig := NewCfgYggdrasilMarkConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_yggdrasil_mark.csv") {
		glog.Info("Load ./shared/csv/data/cfg_yggdrasil_mark.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgYggdrasilObject interface
func (m *ConfigManager) SetCfgYggdrasilObjectConfig(config *CfgYggdrasilObjectConfig) {
	m.CfgYggdrasilObjectConfig = config
}

// reload config CfgYggdrasilObject interface
func (m *ConfigManager) ReloadCfgYggdrasilObjectConfig() *CfgYggdrasilObjectConfig {
	reloadConfig := NewCfgYggdrasilObjectConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_yggdrasil_object.csv") {
		glog.Info("Load ./shared/csv/data/cfg_yggdrasil_object.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgYggdrasilObjectState interface
func (m *ConfigManager) SetCfgYggdrasilObjectStateConfig(config *CfgYggdrasilObjectStateConfig) {
	m.CfgYggdrasilObjectStateConfig = config
}

// reload config CfgYggdrasilObjectState interface
func (m *ConfigManager) ReloadCfgYggdrasilObjectStateConfig() *CfgYggdrasilObjectStateConfig {
	reloadConfig := NewCfgYggdrasilObjectStateConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_yggdrasil_object_state.csv") {
		glog.Info("Load ./shared/csv/data/cfg_yggdrasil_object_state.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgYggdrasilSimpleTalk interface
func (m *ConfigManager) SetCfgYggdrasilSimpleTalkConfig(config *CfgYggdrasilSimpleTalkConfig) {
	m.CfgYggdrasilSimpleTalkConfig = config
}

// reload config CfgYggdrasilSimpleTalk interface
func (m *ConfigManager) ReloadCfgYggdrasilSimpleTalkConfig() *CfgYggdrasilSimpleTalkConfig {
	reloadConfig := NewCfgYggdrasilSimpleTalkConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_yggdrasil_simple_talk.csv") {
		glog.Info("Load ./shared/csv/data/cfg_yggdrasil_simple_talk.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgYggdrasilSubTask interface
func (m *ConfigManager) SetCfgYggdrasilSubTaskConfig(config *CfgYggdrasilSubTaskConfig) {
	m.CfgYggdrasilSubTaskConfig = config
}

// reload config CfgYggdrasilSubTask interface
func (m *ConfigManager) ReloadCfgYggdrasilSubTaskConfig() *CfgYggdrasilSubTaskConfig {
	reloadConfig := NewCfgYggdrasilSubTaskConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_yggdrasil_sub_task.csv") {
		glog.Info("Load ./shared/csv/data/cfg_yggdrasil_sub_task.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgYggdrasilSubTaskEnv interface
func (m *ConfigManager) SetCfgYggdrasilSubTaskEnvConfig(config *CfgYggdrasilSubTaskEnvConfig) {
	m.CfgYggdrasilSubTaskEnvConfig = config
}

// reload config CfgYggdrasilSubTaskEnv interface
func (m *ConfigManager) ReloadCfgYggdrasilSubTaskEnvConfig() *CfgYggdrasilSubTaskEnvConfig {
	reloadConfig := NewCfgYggdrasilSubTaskEnvConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_yggdrasil_sub_task_env.csv") {
		glog.Info("Load ./shared/csv/data/cfg_yggdrasil_sub_task_env.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgYggdrasilTask interface
func (m *ConfigManager) SetCfgYggdrasilTaskConfig(config *CfgYggdrasilTaskConfig) {
	m.CfgYggdrasilTaskConfig = config
}

// reload config CfgYggdrasilTask interface
func (m *ConfigManager) ReloadCfgYggdrasilTaskConfig() *CfgYggdrasilTaskConfig {
	reloadConfig := NewCfgYggdrasilTaskConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_yggdrasil_task.csv") {
		glog.Info("Load ./shared/csv/data/cfg_yggdrasil_task.csv error")
		return nil
	}
	return reloadConfig
}

// set CfgYggdrasilWorld interface
func (m *ConfigManager) SetCfgYggdrasilWorldConfig(config *CfgYggdrasilWorldConfig) {
	m.CfgYggdrasilWorldConfig = config
}

// reload config CfgYggdrasilWorld interface
func (m *ConfigManager) ReloadCfgYggdrasilWorldConfig() *CfgYggdrasilWorldConfig {
	reloadConfig := NewCfgYggdrasilWorldConfig()
	if !reloadConfig.Load("./shared/csv/data/cfg_yggdrasil_world.csv") {
		glog.Info("Load ./shared/csv/data/cfg_yggdrasil_world.csv error")
		return nil
	}
	return reloadConfig
}

// set Gm interface
func (m *ConfigManager) SetGmConfig(config *GmConfig) {
	m.GmConfig = config
}

// reload config Gm interface
func (m *ConfigManager) ReloadGmConfig() *GmConfig {
	reloadConfig := NewGmConfig()
	if !reloadConfig.Load("./shared/csv/data/gm.csv") {
		glog.Info("Load ./shared/csv/data/gm.csv error")
		return nil
	}
	return reloadConfig
}

// set Protocol interface
func (m *ConfigManager) SetProtocolConfig(config *ProtocolConfig) {
	m.ProtocolConfig = config
}

// reload config Protocol interface
func (m *ConfigManager) ReloadProtocolConfig() *ProtocolConfig {
	reloadConfig := NewProtocolConfig()
	if !reloadConfig.Load("./shared/csv/data/protocol.csv") {
		glog.Info("Load ./shared/csv/data/protocol.csv error")
		return nil
	}
	return reloadConfig
}

func (m *ConfigManager) LoadConfig(path string) bool {
	config2Loader := []struct {
		path   string
		loader Config
	}{
		{"cfg_action_unlock.csv", m.CfgActionUnlockConfig},
		{"cfg_activity.csv", m.CfgActivityConfig},
		{"cfg_activity_func.csv", m.CfgActivityFuncConfig},
		{"cfg_battle_level.csv", m.CfgBattleLevelConfig},
		{"cfg_battle_npc.csv", m.CfgBattleNpcConfig},
		{"cfg_battle_position.csv", m.CfgBattlePositionConfig},
		{"cfg_boss_infor.csv", m.CfgBossInforConfig},
		{"cfg_buff.csv", m.CfgBuffConfig},
		{"cfg_bullet.csv", m.CfgBulletConfig},
		{"cfg_chapter_level_achievement.csv", m.CfgChapterLevelAchievementConfig},
		{"cfg_character.csv", m.CfgCharacterConfig},
		{"cfg_character_battle.csv", m.CfgCharacterBattleConfig},
		{"cfg_character_camp.csv", m.CfgCharacterCampConfig},
		{"cfg_character_career.csv", m.CfgCharacterCareerConfig},
		{"cfg_character_data.csv", m.CfgCharacterDataConfig},
		{"cfg_character_feature.csv", m.CfgCharacterFeatureConfig},
		{"cfg_character_level.csv", m.CfgCharacterLevelConfig},
		{"cfg_character_rarity.csv", m.CfgCharacterRarityConfig},
		{"cfg_character_skill.csv", m.CfgCharacterSkillConfig},
		{"cfg_character_skin.csv", m.CfgCharacterSkinConfig},
		{"cfg_character_stage.csv", m.CfgCharacterStageConfig},
		{"cfg_character_star.csv", m.CfgCharacterStarConfig},
		{"cfg_combat_power.csv", m.CfgCombatPowerConfig},
		{"cfg_combat_power_adapt.csv", m.CfgCombatPowerAdaptConfig},
		{"cfg_drop_data.csv", m.CfgDropDataConfig},
		{"cfg_drop_group.csv", m.CfgDropGroupConfig},
		{"cfg_equipLevelUp_data.csv", m.CfgEquipLevelUpDataConfig},
		{"cfg_equip_advance.csv", m.CfgEquipAdvanceConfig},
		{"cfg_equip_attributes_data.csv", m.CfgEquipAttributesDataConfig},
		{"cfg_equip_break.csv", m.CfgEquipBreakConfig},
		{"cfg_equip_data.csv", m.CfgEquipDataConfig},
		{"cfg_equip_rand_attributes.csv", m.CfgEquipRandAttributesConfig},
		{"cfg_equip_rand_attributes_value.csv", m.CfgEquipRandAttributesValueConfig},
		{"cfg_equip_skill.csv", m.CfgEquipSkillConfig},
		{"cfg_everyday_energy_receive.csv", m.CfgEverydayEnergyReceiveConfig},
		{"cfg_explore_chapter.csv", m.CfgExploreChapterConfig},
		{"cfg_explore_chapter_level.csv", m.CfgExploreChapterLevelConfig},
		{"cfg_explore_chapter_level_adapt.csv", m.CfgExploreChapterLevelAdaptConfig},
		{"cfg_explore_chapter_reward.csv", m.CfgExploreChapterRewardConfig},
		{"cfg_explore_event.csv", m.CfgExploreEventConfig},
		{"cfg_explore_fog.csv", m.CfgExploreFogConfig},
		{"cfg_explore_gather.csv", m.CfgExploreGatherConfig},
		{"cfg_explore_map_coordinate.csv", m.CfgExploreMapCoordinateConfig},
		{"cfg_explore_monster.csv", m.CfgExploreMonsterConfig},
		{"cfg_explore_npc.csv", m.CfgExploreNpcConfig},
		{"cfg_explore_resource.csv", m.CfgExploreResourceConfig},
		{"cfg_field_skill.csv", m.CfgFieldSkillConfig},
		{"cfg_fx.csv", m.CfgFxConfig},
		{"cfg_gacha_drop.csv", m.CfgGachaDropConfig},
		{"cfg_gacha_global.csv", m.CfgGachaGlobalConfig},
		{"cfg_gacha_pool.csv", m.CfgGachaPoolConfig},
		{"cfg_global_string.csv", m.CfgGlobalStringConfig},
		{"cfg_graveyard_build.csv", m.CfgGraveyardBuildConfig},
		{"cfg_graveyard_build_level.csv", m.CfgGraveyardBuildLevelConfig},
		{"cfg_graveyard_build_stage.csv", m.CfgGraveyardBuildStageConfig},
		{"cfg_graveyard_main_tower_level.csv", m.CfgGraveyardMainTowerLevelConfig},
		{"cfg_graveyard_produce_buff.csv", m.CfgGraveyardProduceBuffConfig},
		{"cfg_guide_module.csv", m.CfgGuideModuleConfig},
		{"cfg_guide_step.csv", m.CfgGuideStepConfig},
		{"cfg_guild_greetings.csv", m.CfgGuildGreetingsConfig},
		{"cfg_guild_level.csv", m.CfgGuildLevelConfig},
		{"cfg_guild_task.csv", m.CfgGuildTaskConfig},
		{"cfg_hero.csv", m.CfgHeroConfig},
		{"cfg_hero_camp_skill_filter.csv", m.CfgHeroCampSkillFilterConfig},
		{"cfg_hero_data.csv", m.CfgHeroDataConfig},
		{"cfg_hero_extra_data.csv", m.CfgHeroExtraDataConfig},
		{"cfg_hero_level.csv", m.CfgHeroLevelConfig},
		{"cfg_hero_skill.csv", m.CfgHeroSkillConfig},
		{"cfg_hero_task_data.csv", m.CfgHeroTaskDataConfig},
		{"cfg_item_data.csv", m.CfgItemDataConfig},
		{"cfg_mail_template.csv", m.CfgMailTemplateConfig},
		{"cfg_manual.csv", m.CfgManualConfig},
		{"cfg_manual_camp.csv", m.CfgManualCampConfig},
		{"cfg_manual_character_voice.csv", m.CfgManualCharacterVoiceConfig},
		{"cfg_manual_group.csv", m.CfgManualGroupConfig},
		{"cfg_manual_level_attribute_data.csv", m.CfgManualLevelAttributeDataConfig},
		{"cfg_mask_word.csv", m.CfgMaskWordConfig},
		{"cfg_monster_battle.csv", m.CfgMonsterBattleConfig},
		{"cfg_monster_data.csv", m.CfgMonsterDataConfig},
		{"cfg_names_npc.csv", m.CfgNamesNpcConfig},
		{"cfg_passive_effect.csv", m.CfgPassiveEffectConfig},
		{"cfg_passsive_attributes.csv", m.CfgPasssiveAttributesConfig},
		{"cfg_quick_purchase_stamina.csv", m.CfgQuickPurchaseStaminaConfig},
		{"cfg_rarity.csv", m.CfgRarityConfig},
		{"cfg_score_pass_group.csv", m.CfgScorePassGroupConfig},
		{"cfg_score_pass_phase.csv", m.CfgScorePassPhaseConfig},
		{"cfg_score_pass_reward.csv", m.CfgScorePassRewardConfig},
		{"cfg_score_pass_task.csv", m.CfgScorePassTaskConfig},
		{"cfg_signin_data.csv", m.CfgSigninDataConfig},
		{"cfg_skill.csv", m.CfgSkillConfig},
		{"cfg_skill_effect.csv", m.CfgSkillEffectConfig},
		{"cfg_skill_event.csv", m.CfgSkillEventConfig},
		{"cfg_sound.csv", m.CfgSoundConfig},
		{"cfg_store_general.csv", m.CfgStoreGeneralConfig},
		{"cfg_store_goods.csv", m.CfgStoreGoodsConfig},
		{"cfg_store_substore.csv", m.CfgStoreSubstoreConfig},
		{"cfg_store_update.csv", m.CfgStoreUpdateConfig},
		{"cfg_survey_info.csv", m.CfgSurveyInfoConfig},
		{"cfg_survey_question.csv", m.CfgSurveyQuestionConfig},
		{"cfg_targeting.csv", m.CfgTargetingConfig},
		{"cfg_task_data.csv", m.CfgTaskDataConfig},
		{"cfg_task_reward.csv", m.CfgTaskRewardConfig},
		{"cfg_team_level.csv", m.CfgTeamLevelConfig},
		{"cfg_tower.csv", m.CfgTowerConfig},
		{"cfg_tower_level.csv", m.CfgTowerLevelConfig},
		{"cfg_tower_stage.csv", m.CfgTowerStageConfig},
		{"cfg_tp_gate.csv", m.CfgTpGateConfig},
		{"cfg_ultra_skill.csv", m.CfgUltraSkillConfig},
		{"cfg_user_avatar.csv", m.CfgUserAvatarConfig},
		{"cfg_vn.csv", m.CfgVnConfig},
		{"cfg_world_item_advance.csv", m.CfgWorldItemAdvanceConfig},
		{"cfg_world_item_attributes_data.csv", m.CfgWorldItemAttributesDataConfig},
		{"cfg_world_item_data.csv", m.CfgWorldItemDataConfig},
		{"cfg_world_item_level_up_data.csv", m.CfgWorldItemLevelUpDataConfig},
		{"cfg_yggdrasil_area.csv", m.CfgYggdrasilAreaConfig},
		{"cfg_yggdrasil_area_city.csv", m.CfgYggdrasilAreaCityConfig},
		{"cfg_yggdrasil_area_pos.csv", m.CfgYggdrasilAreaPosConfig},
		{"cfg_yggdrasil_building.csv", m.CfgYggdrasilBuildingConfig},
		{"cfg_yggdrasil_cobuilding.csv", m.CfgYggdrasilCobuildingConfig},
		{"cfg_yggdrasil_daily_monster.csv", m.CfgYggdrasilDailyMonsterConfig},
		{"cfg_yggdrasil_dispatch.csv", m.CfgYggdrasilDispatchConfig},
		{"cfg_yggdrasil_mark.csv", m.CfgYggdrasilMarkConfig},
		{"cfg_yggdrasil_object.csv", m.CfgYggdrasilObjectConfig},
		{"cfg_yggdrasil_object_state.csv", m.CfgYggdrasilObjectStateConfig},
		{"cfg_yggdrasil_simple_talk.csv", m.CfgYggdrasilSimpleTalkConfig},
		{"cfg_yggdrasil_sub_task.csv", m.CfgYggdrasilSubTaskConfig},
		{"cfg_yggdrasil_sub_task_env.csv", m.CfgYggdrasilSubTaskEnvConfig},
		{"cfg_yggdrasil_task.csv", m.CfgYggdrasilTaskConfig},
		{"cfg_yggdrasil_world.csv", m.CfgYggdrasilWorldConfig},
		{"gm.csv", m.GmConfig},
		{"protocol.csv", m.ProtocolConfig},
	}
	// range config map.
	for i := 0; i < len(config2Loader); i++ {
		if !config2Loader[i].loader.Load(path + "/" + config2Loader[i].path) {
			glog.Info("Load ", config2Loader[i].path, " file error")
			return false
		} else {
			glog.Info("Load ", config2Loader[i].path, " file success")
		}
	}
	return true
}

var ConfigManagerObj ConfigManager

func init() {
	ConfigManagerObj.Init()
}
