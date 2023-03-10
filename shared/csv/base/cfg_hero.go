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

type CfgHero struct {
	Id             int32
	Show           bool
	Skills         []int32
	CampSkills     []int32
	SkillItem      int32
	MainAttr       int32
	AttrRates      []int32
	Camp           int32
	SkinId         int32
	HeroIcon       string
	HeroStand      string
	HeroSpine      string
	HeroHead       string
	HeroData       string
	UnlockCost     []string
	PledgeUnlockLv []int32
	RecommnedChar  []int32
	PassiveId      int32
	AttributesId   string
}

type CfgHeroConfig struct {
	data map[int32]*CfgHero
}

func NewCfgHeroConfig() *CfgHeroConfig {
	return &CfgHeroConfig{
		data: make(map[int32]*CfgHero),
	}
}

func (c *CfgHeroConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgHero)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgHero.Id field error,value:", vId)
			return false
		}

		/* parse Show field */
		vShow, _ := parse.GetFieldByName(uint32(i), "show")
		var ShowRet bool
		data.Show, ShowRet = String2Bool(vShow)
		if !ShowRet {
			glog.Error("Parse CfgHero.Show field error,value:", vShow)
		}

		/* parse Skills field */
		vecSkills, _ := parse.GetFieldByName(uint32(i), "skills")
		if vecSkills != "" {
			arraySkills := strings.Split(vecSkills, ",")
			for j := 0; j < len(arraySkills); j++ {
				v, ret := String2Int32(arraySkills[j])
				if !ret {
					glog.Error("Parse CfgHero.Skills field error, value:", arraySkills[j])
					return false
				}
				data.Skills = append(data.Skills, v)
			}
		}

		/* parse CampSkills field */
		vecCampSkills, _ := parse.GetFieldByName(uint32(i), "campSkills")
		if vecCampSkills != "" {
			arrayCampSkills := strings.Split(vecCampSkills, ",")
			for j := 0; j < len(arrayCampSkills); j++ {
				v, ret := String2Int32(arrayCampSkills[j])
				if !ret {
					glog.Error("Parse CfgHero.CampSkills field error, value:", arrayCampSkills[j])
					return false
				}
				data.CampSkills = append(data.CampSkills, v)
			}
		}

		/* parse SkillItem field */
		vSkillItem, _ := parse.GetFieldByName(uint32(i), "skillItem")
		var SkillItemRet bool
		data.SkillItem, SkillItemRet = String2Int32(vSkillItem)
		if !SkillItemRet {
			glog.Error("Parse CfgHero.SkillItem field error,value:", vSkillItem)
			return false
		}

		/* parse MainAttr field */
		vMainAttr, _ := parse.GetFieldByName(uint32(i), "mainAttr")
		var MainAttrRet bool
		data.MainAttr, MainAttrRet = String2Int32(vMainAttr)
		if !MainAttrRet {
			glog.Error("Parse CfgHero.MainAttr field error,value:", vMainAttr)
			return false
		}

		/* parse AttrRates field */
		vecAttrRates, _ := parse.GetFieldByName(uint32(i), "attrRates")
		if vecAttrRates != "" {
			arrayAttrRates := strings.Split(vecAttrRates, ",")
			for j := 0; j < len(arrayAttrRates); j++ {
				v, ret := String2Int32(arrayAttrRates[j])
				if !ret {
					glog.Error("Parse CfgHero.AttrRates field error, value:", arrayAttrRates[j])
					return false
				}
				data.AttrRates = append(data.AttrRates, v)
			}
		}

		/* parse Camp field */
		vCamp, _ := parse.GetFieldByName(uint32(i), "camp")
		var CampRet bool
		data.Camp, CampRet = String2Int32(vCamp)
		if !CampRet {
			glog.Error("Parse CfgHero.Camp field error,value:", vCamp)
			return false
		}

		/* parse SkinId field */
		vSkinId, _ := parse.GetFieldByName(uint32(i), "skinId")
		var SkinIdRet bool
		data.SkinId, SkinIdRet = String2Int32(vSkinId)
		if !SkinIdRet {
			glog.Error("Parse CfgHero.SkinId field error,value:", vSkinId)
			return false
		}

		/* parse HeroIcon field */
		data.HeroIcon, _ = parse.GetFieldByName(uint32(i), "heroIcon")

		/* parse HeroStand field */
		data.HeroStand, _ = parse.GetFieldByName(uint32(i), "heroStand")

		/* parse HeroSpine field */
		data.HeroSpine, _ = parse.GetFieldByName(uint32(i), "heroSpine")

		/* parse HeroHead field */
		data.HeroHead, _ = parse.GetFieldByName(uint32(i), "heroHead")

		/* parse HeroData field */
		data.HeroData, _ = parse.GetFieldByName(uint32(i), "heroData")

		/* parse UnlockCost field */
		vecUnlockCost, _ := parse.GetFieldByName(uint32(i), "unlockCost")
		arrayUnlockCost := strings.Split(vecUnlockCost, ",")
		for j := 0; j < len(arrayUnlockCost); j++ {
			v := arrayUnlockCost[j]
			data.UnlockCost = append(data.UnlockCost, v)
		}

		/* parse PledgeUnlockLv field */
		vecPledgeUnlockLv, _ := parse.GetFieldByName(uint32(i), "pledgeUnlockLv")
		if vecPledgeUnlockLv != "" {
			arrayPledgeUnlockLv := strings.Split(vecPledgeUnlockLv, ",")
			for j := 0; j < len(arrayPledgeUnlockLv); j++ {
				v, ret := String2Int32(arrayPledgeUnlockLv[j])
				if !ret {
					glog.Error("Parse CfgHero.PledgeUnlockLv field error, value:", arrayPledgeUnlockLv[j])
					return false
				}
				data.PledgeUnlockLv = append(data.PledgeUnlockLv, v)
			}
		}

		/* parse RecommnedChar field */
		vecRecommnedChar, _ := parse.GetFieldByName(uint32(i), "recommnedChar")
		if vecRecommnedChar != "" {
			arrayRecommnedChar := strings.Split(vecRecommnedChar, ",")
			for j := 0; j < len(arrayRecommnedChar); j++ {
				v, ret := String2Int32(arrayRecommnedChar[j])
				if !ret {
					glog.Error("Parse CfgHero.RecommnedChar field error, value:", arrayRecommnedChar[j])
					return false
				}
				data.RecommnedChar = append(data.RecommnedChar, v)
			}
		}

		/* parse PassiveId field */
		vPassiveId, _ := parse.GetFieldByName(uint32(i), "passiveId")
		var PassiveIdRet bool
		data.PassiveId, PassiveIdRet = String2Int32(vPassiveId)
		if !PassiveIdRet {
			glog.Error("Parse CfgHero.PassiveId field error,value:", vPassiveId)
			return false
		}

		/* parse AttributesId field */
		data.AttributesId, _ = parse.GetFieldByName(uint32(i), "attributesId")

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgHeroConfig) Clear() {
}

func (c *CfgHeroConfig) Find(id int32) (*CfgHero, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgHeroConfig) GetAllData() map[int32]*CfgHero {
	return c.data
}

func (c *CfgHeroConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Show, ",", v.Skills, ",", v.CampSkills, ",", v.SkillItem, ",", v.MainAttr, ",", v.AttrRates, ",", v.Camp, ",", v.SkinId, ",", v.HeroIcon, ",", v.HeroStand, ",", v.HeroSpine, ",", v.HeroHead, ",", v.HeroData, ",", v.UnlockCost, ",", v.PledgeUnlockLv, ",", v.RecommnedChar, ",", v.PassiveId, ",", v.AttributesId)
	}
}
