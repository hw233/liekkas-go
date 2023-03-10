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

type CfgGraveyardMainTowerLevel struct {
	Id               int32
	ExploreCondition int32
	BuildsCondition  []string
	PopulationCount  int32
	BuildsCountLimit []string
	UnlockBuildArea  []string
}

type CfgGraveyardMainTowerLevelConfig struct {
	data map[int32]*CfgGraveyardMainTowerLevel
}

func NewCfgGraveyardMainTowerLevelConfig() *CfgGraveyardMainTowerLevelConfig {
	return &CfgGraveyardMainTowerLevelConfig{
		data: make(map[int32]*CfgGraveyardMainTowerLevel),
	}
}

func (c *CfgGraveyardMainTowerLevelConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgGraveyardMainTowerLevel)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgGraveyardMainTowerLevel.Id field error,value:", vId)
			return false
		}

		/* parse ExploreCondition field */
		vExploreCondition, _ := parse.GetFieldByName(uint32(i), "exploreCondition")
		var ExploreConditionRet bool
		data.ExploreCondition, ExploreConditionRet = String2Int32(vExploreCondition)
		if !ExploreConditionRet {
			glog.Error("Parse CfgGraveyardMainTowerLevel.ExploreCondition field error,value:", vExploreCondition)
			return false
		}

		/* parse BuildsCondition field */
		vecBuildsCondition, _ := parse.GetFieldByName(uint32(i), "buildsCondition")
		arrayBuildsCondition := strings.Split(vecBuildsCondition, ",")
		for j := 0; j < len(arrayBuildsCondition); j++ {
			v := arrayBuildsCondition[j]
			data.BuildsCondition = append(data.BuildsCondition, v)
		}

		/* parse PopulationCount field */
		vPopulationCount, _ := parse.GetFieldByName(uint32(i), "populationCount")
		var PopulationCountRet bool
		data.PopulationCount, PopulationCountRet = String2Int32(vPopulationCount)
		if !PopulationCountRet {
			glog.Error("Parse CfgGraveyardMainTowerLevel.PopulationCount field error,value:", vPopulationCount)
			return false
		}

		/* parse BuildsCountLimit field */
		vecBuildsCountLimit, _ := parse.GetFieldByName(uint32(i), "buildsCountLimit")
		arrayBuildsCountLimit := strings.Split(vecBuildsCountLimit, ",")
		for j := 0; j < len(arrayBuildsCountLimit); j++ {
			v := arrayBuildsCountLimit[j]
			data.BuildsCountLimit = append(data.BuildsCountLimit, v)
		}

		/* parse UnlockBuildArea field */
		vecUnlockBuildArea, _ := parse.GetFieldByName(uint32(i), "unlockBuildArea")
		arrayUnlockBuildArea := strings.Split(vecUnlockBuildArea, ",")
		for j := 0; j < len(arrayUnlockBuildArea); j++ {
			v := arrayUnlockBuildArea[j]
			data.UnlockBuildArea = append(data.UnlockBuildArea, v)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgGraveyardMainTowerLevelConfig) Clear() {
}

func (c *CfgGraveyardMainTowerLevelConfig) Find(id int32) (*CfgGraveyardMainTowerLevel, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgGraveyardMainTowerLevelConfig) GetAllData() map[int32]*CfgGraveyardMainTowerLevel {
	return c.data
}

func (c *CfgGraveyardMainTowerLevelConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.ExploreCondition, ",", v.BuildsCondition, ",", v.PopulationCount, ",", v.BuildsCountLimit, ",", v.UnlockBuildArea)
	}
}
