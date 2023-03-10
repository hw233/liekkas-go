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

type PatchList struct {
	Id              int32
	Channel         string
	AppVersion      string
	ResourceUrl     []string
	ResourceVersion string
}

type PatchListConfig struct {
	data map[int32]*PatchList
}

func NewPatchListConfig() *PatchListConfig {
	return &PatchListConfig{
		data: make(map[int32]*PatchList),
	}
}

func (c *PatchListConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(PatchList)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse PatchList.Id field error,value:", vId)
			return false
		}

		/* parse Channel field */
		data.Channel, _ = parse.GetFieldByName(uint32(i), "channel")

		/* parse AppVersion field */
		data.AppVersion, _ = parse.GetFieldByName(uint32(i), "appVersion")

		/* parse ResourceUrl field */
		vecResourceUrl, _ := parse.GetFieldByName(uint32(i), "resourceUrl")
		arrayResourceUrl := strings.Split(vecResourceUrl, ",")
		for j := 0; j < len(arrayResourceUrl); j++ {
			v := arrayResourceUrl[j]
			data.ResourceUrl = append(data.ResourceUrl, v)
		}

		/* parse ResourceVersion field */
		data.ResourceVersion, _ = parse.GetFieldByName(uint32(i), "resourceVersion")

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *PatchListConfig) Clear() {
}

func (c *PatchListConfig) Find(id int32) (*PatchList, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *PatchListConfig) GetAllData() map[int32]*PatchList {
	return c.data
}

func (c *PatchListConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Channel, ",", v.AppVersion, ",", v.ResourceUrl, ",", v.ResourceVersion)
	}
}
