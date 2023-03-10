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

type CfgGraveyardBuildStage struct {
	Id                   int32
	BuildId              int32
	Stage                int32
	MainTowerLevelLimit  int32
	OtherBuildingLvLimit []string
	StageUpConsume       []string
	StageUpTime          int32
	ProduceCountLimit    int32
	DropId               int32
	SpecialType          bool
	SpecialdropId        int32
	StageUpDropId        int32
}

type CfgGraveyardBuildStageConfig struct {
	data map[int32]*CfgGraveyardBuildStage
}

func NewCfgGraveyardBuildStageConfig() *CfgGraveyardBuildStageConfig {
	return &CfgGraveyardBuildStageConfig{
		data: make(map[int32]*CfgGraveyardBuildStage),
	}
}

func (c *CfgGraveyardBuildStageConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgGraveyardBuildStage)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgGraveyardBuildStage.Id field error,value:", vId)
			return false
		}

		/* parse BuildId field */
		vBuildId, _ := parse.GetFieldByName(uint32(i), "buildId")
		var BuildIdRet bool
		data.BuildId, BuildIdRet = String2Int32(vBuildId)
		if !BuildIdRet {
			glog.Error("Parse CfgGraveyardBuildStage.BuildId field error,value:", vBuildId)
			return false
		}

		/* parse Stage field */
		vStage, _ := parse.GetFieldByName(uint32(i), "stage")
		var StageRet bool
		data.Stage, StageRet = String2Int32(vStage)
		if !StageRet {
			glog.Error("Parse CfgGraveyardBuildStage.Stage field error,value:", vStage)
			return false
		}

		/* parse MainTowerLevelLimit field */
		vMainTowerLevelLimit, _ := parse.GetFieldByName(uint32(i), "mainTowerLevelLimit")
		var MainTowerLevelLimitRet bool
		data.MainTowerLevelLimit, MainTowerLevelLimitRet = String2Int32(vMainTowerLevelLimit)
		if !MainTowerLevelLimitRet {
			glog.Error("Parse CfgGraveyardBuildStage.MainTowerLevelLimit field error,value:", vMainTowerLevelLimit)
			return false
		}

		/* parse OtherBuildingLvLimit field */
		vecOtherBuildingLvLimit, _ := parse.GetFieldByName(uint32(i), "otherBuildingLvLimit")
		arrayOtherBuildingLvLimit := strings.Split(vecOtherBuildingLvLimit, ",")
		for j := 0; j < len(arrayOtherBuildingLvLimit); j++ {
			v := arrayOtherBuildingLvLimit[j]
			data.OtherBuildingLvLimit = append(data.OtherBuildingLvLimit, v)
		}

		/* parse StageUpConsume field */
		vecStageUpConsume, _ := parse.GetFieldByName(uint32(i), "stageUpConsume")
		arrayStageUpConsume := strings.Split(vecStageUpConsume, ",")
		for j := 0; j < len(arrayStageUpConsume); j++ {
			v := arrayStageUpConsume[j]
			data.StageUpConsume = append(data.StageUpConsume, v)
		}

		/* parse StageUpTime field */
		vStageUpTime, _ := parse.GetFieldByName(uint32(i), "stageUpTime")
		var StageUpTimeRet bool
		data.StageUpTime, StageUpTimeRet = String2Int32(vStageUpTime)
		if !StageUpTimeRet {
			glog.Error("Parse CfgGraveyardBuildStage.StageUpTime field error,value:", vStageUpTime)
			return false
		}

		/* parse ProduceCountLimit field */
		vProduceCountLimit, _ := parse.GetFieldByName(uint32(i), "produceCountLimit")
		var ProduceCountLimitRet bool
		data.ProduceCountLimit, ProduceCountLimitRet = String2Int32(vProduceCountLimit)
		if !ProduceCountLimitRet {
			glog.Error("Parse CfgGraveyardBuildStage.ProduceCountLimit field error,value:", vProduceCountLimit)
			return false
		}

		/* parse DropId field */
		vDropId, _ := parse.GetFieldByName(uint32(i), "dropId")
		var DropIdRet bool
		data.DropId, DropIdRet = String2Int32(vDropId)
		if !DropIdRet {
			glog.Error("Parse CfgGraveyardBuildStage.DropId field error,value:", vDropId)
			return false
		}

		/* parse SpecialType field */
		vSpecialType, _ := parse.GetFieldByName(uint32(i), "special_type")
		var SpecialTypeRet bool
		data.SpecialType, SpecialTypeRet = String2Bool(vSpecialType)
		if !SpecialTypeRet {
			glog.Error("Parse CfgGraveyardBuildStage.SpecialType field error,value:", vSpecialType)
		}

		/* parse SpecialdropId field */
		vSpecialdropId, _ := parse.GetFieldByName(uint32(i), "specialdropId")
		var SpecialdropIdRet bool
		data.SpecialdropId, SpecialdropIdRet = String2Int32(vSpecialdropId)
		if !SpecialdropIdRet {
			glog.Error("Parse CfgGraveyardBuildStage.SpecialdropId field error,value:", vSpecialdropId)
			return false
		}

		/* parse StageUpDropId field */
		vStageUpDropId, _ := parse.GetFieldByName(uint32(i), "stageUpDropId")
		var StageUpDropIdRet bool
		data.StageUpDropId, StageUpDropIdRet = String2Int32(vStageUpDropId)
		if !StageUpDropIdRet {
			glog.Error("Parse CfgGraveyardBuildStage.StageUpDropId field error,value:", vStageUpDropId)
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

func (c *CfgGraveyardBuildStageConfig) Clear() {
}

func (c *CfgGraveyardBuildStageConfig) Find(id int32) (*CfgGraveyardBuildStage, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgGraveyardBuildStageConfig) GetAllData() map[int32]*CfgGraveyardBuildStage {
	return c.data
}

func (c *CfgGraveyardBuildStageConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.BuildId, ",", v.Stage, ",", v.MainTowerLevelLimit, ",", v.OtherBuildingLvLimit, ",", v.StageUpConsume, ",", v.StageUpTime, ",", v.ProduceCountLimit, ",", v.DropId, ",", v.SpecialType, ",", v.SpecialdropId, ",", v.StageUpDropId)
	}
}
