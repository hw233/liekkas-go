// ===================================== //
// author:       gavingqf                //
// == Please don'g change me by hand ==  //
//====================================== //

/*you have defined the following interface:
type IConfig interface {
	// load interface
	Load(path string) bool

	// clear interface
	Clear()
}
*/

package base

import (
	"shared/utility/glog"
	"strings"
)

type CfgExploreChapterLevel struct {
	Id                      int32
	ChapterId               int32
	Type                    int32
	UnlockConditions        []string
	TeamExp                 int32
	CharacterExp            int32
	PowerPay                int32
	BattleID                int32
	ShowLevelBoss           int32
	Achievement1ID          int32
	Achievement1DropId      []int32
	Achievement2ID          int32
	Achievement2DropId      []int32
	Target1ID               int32
	Target2ID               int32
	Target3ID               int32
	TargetDrops             []int32
	FirstDrop               []int32
	NormalDrop              []int32
	SweepDrop               []int32
	YggdrasilDrop           []int32
	NextLevel               int32
	ChapterStory            int32
	BattleStory             int32
	WinStory                int32
	LoseStory               int32
	ExploreStory            int32
	FormationLimitCharacter int32
	FormationLimitCamp      []int32
	FormationLimitCarrer    []int32
	FormationLimitRemote    int32
	FormationLimitSex       []int32
	ChallengeTimes          int32
	AdaptionId              int32
	RefreshType             int32
}

type CfgExploreChapterLevelConfig struct {
	data map[int32]*CfgExploreChapterLevel
}

func NewCfgExploreChapterLevelConfig() *CfgExploreChapterLevelConfig {
	return &CfgExploreChapterLevelConfig{
		data: make(map[int32]*CfgExploreChapterLevel),
	}
}

func (c *CfgExploreChapterLevelConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgExploreChapterLevel)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgExploreChapterLevel.Id field error,value:", vId)
			return false
		}

		/* parse ChapterId field */
		vChapterId, _ := parse.GetFieldByName(uint32(i), "chapterId")
		var ChapterIdRet bool
		data.ChapterId, ChapterIdRet = String2Int32(vChapterId)
		if !ChapterIdRet {
			glog.Error("Parse CfgExploreChapterLevel.ChapterId field error,value:", vChapterId)
			return false
		}

		/* parse Type field */
		vType, _ := parse.GetFieldByName(uint32(i), "type")
		var TypeRet bool
		data.Type, TypeRet = String2Int32(vType)
		if !TypeRet {
			glog.Error("Parse CfgExploreChapterLevel.Type field error,value:", vType)
			return false
		}

		/* parse UnlockConditions field */
		vecUnlockConditions, _ := parse.GetFieldByName(uint32(i), "unlockConditions")
		arrayUnlockConditions := strings.Split(vecUnlockConditions, ",")
		for j := 0; j < len(arrayUnlockConditions); j++ {
			v := arrayUnlockConditions[j]
			data.UnlockConditions = append(data.UnlockConditions, v)
		}

		/* parse TeamExp field */
		vTeamExp, _ := parse.GetFieldByName(uint32(i), "teamExp")
		var TeamExpRet bool
		data.TeamExp, TeamExpRet = String2Int32(vTeamExp)
		if !TeamExpRet {
			glog.Error("Parse CfgExploreChapterLevel.TeamExp field error,value:", vTeamExp)
			return false
		}

		/* parse CharacterExp field */
		vCharacterExp, _ := parse.GetFieldByName(uint32(i), "characterExp")
		var CharacterExpRet bool
		data.CharacterExp, CharacterExpRet = String2Int32(vCharacterExp)
		if !CharacterExpRet {
			glog.Error("Parse CfgExploreChapterLevel.CharacterExp field error,value:", vCharacterExp)
			return false
		}

		/* parse PowerPay field */
		vPowerPay, _ := parse.GetFieldByName(uint32(i), "powerPay")
		var PowerPayRet bool
		data.PowerPay, PowerPayRet = String2Int32(vPowerPay)
		if !PowerPayRet {
			glog.Error("Parse CfgExploreChapterLevel.PowerPay field error,value:", vPowerPay)
			return false
		}

		/* parse BattleID field */
		vBattleID, _ := parse.GetFieldByName(uint32(i), "battleID")
		var BattleIDRet bool
		data.BattleID, BattleIDRet = String2Int32(vBattleID)
		if !BattleIDRet {
			glog.Error("Parse CfgExploreChapterLevel.BattleID field error,value:", vBattleID)
			return false
		}

		/* parse ShowLevelBoss field */
		vShowLevelBoss, _ := parse.GetFieldByName(uint32(i), "showLevelBoss")
		var ShowLevelBossRet bool
		data.ShowLevelBoss, ShowLevelBossRet = String2Int32(vShowLevelBoss)
		if !ShowLevelBossRet {
			glog.Error("Parse CfgExploreChapterLevel.ShowLevelBoss field error,value:", vShowLevelBoss)
			return false
		}

		/* parse Achievement1ID field */
		vAchievement1ID, _ := parse.GetFieldByName(uint32(i), "achievement1ID")
		var Achievement1IDRet bool
		data.Achievement1ID, Achievement1IDRet = String2Int32(vAchievement1ID)
		if !Achievement1IDRet {
			glog.Error("Parse CfgExploreChapterLevel.Achievement1ID field error,value:", vAchievement1ID)
			return false
		}

		/* parse Achievement1DropId field */
		vecAchievement1DropId, _ := parse.GetFieldByName(uint32(i), "achievement1DropId")
		if vecAchievement1DropId != "" {
			arrayAchievement1DropId := strings.Split(vecAchievement1DropId, ",")
			for j := 0; j < len(arrayAchievement1DropId); j++ {
				v, ret := String2Int32(arrayAchievement1DropId[j])
				if !ret {
					glog.Error("Parse CfgExploreChapterLevel.Achievement1DropId field error, value:", arrayAchievement1DropId[j])
					return false
				}
				data.Achievement1DropId = append(data.Achievement1DropId, v)
			}
		}

		/* parse Achievement2ID field */
		vAchievement2ID, _ := parse.GetFieldByName(uint32(i), "achievement2ID")
		var Achievement2IDRet bool
		data.Achievement2ID, Achievement2IDRet = String2Int32(vAchievement2ID)
		if !Achievement2IDRet {
			glog.Error("Parse CfgExploreChapterLevel.Achievement2ID field error,value:", vAchievement2ID)
			return false
		}

		/* parse Achievement2DropId field */
		vecAchievement2DropId, _ := parse.GetFieldByName(uint32(i), "achievement2DropId")
		if vecAchievement2DropId != "" {
			arrayAchievement2DropId := strings.Split(vecAchievement2DropId, ",")
			for j := 0; j < len(arrayAchievement2DropId); j++ {
				v, ret := String2Int32(arrayAchievement2DropId[j])
				if !ret {
					glog.Error("Parse CfgExploreChapterLevel.Achievement2DropId field error, value:", arrayAchievement2DropId[j])
					return false
				}
				data.Achievement2DropId = append(data.Achievement2DropId, v)
			}
		}

		/* parse Target1ID field */
		vTarget1ID, _ := parse.GetFieldByName(uint32(i), "target1ID")
		var Target1IDRet bool
		data.Target1ID, Target1IDRet = String2Int32(vTarget1ID)
		if !Target1IDRet {
			glog.Error("Parse CfgExploreChapterLevel.Target1ID field error,value:", vTarget1ID)
			return false
		}

		/* parse Target2ID field */
		vTarget2ID, _ := parse.GetFieldByName(uint32(i), "target2ID")
		var Target2IDRet bool
		data.Target2ID, Target2IDRet = String2Int32(vTarget2ID)
		if !Target2IDRet {
			glog.Error("Parse CfgExploreChapterLevel.Target2ID field error,value:", vTarget2ID)
			return false
		}

		/* parse Target3ID field */
		vTarget3ID, _ := parse.GetFieldByName(uint32(i), "target3ID")
		var Target3IDRet bool
		data.Target3ID, Target3IDRet = String2Int32(vTarget3ID)
		if !Target3IDRet {
			glog.Error("Parse CfgExploreChapterLevel.Target3ID field error,value:", vTarget3ID)
			return false
		}

		/* parse TargetDrops field */
		vecTargetDrops, _ := parse.GetFieldByName(uint32(i), "targetDrops")
		if vecTargetDrops != "" {
			arrayTargetDrops := strings.Split(vecTargetDrops, ",")
			for j := 0; j < len(arrayTargetDrops); j++ {
				v, ret := String2Int32(arrayTargetDrops[j])
				if !ret {
					glog.Error("Parse CfgExploreChapterLevel.TargetDrops field error, value:", arrayTargetDrops[j])
					return false
				}
				data.TargetDrops = append(data.TargetDrops, v)
			}
		}

		/* parse FirstDrop field */
		vecFirstDrop, _ := parse.GetFieldByName(uint32(i), "firstDrop")
		if vecFirstDrop != "" {
			arrayFirstDrop := strings.Split(vecFirstDrop, ",")
			for j := 0; j < len(arrayFirstDrop); j++ {
				v, ret := String2Int32(arrayFirstDrop[j])
				if !ret {
					glog.Error("Parse CfgExploreChapterLevel.FirstDrop field error, value:", arrayFirstDrop[j])
					return false
				}
				data.FirstDrop = append(data.FirstDrop, v)
			}
		}

		/* parse NormalDrop field */
		vecNormalDrop, _ := parse.GetFieldByName(uint32(i), "normalDrop")
		if vecNormalDrop != "" {
			arrayNormalDrop := strings.Split(vecNormalDrop, ",")
			for j := 0; j < len(arrayNormalDrop); j++ {
				v, ret := String2Int32(arrayNormalDrop[j])
				if !ret {
					glog.Error("Parse CfgExploreChapterLevel.NormalDrop field error, value:", arrayNormalDrop[j])
					return false
				}
				data.NormalDrop = append(data.NormalDrop, v)
			}
		}

		/* parse SweepDrop field */
		vecSweepDrop, _ := parse.GetFieldByName(uint32(i), "sweepDrop")
		if vecSweepDrop != "" {
			arraySweepDrop := strings.Split(vecSweepDrop, ",")
			for j := 0; j < len(arraySweepDrop); j++ {
				v, ret := String2Int32(arraySweepDrop[j])
				if !ret {
					glog.Error("Parse CfgExploreChapterLevel.SweepDrop field error, value:", arraySweepDrop[j])
					return false
				}
				data.SweepDrop = append(data.SweepDrop, v)
			}
		}

		/* parse YggdrasilDrop field */
		vecYggdrasilDrop, _ := parse.GetFieldByName(uint32(i), "yggdrasilDrop")
		if vecYggdrasilDrop != "" {
			arrayYggdrasilDrop := strings.Split(vecYggdrasilDrop, ",")
			for j := 0; j < len(arrayYggdrasilDrop); j++ {
				v, ret := String2Int32(arrayYggdrasilDrop[j])
				if !ret {
					glog.Error("Parse CfgExploreChapterLevel.YggdrasilDrop field error, value:", arrayYggdrasilDrop[j])
					return false
				}
				data.YggdrasilDrop = append(data.YggdrasilDrop, v)
			}
		}

		/* parse NextLevel field */
		vNextLevel, _ := parse.GetFieldByName(uint32(i), "nextLevel")
		var NextLevelRet bool
		data.NextLevel, NextLevelRet = String2Int32(vNextLevel)
		if !NextLevelRet {
			glog.Error("Parse CfgExploreChapterLevel.NextLevel field error,value:", vNextLevel)
			return false
		}

		/* parse ChapterStory field */
		vChapterStory, _ := parse.GetFieldByName(uint32(i), "chapterStory")
		var ChapterStoryRet bool
		data.ChapterStory, ChapterStoryRet = String2Int32(vChapterStory)
		if !ChapterStoryRet {
			glog.Error("Parse CfgExploreChapterLevel.ChapterStory field error,value:", vChapterStory)
			return false
		}

		/* parse BattleStory field */
		vBattleStory, _ := parse.GetFieldByName(uint32(i), "battleStory")
		var BattleStoryRet bool
		data.BattleStory, BattleStoryRet = String2Int32(vBattleStory)
		if !BattleStoryRet {
			glog.Error("Parse CfgExploreChapterLevel.BattleStory field error,value:", vBattleStory)
			return false
		}

		/* parse WinStory field */
		vWinStory, _ := parse.GetFieldByName(uint32(i), "winStory")
		var WinStoryRet bool
		data.WinStory, WinStoryRet = String2Int32(vWinStory)
		if !WinStoryRet {
			glog.Error("Parse CfgExploreChapterLevel.WinStory field error,value:", vWinStory)
			return false
		}

		/* parse LoseStory field */
		vLoseStory, _ := parse.GetFieldByName(uint32(i), "loseStory")
		var LoseStoryRet bool
		data.LoseStory, LoseStoryRet = String2Int32(vLoseStory)
		if !LoseStoryRet {
			glog.Error("Parse CfgExploreChapterLevel.LoseStory field error,value:", vLoseStory)
			return false
		}

		/* parse ExploreStory field */
		vExploreStory, _ := parse.GetFieldByName(uint32(i), "exploreStory")
		var ExploreStoryRet bool
		data.ExploreStory, ExploreStoryRet = String2Int32(vExploreStory)
		if !ExploreStoryRet {
			glog.Error("Parse CfgExploreChapterLevel.ExploreStory field error,value:", vExploreStory)
			return false
		}

		/* parse FormationLimitCharacter field */
		vFormationLimitCharacter, _ := parse.GetFieldByName(uint32(i), "formationLimitCharacter")
		var FormationLimitCharacterRet bool
		data.FormationLimitCharacter, FormationLimitCharacterRet = String2Int32(vFormationLimitCharacter)
		if !FormationLimitCharacterRet {
			glog.Error("Parse CfgExploreChapterLevel.FormationLimitCharacter field error,value:", vFormationLimitCharacter)
			return false
		}

		/* parse FormationLimitCamp field */
		vecFormationLimitCamp, _ := parse.GetFieldByName(uint32(i), "formationLimitCamp")
		if vecFormationLimitCamp != "" {
			arrayFormationLimitCamp := strings.Split(vecFormationLimitCamp, ",")
			for j := 0; j < len(arrayFormationLimitCamp); j++ {
				v, ret := String2Int32(arrayFormationLimitCamp[j])
				if !ret {
					glog.Error("Parse CfgExploreChapterLevel.FormationLimitCamp field error, value:", arrayFormationLimitCamp[j])
					return false
				}
				data.FormationLimitCamp = append(data.FormationLimitCamp, v)
			}
		}

		/* parse FormationLimitCarrer field */
		vecFormationLimitCarrer, _ := parse.GetFieldByName(uint32(i), "formationLimitCarrer")
		if vecFormationLimitCarrer != "" {
			arrayFormationLimitCarrer := strings.Split(vecFormationLimitCarrer, ",")
			for j := 0; j < len(arrayFormationLimitCarrer); j++ {
				v, ret := String2Int32(arrayFormationLimitCarrer[j])
				if !ret {
					glog.Error("Parse CfgExploreChapterLevel.FormationLimitCarrer field error, value:", arrayFormationLimitCarrer[j])
					return false
				}
				data.FormationLimitCarrer = append(data.FormationLimitCarrer, v)
			}
		}

		/* parse FormationLimitRemote field */
		vFormationLimitRemote, _ := parse.GetFieldByName(uint32(i), "formationLimitRemote")
		var FormationLimitRemoteRet bool
		data.FormationLimitRemote, FormationLimitRemoteRet = String2Int32(vFormationLimitRemote)
		if !FormationLimitRemoteRet {
			glog.Error("Parse CfgExploreChapterLevel.FormationLimitRemote field error,value:", vFormationLimitRemote)
			return false
		}

		/* parse FormationLimitSex field */
		vecFormationLimitSex, _ := parse.GetFieldByName(uint32(i), "formationLimitSex")
		if vecFormationLimitSex != "" {
			arrayFormationLimitSex := strings.Split(vecFormationLimitSex, ",")
			for j := 0; j < len(arrayFormationLimitSex); j++ {
				v, ret := String2Int32(arrayFormationLimitSex[j])
				if !ret {
					glog.Error("Parse CfgExploreChapterLevel.FormationLimitSex field error, value:", arrayFormationLimitSex[j])
					return false
				}
				data.FormationLimitSex = append(data.FormationLimitSex, v)
			}
		}

		/* parse ChallengeTimes field */
		vChallengeTimes, _ := parse.GetFieldByName(uint32(i), "challengeTimes")
		var ChallengeTimesRet bool
		data.ChallengeTimes, ChallengeTimesRet = String2Int32(vChallengeTimes)
		if !ChallengeTimesRet {
			glog.Error("Parse CfgExploreChapterLevel.ChallengeTimes field error,value:", vChallengeTimes)
			return false
		}

		/* parse AdaptionId field */
		vAdaptionId, _ := parse.GetFieldByName(uint32(i), "adaptionId")
		var AdaptionIdRet bool
		data.AdaptionId, AdaptionIdRet = String2Int32(vAdaptionId)
		if !AdaptionIdRet {
			glog.Error("Parse CfgExploreChapterLevel.AdaptionId field error,value:", vAdaptionId)
			return false
		}

		/* parse RefreshType field */
		vRefreshType, _ := parse.GetFieldByName(uint32(i), "refreshType")
		var RefreshTypeRet bool
		data.RefreshType, RefreshTypeRet = String2Int32(vRefreshType)
		if !RefreshTypeRet {
			glog.Error("Parse CfgExploreChapterLevel.RefreshType field error,value:", vRefreshType)
			return false
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgExploreChapterLevelConfig) Clear() {
}

