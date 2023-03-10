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

type CfgYggdrasilArea struct {
	Id                      int32
	WorldId                 int32
	AreaCityIds             []int32
	AreaPosId               []string
	UnlockCondition         string
	ItemMaxCount            int32
	MessageMaxMatchCount    int32
	MessageMaxCreateCount   int32
	MoveCost                int32
	AMaxMatchCount          int32
	AMaxBuildCount          int32
	BMaxMatchCount          int32
	BMaxBuildCount          int32
	CMaxMatchCount          int32
	CMaxBuildCount          int32
	DMaxMatchCount          int32
	DMaxBuildCount          int32
	EMaxMatchCount          int32
	EMaxBuildCount          int32
	FMaxMatchCount          int32
	FMaxBuildCount          int32
	ExploredProgressPercent []int32
	ExploredProgressDrop    []int32
	DispatchUnlock          int32
	StoreUnlock             int32
	GoStoreId               int32
	PrestigeItemID          int32
	MaxPrestige             int32
	DailyStar               []string
	GuildStar               []int32
	GuildNum                []int32
	GuildListLimit          int32
	SafePos                 []string
}

type CfgYggdrasilAreaConfig struct {
	data map[int32]*CfgYggdrasilArea
}

func NewCfgYggdrasilAreaConfig() *CfgYggdrasilAreaConfig {
	return &CfgYggdrasilAreaConfig{
		data: make(map[int32]*CfgYggdrasilArea),
	}
}

func (c *CfgYggdrasilAreaConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgYggdrasilArea)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgYggdrasilArea.Id field error,value:", vId)
			return false
		}

		/* parse WorldId field */
		vWorldId, _ := parse.GetFieldByName(uint32(i), "worldId")
		var WorldIdRet bool
		data.WorldId, WorldIdRet = String2Int32(vWorldId)
		if !WorldIdRet {
			glog.Error("Parse CfgYggdrasilArea.WorldId field error,value:", vWorldId)
			return false
		}

		/* parse AreaCityIds field */
		vecAreaCityIds, _ := parse.GetFieldByName(uint32(i), "areaCityIds")
		if vecAreaCityIds != "" {
			arrayAreaCityIds := strings.Split(vecAreaCityIds, ",")
			for j := 0; j < len(arrayAreaCityIds); j++ {
				v, ret := String2Int32(arrayAreaCityIds[j])
				if !ret {
					glog.Error("Parse CfgYggdrasilArea.AreaCityIds field error, value:", arrayAreaCityIds[j])
					return false
				}
				data.AreaCityIds = append(data.AreaCityIds, v)
			}
		}

		/* parse AreaPosId field */
		vecAreaPosId, _ := parse.GetFieldByName(uint32(i), "areaPosId")
		arrayAreaPosId := strings.Split(vecAreaPosId, ",")
		for j := 0; j < len(arrayAreaPosId); j++ {
			v := arrayAreaPosId[j]
			data.AreaPosId = append(data.AreaPosId, v)
		}

		/* parse UnlockCondition field */
		data.UnlockCondition, _ = parse.GetFieldByName(uint32(i), "unlockCondition")

		/* parse ItemMaxCount field */
		vItemMaxCount, _ := parse.GetFieldByName(uint32(i), "itemMaxCount")
		var ItemMaxCountRet bool
		data.ItemMaxCount, ItemMaxCountRet = String2Int32(vItemMaxCount)
		if !ItemMaxCountRet {
			glog.Error("Parse CfgYggdrasilArea.ItemMaxCount field error,value:", vItemMaxCount)
			return false
		}

		/* parse MessageMaxMatchCount field */
		vMessageMaxMatchCount, _ := parse.GetFieldByName(uint32(i), "messageMaxMatchCount")
		var MessageMaxMatchCountRet bool
		data.MessageMaxMatchCount, MessageMaxMatchCountRet = String2Int32(vMessageMaxMatchCount)
		if !MessageMaxMatchCountRet {
			glog.Error("Parse CfgYggdrasilArea.MessageMaxMatchCount field error,value:", vMessageMaxMatchCount)
			return false
		}

		/* parse MessageMaxCreateCount field */
		vMessageMaxCreateCount, _ := parse.GetFieldByName(uint32(i), "messageMaxCreateCount")
		var MessageMaxCreateCountRet bool
		data.MessageMaxCreateCount, MessageMaxCreateCountRet = String2Int32(vMessageMaxCreateCount)
		if !MessageMaxCreateCountRet {
			glog.Error("Parse CfgYggdrasilArea.MessageMaxCreateCount field error,value:", vMessageMaxCreateCount)
			return false
		}

		/* parse MoveCost field */
		vMoveCost, _ := parse.GetFieldByName(uint32(i), "moveCost")
		var MoveCostRet bool
		data.MoveCost, MoveCostRet = String2Int32(vMoveCost)
		if !MoveCostRet {
			glog.Error("Parse CfgYggdrasilArea.MoveCost field error,value:", vMoveCost)
			return false
		}

		/* parse AMaxMatchCount field */
		vAMaxMatchCount, _ := parse.GetFieldByName(uint32(i), "aMaxMatchCount")
		var AMaxMatchCountRet bool
		data.AMaxMatchCount, AMaxMatchCountRet = String2Int32(vAMaxMatchCount)
		if !AMaxMatchCountRet {
			glog.Error("Parse CfgYggdrasilArea.AMaxMatchCount field error,value:", vAMaxMatchCount)
			return false
		}

		/* parse AMaxBuildCount field */
		vAMaxBuildCount, _ := parse.GetFieldByName(uint32(i), "aMaxBuildCount")
		var AMaxBuildCountRet bool
		data.AMaxBuildCount, AMaxBuildCountRet = String2Int32(vAMaxBuildCount)
		if !AMaxBuildCountRet {
			glog.Error("Parse CfgYggdrasilArea.AMaxBuildCount field error,value:", vAMaxBuildCount)
			return false
		}

		/* parse BMaxMatchCount field */
		vBMaxMatchCount, _ := parse.GetFieldByName(uint32(i), "bMaxMatchCount")
		var BMaxMatchCountRet bool
		data.BMaxMatchCount, BMaxMatchCountRet = String2Int32(vBMaxMatchCount)
		if !BMaxMatchCountRet {
			glog.Error("Parse CfgYggdrasilArea.BMaxMatchCount field error,value:", vBMaxMatchCount)
			return false
		}

		/* parse BMaxBuildCount field */
		vBMaxBuildCount, _ := parse.GetFieldByName(uint32(i), "bMaxBuildCount")
		var BMaxBuildCountRet bool
		data.BMaxBuildCount, BMaxBuildCountRet = String2Int32(vBMaxBuildCount)
		if !BMaxBuildCountRet {
			glog.Error("Parse CfgYggdrasilArea.BMaxBuildCount field error,value:", vBMaxBuildCount)
			return false
		}

		/* parse CMaxMatchCount field */
		vCMaxMatchCount, _ := parse.GetFieldByName(uint32(i), "cMaxMatchCount")
		var CMaxMatchCountRet bool
		data.CMaxMatchCount, CMaxMatchCountRet = String2Int32(vCMaxMatchCount)
		if !CMaxMatchCountRet {
			glog.Error("Parse CfgYggdrasilArea.CMaxMatchCount field error,value:", vCMaxMatchCount)
			return false
		}

		/* parse CMaxBuildCount field */
		vCMaxBuildCount, _ := parse.GetFieldByName(uint32(i), "cMaxBuildCount")
		var CMaxBuildCountRet bool
		data.CMaxBuildCount, CMaxBuildCountRet = String2Int32(vCMaxBuildCount)
		if !CMaxBuildCountRet {
			glog.Error("Parse CfgYggdrasilArea.CMaxBuildCount field error,value:", vCMaxBuildCount)
			return false
		}

		/* parse DMaxMatchCount field */
		vDMaxMatchCount, _ := parse.GetFieldByName(uint32(i), "dMaxMatchCount")
		var DMaxMatchCountRet bool
		data.DMaxMatchCount, DMaxMatchCountRet = String2Int32(vDMaxMatchCount)
		if !DMaxMatchCountRet {
			glog.Error("Parse CfgYggdrasilArea.DMaxMatchCount field error,value:", vDMaxMatchCount)
			return false
		}

		/* parse DMaxBuildCount field */
		vDMaxBuildCount, _ := parse.GetFieldByName(uint32(i), "dMaxBuildCount")
		var DMaxBuildCountRet bool
		data.DMaxBuildCount, DMaxBuildCountRet = String2Int32(vDMaxBuildCount)
		if !DMaxBuildCountRet {
			glog.Error("Parse CfgYggdrasilArea.DMaxBuildCount field error,value:", vDMaxBuildCount)
			return false
		}

		/* parse EMaxMatchCount field */
		vEMaxMatchCount, _ := parse.GetFieldByName(uint32(i), "eMaxMatchCount")
		var EMaxMatchCountRet bool
		data.EMaxMatchCount, EMaxMatchCountRet = String2Int32(vEMaxMatchCount)
		if !EMaxMatchCountRet {
			glog.Error("Parse CfgYggdrasilArea.EMaxMatchCount field error,value:", vEMaxMatchCount)
			return false
		}

		/* parse EMaxBuildCount field */
		vEMaxBuildCount, _ := parse.GetFieldByName(uint32(i), "eMaxBuildCount")
		var EMaxBuildCountRet bool
		data.EMaxBuildCount, EMaxBuildCountRet = String2Int32(vEMaxBuildCount)
		if !EMaxBuildCountRet {
			glog.Error("Parse CfgYggdrasilArea.EMaxBuildCount field error,value:", vEMaxBuildCount)
			return false
		}

		/* parse FMaxMatchCount field */
		vFMaxMatchCount, _ := parse.GetFieldByName(uint32(i), "fMaxMatchCount")
		var FMaxMatchCountRet bool
		data.FMaxMatchCount, FMaxMatchCountRet = String2Int32(vFMaxMatchCount)
		if !FMaxMatchCountRet {
			glog.Error("Parse CfgYggdrasilArea.FMaxMatchCount field error,value:", vFMaxMatchCount)
			return false
		}

		/* parse FMaxBuildCount field */
		vFMaxBuildCount, _ := parse.GetFieldByName(uint32(i), "fMaxBuildCount")
		var FMaxBuildCountRet bool
		data.FMaxBuildCount, FMaxBuildCountRet = String2Int32(vFMaxBuildCount)
		if !FMaxBuildCountRet {
			glog.Error("Parse CfgYggdrasilArea.FMaxBuildCount field error,value:", vFMaxBuildCount)
			return false
		}

		/* parse ExploredProgressPercent field */
		vecExploredProgressPercent, _ := parse.GetFieldByName(uint32(i), "exploredProgressPercent")
		if vecExploredProgressPercent != "" {
			arrayExploredProgressPercent := strings.Split(vecExploredProgressPercent, ",")
			for j := 0; j < len(arrayExploredProgressPercent); j++ {
				v, ret := String2Int32(arrayExploredProgressPercent[j])
				if !ret {
					glog.Error("Parse CfgYggdrasilArea.ExploredProgressPercent field error, value:", arrayExploredProgressPercent[j])
					return false
				}
				data.ExploredProgressPercent = append(data.ExploredProgressPercent, v)
			}
		}

		/* parse ExploredProgressDrop field */
		vecExploredProgressDrop, _ := parse.GetFieldByName(uint32(i), "exploredProgressDrop")
		if vecExploredProgressDrop != "" {
			arrayExploredProgressDrop := strings.Split(vecExploredProgressDrop, ",")
			for j := 0; j < len(arrayExploredProgressDrop); j++ {
				v, ret := String2Int32(arrayExploredProgressDrop[j])
				if !ret {
					glog.Error("Parse CfgYggdrasilArea.ExploredProgressDrop field error, value:", arrayExploredProgressDrop[j])
					return false
				}
				data.ExploredProgressDrop = append(data.ExploredProgressDrop, v)
			}
		}

		/* parse DispatchUnlock field */
		vDispatchUnlock, _ := parse.GetFieldByName(uint32(i), "dispatchUnlock")
		var DispatchUnlockRet bool
		data.DispatchUnlock, DispatchUnlockRet = String2Int32(vDispatchUnlock)
		if !DispatchUnlockRet {
			glog.Error("Parse CfgYggdrasilArea.DispatchUnlock field error,value:", vDispatchUnlock)
			return false
		}

		/* parse StoreUnlock field */
		vStoreUnlock, _ := parse.GetFieldByName(uint32(i), "storeUnlock")
		var StoreUnlockRet bool
		data.StoreUnlock, StoreUnlockRet = String2Int32(vStoreUnlock)
		if !StoreUnlockRet {
			glog.Error("Parse CfgYggdrasilArea.StoreUnlock field error,value:", vStoreUnlock)
			return false
		}

		/* parse GoStoreId field */
		vGoStoreId, _ := parse.GetFieldByName(uint32(i), "goStoreId")
		var GoStoreIdRet bool
		data.GoStoreId, GoStoreIdRet = String2Int32(vGoStoreId)
		if !GoStoreIdRet {
			glog.Error("Parse CfgYggdrasilArea.GoStoreId field error,value:", vGoStoreId)
			return false
		}

		/* parse PrestigeItemID field */
		vPrestigeItemID, _ := parse.GetFieldByName(uint32(i), "prestigeItemID")
		var PrestigeItemIDRet bool
		data.PrestigeItemID, PrestigeItemIDRet = String2Int32(vPrestigeItemID)
		if !PrestigeItemIDRet {
			glog.Error("Parse CfgYggdrasilArea.PrestigeItemID field error,value:", vPrestigeItemID)
			return false
		}

		/* parse MaxPrestige field */
		vMaxPrestige, _ := parse.GetFieldByName(uint32(i), "maxPrestige")
		var MaxPrestigeRet bool
		data.MaxPrestige, MaxPrestigeRet = String2Int32(vMaxPrestige)
		if !MaxPrestigeRet {
			glog.Error("Parse CfgYggdrasilArea.MaxPrestige field error,value:", vMaxPrestige)
			return false
		}

		/* parse DailyStar field */
		vecDailyStar, _ := parse.GetFieldByName(uint32(i), "dailyStar")
		arrayDailyStar := strings.Split(vecDailyStar, ",")
		for j := 0; j < len(arrayDailyStar); j++ {
			v := arrayDailyStar[j]
			data.DailyStar = append(data.DailyStar, v)
		}

		/* parse GuildStar field */
		vecGuildStar, _ := parse.GetFieldByName(uint32(i), "guildStar")
		if vecGuildStar != "" {
			arrayGuildStar := strings.Split(vecGuildStar, ",")
			for j := 0; j < len(arrayGuildStar); j++ {
				v, ret := String2Int32(arrayGuildStar[j])
				if !ret {
					glog.Error("Parse CfgYggdrasilArea.GuildStar field error, value:", arrayGuildStar[j])
					return false
				}
				data.GuildStar = append(data.GuildStar, v)
			}
		}

		/* parse GuildNum field */
		vecGuildNum, _ := parse.GetFieldByName(uint32(i), "guildNum")
		if vecGuildNum != "" {
			arrayGuildNum := strings.Split(vecGuildNum, ",")
			for j := 0; j < len(arrayGuildNum); j++ {
				v, ret := String2Int32(arrayGuildNum[j])
				if !ret {
					glog.Error("Parse CfgYggdrasilArea.GuildNum field error, value:", arrayGuildNum[j])
					return false
				}
				data.GuildNum = append(data.GuildNum, v)
			}
		}

		/* parse GuildListLimit field */
		vGuildListLimit, _ := parse.GetFieldByName(uint32(i), "guildListLimit")
		var GuildListLimitRet bool
		data.GuildListLimit, GuildListLimitRet = String2Int32(vGuildListLimit)
		if !GuildListLimitRet {
			glog.Error("Parse CfgYggdrasilArea.GuildListLimit field error,value:", vGuildListLimit)
			return false
		}

		/* parse SafePos field */
		vecSafePos, _ := parse.GetFieldByName(uint32(i), "safePos")
		arraySafePos := strings.Split(vecSafePos, ",")
		for j := 0; j < len(arraySafePos); j++ {
			v := arraySafePos[j]
			data.SafePos = append(data.SafePos, v)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgYggdrasilAreaConfig) Clear() {
}

