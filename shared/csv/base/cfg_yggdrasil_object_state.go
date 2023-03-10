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
)

type CfgYggdrasilObjectState struct {
	Id              int32
	ObjectType      int32
	PosType         int32
	NextState       int32
	ServerSubTaskId int32
	ObjectParam     string
}

type CfgYggdrasilObjectStateConfig struct {
	data map[int32]*CfgYggdrasilObjectState
}

func NewCfgYggdrasilObjectStateConfig() *CfgYggdrasilObjectStateConfig {
	return &CfgYggdrasilObjectStateConfig{
		data: make(map[int32]*CfgYggdrasilObjectState),
	}
}

func (c *CfgYggdrasilObjectStateConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgYggdrasilObjectState)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgYggdrasilObjectState.Id field error,value:", vId)
			return false
		}

		/* parse ObjectType field */
		vObjectType, _ := parse.GetFieldByName(uint32(i), "objectType")
		var ObjectTypeRet bool
		data.ObjectType, ObjectTypeRet = String2Int32(vObjectType)
		if !ObjectTypeRet {
			glog.Error("Parse CfgYggdrasilObjectState.ObjectType field error,value:", vObjectType)
			return false
		}

		/* parse PosType field */
		vPosType, _ := parse.GetFieldByName(uint32(i), "posType")
		var PosTypeRet bool
		data.PosType, PosTypeRet = String2Int32(vPosType)
		if !PosTypeRet {
			glog.Error("Parse CfgYggdrasilObjectState.PosType field error,value:", vPosType)
			return false
		}

		/* parse NextState field */
		vNextState, _ := parse.GetFieldByName(uint32(i), "nextState")
		var NextStateRet bool
		data.NextState, NextStateRet = String2Int32(vNextState)
		if !NextStateRet {
			glog.Error("Parse CfgYggdrasilObjectState.NextState field error,value:", vNextState)
			return false
		}

		/* parse ServerSubTaskId field */
		vServerSubTaskId, _ := parse.GetFieldByName(uint32(i), "serverSubTaskId")
		var ServerSubTaskIdRet bool
		data.ServerSubTaskId, ServerSubTaskIdRet = String2Int32(vServerSubTaskId)
		if !ServerSubTaskIdRet {
			glog.Error("Parse CfgYggdrasilObjectState.ServerSubTaskId field error,value:", vServerSubTaskId)
			return false
		}

		/* parse ObjectParam field */
		data.ObjectParam, _ = parse.GetFieldByName(uint32(i), "objectParam")

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgYggdrasilObjectStateConfig) Clear() {
}

func (c *CfgYggdrasilObjectStateConfig) Find(id int32) (*CfgYggdrasilObjectState, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgYggdrasilObjectStateConfig) GetAllData() map[int32]*CfgYggdrasilObjectState {
	return c.data
}

func (c *CfgYggdrasilObjectStateConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.ObjectType, ",", v.PosType, ",", v.NextState, ",", v.ServerSubTaskId, ",", v.ObjectParam)
	}
}
