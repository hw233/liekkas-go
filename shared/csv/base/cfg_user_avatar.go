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

type CfgUserAvatar struct {
	Id         int32
	AvatarType int32
}

type CfgUserAvatarConfig struct {
	data map[int32]*CfgUserAvatar
}

func NewCfgUserAvatarConfig() *CfgUserAvatarConfig {
	return &CfgUserAvatarConfig{
		data: make(map[int32]*CfgUserAvatar),
	}
}

func (c *CfgUserAvatarConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgUserAvatar)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgUserAvatar.Id field error,value:", vId)
			return false
		}

		/* parse AvatarType field */
		vAvatarType, _ := parse.GetFieldByName(uint32(i), "avatarType")
		var AvatarTypeRet bool
		data.AvatarType, AvatarTypeRet = String2Int32(vAvatarType)
		if !AvatarTypeRet {
			glog.Error("Parse CfgUserAvatar.AvatarType field error,value:", vAvatarType)
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

func (c *CfgUserAvatarConfig) Clear() {
}

func (c *CfgUserAvatarConfig) Find(id int32) (*CfgUserAvatar, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgUserAvatarConfig) GetAllData() map[int32]*CfgUserAvatar {
	return c.data
}

func (c *CfgUserAvatarConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.AvatarType)
	}
}
