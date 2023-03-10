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

type CfgMaskWord struct {
	Id   int32
	Word string
}

type CfgMaskWordConfig struct {
	data map[int32]*CfgMaskWord
}

func NewCfgMaskWordConfig() *CfgMaskWordConfig {
	return &CfgMaskWordConfig{
		data: make(map[int32]*CfgMaskWord),
	}
}

func (c *CfgMaskWordConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgMaskWord)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgMaskWord.Id field error,value:", vId)
			return false
		}

		/* parse Word field */
		data.Word, _ = parse.GetFieldByName(uint32(i), "word")

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgMaskWordConfig) Clear() {
}

func (c *CfgMaskWordConfig) Find(id int32) (*CfgMaskWord, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgMaskWordConfig) GetAllData() map[int32]*CfgMaskWord {
	return c.data
}

func (c *CfgMaskWordConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Word)
	}
}
