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

type CfgCharacterSkin struct {
	Id int32
}

type CfgCharacterSkinConfig struct {
	data map[int32]*CfgCharacterSkin
}

func NewCfgCharacterSkinConfig() *CfgCharacterSkinConfig {
	return &CfgCharacterSkinConfig{
		data: make(map[int32]*CfgCharacterSkin),
	}
}

func (c *CfgCharacterSkinConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgCharacterSkin)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgCharacterSkin.Id field error,value:", vId)
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

func (c *CfgCharacterSkinConfig) Clear() {
}

func (c *CfgCharacterSkinConfig) Find(id int32) (*CfgCharacterSkin, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgCharacterSkinConfig) GetAllData() map[int32]*CfgCharacterSkin {
	return c.data
}

func (c *CfgCharacterSkinConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id)
	}
}
