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

type CfgFieldSkill struct {
	Id       int32
	Comment  string
	EventsID []int32
}

type CfgFieldSkillConfig struct {
	data map[int32]*CfgFieldSkill
}

func NewCfgFieldSkillConfig() *CfgFieldSkillConfig {
	return &CfgFieldSkillConfig{
		data: make(map[int32]*CfgFieldSkill),
	}
}

func (c *CfgFieldSkillConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgFieldSkill)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgFieldSkill.Id field error,value:", vId)
			return false
		}

		/* parse Comment field */
		data.Comment, _ = parse.GetFieldByName(uint32(i), "comment")

		/* parse EventsID field */
		vecEventsID, _ := parse.GetFieldByName(uint32(i), "eventsID")
		if vecEventsID != "" {
			arrayEventsID := strings.Split(vecEventsID, ",")
			for j := 0; j < len(arrayEventsID); j++ {
				v, ret := String2Int32(arrayEventsID[j])
				if !ret {
					glog.Error("Parse CfgFieldSkill.EventsID field error, value:", arrayEventsID[j])
					return false
				}
				data.EventsID = append(data.EventsID, v)
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

func (c *CfgFieldSkillConfig) Clear() {
}

func (c *CfgFieldSkillConfig) Find(id int32) (*CfgFieldSkill, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgFieldSkillConfig) GetAllData() map[int32]*CfgFieldSkill {
	return c.data
}

func (c *CfgFieldSkillConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Comment, ",", v.EventsID)
	}
}
