package entry

import (
	"sync"

	"shared/common"
	"shared/csv/static"
	"shared/utility/errors"
	"shared/utility/transfer"
)

type QuestUnlockConditon struct {
	ConditionType int32
	Params        []int32
}

type Quest struct {
	Id              int32
	Module          int32
	UnlockCondition *common.Conditions `rule:"conditions"`
	ConditionType   int32
	ConditionParams []int32
	TargetCount     int32
	IsAutoReceive   bool `rule:"intToBool"`
	AddPoint        int32
	DropId          int32
	Group           int32
	NextQuests      []int32 `ignore:"true"`
	PreQuests       []int32 `ignore:"true"`
	Cname           string
}

type QuestActivity struct {
	Id        int32
	Module    int32
	Condition int32
	DropId    int32
}

type QuestEntry struct {
	sync.RWMutex

	Quests        map[int32]*Quest
	QuestActivity map[int32]*QuestActivity

	AccountLevelUnlocks    map[int32][]int32
	PassLevelUnlocks       map[int32][]int32
	CommonConditionUnlocks map[int32][]int32

	ModuleToQuests        map[int32][]int32
	achievementTotalCount int32
	Groups                map[int32][]int32
}

func NewQuestEntry() *QuestEntry {
	return &QuestEntry{}
}

func (qe *QuestEntry) Check(config *Config) error {
	return nil
}

func (qe *QuestEntry) Reload(config *Config) error {
	qe.Lock()
	defer qe.Unlock()

	quests := map[int32]*Quest{}
	questActivity := map[int32]*QuestActivity{}

	questUnlocks := map[int32][]int32{}
	accountLevelUnlocks := map[int32][]int32{}
	passLevelUnlocks := map[int32][]int32{}
	conditionUnlocks := map[int32][]int32{}

	moduleToQuests := map[int32][]int32{}
	var achievementTotalCount int32
	groups := map[int32][]int32{}

	transferQuestCSV := func(questCSV interface{}) error {
		questCfg := &Quest{
			PreQuests: []int32{},
		}

		err := transfer.Transfer(questCSV, questCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		questId := questCfg.Id

		quests[questId] = questCfg

		module := questCfg.Module
		_, ok := moduleToQuests[module]
		if !ok {
			moduleToQuests[module] = []int32{}
		}

		moduleToQuests[module] = append(moduleToQuests[module], questId)

		groupId := questCfg.Group
		_, ok = groups[groupId]
		if !ok {
			groups[groupId] = []int32{}
		}

		groups[groupId] = append(groups[groupId], questId)

		for _, condition := range *questCfg.UnlockCondition {
			conditionType := condition.ConditionType
			switch conditionType {
			case static.ConditionTypeFinishQuest:
				preQuestId := condition.Params[0]
				unlocks, ok := questUnlocks[preQuestId]
				if !ok {
					unlocks = make([]int32, 0)
				}

				unlocks = append(unlocks, questCfg.Id)
				questUnlocks[preQuestId] = unlocks
				questCfg.PreQuests = append(questCfg.PreQuests, preQuestId)

			case static.ConditionTypeUserLevel:
				level := condition.Params[0]
				unlocks, ok := accountLevelUnlocks[level]
				if !ok {
					unlocks = make([]int32, 0)
				}

				unlocks = append(unlocks, questCfg.Id)
				accountLevelUnlocks[level] = unlocks
			case static.ConditionTypePassLevel:
				levelId := condition.Params[0]
				unlocks, ok := passLevelUnlocks[levelId]
				if !ok {
					unlocks = make([]int32, 0)
				}

				unlocks = append(unlocks, questCfg.Id)
				passLevelUnlocks[levelId] = unlocks

			default:
				unlocks, ok := conditionUnlocks[conditionType]
				if !ok {
					unlocks = make([]int32, 0)
				}

				unlocks = append(unlocks, questCfg.Id)
				conditionUnlocks[conditionType] = unlocks
			}
		}

		if questCfg.Module == static.TaskModuleAchievement {
			achievementTotalCount++
		}

		return nil
	}

	for _, questCSV := range config.CfgTaskDataConfig.GetAllData() {
		err := transferQuestCSV(questCSV)
		if err != nil {
			return err
		}
	}

	for _, questCSV := range config.CfgHeroTaskDataConfig.GetAllData() {
		err := transferQuestCSV(questCSV)
		if err != nil {
			return err
		}
	}

	for _, questCSV := range config.CfgScorePassTaskConfig.GetAllData() {
		err := transferQuestCSV(questCSV)
		if err != nil {
			return err
		}
	}

	for questId, questIds := range questUnlocks {
		questCfg := quests[questId]
		questCfg.NextQuests = []int32{}
		questCfg.NextQuests = append(questCfg.NextQuests, questIds...)
	}

	for _, questAcvtiviyCSV := range config.CfgTaskRewardConfig.GetAllData() {
		questActivityCfg := &QuestActivity{}
		err := transfer.Transfer(questAcvtiviyCSV, questActivityCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}
		questActivity[questAcvtiviyCSV.Id] = questActivityCfg
	}

	qe.Quests = quests
	qe.QuestActivity = questActivity
	qe.AccountLevelUnlocks = accountLevelUnlocks
	qe.PassLevelUnlocks = passLevelUnlocks
	qe.CommonConditionUnlocks = conditionUnlocks
	qe.ModuleToQuests = moduleToQuests
	qe.achievementTotalCount = achievementTotalCount
	qe.Groups = groups

	return nil
}

func (qe *QuestEntry) Get(id int32) (*Quest, error) {
	cfg, ok := qe.Quests[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrQuestConfigNotFound, id)
	}

	return cfg, nil
}

func (qe *QuestEntry) GetQuestActivity(id int32) (*QuestActivity, error) {
	cfg, ok := qe.QuestActivity[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrQuestActivityConfigNotFound, id)
	}

	return cfg, nil
}

func (qe *QuestEntry) GetAll() map[int32]*Quest {
	qe.RLock()
	defer qe.RUnlock()

	return qe.Quests
}

func (qe *QuestEntry) IsAutoRecieve(id int32) bool {
	cfg, ok := qe.Quests[id]
	if !ok {
		return false
	}

	return cfg.IsAutoReceive
}

func (qe *QuestEntry) GetAccountLevelUnlocks(level int32) []int32 {
	questIds, ok := qe.AccountLevelUnlocks[level]
	if !ok {
		return []int32{}
	}

	return questIds
}

func (qe *QuestEntry) GetPassLevelUnlocks(levelId int32) []int32 {
	questIds, ok := qe.PassLevelUnlocks[levelId]
	if !ok {
		return []int32{}
	}

	return questIds
}

func (qe *QuestEntry) GetCommonConditionUnlocks(conditionType int32) []int32 {
	questIds, ok := qe.CommonConditionUnlocks[conditionType]
	if !ok {
		return []int32{}
	}

	return questIds
}

func (qe *QuestEntry) GetQuestIdsByModule(module int32) []int32 {
	questIds, ok := qe.ModuleToQuests[module]
	if !ok {
		return []int32{}
	}

	return questIds
}

func (qe *QuestEntry) GetAchievementTotalCount() int32 {
	qe.RLock()
	defer qe.RUnlock()
	return qe.achievementTotalCount
}

func (qe *QuestEntry) GetQuestIdsByGroup(groupId int32) []int32 {
	questIds, ok := qe.Groups[groupId]
	if !ok {
		return []int32{}
	}

	return questIds
}
