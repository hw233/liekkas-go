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

type CfgSound struct {
	Id     int32
	Unlock []string
}

type CfgSoundConfig struct {
	data map[int32]*CfgSound
}

func NewCfgSoundConfig() *CfgSoundConfig {
	return &CfgSoundConfig{
		data: make(map[int32]*CfgSound),
	}
}

func (c *CfgSoundConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgSound)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgSound.Id field error,value:", vId)
			return false
		}

		/* parse Unlock field */
		vecUnlock, _ := parse.GetFieldByName(uint32(i), "unlock")
		arrayUnlock := strings.Split(vecUnlock, ",")
		for j := 0; j < len(arrayUnlock); j++ {
			v := arrayUnlock[j]
			data.Unlock = append(data.Unlock, v)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgSoundConfig) Clear() {
}

func (c *CfgSoundConfig) Find(id int32) (*CfgSound, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgSoundConfig) GetAllData() map[int32]*CfgSound {
	return c.data
}

func (c *CfgSoundConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Unlock)
	}
}
