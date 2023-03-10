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

type CfgYggdrasilAreaCity struct {
	Id                int32
	CostAp            int32
	WorldId           int32
	CityPos           []string
	CityCenterPosX    int32
	CityCenterPosY    int32
	CityExitPosX      int32
	CityExitPosY      int32
	CityRadius        int32
	CityBanR          int32
	CityMainTasks     []int32
	HostelCharaOffset []float64
	HostelTalkOffset  []float64
}

type CfgYggdrasilAreaCityConfig struct {
	data map[int32]*CfgYggdrasilAreaCity
}

func NewCfgYggdrasilAreaCityConfig() *CfgYggdrasilAreaCityConfig {
	return &CfgYggdrasilAreaCityConfig{
		data: make(map[int32]*CfgYggdrasilAreaCity),
	}
}

func (c *CfgYggdrasilAreaCityConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgYggdrasilAreaCity)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgYggdrasilAreaCity.Id field error,value:", vId)
			return false
		}

		/* parse CostAp field */
		vCostAp, _ := parse.GetFieldByName(uint32(i), "costAp")
		var CostApRet bool
		data.CostAp, CostApRet = String2Int32(vCostAp)
		if !CostApRet {
			glog.Error("Parse CfgYggdrasilAreaCity.CostAp field error,value:", vCostAp)
			return false
		}

		/* parse WorldId field */
		vWorldId, _ := parse.GetFieldByName(uint32(i), "worldId")
		var WorldIdRet bool
		data.WorldId, WorldIdRet = String2Int32(vWorldId)
		if !WorldIdRet {
			glog.Error("Parse CfgYggdrasilAreaCity.WorldId field error,value:", vWorldId)
			return false
		}

		/* parse CityPos field */
		vecCityPos, _ := parse.GetFieldByName(uint32(i), "cityPos")
		arrayCityPos := strings.Split(vecCityPos, ",")
		for j := 0; j < len(arrayCityPos); j++ {
			v := arrayCityPos[j]
			data.CityPos = append(data.CityPos, v)
		}

		/* parse CityCenterPosX field */
		vCityCenterPosX, _ := parse.GetFieldByName(uint32(i), "cityCenterPosX")
		var CityCenterPosXRet bool
		data.CityCenterPosX, CityCenterPosXRet = String2Int32(vCityCenterPosX)
		if !CityCenterPosXRet {
			glog.Error("Parse CfgYggdrasilAreaCity.CityCenterPosX field error,value:", vCityCenterPosX)
			return false
		}

		/* parse CityCenterPosY field */
		vCityCenterPosY, _ := parse.GetFieldByName(uint32(i), "cityCenterPosY")
		var CityCenterPosYRet bool
		data.CityCenterPosY, CityCenterPosYRet = String2Int32(vCityCenterPosY)
		if !CityCenterPosYRet {
			glog.Error("Parse CfgYggdrasilAreaCity.CityCenterPosY field error,value:", vCityCenterPosY)
			return false
		}

		/* parse CityExitPosX field */
		vCityExitPosX, _ := parse.GetFieldByName(uint32(i), "cityExitPosX")
		var CityExitPosXRet bool
		data.CityExitPosX, CityExitPosXRet = String2Int32(vCityExitPosX)
		if !CityExitPosXRet {
			glog.Error("Parse CfgYggdrasilAreaCity.CityExitPosX field error,value:", vCityExitPosX)
			return false
		}

		/* parse CityExitPosY field */
		vCityExitPosY, _ := parse.GetFieldByName(uint32(i), "cityExitPosY")
		var CityExitPosYRet bool
		data.CityExitPosY, CityExitPosYRet = String2Int32(vCityExitPosY)
		if !CityExitPosYRet {
			glog.Error("Parse CfgYggdrasilAreaCity.CityExitPosY field error,value:", vCityExitPosY)
			return false
		}

		/* parse CityRadius field */
		vCityRadius, _ := parse.GetFieldByName(uint32(i), "cityRadius")
		var CityRadiusRet bool
		data.CityRadius, CityRadiusRet = String2Int32(vCityRadius)
		if !CityRadiusRet {
			glog.Error("Parse CfgYggdrasilAreaCity.CityRadius field error,value:", vCityRadius)
			return false
		}

		/* parse CityBanR field */
		vCityBanR, _ := parse.GetFieldByName(uint32(i), "cityBanR")
		var CityBanRRet bool
		data.CityBanR, CityBanRRet = String2Int32(vCityBanR)
		if !CityBanRRet {
			glog.Error("Parse CfgYggdrasilAreaCity.CityBanR field error,value:", vCityBanR)
			return false
		}

		/* parse CityMainTasks field */
		vecCityMainTasks, _ := parse.GetFieldByName(uint32(i), "cityMainTasks")
		if vecCityMainTasks != "" {
			arrayCityMainTasks := strings.Split(vecCityMainTasks, ",")
			for j := 0; j < len(arrayCityMainTasks); j++ {
				v, ret := String2Int32(arrayCityMainTasks[j])
				if !ret {
					glog.Error("Parse CfgYggdrasilAreaCity.CityMainTasks field error, value:", arrayCityMainTasks[j])
					return false
				}
				data.CityMainTasks = append(data.CityMainTasks, v)
			}
		}

		/* parse HostelCharaOffset field */
		vecHostelCharaOffset, _ := parse.GetFieldByName(uint32(i), "hostelCharaOffset")
		arrayHostelCharaOffset := strings.Split(vecHostelCharaOffset, ",")
		for j := 0; j < len(arrayHostelCharaOffset); j++ {
			v, ret := String2Float(arrayHostelCharaOffset[j])
			if !ret {
				glog.Error("Parse CfgYggdrasilAreaCity.HostelCharaOffset field error,value:", arrayHostelCharaOffset[j])
				return false
			}
			data.HostelCharaOffset = append(data.HostelCharaOffset, v)
		}

		/* parse HostelTalkOffset field */
		vecHostelTalkOffset, _ := parse.GetFieldByName(uint32(i), "hostelTalkOffset")
		arrayHostelTalkOffset := strings.Split(vecHostelTalkOffset, ",")
		for j := 0; j < len(arrayHostelTalkOffset); j++ {
			v, ret := String2Float(arrayHostelTalkOffset[j])
			if !ret {
				glog.Error("Parse CfgYggdrasilAreaCity.HostelTalkOffset field error,value:", arrayHostelTalkOffset[j])
				return false
			}
			data.HostelTalkOffset = append(data.HostelTalkOffset, v)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgYggdrasilAreaCityConfig) Clear() {
}

func (c *CfgYggdrasilAreaCityConfig) Find(id int32) (*CfgYggdrasilAreaCity, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgYggdrasilAreaCityConfig) GetAllData() map[int32]*CfgYggdrasilAreaCity {
	return c.data
}

func (c *CfgYggdrasilAreaCityConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.CostAp, ",", v.WorldId, ",", v.CityPos, ",", v.CityCenterPosX, ",", v.CityCenterPosY, ",", v.CityExitPosX, ",", v.CityExitPosY, ",", v.CityRadius, ",", v.CityBanR, ",", v.CityMainTasks, ",", v.HostelCharaOffset, ",", v.HostelTalkOffset)
	}
}
