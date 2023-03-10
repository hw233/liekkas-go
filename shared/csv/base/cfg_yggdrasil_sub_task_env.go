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

type CfgYggdrasilSubTaskEnv struct {
	Id                     int32
	SubTaskId              int32
	SubTaskType            int32
	Target                 []string
	Positions              []string
	Npc                    []int32
	AddTaskItem            []string
	DeleteTaskItem         []string
	ClearPosGroup          []string
	CreateGroupObject      []string
	ChangeObjectState      []int32
	ChangeTerrainState     string
	TerrainStateDeleteAt   int32
	CreateDeleteAtGroupObj []string
	DeleteAt               int32
}

type CfgYggdrasilSubTaskEnvConfig struct {
	data map[int32]*CfgYggdrasilSubTaskEnv
}

func NewCfgYggdrasilSubTaskEnvConfig() *CfgYggdrasilSubTaskEnvConfig {
	return &CfgYggdrasilSubTaskEnvConfig{
		data: make(map[int32]*CfgYggdrasilSubTaskEnv),
	}
}

func (c *CfgYggdrasilSubTaskEnvConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgYggdrasilSubTaskEnv)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgYggdrasilSubTaskEnv.Id field error,value:", vId)
			return false
		}

		/* parse SubTaskId field */
		vSubTaskId, _ := parse.GetFieldByName(uint32(i), "subTaskId")
		var SubTaskIdRet bool
		data.SubTaskId, SubTaskIdRet = String2Int32(vSubTaskId)
		if !SubTaskIdRet {
			glog.Error("Parse CfgYggdrasilSubTaskEnv.SubTaskId field error,value:", vSubTaskId)
			return false
		}

		/* parse SubTaskType field */
		vSubTaskType, _ := parse.GetFieldByName(uint32(i), "subTaskType")
		var SubTaskTypeRet bool
		data.SubTaskType, SubTaskTypeRet = String2Int32(vSubTaskType)
		if !SubTaskTypeRet {
			glog.Error("Parse CfgYggdrasilSubTaskEnv.SubTaskType field error,value:", vSubTaskType)
			return false
		}

		/* parse Target field */
		vecTarget, _ := parse.GetFieldByName(uint32(i), "target")
		arrayTarget := strings.Split(vecTarget, ",")
		for j := 0; j < len(arrayTarget); j++ {
			v := arrayTarget[j]
			data.Target = append(data.Target, v)
		}

		/* parse Positions field */
		vecPositions, _ := parse.GetFieldByName(uint32(i), "positions")
		arrayPositions := strings.Split(vecPositions, ",")
		for j := 0; j < len(arrayPositions); j++ {
			v := arrayPositions[j]
			data.Positions = append(data.Positions, v)
		}

		/* parse Npc field */
		vecNpc, _ := parse.GetFieldByName(uint32(i), "npc")
		if vecNpc != "" {
			arrayNpc := strings.Split(vecNpc, ",")
			for j := 0; j < len(arrayNpc); j++ {
				v, ret := String2Int32(arrayNpc[j])
				if !ret {
					glog.Error("Parse CfgYggdrasilSubTaskEnv.Npc field error, value:", arrayNpc[j])
					return false
				}
				data.Npc = append(data.Npc, v)
			}
		}

		/* parse AddTaskItem field */
		vecAddTaskItem, _ := parse.GetFieldByName(uint32(i), "addTaskItem")
		arrayAddTaskItem := strings.Split(vecAddTaskItem, ",")
		for j := 0; j < len(arrayAddTaskItem); j++ {
			v := arrayAddTaskItem[j]
			data.AddTaskItem = append(data.AddTaskItem, v)
		}

		/* parse DeleteTaskItem field */
		vecDeleteTaskItem, _ := parse.GetFieldByName(uint32(i), "deleteTaskItem")
		arrayDeleteTaskItem := strings.Split(vecDeleteTaskItem, ",")
		for j := 0; j < len(arrayDeleteTaskItem); j++ {
			v := arrayDeleteTaskItem[j]
			data.DeleteTaskItem = append(data.DeleteTaskItem, v)
		}

		/* parse ClearPosGroup field */
		vecClearPosGroup, _ := parse.GetFieldByName(uint32(i), "clearPosGroup")
		arrayClearPosGroup := strings.Split(vecClearPosGroup, ",")
		for j := 0; j < len(arrayClearPosGroup); j++ {
			v := arrayClearPosGroup[j]
			data.ClearPosGroup = append(data.ClearPosGroup, v)
		}

		/* parse CreateGroupObject field */
		vecCreateGroupObject, _ := parse.GetFieldByName(uint32(i), "createGroupObject")
		arrayCreateGroupObject := strings.Split(vecCreateGroupObject, ",")
		for j := 0; j < len(arrayCreateGroupObject); j++ {
			v := arrayCreateGroupObject[j]
			data.CreateGroupObject = append(data.CreateGroupObject, v)
		}

		/* parse ChangeObjectState field */
		vecChangeObjectState, _ := parse.GetFieldByName(uint32(i), "changeObjectState")
		if vecChangeObjectState != "" {
			arrayChangeObjectState := strings.Split(vecChangeObjectState, ",")
			for j := 0; j < len(arrayChangeObjectState); j++ {
				v, ret := String2Int32(arrayChangeObjectState[j])
				if !ret {
					glog.Error("Parse CfgYggdrasilSubTaskEnv.ChangeObjectState field error, value:", arrayChangeObjectState[j])
					return false
				}
				data.ChangeObjectState = append(data.ChangeObjectState, v)
			}
		}

		/* parse ChangeTerrainState field */
		data.ChangeTerrainState, _ = parse.GetFieldByName(uint32(i), "changeTerrainState")

		/* parse TerrainStateDeleteAt field */
		vTerrainStateDeleteAt, _ := parse.GetFieldByName(uint32(i), "terrainStateDeleteAt")
		var TerrainStateDeleteAtRet bool
		data.TerrainStateDeleteAt, TerrainStateDeleteAtRet = String2Int32(vTerrainStateDeleteAt)
		if !TerrainStateDeleteAtRet {
			glog.Error("Parse CfgYggdrasilSubTaskEnv.TerrainStateDeleteAt field error,value:", vTerrainStateDeleteAt)
			return false
		}

		/* parse CreateDeleteAtGroupObj field */
		vecCreateDeleteAtGroupObj, _ := parse.GetFieldByName(uint32(i), "createDeleteAtGroupObj")
		arrayCreateDeleteAtGroupObj := strings.Split(vecCreateDeleteAtGroupObj, ",")
		for j := 0; j < len(arrayCreateDeleteAtGroupObj); j++ {
			v := arrayCreateDeleteAtGroupObj[j]
			data.CreateDeleteAtGroupObj = append(data.CreateDeleteAtGroupObj, v)
		}

		/* parse DeleteAt field */
		vDeleteAt, _ := parse.GetFieldByName(uint32(i), "deleteAt")
		var DeleteAtRet bool
		data.DeleteAt, DeleteAtRet = String2Int32(vDeleteAt)
		if !DeleteAtRet {
			glog.Error("Parse CfgYggdrasilSubTaskEnv.DeleteAt field error,value:", vDeleteAt)
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

func (c *CfgYggdrasilSubTaskEnvConfig) Clear() {
}

func (c *CfgYggdrasilSubTaskEnvConfig) Find(id int32) (*CfgYggdrasilSubTaskEnv, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgYggdrasilSubTaskEnvConfig) GetAllData() map[int32]*CfgYggdrasilSubTaskEnv {
	return c.data
}

func (c *CfgYggdrasilSubTaskEnvConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.SubTaskId, ",", v.SubTaskType, ",", v.Target, ",", v.Positions, ",", v.Npc, ",", v.AddTaskItem, ",", v.DeleteTaskItem, ",", v.ClearPosGroup, ",", v.CreateGroupObject, ",", v.ChangeObjectState, ",", v.ChangeTerrainState, ",", v.TerrainStateDeleteAt, ",", v.CreateDeleteAtGroupObj, ",", v.DeleteAt)
	}
}
