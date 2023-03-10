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

type CfgGraveyardBuildLevel struct {
	Id                  int32
	BuildId             int32
	Level               int32
	MainTowerLevelLimit int32
	LevelUpConsume      []string
	LevelUpTime         int32
	LevelUpDropId       int32
	StoreLimit          int32
	ProduceLimit        int32
	ConsumeResource     []string
	ProduceTime         int32
	UnitCount           []int32
}

type CfgGraveyardBuildLevelConfig struct {
	data map[int32]*CfgGraveyardBuildLevel
}

func NewCfgGraveyardBuildLevelConfig() *CfgGraveyardBuildLevelConfig {
	return &CfgGraveyardBuildLevelConfig{
		data: make(map[int32]*CfgGraveyardBuildLevel),
	}
}

func (c *CfgGraveyardBuildLevelConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgGraveyardBuildLevel)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgGraveyardBuildLevel.Id field error,value:", vId)
			return false
		}

		/* parse BuildId field */
		vBuildId, _ := parse.GetFieldByName(uint32(i), "buildId")
		var BuildIdRet bool
		data.BuildId, BuildIdRet = String2Int32(vBuildId)
		if !BuildIdRet {
			glog.Error("Parse CfgGraveyardBuildLevel.BuildId field error,value:", vBuildId)
			return false
		}

		/* parse Level field */
		vLevel, _ := parse.GetFieldByName(uint32(i), "level")
		var LevelRet bool
		data.Level, LevelRet = String2Int32(vLevel)
		if !LevelRet {
			glog.Error("Parse CfgGraveyardBuildLevel.Level field error,value:", vLevel)
			return false
		}

		/* parse MainTowerLevelLimit field */
		vMainTowerLevelLimit, _ := parse.GetFieldByName(uint32(i), "mainTowerLevelLimit")
		var MainTowerLevelLimitRet bool
		data.MainTowerLevelLimit, MainTowerLevelLimitRet = String2Int32(vMainTowerLevelLimit)
		if !MainTowerLevelLimitRet {
			glog.Error("Parse CfgGraveyardBuildLevel.MainTowerLevelLimit field error,value:", vMainTowerLevelLimit)
			return false
		}

		/* parse LevelUpConsume field */
		vecLevelUpConsume, _ := parse.GetFieldByName(uint32(i), "levelUpConsume")
		arrayLevelUpConsume := strings.Split(vecLevelUpConsume, ",")
		for j := 0; j < len(arrayLevelUpConsume); j++ {
			v := arrayLevelUpConsume[j]
			data.LevelUpConsume = append(data.LevelUpConsume, v)
		}

		/* parse LevelUpTime field */
		vLevelUpTime, _ := parse.GetFieldByName(uint32(i), "levelUpTime")
		var LevelUpTimeRet bool
		data.LevelUpTime, LevelUpTimeRet = String2Int32(vLevelUpTime)
		if !LevelUpTimeRet {
			glog.Error("Parse CfgGraveyardBuildLevel.LevelUpTime field error,value:", vLevelUpTime)
			return false
		}

		/* parse LevelUpDropId field */
		vLevelUpDropId, _ := parse.GetFieldByName(uint32(i), "levelUpDropId")
		var LevelUpDropIdRet bool
		data.LevelUpDropId, LevelUpDropIdRet = String2Int32(vLevelUpDropId)
		if !LevelUpDropIdRet {
			glog.Error("Parse CfgGraveyardBuildLevel.LevelUpDropId field error,value:", vLevelUpDropId)
			return false
		}

		/* parse StoreLimit field */
		vStoreLimit, _ := parse.GetFieldByName(uint32(i), "storeLimit")
		var StoreLimitRet bool
		data.StoreLimit, StoreLimitRet = String2Int32(vStoreLimit)
		if !StoreLimitRet {
			glog.Error("Parse CfgGraveyardBuildLevel.StoreLimit field error,value:", vStoreLimit)
			return false
		}

		/* parse ProduceLimit field */
		vProduceLimit, _ := parse.GetFieldByName(uint32(i), "produceLimit")
		var ProduceLimitRet bool
		data.ProduceLimit, ProduceLimitRet = String2Int32(vProduceLimit)
		if !ProduceLimitRet {
			glog.Error("Parse CfgGraveyardBuildLevel.ProduceLimit field error,value:", vProduceLimit)
			return false
		}

		/* parse ConsumeResource field */
		vecConsumeResource, _ := parse.GetFieldByName(uint32(i), "consumeResource")
		arrayConsumeResource := strings.Split(vecConsumeResource, ",")
		for j := 0; j < len(arrayConsumeResource); j++ {
			v := arrayConsumeResource[j]
			data.ConsumeResource = append(data.ConsumeResource, v)
		}

		/* parse ProduceTime field */
		vProduceTime, _ := parse.GetFieldByName(uint32(i), "produceTime")
		var ProduceTimeRet bool
		data.ProduceTime, ProduceTimeRet = String2Int32(vProduceTime)
		if !ProduceTimeRet {
			glog.Error("Parse CfgGraveyardBuildLevel.ProduceTime field error,value:", vProduceTime)
			return false
		}

		/* parse UnitCount field */
		vecUnitCount, _ := parse.GetFieldByName(uint32(i), "unitCount")
		if vecUnitCount != "" {
			arrayUnitCount := strings.Split(vecUnitCount, ",")
			for j := 0; j < len(arrayUnitCount); j++ {
				v, ret := String2Int32(arrayUnitCount[j])
				if !ret {
					glog.Error("Parse CfgGraveyardBuildLevel.UnitCount field error, value:", arrayUnitCount[j])
					return false
				}
				data.UnitCount = append(data.UnitCount, v)
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

func (c *CfgGraveyardBuildLevelConfig) Clear() {
}

func (c *CfgGraveyardBuildLevelConfig) Find(id int32) (*CfgGraveyardBuildLevel, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgGraveyardBuildLevelConfig) GetAllData() map[int32]*CfgGraveyardBuildLevel {
	return c.data
}

func (c *CfgGraveyardBuildLevelConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.BuildId, ",", v.Level, ",", v.MainTowerLevelLimit, ",", v.LevelUpConsume, ",", v.LevelUpTime, ",", v.LevelUpDropId, ",", v.StoreLimit, ",", v.ProduceLimit, ",", v.ConsumeResource, ",", v.ProduceTime, ",", v.UnitCount)
	}
}
