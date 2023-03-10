package model

import (
	"gamesvr/manager"
	"reflect"
	"shared/common"
	"shared/csv/entry"
	"shared/csv/static"
	"shared/utility/errors"
	"shared/utility/slice"
)

func questConvertParams(paramsObjPtr interface{}, params []interface{}) error {
	objPtrType := reflect.TypeOf(paramsObjPtr)
	if objPtrType.Kind() != reflect.Ptr {
		return errors.New("type of dst isn't ptr")
	}

	objType := objPtrType.Elem()

	fieldNum := objType.NumField()
	if fieldNum != len(params) {
		return errors.Swrapf(common.ErrQuestUpdateParamCountInvalid, fieldNum)
	}

	objMembers := reflect.ValueOf(paramsObjPtr).Elem()

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		paramType := reflect.TypeOf(params[i])
		if field.Type.AssignableTo(paramType) {
			objMember := objMembers.Field(i)
			paramValue := reflect.ValueOf(params[i])

			objMember.Set(paramValue)
		} else if paramType.ConvertibleTo(field.Type) {
			objMember := objMembers.Field(i)
			paramValue := reflect.ValueOf(params[i])
			objMember.Set(paramValue.Convert(field.Type))
		} else {
			return errors.Swrapf(common.ErrQuestUpdateParamInvalid, i, field.Type.Name())
		}

	}

	return nil
}

//----------------------------------------
//QuestUnlockCondition
//----------------------------------------
func (u *User) checkQuestUnlockCondition(questId int32) bool {
	questCfg, err := manager.CSV.Quest.Get(questId)
	if err != nil {
		return false
	}

	if u.CheckUserConditions(questCfg.UnlockCondition) != nil {
		return false
	}

	return true
}

//----------------------------------------
//QuestProgress
//----------------------------------------
func questProgressCounterMaxReplace(oldValue, newValue int32) int32 {
	if newValue > oldValue {
		return newValue
	}

	return oldValue
}

func questProgressCounterAdd(progress, incrument int32) int32 {
	return progress + incrument
}

