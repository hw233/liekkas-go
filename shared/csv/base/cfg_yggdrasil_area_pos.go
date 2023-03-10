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

type CfgYggdrasilAreaPos struct {
	Id            int32
	AreaGroupId   int32
	AreaPosX      int32
	AreaPosY      int32
	AreaPosHeight int32
	PosType       int32
	ObjectId      int32
}

type CfgYggdrasilAreaPosConfig struct {
	data map[int32]*CfgYggdrasilAreaPos
}

func NewCfgYggdrasilAreaPosConfig() *CfgYggdrasilAreaPosConfig {
	return &CfgYggdrasilAreaPosConfig{
		data: make(map[int32]*CfgYggdrasilAreaPos),
	}
}

func (c *CfgYggdrasilAreaPosConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgYggdrasilAreaPos)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgYggdrasilAreaPos.Id field error,value:", vId)
			return false
		}

		/* parse AreaGroupId field */
		vAreaGroupId, _ := parse.GetFieldByName(uint32(i), "areaGroupId")
		var AreaGroupIdRet bool
		data.AreaGroupId, AreaGroupIdRet = String2Int32(vAreaGroupId)
		if !AreaGroupIdRet {
			glog.Error("Parse CfgYggdrasilAreaPos.AreaGroupId field error,value:", vAreaGroupId)
			return false
		}

		/* parse AreaPosX field */
		vAreaPosX, _ := parse.GetFieldByName(uint32(i), "areaPosX")
		var AreaPosXRet bool
		data.AreaPosX, AreaPosXRet = String2Int32(vAreaPosX)
		if !AreaPosXRet {
			glog.Error("Parse CfgYggdrasilAreaPos.AreaPosX field error,value:", vAreaPosX)
			return false
		}

		/* parse AreaPosY field */
		vAreaPosY, _ := parse.GetFieldByName(uint32(i), "areaPosY")
		var AreaPosYRet bool
		data.AreaPosY, AreaPosYRet = String2Int32(vAreaPosY)
		if !AreaPosYRet {
			glog.Error("Parse CfgYggdrasilAreaPos.AreaPosY field error,value:", vAreaPosY)
			return false
		}

		/* parse AreaPosHeight field */
		vAreaPosHeight, _ := parse.GetFieldByName(uint32(i), "areaPosHeight")
		var AreaPosHeightRet bool
		data.AreaPosHeight, AreaPosHeightRet = String2Int32(vAreaPosHeight)
		if !AreaPosHeightRet {
			glog.Error("Parse CfgYggdrasilAreaPos.AreaPosHeight field error,value:", vAreaPosHeight)
			return false
		}

		/* parse PosType field */
		vPosType, _ := parse.GetFieldByName(uint32(i), "posType")
		var PosTypeRet bool
		data.PosType, PosTypeRet = String2Int32(vPosType)
		if !PosTypeRet {
			glog.Error("Parse CfgYggdrasilAreaPos.PosType field error,value:", vPosType)
			return false
		}

		/* parse ObjectId field */
		vObjectId, _ := parse.GetFieldByName(uint32(i), "objectId")
		var ObjectIdRet bool
		data.ObjectId, ObjectIdRet = String2Int32(vObjectId)
		if !ObjectIdRet {
			glog.Error("Parse CfgYggdrasilAreaPos.ObjectId field error,value:", vObjectId)
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

func (c *CfgYggdrasilAreaPosConfig) Clear() {
}

func (c *CfgYggdrasilAreaPosConfig) Find(id int32) (*CfgYggdrasilAreaPos, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgYggdrasilAreaPosConfig) GetAllData() map[int32]*CfgYggdrasilAreaPos {
	return c.data
}

func (c *CfgYggdrasilAreaPosConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.AreaGroupId, ",", v.AreaPosX, ",", v.AreaPosY, ",", v.AreaPosHeight, ",", v.PosType, ",", v.ObjectId)
	}
}
