package common

import (
	"shared/utility/errors"
)

const (
	ErrCodeDefault                = -1
	ErrCodeTokenInvalid           = 5
	ErrCodeUserCacheInvalid       = 11
	ErrCodeUserNotFound           = 12
	ErrCodeWhileList              = 13
	ErrCodeUserLoginInOtherClient = 14
	ErrCodeLoginSdkError          = 15
	ErrCodeUnknownSdk             = 16
	ErrCodeRegisterLimit          = 17
	ErrCodeAccountNotFound        = 18
	ErrCodeAccountRepeat          = 19
	ErrCodeSeverMaintaining       = 20

	ErrCodePatchCfgNotFound = 101
	ErrCodeAppCfgNotFound   = 102

	ErrCodeParamError        = 1001
	ErrCodeNoPermissionError = 1002
	ErrCSVNotFound           = 1003

	ErrCodeItemNotEnough         = 1101
	ErrCodeItemTypeCannotConsume = 1102
	ErrCodeItemTypeCannotUse     = 1103

	// user error code
	ErrCodeUserLevelNotArrival        = 1301
	ErrCodeUserNicknameRepeatedToMany = 1302
	ErrCodeUserNicknameNotChange      = 1303

	// character error code
	ErrCodeCharacterNotFound           = 1501
	ErrCodeCharacterSkillLevelUnlock   = 1502
	ErrCodeCharacterLevelMax           = 1503
	ErrCodeCharacterSkillCannotLevelUp = 1504
	ErrCodeCharacterWearCountInvalid   = 1505
	ErrCodeCharacterWearRepeatPart     = 1506
	ErrCodeCharacterLevelNotArrival    = 1507
	ErrCodeCharacterStarNotArrival     = 1508

	// quest error code
	ErrCodeQuestNotFound                   = 1601
	ErrCodeQuestConfigNotFound             = 1602
	ErrCodeQuestUpdateParamInvalid         = 1603
	ErrCodeQuestNotRegistered              = 1604
	ErrCodeQuestProgressNotArrival         = 1605
	ErrCodeQuestNotProgressing             = 1606
	ErrCodeQuestActivityConfigNotFound     = 1607
	ErrCodeQuestActivityRewardGot          = 1608
	ErrCodeQuestActivityProgressNotArrival = 1609
	ErrCodeQuestUpdateParamCountInvalid    = 1610
	ErrCodeQuestNotComplete                = 1611

	ErrCodeGraveyardTypeCannotBuild               = 1701
	ErrCodeGraveyardNoAreaCannotBuild             = 1702
	ErrCodeGraveyardNumLimitCannotBuild           = 1703
	ErrCodeGraveyardBuildNotExist                 = 1704
	ErrCodeGraveyardBuildNotInTransaction         = 1705
	ErrCodeGraveyardBuildTransactionTimeNotEnough = 1706
	ErrCodeGraveyardBuildInProduct                = 1707
	ErrCodeGraveyardBuildInTransaction            = 1708
	ErrCodeGraveyardMainTowerLvLimit              = 1709
	ErrCodeGraveyardBuildCountLimit               = 1710
	ErrCodeGraveyardProduceTimeLimit              = 1711
	ErrCodeGraveyardBuildCanNotLvUp               = 1712
	ErrCodeGraveyardBuildMaxLv                    = 1713
	ErrCodeGraveyardBuildCanNotStageUp            = 1714
	ErrCodeGraveyardBuildMaxStage                 = 1715
	ErrCodeGraveyardPopulationNumOverflow         = 1716
	ErrCodeGraveyardCharacterNumOverflow          = 1717
	ErrCodeGraveyardTypeCannotProduce             = 1718
	ErrCodeGraveyardProduceTimeNotEnough          = 1719
	ErrCodeGraveyardProduceNumLimit               = 1720
	ErrCodeGraveyardBuildCannotReduceProductTime  = 1721
	ErrCodeGraveyardBuffInUse                     = 1722
	ErrCodeGraveyardCannotUseBuffBuildsNotExist   = 1723
	ErrCodeGraveyardNoRemainedHelpCount           = 1724
	ErrCodeGraveyardBuildInHelpNow                = 1725
	ErrCodeGraveyardCannotPlotRewardNow           = 1726
	ErrCodeGraveyardPlotRewardNumMax              = 1727

	ErrCodeEquipmentNotExist                  = 1801
	ErrCodeEquipmentNotMatchCareer            = 1802
	ErrCodeEquipmentEXPUpToLimit              = 1803
	ErrCodeEquipmentAdvanceMaterialNotMatch   = 1804
	ErrCodeEquipmentAdvanceMaterialNoEnough   = 1805
	ErrCodeEquipmentLocked                    = 1806
	ErrCodeEquipmentAdvanceLevelNotEnough     = 1807
	ErrCodeEquipmentNoCamp                    = 1808
	ErrCodeEquipmentHasNotConfirm             = 1809
	ErrCodeEquipmentNotRecast                 = 1810
	ErrCodeEquipmentStrengthenMaterialInvalid = 1811
	ErrCodeEquipmentStageUpToLimit            = 1812
	ErrCodeEquipmentCantUseSelf               = 1813
	ErrCodeEquipmentWearing                   = 1814

	ErrCodeWorldItemNotExist                = 1901
	ErrCodeWorldItemNotMatchCareer          = 1902
	ErrCodeWorldItemEXPUpToLimit            = 1903
	ErrCodeWorldItemAdvanceMaterialNotMatch = 1904
	ErrCodeWorldItemAdvanceMaterialNoEnough = 1905
	ErrCodeWorldItemLocked                  = 1906
	ErrCodeWorldItemAdvanceLevelNotEnough   = 1907
	ErrCodeWorldItemStageUpToLimit          = 1908
	ErrCodeWorldItemCantUseSelf             = 1909
	ErrCodeWorldItemWearing                 = 1910
	// ErrCodeWorldItemNoCamp                  = 1908
	// ErrCodeWorldItemHasNotConfirm           = 1909

	//signin error code
	ErrCodeSignInWrongIDForSignInGroups = 2001
	ErrCodeSignInWrongIDForSignInData   = 2002
	ErrCodeSignInRepeat                 = 2003
	ErrCodeSignInDropIDMismatchDayCnt   = 2004
	ErrCodeSignInIndexOutOfDropID       = 2005
	ErrCodeSignInUserHasSigned          = 2006
	ErrCodeSignInNoDataForToday         = 2007
	ErrCodeSignActionUnlock             = 2008

	// level error code
	ErrCodeLevelConfigNotFound       = 2101
	ErrCodeLevelNotStarted           = 2103
	ErrCodeOtherLevelStarted         = 2104
	ErrCodeLevelNotPassed            = 2105
	ErrCodeExploreMapObjectNotLevel  = 2106
	ErrCodeLevelPassTimesLimited     = 2107
	ErrCodeLevelTargetNotCompleteAll = 2108
	ErrCodeLevelInvalidSystemParam   = 2109
	ErrCodeLevelInvalidTarget        = 2110
	ErrCodeLevelInvalidAchievement   = 2111

	// explore err code
	ErrCodeChapterConfigNotFound            = 2201
	ErrCodeNotInChapter                     = 2202
	ErrCodeChapterRewardConfigNotFound      = 2203
	ErrCodeChapterScoreNotArrival           = 2204
	ErrCodeChapterRewardReceived            = 2205
	ErrCodeNotCampainLevel                  = 2206
	ErrCodeNotExploreEliteLevel             = 2207
	ErrCodeMapObjectConfigNotFound          = 2212
	ErrCodeExploreEventPointConfigNotFound  = 2223
	ErrCodeExploreInteracted                = 2224
	ErrCodeExploreNotInteracted             = 2225
	ErrCodeExploreNPCConfigNotFound         = 2234
	ErrCodeExploreNPCOptionInvalid          = 2235
	ErrCodeExploreRewardPointConfigNotFound = 2241
	ErrCodeExploreMonsterConfigNotFound     = 2251
	ErrCodeExploreMonsterHasNotLevel        = 2252
	ErrCodeExploreFogConfigNotFound         = 2261
	ErrCodeExploreFogLocked                 = 2262
	ErrCodeExploreResourceConfigNotFound    = 2271
	ErrCodeExploreResourceIsCollecting      = 2272
	ErrCodeExploreResourceCollectTimesLimit = 2273
	ErrCodeExploreResourceNotCollecting     = 2274
	ErrCodeExploreResourceNotCollected      = 2275
	ErrCodeExploreResourceHasNotLevel       = 2276
	ErrCodeExploreResourceLevelPassed       = 2277
	ErrCodeExploreTPGateConfigNotFound      = 2281
	ErrCodeExploreTPGateUseTimesLimited     = 2282

	// survey
	ErrCodeSurveyCntOutOfBitNumber           = 2401
	ErrCodeSurveyDuplicateSurveyData         = 2402
	ErrCodeSurveyQuestionCntBeyondLimit      = 2403
	ErrCodeSurveyAnswerTextBeyondLimit       = 2404
	ErrCodeSurveyUnlockConditionNotSatisfied = 2405
	ErrCodeSurveyValidityPeriodBeyondLimit   = 2406
	ErrCodeSurveyTypeBeyondLimit             = 2407
	ErrCodeSurveyWrongIDForSurveyInfos       = 2408

	// hero
	ErrCodeHeroConfigNotFound          = 2501
	ErrCodeHeroLevelConfigNotFound     = 2502
	ErrCodeHeroExpNotEnough            = 2503
	ErrCodeHeroLevelMax                = 2504
	ErrCodeHeroUnlocked                = 2505
	ErrCodeHeroNotFound                = 2506
	ErrCodeHeroHasNotSkill             = 2507
	ErrCodeHeroSkillConfigNotFound     = 2508
	ErrCodeHeroSkillLevelNotFound      = 2509
	ErrCodeHeroAttendantConfigNotFound = 2510
	ErrCodeHeroLevelNotArrival         = 2511
	ErrCodeHeroSkillLevelNotArrival    = 2512
	ErrCodeHeroSkillItemUsedNotArrival = 2513

	// store
	ErrCodeStoreWrongSubStoreTypeForData            = 2601
	ErrCodeStoreUpdateTimesNotMatchCntOfUpdateCost  = 2602
	ErrCodeStoreWrongStoreIDForData                 = 2603
	ErrCodeStoreWrongSubStoreIDForData              = 2604
	ErrCodeStoreWrongGoodsIDForData                 = 2605
	ErrCodeStoreWrongUpdateRuleIDForData            = 2606
	ErrCodeStoreWrongSubStoreIDInStoreForInfo       = 2607
	ErrCodeStoreIndexOutOfRangeForInfo              = 2608
	ErrCodeStoreWrongGoodsIDInCellForInfo           = 2609
	ErrCodeStoreNumOfRemainsNotEnough               = 2610
	ErrCodeStoreSoldOutConditionNotSatisfied        = 2611
	ErrCodeStoreNoOptionForSelectSubStore           = 2612
	ErrCodeStoreNoOptionForSelectSubStores          = 2613
	ErrCodeStoreUnlockConditionNotSatisfied         = 2614
	ErrCodeStoreSubStoreUnlockConditionNotSatisfied = 2615
	ErrCodeStoreCellUnlockConditionNotSatisfied     = 2616
	ErrCodeStoreGoodsUnlockConditionNotSatisfied    = 2617
	ErrCodeStoreWrongStoreIDForStoreInfo            = 2618
	ErrCodeStoreNotSupportUpdate                    = 2619
	ErrCodeStoreExceedMaxUpdateTimes                = 2620
	ErrCodeStoreCurrencyNotMatchDuringPurchase      = 2621
	ErrCodeStoreCurrencyNotMatchPriceForData        = 2622
	ErrCodeStoreNoGoodsInCellToGenerate             = 2623
	ErrCodeStoreWrongSubStoreIDForInfo              = 2624
	ErrCodeStoreWrongGoodsIDOfQuickPurchaseForData  = 2625

	// quick purchase error code
	ErrCodeQuickPurchaseStaminaIndexOutOfRangeForData = 2701
	ErrCodeQuickPurchaseStaminaExceedMaxTimes         = 2702

	// vn
	ErrCodeVnRewarded      = 2801
	ErrCodeVnRewardIsEmpty = 2802

	// guide
	ErrCodeGuidePassed  = 2901
	ErrCodeGuideNotPass = 2902

	// portal
	ErrCodeTooManyRequest = 3001
	ErrCodeTooManyConn    = 3002

	// tower
	ErrCodeTowerConfigNotFound      = 3101
	ErrCodeTowerStageConfigNotFound = 3102
	ErrCodeTowerStageNotArrival     = 3103
	ErrCodeTowerNotActived          = 3104
	ErrCodeTowerInvalidLevel        = 3105
	ErrCodeTowerGoUpLimited         = 3106
	ErrCodeTowerCharaCampLimited    = 3107

	// mail
	ErrCodeMailNotFound               = 3201
	ErrCodeMailHasNotAttachment       = 3202
	ErrCodeMailIsNotRemovable         = 3203
	ErrCodeMailTemplateConfigNotFound = 3204

	// guild
	ErrCodeGuildHasJoined                 = 3301
	ErrCodeGuildHasNotJoin                = 3302
	ErrCodeGuildQuitCD                    = 3303
	ErrCodeGuildNoPrivilege               = 3304
	ErrCodeGuildNotMember                 = 3305
	ErrCodeGuildChairmanCantQuit          = 3306
	ErrCodeGuildIsFull                    = 3307
	ErrCodeGuildInDissolving              = 3308
	ErrCodeGuildNotFound                  = 3309
	ErrCodeGuildNotInAppliedList          = 3310
	ErrCodeGuildYggWrongIDForCoBuild      = 3311
	ErrCodeGuildDissolved                 = 3312
	ErrCodeGuildMemberEquipmentNotExist   = 3313
	ErrCodeGuildMemberWorldItemNotExist   = 3314
	ErrCodeGuildExceedApplyNumLimit       = 3315
	ErrCodeGuildKickedMemberNotCommon     = 3316
	ErrCodeGuildTaskRewardsHasReceived    = 3317
	ErrCodeGuildTaskConfigNotFound        = 3318
	ErrCodeGuildTaskNotFinish             = 3319
	ErrCodeGuildSendGroupMailInCD         = 3320
	ErrCodeGuildModifyIconIsNull          = 3321
	ErrCodeGuildRepeatName                = 3322
	ErrCodeGuildViceCharimanNumOutOfLimit = 3323
	ErrCodeGuildEliteNumOutOfLimit        = 3324
	ErrCodeGuildInDissolvedCD             = 3325

	// combat power calculate
	ErrCodeCompileLuaFileFailed           = 3401
	ErrCodePCallFunctionProtoFailed       = 3402
	ErrCodeCallGlobalFunctionFailed       = 3403
	ErrCodeRetNotLNumber                  = 3404
	ErrCodePowerAdaptWrongNumberForSymbol = 3405
	ErrCodePowerAdaptWrongCareerIDForParm = 3406
	ErrCodePowerAdaptWrongSymbolForParm   = 3407

	// action unlock error code
	ErrCodeActionUnlockNotsatisfied  = 3501
	ErrCodeActionUnlockWrongActionID = 3502

	// yggdrasil

	ErrCodeYggdrasilInTravel                                 = 3601
	ErrCodeYggdrasilNotInTravel                              = 3602
	ErrCodeYggdrasilCharacterCannotCarry                     = 3603
	ErrCodeYggdrasilCityCannotTravel                         = 3604
	ErrCodeYggdrasilCannotMoveDistanceIllegal                = 3605
	ErrCodeYggdrasilCannotMovePosUnWalkable                  = 3606
	ErrCodeYggdrasilCannotMoveTerrainDiff                    = 3607
	ErrCodeYggdrasilApNotEnough                              = 3609
	ErrCodeYggdrasilNoTravelTime                             = 3610
	ErrCodeYggdrasilCannotReturnCityThisPos                  = 3611
	ErrCodeYggdrasilThisCityCannotTransfer                   = 3612
	ErrCodeYggdrasilInCityNow                                = 3613
	ErrCodeYggdrasilNotInCityNow                             = 3614
	ErrCodeYggdrasilPackGoodsNotFound                        = 3615
	ErrCodeYggdrasilOverlapping                              = 3616
	ErrCodeYggdrasilDiscardGoodsNotFound                     = 3617
	ErrCodeYggdrasilBagIsFull                                = 3618
	ErrCodeYggdrasilObjectNotFound                           = 3619
	ErrCodeYggdrasilTaskProcessError                         = 3620
	ErrCodeYggdrasilCompleteTaskBefore                       = 3621
	ErrCodeYggdrasilAcceptTaskBefore                         = 3622
	ErrCodeYggdrasilSameTaskGroup                            = 3623
	ErrCodeYggdrasilNotDoingTaskNow                          = 3624
	ErrCodeYggdrasilOtherMaxTaskCount                        = 3625
	ErrCodeYggdrasilTaskCannotComplete                       = 3626
	ErrCodeYggdrasilTaskNotInProgress                        = 3627
	ErrCodeYggdrasilPosAlreadyHasEntity                      = 3628
	ErrCodeYggdrasilPosAlreadyHasObject                      = 3629
	ErrCodeYggdrasilPosTypeForbidBuild                       = 3630
	ErrCodeYggdrasilPosInsideCityBanRadius                   = 3631
	ErrCodeYggdrasilPosTooCloseToSameBuild                   = 3632
	ErrCodeYggdrasilWrongUIDForEntity                        = 3633
	ErrCodeYggdrasilCurrentPosNoBuildingToUse                = 3634
	ErrCodeYggdrasilUseCountOutLimit                         = 3635
	ErrCodeYggdrasilWrongGoodsIDForPack                      = 3636
	ErrCodeYggdrasilMessageCountOutOfLimit                   = 3637
	ErrCodeYggdrasilBuildCountOutOfLimit                     = 3638
	ErrCodeYggdrasilPosHasNoMessageToUse                     = 3639
	ErrCodeYggdrasilObjectNotFoundWithParam                  = 3640
	ErrCodeYggdrasilObjectRepeated                           = 3641
	ErrCodeYggdrasilProgressRewardBefore                     = 3642
	ErrCodeYggdrasilBuildNotMatchProtocol                    = 3643
	ErrCodeYggdrasilPosHasNoEntityToUse                      = 3644
	ErrCodeYggdrasilBuildUseOutOfLimit                       = 3645
	ErrCodeYggdrasilPrestigeNotEnough                        = 3646
	ErrCodeYggdrasilNoAreaForPos                             = 3647
	ErrCodeYggdrasilMailNotExist                             = 3648
	ErrCodeYggdrasilMarkNotExist                             = 3649
	ErrCodeYggdrasilMessageRepeated                          = 3650
	ErrCodeYggdrasilCharacterNotExist                        = 3651
	ErrCodeYggdrasilCharacterHpErr                           = 3652
	ErrCodeYggdrasilTransferPortalAlreadyActivated           = 3653
	ErrCodeYggdrasilCoBuildProgressNotEnoughForActivation    = 3654
	ErrCodeYggdrasilTransferPortalNotActivated               = 3655
	ErrCodeYggdrasilPortalTargetNotInList                    = 3656
	ErrCodeYggdrasilTransferPortalThisIDNotTransferPortal    = 3657
	ErrCodeYggdrasilTransferPortalWrongTypeForPortalLocation = 3658
	ErrCodeYggdrasilUnknownMatchEntity                       = 3659
	ErrCodeYggdrasilBuildDestroyedBefore                     = 3660
	ErrCodeYggdrasilMarkCountLimit                           = 3661
	ErrCodeYggdrasilInitPosAreaError                         = 3662
	ErrCodeYggdrasilAreaCoincidence                          = 3663
	ErrCodeYggdrasilTaskCannotAbandon                        = 3664
	ErrCodeYggdrasilBuildCSVNoFindAreaCost                   = 3665
	ErrCodeYggdrasilCannotMoveMonster                        = 3666
	ErrCodeYggdrasilApNoMoreThan                             = 3667
	ErrCodeYggdrasilObjectTypeError                          = 3668

	// manual
	ErrCodeManualNotGet   = 3701
	ErrCodeManualRewarded = 3702

	// Yggdrasil Dispatch
	ErrCodeYggdrasilDispatchTaskNotFound                       = 4001
	ErrCodeYggdrasilDispatchTaskStateNotReadyForMission        = 4002
	ErrCodeYggdrasilDispatchTaskStateNotOnMission              = 4003
	ErrCodeYggdrasilDispatchTaskStateNotReadyForReceiveRewards = 4004
	ErrCodeYggdrasilDispatchNecessaryConditionNotSatisfied     = 4005
	ErrCodeYggdrasilDispatchCharacterIsOnMission               = 4006
	ErrCodeYggdrasilDispatchCharacterLevelNotArrival           = 4007
	ErrCodeYggdrasilDispatchCharacterCampNotArrival            = 4008
	ErrCodeYggdrasilDispatchCharacterCareerNotArrival          = 4009
	ErrCodeYggdrasilDispatchCharacterRarityNotArrival          = 4010
	ErrCodeYggdrasilDispatchCharacterStarNotArrival            = 4011
	ErrCodeYggdrasilDispatchCharacterPowerNotArrival           = 4012
	ErrCodeYggdrasilDispatchCampNotArrival                     = 4013
	ErrCodeYggdrasilDispatchCareerNotArrival                   = 4014
	ErrCodeYggdrasilDispatchSpecificCharacterNotArrival        = 4015
	ErrCodeYggDispatchCSVExtraConditionNotMatchRewards         = 4016
	ErrCodeYggDispatchCSVDispatchTypeNotMatchCloseTime         = 4017
	ErrCodeYggDispatchCSVWrongDispatchType                     = 4018
	ErrCodeYggDispatchCSVWrongLengthOfGuildNum                 = 4019
	ErrCodeYggDispatchWrongTeamSizeForDispatch                 = 4020
	ErrCodeYggDispatchCSVWrongGuildCharacter                   = 4021

	// Greetings
	ErrCodeGreetingsNotFoundInCsv = 4101

	// activity
	ErrCodeActivityCfgNotFound     = 5001
	ErrCodeActivityFuncCfgNotFound = 5002

	//score pass
	ErrCodeScorePassPhaseCfgNotFound  = 6001
	ErrCodeScorePassGroupCfgNotFound  = 6002
	ErrCodeScorePassRewardCfgNotFound = 6003
	ErrCodeScorePassNoSuchSeason      = 6004
	ErrCodeScorePassPhaseNotStart     = 6005

	//battle
	ErrCodeBattleCfgNotFound       = 7001
	ErrCodeBattleNPCCfgNotFound    = 7002
	ErrCodeBattleInvalidNPC        = 7003
	ErrCodeBattleInvalidPosition   = 7004
	ErrCodeBattleDuplicatePosition = 7005
	ErrCodeBattleDuplicateChara    = 7006
	ErrCodeBattleInvalidEndType    = 7007

	// mercenary
	ErrCodeMercenaryWrongSystemType     = 8001
	ErrCodeMercenaryExceedUseLimit      = 8002
	ErrCodeMercenaryNotFound            = 8003
	ErrCodeMercenaryExceedNumLimit      = 8004
	ErrCodeMercenaryApplyExceedLimit    = 8005
	ErrCodeMercenarySendApplyNotFound   = 8006
	ErrCodeMercenaryHandleApplyNotFound = 8007
	ErrCodeMercenaryRepeatedSendApply   = 8008
	ErrCodeMercenaryAlreadyBorrowed     = 8009
	ErrCodeMercenaryAlreadyHad          = 8010

	//formation
	ErrCodeInvalidFormation = 9001
)