//----------------------------------------
//QuestCondition
//----------------------------------------
func (u *User) calcNewQuestProgress(questCfg *entry.Quest, curProgress int32,
	params ...interface{}) (int32, error) {
	progress := curProgress

	var err error = nil

	switch questCfg.ConditionType {
	case static.TaskTypeAccountLevel:
		paramsObj := struct {
			NewLevel int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		progress = questProgressCounterMaxReplace(progress, paramsObj.NewLevel)

	case static.TaskTypeDailyLogin:
		progress = questProgressCounterAdd(progress, 1)

	case static.TaskTypeCostItem:
		paramsObj := struct {
			ItemId int32
			Count  int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}
		triggerParams := []int32{paramsObj.ItemId}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, paramsObj.Count)
		}

	case static.TaskTypeCharacterLevelUpCount:
		paramsObj := struct {
			Chara    *Character
			OldLevel int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{paramsObj.Chara.GetRare()}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, paramsObj.Chara.GetLevel()-paramsObj.OldLevel)
		}

	case static.TaskTypeCharacterLevelCount:
		paramsObj := struct {
			Chara    *Character
			OldLevel int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{paramsObj.Chara.GetRare()}
		if questConditionCheckValueUp(questCfg.ConditionParams, triggerParams,
			paramsObj.OldLevel, paramsObj.Chara.GetLevel()) {
			progress = questProgressCounterAdd(progress, 1)
		}

	case static.TaskTypeCharacterSkillLevelCount:
		paramsObj := struct {
			CharaRare  int32
			SkillType  int32
			SkillLevel int32
			OldLevel   int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{paramsObj.CharaRare, paramsObj.SkillType}
		if questConditionCheckValueUp(questCfg.ConditionParams, triggerParams,
			paramsObj.OldLevel, paramsObj.SkillLevel) {
			progress = questProgressCounterAdd(progress, 1)
		}

	case static.TaskTypeCharacterStageCount:
		paramsObj := struct {
			CharaRare int32
			Stage     int32
			OldStage  int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{paramsObj.CharaRare, paramsObj.Stage}
		if questConditionCheckValueUp(questCfg.ConditionParams, triggerParams,
			paramsObj.OldStage, paramsObj.Stage) {
			progress = questProgressCounterAdd(progress, 1)
		}

	case static.TaskTypeCharacterLevel:
		paramsObj := struct {
			CharaId int32
			Level   int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}
		triggerParams := []int32{paramsObj.CharaId}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterMaxReplace(progress, paramsObj.Level)
		}

	case static.TaskTypeCharacterSkillLevel:
		paramsObj := struct {
			CharaId    int32
			SkillNum   int32
			SkillLevel int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{paramsObj.CharaId, paramsObj.SkillNum}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterMaxReplace(progress, paramsObj.SkillLevel)
		}

	case static.TaskTypeCharacterStage:
		paramsObj := struct {
			CharaId int32
			Stage   int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{paramsObj.CharaId}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterMaxReplace(progress, paramsObj.Stage)
		}

	case static.TaskTypeCharacterStar:
		paramsObj := struct {
			CharaId int32
			Star    int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{paramsObj.CharaId}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterMaxReplace(progress, paramsObj.Star)
		}

	case static.TaskTypeCharacterLevelStrengthenTimes:
		paramsObj := struct {
			CharaRare int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{paramsObj.CharaRare}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, 1)
		}

	case static.TaskTypeCharacterCampCount:
		paramsObj := struct {
			Rare int32
			Camp int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{paramsObj.Rare, paramsObj.Camp}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, 1)
		}

	case static.TaskTypeHasCharacters:
		paramsObj := struct {
			CharaId int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{paramsObj.CharaId}
		if questConditionCheckConfigContains(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, 1)
		}

	case static.TaskTypeEquipmentLevelUpCount:
		paramsObj := struct {
			Rare     int32
			OldLevel int32
			NewLevel int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{paramsObj.Rare}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, paramsObj.NewLevel-paramsObj.OldLevel)
		}

	case static.TaskTypeEquipmentLevelCount:
		paramsObj := struct {
			Rare     int32
			OldLevel int32
			NewLevel int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{paramsObj.Rare}
		if questConditionCheckValueUp(questCfg.ConditionParams, triggerParams,
			paramsObj.OldLevel, paramsObj.NewLevel) {
			progress = questProgressCounterAdd(progress, 1)
		}

	case static.TaskTypeEquipmentExpUpTimes:
		paramsObj := struct {
			Rare int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{paramsObj.Rare}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, 1)
		}

	case static.TaskTypeLevelPass:
		paramsObj := struct {
			LevelId int32
			Score   int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{paramsObj.LevelId, paramsObj.Score}

		if questConditionCheckLastValueBigger(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, 1)
		}

	case static.TaskTypeLevelCharacter:
		paramsObj := struct {
			LevelId  int32
			CharaIds []int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{paramsObj.LevelId}
		triggerParams = append(triggerParams, paramsObj.CharaIds...)

		if questConditionCheckTriggerContains(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, 1)
		}

	case static.TaskTypeLevelEndHero:
		paramsObj := struct {
			HeroId int32
			Times  int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{paramsObj.HeroId}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, paramsObj.Times)
		}

	case static.TaskTypeLevelPassByType,
		static.TaskTypeLevelEndByType:
		paramsObj := struct {
			LevelType int32
			Times     int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{paramsObj.LevelType}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, paramsObj.Times)
		}

	case static.TaskTypeBattleHeroDamage,
		static.TaskTypeBattleHeroCure,
		static.TaskTypeBattleHeroHurt:
		paramsObj := struct {
			HeroId int32
			Value  int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}
		triggerParams := []int32{paramsObj.HeroId}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, paramsObj.Value)
		}

	case static.TaskTypeBattleHeroSkillCast:
		paramsObj := struct {
			HeroId    int32
			SkillId   int32
			CastTimes int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}
		triggerParams := []int32{paramsObj.HeroId, paramsObj.SkillId}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, paramsObj.CastTimes)
		}

	case static.TaskTypeBattleCharaDamage,
		static.TaskTypeBattleCharaCure,
		static.TaskTypeBattleCharaHurt:
		paramsObj := struct {
			CharaId int32
			Value   int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}
		triggerParams := []int32{paramsObj.CharaId}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, paramsObj.Value)
		}

	case static.TaskTypeBattleCharaSkillCast:
		paramsObj := struct {
			CharaId   int32
			SkillNum  int32
			CastTimes int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}
		triggerParams := []int32{paramsObj.CharaId, paramsObj.SkillNum}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, paramsObj.CastTimes)
		}

	case static.TaskTypeHeroLevel:
		paramsObj := struct {
			Level int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{}

		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterMaxReplace(progress, paramsObj.Level)
		}

	case static.TaskTypeHeroSkillUnlockCount,
		static.TaskTypeHeroSkillLevelMaxCount:
		paramsObj := struct {
			HeroId int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{paramsObj.HeroId}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, 1)
		}

	case static.TaskTypeGraveyarProductCollectTimes:
		paramsObj := struct {
			Count int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}
		progress = questProgressCounterAdd(progress, paramsObj.Count)

	case static.TaskTypeGraveyarAccelerateTimes:
		paramsObj := struct {
			BuildId int32
			Times   int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}
		triggerParams := []int32{paramsObj.BuildId}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, paramsObj.Times)
		}

	case static.TaskTypeGraveyarAccelerateTime:
		paramsObj := struct {
			BuildId int32
			AccTime int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}
		triggerParams := []int32{paramsObj.BuildId}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, paramsObj.AccTime)
		}

	case static.TaskTypeGraveyarBuildingLevelCount,
		static.TaskTypeGraveyarBuildingStageCount:
		paramsObj := struct {
			BuildingId int32
			OldValue   int32
			NewValue   int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{paramsObj.BuildingId}
		if questConditionCheckValueUp(questCfg.ConditionParams, triggerParams,
			paramsObj.OldValue, paramsObj.NewValue) {
			progress = questProgressCounterAdd(progress, 1)
		}

	case static.TaskTypeGraveyarDecorationSetting,
		static.TaskTypeGraveyarCharacterSetting:
		paramsObj := struct {
			Count int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		progress = questProgressCounterAdd(progress, paramsObj.Count)

	case static.TaskTypeGraveyarProductRewardType,
		static.TaskTypeGraveyarProductEquipmentPart,
		static.TaskTypeGraveyarProductItem:
		paramsObj := struct {
			Type  int32
			Count int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{paramsObj.Type}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, paramsObj.Count)
		}

	case static.TaskTypeGachaTimes:
		paramsObj := struct {
			GachaType int32
			Times     int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}
		triggerParams := []int32{paramsObj.GachaType}
		if questConditionCheckConfigContains(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, paramsObj.Times)
		}

	case static.TaskTypeGachaRare:
		paramsObj := struct {
			GachaRewards []*common.GachaReward
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		if questConditionGachaRare(questCfg.ConditionParams, paramsObj.GachaRewards) {
			progress = questProgressCounterAdd(progress, 1)
		}

	case static.TaskTypeGachaSame:
		paramsObj := struct {
			GachaRewards []*common.GachaReward
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		var sameTimes int32 = 0
		if questConditionGachaSame(questCfg.ConditionParams, paramsObj.GachaRewards, &sameTimes) {
			progress = questProgressCounterAdd(progress, sameTimes)
		}

	case static.TaskTypeGachaTimesSsr:
		paramsObj := struct {
			GachaRewards []*common.GachaReward
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		for _, gachaReward := range paramsObj.GachaRewards {
			if questConditionGachaTimesSSR(questCfg.ConditionParams, gachaReward) {
				progress = questProgressCounterAdd(progress, 1)
			}
		}

	case static.TaskTypeGachaItemTimes:
		paramsObj := struct {
			GachaRewards []*common.GachaReward
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		for _, gachaReward := range paramsObj.GachaRewards {
			triggerParams := []int32{gachaReward.Type}
			if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
				progress = questProgressCounterAdd(progress, gachaReward.Num)
			}
		}

	case static.TaskTypeWorldItemLevelupTimes:
		paramsObj := struct {
			Rare int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}
		triggerParams := []int32{paramsObj.Rare}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, 1)
		}

	case static.TaskTypeWorldLevelCount,
		static.TaskTypeWorldStarCount:
		paramsObj := struct {
			Rare     int32
			OldValue int32
			NewValue int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		triggerParams := []int32{paramsObj.Rare}
		if questConditionCheckValueUp(questCfg.ConditionParams, triggerParams,
			paramsObj.OldValue, paramsObj.NewValue) {
			progress = questProgressCounterAdd(progress, 1)
		}

	case static.TaskTypeTowerStagePassed:
		paramsObj := struct {
			TowerId int32
			Stage   int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}
		triggerParams := []int32{paramsObj.TowerId}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterMaxReplace(progress, paramsObj.Stage)
		}

	case static.TaskTypeChapterScore:
		paramsObj := struct {
			ChapterType int32
			Score       int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}
		triggerParams := []int32{paramsObj.ChapterType}
		if questConditionCheckSame(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, paramsObj.Score)
		}

	case static.TaskTypeManualCollect:
		paramsObj := struct {
			ManualId int32
		}{}

		err = questConvertParams(&paramsObj, params)
		if err != nil {
			break
		}

		var manaulVersion int32
		manaulVersion, err = manager.CSV.Manual.GetManualVersion(paramsObj.ManualId)
		if err != nil {
			break
		}

		triggerParams := []int32{manaulVersion}
		if questConditionCheckConfigContains(questCfg.ConditionParams, triggerParams) {
			progress = questProgressCounterAdd(progress, 1)
		}

	default:
	}

	return progress, err
}

