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

type CfgManualCharacterVoice struct {
	Id          int32
	CharacterID int32
	Unlock      []string
}

type CfgManualCharacterVoiceConfig struct {
	data map[int32]*CfgManualCharacterVoice
}

func NewCfgManualCharacterVoiceConfig() *CfgManualCharacterVoiceConfig {
	return &CfgManualCharacterVoiceConfig{
		data: make(map[int32]*CfgManualCharacterVoice),
	}
}

func (c *CfgManualCharacterVoiceConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgManualCharacterVoice)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgManualCharacterVoice.Id field error,value:", vId)
			return false
		}

		/* parse CharacterID field */
		vCharacterID, _ := parse.GetFieldByName(uint32(i), "characterID")
		var CharacterIDRet bool
		data.CharacterID, CharacterIDRet = String2Int32(vCharacterID)
		if !CharacterIDRet {
			glog.Error("Parse CfgManualCharacterVoice.CharacterID field error,value:", vCharacterID)
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

func (c *CfgManualCharacterVoiceConfig) Clear() {
}

func (c *CfgManualCharacterVoiceConfig) Find(id int32) (*CfgManualCharacterVoice, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgManualCharacterVoiceConfig) GetAllData() map[int32]*CfgManualCharacterVoice {
	return c.data
}

func (c *CfgManualCharacterVoiceConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.CharacterID, ",", v.Unlock)
	}
}
