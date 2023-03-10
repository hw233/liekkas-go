package entry

import (
	"encoding/json"
	"reflect"
	"shared/common"
	"shared/csv/base"
	"shared/utility/coordinate"
	"shared/utility/errors"
	"time"
)

const (
	GlobalCfgStringConfig = "cfg_global_string_config"
)

type GlobalEntry struct {
	EquipmentMaterialEXPDiscount []EquipmentMaterialEXPDiscount
	EquipmentMaxLevel            map[int8]int32
	WorldItemMaterialEXPDiscount []WorldItemMaterialEXPDiscount
	WorldItemMaxLevel            map[int8]int32
	EquipmentAttrUnlockLevel     []int32
	MainTowerPosition            coordinate.Position
	DailyRefreshHour             int32
	GraveyardPopulationAddition  int32
	InitUserHero                 int32
	InitUserDrop                 int32
	GachaMultipleValue           int32
	GraveyardPlotRewardMaxNum    int32
	GraveyardPlotRewardHour      []int32
	GraveyardPlotRewardDropId    int32
	MailMaxCount                 int32
	YggBagAllCount               int32
	YggApRestoreTime             int64
	YggOtherMaxTaskCount         int32
	YggBuildHelpCount            int32
	YggInitCity                  int32
	YggCharacterRestHour         int32
	YggLightRadius               int32
	YggMarkBarrageLimit          int32
	YggCityBuildLimitRadius      int32
	YggMatchGrabAlgorithmParam   int32
	YggBlockLengthAndWidth       []int32
	YggInitPos                   coordinate.Position
	YggDailyTravelTime           int32
	StaminaRecoverSeconds        int32
	StaminaExpRatio              int32
	YggMarkTotalCount            int32
	YggMailMaxCount              int32
	YggEditTeamLevel             []int32
	YggBagLockLevel              []int32
	YggEditTeamMaxCount          int32
	MailReceiveMaxCount          int32
	NicknameChangeConsume        int32
	CharacterGreetingsLimit      int32
	WorldItemGreetingsLimit      int32
	GreetingsAliveDuration       int32
	GraveyardMinHelpSec          int32
	GraveyardSingleHelpCount     int32
	GraveyardDailyHelpSendCount  int32
	GraveyardMinHelpPercent      int32
	GuildViceChairmenNum         int32
	GuildElitenNum               int32
	GuildContributionExpRatio    int32
	GuildApplyNumLimit           int32
	MercenaryUseLimit            int32
	MercenaryNumLimit            int32
	GuildChatLimit               int32
	GuildTaskRefreshWeekDay      []int32
	GuildCreateCost              []int32
	GachaRecordShowMonth         int32
	GachaRecordStoreMonth        int32
	CBT1RegisterLimit            int64
	DailyRefreshOffset           time.Duration `ignore:"true"`
	EquipmentCountLimit          int
	WorlItemCountLimit           int
	GuildDissolveCD              int64
	GuildQuitCD                  int64
	GuildDissolveDuration        int64
	GuildMemberMaxNum            int32
	SpareMaintainSwitch          bool
}

