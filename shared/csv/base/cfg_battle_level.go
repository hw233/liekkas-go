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

type CfgBattleLevel struct {
	Id               int32
	LevelType        int32
	Param            []int32
	IEnergy          int32
	EnergyRecovery   int32
	ShowLevelBoss    int32
	Npc              []int32
	ControlNpc       []int32
	HostileData1     int32
	HostileData2     int32
	HostileData3     int32
	HostileData4     int32
	HostileData5     int32
	HostileData6     int32
	WinCondition1    []int32
	WinCondition2    []int32
	LevelTime        int32
	LevelSp          int32
	PlayerLimit      int32
	NextBattle       int32
	PlayerPosition   int32
	HostilePosition  int32
	LvPara           []string
	PreloadCharacter []int32
	PreloadFx        []int32
	PreloadBullet    []int32
	BattleCheck      bool
}

type CfgBattleLevelConfig struct {
	data map[int32]*CfgBattleLevel
}

func NewCfgBattleLevelConfig() *CfgBattleLevelConfig {
	return &CfgBattleLevelConfig{
		data: make(map[int32]*CfgBattleLevel),
	}
}

func (c *CfgBattleLevelConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgBattleLevel)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgBattleLevel.Id field error,value:", vId)
			return false
		}

		/* parse LevelType field */
		vLevelType, _ := parse.GetFieldByName(uint32(i), "levelType")
		var LevelTypeRet bool
		data.LevelType, LevelTypeRet = String2Int32(vLevelType)
		if !LevelTypeRet {
			glog.Error("Parse CfgBattleLevel.LevelType field error,value:", vLevelType)
			return false
		}

		/* parse Param field */
		vecParam, _ := parse.GetFieldByName(uint32(i), "param")
		if vecParam != "" {
			arrayParam := strings.Split(vecParam, ",")
			for j := 0; j < len(arrayParam); j++ {
				v, ret := String2Int32(arrayParam[j])
				if !ret {
					glog.Error("Parse CfgBattleLevel.Param field error, value:", arrayParam[j])
					return false
				}
				data.Param = append(data.Param, v)
			}
		}

		/* parse IEnergy field */
		vIEnergy, _ := parse.GetFieldByName(uint32(i), "iEnergy")
		var IEnergyRet bool
		data.IEnergy, IEnergyRet = String2Int32(vIEnergy)
		if !IEnergyRet {
			glog.Error("Parse CfgBattleLevel.IEnergy field error,value:", vIEnergy)
			return false
		}

		/* parse EnergyRecovery field */
		vEnergyRecovery, _ := parse.GetFieldByName(uint32(i), "energyRecovery")
		var EnergyRecoveryRet bool
		data.EnergyRecovery, EnergyRecoveryRet = String2Int32(vEnergyRecovery)
		if !EnergyRecoveryRet {
			glog.Error("Parse CfgBattleLevel.EnergyRecovery field error,value:", vEnergyRecovery)
			return false
		}

		/* parse ShowLevelBoss field */
		vShowLevelBoss, _ := parse.GetFieldByName(uint32(i), "showLevelBoss")
		var ShowLevelBossRet bool
		data.ShowLevelBoss, ShowLevelBossRet = String2Int32(vShowLevelBoss)
		if !ShowLevelBossRet {
			glog.Error("Parse CfgBattleLevel.ShowLevelBoss field error,value:", vShowLevelBoss)
			return false
		}

		/* parse Npc field */
		vecNpc, _ := parse.GetFieldByName(uint32(i), "npc")
		if vecNpc != "" {
			arrayNpc := strings.Split(vecNpc, ",")
			for j := 0; j < len(arrayNpc); j++ {
				v, ret := String2Int32(arrayNpc[j])
				if !ret {
					glog.Error("Parse CfgBattleLevel.Npc field error, value:", arrayNpc[j])
					return false
				}
				data.Npc = append(data.Npc, v)
			}
		}

		/* parse ControlNpc field */
		vecControlNpc, _ := parse.GetFieldByName(uint32(i), "controlNpc")
		if vecControlNpc != "" {
			arrayControlNpc := strings.Split(vecControlNpc, ",")
			for j := 0; j < len(arrayControlNpc); j++ {
				v, ret := String2Int32(arrayControlNpc[j])
				if !ret {
					glog.Error("Parse CfgBattleLevel.ControlNpc field error, value:", arrayControlNpc[j])
					return false
				}
				data.ControlNpc = append(data.ControlNpc, v)
			}
		}

		/* parse HostileData1 field */
		vHostileData1, _ := parse.GetFieldByName(uint32(i), "hostileData1")
		var HostileData1Ret bool
		data.HostileData1, HostileData1Ret = String2Int32(vHostileData1)
		if !HostileData1Ret {
			glog.Error("Parse CfgBattleLevel.HostileData1 field error,value:", vHostileData1)
			return false
		}

		/* parse HostileData2 field */
		vHostileData2, _ := parse.GetFieldByName(uint32(i), "hostileData2")
		var HostileData2Ret bool
		data.HostileData2, HostileData2Ret = String2Int32(vHostileData2)
		if !HostileData2Ret {
			glog.Error("Parse CfgBattleLevel.HostileData2 field error,value:", vHostileData2)
			return false
		}

		/* parse HostileData3 field */
		vHostileData3, _ := parse.GetFieldByName(uint32(i), "hostileData3")
		var HostileData3Ret bool
		data.HostileData3, HostileData3Ret = String2Int32(vHostileData3)
		if !HostileData3Ret {
			glog.Error("Parse CfgBattleLevel.HostileData3 field error,value:", vHostileData3)
			return false
		}

		/* parse HostileData4 field */
		vHostileData4, _ := parse.GetFieldByName(uint32(i), "hostileData4")
		var HostileData4Ret bool
		data.HostileData4, HostileData4Ret = String2Int32(vHostileData4)
		if !HostileData4Ret {
			glog.Error("Parse CfgBattleLevel.HostileData4 field error,value:", vHostileData4)
			return false
		}

		/* parse HostileData5 field */
		vHostileData5, _ := parse.GetFieldByName(uint32(i), "hostileData5")
		var HostileData5Ret bool
		data.HostileData5, HostileData5Ret = String2Int32(vHostileData5)
		if !HostileData5Ret {
			glog.Error("Parse CfgBattleLevel.HostileData5 field error,value:", vHostileData5)
			return false
		}

		/* parse HostileData6 field */
		vHostileData6, _ := parse.GetFieldByName(uint32(i), "hostileData6")
		var HostileData6Ret bool
		data.HostileData6, HostileData6Ret = String2Int32(vHostileData6)
		if !HostileData6Ret {
			glog.Error("Parse CfgBattleLevel.HostileData6 field error,value:", vHostileData6)
			return false
		}

		/* parse WinCondition1 field */
		vecWinCondition1, _ := parse.GetFieldByName(uint32(i), "winCondition1")
		if vecWinCondition1 != "" {
			arrayWinCondition1 := strings.Split(vecWinCondition1, ",")
			for j := 0; j < len(arrayWinCondition1); j++ {
				v, ret := String2Int32(arrayWinCondition1[j])
				if !ret {
					glog.Error("Parse CfgBattleLevel.WinCondition1 field error, value:", arrayWinCondition1[j])
					return false
				}
				data.WinCondition1 = append(data.WinCondition1, v)
			}
		}

		/* parse WinCondition2 field */
		vecWinCondition2, _ := parse.GetFieldByName(uint32(i), "winCondition2")
		if vecWinCondition2 != "" {
			arrayWinCondition2 := strings.Split(vecWinCondition2, ",")
			for j := 0; j < len(arrayWinCondition2); j++ {
				v, ret := String2Int32(arrayWinCondition2[j])
				if !ret {
					glog.Error("Parse CfgBattleLevel.WinCondition2 field error, value:", arrayWinCondition2[j])
					return false
				}
				data.WinCondition2 = append(data.WinCondition2, v)
			}
		}

		/* parse LevelTime field */
		vLevelTime, _ := parse.GetFieldByName(uint32(i), "levelTime")
		var LevelTimeRet bool
		data.LevelTime, LevelTimeRet = String2Int32(vLevelTime)
		if !LevelTimeRet {
			glog.Error("Parse CfgBattleLevel.LevelTime field error,value:", vLevelTime)
			return false
		}

		/* parse LevelSp field */
		vLevelSp, _ := parse.GetFieldByName(uint32(i), "levelSp")
		var LevelSpRet bool
		data.LevelSp, LevelSpRet = String2Int32(vLevelSp)
		if !LevelSpRet {
			glog.Error("Parse CfgBattleLevel.LevelSp field error,value:", vLevelSp)
			return false
		}

		/* parse PlayerLimit field */
		vPlayerLimit, _ := parse.GetFieldByName(uint32(i), "playerLimit")
		var PlayerLimitRet bool
		data.PlayerLimit, PlayerLimitRet = String2Int32(vPlayerLimit)
		if !PlayerLimitRet {
			glog.Error("Parse CfgBattleLevel.PlayerLimit field error,value:", vPlayerLimit)
			return false
		}

		/* parse NextBattle field */
		vNextBattle, _ := parse.GetFieldByName(uint32(i), "nextBattle")
		var NextBattleRet bool
		data.NextBattle, NextBattleRet = String2Int32(vNextBattle)
		if !NextBattleRet {
			glog.Error("Parse CfgBattleLevel.NextBattle field error,value:", vNextBattle)
			return false
		}

		/* parse PlayerPosition field */
		vPlayerPosition, _ := parse.GetFieldByName(uint32(i), "playerPosition")
		var PlayerPositionRet bool
		data.PlayerPosition, PlayerPositionRet = String2Int32(vPlayerPosition)
		if !PlayerPositionRet {
			glog.Error("Parse CfgBattleLevel.PlayerPosition field error,value:", vPlayerPosition)
			return false
		}

		/* parse HostilePosition field */
		vHostilePosition, _ := parse.GetFieldByName(uint32(i), "hostilePosition")
		var HostilePositionRet bool
		data.HostilePosition, HostilePositionRet = String2Int32(vHostilePosition)
		if !HostilePositionRet {
			glog.Error("Parse CfgBattleLevel.HostilePosition field error,value:", vHostilePosition)
			return false
		}

		/* parse LvPara field */
		vecLvPara, _ := parse.GetFieldByName(uint32(i), "lvPara")
		arrayLvPara := strings.Split(vecLvPara, ",")
		for j := 0; j < len(arrayLvPara); j++ {
			v := arrayLvPara[j]
			data.LvPara = append(data.LvPara, v)
		}

		/* parse PreloadCharacter field */
		vecPreloadCharacter, _ := parse.GetFieldByName(uint32(i), "preloadCharacter")
		if vecPreloadCharacter != "" {
			arrayPreloadCharacter := strings.Split(vecPreloadCharacter, ",")
			for j := 0; j < len(arrayPreloadCharacter); j++ {
				v, ret := String2Int32(arrayPreloadCharacter[j])
				if !ret {
					glog.Error("Parse CfgBattleLevel.PreloadCharacter field error, value:", arrayPreloadCharacter[j])
					return false
				}
				data.PreloadCharacter = append(data.PreloadCharacter, v)
			}
		}

		/* parse PreloadFx field */
		vecPreloadFx, _ := parse.GetFieldByName(uint32(i), "preloadFx")
		if vecPreloadFx != "" {
			arrayPreloadFx := strings.Split(vecPreloadFx, ",")
			for j := 0; j < len(arrayPreloadFx); j++ {
				v, ret := String2Int32(arrayPreloadFx[j])
				if !ret {
					glog.Error("Parse CfgBattleLevel.PreloadFx field error, value:", arrayPreloadFx[j])
					return false
				}
				data.PreloadFx = append(data.PreloadFx, v)
			}
		}

		/* parse PreloadBullet field */
		vecPreloadBullet, _ := parse.GetFieldByName(uint32(i), "preloadBullet")
		if vecPreloadBullet != "" {
			arrayPreloadBullet := strings.Split(vecPreloadBullet, ",")
			for j := 0; j < len(arrayPreloadBullet); j++ {
				v, ret := String2Int32(arrayPreloadBullet[j])
				if !ret {
					glog.Error("Parse CfgBattleLevel.PreloadBullet field error, value:", arrayPreloadBullet[j])
					return false
				}
				data.PreloadBullet = append(data.PreloadBullet, v)
			}
		}

		/* parse BattleCheck field */
		vBattleCheck, _ := parse.GetFieldByName(uint32(i), "battleCheck")
		var BattleCheckRet bool
		data.BattleCheck, BattleCheckRet = String2Bool(vBattleCheck)
		if !BattleCheckRet {
			glog.Error("Parse CfgBattleLevel.BattleCheck field error,value:", vBattleCheck)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgBattleLevelConfig) Clear() {
}

func (c *CfgBattleLevelConfig) Find(id int32) (*CfgBattleLevel, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgBattleLevelConfig) GetAllData() map[int32]*CfgBattleLevel {
	return c.data
}

func (c *CfgBattleLevelConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.LevelType, ",", v.Param, ",", v.IEnergy, ",", v.EnergyRecovery, ",", v.ShowLevelBoss, ",", v.Npc, ",", v.ControlNpc, ",", v.HostileData1, ",", v.HostileData2, ",", v.HostileData3, ",", v.HostileData4, ",", v.HostileData5, ",", v.HostileData6, ",", v.WinCondition1, ",", v.WinCondition2, ",", v.LevelTime, ",", v.LevelSp, ",", v.PlayerLimit, ",", v.NextBattle, ",", v.PlayerPosition, ",", v.HostilePosition, ",", v.LvPara, ",", v.PreloadCharacter, ",", v.PreloadFx, ",", v.PreloadBullet, ",", v.BattleCheck)
	}
}
