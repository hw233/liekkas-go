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

type CfgCharacterCamp struct {
	Id            int32
	CampEffective []int32
}

type CfgCharacterCampConfig struct {
	data map[int32]*CfgCharacterCamp
}

func NewCfgCharacterCampConfig() *CfgCharacterCampConfig {
	return &CfgCharacterCampConfig{
		data: make(map[int32]*CfgCharacterCamp),
	}
}

func (c *CfgCharacterCampConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgCharacterCamp)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgCharacterCamp.Id field error,value:", vId)
			return false
		}

		/* parse CampEffective field */
		vecCampEffective, _ := parse.GetFieldByName(uint32(i), "campEffective")
		if vecCampEffective != "" {
			arrayCampEffective := strings.Split(vecCampEffective, ",")
			for j := 0; j < len(arrayCampEffective); j++ {
				v, ret := String2Int32(arrayCampEffective[j])
				if !ret {
					glog.Error("Parse CfgCharacterCamp.CampEffective field error, value:", arrayCampEffective[j])
					return false
				}
				data.CampEffective = append(data.CampEffective, v)
			}
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgCharacterCampConfig) Clear() {
}

func (c *CfgCharacterCampConfig) Find(id int32) (*CfgCharacterCamp, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgCharacterCampConfig) GetAllData() map[int32]*CfgCharacterCamp {
	return c.data
}

func (c *CfgCharacterCampConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.CampEffective)
	}
}