func NewGlobalEntry() *GlobalEntry {
	return &GlobalEntry{
		EquipmentMaterialEXPDiscount: []EquipmentMaterialEXPDiscount{},
		EquipmentMaxLevel:            map[int8]int32{},
		WorldItemMaterialEXPDiscount: []WorldItemMaterialEXPDiscount{},
		WorldItemMaxLevel:            map[int8]int32{},
		EquipmentAttrUnlockLevel:     []int32{},
		MainTowerPosition:            coordinate.Position{},
		DailyRefreshHour:             0,
		GraveyardPopulationAddition:  0,
		InitUserHero:                 0,
		InitUserDrop:                 0,
		GachaMultipleValue:           0,
		GraveyardPlotRewardMaxNum:    0,
		GraveyardPlotRewardHour:      []int32{},
		GraveyardPlotRewardDropId:    0,
		MailMaxCount:                 0,
		YggBagAllCount:               0,
		YggApRestoreTime:             0,
		YggOtherMaxTaskCount:         0,
		YggBuildHelpCount:            0,
		YggInitCity:                  0,
		YggCharacterRestHour:         0,
		YggLightRadius:               0,
		YggMarkBarrageLimit:          0,
		YggCityBuildLimitRadius:      0,
		YggMatchGrabAlgorithmParam:   0,
		YggBlockLengthAndWidth:       []int32{},
		YggInitPos:                   coordinate.Position{},
		YggDailyTravelTime:           0,
		StaminaRecoverSeconds:        0,
		StaminaExpRatio:              0,
		YggMarkTotalCount:            0,
		YggMailMaxCount:              0,
		MailReceiveMaxCount:          0,
		NicknameChangeConsume:        0,
		CharacterGreetingsLimit:      0,
		WorldItemGreetingsLimit:      0,
		GreetingsAliveDuration:       0,
		GraveyardMinHelpSec:          0,
		GraveyardSingleHelpCount:     0,
		GraveyardDailyHelpSendCount:  0,
		GraveyardMinHelpPercent:      0,
		GuildViceChairmenNum:         0,
		GuildElitenNum:               0,
		GuildContributionExpRatio:    0,
		GuildApplyNumLimit:           0,
		MercenaryUseLimit:            0,
		MercenaryNumLimit:            0,
		GuildChatLimit:               0,
		GuildTaskRefreshWeekDay:      []int32{},
		GuildCreateCost:              []int32{},
		GachaRecordShowMonth:         0,
		GachaRecordStoreMonth:        0,
		CBT1RegisterLimit:            0,
		EquipmentCountLimit:          0,
		WorlItemCountLimit:           0,
		GuildDissolveCD:              0,
		GuildDissolveDuration:        0,
		GuildQuitCD:                  0,
		GuildMemberMaxNum:            0,
		SpareMaintainSwitch:          true,
	}
}

func (e *GlobalEntry) Check(config *Config) error {
	kv := map[string]string{}

	for _, global := range config.CfgGlobalStringConfig.GetAllData() {
		kv[global.Key] = global.Value
	}

	typ := reflect.TypeOf(e).Elem()
	val := reflect.ValueOf(e).Elem()

	for i := 0; i < typ.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldTyp := typ.Field(i)

		if fieldTyp.Tag.Get("ignore") == "true" {
			continue
		}

		name := fieldTyp.Name

		v, ok := kv[name]
		if !ok {
			return errors.Swrapf(common.ErrNotFoundInCSV, GlobalCfgStringConfig, name)
		}

		newVal := reflect.New(fieldVal.Type())

		err := json.Unmarshal([]byte(v), &newVal)
		if err != nil {
			return errors.Wrap(err, "Field unmarshal error. (name: %s", name)
		}

		if !fieldVal.CanSet() {
			return errors.New("Field can't set. (name: %s", name)
		}
	}

	return nil
}

func (e *GlobalEntry) Reload(config *base.ConfigManager) error {
	kv := map[string]string{}

	for _, global := range config.CfgGlobalStringConfig.GetAllData() {
		kv[global.Key] = global.Value
	}

	typ := reflect.TypeOf(e).Elem()
	val := reflect.ValueOf(e).Elem()

	for i := 0; i < typ.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldTyp := typ.Field(i)

		if fieldTyp.Tag.Get("ignore") == "true" {
			continue
		}

		name := fieldTyp.Name

		v, ok := kv[name]
		if !ok {
			return errors.Swrapf(common.ErrNotFoundInCSV, GlobalCfgStringConfig, name)
		}

		newVal := reflect.New(fieldVal.Type()).Interface()

		err := json.Unmarshal([]byte(v), &newVal)
		if err != nil {
			return errors.Wrap(err, "Field unmarshal error. (name: %s", name)
		}

		if !fieldVal.CanSet() {
			return errors.New("Field can't set. (name: %s", name)
		}

		fieldVal.Set(reflect.ValueOf(newVal).Elem())
	}

	e.DailyRefreshOffset = time.Hour * time.Duration(e.DailyRefreshHour)

	return nil
}

func (e *GlobalEntry) DailyRefreshTimeOffset() time.Duration {
	return e.DailyRefreshOffset
}
