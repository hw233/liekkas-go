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

type CfgYggdrasilDispatch struct {
	Id                  int32
	AreaId              int32
	Type                int32
	TimeCost            int32
	TaskStarsCount      int32
	TeamSize            int32
	GuildCharacterNum   int32
	GuildCharacterId    int32
	NecessaryConditions []string
	ExtraConditions     []string
	BaseDropId          int32
	ExtraDropId         []int32
	CloseTime           int32
}

type CfgYggdrasilDispatchConfig struct {
	data map[int32]*CfgYggdrasilDispatch
}

func NewCfgYggdrasilDispatchConfig() *CfgYggdrasilDispatchConfig {
	return &CfgYggdrasilDispatchConfig{
		data: make(map[int32]*CfgYggdrasilDispatch),
	}
}

func (c *CfgYggdrasilDispatchConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgYggdrasilDispatch)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgYggdrasilDispatch.Id field error,value:", vId)
			return false
		}

		/* parse AreaId field */
		vAreaId, _ := parse.GetFieldByName(uint32(i), "areaId")
		var AreaIdRet bool
		data.AreaId, AreaIdRet = String2Int32(vAreaId)
		if !AreaIdRet {
			glog.Error("Parse CfgYggdrasilDispatch.AreaId field error,value:", vAreaId)
			return false
		}

		/* parse Type field */
		vType, _ := parse.GetFieldByName(uint32(i), "type")
		var TypeRet bool
		data.Type, TypeRet = String2Int32(vType)
		if !TypeRet {
			glog.Error("Parse CfgYggdrasilDispatch.Type field error,value:", vType)
			return false
		}

		/* parse TimeCost field */
		vTimeCost, _ := parse.GetFieldByName(uint32(i), "timeCost")
		var TimeCostRet bool
		data.TimeCost, TimeCostRet = String2Int32(vTimeCost)
		if !TimeCostRet {
			glog.Error("Parse CfgYggdrasilDispatch.TimeCost field error,value:", vTimeCost)
			return false
		}

		/* parse TaskStarsCount field */
		vTaskStarsCount, _ := parse.GetFieldByName(uint32(i), "taskStarsCount")
		var TaskStarsCountRet bool
		data.TaskStarsCount, TaskStarsCountRet = String2Int32(vTaskStarsCount)
		if !TaskStarsCountRet {
			glog.Error("Parse CfgYggdrasilDispatch.TaskStarsCount field error,value:", vTaskStarsCount)
			return false
		}

		/* parse TeamSize field */
		vTeamSize, _ := parse.GetFieldByName(uint32(i), "teamSize")
		var TeamSizeRet bool
		data.TeamSize, TeamSizeRet = String2Int32(vTeamSize)
		if !TeamSizeRet {
			glog.Error("Parse CfgYggdrasilDispatch.TeamSize field error,value:", vTeamSize)
			return false
		}

		/* parse GuildCharacterNum field */
		vGuildCharacterNum, _ := parse.GetFieldByName(uint32(i), "guildCharacterNum")
		var GuildCharacterNumRet bool
		data.GuildCharacterNum, GuildCharacterNumRet = String2Int32(vGuildCharacterNum)
		if !GuildCharacterNumRet {
			glog.Error("Parse CfgYggdrasilDispatch.GuildCharacterNum field error,value:", vGuildCharacterNum)
			return false
		}

		/* parse GuildCharacterId field */
		vGuildCharacterId, _ := parse.GetFieldByName(uint32(i), "guildCharacterId")
		var GuildCharacterIdRet bool
		data.GuildCharacterId, GuildCharacterIdRet = String2Int32(vGuildCharacterId)
		if !GuildCharacterIdRet {
			glog.Error("Parse CfgYggdrasilDispatch.GuildCharacterId field error,value:", vGuildCharacterId)
			return false
		}

		/* parse NecessaryConditions field */
		vecNecessaryConditions, _ := parse.GetFieldByName(uint32(i), "necessaryConditions")
		arrayNecessaryConditions := strings.Split(vecNecessaryConditions, ",")
		for j := 0; j < len(arrayNecessaryConditions); j++ {
			v := arrayNecessaryConditions[j]
			data.NecessaryConditions = append(data.NecessaryConditions, v)
		}

		/* parse ExtraConditions field */
		vecExtraConditions, _ := parse.GetFieldByName(uint32(i), "extraConditions")
		arrayExtraConditions := strings.Split(vecExtraConditions, ",")
		for j := 0; j < len(arrayExtraConditions); j++ {
			v := arrayExtraConditions[j]
			data.ExtraConditions = append(data.ExtraConditions, v)
		}

		/* parse BaseDropId field */
		vBaseDropId, _ := parse.GetFieldByName(uint32(i), "baseDropId")
		var BaseDropIdRet bool
		data.BaseDropId, BaseDropIdRet = String2Int32(vBaseDropId)
		if !BaseDropIdRet {
			glog.Error("Parse CfgYggdrasilDispatch.BaseDropId field error,value:", vBaseDropId)
			return false
		}

		/* parse ExtraDropId field */
		vecExtraDropId, _ := parse.GetFieldByName(uint32(i), "extraDropId")
		if vecExtraDropId != "" {
			arrayExtraDropId := strings.Split(vecExtraDropId, ",")
			for j := 0; j < len(arrayExtraDropId); j++ {
				v, ret := String2Int32(arrayExtraDropId[j])
				if !ret {
					glog.Error("Parse CfgYggdrasilDispatch.ExtraDropId field error, value:", arrayExtraDropId[j])
					return false
				}
				data.ExtraDropId = append(data.ExtraDropId, v)
			}
		}

		/* parse CloseTime field */
		vCloseTime, _ := parse.GetFieldByName(uint32(i), "closeTime")
		var CloseTimeRet bool
		data.CloseTime, CloseTimeRet = String2Int32(vCloseTime)
		if !CloseTimeRet {
			glog.Error("Parse CfgYggdrasilDispatch.CloseTime field error,value:", vCloseTime)
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

func (c *CfgYggdrasilDispatchConfig) Clear() {
}

func (c *CfgYggdrasilDispatchConfig) Find(id int32) (*CfgYggdrasilDispatch, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgYggdrasilDispatchConfig) GetAllData() map[int32]*CfgYggdrasilDispatch {
	return c.data
}

func (c *CfgYggdrasilDispatchConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.AreaId, ",", v.Type, ",", v.TimeCost, ",", v.TaskStarsCount, ",", v.TeamSize, ",", v.GuildCharacterNum, ",", v.GuildCharacterId, ",", v.NecessaryConditions, ",", v.ExtraConditions, ",", v.BaseDropId, ",", v.ExtraDropId, ",", v.CloseTime)
	}
}
