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

type CfgCharacterRarity struct {
	Id     int32
	Rarity string
}

type CfgCharacterRarityConfig struct {
	data map[int32]*CfgCharacterRarity
}

func NewCfgCharacterRarityConfig() *CfgCharacterRarityConfig {
	return &CfgCharacterRarityConfig{
		data: make(map[int32]*CfgCharacterRarity),
	}
}

func (c *CfgCharacterRarityConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgCharacterRarity)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgCharacterRarity.Id field error,value:", vId)
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

func (c *CfgCharacterRarityConfig) Clear() {
}

func (c *CfgCharacterRarityConfig) Find(id int32) (*CfgCharacterRarity, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgCharacterRarityConfig) GetAllData() map[int32]*CfgCharacterRarity {
	return c.data
}

func (c *CfgCharacterRarityConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Rarity)
	}
}