func init() {
	errors.SetDefaultCode(ErrCodeDefault)
}

var (
	ErrDefault                = errors.NewCode(ErrCodeDefault, "unknown error")
	ErrUserCacheInvalid       = errors.NewCode(ErrCodeUserCacheInvalid, "User cache invalid. (userID: %d")
	ErrUserNotFound           = errors.NewCode(ErrCodeUserNotFound, "User not found. (userID: %d")
	ErrWhileList              = errors.NewCode(ErrCodeWhileList, "white list intercepted")
	ErrUserLoginInOtherClient = errors.NewCode(ErrCodeUserLoginInOtherClient, "user login in other client")
	ErrLoginSdkError          = errors.NewCode(ErrCodeLoginSdkError, "sdk login error, sdk:%s,code:%d,message:%s")
	ErrUnknownSdk             = errors.NewCode(ErrCodeUnknownSdk, "unknown sdk:%s")
	ErrRegisterLimit          = errors.NewCode(ErrCodeRegisterLimit, "register limit")
	ErrAccountNotFound        = errors.NewCode(ErrCodeAccountNotFound, "account not found")
	ErrAccountRepeat          = errors.NewCode(ErrCodeAccountRepeat, "account repeat")
	ErrSeverMaintaining       = errors.NewCode(ErrCodeSeverMaintaining, "sever maintaining")

	ErrCSVFormatInvalid        = errors.NewCode(ErrCodeDefault, "CSV format invalid. (file: %s, id: %d")
	ErrNotFoundInCSV           = errors.NewCode(ErrCSVNotFound, "Not found in csv. (file: %s, id: %d")
	ErrNotFoundInAssociatedCSV = errors.NewCode(ErrCodeDefault, "Not found in associated csv. (csv: %s, id: %d, associated csv:%s, id: %d")
	ErrProtocolCmdNotFound     = errors.NewCode(ErrCodeDefault, "protocol not found in commands, protocol: %s")

	ErrRecursionReward = errors.NewCode(ErrCodeDefault, "Recursion reward")

	ErrPatchCfgNotFound = errors.NewCode(ErrCodePatchCfgNotFound, "patch config not found")
	ErrAppCfgNotFound   = errors.NewCode(ErrCodeAppCfgNotFound, "app config not found")

	ErrParamError        = errors.NewCode(ErrCodeParamError, "Param error")
	ErrNoPermissionError = errors.NewCode(ErrCodeNoPermissionError, "No Permission")

	ErrItemNotEnough         = errors.NewCode(ErrCodeItemNotEnough, "Item not enough. (id: %d, num: %d")
	ErrItemTypeCannotConsume = errors.NewCode(ErrCodeItemTypeCannotConsume, "Item type cannot consume")
	ErrItemTypeCannotUse     = errors.NewCode(ErrCodeItemTypeCannotUse, "Item type cannot use(id: %d)")

	ErrConnectIllegalCmd          = errors.NewCode(ErrCodeDefault, "Illegal commandId connect from portal")
	ErrConnectTokenInvalid        = errors.NewCode(ErrCodeTokenInvalid, "Token is Invalid")
	ErrConnectUserInOtherGameSvr  = errors.NewCode(ErrCodeDefault, "User in other game server")
	ErrConnectProtocolCmdNotFound = errors.NewCode(ErrCodeDefault, "protocol not found in commands, protocol: %s")

	// user errors
	ErrUserLevelNotArrival        = errors.NewCode(ErrCodeUserLevelNotArrival, "user level not arrival")
	ErrUserNicknameRepeatedToMany = errors.NewCode(ErrCodeUserNicknameRepeatedToMany, "nickname repeated to many")
	ErrUserNicknameNotChange      = errors.NewCode(ErrCodeUserNicknameNotChange, "nickname not change")

	// character errors
	ErrCharacterNotFound           = errors.NewCode(ErrCodeCharacterNotFound, "Character not found. (cid: %d")
	ErrCharacterSkillLevelUnlock   = errors.NewCode(ErrCodeCharacterSkillLevelUnlock, "Character skill level unlock. (cid: %d, skillNum: %d")
	ErrCharacterLevelMax           = errors.NewCode(ErrCodeCharacterLevelMax, "Character level max.")
	ErrCharacterCanNotBeStageUp    = errors.NewCode(ErrCodeDefault, "Character can not stage up.")
	ErrCharacterSkillCannotLevelUp = errors.NewCode(ErrCodeCharacterSkillCannotLevelUp, "Character skill can not level up.")
	ErrCharacterWearCountInvalid   = errors.NewCode(ErrCodeCharacterWearCountInvalid, "Character wear count invalid.")
	ErrCharacterWearRepeatPart     = errors.NewCode(ErrCodeCharacterWearRepeatPart, "Character wear repeat part.")
	ErrCharacterLevelNotArrival    = errors.NewCode(ErrCodeCharacterLevelNotArrival, "Character level not arrival, chara id: %d")
	ErrCharacterStarNotArrival     = errors.NewCode(ErrCodeCharacterStarNotArrival, "Character star not arrival, chara id: %d")

	// equipment errors
	ErrEquipmentNotExist                  = errors.NewCode(ErrCodeEquipmentNotExist, "Equipment not exist. (id: %d")
	ErrEquipmentNotMatchCareer            = errors.NewCode(ErrCodeEquipmentNotMatchCareer, "Equipment not match career. (eid: %d, career: %d")
	ErrEquipmentEXPUpToLimit              = errors.NewCode(ErrCodeEquipmentEXPUpToLimit, "Equipment up to limit. (eid: %d")
	ErrEquipmentAdvanceMaterialNotMatch   = errors.NewCode(ErrCodeEquipmentAdvanceMaterialNotMatch, "Equipment advance material not match. (eid: %d, material: %d")
	ErrEquipmentAdvanceMaterialNoEnough   = errors.NewCode(ErrCodeEquipmentAdvanceMaterialNoEnough, "Equipment advance material not enough. (eid: %d, rarity: %d, stage: %d, material: %v")
	ErrEquipmentLocked                    = errors.NewCode(ErrCodeEquipmentLocked, "Equipment locked. (id: %d")
	ErrEquipmentAdvanceLevelNotEnough     = errors.NewCode(ErrCodeEquipmentAdvanceLevelNotEnough, "Equipment advance level not enough. (eid: %d, level: %d, needLevel: %d")
	ErrEquipmentNoCamp                    = errors.NewCode(ErrCodeEquipmentNoCamp, "Equipment no camp. (eid: %d")
	ErrEquipmentHasNotConfirm             = errors.NewCode(ErrCodeEquipmentHasNotConfirm, "Equipment has not confirm. (eid: %d, camp: %d")
	ErrEquipmentNotRecast                 = errors.NewCode(ErrCodeEquipmentNotRecast, "Equipment not recast. (id: %d")
	ErrEquipmentStrengthenMaterialInvalid = errors.NewCode(ErrCodeEquipmentStrengthenMaterialInvalid, "Equipment strengthen material invalid. (material: %d")
	ErrEquipmentStageUpToLimit            = errors.NewCode(ErrCodeEquipmentStageUpToLimit, "Equipment stage up to limit.")
	ErrEquipmentCantUseSelf               = errors.NewCode(ErrCodeEquipmentCantUseSelf, "Equipment can't use self.")
	ErrEquipmentWearing                   = errors.NewCode(ErrCodeEquipmentWearing, "Equipment wearing.")

	// world item errors
	ErrWorldItemNotExist                = errors.NewCode(ErrCodeWorldItemNotExist, "WorldItem not exist. (id: %d")
	ErrWorldItemNotMatchCareer          = errors.NewCode(ErrCodeWorldItemNotMatchCareer, "WorldItem not match career. (wid: %d, career: %d")
	ErrWorldItemEXPUpToLimit            = errors.NewCode(ErrCodeWorldItemEXPUpToLimit, "WorldItem up to limit. (wid: %d")
	ErrWorldItemAdvanceMaterialNotMatch = errors.NewCode(ErrCodeWorldItemAdvanceMaterialNotMatch, "WorldItem advance material not match. (wid: %d, material: %d")
	ErrWorldItemAdvanceMaterialNoEnough = errors.NewCode(ErrCodeWorldItemAdvanceMaterialNoEnough, "WorldItem advance material not enough. (wid: %d, rarity: %d, stage: %d, material: %v")
	ErrWorldItemLocked                  = errors.NewCode(ErrCodeWorldItemLocked, "WorldItem locked. (id: %d")
	ErrWorldItemAdvanceLevelNotEnough   = errors.NewCode(ErrCodeWorldItemAdvanceLevelNotEnough, "WorldItem advance level not enough. (wid: %d, level: %d, needLevel: %d")
	ErrWorldItemStageUpToLimit          = errors.NewCode(ErrCodeWorldItemStageUpToLimit, "WorldItem stage up to limit.")
	ErrWorldItemCantUseSelf             = errors.NewCode(ErrCodeWorldItemCantUseSelf, "WorldItem can't use self.")
	ErrWorldItemWearing                 = errors.NewCode(ErrCodeWorldItemWearing, "WorldItem wearing.")
	// ErrWorldItemNoCamp                  = errors.NewCode(ErrCodeWorldItemNoCamp, "WorldItem no camp. (eid: %d")
	// ErrWorldItemHasNotConfirm           = errors.NewCode(ErrCodeWorldItemHasNotConfirm, "WorldItem has not confirm. (eid: %d, camp: %d")

	// quest errors
	ErrQuestNotFound                   = errors.NewCode(ErrCodeQuestNotFound, "Quest not found, id: %d")
	ErrQuestConfigNotFound             = errors.NewCode(ErrCodeQuestConfigNotFound, "Quest config not found, id: %d")
	ErrQuestUpdateParamInvalid         = errors.NewCode(ErrCodeQuestUpdateParamInvalid, "Quest update param[%d] invalid, need type: %s")
	ErrQuestNotRegistered              = errors.NewCode(ErrCodeQuestNotRegistered, "Quest not registered in progressing, id: %d")
	ErrQuestProgressNotArrival         = errors.NewCode(ErrCodeQuestProgressNotArrival, "Quest progress not arrival, id: %d")
	ErrQuestNotProgressing             = errors.NewCode(ErrCodeQuestNotProgressing, "Quest not progressing, id: %d")
	ErrQuestActivityConfigNotFound     = errors.NewCode(ErrCodeQuestActivityConfigNotFound, "Quest activity config not found, id: %d")
	ErrQuestActivityRewardGot          = errors.NewCode(ErrCodeQuestActivityRewardGot, "Quest activity reward already got, id: %d")
	ErrQuestActivityProgressNotArrival = errors.NewCode(ErrCodeQuestActivityProgressNotArrival, "Quest activity progress not arrival, id: %d")
	ErrQuestUpdateParamCountInvalid    = errors.NewCode(ErrCodeQuestUpdateParamCountInvalid, "Quest update param count invalid, need count: %d")
	ErrQuestNotComplete                = errors.NewCode(ErrCodeQuestNotComplete, "Quest not complete, id: %d")

	// signin errors
	ErrSignInWrongIDForSignInGroups = errors.NewCode(ErrCodeSignInWrongIDForSignInGroups, "SignIn wrong id for SignInGroups, id: %d")
	ErrSignInWrongIDForSignInData   = errors.NewCode(ErrCodeSignInWrongIDForSignInData, "SignIn wrong id for SignInData, id: %d")
	ErrSignInDropIDMismatchDayCnt   = errors.NewCode(ErrCodeSignInDropIDMismatchDayCnt, "SignIn DropID mismatch DayCnt, id: %d")
	ErrSignInIndexOutOfDropID       = errors.NewCode(ErrCodeSignInIndexOutOfDropID, "SignIn index out of DropID, id: %d")
	ErrSignInRepeat                 = errors.NewCode(ErrCodeSignInRepeat, "SignIn Repeat, id: %d")
	ErrSignInUserHasSigned          = errors.NewCode(ErrCodeSignInUserHasSigned, "SignIn user has signed, id: %d")
	ErrSignInNoDataForToday         = errors.NewCode(ErrCodeSignInNoDataForToday, "SignIn data missed for today, date: %d")

	// survey errors
	ErrSurveyCntOutOfBitNumber           = errors.NewCode(ErrCodeSurveyCntOutOfBitNumber, "Survey count is out of bitnumber, id: %d")
	ErrSurveyDuplicateSurveyData         = errors.NewCode(ErrCodeSurveyDuplicateSurveyData, "Survey duplicate surveydata, id: %d")
	ErrSurveyQuestionCntBeyondLimit      = errors.NewCode(ErrCodeSurveyQuestionCntBeyondLimit, "Survey question count is beyond limit, id: %d")
	ErrSurveyAnswerTextBeyondLimit       = errors.NewCode(ErrCodeSurveyAnswerTextBeyondLimit, "Survey answer text is beyond limit, id: %d")
	ErrSurveyUnlockConditionNotSatisfied = errors.NewCode(ErrCodeSurveyUnlockConditionNotSatisfied, "Survey unlock conditions are not satisfied, id: %d")
	ErrSurveyValidityPeriodBeyondLimit   = errors.NewCode(ErrCodeSurveyValidityPeriodBeyondLimit, "Survey validity peroid is beyond limit, id: %d")
	ErrSurveyTypeBeyondLimit             = errors.NewCode(ErrCodeSurveyTypeBeyondLimit, "Survey type is beyond limit, id: %d")
	ErrSurveyWrongIDForSurveyInfos       = errors.NewCode(ErrCodeSurveyWrongIDForSurveyInfos, "Survey id is wrong for SurveyInfos, id: %d")

	// coins errors
	// ErrGoldNotEnough = errors.NewCode(ErrCodeDefault, "Gold not enough. (have: %d, cost: %d")

	// graveyard
	ErrGraveyardTypeCannotBuild               = errors.NewCode(ErrCodeGraveyardTypeCannotBuild, "type can not build")
	ErrGraveyardNoAreaCannotBuild             = errors.NewCode(ErrCodeGraveyardNoAreaCannotBuild, "can not build for no area")
	ErrGraveyardNumLimitCannotBuild           = errors.NewCode(ErrCodeGraveyardNumLimitCannotBuild, "can not build for num limit")
	ErrGraveyardBuildNotExist                 = errors.NewCode(ErrCodeGraveyardBuildNotExist, "build not exist")
	ErrGraveyardBuildNotInTransaction         = errors.NewCode(ErrCodeGraveyardBuildNotInTransaction, "build not in transaction")
	ErrGraveyardBuildTransactionTimeNotEnough = errors.NewCode(ErrCodeGraveyardBuildTransactionTimeNotEnough, "transaction time not enough")
	ErrGraveyardBuildInProduct                = errors.NewCode(ErrCodeGraveyardBuildInProduct, "build in product")
	ErrGraveyardBuildInTransaction            = errors.NewCode(ErrCodeGraveyardBuildInTransaction, "build  in transaction")
	ErrGraveyardMainTowerLvLimit              = errors.NewCode(ErrCodeGraveyardMainTowerLvLimit, "main tower lv limit")
	ErrGraveyardBuildCountLimit               = errors.NewCode(ErrCodeGraveyardBuildCountLimit, "build count limit")
	ErrGraveyardProduceTimeLimit              = errors.NewCode(ErrCodeGraveyardProduceTimeLimit, "produce time limit")
	ErrGraveyardBuildCanNotLvUp               = errors.NewCode(ErrCodeGraveyardBuildCanNotLvUp, "build can not lv up (buildId:%d ,lv:%d)")
	ErrGraveyardBuildMaxLv                    = errors.NewCode(ErrCodeGraveyardBuildMaxLv, "build max Lv (buildId:%d ,lv:%d)")
	ErrGraveyardBuildCanNotStageUp            = errors.NewCode(ErrCodeGraveyardBuildCanNotStageUp, "build can not stage up (buildId:%d ,stage:%d)")
	ErrGraveyardBuildMaxStage                 = errors.NewCode(ErrCodeGraveyardBuildMaxStage, "build max stage (buildId:%d ,stage:%d)")
	ErrGraveyardPopulationNumOverflow         = errors.NewCode(ErrCodeGraveyardPopulationNumOverflow, "build population Overflow (num:%d ,maxNum:%d)")
	ErrGraveyardCharacterNumOverflow          = errors.NewCode(ErrCodeGraveyardCharacterNumOverflow, "build character Overflow (num:%d ,maxNum:%d)")
	ErrGraveyardTypeCannotProduce             = errors.NewCode(ErrCodeGraveyardTypeCannotProduce, "type can not produce")
	ErrGraveyardProduceTimeNotEnough          = errors.NewCode(ErrCodeGraveyardProduceTimeNotEnough, "produce time not enough ")
	ErrGraveyardProduceNumLimit               = errors.NewCode(ErrCodeGraveyardProduceNumLimit, " produce num limit")
	ErrGraveyardBuildCannotReduceProductTime  = errors.NewCode(ErrCodeGraveyardBuildCannotReduceProductTime, "build cannot reduce product time")
	ErrGraveyardBuffInUse                     = errors.NewCode(ErrCodeGraveyardBuffInUse, "buff in use")
	ErrGraveyardCannotUseBuffBuildsNotExist   = errors.NewCode(ErrCodeGraveyardCannotUseBuffBuildsNotExist, "cannot use buff builds not exist")
	ErrGraveyardNoRemainedHelpCount           = errors.NewCode(ErrCodeGraveyardNoRemainedHelpCount, "no remained help count")
	ErrGraveyardBuildInHelpNow                = errors.NewCode(ErrCodeGraveyardBuildInHelpNow, "build in help now")
	ErrGraveyardCannotPlotRewardNow           = errors.NewCode(ErrCodeGraveyardCannotPlotRewardNow, "cannot plot reward now")
	ErrGraveyardPlotRewardNumMax              = errors.NewCode(ErrCodeGraveyardPlotRewardNumMax, "plot reward num max")

	// level
	ErrLevelConfigNotFound       = errors.NewCode(ErrCodeLevelConfigNotFound, "level config not found, id: %d")
	ErrLevelNotStarted           = errors.NewCode(ErrCodeLevelNotStarted, "level not started, id: %d")
	ErrOtherLevelStarted         = errors.NewCode(ErrCodeOtherLevelStarted, "ohter level started, id: %d")
	ErrLevelNotPassed            = errors.NewCode(ErrCodeLevelNotPassed, "level not passed, id: %d")
	ErrExploreMapObjectNotLevel  = errors.NewCode(ErrCodeExploreMapObjectNotLevel, "explore map object not has a level, id: %d")
	ErrLevelPassTimesLimited     = errors.NewCode(ErrCodeLevelPassTimesLimited, "level pass times limited, id: %d")
	ErrLevelTargetNotCompleteAll = errors.NewCode(ErrCodeLevelTargetNotCompleteAll, "level doesn't complete all target, id: %d")
	ErrLevelInvalidSystemParam   = errors.NewCode(ErrCodeLevelInvalidSystemParam, "level invalid system param, id: %d")
	ErrLevelInvalidTarget        = errors.NewCode(ErrCodeLevelInvalidTarget, "level invalid target, id: %d")
	ErrLevelInvalidAchievement   = errors.NewCode(ErrCodeLevelInvalidAchievement, "level invalid system achivement, id: %d")

	// explore
	ErrChapterConfigNotFound            = errors.NewCode(ErrCodeChapterConfigNotFound, "chapter config not found, id: %d")
	ErrNotInChapter                     = errors.NewCode(ErrCodeNotInChapter, "not in chapter, id: %d")
	ErrChapterRewardConfigNotFound      = errors.NewCode(ErrCodeChapterRewardConfigNotFound, "chapter reward config not found, id: %d")
	ErrChapterScoreNotArrival           = errors.NewCode(ErrCodeChapterScoreNotArrival, "chapter score not arrival, id: %d")
	ErrChapterRewardReceived            = errors.NewCode(ErrCodeChapterRewardReceived, "chapter reward received, id: %d, reward id: %d")
	ErrNotCampainLevel                  = errors.NewCode(ErrCodeNotCampainLevel, "level not a campain level, id: %d")
	ErrNotExploreEliteLevel             = errors.NewCode(ErrCodeNotExploreEliteLevel, "level not a elite level, id: %d")
	ErrMapObjectConfigNotFound          = errors.NewCode(ErrCodeMapObjectConfigNotFound, "map object config not found, id: %d")
	ErrExploreEventPointConfigNotFound  = errors.NewCode(ErrCodeExploreEventPointConfigNotFound, "explore event point config not found, id: %d")
	ErrExploreInteracted                = errors.NewCode(ErrCodeExploreInteracted, "explore obj already interacted, id: %d")
	ErrExploreNotInteracted             = errors.NewCode(ErrCodeExploreNotInteracted, "explore obj not interacted, id: %d")
	ErrExploreNPCConfigNotFound         = errors.NewCode(ErrCodeExploreNPCConfigNotFound, "explore npc config not found, id: %d")
	ErrExploreNPCOptionInvalid          = errors.NewCode(ErrCodeExploreNPCOptionInvalid, "explore npc option invalid, id: %d, option: %d")
	ErrExploreRewardPointConfigNotFound = errors.NewCode(ErrCodeExploreRewardPointConfigNotFound, "explore gather config not found, id: %d")
	ErrExploreMonsterConfigNotFound     = errors.NewCode(ErrCodeExploreMonsterConfigNotFound, "explore monster config not found, id: %d")
	ErrExploreMonsterHasNotLevel        = errors.NewCode(ErrCodeExploreMonsterHasNotLevel, "explore monster don't has level, id: %d, level id: %d")
	ErrExploreFogConfigNotFound         = errors.NewCode(ErrCodeExploreFogConfigNotFound, "explore fog config not found, id: %d")
	ErrExploreFogLocked                 = errors.NewCode(ErrCodeExploreFogLocked, "explore fog is locked, id: %d")
	ErrExploreResourceConfigNotFound    = errors.NewCode(ErrCodeExploreResourceConfigNotFound, "explore resource config not found, id: %d")
	ErrExploreResourceIsCollecting      = errors.NewCode(ErrCodeExploreResourceIsCollecting, "explore resource is collecting, id: %d")
	ErrExploreResourceCollectTimesLimit = errors.NewCode(ErrCodeExploreResourceCollectTimesLimit, "explore resource collect time limit, id: %d")
	ErrExploreResourceNotCollecting     = errors.NewCode(ErrCodeExploreResourceNotCollecting, "explore resource is not collecting, id: %d")
	ErrExploreResourceNotCollected      = errors.NewCode(ErrCodeExploreResourceNotCollected, "explore resource is not collected, id: %d")
	ErrExploreResourceHasNotLevel       = errors.NewCode(ErrCodeExploreResourceHasNotLevel, "explore resource don't has level, id: %d, level id: %d")
	ErrExploreResourceLevelPassed       = errors.NewCode(ErrCodeExploreResourceLevelPassed, "explore resource level has passed, id: %d, level id: %d")
	ErrExploreTPGateConfigNotFound      = errors.NewCode(ErrCodeExploreTPGateConfigNotFound, "explore transport gate config not found, id: %d")
	ErrExploreTPGateUseTimesLimited     = errors.NewCode(ErrCodeExploreTPGateUseTimesLimited, "explore transport gate use times limited, id: %d")

	// hero
	ErrHeroConfigNotFound          = errors.NewCode(ErrCodeHeroConfigNotFound, "hero config not found, id: %d")
	ErrHeroLevelConfigNotFound     = errors.NewCode(ErrCodeHeroLevelConfigNotFound, "hero level config not found, level: %d")
	ErrHeroExpNotEnough            = errors.NewCode(ErrCodeHeroExpNotEnough, "hero exp not enough")
	ErrHeroLevelMax                = errors.NewCode(ErrCodeHeroLevelMax, "hero level max")
	ErrHeroUnlocked                = errors.NewCode(ErrCodeHeroUnlocked, "hero already unlocked, id: %d")
	ErrHeroNotFound                = errors.NewCode(ErrCodeHeroNotFound, "hero not found, id: %d")
	ErrHeroHasNotSkill             = errors.NewCode(ErrCodeHeroHasNotSkill, "hero has not skill, id: %d, skill id： %d")
	ErrHeroSkillConfigNotFound     = errors.NewCode(ErrCodeHeroSkillConfigNotFound, "hero skill config not found, id: %d, skill id： %d")
	ErrHeroSkillLevelNotFound      = errors.NewCode(ErrCodeHeroSkillLevelNotFound, "hero skill level not found, id: %d, skill id： %d, level: %d")
	ErrHeroAttendantConfigNotFound = errors.NewCode(ErrCodeHeroAttendantConfigNotFound, "hero attendant config not found, id: %d, slot: %d")
	ErrHeroLevelNotArrival         = errors.NewCode(ErrCodeHeroLevelNotArrival, "hero level not arrival")
	ErrHeroSkillLevelNotArrival    = errors.NewCode(ErrCodeHeroSkillLevelNotArrival, "hero skill level not arrival, hero id: %d, skill id: %d")
	ErrHeroSkillItemUsedNotArrival = errors.NewCode(ErrCodeHeroSkillItemUsedNotArrival, "hero skill item use not arrival, hero id: %d")

	// manual
	ErrManualNotGet   = errors.NewCode(ErrCodeManualNotGet, "manual not get")
	ErrManualRewarded = errors.NewCode(ErrCodeManualRewarded, "manual rewarded")

	// store
	ErrStoreWrongSubStoreTypeForData            = errors.NewCode(ErrCodeStoreWrongSubStoreTypeForData, "store wrong subStore type for storeData, id: %d")
	ErrStoreUpdateTimesNotMatchCntOfUpdateCost  = errors.NewCode(ErrCodeStoreUpdateTimesNotMatchCntOfUpdateCost, "store update times does not match cnt of update cost for storeData, id: %d")
	ErrStoreWrongStoreIDForData                 = errors.NewCode(ErrCodeStoreWrongStoreIDForData, "store wrong storeID for storeData, id: %d")
	ErrStoreWrongSubStoreIDForData              = errors.NewCode(ErrCodeStoreWrongSubStoreIDForData, "store wrong subStoreID for storeData, id: %d")
	ErrStoreWrongGoodsIDForData                 = errors.NewCode(ErrCodeStoreWrongGoodsIDForData, "store wrong goodsID for storeData, id: %d")
	ErrStoreWrongUpdateRuleIDForData            = errors.NewCode(ErrCodeStoreWrongUpdateRuleIDForData, "store wrong updateRuleID for storeData, id: %d")
	ErrStoreWrongSubStoreIDInStoreForInfo       = errors.NewCode(ErrCodeStoreWrongSubStoreIDInStoreForInfo, "store wrong subStoreID in store for storeInfo, id: %d")
	ErrStoreIndexOutOfRangeForInfo              = errors.NewCode(ErrCodeStoreIndexOutOfRangeForInfo, "store cell index out of range for storeInfo, subStoreID: %d, index: %d")
	ErrStoreWrongGoodsIDInCellForInfo           = errors.NewCode(ErrCodeStoreWrongGoodsIDInCellForInfo, "store wrong goodsID in cell for storeInfo, subStoreID: %d, goodsID: %d")
	ErrStoreNumOfRemainsNotEnough               = errors.NewCode(ErrCodeStoreNumOfRemainsNotEnough, "store num of remains of goods not enough, subStoreID: %d, goodsID: %d")
	ErrStoreSoldOutConditionNotSatisfied        = errors.NewCode(ErrCodeStoreSoldOutConditionNotSatisfied, "store sold out condition not satisfied, id: %d")
	ErrStoreNoOptionForSelectSubStore           = errors.NewCode(ErrCodeStoreNoOptionForSelectSubStore, "store no option for SelectSubStore, id: %d")
	ErrStoreNoOptionForSelectSubStores          = errors.NewCode(ErrCodeStoreNoOptionForSelectSubStores, "store no option for SelectSubStores, id: %d")
	ErrStoreUnlockConditionNotSatisfied         = errors.NewCode(ErrCodeStoreUnlockConditionNotSatisfied, "store unlockCondition does not satisfied, id: %d")
	ErrStoreSubStoreUnlockConditionNotSatisfied = errors.NewCode(ErrCodeStoreSubStoreUnlockConditionNotSatisfied, "store subStore unlockCondition does not satisfied, id: %d")
	ErrStoreCellUnlockConditionNotSatisfied     = errors.NewCode(ErrCodeStoreCellUnlockConditionNotSatisfied, "store cell unlockCondition does not satisfied, id: %d")
	ErrStoreGoodsUnlockConditionNotSatisfied    = errors.NewCode(ErrCodeStoreGoodsUnlockConditionNotSatisfied, "store goods unlockCondition does not satisfied, id: %d")
	ErrStoreWrongStoreIDForStoreInfo            = errors.NewCode(ErrCodeStoreWrongStoreIDForStoreInfo, "store wrong storeID for storeInfo, id: %d")
	ErrStoreNotSupportUpdate                    = errors.NewCode(ErrCodeStoreNotSupportUpdate, "store does not support update, id: %d")
	ErrStoreExceedMaxUpdateTimes                = errors.NewCode(ErrCodeStoreExceedMaxUpdateTimes, "store exceed max update times, id: %d")
	ErrStoreCurrencyNotMatchDuringPurchase      = errors.NewCode(ErrCodeStoreCurrencyNotMatchDuringPurchase, "store goods currency does not match during purchase, id: %d")
	ErrStoreCurrencyNotMatchPriceForData        = errors.NewCode(ErrCodeStoreCurrencyNotMatchPriceForData, "store goods currency does not match price for storeData, id: %d")
	ErrStoreNoGoodsInCellToGenerate             = errors.NewCode(ErrCodeStoreNoGoodsInCellToGenerate, "store no goods in cell to generate")
	ErrStoreWrongSubStoreIDForInfo              = errors.NewCode(ErrCodeStoreWrongSubStoreIDForInfo, "store wrong subStoreID for info, id: %d")
	ErrStoreWrongGoodsIDOfQuickPurchaseForData  = errors.NewCode(ErrCodeStoreWrongGoodsIDOfQuickPurchaseForData, "store wrong goodsID of quick-purchase for data, id: %d")

	// quick purchase stamina
	ErrQuickPurchaseStaminaIndexOutOfRangeForData = errors.NewCode(ErrCodeQuickPurchaseStaminaIndexOutOfRangeForData, "stamina index out of range for data, id: %d")
	ErrQuickPurchaseStaminaExceedMaxTimes         = errors.NewCode(ErrCodeQuickPurchaseStaminaExceedMaxTimes, "stamina record exceed maxTimes")

	// vn
	ErrVnRewarded      = errors.NewCode(ErrCodeVnRewarded, "vn rewarded")
	ErrVnRewardIsEmpty = errors.NewCode(ErrCodeVnRewardIsEmpty, "vn reward is empty")

	// guide
	ErrGuidePassed  = errors.NewCode(ErrCodeGuidePassed, "guide passed")
	ErrGuideNotPass = errors.NewCode(ErrCodeGuideNotPass, "guide not pass,id: %d")

	// portal
	ErrTooManyRequest = errors.NewCode(ErrCodeTooManyRequest, "Portal too many request.")
	ErrTooManyConn    = errors.NewCode(ErrCodeTooManyConn, "Portal too many conn.")

	// tower
	ErrTowerConfigNotFound      = errors.NewCode(ErrCodeTowerConfigNotFound, "tower config not found, id: %d")
	ErrTowerStageConfigNotFound = errors.NewCode(ErrCodeTowerStageConfigNotFound, "tower stage config not found, id: %d, stage: %d")
	ErrTowerStageNotArrival     = errors.NewCode(ErrCodeTowerStageNotArrival, "tower stage not arrival, id: %d, stage: %d")
	ErrTowerNotActived          = errors.NewCode(ErrCodeTowerNotActived, "tower not actived, id: %d")
	ErrTowerInvalidLevel        = errors.NewCode(ErrCodeTowerInvalidLevel, "tower invalid level, tower id: %d, level id: %d")
	ErrTowerGoUpLimited         = errors.NewCode(ErrCodeTowerGoUpLimited, "tower go up times limited, tower id: %d")
	ErrTowerCharaCampLimited    = errors.NewCode(ErrCodeTowerCharaCampLimited, "tower chara camp limited, tower id: %d, chara id: %d")

	// mail
	ErrMailNotFound               = errors.NewCode(ErrCodeMailNotFound, "mail not found, id: %d")
	ErrMailHasNotAttachment       = errors.NewCode(ErrCodeMailHasNotAttachment, "mail has not attachment, id: %d")
	ErrMailIsNotRemovable         = errors.NewCode(ErrCodeMailIsNotRemovable, "mail is not removable, id: %d")
	ErrMailTemplateConfigNotFound = errors.NewCode(ErrCodeMailTemplateConfigNotFound, "mail template not found, id: %d")

	// guild
	ErrGuildHasJoined                 = errors.NewCode(ErrCodeGuildHasJoined, "Guild has joined. (guildID: %d")
	ErrGuildHasNotJoin                = errors.NewCode(ErrCodeGuildHasNotJoin, "Guild has not join.")
	ErrGuildQuitCD                    = errors.NewCode(ErrCodeGuildQuitCD, "Guild quit cd. (lastQuitTime: %d")
	ErrGuildNoPrivilege               = errors.NewCode(ErrCodeGuildNoPrivilege, "Guild not privilege.")
	ErrGuildNotMember                 = errors.NewCode(ErrCodeGuildNotMember, "Guild not member. (guildID: %d, userID: %d")
	ErrGuildChairmanCantQuit          = errors.NewCode(ErrCodeGuildChairmanCantQuit, "Guild chairman can't quit.")
	ErrGuildIsFull                    = errors.NewCode(ErrCodeGuildIsFull, "Guild is full.")
	ErrGuildInDissolving              = errors.NewCode(ErrCodeGuildInDissolving, "Guild in dissolving.")
	ErrGuildNotFound                  = errors.NewCode(ErrCodeGuildNotFound, "Guild not found.")
	ErrGuildNotInAppliedList          = errors.NewCode(ErrCodeGuildNotInAppliedList, "Guild not in applied list.")
	ErrGuildYggWrongIDForCoBuild      = errors.NewCode(ErrCodeGuildYggWrongIDForCoBuild, "Guild yggdrasil wrong ID for coBuild, buildID: %d")
	ErrGuildDissolved                 = errors.NewCode(ErrCodeGuildDissolved, "Guild dissolved.")
	ErrGuildMemberEquipmentNotExist   = errors.NewCode(ErrCodeGuildMemberEquipmentNotExist, "Guild equipment of member not exist, userId: %d, equipmentId: %d")
	ErrGuildMemberWorldItemNotExist   = errors.NewCode(ErrCodeGuildMemberWorldItemNotExist, "Guild WorldItem of member not exist, userId: %d, equipmentId: %d")
	ErrGuildExceedApplyNumLimit       = errors.NewCode(ErrCodeGuildExceedApplyNumLimit, "Guild number of apply exceed limit.")
	ErrGuildKickedMemberNotCommon     = errors.NewCode(ErrCodeGuildKickedMemberNotCommon, "Kicked member not common.")
	ErrGuildTaskRewardsHasReceived    = errors.NewCode(ErrCodeGuildTaskRewardsHasReceived, "Guild task rewards have been received, taskId: %d")
	ErrGuildTaskConfigNotFound        = errors.NewCode(ErrCodeGuildTaskConfigNotFound, "Guild task config not found, taskId: %d")
	ErrGuildTaskNotFinish             = errors.NewCode(ErrCodeGuildTaskNotFinish, "Guild task not finish, taskId: %d, current progress: %d, count: %d")
	ErrGuildSendGroupMailInCD         = errors.NewCode(ErrCodeGuildSendGroupMailInCD, "Guild send group mail in cd.")
	ErrGuildModifyIconIsNull          = errors.NewCode(ErrCodeGuildModifyIconIsNull, "Guild icon is nil. GuildId: %d")
	ErrGuildRepeatName                = errors.NewCode(ErrCodeGuildRepeatName, "Guild name repeat.")
	ErrGuildViceCharimanNumOutOfLimit = errors.NewCode(ErrCodeGuildViceCharimanNumOutOfLimit, "Guild num of viceChairman out of limit, guildId: %d")
	ErrGuildEliteNumOutOfLimit        = errors.NewCode(ErrCodeGuildViceCharimanNumOutOfLimit, "Guild num of elite out of limit, guildId: %d")
	ErrGuildInDissolvedCD             = errors.NewCode(ErrCodeGuildInDissolvedCD, "Guild in dissolved CD, GuildId: %d")

	ErrCompileLuaFileFailed           = errors.NewCode(ErrCodeCompileLuaFileFailed, "Combat power calculate compile lua file failed.")
	ErrPCallFunctionProtoFailed       = errors.NewCode(ErrCodePCallFunctionProtoFailed, "Combat power calculate pcall functionproto failed.")
	ErrCallGlobalFunctionFailed       = errors.NewCode(ErrCodeCallGlobalFunctionFailed, "Combat power calcullate call lua function failed.")
	ErrRetNotLNumber                  = errors.NewCode(ErrCodeRetNotLNumber, "Combat power calculate ret is not LNumber.")
	ErrPowerAdaptWrongNumberForSymbol = errors.NewCode(ErrCodePowerAdaptWrongNumberForSymbol, "Power Adapt wrong number for symbol in data, symbol: %d")
	ErrPowerAdaptWrongCareerIDForParm = errors.NewCode(ErrCodePowerAdaptWrongCareerIDForParm, "Power Adapt wrong careerID for parm, careerID: %d")
	ErrPowerAdaptWrongSymbolForParm   = errors.NewCode(ErrCodePowerAdaptWrongSymbolForParm, "Power Adapt wrong symbol for parm, symbol: %d")

	// action unlock
	ErrActionUnlockNotsatisfied  = errors.NewCode(ErrCodeActionUnlockNotsatisfied, "Action unlock conditions not satisfied, ActionID: %d")
	ErrActionUnlockWrongActionID = errors.NewCode(ErrCodeActionUnlockWrongActionID, "Action unlock wrong ActionID for csv, ActionID: %d")
	//yggdrasil

	ErrYggdrasilInTravel                                 = errors.NewCode(ErrCodeYggdrasilInTravel, "in travel now ")
	ErrYggdrasilNotInTravel                              = errors.NewCode(ErrCodeYggdrasilNotInTravel, "not in travel now ")
	ErrYggdrasilCharacterCannotCarry                     = errors.NewCode(ErrCodeYggdrasilCharacterCannotCarry, "character cannot carry")
	ErrYggdrasilCityCannotTravel                         = errors.NewCode(ErrCodeYggdrasilCityCannotTravel, "city cannot travel")
	ErrYggdrasilCannotMoveDistanceIllegal                = errors.NewCode(ErrCodeYggdrasilCannotMoveDistanceIllegal, "cannot move distance illegal")
	ErrYggdrasilCannotMovePosUnWalkable                  = errors.NewCode(ErrCodeYggdrasilCannotMovePosUnWalkable, "cannot move pos unWalkable")
	ErrYggdrasilCannotMoveTerrainDiff                    = errors.NewCode(ErrCodeYggdrasilCannotMoveTerrainDiff, "cannot move terrain diff")
	ErrYggdrasilApNotEnough                              = errors.NewCode(ErrCodeYggdrasilApNotEnough, "ap not enough")
	ErrYggdrasilNoTravelTime                             = errors.NewCode(ErrCodeYggdrasilNoTravelTime, "no travel time")
	ErrYggdrasilCannotReturnCityThisPos                  = errors.NewCode(ErrCodeYggdrasilCannotReturnCityThisPos, "cannot return city this pos")
	ErrYggdrasilThisCityCannotTransfer                   = errors.NewCode(ErrCodeYggdrasilThisCityCannotTransfer, "This City Cannot Transfer to")
	ErrYggdrasilInCityNow                                = errors.NewCode(ErrCodeYggdrasilInCityNow, "in city now")
	ErrYggdrasilNotInCityNow                             = errors.NewCode(ErrCodeYggdrasilNotInCityNow, "not in city now")
	ErrYggdrasilPackGoodsNotFound                        = errors.NewCode(ErrCodeYggdrasilPackGoodsNotFound, "pack goods not found")
	ErrYggdrasilOverlapping                              = errors.NewCode(ErrCodeYggdrasilOverlapping, "Overlapping pos")
	ErrYggdrasilDiscardGoodsNotFound                     = errors.NewCode(ErrCodeYggdrasilDiscardGoodsNotFound, "discard goods not found")
	ErrYggdrasilBagIsFull                                = errors.NewCode(ErrCodeYggdrasilBagIsFull, "bag is full")
	ErrYggdrasilObjectNotFound                           = errors.NewCode(ErrCodeYggdrasilObjectNotFound, "object not found")
	ErrYggdrasilTaskProcessError                         = errors.NewCode(ErrCodeYggdrasilTaskProcessError, "task process error")
	ErrYggdrasilCompleteTaskBefore                       = errors.NewCode(ErrCodeYggdrasilCompleteTaskBefore, "complete task before")
	ErrYggdrasilAcceptTaskBefore                         = errors.NewCode(ErrCodeYggdrasilAcceptTaskBefore, "accept task before")
	ErrYggdrasilSameTaskGroup                            = errors.NewCode(ErrCodeYggdrasilSameTaskGroup, "same task group")
	ErrYggdrasilNotDoingTaskNow                          = errors.NewCode(ErrCodeYggdrasilNotDoingTaskNow, "not doing task now")
	ErrYggdrasilOtherMaxTaskCount                        = errors.NewCode(ErrCodeYggdrasilOtherMaxTaskCount, "other max task count")
	ErrYggdrasilTaskCannotComplete                       = errors.NewCode(ErrCodeYggdrasilTaskCannotComplete, "task cannot complete")
	ErrYggdrasilTaskNotInProgress                        = errors.NewCode(ErrCodeYggdrasilTaskNotInProgress, "task not in progress")
	ErrYggdrasilPosAlreadyHasEntity                      = errors.NewCode(ErrCodeYggdrasilPosAlreadyHasEntity, "Yggdrasil this position already has entity, coordinate: (%d, %d)")
	ErrYggdrasilPosAlreadyHasObject                      = errors.NewCode(ErrCodeYggdrasilPosAlreadyHasObject, "Yggdrasil this position already has object, coordinate: (%d, %d)")
	ErrYggdrasilPosTypeForbidBuild                       = errors.NewCode(ErrCodeYggdrasilPosTypeForbidBuild, "Yggdrasil type of this position is forbid for building, coordinate: (%d, %d)")
	ErrYggdrasilPosInsideCityBanRadius                   = errors.NewCode(ErrCodeYggdrasilPosInsideCityBanRadius, "Yggdrasil this position is inside CityBanRadius, coordinate: (%d, %d)")
	ErrYggdrasilPosTooCloseToSameBuild                   = errors.NewCode(ErrCodeYggdrasilPosTooCloseToSameBuild, "Yggdrasil this position too close to same build, coordinate: (%d, %d)")
	ErrYggdrasilWrongUIDForEntity                        = errors.NewCode(ErrCodeYggdrasilWrongUIDForEntity, "Yggdrasil wrong uid for %s, uid: %d")
	ErrYggdrasilCurrentPosNoBuildingToUse                = errors.NewCode(ErrCodeYggdrasilCurrentPosNoBuildingToUse, "Yggdrasil current position has no building to use, coordinate: (%d, %d)")
	ErrYggdrasilUseCountOutLimit                         = errors.NewCode(ErrCodeYggdrasilUseCountOutLimit, "Yggdrasil usecount for this building is out limit, BuildUseCount: %d")
	ErrYggdrasilWrongGoodsIDForPack                      = errors.NewCode(ErrCodeYggdrasilWrongGoodsIDForPack, "Yggdrasil wrong goodsID for pack, GoodsID: %d")
	ErrYggdrasilMessageCountOutOfLimit                   = errors.NewCode(ErrCodeYggdrasilMessageCountOutOfLimit, "Yggdrasil message count out of limit")
	ErrYggdrasilBuildCountOutOfLimit                     = errors.NewCode(ErrCodeYggdrasilBuildCountOutOfLimit, "Yggdrasil build count out of limit, BuildType: %d")
	ErrYggdrasilPosHasNoMessageToUse                     = errors.NewCode(ErrCodeYggdrasilPosHasNoMessageToUse, "Yggdrasil pos has no message to use, coordinate: (%d, %d)")
	ErrYggdrasilObjectNotFoundWithParam                  = errors.NewCode(ErrCodeYggdrasilObjectNotFoundWithParam, "object not found,objectId: %d")
	ErrYggdrasilObjectRepeated                           = errors.NewCode(ErrCodeYggdrasilObjectRepeated, "object repeated,objectId: %d")
	ErrYggdrasilProgressRewardBefore                     = errors.NewCode(ErrCodeYggdrasilProgressRewardBefore, "progress reward before")
	ErrYggdrasilBuildNoMatchProtocol                     = errors.NewCode(ErrCodeYggdrasilBuildNotMatchProtocol, "Yggdrasil build in this pos does not match protocol")
	ErrYggdrasilPosHasNoEntityToUse                      = errors.NewCode(ErrCodeYggdrasilPosHasNoEntityToUse, "Yggdrasil current position has no %s to use, coordinate:(%d, %d)")
	ErrYggdrasilBuildUseOutOfLimit                       = errors.NewCode(ErrCodeYggdrasilBuildUseOutOfLimit, "Yggdrasil build use out of usingParm limit, buildID: %d")
	ErrYggdrasilPrestigeNotEnough                        = errors.NewCode(ErrCodeYggdrasilPrestigeNotEnough, "prestige not enough")
	ErrYggdrasilNoAreaForPos                             = errors.NewCode(ErrCodeYggdrasilNoAreaForPos, "no area for  pos:%+v ")
	ErrYggdrasilMailNotExist                             = errors.NewCode(ErrCodeYggdrasilMailNotExist, "mail not exist")
	ErrYggdrasilMarkNotExist                             = errors.NewCode(ErrCodeYggdrasilMarkNotExist, "mark not exist")
	ErrYggdrasilMessageRepeated                          = errors.NewCode(ErrCodeYggdrasilMessageRepeated, "message repeated coordinate: (%d, %d)")
	ErrYggdrasilCharacterNotExist                        = errors.NewCode(ErrCodeYggdrasilCharacterNotExist, "character not exist,characterId:%d")
	ErrYggdrasilCharacterHpErr                           = errors.NewCode(ErrCodeYggdrasilCharacterHpErr, "character hp err,characterId:%d,hp:%d")
	ErrYggdrasilTransferPortalAlreadyActivated           = errors.NewCode(ErrCodeYggdrasilTransferPortalAlreadyActivated, "Yggdrasil transfer portal already activated, buildID: %d")
	ErrYggdrasilCoBuildProgressNotEnoughForActivation    = errors.NewCode(ErrCodeYggdrasilCoBuildProgressNotEnoughForActivation, "Yggdrasil coBuild progress not enough for activation, buildID: %d")
	ErrYggdrasilTransferPortalNotActivated               = errors.NewCode(ErrCodeYggdrasilTransferPortalNotActivated, "Yggdrasil transfer portal not activated, buildID: %d")
	ErrYggdrasilPortalTargetNotInList                    = errors.NewCode(ErrCodeYggdrasilPortalTargetNotInList, "Yggdrasil portal target not in list, locationType: %d, locationID: %d")
	ErrYggdrasilTransferPortalThisIDNotTransferPortal    = errors.NewCode(ErrCodeYggdrasilTransferPortalThisIDNotTransferPortal, "Yggdrasil this ID is wrong for transfer portal, buildID: %d")
	ErrYggdrasilTransferPortalWrongTypeForPortalLocation = errors.NewCode(ErrCodeYggdrasilTransferPortalWrongTypeForPortalLocation, "Yggdrasil wrong type for portal location, type: %d")
	ErrYggdrasilUnknownMatchEntity                       = errors.NewCode(ErrCodeYggdrasilUnknownMatchEntity, "unknown match entity:%+v")
	ErrYggdrasilBuildDestroyedBefore                     = errors.NewCode(ErrCodeYggdrasilBuildDestroyedBefore, "build destroyed before uid:%d")
	ErrYggdrasilMarkCountLimit                           = errors.NewCode(ErrCodeYggdrasilMarkCountLimit, "mark count limit")
	ErrYggdrasilInitPosAreaError                         = errors.NewCode(ErrCodeYggdrasilInitPosAreaError, "init pos area error")
	ErrYggdrasilAreaCoincidence                          = errors.NewCode(ErrCodeYggdrasilAreaCoincidence, "area coincidence")
	ErrYggdrasilTaskCannotAbandon                        = errors.NewCode(ErrCodeYggdrasilTaskCannotAbandon, "task cannot abandon")
	ErrYggdrasilBuildCSVNoFindAreaCost                   = errors.NewCode(ErrCodeYggdrasilBuildCSVNoFindAreaCost, "Yggdrasil build can not find build cost for areaId, buildId:  %d, areaId: %d")
	ErrYggdrasilCannotMoveMonster                        = errors.NewCode(ErrCodeYggdrasilCannotMoveMonster, "cannot move monster")
	ErrYggdrasilApNoMoreThan                             = errors.NewCode(ErrCodeYggdrasilApNoMoreThan, "ap no more than,limit:%d,now:%d")
	ErrYggdrasilObjectTypeError                          = errors.NewCode(ErrCodeYggdrasilObjectTypeError, "ygg object type error,target type not in:%v")

	// yggdrasil dispatch
	ErrYggdrasilDispatchTaskNotFound                       = errors.NewCode(ErrCodeYggdrasilDispatchTaskNotFound, "Yggdrasil dispatch task not found, taskId: %d")
	ErrYggdrasilDispatchTaskStateNotReadyForMission        = errors.NewCode(ErrCodeYggdrasilDispatchTaskStateNotReadyForMission, "Yggdrasil dispatch task state not ready for mission, taskId: %d, currentState: %d")
	ErrYggdrasilDispatchTaskStateNotOnMission              = errors.NewCode(ErrCodeYggdrasilDispatchTaskStateNotOnMission, "Yggdrasil dispatch task state not on mission, taskId: %d, currentState: %d")
	ErrYggdrasilDispatchTaskStateNotReadyForReceiveRewards = errors.NewCode(ErrCodeYggdrasilDispatchTaskStateNotReadyForReceiveRewards, "Yggdrasil dispatch task state not ready for rewards, taskId: %d, currentState: %d")
	ErrYggdrasilDispatchNecessaryConditionNotSatisfied     = errors.NewCode(ErrCodeYggdrasilDispatchNecessaryConditionNotSatisfied, "Yggdrasil dispatch necessary condition not satisfied, taskId: %d")
	ErrYggdrasilDispatchCharacterIsOnMission               = errors.NewCode(ErrCodeYggdrasilDispatchCharacterIsOnMission, "Yggdrasil dispatch character is on mission, characterId: %d")
	ErrYggdrasilDispatchCharacterLevelNotArrival           = errors.NewCode(ErrCodeYggdrasilDispatchCharacterLevelNotArrival, "Yggdrasil dispatch character level not arrival")
	ErrYggdrasilDispatchCharacterCampNotArrival            = errors.NewCode(ErrCodeYggdrasilDispatchCharacterCampNotArrival, "Yggdrasil dispatch character camp not arrival")
	ErrYggdrasilDispatchCharacterCareerNotArrival          = errors.NewCode(ErrCodeYggdrasilDispatchCharacterCareerNotArrival, "Yggdrasil dispatch character career not arrival")
	ErrYggdrasilDispatchCharacterRarityNotArrival          = errors.NewCode(ErrCodeYggdrasilDispatchCharacterRarityNotArrival, "Yggdrasil dispatch character rarity not arrival")
	ErrYggdrasilDispatchCharacterStarNotArrival            = errors.NewCode(ErrCodeYggdrasilDispatchCharacterStarNotArrival, "Yggdrasil dispatch character star not arrival")
	ErrYggdrasilDispatchCharacterPowerNotArrival           = errors.NewCode(ErrCodeYggdrasilDispatchCharacterPowerNotArrival, "Yggdrasil dispatch character power not arrival")
	ErrYggdrasilDispatchCampNotArrival                     = errors.NewCode(ErrCodeYggdrasilDispatchCampNotArrival, "Yggdrasil dispatch camp not arrival")
	ErrYggdrasilDispatchCareerNotArrival                   = errors.NewCode(ErrCodeYggdrasilDispatchCareerNotArrival, "Yggdrasil dispatch career not arrival")
	ErrYggdrasilDispatchSpecificCharacterNotArrival        = errors.NewCode(ErrCodeYggdrasilDispatchSpecificCharacterNotArrival, "Yggdrasil dispatch specific character not arrival")
	ErrYggDispatchCSVExtraConditionNotMatchRewards         = errors.NewCode(ErrCodeYggDispatchCSVExtraConditionNotMatchRewards, "Yggdrasil dispatch csv extra condition not match rewards, taskId: %d")
	ErrYggDispatchCSVDispatchTypeNotMatchCloseTime         = errors.NewCode(ErrCodeYggDispatchCSVDispatchTypeNotMatchCloseTime, "Yggdrasil dispatch csv dispatch type not match close time , taskId: %d")
	ErrYggDispatchCSVWrongDispatchType                     = errors.NewCode(ErrCodeYggDispatchCSVWrongDispatchType, "Yggdrasil dispatch csv wrong dispatch type, taskId: %d")
	ErrYggDispatchCSVWrongLengthOfGuildNum                 = errors.NewCode(ErrCodeYggDispatchCSVWrongLengthOfGuildNum, "Yggdrasil dispatch csv wrong length of guild num, areaId: %d")
	ErrYggDispatchWrongTeamSizeForDispatch                 = errors.NewCode(ErrCodeYggDispatchWrongTeamSizeForDispatch, "Yggdrasil dispatch wrong team size for this task, taskId: %d")
	ErrYggDispatchCSVWrongGuildCharacter                   = errors.NewCode(ErrCodeYggDispatchCSVWrongGuildCharacter, "Yggdrasil dispatch csv wrong guild character, taskId: %d, guildCharacNum: %d, guildCharaId: %d")

	// greetings
	ErrGreetingsNotFoundInCsv = errors.NewCode(ErrCodeGreetingsNotFoundInCsv, "Greetings not found in csv, id: %d")

	// activity
	ErrActivityCfgNotFound     = errors.NewCode(ErrCodeActivityCfgNotFound, "activity config not found, id: %d")
	ErrActivityFuncCfgNotFound = errors.NewCode(ErrCodeActivityFuncCfgNotFound, "activity func config not found, id: %d")

	//score pass
	ErrScorePassPhaseCfgNotFound  = errors.NewCode(ErrCodeScorePassPhaseCfgNotFound, "score pass phase cfg not found, id: %d")
	ErrScorePassGroupCfgNotFound  = errors.NewCode(ErrCodeScorePassGroupCfgNotFound, "score pass group cfg not found, id: %d")
	ErrScorePassRewardCfgNotFound = errors.NewCode(ErrCodeScorePassRewardCfgNotFound, "score pass reward cfg not found, id: %d")
	ErrScorePassNoSuchSeason      = errors.NewCode(ErrCodeScorePassNoSuchSeason, "score pass no such season, id: %d")
	ErrScorePassPhaseNotStart     = errors.NewCode(ErrCodeScorePassPhaseNotStart, "score pass phase not start, id: %d")

	//battle
	ErrBattleCfgNotFound       = errors.NewCode(ErrCodeBattleCfgNotFound, "battle config not found, id: %d")
	ErrBattleNPCCfgNotFound    = errors.NewCode(ErrCodeBattleNPCCfgNotFound, "battle npc config not found id: %d")
	ErrBattleInvalidNPC        = errors.NewCode(ErrCodeBattleInvalidNPC, "invalid battle npc: %d")
	ErrBattleInvalidPosition   = errors.NewCode(ErrCodeBattleInvalidPosition, "invalid battle formation position position: %d")
	ErrBattleDuplicatePosition = errors.NewCode(ErrCodeBattleDuplicatePosition, "invalid duplicate battle formation position: %d")
	ErrBattleDuplicateChara    = errors.NewCode(ErrCodeBattleDuplicateChara, "invalid duplicate battle formation character: %d")
	ErrBattleInvalidEndType    = errors.NewCode(ErrCodeBattleInvalidEndType, "invalid battle end type: %d")

	// mercenary
	ErrWrongSystemType              = errors.NewCode(ErrCodeMercenaryWrongSystemType, "mercenary wrong system type, type: %d")
	ErrMercenaryExceedUseLimit      = errors.NewCode(ErrCodeMercenaryExceedUseLimit, "mercenary exceed user limit, characterId: %d, systemType: %d, count: %d")
	ErrMercenaryNotFound            = errors.NewCode(ErrCodeMercenaryNotFound, "mercenary not found. id: %d")
	ErrMercenaryExceedNumLimit      = errors.NewCode(ErrCodeMercenaryExceedNumLimit, "mercenary exceed limit of number")
	ErrMercenaryApplyExceedLimit    = errors.NewCode(ErrCodeMercenaryApplyExceedLimit, "mercenary apply exceed limit of number")
	ErrMercenarySendApplyNotFound   = errors.NewCode(ErrCodeMercenarySendApplyNotFound, "mercenary send apply not found, characterId: %d, userTo: %d")
	ErrMercenaryHandleApplyNotFound = errors.NewCode(ErrCodeMercenaryHandleApplyNotFound, "mercenary handle apply not found, characterId: %d, userTo: %d")
	ErrMercenaryRepeatedSendApply   = errors.NewCode(ErrCodeMercenaryRepeatedSendApply, "mercenary repeated send apply, send to: %d, characterId: %d")
	ErrMercenaryAlreadyBorrowed     = errors.NewCode(ErrCodeMercenaryAlreadyBorrowed, "Mercenary already borrowed, characterId: %d")
	ErrMercenaryAlreadyHad          = errors.NewCode(ErrCodeMercenaryAlreadyHad, "mercenary already had, characterId: %d")
	// formation
	ErrInvalidFormation  = errors.NewCode(ErrCodeInvalidFormation, "invalid formation type: %d, id: %d")
	ErrInvalidFormationS = errors.NewCode(ErrCodeInvalidFormation, "invalid formation")
)
