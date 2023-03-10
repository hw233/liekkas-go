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

type CfgRarity struct {
	Id     int32
	Rarity string
}

type CfgRarityConfig struct {
	data map[int32]*CfgRarity
}

func NewCfgRarityConfig() *CfgRarityConfig {
	return &CfgRarityConfig{
		data: make(map[int32]*CfgRarity),
	}
}

func (c *CfgRarityConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgRarity)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgRarity.Id field error,value:", vId)
			return false
		}

		/* parse Rarity field */
		data.Rarity, _ = parse.GetFieldByName(uint32(i), "rarity")

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgRarityConfig) Clear() {
}

func (c *CfgRarityConfig) Find(id int32) (*CfgRarity, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgRarityConfig) GetAllData() map[int32]*CfgRarity {
	return c.data
}

func (c *CfgRarityConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Rarity)
	}
}
