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

type CfgExploreFog struct {
	Id               int32
	MapId            int32
	UnlockConditions []string
}

type CfgExploreFogConfig struct {
	data map[int32]*CfgExploreFog
}

func NewCfgExploreFogConfig() *CfgExploreFogConfig {
	return &CfgExploreFogConfig{
		data: make(map[int32]*CfgExploreFog),
	}
}

func (c *CfgExploreFogConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgExploreFog)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgExploreFog.Id field error,value:", vId)
			return false
		}

		/* parse MapId field */
		vMapId, _ := parse.GetFieldByName(uint32(i), "mapId")
		var MapIdRet bool
		data.MapId, MapIdRet = String2Int32(vMapId)
		if !MapIdRet {
			glog.Error("Parse CfgExploreFog.MapId field error,value:", vMapId)
			return false
		}

		/* parse UnlockConditions field */
		vecUnlockConditions, _ := parse.GetFieldByName(uint32(i), "unlockConditions")
		arrayUnlockConditions := strings.Split(vecUnlockConditions, ",")
		for j := 0; j < len(arrayUnlockConditions); j++ {
			v := arrayUnlockConditions[j]
			data.UnlockConditions = append(data.UnlockConditions, v)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgExploreFogConfig) Clear() {
}

func (c *CfgExploreFogConfig) Find(id int32) (*CfgExploreFog, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgExploreFogConfig) GetAllData() map[int32]*CfgExploreFog {
	return c.data
}

func (c *CfgExploreFogConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.MapId, ",", v.UnlockConditions)
	}
}
