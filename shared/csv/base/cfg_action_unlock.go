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

type CfgActionUnlock struct {
	Id              int32
	UnlockCondition []string
}

type CfgActionUnlockConfig struct {
	data map[int32]*CfgActionUnlock
}

func NewCfgActionUnlockConfig() *CfgActionUnlockConfig {
	return &CfgActionUnlockConfig{
		data: make(map[int32]*CfgActionUnlock),
	}
}

func (c *CfgActionUnlockConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgActionUnlock)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgActionUnlock.Id field error,value:", vId)
			return false
		}

		/* parse UnlockCondition field */
		vecUnlockCondition, _ := parse.GetFieldByName(uint32(i), "unlockCondition")
		arrayUnlockCondition := strings.Split(vecUnlockCondition, ",")
		for j := 0; j < len(arrayUnlockCondition); j++ {
			v := arrayUnlockCondition[j]
			data.UnlockCondition = append(data.UnlockCondition, v)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgActionUnlockConfig) Clear() {
}

func (c *CfgActionUnlockConfig) Find(id int32) (*CfgActionUnlock, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgActionUnlockConfig) GetAllData() map[int32]*CfgActionUnlock {
	return c.data
}

func (c *CfgActionUnlockConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.UnlockCondition)
	}
}
