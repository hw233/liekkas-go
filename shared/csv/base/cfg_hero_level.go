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

type CfgHeroLevel struct {
	Id  int32
	Exp int32
}

type CfgHeroLevelConfig struct {
	data map[int32]*CfgHeroLevel
}

func NewCfgHeroLevelConfig() *CfgHeroLevelConfig {
	return &CfgHeroLevelConfig{
		data: make(map[int32]*CfgHeroLevel),
	}
}

func (c *CfgHeroLevelConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgHeroLevel)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgHeroLevel.Id field error,value:", vId)
			return false
		}

		/* parse Exp field */
		vExp, _ := parse.GetFieldByName(uint32(i), "exp")
		var ExpRet bool
		data.Exp, ExpRet = String2Int32(vExp)
		if !ExpRet {
			glog.Error("Parse CfgHeroLevel.Exp field error,value:", vExp)
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

func (c *CfgHeroLevelConfig) Clear() {
}

func (c *CfgHeroLevelConfig) Find(id int32) (*CfgHeroLevel, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgHeroLevelConfig) GetAllData() map[int32]*CfgHeroLevel {
	return c.data
}

func (c *CfgHeroLevelConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Exp)
	}
}
