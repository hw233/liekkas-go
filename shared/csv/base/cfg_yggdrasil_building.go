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

type CfgYggdrasilBuilding struct {
	Id               int32
	BuildingType     int32
	UnlockConditions []string
	SingleCost       []string
	UsingCost        []string
	UsingParam       []int32
	UsingTimes       int32
	BuildingR        int32
	MatchParam       int32
	Intimacy         int32
}

type CfgYggdrasilBuildingConfig struct {
	data map[int32]*CfgYggdrasilBuilding
}

func NewCfgYggdrasilBuildingConfig() *CfgYggdrasilBuildingConfig {
	return &CfgYggdrasilBuildingConfig{
		data: make(map[int32]*CfgYggdrasilBuilding),
	}
}

func (c *CfgYggdrasilBuildingConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgYggdrasilBuilding)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgYggdrasilBuilding.Id field error,value:", vId)
			return false
		}

		/* parse BuildingType field */
		vBuildingType, _ := parse.GetFieldByName(uint32(i), "buildingType")
		var BuildingTypeRet bool
		data.BuildingType, BuildingTypeRet = String2Int32(vBuildingType)
		if !BuildingTypeRet {
			glog.Error("Parse CfgYggdrasilBuilding.BuildingType field error,value:", vBuildingType)
			return false
		}

		/* parse UnlockConditions field */
		vecUnlockConditions, _ := parse.GetFieldByName(uint32(i), "unlockConditions")
		arrayUnlockConditions := strings.Split(vecUnlockConditions, ",")
		for j := 0; j < len(arrayUnlockConditions); j++ {
			v := arrayUnlockConditions[j]
			data.UnlockConditions = append(data.UnlockConditions, v)
		}

		/* parse SingleCost field */
		vecSingleCost, _ := parse.GetFieldByName(uint32(i), "singleCost")
		arraySingleCost := strings.Split(vecSingleCost, ",")
		for j := 0; j < len(arraySingleCost); j++ {
			v := arraySingleCost[j]
			data.SingleCost = append(data.SingleCost, v)
		}

		/* parse UsingCost field */
		vecUsingCost, _ := parse.GetFieldByName(uint32(i), "usingCost")
		arrayUsingCost := strings.Split(vecUsingCost, ",")
		for j := 0; j < len(arrayUsingCost); j++ {
			v := arrayUsingCost[j]
			data.UsingCost = append(data.UsingCost, v)
		}

		/* parse UsingParam field */
		vecUsingParam, _ := parse.GetFieldByName(uint32(i), "usingParam")
		if vecUsingParam != "" {
			arrayUsingParam := strings.Split(vecUsingParam, ",")
			for j := 0; j < len(arrayUsingParam); j++ {
				v, ret := String2Int32(arrayUsingParam[j])
				if !ret {
					glog.Error("Parse CfgYggdrasilBuilding.UsingParam field error, value:", arrayUsingParam[j])
					return false
				}
				data.UsingParam = append(data.UsingParam, v)
			}
		}

		/* parse UsingTimes field */
		vUsingTimes, _ := parse.GetFieldByName(uint32(i), "usingTimes")
		var UsingTimesRet bool
		data.UsingTimes, UsingTimesRet = String2Int32(vUsingTimes)
		if !UsingTimesRet {
			glog.Error("Parse CfgYggdrasilBuilding.UsingTimes field error,value:", vUsingTimes)
			return false
		}

		/* parse BuildingR field */
		vBuildingR, _ := parse.GetFieldByName(uint32(i), "buildingR")
		var BuildingRRet bool
		data.BuildingR, BuildingRRet = String2Int32(vBuildingR)
		if !BuildingRRet {
			glog.Error("Parse CfgYggdrasilBuilding.BuildingR field error,value:", vBuildingR)
			return false
		}

		/* parse MatchParam field */
		vMatchParam, _ := parse.GetFieldByName(uint32(i), "matchParam")
		var MatchParamRet bool
		data.MatchParam, MatchParamRet = String2Int32(vMatchParam)
		if !MatchParamRet {
			glog.Error("Parse CfgYggdrasilBuilding.MatchParam field error,value:", vMatchParam)
			return false
		}

		/* parse Intimacy field */
		vIntimacy, _ := parse.GetFieldByName(uint32(i), "intimacy")
		var IntimacyRet bool
		data.Intimacy, IntimacyRet = String2Int32(vIntimacy)
		if !IntimacyRet {
			glog.Error("Parse CfgYggdrasilBuilding.Intimacy field error,value:", vIntimacy)
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

func (c *CfgYggdrasilBuildingConfig) Clear() {
}

func (c *CfgYggdrasilBuildingConfig) Find(id int32) (*CfgYggdrasilBuilding, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgYggdrasilBuildingConfig) GetAllData() map[int32]*CfgYggdrasilBuilding {
	return c.data
}

func (c *CfgYggdrasilBuildingConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.BuildingType, ",", v.UnlockConditions, ",", v.SingleCost, ",", v.UsingCost, ",", v.UsingParam, ",", v.UsingTimes, ",", v.BuildingR, ",", v.MatchParam, ",", v.Intimacy)
	}
}