func (c *CfgYggdrasilAreaConfig) Find(id int32) (*CfgYggdrasilArea, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgYggdrasilAreaConfig) GetAllData() map[int32]*CfgYggdrasilArea {
	return c.data
}

func (c *CfgYggdrasilAreaConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.WorldId, ",", v.AreaCityIds, ",", v.AreaPosId, ",", v.UnlockCondition, ",", v.ItemMaxCount, ",", v.MessageMaxMatchCount, ",", v.MessageMaxCreateCount, ",", v.MoveCost, ",", v.AMaxMatchCount, ",", v.AMaxBuildCount, ",", v.BMaxMatchCount, ",", v.BMaxBuildCount, ",", v.CMaxMatchCount, ",", v.CMaxBuildCount, ",", v.DMaxMatchCount, ",", v.DMaxBuildCount, ",", v.EMaxMatchCount, ",", v.EMaxBuildCount, ",", v.FMaxMatchCount, ",", v.FMaxBuildCount, ",", v.ExploredProgressPercent, ",", v.ExploredProgressDrop, ",", v.DispatchUnlock, ",", v.StoreUnlock, ",", v.GoStoreId, ",", v.PrestigeItemID, ",", v.MaxPrestige, ",", v.DailyStar, ",", v.GuildStar, ",", v.GuildNum, ",", v.GuildListLimit, ",", v.SafePos)
	}
}
