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

type CfgWorldItemData struct {
	Id             int32
	Name           string
	Rarity         int32
	UseType        int32
	UseParam       int32
	SellPrice      string
	EquipCareerLmt []int32
	EquipCampLmt   []int32
	GoldperExp     float64
	IsLocked       bool
	InitialStar    int32
	PassiveSkillID string
	PassiveAttriID string
	BreakItem      []int32
}

type CfgWorldItemDataConfig struct {
	data map[int32]*CfgWorldItemData
}

func NewCfgWorldItemDataConfig() *CfgWorldItemDataConfig {
	return &CfgWorldItemDataConfig{
		data: make(map[int32]*CfgWorldItemData),
	}
}

func (c *CfgWorldItemDataConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgWorldItemData)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgWorldItemData.Id field error,value:", vId)
			return false
		}

		/* parse Name field */
		data.Name, _ = parse.GetFieldByName(uint32(i), "name")

		/* parse Rarity field */
		vRarity, _ := parse.GetFieldByName(uint32(i), "rarity")
		var RarityRet bool
		data.Rarity, RarityRet = String2Int32(vRarity)
		if !RarityRet {
			glog.Error("Parse CfgWorldItemData.Rarity field error,value:", vRarity)
			return false
		}

		/* parse UseType field */
		vUseType, _ := parse.GetFieldByName(uint32(i), "useType")
		var UseTypeRet bool
		data.UseType, UseTypeRet = String2Int32(vUseType)
		if !UseTypeRet {
			glog.Error("Parse CfgWorldItemData.UseType field error,value:", vUseType)
			return false
		}

		/* parse UseParam field */
		vUseParam, _ := parse.GetFieldByName(uint32(i), "useParam")
		var UseParamRet bool
		data.UseParam, UseParamRet = String2Int32(vUseParam)
		if !UseParamRet {
			glog.Error("Parse CfgWorldItemData.UseParam field error,value:", vUseParam)
			return false
		}

		/* parse SellPrice field */
		data.SellPrice, _ = parse.GetFieldByName(uint32(i), "sellPrice")

		/* parse EquipCareerLmt field */
		vecEquipCareerLmt, _ := parse.GetFieldByName(uint32(i), "equipCareerLmt")
		if vecEquipCareerLmt != "" {
			arrayEquipCareerLmt := strings.Split(vecEquipCareerLmt, ",")
			for j := 0; j < len(arrayEquipCareerLmt); j++ {
				v, ret := String2Int32(arrayEquipCareerLmt[j])
				if !ret {
					glog.Error("Parse CfgWorldItemData.EquipCareerLmt field error, value:", arrayEquipCareerLmt[j])
					return false
				}
				data.EquipCareerLmt = append(data.EquipCareerLmt, v)
			}
		}

		/* parse EquipCampLmt field */
		vecEquipCampLmt, _ := parse.GetFieldByName(uint32(i), "equipCampLmt")
		if vecEquipCampLmt != "" {
			arrayEquipCampLmt := strings.Split(vecEquipCampLmt, ",")
			for j := 0; j < len(arrayEquipCampLmt); j++ {
				v, ret := String2Int32(arrayEquipCampLmt[j])
				if !ret {
					glog.Error("Parse CfgWorldItemData.EquipCampLmt field error, value:", arrayEquipCampLmt[j])
					return false
				}
				data.EquipCampLmt = append(data.EquipCampLmt, v)
			}
		}

		/* parse GoldperExp field */
		vGoldperExp, _ := parse.GetFieldByName(uint32(i), "goldperExp")
		var GoldperExpRet bool
		data.GoldperExp, GoldperExpRet = String2Float(vGoldperExp)
		if !GoldperExpRet {
			glog.Error("Parse CfgWorldItemData.GoldperExp field error,value:", vGoldperExp)
		}

		/* parse IsLocked field */
		vIsLocked, _ := parse.GetFieldByName(uint32(i), "isLocked")
		var IsLockedRet bool
		data.IsLocked, IsLockedRet = String2Bool(vIsLocked)
		if !IsLockedRet {
			glog.Error("Parse CfgWorldItemData.IsLocked field error,value:", vIsLocked)
		}

		/* parse InitialStar field */
		vInitialStar, _ := parse.GetFieldByName(uint32(i), "initialStar")
		var InitialStarRet bool
		data.InitialStar, InitialStarRet = String2Int32(vInitialStar)
		if !InitialStarRet {
			glog.Error("Parse CfgWorldItemData.InitialStar field error,value:", vInitialStar)
			return false
		}

		/* parse PassiveSkillID field */
		data.PassiveSkillID, _ = parse.GetFieldByName(uint32(i), "passiveSkillID")

		/* parse PassiveAttriID field */
		data.PassiveAttriID, _ = parse.GetFieldByName(uint32(i), "passiveAttriID")

		/* parse BreakItem field */
		vecBreakItem, _ := parse.GetFieldByName(uint32(i), "breakItem")
		if vecBreakItem != "" {
			arrayBreakItem := strings.Split(vecBreakItem, ",")
			for j := 0; j < len(arrayBreakItem); j++ {
				v, ret := String2Int32(arrayBreakItem[j])
				if !ret {
					glog.Error("Parse CfgWorldItemData.BreakItem field error, value:", arrayBreakItem[j])
					return false
				}
				data.BreakItem = append(data.BreakItem, v)
			}
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgWorldItemDataConfig) Clear() {
}

func (c *CfgWorldItemDataConfig) Find(id int32) (*CfgWorldItemData, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgWorldItemDataConfig) GetAllData() map[int32]*CfgWorldItemData {
	return c.data
}

func (c *CfgWorldItemDataConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Name, ",", v.Rarity, ",", v.UseType, ",", v.UseParam, ",", v.SellPrice, ",", v.EquipCareerLmt, ",", v.EquipCampLmt, ",", v.GoldperExp, ",", v.IsLocked, ",", v.InitialStar, ",", v.PassiveSkillID, ",", v.PassiveAttriID, ",", v.BreakItem)
	}
}
