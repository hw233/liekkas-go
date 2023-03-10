package entry

import (
	"shared/csv/static"
	"sort"
	"sync"

	"shared/common"
	"shared/csv/base"
	"shared/utility/errors"
	"shared/utility/transfer"
)

const (
	CfgCharacter      = "cfg_character"
	CfgCharacterStage = "cfg_character_stage"
	CfgCharacterStar  = "cfg_character_star"
	CfgCharacterSkill = "cfg_character_skill"
	CfgCharacterLevel = "cfg_character_level"

	CharacterStagesCount = 11 // 角色升阶次数
	CharacterStarsCount  = 5  // 角色升星次数

	CharacterSkillsLevelUpCount = 5 // 角色技能升级次数

	CharacterTypeRole = 1 // 角色类型：可以拥有的角色

	CharacterSkillInitLevel = 1 // 技能解锁时初始等级

	CharacterSkillUnlockDefault = 0 // 默认解锁
	CharacterSkillUnlockLevel   = 1 // 角色等级
	CharacterSkillUnlockStage   = 2 // 角色进阶等级
	CharacterSkillUnlockStar    = 3 // 角色星级

)

type Character struct {
	sync.RWMutex

	SkillLevels      map[int32]CharacterSkillLevel
	associatedSkills map[int32]int32
	exps             []int32
	characters       map[int32]CharacterInfo
	allSkills        map[int32]map[int32]*CharacterSkill
}

type CharacterInfo struct {
	Stages []CharacterStage         `ignore:"true"`
	Stars  []CharacterStar          `ignore:"true"`
	Skills map[int32]CharacterSkill `ignore:"true"`

	Career int32 // 职业
	Camp   int32 // 阵营
	Rarity int32
	Avatar int32 `ignore:"true"` // 头像
}

type CharacterLevel struct {
	Exp int32 `src:"RiseExp"`
}

type CharacterStage struct {
	LevelLimit    int32
	Frame         int32
	HpMax         int32
	PhyAtk        int32
	MagAtk        int32
	PhyDfs        int32
	MagDfs        int32
	CritAtkRatio  int32
	CritDfsRatio  int32
	CritAtkValue  int32
	CritDfsValue  int32
	HitRateValue  int32
	EvadeValue    int32
	NormalAtkUp   int32
	NormalDfsUp   int32
	SkillPhyAtkUp int32
	SkillPhyDfsUp int32
	SkillMagAtkUp int32
	SkillMagDfsUp int32
	Cost          *common.Rewards `rule:"rewards"`
}

type CharacterStar struct {
	HpRatio     float64
	PhyAtkRatio float64
	MagAtkRatio float64
	PhyDfsRatio float64
	MagDfsRatio float64
	Cost        *common.Rewards `rule:"rewards"`
	RarityUp    int32
}

type CharacterSkillLevel struct {
	Cost            *common.Rewards `rule:"rewards"`
	SkillUnlock     int32
	UnlockParam     int32
	StarAutoLevelUp int32
	CombatPower     []int32
}
type CharacterSkill struct {
	AssociatedSkillNumber int32 // 有关联技能的技能不可手动升级，会随着关联技能升级而自动升级
	SkillType             int32
}

func NewCharacter() *Character {
	return &Character{
		characters: map[int32]CharacterInfo{},
	}
}