func questConditionCheckSame(cfgParams []int32, triggerParams []int32) bool {
	paramLen := len(cfgParams)
	if paramLen != len(triggerParams) {
		return false
	}

	for i := 0; i < paramLen; i++ {
		if cfgParams[i] != 0 && triggerParams[i] != 0 && cfgParams[i] != triggerParams[i] {
			return false
		}
	}

	return true
}

func questConditionCheckValueUp(cfgParams []int32, triggerParams []int32,
	oldValue, newValue int32) bool {

	paramLen := len(cfgParams)
	for i := 0; i < paramLen-1; i++ {
		if cfgParams[i] != 0 && triggerParams[i] != 0 && cfgParams[i] != triggerParams[i] {
			return false
		}
	}

	lastCfgParam := cfgParams[paramLen-1]
	if newValue < lastCfgParam || oldValue >= lastCfgParam ||
		newValue == oldValue {
		return false
	}

	return true
}

func questConditionCheckLastValueBigger(cfgParams []int32, triggerParams []int32) bool {
	paramLen := len(cfgParams)
	for i := 0; i < paramLen-1; i++ {
		if cfgParams[i] != 0 && triggerParams[i] != 0 && cfgParams[i] != triggerParams[i] {
			return false
		}
	}

	if triggerParams[paramLen-1] < cfgParams[paramLen-1] {
		return false
	}

	return true
}

