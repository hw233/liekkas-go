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

type CfgTowerLevel struct {
	Id                   int32
	BattleId             int32
	IsAchievement        bool
	AchievementCondition []int32
	AchievementID        int32
	AchievementDrop      []int32
	LevelDrop            []int32
	HangUpDrop           []int32
	HangUpSeconds        int32
	LvPara               []string
}

type CfgTowerLevelConfig struct {
	data map[int32]*CfgTowerLevel
}

func NewCfgTowerLevelConfig() *CfgTowerLevelConfig {
	return &CfgTowerLevelConfig{
		data: make(map[int32]*CfgTowerLevel),
	}
}

func (c *CfgTowerLevelConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgTowerLevel)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgTowerLevel.Id field error,value:", vId)
			return false
		}

		/* parse BattleId field */
		vBattleId, _ := parse.GetFieldByName(uint32(i), "battleId")
		var BattleIdRet bool
		data.BattleId, BattleIdRet = String2Int32(vBattleId)
		if !BattleIdRet {
			glog.Error("Parse CfgTowerLevel.BattleId field error,value:", vBattleId)
			return false
		}

		/* parse IsAchievement field */
		vIsAchievement, _ := parse.GetFieldByName(uint32(i), "isAchievement")
		var IsAchievementRet bool
		data.IsAchievement, IsAchievementRet = String2Bool(vIsAchievement)
		if !IsAchievementRet {
			glog.Error("Parse CfgTowerLevel.IsAchievement field error,value:", vIsAchievement)
		}

		/* parse AchievementCondition field */
		vecAchievementCondition, _ := parse.GetFieldByName(uint32(i), "achievementCondition")
		if vecAchievementCondition != "" {
			arrayAchievementCondition := strings.Split(vecAchievementCondition, ",")
			for j := 0; j < len(arrayAchievementCondition); j++ {
				v, ret := String2Int32(arrayAchievementCondition[j])
				if !ret {
					glog.Error("Parse CfgTowerLevel.AchievementCondition field error, value:", arrayAchievementCondition[j])
					return false
				}
				data.AchievementCondition = append(data.AchievementCondition, v)
			}
		}

		/* parse AchievementID field */
		vAchievementID, _ := parse.GetFieldByName(uint32(i), "achievementID")
		var AchievementIDRet bool
		data.AchievementID, AchievementIDRet = String2Int32(vAchievementID)
		if !AchievementIDRet {
			glog.Error("Parse CfgTowerLevel.AchievementID field error,value:", vAchievementID)
			return false
		}

		/* parse AchievementDrop field */
		vecAchievementDrop, _ := parse.GetFieldByName(uint32(i), "achievementDrop")
		if vecAchievementDrop != "" {
			arrayAchievementDrop := strings.Split(vecAchievementDrop, ",")
			for j := 0; j < len(arrayAchievementDrop); j++ {
				v, ret := String2Int32(arrayAchievementDrop[j])
				if !ret {
					glog.Error("Parse CfgTowerLevel.AchievementDrop field error, value:", arrayAchievementDrop[j])
					return false
				}
				data.AchievementDrop = append(data.AchievementDrop, v)
			}
		}

		/* parse LevelDrop field */
		vecLevelDrop, _ := parse.GetFieldByName(uint32(i), "levelDrop")
		if vecLevelDrop != "" {
			arrayLevelDrop := strings.Split(vecLevelDrop, ",")
			for j := 0; j < len(arrayLevelDrop); j++ {
				v, ret := String2Int32(arrayLevelDrop[j])
				if !ret {
					glog.Error("Parse CfgTowerLevel.LevelDrop field error, value:", arrayLevelDrop[j])
					return false
				}
				data.LevelDrop = append(data.LevelDrop, v)
			}
		}

		/* parse HangUpDrop field */
		vecHangUpDrop, _ := parse.GetFieldByName(uint32(i), "hangUpDrop")
		if vecHangUpDrop != "" {
			arrayHangUpDrop := strings.Split(vecHangUpDrop, ",")
			for j := 0; j < len(arrayHangUpDrop); j++ {
				v, ret := String2Int32(arrayHangUpDrop[j])
				if !ret {
					glog.Error("Parse CfgTowerLevel.HangUpDrop field error, value:", arrayHangUpDrop[j])
					return false
				}
				data.HangUpDrop = append(data.HangUpDrop, v)
			}
		}

		/* parse HangUpSeconds field */
		vHangUpSeconds, _ := parse.GetFieldByName(uint32(i), "hangUpSeconds")
		var HangUpSecondsRet bool
		data.HangUpSeconds, HangUpSecondsRet = String2Int32(vHangUpSeconds)
		if !HangUpSecondsRet {
			glog.Error("Parse CfgTowerLevel.HangUpSeconds field error,value:", vHangUpSeconds)
			return false
		}

		/* parse LvPara field */
		vecLvPara, _ := parse.GetFieldByName(uint32(i), "lvPara")
		arrayLvPara := strings.Split(vecLvPara, ",")
		for j := 0; j < len(arrayLvPara); j++ {
			v := arrayLvPara[j]
			data.LvPara = append(data.LvPara, v)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgTowerLevelConfig) Clear() {
}

func (c *CfgTowerLevelConfig) Find(id int32) (*CfgTowerLevel, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgTowerLevelConfig) GetAllData() map[int32]*CfgTowerLevel {
	return c.data
}

func (c *CfgTowerLevelConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.BattleId, ",", v.IsAchievement, ",", v.AchievementCondition, ",", v.AchievementID, ",", v.AchievementDrop, ",", v.LevelDrop, ",", v.HangUpDrop, ",", v.HangUpSeconds, ",", v.LvPara)
	}
}
