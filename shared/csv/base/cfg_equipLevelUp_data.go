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

type CfgEquipLevelUpData struct {
	Id     int32
	Level  int32
	Rarity int32
	Exp    int32
}

type CfgEquipLevelUpDataConfig struct {
	data map[int32]*CfgEquipLevelUpData
}

func NewCfgEquipLevelUpDataConfig() *CfgEquipLevelUpDataConfig {
	return &CfgEquipLevelUpDataConfig{
		data: make(map[int32]*CfgEquipLevelUpData),
	}
}

func (c *CfgEquipLevelUpDataConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgEquipLevelUpData)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgEquipLevelUpData.Id field error,value:", vId)
			return false
		}

		/* parse Level field */
		vLevel, _ := parse.GetFieldByName(uint32(i), "level")
		var LevelRet bool
		data.Level, LevelRet = String2Int32(vLevel)
		if !LevelRet {
			glog.Error("Parse CfgEquipLevelUpData.Level field error,value:", vLevel)
			return false
		}

		/* parse Rarity field */
		vRarity, _ := parse.GetFieldByName(uint32(i), "rarity")
		var RarityRet bool
		data.Rarity, RarityRet = String2Int32(vRarity)
		if !RarityRet {
			glog.Error("Parse CfgEquipLevelUpData.Rarity field error,value:", vRarity)
			return false
		}

		/* parse Exp field */
		vExp, _ := parse.GetFieldByName(uint32(i), "exp")
		var ExpRet bool
		data.Exp, ExpRet = String2Int32(vExp)
		if !ExpRet {
			glog.Error("Parse CfgEquipLevelUpData.Exp field error,value:", vExp)
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

func (c *CfgEquipLevelUpDataConfig) Clear() {
}

func (c *CfgEquipLevelUpDataConfig) Find(id int32) (*CfgEquipLevelUpData, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgEquipLevelUpDataConfig) GetAllData() map[int32]*CfgEquipLevelUpData {
	return c.data
}

func (c *CfgEquipLevelUpDataConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Level, ",", v.Rarity, ",", v.Exp)
	}
}
