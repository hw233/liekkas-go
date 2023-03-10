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

type CfgGlobalString struct {
	Id    int32
	Key   string
	Value string
}

type CfgGlobalStringConfig struct {
	data map[int32]*CfgGlobalString
}

func NewCfgGlobalStringConfig() *CfgGlobalStringConfig {
	return &CfgGlobalStringConfig{
		data: make(map[int32]*CfgGlobalString),
	}
}

func (c *CfgGlobalStringConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgGlobalString)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgGlobalString.Id field error,value:", vId)
			return false
		}

		/* parse Key field */
		data.Key, _ = parse.GetFieldByName(uint32(i), "key")

		/* parse Value field */
		data.Value, _ = parse.GetFieldByName(uint32(i), "value")

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgGlobalStringConfig) Clear() {
}

func (c *CfgGlobalStringConfig) Find(id int32) (*CfgGlobalString, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgGlobalStringConfig) GetAllData() map[int32]*CfgGlobalString {
	return c.data
}

func (c *CfgGlobalStringConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Key, ",", v.Value)
	}
}
