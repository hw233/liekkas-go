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

type AppList struct {
	Id                 int32
	Channel            string
	AppVersion         string
	AppUrl             string
	CompatibleVersions []string
}

type AppListConfig struct {
	data map[int32]*AppList
}

func NewAppListConfig() *AppListConfig {
	return &AppListConfig{
		data: make(map[int32]*AppList),
	}
}

func (c *AppListConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(AppList)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse AppList.Id field error,value:", vId)
			return false
		}

		/* parse Channel field */
		data.Channel, _ = parse.GetFieldByName(uint32(i), "channel")

		/* parse AppVersion field */
		data.AppVersion, _ = parse.GetFieldByName(uint32(i), "appVersion")

		/* parse AppUrl field */
		data.AppUrl, _ = parse.GetFieldByName(uint32(i), "appUrl")

		/* parse CompatibleVersions field */
		vecCompatibleVersions, _ := parse.GetFieldByName(uint32(i), "compatibleVersions")
		arrayCompatibleVersions := strings.Split(vecCompatibleVersions, ",")
		for j := 0; j < len(arrayCompatibleVersions); j++ {
			v := arrayCompatibleVersions[j]
			data.CompatibleVersions = append(data.CompatibleVersions, v)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *AppListConfig) Clear() {
}

func (c *AppListConfig) Find(id int32) (*AppList, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *AppListConfig) GetAllData() map[int32]*AppList {
	return c.data
}

func (c *AppListConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Channel, ",", v.AppVersion, ",", v.AppUrl, ",", v.CompatibleVersions)
	}
}
