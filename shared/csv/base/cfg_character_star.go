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

type CfgCharacterStar struct {
	Id                  int32
	CharID              int32
	Star                int32
	HpRatio             float64
	PhyAtkRatio         float64
	MagAtkRatio         float64
	PhyDfsRatio         float64
	MagDfsRatio         float64
	HpMax               int32
	PhyAtk              int32
	MagAtk              int32
	PhyDfs              int32
	MagDfs              int32
	RarityUp            int32
	Cost                []string
	BuildProducePercent int32
}

type CfgCharacterStarConfig struct {
	data map[int32]*CfgCharacterStar
}

func NewCfgCharacterStarConfig() *CfgCharacterStarConfig {
	return &CfgCharacterStarConfig{
		data: make(map[int32]*CfgCharacterStar),
	}
}

func (c *CfgCharacterStarConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgCharacterStar)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgCharacterStar.Id field error,value:", vId)
			return false
		}

		/* parse CharID field */
		vCharID, _ := parse.GetFieldByName(uint32(i), "charID")
		var CharIDRet bool
		data.CharID, CharIDRet = String2Int32(vCharID)
		if !CharIDRet {
			glog.Error("Parse CfgCharacterStar.CharID field error,value:", vCharID)
			return false
		}

		/* parse Star field */
		vStar, _ := parse.GetFieldByName(uint32(i), "star")
		var StarRet bool
		data.Star, StarRet = String2Int32(vStar)
		if !StarRet {
			glog.Error("Parse CfgCharacterStar.Star field error,value:", vStar)
			return false
		}

		/* parse HpRatio field */
		vHpRatio, _ := parse.GetFieldByName(uint32(i), "hpRatio")
		var HpRatioRet bool
		data.HpRatio, HpRatioRet = String2Float(vHpRatio)
		if !HpRatioRet {
			glog.Error("Parse CfgCharacterStar.HpRatio field error,value:", vHpRatio)
		}

		/* parse PhyAtkRatio field */
		vPhyAtkRatio, _ := parse.GetFieldByName(uint32(i), "phyAtkRatio")
		var PhyAtkRatioRet bool
		data.PhyAtkRatio, PhyAtkRatioRet = String2Float(vPhyAtkRatio)
		if !PhyAtkRatioRet {
			glog.Error("Parse CfgCharacterStar.PhyAtkRatio field error,value:", vPhyAtkRatio)
		}

		/* parse MagAtkRatio field */
		vMagAtkRatio, _ := parse.GetFieldByName(uint32(i), "magAtkRatio")
		var MagAtkRatioRet bool
		data.MagAtkRatio, MagAtkRatioRet = String2Float(vMagAtkRatio)
		if !MagAtkRatioRet {
			glog.Error("Parse CfgCharacterStar.MagAtkRatio field error,value:", vMagAtkRatio)
		}

		/* parse PhyDfsRatio field */
		vPhyDfsRatio, _ := parse.GetFieldByName(uint32(i), "phyDfsRatio")
		var PhyDfsRatioRet bool
		data.PhyDfsRatio, PhyDfsRatioRet = String2Float(vPhyDfsRatio)
		if !PhyDfsRatioRet {
			glog.Error("Parse CfgCharacterStar.PhyDfsRatio field error,value:", vPhyDfsRatio)
		}

		/* parse MagDfsRatio field */
		vMagDfsRatio, _ := parse.GetFieldByName(uint32(i), "magDfsRatio")
		var MagDfsRatioRet bool
		data.MagDfsRatio, MagDfsRatioRet = String2Float(vMagDfsRatio)
		if !MagDfsRatioRet {
			glog.Error("Parse CfgCharacterStar.MagDfsRatio field error,value:", vMagDfsRatio)
		}

		/* parse HpMax field */
		vHpMax, _ := parse.GetFieldByName(uint32(i), "hpMax")
		var HpMaxRet bool
		data.HpMax, HpMaxRet = String2Int32(vHpMax)
		if !HpMaxRet {
			glog.Error("Parse CfgCharacterStar.HpMax field error,value:", vHpMax)
			return false
		}

		/* parse PhyAtk field */
		vPhyAtk, _ := parse.GetFieldByName(uint32(i), "phyAtk")
		var PhyAtkRet bool
		data.PhyAtk, PhyAtkRet = String2Int32(vPhyAtk)
		if !PhyAtkRet {
			glog.Error("Parse CfgCharacterStar.PhyAtk field error,value:", vPhyAtk)
			return false
		}

		/* parse MagAtk field */
		vMagAtk, _ := parse.GetFieldByName(uint32(i), "magAtk")
		var MagAtkRet bool
		data.MagAtk, MagAtkRet = String2Int32(vMagAtk)
		if !MagAtkRet {
			glog.Error("Parse CfgCharacterStar.MagAtk field error,value:", vMagAtk)
			return false
		}

		/* parse PhyDfs field */
		vPhyDfs, _ := parse.GetFieldByName(uint32(i), "phyDfs")
		var PhyDfsRet bool
		data.PhyDfs, PhyDfsRet = String2Int32(vPhyDfs)
		if !PhyDfsRet {
			glog.Error("Parse CfgCharacterStar.PhyDfs field error,value:", vPhyDfs)
			return false
		}

		/* parse MagDfs field */
		vMagDfs, _ := parse.GetFieldByName(uint32(i), "magDfs")
		var MagDfsRet bool
		data.MagDfs, MagDfsRet = String2Int32(vMagDfs)
		if !MagDfsRet {
			glog.Error("Parse CfgCharacterStar.MagDfs field error,value:", vMagDfs)
			return false
		}

		/* parse RarityUp field */
		vRarityUp, _ := parse.GetFieldByName(uint32(i), "rarityUp")
		var RarityUpRet bool
		data.RarityUp, RarityUpRet = String2Int32(vRarityUp)
		if !RarityUpRet {
			glog.Error("Parse CfgCharacterStar.RarityUp field error,value:", vRarityUp)
			return false
		}

		/* parse Cost field */
		vecCost, _ := parse.GetFieldByName(uint32(i), "cost")
		arrayCost := strings.Split(vecCost, ",")
		for j := 0; j < len(arrayCost); j++ {
			v := arrayCost[j]
			data.Cost = append(data.Cost, v)
		}

		/* parse BuildProducePercent field */
		vBuildProducePercent, _ := parse.GetFieldByName(uint32(i), "buildProducePercent")
		var BuildProducePercentRet bool
		data.BuildProducePercent, BuildProducePercentRet = String2Int32(vBuildProducePercent)
		if !BuildProducePercentRet {
			glog.Error("Parse CfgCharacterStar.BuildProducePercent field error,value:", vBuildProducePercent)
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

func (c *CfgCharacterStarConfig) Clear() {
}

func (c *CfgCharacterStarConfig) Find(id int32) (*CfgCharacterStar, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgCharacterStarConfig) GetAllData() map[int32]*CfgCharacterStar {
	return c.data
}

func (c *CfgCharacterStarConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.CharID, ",", v.Star, ",", v.HpRatio, ",", v.PhyAtkRatio, ",", v.MagAtkRatio, ",", v.PhyDfsRatio, ",", v.MagDfsRatio, ",", v.HpMax, ",", v.PhyAtk, ",", v.MagAtk, ",", v.PhyDfs, ",", v.MagDfs, ",", v.RarityUp, ",", v.Cost, ",", v.BuildProducePercent)
	}
}
