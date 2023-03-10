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

type CfgVn struct {
	Id      int32
	Commit  string
	Address string
	Reward1 []string
}

type CfgVnConfig struct {
	data map[int32]*CfgVn
}

func NewCfgVnConfig() *CfgVnConfig {
	return &CfgVnConfig{
		data: make(map[int32]*CfgVn),
	}
}

func (c *CfgVnConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgVn)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgVn.Id field error,value:", vId)
			return false
		}

		/* parse Commit field */
		data.Commit, _ = parse.GetFieldByName(uint32(i), "commit")

		/* parse Address field */
		data.Address, _ = parse.GetFieldByName(uint32(i), "address")

		/* parse Reward1 field */
		vecReward1, _ := parse.GetFieldByName(uint32(i), "reward1")
		arrayReward1 := strings.Split(vecReward1, ",")
		for j := 0; j < len(arrayReward1); j++ {
			v := arrayReward1[j]
			data.Reward1 = append(data.Reward1, v)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgVnConfig) Clear() {
}

func (c *CfgVnConfig) Find(id int32) (*CfgVn, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgVnConfig) GetAllData() map[int32]*CfgVn {
	return c.data
}

func (c *CfgVnConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Commit, ",", v.Address, ",", v.Reward1)
	}
}
