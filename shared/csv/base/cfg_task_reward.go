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

type CfgTaskReward struct {
	Id        int32
	Module    int32
	Condition int32
	DropId    int32
	Reward    []string
}

type CfgTaskRewardConfig struct {
	data map[int32]*CfgTaskReward
}

func NewCfgTaskRewardConfig() *CfgTaskRewardConfig {
	return &CfgTaskRewardConfig{
		data: make(map[int32]*CfgTaskReward),
	}
}

func (c *CfgTaskRewardConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgTaskReward)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgTaskReward.Id field error,value:", vId)
			return false
		}

		/* parse Module field */
		vModule, _ := parse.GetFieldByName(uint32(i), "module")
		var ModuleRet bool
		data.Module, ModuleRet = String2Int32(vModule)
		if !ModuleRet {
			glog.Error("Parse CfgTaskReward.Module field error,value:", vModule)
			return false
		}

		/* parse Condition field */
		vCondition, _ := parse.GetFieldByName(uint32(i), "condition")
		var ConditionRet bool
		data.Condition, ConditionRet = String2Int32(vCondition)
		if !ConditionRet {
			glog.Error("Parse CfgTaskReward.Condition field error,value:", vCondition)
			return false
		}

		/* parse DropId field */
		vDropId, _ := parse.GetFieldByName(uint32(i), "dropId")
		var DropIdRet bool
		data.DropId, DropIdRet = String2Int32(vDropId)
		if !DropIdRet {
			glog.Error("Parse CfgTaskReward.DropId field error,value:", vDropId)
			return false
		}

		/* parse Reward field */
		vecReward, _ := parse.GetFieldByName(uint32(i), "reward")
		arrayReward := strings.Split(vecReward, ",")
		for j := 0; j < len(arrayReward); j++ {
			v := arrayReward[j]
			data.Reward = append(data.Reward, v)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgTaskRewardConfig) Clear() {
}

func (c *CfgTaskRewardConfig) Find(id int32) (*CfgTaskReward, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgTaskRewardConfig) GetAllData() map[int32]*CfgTaskReward {
	return c.data
}

func (c *CfgTaskRewardConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Module, ",", v.Condition, ",", v.DropId, ",", v.Reward)
	}
}
