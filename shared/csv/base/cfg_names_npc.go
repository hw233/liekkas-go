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

type CfgNamesNpc struct {
	Id    int32
	Part1 string
	Part2 string
	Part3 string
}

type CfgNamesNpcConfig struct {
	data map[int32]*CfgNamesNpc
}

func NewCfgNamesNpcConfig() *CfgNamesNpcConfig {
	return &CfgNamesNpcConfig{
		data: make(map[int32]*CfgNamesNpc),
	}
}

func (c *CfgNamesNpcConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgNamesNpc)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgNamesNpc.Id field error,value:", vId)
			return false
		}

		/* parse Part1 field */
		data.Part1, _ = parse.GetFieldByName(uint32(i), "part1")

		/* parse Part2 field */
		data.Part2, _ = parse.GetFieldByName(uint32(i), "part2")

		/* parse Part3 field */
		data.Part3, _ = parse.GetFieldByName(uint32(i), "part3")

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgNamesNpcConfig) Clear() {
}

func (c *CfgNamesNpcConfig) Find(id int32) (*CfgNamesNpc, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgNamesNpcConfig) GetAllData() map[int32]*CfgNamesNpc {
	return c.data
}

func (c *CfgNamesNpcConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Part1, ",", v.Part2, ",", v.Part3)
	}
}
