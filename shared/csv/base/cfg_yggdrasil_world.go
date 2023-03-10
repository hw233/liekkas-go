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

type CfgYggdrasilWorld struct {
	Id              int32
	MapTileDataPath []string
	AreaPosId       []string
}

type CfgYggdrasilWorldConfig struct {
	data map[int32]*CfgYggdrasilWorld
}

func NewCfgYggdrasilWorldConfig() *CfgYggdrasilWorldConfig {
	return &CfgYggdrasilWorldConfig{
		data: make(map[int32]*CfgYggdrasilWorld),
	}
}

func (c *CfgYggdrasilWorldConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgYggdrasilWorld)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgYggdrasilWorld.Id field error,value:", vId)
			return false
		}

		/* parse MapTileDataPath field */
		vecMapTileDataPath, _ := parse.GetFieldByName(uint32(i), "mapTileDataPath")
		arrayMapTileDataPath := strings.Split(vecMapTileDataPath, ",")
		for j := 0; j < len(arrayMapTileDataPath); j++ {
			v := arrayMapTileDataPath[j]
			data.MapTileDataPath = append(data.MapTileDataPath, v)
		}

		/* parse AreaPosId field */
		vecAreaPosId, _ := parse.GetFieldByName(uint32(i), "areaPosId")
		arrayAreaPosId := strings.Split(vecAreaPosId, ",")
		for j := 0; j < len(arrayAreaPosId); j++ {
			v := arrayAreaPosId[j]
			data.AreaPosId = append(data.AreaPosId, v)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgYggdrasilWorldConfig) Clear() {
}

func (c *CfgYggdrasilWorldConfig) Find(id int32) (*CfgYggdrasilWorld, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgYggdrasilWorldConfig) GetAllData() map[int32]*CfgYggdrasilWorld {
	return c.data
}

func (c *CfgYggdrasilWorldConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.MapTileDataPath, ",", v.AreaPosId)
	}
}
