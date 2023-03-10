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

type CfgCharacter struct {
	Id            int32
	CharaName     string
	CharaFullName string
	Visible       bool
	CharaShard    string
	Rarity        int32
	CharaGroup    int32
	Camp          int32
	Feature       []int32
	Career        int32
	RemoteAtk     bool
	Sex           int32
	SkinId        int32
	Type          int32
}

type CfgCharacterConfig struct {
	data map[int32]*CfgCharacter
}

func NewCfgCharacterConfig() *CfgCharacterConfig {
	return &CfgCharacterConfig{
		data: make(map[int32]*CfgCharacter),
	}
}

func (c *CfgCharacterConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgCharacter)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgCharacter.Id field error,value:", vId)
			return false
		}

		/* parse CharaName field */
		data.CharaName, _ = parse.GetFieldByName(uint32(i), "charaName")

		/* parse CharaFullName field */
		data.CharaFullName, _ = parse.GetFieldByName(uint32(i), "charaFullName")

		/* parse Visible field */
		vVisible, _ := parse.GetFieldByName(uint32(i), "visible")
		var VisibleRet bool
		data.Visible, VisibleRet = String2Bool(vVisible)
		if !VisibleRet {
			glog.Error("Parse CfgCharacter.Visible field error,value:", vVisible)
		}

		/* parse CharaShard field */
		data.CharaShard, _ = parse.GetFieldByName(uint32(i), "charaShard")

		/* parse Rarity field */
		vRarity, _ := parse.GetFieldByName(uint32(i), "rarity")
		var RarityRet bool
		data.Rarity, RarityRet = String2Int32(vRarity)
		if !RarityRet {
			glog.Error("Parse CfgCharacter.Rarity field error,value:", vRarity)
			return false
		}

		/* parse CharaGroup field */
		vCharaGroup, _ := parse.GetFieldByName(uint32(i), "charaGroup")
		var CharaGroupRet bool
		data.CharaGroup, CharaGroupRet = String2Int32(vCharaGroup)
		if !CharaGroupRet {
			glog.Error("Parse CfgCharacter.CharaGroup field error,value:", vCharaGroup)
			return false
		}

		/* parse Camp field */
		vCamp, _ := parse.GetFieldByName(uint32(i), "camp")
		var CampRet bool
		data.Camp, CampRet = String2Int32(vCamp)
		if !CampRet {
			glog.Error("Parse CfgCharacter.Camp field error,value:", vCamp)
			return false
		}

		/* parse Feature field */
		vecFeature, _ := parse.GetFieldByName(uint32(i), "feature")
		if vecFeature != "" {
			arrayFeature := strings.Split(vecFeature, ",")
			for j := 0; j < len(arrayFeature); j++ {
				v, ret := String2Int32(arrayFeature[j])
				if !ret {
					glog.Error("Parse CfgCharacter.Feature field error, value:", arrayFeature[j])
					return false
				}
				data.Feature = append(data.Feature, v)
			}
		}

		/* parse Career field */
		vCareer, _ := parse.GetFieldByName(uint32(i), "career")
		var CareerRet bool
		data.Career, CareerRet = String2Int32(vCareer)
		if !CareerRet {
			glog.Error("Parse CfgCharacter.Career field error,value:", vCareer)
			return false
		}

		/* parse RemoteAtk field */
		vRemoteAtk, _ := parse.GetFieldByName(uint32(i), "remoteAtk")
		var RemoteAtkRet bool
		data.RemoteAtk, RemoteAtkRet = String2Bool(vRemoteAtk)
		if !RemoteAtkRet {
			glog.Error("Parse CfgCharacter.RemoteAtk field error,value:", vRemoteAtk)
		}

		/* parse Sex field */
		vSex, _ := parse.GetFieldByName(uint32(i), "sex")
		var SexRet bool
		data.Sex, SexRet = String2Int32(vSex)
		if !SexRet {
			glog.Error("Parse CfgCharacter.Sex field error,value:", vSex)
			return false
		}

		/* parse SkinId field */
		vSkinId, _ := parse.GetFieldByName(uint32(i), "skinId")
		var SkinIdRet bool
		data.SkinId, SkinIdRet = String2Int32(vSkinId)
		if !SkinIdRet {
			glog.Error("Parse CfgCharacter.SkinId field error,value:", vSkinId)
			return false
		}

		/* parse Type field */
		vType, _ := parse.GetFieldByName(uint32(i), "type")
		var TypeRet bool
		data.Type, TypeRet = String2Int32(vType)
		if !TypeRet {
			glog.Error("Parse CfgCharacter.Type field error,value:", vType)
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

func (c *CfgCharacterConfig) Clear() {
}

func (c *CfgCharacterConfig) Find(id int32) (*CfgCharacter, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgCharacterConfig) GetAllData() map[int32]*CfgCharacter {
	return c.data
}

func (c *CfgCharacterConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.CharaName, ",", v.CharaFullName, ",", v.Visible, ",", v.CharaShard, ",", v.Rarity, ",", v.CharaGroup, ",", v.Camp, ",", v.Feature, ",", v.Career, ",", v.RemoteAtk, ",", v.Sex, ",", v.SkinId, ",", v.Type)
	}
}
