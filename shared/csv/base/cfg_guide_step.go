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

type CfgGuideStep struct {
	Id       int32
	ModuleId int32
	StepId   int32
	Comment  string
}

type CfgGuideStepConfig struct {
	data map[int32]*CfgGuideStep
}

func NewCfgGuideStepConfig() *CfgGuideStepConfig {
	return &CfgGuideStepConfig{
		data: make(map[int32]*CfgGuideStep),
	}
}

func (c *CfgGuideStepConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgGuideStep)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgGuideStep.Id field error,value:", vId)
			return false
		}

		/* parse ModuleId field */
		vModuleId, _ := parse.GetFieldByName(uint32(i), "module_id")
		var ModuleIdRet bool
		data.ModuleId, ModuleIdRet = String2Int32(vModuleId)
		if !ModuleIdRet {
			glog.Error("Parse CfgGuideStep.ModuleId field error,value:", vModuleId)
			return false
		}

		/* parse StepId field */
		vStepId, _ := parse.GetFieldByName(uint32(i), "step_id")
		var StepIdRet bool
		data.StepId, StepIdRet = String2Int32(vStepId)
		if !StepIdRet {
			glog.Error("Parse CfgGuideStep.StepId field error,value:", vStepId)
			return false
		}

		/* parse Comment field */
		data.Comment, _ = parse.GetFieldByName(uint32(i), "comment")

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgGuideStepConfig) Clear() {
}

func (c *CfgGuideStepConfig) Find(id int32) (*CfgGuideStep, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgGuideStepConfig) GetAllData() map[int32]*CfgGuideStep {
	return c.data
}

func (c *CfgGuideStepConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.ModuleId, ",", v.StepId, ",", v.Comment)
	}
}
