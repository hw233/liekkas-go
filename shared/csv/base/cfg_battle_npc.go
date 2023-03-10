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

type CfgBattleNpc struct {
	Id        int32
	CharaID   int32
	Level     int32
	Star      int32
	Stage     int32
	Equip1    []int32
	Equip2    []int32
	Equip3    []int32
	Equip4    []int32
	SkillLv   []string
	WorldItem int32
}

type CfgBattleNpcConfig struct {
	data map[int32]*CfgBattleNpc
}

func NewCfgBattleNpcConfig() *CfgBattleNpcConfig {
	return &CfgBattleNpcConfig{
		data: make(map[int32]*CfgBattleNpc),
	}
}

func (c *CfgBattleNpcConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgBattleNpc)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgBattleNpc.Id field error,value:", vId)
			return false
		}

		/* parse CharaID field */
		vCharaID, _ := parse.GetFieldByName(uint32(i), "charaID")
		var CharaIDRet bool
		data.CharaID, CharaIDRet = String2Int32(vCharaID)
		if !CharaIDRet {
			glog.Error("Parse CfgBattleNpc.CharaID field error,value:", vCharaID)
			return false
		}

		/* parse Level field */
		vLevel, _ := parse.GetFieldByName(uint32(i), "level")
		var LevelRet bool
		data.Level, LevelRet = String2Int32(vLevel)
		if !LevelRet {
			glog.Error("Parse CfgBattleNpc.Level field error,value:", vLevel)
			return false
		}

		/* parse Star field */
		vStar, _ := parse.GetFieldByName(uint32(i), "star")
		var StarRet bool
		data.Star, StarRet = String2Int32(vStar)
		if !StarRet {
			glog.Error("Parse CfgBattleNpc.Star field error,value:", vStar)
			return false
		}

		/* parse Stage field */
		vStage, _ := parse.GetFieldByName(uint32(i), "stage")
		var StageRet bool
		data.Stage, StageRet = String2Int32(vStage)
		if !StageRet {
			glog.Error("Parse CfgBattleNpc.Stage field error,value:", vStage)
			return false
		}

		/* parse Equip1 field */
		vecEquip1, _ := parse.GetFieldByName(uint32(i), "equip1")
		if vecEquip1 != "" {
			arrayEquip1 := strings.Split(vecEquip1, ",")
			for j := 0; j < len(arrayEquip1); j++ {
				v, ret := String2Int32(arrayEquip1[j])
				if !ret {
					glog.Error("Parse CfgBattleNpc.Equip1 field error, value:", arrayEquip1[j])
					return false
				}
				data.Equip1 = append(data.Equip1, v)
			}
		}

		/* parse Equip2 field */
		vecEquip2, _ := parse.GetFieldByName(uint32(i), "equip2")
		if vecEquip2 != "" {
			arrayEquip2 := strings.Split(vecEquip2, ",")
			for j := 0; j < len(arrayEquip2); j++ {
				v, ret := String2Int32(arrayEquip2[j])
				if !ret {
					glog.Error("Parse CfgBattleNpc.Equip2 field error, value:", arrayEquip2[j])
					return false
				}
				data.Equip2 = append(data.Equip2, v)
			}
		}

		/* parse Equip3 field */
		vecEquip3, _ := parse.GetFieldByName(uint32(i), "equip3")
		if vecEquip3 != "" {
			arrayEquip3 := strings.Split(vecEquip3, ",")
			for j := 0; j < len(arrayEquip3); j++ {
				v, ret := String2Int32(arrayEquip3[j])
				if !ret {
					glog.Error("Parse CfgBattleNpc.Equip3 field error, value:", arrayEquip3[j])
					return false
				}
				data.Equip3 = append(data.Equip3, v)
			}
		}

		/* parse Equip4 field */
		vecEquip4, _ := parse.GetFieldByName(uint32(i), "equip4")
		if vecEquip4 != "" {
			arrayEquip4 := strings.Split(vecEquip4, ",")
			for j := 0; j < len(arrayEquip4); j++ {
				v, ret := String2Int32(arrayEquip4[j])
				if !ret {
					glog.Error("Parse CfgBattleNpc.Equip4 field error, value:", arrayEquip4[j])
					return false
				}
				data.Equip4 = append(data.Equip4, v)
			}
		}

		/* parse SkillLv field */
		vecSkillLv, _ := parse.GetFieldByName(uint32(i), "skillLv")
		arraySkillLv := strings.Split(vecSkillLv, ",")
		for j := 0; j < len(arraySkillLv); j++ {
			v := arraySkillLv[j]
			data.SkillLv = append(data.SkillLv, v)
		}

		/* parse WorldItem field */
		vWorldItem, _ := parse.GetFieldByName(uint32(i), "worldItem")
		var WorldItemRet bool
		data.WorldItem, WorldItemRet = String2Int32(vWorldItem)
		if !WorldItemRet {
			glog.Error("Parse CfgBattleNpc.WorldItem field error,value:", vWorldItem)
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

func (c *CfgBattleNpcConfig) Clear() {
}

func (c *CfgBattleNpcConfig) Find(id int32) (*CfgBattleNpc, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgBattleNpcConfig) GetAllData() map[int32]*CfgBattleNpc {
	return c.data
}

func (c *CfgBattleNpcConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.CharaID, ",", v.Level, ",", v.Star, ",", v.Stage, ",", v.Equip1, ",", v.Equip2, ",", v.Equip3, ",", v.Equip4, ",", v.SkillLv, ",", v.WorldItem)
	}
}
