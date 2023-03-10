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

type CfgCharacterCareer struct {
	Id int32
}

type CfgCharacterCareerConfig struct {
	data map[int32]*CfgCharacterCareer
}

func NewCfgCharacterCareerConfig() *CfgCharacterCareerConfig {
	return &CfgCharacterCareerConfig{
		data: make(map[int32]*CfgCharacterCareer),
	}
}

func (c *CfgCharacterCareerConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgCharacterCareer)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgCharacterCareer.Id field error,value:", vId)
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

func (c *CfgCharacterCareerConfig) Clear() {
}

func (c *CfgCharacterCareerConfig) Find(id int32) (*CfgCharacterCareer, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgCharacterCareerConfig) GetAllData() map[int32]*CfgCharacterCareer {
	return c.data
}

func (c *CfgCharacterCareerConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id)
	}
}