func (c *Character) Reload(config *Config) error {
	c.Lock()
	defer c.Unlock()

	associatedSkills := map[int32]int32{}
	var exps []int32
	characters := map[int32]CharacterInfo{}
	stages := map[int32][]CharacterStage{}
	stars := map[int32][]CharacterStar{}
	skills := map[int32]map[int32]*CharacterSkill{}
	SkillLevels := map[int32]CharacterSkillLevel{}
	allSkills := map[int32]map[int32]*CharacterSkill{}
	// 初始化角色升阶数据
	for _, stage := range config.CfgCharacterStageConfig.GetAllData() {
		if _, ok := stages[stage.CharId]; !ok {
			stages[stage.CharId] = make([]CharacterStage, CharacterStagesCount, CharacterStagesCount)
		}

		stageIndex := int(stage.Stage)

		if stageIndex >= len(stages[stage.CharId]) {
			return errors.WrapTrace(errors.Swrapf(common.ErrCSVFormatInvalid, CfgCharacter, stage))
		}

		err := transfer.Transfer(stage, &stages[stage.CharId][stageIndex])
		if err != nil {
			return errors.WrapTrace(err)
		}
	}

	// 初始化角色升星数据
	for _, star := range config.CfgCharacterStarConfig.GetAllData() {
		if _, ok := stars[star.CharID]; !ok {
			stars[star.CharID] = make([]CharacterStar, CharacterStarsCount, CharacterStarsCount)
		}

		starIndex := int(star.Star)
		if starIndex > len(stars[star.CharID]) {
			return errors.WrapTrace(errors.Swrapf(common.ErrCSVFormatInvalid, CfgCharacter, star))
		}

		err := transfer.Transfer(star, &stars[star.CharID][starIndex-1])
		if err != nil {
			return errors.WrapTrace(err)
		}
	}

	findSkillBySkillNum := func(all map[int32]*base.CfgCharacterSkill, skillNum int32) bool {
		for _, skill := range all {
			if skill.SkillNumber == skillNum {
				return true
			}
		}
		return false
	}
	// 初始化角色技能数据
	for _, skill := range config.CfgCharacterSkillConfig.GetAllData() {
		if skill.SkillLevel != 0 {
			// 非怪物技能
			if _, ok := skills[skill.RoleID]; !ok {
				skills[skill.RoleID] = map[int32]*CharacterSkill{}

			}

			if _, ok := skills[skill.RoleID][skill.SkillNumber]; !ok {
				skills[skill.RoleID][skill.SkillNumber] = &CharacterSkill{}
			}
			if skill.SkillLevel == 1 {

				if skill.AssociatedSkillNumber > 0 {
					associatedSkills[skill.AssociatedSkillNumber] = skill.SkillNumber
					if !findSkillBySkillNum(config.CfgCharacterSkillConfig.GetAllData(), skill.AssociatedSkillNumber) {
						return errors.New("AssociatedSkillNumber not found ,character_id:%d,skillNum:%d", skill.RoleID, skill.AssociatedSkillNumber)
					}
				}
				// 1级时候读取 关联技能
				err := transfer.Transfer(skill, skills[skill.RoleID][skill.SkillNumber])
				if err != nil {
					return errors.WrapTrace(err)
				}
			}
			skillLevel := &CharacterSkillLevel{}
			//	读取 各等级升级消耗等数据
			err := transfer.Transfer(skill, skillLevel)
			if err != nil {
				return errors.WrapTrace(err)
			}
			SkillLevels[skill.Id] = *skillLevel
		}

		charaId := skill.RoleID
		charaSkills, ok := allSkills[charaId]
		if !ok {
			charaSkills = map[int32]*CharacterSkill{}
			allSkills[charaId] = charaSkills
		}
		skillNum := skill.SkillNumber
		_, ok = charaSkills[skillNum]
		if !ok {
			charaSkills[skillNum] = &CharacterSkill{
				SkillType:             skill.SkillType,
				AssociatedSkillNumber: skill.AssociatedSkillNumber,
			}
		}
	}

	charaAvatar := map[int32]int32{}
	for _, data := range config.CfgItemDataConfig.GetAllData() {
		if data.ItemType == static.ItemTypeCharacterAvatarAuto {
			charaAvatar[data.UseParam[0]] = data.Id
		}
	}

	// 查询玩家数据
	for _, character := range config.CfgCharacterConfig.GetAllData() {
		// 只初始化 图鉴可见 并且 类型是角色 的数据
		if character.Visible && character.Type == CharacterTypeRole {
			stage, ok := stages[character.Id]
			if !ok {
				return errors.WrapTrace(errors.New(" CfgCharacterStage not found,characterId:%d", character.Id))
			}

			star, ok := stars[character.Id]
			if !ok {
				return errors.WrapTrace(errors.New(" CfgCharacterStar not found,characterId:%d", character.Id))
			}

			skill, ok := skills[character.Id]
			if !ok {
				return errors.WrapTrace(errors.New(" CfgCharacterSkill not found,characterId:%d", character.Id))
			}

			Skills := map[int32]CharacterSkill{}
			for k, v := range skill {
				Skills[k] = *v
			}
			avatar, ok := charaAvatar[character.Id]
			if !ok {
				//todo:策划暂时没空配
				//return errors.New(fmt.Sprintf("charaAvatar not found,characterId:%d", character.Id))
			}
			characterInfo := &CharacterInfo{
				Stages: stage,
				Stars:  star,
				Skills: Skills,
				Avatar: avatar,
			}

			err := transfer.Transfer(character, characterInfo)
			if err != nil {
				return errors.WrapTrace(err)
			}
			characters[character.Id] = *characterInfo
		}
	}

	// 侍从的升级所需经验(没有稀有度的区分)

	var cfgs []*base.CfgCharacterLevel
	for _, v := range config.CfgCharacterLevelConfig.GetAllData() {
		cfgs = append(cfgs, v)
	}
	less := func(i, j int) bool {
		return cfgs[i].Id < cfgs[j].Id
	}
	sort.Slice(cfgs, less)
	var characterLevels []CharacterLevel

	for _, level := range cfgs {
		characterLevel := &CharacterLevel{}

		err := transfer.Transfer(level, characterLevel)
		if err != nil {
			return errors.WrapTrace(err)
		}
		characterLevels = append(characterLevels, *characterLevel)
	}
	for _, characterLevel := range characterLevels {
		exps = append(exps, characterLevel.Exp)
	}

	c.associatedSkills = associatedSkills
	c.exps = exps
	c.characters = characters
	c.SkillLevels = SkillLevels
	c.allSkills = allSkills
	return nil
}

