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

type CfgGraveyardProduceBuff struct {
	Id          int32
	Type        int32
	TypeContent int32
	BuildId     int32
	Stage0drop  int32
	Stage1drop  int32
	Stage2drop  int32
	Stage3drop  int32
	Stage4drop  int32
	Stage5drop  int32
}

type CfgGraveyardProduceBuffConfig struct {
	data map[int32]*CfgGraveyardProduceBuff
}

func NewCfgGraveyardProduceBuffConfig() *CfgGraveyardProduceBuffConfig {
	return &CfgGraveyardProduceBuffConfig{
		data: make(map[int32]*CfgGraveyardProduceBuff),
	}
}

func (c *CfgGraveyardProduceBuffConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgGraveyardProduceBuff)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgGraveyardProduceBuff.Id field error,value:", vId)
			return false
		}

		/* parse Type field */
		vType, _ := parse.GetFieldByName(uint32(i), "type")
		var TypeRet bool
		data.Type, TypeRet = String2Int32(vType)
		if !TypeRet {
			glog.Error("Parse CfgGraveyardProduceBuff.Type field error,value:", vType)
			return false
		}

		/* parse TypeContent field */
		vTypeContent, _ := parse.GetFieldByName(uint32(i), "typeContent")
		var TypeContentRet bool
		data.TypeContent, TypeContentRet = String2Int32(vTypeContent)
		if !TypeContentRet {
			glog.Error("Parse CfgGraveyardProduceBuff.TypeContent field error,value:", vTypeContent)
			return false
		}

		/* parse BuildId field */
		vBuildId, _ := parse.GetFieldByName(uint32(i), "buildId")
		var BuildIdRet bool
		data.BuildId, BuildIdRet = String2Int32(vBuildId)
		if !BuildIdRet {
			glog.Error("Parse CfgGraveyardProduceBuff.BuildId field error,value:", vBuildId)
			return false
		}

		/* parse Stage0drop field */
		vStage0drop, _ := parse.GetFieldByName(uint32(i), "stage0drop")
		var Stage0dropRet bool
		data.Stage0drop, Stage0dropRet = String2Int32(vStage0drop)
		if !Stage0dropRet {
			glog.Error("Parse CfgGraveyardProduceBuff.Stage0drop field error,value:", vStage0drop)
			return false
		}

		/* parse Stage1drop field */
		vStage1drop, _ := parse.GetFieldByName(uint32(i), "stage1drop")
		var Stage1dropRet bool
		data.Stage1drop, Stage1dropRet = String2Int32(vStage1drop)
		if !Stage1dropRet {
			glog.Error("Parse CfgGraveyardProduceBuff.Stage1drop field error,value:", vStage1drop)
			return false
		}

		/* parse Stage2drop field */
		vStage2drop, _ := parse.GetFieldByName(uint32(i), "stage2drop")
		var Stage2dropRet bool
		data.Stage2drop, Stage2dropRet = String2Int32(vStage2drop)
		if !Stage2dropRet {
			glog.Error("Parse CfgGraveyardProduceBuff.Stage2drop field error,value:", vStage2drop)
			return false
		}

		/* parse Stage3drop field */
		vStage3drop, _ := parse.GetFieldByName(uint32(i), "stage3drop")
		var Stage3dropRet bool
		data.Stage3drop, Stage3dropRet = String2Int32(vStage3drop)
		if !Stage3dropRet {
			glog.Error("Parse CfgGraveyardProduceBuff.Stage3drop field error,value:", vStage3drop)
			return false
		}

		/* parse Stage4drop field */
		vStage4drop, _ := parse.GetFieldByName(uint32(i), "stage4drop")
		var Stage4dropRet bool
		data.Stage4drop, Stage4dropRet = String2Int32(vStage4drop)
		if !Stage4dropRet {
			glog.Error("Parse CfgGraveyardProduceBuff.Stage4drop field error,value:", vStage4drop)
			return false
		}

		/* parse Stage5drop field */
		vStage5drop, _ := parse.GetFieldByName(uint32(i), "stage5drop")
		var Stage5dropRet bool
		data.Stage5drop, Stage5dropRet = String2Int32(vStage5drop)
		if !Stage5dropRet {
			glog.Error("Parse CfgGraveyardProduceBuff.Stage5drop field error,value:", vStage5drop)
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

func (c *CfgGraveyardProduceBuffConfig) Clear() {
}

func (c *CfgGraveyardProduceBuffConfig) Find(id int32) (*CfgGraveyardProduceBuff, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgGraveyardProduceBuffConfig) GetAllData() map[int32]*CfgGraveyardProduceBuff {
	return c.data
}

func (c *CfgGraveyardProduceBuffConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Type, ",", v.TypeContent, ",", v.BuildId, ",", v.Stage0drop, ",", v.Stage1drop, ",", v.Stage2drop, ",", v.Stage3drop, ",", v.Stage4drop, ",", v.Stage5drop)
	}
}
