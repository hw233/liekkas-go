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

type CfgChapterLevelAchievement struct {
	Id     int32
	Type   int32
	ParamX int32
	ParamY int32
	ParamZ int32
}

type CfgChapterLevelAchievementConfig struct {
	data map[int32]*CfgChapterLevelAchievement
}

func NewCfgChapterLevelAchievementConfig() *CfgChapterLevelAchievementConfig {
	return &CfgChapterLevelAchievementConfig{
		data: make(map[int32]*CfgChapterLevelAchievement),
	}
}

func (c *CfgChapterLevelAchievementConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgChapterLevelAchievement)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgChapterLevelAchievement.Id field error,value:", vId)
			return false
		}

		/* parse Type field */
		vType, _ := parse.GetFieldByName(uint32(i), "type")
		var TypeRet bool
		data.Type, TypeRet = String2Int32(vType)
		if !TypeRet {
			glog.Error("Parse CfgChapterLevelAchievement.Type field error,value:", vType)
			return false
		}

		/* parse ParamX field */
		vParamX, _ := parse.GetFieldByName(uint32(i), "paramX")
		var ParamXRet bool
		data.ParamX, ParamXRet = String2Int32(vParamX)
		if !ParamXRet {
			glog.Error("Parse CfgChapterLevelAchievement.ParamX field error,value:", vParamX)
			return false
		}

		/* parse ParamY field */
		vParamY, _ := parse.GetFieldByName(uint32(i), "paramY")
		var ParamYRet bool
		data.ParamY, ParamYRet = String2Int32(vParamY)
		if !ParamYRet {
			glog.Error("Parse CfgChapterLevelAchievement.ParamY field error,value:", vParamY)
			return false
		}

		/* parse ParamZ field */
		vParamZ, _ := parse.GetFieldByName(uint32(i), "paramZ")
		var ParamZRet bool
		data.ParamZ, ParamZRet = String2Int32(vParamZ)
		if !ParamZRet {
			glog.Error("Parse CfgChapterLevelAchievement.ParamZ field error,value:", vParamZ)
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

func (c *CfgChapterLevelAchievementConfig) Clear() {
}

func (c *CfgChapterLevelAchievementConfig) Find(id int32) (*CfgChapterLevelAchievement, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgChapterLevelAchievementConfig) GetAllData() map[int32]*CfgChapterLevelAchievement {
	return c.data
}

func (c *CfgChapterLevelAchievementConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Type, ",", v.ParamX, ",", v.ParamY, ",", v.ParamZ)
	}
}
