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

type CfgItemData struct {
	Id          int32
	Name        string
	Cname       string
	ItemType    int32
	Rarity      int32
	UseType     int32
	UseParam    []int32
	ExpireType  int32
	ExpireTime  string
	LimitNumber int32
	SellPrice   int32
	ShowInBag   bool
}

type CfgItemDataConfig struct {
	data map[int32]*CfgItemData
}

func NewCfgItemDataConfig() *CfgItemDataConfig {
	return &CfgItemDataConfig{
		data: make(map[int32]*CfgItemData),
	}
}

func (c *CfgItemDataConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgItemData)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgItemData.Id field error,value:", vId)
			return false
		}

		/* parse Name field */
		data.Name, _ = parse.GetFieldByName(uint32(i), "name")

		/* parse Cname field */
		data.Cname, _ = parse.GetFieldByName(uint32(i), "cname")

		/* parse ItemType field */
		vItemType, _ := parse.GetFieldByName(uint32(i), "itemType")
		var ItemTypeRet bool
		data.ItemType, ItemTypeRet = String2Int32(vItemType)
		if !ItemTypeRet {
			glog.Error("Parse CfgItemData.ItemType field error,value:", vItemType)
			return false
		}

		/* parse Rarity field */
		vRarity, _ := parse.GetFieldByName(uint32(i), "rarity")
		var RarityRet bool
		data.Rarity, RarityRet = String2Int32(vRarity)
		if !RarityRet {
			glog.Error("Parse CfgItemData.Rarity field error,value:", vRarity)
			return false
		}

		/* parse UseType field */
		vUseType, _ := parse.GetFieldByName(uint32(i), "useType")
		var UseTypeRet bool
		data.UseType, UseTypeRet = String2Int32(vUseType)
		if !UseTypeRet {
			glog.Error("Parse CfgItemData.UseType field error,value:", vUseType)
			return false
		}

		/* parse UseParam field */
		vecUseParam, _ := parse.GetFieldByName(uint32(i), "useParam")
		if vecUseParam != "" {
			arrayUseParam := strings.Split(vecUseParam, ",")
			for j := 0; j < len(arrayUseParam); j++ {
				v, ret := String2Int32(arrayUseParam[j])
				if !ret {
					glog.Error("Parse CfgItemData.UseParam field error, value:", arrayUseParam[j])
					return false
				}
				data.UseParam = append(data.UseParam, v)
			}
		}

		/* parse ExpireType field */
		vExpireType, _ := parse.GetFieldByName(uint32(i), "expireType")
		var ExpireTypeRet bool
		data.ExpireType, ExpireTypeRet = String2Int32(vExpireType)
		if !ExpireTypeRet {
			glog.Error("Parse CfgItemData.ExpireType field error,value:", vExpireType)
			return false
		}

		/* parse ExpireTime field */
		data.ExpireTime, _ = parse.GetFieldByName(uint32(i), "expireTime")

		/* parse LimitNumber field */
		vLimitNumber, _ := parse.GetFieldByName(uint32(i), "limitNumber")
		var LimitNumberRet bool
		data.LimitNumber, LimitNumberRet = String2Int32(vLimitNumber)
		if !LimitNumberRet {
			glog.Error("Parse CfgItemData.LimitNumber field error,value:", vLimitNumber)
			return false
		}

		/* parse SellPrice field */
		vSellPrice, _ := parse.GetFieldByName(uint32(i), "sellPrice")
		var SellPriceRet bool
		data.SellPrice, SellPriceRet = String2Int32(vSellPrice)
		if !SellPriceRet {
			glog.Error("Parse CfgItemData.SellPrice field error,value:", vSellPrice)
			return false
		}

		/* parse ShowInBag field */
		vShowInBag, _ := parse.GetFieldByName(uint32(i), "showInBag")
		var ShowInBagRet bool
		data.ShowInBag, ShowInBagRet = String2Bool(vShowInBag)
		if !ShowInBagRet {
			glog.Error("Parse CfgItemData.ShowInBag field error,value:", vShowInBag)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgItemDataConfig) Clear() {
}

func (c *CfgItemDataConfig) Find(id int32) (*CfgItemData, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgItemDataConfig) GetAllData() map[int32]*CfgItemData {
	return c.data
}

func (c *CfgItemDataConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Name, ",", v.Cname, ",", v.ItemType, ",", v.Rarity, ",", v.UseType, ",", v.UseParam, ",", v.ExpireType, ",", v.ExpireTime, ",", v.LimitNumber, ",", v.SellPrice, ",", v.ShowInBag)
	}
}
