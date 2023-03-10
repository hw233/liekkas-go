package entry

import (
	"shared/common"
	"shared/csv/static"
	"shared/utility/errors"
	"shared/utility/transfer"
	"sync"
)

type HeroSkill struct {
	Id         int32
	SkillID    int32
	SkillLevel int32
	CostItems  *common.Rewards    `rule:"rewards"`
	Unlock     *common.Conditions `rule:"conditions"`
}

type Hero struct {
	Id               int32
	Show             bool
	Skills           map[int32][]*HeroSkill `ignore:"true"`
	SkillItem        int32
	MainAttr         int32
	UnlockCost       *common.Rewards      `rule:"rewards"`
	AttendantUnlocks []*common.Conditions `ignore:"true"`
}

type HeroLevel struct {
	Id  int32
	Exp int32
}

type HeroEntry struct {
	sync.RWMutex

	Heros      map[int32]*Hero
	HeroLevels map[int32]*HeroLevel

	HeroMaxLevel int32
}

func NewHeroEntry() *HeroEntry {
	return &HeroEntry{
		Heros:      map[int32]*Hero{},
		HeroLevels: map[int32]*HeroLevel{},
	}
}

func (he *HeroEntry) Check(config *Config) error {
	return nil
}

func (he *HeroEntry) Reload(config *Config) error {
	he.Lock()
	defer he.Unlock()

	allSkills := map[int32]map[int32]*HeroSkill{}

	for _, skillCSV := range config.CfgHeroSkillConfig.GetAllData() {
		skill := &HeroSkill{}

		err := transfer.Transfer(skillCSV, skill)
		if err != nil {
			return errors.WrapTrace(err)
		}

		skillId := skill.SkillID
		skills, ok := allSkills[skillId]
		if !ok {
			skills = map[int32]*HeroSkill{}
			allSkills[skillId] = skills
		}

		skills[skill.SkillLevel] = skill
	}

	for _, heroCSV := range config.CfgHeroConfig.GetAllData() {
		hero := &Hero{
			Skills:           map[int32][]*HeroSkill{},
			AttendantUnlocks: []*common.Conditions{},
		}

		err := transfer.Transfer(heroCSV, hero)
		if err != nil {
			return errors.WrapTrace(err)
		}

		for _, unlockLevel := range heroCSV.PledgeUnlockLv {
			conditions := common.NewConditions()
			conditions.AddCondition(common.NewCondition(static.ConditionTypeHeroLevel, unlockLevel))

			hero.AttendantUnlocks = append(hero.AttendantUnlocks, conditions)
		}

		for _, skillId := range heroCSV.Skills {
			skills, ok := allSkills[skillId]
			if !ok {
				return errors.WrapTrace(errors.Swrapf(common.ErrHeroSkillConfigNotFound, hero.Id, skillId))
			}

			maxLevel := len(skills)
			heroSkills := make([]*HeroSkill, 0, maxLevel)
			for level := 1; level <= maxLevel; level++ {
				skill, ok := skills[int32(level)]
				if !ok {
					return errors.WrapTrace(errors.Swrapf(common.ErrHeroSkillLevelNotFound, hero.Id, skillId, level))
				}

				heroSkills = append(heroSkills, skill)
			}

			hero.Skills[skillId] = heroSkills
		}

		he.Heros[hero.Id] = hero
	}

	for _, heroLevelCSV := range config.CfgHeroLevelConfig.GetAllData() {
		heroLevel := &HeroLevel{}
		err := transfer.Transfer(heroLevelCSV, heroLevel)
		if err != nil {
			return errors.WrapTrace(err)
		}

		he.HeroLevels[heroLevel.Id] = heroLevel
	}

	maxLevel := int32(len(he.HeroLevels))
	for level := int32(1); level <= maxLevel; level++ {
		_, ok := he.HeroLevels[level]
		if !ok {
			return errors.WrapTrace(errors.Swrapf(common.ErrHeroLevelConfigNotFound, level))
		}
	}

	he.HeroMaxLevel = maxLevel

	return nil
}

func (he *HeroEntry) GetHero(id int32) (*Hero, error) {
	hero, ok := he.Heros[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrHeroConfigNotFound, id)
	}

	return hero, nil
}

func (he *HeroEntry) GetLevel(level int32) (*HeroLevel, error) {
	cfg, ok := he.HeroLevels[level]

	if !ok {
		return nil, errors.Swrapf(common.ErrHeroLevelConfigNotFound, level)
	}

	return cfg, nil
}

func (he *HeroEntry) GetSkill(heroId, skillId, level int32) (*HeroSkill, error) {
	hero, err := he.GetHero(heroId)
	if err != nil {
		return nil, err
	}

	skills, ok := hero.Skills[skillId]
	if !ok {
		return nil, errors.Swrapf(common.ErrHeroSkillConfigNotFound, hero.Id, skillId)
	}

	if len(skills) < int(level) {
		return nil, errors.WrapTrace(errors.Swrapf(common.ErrHeroSkillLevelNotFound, heroId, skillId, level))
	}

	return skills[level-1], nil
}

func (he *HeroEntry) IsSkillLevelMax(heroId, skillId, level int32) (bool, error) {
	hero, err := he.GetHero(heroId)
	if err != nil {
		return false, err
	}

	skills, ok := hero.Skills[skillId]
	if !ok {
		return false, errors.Swrapf(common.ErrHeroSkillConfigNotFound, hero.Id, skillId)
	}

	return int(level) >= len(skills), nil
}

func (he *HeroEntry) GetAttendantUnlockConditions(heroId, slot int32) (*common.Conditions, error) {
	hero, err := he.GetHero(heroId)
	if err != nil {
		return nil, err
	}

	if slot < 0 || int(slot) >= len(hero.AttendantUnlocks) {
		return nil, errors.Swrapf(common.ErrHeroAttendantConfigNotFound, heroId, slot)
	}

	return hero.AttendantUnlocks[slot], nil
}