func (c *Character) Check(config *Config) error {

	for _, stage := range config.CfgCharacterStageConfig.GetAllData() {
		_, ok := config.CfgItemDataConfig.Find(stage.CharId)
		if !ok {
			return errors.Swrapf(common.ErrNotFoundInAssociatedCSV, CfgCharacterStage, stage.Id, CfgCharacter, stage.CharId)
		}
	}

	for _, star := range config.CfgCharacterStarConfig.GetAllData() {
		_, ok := config.CfgItemDataConfig.Find(star.CharID)
		if !ok {
			return errors.Swrapf(common.ErrNotFoundInAssociatedCSV, CfgCharacterStar, star.Id, CfgCharacter, star.CharID)
		}
	}

	var cfgs []*base.CfgCharacterLevel
	for _, v := range config.CfgCharacterLevelConfig.GetAllData() {
		cfgs = append(cfgs, v)
	}
	less := func(i, j int) bool {
		return cfgs[i].Id < cfgs[j].Id
	}
	sort.Slice(cfgs, less)

	for i, level := range cfgs {

		if int32(i+1) != level.Id {
			return errors.Swrapf(common.ErrNotFoundInCSV, CfgCharacterLevel, level.Id)
		}

	}

	return nil
}

func (c *Character) Star(id, star int32) (*CharacterStar, error) {
	c.RLock()
	defer c.RUnlock()
	charaCfg, ok := c.characters[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgCharacter, id)
	}

	if int(star) > len(charaCfg.Stars) {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgCharacter, id)
	}

	return &charaCfg.Stars[int(star)-1], nil
}

func (c *Character) Stage(id, stage int32) (*CharacterStage, error) {
	c.RLock()
	defer c.RUnlock()
	charaCfg, ok := c.characters[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgCharacter, id)
	}

	if int(stage) > len(charaCfg.Stages) {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgCharacter, id)
	}

	return &charaCfg.Stages[int(stage)], nil
}

func (c *Character) AllSkill(id int32) (map[int32]CharacterSkill, error) {
	c.RLock()
	defer c.RUnlock()
	info, ok := c.characters[id]
	if !ok {
		return nil, errors.WrapTrace(common.ErrNotFoundInCSV)
	}

	return info.Skills, nil
}

func (c *Character) GetSkillByNum(charaId, skillNum int32) (*CharacterSkill, error) {
	c.RLock()
	defer c.RUnlock()

	skills, err := c.AllSkill(charaId)
	if err != nil {
		return nil, err
	}

	skill, ok := skills[skillNum]
	if !ok {
		return nil, errors.WrapTrace(common.ErrNotFoundInCSV)
	}

	return &skill, nil
}

