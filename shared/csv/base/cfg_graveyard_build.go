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

type CfgGraveyardBuild struct {
	Id                   int32
	Type                 int32
	MinProduceCount      int32
	CanReduceProductTime bool
	HorizontalCount      int32
	VerticalCount        int32
	UnlockCondition      []string
}

type CfgGraveyardBuildConfig struct {
	data map[int32]*CfgGraveyardBuild
}

func NewCfgGraveyardBuildConfig() *CfgGraveyardBuildConfig {
	return &CfgGraveyardBuildConfig{
		data: make(map[int32]*CfgGraveyardBuild),
	}
}

func (c *CfgGraveyardBuildConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgGraveyardBuild)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgGraveyardBuild.Id field error,value:", vId)
			return false
		}

		/* parse Type field */
		vType, _ := parse.GetFieldByName(uint32(i), "type")
		var TypeRet bool
		data.Type, TypeRet = String2Int32(vType)
		if !TypeRet {
			glog.Error("Parse CfgGraveyardBuild.Type field error,value:", vType)
			return false
		}

		/* parse MinProduceCount field */
		vMinProduceCount, _ := parse.GetFieldByName(uint32(i), "minProduceCount")
		var MinProduceCountRet bool
		data.MinProduceCount, MinProduceCountRet = String2Int32(vMinProduceCount)
		if !MinProduceCountRet {
			glog.Error("Parse CfgGraveyardBuild.MinProduceCount field error,value:", vMinProduceCount)
			return false
		}

		/* parse CanReduceProductTime field */
		vCanReduceProductTime, _ := parse.GetFieldByName(uint32(i), "canReduceProductTime")
		var CanReduceProductTimeRet bool
		data.CanReduceProductTime, CanReduceProductTimeRet = String2Bool(vCanReduceProductTime)
		if !CanReduceProductTimeRet {
			glog.Error("Parse CfgGraveyardBuild.CanReduceProductTime field error,value:", vCanReduceProductTime)
		}

		/* parse HorizontalCount field */
		vHorizontalCount, _ := parse.GetFieldByName(uint32(i), "horizontalCount")
		var HorizontalCountRet bool
		data.HorizontalCount, HorizontalCountRet = String2Int32(vHorizontalCount)
		if !HorizontalCountRet {
			glog.Error("Parse CfgGraveyardBuild.HorizontalCount field error,value:", vHorizontalCount)
			return false
		}

		/* parse VerticalCount field */
		vVerticalCount, _ := parse.GetFieldByName(uint32(i), "verticalCount")
		var VerticalCountRet bool
		data.VerticalCount, VerticalCountRet = String2Int32(vVerticalCount)
		if !VerticalCountRet {
			glog.Error("Parse CfgGraveyardBuild.VerticalCount field error,value:", vVerticalCount)
			return false
		}

		/* parse UnlockCondition field */
		vecUnlockCondition, _ := parse.GetFieldByName(uint32(i), "unlockCondition")
		arrayUnlockCondition := strings.Split(vecUnlockCondition, ",")
		for j := 0; j < len(arrayUnlockCondition); j++ {
			v := arrayUnlockCondition[j]
			data.UnlockCondition = append(data.UnlockCondition, v)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgGraveyardBuildConfig) Clear() {
}

func (c *CfgGraveyardBuildConfig) Find(id int32) (*CfgGraveyardBuild, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgGraveyardBuildConfig) GetAllData() map[int32]*CfgGraveyardBuild {
	return c.data
}

func (c *CfgGraveyardBuildConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Type, ",", v.MinProduceCount, ",", v.CanReduceProductTime, ",", v.HorizontalCount, ",", v.VerticalCount, ",", v.UnlockCondition)
	}
}
