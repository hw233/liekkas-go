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

type Gm struct {
	Id   int32
	Code string
}

type GmConfig struct {
	data map[int32]*Gm
}

func NewGmConfig() *GmConfig {
	return &GmConfig{
		data: make(map[int32]*Gm),
	}
}

func (c *GmConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(Gm)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse Gm.Id field error,value:", vId)
			return false
		}

		/* parse Code field */
		data.Code, _ = parse.GetFieldByName(uint32(i), "code")

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *GmConfig) Clear() {
}

func (c *GmConfig) Find(id int32) (*Gm, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *GmConfig) GetAllData() map[int32]*Gm {
	return c.data
}

func (c *GmConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Code)
	}
}
