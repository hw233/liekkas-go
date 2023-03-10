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

type CfgManualLevelAttributeData struct {
	Id     int32
	CampId int32
	Lv     int32
	Stars  int32
	HpMax  int32
	PhyAtk int32
	MagAtk int32
	PhyDfs int32
	MagDfs int32
}

type CfgManualLevelAttributeDataConfig struct {
	data map[int32]*CfgManualLevelAttributeData
}

func NewCfgManualLevelAttributeDataConfig() *CfgManualLevelAttributeDataConfig {
	return &CfgManualLevelAttributeDataConfig{
		data: make(map[int32]*CfgManualLevelAttributeData),
	}
}

func (c *CfgManualLevelAttributeDataConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgManualLevelAttributeData)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgManualLevelAttributeData.Id field error,value:", vId)
			return false
		}

		/* parse CampId field */
		vCampId, _ := parse.GetFieldByName(uint32(i), "campId")
		var CampIdRet bool
		data.CampId, CampIdRet = String2Int32(vCampId)
		if !CampIdRet {
			glog.Error("Parse CfgManualLevelAttributeData.CampId field error,value:", vCampId)
			return false
		}

		/* parse Lv field */
		vLv, _ := parse.GetFieldByName(uint32(i), "lv")
		var LvRet bool
		data.Lv, LvRet = String2Int32(vLv)
		if !LvRet {
			glog.Error("Parse CfgManualLevelAttributeData.Lv field error,value:", vLv)
			return false
		}

		/* parse Stars field */
		vStars, _ := parse.GetFieldByName(uint32(i), "stars")
		var StarsRet bool
		data.Stars, StarsRet = String2Int32(vStars)
		if !StarsRet {
			glog.Error("Parse CfgManualLevelAttributeData.Stars field error,value:", vStars)
			return false
		}

		/* parse HpMax field */
		vHpMax, _ := parse.GetFieldByName(uint32(i), "hpMax")
		var HpMaxRet bool
		data.HpMax, HpMaxRet = String2Int32(vHpMax)
		if !HpMaxRet {
			glog.Error("Parse CfgManualLevelAttributeData.HpMax field error,value:", vHpMax)
			return false
		}

		/* parse PhyAtk field */
		vPhyAtk, _ := parse.GetFieldByName(uint32(i), "phyAtk")
		var PhyAtkRet bool
		data.PhyAtk, PhyAtkRet = String2Int32(vPhyAtk)
		if !PhyAtkRet {
			glog.Error("Parse CfgManualLevelAttributeData.PhyAtk field error,value:", vPhyAtk)
			return false
		}

		/* parse MagAtk field */
		vMagAtk, _ := parse.GetFieldByName(uint32(i), "magAtk")
		var MagAtkRet bool
		data.MagAtk, MagAtkRet = String2Int32(vMagAtk)
		if !MagAtkRet {
			glog.Error("Parse CfgManualLevelAttributeData.MagAtk field error,value:", vMagAtk)
			return false
		}

		/* parse PhyDfs field */
		vPhyDfs, _ := parse.GetFieldByName(uint32(i), "phyDfs")
		var PhyDfsRet bool
		data.PhyDfs, PhyDfsRet = String2Int32(vPhyDfs)
		if !PhyDfsRet {
			glog.Error("Parse CfgManualLevelAttributeData.PhyDfs field error,value:", vPhyDfs)
			return false
		}

		/* parse MagDfs field */
		vMagDfs, _ := parse.GetFieldByName(uint32(i), "magDfs")
		var MagDfsRet bool
		data.MagDfs, MagDfsRet = String2Int32(vMagDfs)
		if !MagDfsRet {
			glog.Error("Parse CfgManualLevelAttributeData.MagDfs field error,value:", vMagDfs)
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

func (c *CfgManualLevelAttributeDataConfig) Clear() {
}

func (c *CfgManualLevelAttributeDataConfig) Find(id int32) (*CfgManualLevelAttributeData, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgManualLevelAttributeDataConfig) GetAllData() map[int32]*CfgManualLevelAttributeData {
	return c.data
}

func (c *CfgManualLevelAttributeDataConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.CampId, ",", v.Lv, ",", v.Stars, ",", v.HpMax, ",", v.PhyAtk, ",", v.MagAtk, ",", v.PhyDfs, ",", v.MagDfs)
	}
}