func (c *CfgExploreChapterLevelConfig) Find(id int32) (*CfgExploreChapterLevel, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgExploreChapterLevelConfig) GetAllData() map[int32]*CfgExploreChapterLevel {
	return c.data
}

func (c *CfgExploreChapterLevelConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.ChapterId, ",", v.Type, ",", v.UnlockConditions, ",", v.TeamExp, ",", v.CharacterExp, ",", v.PowerPay, ",", v.BattleID, ",", v.ShowLevelBoss, ",", v.Achievement1ID, ",", v.Achievement1DropId, ",", v.Achievement2ID, ",", v.Achievement2DropId, ",", v.Target1ID, ",", v.Target2ID, ",", v.Target3ID, ",", v.TargetDrops, ",", v.FirstDrop, ",", v.NormalDrop, ",", v.SweepDrop, ",", v.YggdrasilDrop, ",", v.NextLevel, ",", v.ChapterStory, ",", v.BattleStory, ",", v.WinStory, ",", v.LoseStory, ",", v.ExploreStory, ",", v.FormationLimitCharacter, ",", v.FormationLimitCamp, ",", v.FormationLimitCarrer, ",", v.FormationLimitRemote, ",", v.FormationLimitSex, ",", v.ChallengeTimes, ",", v.AdaptionId, ",", v.RefreshType)
	}
}
