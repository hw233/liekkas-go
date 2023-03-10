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

type CfgManual struct {
	Id        int32
	Type      int32
	RelatedId int32
	CampId    int32
	DropId    int32
	Version   int32
}

type CfgManualConfig struct {
	data map[int32]*CfgManual
}

func NewCfgManualConfig() *CfgManualConfig {
	return &CfgManualConfig{
		data: make(map[int32]*CfgManual),
	}
}

func (c *CfgManualConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgManual)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgManual.Id field error,value:", vId)
			return false
		}

		/* parse Type field */
		vType, _ := parse.GetFieldByName(uint32(i), "type")
		var TypeRet bool
		data.Type, TypeRet = String2Int32(vType)
		if !TypeRet {
			glog.Error("Parse CfgManual.Type field error,value:", vType)
			return false
		}

		/* parse RelatedId field */
		vRelatedId, _ := parse.GetFieldByName(uint32(i), "relatedId")
		var RelatedIdRet bool
		data.RelatedId, RelatedIdRet = String2Int32(vRelatedId)
		if !RelatedIdRet {
			glog.Error("Parse CfgManual.RelatedId field error,value:", vRelatedId)
			return false
		}

		/* parse CampId field */
		vCampId, _ := parse.GetFieldByName(uint32(i), "campId")
		var CampIdRet bool
		data.CampId, CampIdRet = String2Int32(vCampId)
		if !CampIdRet {
			glog.Error("Parse CfgManual.CampId field error,value:", vCampId)
			return false
		}

		/* parse DropId field */
		vDropId, _ := parse.GetFieldByName(uint32(i), "dropId")
		var DropIdRet bool
		data.DropId, DropIdRet = String2Int32(vDropId)
		if !DropIdRet {
			glog.Error("Parse CfgManual.DropId field error,value:", vDropId)
			return false
		}

		/* parse Version field */
		vVersion, _ := parse.GetFieldByName(uint32(i), "version")
		var VersionRet bool
		data.Version, VersionRet = String2Int32(vVersion)
		if !VersionRet {
			glog.Error("Parse CfgManual.Version field error,value:", vVersion)
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

func (c *CfgManualConfig) Clear() {
}

func (c *CfgManualConfig) Find(id int32) (*CfgManual, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgManualConfig) GetAllData() map[int32]*CfgManual {
	return c.data
}

func (c *CfgManualConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Type, ",", v.RelatedId, ",", v.CampId, ",", v.DropId, ",", v.Version)
	}
}
