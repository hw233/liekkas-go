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

type CfgYggdrasilSimpleTalk struct {
	Id int32
}

type CfgYggdrasilSimpleTalkConfig struct {
	data map[int32]*CfgYggdrasilSimpleTalk
}

func NewCfgYggdrasilSimpleTalkConfig() *CfgYggdrasilSimpleTalkConfig {
	return &CfgYggdrasilSimpleTalkConfig{
		data: make(map[int32]*CfgYggdrasilSimpleTalk),
	}
}

func (c *CfgYggdrasilSimpleTalkConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgYggdrasilSimpleTalk)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgYggdrasilSimpleTalk.Id field error,value:", vId)
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

func (c *CfgYggdrasilSimpleTalkConfig) Clear() {
}

func (c *CfgYggdrasilSimpleTalkConfig) Find(id int32) (*CfgYggdrasilSimpleTalk, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgYggdrasilSimpleTalkConfig) GetAllData() map[int32]*CfgYggdrasilSimpleTalk {
	return c.data
}

func (c *CfgYggdrasilSimpleTalkConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id)
	}
}
