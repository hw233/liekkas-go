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

type CfgCharacterFeature struct {
	Id          int32
	FeatureName string
	FeatureIcon string
}

type CfgCharacterFeatureConfig struct {
	data map[int32]*CfgCharacterFeature
}

func NewCfgCharacterFeatureConfig() *CfgCharacterFeatureConfig {
	return &CfgCharacterFeatureConfig{
		data: make(map[int32]*CfgCharacterFeature),
	}
}

func (c *CfgCharacterFeatureConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgCharacterFeature)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgCharacterFeature.Id field error,value:", vId)
			return false
		}

		/* parse FeatureName field */
		data.FeatureName, _ = parse.GetFieldByName(uint32(i), "featureName")

		/* parse FeatureIcon field */
		data.FeatureIcon, _ = parse.GetFieldByName(uint32(i), "featureIcon")

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgCharacterFeatureConfig) Clear() {
}

func (c *CfgCharacterFeatureConfig) Find(id int32) (*CfgCharacterFeature, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgCharacterFeatureConfig) GetAllData() map[int32]*CfgCharacterFeature {
	return c.data
}

func (c *CfgCharacterFeatureConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.FeatureName, ",", v.FeatureIcon)
	}
}
