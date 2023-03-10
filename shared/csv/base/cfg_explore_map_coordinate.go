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

type CfgExploreMapCoordinate struct {
	Id        int32
	ChapterId int32
	Name      string
	Type      int32
	TypeParam int32
}

type CfgExploreMapCoordinateConfig struct {
	data map[int32]*CfgExploreMapCoordinate
}

func NewCfgExploreMapCoordinateConfig() *CfgExploreMapCoordinateConfig {
	return &CfgExploreMapCoordinateConfig{
		data: make(map[int32]*CfgExploreMapCoordinate),
	}
}

func (c *CfgExploreMapCoordinateConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgExploreMapCoordinate)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgExploreMapCoordinate.Id field error,value:", vId)
			return false
		}

		/* parse ChapterId field */
		vChapterId, _ := parse.GetFieldByName(uint32(i), "chapterId")
		var ChapterIdRet bool
		data.ChapterId, ChapterIdRet = String2Int32(vChapterId)
		if !ChapterIdRet {
			glog.Error("Parse CfgExploreMapCoordinate.ChapterId field error,value:", vChapterId)
			return false
		}

		/* parse Name field */
		data.Name, _ = parse.GetFieldByName(uint32(i), "name")

		/* parse Type field */
		vType, _ := parse.GetFieldByName(uint32(i), "type")
		var TypeRet bool
		data.Type, TypeRet = String2Int32(vType)
		if !TypeRet {
			glog.Error("Parse CfgExploreMapCoordinate.Type field error,value:", vType)
			return false
		}

		/* parse TypeParam field */
		vTypeParam, _ := parse.GetFieldByName(uint32(i), "typeParam")
		var TypeParamRet bool
		data.TypeParam, TypeParamRet = String2Int32(vTypeParam)
		if !TypeParamRet {
			glog.Error("Parse CfgExploreMapCoordinate.TypeParam field error,value:", vTypeParam)
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

func (c *CfgExploreMapCoordinateConfig) Clear() {
}

func (c *CfgExploreMapCoordinateConfig) Find(id int32) (*CfgExploreMapCoordinate, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgExploreMapCoordinateConfig) GetAllData() map[int32]*CfgExploreMapCoordinate {
	return c.data
}

func (c *CfgExploreMapCoordinateConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.ChapterId, ",", v.Name, ",", v.Type, ",", v.TypeParam)
	}
}
