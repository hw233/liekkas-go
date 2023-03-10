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

type CfgMailTemplate struct {
	Id         int32
	ExpireDays int32
	Publisher  string
	Title      string
	TitleIcon  string
	Content    string
}

type CfgMailTemplateConfig struct {
	data map[int32]*CfgMailTemplate
}

func NewCfgMailTemplateConfig() *CfgMailTemplateConfig {
	return &CfgMailTemplateConfig{
		data: make(map[int32]*CfgMailTemplate),
	}
}

func (c *CfgMailTemplateConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgMailTemplate)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgMailTemplate.Id field error,value:", vId)
			return false
		}

		/* parse ExpireDays field */
		vExpireDays, _ := parse.GetFieldByName(uint32(i), "expireDays")
		var ExpireDaysRet bool
		data.ExpireDays, ExpireDaysRet = String2Int32(vExpireDays)
		if !ExpireDaysRet {
			glog.Error("Parse CfgMailTemplate.ExpireDays field error,value:", vExpireDays)
			return false
		}

		/* parse Publisher field */
		data.Publisher, _ = parse.GetFieldByName(uint32(i), "publisher")

		/* parse Title field */
		data.Title, _ = parse.GetFieldByName(uint32(i), "title")

		/* parse TitleIcon field */
		data.TitleIcon, _ = parse.GetFieldByName(uint32(i), "titleIcon")

		/* parse Content field */
		data.Content, _ = parse.GetFieldByName(uint32(i), "content")

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgMailTemplateConfig) Clear() {
}

func (c *CfgMailTemplateConfig) Find(id int32) (*CfgMailTemplate, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgMailTemplateConfig) GetAllData() map[int32]*CfgMailTemplate {
	return c.data
}

func (c *CfgMailTemplateConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.ExpireDays, ",", v.Publisher, ",", v.Title, ",", v.TitleIcon, ",", v.Content)
	}
}