func (c *Character) FindCanUpgradeSkill(id, skillNum int32) (*CharacterSkill, error) {
	c.RLock()
	defer c.RUnlock()
	info, ok := c.characters[id]
	if !ok {
		return nil, errors.WrapTrace(common.ErrCharacterSkillCannotLevelUp)
	}

	skill, ok := info.Skills[skillNum]
	if !ok {
		return nil, errors.WrapTrace(common.ErrCharacterSkillCannotLevelUp)
	}
	// 有关联技能的技能不可手动升级，会随着关联技能升级而自动升级
	_, ok = c.associatedSkills[skillNum]
	if ok {
		return nil, errors.WrapTrace(common.ErrCharacterSkillCannotLevelUp)
	}

	return &skill, nil
}

func (c *Character) SkillLevel(skillNum, level int32) (*CharacterSkillLevel, error) {
	c.RLock()
	defer c.RUnlock()
	// SkillId 拼接
	skillId := skillNum*10 + level
	skillLevel, ok := c.SkillLevels[skillId]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgCharacterSkill, skillId)
	}
	return &skillLevel, nil
}

func (c *Character) GetSkillsAssociated(characterId int32, skillNum int32) (int32, bool) {
	c.RLock()
	defer c.RUnlock()
	// 有关联技能的技能不可手动升级，会随着关联技能升级而自动升级

	characterInfo, ok := c.characters[characterId]
	if !ok {
		return 0, false
	}
	skill, ok := characterInfo.Skills[skillNum]
	if !ok {
		return 0, false
	}
	if skill.AssociatedSkillNumber == 0 {
		return 0, false
	}
	return skill.AssociatedSkillNumber, true
}

func (c *Character) GetAllCharacterInfo() map[int32]CharacterInfo {
	c.RLock()
	defer c.RUnlock()
	return c.characters
}

func (c *Character) GetNeedExp(nowExp, to int32) int32 {
	if int(to) > len(c.exps) {
		to = int32(len(c.exps))
	}
	return c.exps[to-1] - nowExp

}

func (c *Character) GetExpArr() []int32 {
	c.RLock()
	defer c.RUnlock()
	return c.exps
}

func (c *Character) GetRare(charaId int32) (int32, error) {
	c.RLock()
	defer c.RUnlock()

	charaCfg, ok := c.characters[charaId]
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgCharacter, charaId)
	}

	return charaCfg.Rarity, nil
}

func (c *Character) Career(id int32) (int32, error) {
	c.RLock()
	defer c.RUnlock()
	charaCfg, ok := c.characters[id]
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgCharacter, id)
	}

	return charaCfg.Career, nil
}

func (c *Character) Camp(id int32) (int32, error) {
	c.RLock()
	defer c.RUnlock()
	charaCfg, ok := c.characters[id]
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgCharacter, id)
	}

	return charaCfg.Camp, nil
}

func (c *Character) Avatar(id int32) (int32, error) {
	c.RLock()
	defer c.RUnlock()
	charaCfg, ok := c.characters[id]
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgCharacter, id)
	}

	return charaCfg.Avatar, nil
}

func (c *CharacterSkillLevel) CanSkillUnlock(level, star, stage int32) bool {
	switch c.SkillUnlock {
	case CharacterSkillUnlockDefault:
		return true
	case CharacterSkillUnlockLevel:
		return level >= c.UnlockParam
	case CharacterSkillUnlockStar:
		return star >= c.UnlockParam
	case CharacterSkillUnlockStage:
		return stage >= c.UnlockParam
	}
	return false
}

func (c *CharacterSkillLevel) CanSkillAutoLvUp(star int32) bool {
	return c.StarAutoLevelUp > 0 && star >= c.StarAutoLevelUp
}

func (c *CharacterSkillLevel) CanSkillLevelUp() bool {
	return c.StarAutoLevelUp == 0
}

func (c *Character) GetSkillType(charaId, skillNum int32) (int32, error) {
	charaSkills, ok := c.allSkills[charaId]
	if !ok {
		return 0, errors.WrapTrace(common.ErrNotFoundInCSV)
	}

	skill, ok := charaSkills[skillNum]
	if !ok {
		return 0, errors.WrapTrace(common.ErrNotFoundInCSV)
	}

	return skill.SkillType, nil
}
