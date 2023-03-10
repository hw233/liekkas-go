package entry

import (
	"reflect"

	"shared/csv/base"
	"shared/utility/errors"
	"shared/utility/glog"
)

// todo:reload
type Manager struct {
	// 必须先Reward
	Reward *Reward

	Drop                 *Drop
	Character            *Character
	Protocol             *Protocol
	GM                   *GM
	Equipment            *EquipmentEntry
	WorldItem            *WorldItemEntry
	Quest                *QuestEntry
	Manual               *Manual
	Item                 *Item
	TeamLevelCache       *TeamLevelCache
	GraveyardEntry       *GraveyardEntry
	SignIn               *SignInEntry
	LevelsEntry          *LevelsEntry
	ChapterEntry         *ChapterEntry
	ExploreEntry         *ExploreEntry
	Survey               *SurveyEntry
	Store                *StoreEntry
	QuickPurchaseStamina *StaminaEntry
	GlobalEntry          *GlobalEntry
	Gacha                *GachaEntry
	Hero                 *HeroEntry
	Tower                *TowerEntry
	VN                   *VnEntry
	Guide                *GuideEntry
	Mail                 *MailEntry
	Yggdrasil            *YggdrasilEntry
	ActionUnlock         *ActionUnlockEntry
	PowerAdapt           *PowerAdaptEntry
	UserSetting          *UserSettingEntry
	Activities           *Activities
	ScorePasses          *ScorePasses
	GreetingsEntry       *GreetingsEntry
	Guild                *GuildEntry
	Battle               *Battle
}

func NewManager() *Manager {
	return &Manager{
		Reward:               NewReward(),
		Drop:                 NewDrop(),
		Character:            NewCharacter(),
		Protocol:             NewProtocol(),
		GM:                   NewGM(),
		Equipment:            NewEquipmentEntry(),
		WorldItem:            NewWorldItemEntry(),
		Quest:                NewQuestEntry(),
		Manual:               NewManual(),
		Item:                 NewItem(),
		TeamLevelCache:       NewTeamLevelCache(),
		GraveyardEntry:       NewGraveyardEntry(),
		SignIn:               NewSignEntry(),
		LevelsEntry:          NewLevelsEntry(),
		ChapterEntry:         NewChapterEntry(),
		ExploreEntry:         NewExploreEntry(),
		Survey:               NewSurveyEntry(),
		Store:                NewStoreEntry(),
		QuickPurchaseStamina: NewStaminaEntry(),
		GlobalEntry:          NewGlobalEntry(),
		Gacha:                NewGachaEntry(),
		Hero:                 NewHeroEntry(),
		Tower:                NewTowerEntry(),
		VN:                   NewVnEntry(),
		Guide:                NewGuideEntry(),
		Mail:                 NewMailEntry(),
		Yggdrasil:            NewYggdrasilEntry(),
		ActionUnlock:         NewActionUnlockEntry(),
		PowerAdapt:           NewPowerAdaptEntry(),
		UserSetting:          NewUserSettingEntry(),
		Activities:           NewActivities(),
		ScorePasses:          NewScorePasses(),
		GreetingsEntry:       NewGreetingsEntry(),
		Guild:                NewGuildEntry(),
		Battle:               NewBattle(),
	}
}

type Config struct {
	*GlobalEntry
	*base.ConfigManager
}

type Entry interface {
	// Check(config *Config) error
	Reload(config *Config) error
}

func (m *Manager) Reload(config *base.ConfigManager) error {
	// load global first
	err := m.GlobalEntry.Reload(config)
	if err != nil {
		glog.Errorf("ERROR: load GlobalEntry error: %+v", err)
		return errors.WrapTrace(err)
	}

	v := reflect.ValueOf(m).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		data := field.Interface()

		entry, ok := data.(Entry)
		if ok {
			err := entry.Reload(&Config{
				GlobalEntry:   m.GlobalEntry,
				ConfigManager: config,
			})
			if err != nil {
				glog.Errorf("ERROR: load %s error: %+v", field.Type(), errors.Format(err))
				return errors.WrapTrace(err)
			}

			// log.Printf("reload %+v", entry)
		}
	}

	return nil
}