func questConditionCheckTriggerContains(cfgParams []int32, triggerParams []int32) bool {
	if len(cfgParams) > len(triggerParams) {
		return false
	}

	zeroCount := 0
	for idx, param := range cfgParams {
		triggerParam := triggerParams[idx]
		if param != triggerParam && param*triggerParam == 0 {
			zeroCount = zeroCount + 1
		}
	}

	intersection := slice.GetSliceInt32Intersection(triggerParams, cfgParams)
	if len(intersection)+zeroCount < len(cfgParams) {
		return false
	}

	return true
}

func questConditionCheckConfigContains(cfgParams []int32, triggerParams []int32) bool {
	if len(triggerParams) > len(cfgParams) {
		return false
	}

	zeroCount := 0
	for idx, param := range cfgParams {
		triggerParam := triggerParams[idx]
		if param != triggerParam && param*triggerParam == 0 {
			zeroCount = zeroCount + 1
		}
	}

	intersection := slice.GetSliceInt32Intersection(cfgParams, triggerParams)
	if len(intersection)+zeroCount < len(triggerParams) {
		return false
	}

	return true
}

func questConditionGachaRare(cfgParams []int32, gachaRewards []*common.GachaReward) bool {
	cfgRare := cfgParams[0]
	cfgRewardType := cfgParams[1]
	cfgCount := cfgParams[2]

	var count int32 = 0
	for _, gachaReward := range gachaRewards {
		if gachaReward.Type != cfgRewardType {
			continue
		}

		itemCfg, ok := manager.CSV.Item.GetItem(gachaReward.ID)
		if !ok {
			return false
		}

		if itemCfg.Rarity == cfgRare {
			count = count + 1
		}
	}

	return count >= cfgCount
}

func questConditionGachaSame(cfgParams []int32, gachaRewards []*common.GachaReward, sameTimes *int32) bool {
	cfgRewardType := cfgParams[0]
	cfgCount := cfgParams[1]

	*sameTimes = 0

	statistic := map[int32]int32{}
	for _, gachaReward := range gachaRewards {
		if gachaReward.Type == cfgRewardType {
			continue
		}

		count, ok := statistic[gachaReward.ID]
		if !ok {
			statistic[gachaReward.ID] = 1
		} else {
			statistic[gachaReward.ID] = count + 1
		}
	}

	for _, count := range statistic {
		if count >= cfgCount {
			*sameTimes = *sameTimes + 1
		}
	}

	return *sameTimes > 0
}

func questConditionGachaTimesSSR(cfgParams []int32, gachaReward *common.GachaReward) bool {
	cfgTimes := cfgParams[0]
	cfgRewardType := cfgParams[1]

	if gachaReward.Type != cfgRewardType {
		return false
	}

	itemCfg, ok := manager.CSV.Item.GetItem(gachaReward.ID)
	if !ok {
		return false
	}

	if itemCfg.Rarity == static.RaritySsr && gachaReward.SSRMissCount >= cfgTimes {
		return true
	}

	return false
}
