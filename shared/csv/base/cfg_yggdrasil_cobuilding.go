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

type CfgYggdrasilCobuilding struct {
	Id           int32
	BuildingType int32
	PosX         int32
	PosY         int32
	BuildCost    []string
	UsingCost    []string
	Intimacy     int32
}

type CfgYggdrasilCobuildingConfig struct {
	data map[int32]*CfgYggdrasilCobuilding
}

func NewCfgYggdrasilCobuildingConfig() *CfgYggdrasilCobuildingConfig {
	return &CfgYggdrasilCobuildingConfig{
		data: make(map[int32]*CfgYggdrasilCobuilding),
	}
}

func (c *CfgYggdrasilCobuildingConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgYggdrasilCobuilding)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgYggdrasilCobuilding.Id field error,value:", vId)
			return false
		}

		/* parse BuildingType field */
		vBuildingType, _ := parse.GetFieldByName(uint32(i), "buildingType")
		var BuildingTypeRet bool
		data.BuildingType, BuildingTypeRet = String2Int32(vBuildingType)
		if !BuildingTypeRet {
			glog.Error("Parse CfgYggdrasilCobuilding.BuildingType field error,value:", vBuildingType)
			return false
		}

		/* parse PosX field */
		vPosX, _ := parse.GetFieldByName(uint32(i), "posX")
		var PosXRet bool
		data.PosX, PosXRet = String2Int32(vPosX)
		if !PosXRet {
			glog.Error("Parse CfgYggdrasilCobuilding.PosX field error,value:", vPosX)
			return false
		}

		/* parse PosY field */
		vPosY, _ := parse.GetFieldByName(uint32(i), "posY")
		var PosYRet bool
		data.PosY, PosYRet = String2Int32(vPosY)
		if !PosYRet {
			glog.Error("Parse CfgYggdrasilCobuilding.PosY field error,value:", vPosY)
			return false
		}

		/* parse BuildCost field */
		vecBuildCost, _ := parse.GetFieldByName(uint32(i), "buildCost")
		arrayBuildCost := strings.Split(vecBuildCost, ",")
		for j := 0; j < len(arrayBuildCost); j++ {
			v := arrayBuildCost[j]
			data.BuildCost = append(data.BuildCost, v)
		}

		/* parse UsingCost field */
		vecUsingCost, _ := parse.GetFieldByName(uint32(i), "usingCost")
		arrayUsingCost := strings.Split(vecUsingCost, ",")
		for j := 0; j < len(arrayUsingCost); j++ {
			v := arrayUsingCost[j]
			data.UsingCost = append(data.UsingCost, v)
		}

		/* parse Intimacy field */
		vIntimacy, _ := parse.GetFieldByName(uint32(i), "intimacy")
		var IntimacyRet bool
		data.Intimacy, IntimacyRet = String2Int32(vIntimacy)
		if !IntimacyRet {
			glog.Error("Parse CfgYggdrasilCobuilding.Intimacy field error,value:", vIntimacy)
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

func (c *CfgYggdrasilCobuildingConfig) Clear() {
}

func (c *CfgYggdrasilCobuildingConfig) Find(id int32) (*CfgYggdrasilCobuilding, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgYggdrasilCobuildingConfig) GetAllData() map[int32]*CfgYggdrasilCobuilding {
	return c.data
}

func (c *CfgYggdrasilCobuildingConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.BuildingType, ",", v.PosX, ",", v.PosY, ",", v.BuildCost, ",", v.UsingCost, ",", v.Intimacy)
	}
}
